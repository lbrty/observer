package cmd

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"

	"github.com/lbrty/observer/internal/app"
	"github.com/lbrty/observer/internal/config"
	"github.com/lbrty/observer/internal/database"
	"github.com/lbrty/observer/internal/handler"
	"github.com/lbrty/observer/internal/logger"
	"github.com/lbrty/observer/internal/server"
	"github.com/lbrty/observer/migrations"
)

// ServeCmd starts the HTTP server.
var ServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the HTTP server",
	Long: `Start the Observer HTTP server.

Reads configuration from environment variables (DATABASE_DSN, REDIS_URL, etc.)
or a .env file. In production builds, embedded migrations are applied
automatically on startup. Graceful shutdown on SIGINT/SIGTERM with a
30-second timeout.`,
	Example: `  # Start with defaults (localhost:9000)
  observer serve

  # Custom host and port
  observer serve --host 0.0.0.0 --port 8080

  # With environment configuration
  DATABASE_DSN="postgres://..." REDIS_URL="redis://..." observer serve`,
	RunE: runServe,
}

func init() {
	ServeCmd.Flags().String("host", "", "Server host (overrides SERVER_HOST env)")
	ServeCmd.Flags().Int("port", 0, "Server port (overrides SERVER_PORT env)")
}

func runServe(cmd *cobra.Command, _ []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	if host, _ := cmd.Flags().GetString("host"); host != "" {
		cfg.Server.Host = host
	}
	if port, _ := cmd.Flags().GetInt("port"); port != 0 {
		cfg.Server.Port = port
	}

	log := logger.New(cfg.Log.Level)
	slog.SetDefault(log)

	for _, origin := range cfg.CORS.Origins {
		if strings.Contains(origin, "localhost") || strings.Contains(origin, "127.0.0.1") {
			slog.Warn("CORS_ORIGINS contains a localhost entry — ensure this is intentional in production",
				slog.String("origin", origin))
		}
	}

	if cfg.Sentry.Enabled() {
		err := sentry.Init(sentry.ClientOptions{
			Dsn:              cfg.Sentry.DSN,
			TracesSampleRate: cfg.Sentry.TracesSampleRate,
		})
		if err != nil {
			return fmt.Errorf("sentry init: %w", err)
		}
		defer sentry.Flush(2 * time.Second)
		log.Info("sentry enabled")
	}

	db, err := database.New(cfg.Database.DSN)
	if err != nil {
		return err
	}
	defer db.Close()

	redisOpts, err := redis.ParseURL(cfg.Redis.URL)
	if err != nil {
		return err
	}
	redisClient := redis.NewClient(redisOpts)
	defer redisClient.Close()

	// Auto-run migrations when embedded in production build.
	if migrations.Embedded() {
		if err := autoMigrate(cfg.Database.DSN, log); err != nil {
			return err
		}
	}

	schemaStatus := checkMigrationDrift(cfg.Database.DSN, log)

	container, err := app.NewContainer(cfg, db, redisClient)
	if err != nil {
		return err
	}

	srv := server.New(cfg, db, log, container, schemaStatus)

	vacuumCtx, vacuumCancel := context.WithCancel(context.Background())
	defer vacuumCancel()
	app.StartSessionVacuum(vacuumCtx, container.SessionRepo)

	go func() {
		log.Info("server starting", slog.String("addr", cfg.Server.Host), slog.Int("port", cfg.Server.Port))
		if err := srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("server error", slog.Any("err", err))
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return srv.Shutdown(ctx)
}

// checkMigrationDrift computes how many migrations are pending without applying them.
// On any error it logs a warning and returns a zero-valued status so the server still starts.
func checkMigrationDrift(dsn string, log *slog.Logger) handler.SchemaStatus {
	var s handler.SchemaStatus

	var m *migrate.Migrate
	var err error

	if migrations.Embedded() {
		fsys, ferr := migrations.FS()
		if ferr != nil {
			log.Warn("schema drift check: embedded fs unavailable", slog.Any("err", ferr))
			return s
		}
		d, derr := iofs.New(fsys, ".")
		if derr != nil {
			log.Warn("schema drift check: iofs source", slog.Any("err", derr))
			return s
		}
		m, err = migrate.NewWithSourceInstance("iofs", d, dsn)
		if err != nil {
			log.Warn("schema drift check: migrate init", slog.Any("err", err))
			return s
		}
		defer m.Close()
		s.LatestVersion = scanMaxVersionFS(fsys)
	} else {
		m, err = migrate.New("file://migrations", dsn)
		if err != nil {
			log.Warn("schema drift check: migrate init", slog.Any("err", err))
			return s
		}
		defer m.Close()
		s.LatestVersion = scanMaxVersionDir("migrations")
	}

	v, dirty, verr := m.Version()
	if verr != nil && verr != migrate.ErrNilVersion {
		log.Warn("schema drift check: version query", slog.Any("err", verr))
		return s
	}
	s.CurrentVersion = v
	s.Dirty = dirty

	if s.LatestVersion > s.CurrentVersion {
		s.Pending = int(s.LatestVersion - s.CurrentVersion)
		log.Warn("database schema is behind",
			slog.Uint64("current", uint64(s.CurrentVersion)),
			slog.Uint64("latest", uint64(s.LatestVersion)),
			slog.Int("pending", s.Pending),
		)
	}

	return s
}

func scanMaxVersionFS(fsys fs.FS) uint {
	entries, _ := fs.ReadDir(fsys, ".")
	return maxMigrationSeq(entries)
}

func scanMaxVersionDir(dir string) uint {
	entries, _ := os.ReadDir(dir)
	result := make([]fs.DirEntry, len(entries))
	for i, e := range entries {
		result[i] = e
	}
	return maxMigrationSeq(result)
}

func maxMigrationSeq(entries []fs.DirEntry) uint {
	var max uint
	for _, e := range entries {
		var n uint
		if _, err := fmt.Sscanf(e.Name(), "%06d", &n); err == nil && n > max {
			max = n
		}
	}
	return max
}

func autoMigrate(dsn string, log *slog.Logger) error {
	fsys, err := migrations.FS()
	if err != nil {
		return fmt.Errorf("embedded migrations: %w", err)
	}
	d, err := iofs.New(fsys, ".")
	if err != nil {
		return fmt.Errorf("iofs source: %w", err)
	}
	m, err := migrate.NewWithSourceInstance("iofs", d, dsn)
	if err != nil {
		return fmt.Errorf("migrate init: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migrate up: %w", err)
	}

	v, dirty, _ := m.Version()
	log.Info("migrations applied", slog.Uint64("version", uint64(v)), slog.Bool("dirty", dirty))
	return nil
}

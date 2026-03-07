package cmd

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
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
	"github.com/lbrty/observer/internal/logger"
	"github.com/lbrty/observer/internal/server"
	"github.com/lbrty/observer/migrations"
)

// ServeCmd starts the HTTP server.
var ServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the HTTP server",
	RunE:  runServe,
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

	container, err := app.NewContainer(cfg, db, redisClient)
	if err != nil {
		return err
	}

	srv := server.New(cfg, db, log, container)

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

package cmd

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/lbrty/observer/internal/app"
	"github.com/lbrty/observer/internal/config"
	"github.com/lbrty/observer/internal/database"
	"github.com/lbrty/observer/internal/logger"
	"github.com/lbrty/observer/internal/server"
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

	db, err := database.New(cfg.Database.DSN)
	if err != nil {
		return err
	}
	defer db.Close()

	container, err := app.NewContainer(cfg, db)
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

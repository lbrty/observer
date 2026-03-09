package app

import (
	"context"
	"log/slog"
	"time"

	"github.com/lbrty/observer/internal/repository"
)

// StartSessionVacuum deletes expired sessions every 24 hours.
// Call with a cancellable context; cancel it to stop.
func StartSessionVacuum(ctx context.Context, repo repository.SessionRepository) {
	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				n, err := repo.DeleteExpired(ctx)
				if err != nil {
					slog.Error("session vacuum", slog.Any("err", err))
				} else {
					slog.Info("session vacuum", slog.Int64("deleted", n))
				}
			}
		}
	}()
}

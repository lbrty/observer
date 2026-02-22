package project

import (
	"context"

	"github.com/oklog/ulid/v2"
)

//go:generate mockgen -destination=mock/repository.go -package=mock github.com/lbrty/observer/internal/domain/project PermissionLoader

// PermissionLoader loads project-level permissions for authorization.
type PermissionLoader interface {
	GetPermission(ctx context.Context, userID ulid.ULID, projectID string) (*Permission, error)
	IsProjectOwner(ctx context.Context, userID ulid.ULID, projectID string) (bool, error)
}

package user

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestUserCanLoginRejectsPermanentLock(t *testing.T) {
	now := time.Now()
	u := &User{IsActive: true, LockedPermanentlyAt: &now}

	require.ErrorIs(t, u.CanLogin(), ErrUserNotActive)
}

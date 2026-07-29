package auditctx

import "context"

type key string

const (
	clientIPKey  key = "client_ip"
	userAgentKey key = "user_agent"
	userIDKey    key = "user_id"
)

// WithMetadata enriches a context with audit-relevant request metadata.
func WithMetadata(ctx context.Context, userID, ip, userAgent string) context.Context {
	ctx = context.WithValue(ctx, userIDKey, userID)
	ctx = context.WithValue(ctx, clientIPKey, ip)
	ctx = context.WithValue(ctx, userAgentKey, userAgent)
	return ctx
}

// UserID extracts the user ID from context.
func UserID(ctx context.Context) string {
	value, _ := ctx.Value(userIDKey).(string)
	return value
}

// IP extracts the client IP from context.
func IP(ctx context.Context) string {
	value, _ := ctx.Value(clientIPKey).(string)
	return value
}

// UserAgent extracts the user agent from context.
func UserAgent(ctx context.Context) string {
	value, _ := ctx.Value(userAgentKey).(string)
	return value
}

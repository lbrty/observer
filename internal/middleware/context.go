package middleware

import (
	"context"
)

const (
	ctxClientIP    ctxKey = "client_ip"
	ctxUserAgentKy ctxKey = "user_agent"
	ctxUserIDStr   ctxKey = "user_id_str"
)

// WithAuditContext enriches the context with audit-relevant HTTP metadata.
func WithAuditContext(ctx context.Context, userID, ip, userAgent string) context.Context {
	ctx = context.WithValue(ctx, ctxUserIDStr, userID)
	ctx = context.WithValue(ctx, ctxClientIP, ip)
	ctx = context.WithValue(ctx, ctxUserAgentKy, userAgent)
	return ctx
}

// AuditUserID extracts the user ID string from context.
func AuditUserID(ctx context.Context) string {
	if v, ok := ctx.Value(ctxUserIDStr).(string); ok {
		return v
	}
	return ""
}

// AuditIP extracts the client IP from context.
func AuditIP(ctx context.Context) string {
	if v, ok := ctx.Value(ctxClientIP).(string); ok {
		return v
	}
	return ""
}

// AuditUserAgent extracts the user agent from context.
func AuditUserAgent(ctx context.Context) string {
	if v, ok := ctx.Value(ctxUserAgentKy).(string); ok {
		return v
	}
	return ""
}

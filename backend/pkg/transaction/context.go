package transaction

import (
	"context"
)

// CONTEXT NOT USED YET

type contextKey string

const (
	sessionKey        contextKey = "session"
	sessionManagerKey contextKey = "session-manager"
)

// SetSession установка сессии в контекст
func SetSession(ctx context.Context, s Session) context.Context {
	return context.WithValue(ctx, sessionKey, s)
}

// MustGetSession получение сессии из контекста
func MustGetSession(ctx context.Context) Session {
	return ctx.Value(sessionKey).(Session)
}

// SetSessionManager установка менеджера сессий в контекст
func SetSessionManager(ctx context.Context, sm SessionManager) context.Context {
	return context.WithValue(ctx, sessionManagerKey, sm)
}

// GetSessionManager получение менеджера сессий из контекста
func MustGetSessionManager(ctx context.Context) SessionManager {
	return ctx.Value(sessionManagerKey).(SessionManager)
}

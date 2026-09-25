package audit

import (
	"context"

	"github.com/google/uuid"
)

type contextKey string

const actorContextKey contextKey = "erp_actor_ctx"

// Actor represents the user context executing an action
type Actor struct {
	UserID    uuid.UUID
	CompanyID uuid.UUID
	BranchID  uuid.UUID
	Role      string
	IP        string
	UserAgent string
}

// WithActor injects the Actor into context
func WithActor(ctx context.Context, actor Actor) context.Context {
	return context.WithValue(ctx, actorContextKey, actor)
}

// GetActor extracts the Actor from context
func GetActor(ctx context.Context) (Actor, bool) {
	actor, ok := ctx.Value(actorContextKey).(Actor)
	return actor, ok
}

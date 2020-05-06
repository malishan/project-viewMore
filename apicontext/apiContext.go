package apicontext

import (
	"context"
)

type ctxType string

const (
	// APICtx - defining a separate type to avoid colliding with basic type
	APICtx ctxType = "apiCtx"
)

// APIContext contains context of client
type APIContext struct {
	Token     string // Token is the api token
	RequestID string // RequestID used to track logs across a request-response cycle
	UserID    string // User Id
	UserName  string // Username against the token
	Email     string //Email against the token
	RoleID    string //Role of the user
}

// CustomContext is the combination of native context and APIContext
type CustomContext struct {
	context.Context
	APIContext
}

// GetAPIContext returns the APIContext from the native context provided
func GetAPIContext(ctx context.Context) (APIContext, bool) {
	if ctx == nil {
		return APIContext{}, false
	}
	apiCtx, exists := ctx.Value(APICtx).(APIContext)
	return apiCtx, exists
}

// WithAPIContext returns a new context with the APIContext binded to the native context
func WithAPIContext(ctx context.Context, apictx APIContext) context.Context {
	return context.WithValue(ctx, APICtx, apictx)
}

// UpgradeContext embeds native context and APIContext to form the CustomContext
func UpgradeContext(ctx context.Context) CustomContext {
	var cContext CustomContext
	apiCtx, _ := GetAPIContext(ctx)

	cContext.Context = ctx
	cContext.APIContext = apiCtx
	return cContext
}

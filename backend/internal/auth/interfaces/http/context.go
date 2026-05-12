package http

import (
	"context"

	authdomain "github.com/terraroute/terra-route/backend/internal/auth/domain"
)

type authClaimsContextKey struct{}

func contextWithAuthClaims(ctx context.Context, claims *authdomain.TokenClaims) context.Context {
	return context.WithValue(ctx, authClaimsContextKey{}, claims)
}

func AuthClaimsFromContext(ctx context.Context) (*authdomain.TokenClaims, bool) {
	claims, ok := ctx.Value(authClaimsContextKey{}).(*authdomain.TokenClaims)
	if !ok || claims == nil {
		return nil, false
	}
	return claims, true
}

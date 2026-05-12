package http

import (
	"net/http"

	userdomain "github.com/terraroute/terra-route/backend/internal/users/domain"
)

func RequireRoles(allowedRoles ...string) func(http.Handler) http.Handler {
	allowed := make(map[string]struct{}, len(allowedRoles))
	for _, role := range allowedRoles {
		allowed[role] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := AuthClaimsFromContext(r.Context())
			if !ok {
				writeError(w, http.StatusUnauthorized, "not_authenticated")
				return
			}

			if _, ok := allowed[claims.Role]; !ok {
				writeError(w, http.StatusForbidden, "forbidden")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func RequireCompanyAdminOrOperator(next http.Handler) http.Handler {
	return RequireRoles(userdomain.UserRoleCompanyAdmin, userdomain.UserRoleOperator)(next)
}

func RequireSuperAdmin(next http.Handler) http.Handler {
	return RequireRoles(userdomain.UserRoleSuperAdmin)(next)
}

func RequireDriver(next http.Handler) http.Handler {
	return RequireRoles(userdomain.UserRoleDriver)(next)
}

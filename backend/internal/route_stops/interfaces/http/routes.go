package http

import (
	"net/http"

	authdomain "github.com/terraroute/terra-route/backend/internal/auth/domain"
	authhttp "github.com/terraroute/terra-route/backend/internal/auth/interfaces/http"
)

func RegisterRoutes(mux *http.ServeMux, handler *Handler, tokens authdomain.TokenService) {
	protected := authhttp.AuthMiddleware(tokens)(authhttp.RequireCompanyAdminOrOperator(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.PathValue("stopId") == "":
			handler.List(w, r)
		case r.Method == http.MethodPost && r.PathValue("stopId") == "":
			handler.Add(w, r)
		case r.Method == http.MethodPatch:
			handler.Update(w, r)
		case r.Method == http.MethodDelete:
			handler.Delete(w, r)
		default:
			http.NotFound(w, r)
		}
	})))

	mux.Handle("GET /api/v1/routes/{id}/stops", protected)
	mux.Handle("POST /api/v1/routes/{id}/stops", protected)
	mux.Handle("PATCH /api/v1/routes/{id}/stops/{stopId}", protected)
	mux.Handle("DELETE /api/v1/routes/{id}/stops/{stopId}", protected)
}

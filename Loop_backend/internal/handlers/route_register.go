package handlers

import (
	"Loop_backend/internal/middleware"
	"Loop_backend/internal/services/auth"
	"github.com/gorilla/mux"
	"net/http"
)

// RouteRegister defines the interface for registering routes
type RouteRegister interface {
	RegisterProtectedRoute(path string, handler http.HandlerFunc, dto interface{})
	RegisterPublicRoute(path string, handler http.HandlerFunc, dto interface{})
}

// DefaultRouteRegister implements RouteRegister using mux.Router
type DefaultRouteRegister struct {
	router      *mux.Router
	authService auth.AuthService
}

// NewRouteRegister creates a new route register instance
func NewRouteRegister(router *mux.Router, authService auth.AuthService) RouteRegister {
	return &DefaultRouteRegister{
		router:      router,
		authService: authService,
	}

}

// RegisterProtectedRoute registers a new protected route with auth middleware
func (r *DefaultRouteRegister) RegisterProtectedRoute(path string, handler http.HandlerFunc, dto interface{}) {
	validatedHandler := middleware.ValidateRequest(handler, dto)
	authedHandler := middleware.WithAuth(validatedHandler, r.authService)
	r.router.HandleFunc(path, authedHandler).Methods("GET", "POST", "PUT", "DELETE")
}

// RegisterPublicRoute registers a new public route
func (r *DefaultRouteRegister) RegisterPublicRoute(path string, handler http.HandlerFunc, dto interface{}) {
	validatedHandler := middleware.ValidateRequest(handler, dto)
	r.router.HandleFunc(path, validatedHandler).Methods("GET", "POST", "PUT", "DELETE")
}

package handlers

import (
    "net/http"
    "Loop_backend/internal/middleware"
    "Loop_backend/internal/services"
)

type RouteRegister struct {
    mux         *http.ServeMux
    authService services.AuthService
}

func NewRouteRegister(mux *http.ServeMux, authService services.AuthService) *RouteRegister {
    return &RouteRegister{
        mux:         mux,
        authService: authService,
    }
}

// Protected registers a route that requires authentication
func (r *RouteRegister) RegisterProtectedRoute(path string, handler http.HandlerFunc) {
    r.mux.HandleFunc(path, middleware.WithAuth(handler, r.authService))
}

// Public registers a route that doesn't require authentication
func (r *RouteRegister) RegisterPublicRoute(path string, handler http.HandlerFunc) {
    r.mux.HandleFunc(path, handler)
}


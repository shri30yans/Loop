package handlers

import (
    "net/http"
    "github.com/gorilla/mux"
)

// RouteRegister is an interface for registering routes
type RouteRegister interface {
    HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
    RegisterPublicRouteWithMethod(pattern string, method string, handler http.HandlerFunc)
    RegisterProtectedRouteWithMethod(pattern string, method string, handler http.HandlerFunc)
    // Keeping old methods for backward compatibility
    RegisterPublicRoute(pattern string, handler http.HandlerFunc)
    RegisterProtectedRoute(pattern string, handler http.HandlerFunc)
}

// Handler is an interface for all handlers that can register routes
type Handler interface {
    RegisterRoutes(router RouteRegister)
}


// ExtendedServeMux wraps mux.Router to implement RouteRegister
type ExtendedServeMux struct {
    router        *mux.Router
    authMiddleware func(http.HandlerFunc) http.HandlerFunc
}

// Router returns the underlying mux.Router
func (r *ExtendedServeMux) Router() *mux.Router {
    return r.router
}

// NewRouteRegister creates a new ExtendedServeMux with auth middleware
func NewRouteRegister(authMiddleware func(http.HandlerFunc) http.HandlerFunc) *ExtendedServeMux {
    return &ExtendedServeMux{
        router:        mux.NewRouter(),
        authMiddleware: authMiddleware,
    }
}

// HandleFunc implements the RouteRegister interface
func (r *ExtendedServeMux) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
    r.router.HandleFunc(pattern, handler)
}

// RegisterPublicRouteWithMethod registers a public route with a specific HTTP method
func (r *ExtendedServeMux) RegisterPublicRouteWithMethod(pattern string, method string, handler http.HandlerFunc) {
    r.router.HandleFunc(pattern, handler).Methods(method)
}

// RegisterProtectedRouteWithMethod registers a protected route with a specific HTTP method
func (r *ExtendedServeMux) RegisterProtectedRouteWithMethod(pattern string, method string, handler http.HandlerFunc) {
    r.router.HandleFunc(pattern, r.authMiddleware(handler)).Methods(method)
}

// RegisterPublicRoute registers a public route without auth middleware (defaults to GET)
func (r *ExtendedServeMux) RegisterPublicRoute(pattern string, handler http.HandlerFunc) {
    r.RegisterPublicRouteWithMethod(pattern, http.MethodGet, handler)
}

// RegisterProtectedRoute registers a protected route with auth middleware (defaults to GET)
func (r *ExtendedServeMux) RegisterProtectedRoute(pattern string, handler http.HandlerFunc) {
    r.RegisterProtectedRouteWithMethod(pattern, http.MethodGet, handler)
}

// Ensure all handlers implement the Handler interface
var (
	_ Handler = (*SearchHandler)(nil)
	_ Handler = (*ToolsHandler)(nil)
)

package handlers

import (
	"net/http"

	"Loop_backend/internal/middleware"
	"Loop_backend/internal/response"
	"Loop_backend/internal/services/user"
	"github.com/gorilla/mux"
)

type UserHandler struct {
    userService user.UserService
}

func NewUserHandler(userService user.UserService, authMiddleware func(http.HandlerFunc) http.HandlerFunc) *UserHandler {
    return &UserHandler{
        userService: userService,
    }
}

func (h *UserHandler) RegisterRoutes(r RouteRegister) {
    r.RegisterProtectedRouteWithMethod("/api/user/{user_id:[a-fA-F0-9-]+}", http.MethodGet, h.GetUserInfo)
    r.RegisterProtectedRouteWithMethod("/api/user/delete", http.MethodDelete, h.DeleteAccount)
}

func (h *UserHandler) GetUserInfo(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userID, ok := vars["user_id"]
    if !ok || userID == "" {
        response.RespondWithError(w, http.StatusBadRequest, "uuid parameter is required")
        return
    }

    user, err := h.userService.GetUser(userID)
    if err != nil {
        response.RespondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    response.RespondWithJSON(w, http.StatusOK, user)
}

func (h *UserHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value(middleware.UserIDKey).(string)

    if err := h.userService.DeleteUser(userID); err != nil {
        response.RespondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    response.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Account deleted successfully"})
}

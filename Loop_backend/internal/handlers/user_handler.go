package handlers

import (
	"net/http"

	"Loop_backend/internal/middleware"
	"Loop_backend/internal/response"
	"Loop_backend/internal/services"
	"github.com/gorilla/mux"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) RegisterRoutes(r RouteRegister) {
	r.RegisterProtectedRoute("/api/user/{user_id:[a-fA-F0-9-]+}", h.GetUserInfo, nil)
	r.RegisterProtectedRoute("/api/user/delete", h.DeleteAccount, nil)
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

package handlers

import (
	"net/http"

<<<<<<< HEAD
	"github.com/gorilla/mux"
=======
>>>>>>> 4a2f436bed91636c5c2e3782993f5ab211ecfca7
	"Loop_backend/internal/middleware"
	"Loop_backend/internal/response"
	"Loop_backend/internal/services"
)

type UserHandler struct {
<<<<<<< HEAD
	userService services.UserService
}

func (h *UserHandler) RegisterRoutes(r *RouteRegister) {
	r.RegisterProtectedRoute("/api/user/{user_id:[a-fA-F0-9-]+}", h.GetUserInfo)
	r.RegisterProtectedRoute("/api/user/delete", h.DeleteAccount)
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
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
=======
    userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
    return &UserHandler{
        userService: userService,
    }
}

func (h *UserHandler) GetUserInfo(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value(middleware.UserIDKey).(string)
    
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

func (h *UserHandler) RegisterRoutes(r *RouteRegister) {
    r.RegisterProtectedRoute("/api/users/info", h.GetUserInfo)
    r.RegisterProtectedRoute("/api/users/delete", h.DeleteAccount)

}
>>>>>>> 4a2f436bed91636c5c2e3782993f5ab211ecfca7

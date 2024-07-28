package user

import (
	"fmt"
	"net/http"

	"github.com/ak-ko/ghop.git/service/auth"
	"github.com/ak-ko/ghop.git/types"
	"github.com/ak-ko/ghop.git/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) AuthenticateRoutes(router *mux.Router) {
	router.HandleFunc("/register", h.handleRegister).Methods("POST")
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
}



func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var payload types.RegisterUserPayload

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}

	_,err := h.store.GetUserByEmail(payload.Email); if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user already exists"))
		return
	}

	if err := utils.Validator.Struct(payload); err != nil {
		vError := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, vError)
		return
	}

	hashedPassword, err := auth.HashPassword(payload.Password); if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}

	if err = h.store.CreateUser(types.User{
		Username: payload.Username,
		Email: payload.Email,
		Password: hashedPassword,
	}); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusCreated, nil)
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {

}
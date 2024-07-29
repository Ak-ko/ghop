package user

import (
	"fmt"
	"net/http"

	"github.com/ak-ko/ghop.git/config"
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

	// parsing
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}

	// validating
	if err := utils.Validator.Struct(payload); err != nil {
		vError := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", vError))
		return
	}

	// check user exists
	_,err := h.store.GetUserByEmail(payload.Email)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with email %s already exists", payload.Email))
		return
	}

	// hashing
	hashedPassword, err := auth.HashPassword(payload.Password)
	 if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}

	// creating a user
	if err = h.store.CreateUser(types.User{
		Username: payload.Username,
		Email: payload.Email,
		Password: hashedPassword,
	}); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusCreated, nil) // token
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var payload types.LoginUserPayload

	// parsing
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}

	// validate
	if err := utils.Validator.Struct(payload); err != nil {
		vError := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, vError)
		return
	}

	// check user exists
	user, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid credentials"))
	}

	// password check
	pwMatch := auth.ComparePasswords(user.Password, []byte(payload.Password))
	if !pwMatch {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("incorrect password"))
		return
	}

	// token generate
	secret := []byte(config.ENV.JWT_SECRET)

	token, err := auth.CreateToken(secret, user.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusCreated, map[string]string{"token": token})
}
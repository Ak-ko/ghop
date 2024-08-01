package user

import (
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
	router.HandleFunc("/register", utils.MakeHTTPHandler(h.handleRegister))
	router.HandleFunc("/login", utils.MakeHTTPHandler(h.handleLogin))
}



func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return utils.ApiError{Err: "method not allowed", Status: http.StatusMethodNotAllowed}
	}

	var payload types.RegisterUserPayload

	// parsing
	if err := utils.ParseJSON(r, &payload); err != nil {
		return utils.ApiError{Err: err.Error(), Status: http.StatusBadRequest}
	}

	// validating
	if err := utils.Validator.Struct(payload); err != nil {
		vError := err.(validator.ValidationErrors)
		return utils.ApiError{Err: vError.Error(), Status: http.StatusBadRequest}
	}

	// check user exists
	_,err := h.store.GetUserByEmail(payload.Email)
	if err == nil {
		return utils.ApiError{Err: "user existed", Status: http.StatusBadRequest}
	}

	// hashing
	hashedPassword, err := auth.HashPassword(payload.Password)
	 if err != nil {
		return utils.ApiError{Err: err.Error(), Status: http.StatusInternalServerError}
	}

	// creating a user
	if err = h.store.CreateUser(types.User{
		Username: payload.Username,
		Email: payload.Email,
		Password: hashedPassword,
	}); err != nil {
		return utils.ApiError{Err: err.Error(), Status: http.StatusInternalServerError}
	}

	return utils.WriteJson(w, http.StatusOK, nil) // token
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return utils.ApiError{Err: "method not allowed", Status: http.StatusMethodNotAllowed}
	}

	var payload types.LoginUserPayload

	// parsing
	if err := utils.ParseJSON(r, &payload); err != nil {
		return utils.ApiError{Err: err.Error(), Status: http.StatusBadRequest}
	}

	// validate
	if err := utils.Validator.Struct(payload); err != nil {
		vError := err.(validator.ValidationErrors)
		return utils.ApiError{Err: vError.Error(), Status: http.StatusBadRequest}
	}

	// check user exists
	user, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		return utils.ApiError{Err: "invalid credentials", Status: http.StatusBadRequest}
	}

	// password check
	pwMatch := auth.ComparePasswords(user.Password, []byte(payload.Password))
	if !pwMatch {
		return utils.ApiError{Err: "incorrect password", Status: http.StatusBadRequest}
	}

	// token generate
	secret := []byte(config.ENV.JWT_SECRET)

	token, err := auth.CreateToken(secret, user.ID)
	if err != nil {
		return utils.ApiError{Err: err.Error(), Status: http.StatusInternalServerError}
	}

	return utils.WriteJson(w, http.StatusCreated, map[string]string{"token": token})
}
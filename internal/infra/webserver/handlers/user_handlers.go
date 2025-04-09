package handlers

import (
	"apis/internal/dto"
	"apis/internal/entity"
	"apis/internal/infra/database"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth/v5"
)

type Error struct {
	Message string `json:"message"`
}

type UserHandler struct {
	UserDB database.UserInterface
}

func NewUserHandler(db database.UserInterface, jwt *jwtauth.JWTAuth, jwtExpiresIn int) *UserHandler {
	return &UserHandler{
		UserDB: db,
	}
}

// GenerateJWT godoc
// @Summary Generate JWT
// @Description Get a user JWT
// @Tags users
// @Accept json
// @Produce json
// @Param user body dto.GetJWTInput true "User"
// @Success 200 {object} dto.AccessTokenOutput
// @Failure 404 {object} Error
// @Failure 500 {object} Error
// @Router /users/login [post]
func (h *UserHandler) GenerateJWT(w http.ResponseWriter, r *http.Request) {
	jwt, ok := r.Context().Value("jwt").(*jwtauth.JWTAuth)
	if !ok || jwt == nil {
		writeError(w, http.StatusUnauthorized, "JWT middleware not configured")
		return
	}

	jwtExpiresIn, ok := r.Context().Value("jwtExpiresIn").(int)
	if !ok {
		writeError(w, http.StatusInternalServerError, "JWT expiration config missing")
		return
	}

	var input dto.GetJWTInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := h.UserDB.GetByEmail(input.Email)
	if err != nil {
		writeError(w, http.StatusNotFound, "User not found")
		return
	}

	if err := user.ComparePassword(input.Password); err != nil {
		writeError(w, http.StatusUnauthorized, "Invalid password")
		return
	}

	_, tokenString, _ := jwt.Encode(map[string]interface{}{
		"sub": user.ID.String(),
		"exp": time.Now().Add(time.Second * time.Duration(jwtExpiresIn)).Unix(),
	})

	writeJSON(w, http.StatusOK, dto.AccessTokenOutput{AccessToken: tokenString})
}

// @Summary Create a new user
// @Description Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Param user body dto.CreateUserInput true "User"
// @Success 201 {object} dto.CreateUserInput
// @Failure 500 {object} Error
// @Router /users [post]
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var user dto.CreateUserInput
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	u, err := entity.NewUser(user.Name, user.Email, user.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}
	err = h.UserDB.Create(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Error{Message: msg})
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

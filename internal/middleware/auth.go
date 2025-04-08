package middleware

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth/v5"
)

func ProtectedRoutes(tokenAuth *jwtauth.JWTAuth, register func(r chi.Router)) http.Handler {
	r := chi.NewRouter()
	r.Use(jwtauth.Verifier(tokenAuth))
	r.Use(jwtauth.Authenticator(tokenAuth))
	register(r)
	return r
}

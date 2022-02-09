package rest

import (
	"net/http"

	"tripleoak/auth-api/rest/auth"
)

func Router() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("/v1/auth/login", auth.LoginHandler)
	router.HandleFunc("/v1/auth/logout", auth.LogoutHandler)
	router.HandleFunc("/v1/auth/signup", auth.SignupHandler)
	router.HandleFunc("/v1/auth/forgot-password", auth.ForgotPasswordHandler)

	return router
}

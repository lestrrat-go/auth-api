package auth

import (
	"encoding/json"
	"net/http"
	"strings"
	"tripleoak/auth-api/security"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var lr LoginRequest
		err := json.NewDecoder(r.Body).Decode(&lr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		jwt, err := Login(lr.Username, lr.Password)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
		} else {
			json.NewEncoder(w).Encode(jwt)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	jwtString := strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", -1)
	if jwtString == "" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Missing JWT"))
		return
	}

	jwtToken, err := security.VerifyJWT(jwtString)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid JWT"))
		return
	}

	sub, exists := jwtToken.Get("sub")
	if !exists {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid JWT"))
		return
	}

	switch r.Method {
	case "GET":
		err := Logout(sub.(string))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
	}
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var sr SignupRequest
		err := json.NewDecoder(r.Body).Decode(&sr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		err = Signup(sr.Username, sr.Password)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
	}
}

func ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		// Send email with reset link

		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
	}
}

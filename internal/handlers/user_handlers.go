package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/hoanguyen1998/crypto-payment-system/helpers"
	"github.com/pascaldekloe/jwt"
	"golang.org/x/crypto/bcrypt"
)

type RegisterPayload struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Jwt string `json:"jwt"`
}

func (s *ServerHandler) Register(w http.ResponseWriter, r *http.Request) {
	var userPayload RegisterPayload

	json.NewDecoder(r.Body).Decode(&userPayload)

	user, errRest := s.services.CreateUser(userPayload.Email, userPayload.Name, userPayload.Password)

	if errRest != nil {
		s.writeJSON(w, errRest.Status, errRest)
		return
	}

	s.writeJSON(w, http.StatusOK, user)
}

func (s *ServerHandler) Login(w http.ResponseWriter, r *http.Request) {
	var user LoginPayload

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		return
	}

	userInfo, errRest := s.services.GetUser(user.Email, user.Password)

	if errRest != nil {
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(userInfo.PasswordHash), []byte(user.Password))
	if err != nil {
		fmt.Println("password does not match")
		s.writeJSON(w, http.StatusBadRequest, helpers.NewBadRequestError("password does not match"))
		return
	}

	fmt.Println("password match")
	var claims jwt.Claims
	claims.Subject = fmt.Sprint(userInfo.Id)
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(time.Now().Add(24 * time.Hour))
	claims.Issuer = "mydomain.com"
	claims.Audiences = []string{"mydomain.com"}

	jwtBytes, err := claims.HMACSign(jwt.HS256, []byte("app.config.jwt.secret"))
	if err != nil {
		return
	}

	jwtRes := LoginResponse{Jwt: string(jwtBytes)}

	s.writeJSON(w, http.StatusOK, jwtRes)
}

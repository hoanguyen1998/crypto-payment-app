package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hoanguyen1998/crypto-payment-system/helpers"
	"github.com/pascaldekloe/jwt"
)

type contextKey int

const authUserId contextKey = 0

func (s *ServerHandler) CheckToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			s.writeJSON(w, http.StatusForbidden, helpers.NewForbiddenError("invalid auth header"))
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 {
			s.writeJSON(w, http.StatusForbidden, helpers.NewForbiddenError("invalid auth header"))
			return
		}

		if headerParts[0] != "Bearer" {
			s.writeJSON(w, http.StatusForbidden, helpers.NewForbiddenError("unauthorized - no bearer"))
			return
		}

		token := headerParts[1]

		claims, err := jwt.HMACCheck([]byte(token), []byte("ds")) //app.config.jwt.secret
		if err != nil {
			s.writeJSON(w, http.StatusForbidden, helpers.NewForbiddenError("unauthorized - failed hmac check"))
			return
		}

		if !claims.Valid(time.Now()) {
			s.writeJSON(w, http.StatusForbidden, helpers.NewForbiddenError("unauthorized - token expired"))
			return
		}

		if !claims.AcceptAudience("mydomain.com") {
			s.writeJSON(w, http.StatusForbidden, helpers.NewForbiddenError("unauthorized - invalid audience"))
			return
		}

		if claims.Issuer != "mydomain.com" {
			s.writeJSON(w, http.StatusForbidden, helpers.NewForbiddenError("unauthorized - invalid issuer"))
			return
		}

		userID, err := strconv.ParseInt(claims.Subject, 10, 64)
		if err != nil {
			s.writeJSON(w, http.StatusForbidden, helpers.NewForbiddenError("unauthorized"))
			return
		}

		fmt.Println("Valid user: ", userID)
		ctx := context.WithValue(r.Context(), authUserId, userID)

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

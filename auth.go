package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func WithJWTAuth(handlerFunc http.HandlerFunc, store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := FetchTokenFromRequest(r)

		token, err := validateToken(tokenString)

		if err != nil {
			log.Println("failed to authenticate token")
			denyPermission(w)
			return
		}

		if !token.Valid {
			log.Println("failed to authenticate token")
			denyPermission(w)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		userID := claims["userID"].(string)

		_, err = store.GetUserByID(userID)

		if err != nil {
			log.Println("failed to get user")
			denyPermission(w)
			return
		}

		handlerFunc(w, r)
	}
}

func denyPermission(w http.ResponseWriter) {
	WriteJSON(w, http.StatusUnauthorized, 
		ErrorResponse{Error: fmt.Errorf("permission denied").Error()})
}

func FetchTokenFromRequest(r *http.Request) string {
	authToken := r.Header.Get("Authorisation")
	queryToken := r.URL.Query().Get("token")

	if authToken != "" {
		return authToken
	}

	if queryToken != "" {
		return queryToken
	}

	return ""
}

func validateToken(t string) (*jwt.Token, error) {
	secret := Envs.JWSecret

	return jwt.Parse(t, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(secret), nil
	})
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func CreateJWT(secret []byte, userID int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": strconv.Itoa(int(userID)),
		"expiresAt": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString(secret)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
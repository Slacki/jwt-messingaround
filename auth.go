package main

import (
	"bio-cleaner-api/models"
	"context"
	"crypto/rsa"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	jwtRequest "github.com/dgrijalva/jwt-go/request"
	"golang.org/x/crypto/bcrypt"
)

const (
	privateKeyPath = "keys/app.rsa"
	publicKeyPath  = "keys/app.rsa.pub"
	crtPath        = "keys/server.crt"
	crtKeyPath     = "keys/server.key"
)

var (
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
)

func initKeys() {
	signBytes, err := ioutil.ReadFile(privateKeyPath)
	fatal(err)
	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	fatal(err)

	verifyBytes, err := ioutil.ReadFile(publicKeyPath)
	fatal(err)
	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	fatal(err)
}

type userCredentials struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}

type tokenResponse struct {
	Token string `json:"token"`
}

type jwtClaims struct {
	Identifier string `json:"identifier"`
	jwt.StandardClaims
}

type ctxClaims string

func tlsMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		h.ServeHTTP(w, r)
	})
}

func authenticatedMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := jwtRequest.ParseFromRequestWithClaims(r, jwtRequest.AuthorizationHeaderExtractor, &jwtClaims{}, func(token *jwt.Token) (interface{}, error) {
			return verifyKey, nil
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if claims, ok := token.Claims.(*jwtClaims); ok && token.Valid {
			ctx := context.WithValue(r.Context(), ctxClaims("claims"), claims)
			h.ServeHTTP(w, r.WithContext(ctx))
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
	})
}

func handle_auth(w http.ResponseWriter, r *http.Request) {
	var userCreds userCredentials
	err := json.NewDecoder(r.Body).Decode(&userCreds)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	userCreds.ID = strings.ToLower(userCreds.ID)

	user := models.User{}
	if db.Where("identifier = ?", userCreds.ID).First(&user).RecordNotFound() {
		http.Error(w, "Bad credentials", http.StatusUnauthorized)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(userCreds.Password))
	if err != nil {
		http.Error(w, "Bad credentials", http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS512, jwtClaims{
		user.Identifier,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
			Issuer:    "Bio-Cleaner24 API",
			Subject:   user.Identifier,
		},
	})
	tokenString, err := token.SignedString(signKey)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response := tokenResponse{tokenString}
	jsonResponse(response, w)
}

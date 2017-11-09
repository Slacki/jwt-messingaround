package main

import (
	"bio-cleaner-api/models"
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

//  $ openssl genrsa -out app.rsa 4096
//  $ openssl rsa -in app.rsa -pubout > app.rsa.pub

const (
	privateKeyPath = "keys/app.rsa"
	publicKeyPath  = "keys/app.rsa.pub"
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

func authenticatedMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := jwtRequest.ParseFromRequest(r, jwtRequest.AuthorizationHeaderExtractor, func(token *jwt.Token) (interface{}, error) {
			return verifyKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func auth(w http.ResponseWriter, r *http.Request) {
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

	token := jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(60 * time.Second).Unix(),
		Issuer:    "Bio-Cleaner24 API",
		Subject:   user.Identifier,
	})
	tokenString, err := token.SignedString(signKey)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response := tokenResponse{tokenString}
	jsonResponse(response, w)
}

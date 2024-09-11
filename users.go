package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/DuganChandler/goserver/internal/database"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (cfg *apiConfig) createUsersHandler(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error decoding request body")
		return
	}

	password, err := bcrypt.GenerateFromPassword([]byte(params.Password), 14)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to hash password")
		return
	}

	user, err := cfg.DB.CreateUsers(params.Email, string(password))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating user")
		return
	}

	responseWithJSON(w, http.StatusCreated, database.User{
		Id:    user.Id,
		Email: user.Email,
	})
}

func (cfg *apiConfig) updateUsersLoginHandler(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	tokenString := strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer ")

	log.Print(tokenString)
	type customClaims struct {
		jwt.RegisteredClaims
	}

	token, err := jwt.ParseWithClaims(tokenString, &customClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWTSecret), nil
	})

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	subject, err := token.Claims.GetSubject()
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unable to get subect from token")
		return
	}

	userID, err := strconv.Atoi(subject)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unable to turn id into int")
		return
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error decoding request body")
		return
	}

	password, err := bcrypt.GenerateFromPassword([]byte(params.Password), 14)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to hash password")
		return
	}

	user, err := cfg.DB.UpdateUserLogin(params.Email, string(password), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to update user login info")
		return
	}

	responseWithJSON(w, http.StatusOK, database.User{
		Id:    user.Id,
		Email: user.Email,
		Token: user.Token,
	})
}

func (cfg *apiConfig) loginUsersHadler(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error decoding request body")
		return
	}

	user, err := cfg.DB.GetUserByEmail(params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "you are unauthorized")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password))
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "you are unauthorized")
		return
	}

    expiration := time.Duration(0 * time.Second)
	if params.ExpiresInSeconds == 0 || params.ExpiresInSeconds > 24 {
		expiration = time.Duration(24 * time.Hour)
	} else {
		expiration = time.Duration(params.ExpiresInSeconds * int(time.Second))
	}

	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration).UTC()),
		Subject:   strconv.Itoa(user.Id),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signature, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unable to create jwt signature")
	}

	responseWithJSON(w, http.StatusOK, database.User{
		Id:    user.Id,
		Email: user.Email,
		Token: signature,
	})
}

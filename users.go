package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/DuganChandler/goserver/internal/auth"
	"github.com/DuganChandler/goserver/internal/database"
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

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
	}

	user, err := cfg.DB.CreateUsers(params.Email, string(hashedPassword))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating user")
		return
	}

	responseWithJSON(w, http.StatusCreated, database.User{
		Id:          user.Id,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	})
}

func (cfg *apiConfig) updateUsersLoginHandler(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	tokenString, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "could not find JWT")
		return
	}

	subject, err := auth.VerifyJWT(tokenString, cfg.JWTSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("unable to verify jwt: %s", err))
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

	userIDInt, err := strconv.Atoi(subject)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not parse user ID")
		return
	}

	user, err := cfg.DB.UpdateUserLogin(params.Email, string(password), userIDInt)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to update user login info")
		return
	}

	responseWithJSON(w, http.StatusOK, database.User{
		Id:          user.Id,
		Email:       user.Email,
		Token:       user.Token,
		IsChirpyRed: user.IsChirpyRed,
	})
}

func (cfg *apiConfig) loginUsersHadler(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		ID           int    `json:"id"`
		Email        string `json:"email"`
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
		IsChirpyRed  bool   `json:"is_chirpy_red"`
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

	err = auth.CheckPasswordHash(params.Password, user.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "you are unauthorized")
		return
	}

	duration := time.Duration(1 * time.Hour)
	jwtToken, err := auth.MakeJWT(user.Id, cfg.JWTSecret, duration)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unable to create jwt signature")
	}

	refreshToken, err := auth.CreateNewRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unable to create new refresh token")
	}

	err = cfg.DB.StoreRefreshToken(refreshToken, user.Id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to store refresh token")
		return
	}

	responseWithJSON(w, http.StatusOK, response{
		ID:           user.Id,
		Email:        user.Email,
		Token:        jwtToken,
		RefreshToken: refreshToken,
        IsChirpyRed: user.IsChirpyRed,
	})
}


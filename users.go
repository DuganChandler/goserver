package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
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
		Id:    user.Id,
		Email: user.Email,
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
		Id:    user.Id,
		Email: user.Email,
		Token: user.Token,
	})
}

func (cfg *apiConfig) loginUsersHadler(w http.ResponseWriter, req *http.Request) {
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

	rToken := make([]byte, 32)
	_, err = rand.Read(rToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unable to create refresh token")
	}
	rTokenEncoded := hex.EncodeToString(rToken)

	refreshToken := database.RefreshToken{
		Token:    rTokenEncoded,
		Duration: time.Duration((24 * 60) * time.Hour),
	}

	_, err = cfg.DB.StoreRefreshToken(refreshToken, user.Id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to store refresh token")
		return
	}

	responseWithJSON(w, http.StatusOK, database.User{
		Id:           user.Id,
		Email:        user.Email,
		Token:        jwtToken,
		RefreshToken: database.RefreshToken{Token: rTokenEncoded},
	})
}

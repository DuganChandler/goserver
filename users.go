package main

import (
	"encoding/json"
	"net/http"

    "golang.org/x/crypto/bcrypt"
	"github.com/DuganChandler/goserver/internal/database"
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

    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password))    
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "you are unauthorized")
        return
    }

    responseWithJSON(w, http.StatusOK, database.User{
        Id: user.Id,
        Email: user.Email,
    })

}

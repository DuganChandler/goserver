package main

import (
	"net/http"
	"time"

	"github.com/DuganChandler/goserver/internal/auth"
)

func (cfg *apiConfig) refreshTokenHandler(w http.ResponseWriter, req *http.Request) {
    type response struct {
        Token string `json:"token"`
    }

	refreshTokenString, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "could not find refresh token")
		return
	}

	user, err := cfg.DB.GetUserByRefreshToken(refreshTokenString)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "no user matching refresh token")
		return
	}

	newToken, err := auth.MakeJWT(user.Id, cfg.JWTSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unable to create new JWT token")
		return
	}

	responseWithJSON(w, http.StatusOK, response{
        Token: newToken,
    })

}

func (cfg *apiConfig) revokeTokenHandler(w http.ResponseWriter, req *http.Request) {
	refreshTokenString, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "could not find refresh token")
		return
	}


    err = cfg.DB.RevokeRefreshToken(refreshTokenString)
    if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unable to revoke jwt")
		return
    }

    w.WriteHeader(http.StatusNoContent)
}

package database

import (
	"fmt"
	"time"
)

func (db *DB) StoreRefreshToken(token string, userID int) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	refreshToken := RefreshToken{
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(time.Hour),
	}
	dbStructure.RefreshTokens[token] = refreshToken

	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) RevokeRefreshToken(token string) error {
    dbStructure, err := db.loadDB()
    if err != nil {
        return err
    }
    delete(dbStructure.RefreshTokens, token)
    

    err = db.writeDB(dbStructure)
    if err != nil {
        return fmt.Errorf("unable to write revoked token to db")
    }
    
    return nil
}

func (db *DB) GetUserByRefreshToken(tokenString string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	refreshToken, ok := dbStructure.RefreshTokens[tokenString]
	if !ok {
		return User{}, fmt.Errorf("Token does not exist")
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		return User{}, fmt.Errorf("token has expired")
	}

	user, err := db.GetUserByID(refreshToken.UserID)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

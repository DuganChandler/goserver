package database

import (
	"fmt"
)

func (db *DB) CreateUsers(email, password string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	_, ok := searchUserByEmail(dbStructure, email)
	if ok {
		return User{}, fmt.Errorf("user with email already exists")

	}

	userID := len(dbStructure.Users) + 1
	user := User{
		Id:          userID,
		Email:       email,
		Password:    password,
		IsChirpyRed: false,
	}

	dbStructure.Users[userID] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *DB) GetUserByID(id int) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user, ok := dbStructure.Users[id]
	if !ok {
		return User{}, fmt.Errorf("user does not exist")
	}

	return user, nil
}

func (db *DB) UpdateUserLogin(email, password string, id int) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user, ok := dbStructure.Users[id]
	if !ok {
		return User{}, fmt.Errorf("user does not exist")
	}

	updatedUser := User{
		Id:          user.Id,
		Email:       email,
		Password:    password,
		Token:       user.Token,
		IsChirpyRed: user.IsChirpyRed,
	}

	dbStructure.Users[id] = updatedUser

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}
	return updatedUser, nil
}

func (db *DB) GetUserByEmail(email string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user, ok := searchUserByEmail(dbStructure, email)
	if !ok {
		return User{}, fmt.Errorf("user does not exist")
	}

	return user, nil
}

func (db *DB) UpgradeUser(userID int) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	user, ok := dbStructure.Users[userID]
	if !ok {
		return fmt.Errorf("user does not exist")
	}

	updatedUser := User{
		Id:          user.Id,
		Email:       user.Email,
		Password:    user.Password,
		Token:       user.Email,
		IsChirpyRed: true,
	}

	dbStructure.Users[userID] = updatedUser

	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}

	return nil
}

func searchUserByEmail(dbStructure DBStructure, email string) (User, bool) {
	for _, user := range dbStructure.Users {
		if user.Email == email {
			return user, true
		}
	}
	return User{}, false
}

package database

import (
	"fmt"
)

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string, userID int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, fmt.Errorf("unable to load db: %s", err)
	}

	id := len(dbStructure.Chirps) + 1

	chirp := Chirp{
		Id:       id,
		Body:     body,
		AuthorID: userID,
	}
	dbStructure.Chirps[id] = chirp

	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, fmt.Errorf("unable to write to db: %s", err)
	}

	return chirp, nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return []Chirp{}, err
	}

	chirps := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, val := range dbStructure.Chirps {
		chirps = append(chirps, val)
	}

	return chirps, nil
}

func (db *DB) GetChirpByID(id int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	chirp, ok := dbStructure.Chirps[id]
	if !ok {
		return Chirp{}, fmt.Errorf("No chirp matching the id: %d", id)
	}

	return chirp, nil
}

func (db *DB) GetChirpsByAuthor(id int) ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return []Chirp{}, err
	}

	chirps := make([]Chirp, 0)
	for _, val := range dbStructure.Chirps {
		if val.AuthorID == id {
			chirps = append(chirps, val)
		}
	}

	if len(chirps) == 0 {
		return chirps, fmt.Errorf("no chirps matching that author_id")
	}

	return chirps, nil
}

func (db *DB) DeleteChirpByID(id int) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	_, ok := dbStructure.Chirps[id]
	if !ok {
		return fmt.Errorf("unable to find chirp with provided id")
	}

	delete(dbStructure.Chirps, id)

	for key, val := range dbStructure.Chirps {
		if key > id {
			dbStructure.Chirps[key-1] = val
			delete(dbStructure.Chirps, key)
		}
	}

	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}

	return nil
}

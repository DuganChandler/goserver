package database

import (
	"fmt"
)

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, fmt.Errorf("unable to load db: %s", err)
	}

	id := len(dbStructure.Chirps) + 1

	chirp := Chirp{
		Id:   id,
		Body: body,
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

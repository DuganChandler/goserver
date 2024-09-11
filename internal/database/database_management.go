package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
)

// create new db
func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}
	err := db.ensureDB()

	return db, err
}

// ensure db exists
func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
    if errors.Is(err, os.ErrNotExist) {
        db.createDB()
    }

	return err
}

func (db *DB) createDB() error {
	dbStructure := DBStructure{
		Chirps: map[int]Chirp{},
        Users: map[int]User{},
	}
	return db.writeDB(dbStructure)
}

func (db *DB) loadDB() (DBStructure, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbStructure := DBStructure{}
	data, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return dbStructure, err
	}

	err = json.Unmarshal(data, &dbStructure)
	if err != nil {
		return dbStructure, fmt.Errorf("unable to unmarshal json while loading db: %s", err)
	}

	return dbStructure, nil
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	data, err := json.Marshal(dbStructure)
	if err != nil {
		return fmt.Errorf("unable to marshal json while writing to db: %s", err)
	}

	err = os.WriteFile(db.path, data, 0600)
	if err != nil {
		return fmt.Errorf("unable to write to file: %s", err)
	}

	return nil
}

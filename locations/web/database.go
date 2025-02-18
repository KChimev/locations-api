package main

import (
	"database/sql"

	"github.com/kchimev/locations-api/locations/internal/constants"
	_ "github.com/lib/pq"
)

func GetDBPool() (*sql.DB, error) {
	db, err := sql.Open("postgres", constants.DBConnectionString)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func GetPOIDbPool() (*sql.DB, error) {
	db, err := sql.Open("postgres", constants.POIConnectionString)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

package model

import "github.com/chrissnell/chickenlittle/db"

// Model contains the data model with the associated DB connection handle
type Model struct {
	db *db.DB
}

// New creates a new data model with the given DB connection handle
func New(db *db.DB) *Model {
	m := &Model{
		db: db,
	}
	return m
}

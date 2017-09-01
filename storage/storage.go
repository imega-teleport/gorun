package storage

import (
	"database/sql"
)

type Storage interface {
	GetGroups(out chan<- interface{}, e chan<- error)
	GetProducts(out chan<- interface{}, e chan<- error)
	GetProductsGroups(out chan<- interface{}, e chan<- error)
	GetProductsProperties(out chan<- interface{}, e chan<- error, condition []string)
	GetProductsPropertiesSpecial(out chan<- interface{}, e chan<- error, condition []string)
}

type storage struct {
	db    *sql.DB
	limit int
}

func NewStorage(db *sql.DB, limit int) Storage {
	return &storage{
		db:    db,
		limit: limit,
	}
}

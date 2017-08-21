package storage

import (
	"database/sql"
	//"sync"
)

type Storage interface {
	GetGroups(out chan<- interface{}, e chan error)
	GetProducts(out chan<- interface{}, e chan error)
}

type storage struct {
	db *sql.DB
	//wg *sync.WaitGroup
}

func NewStorage(db *sql.DB) Storage {
	return &storage{
		db: db,
		//wg: wg,
	}
}

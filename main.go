package main

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/imega-teleport/gorun/storage"
	log "github.com/sirupsen/logrus"
	"github.com/imega-teleport/gorun/packer"
)

func main() {
	log.Info("Start")
	dsn := fmt.Sprintf("mysql://%s:%s@tcp(%s)/%s", "", "", "10.0.3.102:3306", "test_teleport")
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Could not connect db, %s", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Fail ping db, %s", err)
	}
	defer func() {
		err = db.Close()
		if err != nil {
			log.Fatalf("Fail closes db connection, %s", err)
		}
		log.Info("Closed db connection")
	}()

	wg := sync.WaitGroup{}
	s := storage.NewStorage(db)

	dataChan := make(chan interface{}, 10)
	errChan := make(chan error)

	p := packer.New(500)

	wg.Add(1)
	go func() {
		defer wg.Done()
		s.GetGroups(dataChan, errChan)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		s.GetProducts(dataChan, errChan)
	}()

	go func() {
		p.Listen(dataChan, errChan)
	}()

	go func() {
		wg.Wait()
		close(dataChan)
		close(errChan)
	}()

	if err := <-errChan; err != nil {
		log.Fatalf("%s", err)
		close(dataChan)
		close(errChan)
		return
	}
}

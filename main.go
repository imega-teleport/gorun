package main

import (
	"database/sql"
	"fmt"
	"os"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/imega-teleport/gorun/storage"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("Start")
	dsn := fmt.Sprintf("mysql://%s:%s@tcp(%s)/%s", "", "", "10.0.3.32:3306", "test_teleport")
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("error: %s", err)
		os.Exit(1)
	}

	err = db.Ping()
	if err != nil {
		fmt.Printf("error: %s", err)
		os.Exit(1)
	}
	defer func() {
		err = db.Close()
		if err != nil {
			fmt.Printf("error: %s", err)
			os.Exit(1)
		}
		fmt.Println("Closed db connection")
	}()

	//wg := &sync.WaitGroup{}
	var wg sync.WaitGroup
	s := storage.NewStorage(db)

	dataChan := make(chan interface{}, 10)
	errChan := make(chan error, 1)

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
		err := errHandle(errChan)
		if err != nil {
			log.Fatalf("%s", err)
		}
	}()

	go printer(dataChan)

	//go func() {
	//
	//	fmt.Println("=====")
	//	close(dataChan)
	//	}()
	wg.Wait()
}

func printer(in <-chan interface{}) {
	for v := range in {
		//time.Sleep(time.Second)
		switch v.(type) {
		case storage.Product:
			fmt.Println("Product: ", v.(storage.Product).Name)
		case storage.Group:
			fmt.Println("Group: ", v.(storage.Group).Name)
		}
	}
}

func errHandle(in <-chan error) error {
	for v := range in {
		return v
	}
	return nil
}

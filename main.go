package main

import (
	"database/sql"
	"fmt"
	"os"

	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/imega-teleport/gorun/storage"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("Start")
	dsn := fmt.Sprintf("mysql://%s:%s@tcp(%s)/%s", "", "", "10.0.3.102:3306", "test_teleport")
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

	s := storage.NewStorage(db)

	groupChan := make(chan interface{}, 10)
	errChan := make(chan error)

	go s.GetGroups(groupChan, errChan)

	printer(groupChan)
}

var groups = struct {
	items  []storage.Group
	length int
}{}

func printer(in <-chan interface{}) {
	for v := range in {
		time.Sleep(time.Second)

		groups.length = groups.length + len(v.(storage.Group).ID) + len(v.(storage.Group).Name) + len(v.(storage.Group).ParentID)

		groups.items = append(groups.items, v.(storage.Group))
		//fmt.Println(groups)
		fmt.Println(groups.length)
	}
}

package main // import "github.com/imega-teleport/gorun"

import (
	"database/sql"
	"fmt"
	"sync"
	"flag"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/imega-teleport/gorun/packer"
	"github.com/imega-teleport/gorun/storage"
	log "github.com/sirupsen/logrus"
)

func main() {
	user, pass, host := os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST")

	dbname := flag.String("db", "test_teleport", "Database name")
	path := flag.String("path", "/tmp", "Save to path")
	limit := flag.Int("limit", 500000, "Limit bytes")
	prefixTable := flag.String("ptable", "wp_", "Prefix table name")
	prefixFile := flag.String("pfile", "out", "Prefix file name")
	flag.Parse()

	dsn := fmt.Sprintf("mysql://%s:%s@tcp(%s)/%s", user, pass, host, *dbname)

	log.Info("Start")
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

	p := packer.New(packer.Options{
		MaxBytes:        *limit,
		PrefixFileName:  *prefixFile,
		PathToSave:      *path,
		PrefixTableName: *prefixTable,
	})

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

	wg.Add(1)
	go func() {
		defer wg.Done()
		s.GetProductsGroups(dataChan, errChan)
	}()

	go func() {
		p.Listen(dataChan, errChan)
	}()

	go func() {
		wg.Wait()
		p.SaveToFile()
		p.SecondSaveToFile()
		p.ThirdPackSaveToFile()
		close(dataChan)
		close(errChan)
		log.Info("End work!")
	}()

	if err := <-errChan; err != nil {
		log.Fatalf("%s", err)
		close(dataChan)
		close(errChan)
		return
	}
}

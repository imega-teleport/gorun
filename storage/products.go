package storage

import "database/sql"

type Product struct {
	ID          string
	Name        string
	Description string
	Barcode     string
	Article     string
	FullName    string
	Country     string
	Brand       string
}

func (s *storage) GetProducts(out chan<- interface{}, e chan error) {
	/*for v := range e {
		e <- v
		return
	}*/
	s.getRecords(out, e, "select id, name, description, barcode, article, full_name, country, brand from products", func(rows *sql.Rows) (interface{}, error) {
		item := Product{}
		err := rows.Scan(&item.ID, &item.Name, &item.Description, &item.Barcode, &item.Article, &item.FullName, &item.Country, &item.Brand)
		return item, err
	})
}

package storage

import "database/sql"

type ProductsGroups struct {
	ProductID string
	GroupID   string
}

func (s *storage) GetProductsGroups(out chan<- interface{}, e chan<- error) {
	s.getRecords(out, e, "select product_id, group_id from products_groups", func(rows *sql.Rows) (interface{}, error) {
		i := ProductsGroups{}
		err := rows.Scan(&i.ProductID, &i.GroupID)
		return i, err
	})
}

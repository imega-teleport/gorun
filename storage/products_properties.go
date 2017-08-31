package storage

import (
	"database/sql"
	"fmt"
	"strings"
)

type ProductsProperties struct {
	ProductID    string
	PropertyID   string
	Value        string
	PropertyName string
	PropertyType int
}

type ProductsPropertiesSpecial struct {
	ProductID    string
	PropertyID   string
	Value        string
	PropertyName string
	PropertyType int
}

func (s *storage) GetProductsProperties(out chan<- interface{}, e chan<- error, condition []string) {
	query := "select pp.parent_id product_id, pp.id property_id, pp.value value, i.name property_name, i.type property_type from products_properties pp join properties i on pp.id=i.id"
	cond := []string{}

	for _, v := range condition {
		if v != "" {
			cond = append(cond, fmt.Sprintf("pp.id<>'%s'", v))
		}
	}

	if len(cond) > 0 {
		query = fmt.Sprintf("%s where %s", query, strings.Join(cond, " AND "))
	}

	s.getRecords(out, e, query, func(rows *sql.Rows) (interface{}, error) {
		i := ProductsProperties{}
		err := rows.Scan(&i.ProductID, &i.PropertyID, &i.Value, &i.PropertyName, &i.PropertyType)
		return i, err
	})
}

func (s *storage) GetProductsPropertiesSpecial(out chan<- interface{}, e chan<- error, condition []string) {
	query := "select pp.parent_id product_id, pp.id property_id, pp.value value, i.name property_name, i.type property_type from products_properties pp join properties i on pp.id=i.id"
	cond := []string{}

	for _, v := range condition {
		if v != "" {
			cond = append(cond, fmt.Sprintf("pp.id='%s'", v))
		}
	}

	if len(cond) > 0 {
		return
	}

	query = fmt.Sprintf("%s where %s", query, strings.Join(cond, " OR "))

	s.getRecords(out, e, query, func(rows *sql.Rows) (interface{}, error) {
		i := ProductsPropertiesSpecial{}
		err := rows.Scan(&i.ProductID, &i.PropertyID, &i.Value, &i.PropertyName, &i.PropertyType)
		return i, err
	})
}

package storage

import (
	"database/sql"
	"fmt"
)

func (s *storage) getRecords(out chan<- interface{}, e chan<- error, dql string, scan func(rows *sql.Rows) (interface{}, error)) {
	i := 0
	for {
		hasResults := false
		rows, err := s.db.Query(fmt.Sprintf("%s limit %d, %d", dql, i, s.limit))
		if err != nil {
			e <- err
			break
		}
		for rows.Next() {
			hasResults = true
			data, err := scan(rows)
			if err != nil {
				e <- err
			}
			out <- data
		}

		if hasResults {
			i = i + s.limit
		} else {
			break
		}
	}
}

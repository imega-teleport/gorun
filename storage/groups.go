package storage

import "database/sql"

type Group struct {
	ID       string
	ParentID string
	Name     string
}

func (s *storage) GetGroups(out chan<- interface{}, e chan error) {
	/*for v := range e {
		e <- v
		return
	}*/
	s.getRecords(out, e, "select id, parent_id, name from groups", func(rows *sql.Rows) (interface{}, error) {
		g := Group{}
		err := rows.Scan(&g.ID, &g.ParentID, &g.Name)
		return g, err
	})
}

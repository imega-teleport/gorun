package storage

import (
	"database/sql"
)

type Terms struct {
	ID    string `teleport:"term_id"`
	Name  string `teleport:"name"`
	Slug  string `teleport:"slug"`
	Group string `teleport:"term_group"`
}

func (t Terms) SizeOf() int {
	return len(t.ID) + len(t.Name) + len(t.Slug) + len(t.Group)
}

func (s *storage) GetGroups(out chan<- interface{}, e chan<- error) {
	s.getRecords(out, e, "select id,parent_id,name from groups", func(rows *sql.Rows) (interface{}, error) {
		g := struct {
			ID       string
			ParentID string
			Name     string
		}{}
		err := rows.Scan(&g.ID, &g.ParentID, &g.Name)
		return Group{
			ID:       g.ID,
			ParentID: g.ParentID,
			Name:     g.Name,
		}, err
	})
}

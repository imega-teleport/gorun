package teleport

import (
	"fmt"

	slugmaker "github.com/gosimple/slug"
	"gopkg.in/Masterminds/squirrel.v1"
)

type Slug string

func (s Slug) String() string {
	return slugmaker.Make(string(s))
}

type Term struct {
	ID    UUID
	Name  string
	Slug  Slug
	Group string
}

func (t Term) SizeOf() int {
	return len(t.ID) + len(t.Name) + len(t.Slug) + len(t.Group)
}

func (w *Wpwc) BuilderTerm() builder {
	return builder{
		squirrel.Insert(w.Prefix+"terms").Columns("term_id", "name", "slug", "term_group"),
	}
}

func (b *builder) AddTerm(t Term) {
	*b = builder{
		b.Values(squirrel.Expr(fmt.Sprintf("@max_term_id+%s", t.ID.ToVar())), t.Name, t.Slug.String(), 0),
	}
}

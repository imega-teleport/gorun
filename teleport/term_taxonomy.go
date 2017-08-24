package teleport

import (
	"fmt"

	squirrel "gopkg.in/Masterminds/squirrel.v1"
)

type taxonomyID int

func (i taxonomyID) String() string {
	ret := fmt.Sprintf("@max_term_taxonomy_id+%d", i)
	if i == 0 {
		ret = "0"
	}
	return ret
}

type TermTaxonomy struct {
	ID          taxonomyID
	TermID      UUID
	Taxonomy    string
	Description string
	Parent      UUID
}

func (t TermTaxonomy) SizeOf() int {
	return lengthDefineIndex + len(t.TermID) + len(t.Taxonomy) + len(t.Description) + len(t.Parent)
}

func (w *Wpwc) BuilderTermTaxonomy() builder {
	return builder{
		squirrel.Insert(w.Prefix+"term_taxonomy").Columns("term_taxonomy_id", "term_id", "taxonomy", "description", "parent"),
	}
}

func (b *builder) AddTermTaxonomy(t TermTaxonomy) {
	var parent interface{}
	if t.Parent != "" {
		parent = squirrel.Expr(t.Parent.ToVar())
	} else {
		parent = ""
	}

	*b = builder{
		b.Values(squirrel.Expr(t.ID.String()), squirrel.Expr(t.TermID.ToVar()), t.Taxonomy, t.Description, parent),
	}
}

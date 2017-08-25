package teleport

import squirrel "gopkg.in/Masterminds/squirrel.v1"

type TermRelationship struct {
	ObjectID       UUID
	TermTaxonomyID UUID
	TermOrder      int
}

func (t TermRelationship) SizeOf() int {
	return len(t.ObjectID) + len(t.TermTaxonomyID) + lengthDefineIndex
}

func (w *Wpwc) BuilderTermRelationships() builder {
	return builder{
		squirrel.Insert(w.Prefix+"term_relationships").Columns("object_id", "term_taxonomy_id", "term_order"),
	}
}

func (b *builder) AddTermRelationships(r TermRelationship) {
	*b = builder{
		b.Values(
			squirrel.Expr(r.ObjectID.ToVar()),
			squirrel.Expr(r.TermTaxonomyID.ToVar()),
			r.TermOrder,
		),
	}
}

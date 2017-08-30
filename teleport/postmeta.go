package teleport

import (
	"fmt"

	"gopkg.in/Masterminds/squirrel.v1"
)

type PostMeta struct {
	PostID UUID
	Key    string
	Value  string
}

func (p PostMeta) SizeOf() int {
	return len(p.PostID) + len(p.Key) + len(p.Value)
}

func (w *Wpwc) BuilderPostMeta() builder {
	return builder{
		squirrel.Insert(fmt.Sprintf("%spostmeta", w.Prefix)).Columns(
			"post_id",
			"meta_key",
			"meta_value",
		),
	}
}

func (b *builder) AddrPostMeta(i PostMeta) {
	*b = builder{
		b.Values(
			squirrel.Expr(i.PostID.ToVar()),
			i.Key,
			i.Value,
		),
	}
}

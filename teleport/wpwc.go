package teleport

import (
	"fmt"
	"strings"
	"time"

	slugmaker "github.com/gosimple/slug"
	"gopkg.in/Masterminds/squirrel.v1"
)

type Package struct {
	TeleportItem []TeleportItem
	Term         []Term
	Post         []Post
	Length       int
}

func (p *Package) AddItem(item interface{}) {
	switch item.(type) {
	case TeleportItem:
		p.Length = p.Length + item.(TeleportItem).SizeOf()
		p.TeleportItem = append(p.TeleportItem, item.(TeleportItem))
	case Term:
		p.Length = p.Length + item.(Term).SizeOf()
		p.Term = append(p.Term, item.(Term))
	case Post:
		p.Length = p.Length + item.(Post).SizeOf()
		p.Post = append(p.Post, item.(Post))
	}
}

type UUID string

func (id UUID) ToVar() string {
	return "@" + strings.Replace(slugmaker.Make(string(id)), "-", "", -1)
}

func (id UUID) String() string {
	return strings.Replace(string(id), "-", "", -1)
}

type Wpwc struct {
	Prefix string
}

type Term struct {
	ID    UUID
	Name  string
	Slug  string
	Group string
}

func (t Term) SizeOf() int {
	return len(t.ID) + len(t.Name) + len(t.Slug) + len(t.Group)
}

var dateLen = 19
var intLen = 5

type Post struct {
	ID       UUID
	AuthorID int
	Date     time.Time
	Content  string
	Title    string
	Excerpt  string
	Name     string
	Modified time.Time
}

func (p Post) SizeOf() int {
	return len(p.ID) + (dateLen * 4) + len(p.Name) + len(p.Title) + len(p.Content) + len(p.Excerpt)
}

type TeleportItem struct {
	GUID UUID
	Type string
	ID   int
	Date time.Time
}

func (t TeleportItem) SizeOf() int {
	return len(t.GUID) + len(t.Type) + intLen + dateLen
}

type builder struct {
	squirrel.InsertBuilder
}

func (w *Wpwc) BuilderTerm() builder {
	return builder{
		squirrel.Insert(w.Prefix+"terms").Columns("term_id", "name", "slug", "term_group"),
	}
}

func (b *builder) AddTerm(t Term) {
	*b = builder{
		b.Values(squirrel.Expr(t.ID.ToVar()), t.Name, t.Slug, 0),
	}
}

func (w *Wpwc) BuilderPost() builder {
	return builder{
		squirrel.Insert(fmt.Sprintf("%sposts", w.Prefix)).Columns(
			"id",
			"post_author",
			"post_date",
			"post_date_gmt",
			"post_content",
			"post_title",
			"post_excerpt",
			"post_name",
			"post_modified",
			"post_modified_gmt",
		),
	}
}

func (b *builder) AddPost(post Post) {
	*b = builder{
		b.Values(
			squirrel.Expr(post.ID.ToVar()),
			squirrel.Expr("1"),
			post.Date.String(),
			post.Date.UTC().String(),
			post.Content,
			post.Title,
			post.Excerpt,
			post.Name,
			post.Modified.String(),
			post.Modified.UTC().String(),
		),
	}
}

func (w *Wpwc) builderTeleportItem(prefix string) builder {
	return builder{
		squirrel.Insert(fmt.Sprintf("%steleport_item", w.Prefix)).Columns(
			"guid",
			"type",
			"id",
			"date",
		),
	}
}

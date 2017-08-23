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

const (
	lengthDefineVariable    = 0 //44 // ex. "set @d913f8c063a711e6a562005056b9f84b=949;"
	lengthDefineDate        = 22
	lengthDefineIndex       = 5
	lengthDefineSyntax      = 140
	lengthDefineValueSyntax = 13
)

func (p *Package) AddItem(item interface{}) {
	switch item.(type) {
	case TeleportItem:
		p.Length = p.Length + item.(TeleportItem).SizeOf() + lengthDefineValueSyntax
		p.TeleportItem = append(p.TeleportItem, item.(TeleportItem))
	case Term:
		p.Length = p.Length + item.(Term).SizeOf() + lengthDefineVariable + lengthDefineValueSyntax
		p.Term = append(p.Term, item.(Term))
	case Post:
		p.Length = p.Length + item.(Post).SizeOf() + lengthDefineVariable + lengthDefineValueSyntax
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
	return len(p.ID) + (lengthDefineDate * 4) + len(p.Name) + len(p.Title) + len(p.Content) + len(p.Excerpt) + lengthDefineIndex
}

type TeleportItem struct {
	GUID UUID
	Type string
	ID   int
	Date time.Time
}

func (t TeleportItem) SizeOf() int {
	return len(t.GUID) + len(t.Type) + lengthDefineIndex + lengthDefineDate
}

type builder struct {
	squirrel.InsertBuilder
}

func (w *Wpwc) BuilderTeleportItem() builder {
	return builder{
		squirrel.Insert(w.Prefix+"teleport_item").Columns("guid", "type", "id", "date"),
	}
}

func (b *builder) AddTeleportItem(i TeleportItem) {
	*b = builder{
		b.Values(i.GUID, i.Type, squirrel.Expr(fmt.Sprintf("@max_post_id+%s",i.GUID.ToVar())), i.Date.Format("2006-01-02 15:04:05")),
	}
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
			squirrel.Expr(fmt.Sprintf("@max_post_id+%s", post.ID.ToVar())),
			squirrel.Expr("1"),
			post.Date.Format("2006-01-02 15:04:05"),
			post.Date.UTC().Format("2006-01-02 15:04:05"),
			post.Content,
			post.Title,
			post.Excerpt,
			post.Name,
			post.Modified.Format("2006-01-02 15:04:05"),
			post.Modified.UTC().Format("2006-01-02 15:04:05"),
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

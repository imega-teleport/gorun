package teleport

import (
	"fmt"
	"strings"
	"time"

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
	}
}

type uuid string

func (id uuid) ToVar() string {
	return "@" + strings.Replace(slugmaker.Make(string(id)), "-", "", -1)
}

type Wpwc struct {
	Prefix string
}

type Term struct {
	ID    uuid
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
	ID       uuid
	AuthorID int
	Date     time.Time
	Content  string
	Title    string
	Excerpt  string
	Name     string
	Modified time.Time
}

type TeleportItem struct {
	GUID string
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

func (w *Wpwc) builderPost(prefix string) builder {
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

package teleport

import (
	"strings"

	slugmaker "github.com/gosimple/slug"
	"gopkg.in/Masterminds/squirrel.v1"
)

const (
	lengthDefineVariable    = 44 // ex. "set @d913f8c063a711e6a562005056b9f84b=949;"
	lengthDefineDate        = 22
	lengthDefineIndex       = 5
	lengthDefineSyntax      = 140
	lengthDefineValueSyntax = 13
)

type PrimaryPackage struct {
	TeleportItem []TeleportItem
	Term         []Term
	Post         []Post
	Length       int
}

func (p *PrimaryPackage) AddItem(item interface{}) {
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

type SecondaryPackage struct {
	TermTaxonomy []TermTaxonomy
	Length       int
}

func (p *SecondaryPackage) AddItem(item interface{}) {
	switch item.(type) {
	case TermTaxonomy:
		p.Length = p.Length + item.(TermTaxonomy).SizeOf() + (lengthDefineVariable * 2) + lengthDefineValueSyntax
		p.TermTaxonomy = append(p.TermTaxonomy, item.(TermTaxonomy))
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

type builder struct {
	squirrel.InsertBuilder
}

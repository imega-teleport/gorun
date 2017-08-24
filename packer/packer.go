package packer

import (
	"fmt"

	"time"

	"github.com/imega-teleport/gorun/indexer"
	"github.com/imega-teleport/gorun/storage"
	"github.com/imega-teleport/gorun/teleport"
	"github.com/imega-teleport/gorun/writer"
	"gopkg.in/Masterminds/squirrel.v1"
)

// Packer is interface
type Packer interface {
	Listen(in <-chan interface{}, e chan<- error)
	SaveToFile() error
	SecondSaveToFile() error
}

type Options struct {
	MaxBytes        int
	PrefixFileName  string
	PathToSave      string
	PrefixTableName string
}

type pkg struct {
	Options       Options
	PrimaryPack   teleport.PrimaryPackage
	SecondPack    teleport.SecondaryPackage
	Indexer       indexer.Indexer
	PackQty       int
	SecondPackQty int
	Content       string
}

// New instance packer
func New(opt Options) Packer {
	return &pkg{
		Options:       opt,
		Indexer:       indexer.NewIndexer(),
		PackQty:       1,
		SecondPackQty: 1,
	}
}

func (p *pkg) Listen(in <-chan interface{}, e chan<- error) {
	for v := range in {
		if p.IsFull(p.PrimaryPack) {
			p.SaveToFile()
			pack := teleport.PrimaryPackage{}
			p.Content = ""
			p.PrimaryPack = pack
			p.PackQty++
		}

		if p.SecondIsFull(p.SecondPack) {
			p.SecondSaveToFile()
			pack := teleport.SecondaryPackage{}
			p.SecondPack = pack
			p.SecondPackQty++
		}

		switch v.(type) {
		case storage.Product:
			p.Indexer.Set(teleport.UUID(v.(storage.Product).ID).String())
			p.PrimaryPack.AddItem(teleport.Post{
				ID:       teleport.UUID(v.(storage.Product).ID),
				AuthorID: 1,
				Date:     time.Now(),
				Content:  v.(storage.Product).Description,
				Title:    v.(storage.Product).Name,
				Excerpt:  "",
				Name:     v.(storage.Product).Name,
				Modified: time.Now(),
			})
			p.PrimaryPack.AddItem(teleport.TeleportItem{
				GUID: teleport.UUID(v.(storage.Product).ID),
				Type: "post",
				Date: time.Now(),
			})
		case storage.Group:
			p.Indexer.Set(teleport.UUID(v.(storage.Group).ID).String())
			p.PrimaryPack.AddItem(teleport.Term{
				ID:    teleport.UUID(v.(storage.Group).ID),
				Name:  v.(storage.Group).Name,
				Slug:  teleport.Slug(v.(storage.Group).Name),
				Group: "0",
			})
			p.PrimaryPack.AddItem(teleport.TeleportItem{
				GUID: teleport.UUID(v.(storage.Group).ID),
				Type: "term",
				Date: time.Now(),
			})
			p.SecondPack.AddItem(teleport.TermTaxonomy{
				TermID:      teleport.UUID(v.(storage.Group).ID),
				Taxonomy:    "product_cat",
				Description: v.(storage.Group).Name,
				Parent:      teleport.UUID(v.(storage.Group).ParentID),
			})
		}
	}
}

func (p *pkg) IsFull(pack teleport.PrimaryPackage) bool {
	return pack.Length >= p.Options.MaxBytes+p.Indexer.GetLength()+2000
}

func (p *pkg) SecondIsFull(pack teleport.SecondaryPackage) bool {
	return pack.Length >= p.Options.MaxBytes+2000
}

func (p *pkg) PreContent(s string) {
	p.Content = s + p.Content
}

func (p *pkg) AddContent(s string) {
	p.Content = p.Content + s
}

func (p *pkg) ClearContent() {
	p.Content = ""
}

func (p *pkg) SaveToFile() error {
	w := writer.NewWriter(p.Options.PrefixFileName, p.Options.PathToSave)
	fileName := w.GetFileName(p.PackQty)
	fmt.Println(fileName)
	wpwc := teleport.Wpwc{
		Prefix: p.Options.PrefixTableName,
	}

	if p.PackQty == 1 {
		p.AddContent(fmt.Sprintf("create table if not exists %steleport_item(guid char(32)not null,type char(8)not null,id bigint,date datetime,key id(`id`))engine=innodb default charset=utf8;", p.Options.PrefixTableName))
	}

	p.AddContent("start transaction;")
	p.AddContent(fmt.Sprintf("set @max_term_id=(select ifnull(max(term_id),0)from %sterms);", p.Options.PrefixTableName))
	p.AddContent(fmt.Sprintf("set @max_term_taxonomy_id=(select ifnull(max(term_taxonomy_id),0)from %sterm_taxonomy);", p.Options.PrefixTableName))
	p.AddContent(fmt.Sprintf("set @max_post_id=(select ifnull(max(id),0)from %sposts);", p.Options.PrefixTableName))
	p.AddContent(fmt.Sprintf("set @author_id=%d;", 1)) //todo author

	if len(p.Indexer.GetAll()) > 0 {
		for k, v := range p.Indexer.GetAll() {
			p.AddContent(fmt.Sprintf("set @%s=%d;", k, v))
		}
	}

	if len(p.PrimaryPack.Term) > 0 {
		builder := wpwc.BuilderTerm()
		for _, v := range p.PrimaryPack.Term {
			builder.AddTerm(v)
		}
		p.AddContent(fmt.Sprintf("%s;", squirrel.DebugSqlizer(builder)))
	}

	if len(p.PrimaryPack.Post) > 0 {
		builder := wpwc.BuilderPost()
		for _, v := range p.PrimaryPack.Post {
			builder.AddPost(v)
		}
		p.AddContent(fmt.Sprintf("%s;", squirrel.DebugSqlizer(builder)))
	}

	if len(p.Indexer.GetAll()) > 0 {
		builder := wpwc.BuilderTeleportItem()
		for _, v := range p.PrimaryPack.TeleportItem {
			builder.AddTeleportItem(v)
		}
		p.AddContent(fmt.Sprintf("%s;", squirrel.DebugSqlizer(builder)))
	}

	p.AddContent("commit;")
	fmt.Printf("%s\n", p.Content)
	//fmt.Println(p.PrimaryPack.Length)
	//fmt.Println(p.Indexer.GetLength())

	//w.WriteFile(fileName, content)
	return nil
}

func (p *pkg) SecondSaveToFile() error {
	p.ClearContent()
	w := writer.NewWriter(fmt.Sprintf("sec/%s", p.Options.PrefixFileName), p.Options.PathToSave)
	fileName := w.GetFileName(p.SecondPackQty)
	fmt.Println(fileName)
	wpwc := teleport.Wpwc{
		Prefix: p.Options.PrefixTableName,
	}

	idx := indexer.NewIndexer()

	if len(p.SecondPack.TermTaxonomy) > 0 {
		builder := wpwc.BuilderTermTaxonomy()
		for _, v := range p.SecondPack.TermTaxonomy {
			idx.Set(v.TermID.String())
			idx.Set(v.Parent.String())
			builder.AddTermTaxonomy(v)
		}
		p.AddContent(squirrel.DebugSqlizer(builder))
	}

	if len(idx.GetAll()) > 0 {
		for k, _ := range idx.GetAll() {
			if k != "" {
				p.PreContent(fmt.Sprintf("set @%s=(select id from %steleport_item where guid='%s');", k, wpwc.Prefix, k))
			}
		}
	}

	fmt.Println(p.Content)

	return nil
}

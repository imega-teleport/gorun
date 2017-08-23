package packer

import (
	"fmt"

	"time"

	"github.com/imega-teleport/gorun/indexer"
	"github.com/imega-teleport/gorun/storage"
	"github.com/imega-teleport/gorun/teleport"
	"github.com/imega-teleport/gorun/writer"
	squirrel "gopkg.in/Masterminds/squirrel.v1"
)

// Packer is interface
type Packer interface {
	Listen(in <-chan interface{}, e chan<- error)
	SaveToFile() error
}

type Options struct {
	MaxBytes        int
	PrefixFileName  string
	PathToSave      string
	PrefixTableName string
}

type pkg struct {
	Options Options
	Pack    teleport.Package
	Indexer indexer.Indexer
	PackQty int
	Content string
}

// New instance packer
func New(opt Options) Packer {
	return &pkg{
		Options: opt,
		Indexer: indexer.NewIndexer(),
		PackQty: 1,
	}
}

func (p *pkg) Listen(in <-chan interface{}, e chan<- error) {
	for v := range in {
		if p.IsFull(p.Pack) {
			p.SaveToFile()
			pack := teleport.Package{}
			p.Content = ""
			p.Pack = pack
			p.PackQty++
		}

		switch v.(type) {
		case storage.Product:
			p.Indexer.Set(teleport.UUID(v.(storage.Product).ID).String())
			p.Pack.AddItem(teleport.Post{
				ID:       teleport.UUID(v.(storage.Product).ID),
				AuthorID: 1,
				Date:     time.Now(),
				Content:  v.(storage.Product).Description,
				Title:    v.(storage.Product).Name,
				Excerpt:  "",
				Name:     v.(storage.Product).Name,
				Modified: time.Now(),
			})
			p.Pack.AddItem(teleport.TeleportItem{
				GUID: teleport.UUID(v.(storage.Product).ID),
				Type: "post",
				Date: time.Now(),
			})
		case storage.Group:
			p.Indexer.Set(teleport.UUID(v.(storage.Group).ID).String())
			p.Pack.AddItem(teleport.Term{
				ID:    teleport.UUID(v.(storage.Group).ID),
				Name:  v.(storage.Group).Name,
				Slug:  teleport.Slug(v.(storage.Group).Name),
				Group: "0",
			})
			p.Pack.AddItem(teleport.TeleportItem{
				GUID: teleport.UUID(v.(storage.Group).ID),
				Type: "term",
				Date: time.Now(),
			})
		}
	}
}

func (p *pkg) IsFull(pack teleport.Package) bool {
	return pack.Length >= p.Options.MaxBytes+p.Indexer.GetLength()+2000
}

func (p *pkg) AddContent(s string) {
	p.Content = p.Content + s
}

func (p *pkg) SaveToFile() error {
	w := writer.NewWriter(p.Options.PrefixFileName, p.Options.PathToSave)
	fileName := w.GetFileName(p.PackQty)
	fmt.Println(fileName)
	wpwc := teleport.Wpwc{
		Prefix: p.Options.PrefixTableName,
	}

	if p.PackQty == 1 {
		p.AddContent("create table if not exists teleport_item(guid CHAR(32) NOT NULL,type CHAR(8) NOT NULL, id bigint, date datetime, KEY id (`id`)) ENGINE=InnoDB DEFAULT CHARSET=utf8;")
	}

	p.AddContent("start transaction;")
	p.AddContent(fmt.Sprintf("set @max_term_id=(select max(term_id) from %sterms);", p.Options.PrefixTableName))
	p.AddContent(fmt.Sprintf("set @max_term_taxonomy_id=(select max(term_taxonomy_id) from %sterm_taxonomy);", p.Options.PrefixTableName))
	p.AddContent(fmt.Sprintf("set @max_post_id=(select max(id) from %sposts);", p.Options.PrefixTableName))
	p.AddContent(fmt.Sprintf("set @author_id=%d;", 1)) //todo author


	if len(p.Indexer.GetAll()) > 0 {
		for k, v := range p.Indexer.GetAll() {
			p.AddContent(fmt.Sprintf("set @%s=%d;", k, v))
		}
	}

	if len(p.Pack.Term) > 0 {
		builder := wpwc.BuilderTerm()
		for _, v := range p.Pack.Term {
			builder.AddTerm(v)
		}
		p.AddContent(fmt.Sprintf("%s;", squirrel.DebugSqlizer(builder)))
	}

	if len(p.Pack.Post) > 0 {
		builder := wpwc.BuilderPost()
		for _, v := range p.Pack.Post {
			builder.AddPost(v)
		}
		p.AddContent(fmt.Sprintf("%s;", squirrel.DebugSqlizer(builder)))
	}

	if len(p.Indexer.GetAll()) > 0 {
		builder := wpwc.BuilderTeleportItem()
		for _, v := range p.Pack.TeleportItem {
			builder.AddTeleportItem(v)
		}
		p.AddContent(fmt.Sprintf("%s;", squirrel.DebugSqlizer(builder)))
	}

	p.AddContent("commit;")
	fmt.Printf("%s\n", p.Content)
	//fmt.Println(p.Pack.Length)
	//fmt.Println(p.Indexer.GetLength())

	//w.WriteFile(fileName, content)
	return nil
}

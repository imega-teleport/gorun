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
			p.Pack = pack
			p.PackQty++
		}

		switch v.(type) {
		case storage.Product:
			//fmt.Println("Product: ", v.(storage.Product).Name)
			p.Indexer.Set(v.(storage.Product).ID)
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
			p.Indexer.Set(v.(storage.Group).ID)
			p.Pack.AddItem(teleport.Term{
				ID:    teleport.UUID(v.(storage.Group).ID),
				Name:  v.(storage.Group).Name,
				Slug:  v.(storage.Group).Name,
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
	return pack.Length >= p.Options.MaxBytes+500
}

func (p *pkg) SaveToFile() error {
	w := writer.NewWriter(p.Options.PrefixFileName, p.Options.PathToSave)
	fileName := w.GetFileName(p.PackQty)
	fmt.Println(fileName)
	wpwc := teleport.Wpwc{
		Prefix: p.Options.PrefixTableName,
	}

	if len(p.Pack.Term) > 0 {
		builder := wpwc.BuilderTerm()
		for _, v := range p.Pack.Term {
			builder.AddTerm(v)
		}
		fmt.Println(squirrel.DebugSqlizer(builder))
	}


	if len(p.Pack.Post) > 0 {
		builder1 := wpwc.BuilderPost()
		for _, v := range p.Pack.Post {
			builder1.AddPost(v)
		}
		fmt.Println(squirrel.DebugSqlizer(builder1))
	}
	//w.WriteFile(fileName, content)
	return nil
}

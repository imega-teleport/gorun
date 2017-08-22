package packer

import (
	"fmt"

	"github.com/imega-teleport/gorun/storage"
	"github.com/imega-teleport/gorun/teleport"
)

// Packer is interface
type Packer interface {
	Listen(in <-chan interface{}, e chan<- error)
}

type pkg struct {
	MaxBytes int
	Pack     teleport.Package
}

// New instance packer
func New(maxBytes int) Packer {
	return &pkg{
		MaxBytes: maxBytes,
	}
}

func (p *pkg) Listen(in <-chan interface{}, e chan<- error) {
	for v := range in {
		if p.IsFull(p.Pack) {
			pack := teleport.Package{}
			p.Pack = pack
		}

		switch v.(type) {
		case storage.Product:
			fmt.Println("Product: ", v.(storage.Product).Name)
		case storage.Group:
			p.Pack.AddItem(teleport.Term{
				ID:    v.(storage.Group).ID,
				Name:  v.(storage.Group).Name,
				Slug:  v.(storage.Group).Name,
				Group: "0",
			})
		}
	}
}

func (p *pkg) IsFull(pack teleport.Package) bool {
	return pack.Length >= p.MaxBytes+500
}

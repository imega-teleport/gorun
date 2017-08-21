package packer

import "github.com/imega-teleport/gorun/teleport"

// Packer is interface
type Packer interface {
	Listen(in <-chan interface{}, out <-chan error)
}

type pkg struct {
	MaxBytes int
	Pack     teleport.Package
}

// New instance packer
func New(maxBytes int) 	Packer {
	return &pkg{
		MaxBytes: maxBytes,
	}
}

func (p *pkg) Listen(in <-chan interface{}, out <-chan error) {
	/*for v := range in {

		teleport.
			fmt.Println(groups.length)
	}*/
}

func (p *pkg) IsFull(pack teleport.Package) bool {
	return pack.Length >= p.MaxBytes+500
}

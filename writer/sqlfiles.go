package writer

import (
	"fmt"
	"os"
)

// Writer is interface
type Writer interface {
	Listen(in <-chan string)
}

type writerFiles struct {
	path   string
	count  int
	prefix string
}

// NewWriter get new instance
func NewWriter(PrefixFileName string, path string) Writer {
	return &writerFiles{
		path:   path,
		prefix: PrefixFileName,
	}
}

func (w *writerFiles) Listen(in <-chan string, errOut chan<- error) {
	for v := range in {
		w.count++
		fileName := fmt.Sprintf("%s%c%s_%d.sql", w.path, os.PathSeparator, w.prefix, w.count)
		errOut <- writeFile(fileName, v)
	}
}

func writeFile(fileName, content string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			return err
		}
	}()
	_, err := file.WriteString(in)
	if err != nil {
		return err
	}
}

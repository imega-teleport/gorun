package writer

import (
	"fmt"
	"os"
)

// Writer is interface
type Writer interface {
	Listen(in <-chan string, errOut chan<- error)
	WriteFile(fileName, content string) (err error)
	GetFileName(count int) string
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
		errOut <- w.WriteFile(fileName, v)
	}
}

func (w *writerFiles) WriteFile(fileName, content string) (err error) {
	file, err := os.Create(fileName)
	if err != nil {
		return
	}
	defer func() {
		err = file.Close()
	}()
	_, err = file.WriteString(content)
	return
}

func (w *writerFiles) GetFileName(count int) string {
	return fmt.Sprintf("%s%c%s_%d.sql", w.path, os.PathSeparator, w.prefix, count)
}

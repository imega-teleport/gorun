package writer

import (
	"os"
	"fmt"
	"google.golang.org/appengine/file"
)

type Writer interface {
	Listen(in <-chan string)
}

type writerFiles struct {
	path   string
	count  int
	prefix string
	File   os.File
}

func NewWriter(PrefixFileName string, path string) Writer {
	return &writerFiles{
		path:   path,
		prefix: PrefixFileName,
	}
}

func (w *writerFiles) Listen(in <-chan string) {
	for v := range in {
		file, _ := os.Create(fmt.Sprintf("%s%c%s_%d.sql", w.path, os.PathSeparator, w.prefix, w.count))
		/*if err != nil {
			fmt.Printf("Could not create file: %s", err)
			os.Exit(1)
		}*/
		defer func() {
			file.Close()
			/*if err := file.Close(); err != nil {
				fmt.Printf("Fail close file: %s", err)
				os.Exit(1)
			}*/
		}()
	}
}

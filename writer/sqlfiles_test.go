package writer

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSqlFiles_writeFile(t *testing.T) {
	fileName := "/tmp/testsqlfiles_writefile"
	expected := "content"
	w := writerFiles{}
	err := w.WriteFile(fileName, expected)
	assert.NoError(t, err)

	data, err := ioutil.ReadFile(fileName)
	assert.NoError(t, err)

	assert.Equal(t, expected, string(data))

	err = os.Remove(fileName)
	assert.NoError(t, err)
}

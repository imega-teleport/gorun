package writer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSqlFiles_writeFile(t *testing.T) {

	writeFile("/tmp/testsqlfiles_writefile", "content")
	assert.NoError(t)
}

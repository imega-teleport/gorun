package indexer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIdx_Set_OneItem_ReturnsOne(t *testing.T) {
	i := NewIndexer()

	i.Set("test1")

	assert.Equal(t, 1, i.Get("test1"))
}

func TestIdx_Set_TwoItem_ReturnsTwo(t *testing.T) {
	i := NewIndexer()

	i.Set("test1")
	i.Set("test2")

	assert.Equal(t, 2, i.Get("test2"))
}

func TestIdx_Set_ExistsItem_ReturnsOne(t *testing.T) {
	i := NewIndexer()

	i.Set("test1")
	i.Set("test1")

	assert.Equal(t, 1, i.Get("test1"))
}

func TestIdx_Get_NotExistsItem_ReturnsZero(t *testing.T) {
	i := NewIndexer()

	assert.Equal(t, 0, i.Get("test1"))
}

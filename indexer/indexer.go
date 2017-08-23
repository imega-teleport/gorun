package indexer

type Indexer interface {
	Get(uid string) int
	Set(uid string)
	GetAll() map[string]int
	GetLength() int
}

type idx struct {
	values map[string]int
	length int
}

func NewIndexer() Indexer {
	return &idx{
		values: make(map[string]int),
	}
}

func (i *idx) Get(uid string) int {
	return i.values[uid]
}

func (i *idx) Set(uid string) {
	if i.values[uid] == 0 {
		i.length = i.length + len(uid) + 12
		i.values[uid] = len(i.values) + 1
	}
}

func (i *idx) GetAll() map[string]int {
	return i.values
}

func (i *idx) GetLength() int {
	return i.length
}

package indexer

type Indexer interface {
	Get(uid string) int
	Set(uid string)
	GetAll() map[string]int
}

type idx struct {
	values map[string]int
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
		i.values[uid] = len(i.values) + 1
	}
}

func (i *idx) GetAll() map[string]int {
	return i.values
}

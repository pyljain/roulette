package lru

type LRUCache interface {
	Get(fileName string) ([]byte, error)
	Add(filename string, contents []byte) error
}

type CacheItem struct {
	Filename string
	Contents []byte
	SizeInKB float64
}

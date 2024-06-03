package cache

type Cache interface {
	Has(string) (bool, error)
	Get(string) (interface{}, error)
	Set(string, interface{}, ...int) error
	Delete(string) error
	EmptyByMatch(string) error
	Prune() error
}

// Entry is a map of string to interface
// This is what will be stored in the cache after serializing it.
type Entry map[string]interface{}

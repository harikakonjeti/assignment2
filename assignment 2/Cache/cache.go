package Cache

type Data struct {
	Key string
	Ele interface{}
}

type Cache interface {
	Put(d Data) error
	Get(key string) (Data, error)
	Delete(key string) error
	Purge() error
	GetList() ([]Data, error)
	Type() string
}
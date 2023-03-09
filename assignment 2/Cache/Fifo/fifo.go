package Fifo

import (
	"lru/Cache"
	"lru/Cache/Redis"
	"container/list"
	"errors"
)

type Fifo struct {
	que         list.List
	mp          map[string]*list.Element
	size        int
	policy      string
	nextservice Cache.Cache
}

func CreateFifoCache(csize int) (Cache.Cache, error) {
	if csize <= 0 {
		return nil, errors.New("invalid cache size")
	}
	return &Fifo{
		que:         list.List{},
		mp:          make(map[string]*list.Element),
		size:        csize,
		policy:      "FIFO",
		nextservice: Redis.NewRedis(),
	}, nil
}

func (f *Fifo) Put(d Cache.Data) (err error) {
	_, ok := f.mp[d.Key]
	if !ok {

		if f.que.Len() == f.size {

			last := f.que.Back()
			d := f.que.Remove(last)
			delete(f.mp, d.(Cache.Data).Key)

		}
	} else {
		f.que.Remove(f.mp[d.Key])
	}
	f.mp[d.Key] = f.que.PushFront(d)
	return nil
}

func (f *Fifo) Type() string {
	return f.policy
}

func (f *Fifo) Get(key string) (data Cache.Data, err error) {

	_, ok := f.mp[key]
	if !ok {
		d, err := f.nextservice.Get(key)
		if err == nil {
			f.Put(d)
		}
		return d, err
	}
	return f.mp[key].Value.(Cache.Data), nil

}

func (f *Fifo) Delete(key string) error {
	_, ok := f.mp[key]
	if !ok {
		return errors.New("key is invalid")
	}
	f.que.Remove(f.mp[key])
	delete(f.mp, key)
	return nil
}

func (f *Fifo) Purge() error {
	if f.que.Len() == 0 {
		return errors.New("cache is already empty")
	}
	for e := f.que.Front(); e != nil; e = e.Next() {
		f.que.Remove(e)
	}
	return nil
}

func (f *Fifo) GetList() (dataList []Cache.Data, err error) {

	var dt []Cache.Data
	if f.que.Len() == 0 {
		return dt, errors.New("the list is empty")
	}

	for e := f.que.Front(); e != nil; e = e.Next() {
		dt = append(dt, e.Value.(Cache.Data))
	}
	return dt, nil
}
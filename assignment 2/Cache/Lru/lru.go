package Lru

import (
	"lru/Cache"
	"lru/Cache/Redis"
	"container/list"
	"errors"
)

type Lru struct {
	que         list.List
	mp          map[string]*list.Element
	size        int
	policy      string
	nextservice Cache.Cache
}

func CreateLruCache(csize int) (Cache.Cache,error) {
	if(csize <=0){
		return nil,errors.New("invalid cache size")
	}
	return &Lru{
		que:         list.List{},
		mp:          make(map[string]*list.Element),
		size:        csize,
		policy:      "LRU",
		nextservice: Redis.NewRedis(),
	},nil
}

func (l*Lru) Type() string{
	return l.policy
}

func (l *Lru) Put(d Cache.Data) (err error) {
	_, ok := l.mp[d.Key]
	if !ok {
		if l.que.Len() == l.size {
			last := l.que.Back()
			d := l.que.Remove(last)
			delete(l.mp, d.(Cache.Data).Key)

		}
	} else {
		l.que.Remove(l.mp[d.Key])
	}
	l.mp[d.Key] = l.que.PushFront(d)
	return nil
}

func (l *Lru) Get(key string) (data Cache.Data, err error) {

	_, ok := l.mp[key]
	if !ok {
		d, err := l.nextservice.Get(key)
		if err == nil {
			l.Put(d)
		}
		return d, err
	}
	l.que.MoveToFront(l.mp[key])
	l.mp[key] = l.que.Front()

	return l.que.Front().Value.(Cache.Data), nil
}

func (l *Lru) Delete(key string) error {
	_, ok := l.mp[key]
	if !ok {
		return errors.New("key is invalid")
	}
	l.que.Remove(l.mp[key])
	delete(l.mp, key)
	return nil
}

func (l *Lru) Purge() error {
	if l.que.Len() == 0 {
		return errors.New("cache is already empty")
	}
	for l.que.Len() != 0 {
		l.que.Remove(l.que.Front())
	}
	return nil
}

func (l *Lru) GetList() (dataList []Cache.Data, err error) {

	var dt []Cache.Data
	if l.que.Len() == 0 {
		return nil, errors.New("the list is empty")
	}

	for e := l.que.Front(); e != nil; e = e.Next() {
		dt = append(dt, e.Value.(Cache.Data))
	}
	return dt, nil
}
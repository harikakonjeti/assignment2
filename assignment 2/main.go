package main

import (
	"lru/Cache"
	"lru/Cache/Lru"
	"fmt"
	"log"
)

func main() {

	l, err := Lru.CreateLruCache(1)
	if err != nil {
		log.Fatal(err)
	}
	var data = []Cache.Data{
		{"1", "one"},
		{"2", "two"},
		{"3", "three"},
	}
	for _, d := range data {
		l.Put(d)
	}

	val, err := l.Get("2")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("The value for key : ", val)

}

// package main

// import "container/list"

// type lru struct {
// 	queue list.List
// 	size  int
// }
// type LRU interface{
//    lrucache()
// }

// func main() {
//    l:=&lru{}


// 	l.set(3)
//    l.lrucache(1)
//    l.lrucache(2)
//    l.lrucache(3)
//    l.lrucache(2)
//    l.lrucache(5)
// }


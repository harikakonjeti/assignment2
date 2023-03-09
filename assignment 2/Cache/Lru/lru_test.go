package Lru

import (
	"lru/Cache"
	"lru/Cache/Redis"
	"reflect"
	"testing"
)

func TestCreateLruCache(t *testing.T) {
	type args struct {
		csize int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    string
	}{
		{
			name:    "Test1",
			args:    args{5},
			wantErr: false,
		},
		{
			name:    "Test2",
			args:    args{-1},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		_, err := CreateLruCache(tt.args.csize)
		if (err != nil) != tt.wantErr {
			t.Error("Error in Cache size handling")
		}
	}
}
func TestLru_Get(t *testing.T) {
	c, _ := CreateLruCache(3)
	l := c.(*Lru)
	l.Put(Cache.Data{"1", "one"})
	l.Put(Cache.Data{"2", "two"})
	r := Redis.NewRedis()
	r.Put(Cache.Data{"3", "three"})
	type args struct {
		key string
	}
	tests := []struct {
		name     string
	
		args     args
		wantData Cache.Data
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			name: "Test1",
			args: args{
				key: "2",
			},
			wantData: Cache.Data{"2", "two"},
			wantErr:  false,
		},
		{
			name: "Test2",
			args: args{
				key: "3",
			},
			wantData: Cache.Data{"3", "three"},
			wantErr:  false,
		},
		{
			name: "Test3",
			args: args{
				key: "abcdefg",
			},
			wantData: Cache.Data{},
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			
			gotData, err := l.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Lru.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotData, tt.wantData) {
				t.Errorf("Lru.Get() = %v, want %v", gotData, tt.wantData)
			}
			if err == nil && l.que.Front().Value != tt.wantData {
				t.Error("Queue is not updating")
			}
			if !reflect.DeepEqual(gotData, tt.wantData) {
				t.Errorf("Lru.Get() = %v, want %v", gotData, tt.wantData)
			}
		})
	}
}




func TestLru_Put(t *testing.T) {
	c, _ := CreateLruCache(3)
	l := c.(*Lru)
	type args struct {
		d Cache.Data
	}
	tests := []struct {
		name    string
		args    args
	}{
		{
			args: args{
				d: Cache.Data{"4", "four"},
			},
		},
		{
			args: args{
				d: Cache.Data{"3", "three"},
			},
		},
	}
	for i, tt := range tests {
		l.Put(tt.args.d)
		if l.mp[tt.args.d.Key] != l.que.Front() {
			t.Error("Map is not updating")
		}
		if l.que.Front().Value != tt.args.d {
			t.Error("Queue is not updating")
		}
		_, ok := l.mp[tt.args.d.Key]
		if !ok {
			if i >= l.size-1 && l.que.Back().Value != tests[i-l.size+1].args.d {
				t.Error("Error in Removing old Data")
			}
		}
	}
}

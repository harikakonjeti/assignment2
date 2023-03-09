package Fifo

import (
	"lru/Cache"
	"lru/Cache/Redis"
	"reflect"
	"testing"
)

func TestCreateFifoCache(t *testing.T) {
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
		_, err := CreateFifoCache(tt.args.csize)
		if (err != nil) != tt.wantErr {
			t.Error("Error in Cache size handling")
		}
	}
}
func TestFifo_Put(t *testing.T) {
	c, _ := CreateFifoCache(3)
	f := c.(*Fifo)
	type args struct {
		d Cache.Data
	}
	tests := []struct {
		name    string
	
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := f.Put(tt.args.d); (err != nil) != tt.wantErr {
				t.Errorf("Fifo.Put() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFifo_Get(t *testing.T) {
	c, _ := CreateFifoCache(3)
	f := c.(*Fifo)
	f.Put(Cache.Data{"1", "one"})
	f.Put(Cache.Data{"2", "two"})
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

		gotData, err := f.Get(tt.args.key)
		if (err != nil) != tt.wantErr {
			t.Errorf("Fifo.Get() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(gotData, tt.wantData) {
			t.Errorf("Fifo.Get() = %v, want %v", gotData, tt.wantData)
		}
	}
}

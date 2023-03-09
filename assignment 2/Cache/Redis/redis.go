package Redis

import (
	"lru/Cache"
	"context"
	"errors"
	"log"

	"github.com/go-redis/redis/v8"
)

type Rds struct {
	client *redis.Client
	policy string
}

var ctx = context.Background()

func NewRedis() Cache.Cache {
	return &Rds{
		client: redis.NewClient(&redis.Options{
			Addr:     "redis-15761.c44.us-east-1-2.ec2.cloud.redislabs.com:15761",
			Password: "ckDGRXV67OLZBBpbB1A2ASE4ElBYbV5f",
		}),
		policy: "RDS",
	}
}

func (r *Rds) Put(d Cache.Data) (err error) {
	return r.client.Set(ctx, d.Key, d.Ele, 0).Err()
}
func (r *Rds) Type() string {
	return r.policy
}
func (r *Rds) Get(key string) (data Cache.Data, err error) {

	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return Cache.Data{}, errors.New("key not found")
	} else if err != nil {
		log.Printf("Error in Redix %v", err)
		return Cache.Data{}, err
	}
	return Cache.Data{key, val}, err
}

func (r *Rds) Delete(key string) (err error) {

	err = r.client.Del(ctx, key).Err()
	if err == redis.Nil {
		return errors.New("key not found")
	} else if err != nil {
		log.Printf("Error in Redix %v", err)
		return err
	}
	return err
}

func (r *Rds) Purge() (err error) {

	err = r.client.FlushDB(ctx).Err()
	if err != nil {
		log.Printf("Error in Redix %v", err)
		return err
	}
	return err
}

func (r *Rds) GetList() (dataList []Cache.Data, err error) {

	keys, err := r.client.Keys(ctx, "*").Result()
	if err != nil {
		log.Println("Error in scanning Database in redix")
		return nil, err
	}

	for _, k := range keys {
		val, _ := r.client.Get(ctx, k).Result()
		dataList = append(dataList, Cache.Data{k, val})
	}
	return dataList, nil

}
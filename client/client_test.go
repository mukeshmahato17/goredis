package client

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func TestNewRedisClient(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:4000",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	fmt.Println(rdb)
	fmt.Println("this is working")

	err := rdb.Set(context.Background(), "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	// val, err := rdb.Get(context.Background(), "key").Result()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("key", val)
}

func TestNewClient1(t *testing.T) {
	client, err := New("localhost:4000")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	time.Sleep(time.Second)
	if err := client.Set((context.TODO()), "foo", 1); err != nil {
		log.Fatal(err)
	}
	val, err := client.Get((context.TODO()), "foo")
	if err != nil {
		log.Fatal(err)
	}

	n, _ := strconv.Atoi(val)
	fmt.Println(n)
	fmt.Println("GET =>", val)
}

func TestNewClient(t *testing.T) {
	client, err := New("localhost:4000")
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		fmt.Println("SET =>", fmt.Sprintf("bar_%d", i))
		if err := client.Set((context.TODO()), fmt.Sprintf("foo_%d", i), fmt.Sprintf("bar_%d", i)); err != nil {
			log.Fatal(err)
		}
		val, err := client.Get((context.TODO()), fmt.Sprintf("foo_%d", i))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("GET =>", val)
	}
}

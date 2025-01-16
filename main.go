package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mukeshmahato17/goredis/client"
)

func main() {
	cfg := Config{
		ListenAddr: ":4000",
	}
	server := NewServer(cfg)

	go func() {
		server.Start()
	}()
	time.Sleep(time.Second)

	for i := 0; i < 10; i++ {

		client := client.New("localhost:4000")
		if err := client.Set((context.TODO()), fmt.Sprintf("foo_%d", i), fmt.Sprintf("bar_%d", i)); err != nil {
			log.Fatal(err)
		}
		val, err := client.Get((context.TODO()), fmt.Sprintf("foo_%d", i))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(val)
	}

	time.Sleep(time.Second)
	fmt.Println(server.kv.data)
}

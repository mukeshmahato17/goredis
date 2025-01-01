package main

import (
	"context"
	"log"
	"time"

	"github.com/mukeshmahato17/goredis/client"
)

func main() {
	cfg := Config{
		ListenAddr: ":4000",
	}
	go func() {
		server := NewServer(cfg)
		server.Start()
	}()
	time.Sleep(time.Second)

	for i := 0; i < 10; i++ {

		client := client.New("localhost:4000")
		if err := client.Set(context.Background(), "foo", "bar"); err != nil {
			log.Fatal(err)
		}
	}

	time.Sleep(time.Second)
}

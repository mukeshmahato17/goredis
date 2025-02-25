package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/mukeshmahato17/goredis/client"
)

func TestFooBar(t *testing.T) {
	in := map[string]string{
		"first":  "1",
		"second": "2",
	}
	out := respWriteMap(in)
	fmt.Println(out)
}

func TestServerWithNewClients(t *testing.T) {

	server := NewServer(Config{})
	go func() {
		log.Fatal(server.Start())
	}()
	time.Sleep(time.Second)

	nClients := 10
	wg := sync.WaitGroup{}
	wg.Add(nClients)

	for i := 0; i < nClients; i++ {
		go func(it int) {
			c, err := client.New("localhost:4000")
			if err != nil {
				log.Fatal(err)
			}
			defer c.Close()

			key := fmt.Sprintf("client_foo_%d", i)
			value := fmt.Sprintf("client_bar_%d", i)
			if err := c.Set((context.TODO()), key, value); err != nil {
				log.Fatal(err)
			}
			val, err := c.Get((context.TODO()), key)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("client %d got this value back => %s\n", it, val)
			wg.Done()
		}(i)
	}

	wg.Wait()

	time.Sleep(time.Second)
	if len(server.peers) != 0 {
		t.Fatalf("expected 0 but got %d", len(server.peers))
	}
}

package client

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client, err := New("localhost:4000")
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		fmt.Println("set =>", fmt.Sprintf("bar_%d", i))
		if err := client.Set((context.TODO()), fmt.Sprintf("foo_%d", i), fmt.Sprintf("bar_%d", i)); err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Second)
		val, err := client.Get((context.TODO()), fmt.Sprintf("foo_%d", i))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("get", val)
	}
}

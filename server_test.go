package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/adityapandey23/redis-clone/client"
)

func TestServerWithMultiClients(t *testing.T) {

	server := NewServer(Config{})
	go func() {
		log.Fatal(server.Start())
	}()

	time.Sleep(time.Second) // For server to boot up

	nClients := 10
	wg := sync.WaitGroup{}
	wg.Add(nClients)

	for i := range nClients {
		go func(it int) {
			c, err := client.New("localhost:5000")
			if err != nil {
				log.Fatal("client error", "error", err)
			}

			defer c.Close()

			key := fmt.Sprintf("foo_%d", i)
			value := fmt.Sprintf("bar_%d", i)

			if err := c.Set(context.TODO(), key, value); err != nil {
				log.Fatal(err)
			}

			val, err := c.Get(context.TODO(), key)

			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("client %d got this %s back\n", it, val)
			wg.Done()
		}(i)
	}

	wg.Wait()

	time.Sleep(time.Second)

	if len(server.peers) != 0 {
		t.Fatalf("expected 0 peers but got %d", len(server.peers))
	}

}

func TestYadaYada(t *testing.T) {
	in := map[string]string{
		"first":  "1",
		"second": "2",
	}

	out := respWriteMap(in)

	fmt.Println(out)

}

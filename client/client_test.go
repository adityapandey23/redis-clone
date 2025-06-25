package client

import (
	"context"
	"fmt"
	"log"
	"sync"
	"testing"
)

func TestNewClients(t *testing.T) {
	nClients := 10
	wg := sync.WaitGroup{}
	wg.Add(nClients)

	for i := range nClients {
		go func(it int) {
			c, err := New("localhost:5002")
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

}

func TestNewClient(t *testing.T) {
	for i := range 10 {
		c, err := New("localhost:5002")
		if err != nil {
			log.Fatal("client error", "error", err)
		}
		if err := c.Set(context.TODO(), fmt.Sprintf("foo_%d", i), fmt.Sprintf("bar_%d", i)); err != nil {
			log.Fatal(err)
		}

		val, err := c.Get(context.TODO(), fmt.Sprintf("foo_%d", i))

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("got this back", val)

	}

}

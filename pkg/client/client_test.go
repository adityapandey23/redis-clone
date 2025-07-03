package client

import (
	"context"
	"fmt"
	"log"
	"testing"
)

func TestNewClient1(t *testing.T) {
	for i := range 10 {
		c, err := New("localhost:5000")
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

func TestNewClient2(t *testing.T) {
	c, err := New("localhost:5000")
	if err != nil {
		log.Fatal("client error", "error", err)
	}

	defer c.Close()

	key := "foo"
	value := 2

	if err := c.Set(context.TODO(), key, value); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("the type of key %T, the type of value %T\n", key, value)

	fmt.Printf("setting key as %v and val as %v\n", key, value)

	val, err := c.Get(context.TODO(), key)

	fmt.Printf("getting %v for the key %v\n", val, key)

	if err != nil {
		log.Fatal(err)
	}

}

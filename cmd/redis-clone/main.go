package main

import (
	"flag"
	"log"
)
import "github.com/adityapandey23/redis-clone/pkg/server"

func main() {
	listenAddr := flag.String("listenAddr", server.DefaultListenAddr, "listening address of the server")

	flag.Parse()
	s := server.NewServer(server.Config{
		ListenAddr: *listenAddr,
	})
	log.Fatal(s.Start())
}

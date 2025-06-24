package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"time"

	"github.com/adityapandey23/redis-clone/client"
)

const defaultListenAddr = ":5001"

type Config struct {
	ListenAddr string
}

type Message struct {
	data []byte
	peer *Peer
}

type Server struct {
	Config    Config
	peers     map[*Peer]bool
	ln        net.Listener
	addPeerCh chan *Peer
	quitChn   chan struct{}
	msgChn    chan Message
	kv        *KV
}

func NewServer(cfg Config) *Server {
	if len(cfg.ListenAddr) == 0 {
		cfg.ListenAddr = defaultListenAddr
	}

	return &Server{
		Config:    cfg,
		peers:     make(map[*Peer]bool),
		addPeerCh: make(chan *Peer),
		quitChn:   make(chan struct{}),
		msgChn:    make(chan Message),
		kv:        NewKV(),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.Config.ListenAddr)

	if err != nil {
		return err
	}

	s.ln = ln

	go s.loop()

	slog.Info("server running", "listenAddr", s.Config.ListenAddr)

	return s.acceptLoop()
}

func (s *Server) handleMessage(msg Message) error {
	cmd, err := parseCommand(string(msg.data))
	if err != nil {
		return err
	}

	switch v := cmd.(type) {
	case SetCommand:
		return s.kv.Set(v.key, v.val)

	case GetCommand:
		val, ok := s.kv.Get(v.key)
		if !ok {
			return fmt.Errorf("key not found")
		}

		_, err := msg.peer.Send(val)
		if err != nil {
			slog.Error("peer send error", "err", err)
		}

	}

	return nil
}

func (s *Server) loop() {
	for {
		select {
		case msg := <-s.msgChn:
			if err := s.handleMessage(msg); err != nil {
				slog.Error("raw message error", "error", err.Error())
			}
		case <-s.quitChn:
			return
		case peer := <-s.addPeerCh:
			s.peers[peer] = true
		}
	}
}

func (s *Server) acceptLoop() error {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			slog.Error("accept error", "err", err)
			continue
		}
		// slog.Info("peer connection initiated")
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	peer := NewPeer(conn, s.msgChn)
	s.addPeerCh <- peer
	// slog.Info("peer connected", "remoteAddr", conn.RemoteAddr())
	go peer.readLoop()
}

func main() {
	server := NewServer(Config{})
	go func() {
		log.Fatal(server.Start())
	}()

	time.Sleep(time.Second)

	fmt.Println("This is before", server.kv.data)

	for i := range 10 {
		c := client.New("localhost:5001") // this address is server's address
		if err := c.Set(context.TODO(), fmt.Sprintf("foo_%d", i), fmt.Sprintf("bar_%d", i)); err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Second)
		// Alert: Here we can face the issue of, Getting a thing before setting
		val, err := c.Get(context.TODO(), fmt.Sprintf("foo_%d", i))

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("got this back", val)

	}

	time.Sleep(time.Second)

	fmt.Println("This is after", server.kv.data)
}

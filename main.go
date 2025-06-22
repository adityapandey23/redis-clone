package main

import (
	"context"
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

type Server struct {
	Config    Config
	peers     map[*Peer]bool
	ln        net.Listener
	addPeerCh chan *Peer
	quitChn   chan struct{}
	msgChn    chan []byte
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
		msgChn:    make(chan []byte),
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

func (s *Server) handleRawMessage(rawMsg []byte) error {
	cmd, err := parseCommand(string(rawMsg))
	if err != nil {
		return err
	}

	switch v := cmd.(type) {
	case SetCommand:
		slog.Info("set command received", "key", v.key, "val", v.val)
	}

	return nil
}

func (s *Server) loop() {
	for {
		select {
		case rawMsg := <-s.msgChn:
			if err := s.handleRawMessage(rawMsg); err != nil {
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
		slog.Info("peer connection initiated")
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	peer := NewPeer(conn, s.msgChn)
	s.addPeerCh <- peer
	slog.Info("peer connected", "remoteAddr", conn.RemoteAddr())
	go peer.readLoop()
}

func main() {
	go func() {
		server := NewServer(Config{})
		log.Fatal(server.Start())
	}()

	time.Sleep(time.Second)

	client := client.New("localhost:5001")

	if err := client.Set(context.TODO(), "foo", "bar"); err != nil {
		log.Fatal(err)
	}

	select {} // we are blocking here
}

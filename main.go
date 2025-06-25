package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net"
)

const defaultListenAddr = ":5001"

type Config struct {
	ListenAddr string
}

type Message struct {
	cmd  Command
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

	switch v := msg.cmd.(type) {
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
			slog.Info("new peer connected", "remoteAddr", peer.conn.RemoteAddr())
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
	listenAddr := flag.String("listenAddr", defaultListenAddr, "listening address of the server")
	flag.Parse()
	server := NewServer(Config{
		ListenAddr: *listenAddr, // Derefrencing
	})
	log.Fatal(server.Start())
}

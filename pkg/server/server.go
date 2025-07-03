package server

import (
	"fmt"
	"github.com/adityapandey23/redis-clone/internal"
	"log/slog"
	"net"
)

const DefaultListenAddr = ":5000"

type Config struct {
	ListenAddr string
}

type Server struct {
	Config    Config
	peers     map[*internal.Peer]bool
	ln        net.Listener
	addPeerCh chan *internal.Peer
	delPeerCh chan *internal.Peer
	msgChn    chan internal.Message
	kv        *internal.KV
}

func NewServer(cfg Config) *Server {
	if len(cfg.ListenAddr) == 0 {
		cfg.ListenAddr = DefaultListenAddr
	}

	return &Server{
		Config:    cfg,
		peers:     make(map[*internal.Peer]bool),
		addPeerCh: make(chan *internal.Peer),
		delPeerCh: make(chan *internal.Peer),
		msgChn:    make(chan internal.Message),
		kv:        internal.NewKV(),
	}
}

func (s *Server) handleMessage(msg internal.Message) error {
	switch v := msg.Cmd.(type) {
	case internal.SetCommand:
		return s.kv.Set(v.Key, v.Val)

	case internal.GetCommand:
		val, ok := s.kv.Get(v.Key)
		if !ok {
			return fmt.Errorf("key not found")
		}

		_, err := msg.Peer.Send(val)
		if err != nil {
			return fmt.Errorf("error: %s", err.Error())
		}
	}
	return nil
}

func (s *Server) loop() {
	for {
		select {
		case msg := <-s.msgChn:
			if err := s.handleMessage(msg); err != nil {
				slog.Error("message error", "error", err.Error())
			}
		case peer := <-s.delPeerCh:
			slog.Info("peer disconnected", "remoteAddr", peer.Conn.RemoteAddr())
			delete(s.peers, peer)
		case peer := <-s.addPeerCh:
			slog.Info("peer connected", "remoteAddr", peer.Conn.RemoteAddr())
			s.peers[peer] = true
		}
	}
}

func (s *Server) handleConn(conn net.Conn) {
	peer := internal.NewPeer(conn, s.msgChn, s.delPeerCh)
	s.addPeerCh <- peer
	slog.Info("peer connected", "remoteAddr", conn.RemoteAddr())
	go func() {
		err := peer.ReadLoop()
		if err != nil {
			slog.Error("unable to read", "err", err)
		}
	}()
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

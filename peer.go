package main

import "net"

type Peer struct {
	conn   net.Conn
	msgChn chan []byte
}

func NewPeer(conn net.Conn, msgChn chan []byte) *Peer {
	return &Peer{
		conn:   conn,
		msgChn: msgChn,
	}
}

func (p *Peer) readLoop() error {
	buf := make([]byte, 1024)

	for {
		n, err := p.conn.Read(buf)
		if err != nil {
			return err
		}
		msgBuf := make([]byte, n)
		copy(msgBuf, buf[:n])
		p.msgChn <- msgBuf
	}
}

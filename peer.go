package main

import (
	"fmt"
	"io"
	"net"

	"github.com/tidwall/resp"
)

type Peer struct {
	conn   net.Conn
	msgChn chan Message
}

func NewPeer(conn net.Conn, msgChn chan Message) *Peer {
	return &Peer{
		conn:   conn,
		msgChn: msgChn,
	}
}

func (p *Peer) Send(msg []byte) (int, error) {
	return p.conn.Write(msg)
}

func (p *Peer) readLoop() error {
	rd := resp.NewReader(p.conn)

	for {
		v, _, err := rd.ReadValue()

		if err == io.EOF {
			break
		}

		if v.Type() == resp.Array {
			for _, value := range v.Array() {
				switch value.String() {
				case CommandSet:
					if len(v.Array()) != 3 {
						return fmt.Errorf("invalid number of variables for SET command")
					}
					cmd := SetCommand{
						key: v.Array()[1].Bytes(),
						val: v.Array()[2].Bytes(),
					}

					p.msgChn <- Message{
						cmd:  cmd,
						peer: p,
					}

				case CommandGet:
					if len(v.Array()) != 2 {
						return fmt.Errorf("invalid number of variables for GET command")
					}
					cmd := GetCommand{
						key: v.Array()[1].Bytes(),
					}

					p.msgChn <- Message{
						cmd:  cmd,
						peer: p,
					}

				}
			}
		}
	}

	return nil
}

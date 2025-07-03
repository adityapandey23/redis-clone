package internal

import (
	"fmt"
	"github.com/tidwall/resp"
	"io"
	"net"
)

type Peer struct {
	Conn   net.Conn
	msgChn chan Message
	delChn chan *Peer
}

func NewPeer(conn net.Conn, msgChn chan Message, delChn chan *Peer) *Peer {
	return &Peer{
		Conn:   conn,
		msgChn: msgChn,
		delChn: delChn,
	}
}

func (p *Peer) Send(msg []byte) (int, error) {
	return p.Conn.Write(msg)
}

func (p *Peer) ReadLoop() error {
	rd := resp.NewReader(p.Conn)

	for {
		v, _, err := rd.ReadValue()

		if err == io.EOF {
			p.delChn <- p
			break
		}

		if v.Type() == resp.Array {
			for _, value := range v.Array() {
				var cmd Command
				switch value.String() {
				case CommandSet:
					if len(v.Array()) != 3 {
						return fmt.Errorf("invalid number of variables for SET command")
					}
					cmd = SetCommand{
						Key: v.Array()[1].Bytes(),
						Val: v.Array()[2].Bytes(),
					}
				case CommandGet:
					if len(v.Array()) != 2 {
						return fmt.Errorf("invalid number of variables for GET command")
					}
					cmd = GetCommand{
						Key: v.Array()[1].Bytes(),
					}
				}

				p.msgChn <- Message{
					Cmd:  cmd,
					Peer: p,
				}

			}
		}
	}
	return nil
}

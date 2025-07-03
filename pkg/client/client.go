package client

import (
	"bytes"
	"context"
	"io"
	"net"

	"github.com/tidwall/resp"
)

type Client struct {
	addr string
	conn net.Conn
}

func New(addr string) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Client{
		addr: addr,
		conn: conn,
	}, nil
}

func (c *Client) Set(ctx context.Context, key interface{}, val interface{}) error {

	buf := &bytes.Buffer{} // Didn't use var buf *bytes.Buffer because nil pointer, this will lead to runtime panic when we try to write something!
	wr := resp.NewWriter(buf)
	wr.WriteArray([]resp.Value{
		resp.AnyValue("SET"),
		resp.AnyValue(key),
		resp.AnyValue(val),
	})

	_, err := io.Copy(c.conn, buf) // Now whenever we dial the client, we can reuse the connection

	return err
}

func (c *Client) Get(ctx context.Context, key interface{}) (string, error) {

	buf := &bytes.Buffer{} // Didn't use var buf *bytes.Buffer because nil pointer, this will lead to runtime panic when we try to write something!
	wr := resp.NewWriter(buf)
	wr.WriteArray([]resp.Value{
		resp.AnyValue("GET"),
		resp.AnyValue(key),
	})

	_, err := io.Copy(c.conn, buf)

	if err != nil {
		return "", err
	}

	b := make([]byte, 1024)
	n, err := c.conn.Read(b)
	return string(b[:n]), err
}

func (c *Client) Close() error {
	return c.conn.Close()
}

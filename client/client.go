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
}

func New(addr string) *Client {
	return &Client{
		addr: addr,
	}
}

func (c *Client) Set(ctx context.Context, key string, val string) error {
	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		return err
	}

	buf := &bytes.Buffer{} // Didn't use var buf *bytes.Buffer because nil pointer, this will lead to runtime panic when we try to write something!
	wr := resp.NewWriter(buf)
	wr.WriteArray([]resp.Value{
		resp.StringValue("SET"),
		resp.StringValue(key),
		resp.StringValue(val),
	})

	_, err = io.Copy(conn, buf)

	return err
}

func (c *Client) Get(ctx context.Context, key string) (string, error) {
	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		return "", err
	}

	buf := &bytes.Buffer{} // Didn't use var buf *bytes.Buffer because nil pointer, this will lead to runtime panic when we try to write something!
	wr := resp.NewWriter(buf)
	wr.WriteArray([]resp.Value{
		resp.StringValue("GET"),
		resp.StringValue(key),
	})

	_, err = io.Copy(conn, buf)

	if err != nil {
		return "", err
	}

	b := make([]byte, 1024)
	n, err := conn.Read(b)
	return string(b[:n]), err
}

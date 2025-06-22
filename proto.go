package main

import (
	"bytes"
	"fmt"
	"io"

	"github.com/tidwall/resp"
)

const (
	CommandSet = "SET"
)

type Command interface {
}

type SetCommand struct {
	key, val string
}

func parseCommand(raw string) (Command, error) {
	rd := resp.NewReader(bytes.NewBufferString(raw))
	for {
		v, _, err := rd.ReadValue()

		if err == io.EOF {
			return nil, err
		}

		if v.Type() == resp.Array {
			for _, value := range v.Array() {
				switch value.String() {
				case CommandSet:
					if len(v.Array()) != 3 {
						return nil, fmt.Errorf("invalid number of variables for SET command")
					}
					cmd := SetCommand{
						key: v.Array()[1].String(),
						val: v.Array()[2].String(),
					}
					return cmd, nil
				}
			}
		}
		return nil, fmt.Errorf("invalid or unknown command recieved HALO: %s", raw)
	}
	return nil, fmt.Errorf("invalid or unknown command recieved HALO: %s", raw)
}

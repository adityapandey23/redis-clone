package internal

import (
	"bytes"
	"fmt"
)

const (
	CommandSet = "SET"
	CommandGet = "GET"
)

type Command interface{}

type SetCommand struct {
	Key, Val []byte
}

type GetCommand struct {
	Key []byte
}

func respWriteMap(m map[string]string) string {
	buf := bytes.Buffer{}
	buf.WriteString("%" + fmt.Sprintf("%d\r\n", len(m)))
	for k, v := range m {
		buf.WriteString(fmt.Sprintf("+%s\r\n", k))
		buf.WriteString(fmt.Sprintf(":%s\r\n", v))
	}

	return buf.String()
}

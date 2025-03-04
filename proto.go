package main

import (
	"bytes"
	"fmt"
	"io"
	"log"

	"github.com/tidwall/resp"
)

const (
	CommandSET   = "SET"
	CommandGET   = "GET"
	CommandHELLO = "hello"
)

type Command interface{}

type SetCommand struct {
	key, val []byte
}

type HelloCommand struct {
	value string
}

type GetCommand struct {
	key []byte
}

func parseCommand(raw string) (Command, error) {
	rd := resp.NewReader(bytes.NewBufferString(raw))

	for {
		v, _, err := rd.ReadValue()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if v.Type() == resp.Array {
			for _, value := range v.Array() {
				switch value.String() {
				case CommandGET:
					if len(v.Array()) != 2 {
						return nil, fmt.Errorf("invalid number of variables for GET command")
					}
					cmd := GetCommand{
						key: []byte(v.Array()[1].Bytes()),
					}
					return cmd, nil
				case CommandHELLO:
					cmd := HelloCommand{
						value: v.Array()[1].String(),
					}
					return cmd, nil
				}
			}
		}
	}
	return nil, fmt.Errorf("invalid or unknown command recieved: %s", raw)
}

func respWriteMap(m map[string]string) []byte {
	buf := bytes.Buffer{}
	buf.WriteString("%" + fmt.Sprintf("%d\n\r", len(m)))
	for k, v := range m {
		buf.WriteString(fmt.Sprintf("+%s\n\r", k))
		buf.WriteString(fmt.Sprintf(":%s\n\r", v))
	}
	return buf.Bytes()
}

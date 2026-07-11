package llmsdk

import (
	"io"
	"strings"
)

type Role string

const (
	User      Role = "user"
	System    Role = "system"
	Assistant Role = "assistant"
)

type Message struct {
	Role    Role
	Content string
}

type LLMResponse struct {
	Content strings.Builder
	Tokens  uint64
	data    []byte
	pos     int
}

func (resp *LLMResponse) Read(buf []byte) (int, error) {
	if resp.pos >= len(resp.data) {
		return 0, io.EOF
	}

	n := copy(buf, resp.data[resp.pos:])
	resp.pos = n

	return n, nil
}

type LLMProvider interface {
	Send(messages []Message) (LLMResponse, error)
	IsActive() bool
	GetModelInf() ModelInf
}

type ModelInf struct {
	Name string
}

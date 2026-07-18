package llmsdk

import (
	"io"
)

type Role string

const (
	User      Role = "user"
	System    Role = "system"
	Assistant Role = "assistant"
)

type Message struct {
	Role    Role   `json:"role"`
	Content string `json:"content"`
}

type LLMResponse struct {
	Message Message
	Tokens  uint64
	Done    bool
	Reader  io.Reader
}

type Settings struct {
	Temperature float64 `json:"temperature"`
	Seed        int     `json:"seed"`
	Stream      bool    `json:"stream"`
	TopK        int     `json:"top_k"`
	TopP        float64 `json:"top_p"`
	MinP        float64 `json:"min_p"`
	Stop        string  `json:"stop"`
	NumCtx      int     `json:"num_ctx"`
	NumPredict  int     `json:"num_predict"`
	Think       bool    `json:"think"`
	KeepAlive   int     `json:"keep_alive"`
}

type LLMProvider interface {
	Send(messages []Message) (*LLMResponse, error)
	GetModelInf() ModelInf
	SetSettings(Settings)
}

type ModelInf struct {
	Name string `json:"name"`
}

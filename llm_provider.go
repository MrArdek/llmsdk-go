package llmsdk

import (
	"io"
)

type Role string

const (
	User      Role = "user"
	System    Role = "system"
	Assistant Role = "assistant"
	Tool      Role = "tool"
)

type Message struct {
	Role    Role   `json:"role"`
	Content string `json:"content"`
}

type ChuckedMessage struct {
	Model string `json:"model"`
}

type LLMResponse struct {
	Message Message
	Tokens  uint64
	Done    bool
	Reader  io.Reader
	buf     []byte
}

func (r *LLMResponse) Read(buf []byte) (int, error) {
	n, err := r.Reader.Read(buf)

	if r.buf == nil {
		r.buf = make([]byte, 0, 1024)
	}

	r.buf = append(r.buf, buf[:n]...)

	if err != nil {
		if err == io.EOF {
			r.Message.Content = string(r.buf)
			r.Done = true
			r.Message.Role = Assistant
			return n, io.EOF
		}
		return 0, err
	}

	return n, nil
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
	GetModelInfo() ModelInfo
	SetSettings(*Settings)
}

type ModelInfo struct {
	Name string `json:"name"`
}

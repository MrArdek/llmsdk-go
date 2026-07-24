package llmsdk

import (
	"encoding/json"
	"io"
)

type (
	Role     string
	ToolType string
)

const (
	User         Role     = "user"
	System       Role     = "system"
	Assistant    Role     = "assistant"
	ToolRole     Role     = "tool"
	FunctionTool ToolType = "function"
)

type Message struct {
	Role       Role       `json:"role"`
	Content    string     `json:"content"`
	ToolCalls  []ToolCall `json:"tool_calls"`
	ToolCallID string     `json:"tool_call_id"`
}

type Tool struct {
	Type        ToolType            `json:"type"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	HandlerFunc func([]byte) []byte `json:"-"`
	Parameters  []Parameter         `json:"parameters"`
}

type ToolCall struct {
	ID        string          `json:"id"`
	CallID    string          `json:"call_id"`
	Name      string          `json:"name"`
	Type      string          `json:"type"`
	Arguments json.RawMessage `json:"arguments"`
}

type Parameter struct {
	Type        string      `json:"type"`
	Name        string      `json:"name,omitempty"`
	Description string      `json:"description,omitempty"`
	Properties  []Parameter `json:"properties,omitempty"`
	Required    []string    `json:"required,omitempty"`
}

type ChuckedMessage struct {
	Model  string `json:"model"`
	Done   bool   `json:"done"`
	Tokens uint64 `json:"tokens"`
}

type LLMResponse struct {
	Message Message
	Tokens  uint64
	Done    bool
	Reader  io.Reader
	buf     []byte
}

type LLMRequest struct {
	Messages []Message
	Tools    []Tool
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
	GetModel() Model
	SetModel(Model)
	SetSettings(*Settings)
}

type Model struct {
	Name string `json:"name"`
}

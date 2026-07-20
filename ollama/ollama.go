package ollama

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"llmsdk"
)

// Сделано для удобства, непредпологается что URL будет менятся часто
const baseURL = "http://localhost:11434/api/chat"

func GetOllamaProvider(modelName string, client *http.Client, settings llmsdk.Settings) llmsdk.LLMProvider {
	ol := Ollama{Model: llmsdk.ModelInf{Name: modelName}, HTTPClient: client, Settings: &settings}
	return &ol
}

type Ollama struct {
	Model      llmsdk.ModelInf
	HTTPClient *http.Client
	Settings   *llmsdk.Settings
}
type ollamaRequest struct {
	Model    string           `json:"model"`
	Messages []llmsdk.Message `json:"messages"`
	Stream   bool             `json:"stream"`
	// Options   Options          `json:"options"`
	Think     bool `json:"think"`
	KeepAlive int  `json:"keep_alive"`
	// LogProbs bool `json:"logprobs"`
	// TopLogprobs int `json:"top_logprobs"`
	// I think we don't need this
}
type ollamaResponse struct {
	Model              string         `json:"model"`
	CreatedAt          string         `json:"created_at"`
	Message            llmsdk.Message `json:"message"`
	Done               bool           `json:"done"`
	DoneReason         string         `json:"done_reason"`
	TotalDuration      int            `json:"total_duration"`
	LoadDuration       int            `json:"load_duration"`
	PromptEvalCount    int            `json:"prompt_eval_count"`
	PromptEvalDuration int            `json:"prompt_eval_duration"`
	EvalCount          int            `json:"eval_count"`
	EvalDuration       int            `json:"eval_duration"`
	// logprobs... I think we don't need this now
}
type contentReader struct {
	Data []byte
	Dec  *json.Decoder
	Pos  int
	Resp ollamaResponse
}
type Options struct {
	Seed        int     `json:"seed"`
	Temperature float64 `json:"temperature"`
	TopK        int     `json:"top_k"`
	TopP        float64 `json:"top_p"`
	MinP        float64 `json:"min_p"`
	Stop        string  `json:"stop"`
	NumCtx      int     `json:"num_ctx"`
	NumPredict  int     `json:"num_predict"`
}

func (c *contentReader) Read(buf []byte) (int, error) {
	if c.Pos >= len(c.Data) {
		err := c.Dec.Decode(&c.Resp)
		if err != nil {
			if err == io.EOF {
				return 0, io.EOF
			}

			return 0, err
		}

		c.Data = []byte(c.Resp.Message.Content)
		c.Pos = 0
	}

	n := copy(buf, c.Data[c.Pos:])
	c.Pos += n
	return n, nil
}

func (l *Ollama) SetSettings(settings llmsdk.Settings) {
	l.Settings = &settings
}

func (l *Ollama) Send(messages []llmsdk.Message) (*llmsdk.LLMResponse, error) {
	//opt := Options{
	//	Seed:        l.Settings.Seed,
	//	Temperature: l.Settings.Temperature,
	//	TopK:        l.Settings.TopK,
	//	TopP:        l.Settings.TopP,
	//	MinP:        l.Settings.MinP,
	//	Stop:        l.Settings.Stop,
	//	NumCtx:      l.Settings.NumCtx,
	//	NumPredict:  l.Settings.NumPredict,
	//}

	req := ollamaRequest{
		Model:    l.Model.Name,
		Messages: messages,
		Stream:   l.Settings.Stream,
		// Options:   opt,
		Think:     l.Settings.Think,
		KeepAlive: l.Settings.KeepAlive,
	}

	reqJSON, err := json.Marshal(req)
	if err != nil {
		log.Println(err.Error())
		return &llmsdk.LLMResponse{}, err
	}

	resp, err := l.HTTPClient.Post(baseURL, "application/json", bytes.NewReader(reqJSON))
	if err != nil {
		log.Println("error while getting resp")
		log.Println(err.Error())
		return &llmsdk.LLMResponse{}, err
	}
	log.Println(resp)

	dec := json.NewDecoder(resp.Body)
	defer resp.Body.Close()

	llmResp := llmsdk.LLMResponse{
		Message: llmsdk.Message{},
		Reader:  &contentReader{Data: make([]byte, 0), Dec: dec, Pos: 0},
	}
	return &llmResp, nil
}

func (l *Ollama) GetModelInf() llmsdk.ModelInf {
	return l.Model
}

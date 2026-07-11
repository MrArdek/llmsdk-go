package llmsdk

import "io"

type LLMSession struct {
	Chat        Chat
	LLMProvider LLMProvider
}
type Chat struct {
	Messages     []Message
	InputTokens  uint64
	OutputTokens uint64
}
type BufferedResponse struct {
	llmResp io.Reader
	bufR    []byte
	chat    *Chat
}

func (bresp *BufferedResponse) ReadB(buf []byte) (int, error) {
	n, err := bresp.llmResp.Read(buf)
	bresp.bufR = append(bresp.bufR, buf[:n]...)

	if err != nil {
		if err == io.EOF {
			bresp.chat.Messages = append(bresp.chat.Messages, Message{Role: Assistant, Content: string(bresp.bufR)})

			return n, err
		}

		return 0, err
	}

	return n, nil
}

func (chat *Chat) ClearChat() {
	*chat = Chat{}
}

func (chat *Chat) AddMessage(message Message) {
	chat.Messages = append(chat.Messages, message)
}

func (llm *LLMSession) LLMSend(message Message) (BufferedResponse, error) {
	llm.Chat.AddMessage(message)

	resp, err := llm.LLMProvider.Send(llm.Chat.Messages)
	if err != nil {
		return BufferedResponse{}, err
	}

	bufResp := BufferedResponse{llmResp: &resp, chat: &llm.Chat}

	return bufResp, nil
}

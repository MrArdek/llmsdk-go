package llmsdk

import (
	"fmt"
	"log"
)

type LLMSession struct {
	Chat        Chat
	LLMProvider LLMProvider
}
type Chat struct {
	Messages     []Message
	InputTokens  uint64
	OutputTokens uint64
}
type InvalidResponseError struct {
	Field  string
	Reason string
}

func (e InvalidResponseError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Reason)
}

func (chat *Chat) ClearChat() {
	chat.Messages = chat.Messages[:0]
}

func (chat *Chat) AddMessage(message Message) {
	chat.Messages = append(chat.Messages, message)
}

func (llm *LLMSession) LLMSend() (*LLMResponse, error) {
	if len(llm.Chat.Messages) <= 0 {
		return nil, InvalidResponseError{Field: "Messages", Reason: "cannot send request with empty messages!"}
	}

	resp, err := llm.LLMProvider.Send(llm.Chat.Messages)
	if err != nil {
		log.Println("error while sending message!")
		log.Println(err)
		return nil, err
	}

	return resp, nil
}

func (llm *LLMSession) LLMSendMessage(message Message) (*LLMResponse, error) {
	llm.Chat.AddMessage(message)

	resp, err := llm.LLMProvider.Send(llm.Chat.Messages)
	if err != nil {
		log.Println("error while sending message!")
		log.Println(err.Error())
		return nil, err
	}

	return resp, nil
}

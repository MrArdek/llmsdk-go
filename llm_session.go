package llmsdk

import (
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

func (chat *Chat) ClearChat() {
	chat.Messages = chat.Messages[:0]
}

func (chat *Chat) AddMessage(message Message) {
	chat.Messages = append(chat.Messages, message)
}

func (llm *LLMSession) LLMSend(message Message) (*LLMResponse, error) {
	llm.Chat.AddMessage(message)

	resp, err := llm.LLMProvider.Send(llm.Chat.Messages)
	if err != nil {
		log.Println("error while sending message!")
		log.Println(err.Error())
		return nil, err
	}

	return resp, nil
}

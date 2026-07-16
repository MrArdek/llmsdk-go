package main

import (
	"net/http"

	"llmsdk"
	"llmsdk/ollama"
)

func main() {
	ollamaPr := ollama.GetOllamaProvider("qwen2.5:14b", &http.Client{}, llmsdk.Settings{Stream: false})
	// Да помогут тебе утром все силы мира...

	llmSession := llmsdk.LLMSession{llmsdk.Chat{make([]llmsdk.Message, 0, 20), 0, 0}, ollamaPr}

	llmSession.LLMSend(llmsdk.Message{llmsdk.Assistant, "Привут! Ты"})
}

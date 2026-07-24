package openaicompatible

import (
	"net/http"

	"llmsdk"
)

type OpenAICompatibleProvider struct {
	baseURL    string
	Model      llmsdk.Model
	HTTPClient *http.Client
	Settings   *llmsdk.Settings
}

type openAIResponse struct{}

type openAIRequest struct{}

func (o *OpenAICompatibleProvider) Send(messages []llmsdk.Message) (*llmsdk.LLMResponse, error) {
	return nil, nil
}

func (o *OpenAICompatibleProvider) GetModel() llmsdk.Model {
	return o.Model
}

func (o *OpenAICompatibleProvider) SetModel(m llmsdk.Model) {
	o.Model = m
}

func (o *OpenAICompatibleProvider) SetSettings(settings *llmsdk.Settings) {
	o.Settings = settings
}

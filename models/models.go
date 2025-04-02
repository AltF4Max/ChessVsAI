package models

type APIResponse struct {
	ID       string   `json:"id"`
	Provider string   `json:"provider"`
	Model    string   `json:"model"`
	Object   string   `json:"object"`
	Created  int64    `json:"created"`
	Choices  []Choice `json:"choices"`
	Usage    Usage    `json:"usage"`
}

type Choice struct {
	Logprobs           interface{} `json:"logprobs"`
	FinishReason       string      `json:"finish_reason"`
	NativeFinishReason string      `json:"native_finish_reason"`
	Index              int         `json:"index"`
	Message            Message     `json:"message"`
}

type Message struct {
	Role      string  `json:"role"`
	Content   string  `json:"content"`
	Refusal   *string `json:"refusal"`
	Reasoning *string `json:"reasoning"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

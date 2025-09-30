package config

type LLMConfig struct {
	ApiUrl             string   `json:"api_url"`
	ApiKey             string   `json:"api_key"`
	MaxConcurrentTasks int      `json:"max_concurrent_tasks"`
	AvailableModels    []string `json:"available_models"`
}

type ChatParams struct {
	MaxTokens        int     `json:"max_tokens"`
	Temperature      float64 `json:"temperature"`
	TopP             float64 `json:"top_p"`
	FrequencyPenalty float64 `json:"frequency_penalty"`
	PresencePenalty  float64 `json:"presence_penalty"`
	RepeatPenalty    float64 `json:"repeat_penalty"`
}

var (
	LLMConfigure   LLMConfig
	ChatParameters ChatParams
)

func InitLLMConfig() {
	LLMConfigure = LLMConfig{
		ApiUrl:             getEnv("LLM_API_URL", "http://ai.api.maybered.com/AI-VMZ-8B/v1/chat/completions"),
		ApiKey:             getEnv("LLM_API_KEY", ""),
		MaxConcurrentTasks: getEnvInt("LLM_MAX_CONCURRENT_TASKS", 4),
		AvailableModels:    []string{"VMZ-8B"},
	}
	ChatParameters = ChatParams{
		MaxTokens:        getEnvInt("LLM_MAX_TOKENS", 4096),
		Temperature:      getEnvFloat("LLM_TEMPERATURE", 0.3),
		TopP:             getEnvFloat("LLM_TOP_P", 0.7),
		FrequencyPenalty: getEnvFloat("LLM_FREQUENCY_PENALTY", 0.0),
		PresencePenalty:  getEnvFloat("LLM_PRESENCE_PENALTY", 0.0),
		RepeatPenalty:    getEnvFloat("LLM_REPEAT_PENALTY", 1.05),
	}
}

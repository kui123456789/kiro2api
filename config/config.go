package config

import (
	"os"
	"strconv"
)

// ModelInfo 描述一个可请求的 Kiro 模型或兼容别名。
type ModelInfo struct {
	ID            string
	KiroModelID   string
	ContextWindow int
	Public        bool
}

var models = []ModelInfo{
	{ID: "claude-opus-4-7", KiroModelID: "claude-opus-4.7", ContextWindow: 1000000, Public: true},
	{ID: "claude-opus-4-6", KiroModelID: "claude-opus-4.6", ContextWindow: 1000000, Public: true},
	{ID: "claude-opus-4-5", KiroModelID: "claude-opus-4.5", ContextWindow: 200000, Public: true},
	{ID: "claude-sonnet-4-6", KiroModelID: "claude-sonnet-4.6", ContextWindow: 1000000, Public: true},
	{ID: "claude-sonnet-4-5", KiroModelID: "claude-sonnet-4.5", ContextWindow: 200000, Public: true},
	{ID: "claude-sonnet-4-0", KiroModelID: "CLAUDE_SONNET_4_20250514_V1_0", ContextWindow: 200000, Public: true},
	{ID: "claude-haiku-4-5", KiroModelID: "claude-haiku-4.5", ContextWindow: 200000, Public: true},
	{ID: "deepseek-3-2", KiroModelID: "deepseek-3.2", ContextWindow: 128000, Public: true},
	{ID: "minimax-m2-5", KiroModelID: "minimax-m2.5", ContextWindow: 200000, Public: true},
	{ID: "glm-5", KiroModelID: "glm-5", ContextWindow: 200000, Public: true},
	{ID: "minimax-m2-1", KiroModelID: "minimax-m2.1", ContextWindow: 200000, Public: true},
	{ID: "qwen3-coder-next", KiroModelID: "qwen3-coder-next", ContextWindow: 256000, Public: true},
	{ID: "claude-sonnet-4-5-20250929", KiroModelID: "CLAUDE_SONNET_4_5_20250929_V1_0", ContextWindow: 200000, Public: false},
	{ID: "claude-sonnet-4-20250514", KiroModelID: "CLAUDE_SONNET_4_20250514_V1_0", ContextWindow: 200000, Public: false},
	{ID: "claude-3-7-sonnet-20250219", KiroModelID: "CLAUDE_3_7_SONNET_20250219_V1_0", ContextWindow: 200000, Public: false},
	{ID: "claude-3-5-haiku-20241022", KiroModelID: "auto", ContextWindow: 200000, Public: false},
	{ID: "claude-haiku-4-5-20251001", KiroModelID: "claude-haiku-4.5", ContextWindow: 200000, Public: false},
}

var (
	modelMap            = buildModelMap(models)
	modelContextWindows = buildModelContextWindows(models)
	publicModels        = buildPublicModels(models)
)

// PublicModels returns the Kiro models exposed by /v1/models.
func PublicModels() []ModelInfo {
	return append([]ModelInfo(nil), publicModels...)
}

// LookupModel returns the upstream Kiro model ID for a public model or compatibility alias.
func LookupModel(modelID string) (string, bool) {
	kiroModelID, exists := modelMap[modelID]
	return kiroModelID, exists
}

// ContextWindow returns the context window for a public model or compatibility alias.
func ContextWindow(modelID string) (int, bool) {
	contextWindow, exists := modelContextWindows[modelID]
	return contextWindow, exists
}

func buildModelMap(models []ModelInfo) map[string]string {
	builtModelMap := make(map[string]string, len(models))
	for _, model := range models {
		builtModelMap[model.ID] = model.KiroModelID
	}
	return builtModelMap
}

func buildModelContextWindows(models []ModelInfo) map[string]int {
	contextWindows := make(map[string]int, len(models))
	for _, model := range models {
		contextWindows[model.ID] = model.ContextWindow
	}
	return contextWindows
}

func buildPublicModels(models []ModelInfo) []ModelInfo {
	publicModels := make([]ModelInfo, 0, len(models))
	for _, model := range models {
		if model.Public {
			publicModels = append(publicModels, model)
		}
	}
	return publicModels
}

// RefreshTokenURL 刷新token的URL (social方式)
const RefreshTokenURL = "https://prod.us-east-1.auth.desktop.kiro.dev/refreshToken"

// IdcRefreshTokenURL IdC认证方式的刷新token URL
const IdcRefreshTokenURL = "https://oidc.us-east-1.amazonaws.com/token"

// CodeWhispererURL CodeWhisperer API的URL
const CodeWhispererURL = "https://codewhisperer.us-east-1.amazonaws.com/generateAssistantResponse"

// MaxToolDescriptionLength 工具描述的最大长度（字符数）
// 可通过环境变量 MAX_TOOL_DESCRIPTION_LENGTH 配置，默认 10000
var MaxToolDescriptionLength = getEnvIntWithDefault("MAX_TOOL_DESCRIPTION_LENGTH", 10000)

// getEnvIntWithDefault 获取整数类型环境变量（带默认值）
func getEnvIntWithDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

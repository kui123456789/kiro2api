package config

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPublicModels_ContainsOfficialKiroModelsExceptAuto(t *testing.T) {
	expectedModels := map[string]struct {
		kiroModelID string
		maxTokens   int
	}{
		"claude-opus-4-7":   {"claude-opus-4.7", 1000000},
		"claude-opus-4-6":   {"claude-opus-4.6", 1000000},
		"claude-opus-4-5":   {"claude-opus-4.5", 200000},
		"claude-sonnet-4-6": {"claude-sonnet-4.6", 1000000},
		"claude-sonnet-4-5": {"claude-sonnet-4.5", 200000},
		"claude-sonnet-4-0": {"CLAUDE_SONNET_4_20250514_V1_0", 200000},
		"claude-haiku-4-5":  {"claude-haiku-4.5", 200000},
		"deepseek-3-2":      {"deepseek-3.2", 128000},
		"minimax-m2-5":      {"minimax-m2.5", 200000},
		"glm-5":             {"glm-5", 200000},
		"minimax-m2-1":      {"minimax-m2.1", 200000},
		"qwen3-coder-next":  {"qwen3-coder-next", 256000},
	}

	publicModels := PublicModels()
	actualModels := make(map[string]ModelInfo, len(publicModels))
	for _, model := range publicModels {
		_, exists := actualModels[model.ID]
		assert.False(t, exists, "Model %s should not be duplicated", model.ID)
		actualModels[model.ID] = model
	}

	assert.Len(t, actualModels, len(expectedModels))
	for modelID, expected := range expectedModels {
		model, exists := actualModels[modelID]
		assert.True(t, exists, "Model %s should be an official public Kiro model", modelID)
		assert.Equal(t, expected.kiroModelID, model.KiroModelID, "Model %s should map to the expected Kiro model ID", modelID)
		assert.Equal(t, expected.maxTokens, model.ContextWindow, "Model %s should expose the official context window", modelID)
		assert.True(t, model.Public)
	}
}

func TestModelMap_PreservesLegacyAliases(t *testing.T) {
	expectedAliases := map[string]string{
		"claude-sonnet-4-5-20250929": "CLAUDE_SONNET_4_5_20250929_V1_0",
		"claude-sonnet-4-20250514":   "CLAUDE_SONNET_4_20250514_V1_0",
		"claude-3-7-sonnet-20250219": "CLAUDE_3_7_SONNET_20250219_V1_0",
		"claude-3-5-haiku-20241022":  "auto",
		"claude-haiku-4-5-20251001":  "claude-haiku-4.5",
	}

	for model, expectedKiroModelID := range expectedAliases {
		kiroModelID, exists := LookupModel(model)
		assert.True(t, exists, "Legacy alias %s should remain accepted", model)
		assert.Equal(t, expectedKiroModelID, kiroModelID)
	}
}

func TestPublicModels_ExcludesAutoModel(t *testing.T) {
	for _, model := range PublicModels() {
		assert.NotEqual(t, "auto", model.ID)
		assert.NotEqual(t, "auto", model.KiroModelID)
	}
}

func TestModelMap_NonExistentModel(t *testing.T) {
	_, exists := LookupModel("non-existent-model")
	assert.False(t, exists)
}

func TestModelMap_AllModelsHaveContextWindow(t *testing.T) {
	for model := range modelMap {
		maxTokens, exists := ContextWindow(model)
		assert.True(t, exists, "Model %s should have a context window", model)
		assert.Greater(t, maxTokens, 0, "Model %s should have a positive context window", model)
	}
}

func TestModelMap_MappingsAreCorrectFormat(t *testing.T) {
	kiroModelIDs := map[string]struct{}{
		"claude-opus-4.5":   {},
		"claude-opus-4.6":   {},
		"claude-opus-4.7":   {},
		"claude-sonnet-4.5": {},
		"claude-sonnet-4.6": {},
		"claude-haiku-4.5":  {},
		"deepseek-3.2":      {},
		"minimax-m2.5":      {},
		"glm-5":             {},
		"minimax-m2.1":      {},
		"qwen3-coder-next":  {},
		"auto":              {},
	}

	for inputModel, outputModel := range modelMap {
		if _, ok := kiroModelIDs[outputModel]; ok {
			continue
		}

		assert.False(t, strings.HasPrefix(outputModel, "claude-"),
			"Model mapping for %s should use a known Kiro model ID", inputModel)
		assert.Contains(t, outputModel, "CLAUDE",
			"Model mapping for %s should contain 'CLAUDE'", inputModel)
		assert.Contains(t, outputModel, "_V1_0",
			"Model mapping for %s should contain '_V1_0'", inputModel)
	}
}

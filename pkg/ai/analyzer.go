package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// AnalysisResult represents the structured output from AI
type AnalysisResult struct {
	PrimaryIssue         string         `json:"primary_issue"`
	SecondaryEffects     []string       `json:"secondary_effects"`
	FirstSeen            string         `json:"first_seen"`
	LikelyCauses         []Cause        `json:"likely_causes"`
	RecommendedActions   []string       `json:"recommended_actions"`
	SimilarPastIncidents []PastIncident `json:"similar_past_incidents"`
}

type Cause struct {
	Cause      string  `json:"cause"`
	Confidence float64 `json:"confidence"`
}

type PastIncident struct {
	Date       string `json:"date"`
	Resolution string `json:"resolution"`
}

// AnalyzeLogs sends logs to Google AI Studio and returns structured analysis
func AnalyzeLogs(logText string) (*AnalysisResult, error) {
	apiKey := os.Getenv("GOOGLE_AI_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GOOGLE_AI_KEY environment variable not set")
	}

	// Build the prompt
	prompt := buildPrompt(logText)

	// Call Google AI Studio API
	response, err := callGoogleAI(apiKey, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to call Google AI: %w", err)
	}

	// Parse the response
	var result AnalysisResult
	if err := json.Unmarshal([]byte(response), &result); err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}

	return &result, nil
}

// buildPrompt creates the prompt for the AI
func buildPrompt(logText string) string {
	systemPrompt := `You are an incident triage system.
Analyze the logs given. Respond ONLY in JSON format matching the schema.
Do not explain casually. Be concise, factual, and assign confidence to each likely cause.

Schema:
{
  "primary_issue": "string",
  "secondary_effects": ["string"],
  "first_seen": "ISO 8601 timestamp",
  "likely_causes": [{"cause": "string", "confidence": 0.0-1.0}],
  "recommended_actions": ["string"],
  "similar_past_incidents": [{"date": "YYYY-MM-DD", "resolution": "string"}]
}

Logs to analyze:
`
	return systemPrompt + "\n" + logText
}

// callGoogleAI makes the API call to Google AI Studio
func callGoogleAI(apiKey, prompt string) (string, error) {
	// Using Google AI Studio Gemini API
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent?key=%s", apiKey)

	requestBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{"text": prompt},
				},
			},
		},
		"generationConfig": map[string]interface{}{
			"temperature":     0.2,
			"maxOutputTokens": 2048,
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Parse Google AI response
	var aiResponse struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.Unmarshal(body, &aiResponse); err != nil {
		return "", err
	}

	if len(aiResponse.Candidates) == 0 || len(aiResponse.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response from AI")
	}

	return aiResponse.Candidates[0].Content.Parts[0].Text, nil
}

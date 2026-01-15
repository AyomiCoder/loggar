package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// AnalysisResult represents the structured output from AI
type AnalysisResult struct {
	Summary  string    `json:"summary"`
	Sections []Section `json:"sections"`
}

type Section struct {
	Title   string   `json:"title"`
	Content []string `json:"content"`
}

// ... (Cause and PastIncident removed if not used in new schema, but let's keep it simple for now)

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
	systemPrompt := `You are a world-class Principal Software Engineer (L8+). 
Analyze the provided logs with surgical precision. Your triage must be high-signal, concise, and intellectually punchy.

Rules:
1. "summary": A single, dense paragraph. Synthesize the failure into a technical narrative. No filler.
2. "sections": Provide exactly 2-3 sections. Use titles that reflect high-level architecture (e.g., "CORE DIAGNOSIS", "IMMEDIATE RESOLUTION").
3. Expert Parsimony: Use fewer words to say more. Avoid generic "potential causes" list. Focus on the most probable architectural or code-level failure.
4. Technical Depth: If you see a stack trace or code path, call out the exact point of failure and why it's likely occurring (e.g., "race condition in connection pooling handler").
5. Respond ONLY in JSON.

Schema:
{
  "summary": "string",
  "sections": [
    {
      "title": "string",
      "content": ["string"]
    }
  ]
}

Logs to analyze:
`
	return systemPrompt + "\n" + logText
}

// callGoogleAI makes the API call to Google AI Studio with exponential backoff retry
func callGoogleAI(apiKey, prompt string) (string, error) {
	// Using Google AI Studio Gemini 3 Flash Preview (Experimental model with high quota)
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-3-flash-preview:generateContent?key=%s", apiKey)

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

	var lastErr error
	var body []byte

	// Retry mechanism: 5 attempts with exponential backoff
	for i := 0; i < 5; i++ {
		if i > 0 {
			// Exponential backoff: 1s, 2s, 4s, 8s, 10s (max)
			delay := time.Duration(1<<uint(i-1)) * time.Second
			if delay > 10*time.Second {
				delay = 10 * time.Second
			}
			time.Sleep(delay)
		}

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			lastErr = fmt.Errorf("network error: %w", err)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			body, _ = io.ReadAll(resp.Body)
			resp.Body.Close()

			// Only retry for rate limits (429) or server errors (500/503)
			if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500 {
				lastErr = fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
				continue
			}
			return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
		}

		body, err = io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = fmt.Errorf("failed to read response body: %w", err)
			continue
		}

		// If we got here, we have a successful response
		break
	}

	if body == nil {
		return "", fmt.Errorf("all retry attempts failed: %w", lastErr)
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
		return "", fmt.Errorf("no response from AI after parsing")
	}

	text := aiResponse.Candidates[0].Content.Parts[0].Text

	// Sanitize response (remove possible markdown code blocks)
	text = sanitizeJSON(text)

	return text, nil
}

func sanitizeJSON(input string) string {
	input = strings.TrimSpace(input)
	// Remove ```json and ``` if present
	input = strings.TrimPrefix(input, "```json")
	input = strings.TrimPrefix(input, "```")
	input = strings.TrimSuffix(input, "```")
	return strings.TrimSpace(input)
}

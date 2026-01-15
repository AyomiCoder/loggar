package output

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
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

// PrintAnalysis prints the analysis result in a pretty terminal format
func PrintAnalysis(result *AnalysisResult) {
	separator := strings.Repeat("━", 50)

	// Colors
	titleColor := color.New(color.FgCyan, color.Bold)
	primaryColor := color.New(color.FgRed, color.Bold)
	bulletColor := color.New(color.FgYellow)
	actionColor := color.New(color.FgGreen)

	// Primary Issue
	fmt.Println(separator)
	titleColor.Println("PRIMARY ISSUE")
	primaryColor.Printf("→ %s\n", result.PrimaryIssue)

	// Secondary Effects
	if len(result.SecondaryEffects) > 0 {
		fmt.Println(separator)
		titleColor.Println("SECONDARY EFFECTS")
		for _, effect := range result.SecondaryEffects {
			bulletColor.Printf("• %s\n", effect)
		}
	}

	// Likely Causes
	if len(result.LikelyCauses) > 0 {
		fmt.Println(separator)
		titleColor.Println("LIKELY CAUSES")
		for i, cause := range result.LikelyCauses {
			confidence := int(cause.Confidence * 100)
			fmt.Printf("%d. %s (%d%%)\n", i+1, cause.Cause, confidence)
		}
	}

	// Recommended Actions
	if len(result.RecommendedActions) > 0 {
		fmt.Println(separator)
		titleColor.Println("RECOMMENDED ACTIONS")
		for _, action := range result.RecommendedActions {
			actionColor.Printf("• %s\n", action)
		}
	}

	// Similar Past Incidents
	if len(result.SimilarPastIncidents) > 0 {
		fmt.Println(separator)
		titleColor.Println("SIMILAR PAST INCIDENTS")
		for _, incident := range result.SimilarPastIncidents {
			fmt.Printf("• %s – %s\n", incident.Date, incident.Resolution)
		}
	}

	fmt.Println(separator)
}

// PrintJSON prints the raw JSON output
func PrintJSON(jsonData string) {
	fmt.Println(jsonData)
}

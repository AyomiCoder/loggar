package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"golang.org/x/term"
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

// getTermWidth returns a comfortable reading width, clamped between 80 and 120
func getTermWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	// If detection fails or width is too small, use a solid standard
	if err != nil || width < 80 {
		return 110 // "Sweet spot" for technical reading
	}
	// If the terminal is very wide, clamp to 120 so lines don't get too long to read
	if width > 120 {
		return 120
	}
	// Otherwise, use almost the full width
	return width - 2
}

// wrapText wraps the given text to the specified width and adds an optional indent for subsequent lines
func wrapText(text string, width int, indent string) []string {
	if width <= 0 {
		return []string{text}
	}
	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{""}
	}

	var lines []string
	var currentLine strings.Builder
	currentWidth := 0

	for i, word := range words {
		wordWidth := len(word)
		// Account for space between words
		spaceWidth := 0
		if currentLine.Len() > 0 {
			spaceWidth = 1
		}

		// Check if adding this word (and a space) exceeds the width
		// We use i > 0 to account for initial indent if needed,
		// but here we handle physical wrapping.
		if currentWidth+spaceWidth+wordWidth > width {
			lines = append(lines, currentLine.String())
			currentLine.Reset()
			currentLine.WriteString(indent)
			currentWidth = len(indent)
		} else if i > 0 && currentLine.Len() > 0 {
			currentLine.WriteString(" ")
			currentWidth++
		}

		currentLine.WriteString(word)
		currentWidth += wordWidth
	}
	if currentLine.Len() > 0 {
		lines = append(lines, currentLine.String())
	}
	return lines
}

// slowPrintWrapped prints wrapped lines character by character
func slowPrintWrapped(lines []string, delay time.Duration) {
	for i, line := range lines {
		for _, c := range line {
			fmt.Printf("%c", c)
			time.Sleep(delay)
		}
		// If there's a next line, print a newline.
		// If it was the last line, PrintAnalysis adds its own newlines to separate sections.
		if i < len(lines)-1 {
			fmt.Println()
		}
	}
	fmt.Println()
}

// lineDelay adds a small pause between sections
func lineDelay() {
	time.Sleep(150 * time.Millisecond)
}

// PrintAnalysis prints the analysis result in a pretty terminal format with a progressive effect
func PrintAnalysis(result *AnalysisResult) {
	// Colors
	titleColor := color.New(color.FgCyan, color.Bold)
	summaryColor := color.New(color.FgHiWhite, color.Italic)
	bulletColor := color.New(color.FgHiWhite)               // Toned down from yellow to white
	versionColor := color.New(color.FgHiBlack, color.Faint) // Faint for footer

	termWidth := getTermWidth()

	fmt.Println() // Start with a newline

	// Summary
	if result.Summary != "" {
		summaryColor.Print("💡 ")
		lines := wrapText(result.Summary, termWidth-3, "   ")
		slowPrintWrapped(lines, 12*time.Millisecond)
		lineDelay()
	}

	// Dynamic Sections
	for _, section := range result.Sections {
		titleColor.Println(strings.ToUpper(section.Title))
		for _, item := range section.Content {
			// Determine bullet and clean item
			bullet := "• "
			cleanItem := item
			if strings.HasPrefix(item, "→") {
				bullet = "→ "
				cleanItem = strings.TrimSpace(strings.TrimPrefix(item, "→"))
			} else if strings.HasPrefix(item, "•") {
				bullet = "• "
				cleanItem = strings.TrimSpace(strings.TrimPrefix(item, "•"))
			}

			bulletColor.Print(bullet)
			// Wrap item content with indent to align with bullet
			lines := wrapText(cleanItem, termWidth-3, "  ")
			slowPrintWrapped(lines, 40*time.Millisecond)
		}
		lineDelay()
		fmt.Println()
	}

	// Footer
	versionColor.Println("loggar v1.0.0")
	fmt.Println()
}

// PrintJSON prints the raw JSON output with indentation
func PrintJSON(jsonData string) {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(jsonData), "", "  ")
	if err != nil {
		fmt.Println(jsonData)
		return
	}
	fmt.Println(out.String())
}

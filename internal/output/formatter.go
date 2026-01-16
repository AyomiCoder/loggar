package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
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

// highlightLine applies colors to specific patterns in the text
func highlightLine(text string) string {
	// Highlighting rules
	// 1. Durations and Percentages (e.g. 104ms, 10-15s, 98%) -> Yellow
	text = replacePattern(text, `\b\d+(?:\.\d+)?(?:ms|s|%|kb|mb)\b`, color.New(color.FgHiYellow).SprintFunc())

	// 2. IDs and Codes (e.g. TX_9921, user_99a82, /v1/payment_intents) -> Cyan
	text = replacePattern(text, `\b[A-Za-z0-9_/-]{4,}\d+[A-Za-z0-9_/-]*\b`, color.New(color.FgHiCyan).SprintFunc())

	// 3. Quoted text -> Cyan
	text = replacePattern(text, `"[^"]+"`, color.New(color.FgHiCyan).SprintFunc())

	// 4. Severity words -> Red
	text = replacePattern(text, `(?i)\b(timeout|failed|failure|error|critical|collapsed|prohibited|refused)\b`, color.New(color.FgHiRed).SprintFunc())

	// 5. Success words -> Green
	text = replacePattern(text, `(?i)\b(success|resolved|healthy|stable|ok)\b`, color.New(color.FgHiGreen).SprintFunc())

	return text
}

// replacePattern is a helper to replace regex matches with a colored version
func replacePattern(text string, pattern string, colorFunc func(a ...interface{}) string) string {
	re := regexp.MustCompile(pattern)
	return re.ReplaceAllStringFunc(text, func(match string) string {
		return colorFunc(match)
	})
}

// PrintAnalysis prints the analysis result in a pretty terminal format with a progressive effect
func PrintAnalysis(result *AnalysisResult) {
	// Structural Colors
	titleColor := color.New(color.FgHiMagenta, color.Bold)
	summaryColor := color.New(color.FgWhite) // Cleaner white
	bulletColor := color.New(color.FgHiYellow)
	arrowColor := color.New(color.FgHiRed)
	versionColor := color.New(color.FgHiBlack, color.Faint)

	termWidth := getTermWidth()

	fmt.Println()

	// Summary
	if result.Summary != "" {
		summaryColor.Print("💡 ")
		lines := wrapText(result.Summary, termWidth-3, "   ")
		// Highlight summary content
		for i, line := range lines {
			lines[i] = highlightLine(line)
		}
		slowPrintWrapped(lines, 5*time.Millisecond) // Faster typing
		lineDelay()
	}

	// Dynamic Sections
	for _, section := range result.Sections {
		// Colorize title based on content?
		tColor := titleColor
		if strings.Contains(strings.ToUpper(section.Title), "RESOLUTION") {
			tColor = color.New(color.FgHiBlue, color.Bold)
		}

		tColor.Println(strings.ToUpper(section.Title))

		for _, item := range section.Content {
			bullet := "• "
			cleanItem := item
			var bColor *color.Color = bulletColor

			if strings.HasPrefix(item, "→") {
				bullet = "→ "
				cleanItem = strings.TrimSpace(strings.TrimPrefix(item, "→"))
				bColor = arrowColor
			} else if strings.HasPrefix(item, "•") {
				bullet = "• "
				cleanItem = strings.TrimSpace(strings.TrimPrefix(item, "•"))
			}

			bColor.Print(bullet)

			// Wrap item content
			lines := wrapText(cleanItem, termWidth-3, "  ")

			// Highlight content in lines
			for i := range lines {
				lines[i] = highlightLine(lines[i])
			}

			slowPrintWrapped(lines, 15*time.Millisecond) // Faster typing
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

package cmd

import (
	"fmt"
	"syscall"
	"time"

	"github.com/fatih/color"
	"golang.org/x/term"
)

// readPassword reads a password from the terminal without echoing
func readPassword(prompt string) (string, error) {
	fmt.Print(prompt)
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println() // Print a newline after password entry
	if err != nil {
		return "", err
	}
	return string(bytePassword), nil
}

// Progress represents an animated progress indicator
type Progress struct {
	message string
	stop    chan bool
}

// StartProgress starts an animated progress indicator
func StartProgress(msg string) *Progress {
	p := &Progress{
		message: msg,
		stop:    make(chan bool),
	}

	go func() {
		ticker := time.NewTicker(400 * time.Millisecond)
		defer ticker.Stop()
		dots := 0
		start := time.Now()
		msgColor := color.New(color.FgCyan)

		// Print initial message immediately
		fmt.Printf("\r\033[K")
		msgColor.Printf("→ %s", msg)

		for {
			select {
			case <-p.stop:
				return
			case <-ticker.C:
				dots = (dots + 1) % 4
				dotStr := ""
				for i := 0; i < dots; i++ {
					dotStr += "."
				}

				// Add an "almost done" message after 3 seconds
				displayMsg := p.message
				if time.Since(start) > 3*time.Second {
					displayMsg = "Almost done..."
				}

				// Clear line and print
				fmt.Printf("\r\033[K") // Clear line
				msgColor.Printf("→ %s%-3s", displayMsg, dotStr)
			}
		}
	}()

	return p
}

// Stop stops the progress indicator
func (p *Progress) Stop() {
	// Use a non-blocking send or close to avoid hanging if already stopped
	select {
	case p.stop <- true:
	default:
	}
	fmt.Printf("\r\033[K") // Clear the progress line
}

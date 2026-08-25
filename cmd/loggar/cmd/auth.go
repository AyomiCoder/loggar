package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AyomiCoder/loggar/internal/config"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate via GitHub/Google",
	Long:  "Opens your browser to authenticate with GitHub or Google and saves your access token.",
	Run:   runAuth,
}

var resetFlag bool

func init() {
	authCmd.Flags().BoolVar(&resetFlag, "reset", false, "Clear saved token")
	rootCmd.AddCommand(authCmd)
}

func runAuth(cmd *cobra.Command, args []string) {
	// Handle reset flag
	if resetFlag {
		if err := config.ClearToken(); err != nil {
			color.Red("✗ Failed to clear token: %v", err)
			os.Exit(1)
		}
		color.Green("✓ Token cleared successfully")
		return
	}

	// Branding Header
	headerColor := color.New(color.FgHiCyan, color.Bold)
	headerColor.Println(`
    __    ____  ___________  ___  ____ 
   / /   / __ \/ ____/ ____|/   |/ __ \
  / /   / / / / / __/ / __ / /| / /_/ /
 / /___/ /_/ / /_/ / /_/ // ___ / _, _/ 
/_____/\____/\____/\____//_/  |/_/ |_|  
	`)
	color.White("           AI-POWERED LOG TRIAGE\n\n")

	// 1. Select Provider (Interactive)
	var choice string
	prompt := &survey.Select{
		Message: "Select Identity Provider:",
		Options: []string{"GitHub", "Google"},
		Help:    "Use arrow keys to select your preferred auth provider",
	}
	err := survey.AskOne(prompt, &choice)
	if err != nil {
		color.Red("✗ Selection interrupted")
		os.Exit(1)
	}

	var providerPath string
	switch choice {
	case "GitHub":
		providerPath = "/auth/github"
	case "Google":
		providerPath = "/auth/google"
	}

	// 2. Start local listener
	tokenChan := make(chan string)
	emailChan := make(chan string)

	// Create a multiplexer
	mux := http.NewServeMux()
	server := &http.Server{Addr: ":10999", Handler: mux}

	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		email := r.URL.Query().Get("email")

		if token != "" {
			tokenChan <- token
			emailChan <- email

			// Show success page
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprint(w, `
				<!DOCTYPE html>
				<html>
				<head>
					<title>Authentication Successful - Loggar</title>
					<link href="https://fonts.googleapis.com/css2?family=Montserrat:wght@300;400;600;700&display=swap" rel="stylesheet">
					<style>
						body {
							background-color: #000;
							color: #fff;
							font-family: 'Montserrat', sans-serif;
							display: flex;
							align-items: center;
							justify-content: center;
							height: 100vh;
							margin: 0;
						}
						.card {
							background: #0a0a0a;
							border: 1px solid #1a1a1a;
							border-radius: 12px;
							padding: 3rem 4rem;
							text-align: center;
							box-shadow: 0 30px 60px rgba(0,0,0,0.5);
							max-width: 400px;
							width: 100%;
						}
						.logo {
							font-size: 1.5rem;
							font-weight: 700;
							margin-bottom: 2rem;
							letter-spacing: -0.02em;
						}
						.logo span { color: #444; }
						.icon {
							color: #27c93f;
							margin-bottom: 1.5rem;
						}
						h1 {
							font-size: 1.5rem;
							margin-bottom: 1rem;
							font-weight: 600;
						}
						p {
							color: #666;
							font-size: 0.9rem;
							line-height: 1.6;
						}
						.cmd {
							font-family: 'JetBrains Mono', monospace;
							background: #111;
							padding: 0.5rem;
							border-radius: 4px;
							color: #888;
							font-size: 0.8rem;
							margin-top: 2rem;
							display: inline-block;
						}
					</style>
				</head>
				<body>
					<div class="card">
						<div class="logo">loggar</div>
						<div class="icon">
							<svg width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
								<path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path>
								<polyline points="22 4 12 14.01 9 11.01"></polyline>
							</svg>
						</div>
						<h1>Authenticated</h1>
						<p>You have successfully logged in.<br>You can now close this tab and return to your terminal.</p>
						<div class="cmd">session active</div>
					</div>
				</body>
				</html>
			`)
		} else {
			http.Error(w, "Missing token in callback", http.StatusBadRequest)
		}
	})

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			color.Red("✗ Failed to start local listener: %v", err)
			os.Exit(1)
		}
	}()

	// 3. Open Browser
	apiURL := os.Getenv("API_URL")
	if apiURL == "" {
		apiURL = "https://api.loggar.space"
	}

	authURL := fmt.Sprintf("%s%s?cli_port=10999", apiURL, providerPath)
	color.Cyan("Opening browser to authenticate...")
	color.White("If it doesn't open, visit: %s", authURL)

	if err := openBrowser(authURL); err != nil {
		color.Yellow("Warning: Failed to open browser automatically: %v", err)
	}

	// 3. Wait for token
	prog := StartProgress("Waiting for authentication...")

	select {
	case token := <-tokenChan:
		email := <-emailChan
		prog.Stop()

		// 4. Save Token
		if err := config.SaveToken(token, email); err != nil {
			color.Red("\n✗ Failed to save token: %v", err)
			os.Exit(1)
		}

		color.Green("\n✓ Successfully authenticated as %s", email)

		// Shutdown server gracefully
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		server.Shutdown(ctx)

	case <-time.After(5 * time.Minute):
		prog.Stop()
		color.Red("\n✗ Authentication timed out. Please try again.")
		os.Exit(1)
	}
}

func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/AyomiCoder/loggar/internal/config"
	"github.com/AyomiCoder/loggar/internal/output"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze [file]",
	Short: "Analyze log files with AI",
	Long:  "Analyze log files or stdin to identify issues and get recommendations",
	Run:   runAnalyze,
}

var (
	jsonOutput bool
	stdinInput bool
)

func init() {
	analyzeCmd.Flags().BoolVar(&jsonOutput, "json", false, "Output raw JSON")
	analyzeCmd.Flags().BoolVar(&stdinInput, "stdin", false, "Read from stdin")
	rootCmd.AddCommand(analyzeCmd)
}

func runAnalyze(cmd *cobra.Command, args []string) {
	var logContent string
	var err error

	// Read logs from file or stdin
	if stdinInput || len(args) == 0 {
		// Read from stdin
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			color.Red("✗ Failed to read from stdin: %v", err)
			os.Exit(1)
		}
		logContent = string(data)
	} else {
		// Read from file
		data, err := os.ReadFile(args[0])
		if err != nil {
			color.Red("✗ Failed to read file: %v", err)
			os.Exit(1)
		}
		logContent = string(data)
	}

	if logContent == "" {
		color.Red("✗ No log content provided")
		os.Exit(1)
	}

	prog := StartProgress(fmt.Sprintf("Analyzing %d bytes of logs", len(logContent)))

	// Load token
	cfg, err := config.LoadToken()
	if err != nil {
		prog.Stop()
		color.Red("✗ Not authenticated. Run 'loggar auth' first")
		os.Exit(1)
	}

	// Make API request
	apiURL := os.Getenv("API_URL")
	if apiURL == "" {
		apiURL = "https://api.loggar.space"
	}

	requestData := map[string]string{
		"logs": logContent,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		prog.Stop()
		color.Red("✗ Failed to prepare request: %v", err)
		os.Exit(1)
	}

	req, err := http.NewRequest("POST", apiURL+"/api/analyze", bytes.NewBuffer(jsonData))
	if err != nil {
		prog.Stop()
		color.Red("✗ Failed to create request: %v", err)
		os.Exit(1)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		prog.Stop()
		color.Red("✗ Failed to connect to server: %v", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		prog.Stop()
		color.Red("✗ Failed to read response: %v", err)
		os.Exit(1)
	}

	prog.Stop()

	if resp.StatusCode != http.StatusOK {
		color.Red("✗ Analysis failed: %s", string(body))
		os.Exit(1)
	}

	// Output results
	if jsonOutput {
		output.PrintJSON(string(body))
	} else {
		var result output.AnalysisResult
		if err := json.Unmarshal(body, &result); err != nil {
			color.Red("✗ Failed to parse response: %v", err)
			os.Exit(1)
		}
		output.PrintAnalysis(&result)
	}
}

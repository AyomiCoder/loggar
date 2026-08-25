package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "loggar",
	Short: "AI-powered log triage for Mac Terminal",
	Long: `Loggar.dev - Analyze server logs and identify root causes with AI

Loggar helps developers quickly analyze and triage server logs,
identifying root causes, secondary effects, and recommended actions.`,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

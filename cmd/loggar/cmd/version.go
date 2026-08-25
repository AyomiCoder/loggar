package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version = "0.2.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  "Print the version number of Loggar CLI",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Loggar CLI v%s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "forgego",
	Short: "ForgeGo - A powerful Go project scaffolding tool",
	Long: `ForgeGo is a CLI tool that helps developers quickly scaffold 
production-ready Go projects for different purposes including 
RESTful API backends, CLI applications, and Go libraries.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Add global flags here if needed
}

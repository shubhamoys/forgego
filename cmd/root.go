package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "forgego",
	Short: "ForgeGo is a CLI tool for scaffolding Go projects",
	Long: `ForgeGo is a developer-friendly CLI tool to bootstrap production-ready Go projects,
including RESTful APIs, CLI applications, and libraries, with idiomatic folder structures
and optional configurations like Docker and Git.`,
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

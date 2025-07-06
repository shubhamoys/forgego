package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new Go project",
	Long:  "Create a new Go project with interactive prompts or flags",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Scaffolding new Go project...")
	},
}

func init() {
	rootCmd.AddCommand(newCmd)

}

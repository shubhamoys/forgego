package cmd

import (
	"fmt"

	"github.com/shubhamoys/forgego/internal/prompts"
	"github.com/shubhamoys/forgego/internal/scaffold"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new Go project",
	Long:  "Create a new Go project with interactive prompts or flags",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := prompts.PromptNewProject()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		if err := scaffold.ScaffoldProject(config); err != nil {
			fmt.Printf("Error scaffolding project: %v\n", err)
			return
		}
		fmt.Printf("Successfully scaffolded project '%s' at %s\n", config.ProjectName, config.ProjectName)
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}

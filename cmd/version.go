package cmd

import (
	"fmt"

	"github.com/shubhamoys/forgego/internal/constants"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of ForgeGo",
	Long:  "Print the version number of ForgeGo",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("ForgeGo v%s\n", constants.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

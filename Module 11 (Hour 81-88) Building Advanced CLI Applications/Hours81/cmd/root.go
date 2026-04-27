package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mycli",
	Short: "My CLI Application",
	Long:  "A production-ready CLI built using Cobra",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to my CLI")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
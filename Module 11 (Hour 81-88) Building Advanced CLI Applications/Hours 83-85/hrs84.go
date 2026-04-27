package main

import (
	"fmt"
	"log"

	"github.com/manifoldco/promptui"//go get github.com/manifoldco/promptui
	"github.com/spf13/cobra"
)

func main() {

	rootCmd := &cobra.Command{
		Use: "app",
		Run: func(cmd *cobra.Command, args []string) {

			prompt := promptui.Prompt{
				Label: "Enter your name",
			}

			result, err := prompt.Run()
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("Hello", result)
		},
	}

	rootCmd.Execute()
}
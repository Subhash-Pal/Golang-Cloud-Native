package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var name string

func main() {

	// Load config
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Config not found")
	}

	rootCmd := &cobra.Command{
		Use: "app",
		Run: func(cmd *cobra.Command, args []string) {

			// Priority: Flag > Config
			if name == "" {
				name = viper.GetString("default_user")
			}

			fmt.Println("App:", viper.GetString("app_name"))
			fmt.Println("Hello", name)
		},
	}

	rootCmd.Flags().StringVarP(&name, "name", "n", "", "User name")

	rootCmd.Execute()
}
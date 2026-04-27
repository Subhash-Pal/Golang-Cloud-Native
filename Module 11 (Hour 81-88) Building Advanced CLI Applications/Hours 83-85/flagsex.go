package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)
//need advance cobra flags demo using golang 
var name string
func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Config not found")
	}

	var rootCmd = &cobra.Command{
		Use: "app",			
		Run: func(cmd *cobra.Command, args []string) {

				fmt.Println("App:", viper.GetString("app_name"))
				fmt.Println("Hello", name)
		},
	}	
	
	rootCmd.Flags().StringVarP(&name, "name", "n", "", "User name")
	rootCmd.Execute()
	}

	//how to run and use 
	//go run flagsex.go --name Shubh
		
		


		

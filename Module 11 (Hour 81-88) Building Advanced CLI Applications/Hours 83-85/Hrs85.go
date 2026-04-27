package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

// Struct for API response
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var userID int

func main() {

	var rootCmd = &cobra.Command{
		Use:   "app",
		Short: "CLI API Client",
	}

	// USERS COMMAND
	var usersCmd = &cobra.Command{
		Use:   "users",
		Short: "User operations",
	}

	// LIST USERS
	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List all users",
		Run: func(cmd *cobra.Command, args []string) {

			resp, err := http.Get("https://jsonplaceholder.typicode.com/users")
			if err != nil {
				log.Fatal("Request failed:", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				log.Fatalf("API error: %s", resp.Status)
			}

			var users []User
			err = json.NewDecoder(resp.Body).Decode(&users)
			if err != nil {
				log.Fatal("JSON decode error:", err)
			}

			// Pretty output
			fmt.Println("ID\tNAME\t\tEMAIL")
			for _, u := range users {
				fmt.Printf("%d\t%s\t%s\n", u.ID, u.Name, u.Email)
			}
		},
	}

	// GET USER BY ID
	var getCmd = &cobra.Command{
		Use:   "get",
		Short: "Get user by ID",
		Run: func(cmd *cobra.Command, args []string) {

			if userID == 0 {
				log.Fatal("Please provide --id")
			}

			url := fmt.Sprintf("https://jsonplaceholder.typicode.com/users/%d", userID)

			resp, err := http.Get(url)
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				log.Fatalf("User not found (status: %d)", resp.StatusCode)
			}

			var user User
			err = json.NewDecoder(resp.Body).Decode(&user)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("User Details:")
			fmt.Println("ID   :", user.ID)
			fmt.Println("Name :", user.Name)
			fmt.Println("Email:", user.Email)
		},
	}

	// FLAG
	getCmd.Flags().IntVarP(&userID, "id", "i", 0, "User ID")

	// COMMAND HIERARCHY
	rootCmd.AddCommand(usersCmd)
	usersCmd.AddCommand(listCmd)
	usersCmd.AddCommand(getCmd)

	// EXECUTE
	rootCmd.Execute()
}//go run hrs85.go users list
//Give me more commands

/*
Give me more option to  


*/


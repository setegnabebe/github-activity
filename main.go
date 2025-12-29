package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Event represents the parts of the GitHub event we want
type Event struct {
	Type string `json:"type"`
	Repo struct {
		Name string `json:"name"`
	} `json:"repo"`
}

func main() {
	var username string

	// Ask the user for the GitHub username
	fmt.Print("Enter GitHub username: ")
	fmt.Scanln(&username)

	if username == "" {
		fmt.Println("Username cannot be empty")
		return
	}

	// Build the GitHub API URL
	url := "https://api.github.com/users/" + username + "/events/public"

	// Fetch data from GitHub
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching data:", err)
		return
	}
	defer resp.Body.Close()

	// Decode the JSON into our Event struct
	var events []Event
	err = json.NewDecoder(resp.Body).Decode(&events)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	if len(events) == 0 {
		fmt.Println("No recent public activity found for this user.")
		return
	}

	// Print the results
	fmt.Println("\nRecent public activity:")
	for _, event := range events {
		fmt.Printf("- %-15s in %s\n", event.Type, event.Repo.Name)
	}
}

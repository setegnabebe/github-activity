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
	Payload struct {
		Commits []struct {
			Message string `json:"message"`
		} `json:"commits"`
		Action string `json:"action"`
	} `json:"payload"`
}

func main() {
	var username string

	// Ask the user to enter their GitHub username
	fmt.Print("Enter GitHub username: ")
	fmt.Scanln(&username)

	if username == "" {
		fmt.Println("Username cannot be empty")
		return
	}

	url := "https://api.github.com/users/" + username + "/events"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching data:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		fmt.Println("User not found.")
		return
	}
	if resp.StatusCode != 200 {
		fmt.Println("GitHub API returned:", resp.Status)
		return
	}

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

	fmt.Printf("\nRecent activity for %s:\n\n", username)
	for _, event := range events {
		switch event.Type {
		case "PushEvent":
			commitCount := len(event.Payload.Commits)
			fmt.Printf("- Pushed %d commit(s) to %s\n", commitCount, event.Repo.Name)
		case "IssuesEvent":
			action := event.Payload.Action
			if action == "" {
				action = "did something with"
			}
			fmt.Printf("- %s an issue in %s\n", capitalize(action), event.Repo.Name)
		case "PullRequestEvent":
			action := event.Payload.Action
			if action == "" {
				action = "did something with"
			}
			fmt.Printf("- %s a pull request in %s\n", capitalize(action), event.Repo.Name)
		case "WatchEvent":
			fmt.Printf("- Starred %s\n", event.Repo.Name)
		default:
			fmt.Printf("- %s in %s\n", event.Type, event.Repo.Name)
		}
	}
}

// capitalize the first letter of a string
func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return string(s[0]-32) + s[1:]
}

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {

		getActivity(scanner.Text())
	}

}

func getActivity(username string) {
	url := fmt.Sprintf("https://api.github.com/users/%s/events", username)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetching activity: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("API returned status %d", resp.StatusCode)
		return
	}

	var events []GitHubResponse
	if err = json.NewDecoder(resp.Body).Decode(&events); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		return
	}

	eventCounter(events)

}

func eventCounter(events []GitHubResponse) {
	if events == nil {
		fmt.Println("No events found")
		return
	}
	commits := make(map[string]int)
	issues := make(map[string]int)
	var starred []string

	for _, v := range events {
		repoName := v.Repo.Url

		switch v.Tipo {
		case "PushEvent":
			commits[repoName]++
		case "IssuesEvent":
			issues[repoName]++
		case "WatchEvent":
			starred = append(starred, repoName)
		}
	}

	for repo, count := range commits {
		fmt.Printf("- Pushed %d commits to %s\n", count, repo)

	}
	for repo, count := range issues {
		fmt.Printf("- Opened %d issue in %s\n", count, repo)
	}
	for _, repo := range starred {
		fmt.Printf("- Starred %s\n", repo)
	}
}

type GitHubResponse struct {
	Tipo string `json:"type"`
	Repo struct {
		Url string `json:"url"`
	}
}

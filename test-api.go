package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	//"os"
	"strings"
)

// Struct for testing a team roster endpoint
type TeamRosterResponse struct {
	Items []struct {
		Ref string `json:"$ref"`
	} `json:"items"`
}

// Struct for testing player overview
type PlayerDetails struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	FullName  string `json:"fullName"`

	Position struct {
		Name         string `json:"name"`
		Abbreviation string `json:"abbreviation"`
	} `json:"position"`

	College struct {
		Name string `json:"name"`
	} `json:"college"`

	Status     string `json:"status"`
	Experience int    `json:"experience"`
	Jersey     string `json:"jersey"`

	Team struct {
		ID          string `json:"id"`
		DisplayName string `json:"displayName"`
	} `json:"team"`

	Draft struct {
		Year      int `json:"year"`
		Round     int `json:"round"`
		Selection int `json:"selection"`
	} `json:"draft"`
}

func runTestAPI() {
	ctx := context.Background()

	// Choose your test case below:
	teamRosterTest(ctx)
	playerOverviewTest(ctx, "4038943") // ← Replace with any player ID to test
}

// ========== Test 1: Team Roster ==========

func teamRosterTest(ctx context.Context) {
	url := "https://sports.core.api.espn.com/v2/sports/football/leagues/nfl/seasons/2024/teams/22/athletes?limit=200"

	body, err := makeRequest(ctx, url)
	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}

	var roster TeamRosterResponse
	if err := json.Unmarshal(body, &roster); err != nil {
		log.Fatalf("JSON unmarshal failed: %v", err)
	}

	for _, item := range roster.Items {
		playerID := strings.Split(strings.Split(item.Ref, "/")[len(strings.Split(item.Ref, "/"))-1], "?")[0]
		fmt.Println("Player ID:", playerID)
	}
}

// ========== Test 2: Player Overview ==========

func playerOverviewTest(ctx context.Context, playerID string) {
	url := fmt.Sprintf("https://site.web.api.espn.com/apis/common/v3/sports/football/nfl/athletes/%s/overview", playerID)

	body, err := makeRequest(ctx, url)
	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}

	var result struct {
		Athlete PlayerDetails `json:"athlete"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		log.Fatalf("JSON unmarshal failed: %v", err)
	}

	player := result.Athlete
	fmt.Printf("✅ Parsed player %s %s (%s)\n", player.FirstName, player.LastName, player.ID)
	fmt.Printf("Position: %s\n", player.Position.Abbreviation)
	fmt.Printf("College: %s\n", player.College.Name)
	fmt.Printf("Team: %s (ID: %s)\n", player.Team.DisplayName, player.Team.ID)
	fmt.Printf("Jersey: %s | Status: %s | Draft: %d Rd %d Pick %d\n",
		player.Jersey, player.Status, player.Draft.Year, player.Draft.Round, player.Draft.Selection)
}

// ========== Reusable HTTP Request ==========

func makeRequest(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %d %s", resp.StatusCode, body)
	}

	return io.ReadAll(resp.Body)
}

// Temporary hook for testing
func init() {
	runTestAPI()
}

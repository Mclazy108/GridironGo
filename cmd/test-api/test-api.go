package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// Command line flags
var (
	teamID     = flag.String("team", "22", "Team ID to test (default: 22 - Cleveland Browns)")
	playerID   = flag.String("player", "", "Player ID to test (if empty, will use first player from team)")
	allTeams   = flag.Bool("all-teams", false, "Test getting all teams")
	testRoster = flag.Bool("roster", true, "Test team roster endpoint")
	testPlayer = flag.Bool("player-details", true, "Test player details endpoint")
	verbose    = flag.Bool("verbose", false, "Print verbose output including raw JSON")
	outputJSON = flag.Bool("output-json", false, "Save raw JSON responses to files")
)

// ====== API Response Types ======

// TeamListResponse for all teams endpoint
type TeamListResponse struct {
	Items []struct {
		Ref string `json:"$ref"`
		ID  string `json:"id"`
	} `json:"items"`
}

// TeamRosterResponse for team roster endpoint
type TeamRosterResponse struct {
	Items []struct {
		Ref string `json:"$ref"`
	} `json:"items"`
}

// PlayerDetailsResponse for player details endpoint
type PlayerDetailsResponse struct {
	Athlete struct {
		ID          string `json:"id"`
		FirstName   string `json:"firstName"`
		LastName    string `json:"lastName"`
		FullName    string `json:"fullName"`
		DisplayName string `json:"displayName"`
		ShortName   string `json:"shortName"`
		Weight      int    `json:"weight"`
		Height      int    `json:"height"`
		Jersey      string `json:"jersey"`
		Position    struct {
			Name         string `json:"name"`
			Abbreviation string `json:"abbreviation"`
		} `json:"position"`
		Active     bool `json:"active"`
		Experience int  `json:"experience"`
		College    struct {
			Name string `json:"name"`
		} `json:"college"`
		Status string `json:"status"`
		Team   struct {
			ID           string `json:"id"`
			DisplayName  string `json:"displayName"`
			Abbreviation string `json:"abbreviation"`
		} `json:"team"`
		Draft struct {
			Year      int `json:"year"`
			Round     int `json:"round"`
			Selection int `json:"selection"`
		} `json:"draft"`
	} `json:"athlete"`
	Headshot struct {
		Href string `json:"href"`
	} `json:"headshot"`
}

func main() {
	// Parse command line flags
	flag.Parse()

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Println("🏈 ESPN API Testing Tool for GridironGo 🏈")
	fmt.Println("==========================================")

	// Test fetching all teams if requested
	if *allTeams {
		fmt.Println("\n⚙️ TEST: Fetching all NFL teams...")
		teams, err := fetchAllTeams(ctx)
		if err != nil {
			log.Fatalf("❌ ERROR: Failed to fetch teams: %v", err)
		}

		fmt.Printf("✅ SUCCESS: Found %d teams\n", len(teams))

		// Print teams
		for i, team := range teams {
			fmt.Printf("  %d. Team ID: %s\n", i+1, team)
			if i >= 4 && !*verbose {
				fmt.Printf("  ... and %d more teams (use --verbose to see all)\n", len(teams)-5)
				break
			}
		}
	}

	// Test the team roster endpoint
	var firstPlayerID string
	if *testRoster {
		fmt.Printf("\n⚙️ TEST: Fetching roster for team ID %s...\n", *teamID)
		playerIDs, err := fetchTeamRoster(ctx, *teamID)
		if err != nil {
			log.Fatalf("❌ ERROR: Failed to fetch team roster: %v", err)
		}

		fmt.Printf("✅ SUCCESS: Found %d players on roster\n", len(playerIDs))

		// Print player IDs
		for i, id := range playerIDs {
			fmt.Printf("  %d. Player ID: %s\n", i+1, id)
			if i == 0 {
				firstPlayerID = id
			}
			if i >= 4 && !*verbose {
				fmt.Printf("  ... and %d more players (use --verbose to see all)\n", len(playerIDs)-5)
				break
			}
		}
	}

	// Test the player details endpoint
	if *testPlayer {
		// Determine which player ID to use
		testPlayerID := *playerID
		if testPlayerID == "" {
			if firstPlayerID != "" {
				testPlayerID = firstPlayerID
				fmt.Printf("\n⚙️ TEST: Using first player from roster (ID: %s)...\n", testPlayerID)
			} else {
				log.Fatal("❌ ERROR: No player ID provided and roster test disabled")
			}
		} else {
			fmt.Printf("\n⚙️ TEST: Fetching details for player ID %s...\n", testPlayerID)
		}

		playerDetails, rawJSON, err := fetchPlayerDetails(ctx, testPlayerID)
		if err != nil {
			log.Fatalf("❌ ERROR: Failed to fetch player details: %v", err)
		}

		// Print player details
		player := playerDetails.Athlete
		fmt.Println("✅ SUCCESS: Player details found")
		fmt.Printf("  Name: %s %s (%s)\n", player.FirstName, player.LastName, player.FullName)
		fmt.Printf("  Position: %s (%s)\n", player.Position.Name, player.Position.Abbreviation)
		fmt.Printf("  Team: %s (ID: %s)\n", player.Team.DisplayName, player.Team.ID)
		fmt.Printf("  Jersey: %s | Status: %s | Active: %v\n", player.Jersey, player.Status, player.Active)
		fmt.Printf("  Height: %d | Weight: %d | Experience: %d years\n", player.Height, player.Weight, player.Experience)
		fmt.Printf("  College: %s\n", player.College.Name)
		fmt.Printf("  Draft: Year %d, Round %d, Pick %d\n", player.Draft.Year, player.Draft.Round, player.Draft.Selection)
		fmt.Printf("  Headshot URL: %s\n", playerDetails.Headshot.Href)

		// Save raw JSON if requested
		if *outputJSON {
			filename := fmt.Sprintf("player_%s.json", testPlayerID)
			if err := saveToFile(filename, rawJSON); err != nil {
				fmt.Printf("⚠️ WARNING: Could not save JSON to file: %v\n", err)
			} else {
				fmt.Printf("📄 Saved raw JSON to %s\n", filename)
			}
		}

		// Print raw JSON if verbose mode is enabled
		if *verbose {
			fmt.Println("\n📊 Raw JSON response:")
			fmt.Println(string(rawJSON))
		}
	}

	fmt.Println("\n✅ All tests completed successfully")
}

// ====== API Functions ======

// fetchAllTeams fetches the list of all NFL teams
func fetchAllTeams(ctx context.Context) ([]string, error) {
	url := "https://sports.core.api.espn.com/v2/sports/football/leagues/nfl/teams"

	body, err := makeRequest(ctx, url)
	if err != nil {
		return nil, err
	}

	// Save raw JSON if requested
	if *outputJSON {
		if err := saveToFile("all_teams.json", body); err != nil {
			fmt.Printf("⚠️ WARNING: Could not save JSON to file: %v\n", err)
		}
	}

	var response TeamListResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("JSON unmarshal failed: %w", err)
	}

	var teamIDs []string
	for _, item := range response.Items {
		// Try to use the ID field first
		if item.ID != "" {
			teamIDs = append(teamIDs, item.ID)
			continue
		}

		// If ID is empty, try to extract from the reference URL
		if item.Ref != "" {
			parts := strings.Split(item.Ref, "/")
			if len(parts) > 0 {
				lastPart := parts[len(parts)-1]
				teamID := strings.Split(lastPart, "?")[0]
				if teamID != "" {
					teamIDs = append(teamIDs, teamID)
				}
			}
		}
	}

	return teamIDs, nil
}

// fetchTeamRoster fetches the roster for a specific team
func fetchTeamRoster(ctx context.Context, teamID string) ([]string, error) {
	url := fmt.Sprintf("https://sports.core.api.espn.com/v2/sports/football/leagues/nfl/seasons/2024/teams/%s/athletes?limit=200", teamID)

	body, err := makeRequest(ctx, url)
	if err != nil {
		return nil, err
	}

	// Save raw JSON if requested
	if *outputJSON {
		if err := saveToFile(fmt.Sprintf("team_%s_roster.json", teamID), body); err != nil {
			fmt.Printf("⚠️ WARNING: Could not save JSON to file: %v\n", err)
		}
	}

	var response TeamRosterResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("JSON unmarshal failed: %w", err)
	}

	var playerIDs []string
	for _, item := range response.Items {
		if item.Ref != "" {
			parts := strings.Split(item.Ref, "/")
			if len(parts) > 0 {
				lastPart := parts[len(parts)-1]
				playerID := strings.Split(lastPart, "?")[0]
				if playerID != "" {
					playerIDs = append(playerIDs, playerID)
				}
			}
		}
	}

	return playerIDs, nil
}

// fetchPlayerDetails fetches details for a specific player
func fetchPlayerDetails(ctx context.Context, playerID string) (*PlayerDetailsResponse, []byte, error) {
	url := fmt.Sprintf("https://site.web.api.espn.com/apis/common/v3/sports/football/nfl/athletes/%s/overview", playerID)

	body, err := makeRequest(ctx, url)
	if err != nil {
		return nil, nil, err
	}

	var response PlayerDetailsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, body, fmt.Errorf("JSON unmarshal failed: %w", err)
	}

	return &response, body, nil
}

// ====== Helper Functions ======

// makeRequest performs an HTTP request and returns the response body
func makeRequest(ctx context.Context, url string) ([]byte, error) {
	// Create a new request with context
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Set headers to make it look like a browser request
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
	req.Header.Set("Accept", "application/json")

	// Send the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return body, nil
}

// saveToFile saves content to a file
func saveToFile(filename string, content []byte) error {
	return os.WriteFile(filename, content, 0644)
}

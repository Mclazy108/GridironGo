package scraper

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Mclazy108/GridironGo/internals/data"
	"github.com/Mclazy108/GridironGo/internals/data/sqlc"
)

// PlayerScraper handles fetching and storing NFL player data
type PlayerScraper struct {
	DB *data.DB
}

// NewPlayerScraper creates a new scraper for NFL player data
func NewPlayerScraper(db *data.DB) *PlayerScraper {
	return &PlayerScraper{
		DB: db,
	}
}

// ESPNPlayerResponse represents the direct player data structure from ESPN API
type ESPNPlayerResponse struct {
	ID          string  `json:"id"`
	UID         string  `json:"uid"`
	GUID        string  `json:"guid"`
	FirstName   string  `json:"firstName"`
	LastName    string  `json:"lastName"`
	FullName    string  `json:"fullName"`
	DisplayName string  `json:"displayName"`
	ShortName   string  `json:"shortName"`
	Weight      float64 `json:"weight"` // Using float64 as API returns decimal values like 213.0
	Height      float64 `json:"height"` // Using float64 as API returns decimal values like 74.0
	Jersey      string  `json:"jersey"`
	Age         int     `json:"age,omitempty"` // Age is optional as it might not always be present
	DateOfBirth string  `json:"dateOfBirth,omitempty"`
	Position    struct {
		ID           string `json:"id"`
		Name         string `json:"name"`
		Abbreviation string `json:"abbreviation"`
		DisplayName  string `json:"displayName"`
	} `json:"position"`
	Team struct {
		ID           string `json:"id"`
		UID          string `json:"uid"`
		Slug         string `json:"slug"`
		Location     string `json:"location"`
		Name         string `json:"name"`
		Abbreviation string `json:"abbreviation"`
		DisplayName  string `json:"displayName"`
		ShortName    string `json:"shortName"`
		Color        string `json:"color"`
	} `json:"team"`
	College struct {
		ID           string `json:"id"`
		Name         string `json:"name"`
		Abbreviation string `json:"abbreviation"`
	} `json:"college"`
	Status struct {
		ID           string `json:"id"`
		Name         string `json:"name"`
		Type         string `json:"type"`
		Abbreviation string `json:"abbreviation"`
	} `json:"status"`
	Experience struct {
		Years int `json:"years"`
	} `json:"experience"`
	Active          bool   `json:"active"`
	HeadshotImgURL  string `json:"headshotImgUrl,omitempty"`  // Image URL may be at this level
	HeadshotImgHref string `json:"headshotImgHref,omitempty"` // Or might be here
	Draft           struct {
		Year      int `json:"year"`
		Round     int `json:"round"`
		Selection int `json:"selection"`
		Team      struct {
			ID           string `json:"id"`
			DisplayName  string `json:"displayName"`
			Abbreviation string `json:"abbreviation"`
		} `json:"team"`
	} `json:"draft"`
	Headshot struct {
		Href string `json:"href"` // Or might be in this nested structure
		Alt  string `json:"alt"`
	} `json:"headshot"`
}

// ScrapeNFLPlayers fetches and stores NFL player data
func (s *PlayerScraper) ScrapeNFLPlayers(ctx context.Context) error {
	log.Println("Starting NFL players scraping process...")

	// First, get all teams from the database
	teams, err := s.DB.Queries.GetAllNFLTeams(ctx)
	if err != nil || len(teams) == 0 {
		return fmt.Errorf("failed to fetch NFL teams from database: %w", err)
	}

	log.Printf("Found %d teams. Will fetch player data from team rosters", len(teams))

	// Track processed players to avoid duplicates
	processedPlayers := make(map[string]bool)
	totalPlayers := 0
	successfulPlayers := 0

	// Process each team
	for _, team := range teams {
		// Check if context was cancelled
		select {
		case <-ctx.Done():
			log.Println("Scraping cancelled by user")
			return ctx.Err()
		default:
			// Continue processing
		}

		log.Printf("Processing team: %s (%s)", team.DisplayName, team.TeamID)

		// Fetch team roster
		playerIDs, err := s.fetchTeamRoster(ctx, team.TeamID)
		if err != nil {
			log.Printf("Error fetching roster for team %s: %v", team.DisplayName, err)
			continue
		}

		log.Printf("Found %d players on %s roster", len(playerIDs), team.DisplayName)

		// Process each player in the roster
		for _, playerID := range playerIDs {
			// Skip if we've already processed this player
			if processedPlayers[playerID] {
				continue
			}

			// Mark this player as processed
			processedPlayers[playerID] = true
			totalPlayers++

			// Check if context was cancelled
			select {
			case <-ctx.Done():
				log.Println("Scraping cancelled by user")
				return ctx.Err()
			default:
				// Continue processing
			}

			// Process the player details
			if err := s.processPlayer(ctx, playerID); err != nil {
				log.Printf("Error processing player ID %s: %v", playerID, err)
			} else {
				successfulPlayers++
				if successfulPlayers%10 == 0 {
					log.Printf("Successfully processed %d/%d players", successfulPlayers, totalPlayers)
				}
			}

			// Sleep briefly to avoid rate limiting
			time.Sleep(100 * time.Millisecond)
		}
	}

	log.Printf("Processed %d unique players from team rosters (%d successful, %d failed)",
		totalPlayers, successfulPlayers, totalPlayers-successfulPlayers)
	log.Println("NFL players scraping completed")
	return nil
}

// fetchTeamRoster fetches the roster for a specific team
func (s *PlayerScraper) fetchTeamRoster(ctx context.Context, teamID string) ([]string, error) {
	// Construct the API URL for team roster
	url := fmt.Sprintf("https://sports.core.api.espn.com/v2/sports/football/leagues/nfl/seasons/2024/teams/%s/athletes?limit=200", teamID)

	// Create a new request with context
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers to make it look like a browser request
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
	req.Header.Set("Accept", "application/json")

	// Send HTTP GET request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return nil, fmt.Errorf("request cancelled: %w", ctx.Err())
		}
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned non-OK status: %d. Response: %s", resp.StatusCode, string(body))
	}

	// Parse the JSON response
	var rosterResponse struct {
		Items []struct {
			Ref string `json:"$ref"`
		} `json:"items"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&rosterResponse); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	// Extract player IDs from the response
	var playerIDs []string
	for _, item := range rosterResponse.Items {
		if item.Ref != "" {
			// The $ref URL is typically in the format ".../athletes/{playerID}?..."
			// We need to extract just the playerID part
			parts := strings.Split(item.Ref, "/")
			if len(parts) > 0 {
				// Get the last part of the URL (which contains the player ID)
				lastPart := parts[len(parts)-1]
				// If there's a query string, remove it
				playerID := strings.Split(lastPart, "?")[0]
				if playerID != "" {
					playerIDs = append(playerIDs, playerID)
				}
			}
		}
	}

	return playerIDs, nil
}

// fetchPlayerDetails fetches detailed information for a specific player
func (s *PlayerScraper) fetchPlayerDetails(ctx context.Context, playerID string) (*ESPNPlayerResponse, error) {
	// Construct the API URL for player details
	// Using the direct athlete endpoint from ESPN API
	url := fmt.Sprintf("https://sports.core.api.espn.com/v2/sports/football/leagues/nfl/athletes/%s", playerID)

	// Create a new request with context
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers to make it look like a browser request
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
	req.Header.Set("Accept", "application/json")

	// Send HTTP GET request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return nil, fmt.Errorf("request cancelled: %w", ctx.Err())
		}
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned non-OK status: %d. Response: %s", resp.StatusCode, string(body))
	}

	// For debugging, optionally log the raw response
	// body, _ := io.ReadAll(resp.Body)
	// log.Printf("Raw player response: %s", string(body))
	// resp.Body = io.NopCloser(bytes.NewBuffer(body))

	// Parse the JSON response
	var playerResponse ESPNPlayerResponse
	if err := json.NewDecoder(resp.Body).Decode(&playerResponse); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	return &playerResponse, nil
}

// processPlayer handles a single player (fetch details and update database)
func (s *PlayerScraper) processPlayer(ctx context.Context, playerID string) error {
	// Fetch player details from the ESPN API
	playerResponse, err := s.fetchPlayerDetails(ctx, playerID)
	if err != nil {
		return fmt.Errorf("error fetching player details: %w", err)
	}

	// Determine which URL to use for the player's image
	imageURL := playerResponse.HeadshotImgURL
	if imageURL == "" {
		imageURL = playerResponse.HeadshotImgHref
	}
	if imageURL == "" && playerResponse.Headshot.Href != "" {
		imageURL = playerResponse.Headshot.Href
	}

	// Skip players without position data
	if playerResponse.Position.Abbreviation == "" {
		log.Printf("Skipping player %s - no position data", playerResponse.FullName)
		return nil
	}

	// Get status value (use Name field from the Status struct)
	statusValue := ""
	if playerResponse.Status.Name != "" {
		statusValue = playerResponse.Status.Name
	}

	// Get experience value (either from Years field or the whole int)
	experienceValue := 0
	if playerResponse.Experience.Years > 0 {
		experienceValue = playerResponse.Experience.Years
	}

	// Check if player already exists in database
	_, err = s.DB.Queries.GetNFLPlayer(ctx, playerID)

	// Set up database parameters
	playerParams := sqlc.CreateNFLPlayerParams{
		PlayerID:   playerID,
		FirstName:  playerResponse.FirstName,
		LastName:   playerResponse.LastName,
		FullName:   playerResponse.FullName,
		Position:   playerResponse.Position.Abbreviation,
		TeamID:     sql.NullString{String: playerResponse.Team.ID, Valid: playerResponse.Team.ID != ""},
		Jersey:     sql.NullString{String: playerResponse.Jersey, Valid: playerResponse.Jersey != ""},
		Height:     sql.NullInt64{Int64: int64(playerResponse.Height), Valid: playerResponse.Height > 0},
		Weight:     sql.NullInt64{Int64: int64(playerResponse.Weight), Valid: playerResponse.Weight > 0},
		Active:     playerResponse.Active,
		College:    sql.NullString{String: playerResponse.College.Name, Valid: playerResponse.College.Name != ""},
		Experience: sql.NullInt64{Int64: int64(experienceValue), Valid: experienceValue >= 0},
		DraftYear:  sql.NullInt64{Int64: int64(playerResponse.Draft.Year), Valid: playerResponse.Draft.Year > 0},
		DraftRound: sql.NullInt64{Int64: int64(playerResponse.Draft.Round), Valid: playerResponse.Draft.Round > 0},
		DraftPick:  sql.NullInt64{Int64: int64(playerResponse.Draft.Selection), Valid: playerResponse.Draft.Selection > 0},
		Status:     sql.NullString{String: statusValue, Valid: statusValue != ""},
		ImageUrl:   sql.NullString{String: imageURL, Valid: imageURL != ""},
	}

	if err == nil {
		// Player exists, update
		updateParams := sqlc.UpdateNFLPlayerParams{
			PlayerID:   playerParams.PlayerID,
			FirstName:  playerParams.FirstName,
			LastName:   playerParams.LastName,
			FullName:   playerParams.FullName,
			Position:   playerParams.Position,
			TeamID:     playerParams.TeamID,
			Jersey:     playerParams.Jersey,
			Height:     playerParams.Height,
			Weight:     playerParams.Weight,
			Active:     playerParams.Active,
			College:    playerParams.College,
			Experience: playerParams.Experience,
			DraftYear:  playerParams.DraftYear,
			DraftRound: playerParams.DraftRound,
			DraftPick:  playerParams.DraftPick,
			Status:     playerParams.Status,
			ImageUrl:   playerParams.ImageUrl,
		}

		if err := s.DB.Queries.UpdateNFLPlayer(ctx, updateParams); err != nil {
			return fmt.Errorf("error updating player in database: %w", err)
		}
	} else {
		// Player doesn't exist, insert
		if err := s.DB.Queries.CreateNFLPlayer(ctx, playerParams); err != nil {
			return fmt.Errorf("error inserting player into database: %w", err)
		}
	}

	return nil
}

// fetchPlayerImageURL attempts to fetch a player's image URL from the ESPN API
func (s *PlayerScraper) fetchPlayerImageURL(ctx context.Context, playerID string) (string, error) {
	// Alternative endpoint that may contain image URLs
	url := fmt.Sprintf("https://site.web.api.espn.com/apis/common/v3/sports/football/nfl/athletes/%s/overview", playerID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var response struct {
		Athlete struct {
			HeadShot struct {
				Href string `json:"href"`
			} `json:"headshot"`
		} `json:"athlete"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}

	return response.Athlete.HeadShot.Href, nil
}

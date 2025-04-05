package scraper

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/Mclazy108/GridironGo/internals/data"
	"github.com/Mclazy108/GridironGo/internals/data/sqlc"
	"io"
	"log"
	"net/http"
	"time"
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

// PlayerListResponse represents the top-level response from the ESPN players API
type PlayerListResponse struct {
	Items []struct {
		ID  string `json:"id"`
		Ref string `json:"$ref"`
	} `json:"items"`
}

// PlayerResponse represents the detailed player information from the ESPN API
type PlayerResponse struct {
	Player PlayerDetails `json:"athlete"`
}

// PlayerDetails represents an NFL player's details
type PlayerDetails struct {
	ID          string    `json:"id"`
	UID         string    `json:"uid"`
	GUID        string    `json:"guid"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	FullName    string    `json:"fullName"`
	DisplayName string    `json:"displayName"`
	ShortName   string    `json:"shortName"`
	Weight      int       `json:"weight"`
	Height      int       `json:"height"`
	Age         int       `json:"age"`
	DateOfBirth string    `json:"dateOfBirth"`
	Jersey      string    `json:"jersey"`
	Position    string    `json:"position"`
	Positions   []string  `json:"positions"`
	Active      bool      `json:"active"`
	DebutYear   int       `json:"debutYear"`
	ImageURL    string    `json:"headshot"`
	Status      string    `json:"status"`
	Experience  int       `json:"experience"`
	College     string    `json:"college"`
	Team        Team      `json:"team"`
	DraftTeam   Team      `json:"draftTeam,omitempty"`
	DraftInfo   DraftInfo `json:"draft,omitempty"`
}

// Team represents a player's team information
type Team struct {
	ID           string `json:"id"`
	DisplayName  string `json:"displayName"`
	Abbreviation string `json:"abbreviation"`
	Name         string `json:"name"`
	Nickname     string `json:"nickname"`
	Location     string `json:"location"`
}

// DraftInfo represents a player's draft information
type DraftInfo struct {
	Year      int `json:"year"`
	Round     int `json:"round"`
	Selection int `json:"selection"`
}

// ScrapeNFLPlayers fetches and stores NFL player data
func (s *PlayerScraper) ScrapeNFLPlayers(ctx context.Context) error {
	log.Println("Starting NFL players scraping process...")
	log.Println("Press Ctrl+C to cancel the scraping process gracefully")

	// First, try to get team rosters as an alternative approach
	teams, err := s.DB.Queries.GetAllNFLTeams(ctx)
	if err == nil && len(teams) > 0 {
		log.Printf("Found %d teams. Fetching player data from team rosters", len(teams))

		// A map to track unique player IDs we've already processed
		processedPlayers := make(map[string]bool)

		// For each team, fetch their roster
		for _, team := range teams {
			// Fetch team roster
			playerIDs, err := s.fetchTeamRoster(ctx, team.TeamID)
			if err != nil {
				log.Printf("Warning: Error fetching roster for team %s (%s): %v", team.DisplayName, team.TeamID, err)
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

				// Process player details
				if err := s.processPlayer(ctx, playerID); err != nil {
					log.Printf("Error processing player ID %s: %v", playerID, err)
				}
			}
		}

		log.Printf("Processed %d unique players from team rosters", len(processedPlayers))
		log.Println("NFL players scraping completed successfully")
		return nil
	}

	// Fall back to the original approach if team roster method fails
	log.Println("Falling back to direct player list fetch method")

	// Start with getting the list of player IDs
	playerIDs, err := s.fetchPlayerIDs(ctx)
	if err != nil {
		return fmt.Errorf("error fetching player list: %w", err)
	}

	log.Printf("Found %d NFL players to process", len(playerIDs))

	// Limit to process only a portion of players if needed for testing
	// (Can be removed for production)
	/*
		limit := 50
		if len(playerIDs) > limit {
			log.Printf("Limiting to first %d players for testing", limit)
			playerIDs = playerIDs[:limit]
		}
	*/

	// Process each player to get detailed information
	for i, playerID := range playerIDs {
		// Check if context was cancelled
		select {
		case <-ctx.Done():
			log.Println("Scraping cancelled by user")
			return nil
		default:
			// Continue processing
		}

		// Only log every 10th player to avoid log spam
		if i%10 == 0 {
			log.Printf("Processing player %d of %d (ID: %s)...", i+1, len(playerIDs), playerID)
		}

		// Fetch detailed player information
		playerDetails, err := s.fetchPlayerDetails(ctx, playerID)
		if err != nil {
			log.Printf("Error fetching details for player ID %s: %v", playerID, err)
			continue
		}

		// Skip players without position data
		if playerDetails.Position == "" {
			log.Printf("Skipping player %s - no position data", playerDetails.FullName)
			continue
		}

		// Check if the player already exists in the database
		existingPlayer, err := s.DB.Queries.GetNFLPlayer(ctx, playerID)
		if err == nil {
			// Player exists, check if we need to update it
			if i%50 == 0 {
				log.Printf("Player with ID %s already exists: %s", playerID, existingPlayer.FullName)
			}

			teamIDChanged := (existingPlayer.TeamID.Valid && existingPlayer.TeamID.String != playerDetails.Team.ID) ||
				(!existingPlayer.TeamID.Valid && playerDetails.Team.ID != "")

			jerseyChanged := (existingPlayer.Jersey.Valid && existingPlayer.Jersey.String != playerDetails.Jersey) ||
				(!existingPlayer.Jersey.Valid && playerDetails.Jersey != "")

			imageURLChanged := (existingPlayer.ImageUrl.Valid && existingPlayer.ImageUrl.String != playerDetails.ImageURL) ||
				(!existingPlayer.ImageUrl.Valid && playerDetails.ImageURL != "")

			statusChanged := (existingPlayer.Status.Valid && existingPlayer.Status.String != playerDetails.Status) ||
				(!existingPlayer.Status.Valid && playerDetails.Status != "")

			if existingPlayer.FirstName != playerDetails.FirstName ||
				existingPlayer.LastName != playerDetails.LastName ||
				existingPlayer.Position != playerDetails.Position ||
				existingPlayer.Active != playerDetails.Active ||
				teamIDChanged ||
				jerseyChanged ||
				imageURLChanged ||
				statusChanged {

				// Update the player
				updateParams := sqlc.UpdateNFLPlayerParams{
					PlayerID:   playerID,
					FirstName:  playerDetails.FirstName,
					LastName:   playerDetails.LastName,
					FullName:   playerDetails.FullName,
					Position:   playerDetails.Position,
					TeamID:     sql.NullString{String: playerDetails.Team.ID, Valid: playerDetails.Team.ID != ""},
					Jersey:     sql.NullString{String: playerDetails.Jersey, Valid: playerDetails.Jersey != ""},
					Height:     sql.NullInt64{Int64: int64(playerDetails.Height), Valid: playerDetails.Height > 0},
					Weight:     sql.NullInt64{Int64: int64(playerDetails.Weight), Valid: playerDetails.Weight > 0},
					Active:     playerDetails.Active,
					College:    sql.NullString{String: playerDetails.College, Valid: playerDetails.College != ""},
					Experience: sql.NullInt64{Int64: int64(playerDetails.Experience), Valid: playerDetails.Experience >= 0},
					DraftYear:  sql.NullInt64{Int64: int64(playerDetails.DraftInfo.Year), Valid: playerDetails.DraftInfo.Year > 0},
					DraftRound: sql.NullInt64{Int64: int64(playerDetails.DraftInfo.Round), Valid: playerDetails.DraftInfo.Round > 0},
					DraftPick:  sql.NullInt64{Int64: int64(playerDetails.DraftInfo.Selection), Valid: playerDetails.DraftInfo.Selection > 0},
					Status:     sql.NullString{String: playerDetails.Status, Valid: playerDetails.Status != ""},
					ImageUrl:   sql.NullString{String: playerDetails.ImageURL, Valid: playerDetails.ImageURL != ""},
				}

				err = s.DB.Queries.UpdateNFLPlayer(ctx, updateParams)
				if err != nil {
					log.Printf("Error updating player with ID %s: %v", playerID, err)
				} else if i%50 == 0 {
					log.Printf("Updated player: %s (ID: %s)", playerDetails.FullName, playerID)
				}
			}
			continue
		}

		// Insert new player into database
		params := sqlc.CreateNFLPlayerParams{
			PlayerID:   playerID,
			FirstName:  playerDetails.FirstName,
			LastName:   playerDetails.LastName,
			FullName:   playerDetails.FullName,
			Position:   playerDetails.Position,
			TeamID:     sql.NullString{String: playerDetails.Team.ID, Valid: playerDetails.Team.ID != ""},
			Jersey:     sql.NullString{String: playerDetails.Jersey, Valid: playerDetails.Jersey != ""},
			Height:     sql.NullInt64{Int64: int64(playerDetails.Height), Valid: playerDetails.Height > 0},
			Weight:     sql.NullInt64{Int64: int64(playerDetails.Weight), Valid: playerDetails.Weight > 0},
			Active:     playerDetails.Active,
			College:    sql.NullString{String: playerDetails.College, Valid: playerDetails.College != ""},
			Experience: sql.NullInt64{Int64: int64(playerDetails.Experience), Valid: playerDetails.Experience >= 0},
			DraftYear:  sql.NullInt64{Int64: int64(playerDetails.DraftInfo.Year), Valid: playerDetails.DraftInfo.Year > 0},
			DraftRound: sql.NullInt64{Int64: int64(playerDetails.DraftInfo.Round), Valid: playerDetails.DraftInfo.Round > 0},
			DraftPick:  sql.NullInt64{Int64: int64(playerDetails.DraftInfo.Selection), Valid: playerDetails.DraftInfo.Selection > 0},
			Status:     sql.NullString{String: playerDetails.Status, Valid: playerDetails.Status != ""},
			ImageUrl:   sql.NullString{String: playerDetails.ImageURL, Valid: playerDetails.ImageURL != ""},
		}

		err = s.DB.Queries.CreateNFLPlayer(ctx, params)
		if err != nil {
			log.Printf("Error inserting player with ID %s: %v", playerID, err)
			continue
		}

		if i%50 == 0 {
			log.Printf("Inserted player: %s (ID: %s)", playerDetails.FullName, playerID)
		}

		// Sleep to avoid rate limiting, but make it interruptible
		select {
		case <-ctx.Done():
			log.Println("Scraping cancelled by user during rate limit sleep")
			return nil
		case <-time.After(100 * time.Millisecond):
			// Continue with the next player
		}
	}

	log.Println("NFL players scraping completed successfully")
	return nil
}

// fetchPlayerIDs fetches the list of active NFL player IDs from the ESPN API
func (s *PlayerScraper) fetchPlayerIDs(ctx context.Context) ([]string, error) {
	// Construct the API URL - using the Player Metadata endpoint from the CSV
	url := "https://sports.core.api.espn.com/v2/sports/football/leagues/nfl/athletes?limit=1000&active=true"

	// Create a new request with context
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers to make the request more like a browser
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	// Debug log
	log.Printf("Sending request to %s", url)

	// Send HTTP GET request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// Check if the error was due to context cancellation
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

	// Debug log
	log.Printf("Received response from %s", url)

	// Parse the JSON response
	var playerListResponse PlayerListResponse
	err = json.NewDecoder(resp.Body).Decode(&playerListResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	// Extract player IDs
	var playerIDs []string
	for _, item := range playerListResponse.Items {
		if item.ID != "" {
			playerIDs = append(playerIDs, item.ID)
		}
	}

	log.Printf("Found %d player IDs in response", len(playerIDs))

	return playerIDs, nil
}

// fetchPlayerDetails fetches detailed information for a specific player
func (s *PlayerScraper) fetchPlayerDetails(ctx context.Context, playerID string) (*PlayerDetails, error) {
	// Construct the API URL for detailed player information
	url := fmt.Sprintf("https://site.web.api.espn.com/apis/common/v3/sports/football/nfl/athletes/%s/overview", playerID)

	// Create a new request with context
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers to make the request more like a browser
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	// Send HTTP GET request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// Check if the error was due to context cancellation
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
	var response struct {
		Player struct {
			ID          string `json:"id"`
			UID         string `json:"uid"`
			GUID        string `json:"guid"`
			FirstName   string `json:"firstName"`
			LastName    string `json:"lastName"`
			FullName    string `json:"fullName"`
			DisplayName string `json:"displayName"`
			ShortName   string `json:"shortName"`
			Weight      int    `json:"weight"`
			Height      int    `json:"height"`
			Age         int    `json:"age"`
			Jersey      string `json:"jersey"`
			Position    string `json:"position"`
			Active      bool   `json:"active"`
			College     string `json:"college"`
			Experience  int    `json:"experience"`
			Status      string `json:"status"`
			Team        struct {
				ID           string `json:"id"`
				DisplayName  string `json:"displayName"`
				Abbreviation string `json:"abbreviation"`
				Name         string `json:"name"`
				Nickname     string `json:"nickname"`
				Location     string `json:"location"`
			} `json:"team"`
			DraftInfo struct {
				Year      int `json:"year"`
				Round     int `json:"round"`
				Selection int `json:"selection"`
			} `json:"draft"`
			HeadShot struct {
				Href string `json:"href"`
			} `json:"headshot"`
		} `json:"athlete"`
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	// Convert the response to our PlayerDetails struct
	details := &PlayerDetails{
		ID:          response.Player.ID,
		UID:         response.Player.UID,
		GUID:        response.Player.GUID,
		FirstName:   response.Player.FirstName,
		LastName:    response.Player.LastName,
		FullName:    response.Player.FullName,
		DisplayName: response.Player.DisplayName,
		ShortName:   response.Player.ShortName,
		Weight:      response.Player.Weight,
		Height:      response.Player.Height,
		Age:         response.Player.Age,
		Jersey:      response.Player.Jersey,
		Position:    response.Player.Position,
		Active:      response.Player.Active,
		College:     response.Player.College,
		Experience:  response.Player.Experience,
		Status:      response.Player.Status,
		Team: Team{
			ID:           response.Player.Team.ID,
			DisplayName:  response.Player.Team.DisplayName,
			Abbreviation: response.Player.Team.Abbreviation,
			Name:         response.Player.Team.Name,
			Nickname:     response.Player.Team.Nickname,
			Location:     response.Player.Team.Location,
		},
		DraftInfo: DraftInfo{
			Year:      response.Player.DraftInfo.Year,
			Round:     response.Player.DraftInfo.Round,
			Selection: response.Player.DraftInfo.Selection,
		},
		ImageURL: response.Player.HeadShot.Href,
	}

	return details, nil
}

// fetchTeamRoster fetches a team's roster using the API endpoint from the CSV file
func (s *PlayerScraper) fetchTeamRoster(ctx context.Context, teamID string) ([]string, error) {
	// Construct the API URL for team roster
	// Based on "Team Rosters" endpoint in api.csv:
	// https://sports.core.api.espn.com/v2/sports/football/leagues/nfl/seasons/{YEAR}/teams/{TEAM_ID}/athletes?limit=200
	url := fmt.Sprintf("https://sports.core.api.espn.com/v2/sports/football/leagues/nfl/seasons/2024/teams/%s/athletes?limit=200", teamID)

	// Create a new request with context
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	log.Printf("Fetching roster for team ID %s", teamID)

	// Send HTTP GET request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// Check if the error was due to context cancellation
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
	var rosterResponse PlayerListResponse
	err = json.NewDecoder(resp.Body).Decode(&rosterResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	// Extract player IDs
	var playerIDs []string
	for _, item := range rosterResponse.Items {
		if item.ID != "" {
			playerIDs = append(playerIDs, item.ID)
		}
	}

	return playerIDs, nil
}

// processPlayer fetches and stores data for a single player
func (s *PlayerScraper) processPlayer(ctx context.Context, playerID string) error {
	// Fetch player details
	playerDetails, err := s.fetchPlayerDetails(ctx, playerID)
	if err != nil {
		return fmt.Errorf("error fetching details: %w", err)
	}

	// Skip players without position data
	if playerDetails.Position == "" {
		log.Printf("Skipping player %s - no position data", playerDetails.FullName)
		return nil
	}

	// Check if player exists in database
	existingPlayer, err := s.DB.Queries.GetNFLPlayer(ctx, playerID)
	if err == nil {
		// Player exists, check if we need to update
		teamIDChanged := (existingPlayer.TeamID.Valid && existingPlayer.TeamID.String != playerDetails.Team.ID) ||
			(!existingPlayer.TeamID.Valid && playerDetails.Team.ID != "")

		jerseyChanged := (existingPlayer.Jersey.Valid && existingPlayer.Jersey.String != playerDetails.Jersey) ||
			(!existingPlayer.Jersey.Valid && playerDetails.Jersey != "")

		imageURLChanged := (existingPlayer.ImageUrl.Valid && existingPlayer.ImageUrl.String != playerDetails.ImageURL) ||
			(!existingPlayer.ImageUrl.Valid && playerDetails.ImageURL != "")

		statusChanged := (existingPlayer.Status.Valid && existingPlayer.Status.String != playerDetails.Status) ||
			(!existingPlayer.Status.Valid && playerDetails.Status != "")

		if existingPlayer.FirstName != playerDetails.FirstName ||
			existingPlayer.LastName != playerDetails.LastName ||
			existingPlayer.Position != playerDetails.Position ||
			existingPlayer.Active != playerDetails.Active ||
			teamIDChanged ||
			jerseyChanged ||
			imageURLChanged ||
			statusChanged {

			// Update the player
			updateParams := sqlc.UpdateNFLPlayerParams{
				PlayerID:   playerID,
				FirstName:  playerDetails.FirstName,
				LastName:   playerDetails.LastName,
				FullName:   playerDetails.FullName,
				Position:   playerDetails.Position,
				TeamID:     sql.NullString{String: playerDetails.Team.ID, Valid: playerDetails.Team.ID != ""},
				Jersey:     sql.NullString{String: playerDetails.Jersey, Valid: playerDetails.Jersey != ""},
				Height:     sql.NullInt64{Int64: int64(playerDetails.Height), Valid: playerDetails.Height > 0},
				Weight:     sql.NullInt64{Int64: int64(playerDetails.Weight), Valid: playerDetails.Weight > 0},
				Active:     playerDetails.Active,
				College:    sql.NullString{String: playerDetails.College, Valid: playerDetails.College != ""},
				Experience: sql.NullInt64{Int64: int64(playerDetails.Experience), Valid: playerDetails.Experience >= 0},
				DraftYear:  sql.NullInt64{Int64: int64(playerDetails.DraftInfo.Year), Valid: playerDetails.DraftInfo.Year > 0},
				DraftRound: sql.NullInt64{Int64: int64(playerDetails.DraftInfo.Round), Valid: playerDetails.DraftInfo.Round > 0},
				DraftPick:  sql.NullInt64{Int64: int64(playerDetails.DraftInfo.Selection), Valid: playerDetails.DraftInfo.Selection > 0},
				Status:     sql.NullString{String: playerDetails.Status, Valid: playerDetails.Status != ""},
				ImageUrl:   sql.NullString{String: playerDetails.ImageURL, Valid: playerDetails.ImageURL != ""},
			}

			return s.DB.Queries.UpdateNFLPlayer(ctx, updateParams)
		}

		return nil // No update needed
	}

	// Insert new player
	params := sqlc.CreateNFLPlayerParams{
		PlayerID:   playerID,
		FirstName:  playerDetails.FirstName,
		LastName:   playerDetails.LastName,
		FullName:   playerDetails.FullName,
		Position:   playerDetails.Position,
		TeamID:     sql.NullString{String: playerDetails.Team.ID, Valid: playerDetails.Team.ID != ""},
		Jersey:     sql.NullString{String: playerDetails.Jersey, Valid: playerDetails.Jersey != ""},
		Height:     sql.NullInt64{Int64: int64(playerDetails.Height), Valid: playerDetails.Height > 0},
		Weight:     sql.NullInt64{Int64: int64(playerDetails.Weight), Valid: playerDetails.Weight > 0},
		Active:     playerDetails.Active,
		College:    sql.NullString{String: playerDetails.College, Valid: playerDetails.College != ""},
		Experience: sql.NullInt64{Int64: int64(playerDetails.Experience), Valid: playerDetails.Experience >= 0},
		DraftYear:  sql.NullInt64{Int64: int64(playerDetails.DraftInfo.Year), Valid: playerDetails.DraftInfo.Year > 0},
		DraftRound: sql.NullInt64{Int64: int64(playerDetails.DraftInfo.Round), Valid: playerDetails.DraftInfo.Round > 0},
		DraftPick:  sql.NullInt64{Int64: int64(playerDetails.DraftInfo.Selection), Valid: playerDetails.DraftInfo.Selection > 0},
		Status:     sql.NullString{String: playerDetails.Status, Valid: playerDetails.Status != ""},
		ImageUrl:   sql.NullString{String: playerDetails.ImageURL, Valid: playerDetails.ImageURL != ""},
	}

	return s.DB.Queries.CreateNFLPlayer(ctx, params)
}

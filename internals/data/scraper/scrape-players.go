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

type PlayerScraper struct {
	DB *data.DB
}

func NewPlayerScraper(db *data.DB) *PlayerScraper {
	return &PlayerScraper{DB: db}
}

// TeamRosterResponse represents the response format from ESPN team roster API
type TeamRosterResponse struct {
	Items []struct {
		Ref string `json:"$ref"`
	} `json:"items"`
}

// PlayerDetails represents the response format from ESPN player overview API
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

// ScrapeNFLPlayers fetches and stores NFL player data
func (s *PlayerScraper) ScrapeNFLPlayers(ctx context.Context) error {
	log.Println("Starting NFL players scraping process...")

	// First, get all teams from the database
	teams, err := s.DB.Queries.GetAllNFLTeams(ctx)
	if err != nil || len(teams) == 0 {
		return fmt.Errorf("failed to fetch NFL teams from database: %w", err)
	}

	log.Printf("Found %d teams. Fetching player data from team rosters", len(teams))

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

// fetchTeamRoster fetches the list of player IDs for a specific team
func (s *PlayerScraper) fetchTeamRoster(ctx context.Context, teamID string) ([]string, error) {
	// Construct the API URL for team roster
	url := fmt.Sprintf("https://sports.core.api.espn.com/v2/sports/football/leagues/nfl/seasons/2024/teams/%s/athletes?limit=200", teamID)

	// Create a new request with context
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
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
	var rosterResponse TeamRosterResponse
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
func (s *PlayerScraper) fetchPlayerDetails(ctx context.Context, playerID string) (*PlayerDetailsResponse, error) {
	// Construct the API URL for player details
	url := fmt.Sprintf("https://site.web.api.espn.com/apis/common/v3/sports/football/nfl/athletes/%s/overview", playerID)

	// Create a new request with context
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
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
	var playerResponse PlayerDetailsResponse
	if err := json.NewDecoder(resp.Body).Decode(&playerResponse); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	return &playerResponse, nil
}

// processPlayer processes a single player (fetch details and update/insert in database)
func (s *PlayerScraper) processPlayer(ctx context.Context, playerID string) error {
	// Fetch player details
	playerResponse, err := s.fetchPlayerDetails(ctx, playerID)
	if err != nil {
		return fmt.Errorf("error fetching player details: %w", err)
	}

	// Extract player details from response
	player := playerResponse.Athlete

	// Skip players without position data
	if player.Position.Abbreviation == "" {
		log.Printf("Skipping player %s - no position data", player.FullName)
		return nil
	}

	// Check if player already exists in the database
	_, err = s.DB.Queries.GetNFLPlayer(ctx, playerID)
	if err == nil {
		// Player exists, update it
		updateParams := sqlc.UpdateNFLPlayerParams{
			PlayerID:   playerID,
			FirstName:  player.FirstName,
			LastName:   player.LastName,
			FullName:   player.FullName,
			Position:   player.Position.Abbreviation,
			TeamID:     sql.NullString{String: player.Team.ID, Valid: player.Team.ID != ""},
			Jersey:     sql.NullString{String: player.Jersey, Valid: player.Jersey != ""},
			Height:     sql.NullInt64{Int64: int64(player.Height), Valid: player.Height > 0},
			Weight:     sql.NullInt64{Int64: int64(player.Weight), Valid: player.Weight > 0},
			Active:     player.Active,
			College:    sql.NullString{String: player.College.Name, Valid: player.College.Name != ""},
			Experience: sql.NullInt64{Int64: int64(player.Experience), Valid: player.Experience >= 0},
			DraftYear:  sql.NullInt64{Int64: int64(player.Draft.Year), Valid: player.Draft.Year > 0},
			DraftRound: sql.NullInt64{Int64: int64(player.Draft.Round), Valid: player.Draft.Round > 0},
			DraftPick:  sql.NullInt64{Int64: int64(player.Draft.Selection), Valid: player.Draft.Selection > 0},
			Status:     sql.NullString{String: player.Status, Valid: player.Status != ""},
			ImageUrl:   sql.NullString{String: playerResponse.Headshot.Href, Valid: playerResponse.Headshot.Href != ""},
		}

		if err := s.DB.Queries.UpdateNFLPlayer(ctx, updateParams); err != nil {
			return fmt.Errorf("error updating player in database: %w", err)
		}
	} else {
		// Player doesn't exist, insert it
		insertParams := sqlc.CreateNFLPlayerParams{
			PlayerID:   playerID,
			FirstName:  player.FirstName,
			LastName:   player.LastName,
			FullName:   player.FullName,
			Position:   player.Position.Abbreviation,
			TeamID:     sql.NullString{String: player.Team.ID, Valid: player.Team.ID != ""},
			Jersey:     sql.NullString{String: player.Jersey, Valid: player.Jersey != ""},
			Height:     sql.NullInt64{Int64: int64(player.Height), Valid: player.Height > 0},
			Weight:     sql.NullInt64{Int64: int64(player.Weight), Valid: player.Weight > 0},
			Active:     player.Active,
			College:    sql.NullString{String: player.College.Name, Valid: player.College.Name != ""},
			Experience: sql.NullInt64{Int64: int64(player.Experience), Valid: player.Experience >= 0},
			DraftYear:  sql.NullInt64{Int64: int64(player.Draft.Year), Valid: player.Draft.Year > 0},
			DraftRound: sql.NullInt64{Int64: int64(player.Draft.Round), Valid: player.Draft.Round > 0},
			DraftPick:  sql.NullInt64{Int64: int64(player.Draft.Selection), Valid: player.Draft.Selection > 0},
			Status:     sql.NullString{String: player.Status, Valid: player.Status != ""},
			ImageUrl:   sql.NullString{String: playerResponse.Headshot.Href, Valid: playerResponse.Headshot.Href != ""},
		}

		if err := s.DB.Queries.CreateNFLPlayer(ctx, insertParams); err != nil {
			return fmt.Errorf("error inserting player into database: %w", err)
		}
	}

	return nil
}


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
	"sync"
	"sync/atomic"

	"github.com/Mclazy108/GridironGo/internals/data"
	"github.com/Mclazy108/GridironGo/internals/data/sqlc"
	"golang.org/x/time/rate"
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
	Linked bool `json:"linked"`
}

// PlayerData holds a player's information and their team ID
type PlayerData struct {
	PlayerID   string
	TeamID     string
	PlayerInfo *ESPNPlayerResponse
}

// ScrapeNFLPlayers fetches and stores NFL player data with team-based batching
func (s *PlayerScraper) ScrapeNFLPlayers(ctx context.Context) error {
	log.Println("Starting NFL players scraping process...")

	// First, get all teams from the database
	teams, err := s.DB.Queries.GetAllNFLTeams(ctx)
	if err != nil || len(teams) == 0 {
		return fmt.Errorf("failed to fetch NFL teams from database: %w", err)
	}

	log.Printf("Found %d teams. Will fetch player data from team rosters", len(teams))

	// Track total players for stats
	var totalPlayers int32 = 0
	var processedTeams int32 = 0
	var failedPlayers int32 = 0

	// Create a rate limiter to avoid overwhelming the API
	// Limit to 10 requests per second (adjust as needed)
	limiter := rate.NewLimiter(10, 1)

	// Create a wait group to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Create a channel to process teams
	teamChan := make(chan *sqlc.NflTeam, len(teams))

	// Number of worker goroutines to process teams
	numWorkers := 25
	log.Printf("Starting %d team worker goroutines", numWorkers)

	// Launch worker goroutines
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for team := range teamChan {
				log.Printf("Worker %d: Processing team %s", workerID, team.DisplayName)

				// Process the team roster as a batch
				playersProcessed, err := s.processTeamRoster(ctx, *team, limiter)

				if err != nil {
					log.Printf("Worker %d: Error processing team %s: %v",
						workerID, team.DisplayName, err)
				} else {
					log.Printf("Worker %d: Successfully processed %d players for team %s",
						workerID, playersProcessed, team.DisplayName)

					// Increment processed teams counter
					numProcessed := atomic.AddInt32(&processedTeams, 1)
					totalPlayers := atomic.AddInt32(&totalPlayers, int32(playersProcessed))

					log.Printf("Progress: %d/%d teams processed, %d total players",
						numProcessed, len(teams), totalPlayers)
				}
			}
			log.Printf("Worker %d finished", workerID)
		}(i)
	}

	// Send teams to workers
	for i := range teams {
		teamChan <- teams[i]
	}

	// Close the team channel when done
	close(teamChan)

	// Wait for all team workers to finish
	log.Println("Waiting for all team workers to finish...")
	wg.Wait()

	log.Printf("Processed %d teams with %d total players (%d failed)",
		len(teams), totalPlayers, failedPlayers)
	log.Println("NFL players scraping completed")
	return nil
}

// processTeamRoster fetches and processes an entire team's roster in a single transaction
func (s *PlayerScraper) processTeamRoster(ctx context.Context, team sqlc.NflTeam, limiter *rate.Limiter) (int, error) {
	teamID := team.TeamID
	teamName := team.DisplayName

	// Fetch the team's roster
	if err := limiter.Wait(ctx); err != nil {
		return 0, fmt.Errorf("rate limiter error: %w", err)
	}

	playerIDs, err := s.fetchTeamRoster(ctx, teamID)
	if err != nil {
		return 0, fmt.Errorf("error fetching roster for team %s: %w", teamName, err)
	}

	log.Printf("Found %d players on %s roster, fetching player details", len(playerIDs), teamName)

	// Collect player data
	playerDataList := make([]PlayerData, 0, len(playerIDs))

	// Fetch details for all players in the roster
	for _, playerID := range playerIDs {
		// Rate limit API calls
		if err := limiter.Wait(ctx); err != nil {
			log.Printf("Rate limiter error while fetching player %s: %v", playerID, err)
			continue
		}

		// Fetch player details
		playerResponse, err := s.fetchPlayerDetails(ctx, playerID)
		if err != nil {
			log.Printf("Error fetching details for player ID %s: %v", playerID, err)
			continue
		}

		// Skip players without position data
		if playerResponse.Position.Abbreviation == "" {
			log.Printf("Skipping player %s - no position data", playerResponse.FullName)
			continue
		}

		// Add to player data list
		playerDataList = append(playerDataList, PlayerData{
			PlayerID:   playerID,
			TeamID:     teamID,
			PlayerInfo: playerResponse,
		})
	}

	// If no players to process, return early
	if len(playerDataList) == 0 {
		log.Printf("No valid players found for team %s", teamName)
		return 0, nil
	}

	log.Printf("Processed details for %d players on %s roster, saving to database...",
		len(playerDataList), teamName)

	// Start a transaction for the batch insert/update
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Create a queries object with the transaction
	q := s.DB.Queries.WithTx(tx)

	// Counter for successful inserts/updates
	successCount := 0

	// Process each player in the transaction
	for _, playerData := range playerDataList {
		playerID := playerData.PlayerID
		teamID := playerData.TeamID
		playerResponse := playerData.PlayerInfo

		// Determine which URL to use for the player's image
		imageURL := playerResponse.HeadshotImgURL
		if imageURL == "" {
			imageURL = playerResponse.HeadshotImgHref
		}
		if imageURL == "" && playerResponse.Headshot.Href != "" {
			imageURL = playerResponse.Headshot.Href
		}

		// Get status value
		statusValue := ""
		if playerResponse.Status.Name != "" {
			statusValue = playerResponse.Status.Name
		}

		// Get experience value
		experienceValue := 0
		if playerResponse.Experience.Years > 0 {
			experienceValue = playerResponse.Experience.Years
		}

		// Create player parameters
		playerParams := sqlc.CreateNFLPlayerParams{
			PlayerID:   playerID,
			FirstName:  playerResponse.FirstName,
			LastName:   playerResponse.LastName,
			FullName:   playerResponse.FullName,
			Position:   playerResponse.Position.Abbreviation,
			TeamID:     sql.NullString{String: teamID, Valid: teamID != ""},
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

		// Check if player exists
		_, err = q.GetNFLPlayer(ctx, playerID)
		if err == nil {
			// Update existing player
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

			if err := q.UpdateNFLPlayer(ctx, updateParams); err != nil {
				log.Printf("Error updating player %s: %v", playerID, err)
			} else {
				successCount++
			}
		} else {
			// Insert new player
			if err := q.CreateNFLPlayer(ctx, playerParams); err != nil {
				log.Printf("Error inserting player %s: %v", playerID, err)
			} else {
				successCount++
			}
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return successCount, fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("Successfully saved %d/%d players for team %s to database",
		successCount, len(playerDataList), teamName)

	return successCount, nil
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

	// Parse the JSON response
	var playerResponse ESPNPlayerResponse
	if err := json.NewDecoder(resp.Body).Decode(&playerResponse); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	return &playerResponse, nil
}

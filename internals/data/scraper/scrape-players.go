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
	SeasonYear int
	PlayerInfo *ESPNPlayerResponse
}

// ScrapeNFLPlayers fetches and stores NFL player data with team-based batching
func (s *PlayerScraper) ScrapeNFLPlayers(ctx context.Context, seasons []int) error {
	log.Println("Starting NFL players scraping process...")

	// First, get all teams from the database
	teams, err := s.DB.Queries.GetAllNFLTeams(ctx)
	if err != nil || len(teams) == 0 {
		return fmt.Errorf("failed to fetch NFL teams from database: %w", err)
	}

	log.Printf("Found %d teams. Will fetch player data from team rosters for seasons: %v", len(teams), seasons)

	// Track total players for stats
	var totalPlayers int32 = 0
	var processedTeams int32 = 0
	var failedPlayers int32 = 0

	// Create a rate limiter to avoid overwhelming the API
	limiter := rate.NewLimiter(2500, 1)

	// Create a wait group to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Create a channel for team-season combinations
	type TeamSeason struct {
		Team       *sqlc.NflTeam
		SeasonYear int
	}
	teamSeasonChan := make(chan TeamSeason, len(teams)*len(seasons))

	// Number of worker goroutines to process teams
	numWorkers := 32
	log.Printf("Starting %d team worker goroutines", numWorkers)

	// Launch worker goroutines
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for teamSeason := range teamSeasonChan {
				team := teamSeason.Team
				seasonYear := teamSeason.SeasonYear

				log.Printf("Worker %d: Processing team %s for season %d", workerID, team.DisplayName, seasonYear)

				// Process the team roster as a batch
				playersProcessed, err := s.processTeamRoster(ctx, *team, seasonYear, limiter)

				if err != nil {
					log.Printf("Worker %d: Error processing team %s for season %d: %v",
						workerID, team.DisplayName, seasonYear, err)
				} else {
					log.Printf("Worker %d: Successfully processed %d players for team %s in season %d",
						workerID, playersProcessed, team.DisplayName, seasonYear)

					// Increment processed teams counter
					numProcessed := atomic.AddInt32(&processedTeams, 1)
					totalPlayersCount := atomic.AddInt32(&totalPlayers, int32(playersProcessed))

					log.Printf("Progress: %d/%d team-seasons processed, %d total players",
						numProcessed, len(teams)*len(seasons), totalPlayersCount)
				}
			}
			log.Printf("Worker %d finished", workerID)
		}(i)
	}

	// Send team-season combinations to workers
	for _, season := range seasons {
		for i := range teams {
			teamSeasonChan <- TeamSeason{
				Team:       teams[i],
				SeasonYear: season,
			}
		}
	}

	// Close the team-season channel when done
	close(teamSeasonChan)

	// Wait for all team workers to finish
	log.Println("Waiting for all team workers to finish...")
	wg.Wait()

	log.Printf("Processed %d team-seasons with %d total players (%d failed)",
		len(teams)*len(seasons), totalPlayers, failedPlayers)
	log.Println("NFL players scraping completed")
	return nil
}

// insertNFLPlayersBulk inserts or updates players in the nfl_players table
func insertNFLPlayersBulk(ctx context.Context, tx *sql.Tx, players []PlayerData) error {
	if len(players) == 0 {
		return nil
	}

	// First prepare to gather unique players for the base table
	playerMap := make(map[string]PlayerData)

	// Use the most recent data for each player
	for _, p := range players {
		existing, ok := playerMap[p.PlayerID]
		if !ok || p.SeasonYear > existing.SeasonYear {
			playerMap[p.PlayerID] = p
		}
	}

	// Build query for nfl_players table
	query := `INSERT INTO nfl_players (
		player_id, first_name, last_name, full_name, position, team_id,
		jersey, height, weight, active, college, experience,
		draft_year, draft_round, draft_pick, status, image_url
	) VALUES `

	// Collect value placeholders
	valueStrings := make([]string, 0, len(playerMap))
	valueArgs := make([]interface{}, 0, len(playerMap)*17)

	for _, p := range playerMap {
		playerInfo := p.PlayerInfo
		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")

		// Get image URL
		imageURL := playerInfo.HeadshotImgURL
		if imageURL == "" {
			imageURL = playerInfo.HeadshotImgHref
		}
		if imageURL == "" && playerInfo.Headshot.Href != "" {
			imageURL = playerInfo.Headshot.Href
		}

		// Get experience value
		experienceValue := 0
		if playerInfo.Experience.Years > 0 {
			experienceValue = playerInfo.Experience.Years
		}

		// Get status value
		statusValue := ""
		if playerInfo.Status.Name != "" {
			statusValue = playerInfo.Status.Name
		}

		valueArgs = append(valueArgs,
			p.PlayerID,
			playerInfo.FirstName,
			playerInfo.LastName,
			playerInfo.FullName,
			playerInfo.Position.Abbreviation,
			p.TeamID, // Current team
			playerInfo.Jersey,
			int(playerInfo.Height),
			int(playerInfo.Weight),
			playerInfo.Active,
			playerInfo.College.Name,
			experienceValue,
			playerInfo.Draft.Year,
			playerInfo.Draft.Round,
			playerInfo.Draft.Selection,
			statusValue,
			imageURL,
		)
	}

	query += strings.Join(valueStrings, ",") +
		` ON CONFLICT(player_id) DO UPDATE SET
			first_name = excluded.first_name,
			last_name = excluded.last_name,
			full_name = excluded.full_name,
			position = excluded.position,
			team_id = excluded.team_id,
			jersey = excluded.jersey,
			height = excluded.height,
			weight = excluded.weight,
			active = excluded.active,
			college = excluded.college,
			experience = excluded.experience,
			draft_year = excluded.draft_year,
			draft_round = excluded.draft_round,
			draft_pick = excluded.draft_pick,
			status = excluded.status,
			image_url = excluded.image_url`

	// Execute the query for nfl_players
	_, err := tx.ExecContext(ctx, query, valueArgs...)
	if err != nil {
		return fmt.Errorf("error inserting players: %w", err)
	}

	return nil
}

// insertPlayerSeasonsBulk inserts or updates player seasons data
func insertPlayerSeasonsBulk(ctx context.Context, tx *sql.Tx, players []PlayerData) error {
	if len(players) == 0 {
		return nil
	}

	// Build query for nfl_player_seasons table
	query := `INSERT INTO nfl_player_seasons (
		player_id, season_year, team_id, jersey, active, experience, status
	) VALUES `

	// Collect value placeholders
	valueStrings := make([]string, 0, len(players))
	valueArgs := make([]interface{}, 0, len(players)*7)

	for _, p := range players {
		playerInfo := p.PlayerInfo
		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?)")

		// Get experience value
		experienceValue := 0
		if playerInfo.Experience.Years > 0 {
			experienceValue = playerInfo.Experience.Years
		}

		// Get status value
		statusValue := ""
		if playerInfo.Status.Name != "" {
			statusValue = playerInfo.Status.Name
		}

		valueArgs = append(valueArgs,
			p.PlayerID,
			p.SeasonYear,
			p.TeamID,
			playerInfo.Jersey,
			playerInfo.Active,
			experienceValue,
			statusValue,
		)
	}

	query += strings.Join(valueStrings, ",") +
		` ON CONFLICT(player_id, season_year) DO UPDATE SET
			team_id = excluded.team_id,
			jersey = excluded.jersey,
			active = excluded.active,
			experience = excluded.experience,
			status = excluded.status`

	// Execute the query for nfl_player_seasons
	_, err := tx.ExecContext(ctx, query, valueArgs...)
	if err != nil {
		return fmt.Errorf("error inserting player seasons: %w", err)
	}

	return nil
}

// processTeamRoster fetches and processes an entire team's roster for a specific season
func (s *PlayerScraper) processTeamRoster(ctx context.Context, team sqlc.NflTeam, seasonYear int, limiter *rate.Limiter) (int, error) {
	teamID := team.TeamID
	teamName := team.DisplayName

	if err := limiter.Wait(ctx); err != nil {
		return 0, fmt.Errorf("rate limiter error: %w", err)
	}

	playerIDs, err := s.fetchTeamRoster(ctx, teamID, seasonYear, limiter)
	if err != nil {
		return 0, fmt.Errorf("error fetching roster for team %s in season %d: %w", teamName, seasonYear, err)
	}

	log.Printf("Found %d players on %s roster for season %d, fetching player details", len(playerIDs), teamName, seasonYear)

	playerDataList := make([]PlayerData, 0, len(playerIDs))

	for _, playerID := range playerIDs {
		if err := limiter.Wait(ctx); err != nil {
			log.Printf("Rate limiter error while fetching player %s: %v", playerID, err)
			continue
		}

		playerResponse, err := s.fetchPlayerDetails(ctx, playerID, limiter)
		if err != nil {
			log.Printf("Error fetching details for player ID %s: %v", playerID, err)
			continue
		}

		if playerResponse.Position.Abbreviation == "" {
			log.Printf("Skipping player %s - no position data", playerResponse.FullName)
			continue
		}

		playerDataList = append(playerDataList, PlayerData{
			PlayerID:   playerID,
			TeamID:     teamID,
			SeasonYear: seasonYear,
			PlayerInfo: playerResponse,
		})
	}

	if len(playerDataList) == 0 {
		log.Printf("No valid players found for team %s in season %d", teamName, seasonYear)
		return 0, nil
	}

	log.Printf("Processed details for %d players on %s roster for season %d, saving to database...",
		len(playerDataList), teamName, seasonYear)

	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}

	// First insert/update the base player records
	err = insertNFLPlayersBulk(ctx, tx, playerDataList)
	if err != nil {
		_ = tx.Rollback()
		return 0, fmt.Errorf("error inserting base player records: %w", err)
	}

	// Then insert/update the player seasons records
	err = insertPlayerSeasonsBulk(ctx, tx, playerDataList)
	if err != nil {
		_ = tx.Rollback()
		return 0, fmt.Errorf("error inserting player seasons: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("Successfully saved %d players for team %s in season %d to database",
		len(playerDataList), teamName, seasonYear)

	return len(playerDataList), nil
}

// fetchTeamRoster fetches the roster for a specific team and season
func (s *PlayerScraper) fetchTeamRoster(ctx context.Context, teamID string, seasonYear int, limiter *rate.Limiter) ([]string, error) {
	// Construct the API URL for team roster with the specific season
	url := fmt.Sprintf("https://sports.core.api.espn.com/v2/sports/football/leagues/nfl/seasons/%d/teams/%s/athletes?limit=200",
		seasonYear, teamID)

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
func (s *PlayerScraper) fetchPlayerDetails(ctx context.Context, playerID string, limiter *rate.Limiter) (*ESPNPlayerResponse, error) {
	// Construct the API URL for player details
	// Using the direct athlete endpoint from ESPN API
	url := fmt.Sprintf("https://sports.core.api.espn.com/v2/sports/football/leagues/nfl/athletes/%s", playerID)

	// Wait for rate limiter
	if err := limiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limiter error: %w", err)
	}

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


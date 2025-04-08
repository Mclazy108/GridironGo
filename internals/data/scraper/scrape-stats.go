package scraper

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	//"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/Mclazy108/GridironGo/internals/data"
	"github.com/Mclazy108/GridironGo/internals/data/sqlc"
	"golang.org/x/time/rate"
)

// StatScraper handles fetching and storing NFL game statistics
type StatScraper struct {
	DB *data.DB
}

// NewStatScraper creates a new scraper for NFL game statistics
func NewStatScraper(db *data.DB) *StatScraper {
	return &StatScraper{
		DB: db,
	}
}

// ESPNGameSummaryResponse represents the JSON structure from the ESPN summary API
type ESPNGameSummaryResponse struct {
	Header struct {
		ID string `json:"id"`
	} `json:"header"`
	Boxscore struct {
		Teams []struct {
			Team struct {
				ID           string `json:"id"`
				Abbreviation string `json:"abbreviation"`
			} `json:"team"`
			Statistics []struct {
				Name         string   `json:"name"`
				DisplayName  string   `json:"displayName"`
				Keys         []string `json:"keys"`
				Labels       []string `json:"labels"`
				Descriptions []string `json:"descriptions"`
			} `json:"statistics"`
		} `json:"teams"`
		Players []struct {
			Team struct {
				ID string `json:"id"`
			} `json:"team"`
			Statistics []struct {
				Name         string   `json:"name"`
				Keys         []string `json:"keys"`
				Labels       []string `json:"labels"`
				Descriptions []string `json:"descriptions"`
				Athletes     []struct {
					Athlete struct {
						ID          string `json:"id"`
						DisplayName string `json:"displayName"`
					} `json:"athlete"`
					Stats []string `json:"stats"`
				} `json:"athletes"`
			} `json:"statistics"`
		} `json:"players"`
	} `json:"boxscore"`
	Leaders []struct {
		Team struct {
			ID string `json:"id"`
		} `json:"team"`
		Leaders []struct {
			Name        string `json:"name"`
			DisplayName string `json:"displayName"`
			Leaders     []struct {
				Athlete struct {
					ID   string `json:"id"`
					Name string `json:"name"`
				} `json:"athlete"`
				Value float64 `json:"value"`
			} `json:"leaders"`
		} `json:"leaders"`
	} `json:"leaders"`
}

// StatData represents processed statistics ready to be stored
type StatData struct {
	GameID    int64
	PlayerID  string
	TeamID    string
	Category  string
	StatType  string
	StatValue float64
}

// ScrapeNFLGameStats fetches and stores NFL game statistics
func (s *StatScraper) ScrapeNFLGameStats(ctx context.Context, seasons []int) error {
	log.Println("Starting NFL game statistics scraping process...")

	// Fetch all games for the specified seasons
	var games []*sqlc.NflGame
	for _, season := range seasons {
		seasonGames, err := s.DB.Queries.GetGamesBySeason(ctx, int64(season))
		if err != nil {
			return fmt.Errorf("failed to fetch NFL games for season %d: %w", season, err)
		}
		games = append(games, seasonGames...)
	}

	if len(games) == 0 {
		return fmt.Errorf("no games found for specified seasons: %v", seasons)
	}

	log.Printf("Found %d games across specified seasons. Will fetch game statistics", len(games))

	// Debug the API response for the first game
	/*
		if len(games) > 0 {
			err := s.debugAPIResponse(ctx, games[0].EventID)
			if err != nil {
				log.Printf("Warning: Failed to debug API response: %v", err)
			}
		}
	*/

	// Track processing statistics
	var totalStats int32 = 0
	var processedGames int32 = 0
	var failedGames int32 = 0

	limiter := rate.NewLimiter(500, 1)

	// Create a wait group to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Create a channel to process games
	gameChan := make(chan *sqlc.NflGame, len(games))

	// Number of worker goroutines to process games
	numWorkers := 15
	log.Printf("Starting %d game worker goroutines", numWorkers)

	// Launch worker goroutines
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for game := range gameChan {
				select {
				case <-ctx.Done():
					log.Printf("Worker %d: Stopping due to context cancellation", workerID)
					return
				default:
					log.Printf("Worker %d: Processing game %s (ID: %d)", workerID, game.Name, game.EventID)

					// Process the game statistics
					statsProcessed, err := s.processGameStats(ctx, *game, limiter)

					if err != nil {
						log.Printf("Worker %d: Error processing game %s: %v",
							workerID, game.Name, err)
						atomic.AddInt32(&failedGames, 1)
					} else {
						log.Printf("Worker %d: Successfully processed %d stats for game %s",
							workerID, statsProcessed, game.Name)

						// Increment processed games counter
						atomic.AddInt32(&processedGames, 1)
						atomic.AddInt32(&totalStats, int32(statsProcessed))

						// Log progress periodically
						if atomic.LoadInt32(&processedGames)%10 == 0 {
							log.Printf("Progress: %d/%d games processed, %d total stats",
								atomic.LoadInt32(&processedGames), len(games), atomic.LoadInt32(&totalStats))
						}
					}
				}
			}
			log.Printf("Worker %d finished", workerID)
		}(i)
	}

	// Send games to workers
	for i := range games {
		select {
		case <-ctx.Done():
			log.Println("Stopping game distribution due to context cancellation")
			break
		default:
			gameChan <- games[i]
		}
	}

	// Close the game channel when done
	close(gameChan)

	// Wait for all game workers to finish
	log.Println("Waiting for all game workers to finish...")
	wg.Wait()

	log.Printf("Processed %d/%d games with %d total stats (%d failed)",
		atomic.LoadInt32(&processedGames), len(games), atomic.LoadInt32(&totalStats), atomic.LoadInt32(&failedGames))
	log.Println("NFL game statistics scraping completed")
	return nil
}

// debugAPIResponse fetches and saves a complete API response for debugging
/*
func (s *StatScraper) debugAPIResponse(ctx context.Context, gameID int64) error {
	// Construct the API URL for game summary
	url := fmt.Sprintf("https://site.api.espn.com/apis/site/v2/sports/football/nfl/summary?event=%d", gameID)

	// Create a new request with context
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers to make it look like a browser request
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
	req.Header.Set("Accept", "application/json")

	// Send HTTP GET request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return fmt.Errorf("request cancelled: %w", ctx.Err())
		}
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API returned non-OK status: %d. Response: %s", resp.StatusCode, string(body))
	}

	// Read the body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Save the response to a file for inspection
	err = os.WriteFile(fmt.Sprintf("game_%d_response.json", gameID), body, 0644)
	if err != nil {
		return fmt.Errorf("failed to save response to file: %w", err)
	}

	log.Printf("Saved API response for game %d to game_%d_response.json", gameID, gameID)
	return nil
}
*/
// processGameStats fetches and processes statistics for a single game
func (s *StatScraper) processGameStats(ctx context.Context, game sqlc.NflGame, limiter *rate.Limiter) (int, error) {
	// Wait for rate limiter
	if err := limiter.Wait(ctx); err != nil {
		return 0, fmt.Errorf("rate limiter error: %w", err)
	}

	// Fetch game summary from ESPN API
	gameSummary, err := s.fetchGameSummary(ctx, game.EventID)
	if err != nil {
		return 0, fmt.Errorf("error fetching summary for game %d: %w", game.EventID, err)
	}

	// Process the game summary to extract stats
	stats, err := s.extractGameStats(ctx, gameSummary, game.EventID)
	if err != nil {
		return 0, fmt.Errorf("error extracting stats for game %d: %w", game.EventID, err)
	}

	if len(stats) == 0 {
		log.Printf("No stats found for game %d (%s)", game.EventID, game.Name)
		return 0, nil
	}

	log.Printf("Extracted %d stats for game %d, saving to database...", len(stats), game.EventID)

	// Save stats to database in a transaction
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Insert each stat
	insertCount := 0
	for _, stat := range stats {

		if stat.StatValue == 0 {
			continue
		}
		// First check if this stat already exists
		var statID int64
		err := tx.QueryRowContext(ctx, `
			SELECT stat_id FROM nfl_stats 
			WHERE game_id = ? AND player_id = ? AND team_id = ? AND category = ? AND stat_type = ?
		`, stat.GameID, stat.PlayerID, stat.TeamID, stat.Category, stat.StatType).Scan(&statID)

		if err != nil && err != sql.ErrNoRows {
			// Query error
			_ = tx.Rollback()
			return 0, fmt.Errorf("error checking if stat exists: %w", err)
		}

		if err == sql.ErrNoRows {
			// Stat doesn't exist, insert it
			_, err = tx.ExecContext(ctx, `
				INSERT INTO nfl_stats (game_id, player_id, team_id, category, stat_type, stat_value)
				VALUES (?, ?, ?, ?, ?, ?)
			`, stat.GameID, stat.PlayerID, stat.TeamID, stat.Category, stat.StatType, stat.StatValue)
		} else {
			// Stat exists, update it
			_, err = tx.ExecContext(ctx, `
				UPDATE nfl_stats
				SET stat_value = ?
				WHERE stat_id = ?
			`, stat.StatValue, statID)
		}

		if err != nil {
			// Insert/update error
			_ = tx.Rollback()
			return 0, fmt.Errorf("error inserting/updating stat: %w", err)
		}

		insertCount++
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return insertCount, nil
}

// fetchGameSummary retrieves game summary data from the ESPN API
func (s *StatScraper) fetchGameSummary(ctx context.Context, gameID int64) (*ESPNGameSummaryResponse, error) {
	// Construct the API URL for game summary
	url := fmt.Sprintf("https://site.api.espn.com/apis/site/v2/sports/football/nfl/summary?event=%d", gameID)

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
	var summaryResponse ESPNGameSummaryResponse
	if err := json.NewDecoder(resp.Body).Decode(&summaryResponse); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	return &summaryResponse, nil
}

// extractGameStats processes a game summary to extract all player statistics
func (s *StatScraper) extractGameStats(ctx context.Context, summary *ESPNGameSummaryResponse, gameID int64) ([]StatData, error) {
	log.Printf("Debug: Game ID: %d", gameID)

	var stats []StatData

	// Verify game exists first
	var gameExists bool
	err := s.DB.DB.QueryRowContext(ctx, "SELECT 1 FROM nfl_games WHERE event_id = ?", gameID).Scan(&gameExists)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("error checking if game exists: %w", err)
	}
	if !gameExists {
		return nil, fmt.Errorf("game with ID %d does not exist in database", gameID)
	}

	// Get all valid player IDs and team IDs from database to verify them
	rows, err := s.DB.DB.QueryContext(ctx, "SELECT player_id FROM nfl_players")
	if err != nil {
		return nil, fmt.Errorf("error fetching player IDs: %w", err)
	}
	defer rows.Close()

	validPlayerIDs := make(map[string]bool)
	for rows.Next() {
		var playerID string
		if err := rows.Scan(&playerID); err != nil {
			return nil, fmt.Errorf("error scanning player ID: %w", err)
		}
		validPlayerIDs[playerID] = true
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating player IDs: %w", err)
	}

	// Get all valid team IDs
	rows, err = s.DB.DB.QueryContext(ctx, "SELECT team_id FROM nfl_teams")
	if err != nil {
		return nil, fmt.Errorf("error fetching team IDs: %w", err)
	}
	defer rows.Close()

	validTeamIDs := make(map[string]bool)
	for rows.Next() {
		var teamID string
		if err := rows.Scan(&teamID); err != nil {
			return nil, fmt.Errorf("error scanning team ID: %w", err)
		}
		validTeamIDs[teamID] = true
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating team IDs: %w", err)
	}

	// Process boxscore player statistics
	for _, teamPlayers := range summary.Boxscore.Players {
		teamID := teamPlayers.Team.ID
		if !validTeamIDs[teamID] {
			continue
		}

		for _, statCategory := range teamPlayers.Statistics {
			category := statCategory.Name
			for keyIndex, key := range statCategory.Keys {
				for _, athlete := range statCategory.Athletes {
					playerID := athlete.Athlete.ID
					if !validPlayerIDs[playerID] {
						continue
					}

					if keyIndex >= len(athlete.Stats) {
						log.Printf("Skipping stat: index %d out of bounds for player %s (ID: %s) — category: %s, key: %s",
							keyIndex, athlete.Athlete.DisplayName, playerID, category, key)
						continue
					}

					rawStat := athlete.Stats[keyIndex]
					if rawStat == "" || rawStat == "--" {
						continue
					}

					statValue, err := parseStatValue(rawStat)
					if err != nil {
						log.Printf("Could not parse stat value '%s' for player %s (ID: %s), category: %s, key: %s — error: %v",
							rawStat, athlete.Athlete.DisplayName, playerID, category, key, err)
						continue
					}

					stats = append(stats, StatData{
						GameID:    gameID,
						PlayerID:  playerID,
						TeamID:    teamID,
						Category:  category,
						StatType:  key,
						StatValue: statValue,
					})
				}
			}
		}
	}

	// Process leaders data with the same validation
	/*
		for _, teamLeader := range summary.Leaders {
			teamID := teamLeader.Team.ID

			// Skip if team ID is not valid
			if !validTeamIDs[teamID] {
				log.Printf("Warning: skipping leader stats for team ID %s (not found in database)", teamID)
				continue
			}

			for _, leaderCategory := range teamLeader.Leaders {
				category := leaderCategory.Name
				log.Printf("Debug: Leader Category: %s", category)

				for _, leader := range leaderCategory.Leaders {
					playerID := leader.Athlete.ID
					log.Printf("Debug: Leader: %s, Value: %f", leader.Athlete.Name, leader.Value)

					// Skip if player ID is not valid
					if !validPlayerIDs[playerID] {
						log.Printf("Warning: skipping leader stats for player ID %s (not found in database)", playerID)
						continue
					}

					statValue := leader.Value

					// Use DisplayName as stat type if available
					statType := leaderCategory.DisplayName
					if statType == "" {
						statType = category
					}

					log.Printf("Debug: Extracted leader stat - Player: %s, Category: %s, Type: %s, Value: %f",
						leader.Athlete.Name, category, statType, statValue)

					stats = append(stats, StatData{
						GameID:    gameID,
						PlayerID:  playerID,
						TeamID:    teamID,
						Category:  category,
						StatType:  statType,
						StatValue: statValue,
					})
				}
			}
		}
	*/
	///

	return stats, nil
}

// parseStatValue converts a string stat value to a float64

func parseStatValue(raw string) (float64, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" || raw == "--" {
		return 0, nil
	}
	// Remove trailing percentage if it exists.
	raw = strings.TrimSuffix(raw, "%")

	// Handle time format like "29:30"
	if strings.Contains(raw, ":") {
		parts := strings.Split(raw, ":")
		if len(parts) == 2 {
			minutes, err := strconv.ParseFloat(parts[0], 64)
			if err != nil {
				return 0, err
			}
			seconds, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				return 0, err
			}
			// Example: return total seconds
			return minutes*60 + seconds, nil
		}
	}

	// Handle slashes (e.g. "21/31")
	if strings.Contains(raw, "/") {
		parts := strings.Split(raw, "/")
		if len(parts) > 0 {
			return strconv.ParseFloat(parts[0], 64)
		}
	}

	// Only split on '-' if it is not at the start (to avoid negative numbers)
	if strings.Contains(raw, "-") && !strings.HasPrefix(raw, "-") {
		parts := strings.Split(raw, "-")
		if len(parts) > 0 {
			return strconv.ParseFloat(parts[0], 64)
		}
	}

	// Default: attempt to parse the raw value as a float.
	return strconv.ParseFloat(raw, 64)
}

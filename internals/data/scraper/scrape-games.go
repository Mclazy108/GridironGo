package scraper

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/Mclazy108/GridironGo/internals/data"
	//"github.com/Mclazy108/GridironGo/internals/data/sqlc"
	"golang.org/x/time/rate"
)

// NFLScraper handles fetching and storing NFL data
type NFLScraper struct {
	DB *data.DB
}

// NewScraper creates a new scraper that can populate the database
func NewScraper(db *data.DB) *NFLScraper {
	return &NFLScraper{
		DB: db,
	}
}

// ScoreboardResponse represents the JSON structure returned by the ESPN API
type ScoreboardResponse struct {
	Events []Event `json:"events"`
}

// Event represents a game event from the ESPN API
type Event struct {
	ID        string `json:"id"`
	Date      string `json:"date"`
	Name      string `json:"name"`
	ShortName string `json:"shortName"`
	Season    Season `json:"season"`
	Week      Week   `json:"week"`
}

// Season represents season information from the ESPN API
type Season struct {
	Year int `json:"year"`
}

// Week represents week information from the ESPN API
type Week struct {
	Number int `json:"number"`
}

// GameWeekJob represents a job to scrape games for a specific season and week
type GameWeekJob struct {
	Season int
	Week   int
}

// GameData holds processed game data ready for database insertion
type GameData struct {
	EventID   int64
	Date      string
	Name      string
	ShortName string
	Season    int64
	Week      int64
	AwayTeam  string
	HomeTeam  string
}

// insertNFLGamesBulk performs a bulk insert of game data
func insertNFLGamesBulk(ctx context.Context, tx *sql.Tx, games []GameData) error {
	if len(games) == 0 {
		return nil
	}

	// Build query
	query := `INSERT INTO nfl_games (
		event_id, date, name, short_name, season, week, away_team, home_team
	) VALUES `

	// Collect value placeholders like (?, ?, ?, ...), (?, ?, ?, ...), ...
	valueStrings := make([]string, 0, len(games))
	valueArgs := make([]interface{}, 0, len(games)*8)

	for _, g := range games {
		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?, ?)")
		valueArgs = append(valueArgs,
			g.EventID, g.Date, g.Name, g.ShortName, g.Season, g.Week, g.AwayTeam, g.HomeTeam,
		)
	}

	query += strings.Join(valueStrings, ",") +
		` ON CONFLICT(event_id) DO UPDATE SET
			date = excluded.date,
			name = excluded.name,
			short_name = excluded.short_name,
			season = excluded.season,
			week = excluded.week,
			away_team = excluded.away_team,
			home_team = excluded.home_team`

	// Prepare + exec
	_, err := tx.ExecContext(ctx, query, valueArgs...)
	return err
}

// ScrapeNFLGames fetches and stores NFL game data for multiple seasons using worker goroutines
func (s *NFLScraper) ScrapeNFLGames(ctx context.Context, seasons []int) error {
	log.Println("Starting NFL games scraping process with parallel workers...")
	log.Println("Press Ctrl+C to cancel the scraping process gracefully")

	// Create a rate limiter to avoid overwhelming the API
	// Limit to 10 requests per second (adjust as needed)
	limiter := rate.NewLimiter(10, 1)

	// Create a wait group to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Create job channel
	jobChan := make(chan GameWeekJob, 100)

	// Track stats
	var processedWeeks int32 = 0
	var totalGames int32 = 0
	var failedWeeks int32 = 0

	// Number of worker goroutines to process weeks
	numWorkers := 18 // One worker per week of NFL season
	log.Printf("Starting %d week worker goroutines", numWorkers)

	// Launch worker goroutines
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for job := range jobChan {
				log.Printf("Worker %d: Processing Season %d, Week %d",
					workerID, job.Season, job.Week)

				// Fetch games for this week and year
				events, err := s.fetchEvents(ctx, job.Season, job.Week, limiter)
				if err != nil {
					log.Printf("Worker %d: Error fetching games for Season %d, Week %d: %v",
						workerID, job.Season, job.Week, err)
					atomic.AddInt32(&failedWeeks, 1)
					continue
				}

				if len(events) == 0 {
					log.Printf("Worker %d: No games found for Season %d, Week %d",
						workerID, job.Season, job.Week)
					atomic.AddInt32(&processedWeeks, 1)
					continue
				}

				// Process events and prepare for bulk insert
				gameData := make([]GameData, 0, len(events))
				for _, event := range events {
					// Convert event ID string to int64
					var eventID int64
					_, err := fmt.Sscanf(event.ID, "%d", &eventID)
					if err != nil {
						// Try alternate parsing if simple scanf fails
						var temp int64
						for i := 0; i < len(event.ID); i++ {
							if event.ID[i] >= '0' && event.ID[i] <= '9' {
								temp = temp*10 + int64(event.ID[i]-'0')
							}
						}

						if temp > 0 {
							eventID = temp
						} else {
							log.Printf("Worker %d: Error parsing event ID '%s': %v",
								workerID, event.ID, err)
							continue
						}
					}

					// Extract home and away teams from name
					awayTeam, homeTeam := extractTeams(event.Name)

					// Skip games where team extraction failed
					if awayTeam == "" || homeTeam == "" {
						log.Printf("Worker %d: Skipping game with ID %s due to missing team information",
							workerID, event.ID)
						continue
					}

					// Format date for better readability in database
					formattedDate := event.Date
					if len(event.Date) >= 10 {
						formattedDate = event.Date[:10]
					}

					gameData = append(gameData, GameData{
						EventID:   eventID,
						Date:      formattedDate,
						Name:      event.Name,
						ShortName: event.ShortName,
						Season:    int64(event.Season.Year),
						Week:      int64(event.Week.Number),
						AwayTeam:  awayTeam,
						HomeTeam:  homeTeam,
					})
				}

				// Insert into database in a single transaction
				if len(gameData) > 0 {
					tx, err := s.DB.BeginTx(ctx, nil)
					if err != nil {
						log.Printf("Worker %d: Failed to begin transaction: %v", workerID, err)
						atomic.AddInt32(&failedWeeks, 1)
						continue
					}

					err = insertNFLGamesBulk(ctx, tx, gameData)
					if err != nil {
						log.Printf("Worker %d: Error inserting bulk games for Season %d, Week %d: %v",
							workerID, job.Season, job.Week, err)
						_ = tx.Rollback()
						atomic.AddInt32(&failedWeeks, 1)
						continue
					}

					if err := tx.Commit(); err != nil {
						log.Printf("Worker %d: Failed to commit transaction: %v", workerID, err)
						atomic.AddInt32(&failedWeeks, 1)
						continue
					}

					gamesInserted := len(gameData)
					atomic.AddInt32(&totalGames, int32(gamesInserted))
					log.Printf("Worker %d: Successfully inserted %d games for Season %d, Week %d",
						workerID, gamesInserted, job.Season, job.Week)
				}

				// Increment counter for processed weeks
				weeksProcessed := atomic.AddInt32(&processedWeeks, 1)
				totalJobs := len(seasons) * 18 // Assuming 18 weeks per season
				log.Printf("Progress: %d/%d weeks processed, %d total games",
					weeksProcessed, totalJobs, atomic.LoadInt32(&totalGames))
			}
			log.Printf("Worker %d finished", workerID)
		}(i)
	}

	// Create jobs for all seasons and weeks
	totalJobs := 0
	for _, year := range seasons {
		// Scrape regular season weeks 1-18
		for week := 1; week <= 18; week++ {
			select {
			case <-ctx.Done():
				log.Println("Job creation cancelled by user")
				close(jobChan)
				return ctx.Err()
			default:
				jobChan <- GameWeekJob{Season: year, Week: week}
				totalJobs++
			}
		}
	}

	log.Printf("Created %d jobs for %d seasons", totalJobs, len(seasons))

	// Close the job channel when all jobs are queued
	close(jobChan)

	// Wait for all workers to finish
	log.Println("Waiting for all workers to finish...")
	wg.Wait()

	log.Printf("NFL games scraping completed: %d/%d weeks processed, %d total games, %d failed weeks",
		processedWeeks, totalJobs, totalGames, failedWeeks)
	return nil
}

// fetchEvents fetches NFL games for a specific year and week from the ESPN API
func (s *NFLScraper) fetchEvents(ctx context.Context, year int, week int, limiter *rate.Limiter) ([]Event, error) {
	// Wait for rate limiter
	if err := limiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limiter error: %w", err)
	}

	// Construct the API URL
	url := fmt.Sprintf("https://site.api.espn.com/apis/site/v2/sports/football/nfl/scoreboard?dates=%d&seasontype=2&week=%d", year, week)

	// Create a new request with context
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

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
		return nil, fmt.Errorf("API returned non-OK status: %d", resp.StatusCode)
	}

	// Parse the JSON response
	var scoreboardResponse ScoreboardResponse
	err = json.NewDecoder(resp.Body).Decode(&scoreboardResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	return scoreboardResponse.Events, nil
}

// extractTeams extracts away and home teams from the game name
func extractTeams(gameName string) (string, string) {
	// Standard format in ESPN API is "Team A at Team B"
	separator := " at "
	pos := -1

	// Find the position of the separator
	for i := 0; i <= len(gameName)-len(separator); i++ {
		if gameName[i:i+len(separator)] == separator {
			pos = i
			break
		}
	}

	// If separator found, extract teams
	if pos > 0 {
		awayTeam := gameName[:pos]
		homeTeam := gameName[pos+len(separator):]
		return awayTeam, homeTeam
	}

	// If we couldn't parse, log and return empty strings
	log.Printf("Failed to extract teams from game name: %s", gameName)
	return "", ""
}


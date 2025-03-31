package data

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Mclazy108/GridironGo/internals/data/sqlc"
)

// NFLScraper handles fetching and storing NFL data
type NFLScraper struct {
	DB *DB
}

// NewScraper creates a new scraper that can populate the database
func NewScraper(db *DB) *NFLScraper {
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

// ScrapeNFLGames fetches and stores NFL game data for multiple seasons
func (s *NFLScraper) ScrapeNFLGames(ctx context.Context, seasons []int) error {
	log.Println("Starting NFL games scraping process...")
	log.Println("Press Ctrl+C to cancel the scraping process gracefully")

	for _, year := range seasons {
		// Check if context was cancelled
		select {
		case <-ctx.Done():
			log.Println("Scraping cancelled by user")
			return nil
		default:
			// Continue processing
		}

		log.Printf("Fetching games for season %d...", year)

		// Scrape regular season weeks 1-18
		for week := 1; week <= 18; week++ {
			// Check if context was cancelled
			select {
			case <-ctx.Done():
				log.Println("Scraping cancelled by user")
				return nil
			default:
				// Continue processing
			}

			log.Printf("Fetching week %d of %d season...", week, year)

			// Fetch games for this week and year
			events, err := s.fetchEvents(ctx, year, week)
			if err != nil {
				log.Printf("Error fetching games for Week %d, Year %d: %v", week, year, err)
				continue
			}

			// Insert events into database
			for _, event := range events {
				// Check if context was cancelled
				select {
				case <-ctx.Done():
					log.Println("Scraping cancelled by user")
					return nil
				default:
					// Continue processing
				}

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
						log.Printf("Error parsing event ID '%s': %v", event.ID, err)
						continue
					}
				}

				// Extract home and away teams from name
				awayTeam, homeTeam := extractTeams(event.Name)

				// Skip games where team extraction failed
				if awayTeam == "" || homeTeam == "" {
					log.Printf("Skipping game with ID %s due to missing team information", event.ID)
					continue
				}

				// Format date for better readability in database
				formattedDate := event.Date
				if len(event.Date) >= 10 {
					formattedDate = event.Date[:10]
				}

				// Create database parameters
				params := sqlc.CreateGameParams{
					EventID:   eventID,
					Date:      formattedDate,
					Name:      event.Name,
					ShortName: event.ShortName,
					Season:    int64(event.Season.Year),
					Week:      int64(event.Week.Number),
					AwayTeam:  awayTeam,
					HomeTeam:  homeTeam,
				}

				// Try to get the game first to see if it exists
				existingGame, err := s.DB.Queries.GetGame(ctx, eventID)
				if err == nil {
					// Game exists, check if we need to update it
					log.Printf("Game with ID %d already exists: %s", eventID, existingGame.Name)

					// Check if game data has changed
					if existingGame.Date != formattedDate ||
						existingGame.Name != event.Name ||
						existingGame.ShortName != event.ShortName {

						// Format date for the update as well
						formattedDate := event.Date
						if len(event.Date) >= 10 {
							formattedDate = event.Date[:10]
						}

						// Update the game
						updateParams := sqlc.UpdateGameParams{
							EventID:   eventID,
							Date:      formattedDate,
							Name:      event.Name,
							ShortName: event.ShortName,
							Season:    int64(event.Season.Year),
							Week:      int64(event.Week.Number),
							AwayTeam:  awayTeam,
							HomeTeam:  homeTeam,
						}

						err = s.DB.Queries.UpdateGame(ctx, updateParams)
						if err != nil {
							log.Printf("Error updating game with ID %d: %v", eventID, err)
						} else {
							log.Printf("Updated game: %s (ID: %d)", event.Name, eventID)
						}
					}

					continue
				}

				// Insert new game into database
				err = s.DB.Queries.CreateGame(ctx, params)
				if err != nil {
					log.Printf("Error inserting game with ID %d: %v", eventID, err)
					continue
				}

				log.Printf("Inserted game: %s (ID: %d)", event.Name, eventID)
			}

			// Sleep to avoid rate limiting, but make it interruptible
			select {
			case <-ctx.Done():
				log.Println("Scraping cancelled by user during rate limit sleep")
				return nil
			case <-time.After(500 * time.Millisecond):
				// Continue with the next week
			}
		}
	}

	log.Println("NFL games scraping completed successfully")
	return nil
}

// fetchEvents fetches NFL games for a specific year and week from the ESPN API
func (s *NFLScraper) fetchEvents(ctx context.Context, year int, week int) ([]Event, error) {
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


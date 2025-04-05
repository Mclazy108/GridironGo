package scraper

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/Mclazy108/GridironGo/internals/data"
	"github.com/Mclazy108/GridironGo/internals/data/sqlc"
	"log"
	"net/http"
	"strings"
	"time"
)

// TeamScraper handles fetching and storing NFL team data
type TeamScraper struct {
	DB *data.DB
}

// TeamListResponse represents the top-level response from the ESPN teams API
type TeamListResponse struct {
	Items []TeamItem `json:"items"`
}

// TeamItem represents a reference to a team in the ESPN API
type TeamItem struct {
	Ref string `json:"$ref"`
	ID  string `json:"id"`
}

// TeamDetails represents the detailed team information from the ESPN API
type TeamDetails struct {
	ID             string       `json:"id"`
	UID            string       `json:"uid"`
	Slug           string       `json:"slug"`
	Abbreviation   string       `json:"abbreviation"`
	DisplayName    string       `json:"displayName"`
	ShortName      string       `json:"shortName"`
	Name           string       `json:"name"`
	Nickname       string       `json:"nickname"`
	Location       string       `json:"location"`
	Color          string       `json:"color"`
	AlternateColor string       `json:"alternateColor"`
	IsActive       bool         `json:"isActive"`
	IsAllStar      bool         `json:"isAllStar"`
	Logo           string       `json:"logo"`
	Links          []TeamLink   `json:"links"`
	Venue          TeamVenue    `json:"venue"`
	Conference     TeamCategory `json:"conference"`
	Division       TeamCategory `json:"division"`
}

// TeamLink represents a link related to a team
type TeamLink struct {
	Rel        []string `json:"rel"`
	Href       string   `json:"href"`
	Text       string   `json:"text"`
	IsExternal bool     `json:"isExternal"`
	IsPremium  bool     `json:"isPremium"`
}

// TeamVenue represents the venue information for a team
type TeamVenue struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// TeamCategory represents a division or conference
type TeamCategory struct {
	ID           string `json:"id"`
	UID          string `json:"uid"`
	Name         string `json:"name"`
	Abbreviation string `json:"abbreviation"`
}

// ScrapeNFLTeams fetches and stores NFL team data
func (s *TeamScraper) ScrapeNFLTeams(ctx context.Context) error {
	log.Println("Starting NFL teams scraping process...")
	log.Println("Press Ctrl+C to cancel the scraping process gracefully")

	// First, get the list of team IDs from the main endpoint
	teamItems, err := s.fetchTeamList(ctx)
	if err != nil {
		return fmt.Errorf("error fetching team list: %w", err)
	}

	log.Printf("Found %d NFL teams to process", len(teamItems))

	// Process each team to get detailed information
	for i, teamItem := range teamItems {
		// Check if context was cancelled
		select {
		case <-ctx.Done():
			log.Println("Scraping cancelled by user")
			return nil
		default:
			// Continue processing
		}

		log.Printf("Processing team %d of %d (ID: %s)...", i+1, len(teamItems), teamItem.ID)

		// Fetch detailed team information
		teamDetails, err := s.fetchTeamDetails(ctx, teamItem.ID)
		if err != nil {
			log.Printf("Error fetching details for team ID %s: %v", teamItem.ID, err)
			continue
		}

		// Check if the team already exists in the database
		existingTeam, err := s.DB.Queries.GetNFLTeam(ctx, teamItem.ID)
		if err == nil {
			// Team exists, check if we need to update it
			log.Printf("Team with ID %s already exists: %s", teamItem.ID, existingTeam.DisplayName)

			// Check if team data has changed
			primaryColorChanged := (existingTeam.PrimaryColor.Valid && existingTeam.PrimaryColor.String != teamDetails.Color) ||
				(!existingTeam.PrimaryColor.Valid && teamDetails.Color != "")

			secondaryColorChanged := (existingTeam.SecondaryColor.Valid && existingTeam.SecondaryColor.String != teamDetails.AlternateColor) ||
				(!existingTeam.SecondaryColor.Valid && teamDetails.AlternateColor != "")

			logoUrlChanged := (existingTeam.LogoUrl.Valid && existingTeam.LogoUrl.String != teamDetails.Logo) ||
				(!existingTeam.LogoUrl.Valid && teamDetails.Logo != "")

			if existingTeam.DisplayName != teamDetails.DisplayName ||
				existingTeam.Abbreviation != teamDetails.Abbreviation ||
				existingTeam.Location != teamDetails.Location ||
				existingTeam.Nickname != teamDetails.Nickname ||
				primaryColorChanged ||
				secondaryColorChanged ||
				logoUrlChanged {

				// Update the team
				updateParams := sqlc.UpdateNFLTeamParams{
					TeamID:         teamItem.ID,
					DisplayName:    teamDetails.DisplayName,
					Abbreviation:   teamDetails.Abbreviation,
					ShortName:      teamDetails.ShortName,
					Location:       teamDetails.Location,
					Nickname:       teamDetails.Nickname,
					Conference:     teamDetails.Conference.Name,
					Division:       teamDetails.Division.Name,
					PrimaryColor:   sql.NullString{String: teamDetails.Color, Valid: teamDetails.Color != ""},
					SecondaryColor: sql.NullString{String: teamDetails.AlternateColor, Valid: teamDetails.AlternateColor != ""},
					LogoUrl:        sql.NullString{String: teamDetails.Logo, Valid: teamDetails.Logo != ""},
				}

				err = s.DB.Queries.UpdateNFLTeam(ctx, updateParams)
				if err != nil {
					log.Printf("Error updating team with ID %s: %v", teamItem.ID, err)
				} else {
					log.Printf("Updated team: %s (ID: %s)", teamDetails.DisplayName, teamItem.ID)
				}
			}
			continue
		}

		// Insert new team into database
		params := sqlc.CreateNFLTeamParams{
			TeamID:         teamItem.ID,
			DisplayName:    teamDetails.DisplayName,
			Abbreviation:   teamDetails.Abbreviation,
			ShortName:      teamDetails.ShortName,
			Location:       teamDetails.Location,
			Nickname:       teamDetails.Nickname,
			Conference:     teamDetails.Conference.Name,
			Division:       teamDetails.Division.Name,
			PrimaryColor:   sql.NullString{String: teamDetails.Color, Valid: teamDetails.Color != ""},
			SecondaryColor: sql.NullString{String: teamDetails.AlternateColor, Valid: teamDetails.AlternateColor != ""},
			LogoUrl:        sql.NullString{String: teamDetails.Logo, Valid: teamDetails.Logo != ""},
		}

		err = s.DB.Queries.CreateNFLTeam(ctx, params)
		if err != nil {
			log.Printf("Error inserting team with ID %s: %v", teamItem.ID, err)
			continue
		}

		log.Printf("Inserted team: %s (ID: %s)", teamDetails.DisplayName, teamItem.ID)

		// Sleep to avoid rate limiting, but make it interruptible
		select {
		case <-ctx.Done():
			log.Println("Scraping cancelled by user during rate limit sleep")
			return nil
		case <-time.After(300 * time.Millisecond):
			// Continue with the next team
		}
	}

	log.Println("NFL teams scraping completed successfully")
	return nil
}

// fetchTeamList fetches the list of NFL teams from the ESPN API
func (s *TeamScraper) fetchTeamList(ctx context.Context) ([]TeamItem, error) {
	// Construct the API URL
	url := "https://sports.core.api.espn.com/v2/sports/football/leagues/nfl/teams"

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
	var teamListResponse TeamListResponse
	err = json.NewDecoder(resp.Body).Decode(&teamListResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	// Process the items to extract team IDs from URLs
	var teams []TeamItem
	for _, item := range teamListResponse.Items {
		// The $ref URL typically has the format ".../teams/{team_id}"
		// Extract the team ID from the reference URL if it's not directly available
		teamID := item.ID
		if teamID == "" && item.Ref != "" {
			parts := strings.Split(item.Ref, "/")
			if len(parts) > 0 {
				teamID = parts[len(parts)-1] // Get the last part of the URL
			}
		}

		// Clean up the team ID - remove any query parameters
		if teamID != "" {
			// Split on ? and take just the first part
			cleanID := strings.Split(teamID, "?")[0]
			teams = append(teams, TeamItem{
				Ref: item.Ref,
				ID:  cleanID,
			})
		} else {
			log.Printf("Warning: Could not extract team ID from reference: %s", item.Ref)
		}
	}

	return teams, nil
}

// fetchTeamDetails fetches detailed information for a specific team
func (s *TeamScraper) fetchTeamDetails(ctx context.Context, teamID string) (*TeamDetails, error) {
	// Construct the API URL for detailed team information
	url := fmt.Sprintf("https://site.api.espn.com/apis/site/v2/sports/football/nfl/teams/%s", teamID)

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
	var response struct {
		Team TeamDetails `json:"team"`
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	return &response.Team, nil
}

// NewTeamScraper creates a new scraper for NFL team data
func NewTeamScraper(db *data.DB) *TeamScraper {
	return &TeamScraper{
		DB: db,
	}
}

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
	"strings"
)

type PlayerScraper struct {
	DB *data.DB
}

func NewPlayerScraper(db *data.DB) *PlayerScraper {
	return &PlayerScraper{DB: db}
}

type PlayerDetails struct {
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
	DateOfBirth string `json:"dateOfBirth"`
	Jersey      string `json:"jersey"`
	Position    struct {
		Name         string `json:"name"`
		Abbreviation string `json:"abbreviation"`
	} `json:"position"`
	Active     bool `json:"active"`
	DebutYear  int  `json:"debutYear"`
	Status     string
	ImageURL   string
	Experience int `json:"experience"`
	College    struct {
		Name string `json:"name"`
	} `json:"college"`
	Team struct {
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
}

func (s *PlayerScraper) fetchTeamRoster(ctx context.Context, teamID string) ([]string, error) {
	url := fmt.Sprintf("https://sports.core.api.espn.com/v2/sports/football/leagues/nfl/seasons/2024/teams/%s/athletes?limit=200", teamID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return nil, fmt.Errorf("request cancelled: %w", ctx.Err())
		}
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned non-OK status: %d. Response: %s", resp.StatusCode, string(body))
	}

	var rosterResponse struct {
		Items []struct {
			Ref string `json:"$ref"`
		} `json:"items"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&rosterResponse); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	var playerIDs []string
	for _, item := range rosterResponse.Items {
		if item.Ref != "" {
			parts := strings.Split(item.Ref, "/")
			lastSegment := parts[len(parts)-1]
			playerID := strings.Split(lastSegment, "?")[0]
			if playerID != "" {
				playerIDs = append(playerIDs, playerID)
			}
		}
	}
	return playerIDs, nil
}

func (s *PlayerScraper) fetchPlayerDetails(ctx context.Context, playerID string) (*PlayerDetails, error) {
	url := fmt.Sprintf("https://site.web.api.espn.com/apis/common/v3/sports/football/nfl/athletes/%s/overview", playerID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return nil, fmt.Errorf("request cancelled: %w", ctx.Err())
		}
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned non-OK status: %d. Response: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Player   PlayerDetails `json:"athlete"`
		Headshot struct {
			Href string `json:"href"`
		} `json:"headshot"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	result.Player.ImageURL = result.Headshot.Href
	return &result.Player, nil
}

func (s *PlayerScraper) ScrapeNFLPlayers(ctx context.Context) error {
	log.Println("Starting NFL players scraping process...")

	teams, err := s.DB.Queries.GetAllNFLTeams(ctx)
	if err != nil || len(teams) == 0 {
		return fmt.Errorf("failed to fetch NFL teams from DB: %w", err)
	}

	log.Printf("Found %d teams. Fetching player data from team rosters", len(teams))
	processed := make(map[string]bool)

	for _, team := range teams {
		playerIDs, err := s.fetchTeamRoster(ctx, team.TeamID)
		if err != nil {
			log.Printf("Error fetching roster for team %s: %v", team.DisplayName, err)
			continue
		}

		log.Printf("Found %d players on %s roster", len(playerIDs), team.DisplayName)

		for _, playerID := range playerIDs {
			if processed[playerID] {
				continue
			}
			processed[playerID] = true

			if err := s.processPlayer(ctx, playerID); err != nil {
				log.Printf("Error processing player ID %s: %v", playerID, err)
			}
		}
	}

	log.Printf("Processed %d unique players from team rosters", len(processed))
	log.Println("NFL players scraping completed successfully")
	return nil
}

func (s *PlayerScraper) processPlayer(ctx context.Context, playerID string) error {
	playerDetails, err := s.fetchPlayerDetails(ctx, playerID)
	if err != nil {
		return fmt.Errorf("error fetching details: %w", err)
	}

	if playerDetails.Position.Abbreviation == "" {
		log.Printf("Skipping player %s - no position data", playerDetails.FullName)
		return nil
	}

	_, err = s.DB.Queries.GetNFLPlayer(ctx, playerID)
	//err = s.DB.Queries.GetNFLPlayer(ctx, playerID)
	//_, err := s.DB.Queries.GetNFLPlayer(ctx, playerID)
	//existingPlayer, err := s.DB.Queries.GetNFLPlayer(ctx, playerID)
	if err == nil {
		updateParams := sqlc.UpdateNFLPlayerParams{
			PlayerID:   playerID,
			FirstName:  playerDetails.FirstName,
			LastName:   playerDetails.LastName,
			FullName:   playerDetails.FullName,
			Position:   playerDetails.Position.Abbreviation,
			TeamID:     sql.NullString{String: playerDetails.Team.ID, Valid: playerDetails.Team.ID != ""},
			Jersey:     sql.NullString{String: playerDetails.Jersey, Valid: playerDetails.Jersey != ""},
			Height:     sql.NullInt64{Int64: int64(playerDetails.Height), Valid: playerDetails.Height > 0},
			Weight:     sql.NullInt64{Int64: int64(playerDetails.Weight), Valid: playerDetails.Weight > 0},
			Active:     playerDetails.Active,
			College:    sql.NullString{String: playerDetails.College.Name, Valid: playerDetails.College.Name != ""},
			Experience: sql.NullInt64{Int64: int64(playerDetails.Experience), Valid: playerDetails.Experience >= 0},
			DraftYear:  sql.NullInt64{Int64: int64(playerDetails.DraftInfo.Year), Valid: playerDetails.DraftInfo.Year > 0},
			DraftRound: sql.NullInt64{Int64: int64(playerDetails.DraftInfo.Round), Valid: playerDetails.DraftInfo.Round > 0},
			DraftPick:  sql.NullInt64{Int64: int64(playerDetails.DraftInfo.Selection), Valid: playerDetails.DraftInfo.Selection > 0},
			Status:     sql.NullString{String: playerDetails.Status, Valid: playerDetails.Status != ""},
			ImageUrl:   sql.NullString{String: playerDetails.ImageURL, Valid: playerDetails.ImageURL != ""},
		}
		return s.DB.Queries.UpdateNFLPlayer(ctx, updateParams)
	}

	insertParams := sqlc.CreateNFLPlayerParams{
		PlayerID:   playerID,
		FirstName:  playerDetails.FirstName,
		LastName:   playerDetails.LastName,
		FullName:   playerDetails.FullName,
		Position:   playerDetails.Position.Abbreviation,
		TeamID:     sql.NullString{String: playerDetails.Team.ID, Valid: playerDetails.Team.ID != ""},
		Jersey:     sql.NullString{String: playerDetails.Jersey, Valid: playerDetails.Jersey != ""},
		Height:     sql.NullInt64{Int64: int64(playerDetails.Height), Valid: playerDetails.Height > 0},
		Weight:     sql.NullInt64{Int64: int64(playerDetails.Weight), Valid: playerDetails.Weight > 0},
		Active:     playerDetails.Active,
		College:    sql.NullString{String: playerDetails.College.Name, Valid: playerDetails.College.Name != ""},
		Experience: sql.NullInt64{Int64: int64(playerDetails.Experience), Valid: playerDetails.Experience >= 0},
		DraftYear:  sql.NullInt64{Int64: int64(playerDetails.DraftInfo.Year), Valid: playerDetails.DraftInfo.Year > 0},
		DraftRound: sql.NullInt64{Int64: int64(playerDetails.DraftInfo.Round), Valid: playerDetails.DraftInfo.Round > 0},
		DraftPick:  sql.NullInt64{Int64: int64(playerDetails.DraftInfo.Selection), Valid: playerDetails.DraftInfo.Selection > 0},
		Status:     sql.NullString{String: playerDetails.Status, Valid: playerDetails.Status != ""},
		ImageUrl:   sql.NullString{String: playerDetails.ImageURL, Valid: playerDetails.ImageURL != ""},
	}
	return s.DB.Queries.CreateNFLPlayer(ctx, insertParams)
}

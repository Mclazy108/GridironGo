package data

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Mclazy108/GridironGo/internals/data/sqlc"
	"github.com/gocolly/colly/v2"
)

type NFLScraper struct {
	DB *DB
}

// NewScraper creates a new scraper that can populate the database
func NewScraper(db *DB) *NFLScraper {
	return &NFLScraper{
		DB: db,
	}
}

// ScrapeTeams scrapes NFL team data from Pro Football Reference
func (s *NFLScraper) ScrapeTeams(ctx context.Context) error {
	log.Println("Starting team scraping process...")

	// Set up logging to file
	logFile, err := os.OpenFile("scrape.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer logFile.Close()

	// Create a multi-writer to write to both console and file
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	logger := log.New(multiWriter, "", log.LstdFlags)

	// Define the target URL for scraping team info
	targetURL := "https://www.pro-football-reference.com/teams/"
	logger.Printf("Target URL for team scraping: %s", targetURL)

	// Parse the URL to verify it's valid
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		logger.Printf("WARNING: URL parsing error: %v", err)
	} else {
		logger.Printf("Parsed URL - Scheme: %s, Host: %s, Path: %s",
			parsedURL.Scheme, parsedURL.Host, parsedURL.Path)
	}

	c := colly.NewCollector(
		colly.AllowedDomains("www.pro-football-reference.com"),
		colly.UserAgent("GridironGo Fantasy Football App v1.0"),
	)

	// Set some limits and timeouts
	c.SetRequestTimeout(120 * time.Second)
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		RandomDelay: 5 * time.Second,
	})

	type Team struct {
		Name         string
		City         string
		Abbreviation string
		Conference   string
		Division     string
	}

	teams := make(map[string]Team)
	logger.Println("Initialized team collector")

	// Process the team links directly - this is much simpler and more reliable
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		href := e.Attr("href")

		// Only process direct team links - they follow the pattern /teams/XXX/
		// But exclude any links with years in them like /teams/buf/2024.htm
		if strings.HasPrefix(href, "/teams/") && !strings.Contains(href, ".htm") && len(href) > 7 {
			teamName := e.Text
			logger.Printf("Processing team link: %s (%s)", teamName, href)

			// Extract abbreviation from the href
			teamAbbr := ""
			path := strings.TrimSuffix(href, "/")
			path = strings.TrimPrefix(path, "/teams/")
			if path != "" && path != "index" {
				teamAbbr = path
				logger.Printf("  Team abbreviation: %s", teamAbbr)

				// Split the team name into city and nickname
				cityName, teamNickname := splitTeamName(teamName)
				logger.Printf("  City: %s, Team name: %s", cityName, teamNickname)

				// Determine conference and division
				conference, division := getTeamConferenceAndDivision(teamAbbr)
				logger.Printf("  Conference: %s, Division: %s", conference, division)

				// Store in our map
				teams[teamAbbr] = Team{
					Name:         teamNickname,
					City:         cityName,
					Abbreviation: teamAbbr,
					Conference:   conference,
					Division:     division,
				}

				logger.Printf("  Team %s (%s %s) added to collection", teamAbbr, cityName, teamNickname)
			}
		}
	})

	c.OnRequest(func(r *colly.Request) {
		logger.Printf("Making HTTP request to: %s", r.URL.String())
	})

	c.OnResponse(func(r *colly.Response) {
		logger.Printf("Received response from %s: status=%d, length=%d bytes",
			r.Request.URL, r.StatusCode, len(r.Body))
	})

	c.OnError(func(r *colly.Response, err error) {
		logger.Printf("Error when visiting %s: %v", r.Request.URL, err)
		logger.Printf("Response status code: %d", r.StatusCode)
	})

	c.OnScraped(func(r *colly.Response) {
		logger.Printf("Finished scraping: %s", r.Request.URL)
		logger.Printf("Total teams found: %d", len(teams))
	})

	// Start scraping
	logger.Println("Starting to visit teams page...")
	err = c.Visit(targetURL)
	if err != nil {
		return fmt.Errorf("failed to visit teams page: %w", err)
	}

	// Log all found teams
	logger.Println("Teams found during scraping:")
	for abbr, team := range teams {
		logger.Printf("- %s: %s %s (%s, %s)",
			abbr, team.City, team.Name, team.Conference, team.Division)
	}

	logger.Println("Scraping completed. Now inserting teams into database...")

	// Insert teams into the database
	for abbr, team := range teams {
		logger.Printf("Processing team: %s %s (%s)", team.City, team.Name, abbr)

		// Check if team already exists
		existingTeam, err := s.DB.Queries.GetNFLTeamByAbbreviation(ctx, team.Abbreviation)
		if err == nil {
			logger.Printf("Team %s already exists in database with ID %d, skipping...",
				team.Abbreviation, existingTeam.ID)
			continue
		} else {
			logger.Printf("Team %s not found in database (error: %v), will insert",
				team.Abbreviation, err)
		}

		// Insert team
		params := sqlc.InsertTeamParams{
			Name:         team.Name,
			City:         team.City,
			Abbreviation: team.Abbreviation,
			Conference:   team.Conference,
			Division:     team.Division,
		}

		logger.Printf("Inserting team with params: %+v", params)

		id, err := s.DB.Queries.InsertTeam(ctx, params)
		if err != nil {
			logger.Printf("ERROR: Failed to insert team %s: %v", team.Name, err)
		} else {
			logger.Printf("SUCCESS: Inserted team: %s %s with ID %d", team.City, team.Name, id)
		}
	}

	logger.Println("Team scraping and database insertion completed")
	return nil
}

// ScrapePlayers scrapes NFL player data from Pro Football Reference for specified seasons
func (s *NFLScraper) ScrapePlayers(ctx context.Context, seasons []int) error {
	// Set up logging to file
	logFile, err := os.OpenFile("scrape.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer logFile.Close()

	// Create a multi-writer to write to both console and file
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	logger := log.New(multiWriter, "", log.LstdFlags)

	logger.Println("Starting player scraping process...")

	for _, year := range seasons {
		logger.Printf("Scraping players for %d season...", year)

		// Find or create the season in our database
		var seasonID int64
		season, err := s.DB.Queries.GetSeasonByYear(ctx, int64(year))
		if err != nil {
			logger.Printf("Season %d not found, creating it...", year)
			// Create season if it doesn't exist
			startDate := fmt.Sprintf("%d-09-01", year) // Approximate NFL season start
			endDate := fmt.Sprintf("%d-02-15", year+1) // Approximate NFL season end

			seasonParams := sqlc.CreateSeasonParams{
				Year:      int64(year),
				StartDate: startDate,
				EndDate:   endDate,
				Current:   sql.NullInt64{Int64: 0, Valid: true},
			}

			// If it's 2024, mark as current season
			if year == 2024 {
				seasonParams.Current = sql.NullInt64{Int64: 1, Valid: true}
			}

			seasonID, err = s.DB.Queries.CreateSeason(ctx, seasonParams)
			if err != nil {
				return fmt.Errorf("failed to create season %d: %w", year, err)
			}
			logger.Printf("Created season %d with ID %d", year, seasonID)
		} else {
			seasonID = season.ID
			logger.Printf("Found existing season %d with ID %d", year, seasonID)
		}

		// Define the target URLs for the teams
		teams, err := s.DB.Queries.GetAllNFLTeams(ctx)
		if err != nil {
			return fmt.Errorf("failed to get NFL teams: %w", err)
		}

		for _, team := range teams {
			// URL format example: https://www.pro-football-reference.com/teams/kan/2023_roster.htm
			targetURL := fmt.Sprintf("https://www.pro-football-reference.com/teams/%s/%d_roster.htm",
				strings.ToLower(team.Abbreviation), year)

			logger.Printf("Scraping roster for %s %s (%d)", team.City, team.Name, year)
			logger.Printf("Target URL: %s", targetURL)

			err = s.scrapeTeamRoster(ctx, targetURL, team.ID, seasonID)
			if err != nil {
				logger.Printf("Error scraping roster for %s: %v", team.Abbreviation, err)
				continue
			}

			// Sleep to avoid hammering the server
			time.Sleep(1 * time.Second)
		}
	}

	logger.Println("Player scraping completed")
	return nil
}

// scrapeTeamRoster scrapes the roster for a specific team and season
func (s *NFLScraper) scrapeTeamRoster(ctx context.Context, url string, teamID int64, seasonID int64) error {
	// Set up logging to file
	logFile, err := os.OpenFile("scrape.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer logFile.Close()

	// Create a multi-writer to write to both console and file
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	logger := log.New(multiWriter, "", log.LstdFlags)

	// Use the logger instead of the global log
	logger.Printf("Starting to scrape roster for team ID %d (season ID %d)", teamID, seasonID)
	logger.Printf("Target URL: %s", url)

	c := colly.NewCollector(
		colly.AllowedDomains("www.pro-football-reference.com"),
		colly.UserAgent("GridironGo Fantasy Football App v1.0"),
	)

	// Set limits and timeouts
	c.SetRequestTimeout(60 * time.Second)
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		RandomDelay: 3 * time.Second,
	})

	playersCount := 0

	// Debug the entire HTML structure
	c.OnHTML("html", func(e *colly.HTMLElement) {
		logger.Printf("Investigating HTML structure of roster page")

		// Look for tables and print their IDs
		e.ForEach("table", func(i int, table *colly.HTMLElement) {
			id := table.Attr("id")
			logger.Printf("Table #%d has ID: '%s'", i+1, id)

			// Check for table headers to understand structure
			var headerCount = 0
			table.ForEach("th", func(_ int, header *colly.HTMLElement) {
				headerCount++
				text := header.Text
				datastat := header.Attr("data-stat")
				logger.Printf("    Header #%d: Text='%s', data-stat='%s'", headerCount, text, datastat)
			})
			logger.Printf("  Table has %d header cells", headerCount)

			// Check a sample row
			var cellCount = 0
			table.ForEach("tr:first-child td", func(_ int, cell *colly.HTMLElement) {
				cellCount++
				text := cell.Text
				datastat := cell.Attr("data-stat")
				logger.Printf("    Cell #%d: Text='%s', data-stat='%s'", cellCount, text, datastat)
			})
			logger.Printf("  First row has %d cells", cellCount)
		})

		// Check if there's a message about the season data
		var found = false
		e.ForEach("div", func(_ int, div *colly.HTMLElement) {
			if strings.Contains(div.Text, "roster") && (strings.Contains(div.Text, "2023") || strings.Contains(div.Text, "2024")) {
				logger.Printf("Found message about roster data: %s", div.Text)
				found = true
			}
		})
		if !found {
			logger.Printf("No messages about 'roster' found in div elements")
		}

		// Look for possible player containers
		var containerCount = 0
		e.ForEach("div, section", func(_ int, container *colly.HTMLElement) {
			if strings.Contains(strings.ToLower(container.Text), "roster") {
				containerCount++
				logger.Printf("Found potential roster container: %s", container.Attr("class")+" "+container.Attr("id"))
			}
		})
		logger.Printf("Found %d potential roster containers", containerCount)
	})

	// Process any kind of player data we can find
	c.OnHTML("tr", func(e *colly.HTMLElement) {
		// Skip header rows
		if e.Attr("class") == "thead" || e.Attr("class") == "over_header" {
			return
		}

		// Try to extract data from different possible structures
		playerName := e.ChildText("td[data-stat='player'], a, td:nth-child(2)")
		position := e.ChildText("td[data-stat='pos'], td:nth-child(3)")
		playerNum := e.ChildText("th[data-stat='number'], td:nth-child(1)")

		if playerName != "" {
			logger.Printf("Found potential player: %s, Position: %s, Number: %s",
				playerName, position, playerNum)
		}

		// Only process if we have a name and position
		if playerName == "" || position == "" {
			return
		}

		// Rest of the player processing code (same as before)
		jerseyNumber := sql.NullInt64{Valid: false}
		if playerNum != "" {
			num, err := strconv.ParseInt(playerNum, 10, 64)
			if err == nil {
				jerseyNumber = sql.NullInt64{Int64: num, Valid: true}
			}
		}

		// Check if player already exists
		var playerID int64 = 0
		players, err := s.DB.Queries.SearchPlayers(ctx, sql.NullString{String: playerName, Valid: true})
		if err == nil && len(players) > 0 {
			// Check if we have a matching player
			for _, p := range players {
				if p.Name == playerName && p.Position == position {
					playerID = p.ID
					break
				}
			}
		}

		if playerID == 0 {
			// Insert new player
			params := sqlc.InsertPlayerParams{
				Name:         playerName,
				Position:     position,
				TeamID:       sql.NullInt64{Int64: teamID, Valid: true},
				JerseyNumber: jerseyNumber,
				Status:       sql.NullString{String: "Active", Valid: true},
			}

			id, err := s.DB.Queries.InsertPlayer(ctx, params)
			if err != nil {
				logger.Printf("Failed to insert player %s: %v", playerName, err)
				return
			}
			playerID = id
			logger.Printf("Inserted new player: %s (%s) with ID %d", playerName, position, playerID)
		} else {
			// Update existing player
			params := sqlc.UpdatePlayerParams{
				ID:           playerID,
				Name:         playerName,
				Position:     position,
				TeamID:       sql.NullInt64{Int64: teamID, Valid: true},
				JerseyNumber: jerseyNumber,
				Status:       sql.NullString{String: "Active", Valid: true},
			}

			err := s.DB.Queries.UpdatePlayer(ctx, params)
			if err != nil {
				logger.Printf("Failed to update player %s: %v", playerName, err)
				return
			}
			logger.Printf("Updated existing player: %s (%s) with ID %d", playerName, position, playerID)
		}

		playersCount++
	})

	// Log a larger sample of the response body
	c.OnResponse(func(r *colly.Response) {
		logger.Printf("Received response from %s: status=%d, length=%d bytes",
			r.Request.URL, r.StatusCode, len(r.Body))

		// Save a sample of the HTML for debugging
		sample := string(r.Body)
		if len(sample) > 1000 {
			sample = sample[:1000]
		}
		logger.Printf("Sample HTML (first 1000 chars): %s", sample)

		// Look for some key patterns
		htmlStr := string(r.Body)
		if strings.Contains(htmlStr, "table id=\"roster\"") {
			logger.Printf("Found table#roster in HTML")
			idx := strings.Index(htmlStr, "table id=\"roster\"")
			context := htmlStr[idx:min(idx+200, len(htmlStr))]
			logger.Printf("Context around table#roster: %s", context)
		} else {
			logger.Printf("NO table#roster found in HTML")
		}

		// Check for other possible roster table identifiers
		possibleTableIDs := []string{"games_played_team", "games_played", "team_roster", "players", "roster_table"}
		for _, id := range possibleTableIDs {
			if strings.Contains(htmlStr, fmt.Sprintf("table id=\"%s\"", id)) {
				logger.Printf("Found table#%s in HTML", id)
				idx := strings.Index(htmlStr, fmt.Sprintf("table id=\"%s\"", id))
				context := htmlStr[idx:min(idx+200, len(htmlStr))]
				logger.Printf("Context around table#%s: %s", id, context)
			}
		}
	})

	// Standard request/error handling
	c.OnRequest(func(r *colly.Request) {
		logger.Printf("Making HTTP request to: %s", r.URL.String())
	})

	c.OnError(func(r *colly.Response, err error) {
		logger.Printf("Error when visiting %s: %v", r.Request.URL, err)
		logger.Printf("Response status code: %d", r.StatusCode)
	})

	c.OnScraped(func(r *colly.Response) {
		logger.Printf("Finished scraping: %s", r.Request.URL)
		logger.Printf("Total players found for team: %d", playersCount)
	})

	// Start scraping
	err = c.Visit(url)
	if err != nil {
		return fmt.Errorf("failed to visit URL %s: %w", url, err)
	}

	return nil
}

// ScrapePlayerStats scrapes player stats for specified seasons
func (s *NFLScraper) ScrapePlayerStats(ctx context.Context, seasons []int) error {
	// Set up logging to file
	logFile, err := os.OpenFile("scrape.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer logFile.Close()

	// Create a multi-writer to write to both console and file
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	logger := log.New(multiWriter, "", log.LstdFlags)

	logger.Println("Starting player stats scraping process...")

	for _, year := range seasons {
		logger.Printf("Scraping player stats for %d season...", year)

		// Get season ID
		season, err := s.DB.Queries.GetSeasonByYear(ctx, int64(year))
		if err != nil {
			return fmt.Errorf("season %d not found: %w", year, err)
		}

		// Scrape different stat categories
		statTypes := []string{"passing", "rushing", "receiving", "kicking", "defense"}

		for _, statType := range statTypes {
			logger.Printf("Scraping %s stats for %d season...", statType, year)

			// URL format example: https://www.pro-football-reference.com/years/2023/passing.htm
			targetURL := fmt.Sprintf("https://www.pro-football-reference.com/years/%d/%s.htm",
				year, statType)

			logger.Printf("Target URL: %s", targetURL)

			err = s.scrapeStatCategory(ctx, targetURL, statType, season.ID)
			if err != nil {
				logger.Printf("Error scraping %s stats: %v", statType, err)
				continue
			}

			// Sleep to avoid hammering the server
			time.Sleep(2 * time.Second)
		}
	}

	logger.Println("Player stats scraping completed")
	return nil
}

// scrapeStatCategory scrapes stats for a specific category and season
func (s *NFLScraper) scrapeStatCategory(ctx context.Context, url string, statType string, seasonID int64) error {
	// Set up logging to file
	logFile, err := os.OpenFile("scrape.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer logFile.Close()

	// Create a multi-writer to write to both console and file
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	logger := log.New(multiWriter, "", log.LstdFlags)

	c := colly.NewCollector(
		colly.AllowedDomains("www.pro-football-reference.com"),
		colly.UserAgent("GridironGo Fantasy Football App v1.0"),
	)

	// Set limits and timeouts
	c.SetRequestTimeout(60 * time.Second)
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		RandomDelay: 3 * time.Second,
	})

	statsCount := 0

	// Debug the entire HTML structure
	c.OnHTML("html", func(e *colly.HTMLElement) {
		logger.Printf("Investigating HTML structure of %s stats page", statType)

		// Look for tables and print their IDs
		e.ForEach("table", func(i int, table *colly.HTMLElement) {
			id := table.Attr("id")
			logger.Printf("Table #%d has ID: '%s'", i+1, id)

			// Check for table headers to understand structure
			var headerCount = 0
			table.ForEach("th", func(_ int, header *colly.HTMLElement) {
				headerCount++
				text := header.Text
				datastat := header.Attr("data-stat")
				logger.Printf("    Header #%d: Text='%s', data-stat='%s'", headerCount, text, datastat)
			})
			logger.Printf("  Table has %d header cells", headerCount)

			// Check a sample row
			var cellCount = 0
			table.ForEach("tr:first-child td", func(_ int, cell *colly.HTMLElement) {
				cellCount++
				text := cell.Text
				datastat := cell.Attr("data-stat")
				logger.Printf("    Cell #%d: Text='%s', data-stat='%s'", cellCount, text, datastat)
			})
			logger.Printf("  First row has %d cells", cellCount)
		})
	})

	// Process the stats table
	selector := fmt.Sprintf("table#%s tbody tr, table#%s_advanced tbody tr", statType, statType)
	c.OnHTML(selector, func(e *colly.HTMLElement) {
		// Skip header rows or divider rows
		if e.Attr("class") == "thead" || e.Attr("class") == "divider" || e.Attr("class") == "over_header" {
			return
		}

		// Get player name and find player ID
		playerName := e.ChildText("td[data-stat='player'], td[data-stat='player_name']")
		if playerName == "" {
			return
		}

		logger.Printf("Found player in %s stats: %s", statType, playerName)

		// Find the player in our database
		players, err := s.DB.Queries.SearchPlayers(ctx, sql.NullString{String: playerName, Valid: true})
		if err != nil || len(players) == 0 {
			logger.Printf("Player not found: %s", playerName)
			return
		}

		// For simplicity, use the first matching player
		// In a more robust implementation, you might want to match by team too
		playerID := players[0].ID

		// Create a game for the season stats (we'll use game_id = 0 for season totals)
		// This is a simplification - in a real implementation, you'd tie stats to actual games
		gameID := int64(0)
		week := int64(0)

		// Create player stats based on stat type
		switch statType {
		case "passing":
			attempts := parseNullInt(e.ChildText("td[data-stat='pass_att']"))
			completions := parseNullInt(e.ChildText("td[data-stat='pass_cmp']"))
			yards := parseNullInt(e.ChildText("td[data-stat='pass_yds']"))
			tds := parseNullInt(e.ChildText("td[data-stat='pass_td']"))
			ints := parseNullInt(e.ChildText("td[data-stat='pass_int']"))

			// Create stats record
			params := sqlc.CreatePlayerStatsParams{
				PlayerID:             playerID,
				GameID:               gameID,
				SeasonID:             seasonID,
				Week:                 week,
				PassingAttempts:      attempts,
				PassingCompletions:   completions,
				PassingYards:         yards,
				PassingTouchdowns:    tds,
				PassingInterceptions: ints,
			}

			_, err := s.DB.Queries.CreatePlayerStats(ctx, params)
			if err != nil {
				logger.Printf("Failed to insert passing stats for %s: %v", playerName, err)
				return
			}
			logger.Printf("Inserted passing stats for %s", playerName)

		case "rushing":
			attempts := parseNullInt(e.ChildText("td[data-stat='rush_att']"))
			yards := parseNullInt(e.ChildText("td[data-stat='rush_yds']"))
			tds := parseNullInt(e.ChildText("td[data-stat='rush_td']"))

			// Create stats record
			params := sqlc.CreatePlayerStatsParams{
				PlayerID:          playerID,
				GameID:            gameID,
				SeasonID:          seasonID,
				Week:              week,
				RushingAttempts:   attempts,
				RushingYards:      yards,
				RushingTouchdowns: tds,
			}

			_, err := s.DB.Queries.CreatePlayerStats(ctx, params)
			if err != nil {
				logger.Printf("Failed to insert rushing stats for %s: %v", playerName, err)
				return
			}
			logger.Printf("Inserted rushing stats for %s", playerName)

		case "receiving":
			targets := parseNullInt(e.ChildText("td[data-stat='targets']"))
			receptions := parseNullInt(e.ChildText("td[data-stat='rec']"))
			yards := parseNullInt(e.ChildText("td[data-stat='rec_yds']"))
			tds := parseNullInt(e.ChildText("td[data-stat='rec_td']"))

			// Create stats record
			params := sqlc.CreatePlayerStatsParams{
				PlayerID:            playerID,
				GameID:              gameID,
				SeasonID:            seasonID,
				Week:                week,
				Targets:             targets,
				Receptions:          receptions,
				ReceivingYards:      yards,
				ReceivingTouchdowns: tds,
			}

			_, err := s.DB.Queries.CreatePlayerStats(ctx, params)
			if err != nil {
				logger.Printf("Failed to insert receiving stats for %s: %v", playerName, err)
				return
			}
			logger.Printf("Inserted receiving stats for %s", playerName)

		case "kicking":
			fgMade := parseNullInt(e.ChildText("td[data-stat='fgm']"))
			fgAtt := parseNullInt(e.ChildText("td[data-stat='fga']"))
			xpMade := parseNullInt(e.ChildText("td[data-stat='xpm']"))
			xpAtt := parseNullInt(e.ChildText("td[data-stat='xpa']"))

			// Create stats record
			params := sqlc.CreatePlayerStatsParams{
				PlayerID:             playerID,
				GameID:               gameID,
				SeasonID:             seasonID,
				Week:                 week,
				FieldGoalsMade:       fgMade,
				FieldGoalsAttempted:  fgAtt,
				ExtraPointsMade:      xpMade,
				ExtraPointsAttempted: xpAtt,
			}

			_, err := s.DB.Queries.CreatePlayerStats(ctx, params)
			if err != nil {
				logger.Printf("Failed to insert kicking stats for %s: %v", playerName, err)
				return
			}
			logger.Printf("Inserted kicking stats for %s", playerName)

		case "defense":
			sacks := parseNullFloat(e.ChildText("td[data-stat='sacks']"))
			ints := parseNullInt(e.ChildText("td[data-stat='def_int']"))
			fumRec := parseNullInt(e.ChildText("td[data-stat='fumbles_rec']"))
			defTD := parseNullInt(e.ChildText("td[data-stat='def_td']"))

			// Create stats record
			params := sqlc.CreatePlayerStatsParams{
				PlayerID:            playerID,
				GameID:              gameID,
				SeasonID:            seasonID,
				Week:                week,
				Sacks:               sacks,
				Interceptions:       ints,
				FumbleRecoveries:    fumRec,
				DefensiveTouchdowns: defTD,
			}

			_, err := s.DB.Queries.CreatePlayerStats(ctx, params)
			if err != nil {
				logger.Printf("Failed to insert defense stats for %s: %v", playerName, err)
				return
			}
			logger.Printf("Inserted defense stats for %s", playerName)
		}

		statsCount++
	})

	// Log responses for debugging
	c.OnResponse(func(r *colly.Response) {
		logger.Printf("Received response from %s: status=%d, length=%d bytes",
			r.Request.URL, r.StatusCode, len(r.Body))

		// Check for table IDs
		htmlStr := string(r.Body)
		expectedTableID := fmt.Sprintf("table id=\"%s\"", statType)
		if strings.Contains(htmlStr, expectedTableID) {
			logger.Printf("Found %s in HTML", expectedTableID)
		} else {
			logger.Printf("NO %s found in HTML", expectedTableID)
			// Check for alternative table IDs
			alternateID := fmt.Sprintf("%s_advanced", statType)
			if strings.Contains(htmlStr, alternateID) {
				logger.Printf("Found alternative table ID: %s", alternateID)
			}
		}
	})

	// Request callbacks
	c.OnRequest(func(r *colly.Request) {
		logger.Printf("Making HTTP request to: %s", r.URL.String())
	})

	c.OnError(func(r *colly.Response, err error) {
		logger.Printf("Error when visiting %s: %v", r.Request.URL, err)
		logger.Printf("Response status code: %d", r.StatusCode)
	})

	c.OnScraped(func(r *colly.Response) {
		logger.Printf("Finished scraping: %s", r.Request.URL)
		logger.Printf("Total %s stats records: %d", statType, statsCount)
	})

	// Start scraping
	err = c.Visit(url)
	if err != nil {
		return fmt.Errorf("failed to visit URL %s: %w", url, err)
	}

	return nil
}

// Helper functions for parsing stats
func parseNullInt(s string) sql.NullInt64 {
	if s == "" {
		return sql.NullInt64{Valid: false}
	}

	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return sql.NullInt64{Valid: false}
	}

	return sql.NullInt64{Int64: val, Valid: true}
}

func parseNullFloat(s string) sql.NullFloat64 {
	if s == "" {
		return sql.NullFloat64{Valid: false}
	}

	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return sql.NullFloat64{Valid: false}
	}

	return sql.NullFloat64{Float64: val, Valid: true}
}

// splitTeamName splits a full team name into city and nickname
func splitTeamName(fullName string) (city, nickname string) {
	// Handle special cases first
	if strings.Contains(fullName, "Washington Commanders") {
		return "Washington", "Commanders"
	} else if strings.Contains(fullName, "Tampa Bay") {
		return "Tampa Bay", "Buccaneers"
	} else if strings.Contains(fullName, "Green Bay") {
		return "Green Bay", "Packers"
	} else if strings.Contains(fullName, "New England") {
		return "New England", "Patriots"
	} else if strings.Contains(fullName, "New Orleans") {
		return "New Orleans", "Saints"
	} else if strings.Contains(fullName, "New York Giants") {
		return "New York", "Giants"
	} else if strings.Contains(fullName, "New York Jets") {
		return "New York", "Jets"
	} else if strings.Contains(fullName, "Las Vegas") {
		return "Las Vegas", "Raiders"
	} else if strings.Contains(fullName, "Los Angeles Rams") {
		return "Los Angeles", "Rams"
	} else if strings.Contains(fullName, "Los Angeles Chargers") {
		return "Los Angeles", "Chargers"
	} else if strings.Contains(fullName, "Kansas City") {
		return "Kansas City", "Chiefs"
	} else if strings.Contains(fullName, "San Francisco") {
		return "San Francisco", "49ers"
	}

	// For other teams, split on the last space
	parts := strings.Split(fullName, " ")
	if len(parts) == 1 {
		return "", parts[0] // Just nickname
	}

	nickname = parts[len(parts)-1]
	city = strings.Join(parts[:len(parts)-1], " ")
	return city, nickname
}

// getTeamConferenceAndDivision returns the conference and division for a team
func getTeamConferenceAndDivision(abbr string) (conference, division string) {
	// This is based on current NFL alignments as of 2024
	switch abbr {
	// AFC East
	case "buf", "mia", "nwe", "nyj":
		return "AFC", "East"
	// AFC North
	case "bal", "cin", "cle", "pit", "rav": // Added "rav" for Ravens
		return "AFC", "North"
	// AFC South
	case "hou", "ind", "jax", "ten", "htx", "clt", "oti": // Added "htx", "clt", "oti" for Texans, Colts, and Titans
		return "AFC", "South"
	// AFC West
	case "den", "kan", "lac", "lvr", "rai", "sdg", "oak": // Added "lvr" and "lac" for Raiders and Chargers
		return "AFC", "West"
	// NFC East
	case "dal", "nyg", "phi", "was": // Updated "was" for Washington
		return "NFC", "East"
	// NFC North
	case "chi", "det", "gnb", "min":
		return "NFC", "North"
	// NFC South
	case "atl", "car", "nor", "tam", "tb":
		return "NFC", "South"
	// NFC West
	case "ari", "ram", "lar", "sea", "sfo", "crd": // Added "crd" for Arizona Cardinals
		return "NFC", "West"
	default:
		return "Unknown", "Unknown"
	}
}

// Helper function for string slicing
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

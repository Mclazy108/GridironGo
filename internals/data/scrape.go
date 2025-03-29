package data

import (
	"context"
	"fmt"
	"log"
	"net/url"
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

	// Define the target URL for scraping team info
	targetURL := "https://www.pro-football-reference.com/teams/"
	log.Printf("Target URL for team scraping: %s", targetURL)

	// Parse the URL to verify it's valid
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		log.Printf("WARNING: URL parsing error: %v", err)
	} else {
		log.Printf("Parsed URL - Scheme: %s, Host: %s, Path: %s",
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
	log.Println("Initialized team collector")

	// Process the team links directly - this is much simpler and more reliable
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		href := e.Attr("href")

		// Only process direct team links - they follow the pattern /teams/XXX/
		// But exclude any links with years in them like /teams/buf/2024.htm
		if strings.HasPrefix(href, "/teams/") && !strings.Contains(href, ".htm") && len(href) > 7 {
			teamName := e.Text
			log.Printf("Processing team link: %s (%s)", teamName, href)

			// Extract abbreviation from the href
			teamAbbr := ""
			path := strings.TrimSuffix(href, "/")
			path = strings.TrimPrefix(path, "/teams/")
			if path != "" && path != "index" {
				teamAbbr = path
				log.Printf("  Team abbreviation: %s", teamAbbr)

				// Split the team name into city and nickname
				cityName, teamNickname := splitTeamName(teamName)
				log.Printf("  City: %s, Team name: %s", cityName, teamNickname)

				// Determine conference and division
				conference, division := getTeamConferenceAndDivision(teamAbbr)
				log.Printf("  Conference: %s, Division: %s", conference, division)

				// Store in our map
				teams[teamAbbr] = Team{
					Name:         teamNickname,
					City:         cityName,
					Abbreviation: teamAbbr,
					Conference:   conference,
					Division:     division,
				}

				log.Printf("  Team %s (%s %s) added to collection", teamAbbr, cityName, teamNickname)
			}
		}
	})

	c.OnRequest(func(r *colly.Request) {
		log.Printf("Making HTTP request to: %s", r.URL.String())
	})

	c.OnResponse(func(r *colly.Response) {
		log.Printf("Received response from %s: status=%d, length=%d bytes",
			r.Request.URL, r.StatusCode, len(r.Body))
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Error when visiting %s: %v", r.Request.URL, err)
		log.Printf("Response status code: %d", r.StatusCode)
	})

	c.OnScraped(func(r *colly.Response) {
		log.Printf("Finished scraping: %s", r.Request.URL)
		log.Printf("Total teams found: %d", len(teams))
	})

	// Start scraping
	log.Println("Starting to visit teams page...")
	err = c.Visit(targetURL)
	if err != nil {
		return fmt.Errorf("failed to visit teams page: %w", err)
	}

	// Log all found teams
	log.Println("Teams found during scraping:")
	for abbr, team := range teams {
		log.Printf("- %s: %s %s (%s, %s)",
			abbr, team.City, team.Name, team.Conference, team.Division)
	}

	log.Println("Scraping completed. Now inserting teams into database...")

	// Insert teams into the database
	for abbr, team := range teams {
		log.Printf("Processing team: %s %s (%s)", team.City, team.Name, abbr)

		// Check if team already exists
		existingTeam, err := s.DB.Queries.GetNFLTeamByAbbreviation(ctx, team.Abbreviation)
		if err == nil {
			log.Printf("Team %s already exists in database with ID %d, skipping...",
				team.Abbreviation, existingTeam.ID)
			continue
		} else {
			log.Printf("Team %s not found in database (error: %v), will insert",
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

		log.Printf("Inserting team with params: %+v", params)

		id, err := s.DB.Queries.InsertTeam(ctx, params)
		if err != nil {
			log.Printf("ERROR: Failed to insert team %s: %v", team.Name, err)
		} else {
			log.Printf("SUCCESS: Inserted team: %s %s with ID %d", team.City, team.Name, id)
		}
	}

	log.Println("Team scraping and database insertion completed")
	return nil
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

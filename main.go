package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Mclazy108/GridironGo/internals/data"
	"github.com/Mclazy108/GridironGo/internals/data/scraper"
)

func main() {
	// Define command-line flags
	scrapeGames := flag.Bool("scrape-games", false, "Scrape NFL game data for seasons 2022-2024")
	scrapeTeams := flag.Bool("scrape-teams", false, "Scrape NFL team data")
	scrapePlayers := flag.Bool("scrape-players", false, "Scrape NFL player data")
	scrapeStats := flag.Bool("scrape-stats", false, "Scrape NFL game statistics")
	dbPath := flag.String("db", "./GridironGo.db", "Path to SQLite database (default: ./GridironGo.db)")

	// Add specific season flags
	seasons := flag.String("seasons", "2022,2023,2024", "Comma-separated list of seasons to scrape games for")

	// Track durations for summary
	var (
		gameDuration   time.Duration
		teamDuration   time.Duration
		playerDuration time.Duration
		statDuration   time.Duration
	)
	// Parse flags
	flag.Parse()

	// Check and log database path
	if *dbPath == "" {
		*dbPath = "./GridironGo.db"
		log.Printf("Empty database path provided, using default: %s", *dbPath)
	} else {
		log.Printf("Using database path: %s", *dbPath)
	}

	// Create database connection
	log.Println("Initializing database connection...")
	db, err := data.NewDB(&data.DBConfig{
		Path: *dbPath,
	})
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()
	log.Println("Database connection established")

	// Create a context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling for graceful cancellation
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Handle signals in a separate goroutine
	go func() {
		sig := <-sigCh
		log.Printf("Received signal %v, gracefully shutting down...", sig)
		cancel() // This will cancel the context

		// If a second signal is received, force exit
		sig = <-sigCh
		log.Printf("Received second signal %v, forcing immediate exit", sig)
		os.Exit(1)
	}()

	// Check if no specific scraping flags were provided
	runDefaultScraping := !*scrapeGames && !*scrapeTeams && !*scrapePlayers && !*scrapeStats && len(flag.Args()) == 0

	// If running default scraping, run them in the correct order
	if runDefaultScraping {
		// Run game scraping first
		start := time.Now()
		err := runGameScraper(ctx, db, *seasons)
		gameDuration = time.Since(start)
		if err != nil {
			log.Printf("Error during game scraping: %v", err)
		}

		// Run team scraping second
		start = time.Now()
		err = runTeamScraper(ctx, db)
		teamDuration = time.Since(start)
		if err != nil {
			log.Printf("Error during team scraping: %v", err)
		}

		// Run player scraping third
		start = time.Now()
		err = runPlayerScraper(ctx, db, *seasons) // Pass seasons to player scraper
		playerDuration = time.Since(start)
		if err != nil {
			log.Printf("Error during player scraping: %v", err)
		}

		// Run stats scraping last
		start = time.Now()
		err = runStatScraper(ctx, db, *seasons)
		statDuration = time.Since(start)
		if err != nil {
			log.Printf("Error during stat scraping: %v", err)
		}
	} else {
		// Otherwise, run only the explicitly requested scrapers
		// (in the correct order when multiple are selected)

		// Run game scraping first if requested
		if *scrapeGames {
			start := time.Now()
			err := runGameScraper(ctx, db, *seasons)
			gameDuration = time.Since(start)
			if err != nil {
				log.Printf("Error during game scraping: %v", err)
			}
		}

		// Run team scraping second if requested
		if *scrapeTeams {
			start := time.Now()
			err := runTeamScraper(ctx, db)
			teamDuration = time.Since(start)
			if err != nil {
				log.Printf("Error during team scraping: %v", err)
			}
		}

		// Run player scraping third if requested
		if *scrapePlayers {
			start := time.Now()
			err = runPlayerScraper(ctx, db, *seasons) // Pass seasons to player scraper
			playerDuration = time.Since(start)
			if err != nil {
				log.Printf("Error during player scraping: %v", err)
			}
		}

		// Run stats scraping last if requested
		if *scrapeStats {
			start := time.Now()
			err := runStatScraper(ctx, db, *seasons)
			statDuration = time.Since(start)
			if err != nil {
				log.Printf("Error during stat scraping: %v", err)
			}
		}
	}

	// If specific scraping flags were provided or default scraping was run, exit
	if *scrapeGames || *scrapeTeams || *scrapePlayers || *scrapeStats || runDefaultScraping {
		// Get record counts for summary
		gameCount, _ := getGameCount(ctx, db)
		teamCount, _ := getTeamCount(ctx, db)
		playerCount, _ := getPlayerCount(ctx, db)
		playerSeasonCount, _ := getPlayerSeasonCount(ctx, db)
		statCount, _ := getStatCount(ctx, db)

		log.Println("------------------------------------------------")
		log.Println("🏁 Scraping Summary:")
		if *scrapeGames || runDefaultScraping {
			log.Printf("⏱  Games scraped in:   %s (Total records: %d)", gameDuration, gameCount)
		}
		if *scrapeTeams || runDefaultScraping {
			log.Printf("⏱  Teams scraped in:   %s (Total records: %d)", teamDuration, teamCount)
		}
		if *scrapePlayers || runDefaultScraping {
			log.Printf("⏱  Players scraped in: %s (Total players: %d, Total player-seasons: %d)",
				playerDuration, playerCount, playerSeasonCount)
		}
		if *scrapeStats || runDefaultScraping {
			log.Printf("⏱  Stats scraped in:   %s (Total records: %d)", statDuration, statCount)
		}
		log.Println("------------------------------------------------")
		return
	}

	// Otherwise, start the TUI application
	fmt.Println("GridironGo - Fantasy Football CLI App")
	fmt.Println("Starting the application...")
	// Here you would normally start your TUI application
}

// Get count of records in the nfl_games table
func getGameCount(ctx context.Context, db *data.DB) (int, error) {
	var count int
	err := db.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM nfl_games").Scan(&count)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	return count, nil
}

// Get count of records in the nfl_teams table
func getTeamCount(ctx context.Context, db *data.DB) (int, error) {
	var count int
	err := db.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM nfl_teams").Scan(&count)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	return count, nil
}

// Get count of records in the nfl_players table
func getPlayerCount(ctx context.Context, db *data.DB) (int, error) {
	var count int
	err := db.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM nfl_players").Scan(&count)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	return count, nil
}

// Get count of records in the nfl_player_seasons table
func getPlayerSeasonCount(ctx context.Context, db *data.DB) (int, error) {
	var count int
	err := db.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM nfl_player_seasons").Scan(&count)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	return count, nil
}

// Get count of records in the nfl_stats table
func getStatCount(ctx context.Context, db *data.DB) (int, error) {
	var count int
	err := db.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM nfl_stats").Scan(&count)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	return count, nil
}

// Parse comma-separated seasons string into slice of integers
func parseSeasons(seasonsStr string) []int {
	var seasonsInt []int
	var currentNum int

	for i := 0; i < len(seasonsStr); i++ {
		c := seasonsStr[i]

		// If we find a digit, process it
		if c >= '0' && c <= '9' {
			currentNum = currentNum*10 + int(c-'0')
		} else if c == ',' {
			// Add the current number to our list and reset
			if currentNum > 0 {
				seasonsInt = append(seasonsInt, currentNum)
				currentNum = 0
			}
		}
	}

	// Don't forget the last number if there is one
	if currentNum > 0 {
		seasonsInt = append(seasonsInt, currentNum)
	}

	// Default to 2022-2024 if no valid seasons were provided
	if len(seasonsInt) == 0 {
		return []int{2022, 2023, 2024}
	}

	return seasonsInt
}

// runGameScraper handles the game scraping process
func runGameScraper(ctx context.Context, db *data.DB, seasonsStr string) error {
	log.Println("Starting NFL game data scraping...")
	log.Println("Press Ctrl+C for graceful cancellation")

	seasons := parseSeasons(seasonsStr)
	log.Printf("Will scrape games for seasons: %v", seasons)

	scraperInstance := scraper.NewScraper(db)

	// Count games before scraping
	gameCount, err := getGameCount(ctx, db)
	if err != nil {
		log.Printf("Warning: Could not get existing game count: %v", err)
	} else {
		log.Printf("Found %d existing games in database before scraping", gameCount)
	}

	// Perform scraping with cancellable context
	err = scraperInstance.ScrapeNFLGames(ctx, seasons)

	// Check if the operation was cancelled by the user
	if ctx.Err() != nil {
		log.Println("Scraping was cancelled by the user")
		return ctx.Err()
	}

	if err != nil {
		return fmt.Errorf("error scraping NFL games: %w", err)
	}

	// Count games after scraping
	gameCount, err = getGameCount(ctx, db)
	if err != nil {
		log.Printf("Warning: Could not get updated game count: %v", err)
	} else {
		log.Printf("Database now contains %d games after scraping", gameCount)
	}

	// Report success
	log.Println("NFL game data scraping completed successfully")
	return nil
}

// runTeamScraper handles the team scraping process
func runTeamScraper(ctx context.Context, db *data.DB) error {
	log.Println("Starting NFL team data scraping...")
	log.Println("Press Ctrl+C for graceful cancellation")

	teamScraperInstance := scraper.NewTeamScraper(db)

	// Count teams before scraping
	teamCount, err := getTeamCount(ctx, db)
	if err != nil {
		log.Printf("Warning: Could not get existing team count: %v", err)
	} else {
		log.Printf("Found %d existing teams in database before scraping", teamCount)
	}

	// Perform scraping with cancellable context
	err = teamScraperInstance.ScrapeNFLTeams(ctx)

	// Check if the operation was cancelled by the user
	if ctx.Err() != nil {
		log.Println("Team scraping was cancelled by the user")
		return ctx.Err()
	}

	if err != nil {
		return fmt.Errorf("error scraping NFL teams: %w", err)
	}

	// Count teams after scraping
	teamCount, err = getTeamCount(ctx, db)
	if err != nil {
		log.Printf("Warning: Could not get updated team count: %v", err)
	} else {
		log.Printf("Database now contains %d teams after scraping", teamCount)
	}

	// Report success
	log.Println("NFL team data scraping completed successfully")
	return nil
}

// runPlayerScraper handles the player scraping process
func runPlayerScraper(ctx context.Context, db *data.DB, seasonsStr string) error {
	log.Println("Starting NFL player data scraping...")
	log.Println("Press Ctrl+C for graceful cancellation")

	// Parse seasons
	seasons := parseSeasons(seasonsStr)
	log.Printf("Will scrape players for seasons: %v", seasons)

	playerScraperInstance := scraper.NewPlayerScraper(db)

	// Count players before scraping
	playerCount, err := getPlayerCount(ctx, db)
	if err != nil {
		log.Printf("Warning: Could not get existing player count: %v", err)
	} else {
		log.Printf("Found %d existing players in database before scraping", playerCount)
	}

	// Count player seasons before scraping
	seasonCount, err := getPlayerSeasonCount(ctx, db)
	if err != nil {
		log.Printf("Warning: Could not get existing player season count: %v", err)
	} else {
		log.Printf("Found %d existing player-season records in database before scraping", seasonCount)
	}

	// Perform scraping with cancellable context and specified seasons
	err = playerScraperInstance.ScrapeNFLPlayers(ctx, seasons)

	// Check if the operation was cancelled by the user
	if ctx.Err() != nil {
		log.Println("Player scraping was cancelled by the user")
		return ctx.Err()
	}

	if err != nil {
		return fmt.Errorf("error scraping NFL players: %w", err)
	}

	// Count players after scraping
	playerCount, err = getPlayerCount(ctx, db)
	if err != nil {
		log.Printf("Warning: Could not get updated player count: %v", err)
	} else {
		log.Printf("Database now contains %d unique players after scraping", playerCount)
	}

	// Count player seasons after scraping
	seasonCount, err = getPlayerSeasonCount(ctx, db)
	if err != nil {
		log.Printf("Warning: Could not get updated player season count: %v", err)
	} else {
		log.Printf("Database now contains %d player-season records after scraping", seasonCount)
	}

	// Report success
	log.Println("NFL player data scraping completed successfully")
	return nil
}

// runStatScraper handles the game statistics scraping process
func runStatScraper(ctx context.Context, db *data.DB, seasonsStr string) error {
	log.Println("Starting NFL game statistics scraping...")
	log.Println("Press Ctrl+C for graceful cancellation")

	seasons := parseSeasons(seasonsStr)
	log.Printf("Will scrape statistics for seasons: %v", seasons)

	statScraperInstance := scraper.NewStatScraper(db)

	// Count stats before scraping
	statCount, err := getStatCount(ctx, db)
	if err != nil {
		log.Printf("Warning: Could not get existing stats count: %v", err)
	} else {
		log.Printf("Found %d existing statistics in database before scraping", statCount)
	}

	// Perform scraping with cancellable context
	err = statScraperInstance.ScrapeNFLGameStats(ctx, seasons)

	// Check if the operation was cancelled by the user
	if ctx.Err() != nil {
		log.Println("Scraping was cancelled by the user")
		return ctx.Err()
	}

	if err != nil {
		return fmt.Errorf("error scraping NFL game statistics: %w", err)
	}

	// Count stats after scraping
	statCount, err = getStatCount(ctx, db)
	if err != nil {
		log.Printf("Warning: Could not get updated stats count: %v", err)
	} else {
		log.Printf("Database now contains %d statistics after scraping", statCount)
	}

	// Report success
	log.Println("NFL game statistics scraping completed successfully")
	return nil
}


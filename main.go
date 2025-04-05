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

	"github.com/Mclazy108/GridironGo/internals/data"
	"github.com/Mclazy108/GridironGo/internals/data/scraper"
)

func main() {
	// Define command-line flags
	scrapeGames := flag.Bool("scrape-games", false, "Scrape NFL game data for seasons 2022-2024")
	scrapeTeams := flag.Bool("scrape-teams", false, "Scrape NFL team data")
	scrapePlayers := flag.Bool("scrape-players", false, "Scrape NFL player data")
	dbPath := flag.String("db", "./GridironGo.db", "Path to SQLite database (default: ./GridironGo.db)")

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
	runDefaultScraping := !*scrapeGames && !*scrapeTeams && !*scrapePlayers && len(flag.Args()) == 0

	// Run game scraping if explicitly requested or running default
	if *scrapeGames || runDefaultScraping {
		err := runGameScraper(ctx, db)
		if err != nil {
			log.Printf("Error during game scraping: %v", err)
		}
	}

	// Run team scraping if explicitly requested or running default
	if *scrapeTeams || runDefaultScraping {
		err := runTeamScraper(ctx, db)
		if err != nil {
			log.Printf("Error during team scraping: %v", err)
		}
	}

	// Run player scraping if explicitly requested or running default
	if *scrapePlayers || runDefaultScraping {
		err := runPlayerScraper(ctx, db)
		if err != nil {
			log.Printf("Error during player scraping: %v", err)
		}
	}

	// If specific scraping flags were provided or default scraping was run, exit
	if *scrapeGames || *scrapeTeams || *scrapePlayers || runDefaultScraping {
		return
	}

	// Otherwise, start the TUI application
	fmt.Println("GridironGo - Fantasy Football CLI App")
	fmt.Println("Starting the application...")
	// Here you would normally start your TUI application
}

// runGameScraper handles the game scraping process
func runGameScraper(ctx context.Context, db *data.DB) error {
	log.Println("Starting NFL game data scraping...")
	log.Println("Press Ctrl+C for graceful cancellation")

	scraperInstance := scraper.NewScraper(db)

	// Scrape games for seasons 2022-2024
	seasons := []int{2022, 2023, 2024}

	// Count games before scraping
	games, err := db.Queries.GetAllGames(ctx)
	if err != nil {
		log.Printf("Warning: Could not get existing game count: %v", err)
	} else {
		log.Printf("Found %d existing games in database before scraping", len(games))
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
	games, err = db.Queries.GetAllGames(ctx)
	if err != nil {
		log.Printf("Warning: Could not get updated game count: %v", err)
	} else {
		log.Printf("Database now contains %d games after scraping", len(games))
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
	teams, err := db.Queries.GetAllNFLTeams(ctx)
	if err != nil {
		log.Printf("Warning: Could not get existing team count: %v", err)
	} else {
		log.Printf("Found %d existing teams in database before scraping", len(teams))
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
	teams, err = db.Queries.GetAllNFLTeams(ctx)
	if err != nil {
		log.Printf("Warning: Could not get updated team count: %v", err)
	} else {
		log.Printf("Database now contains %d teams after scraping", len(teams))
	}

	// Report success
	log.Println("NFL team data scraping completed successfully")
	return nil
}

// runPlayerScraper handles the player scraping process
func runPlayerScraper(ctx context.Context, db *data.DB) error {
	log.Println("Starting NFL player data scraping...")
	log.Println("Press Ctrl+C for graceful cancellation")

	playerScraperInstance := scraper.NewPlayerScraper(db)

	// Count players before scraping
	players, err := db.Queries.GetAllNFLPlayers(ctx)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Warning: Could not get existing player count: %v", err)
	} else {
		log.Printf("Found %d existing players in database before scraping", len(players))
	}

	// Perform scraping with cancellable context
	err = playerScraperInstance.ScrapeNFLPlayers(ctx)

	// Check if the operation was cancelled by the user
	if ctx.Err() != nil {
		log.Println("Player scraping was cancelled by the user")
		return ctx.Err()
	}

	if err != nil {
		return fmt.Errorf("error scraping NFL players: %w", err)
	}

	// Count players after scraping
	players, err = db.Queries.GetAllNFLPlayers(ctx)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Warning: Could not get updated player count: %v", err)
	} else {
		log.Printf("Database now contains %d players after scraping", len(players))
	}

	// Report success
	log.Println("NFL player data scraping completed successfully")
	return nil
}

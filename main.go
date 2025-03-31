package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Mclazy108/GridironGo/internals/data"
)

func main() {
	// Define command-line flags
	scrapeGames := flag.Bool("scrape-games", false, "Scrape NFL game data for seasons 2022-2024")
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

	// If scraping games is requested, do that and exit
	if *scrapeGames {
		log.Println("Starting NFL game data scraping...")
		log.Println("Press Ctrl+C for graceful cancellation")

		// Create scraper
		scraper := data.NewScraper(db)

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
		err = scraper.ScrapeNFLGames(ctx, seasons)

		// Check if the operation was cancelled by the user
		if ctx.Err() != nil {
			log.Println("Scraping was cancelled by the user")
			return
		}

		if err != nil {
			log.Fatalf("Error scraping NFL games: %v", err)
		}

		// Count games after scraping
		games, err = db.Queries.GetAllGames(ctx)
		if err != nil {
			log.Printf("Warning: Could not get updated game count: %v", err)
		} else {
			log.Printf("Database now contains %d games after scraping", len(games))
		}

		// Report success and exit
		log.Println("NFL game data scraping completed successfully")
		return
	}

	// Here you would normally start your TUI application
	// For now, just print a message if no actions were specified
	fmt.Println("GridironGo - Fantasy Football CLI App")
	fmt.Println("Use -h flag to see available options")
}


package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Mclazy108/GridironGo/internals/data"
)

func main() {
	log.Println("Starting GridironGo...")

	// Set up database connection
	log.Println("Connecting to database...")
	db, err := data.NewDB(nil) // Use default configuration
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Database connection established successfully")

	// Test query - get current season if it exists
	ctx := context.Background()
	season, err := db.Queries.GetCurrentSeason(ctx)
	if err != nil {
		log.Printf("No current season found: %v", err)
	} else {
		log.Printf("Current season: %d", season.Year)
	}

	// Test query - get all NFL teams count
	teams, err := db.Queries.GetAllNFLTeams(ctx)
	if err != nil {
		log.Printf("Error fetching teams: %v", err)
	} else {
		log.Printf("Found %d NFL teams in database", len(teams))
	}

	// Wait for interrupt signal (Ctrl+C) to gracefully shut down
	log.Println("Application running. Press Ctrl+C to exit.")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	fmt.Println("\nShutting down...")
}

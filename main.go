package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
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

	// Create a scraper
	scraper := data.NewScraper(db)

	// Start scraping
	ctx := context.Background()

	// Command line arguments
	if len(os.Args) > 1 {
		cmd := os.Args[1]
		switch cmd {
		case "scrape-teams":
			// Scrape NFL teams
			log.Println("Scraping NFL teams...")
			err = scraper.ScrapeTeams(ctx)
			if err != nil {
				log.Fatalf("Failed to scrape teams: %v", err)
			}
			log.Println("Team scraping completed")
			return

		case "scrape-players":
			// First ensure we have teams
			teams, err := db.Queries.GetAllNFLTeams(ctx)
			if err != nil || len(teams) < 32 {
				log.Println("Teams data appears incomplete. Scraping teams first...")
				err = scraper.ScrapeTeams(ctx)
				if err != nil {
					log.Fatalf("Failed to scrape teams: %v", err)
				}
				log.Println("Team scraping completed")
			} else {
				log.Printf("Found %d NFL teams in database, proceeding with player scraping", len(teams))
			}

			// Default seasons to scrape if none specified
			seasons := []int{2022, 2023, 2024}

			// Check if specific seasons are provided
			if len(os.Args) > 2 {
				seasons = []int{}
				for i := 2; i < len(os.Args); i++ {
					year, err := strconv.Atoi(os.Args[i])
					if err != nil {
						log.Fatalf("Invalid season year: %s", os.Args[i])
					}
					seasons = append(seasons, year)
				}
			}

			log.Printf("Scraping players for seasons: %v", seasons)
			err = scraper.ScrapePlayers(ctx, seasons)
			if err != nil {
				log.Fatalf("Failed to scrape players: %v", err)
			}
			log.Println("Player scraping completed")
			return

		case "scrape-stats":
			// First ensure we have teams and players
			teams, err := db.Queries.GetAllNFLTeams(ctx)
			if err != nil || len(teams) < 32 {
				log.Println("Teams data appears incomplete. Scraping teams first...")
				err = scraper.ScrapeTeams(ctx)
				if err != nil {
					log.Fatalf("Failed to scrape teams: %v", err)
				}
				log.Println("Team scraping completed")
			} else {
				log.Printf("Found %d NFL teams in database", len(teams))
			}

			players, err := db.Queries.GetAllPlayers(ctx)
			if err != nil || len(players) < 100 {
				log.Println("Players data appears incomplete. Scraping players first...")
				// Default seasons to scrape
				playerSeasons := []int{2022, 2023, 2024}
				err = scraper.ScrapePlayers(ctx, playerSeasons)
				if err != nil {
					log.Fatalf("Failed to scrape players: %v", err)
				}
				log.Println("Player scraping completed")
			} else {
				log.Printf("Found %d players in database", len(players))
			}

			// Default seasons to scrape if none specified
			seasons := []int{2022, 2023, 2024}

			// Check if specific seasons are provided
			if len(os.Args) > 2 {
				seasons = []int{}
				for i := 2; i < len(os.Args); i++ {
					year, err := strconv.Atoi(os.Args[i])
					if err != nil {
						log.Fatalf("Invalid season year: %s", os.Args[i])
					}
					seasons = append(seasons, year)
				}
			}

			log.Printf("Scraping player stats for seasons: %v", seasons)
			err = scraper.ScrapePlayerStats(ctx, seasons)
			if err != nil {
				log.Fatalf("Failed to scrape player stats: %v", err)
			}
			log.Println("Player stats scraping completed")
			return

		case "scrape-all":
			// Default seasons to scrape if none specified
			seasons := []int{2022, 2023, 2024}

			// Check if specific seasons are provided
			if len(os.Args) > 2 {
				seasons = []int{}
				for i := 2; i < len(os.Args); i++ {
					year, err := strconv.Atoi(os.Args[i])
					if err != nil {
						log.Fatalf("Invalid season year: %s", os.Args[i])
					}
					seasons = append(seasons, year)
				}
			}

			// Scrape teams first
			log.Println("Scraping NFL teams...")
			err = scraper.ScrapeTeams(ctx)
			if err != nil {
				log.Fatalf("Failed to scrape teams: %v", err)
			}
			log.Println("Team scraping completed")

			// Then scrape players
			log.Printf("Scraping players for seasons: %v", seasons)
			err = scraper.ScrapePlayers(ctx, seasons)
			if err != nil {
				log.Fatalf("Failed to scrape players: %v", err)
			}
			log.Println("Player scraping completed")

			// Finally scrape stats
			log.Printf("Scraping player stats for seasons: %v", seasons)
			err = scraper.ScrapePlayerStats(ctx, seasons)
			if err != nil {
				log.Fatalf("Failed to scrape player stats: %v", err)
			}
			log.Println("Player stats scraping completed")

			log.Println("All scraping operations completed successfully")
			return
		}
	}

	// Test query - get current season if it exists
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

	// Test query - get player count
	players, err := db.Queries.GetAllPlayers(ctx)
	if err != nil {
		log.Printf("Error fetching players: %v", err)
	} else {
		log.Printf("Found %d players in database", len(players))
	}

	// Wait for interrupt signal (Ctrl+C) to gracefully shut down
	log.Println("Application running. Press Ctrl+C to exit.")
	log.Println("Available commands:")
	log.Println("  - scrape-teams: Scrape NFL teams data")
	log.Println("  - scrape-players [year...]: Scrape NFL players for specified seasons (default: 2022, 2023, 2024)")
	log.Println("  - scrape-stats [year...]: Scrape NFL player stats for specified seasons (default: 2022, 2023, 2024)")
	log.Println("  - scrape-all [year...]: Scrape all NFL data for specified seasons (default: 2022, 2023, 2024)")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	fmt.Println("\nShutting down...")
}


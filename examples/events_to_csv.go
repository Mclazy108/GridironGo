package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

type ScoreboardResponse struct {
	Events []Event `json:"events"`
}

type Event struct {
	ID        string `json:"id"`
	Date      string `json:"date"`
	Name      string `json:"name"`
	ShortName string `json:"shortName"`
	Season    Season `json:"season"`
	Week      Week   `json:"week"`
}

type Season struct {
	Year int `json:"year"`
}

type Week struct {
	Number int `json:"number"`
}

func fetchEvents(year int, week int) ([]Event, error) {
	// Construct the API URL
	url := fmt.Sprintf("https://site.api.espn.com/apis/site/v2/sports/football/nfl/scoreboard?dates=%d&seasontype=2&week=%d", year, week)

	// Send HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse the JSON response
	var scoreboardResponse ScoreboardResponse
	err = json.NewDecoder(resp.Body).Decode(&scoreboardResponse)
	if err != nil {
		return nil, err
	}

	return scoreboardResponse.Events, nil
}

func saveEventsToCSV(events []Event, fileName string) error {
	// Create CSV file
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	err = writer.Write([]string{"EventID", "Date", "Name", "ShortName", "SeasonYear", "WeekNumber", "AwayTeam", "HomeTeam"})
	if err != nil {
		return err
	}

	// Write event details to CSV
	for _, event := range events {
		// Extract Home Team and Away Team from the Name field
		awayTeam, homeTeam := extractTeams(event.Name)

		err := writer.Write([]string{
			event.ID,
			event.Date,
			event.Name,
			event.ShortName,
			fmt.Sprintf("%d", event.Season.Year),
			fmt.Sprintf("%d", event.Week.Number),
			awayTeam,
			homeTeam,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// Function to extract the Away Team and Home Team from the game name
func extractTeams(gameName string) (string, string) {
	// Split the name by " at " to separate away and home teams
	teams := strings.Split(gameName, " at ")
	if len(teams) == 2 {
		awayTeam := teams[0]
		homeTeam := teams[1]
		return awayTeam, homeTeam
	}
	return "", "" // Return empty strings if the format doesn't match
}

func main() {
	var allEvents []Event

	// Loop through years 2022 to 2024
	for year := 2022; year <= 2024; year++ {
		// Loop through weeks 1 to 18
		for week := 1; week <= 18; week++ {
			fmt.Printf("Fetching data for Week %d, Year %d...\n", week, year)
			events, err := fetchEvents(year, week)
			if err != nil {
				log.Printf("Error fetching data for Week %d, Year %d: %v", week, year, err)
				continue
			}
			allEvents = append(allEvents, events...)
		}
	}

	// Save all events to CSV file
	err := saveEventsToCSV(allEvents, "nfl_events_2022_2024.csv")
	if err != nil {
		log.Fatalf("Error saving data to CSV: %v", err)
	}

	fmt.Println("Data saved to nfl_events_2022_2024.csv")
}


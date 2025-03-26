# GridironGo

GridironGo is a Fantasy Football CLI app built in Go using the [Bubble Tea](https://github.com/charmbracelet/bubbletea) TUI framework.

## Project Structure

```
.
├── assets                      # Directory for static assets like JSON data, logos, or templates
├── go.mod                      # Defines the module name and dependencies for the Go project
├── go.sum                      # Locks dependency versions to ensure reproducible builds
├── internals                   # Contains core application logic split into sub-packages
│   ├── data
│   │   ├── db.go               # Handles SQLite database connections and queries
│   │   ├── models.go           # Defines the data models for players, teams, and matches
│   │   └── scrape.go           # Scrapes football stats and schedules for the last 3 seasons
│   ├── league
│   │   ├── league.go           # Manages fantasy league setup and operations
│   │   ├── rules.go            # Handles league rules including scoring and configurations
│   │   ├── schedule.go         # Generates and manages league schedules, including playoffs
│   │   └── team.go             # Manages fantasy teams including bot teams and user team
│   ├── rules
│   │   ├── draft.go            # Handles the drafting logic and player selection process
│   │   └── scoring.go          # Implements fantasy football scoring rules and calculations
│   └── tui
│       ├── league_menu.go      # TUI logic for the fantasy league menu and its options
│       ├── menu.go             # Main TUI entry point with initial menu options
│       ├── player_menu.go      # TUI logic for viewing players and selecting them
│       └── schedule_menu.go    # TUI logic for viewing the real and fantasy schedules
├── LICENSE                     # Project license file
├── main.go                     # Entry point for the application. Initializes TUI and loads data
└── README.md                   # Project documentation and setup instructions
```

## Features
- Create and manage a Fantasy Football League with bot and user teams
- View and draft players using previous season stats
- Customize league rules including PPR, scoring, and roster settings
- Full TUI experience using Bubble Tea
- SQLite backend for storing player stats and league data

## Getting Started
1. Clone the repo
2. Run `go mod tidy` to install dependencies
3. Run `go run main.go` to launch the TUI

## License
MIT

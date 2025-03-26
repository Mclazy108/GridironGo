# GridironGo

GridironGo is a Fantasy Football CLI app built in Go using the [Bubble Tea](https://github.com/charmbracelet/bubbletea) TUI framework.

## Project Structure

```
.
├── assets
├── go.mod
├── go.sum
├── internals
│   ├── data
│   │   ├── db.go
│   │   ├── models.go
│   │   └── scrape.go
│   ├── league
│   │   ├── league.go
│   │   ├── rules.go
│   │   ├── schedule.go
│   │   └── team.go
│   ├── rules
│   │   ├── draft.go
│   │   └── scoring.go
│   └── tui
│       ├── league_menu.go
│       ├── menu.go
│       ├── player_menu.go
│       └── schedule_menu.go
├── LICENSE
├── main.go
└── README.md
```

### File Descriptions

- **`assets/`**: Directory for static assets like JSON data, logos, or templates.
- **`go.mod`**: Defines the module name and dependencies for the Go project.
- **`go.sum`**: Locks dependency versions to ensure reproducible builds.
- **`internals/`**: Contains core application logic split into sub-packages.
- **`internals/data/db.go`**: Handles SQLite database connections and queries.
- **`internals/data/models.go`**: Defines the data models for players, teams, and matches.
- **`internals/data/scrape.go`**: Scrapes football stats and schedules for the last 3 seasons.
- **`internals/league/league.go`**: Manages fantasy league setup and operations.
- **`internals/league/rules.go`**: Handles league rules including scoring and configurations.
- **`internals/league/schedule.go`**: Generates and manages league schedules, including playoffs.
- **`internals/league/team.go`**: Manages fantasy teams including bot teams and user team.
- **`internals/rules/draft.go`**: Handles the drafting logic and player selection process.
- **`internals/rules/scoring.go`**: Implements fantasy football scoring rules and calculations.
- **`internals/tui/league_menu.go`**: TUI logic for the fantasy league menu and its options.
- **`internals/tui/menu.go`**: Main TUI entry point with initial menu options.
- **`internals/tui/player_menu.go`**: TUI logic for viewing players and selecting them.
- **`internals/tui/schedule_menu.go`**: TUI logic for viewing the real and fantasy schedules.
- **`LICENSE`**: Project license file.
- **`main.go`**: Entry point for the application. Initializes TUI and loads data.
- **`README.md`**: Project documentation and setup instructions.

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

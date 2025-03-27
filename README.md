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
│   │   ├── scrape.go           # Scrapes football stats and schedules for the last 3 seasons
│   │   └── queries             # Directory for SQL queries used by sqlc
│   │       ├── schema.sql      # Database schema definition with tables and indexes
│   │       ├── player.sql      # Player-related queries (stats, fantasy points, searching)
│   │       ├── team.sql        # Team management queries (roster, standings, updates)
│   │       ├── league.sql      # League operations queries (rules, matchups, settings)
│   │       ├── season.sql      # Season data queries (schedules, games, tracking progress)
│   │       ├── draft.sql       # Draft functionality queries (picks, available players, tracking)
│   │       └── score.sql       # Scoring system queries (weekly scores, standings updates)
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
├── sqlc.yaml                   # Configuration file for sqlc code generation
└── README.md                   # Project documentation and setup instructions
```

## Database Schema
The application uses SQLite with sqlc for type-safe database operations. Key tables include:

- **NFL Data**: 
  - `seasons` - NFL seasons information
  - `nfl_teams` - All NFL teams data
  - `nfl_players` - NFL players with their positions
  - `nfl_games` - NFL game schedule and results
  - `player_stats` - Player statistics for each game

- **Fantasy Football**:
  - `fantasy_leagues` - League settings and scoring rules
  - `fantasy_teams` - Teams in each league (user and bots)
  - `fantasy_rosters` - Players on each team's roster
  - `fantasy_matchups` - Weekly matchups between teams
  - `fantasy_drafts` - Draft history and picks
  - `fantasy_player_scores` - Weekly fantasy scores for players

## Features
- Create and manage a Fantasy Football League with bot and user teams
- View and draft players using previous season stats
- Customize league rules including PPR, scoring, and roster settings
- Full TUI experience using Bubble Tea
- SQLite backend for storing player stats and league data
- Data scraping for the last 3 seasons of NFL stats

## Fantasy League Features
- Customizable roster positions (QB, RB, WR, TE, FLEX, K, DST)
- PPR (Points Per Reception) option
- Customizable scoring settings for all stat categories
- Automatic schedule generation
- Regular season (weeks 1-14) and playoffs (weeks 15-16)
- Top 4 teams make playoffs based on record and points
- Full draft system with player rankings based on historical performance

## Getting Started
1. Clone the repo
2. Run `go mod tidy` to install dependencies
3. Run `sqlc generate` to generate database code
4. Run `go run main.go` to launch the TUI

## License
MIT

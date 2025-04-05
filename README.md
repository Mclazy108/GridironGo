# GridironGo
GridironGo is a Fantasy Football CLI app built in Go using the [Bubble Tea](https://github.com/charmbracelet/bubbletea) TUI framework.

![Lines of Code](https://tokei.rs/b1/github/Mclazy108/GridironGo)


## Project Structure
```
.
├── assets                      # Directory for static assets like JSON data, logos, or templates
├── build.sh                    # Script to build executables for different platforms
├── go.mod                      # Defines the module name and dependencies for the Go project
├── go.sum                      # Locks dependency versions to ensure reproducible builds
├── internals                   # Contains core application logic split into sub-packages
│   ├── data
│   │   ├── database.go         # Handles SQLite database connections and queries
│   │   ├── migrations          # Directory for SQL migrations
│   │   │   └── schema.sql      # Database schema definition with tables and indexes
│   │   ├── queries             # Directory for SQL queries used by sqlc
│   │   │   ├── draft.sql       # Draft functionality queries (picks, available players, tracking)
│   │   │   ├── league.sql      # League operations queries (rules, matchups, settings)
│   │   │   ├── player.sql      # Player-related queries (stats, fantasy points, searching)
│   │   │   ├── score.sql       # Scoring system queries (weekly scores, standings updates)
│   │   │   ├── season.sql      # Season data queries (schedules, games, tracking progress)
│   │   │   ├── teams.sql       # Team management queries (roster, standings, updates)
│   │   │   └── games.sql       # Game schedule queries
│   │   ├── scraper             # Data scrapers for NFL data
│   │   │   ├── scrape-games.go # Scrapes NFL game schedules from ESPN API
│   │   │   └── scrape-teams.go # Scrapes NFL team data from ESPN API
│   │   └── sqlc                # Generated SQL code by sqlc
│   │       ├── db.go           # Database connection and query execution
│   │       ├── games.sql.go    # Generated code for game queries
│   │       ├── teams.sql.go    # Generated code for team queries
│   │       ├── models.go       # Generated data models
│   │       └── querier.go      # Interface defining all available queries
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
├── README.md                   # Project documentation and setup instructions
└── sqlc.yaml                   # Configuration file for sqlc code generation
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
   ```bash
   git clone https://github.com/Mclazy108/GridironGo.git
   cd GridironGo
   ```

2. Install dependencies
   ```bash
   go mod tidy
   ```

3. Generate SQLc code
   ```bash
   sqlc generate
   ```

4. Scrape NFL game data
   ```bash
   go run main.go -scrape-games
   ```

5. Scrape NFL team data
   ```bash
   go run main.go -scrape-teams
   ```

6. Run the application
   ```bash
   go run main.go
   ```

7. Or build the application for your platform
   ```bash
   go build -o GridironGo
   ```

## Command Line Options
GridironGo supports several command line flags:

- `-db`: Specify the path to the SQLite database (default: "./GridironGo.db")
- `-scrape-games`: Scrape NFL game data for seasons 2022-2024
- `-scrape-teams`: Scrape NFL team data (team names, divisions, colors, etc.)

## Building Executables
The project includes a build script that creates executables for multiple platforms:

```bash
# Make the build script executable
chmod +x build.sh

# Run the build script
./build.sh
```

This will create executables for Linux, macOS, and Windows in the `executables` directory.

## Data Sources
The application uses ESPN API endpoints to fetch NFL data:
- Game schedules: `https://site.api.espn.com/apis/site/v2/sports/football/nfl/scoreboard`
- Team information: `https://sports.core.api.espn.com/v2/sports/football/leagues/nfl/teams`
- Detailed team data: `https://site.api.espn.com/apis/site/v2/sports/football/nfl/teams/{team_id}`

## License
MIT

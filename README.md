# GridironGo
GridironGo is a Fantasy Football CLI app built in Go using the [Bubble Tea](https://github.com/charmbracelet/bubbletea) TUI framework.

![Lines of Code](https://tokei.rs/b1/github/Mclazy108/GridironGo)

---

## Project Structure
```
.
├── LICENSE                     		# Project license file
├── README.md                   		# Project documentation and setup instructions
├── build.sh                    		# Script to build executables for different platforms
├── examples                    		# Directory with example code and API documentation
│   └── api.csv                 		# CSV listing all available ESPN API endpoints with descriptions
├── executables                 		# Directory for compiled binaries
│   ├── GridironGo-linux        		# Linux executable
│   ├── GridironGo-mac          		# macOS executable
│   └── GridironGo-windows.exe  		# Windows executable
├── go.mod                      		# Defines the module name and dependencies for the Go project
├── go.sum                      		# Locks dependency versions to ensure reproducible builds
├── internals                   		# Contains core application logic split into sub-packages
│   ├── data                    		# Data layer for database operations and scraping
│   │   ├── database.go         		# Handles SQLite database connections and queries
│   │   ├── migrations          		# Directory for SQL migrations
│   │   │   └── schema.sql      		# Database schema definition with tables and indexes
│   │   ├── queries             		# Directory for SQL queries used by sqlc
│   │   │   ├── games.sql       		# Game schedule queries
│   │   │   ├── player_seasons.sql 		# Player season tracking queries
│   │   │   ├── players.sql     		# Player-related queries (stats, fantasy points, searching)
│   │   │   ├── stats.sql       		# Statistics and scoring system queries
│   │   │   └── teams.sql       		# Team management queries (roster, standings, updates)
│   │   ├── scraper             		# Data scrapers for NFL data
│   │   │   ├── scrape-games.go 		# Scrapes NFL game schedules from ESPN API
│   │   │   ├── scrape-players.go 		# Scrapes NFL player data from ESPN API
│   │   │   ├── scrape-stats.go 		# Scrapes NFL player and game statistics from ESPN API
│   │   │   └── scrape-teams.go 		# Scrapes NFL team data from ESPN API
│   │   └── sqlc                		# Generated SQL code by sqlc
│   │       ├── db.go           		# Database connection and query execution
│   │       ├── games.sql.go    		# Generated code for game queries
│   │       ├── models.go       		# Generated data models
│   │       ├── player_seasons.sql.go 	# Generated code for player seasons queries
│   │       ├── players.sql.go  		# Generated code for player queries
│   │       ├── querier.go      		# Interface defining all available queries
│   │       ├── stats.sql.go    		# Generated code for statistics queries
│   │       └── teams.sql.go    		# Generated code for team queries
│   ├── league                  		# Fantasy league management
│   │   ├── league.go           		# Manages fantasy league setup and operations
│   │   ├── rules.go            		# Handles league rules including scoring and configurations
│   │   ├── schedule.go         		# Generates and manages league schedules, including playoffs
│   │   ├── team.go             		# Manages fantasy teams including bot teams and user team
│   │   ├── draft.go            		# Handles the drafting logic and player selection process
│   │   └── scoring.go          		# Implements fantasy football scoring rules and calculations
│   └── tui                     		# Terminal User Interface components
│       ├── league_menu.go      		# TUI logic for the fantasy league menu and its options
│       ├── menu.go             		# Main TUI entry point with initial menu options
│       ├── player_menu.go      		# TUI logic for viewing players and selecting them
│       └── schedule_menu.go    		# TUI logic for viewing the real and fantasy schedules
├── main.go                     		# Entry point for the application
├── planning.txt                		# Project planning notes and roadmap
└── sqlc.yaml                   		# Configuration file for sqlc code generation
```

---

## Features
- Create and manage a Fantasy Football League with bot and user teams
- View and draft players using previous season stats
- Customize league rules including PPR, scoring, and roster settings
- Full TUI experience using Bubble Tea
- SQLite backend for storing player stats and league data
- Data scraping for NFL teams, players, and schedules

## Fantasy League Features
- Customizable roster positions (QB, RB, WR, TE, FLEX, K, DST)
- PPR (Points Per Reception) option
- Customizable scoring settings for all stat categories
- Automatic schedule generation
- Regular season (weeks 1–14) and playoffs (weeks 15–16)
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

4. Scrape NFL data
   ```bash
   go run main.go -scrape-teams
   go run main.go -scrape-players
   go run main.go -scrape-games
   go run main.go -scrape-stats
   ```

5. Run the application
   ```bash
   go run main.go
   ```

6. Or build the application for your platform
   ```bash
   go build -o GridironGo
   ```

## Command Line Options
- `-db`: Specify path to SQLite database (default: "./GridironGo.db")
- `-scrape-games`: Scrape NFL game data
- `-scrape-teams`: Scrape NFL team data
- `-scrape-players`: Scrape NFL player data
- `-scrape-stats`: Scrape NFL game statistics
- `-seasons`: Comma-separated list of seasons to scrape data for (default: "2022,2023,2024")

## Scraping Examples
```bash
# Scrape all teams
go run main.go -scrape-teams

# Scrape games for default seasons (2022-2024)
go run main.go -scrape-games

# Scrape games for specific seasons
go run main.go -scrape-games -seasons="2023,2024"

# Scrape players for all teams for specific seasons
go run main.go -scrape-players -seasons="2023,2024"

# Scrape player stats
go run main.go -scrape-stats -seasons="2023"

# Scrape everything with custom database path
go run main.go -scrape-teams -scrape-games -scrape-players -scrape-stats -db="./data/nfl.db"
```

## Building Executables
The project includes a build script that creates executables for multiple platforms:

```bash
chmod +x build.sh
./build.sh
```

Executables will be saved in the `executables` directory.

## Data Sources
This app uses the following ESPN APIs:

### Currently Used APIs
The following APIs are actively used in the current codebase:

- 🏈 **Game Schedules**
  `https://site.api.espn.com/apis/site/v2/sports/football/nfl/scoreboard?dates={year}&seasontype=2&week={week}`

- 👥 **NFL Teams List**
  `https://sports.core.api.espn.com/v2/sports/football/leagues/nfl/teams`

- 🧾 **Team Details**
  `https://site.api.espn.com/apis/site/v2/sports/football/nfl/teams/{team_id}`

- 👤 **Team Roster by Year**
  `https://sports.core.api.espn.com/v2/sports/football/leagues/nfl/seasons/{year}/teams/{team_id}/athletes`

- 📋 **Player Detail Lookup**
  `https://sports.core.api.espn.com/v2/sports/football/leagues/nfl/athletes/{player_id}`

- 📊 **Game Summary with Stats**
  `https://site.api.espn.com/apis/site/v2/sports/football/nfl/summary?event={event_id}`

## Database Schema
The application uses SQLite with the following tables:

- `nfl_games` - Store NFL game information (teams, dates, seasons)
- `nfl_teams` - Store NFL team information (names, abbreviations, divisions)
- `nfl_players` - Store NFL player information (names, positions, stats)
- `nfl_player_seasons` - Store player information for specific seasons
- `nfl_stats` - Store game statistics for players and teams

## License
MIT

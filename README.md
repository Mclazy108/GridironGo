# GridironGo
GridironGo is a Fantasy Football CLI app built in Go using the [Bubble Tea](https://github.com/charmbracelet/bubbletea) TUI framework.

![Lines of Code](https://tokei.rs/b1/github/Mclazy108/GridironGo)

---

## Project Structure
```
.
├── LICENSE
├── README.md
├── build.sh
├── examples
│   ├── api.csv
│   └── events_to_csv.go
├── executables
│   ├── GridironGo-linux
│   ├── GridironGo-mac
│   └── GridironGo-windows.exe
├── go.mod
├── go.sum
├── internals
│   ├── data
│   │   ├── database.go
│   │   ├── migrations
│   │   │   └── schema.sql
│   │   ├── queries
│   │   │   ├── games.sql
│   │   │   ├── player_seasons.sql
│   │   │   ├── players.sql
│   │   │   ├── stats.sql
│   │   │   └── teams.sql
│   │   ├── scraper
│   │   │   ├── scrape-games.go
│   │   │   ├── scrape-players.go
│   │   │   ├── scrape-stats.go
│   │   │   └── scrape-teams.go
│   │   └── sqlc
│   │       ├── db.go
│   │       ├── games.sql.go
│   │       ├── models.go
│   │       ├── player_seasons.sql.go
│   │       ├── players.sql.go
│   │       ├── querier.go
│   │       ├── stats.sql.go
│   │       └── teams.sql.go
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
├── main.go
├── planning.txt
└── sqlc.yaml
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

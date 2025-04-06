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
│   ├── api.csv
│   └── events_to_csv.go
├── executables
│   ├── GridironGo-linux
│   ├── GridironGo-mac
│   └── GridironGo-windows.exe
├── go.mod
├── go.sum
├── internals
│   ├── data
│   │   ├── database.go
│   │   ├── migrations
│   │   │   └── schema.sql
│   │   ├── queries
│   │   │   ├── games.sql
│   │   │   ├── players.sql
│   │   │   └── teams.sql
│   │   ├── scraper
│   │   │   ├── scrape-games.go
│   │   │   ├── scrape-players.go
│   │   │   └── scrape-teams.go
│   │   └── sqlc
│   │       ├── db.go
│   │       ├── games.sql.go
│   │       ├── players.sql.go
│   │       ├── teams.sql.go
│   │       ├── models.go
│   │       └── querier.go
│   ├── league
│   │   ├── league.go
│   │   ├── rules.go
│   │   ├── schedule.go
│   │   └── team.go
│   ├── rules
│   │   ├── draft.go
│   │   └── scoring.go
│   └── tui
│       ├── league_menu.go
│       ├── menu.go
│       ├── player_menu.go
│       └── schedule_menu.go
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

## Building Executables
The project includes a build script that creates executables for multiple platforms:

```bash
chmod +x build.sh
./build.sh
```

Executables will be saved in the `executables` directory.

## Data Sources
This app uses the following ESPN APIs:

- 🏈 **Game Schedules**  
  `https://site.api.espn.com/apis/site/v2/sports/football/nfl/scoreboard?dates={year}&seasontype=2&week={week}`

- 👥 **NFL Teams List**  
  `https://sports.core.api.espn.com/v2/sports/football/leagues/nfl/teams`

- 🧾 **Team Details**  
  `https://site.api.espn.com/apis/site/v2/sports/football/nfl/teams/{team_id}`

- 👤 **Team Roster by Year**  
  `https://sports.core.api.espn.com/v2/sports/football/leagues/nfl/seasons/2024/teams/{team_id}/athletes`

- 📋 **Player Detail Lookup**  
  `https://sports.core.api.espn.com/v2/sports/football/leagues/nfl/athletes/{player_id}`

## License
MIT

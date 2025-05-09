// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package sqlc

import (
	"database/sql"
)

type NflGame struct {
	EventID   int64  `json:"event_id"`
	Date      string `json:"date"`
	Name      string `json:"name"`
	ShortName string `json:"short_name"`
	Season    int64  `json:"season"`
	Week      int64  `json:"week"`
	AwayTeam  string `json:"away_team"`
	HomeTeam  string `json:"home_team"`
}

type NflPlayer struct {
	PlayerID   string         `json:"player_id"`
	FirstName  string         `json:"first_name"`
	LastName   string         `json:"last_name"`
	FullName   string         `json:"full_name"`
	Position   string         `json:"position"`
	TeamID     sql.NullString `json:"team_id"`
	Jersey     sql.NullString `json:"jersey"`
	Height     sql.NullInt64  `json:"height"`
	Weight     sql.NullInt64  `json:"weight"`
	Active     bool           `json:"active"`
	College    sql.NullString `json:"college"`
	Experience sql.NullInt64  `json:"experience"`
	DraftYear  sql.NullInt64  `json:"draft_year"`
	DraftRound sql.NullInt64  `json:"draft_round"`
	DraftPick  sql.NullInt64  `json:"draft_pick"`
	Status     sql.NullString `json:"status"`
	ImageUrl   sql.NullString `json:"image_url"`
}

type NflPlayerSeason struct {
	PlayerID   string         `json:"player_id"`
	SeasonYear int64          `json:"season_year"`
	TeamID     sql.NullString `json:"team_id"`
	Jersey     sql.NullString `json:"jersey"`
	Active     bool           `json:"active"`
	Experience sql.NullInt64  `json:"experience"`
	Status     sql.NullString `json:"status"`
}

type NflStat struct {
	StatID    int64   `json:"stat_id"`
	GameID    int64   `json:"game_id"`
	PlayerID  string  `json:"player_id"`
	TeamID    string  `json:"team_id"`
	Category  string  `json:"category"`
	StatType  string  `json:"stat_type"`
	StatValue float64 `json:"stat_value"`
}

type NflTeam struct {
	TeamID         string         `json:"team_id"`
	DisplayName    string         `json:"display_name"`
	Abbreviation   string         `json:"abbreviation"`
	ShortName      string         `json:"short_name"`
	Location       string         `json:"location"`
	Nickname       string         `json:"nickname"`
	Conference     string         `json:"conference"`
	Division       string         `json:"division"`
	PrimaryColor   sql.NullString `json:"primary_color"`
	SecondaryColor sql.NullString `json:"secondary_color"`
	LogoUrl        sql.NullString `json:"logo_url"`
}

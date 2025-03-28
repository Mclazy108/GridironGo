// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: team.sql

package sqlc

import (
	"context"
	"database/sql"
)

const addPlayerToFantasyTeam = `-- name: AddPlayerToFantasyTeam :one
INSERT INTO fantasy_rosters (
    team_id,
    player_id,
    position,
    is_starter
) VALUES (?, ?, ?, ?)
RETURNING id
`

type AddPlayerToFantasyTeamParams struct {
	TeamID    int64         `json:"team_id"`
	PlayerID  int64         `json:"player_id"`
	Position  string        `json:"position"`
	IsStarter sql.NullInt64 `json:"is_starter"`
}

func (q *Queries) AddPlayerToFantasyTeam(ctx context.Context, arg AddPlayerToFantasyTeamParams) (int64, error) {
	row := q.queryRow(ctx, q.addPlayerToFantasyTeamStmt, addPlayerToFantasyTeam,
		arg.TeamID,
		arg.PlayerID,
		arg.Position,
		arg.IsStarter,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const createFantasyTeam = `-- name: CreateFantasyTeam :one
INSERT INTO fantasy_teams (
    league_id, 
    name, 
    owner_name, 
    is_user, 
    draft_position, 
    wins, 
    losses, 
    tie_games, 
    points_for, 
    points_against
) VALUES (?, ?, ?, ?, ?, 0, 0, 0, 0, 0)
RETURNING id
`

type CreateFantasyTeamParams struct {
	LeagueID      int64         `json:"league_id"`
	Name          string        `json:"name"`
	OwnerName     string        `json:"owner_name"`
	IsUser        sql.NullInt64 `json:"is_user"`
	DraftPosition sql.NullInt64 `json:"draft_position"`
}

func (q *Queries) CreateFantasyTeam(ctx context.Context, arg CreateFantasyTeamParams) (int64, error) {
	row := q.queryRow(ctx, q.createFantasyTeamStmt, createFantasyTeam,
		arg.LeagueID,
		arg.Name,
		arg.OwnerName,
		arg.IsUser,
		arg.DraftPosition,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const getAllFantasyTeams = `-- name: GetAllFantasyTeams :many
SELECT 
    ft.id, 
    ft.name, 
    ft.owner_name, 
    ft.is_user, 
    ft.draft_position, 
    ft.wins, 
    ft.losses, 
    ft.tie_games, 
    ft.points_for, 
    ft.points_against,
    fl.name as league_name
FROM fantasy_teams ft
JOIN fantasy_leagues fl ON ft.league_id = fl.id
WHERE ft.league_id = ?
ORDER BY (ft.wins * 2 + ft.tie_games) DESC, ft.points_for DESC
`

type GetAllFantasyTeamsRow struct {
	ID            int64           `json:"id"`
	Name          string          `json:"name"`
	OwnerName     string          `json:"owner_name"`
	IsUser        sql.NullInt64   `json:"is_user"`
	DraftPosition sql.NullInt64   `json:"draft_position"`
	Wins          sql.NullInt64   `json:"wins"`
	Losses        sql.NullInt64   `json:"losses"`
	TieGames      sql.NullInt64   `json:"tie_games"`
	PointsFor     sql.NullFloat64 `json:"points_for"`
	PointsAgainst sql.NullFloat64 `json:"points_against"`
	LeagueName    string          `json:"league_name"`
}

func (q *Queries) GetAllFantasyTeams(ctx context.Context, leagueID int64) ([]*GetAllFantasyTeamsRow, error) {
	rows, err := q.query(ctx, q.getAllFantasyTeamsStmt, getAllFantasyTeams, leagueID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetAllFantasyTeamsRow{}
	for rows.Next() {
		var i GetAllFantasyTeamsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.OwnerName,
			&i.IsUser,
			&i.DraftPosition,
			&i.Wins,
			&i.Losses,
			&i.TieGames,
			&i.PointsFor,
			&i.PointsAgainst,
			&i.LeagueName,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllNFLTeams = `-- name: GetAllNFLTeams :many
SELECT id, name, city, abbreviation, conference, division
FROM nfl_teams
ORDER BY conference, division, name
`

func (q *Queries) GetAllNFLTeams(ctx context.Context) ([]*NflTeam, error) {
	rows, err := q.query(ctx, q.getAllNFLTeamsStmt, getAllNFLTeams)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*NflTeam{}
	for rows.Next() {
		var i NflTeam
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.City,
			&i.Abbreviation,
			&i.Conference,
			&i.Division,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAvailablePlayers = `-- name: GetAvailablePlayers :many
SELECT 
    p.id,
    p.name,
    p.position,
    p.team_id,
    t.name as team_name,
    t.abbreviation as team_abbreviation
FROM nfl_players p
JOIN nfl_teams t ON p.team_id = t.id
WHERE p.id NOT IN (
    SELECT player_id FROM fantasy_rosters fr
    JOIN fantasy_teams ft ON fr.team_id = ft.id
    WHERE ft.league_id = ?
)
AND p.position = ?
ORDER BY p.name
LIMIT ?
`

type GetAvailablePlayersParams struct {
	LeagueID int64  `json:"league_id"`
	Position string `json:"position"`
	Limit    int64  `json:"limit"`
}

type GetAvailablePlayersRow struct {
	ID               int64         `json:"id"`
	Name             string        `json:"name"`
	Position         string        `json:"position"`
	TeamID           sql.NullInt64 `json:"team_id"`
	TeamName         string        `json:"team_name"`
	TeamAbbreviation string        `json:"team_abbreviation"`
}

func (q *Queries) GetAvailablePlayers(ctx context.Context, arg GetAvailablePlayersParams) ([]*GetAvailablePlayersRow, error) {
	rows, err := q.query(ctx, q.getAvailablePlayersStmt, getAvailablePlayers, arg.LeagueID, arg.Position, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetAvailablePlayersRow{}
	for rows.Next() {
		var i GetAvailablePlayersRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Position,
			&i.TeamID,
			&i.TeamName,
			&i.TeamAbbreviation,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getFantasyTeamById = `-- name: GetFantasyTeamById :one
SELECT 
    ft.id, 
    ft.league_id,
    ft.name, 
    ft.owner_name, 
    ft.is_user, 
    ft.draft_position, 
    ft.wins, 
    ft.losses, 
    ft.tie_games, 
    ft.points_for, 
    ft.points_against,
    fl.name as league_name
FROM fantasy_teams ft
JOIN fantasy_leagues fl ON ft.league_id = fl.id
WHERE ft.id = ?
`

type GetFantasyTeamByIdRow struct {
	ID            int64           `json:"id"`
	LeagueID      int64           `json:"league_id"`
	Name          string          `json:"name"`
	OwnerName     string          `json:"owner_name"`
	IsUser        sql.NullInt64   `json:"is_user"`
	DraftPosition sql.NullInt64   `json:"draft_position"`
	Wins          sql.NullInt64   `json:"wins"`
	Losses        sql.NullInt64   `json:"losses"`
	TieGames      sql.NullInt64   `json:"tie_games"`
	PointsFor     sql.NullFloat64 `json:"points_for"`
	PointsAgainst sql.NullFloat64 `json:"points_against"`
	LeagueName    string          `json:"league_name"`
}

func (q *Queries) GetFantasyTeamById(ctx context.Context, id int64) (*GetFantasyTeamByIdRow, error) {
	row := q.queryRow(ctx, q.getFantasyTeamByIdStmt, getFantasyTeamById, id)
	var i GetFantasyTeamByIdRow
	err := row.Scan(
		&i.ID,
		&i.LeagueID,
		&i.Name,
		&i.OwnerName,
		&i.IsUser,
		&i.DraftPosition,
		&i.Wins,
		&i.Losses,
		&i.TieGames,
		&i.PointsFor,
		&i.PointsAgainst,
		&i.LeagueName,
	)
	return &i, err
}

const getFantasyTeamForWeek = `-- name: GetFantasyTeamForWeek :one
SELECT 
    ft.id, 
    ft.name, 
    ft.owner_name,
    ft.wins,
    ft.losses,
    ft.tie_games,
    SUM(CASE WHEN fm.home_team_id = ft.id THEN fm.home_score ELSE fm.away_score END) as points_for_week
FROM fantasy_teams ft
JOIN fantasy_matchups fm ON (fm.home_team_id = ft.id OR fm.away_team_id = ft.id)
WHERE ft.id = ? AND fm.week = ? AND fm.league_id = ?
GROUP BY ft.id, ft.name, ft.owner_name, ft.wins, ft.losses, ft.tie_games
`

type GetFantasyTeamForWeekParams struct {
	ID       int64 `json:"id"`
	Week     int64 `json:"week"`
	LeagueID int64 `json:"league_id"`
}

type GetFantasyTeamForWeekRow struct {
	ID            int64           `json:"id"`
	Name          string          `json:"name"`
	OwnerName     string          `json:"owner_name"`
	Wins          sql.NullInt64   `json:"wins"`
	Losses        sql.NullInt64   `json:"losses"`
	TieGames      sql.NullInt64   `json:"tie_games"`
	PointsForWeek sql.NullFloat64 `json:"points_for_week"`
}

func (q *Queries) GetFantasyTeamForWeek(ctx context.Context, arg GetFantasyTeamForWeekParams) (*GetFantasyTeamForWeekRow, error) {
	row := q.queryRow(ctx, q.getFantasyTeamForWeekStmt, getFantasyTeamForWeek, arg.ID, arg.Week, arg.LeagueID)
	var i GetFantasyTeamForWeekRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.OwnerName,
		&i.Wins,
		&i.Losses,
		&i.TieGames,
		&i.PointsForWeek,
	)
	return &i, err
}

const getFantasyTeamRoster = `-- name: GetFantasyTeamRoster :many
SELECT 
    fr.id as roster_id,
    p.id as player_id,
    p.name as player_name,
    fr.position as roster_position,
    p.position as player_position,
    fr.is_starter,
    nt.name as nfl_team_name
FROM fantasy_rosters fr
JOIN nfl_players p ON fr.player_id = p.id
LEFT JOIN nfl_teams nt ON p.team_id = nt.id
WHERE fr.team_id = ?
ORDER BY fr.is_starter DESC, fr.position, p.name
`

type GetFantasyTeamRosterRow struct {
	RosterID       int64          `json:"roster_id"`
	PlayerID       int64          `json:"player_id"`
	PlayerName     string         `json:"player_name"`
	RosterPosition string         `json:"roster_position"`
	PlayerPosition string         `json:"player_position"`
	IsStarter      sql.NullInt64  `json:"is_starter"`
	NflTeamName    sql.NullString `json:"nfl_team_name"`
}

func (q *Queries) GetFantasyTeamRoster(ctx context.Context, teamID int64) ([]*GetFantasyTeamRosterRow, error) {
	rows, err := q.query(ctx, q.getFantasyTeamRosterStmt, getFantasyTeamRoster, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetFantasyTeamRosterRow{}
	for rows.Next() {
		var i GetFantasyTeamRosterRow
		if err := rows.Scan(
			&i.RosterID,
			&i.PlayerID,
			&i.PlayerName,
			&i.RosterPosition,
			&i.PlayerPosition,
			&i.IsStarter,
			&i.NflTeamName,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getNFLTeamByAbbreviation = `-- name: GetNFLTeamByAbbreviation :one
SELECT id, name, city, abbreviation, conference, division
FROM nfl_teams
WHERE abbreviation = ?
`

func (q *Queries) GetNFLTeamByAbbreviation(ctx context.Context, abbreviation string) (*NflTeam, error) {
	row := q.queryRow(ctx, q.getNFLTeamByAbbreviationStmt, getNFLTeamByAbbreviation, abbreviation)
	var i NflTeam
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.City,
		&i.Abbreviation,
		&i.Conference,
		&i.Division,
	)
	return &i, err
}

const getNFLTeamById = `-- name: GetNFLTeamById :one
SELECT id, name, city, abbreviation, conference, division
FROM nfl_teams
WHERE id = ?
`

func (q *Queries) GetNFLTeamById(ctx context.Context, id int64) (*NflTeam, error) {
	row := q.queryRow(ctx, q.getNFLTeamByIdStmt, getNFLTeamById, id)
	var i NflTeam
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.City,
		&i.Abbreviation,
		&i.Conference,
		&i.Division,
	)
	return &i, err
}

const getTeamRoster = `-- name: GetTeamRoster :many
SELECT 
    p.id, 
    p.name, 
    p.position, 
    p.jersey_number, 
    p.status
FROM nfl_players p
WHERE p.team_id = ?
ORDER BY p.position, p.name
`

type GetTeamRosterRow struct {
	ID           int64          `json:"id"`
	Name         string         `json:"name"`
	Position     string         `json:"position"`
	JerseyNumber sql.NullInt64  `json:"jersey_number"`
	Status       sql.NullString `json:"status"`
}

func (q *Queries) GetTeamRoster(ctx context.Context, teamID sql.NullInt64) ([]*GetTeamRosterRow, error) {
	rows, err := q.query(ctx, q.getTeamRosterStmt, getTeamRoster, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetTeamRosterRow{}
	for rows.Next() {
		var i GetTeamRosterRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Position,
			&i.JerseyNumber,
			&i.Status,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTeamsByConference = `-- name: GetTeamsByConference :many
SELECT id, name, city, abbreviation, conference, division
FROM nfl_teams
WHERE conference = ?
ORDER BY division, name
`

func (q *Queries) GetTeamsByConference(ctx context.Context, conference string) ([]*NflTeam, error) {
	rows, err := q.query(ctx, q.getTeamsByConferenceStmt, getTeamsByConference, conference)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*NflTeam{}
	for rows.Next() {
		var i NflTeam
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.City,
			&i.Abbreviation,
			&i.Conference,
			&i.Division,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTeamsByDivision = `-- name: GetTeamsByDivision :many
SELECT id, name, city, abbreviation, conference, division
FROM nfl_teams
WHERE division = ?
ORDER BY name
`

func (q *Queries) GetTeamsByDivision(ctx context.Context, division string) ([]*NflTeam, error) {
	rows, err := q.query(ctx, q.getTeamsByDivisionStmt, getTeamsByDivision, division)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*NflTeam{}
	for rows.Next() {
		var i NflTeam
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.City,
			&i.Abbreviation,
			&i.Conference,
			&i.Division,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const insertTeam = `-- name: InsertTeam :one
INSERT INTO nfl_teams (name, city, abbreviation, conference, division)
VALUES (?, ?, ?, ?, ?)
RETURNING id
`

type InsertTeamParams struct {
	Name         string `json:"name"`
	City         string `json:"city"`
	Abbreviation string `json:"abbreviation"`
	Conference   string `json:"conference"`
	Division     string `json:"division"`
}

func (q *Queries) InsertTeam(ctx context.Context, arg InsertTeamParams) (int64, error) {
	row := q.queryRow(ctx, q.insertTeamStmt, insertTeam,
		arg.Name,
		arg.City,
		arg.Abbreviation,
		arg.Conference,
		arg.Division,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const removePlayerFromFantasyTeam = `-- name: RemovePlayerFromFantasyTeam :exec
DELETE FROM fantasy_rosters
WHERE team_id = ? AND player_id = ?
`

type RemovePlayerFromFantasyTeamParams struct {
	TeamID   int64 `json:"team_id"`
	PlayerID int64 `json:"player_id"`
}

func (q *Queries) RemovePlayerFromFantasyTeam(ctx context.Context, arg RemovePlayerFromFantasyTeamParams) error {
	_, err := q.exec(ctx, q.removePlayerFromFantasyTeamStmt, removePlayerFromFantasyTeam, arg.TeamID, arg.PlayerID)
	return err
}

const updateFantasyRoster = `-- name: UpdateFantasyRoster :exec
UPDATE fantasy_rosters
SET is_starter = ?
WHERE id = ?
`

type UpdateFantasyRosterParams struct {
	IsStarter sql.NullInt64 `json:"is_starter"`
	ID        int64         `json:"id"`
}

func (q *Queries) UpdateFantasyRoster(ctx context.Context, arg UpdateFantasyRosterParams) error {
	_, err := q.exec(ctx, q.updateFantasyRosterStmt, updateFantasyRoster, arg.IsStarter, arg.ID)
	return err
}

const updateFantasyTeamRecord = `-- name: UpdateFantasyTeamRecord :exec
UPDATE fantasy_teams
SET 
    wins = ?,
    losses = ?,
    tie_games = ?,
    points_for = ?,
    points_against = ?
WHERE id = ?
`

type UpdateFantasyTeamRecordParams struct {
	Wins          sql.NullInt64   `json:"wins"`
	Losses        sql.NullInt64   `json:"losses"`
	TieGames      sql.NullInt64   `json:"tie_games"`
	PointsFor     sql.NullFloat64 `json:"points_for"`
	PointsAgainst sql.NullFloat64 `json:"points_against"`
	ID            int64           `json:"id"`
}

func (q *Queries) UpdateFantasyTeamRecord(ctx context.Context, arg UpdateFantasyTeamRecordParams) error {
	_, err := q.exec(ctx, q.updateFantasyTeamRecordStmt, updateFantasyTeamRecord,
		arg.Wins,
		arg.Losses,
		arg.TieGames,
		arg.PointsFor,
		arg.PointsAgainst,
		arg.ID,
	)
	return err
}

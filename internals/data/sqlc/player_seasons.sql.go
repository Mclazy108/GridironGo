// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: player_seasons.sql

package sqlc

import (
	"context"
	"database/sql"
)

const createPlayerSeason = `-- name: CreatePlayerSeason :exec
INSERT INTO nfl_player_seasons (
  player_id, season_year, team_id, jersey, active, experience, status
) VALUES (
  ?, ?, ?, ?, ?, ?, ?
)
`

type CreatePlayerSeasonParams struct {
	PlayerID   string         `json:"player_id"`
	SeasonYear int64          `json:"season_year"`
	TeamID     sql.NullString `json:"team_id"`
	Jersey     sql.NullString `json:"jersey"`
	Active     bool           `json:"active"`
	Experience sql.NullInt64  `json:"experience"`
	Status     sql.NullString `json:"status"`
}

func (q *Queries) CreatePlayerSeason(ctx context.Context, arg CreatePlayerSeasonParams) error {
	_, err := q.exec(ctx, q.createPlayerSeasonStmt, createPlayerSeason,
		arg.PlayerID,
		arg.SeasonYear,
		arg.TeamID,
		arg.Jersey,
		arg.Active,
		arg.Experience,
		arg.Status,
	)
	return err
}

const deletePlayerSeason = `-- name: DeletePlayerSeason :exec
DELETE FROM nfl_player_seasons
WHERE player_id = ? AND season_year = ?
`

type DeletePlayerSeasonParams struct {
	PlayerID   string `json:"player_id"`
	SeasonYear int64  `json:"season_year"`
}

func (q *Queries) DeletePlayerSeason(ctx context.Context, arg DeletePlayerSeasonParams) error {
	_, err := q.exec(ctx, q.deletePlayerSeasonStmt, deletePlayerSeason, arg.PlayerID, arg.SeasonYear)
	return err
}

const getActivePlayerSeasonsByYear = `-- name: GetActivePlayerSeasonsByYear :many
SELECT ps.player_id, ps.season_year, ps.team_id, ps.jersey, ps.active, ps.experience, ps.status
FROM nfl_player_seasons ps
WHERE ps.active = true AND ps.season_year = ?
ORDER BY ps.player_id
`

func (q *Queries) GetActivePlayerSeasonsByYear(ctx context.Context, seasonYear int64) ([]*NflPlayerSeason, error) {
	rows, err := q.query(ctx, q.getActivePlayerSeasonsByYearStmt, getActivePlayerSeasonsByYear, seasonYear)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*NflPlayerSeason{}
	for rows.Next() {
		var i NflPlayerSeason
		if err := rows.Scan(
			&i.PlayerID,
			&i.SeasonYear,
			&i.TeamID,
			&i.Jersey,
			&i.Active,
			&i.Experience,
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

const getAllPlayerSeasons = `-- name: GetAllPlayerSeasons :many
SELECT player_id, season_year, team_id, jersey, active, experience, status FROM nfl_player_seasons
ORDER BY season_year DESC, player_id
`

func (q *Queries) GetAllPlayerSeasons(ctx context.Context) ([]*NflPlayerSeason, error) {
	rows, err := q.query(ctx, q.getAllPlayerSeasonsStmt, getAllPlayerSeasons)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*NflPlayerSeason{}
	for rows.Next() {
		var i NflPlayerSeason
		if err := rows.Scan(
			&i.PlayerID,
			&i.SeasonYear,
			&i.TeamID,
			&i.Jersey,
			&i.Active,
			&i.Experience,
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

const getPlayerSeason = `-- name: GetPlayerSeason :one
SELECT player_id, season_year, team_id, jersey, active, experience, status FROM nfl_player_seasons
WHERE player_id = ? AND season_year = ?
`

type GetPlayerSeasonParams struct {
	PlayerID   string `json:"player_id"`
	SeasonYear int64  `json:"season_year"`
}

func (q *Queries) GetPlayerSeason(ctx context.Context, arg GetPlayerSeasonParams) (*NflPlayerSeason, error) {
	row := q.queryRow(ctx, q.getPlayerSeasonStmt, getPlayerSeason, arg.PlayerID, arg.SeasonYear)
	var i NflPlayerSeason
	err := row.Scan(
		&i.PlayerID,
		&i.SeasonYear,
		&i.TeamID,
		&i.Jersey,
		&i.Active,
		&i.Experience,
		&i.Status,
	)
	return &i, err
}

const getPlayerSeasonsByTeam = `-- name: GetPlayerSeasonsByTeam :many
SELECT ps.player_id, ps.season_year, ps.team_id, ps.jersey, ps.active, ps.experience, ps.status
FROM nfl_player_seasons ps
JOIN nfl_players p ON ps.player_id = p.player_id
WHERE ps.team_id = ? AND ps.season_year = ?
ORDER BY p.position, p.last_name, p.first_name
`

type GetPlayerSeasonsByTeamParams struct {
	TeamID     sql.NullString `json:"team_id"`
	SeasonYear int64          `json:"season_year"`
}

func (q *Queries) GetPlayerSeasonsByTeam(ctx context.Context, arg GetPlayerSeasonsByTeamParams) ([]*NflPlayerSeason, error) {
	rows, err := q.query(ctx, q.getPlayerSeasonsByTeamStmt, getPlayerSeasonsByTeam, arg.TeamID, arg.SeasonYear)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*NflPlayerSeason{}
	for rows.Next() {
		var i NflPlayerSeason
		if err := rows.Scan(
			&i.PlayerID,
			&i.SeasonYear,
			&i.TeamID,
			&i.Jersey,
			&i.Active,
			&i.Experience,
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

const getPlayerSeasonsByYear = `-- name: GetPlayerSeasonsByYear :many
SELECT player_id, season_year, team_id, jersey, active, experience, status FROM nfl_player_seasons
WHERE season_year = ?
ORDER BY player_id
`

func (q *Queries) GetPlayerSeasonsByYear(ctx context.Context, seasonYear int64) ([]*NflPlayerSeason, error) {
	rows, err := q.query(ctx, q.getPlayerSeasonsByYearStmt, getPlayerSeasonsByYear, seasonYear)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*NflPlayerSeason{}
	for rows.Next() {
		var i NflPlayerSeason
		if err := rows.Scan(
			&i.PlayerID,
			&i.SeasonYear,
			&i.TeamID,
			&i.Jersey,
			&i.Active,
			&i.Experience,
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

const updatePlayerSeason = `-- name: UpdatePlayerSeason :exec
UPDATE nfl_player_seasons
SET team_id = ?,
    jersey = ?,
    active = ?,
    experience = ?,
    status = ?
WHERE player_id = ? AND season_year = ?
`

type UpdatePlayerSeasonParams struct {
	TeamID     sql.NullString `json:"team_id"`
	Jersey     sql.NullString `json:"jersey"`
	Active     bool           `json:"active"`
	Experience sql.NullInt64  `json:"experience"`
	Status     sql.NullString `json:"status"`
	PlayerID   string         `json:"player_id"`
	SeasonYear int64          `json:"season_year"`
}

func (q *Queries) UpdatePlayerSeason(ctx context.Context, arg UpdatePlayerSeasonParams) error {
	_, err := q.exec(ctx, q.updatePlayerSeasonStmt, updatePlayerSeason,
		arg.TeamID,
		arg.Jersey,
		arg.Active,
		arg.Experience,
		arg.Status,
		arg.PlayerID,
		arg.SeasonYear,
	)
	return err
}

const upsertPlayerSeason = `-- name: UpsertPlayerSeason :exec
INSERT INTO nfl_player_seasons (
  player_id, season_year, team_id, jersey, active, experience, status
) VALUES (
  ?, ?, ?, ?, ?, ?, ?
) ON CONFLICT(player_id, season_year) DO UPDATE SET
  team_id = excluded.team_id,
  jersey = excluded.jersey,
  active = excluded.active,
  experience = excluded.experience,
  status = excluded.status
`

type UpsertPlayerSeasonParams struct {
	PlayerID   string         `json:"player_id"`
	SeasonYear int64          `json:"season_year"`
	TeamID     sql.NullString `json:"team_id"`
	Jersey     sql.NullString `json:"jersey"`
	Active     bool           `json:"active"`
	Experience sql.NullInt64  `json:"experience"`
	Status     sql.NullString `json:"status"`
}

func (q *Queries) UpsertPlayerSeason(ctx context.Context, arg UpsertPlayerSeasonParams) error {
	_, err := q.exec(ctx, q.upsertPlayerSeasonStmt, upsertPlayerSeason,
		arg.PlayerID,
		arg.SeasonYear,
		arg.TeamID,
		arg.Jersey,
		arg.Active,
		arg.Experience,
		arg.Status,
	)
	return err
}

// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: player.sql

package data

import (
	"context"
	"database/sql"
)

const getAllPlayers = `-- name: GetAllPlayers :many
SELECT id, name, position, team_id, jersey_number, status
FROM nfl_players
ORDER BY name
`

func (q *Queries) GetAllPlayers(ctx context.Context) ([]*NflPlayer, error) {
	rows, err := q.query(ctx, q.getAllPlayersStmt, getAllPlayers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*NflPlayer{}
	for rows.Next() {
		var i NflPlayer
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Position,
			&i.TeamID,
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

const getPlayerById = `-- name: GetPlayerById :one
SELECT id, name, position, team_id, jersey_number, status
FROM nfl_players
WHERE id = ?
`

func (q *Queries) GetPlayerById(ctx context.Context, id int64) (*NflPlayer, error) {
	row := q.queryRow(ctx, q.getPlayerByIdStmt, getPlayerById, id)
	var i NflPlayer
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Position,
		&i.TeamID,
		&i.JerseyNumber,
		&i.Status,
	)
	return &i, err
}

const getPlayerFantasyPoints = `-- name: GetPlayerFantasyPoints :many
SELECT 
    fps.week,
    fps.points as fantasy_points,
    ps.id, ps.player_id, ps.game_id, ps.season_id, ps.week, ps.passing_attempts, ps.passing_completions, ps.passing_yards, ps.passing_touchdowns, ps.passing_interceptions, ps.rushing_attempts, ps.rushing_yards, ps.rushing_touchdowns, ps.targets, ps.receptions, ps.receiving_yards, ps.receiving_touchdowns, ps.field_goals_made, ps.field_goals_attempted, ps.extra_points_made, ps.extra_points_attempted, ps.sacks, ps.interceptions, ps.fumble_recoveries, ps.defensive_touchdowns, ps.safeties, ps.fumbles_lost, ps.two_point_conversions
FROM fantasy_player_scores fps
JOIN player_stats ps ON fps.player_id = ps.player_id AND fps.week = ps.week AND fps.season_id = ps.season_id
WHERE fps.player_id = ? AND fps.season_id = ? AND fps.league_id = ?
ORDER BY fps.week
`

type GetPlayerFantasyPointsParams struct {
	PlayerID int64 `json:"player_id"`
	SeasonID int64 `json:"season_id"`
	LeagueID int64 `json:"league_id"`
}

type GetPlayerFantasyPointsRow struct {
	Week                 int64           `json:"week"`
	FantasyPoints        sql.NullFloat64 `json:"fantasy_points"`
	ID                   int64           `json:"id"`
	PlayerID             int64           `json:"player_id"`
	GameID               int64           `json:"game_id"`
	SeasonID             int64           `json:"season_id"`
	Week_2               int64           `json:"week_2"`
	PassingAttempts      sql.NullInt64   `json:"passing_attempts"`
	PassingCompletions   sql.NullInt64   `json:"passing_completions"`
	PassingYards         sql.NullInt64   `json:"passing_yards"`
	PassingTouchdowns    sql.NullInt64   `json:"passing_touchdowns"`
	PassingInterceptions sql.NullInt64   `json:"passing_interceptions"`
	RushingAttempts      sql.NullInt64   `json:"rushing_attempts"`
	RushingYards         sql.NullInt64   `json:"rushing_yards"`
	RushingTouchdowns    sql.NullInt64   `json:"rushing_touchdowns"`
	Targets              sql.NullInt64   `json:"targets"`
	Receptions           sql.NullInt64   `json:"receptions"`
	ReceivingYards       sql.NullInt64   `json:"receiving_yards"`
	ReceivingTouchdowns  sql.NullInt64   `json:"receiving_touchdowns"`
	FieldGoalsMade       sql.NullInt64   `json:"field_goals_made"`
	FieldGoalsAttempted  sql.NullInt64   `json:"field_goals_attempted"`
	ExtraPointsMade      sql.NullInt64   `json:"extra_points_made"`
	ExtraPointsAttempted sql.NullInt64   `json:"extra_points_attempted"`
	Sacks                sql.NullFloat64 `json:"sacks"`
	Interceptions        sql.NullInt64   `json:"interceptions"`
	FumbleRecoveries     sql.NullInt64   `json:"fumble_recoveries"`
	DefensiveTouchdowns  sql.NullInt64   `json:"defensive_touchdowns"`
	Safeties             sql.NullInt64   `json:"safeties"`
	FumblesLost          sql.NullInt64   `json:"fumbles_lost"`
	TwoPointConversions  sql.NullInt64   `json:"two_point_conversions"`
}

func (q *Queries) GetPlayerFantasyPoints(ctx context.Context, arg GetPlayerFantasyPointsParams) ([]*GetPlayerFantasyPointsRow, error) {
	rows, err := q.query(ctx, q.getPlayerFantasyPointsStmt, getPlayerFantasyPoints, arg.PlayerID, arg.SeasonID, arg.LeagueID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetPlayerFantasyPointsRow{}
	for rows.Next() {
		var i GetPlayerFantasyPointsRow
		if err := rows.Scan(
			&i.Week,
			&i.FantasyPoints,
			&i.ID,
			&i.PlayerID,
			&i.GameID,
			&i.SeasonID,
			&i.Week_2,
			&i.PassingAttempts,
			&i.PassingCompletions,
			&i.PassingYards,
			&i.PassingTouchdowns,
			&i.PassingInterceptions,
			&i.RushingAttempts,
			&i.RushingYards,
			&i.RushingTouchdowns,
			&i.Targets,
			&i.Receptions,
			&i.ReceivingYards,
			&i.ReceivingTouchdowns,
			&i.FieldGoalsMade,
			&i.FieldGoalsAttempted,
			&i.ExtraPointsMade,
			&i.ExtraPointsAttempted,
			&i.Sacks,
			&i.Interceptions,
			&i.FumbleRecoveries,
			&i.DefensiveTouchdowns,
			&i.Safeties,
			&i.FumblesLost,
			&i.TwoPointConversions,
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

const getPlayerFantasyTotalPoints = `-- name: GetPlayerFantasyTotalPoints :one
SELECT 
    p.id,
    p.name,
    p.position,
    SUM(fps.points) as total_fantasy_points,
    AVG(fps.points) as avg_fantasy_points_per_game,
    COUNT(fps.week) as games_played
FROM nfl_players p
JOIN fantasy_player_scores fps ON p.id = fps.player_id
WHERE fps.season_id = ? AND fps.league_id = ?
GROUP BY p.id, p.name, p.position
`

type GetPlayerFantasyTotalPointsParams struct {
	SeasonID int64 `json:"season_id"`
	LeagueID int64 `json:"league_id"`
}

type GetPlayerFantasyTotalPointsRow struct {
	ID                      int64           `json:"id"`
	Name                    string          `json:"name"`
	Position                string          `json:"position"`
	TotalFantasyPoints      sql.NullFloat64 `json:"total_fantasy_points"`
	AvgFantasyPointsPerGame sql.NullFloat64 `json:"avg_fantasy_points_per_game"`
	GamesPlayed             int64           `json:"games_played"`
}

func (q *Queries) GetPlayerFantasyTotalPoints(ctx context.Context, arg GetPlayerFantasyTotalPointsParams) (*GetPlayerFantasyTotalPointsRow, error) {
	row := q.queryRow(ctx, q.getPlayerFantasyTotalPointsStmt, getPlayerFantasyTotalPoints, arg.SeasonID, arg.LeagueID)
	var i GetPlayerFantasyTotalPointsRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Position,
		&i.TotalFantasyPoints,
		&i.AvgFantasyPointsPerGame,
		&i.GamesPlayed,
	)
	return &i, err
}

const getPlayerSeasonStats = `-- name: GetPlayerSeasonStats :one
SELECT 
    SUM(passing_attempts) as total_passing_attempts,
    SUM(passing_completions) as total_passing_completions,
    SUM(passing_yards) as total_passing_yards,
    SUM(passing_touchdowns) as total_passing_touchdowns,
    SUM(passing_interceptions) as total_passing_interceptions,
    SUM(rushing_attempts) as total_rushing_attempts,
    SUM(rushing_yards) as total_rushing_yards,
    SUM(rushing_touchdowns) as total_rushing_touchdowns,
    SUM(targets) as total_targets,
    SUM(receptions) as total_receptions,
    SUM(receiving_yards) as total_receiving_yards,
    SUM(receiving_touchdowns) as total_receiving_touchdowns,
    SUM(field_goals_made) as total_field_goals_made,
    SUM(field_goals_attempted) as total_field_goals_attempted,
    SUM(extra_points_made) as total_extra_points_made,
    SUM(extra_points_attempted) as total_extra_points_attempted,
    SUM(sacks) as total_sacks,
    SUM(interceptions) as total_interceptions,
    SUM(fumble_recoveries) as total_fumble_recoveries,
    SUM(defensive_touchdowns) as total_defensive_touchdowns,
    SUM(safeties) as total_safeties,
    SUM(fumbles_lost) as total_fumbles_lost,
    SUM(two_point_conversions) as total_two_point_conversions
FROM player_stats
WHERE player_id = ? AND season_id = ?
`

type GetPlayerSeasonStatsParams struct {
	PlayerID int64 `json:"player_id"`
	SeasonID int64 `json:"season_id"`
}

type GetPlayerSeasonStatsRow struct {
	TotalPassingAttempts      sql.NullFloat64 `json:"total_passing_attempts"`
	TotalPassingCompletions   sql.NullFloat64 `json:"total_passing_completions"`
	TotalPassingYards         sql.NullFloat64 `json:"total_passing_yards"`
	TotalPassingTouchdowns    sql.NullFloat64 `json:"total_passing_touchdowns"`
	TotalPassingInterceptions sql.NullFloat64 `json:"total_passing_interceptions"`
	TotalRushingAttempts      sql.NullFloat64 `json:"total_rushing_attempts"`
	TotalRushingYards         sql.NullFloat64 `json:"total_rushing_yards"`
	TotalRushingTouchdowns    sql.NullFloat64 `json:"total_rushing_touchdowns"`
	TotalTargets              sql.NullFloat64 `json:"total_targets"`
	TotalReceptions           sql.NullFloat64 `json:"total_receptions"`
	TotalReceivingYards       sql.NullFloat64 `json:"total_receiving_yards"`
	TotalReceivingTouchdowns  sql.NullFloat64 `json:"total_receiving_touchdowns"`
	TotalFieldGoalsMade       sql.NullFloat64 `json:"total_field_goals_made"`
	TotalFieldGoalsAttempted  sql.NullFloat64 `json:"total_field_goals_attempted"`
	TotalExtraPointsMade      sql.NullFloat64 `json:"total_extra_points_made"`
	TotalExtraPointsAttempted sql.NullFloat64 `json:"total_extra_points_attempted"`
	TotalSacks                sql.NullFloat64 `json:"total_sacks"`
	TotalInterceptions        sql.NullFloat64 `json:"total_interceptions"`
	TotalFumbleRecoveries     sql.NullFloat64 `json:"total_fumble_recoveries"`
	TotalDefensiveTouchdowns  sql.NullFloat64 `json:"total_defensive_touchdowns"`
	TotalSafeties             sql.NullFloat64 `json:"total_safeties"`
	TotalFumblesLost          sql.NullFloat64 `json:"total_fumbles_lost"`
	TotalTwoPointConversions  sql.NullFloat64 `json:"total_two_point_conversions"`
}

func (q *Queries) GetPlayerSeasonStats(ctx context.Context, arg GetPlayerSeasonStatsParams) (*GetPlayerSeasonStatsRow, error) {
	row := q.queryRow(ctx, q.getPlayerSeasonStatsStmt, getPlayerSeasonStats, arg.PlayerID, arg.SeasonID)
	var i GetPlayerSeasonStatsRow
	err := row.Scan(
		&i.TotalPassingAttempts,
		&i.TotalPassingCompletions,
		&i.TotalPassingYards,
		&i.TotalPassingTouchdowns,
		&i.TotalPassingInterceptions,
		&i.TotalRushingAttempts,
		&i.TotalRushingYards,
		&i.TotalRushingTouchdowns,
		&i.TotalTargets,
		&i.TotalReceptions,
		&i.TotalReceivingYards,
		&i.TotalReceivingTouchdowns,
		&i.TotalFieldGoalsMade,
		&i.TotalFieldGoalsAttempted,
		&i.TotalExtraPointsMade,
		&i.TotalExtraPointsAttempted,
		&i.TotalSacks,
		&i.TotalInterceptions,
		&i.TotalFumbleRecoveries,
		&i.TotalDefensiveTouchdowns,
		&i.TotalSafeties,
		&i.TotalFumblesLost,
		&i.TotalTwoPointConversions,
	)
	return &i, err
}

const getPlayerStats = `-- name: GetPlayerStats :many
SELECT 
    ps.id, ps.player_id, ps.game_id, ps.season_id, ps.week, ps.passing_attempts, ps.passing_completions, ps.passing_yards, ps.passing_touchdowns, ps.passing_interceptions, ps.rushing_attempts, ps.rushing_yards, ps.rushing_touchdowns, ps.targets, ps.receptions, ps.receiving_yards, ps.receiving_touchdowns, ps.field_goals_made, ps.field_goals_attempted, ps.extra_points_made, ps.extra_points_attempted, ps.sacks, ps.interceptions, ps.fumble_recoveries, ps.defensive_touchdowns, ps.safeties, ps.fumbles_lost, ps.two_point_conversions, 
    g.week,
    g.game_date,
    ht.name as home_team_name,
    at.name as away_team_name,
    g.home_score,
    g.away_score
FROM player_stats ps
JOIN nfl_games g ON ps.game_id = g.id
JOIN nfl_teams ht ON g.home_team_id = ht.id
JOIN nfl_teams at ON g.away_team_id = at.id
WHERE ps.player_id = ? AND ps.season_id = ?
ORDER BY g.week
`

type GetPlayerStatsParams struct {
	PlayerID int64 `json:"player_id"`
	SeasonID int64 `json:"season_id"`
}

type GetPlayerStatsRow struct {
	ID                   int64           `json:"id"`
	PlayerID             int64           `json:"player_id"`
	GameID               int64           `json:"game_id"`
	SeasonID             int64           `json:"season_id"`
	Week                 int64           `json:"week"`
	PassingAttempts      sql.NullInt64   `json:"passing_attempts"`
	PassingCompletions   sql.NullInt64   `json:"passing_completions"`
	PassingYards         sql.NullInt64   `json:"passing_yards"`
	PassingTouchdowns    sql.NullInt64   `json:"passing_touchdowns"`
	PassingInterceptions sql.NullInt64   `json:"passing_interceptions"`
	RushingAttempts      sql.NullInt64   `json:"rushing_attempts"`
	RushingYards         sql.NullInt64   `json:"rushing_yards"`
	RushingTouchdowns    sql.NullInt64   `json:"rushing_touchdowns"`
	Targets              sql.NullInt64   `json:"targets"`
	Receptions           sql.NullInt64   `json:"receptions"`
	ReceivingYards       sql.NullInt64   `json:"receiving_yards"`
	ReceivingTouchdowns  sql.NullInt64   `json:"receiving_touchdowns"`
	FieldGoalsMade       sql.NullInt64   `json:"field_goals_made"`
	FieldGoalsAttempted  sql.NullInt64   `json:"field_goals_attempted"`
	ExtraPointsMade      sql.NullInt64   `json:"extra_points_made"`
	ExtraPointsAttempted sql.NullInt64   `json:"extra_points_attempted"`
	Sacks                sql.NullFloat64 `json:"sacks"`
	Interceptions        sql.NullInt64   `json:"interceptions"`
	FumbleRecoveries     sql.NullInt64   `json:"fumble_recoveries"`
	DefensiveTouchdowns  sql.NullInt64   `json:"defensive_touchdowns"`
	Safeties             sql.NullInt64   `json:"safeties"`
	FumblesLost          sql.NullInt64   `json:"fumbles_lost"`
	TwoPointConversions  sql.NullInt64   `json:"two_point_conversions"`
	Week_2               int64           `json:"week_2"`
	GameDate             string          `json:"game_date"`
	HomeTeamName         string          `json:"home_team_name"`
	AwayTeamName         string          `json:"away_team_name"`
	HomeScore            sql.NullInt64   `json:"home_score"`
	AwayScore            sql.NullInt64   `json:"away_score"`
}

func (q *Queries) GetPlayerStats(ctx context.Context, arg GetPlayerStatsParams) ([]*GetPlayerStatsRow, error) {
	rows, err := q.query(ctx, q.getPlayerStatsStmt, getPlayerStats, arg.PlayerID, arg.SeasonID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetPlayerStatsRow{}
	for rows.Next() {
		var i GetPlayerStatsRow
		if err := rows.Scan(
			&i.ID,
			&i.PlayerID,
			&i.GameID,
			&i.SeasonID,
			&i.Week,
			&i.PassingAttempts,
			&i.PassingCompletions,
			&i.PassingYards,
			&i.PassingTouchdowns,
			&i.PassingInterceptions,
			&i.RushingAttempts,
			&i.RushingYards,
			&i.RushingTouchdowns,
			&i.Targets,
			&i.Receptions,
			&i.ReceivingYards,
			&i.ReceivingTouchdowns,
			&i.FieldGoalsMade,
			&i.FieldGoalsAttempted,
			&i.ExtraPointsMade,
			&i.ExtraPointsAttempted,
			&i.Sacks,
			&i.Interceptions,
			&i.FumbleRecoveries,
			&i.DefensiveTouchdowns,
			&i.Safeties,
			&i.FumblesLost,
			&i.TwoPointConversions,
			&i.Week_2,
			&i.GameDate,
			&i.HomeTeamName,
			&i.AwayTeamName,
			&i.HomeScore,
			&i.AwayScore,
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

const getPlayersByPosition = `-- name: GetPlayersByPosition :many
SELECT id, name, position, team_id, jersey_number, status
FROM nfl_players
WHERE position = ?
ORDER BY name
`

func (q *Queries) GetPlayersByPosition(ctx context.Context, position string) ([]*NflPlayer, error) {
	rows, err := q.query(ctx, q.getPlayersByPositionStmt, getPlayersByPosition, position)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*NflPlayer{}
	for rows.Next() {
		var i NflPlayer
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Position,
			&i.TeamID,
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

const getPlayersByTeam = `-- name: GetPlayersByTeam :many
SELECT p.id, p.name, p.position, p.team_id, p.jersey_number, p.status
FROM nfl_players p
JOIN nfl_teams t ON p.team_id = t.id
WHERE t.id = ?
ORDER BY p.position, p.name
`

func (q *Queries) GetPlayersByTeam(ctx context.Context, id int64) ([]*NflPlayer, error) {
	rows, err := q.query(ctx, q.getPlayersByTeamStmt, getPlayersByTeam, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*NflPlayer{}
	for rows.Next() {
		var i NflPlayer
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Position,
			&i.TeamID,
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

const getTopPlayersByPositionAndSeason = `-- name: GetTopPlayersByPositionAndSeason :many
SELECT 
    p.id,
    p.name,
    p.position,
    p.team_id,
    t.name as team_name,
    SUM(fps.points) as total_fantasy_points,
    AVG(fps.points) as avg_fantasy_points_per_game,
    COUNT(fps.week) as games_played
FROM nfl_players p
JOIN fantasy_player_scores fps ON p.id = fps.player_id
JOIN nfl_teams t ON p.team_id = t.id
WHERE fps.season_id = ? AND fps.league_id = ? AND p.position = ?
GROUP BY p.id, p.name, p.position, p.team_id, t.name
ORDER BY total_fantasy_points DESC
LIMIT ?
`

type GetTopPlayersByPositionAndSeasonParams struct {
	SeasonID int64  `json:"season_id"`
	LeagueID int64  `json:"league_id"`
	Position string `json:"position"`
	Limit    int64  `json:"limit"`
}

type GetTopPlayersByPositionAndSeasonRow struct {
	ID                      int64           `json:"id"`
	Name                    string          `json:"name"`
	Position                string          `json:"position"`
	TeamID                  sql.NullInt64   `json:"team_id"`
	TeamName                string          `json:"team_name"`
	TotalFantasyPoints      sql.NullFloat64 `json:"total_fantasy_points"`
	AvgFantasyPointsPerGame sql.NullFloat64 `json:"avg_fantasy_points_per_game"`
	GamesPlayed             int64           `json:"games_played"`
}

func (q *Queries) GetTopPlayersByPositionAndSeason(ctx context.Context, arg GetTopPlayersByPositionAndSeasonParams) ([]*GetTopPlayersByPositionAndSeasonRow, error) {
	rows, err := q.query(ctx, q.getTopPlayersByPositionAndSeasonStmt, getTopPlayersByPositionAndSeason,
		arg.SeasonID,
		arg.LeagueID,
		arg.Position,
		arg.Limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetTopPlayersByPositionAndSeasonRow{}
	for rows.Next() {
		var i GetTopPlayersByPositionAndSeasonRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Position,
			&i.TeamID,
			&i.TeamName,
			&i.TotalFantasyPoints,
			&i.AvgFantasyPointsPerGame,
			&i.GamesPlayed,
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

const insertPlayer = `-- name: InsertPlayer :one
INSERT INTO nfl_players (name, position, team_id, jersey_number, status)
VALUES (?, ?, ?, ?, ?)
RETURNING id
`

type InsertPlayerParams struct {
	Name         string         `json:"name"`
	Position     string         `json:"position"`
	TeamID       sql.NullInt64  `json:"team_id"`
	JerseyNumber sql.NullInt64  `json:"jersey_number"`
	Status       sql.NullString `json:"status"`
}

func (q *Queries) InsertPlayer(ctx context.Context, arg InsertPlayerParams) (int64, error) {
	row := q.queryRow(ctx, q.insertPlayerStmt, insertPlayer,
		arg.Name,
		arg.Position,
		arg.TeamID,
		arg.JerseyNumber,
		arg.Status,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const searchPlayers = `-- name: SearchPlayers :many
SELECT id, name, position, team_id, jersey_number, status
FROM nfl_players
WHERE name LIKE '%' || ? || '%'
ORDER BY name
LIMIT 20
`

func (q *Queries) SearchPlayers(ctx context.Context, dollar_1 sql.NullString) ([]*NflPlayer, error) {
	rows, err := q.query(ctx, q.searchPlayersStmt, searchPlayers, dollar_1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*NflPlayer{}
	for rows.Next() {
		var i NflPlayer
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Position,
			&i.TeamID,
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

const updatePlayer = `-- name: UpdatePlayer :exec
UPDATE nfl_players
SET 
    name = ?,
    position = ?,
    team_id = ?,
    jersey_number = ?,
    status = ?
WHERE id = ?
`

type UpdatePlayerParams struct {
	Name         string         `json:"name"`
	Position     string         `json:"position"`
	TeamID       sql.NullInt64  `json:"team_id"`
	JerseyNumber sql.NullInt64  `json:"jersey_number"`
	Status       sql.NullString `json:"status"`
	ID           int64          `json:"id"`
}

func (q *Queries) UpdatePlayer(ctx context.Context, arg UpdatePlayerParams) error {
	_, err := q.exec(ctx, q.updatePlayerStmt, updatePlayer,
		arg.Name,
		arg.Position,
		arg.TeamID,
		arg.JerseyNumber,
		arg.Status,
		arg.ID,
	)
	return err
}

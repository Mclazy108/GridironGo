-- name: GetAllSeasons :many
SELECT id, year, start_date, end_date, current
FROM seasons
ORDER BY year DESC;

-- name: GetSeasonById :one
SELECT id, year, start_date, end_date, current
FROM seasons
WHERE id = ?;

-- name: GetSeasonByYear :one
SELECT id, year, start_date, end_date, current
FROM seasons
WHERE year = ?;

-- name: GetCurrentSeason :one
SELECT id, year, start_date, end_date, current
FROM seasons
WHERE current = 1
LIMIT 1;

-- name: GetPreviousSeasons :many
SELECT id, year, start_date, end_date, current
FROM seasons
WHERE year < (SELECT year FROM seasons WHERE current = 1)
ORDER BY year DESC
LIMIT ?;

-- name: GetNFLScheduleForSeason :many
SELECT 
    g.id,
    g.week,
    g.home_team_id,
    ht.name as home_team_name,
    ht.abbreviation as home_team_abbr,
    g.away_team_id,
    at.name as away_team_name,
    at.abbreviation as away_team_abbr,
    g.home_score,
    g.away_score,
    g.game_date,
    g.game_time,
    g.status
FROM nfl_games g
JOIN nfl_teams ht ON g.home_team_id = ht.id
JOIN nfl_teams at ON g.away_team_id = at.id
WHERE g.season_id = ?
ORDER BY g.week, g.game_date, g.game_time;

-- name: GetNFLScheduleForWeek :many
SELECT 
    g.id,
    g.week,
    g.home_team_id,
    ht.name as home_team_name,
    ht.abbreviation as home_team_abbr,
    g.away_team_id,
    at.name as away_team_name,
    at.abbreviation as away_team_abbr,
    g.home_score,
    g.away_score,
    g.game_date,
    g.game_time,
    g.status
FROM nfl_games g
JOIN nfl_teams ht ON g.home_team_id = ht.id
JOIN nfl_teams at ON g.away_team_id = at.id
WHERE g.season_id = ? AND g.week = ?
ORDER BY g.game_date, g.game_time;

-- name: CreateSeason :one
INSERT INTO seasons (
    year,
    start_date,
    end_date,
    current
) VALUES (?, ?, ?, ?)
RETURNING id;

-- name: UpdateSeason :exec
UPDATE seasons
SET 
    start_date = ?,
    end_date = ?,
    current = ?
WHERE id = ?;

-- name: SetCurrentSeason :exec
UPDATE seasons
SET current = 0
WHERE current = 1;

UPDATE seasons
SET current = 1
WHERE id = ?;

-- name: CreateNFLGame :one
INSERT INTO nfl_games (
    season_id,
    week,
    home_team_id,
    away_team_id,
    home_score,
    away_score,
    game_date,
    game_time,
    status
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING id;

-- name: UpdateNFLGameScore :exec
UPDATE nfl_games
SET
    home_score = ?,
    away_score = ?,
    status = ?
WHERE id = ?;

-- name: GetTotalWeeksInSeason :one
SELECT MAX(week) as total_weeks
FROM nfl_games
WHERE season_id = ?;

-- name: GetHistoricalPlayerStats :many
SELECT 
    s.year,
    COUNT(DISTINCT ps.game_id) as games_played,
    SUM(ps.passing_yards) as passing_yards,
    SUM(ps.passing_touchdowns) as passing_touchdowns,
    SUM(ps.rushing_yards) as rushing_yards,
    SUM(ps.rushing_touchdowns) as rushing_touchdowns,
    SUM(ps.receptions) as receptions,
    SUM(ps.receiving_yards) as receiving_yards,
    SUM(ps.receiving_touchdowns) as receiving_touchdowns
FROM player_stats ps
JOIN seasons s ON ps.season_id = s.id
WHERE ps.player_id = ?
GROUP BY s.year
ORDER BY s.year DESC;

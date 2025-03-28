-- name: GetAllPlayers :many
SELECT id, name, position, team_id, jersey_number, status
FROM nfl_players
ORDER BY name;

-- name: GetPlayersByPosition :many
SELECT id, name, position, team_id, jersey_number, status
FROM nfl_players
WHERE position = ?
ORDER BY name;

-- name: GetPlayersByTeam :many
SELECT p.id, p.name, p.position, p.team_id, p.jersey_number, p.status
FROM nfl_players p
JOIN nfl_teams t ON p.team_id = t.id
WHERE t.id = ?
ORDER BY p.position, p.name;

-- name: GetPlayerById :one
SELECT id, name, position, team_id, jersey_number, status
FROM nfl_players
WHERE id = ?;

-- name: GetPlayerStats :many
SELECT 
    ps.*, 
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
ORDER BY g.week;

-- name: GetPlayerSeasonStats :one
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
WHERE player_id = ? AND season_id = ?;

-- name: SearchPlayers :many
SELECT id, name, position, team_id, jersey_number, status
FROM nfl_players
WHERE name LIKE '%' || ? || '%'
ORDER BY name
LIMIT 20;

-- name: InsertPlayer :one
INSERT INTO nfl_players (name, position, team_id, jersey_number, status)
VALUES (?, ?, ?, ?, ?)
RETURNING id;

-- name: UpdatePlayer :exec
UPDATE nfl_players
SET 
    name = ?,
    position = ?,
    team_id = ?,
    jersey_number = ?,
    status = ?
WHERE id = ?;

-- name: GetPlayerFantasyPoints :many
SELECT 
    fps.week,
    fps.points as fantasy_points,
    ps.*
FROM fantasy_player_scores fps
JOIN player_stats ps ON fps.player_id = ps.player_id AND fps.week = ps.week AND fps.season_id = ps.season_id
WHERE fps.player_id = ? AND fps.season_id = ? AND fps.league_id = ?
ORDER BY fps.week;

-- name: GetPlayerFantasyTotalPoints :one
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
GROUP BY p.id, p.name, p.position;

-- name: GetTopPlayersByPositionAndSeason :many
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
LIMIT ?;

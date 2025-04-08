-- name: CreateNFLStat :exec
INSERT INTO nfl_stats (
  game_id, player_id, team_id, category, stat_type, stat_value
) VALUES (
  ?, ?, ?, ?, ?, ?
);

-- name: GetStatsByGame :many
SELECT * FROM nfl_stats
WHERE game_id = ?
ORDER BY player_id, category, stat_type;

-- name: GetStatsByPlayer :many
SELECT * FROM nfl_stats
WHERE player_id = ?
ORDER BY game_id, category, stat_type;

-- name: GetStatsByTeam :many
SELECT * FROM nfl_stats
WHERE team_id = ?
ORDER BY game_id, player_id, category, stat_type;

-- name: GetStatsByGameAndPlayer :many
SELECT * FROM nfl_stats
WHERE game_id = ? AND player_id = ?
ORDER BY category, stat_type;

-- name: GetStatsByCategory :many
SELECT * FROM nfl_stats
WHERE category = ?
ORDER BY game_id, player_id, stat_type;

-- name: GetStatsByStatType :many
SELECT * FROM nfl_stats
WHERE stat_type = ?
ORDER BY game_id, player_id, category;

-- name: UpdateNFLStat :exec
UPDATE nfl_stats
SET stat_value = ?
WHERE stat_id = ?;

-- name: DeleteNFLStat :exec
DELETE FROM nfl_stats
WHERE stat_id = ?;

-- name: UpsertNFLStat :exec
INSERT INTO nfl_stats (
  game_id, player_id, team_id, category, stat_type, stat_value
) VALUES (
  ?, ?, ?, ?, ?, ?
) ON CONFLICT(stat_id) DO UPDATE SET
  stat_value = excluded.stat_value;

-- CORE AGGREGATION FUNCTIONS

-- name: GetPlayerTotalStatByType :one
-- Get the total of a specific stat type for a player
SELECT 
  SUM(stat_value) as total_value
FROM 
  nfl_stats
WHERE 
  player_id = ? AND stat_type = ?;

-- name: GetPlayerTotalStatByTypeForSeason :one
-- Get the total of a specific stat type for a player in a specific season
SELECT 
  SUM(s.stat_value) as total_value
FROM 
  nfl_stats s
JOIN 
  nfl_games g ON s.game_id = g.event_id
WHERE 
  s.player_id = ? AND s.stat_type = ? AND g.season = ?;

-- name: GetPlayerWeeklyStatByType :many
-- Get weekly stats of a specific type for a player in a season
SELECT 
  g.week,
  COALESCE(SUM(s.stat_value), 0) as stat_value
FROM 
  nfl_games g
LEFT JOIN 
  nfl_stats s ON g.event_id = s.game_id AND s.player_id = ? AND s.stat_type = ?
WHERE 
  g.season = ?
GROUP BY 
  g.week
ORDER BY 
  g.week;

-- name: GetPlayerStatAverage :one
-- Get the average of a specific stat type for a player
SELECT 
  AVG(stat_value) as avg_value
FROM 
  nfl_stats
WHERE 
  player_id = ? AND stat_type = ?;

-- name: GetPlayerStatsByGame :many
-- Get all stats for a player in a specific game
SELECT 
  s.category,
  s.stat_type,
  s.stat_value
FROM 
  nfl_stats s
WHERE 
  s.player_id = ? AND s.game_id = ?
ORDER BY 
  s.category, s.stat_type;

-- name: GetTopPlayersByStat :many
-- Get top N players for a specific stat type in a season
SELECT 
  p.player_id,
  p.full_name,
  p.position,
  SUM(s.stat_value) as total_value
FROM 
  nfl_stats s
JOIN 
  nfl_players p ON s.player_id = p.player_id
JOIN 
  nfl_games g ON s.game_id = g.event_id
WHERE 
  s.stat_type = ? AND g.season = ?
GROUP BY 
  p.player_id, p.full_name, p.position
ORDER BY 
  total_value DESC
LIMIT ?;

-- name: GetTeamStatsByCategory :many
-- Get team-level stats for a category
SELECT 
  t.team_id,
  t.display_name,
  s.stat_type,
  SUM(s.stat_value) as total_value,
  AVG(s.stat_value) as avg_value_per_game
FROM 
  nfl_stats s
JOIN 
  nfl_teams t ON s.team_id = t.team_id
JOIN 
  nfl_games g ON s.game_id = g.event_id
WHERE 
  s.category = ? AND g.season = ?
GROUP BY 
  t.team_id, t.display_name, s.stat_type
ORDER BY 
  s.stat_type, total_value DESC;

-- name: GetPlayerTotalStatsByPosition :many
-- Get total stats for all players of a specific position in a season
SELECT 
  p.player_id,
  p.full_name,
  s.stat_type,
  SUM(s.stat_value) as total_value
FROM 
  nfl_stats s
JOIN 
  nfl_players p ON s.player_id = p.player_id
JOIN 
  nfl_games g ON s.game_id = g.event_id
WHERE 
  p.position = ? AND g.season = ?
GROUP BY 
  p.player_id, p.full_name, s.stat_type
ORDER BY 
  s.stat_type, total_value DESC;

-- name: GetPlayerStatsByWeek :many
-- Get a player's stats for a specific week in a season
SELECT 
  s.category,
  s.stat_type,
  s.stat_value
FROM 
  nfl_stats s
JOIN 
  nfl_games g ON s.game_id = g.event_id
WHERE 
  s.player_id = ? AND g.season = ? AND g.week = ?
ORDER BY 
  s.category, s.stat_type;

-- name: GetPlayerTotalStatsBySeason :many
-- Get a player's total stats for each stat type in a season
SELECT 
  s.stat_type,
  SUM(s.stat_value) as total_value
FROM 
  nfl_stats s
JOIN 
  nfl_games g ON s.game_id = g.event_id
WHERE 
  s.player_id = ? AND g.season = ?
GROUP BY 
  s.stat_type
ORDER BY 
  s.stat_type;

-- name: CreateGame :exec
INSERT INTO nfl_games (
  event_id, date, name, short_name, season, week, away_team, home_team
) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?
);

-- name: GetGame :one
SELECT * FROM nfl_games
WHERE event_id = ?;

-- name: GetAllGames :many
SELECT * FROM nfl_games
ORDER BY date DESC;

-- name: DeleteGame :exec
DELETE FROM nfl_games
WHERE event_id = ?;

-- name: UpdateGame :exec
UPDATE nfl_games
SET date = ?,
    name = ?,
    short_name = ?,
    season = ?,
    week = ?,
    away_team = ?,
    home_team = ?
WHERE event_id = ?;

-- name: GetAllGamesBySeasonAndWeek :many
SELECT * FROM nfl_games
WHERE season = ? AND week = ?
ORDER BY date ASC;

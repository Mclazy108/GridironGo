-- name: CreateNFLPlayer :exec
INSERT INTO nfl_players (
  player_id, first_name, last_name, full_name, position, team_id, jersey,
  height, weight, active, college, experience, draft_year, draft_round, draft_pick,
  status, image_url
) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
);

-- name: GetNFLPlayer :one
SELECT * FROM nfl_players
WHERE player_id = ?;

-- name: GetAllNFLPlayers :many
SELECT * FROM nfl_players
ORDER BY last_name, first_name ASC;

-- name: GetActiveNFLPlayers :many
SELECT * FROM nfl_players
WHERE active = true
ORDER BY last_name, first_name ASC;

-- name: DeleteNFLPlayer :exec
DELETE FROM nfl_players
WHERE player_id = ?;

-- name: UpdateNFLPlayer :exec
UPDATE nfl_players
SET first_name = ?,
    last_name = ?,
    full_name = ?,
    position = ?,
    team_id = ?,
    jersey = ?,
    height = ?,
    weight = ?,
    active = ?,
    college = ?,
    experience = ?,
    draft_year = ?,
    draft_round = ?,
    draft_pick = ?,
    status = ?,
    image_url = ?
WHERE player_id = ?;

-- name: GetPlayersByTeam :many
SELECT * FROM nfl_players
WHERE team_id = ?
ORDER BY position, last_name, first_name;

-- name: GetPlayersByPosition :many
SELECT * FROM nfl_players
WHERE position = ? AND active = true
ORDER BY last_name, first_name;

-- name: SearchPlayers :many
SELECT * FROM nfl_players
WHERE (full_name LIKE ? OR last_name LIKE ?) AND active = true
ORDER BY last_name, first_name
LIMIT 50;

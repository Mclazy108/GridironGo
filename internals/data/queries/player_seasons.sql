-- name: CreatePlayerSeason :exec
INSERT INTO nfl_player_seasons (
  player_id, season_year, team_id, jersey, active, experience, status
) VALUES (
  ?, ?, ?, ?, ?, ?, ?
);

-- name: GetPlayerSeason :one
SELECT * FROM nfl_player_seasons
WHERE player_id = ? AND season_year = ?;

-- name: GetAllPlayerSeasons :many
SELECT * FROM nfl_player_seasons
ORDER BY season_year DESC, player_id;

-- name: GetPlayerSeasonsByYear :many
SELECT * FROM nfl_player_seasons
WHERE season_year = ?
ORDER BY player_id;

-- name: GetPlayerSeasonsByTeam :many
SELECT ps.*
FROM nfl_player_seasons ps
JOIN nfl_players p ON ps.player_id = p.player_id
WHERE ps.team_id = ? AND ps.season_year = ?
ORDER BY p.position, p.last_name, p.first_name;

-- name: GetActivePlayerSeasonsByYear :many
SELECT ps.*
FROM nfl_player_seasons ps
WHERE ps.active = true AND ps.season_year = ?
ORDER BY ps.player_id;

-- name: DeletePlayerSeason :exec
DELETE FROM nfl_player_seasons
WHERE player_id = ? AND season_year = ?;

-- name: UpdatePlayerSeason :exec
UPDATE nfl_player_seasons
SET team_id = ?,
    jersey = ?,
    active = ?,
    experience = ?,
    status = ?
WHERE player_id = ? AND season_year = ?;

-- name: UpsertPlayerSeason :exec
INSERT INTO nfl_player_seasons (
  player_id, season_year, team_id, jersey, active, experience, status
) VALUES (
  ?, ?, ?, ?, ?, ?, ?
) ON CONFLICT(player_id, season_year) DO UPDATE SET
  team_id = excluded.team_id,
  jersey = excluded.jersey,
  active = excluded.active,
  experience = excluded.experience,
  status = excluded.status;

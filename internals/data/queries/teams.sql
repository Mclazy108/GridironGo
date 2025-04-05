-- name: CreateNFLTeam :exec
INSERT INTO nfl_teams (
  team_id, display_name, abbreviation, short_name, location, nickname,
  conference, division, primary_color, secondary_color, logo_url
) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
);

-- name: GetNFLTeam :one
SELECT * FROM nfl_teams
WHERE team_id = ?;

-- name: GetAllNFLTeams :many
SELECT * FROM nfl_teams
ORDER BY display_name ASC;

-- name: DeleteNFLTeam :exec
DELETE FROM nfl_teams
WHERE team_id = ?;

-- name: UpdateNFLTeam :exec
UPDATE nfl_teams
SET display_name = ?,
    abbreviation = ?,
    short_name = ?,
    location = ?,
    nickname = ?,
    conference = ?,
    division = ?,
    primary_color = ?,
    secondary_color = ?,
    logo_url = ?
WHERE team_id = ?;

-- name: GetTeamsByConference :many
SELECT * FROM nfl_teams
WHERE conference = ?
ORDER BY division, display_name;

-- name: GetTeamsByDivision :many
SELECT * FROM nfl_teams
WHERE division = ?
ORDER BY display_name;

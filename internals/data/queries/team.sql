-- name: GetAllNFLTeams :many
SELECT id, name, city, abbreviation, conference, division
FROM nfl_teams
ORDER BY conference, division, name;

-- name: GetNFLTeamById :one
SELECT id, name, city, abbreviation, conference, division
FROM nfl_teams
WHERE id = ?;

-- name: GetNFLTeamByAbbreviation :one
SELECT id, name, city, abbreviation, conference, division
FROM nfl_teams
WHERE abbreviation = ?;

-- name: GetTeamsByConference :many
SELECT id, name, city, abbreviation, conference, division
FROM nfl_teams
WHERE conference = ?
ORDER BY division, name;

-- name: GetTeamsByDivision :many
SELECT id, name, city, abbreviation, conference, division
FROM nfl_teams
WHERE division = ?
ORDER BY name;

-- name: GetTeamRoster :many
SELECT 
    p.id, 
    p.name, 
    p.position, 
    p.jersey_number, 
    p.status
FROM nfl_players p
WHERE p.team_id = ?
ORDER BY p.position, p.name;

-- name: GetAllFantasyTeams :many
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
ORDER BY (ft.wins * 2 + ft.tie_games) DESC, ft.points_for DESC;

-- name: GetFantasyTeamById :one
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
WHERE ft.id = ?;

-- name: GetFantasyTeamRoster :many
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
ORDER BY fr.is_starter DESC, fr.position, p.name;

-- name: GetFantasyTeamForWeek :one
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
GROUP BY ft.id, ft.name, ft.owner_name, ft.wins, ft.losses, ft.tie_games;

-- name: CreateFantasyTeam :one
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
RETURNING id;

-- name: UpdateFantasyTeamRecord :exec
UPDATE fantasy_teams
SET 
    wins = ?,
    losses = ?,
    tie_games = ?,
    points_for = ?,
    points_against = ?
WHERE id = ?;

-- name: AddPlayerToFantasyTeam :one
INSERT INTO fantasy_rosters (
    team_id,
    player_id,
    position,
    is_starter
) VALUES (?, ?, ?, ?)
RETURNING id;

-- name: RemovePlayerFromFantasyTeam :exec
DELETE FROM fantasy_rosters
WHERE team_id = ? AND player_id = ?;

-- name: UpdateFantasyRoster :exec
UPDATE fantasy_rosters
SET is_starter = ?
WHERE id = ?;

-- name: GetAvailablePlayers :many
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
LIMIT ?;

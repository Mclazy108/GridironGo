-- name: GetDraftPicksForLeague :many
SELECT 
    fd.id,
    fd.league_id,
    fd.team_id,
    ft.name as team_name,
    fd.player_id,
    p.name as player_name,
    p.position as player_position,
    t.abbreviation as team_abbreviation,
    fd.round,
    fd.pick_number,
    fd.draft_time
FROM fantasy_drafts fd
JOIN fantasy_teams ft ON fd.team_id = ft.id
JOIN nfl_players p ON fd.player_id = p.id
LEFT JOIN nfl_teams t ON p.team_id = t.id
WHERE fd.league_id = ?
ORDER BY fd.pick_number;

-- name: GetDraftPicksByTeam :many
SELECT 
    fd.id,
    fd.league_id,
    fd.team_id,
    fd.player_id,
    p.name as player_name,
    p.position as player_position,
    t.abbreviation as team_abbreviation,
    fd.round,
    fd.pick_number,
    fd.draft_time
FROM fantasy_drafts fd
JOIN nfl_players p ON fd.player_id = p.id
LEFT JOIN nfl_teams t ON p.team_id = t.id
WHERE fd.league_id = ? AND fd.team_id = ?
ORDER BY fd.pick_number;

-- name: GetLastDraftPick :one
SELECT 
    MAX(pick_number) as last_pick_number,
    (MAX(pick_number) / teams_count) as last_completed_round
FROM fantasy_drafts fd
JOIN fantasy_leagues fl ON fd.league_id = fl.id
WHERE fd.league_id = ?;

-- name: GetDraftOrder :many
SELECT
    ft.id as team_id,
    ft.name as team_name,
    ft.owner_name,
    ft.draft_position
FROM fantasy_teams ft
WHERE ft.league_id = ?
ORDER BY ft.draft_position;

-- name: GetTeamAtDraftPosition :one
SELECT 
    id,
    name,
    owner_name,
    is_user
FROM fantasy_teams
WHERE league_id = ? AND draft_position = ?;

-- name: GetNextDraftingTeam :one
WITH league_info AS (
    SELECT 
        fl.id,
        fl.teams_count,
        COUNT(fd.id) as picks_made
    FROM fantasy_leagues fl
    LEFT JOIN fantasy_drafts fd ON fl.id = fd.league_id
    WHERE fl.id = ?
    GROUP BY fl.id, fl.teams_count
),
draft_position AS (
    SELECT
        li.picks_made + 1 as next_pick_number,
        CEIL((li.picks_made + 1) * 1.0 / li.teams_count) as next_round,
        MOD(li.picks_made, li.teams_count) + 1 as position_in_round
    FROM league_info li
),
snake_position AS (
    SELECT
        dp.next_pick_number,
        dp.next_round,
        CASE 
            WHEN dp.next_round % 2 = 1 THEN dp.position_in_round
            ELSE li.teams_count - dp.position_in_round + 1
        END as snake_draft_position
    FROM draft_position dp, league_info li
)
SELECT
    ft.id,
    ft.name,
    ft.owner_name,
    ft.is_user,
    sp.next_pick_number as pick_number,
    sp.next_round as round
FROM snake_position sp
JOIN fantasy_teams ft ON ft.draft_position = sp.snake_draft_position
WHERE ft.league_id = ?;

-- name: GetAvailableDraftPlayers :many
SELECT 
    p.id,
    p.name,
    p.position,
    t.abbreviation as team_abbreviation,
    COALESCE(
        (SELECT SUM(fps.points) 
         FROM fantasy_player_scores fps 
         WHERE fps.player_id = p.id AND fps.league_id = ?1 AND fps.season_id = ?2), 
        0
    ) as previous_season_points,
    COALESCE(
        (SELECT AVG(fps.points) 
         FROM fantasy_player_scores fps 
         WHERE fps.player_id = p.id AND fps.league_id = ?1 AND fps.season_id = ?2), 
        0
    ) as avg_points_per_game
FROM nfl_players p
LEFT JOIN nfl_teams t ON p.team_id = t.id
WHERE p.id NOT IN (
    SELECT player_id FROM fantasy_drafts WHERE league_id = ?1
)
AND p.position = ?3
ORDER BY previous_season_points DESC
LIMIT ?4;

-- name: AddDraftPick :one
INSERT INTO fantasy_drafts (
    league_id,
    team_id,
    player_id,
    round,
    pick_number,
    draft_time
) VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
RETURNING id;

-- name: GetDraftSummary :many
SELECT 
    p.position,
    COUNT(*) as players_drafted
FROM fantasy_drafts fd
JOIN nfl_players p ON fd.player_id = p.id
WHERE fd.league_id = ? AND fd.team_id = ?
GROUP BY p.position
ORDER BY p.position;

-- name: ClearDraft :exec
DELETE FROM fantasy_drafts
WHERE league_id = ?;

-- name: GetAllPlayersWithFantasyPoints :many
SELECT 
    p.id,
    p.name,
    p.position,
    t.name as team_name,
    t.abbreviation as team_abbreviation,
    COALESCE(
        (SELECT SUM(fps.points) 
         FROM fantasy_player_scores fps 
         WHERE fps.player_id = p.id AND fps.league_id = ?1 AND fps.season_id = ?2), 
        0
    ) as previous_season_points,
    COALESCE(
        (SELECT AVG(fps.points) 
         FROM fantasy_player_scores fps 
         WHERE fps.player_id = p.id AND fps.league_id = ?1 AND fps.season_id = ?2 AND fps.points > 0), 
        0
    ) as avg_points_per_game,
    COALESCE(
        (SELECT COUNT(fps.week) 
         FROM fantasy_player_scores fps 
         WHERE fps.player_id = p.id AND fps.league_id = ?1 AND fps.season_id = ?2 AND fps.points > 0), 
        0
    ) as games_with_points
FROM nfl_players p
LEFT JOIN nfl_teams t ON p.team_id = t.id
WHERE p.position = ?3
ORDER BY previous_season_points DESC, avg_points_per_game DESC
LIMIT ?4;

-- name: GetBestAvailablePlayers :many
SELECT 
    p.id,
    p.name,
    p.position,
    t.abbreviation as team_abbreviation,
    COALESCE(
        (SELECT SUM(fps.points) 
         FROM fantasy_player_scores fps 
         WHERE fps.player_id = p.id AND fps.league_id = ?1 AND fps.season_id = ?2), 
        0
    ) as previous_season_points
FROM nfl_players p
LEFT JOIN nfl_teams t ON p.team_id = t.id
WHERE p.id NOT IN (
    SELECT player_id FROM fantasy_drafts WHERE league_id = ?1
)
AND p.position = ?3
ORDER BY p.position, previous_season_points DESC
LIMIT ?4;

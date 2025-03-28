-- name: GetPlayerStatsByGame :one
SELECT * FROM player_stats
WHERE player_id = ? AND game_id = ?;

-- name: GetPlayerStatsByWeek :many
SELECT 
    ps.*,
    p.name as player_name,
    p.position as player_position,
    ht.name as home_team_name,
    at.name as away_team_name
FROM player_stats ps
JOIN nfl_players p ON ps.player_id = p.id
JOIN nfl_games g ON ps.game_id = g.id
JOIN nfl_teams ht ON g.home_team_id = ht.id
JOIN nfl_teams at ON g.away_team_id = at.id
WHERE ps.season_id = ? AND ps.week = ? AND ps.player_id = ?;

-- name: CreatePlayerStats :one
INSERT INTO player_stats (
    player_id,
    game_id,
    season_id,
    week,
    passing_attempts,
    passing_completions,
    passing_yards,
    passing_touchdowns,
    passing_interceptions,
    rushing_attempts,
    rushing_yards,
    rushing_touchdowns,
    targets,
    receptions,
    receiving_yards,
    receiving_touchdowns,
    field_goals_made,
    field_goals_attempted,
    extra_points_made,
    extra_points_attempted,
    sacks,
    interceptions,
    fumble_recoveries,
    defensive_touchdowns,
    safeties,
    fumbles_lost,
    two_point_conversions
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
    ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
    ?, ?, ?, ?, ?, ?, ?
)
RETURNING id;

-- name: UpdatePlayerStats :exec
UPDATE player_stats
SET
    passing_attempts = ?,
    passing_completions = ?,
    passing_yards = ?,
    passing_touchdowns = ?,
    passing_interceptions = ?,
    rushing_attempts = ?,
    rushing_yards = ?,
    rushing_touchdowns = ?,
    targets = ?,
    receptions = ?,
    receiving_yards = ?,
    receiving_touchdowns = ?,
    field_goals_made = ?,
    field_goals_attempted = ?,
    extra_points_made = ?,
    extra_points_attempted = ?,
    sacks = ?,
    interceptions = ?,
    fumble_recoveries = ?,
    defensive_touchdowns = ?,
    safeties = ?,
    fumbles_lost = ?,
    two_point_conversions = ?
WHERE id = ?;

-- name: CalculatePlayerFantasyPoints :one
WITH player_stats_data AS (
    SELECT
        ps.*,
        fsr.passing_yards_per_point,
        fsr.passing_touchdown_points,
        fsr.passing_interception_points,
        fsr.rushing_yards_per_point,
        fsr.rushing_touchdown_points,
        fsr.receiving_yards_per_point,
        fsr.receiving_touchdown_points,
        fsr.reception_points,
        fsr.field_goal_0_39_points,
        fsr.field_goal_40_49_points,
        fsr.field_goal_50_plus_points,
        fsr.extra_point_points,
        fsr.sack_points,
        fsr.interception_points,
        fsr.fumble_recovery_points,
        fsr.defensive_touchdown_points,
        fsr.safety_points,
        fsr.two_point_conversion_points,
        fsr.fumble_lost_points
    FROM player_stats ps
    CROSS JOIN fantasy_scoring_rules fsr
    WHERE ps.player_id = ? AND ps.season_id = ? AND ps.week = ? AND fsr.league_id = ?
)
SELECT
    -- Passing points
    (passing_yards / passing_yards_per_point) +
    (passing_touchdowns * passing_touchdown_points) +
    (passing_interceptions * passing_interception_points) +
    
    -- Rushing points
    (rushing_yards / rushing_yards_per_point) +
    (rushing_touchdowns * rushing_touchdown_points) +
    
    -- Receiving points
    (receiving_yards / receiving_yards_per_point) +
    (receiving_touchdowns * receiving_touchdown_points) +
    (receptions * reception_points) +
    
    -- Kicking points (simplified - would need more data for FG distance)
    (field_goals_made * field_goal_0_39_points) +
    (extra_points_made * extra_point_points) +
    
    -- Defense points
    (sacks * sack_points) +
    (interceptions * interception_points) +
    (fumble_recoveries * fumble_recovery_points) +
    (defensive_touchdowns * defensive_touchdown_points) +
    (safeties * safety_points) +
    
    -- Misc points
    (two_point_conversions * two_point_conversion_points) +
    (fumbles_lost * fumble_lost_points) AS fantasy_points
FROM player_stats_data;

-- name: GetTopScorersForWeek :many
SELECT 
    fps.player_id,
    p.name as player_name,
    p.position as player_position,
    t.abbreviation as team_abbreviation,
    fps.points as fantasy_points,
    fps.team_id as fantasy_team_id,
    ft.name as fantasy_team_name
FROM fantasy_player_scores fps
JOIN nfl_players p ON fps.player_id = p.id
LEFT JOIN nfl_teams t ON p.team_id = t.id
LEFT JOIN fantasy_teams ft ON fps.team_id = ft.id
WHERE fps.league_id = ? AND fps.week = ? AND fps.season_id = ?
ORDER BY fps.points DESC
LIMIT ?;

-- name: GetTopScorersForSeason :many
SELECT 
    p.id as player_id,
    p.name as player_name,
    p.position as player_position,
    t.abbreviation as team_abbreviation,
    SUM(fps.points) as total_fantasy_points,
    AVG(fps.points) as avg_fantasy_points,
    COUNT(fps.week) as games_played
FROM fantasy_player_scores fps
JOIN nfl_players p ON fps.player_id = p.id
LEFT JOIN nfl_teams t ON p.team_id = t.id
WHERE fps.league_id = ? AND fps.season_id = ?
GROUP BY p.id, p.name, p.position, t.abbreviation
ORDER BY total_fantasy_points DESC
LIMIT ?;

-- name: GetTeamScoreForWeek :one
SELECT 
    ft.id as team_id,
    ft.name as team_name,
    SUM(fps.points) as total_points
FROM fantasy_player_scores fps
JOIN fantasy_teams ft ON fps.team_id = ft.id
JOIN fantasy_rosters fr ON fps.player_id = fr.player_id AND fr.team_id = ft.id
WHERE fps.league_id = ? AND fps.week = ? AND fps.season_id = ? AND fr.is_starter = 1 AND ft.id = ?
GROUP BY ft.id, ft.name;

-- name: GetTeamsForScoringUpdate :many
SELECT 
    ft.id as team_id,
    fm.id as matchup_id,
    fm.week,
    CASE WHEN fm.home_team_id = ft.id THEN 'home' ELSE 'away' END as home_or_away
FROM fantasy_teams ft
JOIN fantasy_matchups fm ON (fm.home_team_id = ft.id OR fm.away_team_id = ft.id)
WHERE fm.league_id = ? AND fm.week = ? AND fm.completed = 0;

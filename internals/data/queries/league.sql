-- name: GetAllLeagues :many
SELECT 
    fl.id, 
    fl.name, 
    fl.season_id, 
    fl.teams_count, 
    fl.qb_count, 
    fl.rb_count, 
    fl.wr_count, 
    fl.te_count, 
    fl.flex_count, 
    fl.k_count, 
    fl.dst_count, 
    fl.bench_count, 
    fl.ppr,
    fl.created_at,
    s.year as season_year
FROM fantasy_leagues fl
JOIN seasons s ON fl.season_id = s.id
ORDER BY fl.created_at DESC;

-- name: GetLeagueById :one
SELECT 
    fl.id, 
    fl.name, 
    fl.season_id, 
    fl.teams_count, 
    fl.qb_count, 
    fl.rb_count, 
    fl.wr_count, 
    fl.te_count, 
    fl.flex_count, 
    fl.k_count, 
    fl.dst_count, 
    fl.bench_count, 
    fl.ppr,
    fl.created_at,
    s.year as season_year
FROM fantasy_leagues fl
JOIN seasons s ON fl.season_id = s.id
WHERE fl.id = ?;

-- name: GetLeagueStandings :many
SELECT 
    ft.id,
    ft.name,
    ft.owner_name,
    ft.is_user,
    ft.wins,
    ft.losses,
    ft.tie_games,
    ft.points_for,
    ft.points_against,
    (ft.wins * 2 + ft.tie_games) as total_points,
    (ft.points_for - ft.points_against) as point_differential
FROM fantasy_teams ft
WHERE ft.league_id = ?
ORDER BY total_points DESC, point_differential DESC, ft.points_for DESC;

-- name: GetLeagueMatchups :many
SELECT 
    fm.id,
    fm.week,
    fm.home_team_id,
    ht.name as home_team_name,
    fm.away_team_id,
    at.name as away_team_name,
    fm.home_score,
    fm.away_score,
    fm.is_playoff,
    fm.is_championship,
    fm.completed
FROM fantasy_matchups fm
JOIN fantasy_teams ht ON fm.home_team_id = ht.id
JOIN fantasy_teams at ON fm.away_team_id = at.id
WHERE fm.league_id = ?
ORDER BY fm.week, fm.id;

-- name: GetLeagueMatchupsByWeek :many
SELECT 
    fm.id,
    fm.week,
    fm.home_team_id,
    ht.name as home_team_name,
    ht.owner_name as home_owner_name,
    fm.away_team_id,
    at.name as away_team_name,
    at.owner_name as away_owner_name,
    fm.home_score,
    fm.away_score,
    fm.is_playoff,
    fm.is_championship,
    fm.completed
FROM fantasy_matchups fm
JOIN fantasy_teams ht ON fm.home_team_id = ht.id
JOIN fantasy_teams at ON fm.away_team_id = at.id
WHERE fm.league_id = ? AND fm.week = ?
ORDER BY fm.id;

-- name: GetLeagueScoringRules :one
SELECT 
    id,
    league_id,
    passing_yards_per_point,
    passing_touchdown_points,
    passing_interception_points,
    rushing_yards_per_point,
    rushing_touchdown_points,
    receiving_yards_per_point,
    receiving_touchdown_points,
    reception_points,
    field_goal_0_39_points,
    field_goal_40_49_points,
    field_goal_50_plus_points,
    extra_point_points,
    sack_points,
    interception_points,
    fumble_recovery_points,
    defensive_touchdown_points,
    safety_points,
    two_point_conversion_points,
    fumble_lost_points
FROM fantasy_scoring_rules
WHERE league_id = ?;

-- name: CreateLeague :one
INSERT INTO fantasy_leagues (
    name,
    season_id,
    teams_count,
    qb_count,
    rb_count,
    wr_count,
    te_count,
    flex_count,
    k_count,
    dst_count,
    bench_count,
    ppr
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING id;

-- name: CreateScoringRules :one
INSERT INTO fantasy_scoring_rules (
    league_id,
    passing_yards_per_point,
    passing_touchdown_points,
    passing_interception_points,
    rushing_yards_per_point,
    rushing_touchdown_points,
    receiving_yards_per_point,
    receiving_touchdown_points,
    reception_points,
    field_goal_0_39_points,
    field_goal_40_49_points,
    field_goal_50_plus_points,
    extra_point_points,
    sack_points,
    interception_points,
    fumble_recovery_points,
    defensive_touchdown_points,
    safety_points,
    two_point_conversion_points,
    fumble_lost_points
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
    ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
)
RETURNING id;

-- name: UpdateScoringRules :exec
UPDATE fantasy_scoring_rules
SET
    passing_yards_per_point = ?,
    passing_touchdown_points = ?,
    passing_interception_points = ?,
    rushing_yards_per_point = ?,
    rushing_touchdown_points = ?,
    receiving_yards_per_point = ?,
    receiving_touchdown_points = ?,
    reception_points = ?,
    field_goal_0_39_points = ?,
    field_goal_40_49_points = ?,
    field_goal_50_plus_points = ?,
    extra_point_points = ?,
    sack_points = ?,
    interception_points = ?,
    fumble_recovery_points = ?,
    defensive_touchdown_points = ?,
    safety_points = ?,
    two_point_conversion_points = ?,
    fumble_lost_points = ?
WHERE league_id = ?;

-- name: CreateMatchup :one
INSERT INTO fantasy_matchups (
    league_id,
    week,
    home_team_id,
    away_team_id,
    is_playoff,
    is_championship,
    completed
) VALUES (?, ?, ?, ?, ?, ?, 0)
RETURNING id;

-- name: UpdateMatchupScore :exec
UPDATE fantasy_matchups
SET
    home_score = ?,
    away_score = ?,
    completed = ?
WHERE id = ?;

-- name: GetPlayoffTeams :many
SELECT 
    ft.id,
    ft.name,
    ft.owner_name,
    ft.is_user,
    ft.wins,
    ft.losses,
    ft.tie_games,
    ft.points_for,
    ft.points_against,
    (ft.wins * 2 + ft.tie_games) as total_points
FROM fantasy_teams ft
WHERE ft.league_id = ?
ORDER BY total_points DESC, ft.points_for DESC
LIMIT 4;

-- name: GetLeagueWeeklyScores :many
SELECT 
    fps.week,
    fps.player_id,
    np.name as player_name,
    np.position as player_position,
    ft.id as team_id,
    ft.name as team_name,
    fps.points as fantasy_points
FROM fantasy_player_scores fps
JOIN nfl_players np ON fps.player_id = np.id
JOIN fantasy_teams ft ON fps.team_id = ft.id
WHERE fps.league_id = ? AND fps.week = ?
ORDER BY ft.id, fps.points DESC;

-- name: CalculateFantasyScore :exec
INSERT INTO fantasy_player_scores (
    league_id,
    team_id,
    player_id,
    week,
    season_id,
    points
) VALUES (?, ?, ?, ?, ?, ?)
ON CONFLICT (league_id, player_id, week, season_id)
DO UPDATE SET
    team_id = excluded.team_id,
    points = excluded.points;

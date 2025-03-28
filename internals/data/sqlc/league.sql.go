// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: league.sql

package data

import (
	"context"
	"database/sql"
)

const calculateFantasyScore = `-- name: CalculateFantasyScore :exec
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
    points = excluded.points
`

type CalculateFantasyScoreParams struct {
	LeagueID int64           `json:"league_id"`
	TeamID   sql.NullInt64   `json:"team_id"`
	PlayerID int64           `json:"player_id"`
	Week     int64           `json:"week"`
	SeasonID int64           `json:"season_id"`
	Points   sql.NullFloat64 `json:"points"`
}

func (q *Queries) CalculateFantasyScore(ctx context.Context, arg CalculateFantasyScoreParams) error {
	_, err := q.exec(ctx, q.calculateFantasyScoreStmt, calculateFantasyScore,
		arg.LeagueID,
		arg.TeamID,
		arg.PlayerID,
		arg.Week,
		arg.SeasonID,
		arg.Points,
	)
	return err
}

const createLeague = `-- name: CreateLeague :one
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
RETURNING id
`

type CreateLeagueParams struct {
	Name       string        `json:"name"`
	SeasonID   int64         `json:"season_id"`
	TeamsCount int64         `json:"teams_count"`
	QbCount    sql.NullInt64 `json:"qb_count"`
	RbCount    sql.NullInt64 `json:"rb_count"`
	WrCount    sql.NullInt64 `json:"wr_count"`
	TeCount    sql.NullInt64 `json:"te_count"`
	FlexCount  sql.NullInt64 `json:"flex_count"`
	KCount     sql.NullInt64 `json:"k_count"`
	DstCount   sql.NullInt64 `json:"dst_count"`
	BenchCount sql.NullInt64 `json:"bench_count"`
	Ppr        sql.NullInt64 `json:"ppr"`
}

func (q *Queries) CreateLeague(ctx context.Context, arg CreateLeagueParams) (int64, error) {
	row := q.queryRow(ctx, q.createLeagueStmt, createLeague,
		arg.Name,
		arg.SeasonID,
		arg.TeamsCount,
		arg.QbCount,
		arg.RbCount,
		arg.WrCount,
		arg.TeCount,
		arg.FlexCount,
		arg.KCount,
		arg.DstCount,
		arg.BenchCount,
		arg.Ppr,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const createMatchup = `-- name: CreateMatchup :one
INSERT INTO fantasy_matchups (
    league_id,
    week,
    home_team_id,
    away_team_id,
    is_playoff,
    is_championship,
    completed
) VALUES (?, ?, ?, ?, ?, ?, 0)
RETURNING id
`

type CreateMatchupParams struct {
	LeagueID       int64         `json:"league_id"`
	Week           int64         `json:"week"`
	HomeTeamID     int64         `json:"home_team_id"`
	AwayTeamID     int64         `json:"away_team_id"`
	IsPlayoff      sql.NullInt64 `json:"is_playoff"`
	IsChampionship sql.NullInt64 `json:"is_championship"`
}

func (q *Queries) CreateMatchup(ctx context.Context, arg CreateMatchupParams) (int64, error) {
	row := q.queryRow(ctx, q.createMatchupStmt, createMatchup,
		arg.LeagueID,
		arg.Week,
		arg.HomeTeamID,
		arg.AwayTeamID,
		arg.IsPlayoff,
		arg.IsChampionship,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const createScoringRules = `-- name: CreateScoringRules :one
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
RETURNING id
`

type CreateScoringRulesParams struct {
	LeagueID                  int64           `json:"league_id"`
	PassingYardsPerPoint      sql.NullFloat64 `json:"passing_yards_per_point"`
	PassingTouchdownPoints    sql.NullFloat64 `json:"passing_touchdown_points"`
	PassingInterceptionPoints sql.NullFloat64 `json:"passing_interception_points"`
	RushingYardsPerPoint      sql.NullFloat64 `json:"rushing_yards_per_point"`
	RushingTouchdownPoints    sql.NullFloat64 `json:"rushing_touchdown_points"`
	ReceivingYardsPerPoint    sql.NullFloat64 `json:"receiving_yards_per_point"`
	ReceivingTouchdownPoints  sql.NullFloat64 `json:"receiving_touchdown_points"`
	ReceptionPoints           sql.NullFloat64 `json:"reception_points"`
	FieldGoal039Points        sql.NullFloat64 `json:"field_goal_0_39_points"`
	FieldGoal4049Points       sql.NullFloat64 `json:"field_goal_40_49_points"`
	FieldGoal50PlusPoints     sql.NullFloat64 `json:"field_goal_50_plus_points"`
	ExtraPointPoints          sql.NullFloat64 `json:"extra_point_points"`
	SackPoints                sql.NullFloat64 `json:"sack_points"`
	InterceptionPoints        sql.NullFloat64 `json:"interception_points"`
	FumbleRecoveryPoints      sql.NullFloat64 `json:"fumble_recovery_points"`
	DefensiveTouchdownPoints  sql.NullFloat64 `json:"defensive_touchdown_points"`
	SafetyPoints              sql.NullFloat64 `json:"safety_points"`
	TwoPointConversionPoints  sql.NullFloat64 `json:"two_point_conversion_points"`
	FumbleLostPoints          sql.NullFloat64 `json:"fumble_lost_points"`
}

func (q *Queries) CreateScoringRules(ctx context.Context, arg CreateScoringRulesParams) (int64, error) {
	row := q.queryRow(ctx, q.createScoringRulesStmt, createScoringRules,
		arg.LeagueID,
		arg.PassingYardsPerPoint,
		arg.PassingTouchdownPoints,
		arg.PassingInterceptionPoints,
		arg.RushingYardsPerPoint,
		arg.RushingTouchdownPoints,
		arg.ReceivingYardsPerPoint,
		arg.ReceivingTouchdownPoints,
		arg.ReceptionPoints,
		arg.FieldGoal039Points,
		arg.FieldGoal4049Points,
		arg.FieldGoal50PlusPoints,
		arg.ExtraPointPoints,
		arg.SackPoints,
		arg.InterceptionPoints,
		arg.FumbleRecoveryPoints,
		arg.DefensiveTouchdownPoints,
		arg.SafetyPoints,
		arg.TwoPointConversionPoints,
		arg.FumbleLostPoints,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const getAllLeagues = `-- name: GetAllLeagues :many
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
ORDER BY fl.created_at DESC
`

type GetAllLeaguesRow struct {
	ID         int64          `json:"id"`
	Name       string         `json:"name"`
	SeasonID   int64          `json:"season_id"`
	TeamsCount int64          `json:"teams_count"`
	QbCount    sql.NullInt64  `json:"qb_count"`
	RbCount    sql.NullInt64  `json:"rb_count"`
	WrCount    sql.NullInt64  `json:"wr_count"`
	TeCount    sql.NullInt64  `json:"te_count"`
	FlexCount  sql.NullInt64  `json:"flex_count"`
	KCount     sql.NullInt64  `json:"k_count"`
	DstCount   sql.NullInt64  `json:"dst_count"`
	BenchCount sql.NullInt64  `json:"bench_count"`
	Ppr        sql.NullInt64  `json:"ppr"`
	CreatedAt  sql.NullString `json:"created_at"`
	SeasonYear int64          `json:"season_year"`
}

func (q *Queries) GetAllLeagues(ctx context.Context) ([]*GetAllLeaguesRow, error) {
	rows, err := q.query(ctx, q.getAllLeaguesStmt, getAllLeagues)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetAllLeaguesRow{}
	for rows.Next() {
		var i GetAllLeaguesRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.SeasonID,
			&i.TeamsCount,
			&i.QbCount,
			&i.RbCount,
			&i.WrCount,
			&i.TeCount,
			&i.FlexCount,
			&i.KCount,
			&i.DstCount,
			&i.BenchCount,
			&i.Ppr,
			&i.CreatedAt,
			&i.SeasonYear,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getLeagueById = `-- name: GetLeagueById :one
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
WHERE fl.id = ?
`

type GetLeagueByIdRow struct {
	ID         int64          `json:"id"`
	Name       string         `json:"name"`
	SeasonID   int64          `json:"season_id"`
	TeamsCount int64          `json:"teams_count"`
	QbCount    sql.NullInt64  `json:"qb_count"`
	RbCount    sql.NullInt64  `json:"rb_count"`
	WrCount    sql.NullInt64  `json:"wr_count"`
	TeCount    sql.NullInt64  `json:"te_count"`
	FlexCount  sql.NullInt64  `json:"flex_count"`
	KCount     sql.NullInt64  `json:"k_count"`
	DstCount   sql.NullInt64  `json:"dst_count"`
	BenchCount sql.NullInt64  `json:"bench_count"`
	Ppr        sql.NullInt64  `json:"ppr"`
	CreatedAt  sql.NullString `json:"created_at"`
	SeasonYear int64          `json:"season_year"`
}

func (q *Queries) GetLeagueById(ctx context.Context, id int64) (*GetLeagueByIdRow, error) {
	row := q.queryRow(ctx, q.getLeagueByIdStmt, getLeagueById, id)
	var i GetLeagueByIdRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.SeasonID,
		&i.TeamsCount,
		&i.QbCount,
		&i.RbCount,
		&i.WrCount,
		&i.TeCount,
		&i.FlexCount,
		&i.KCount,
		&i.DstCount,
		&i.BenchCount,
		&i.Ppr,
		&i.CreatedAt,
		&i.SeasonYear,
	)
	return &i, err
}

const getLeagueMatchups = `-- name: GetLeagueMatchups :many
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
ORDER BY fm.week, fm.id
`

type GetLeagueMatchupsRow struct {
	ID             int64           `json:"id"`
	Week           int64           `json:"week"`
	HomeTeamID     int64           `json:"home_team_id"`
	HomeTeamName   string          `json:"home_team_name"`
	AwayTeamID     int64           `json:"away_team_id"`
	AwayTeamName   string          `json:"away_team_name"`
	HomeScore      sql.NullFloat64 `json:"home_score"`
	AwayScore      sql.NullFloat64 `json:"away_score"`
	IsPlayoff      sql.NullInt64   `json:"is_playoff"`
	IsChampionship sql.NullInt64   `json:"is_championship"`
	Completed      sql.NullInt64   `json:"completed"`
}

func (q *Queries) GetLeagueMatchups(ctx context.Context, leagueID int64) ([]*GetLeagueMatchupsRow, error) {
	rows, err := q.query(ctx, q.getLeagueMatchupsStmt, getLeagueMatchups, leagueID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetLeagueMatchupsRow{}
	for rows.Next() {
		var i GetLeagueMatchupsRow
		if err := rows.Scan(
			&i.ID,
			&i.Week,
			&i.HomeTeamID,
			&i.HomeTeamName,
			&i.AwayTeamID,
			&i.AwayTeamName,
			&i.HomeScore,
			&i.AwayScore,
			&i.IsPlayoff,
			&i.IsChampionship,
			&i.Completed,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getLeagueMatchupsByWeek = `-- name: GetLeagueMatchupsByWeek :many
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
ORDER BY fm.id
`

type GetLeagueMatchupsByWeekParams struct {
	LeagueID int64 `json:"league_id"`
	Week     int64 `json:"week"`
}

type GetLeagueMatchupsByWeekRow struct {
	ID             int64           `json:"id"`
	Week           int64           `json:"week"`
	HomeTeamID     int64           `json:"home_team_id"`
	HomeTeamName   string          `json:"home_team_name"`
	HomeOwnerName  string          `json:"home_owner_name"`
	AwayTeamID     int64           `json:"away_team_id"`
	AwayTeamName   string          `json:"away_team_name"`
	AwayOwnerName  string          `json:"away_owner_name"`
	HomeScore      sql.NullFloat64 `json:"home_score"`
	AwayScore      sql.NullFloat64 `json:"away_score"`
	IsPlayoff      sql.NullInt64   `json:"is_playoff"`
	IsChampionship sql.NullInt64   `json:"is_championship"`
	Completed      sql.NullInt64   `json:"completed"`
}

func (q *Queries) GetLeagueMatchupsByWeek(ctx context.Context, arg GetLeagueMatchupsByWeekParams) ([]*GetLeagueMatchupsByWeekRow, error) {
	rows, err := q.query(ctx, q.getLeagueMatchupsByWeekStmt, getLeagueMatchupsByWeek, arg.LeagueID, arg.Week)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetLeagueMatchupsByWeekRow{}
	for rows.Next() {
		var i GetLeagueMatchupsByWeekRow
		if err := rows.Scan(
			&i.ID,
			&i.Week,
			&i.HomeTeamID,
			&i.HomeTeamName,
			&i.HomeOwnerName,
			&i.AwayTeamID,
			&i.AwayTeamName,
			&i.AwayOwnerName,
			&i.HomeScore,
			&i.AwayScore,
			&i.IsPlayoff,
			&i.IsChampionship,
			&i.Completed,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getLeagueScoringRules = `-- name: GetLeagueScoringRules :one
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
WHERE league_id = ?
`

func (q *Queries) GetLeagueScoringRules(ctx context.Context, leagueID int64) (*FantasyScoringRule, error) {
	row := q.queryRow(ctx, q.getLeagueScoringRulesStmt, getLeagueScoringRules, leagueID)
	var i FantasyScoringRule
	err := row.Scan(
		&i.ID,
		&i.LeagueID,
		&i.PassingYardsPerPoint,
		&i.PassingTouchdownPoints,
		&i.PassingInterceptionPoints,
		&i.RushingYardsPerPoint,
		&i.RushingTouchdownPoints,
		&i.ReceivingYardsPerPoint,
		&i.ReceivingTouchdownPoints,
		&i.ReceptionPoints,
		&i.FieldGoal039Points,
		&i.FieldGoal4049Points,
		&i.FieldGoal50PlusPoints,
		&i.ExtraPointPoints,
		&i.SackPoints,
		&i.InterceptionPoints,
		&i.FumbleRecoveryPoints,
		&i.DefensiveTouchdownPoints,
		&i.SafetyPoints,
		&i.TwoPointConversionPoints,
		&i.FumbleLostPoints,
	)
	return &i, err
}

const getLeagueStandings = `-- name: GetLeagueStandings :many
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
ORDER BY total_points DESC, point_differential DESC, ft.points_for DESC
`

type GetLeagueStandingsRow struct {
	ID                int64           `json:"id"`
	Name              string          `json:"name"`
	OwnerName         string          `json:"owner_name"`
	IsUser            sql.NullInt64   `json:"is_user"`
	Wins              sql.NullInt64   `json:"wins"`
	Losses            sql.NullInt64   `json:"losses"`
	TieGames          sql.NullInt64   `json:"tie_games"`
	PointsFor         sql.NullFloat64 `json:"points_for"`
	PointsAgainst     sql.NullFloat64 `json:"points_against"`
	TotalPoints       interface{}     `json:"total_points"`
	PointDifferential interface{}     `json:"point_differential"`
}

func (q *Queries) GetLeagueStandings(ctx context.Context, leagueID int64) ([]*GetLeagueStandingsRow, error) {
	rows, err := q.query(ctx, q.getLeagueStandingsStmt, getLeagueStandings, leagueID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetLeagueStandingsRow{}
	for rows.Next() {
		var i GetLeagueStandingsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.OwnerName,
			&i.IsUser,
			&i.Wins,
			&i.Losses,
			&i.TieGames,
			&i.PointsFor,
			&i.PointsAgainst,
			&i.TotalPoints,
			&i.PointDifferential,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getLeagueWeeklyScores = `-- name: GetLeagueWeeklyScores :many
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
ORDER BY ft.id, fps.points DESC
`

type GetLeagueWeeklyScoresParams struct {
	LeagueID int64 `json:"league_id"`
	Week     int64 `json:"week"`
}

type GetLeagueWeeklyScoresRow struct {
	Week           int64           `json:"week"`
	PlayerID       int64           `json:"player_id"`
	PlayerName     string          `json:"player_name"`
	PlayerPosition string          `json:"player_position"`
	TeamID         int64           `json:"team_id"`
	TeamName       string          `json:"team_name"`
	FantasyPoints  sql.NullFloat64 `json:"fantasy_points"`
}

func (q *Queries) GetLeagueWeeklyScores(ctx context.Context, arg GetLeagueWeeklyScoresParams) ([]*GetLeagueWeeklyScoresRow, error) {
	rows, err := q.query(ctx, q.getLeagueWeeklyScoresStmt, getLeagueWeeklyScores, arg.LeagueID, arg.Week)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetLeagueWeeklyScoresRow{}
	for rows.Next() {
		var i GetLeagueWeeklyScoresRow
		if err := rows.Scan(
			&i.Week,
			&i.PlayerID,
			&i.PlayerName,
			&i.PlayerPosition,
			&i.TeamID,
			&i.TeamName,
			&i.FantasyPoints,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPlayoffTeams = `-- name: GetPlayoffTeams :many
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
LIMIT 4
`

type GetPlayoffTeamsRow struct {
	ID            int64           `json:"id"`
	Name          string          `json:"name"`
	OwnerName     string          `json:"owner_name"`
	IsUser        sql.NullInt64   `json:"is_user"`
	Wins          sql.NullInt64   `json:"wins"`
	Losses        sql.NullInt64   `json:"losses"`
	TieGames      sql.NullInt64   `json:"tie_games"`
	PointsFor     sql.NullFloat64 `json:"points_for"`
	PointsAgainst sql.NullFloat64 `json:"points_against"`
	TotalPoints   interface{}     `json:"total_points"`
}

func (q *Queries) GetPlayoffTeams(ctx context.Context, leagueID int64) ([]*GetPlayoffTeamsRow, error) {
	rows, err := q.query(ctx, q.getPlayoffTeamsStmt, getPlayoffTeams, leagueID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetPlayoffTeamsRow{}
	for rows.Next() {
		var i GetPlayoffTeamsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.OwnerName,
			&i.IsUser,
			&i.Wins,
			&i.Losses,
			&i.TieGames,
			&i.PointsFor,
			&i.PointsAgainst,
			&i.TotalPoints,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateMatchupScore = `-- name: UpdateMatchupScore :exec
UPDATE fantasy_matchups
SET
    home_score = ?,
    away_score = ?,
    completed = ?
WHERE id = ?
`

type UpdateMatchupScoreParams struct {
	HomeScore sql.NullFloat64 `json:"home_score"`
	AwayScore sql.NullFloat64 `json:"away_score"`
	Completed sql.NullInt64   `json:"completed"`
	ID        int64           `json:"id"`
}

func (q *Queries) UpdateMatchupScore(ctx context.Context, arg UpdateMatchupScoreParams) error {
	_, err := q.exec(ctx, q.updateMatchupScoreStmt, updateMatchupScore,
		arg.HomeScore,
		arg.AwayScore,
		arg.Completed,
		arg.ID,
	)
	return err
}

const updateScoringRules = `-- name: UpdateScoringRules :exec
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
WHERE league_id = ?
`

type UpdateScoringRulesParams struct {
	PassingYardsPerPoint      sql.NullFloat64 `json:"passing_yards_per_point"`
	PassingTouchdownPoints    sql.NullFloat64 `json:"passing_touchdown_points"`
	PassingInterceptionPoints sql.NullFloat64 `json:"passing_interception_points"`
	RushingYardsPerPoint      sql.NullFloat64 `json:"rushing_yards_per_point"`
	RushingTouchdownPoints    sql.NullFloat64 `json:"rushing_touchdown_points"`
	ReceivingYardsPerPoint    sql.NullFloat64 `json:"receiving_yards_per_point"`
	ReceivingTouchdownPoints  sql.NullFloat64 `json:"receiving_touchdown_points"`
	ReceptionPoints           sql.NullFloat64 `json:"reception_points"`
	FieldGoal039Points        sql.NullFloat64 `json:"field_goal_0_39_points"`
	FieldGoal4049Points       sql.NullFloat64 `json:"field_goal_40_49_points"`
	FieldGoal50PlusPoints     sql.NullFloat64 `json:"field_goal_50_plus_points"`
	ExtraPointPoints          sql.NullFloat64 `json:"extra_point_points"`
	SackPoints                sql.NullFloat64 `json:"sack_points"`
	InterceptionPoints        sql.NullFloat64 `json:"interception_points"`
	FumbleRecoveryPoints      sql.NullFloat64 `json:"fumble_recovery_points"`
	DefensiveTouchdownPoints  sql.NullFloat64 `json:"defensive_touchdown_points"`
	SafetyPoints              sql.NullFloat64 `json:"safety_points"`
	TwoPointConversionPoints  sql.NullFloat64 `json:"two_point_conversion_points"`
	FumbleLostPoints          sql.NullFloat64 `json:"fumble_lost_points"`
	LeagueID                  int64           `json:"league_id"`
}

func (q *Queries) UpdateScoringRules(ctx context.Context, arg UpdateScoringRulesParams) error {
	_, err := q.exec(ctx, q.updateScoringRulesStmt, updateScoringRules,
		arg.PassingYardsPerPoint,
		arg.PassingTouchdownPoints,
		arg.PassingInterceptionPoints,
		arg.RushingYardsPerPoint,
		arg.RushingTouchdownPoints,
		arg.ReceivingYardsPerPoint,
		arg.ReceivingTouchdownPoints,
		arg.ReceptionPoints,
		arg.FieldGoal039Points,
		arg.FieldGoal4049Points,
		arg.FieldGoal50PlusPoints,
		arg.ExtraPointPoints,
		arg.SackPoints,
		arg.InterceptionPoints,
		arg.FumbleRecoveryPoints,
		arg.DefensiveTouchdownPoints,
		arg.SafetyPoints,
		arg.TwoPointConversionPoints,
		arg.FumbleLostPoints,
		arg.LeagueID,
	)
	return err
}

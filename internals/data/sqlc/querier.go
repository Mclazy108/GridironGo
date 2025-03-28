// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package sqlc

import (
	"context"
	"database/sql"
)

type Querier interface {
	AddDraftPick(ctx context.Context, arg AddDraftPickParams) (int64, error)
	AddPlayerToFantasyTeam(ctx context.Context, arg AddPlayerToFantasyTeamParams) (int64, error)
	CalculateFantasyScore(ctx context.Context, arg CalculateFantasyScoreParams) error
	CalculatePlayerFantasyPoints(ctx context.Context, arg CalculatePlayerFantasyPointsParams) (int64, error)
	ClearDraft(ctx context.Context, leagueID int64) error
	CreateFantasyTeam(ctx context.Context, arg CreateFantasyTeamParams) (int64, error)
	CreateLeague(ctx context.Context, arg CreateLeagueParams) (int64, error)
	CreateMatchup(ctx context.Context, arg CreateMatchupParams) (int64, error)
	CreateNFLGame(ctx context.Context, arg CreateNFLGameParams) (int64, error)
	CreatePlayerStats(ctx context.Context, arg CreatePlayerStatsParams) (int64, error)
	CreateScoringRules(ctx context.Context, arg CreateScoringRulesParams) (int64, error)
	CreateSeason(ctx context.Context, arg CreateSeasonParams) (int64, error)
	GetAllFantasyTeams(ctx context.Context, leagueID int64) ([]*GetAllFantasyTeamsRow, error)
	GetAllLeagues(ctx context.Context) ([]*GetAllLeaguesRow, error)
	GetAllNFLTeams(ctx context.Context) ([]*NflTeam, error)
	GetAllPlayers(ctx context.Context) ([]*NflPlayer, error)
	GetAllPlayersWithFantasyPoints(ctx context.Context, arg GetAllPlayersWithFantasyPointsParams) ([]*GetAllPlayersWithFantasyPointsRow, error)
	GetAllSeasons(ctx context.Context) ([]*Season, error)
	GetAvailableDraftPlayers(ctx context.Context, arg GetAvailableDraftPlayersParams) ([]*GetAvailableDraftPlayersRow, error)
	GetAvailablePlayers(ctx context.Context, arg GetAvailablePlayersParams) ([]*GetAvailablePlayersRow, error)
	GetBestAvailablePlayers(ctx context.Context, arg GetBestAvailablePlayersParams) ([]*GetBestAvailablePlayersRow, error)
	GetCurrentSeason(ctx context.Context) (*Season, error)
	GetDraftOrder(ctx context.Context, leagueID int64) ([]*GetDraftOrderRow, error)
	GetDraftPicksByTeam(ctx context.Context, arg GetDraftPicksByTeamParams) ([]*GetDraftPicksByTeamRow, error)
	GetDraftPicksForLeague(ctx context.Context, leagueID int64) ([]*GetDraftPicksForLeagueRow, error)
	GetDraftSummary(ctx context.Context, arg GetDraftSummaryParams) ([]*GetDraftSummaryRow, error)
	GetFantasyTeamById(ctx context.Context, id int64) (*GetFantasyTeamByIdRow, error)
	GetFantasyTeamForWeek(ctx context.Context, arg GetFantasyTeamForWeekParams) (*GetFantasyTeamForWeekRow, error)
	GetFantasyTeamRoster(ctx context.Context, teamID int64) ([]*GetFantasyTeamRosterRow, error)
	GetHistoricalPlayerStats(ctx context.Context, playerID int64) ([]*GetHistoricalPlayerStatsRow, error)
	GetLastDraftPick(ctx context.Context, leagueID int64) (*GetLastDraftPickRow, error)
	GetLeagueById(ctx context.Context, id int64) (*GetLeagueByIdRow, error)
	GetLeagueMatchups(ctx context.Context, leagueID int64) ([]*GetLeagueMatchupsRow, error)
	GetLeagueMatchupsByWeek(ctx context.Context, arg GetLeagueMatchupsByWeekParams) ([]*GetLeagueMatchupsByWeekRow, error)
	GetLeagueScoringRules(ctx context.Context, leagueID int64) (*FantasyScoringRule, error)
	GetLeagueStandings(ctx context.Context, leagueID int64) ([]*GetLeagueStandingsRow, error)
	GetLeagueWeeklyScores(ctx context.Context, arg GetLeagueWeeklyScoresParams) ([]*GetLeagueWeeklyScoresRow, error)
	GetNFLScheduleForSeason(ctx context.Context, seasonID int64) ([]*GetNFLScheduleForSeasonRow, error)
	GetNFLScheduleForWeek(ctx context.Context, arg GetNFLScheduleForWeekParams) ([]*GetNFLScheduleForWeekRow, error)
	GetNFLTeamByAbbreviation(ctx context.Context, abbreviation string) (*NflTeam, error)
	GetNFLTeamById(ctx context.Context, id int64) (*NflTeam, error)
	GetNextDraftingTeam(ctx context.Context, arg GetNextDraftingTeamParams) (*GetNextDraftingTeamRow, error)
	GetPlayerById(ctx context.Context, id int64) (*NflPlayer, error)
	GetPlayerFantasyPoints(ctx context.Context, arg GetPlayerFantasyPointsParams) ([]*GetPlayerFantasyPointsRow, error)
	GetPlayerFantasyTotalPoints(ctx context.Context, arg GetPlayerFantasyTotalPointsParams) (*GetPlayerFantasyTotalPointsRow, error)
	GetPlayerSeasonStats(ctx context.Context, arg GetPlayerSeasonStatsParams) (*GetPlayerSeasonStatsRow, error)
	GetPlayerStats(ctx context.Context, arg GetPlayerStatsParams) ([]*GetPlayerStatsRow, error)
	GetPlayerStatsByGame(ctx context.Context, arg GetPlayerStatsByGameParams) (*PlayerStat, error)
	GetPlayerStatsByWeek(ctx context.Context, arg GetPlayerStatsByWeekParams) ([]*GetPlayerStatsByWeekRow, error)
	GetPlayersByPosition(ctx context.Context, position string) ([]*NflPlayer, error)
	GetPlayersByTeam(ctx context.Context, id int64) ([]*NflPlayer, error)
	GetPlayoffTeams(ctx context.Context, leagueID int64) ([]*GetPlayoffTeamsRow, error)
	GetPreviousSeasons(ctx context.Context, limit int64) ([]*Season, error)
	GetSeasonById(ctx context.Context, id int64) (*Season, error)
	GetSeasonByYear(ctx context.Context, year int64) (*Season, error)
	GetTeamAtDraftPosition(ctx context.Context, arg GetTeamAtDraftPositionParams) (*GetTeamAtDraftPositionRow, error)
	GetTeamRoster(ctx context.Context, teamID sql.NullInt64) ([]*GetTeamRosterRow, error)
	GetTeamScoreForWeek(ctx context.Context, arg GetTeamScoreForWeekParams) (*GetTeamScoreForWeekRow, error)
	GetTeamsByConference(ctx context.Context, conference string) ([]*NflTeam, error)
	GetTeamsByDivision(ctx context.Context, division string) ([]*NflTeam, error)
	GetTeamsForScoringUpdate(ctx context.Context, arg GetTeamsForScoringUpdateParams) ([]*GetTeamsForScoringUpdateRow, error)
	GetTopPlayersByPositionAndSeason(ctx context.Context, arg GetTopPlayersByPositionAndSeasonParams) ([]*GetTopPlayersByPositionAndSeasonRow, error)
	GetTopScorersForSeason(ctx context.Context, arg GetTopScorersForSeasonParams) ([]*GetTopScorersForSeasonRow, error)
	GetTopScorersForWeek(ctx context.Context, arg GetTopScorersForWeekParams) ([]*GetTopScorersForWeekRow, error)
	GetTotalWeeksInSeason(ctx context.Context, seasonID int64) (interface{}, error)
	InsertPlayer(ctx context.Context, arg InsertPlayerParams) (int64, error)
	InsertTeam(ctx context.Context, arg InsertTeamParams) (int64, error)
	RemovePlayerFromFantasyTeam(ctx context.Context, arg RemovePlayerFromFantasyTeamParams) error
	SearchPlayers(ctx context.Context, dollar_1 sql.NullString) ([]*NflPlayer, error)
	SetCurrentSeason(ctx context.Context) error
	UpdateFantasyRoster(ctx context.Context, arg UpdateFantasyRosterParams) error
	UpdateFantasyTeamRecord(ctx context.Context, arg UpdateFantasyTeamRecordParams) error
	UpdateMatchupScore(ctx context.Context, arg UpdateMatchupScoreParams) error
	UpdateNFLGameScore(ctx context.Context, arg UpdateNFLGameScoreParams) error
	UpdatePlayer(ctx context.Context, arg UpdatePlayerParams) error
	UpdatePlayerStats(ctx context.Context, arg UpdatePlayerStatsParams) error
	UpdateScoringRules(ctx context.Context, arg UpdateScoringRulesParams) error
	UpdateSeason(ctx context.Context, arg UpdateSeasonParams) error
}

var _ Querier = (*Queries)(nil)

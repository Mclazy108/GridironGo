// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package data

import (
	"context"
	"database/sql"
	"fmt"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

func Prepare(ctx context.Context, db DBTX) (*Queries, error) {
	q := Queries{db: db}
	var err error
	if q.addDraftPickStmt, err = db.PrepareContext(ctx, addDraftPick); err != nil {
		return nil, fmt.Errorf("error preparing query AddDraftPick: %w", err)
	}
	if q.addPlayerToFantasyTeamStmt, err = db.PrepareContext(ctx, addPlayerToFantasyTeam); err != nil {
		return nil, fmt.Errorf("error preparing query AddPlayerToFantasyTeam: %w", err)
	}
	if q.calculateFantasyScoreStmt, err = db.PrepareContext(ctx, calculateFantasyScore); err != nil {
		return nil, fmt.Errorf("error preparing query CalculateFantasyScore: %w", err)
	}
	if q.calculatePlayerFantasyPointsStmt, err = db.PrepareContext(ctx, calculatePlayerFantasyPoints); err != nil {
		return nil, fmt.Errorf("error preparing query CalculatePlayerFantasyPoints: %w", err)
	}
	if q.clearDraftStmt, err = db.PrepareContext(ctx, clearDraft); err != nil {
		return nil, fmt.Errorf("error preparing query ClearDraft: %w", err)
	}
	if q.createFantasyTeamStmt, err = db.PrepareContext(ctx, createFantasyTeam); err != nil {
		return nil, fmt.Errorf("error preparing query CreateFantasyTeam: %w", err)
	}
	if q.createLeagueStmt, err = db.PrepareContext(ctx, createLeague); err != nil {
		return nil, fmt.Errorf("error preparing query CreateLeague: %w", err)
	}
	if q.createMatchupStmt, err = db.PrepareContext(ctx, createMatchup); err != nil {
		return nil, fmt.Errorf("error preparing query CreateMatchup: %w", err)
	}
	if q.createNFLGameStmt, err = db.PrepareContext(ctx, createNFLGame); err != nil {
		return nil, fmt.Errorf("error preparing query CreateNFLGame: %w", err)
	}
	if q.createPlayerStatsStmt, err = db.PrepareContext(ctx, createPlayerStats); err != nil {
		return nil, fmt.Errorf("error preparing query CreatePlayerStats: %w", err)
	}
	if q.createScoringRulesStmt, err = db.PrepareContext(ctx, createScoringRules); err != nil {
		return nil, fmt.Errorf("error preparing query CreateScoringRules: %w", err)
	}
	if q.createSeasonStmt, err = db.PrepareContext(ctx, createSeason); err != nil {
		return nil, fmt.Errorf("error preparing query CreateSeason: %w", err)
	}
	if q.getAllFantasyTeamsStmt, err = db.PrepareContext(ctx, getAllFantasyTeams); err != nil {
		return nil, fmt.Errorf("error preparing query GetAllFantasyTeams: %w", err)
	}
	if q.getAllLeaguesStmt, err = db.PrepareContext(ctx, getAllLeagues); err != nil {
		return nil, fmt.Errorf("error preparing query GetAllLeagues: %w", err)
	}
	if q.getAllNFLTeamsStmt, err = db.PrepareContext(ctx, getAllNFLTeams); err != nil {
		return nil, fmt.Errorf("error preparing query GetAllNFLTeams: %w", err)
	}
	if q.getAllPlayersStmt, err = db.PrepareContext(ctx, getAllPlayers); err != nil {
		return nil, fmt.Errorf("error preparing query GetAllPlayers: %w", err)
	}
	if q.getAllPlayersWithFantasyPointsStmt, err = db.PrepareContext(ctx, getAllPlayersWithFantasyPoints); err != nil {
		return nil, fmt.Errorf("error preparing query GetAllPlayersWithFantasyPoints: %w", err)
	}
	if q.getAllSeasonsStmt, err = db.PrepareContext(ctx, getAllSeasons); err != nil {
		return nil, fmt.Errorf("error preparing query GetAllSeasons: %w", err)
	}
	if q.getAvailableDraftPlayersStmt, err = db.PrepareContext(ctx, getAvailableDraftPlayers); err != nil {
		return nil, fmt.Errorf("error preparing query GetAvailableDraftPlayers: %w", err)
	}
	if q.getAvailablePlayersStmt, err = db.PrepareContext(ctx, getAvailablePlayers); err != nil {
		return nil, fmt.Errorf("error preparing query GetAvailablePlayers: %w", err)
	}
	if q.getBestAvailablePlayersStmt, err = db.PrepareContext(ctx, getBestAvailablePlayers); err != nil {
		return nil, fmt.Errorf("error preparing query GetBestAvailablePlayers: %w", err)
	}
	if q.getCurrentSeasonStmt, err = db.PrepareContext(ctx, getCurrentSeason); err != nil {
		return nil, fmt.Errorf("error preparing query GetCurrentSeason: %w", err)
	}
	if q.getDraftOrderStmt, err = db.PrepareContext(ctx, getDraftOrder); err != nil {
		return nil, fmt.Errorf("error preparing query GetDraftOrder: %w", err)
	}
	if q.getDraftPicksByTeamStmt, err = db.PrepareContext(ctx, getDraftPicksByTeam); err != nil {
		return nil, fmt.Errorf("error preparing query GetDraftPicksByTeam: %w", err)
	}
	if q.getDraftPicksForLeagueStmt, err = db.PrepareContext(ctx, getDraftPicksForLeague); err != nil {
		return nil, fmt.Errorf("error preparing query GetDraftPicksForLeague: %w", err)
	}
	if q.getDraftSummaryStmt, err = db.PrepareContext(ctx, getDraftSummary); err != nil {
		return nil, fmt.Errorf("error preparing query GetDraftSummary: %w", err)
	}
	if q.getFantasyTeamByIdStmt, err = db.PrepareContext(ctx, getFantasyTeamById); err != nil {
		return nil, fmt.Errorf("error preparing query GetFantasyTeamById: %w", err)
	}
	if q.getFantasyTeamForWeekStmt, err = db.PrepareContext(ctx, getFantasyTeamForWeek); err != nil {
		return nil, fmt.Errorf("error preparing query GetFantasyTeamForWeek: %w", err)
	}
	if q.getFantasyTeamRosterStmt, err = db.PrepareContext(ctx, getFantasyTeamRoster); err != nil {
		return nil, fmt.Errorf("error preparing query GetFantasyTeamRoster: %w", err)
	}
	if q.getHistoricalPlayerStatsStmt, err = db.PrepareContext(ctx, getHistoricalPlayerStats); err != nil {
		return nil, fmt.Errorf("error preparing query GetHistoricalPlayerStats: %w", err)
	}
	if q.getLastDraftPickStmt, err = db.PrepareContext(ctx, getLastDraftPick); err != nil {
		return nil, fmt.Errorf("error preparing query GetLastDraftPick: %w", err)
	}
	if q.getLeagueByIdStmt, err = db.PrepareContext(ctx, getLeagueById); err != nil {
		return nil, fmt.Errorf("error preparing query GetLeagueById: %w", err)
	}
	if q.getLeagueMatchupsStmt, err = db.PrepareContext(ctx, getLeagueMatchups); err != nil {
		return nil, fmt.Errorf("error preparing query GetLeagueMatchups: %w", err)
	}
	if q.getLeagueMatchupsByWeekStmt, err = db.PrepareContext(ctx, getLeagueMatchupsByWeek); err != nil {
		return nil, fmt.Errorf("error preparing query GetLeagueMatchupsByWeek: %w", err)
	}
	if q.getLeagueScoringRulesStmt, err = db.PrepareContext(ctx, getLeagueScoringRules); err != nil {
		return nil, fmt.Errorf("error preparing query GetLeagueScoringRules: %w", err)
	}
	if q.getLeagueStandingsStmt, err = db.PrepareContext(ctx, getLeagueStandings); err != nil {
		return nil, fmt.Errorf("error preparing query GetLeagueStandings: %w", err)
	}
	if q.getLeagueWeeklyScoresStmt, err = db.PrepareContext(ctx, getLeagueWeeklyScores); err != nil {
		return nil, fmt.Errorf("error preparing query GetLeagueWeeklyScores: %w", err)
	}
	if q.getNFLScheduleForSeasonStmt, err = db.PrepareContext(ctx, getNFLScheduleForSeason); err != nil {
		return nil, fmt.Errorf("error preparing query GetNFLScheduleForSeason: %w", err)
	}
	if q.getNFLScheduleForWeekStmt, err = db.PrepareContext(ctx, getNFLScheduleForWeek); err != nil {
		return nil, fmt.Errorf("error preparing query GetNFLScheduleForWeek: %w", err)
	}
	if q.getNFLTeamByAbbreviationStmt, err = db.PrepareContext(ctx, getNFLTeamByAbbreviation); err != nil {
		return nil, fmt.Errorf("error preparing query GetNFLTeamByAbbreviation: %w", err)
	}
	if q.getNFLTeamByIdStmt, err = db.PrepareContext(ctx, getNFLTeamById); err != nil {
		return nil, fmt.Errorf("error preparing query GetNFLTeamById: %w", err)
	}
	if q.getNextDraftingTeamStmt, err = db.PrepareContext(ctx, getNextDraftingTeam); err != nil {
		return nil, fmt.Errorf("error preparing query GetNextDraftingTeam: %w", err)
	}
	if q.getPlayerByIdStmt, err = db.PrepareContext(ctx, getPlayerById); err != nil {
		return nil, fmt.Errorf("error preparing query GetPlayerById: %w", err)
	}
	if q.getPlayerFantasyPointsStmt, err = db.PrepareContext(ctx, getPlayerFantasyPoints); err != nil {
		return nil, fmt.Errorf("error preparing query GetPlayerFantasyPoints: %w", err)
	}
	if q.getPlayerFantasyTotalPointsStmt, err = db.PrepareContext(ctx, getPlayerFantasyTotalPoints); err != nil {
		return nil, fmt.Errorf("error preparing query GetPlayerFantasyTotalPoints: %w", err)
	}
	if q.getPlayerSeasonStatsStmt, err = db.PrepareContext(ctx, getPlayerSeasonStats); err != nil {
		return nil, fmt.Errorf("error preparing query GetPlayerSeasonStats: %w", err)
	}
	if q.getPlayerStatsStmt, err = db.PrepareContext(ctx, getPlayerStats); err != nil {
		return nil, fmt.Errorf("error preparing query GetPlayerStats: %w", err)
	}
	if q.getPlayerStatsByGameStmt, err = db.PrepareContext(ctx, getPlayerStatsByGame); err != nil {
		return nil, fmt.Errorf("error preparing query GetPlayerStatsByGame: %w", err)
	}
	if q.getPlayerStatsByWeekStmt, err = db.PrepareContext(ctx, getPlayerStatsByWeek); err != nil {
		return nil, fmt.Errorf("error preparing query GetPlayerStatsByWeek: %w", err)
	}
	if q.getPlayersByPositionStmt, err = db.PrepareContext(ctx, getPlayersByPosition); err != nil {
		return nil, fmt.Errorf("error preparing query GetPlayersByPosition: %w", err)
	}
	if q.getPlayersByTeamStmt, err = db.PrepareContext(ctx, getPlayersByTeam); err != nil {
		return nil, fmt.Errorf("error preparing query GetPlayersByTeam: %w", err)
	}
	if q.getPlayoffTeamsStmt, err = db.PrepareContext(ctx, getPlayoffTeams); err != nil {
		return nil, fmt.Errorf("error preparing query GetPlayoffTeams: %w", err)
	}
	if q.getPreviousSeasonsStmt, err = db.PrepareContext(ctx, getPreviousSeasons); err != nil {
		return nil, fmt.Errorf("error preparing query GetPreviousSeasons: %w", err)
	}
	if q.getSeasonByIdStmt, err = db.PrepareContext(ctx, getSeasonById); err != nil {
		return nil, fmt.Errorf("error preparing query GetSeasonById: %w", err)
	}
	if q.getSeasonByYearStmt, err = db.PrepareContext(ctx, getSeasonByYear); err != nil {
		return nil, fmt.Errorf("error preparing query GetSeasonByYear: %w", err)
	}
	if q.getTeamAtDraftPositionStmt, err = db.PrepareContext(ctx, getTeamAtDraftPosition); err != nil {
		return nil, fmt.Errorf("error preparing query GetTeamAtDraftPosition: %w", err)
	}
	if q.getTeamRosterStmt, err = db.PrepareContext(ctx, getTeamRoster); err != nil {
		return nil, fmt.Errorf("error preparing query GetTeamRoster: %w", err)
	}
	if q.getTeamScoreForWeekStmt, err = db.PrepareContext(ctx, getTeamScoreForWeek); err != nil {
		return nil, fmt.Errorf("error preparing query GetTeamScoreForWeek: %w", err)
	}
	if q.getTeamsByConferenceStmt, err = db.PrepareContext(ctx, getTeamsByConference); err != nil {
		return nil, fmt.Errorf("error preparing query GetTeamsByConference: %w", err)
	}
	if q.getTeamsByDivisionStmt, err = db.PrepareContext(ctx, getTeamsByDivision); err != nil {
		return nil, fmt.Errorf("error preparing query GetTeamsByDivision: %w", err)
	}
	if q.getTeamsForScoringUpdateStmt, err = db.PrepareContext(ctx, getTeamsForScoringUpdate); err != nil {
		return nil, fmt.Errorf("error preparing query GetTeamsForScoringUpdate: %w", err)
	}
	if q.getTopPlayersByPositionAndSeasonStmt, err = db.PrepareContext(ctx, getTopPlayersByPositionAndSeason); err != nil {
		return nil, fmt.Errorf("error preparing query GetTopPlayersByPositionAndSeason: %w", err)
	}
	if q.getTopScorersForSeasonStmt, err = db.PrepareContext(ctx, getTopScorersForSeason); err != nil {
		return nil, fmt.Errorf("error preparing query GetTopScorersForSeason: %w", err)
	}
	if q.getTopScorersForWeekStmt, err = db.PrepareContext(ctx, getTopScorersForWeek); err != nil {
		return nil, fmt.Errorf("error preparing query GetTopScorersForWeek: %w", err)
	}
	if q.getTotalWeeksInSeasonStmt, err = db.PrepareContext(ctx, getTotalWeeksInSeason); err != nil {
		return nil, fmt.Errorf("error preparing query GetTotalWeeksInSeason: %w", err)
	}
	if q.insertPlayerStmt, err = db.PrepareContext(ctx, insertPlayer); err != nil {
		return nil, fmt.Errorf("error preparing query InsertPlayer: %w", err)
	}
	if q.removePlayerFromFantasyTeamStmt, err = db.PrepareContext(ctx, removePlayerFromFantasyTeam); err != nil {
		return nil, fmt.Errorf("error preparing query RemovePlayerFromFantasyTeam: %w", err)
	}
	if q.searchPlayersStmt, err = db.PrepareContext(ctx, searchPlayers); err != nil {
		return nil, fmt.Errorf("error preparing query SearchPlayers: %w", err)
	}
	if q.setCurrentSeasonStmt, err = db.PrepareContext(ctx, setCurrentSeason); err != nil {
		return nil, fmt.Errorf("error preparing query SetCurrentSeason: %w", err)
	}
	if q.updateFantasyRosterStmt, err = db.PrepareContext(ctx, updateFantasyRoster); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateFantasyRoster: %w", err)
	}
	if q.updateFantasyTeamRecordStmt, err = db.PrepareContext(ctx, updateFantasyTeamRecord); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateFantasyTeamRecord: %w", err)
	}
	if q.updateMatchupScoreStmt, err = db.PrepareContext(ctx, updateMatchupScore); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateMatchupScore: %w", err)
	}
	if q.updateNFLGameScoreStmt, err = db.PrepareContext(ctx, updateNFLGameScore); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateNFLGameScore: %w", err)
	}
	if q.updatePlayerStmt, err = db.PrepareContext(ctx, updatePlayer); err != nil {
		return nil, fmt.Errorf("error preparing query UpdatePlayer: %w", err)
	}
	if q.updatePlayerStatsStmt, err = db.PrepareContext(ctx, updatePlayerStats); err != nil {
		return nil, fmt.Errorf("error preparing query UpdatePlayerStats: %w", err)
	}
	if q.updateScoringRulesStmt, err = db.PrepareContext(ctx, updateScoringRules); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateScoringRules: %w", err)
	}
	if q.updateSeasonStmt, err = db.PrepareContext(ctx, updateSeason); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateSeason: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.addDraftPickStmt != nil {
		if cerr := q.addDraftPickStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing addDraftPickStmt: %w", cerr)
		}
	}
	if q.addPlayerToFantasyTeamStmt != nil {
		if cerr := q.addPlayerToFantasyTeamStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing addPlayerToFantasyTeamStmt: %w", cerr)
		}
	}
	if q.calculateFantasyScoreStmt != nil {
		if cerr := q.calculateFantasyScoreStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing calculateFantasyScoreStmt: %w", cerr)
		}
	}
	if q.calculatePlayerFantasyPointsStmt != nil {
		if cerr := q.calculatePlayerFantasyPointsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing calculatePlayerFantasyPointsStmt: %w", cerr)
		}
	}
	if q.clearDraftStmt != nil {
		if cerr := q.clearDraftStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing clearDraftStmt: %w", cerr)
		}
	}
	if q.createFantasyTeamStmt != nil {
		if cerr := q.createFantasyTeamStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createFantasyTeamStmt: %w", cerr)
		}
	}
	if q.createLeagueStmt != nil {
		if cerr := q.createLeagueStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createLeagueStmt: %w", cerr)
		}
	}
	if q.createMatchupStmt != nil {
		if cerr := q.createMatchupStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createMatchupStmt: %w", cerr)
		}
	}
	if q.createNFLGameStmt != nil {
		if cerr := q.createNFLGameStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createNFLGameStmt: %w", cerr)
		}
	}
	if q.createPlayerStatsStmt != nil {
		if cerr := q.createPlayerStatsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createPlayerStatsStmt: %w", cerr)
		}
	}
	if q.createScoringRulesStmt != nil {
		if cerr := q.createScoringRulesStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createScoringRulesStmt: %w", cerr)
		}
	}
	if q.createSeasonStmt != nil {
		if cerr := q.createSeasonStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createSeasonStmt: %w", cerr)
		}
	}
	if q.getAllFantasyTeamsStmt != nil {
		if cerr := q.getAllFantasyTeamsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getAllFantasyTeamsStmt: %w", cerr)
		}
	}
	if q.getAllLeaguesStmt != nil {
		if cerr := q.getAllLeaguesStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getAllLeaguesStmt: %w", cerr)
		}
	}
	if q.getAllNFLTeamsStmt != nil {
		if cerr := q.getAllNFLTeamsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getAllNFLTeamsStmt: %w", cerr)
		}
	}
	if q.getAllPlayersStmt != nil {
		if cerr := q.getAllPlayersStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getAllPlayersStmt: %w", cerr)
		}
	}
	if q.getAllPlayersWithFantasyPointsStmt != nil {
		if cerr := q.getAllPlayersWithFantasyPointsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getAllPlayersWithFantasyPointsStmt: %w", cerr)
		}
	}
	if q.getAllSeasonsStmt != nil {
		if cerr := q.getAllSeasonsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getAllSeasonsStmt: %w", cerr)
		}
	}
	if q.getAvailableDraftPlayersStmt != nil {
		if cerr := q.getAvailableDraftPlayersStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getAvailableDraftPlayersStmt: %w", cerr)
		}
	}
	if q.getAvailablePlayersStmt != nil {
		if cerr := q.getAvailablePlayersStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getAvailablePlayersStmt: %w", cerr)
		}
	}
	if q.getBestAvailablePlayersStmt != nil {
		if cerr := q.getBestAvailablePlayersStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getBestAvailablePlayersStmt: %w", cerr)
		}
	}
	if q.getCurrentSeasonStmt != nil {
		if cerr := q.getCurrentSeasonStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getCurrentSeasonStmt: %w", cerr)
		}
	}
	if q.getDraftOrderStmt != nil {
		if cerr := q.getDraftOrderStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getDraftOrderStmt: %w", cerr)
		}
	}
	if q.getDraftPicksByTeamStmt != nil {
		if cerr := q.getDraftPicksByTeamStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getDraftPicksByTeamStmt: %w", cerr)
		}
	}
	if q.getDraftPicksForLeagueStmt != nil {
		if cerr := q.getDraftPicksForLeagueStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getDraftPicksForLeagueStmt: %w", cerr)
		}
	}
	if q.getDraftSummaryStmt != nil {
		if cerr := q.getDraftSummaryStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getDraftSummaryStmt: %w", cerr)
		}
	}
	if q.getFantasyTeamByIdStmt != nil {
		if cerr := q.getFantasyTeamByIdStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getFantasyTeamByIdStmt: %w", cerr)
		}
	}
	if q.getFantasyTeamForWeekStmt != nil {
		if cerr := q.getFantasyTeamForWeekStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getFantasyTeamForWeekStmt: %w", cerr)
		}
	}
	if q.getFantasyTeamRosterStmt != nil {
		if cerr := q.getFantasyTeamRosterStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getFantasyTeamRosterStmt: %w", cerr)
		}
	}
	if q.getHistoricalPlayerStatsStmt != nil {
		if cerr := q.getHistoricalPlayerStatsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getHistoricalPlayerStatsStmt: %w", cerr)
		}
	}
	if q.getLastDraftPickStmt != nil {
		if cerr := q.getLastDraftPickStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getLastDraftPickStmt: %w", cerr)
		}
	}
	if q.getLeagueByIdStmt != nil {
		if cerr := q.getLeagueByIdStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getLeagueByIdStmt: %w", cerr)
		}
	}
	if q.getLeagueMatchupsStmt != nil {
		if cerr := q.getLeagueMatchupsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getLeagueMatchupsStmt: %w", cerr)
		}
	}
	if q.getLeagueMatchupsByWeekStmt != nil {
		if cerr := q.getLeagueMatchupsByWeekStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getLeagueMatchupsByWeekStmt: %w", cerr)
		}
	}
	if q.getLeagueScoringRulesStmt != nil {
		if cerr := q.getLeagueScoringRulesStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getLeagueScoringRulesStmt: %w", cerr)
		}
	}
	if q.getLeagueStandingsStmt != nil {
		if cerr := q.getLeagueStandingsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getLeagueStandingsStmt: %w", cerr)
		}
	}
	if q.getLeagueWeeklyScoresStmt != nil {
		if cerr := q.getLeagueWeeklyScoresStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getLeagueWeeklyScoresStmt: %w", cerr)
		}
	}
	if q.getNFLScheduleForSeasonStmt != nil {
		if cerr := q.getNFLScheduleForSeasonStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getNFLScheduleForSeasonStmt: %w", cerr)
		}
	}
	if q.getNFLScheduleForWeekStmt != nil {
		if cerr := q.getNFLScheduleForWeekStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getNFLScheduleForWeekStmt: %w", cerr)
		}
	}
	if q.getNFLTeamByAbbreviationStmt != nil {
		if cerr := q.getNFLTeamByAbbreviationStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getNFLTeamByAbbreviationStmt: %w", cerr)
		}
	}
	if q.getNFLTeamByIdStmt != nil {
		if cerr := q.getNFLTeamByIdStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getNFLTeamByIdStmt: %w", cerr)
		}
	}
	if q.getNextDraftingTeamStmt != nil {
		if cerr := q.getNextDraftingTeamStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getNextDraftingTeamStmt: %w", cerr)
		}
	}
	if q.getPlayerByIdStmt != nil {
		if cerr := q.getPlayerByIdStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getPlayerByIdStmt: %w", cerr)
		}
	}
	if q.getPlayerFantasyPointsStmt != nil {
		if cerr := q.getPlayerFantasyPointsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getPlayerFantasyPointsStmt: %w", cerr)
		}
	}
	if q.getPlayerFantasyTotalPointsStmt != nil {
		if cerr := q.getPlayerFantasyTotalPointsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getPlayerFantasyTotalPointsStmt: %w", cerr)
		}
	}
	if q.getPlayerSeasonStatsStmt != nil {
		if cerr := q.getPlayerSeasonStatsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getPlayerSeasonStatsStmt: %w", cerr)
		}
	}
	if q.getPlayerStatsStmt != nil {
		if cerr := q.getPlayerStatsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getPlayerStatsStmt: %w", cerr)
		}
	}
	if q.getPlayerStatsByGameStmt != nil {
		if cerr := q.getPlayerStatsByGameStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getPlayerStatsByGameStmt: %w", cerr)
		}
	}
	if q.getPlayerStatsByWeekStmt != nil {
		if cerr := q.getPlayerStatsByWeekStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getPlayerStatsByWeekStmt: %w", cerr)
		}
	}
	if q.getPlayersByPositionStmt != nil {
		if cerr := q.getPlayersByPositionStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getPlayersByPositionStmt: %w", cerr)
		}
	}
	if q.getPlayersByTeamStmt != nil {
		if cerr := q.getPlayersByTeamStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getPlayersByTeamStmt: %w", cerr)
		}
	}
	if q.getPlayoffTeamsStmt != nil {
		if cerr := q.getPlayoffTeamsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getPlayoffTeamsStmt: %w", cerr)
		}
	}
	if q.getPreviousSeasonsStmt != nil {
		if cerr := q.getPreviousSeasonsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getPreviousSeasonsStmt: %w", cerr)
		}
	}
	if q.getSeasonByIdStmt != nil {
		if cerr := q.getSeasonByIdStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getSeasonByIdStmt: %w", cerr)
		}
	}
	if q.getSeasonByYearStmt != nil {
		if cerr := q.getSeasonByYearStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getSeasonByYearStmt: %w", cerr)
		}
	}
	if q.getTeamAtDraftPositionStmt != nil {
		if cerr := q.getTeamAtDraftPositionStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getTeamAtDraftPositionStmt: %w", cerr)
		}
	}
	if q.getTeamRosterStmt != nil {
		if cerr := q.getTeamRosterStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getTeamRosterStmt: %w", cerr)
		}
	}
	if q.getTeamScoreForWeekStmt != nil {
		if cerr := q.getTeamScoreForWeekStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getTeamScoreForWeekStmt: %w", cerr)
		}
	}
	if q.getTeamsByConferenceStmt != nil {
		if cerr := q.getTeamsByConferenceStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getTeamsByConferenceStmt: %w", cerr)
		}
	}
	if q.getTeamsByDivisionStmt != nil {
		if cerr := q.getTeamsByDivisionStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getTeamsByDivisionStmt: %w", cerr)
		}
	}
	if q.getTeamsForScoringUpdateStmt != nil {
		if cerr := q.getTeamsForScoringUpdateStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getTeamsForScoringUpdateStmt: %w", cerr)
		}
	}
	if q.getTopPlayersByPositionAndSeasonStmt != nil {
		if cerr := q.getTopPlayersByPositionAndSeasonStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getTopPlayersByPositionAndSeasonStmt: %w", cerr)
		}
	}
	if q.getTopScorersForSeasonStmt != nil {
		if cerr := q.getTopScorersForSeasonStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getTopScorersForSeasonStmt: %w", cerr)
		}
	}
	if q.getTopScorersForWeekStmt != nil {
		if cerr := q.getTopScorersForWeekStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getTopScorersForWeekStmt: %w", cerr)
		}
	}
	if q.getTotalWeeksInSeasonStmt != nil {
		if cerr := q.getTotalWeeksInSeasonStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getTotalWeeksInSeasonStmt: %w", cerr)
		}
	}
	if q.insertPlayerStmt != nil {
		if cerr := q.insertPlayerStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing insertPlayerStmt: %w", cerr)
		}
	}
	if q.removePlayerFromFantasyTeamStmt != nil {
		if cerr := q.removePlayerFromFantasyTeamStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing removePlayerFromFantasyTeamStmt: %w", cerr)
		}
	}
	if q.searchPlayersStmt != nil {
		if cerr := q.searchPlayersStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing searchPlayersStmt: %w", cerr)
		}
	}
	if q.setCurrentSeasonStmt != nil {
		if cerr := q.setCurrentSeasonStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing setCurrentSeasonStmt: %w", cerr)
		}
	}
	if q.updateFantasyRosterStmt != nil {
		if cerr := q.updateFantasyRosterStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateFantasyRosterStmt: %w", cerr)
		}
	}
	if q.updateFantasyTeamRecordStmt != nil {
		if cerr := q.updateFantasyTeamRecordStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateFantasyTeamRecordStmt: %w", cerr)
		}
	}
	if q.updateMatchupScoreStmt != nil {
		if cerr := q.updateMatchupScoreStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateMatchupScoreStmt: %w", cerr)
		}
	}
	if q.updateNFLGameScoreStmt != nil {
		if cerr := q.updateNFLGameScoreStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateNFLGameScoreStmt: %w", cerr)
		}
	}
	if q.updatePlayerStmt != nil {
		if cerr := q.updatePlayerStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updatePlayerStmt: %w", cerr)
		}
	}
	if q.updatePlayerStatsStmt != nil {
		if cerr := q.updatePlayerStatsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updatePlayerStatsStmt: %w", cerr)
		}
	}
	if q.updateScoringRulesStmt != nil {
		if cerr := q.updateScoringRulesStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateScoringRulesStmt: %w", cerr)
		}
	}
	if q.updateSeasonStmt != nil {
		if cerr := q.updateSeasonStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateSeasonStmt: %w", cerr)
		}
	}
	return err
}

func (q *Queries) exec(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (sql.Result, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).ExecContext(ctx, args...)
	case stmt != nil:
		return stmt.ExecContext(ctx, args...)
	default:
		return q.db.ExecContext(ctx, query, args...)
	}
}

func (q *Queries) query(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (*sql.Rows, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryContext(ctx, args...)
	default:
		return q.db.QueryContext(ctx, query, args...)
	}
}

func (q *Queries) queryRow(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) *sql.Row {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryRowContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryRowContext(ctx, args...)
	default:
		return q.db.QueryRowContext(ctx, query, args...)
	}
}

type Queries struct {
	db                                   DBTX
	tx                                   *sql.Tx
	addDraftPickStmt                     *sql.Stmt
	addPlayerToFantasyTeamStmt           *sql.Stmt
	calculateFantasyScoreStmt            *sql.Stmt
	calculatePlayerFantasyPointsStmt     *sql.Stmt
	clearDraftStmt                       *sql.Stmt
	createFantasyTeamStmt                *sql.Stmt
	createLeagueStmt                     *sql.Stmt
	createMatchupStmt                    *sql.Stmt
	createNFLGameStmt                    *sql.Stmt
	createPlayerStatsStmt                *sql.Stmt
	createScoringRulesStmt               *sql.Stmt
	createSeasonStmt                     *sql.Stmt
	getAllFantasyTeamsStmt               *sql.Stmt
	getAllLeaguesStmt                    *sql.Stmt
	getAllNFLTeamsStmt                   *sql.Stmt
	getAllPlayersStmt                    *sql.Stmt
	getAllPlayersWithFantasyPointsStmt   *sql.Stmt
	getAllSeasonsStmt                    *sql.Stmt
	getAvailableDraftPlayersStmt         *sql.Stmt
	getAvailablePlayersStmt              *sql.Stmt
	getBestAvailablePlayersStmt          *sql.Stmt
	getCurrentSeasonStmt                 *sql.Stmt
	getDraftOrderStmt                    *sql.Stmt
	getDraftPicksByTeamStmt              *sql.Stmt
	getDraftPicksForLeagueStmt           *sql.Stmt
	getDraftSummaryStmt                  *sql.Stmt
	getFantasyTeamByIdStmt               *sql.Stmt
	getFantasyTeamForWeekStmt            *sql.Stmt
	getFantasyTeamRosterStmt             *sql.Stmt
	getHistoricalPlayerStatsStmt         *sql.Stmt
	getLastDraftPickStmt                 *sql.Stmt
	getLeagueByIdStmt                    *sql.Stmt
	getLeagueMatchupsStmt                *sql.Stmt
	getLeagueMatchupsByWeekStmt          *sql.Stmt
	getLeagueScoringRulesStmt            *sql.Stmt
	getLeagueStandingsStmt               *sql.Stmt
	getLeagueWeeklyScoresStmt            *sql.Stmt
	getNFLScheduleForSeasonStmt          *sql.Stmt
	getNFLScheduleForWeekStmt            *sql.Stmt
	getNFLTeamByAbbreviationStmt         *sql.Stmt
	getNFLTeamByIdStmt                   *sql.Stmt
	getNextDraftingTeamStmt              *sql.Stmt
	getPlayerByIdStmt                    *sql.Stmt
	getPlayerFantasyPointsStmt           *sql.Stmt
	getPlayerFantasyTotalPointsStmt      *sql.Stmt
	getPlayerSeasonStatsStmt             *sql.Stmt
	getPlayerStatsStmt                   *sql.Stmt
	getPlayerStatsByGameStmt             *sql.Stmt
	getPlayerStatsByWeekStmt             *sql.Stmt
	getPlayersByPositionStmt             *sql.Stmt
	getPlayersByTeamStmt                 *sql.Stmt
	getPlayoffTeamsStmt                  *sql.Stmt
	getPreviousSeasonsStmt               *sql.Stmt
	getSeasonByIdStmt                    *sql.Stmt
	getSeasonByYearStmt                  *sql.Stmt
	getTeamAtDraftPositionStmt           *sql.Stmt
	getTeamRosterStmt                    *sql.Stmt
	getTeamScoreForWeekStmt              *sql.Stmt
	getTeamsByConferenceStmt             *sql.Stmt
	getTeamsByDivisionStmt               *sql.Stmt
	getTeamsForScoringUpdateStmt         *sql.Stmt
	getTopPlayersByPositionAndSeasonStmt *sql.Stmt
	getTopScorersForSeasonStmt           *sql.Stmt
	getTopScorersForWeekStmt             *sql.Stmt
	getTotalWeeksInSeasonStmt            *sql.Stmt
	insertPlayerStmt                     *sql.Stmt
	removePlayerFromFantasyTeamStmt      *sql.Stmt
	searchPlayersStmt                    *sql.Stmt
	setCurrentSeasonStmt                 *sql.Stmt
	updateFantasyRosterStmt              *sql.Stmt
	updateFantasyTeamRecordStmt          *sql.Stmt
	updateMatchupScoreStmt               *sql.Stmt
	updateNFLGameScoreStmt               *sql.Stmt
	updatePlayerStmt                     *sql.Stmt
	updatePlayerStatsStmt                *sql.Stmt
	updateScoringRulesStmt               *sql.Stmt
	updateSeasonStmt                     *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                                   tx,
		tx:                                   tx,
		addDraftPickStmt:                     q.addDraftPickStmt,
		addPlayerToFantasyTeamStmt:           q.addPlayerToFantasyTeamStmt,
		calculateFantasyScoreStmt:            q.calculateFantasyScoreStmt,
		calculatePlayerFantasyPointsStmt:     q.calculatePlayerFantasyPointsStmt,
		clearDraftStmt:                       q.clearDraftStmt,
		createFantasyTeamStmt:                q.createFantasyTeamStmt,
		createLeagueStmt:                     q.createLeagueStmt,
		createMatchupStmt:                    q.createMatchupStmt,
		createNFLGameStmt:                    q.createNFLGameStmt,
		createPlayerStatsStmt:                q.createPlayerStatsStmt,
		createScoringRulesStmt:               q.createScoringRulesStmt,
		createSeasonStmt:                     q.createSeasonStmt,
		getAllFantasyTeamsStmt:               q.getAllFantasyTeamsStmt,
		getAllLeaguesStmt:                    q.getAllLeaguesStmt,
		getAllNFLTeamsStmt:                   q.getAllNFLTeamsStmt,
		getAllPlayersStmt:                    q.getAllPlayersStmt,
		getAllPlayersWithFantasyPointsStmt:   q.getAllPlayersWithFantasyPointsStmt,
		getAllSeasonsStmt:                    q.getAllSeasonsStmt,
		getAvailableDraftPlayersStmt:         q.getAvailableDraftPlayersStmt,
		getAvailablePlayersStmt:              q.getAvailablePlayersStmt,
		getBestAvailablePlayersStmt:          q.getBestAvailablePlayersStmt,
		getCurrentSeasonStmt:                 q.getCurrentSeasonStmt,
		getDraftOrderStmt:                    q.getDraftOrderStmt,
		getDraftPicksByTeamStmt:              q.getDraftPicksByTeamStmt,
		getDraftPicksForLeagueStmt:           q.getDraftPicksForLeagueStmt,
		getDraftSummaryStmt:                  q.getDraftSummaryStmt,
		getFantasyTeamByIdStmt:               q.getFantasyTeamByIdStmt,
		getFantasyTeamForWeekStmt:            q.getFantasyTeamForWeekStmt,
		getFantasyTeamRosterStmt:             q.getFantasyTeamRosterStmt,
		getHistoricalPlayerStatsStmt:         q.getHistoricalPlayerStatsStmt,
		getLastDraftPickStmt:                 q.getLastDraftPickStmt,
		getLeagueByIdStmt:                    q.getLeagueByIdStmt,
		getLeagueMatchupsStmt:                q.getLeagueMatchupsStmt,
		getLeagueMatchupsByWeekStmt:          q.getLeagueMatchupsByWeekStmt,
		getLeagueScoringRulesStmt:            q.getLeagueScoringRulesStmt,
		getLeagueStandingsStmt:               q.getLeagueStandingsStmt,
		getLeagueWeeklyScoresStmt:            q.getLeagueWeeklyScoresStmt,
		getNFLScheduleForSeasonStmt:          q.getNFLScheduleForSeasonStmt,
		getNFLScheduleForWeekStmt:            q.getNFLScheduleForWeekStmt,
		getNFLTeamByAbbreviationStmt:         q.getNFLTeamByAbbreviationStmt,
		getNFLTeamByIdStmt:                   q.getNFLTeamByIdStmt,
		getNextDraftingTeamStmt:              q.getNextDraftingTeamStmt,
		getPlayerByIdStmt:                    q.getPlayerByIdStmt,
		getPlayerFantasyPointsStmt:           q.getPlayerFantasyPointsStmt,
		getPlayerFantasyTotalPointsStmt:      q.getPlayerFantasyTotalPointsStmt,
		getPlayerSeasonStatsStmt:             q.getPlayerSeasonStatsStmt,
		getPlayerStatsStmt:                   q.getPlayerStatsStmt,
		getPlayerStatsByGameStmt:             q.getPlayerStatsByGameStmt,
		getPlayerStatsByWeekStmt:             q.getPlayerStatsByWeekStmt,
		getPlayersByPositionStmt:             q.getPlayersByPositionStmt,
		getPlayersByTeamStmt:                 q.getPlayersByTeamStmt,
		getPlayoffTeamsStmt:                  q.getPlayoffTeamsStmt,
		getPreviousSeasonsStmt:               q.getPreviousSeasonsStmt,
		getSeasonByIdStmt:                    q.getSeasonByIdStmt,
		getSeasonByYearStmt:                  q.getSeasonByYearStmt,
		getTeamAtDraftPositionStmt:           q.getTeamAtDraftPositionStmt,
		getTeamRosterStmt:                    q.getTeamRosterStmt,
		getTeamScoreForWeekStmt:              q.getTeamScoreForWeekStmt,
		getTeamsByConferenceStmt:             q.getTeamsByConferenceStmt,
		getTeamsByDivisionStmt:               q.getTeamsByDivisionStmt,
		getTeamsForScoringUpdateStmt:         q.getTeamsForScoringUpdateStmt,
		getTopPlayersByPositionAndSeasonStmt: q.getTopPlayersByPositionAndSeasonStmt,
		getTopScorersForSeasonStmt:           q.getTopScorersForSeasonStmt,
		getTopScorersForWeekStmt:             q.getTopScorersForWeekStmt,
		getTotalWeeksInSeasonStmt:            q.getTotalWeeksInSeasonStmt,
		insertPlayerStmt:                     q.insertPlayerStmt,
		removePlayerFromFantasyTeamStmt:      q.removePlayerFromFantasyTeamStmt,
		searchPlayersStmt:                    q.searchPlayersStmt,
		setCurrentSeasonStmt:                 q.setCurrentSeasonStmt,
		updateFantasyRosterStmt:              q.updateFantasyRosterStmt,
		updateFantasyTeamRecordStmt:          q.updateFantasyTeamRecordStmt,
		updateMatchupScoreStmt:               q.updateMatchupScoreStmt,
		updateNFLGameScoreStmt:               q.updateNFLGameScoreStmt,
		updatePlayerStmt:                     q.updatePlayerStmt,
		updatePlayerStatsStmt:                q.updatePlayerStatsStmt,
		updateScoringRulesStmt:               q.updateScoringRulesStmt,
		updateSeasonStmt:                     q.updateSeasonStmt,
	}
}

// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package sqlc

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
	if q.createGameStmt, err = db.PrepareContext(ctx, createGame); err != nil {
		return nil, fmt.Errorf("error preparing query CreateGame: %w", err)
	}
	if q.createNFLPlayerStmt, err = db.PrepareContext(ctx, createNFLPlayer); err != nil {
		return nil, fmt.Errorf("error preparing query CreateNFLPlayer: %w", err)
	}
	if q.createNFLStatStmt, err = db.PrepareContext(ctx, createNFLStat); err != nil {
		return nil, fmt.Errorf("error preparing query CreateNFLStat: %w", err)
	}
	if q.createNFLTeamStmt, err = db.PrepareContext(ctx, createNFLTeam); err != nil {
		return nil, fmt.Errorf("error preparing query CreateNFLTeam: %w", err)
	}
	if q.deleteGameStmt, err = db.PrepareContext(ctx, deleteGame); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteGame: %w", err)
	}
	if q.deleteNFLPlayerStmt, err = db.PrepareContext(ctx, deleteNFLPlayer); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteNFLPlayer: %w", err)
	}
	if q.deleteNFLStatStmt, err = db.PrepareContext(ctx, deleteNFLStat); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteNFLStat: %w", err)
	}
	if q.deleteNFLTeamStmt, err = db.PrepareContext(ctx, deleteNFLTeam); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteNFLTeam: %w", err)
	}
	if q.getActiveNFLPlayersStmt, err = db.PrepareContext(ctx, getActiveNFLPlayers); err != nil {
		return nil, fmt.Errorf("error preparing query GetActiveNFLPlayers: %w", err)
	}
	if q.getAllGamesStmt, err = db.PrepareContext(ctx, getAllGames); err != nil {
		return nil, fmt.Errorf("error preparing query GetAllGames: %w", err)
	}
	if q.getAllGamesBySeasonAndWeekStmt, err = db.PrepareContext(ctx, getAllGamesBySeasonAndWeek); err != nil {
		return nil, fmt.Errorf("error preparing query GetAllGamesBySeasonAndWeek: %w", err)
	}
	if q.getAllNFLPlayersStmt, err = db.PrepareContext(ctx, getAllNFLPlayers); err != nil {
		return nil, fmt.Errorf("error preparing query GetAllNFLPlayers: %w", err)
	}
	if q.getAllNFLTeamsStmt, err = db.PrepareContext(ctx, getAllNFLTeams); err != nil {
		return nil, fmt.Errorf("error preparing query GetAllNFLTeams: %w", err)
	}
	if q.getGameStmt, err = db.PrepareContext(ctx, getGame); err != nil {
		return nil, fmt.Errorf("error preparing query GetGame: %w", err)
	}
	if q.getNFLPlayerStmt, err = db.PrepareContext(ctx, getNFLPlayer); err != nil {
		return nil, fmt.Errorf("error preparing query GetNFLPlayer: %w", err)
	}
	if q.getNFLTeamStmt, err = db.PrepareContext(ctx, getNFLTeam); err != nil {
		return nil, fmt.Errorf("error preparing query GetNFLTeam: %w", err)
	}
	if q.getPlayerStatAverageStmt, err = db.PrepareContext(ctx, getPlayerStatAverage); err != nil {
		return nil, fmt.Errorf("error preparing query GetPlayerStatAverage: %w", err)
	}
	if q.getPlayerStatsByGameStmt, err = db.PrepareContext(ctx, getPlayerStatsByGame); err != nil {
		return nil, fmt.Errorf("error preparing query GetPlayerStatsByGame: %w", err)
	}
	if q.getPlayerStatsByWeekStmt, err = db.PrepareContext(ctx, getPlayerStatsByWeek); err != nil {
		return nil, fmt.Errorf("error preparing query GetPlayerStatsByWeek: %w", err)
	}
	if q.getPlayerTotalStatByTypeStmt, err = db.PrepareContext(ctx, getPlayerTotalStatByType); err != nil {
		return nil, fmt.Errorf("error preparing query GetPlayerTotalStatByType: %w", err)
	}
	if q.getPlayerTotalStatByTypeForSeasonStmt, err = db.PrepareContext(ctx, getPlayerTotalStatByTypeForSeason); err != nil {
		return nil, fmt.Errorf("error preparing query GetPlayerTotalStatByTypeForSeason: %w", err)
	}
	if q.getPlayerTotalStatsByPositionStmt, err = db.PrepareContext(ctx, getPlayerTotalStatsByPosition); err != nil {
		return nil, fmt.Errorf("error preparing query GetPlayerTotalStatsByPosition: %w", err)
	}
	if q.getPlayerTotalStatsBySeasonStmt, err = db.PrepareContext(ctx, getPlayerTotalStatsBySeason); err != nil {
		return nil, fmt.Errorf("error preparing query GetPlayerTotalStatsBySeason: %w", err)
	}
	if q.getPlayerWeeklyStatByTypeStmt, err = db.PrepareContext(ctx, getPlayerWeeklyStatByType); err != nil {
		return nil, fmt.Errorf("error preparing query GetPlayerWeeklyStatByType: %w", err)
	}
	if q.getPlayersByPositionStmt, err = db.PrepareContext(ctx, getPlayersByPosition); err != nil {
		return nil, fmt.Errorf("error preparing query GetPlayersByPosition: %w", err)
	}
	if q.getPlayersByTeamStmt, err = db.PrepareContext(ctx, getPlayersByTeam); err != nil {
		return nil, fmt.Errorf("error preparing query GetPlayersByTeam: %w", err)
	}
	if q.getStatsByCategoryStmt, err = db.PrepareContext(ctx, getStatsByCategory); err != nil {
		return nil, fmt.Errorf("error preparing query GetStatsByCategory: %w", err)
	}
	if q.getStatsByGameStmt, err = db.PrepareContext(ctx, getStatsByGame); err != nil {
		return nil, fmt.Errorf("error preparing query GetStatsByGame: %w", err)
	}
	if q.getStatsByGameAndPlayerStmt, err = db.PrepareContext(ctx, getStatsByGameAndPlayer); err != nil {
		return nil, fmt.Errorf("error preparing query GetStatsByGameAndPlayer: %w", err)
	}
	if q.getStatsByPlayerStmt, err = db.PrepareContext(ctx, getStatsByPlayer); err != nil {
		return nil, fmt.Errorf("error preparing query GetStatsByPlayer: %w", err)
	}
	if q.getStatsByStatTypeStmt, err = db.PrepareContext(ctx, getStatsByStatType); err != nil {
		return nil, fmt.Errorf("error preparing query GetStatsByStatType: %w", err)
	}
	if q.getStatsByTeamStmt, err = db.PrepareContext(ctx, getStatsByTeam); err != nil {
		return nil, fmt.Errorf("error preparing query GetStatsByTeam: %w", err)
	}
	if q.getTeamStatsByCategoryStmt, err = db.PrepareContext(ctx, getTeamStatsByCategory); err != nil {
		return nil, fmt.Errorf("error preparing query GetTeamStatsByCategory: %w", err)
	}
	if q.getTeamsByConferenceStmt, err = db.PrepareContext(ctx, getTeamsByConference); err != nil {
		return nil, fmt.Errorf("error preparing query GetTeamsByConference: %w", err)
	}
	if q.getTeamsByDivisionStmt, err = db.PrepareContext(ctx, getTeamsByDivision); err != nil {
		return nil, fmt.Errorf("error preparing query GetTeamsByDivision: %w", err)
	}
	if q.getTopPlayersByStatStmt, err = db.PrepareContext(ctx, getTopPlayersByStat); err != nil {
		return nil, fmt.Errorf("error preparing query GetTopPlayersByStat: %w", err)
	}
	if q.searchPlayersStmt, err = db.PrepareContext(ctx, searchPlayers); err != nil {
		return nil, fmt.Errorf("error preparing query SearchPlayers: %w", err)
	}
	if q.updateGameStmt, err = db.PrepareContext(ctx, updateGame); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateGame: %w", err)
	}
	if q.updateNFLPlayerStmt, err = db.PrepareContext(ctx, updateNFLPlayer); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateNFLPlayer: %w", err)
	}
	if q.updateNFLStatStmt, err = db.PrepareContext(ctx, updateNFLStat); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateNFLStat: %w", err)
	}
	if q.updateNFLTeamStmt, err = db.PrepareContext(ctx, updateNFLTeam); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateNFLTeam: %w", err)
	}
	if q.upsertGameStmt, err = db.PrepareContext(ctx, upsertGame); err != nil {
		return nil, fmt.Errorf("error preparing query UpsertGame: %w", err)
	}
	if q.upsertNFLPlayerStmt, err = db.PrepareContext(ctx, upsertNFLPlayer); err != nil {
		return nil, fmt.Errorf("error preparing query UpsertNFLPlayer: %w", err)
	}
	if q.upsertNFLStatStmt, err = db.PrepareContext(ctx, upsertNFLStat); err != nil {
		return nil, fmt.Errorf("error preparing query UpsertNFLStat: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.createGameStmt != nil {
		if cerr := q.createGameStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createGameStmt: %w", cerr)
		}
	}
	if q.createNFLPlayerStmt != nil {
		if cerr := q.createNFLPlayerStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createNFLPlayerStmt: %w", cerr)
		}
	}
	if q.createNFLStatStmt != nil {
		if cerr := q.createNFLStatStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createNFLStatStmt: %w", cerr)
		}
	}
	if q.createNFLTeamStmt != nil {
		if cerr := q.createNFLTeamStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createNFLTeamStmt: %w", cerr)
		}
	}
	if q.deleteGameStmt != nil {
		if cerr := q.deleteGameStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteGameStmt: %w", cerr)
		}
	}
	if q.deleteNFLPlayerStmt != nil {
		if cerr := q.deleteNFLPlayerStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteNFLPlayerStmt: %w", cerr)
		}
	}
	if q.deleteNFLStatStmt != nil {
		if cerr := q.deleteNFLStatStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteNFLStatStmt: %w", cerr)
		}
	}
	if q.deleteNFLTeamStmt != nil {
		if cerr := q.deleteNFLTeamStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteNFLTeamStmt: %w", cerr)
		}
	}
	if q.getActiveNFLPlayersStmt != nil {
		if cerr := q.getActiveNFLPlayersStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getActiveNFLPlayersStmt: %w", cerr)
		}
	}
	if q.getAllGamesStmt != nil {
		if cerr := q.getAllGamesStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getAllGamesStmt: %w", cerr)
		}
	}
	if q.getAllGamesBySeasonAndWeekStmt != nil {
		if cerr := q.getAllGamesBySeasonAndWeekStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getAllGamesBySeasonAndWeekStmt: %w", cerr)
		}
	}
	if q.getAllNFLPlayersStmt != nil {
		if cerr := q.getAllNFLPlayersStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getAllNFLPlayersStmt: %w", cerr)
		}
	}
	if q.getAllNFLTeamsStmt != nil {
		if cerr := q.getAllNFLTeamsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getAllNFLTeamsStmt: %w", cerr)
		}
	}
	if q.getGameStmt != nil {
		if cerr := q.getGameStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getGameStmt: %w", cerr)
		}
	}
	if q.getNFLPlayerStmt != nil {
		if cerr := q.getNFLPlayerStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getNFLPlayerStmt: %w", cerr)
		}
	}
	if q.getNFLTeamStmt != nil {
		if cerr := q.getNFLTeamStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getNFLTeamStmt: %w", cerr)
		}
	}
	if q.getPlayerStatAverageStmt != nil {
		if cerr := q.getPlayerStatAverageStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getPlayerStatAverageStmt: %w", cerr)
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
	if q.getPlayerTotalStatByTypeStmt != nil {
		if cerr := q.getPlayerTotalStatByTypeStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getPlayerTotalStatByTypeStmt: %w", cerr)
		}
	}
	if q.getPlayerTotalStatByTypeForSeasonStmt != nil {
		if cerr := q.getPlayerTotalStatByTypeForSeasonStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getPlayerTotalStatByTypeForSeasonStmt: %w", cerr)
		}
	}
	if q.getPlayerTotalStatsByPositionStmt != nil {
		if cerr := q.getPlayerTotalStatsByPositionStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getPlayerTotalStatsByPositionStmt: %w", cerr)
		}
	}
	if q.getPlayerTotalStatsBySeasonStmt != nil {
		if cerr := q.getPlayerTotalStatsBySeasonStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getPlayerTotalStatsBySeasonStmt: %w", cerr)
		}
	}
	if q.getPlayerWeeklyStatByTypeStmt != nil {
		if cerr := q.getPlayerWeeklyStatByTypeStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getPlayerWeeklyStatByTypeStmt: %w", cerr)
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
	if q.getStatsByCategoryStmt != nil {
		if cerr := q.getStatsByCategoryStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getStatsByCategoryStmt: %w", cerr)
		}
	}
	if q.getStatsByGameStmt != nil {
		if cerr := q.getStatsByGameStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getStatsByGameStmt: %w", cerr)
		}
	}
	if q.getStatsByGameAndPlayerStmt != nil {
		if cerr := q.getStatsByGameAndPlayerStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getStatsByGameAndPlayerStmt: %w", cerr)
		}
	}
	if q.getStatsByPlayerStmt != nil {
		if cerr := q.getStatsByPlayerStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getStatsByPlayerStmt: %w", cerr)
		}
	}
	if q.getStatsByStatTypeStmt != nil {
		if cerr := q.getStatsByStatTypeStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getStatsByStatTypeStmt: %w", cerr)
		}
	}
	if q.getStatsByTeamStmt != nil {
		if cerr := q.getStatsByTeamStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getStatsByTeamStmt: %w", cerr)
		}
	}
	if q.getTeamStatsByCategoryStmt != nil {
		if cerr := q.getTeamStatsByCategoryStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getTeamStatsByCategoryStmt: %w", cerr)
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
	if q.getTopPlayersByStatStmt != nil {
		if cerr := q.getTopPlayersByStatStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getTopPlayersByStatStmt: %w", cerr)
		}
	}
	if q.searchPlayersStmt != nil {
		if cerr := q.searchPlayersStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing searchPlayersStmt: %w", cerr)
		}
	}
	if q.updateGameStmt != nil {
		if cerr := q.updateGameStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateGameStmt: %w", cerr)
		}
	}
	if q.updateNFLPlayerStmt != nil {
		if cerr := q.updateNFLPlayerStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateNFLPlayerStmt: %w", cerr)
		}
	}
	if q.updateNFLStatStmt != nil {
		if cerr := q.updateNFLStatStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateNFLStatStmt: %w", cerr)
		}
	}
	if q.updateNFLTeamStmt != nil {
		if cerr := q.updateNFLTeamStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateNFLTeamStmt: %w", cerr)
		}
	}
	if q.upsertGameStmt != nil {
		if cerr := q.upsertGameStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing upsertGameStmt: %w", cerr)
		}
	}
	if q.upsertNFLPlayerStmt != nil {
		if cerr := q.upsertNFLPlayerStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing upsertNFLPlayerStmt: %w", cerr)
		}
	}
	if q.upsertNFLStatStmt != nil {
		if cerr := q.upsertNFLStatStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing upsertNFLStatStmt: %w", cerr)
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
	db                                    DBTX
	tx                                    *sql.Tx
	createGameStmt                        *sql.Stmt
	createNFLPlayerStmt                   *sql.Stmt
	createNFLStatStmt                     *sql.Stmt
	createNFLTeamStmt                     *sql.Stmt
	deleteGameStmt                        *sql.Stmt
	deleteNFLPlayerStmt                   *sql.Stmt
	deleteNFLStatStmt                     *sql.Stmt
	deleteNFLTeamStmt                     *sql.Stmt
	getActiveNFLPlayersStmt               *sql.Stmt
	getAllGamesStmt                       *sql.Stmt
	getAllGamesBySeasonAndWeekStmt        *sql.Stmt
	getAllNFLPlayersStmt                  *sql.Stmt
	getAllNFLTeamsStmt                    *sql.Stmt
	getGameStmt                           *sql.Stmt
	getNFLPlayerStmt                      *sql.Stmt
	getNFLTeamStmt                        *sql.Stmt
	getPlayerStatAverageStmt              *sql.Stmt
	getPlayerStatsByGameStmt              *sql.Stmt
	getPlayerStatsByWeekStmt              *sql.Stmt
	getPlayerTotalStatByTypeStmt          *sql.Stmt
	getPlayerTotalStatByTypeForSeasonStmt *sql.Stmt
	getPlayerTotalStatsByPositionStmt     *sql.Stmt
	getPlayerTotalStatsBySeasonStmt       *sql.Stmt
	getPlayerWeeklyStatByTypeStmt         *sql.Stmt
	getPlayersByPositionStmt              *sql.Stmt
	getPlayersByTeamStmt                  *sql.Stmt
	getStatsByCategoryStmt                *sql.Stmt
	getStatsByGameStmt                    *sql.Stmt
	getStatsByGameAndPlayerStmt           *sql.Stmt
	getStatsByPlayerStmt                  *sql.Stmt
	getStatsByStatTypeStmt                *sql.Stmt
	getStatsByTeamStmt                    *sql.Stmt
	getTeamStatsByCategoryStmt            *sql.Stmt
	getTeamsByConferenceStmt              *sql.Stmt
	getTeamsByDivisionStmt                *sql.Stmt
	getTopPlayersByStatStmt               *sql.Stmt
	searchPlayersStmt                     *sql.Stmt
	updateGameStmt                        *sql.Stmt
	updateNFLPlayerStmt                   *sql.Stmt
	updateNFLStatStmt                     *sql.Stmt
	updateNFLTeamStmt                     *sql.Stmt
	upsertGameStmt                        *sql.Stmt
	upsertNFLPlayerStmt                   *sql.Stmt
	upsertNFLStatStmt                     *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                                    tx,
		tx:                                    tx,
		createGameStmt:                        q.createGameStmt,
		createNFLPlayerStmt:                   q.createNFLPlayerStmt,
		createNFLStatStmt:                     q.createNFLStatStmt,
		createNFLTeamStmt:                     q.createNFLTeamStmt,
		deleteGameStmt:                        q.deleteGameStmt,
		deleteNFLPlayerStmt:                   q.deleteNFLPlayerStmt,
		deleteNFLStatStmt:                     q.deleteNFLStatStmt,
		deleteNFLTeamStmt:                     q.deleteNFLTeamStmt,
		getActiveNFLPlayersStmt:               q.getActiveNFLPlayersStmt,
		getAllGamesStmt:                       q.getAllGamesStmt,
		getAllGamesBySeasonAndWeekStmt:        q.getAllGamesBySeasonAndWeekStmt,
		getAllNFLPlayersStmt:                  q.getAllNFLPlayersStmt,
		getAllNFLTeamsStmt:                    q.getAllNFLTeamsStmt,
		getGameStmt:                           q.getGameStmt,
		getNFLPlayerStmt:                      q.getNFLPlayerStmt,
		getNFLTeamStmt:                        q.getNFLTeamStmt,
		getPlayerStatAverageStmt:              q.getPlayerStatAverageStmt,
		getPlayerStatsByGameStmt:              q.getPlayerStatsByGameStmt,
		getPlayerStatsByWeekStmt:              q.getPlayerStatsByWeekStmt,
		getPlayerTotalStatByTypeStmt:          q.getPlayerTotalStatByTypeStmt,
		getPlayerTotalStatByTypeForSeasonStmt: q.getPlayerTotalStatByTypeForSeasonStmt,
		getPlayerTotalStatsByPositionStmt:     q.getPlayerTotalStatsByPositionStmt,
		getPlayerTotalStatsBySeasonStmt:       q.getPlayerTotalStatsBySeasonStmt,
		getPlayerWeeklyStatByTypeStmt:         q.getPlayerWeeklyStatByTypeStmt,
		getPlayersByPositionStmt:              q.getPlayersByPositionStmt,
		getPlayersByTeamStmt:                  q.getPlayersByTeamStmt,
		getStatsByCategoryStmt:                q.getStatsByCategoryStmt,
		getStatsByGameStmt:                    q.getStatsByGameStmt,
		getStatsByGameAndPlayerStmt:           q.getStatsByGameAndPlayerStmt,
		getStatsByPlayerStmt:                  q.getStatsByPlayerStmt,
		getStatsByStatTypeStmt:                q.getStatsByStatTypeStmt,
		getStatsByTeamStmt:                    q.getStatsByTeamStmt,
		getTeamStatsByCategoryStmt:            q.getTeamStatsByCategoryStmt,
		getTeamsByConferenceStmt:              q.getTeamsByConferenceStmt,
		getTeamsByDivisionStmt:                q.getTeamsByDivisionStmt,
		getTopPlayersByStatStmt:               q.getTopPlayersByStatStmt,
		searchPlayersStmt:                     q.searchPlayersStmt,
		updateGameStmt:                        q.updateGameStmt,
		updateNFLPlayerStmt:                   q.updateNFLPlayerStmt,
		updateNFLStatStmt:                     q.updateNFLStatStmt,
		updateNFLTeamStmt:                     q.updateNFLTeamStmt,
		upsertGameStmt:                        q.upsertGameStmt,
		upsertNFLPlayerStmt:                   q.upsertNFLPlayerStmt,
		upsertNFLStatStmt:                     q.upsertNFLStatStmt,
	}
}

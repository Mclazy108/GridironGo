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
	if q.createNFLTeamStmt, err = db.PrepareContext(ctx, createNFLTeam); err != nil {
		return nil, fmt.Errorf("error preparing query CreateNFLTeam: %w", err)
	}
	if q.deleteGameStmt, err = db.PrepareContext(ctx, deleteGame); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteGame: %w", err)
	}
	if q.deleteNFLPlayerStmt, err = db.PrepareContext(ctx, deleteNFLPlayer); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteNFLPlayer: %w", err)
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
	if q.getPlayersByPositionStmt, err = db.PrepareContext(ctx, getPlayersByPosition); err != nil {
		return nil, fmt.Errorf("error preparing query GetPlayersByPosition: %w", err)
	}
	if q.getPlayersByTeamStmt, err = db.PrepareContext(ctx, getPlayersByTeam); err != nil {
		return nil, fmt.Errorf("error preparing query GetPlayersByTeam: %w", err)
	}
	if q.getTeamsByConferenceStmt, err = db.PrepareContext(ctx, getTeamsByConference); err != nil {
		return nil, fmt.Errorf("error preparing query GetTeamsByConference: %w", err)
	}
	if q.getTeamsByDivisionStmt, err = db.PrepareContext(ctx, getTeamsByDivision); err != nil {
		return nil, fmt.Errorf("error preparing query GetTeamsByDivision: %w", err)
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
	if q.updateNFLTeamStmt, err = db.PrepareContext(ctx, updateNFLTeam); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateNFLTeam: %w", err)
	}
	if q.upsertGameStmt, err = db.PrepareContext(ctx, upsertGame); err != nil {
		return nil, fmt.Errorf("error preparing query UpsertGame: %w", err)
	}
	if q.upsertNFLPlayerStmt, err = db.PrepareContext(ctx, upsertNFLPlayer); err != nil {
		return nil, fmt.Errorf("error preparing query UpsertNFLPlayer: %w", err)
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
	db                             DBTX
	tx                             *sql.Tx
	createGameStmt                 *sql.Stmt
	createNFLPlayerStmt            *sql.Stmt
	createNFLTeamStmt              *sql.Stmt
	deleteGameStmt                 *sql.Stmt
	deleteNFLPlayerStmt            *sql.Stmt
	deleteNFLTeamStmt              *sql.Stmt
	getActiveNFLPlayersStmt        *sql.Stmt
	getAllGamesStmt                *sql.Stmt
	getAllGamesBySeasonAndWeekStmt *sql.Stmt
	getAllNFLPlayersStmt           *sql.Stmt
	getAllNFLTeamsStmt             *sql.Stmt
	getGameStmt                    *sql.Stmt
	getNFLPlayerStmt               *sql.Stmt
	getNFLTeamStmt                 *sql.Stmt
	getPlayersByPositionStmt       *sql.Stmt
	getPlayersByTeamStmt           *sql.Stmt
	getTeamsByConferenceStmt       *sql.Stmt
	getTeamsByDivisionStmt         *sql.Stmt
	searchPlayersStmt              *sql.Stmt
	updateGameStmt                 *sql.Stmt
	updateNFLPlayerStmt            *sql.Stmt
	updateNFLTeamStmt              *sql.Stmt
	upsertGameStmt                 *sql.Stmt
	upsertNFLPlayerStmt            *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                             tx,
		tx:                             tx,
		createGameStmt:                 q.createGameStmt,
		createNFLPlayerStmt:            q.createNFLPlayerStmt,
		createNFLTeamStmt:              q.createNFLTeamStmt,
		deleteGameStmt:                 q.deleteGameStmt,
		deleteNFLPlayerStmt:            q.deleteNFLPlayerStmt,
		deleteNFLTeamStmt:              q.deleteNFLTeamStmt,
		getActiveNFLPlayersStmt:        q.getActiveNFLPlayersStmt,
		getAllGamesStmt:                q.getAllGamesStmt,
		getAllGamesBySeasonAndWeekStmt: q.getAllGamesBySeasonAndWeekStmt,
		getAllNFLPlayersStmt:           q.getAllNFLPlayersStmt,
		getAllNFLTeamsStmt:             q.getAllNFLTeamsStmt,
		getGameStmt:                    q.getGameStmt,
		getNFLPlayerStmt:               q.getNFLPlayerStmt,
		getNFLTeamStmt:                 q.getNFLTeamStmt,
		getPlayersByPositionStmt:       q.getPlayersByPositionStmt,
		getPlayersByTeamStmt:           q.getPlayersByTeamStmt,
		getTeamsByConferenceStmt:       q.getTeamsByConferenceStmt,
		getTeamsByDivisionStmt:         q.getTeamsByDivisionStmt,
		searchPlayersStmt:              q.searchPlayersStmt,
		updateGameStmt:                 q.updateGameStmt,
		updateNFLPlayerStmt:            q.updateNFLPlayerStmt,
		updateNFLTeamStmt:              q.updateNFLTeamStmt,
		upsertGameStmt:                 q.upsertGameStmt,
		upsertNFLPlayerStmt:            q.upsertNFLPlayerStmt,
	}
}

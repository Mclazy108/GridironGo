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
	if q.createNFLTeamStmt, err = db.PrepareContext(ctx, createNFLTeam); err != nil {
		return nil, fmt.Errorf("error preparing query CreateNFLTeam: %w", err)
	}
	if q.deleteGameStmt, err = db.PrepareContext(ctx, deleteGame); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteGame: %w", err)
	}
	if q.deleteNFLTeamStmt, err = db.PrepareContext(ctx, deleteNFLTeam); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteNFLTeam: %w", err)
	}
	if q.getAllGamesStmt, err = db.PrepareContext(ctx, getAllGames); err != nil {
		return nil, fmt.Errorf("error preparing query GetAllGames: %w", err)
	}
	if q.getAllGamesBySeasonAndWeekStmt, err = db.PrepareContext(ctx, getAllGamesBySeasonAndWeek); err != nil {
		return nil, fmt.Errorf("error preparing query GetAllGamesBySeasonAndWeek: %w", err)
	}
	if q.getAllNFLTeamsStmt, err = db.PrepareContext(ctx, getAllNFLTeams); err != nil {
		return nil, fmt.Errorf("error preparing query GetAllNFLTeams: %w", err)
	}
	if q.getGameStmt, err = db.PrepareContext(ctx, getGame); err != nil {
		return nil, fmt.Errorf("error preparing query GetGame: %w", err)
	}
	if q.getNFLTeamStmt, err = db.PrepareContext(ctx, getNFLTeam); err != nil {
		return nil, fmt.Errorf("error preparing query GetNFLTeam: %w", err)
	}
	if q.getTeamsByConferenceStmt, err = db.PrepareContext(ctx, getTeamsByConference); err != nil {
		return nil, fmt.Errorf("error preparing query GetTeamsByConference: %w", err)
	}
	if q.getTeamsByDivisionStmt, err = db.PrepareContext(ctx, getTeamsByDivision); err != nil {
		return nil, fmt.Errorf("error preparing query GetTeamsByDivision: %w", err)
	}
	if q.updateGameStmt, err = db.PrepareContext(ctx, updateGame); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateGame: %w", err)
	}
	if q.updateNFLTeamStmt, err = db.PrepareContext(ctx, updateNFLTeam); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateNFLTeam: %w", err)
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
	if q.deleteNFLTeamStmt != nil {
		if cerr := q.deleteNFLTeamStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteNFLTeamStmt: %w", cerr)
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
	if q.getNFLTeamStmt != nil {
		if cerr := q.getNFLTeamStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getNFLTeamStmt: %w", cerr)
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
	if q.updateGameStmt != nil {
		if cerr := q.updateGameStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateGameStmt: %w", cerr)
		}
	}
	if q.updateNFLTeamStmt != nil {
		if cerr := q.updateNFLTeamStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateNFLTeamStmt: %w", cerr)
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
	createNFLTeamStmt              *sql.Stmt
	deleteGameStmt                 *sql.Stmt
	deleteNFLTeamStmt              *sql.Stmt
	getAllGamesStmt                *sql.Stmt
	getAllGamesBySeasonAndWeekStmt *sql.Stmt
	getAllNFLTeamsStmt             *sql.Stmt
	getGameStmt                    *sql.Stmt
	getNFLTeamStmt                 *sql.Stmt
	getTeamsByConferenceStmt       *sql.Stmt
	getTeamsByDivisionStmt         *sql.Stmt
	updateGameStmt                 *sql.Stmt
	updateNFLTeamStmt              *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                             tx,
		tx:                             tx,
		createGameStmt:                 q.createGameStmt,
		createNFLTeamStmt:              q.createNFLTeamStmt,
		deleteGameStmt:                 q.deleteGameStmt,
		deleteNFLTeamStmt:              q.deleteNFLTeamStmt,
		getAllGamesStmt:                q.getAllGamesStmt,
		getAllGamesBySeasonAndWeekStmt: q.getAllGamesBySeasonAndWeekStmt,
		getAllNFLTeamsStmt:             q.getAllNFLTeamsStmt,
		getGameStmt:                    q.getGameStmt,
		getNFLTeamStmt:                 q.getNFLTeamStmt,
		getTeamsByConferenceStmt:       q.getTeamsByConferenceStmt,
		getTeamsByDivisionStmt:         q.getTeamsByDivisionStmt,
		updateGameStmt:                 q.updateGameStmt,
		updateNFLTeamStmt:              q.updateNFLTeamStmt,
	}
}

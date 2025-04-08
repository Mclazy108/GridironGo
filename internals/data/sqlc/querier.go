// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package sqlc

import (
	"context"
	"database/sql"
)

type Querier interface {
	CreateGame(ctx context.Context, arg CreateGameParams) error
	CreateNFLPlayer(ctx context.Context, arg CreateNFLPlayerParams) error
	CreateNFLTeam(ctx context.Context, arg CreateNFLTeamParams) error
	DeleteGame(ctx context.Context, eventID int64) error
	DeleteNFLPlayer(ctx context.Context, playerID string) error
	DeleteNFLTeam(ctx context.Context, teamID string) error
	GetActiveNFLPlayers(ctx context.Context) ([]*NflPlayer, error)
	GetAllGames(ctx context.Context) ([]*NflGame, error)
	GetAllGamesBySeasonAndWeek(ctx context.Context, arg GetAllGamesBySeasonAndWeekParams) ([]*NflGame, error)
	GetAllNFLPlayers(ctx context.Context) ([]*NflPlayer, error)
	GetAllNFLTeams(ctx context.Context) ([]*NflTeam, error)
	GetGame(ctx context.Context, eventID int64) (*NflGame, error)
	GetNFLPlayer(ctx context.Context, playerID string) (*NflPlayer, error)
	GetNFLTeam(ctx context.Context, teamID string) (*NflTeam, error)
	GetPlayersByPosition(ctx context.Context, position string) ([]*NflPlayer, error)
	GetPlayersByTeam(ctx context.Context, teamID sql.NullString) ([]*NflPlayer, error)
	GetTeamsByConference(ctx context.Context, conference string) ([]*NflTeam, error)
	GetTeamsByDivision(ctx context.Context, division string) ([]*NflTeam, error)
	SearchPlayers(ctx context.Context, arg SearchPlayersParams) ([]*NflPlayer, error)
	UpdateGame(ctx context.Context, arg UpdateGameParams) error
	UpdateNFLPlayer(ctx context.Context, arg UpdateNFLPlayerParams) error
	UpdateNFLTeam(ctx context.Context, arg UpdateNFLTeamParams) error
	UpsertGame(ctx context.Context, arg UpsertGameParams) error
	UpsertNFLPlayer(ctx context.Context, arg UpsertNFLPlayerParams) error
}

var _ Querier = (*Queries)(nil)

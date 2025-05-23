// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: teams.sql

package sqlc

import (
	"context"
	"database/sql"
)

const createNFLTeam = `-- name: CreateNFLTeam :exec
INSERT INTO nfl_teams (
  team_id, display_name, abbreviation, short_name, location, nickname,
  conference, division, primary_color, secondary_color, logo_url
) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
)
`

type CreateNFLTeamParams struct {
	TeamID         string         `json:"team_id"`
	DisplayName    string         `json:"display_name"`
	Abbreviation   string         `json:"abbreviation"`
	ShortName      string         `json:"short_name"`
	Location       string         `json:"location"`
	Nickname       string         `json:"nickname"`
	Conference     string         `json:"conference"`
	Division       string         `json:"division"`
	PrimaryColor   sql.NullString `json:"primary_color"`
	SecondaryColor sql.NullString `json:"secondary_color"`
	LogoUrl        sql.NullString `json:"logo_url"`
}

func (q *Queries) CreateNFLTeam(ctx context.Context, arg CreateNFLTeamParams) error {
	_, err := q.exec(ctx, q.createNFLTeamStmt, createNFLTeam,
		arg.TeamID,
		arg.DisplayName,
		arg.Abbreviation,
		arg.ShortName,
		arg.Location,
		arg.Nickname,
		arg.Conference,
		arg.Division,
		arg.PrimaryColor,
		arg.SecondaryColor,
		arg.LogoUrl,
	)
	return err
}

const deleteNFLTeam = `-- name: DeleteNFLTeam :exec
DELETE FROM nfl_teams
WHERE team_id = ?
`

func (q *Queries) DeleteNFLTeam(ctx context.Context, teamID string) error {
	_, err := q.exec(ctx, q.deleteNFLTeamStmt, deleteNFLTeam, teamID)
	return err
}

const getAllNFLTeams = `-- name: GetAllNFLTeams :many
SELECT team_id, display_name, abbreviation, short_name, location, nickname, conference, division, primary_color, secondary_color, logo_url FROM nfl_teams
ORDER BY display_name ASC
`

func (q *Queries) GetAllNFLTeams(ctx context.Context) ([]*NflTeam, error) {
	rows, err := q.query(ctx, q.getAllNFLTeamsStmt, getAllNFLTeams)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*NflTeam{}
	for rows.Next() {
		var i NflTeam
		if err := rows.Scan(
			&i.TeamID,
			&i.DisplayName,
			&i.Abbreviation,
			&i.ShortName,
			&i.Location,
			&i.Nickname,
			&i.Conference,
			&i.Division,
			&i.PrimaryColor,
			&i.SecondaryColor,
			&i.LogoUrl,
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

const getNFLTeam = `-- name: GetNFLTeam :one
SELECT team_id, display_name, abbreviation, short_name, location, nickname, conference, division, primary_color, secondary_color, logo_url FROM nfl_teams
WHERE team_id = ?
`

func (q *Queries) GetNFLTeam(ctx context.Context, teamID string) (*NflTeam, error) {
	row := q.queryRow(ctx, q.getNFLTeamStmt, getNFLTeam, teamID)
	var i NflTeam
	err := row.Scan(
		&i.TeamID,
		&i.DisplayName,
		&i.Abbreviation,
		&i.ShortName,
		&i.Location,
		&i.Nickname,
		&i.Conference,
		&i.Division,
		&i.PrimaryColor,
		&i.SecondaryColor,
		&i.LogoUrl,
	)
	return &i, err
}

const getTeamsByConference = `-- name: GetTeamsByConference :many
SELECT team_id, display_name, abbreviation, short_name, location, nickname, conference, division, primary_color, secondary_color, logo_url FROM nfl_teams
WHERE conference = ?
ORDER BY division, display_name
`

func (q *Queries) GetTeamsByConference(ctx context.Context, conference string) ([]*NflTeam, error) {
	rows, err := q.query(ctx, q.getTeamsByConferenceStmt, getTeamsByConference, conference)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*NflTeam{}
	for rows.Next() {
		var i NflTeam
		if err := rows.Scan(
			&i.TeamID,
			&i.DisplayName,
			&i.Abbreviation,
			&i.ShortName,
			&i.Location,
			&i.Nickname,
			&i.Conference,
			&i.Division,
			&i.PrimaryColor,
			&i.SecondaryColor,
			&i.LogoUrl,
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

const getTeamsByDivision = `-- name: GetTeamsByDivision :many
SELECT team_id, display_name, abbreviation, short_name, location, nickname, conference, division, primary_color, secondary_color, logo_url FROM nfl_teams
WHERE division = ?
ORDER BY display_name
`

func (q *Queries) GetTeamsByDivision(ctx context.Context, division string) ([]*NflTeam, error) {
	rows, err := q.query(ctx, q.getTeamsByDivisionStmt, getTeamsByDivision, division)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*NflTeam{}
	for rows.Next() {
		var i NflTeam
		if err := rows.Scan(
			&i.TeamID,
			&i.DisplayName,
			&i.Abbreviation,
			&i.ShortName,
			&i.Location,
			&i.Nickname,
			&i.Conference,
			&i.Division,
			&i.PrimaryColor,
			&i.SecondaryColor,
			&i.LogoUrl,
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

const updateNFLTeam = `-- name: UpdateNFLTeam :exec
UPDATE nfl_teams
SET display_name = ?,
    abbreviation = ?,
    short_name = ?,
    location = ?,
    nickname = ?,
    conference = ?,
    division = ?,
    primary_color = ?,
    secondary_color = ?,
    logo_url = ?
WHERE team_id = ?
`

type UpdateNFLTeamParams struct {
	DisplayName    string         `json:"display_name"`
	Abbreviation   string         `json:"abbreviation"`
	ShortName      string         `json:"short_name"`
	Location       string         `json:"location"`
	Nickname       string         `json:"nickname"`
	Conference     string         `json:"conference"`
	Division       string         `json:"division"`
	PrimaryColor   sql.NullString `json:"primary_color"`
	SecondaryColor sql.NullString `json:"secondary_color"`
	LogoUrl        sql.NullString `json:"logo_url"`
	TeamID         string         `json:"team_id"`
}

func (q *Queries) UpdateNFLTeam(ctx context.Context, arg UpdateNFLTeamParams) error {
	_, err := q.exec(ctx, q.updateNFLTeamStmt, updateNFLTeam,
		arg.DisplayName,
		arg.Abbreviation,
		arg.ShortName,
		arg.Location,
		arg.Nickname,
		arg.Conference,
		arg.Division,
		arg.PrimaryColor,
		arg.SecondaryColor,
		arg.LogoUrl,
		arg.TeamID,
	)
	return err
}

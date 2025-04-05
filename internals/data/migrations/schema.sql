-- NFL Data Tables
CREATE TABLE nfl_games (
	event_id INTEGER PRIMARY KEY,
	date TEXT NOT NULL,
	name TEXT NOT NULL,
	short_name TEXT NOT NULL,
	season INTEGER NOT NULL,
	week INTEGER NOT NULL,
	away_team TEXT NOT NULL,
	home_team TEXT NOT NULL
);

CREATE TABLE nfl_teams (
    team_id TEXT PRIMARY KEY,
    display_name TEXT NOT NULL,
    abbreviation TEXT NOT NULL,
    short_name TEXT NOT NULL,
    location TEXT NOT NULL,
    nickname TEXT NOT NULL,
    conference TEXT NOT NULL,
    division TEXT NOT NULL,
    primary_color TEXT,
    secondary_color TEXT,
    logo_url TEXT
);

CREATE TABLE nfl_players (
    player_id TEXT PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    full_name TEXT NOT NULL,
    position TEXT NOT NULL,
    team_id TEXT,
    jersey TEXT,
    height INTEGER,
    weight INTEGER,
    active BOOLEAN NOT NULL,
    college TEXT,
    experience INTEGER,
    draft_year INTEGER,
    draft_round INTEGER,
    draft_pick INTEGER,
    status TEXT,
    image_url TEXT,
    FOREIGN KEY (team_id) REFERENCES nfl_teams (team_id)
);

CREATE INDEX idx_nfl_players_team ON nfl_players (team_id);
CREATE INDEX idx_nfl_players_position ON nfl_players (position);
CREATE INDEX idx_nfl_players_active ON nfl_players (active);

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

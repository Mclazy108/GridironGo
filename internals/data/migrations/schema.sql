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

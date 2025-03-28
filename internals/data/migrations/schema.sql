-- NFL Data Tables
CREATE TABLE seasons (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    year INTEGER NOT NULL UNIQUE,
    start_date TEXT NOT NULL,
    end_date TEXT NOT NULL,
    current INTEGER DEFAULT 0
);

CREATE TABLE nfl_teams (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    city TEXT NOT NULL,
    abbreviation TEXT NOT NULL UNIQUE,
    conference TEXT NOT NULL,
    division TEXT NOT NULL
);

CREATE TABLE nfl_players (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    position TEXT NOT NULL,
    team_id INTEGER,
    jersey_number INTEGER,
    status TEXT,
    FOREIGN KEY (team_id) REFERENCES nfl_teams(id)
);

CREATE TABLE nfl_games (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    season_id INTEGER NOT NULL,
    week INTEGER NOT NULL,
    home_team_id INTEGER NOT NULL,
    away_team_id INTEGER NOT NULL,
    home_score INTEGER,
    away_score INTEGER,
    game_date TEXT NOT NULL,
    game_time TEXT NOT NULL,
    status TEXT DEFAULT 'scheduled',
    FOREIGN KEY (season_id) REFERENCES seasons(id),
    FOREIGN KEY (home_team_id) REFERENCES nfl_teams(id),
    FOREIGN KEY (away_team_id) REFERENCES nfl_teams(id)
);

CREATE TABLE player_stats (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    player_id INTEGER NOT NULL,
    game_id INTEGER NOT NULL,
    season_id INTEGER NOT NULL,
    week INTEGER NOT NULL,
    
    -- Passing Stats
    passing_attempts INTEGER DEFAULT 0,
    passing_completions INTEGER DEFAULT 0,
    passing_yards INTEGER DEFAULT 0,
    passing_touchdowns INTEGER DEFAULT 0,
    passing_interceptions INTEGER DEFAULT 0,
    
    -- Rushing Stats
    rushing_attempts INTEGER DEFAULT 0,
    rushing_yards INTEGER DEFAULT 0,
    rushing_touchdowns INTEGER DEFAULT 0,
    
    -- Receiving Stats
    targets INTEGER DEFAULT 0,
    receptions INTEGER DEFAULT 0,
    receiving_yards INTEGER DEFAULT 0,
    receiving_touchdowns INTEGER DEFAULT 0,
    
    -- Kicking Stats
    field_goals_made INTEGER DEFAULT 0,
    field_goals_attempted INTEGER DEFAULT 0,
    extra_points_made INTEGER DEFAULT 0,
    extra_points_attempted INTEGER DEFAULT 0,
    
    -- Defense Stats
    sacks REAL DEFAULT 0,
    interceptions INTEGER DEFAULT 0,
    fumble_recoveries INTEGER DEFAULT 0,
    defensive_touchdowns INTEGER DEFAULT 0,
    safeties INTEGER DEFAULT 0,
    
    -- Misc Stats
    fumbles_lost INTEGER DEFAULT 0,
    two_point_conversions INTEGER DEFAULT 0,
    
    FOREIGN KEY (player_id) REFERENCES nfl_players(id),
    FOREIGN KEY (game_id) REFERENCES nfl_games(id),
    FOREIGN KEY (season_id) REFERENCES seasons(id),
    UNIQUE(player_id, game_id)
);

-- Fantasy Football Tables
CREATE TABLE fantasy_leagues (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    season_id INTEGER NOT NULL,
    teams_count INTEGER NOT NULL,
    qb_count INTEGER DEFAULT 1,
    rb_count INTEGER DEFAULT 2,
    wr_count INTEGER DEFAULT 2,
    te_count INTEGER DEFAULT 1,
    flex_count INTEGER DEFAULT 1,
    k_count INTEGER DEFAULT 1,
    dst_count INTEGER DEFAULT 1,
    bench_count INTEGER DEFAULT 6,
    ppr INTEGER DEFAULT 1,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (season_id) REFERENCES seasons(id)
);

CREATE TABLE fantasy_scoring_rules (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    league_id INTEGER NOT NULL,
    
    -- Passing
    passing_yards_per_point REAL DEFAULT 25.0,
    passing_touchdown_points REAL DEFAULT 4.0,
    passing_interception_points REAL DEFAULT -2.0,
    
    -- Rushing
    rushing_yards_per_point REAL DEFAULT 10.0,
    rushing_touchdown_points REAL DEFAULT 6.0,
    
    -- Receiving
    receiving_yards_per_point REAL DEFAULT 10.0,
    receiving_touchdown_points REAL DEFAULT 6.0,
    reception_points REAL DEFAULT 1.0,
    
    -- Kicking
    field_goal_0_39_points REAL DEFAULT 3.0,
    field_goal_40_49_points REAL DEFAULT 4.0,
    field_goal_50_plus_points REAL DEFAULT 5.0,
    extra_point_points REAL DEFAULT 1.0,
    
    -- Defense
    sack_points REAL DEFAULT 1.0,
    interception_points REAL DEFAULT 2.0,
    fumble_recovery_points REAL DEFAULT 2.0,
    defensive_touchdown_points REAL DEFAULT 6.0,
    safety_points REAL DEFAULT 2.0,
    
    -- Misc
    two_point_conversion_points REAL DEFAULT 2.0,
    fumble_lost_points REAL DEFAULT -2.0,
    
    FOREIGN KEY (league_id) REFERENCES fantasy_leagues(id),
    UNIQUE(league_id)
);

-- Modified table that was causing issues
CREATE TABLE fantasy_teams (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    league_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    owner_name TEXT NOT NULL,
    is_user INTEGER,
    draft_position INTEGER,
    wins INTEGER,
    losses INTEGER,
    tie_games INTEGER,
    points_for REAL,
    points_against REAL,
    FOREIGN KEY (league_id) REFERENCES fantasy_leagues(id),
    UNIQUE(league_id, name)
);

CREATE TABLE fantasy_rosters (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    team_id INTEGER NOT NULL,
    player_id INTEGER NOT NULL,
    position TEXT NOT NULL,
    is_starter INTEGER,
    FOREIGN KEY (team_id) REFERENCES fantasy_teams(id),
    FOREIGN KEY (player_id) REFERENCES nfl_players(id),
    UNIQUE(team_id, player_id)
);

CREATE TABLE fantasy_matchups (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    league_id INTEGER NOT NULL,
    week INTEGER NOT NULL,
    home_team_id INTEGER NOT NULL,
    away_team_id INTEGER NOT NULL,
    home_score REAL,
    away_score REAL,
    is_playoff INTEGER,
    is_championship INTEGER,
    completed INTEGER,
    FOREIGN KEY (league_id) REFERENCES fantasy_leagues(id),
    FOREIGN KEY (home_team_id) REFERENCES fantasy_teams(id),
    FOREIGN KEY (away_team_id) REFERENCES fantasy_teams(id),
    UNIQUE(league_id, week, home_team_id),
    UNIQUE(league_id, week, away_team_id)
);

CREATE TABLE fantasy_drafts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    league_id INTEGER NOT NULL,
    team_id INTEGER NOT NULL,
    player_id INTEGER NOT NULL,
    round INTEGER NOT NULL,
    pick_number INTEGER NOT NULL,
    draft_time TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (league_id) REFERENCES fantasy_leagues(id),
    FOREIGN KEY (team_id) REFERENCES fantasy_teams(id),
    FOREIGN KEY (player_id) REFERENCES nfl_players(id),
    UNIQUE(league_id, player_id)
);

CREATE TABLE fantasy_player_scores (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    league_id INTEGER NOT NULL,
    team_id INTEGER,
    player_id INTEGER NOT NULL,
    week INTEGER NOT NULL,
    season_id INTEGER NOT NULL,
    points REAL,
    FOREIGN KEY (league_id) REFERENCES fantasy_leagues(id),
    FOREIGN KEY (team_id) REFERENCES fantasy_teams(id),
    FOREIGN KEY (player_id) REFERENCES nfl_players(id),
    FOREIGN KEY (season_id) REFERENCES seasons(id),
    UNIQUE(league_id, player_id, week, season_id)
);

-- Create indexes for better performance
CREATE INDEX idx_nfl_players_position ON nfl_players(position);
CREATE INDEX idx_nfl_players_team ON nfl_players(team_id);
CREATE INDEX idx_nfl_games_season_week ON nfl_games(season_id, week);
CREATE INDEX idx_player_stats_player ON player_stats(player_id);
CREATE INDEX idx_player_stats_game ON player_stats(game_id);
CREATE INDEX idx_player_stats_season_week ON player_stats(season_id, week);
CREATE INDEX idx_fantasy_rosters_team ON fantasy_rosters(team_id);
CREATE INDEX idx_fantasy_rosters_player ON fantasy_rosters(player_id);
CREATE INDEX idx_fantasy_matchups_league_week ON fantasy_matchups(league_id, week);
CREATE INDEX idx_fantasy_player_scores_player_week ON fantasy_player_scores(player_id, week);

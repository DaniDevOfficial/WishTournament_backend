-- Drop the database if it exists and create a new one
DROP DATABASE IF EXISTS wishtournament;
CREATE DATABASE wishtournament;
\c wishtournament;  -- Connect to the newly created database

-- Create extension for UUID generation (PostgreSQL specific)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Table: users
CREATE TABLE users (
                       id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                       username VARCHAR(255) NOT NULL,
                       email VARCHAR(255) UNIQUE NOT NULL,
                       password VARCHAR(255) NOT NULL,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table: tournaments
CREATE TABLE tournaments (
                             id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                             name VARCHAR(255) NOT NULL,
                             host_id UUID REFERENCES users(id) ON DELETE CASCADE,
                             point_type VARCHAR(50) NOT NULL,
                             scoring_method VARCHAR(50) NOT NULL,
                             can_tie BOOLEAN DEFAULT false,
                             start_date_time TIMESTAMP NOT NULL,
                             is_started BOOLEAN DEFAULT false,
                             created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table: teams
CREATE TABLE teams (
                       id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                       name VARCHAR(255) NOT NULL,
                       tournament_id UUID REFERENCES tournaments(id) ON DELETE CASCADE,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table: players
CREATE TABLE players (
                         id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                         name VARCHAR(255) NOT NULL,
                         is_registered BOOLEAN DEFAULT false,
                         user_id UUID REFERENCES users(id) ON DELETE SET NULL,
                         team_id UUID REFERENCES teams(id) ON DELETE CASCADE,
                         tournament_id UUID REFERENCES tournaments(id) ON DELETE CASCADE,
                         created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table: brackets
CREATE TABLE brackets (
                          id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                          tournament_id UUID REFERENCES tournaments(id) ON DELETE CASCADE,
                          type VARCHAR(50) NOT NULL,
                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table: rounds
CREATE TABLE rounds (
                        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                        bracket_id UUID REFERENCES brackets(id) ON DELETE CASCADE,
                        round_number INT NOT NULL,
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table: games
CREATE TABLE games (
                       id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                       round_id UUID REFERENCES rounds(id) ON DELETE CASCADE,
                       team1_id UUID REFERENCES teams(id) ON DELETE CASCADE,
                       team2_id UUID REFERENCES teams(id) ON DELETE CASCADE,
                       team1_score FLOAT DEFAULT 0,
                       team2_score FLOAT DEFAULT 0,
                       winner_id UUID REFERENCES teams(id) ON DELETE SET NULL,
                       is_tie BOOLEAN DEFAULT false,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Mock Data Insertions

-- Users
INSERT INTO users (id, username, email, password) VALUES
                                                      (uuid_generate_v4(), 'hostuser', 'host@tournament.com', 'passwordhash1'),
                                                      (uuid_generate_v4(), 'player1', 'player1@game.com', 'passwordhash2'),
                                                      (uuid_generate_v4(), 'player2', 'player2@game.com', 'passwordhash3');

-- Tournaments
INSERT INTO tournaments (id, name, host_id, point_type, scoring_method, can_tie, start_date_time, is_started) VALUES
                                                                                                                  (uuid_generate_v4(), 'Summer Tournament', (SELECT id FROM users WHERE username = 'hostuser'), 'points', 'highest', false, '2024-11-10 10:00:00', false),
                                                                                                                  (uuid_generate_v4(), 'Winter Tournament', (SELECT id FROM users WHERE username = 'hostuser'), 'time', 'lowest', true, '2024-12-20 14:00:00', false);

-- Teams
INSERT INTO teams (id, name, tournament_id) VALUES
                                                (uuid_generate_v4(), 'Team Alpha', (SELECT id FROM tournaments WHERE name = 'Summer Tournament')),
                                                (uuid_generate_v4(), 'Team Beta', (SELECT id FROM tournaments WHERE name = 'Summer Tournament')),
                                                (uuid_generate_v4(), 'Team Gamma', (SELECT id FROM tournaments WHERE name = 'Winter Tournament'));

-- Players (anonymous and registered)
INSERT INTO players (id, name, is_registered, user_id, team_id, tournament_id) VALUES
                                                                                   (uuid_generate_v4(), 'Anonymous Player 1', false, NULL, (SELECT id FROM teams WHERE name = 'Team Alpha'), (SELECT id FROM tournaments WHERE name = 'Summer Tournament')),
                                                                                   (uuid_generate_v4(), 'player1', true, (SELECT id FROM users WHERE username = 'player1'), (SELECT id FROM teams WHERE name = 'Team Beta'), (SELECT id FROM tournaments WHERE name = 'Summer Tournament')),
                                                                                   (uuid_generate_v4(), 'player2', true, (SELECT id FROM users WHERE username = 'player2'), (SELECT id FROM teams WHERE name = 'Team Gamma'), (SELECT id FROM tournaments WHERE name = 'Winter Tournament'));

-- Brackets
INSERT INTO brackets (id, tournament_id, type) VALUES
                                                   (uuid_generate_v4(), (SELECT id FROM tournaments WHERE name = 'Summer Tournament'), 'single_elimination'),
                                                   (uuid_generate_v4(), (SELECT id FROM tournaments WHERE name = 'Winter Tournament'), 'single_elimination');

-- Rounds
INSERT INTO rounds (id, bracket_id, round_number) VALUES
                                                      (uuid_generate_v4(), (SELECT id FROM brackets WHERE tournament_id = (SELECT id FROM tournaments WHERE name = 'Summer Tournament')), 1),
                                                      (uuid_generate_v4(), (SELECT id FROM brackets WHERE tournament_id = (SELECT id FROM tournaments WHERE name = 'Winter Tournament')), 1);

-- Games (Game 1: Team Alpha vs Team Beta)
INSERT INTO games (id, round_id, team1_id, team2_id, team1_score, team2_score, winner_id, is_tie) VALUES
    (uuid_generate_v4(), (SELECT id FROM rounds WHERE bracket_id = (SELECT id FROM brackets WHERE tournament_id = (SELECT id FROM tournaments WHERE name = 'Summer Tournament'))),
     (SELECT id FROM teams WHERE name = 'Team Alpha'), (SELECT id FROM teams WHERE name = 'Team Beta'), 95, 100, (SELECT id FROM teams WHERE name = 'Team Beta'), false);

-- Games (Game 2: Team Gamma vs Bye)
INSERT INTO games (id, round_id, team1_id, team2_id, team1_score, team2_score, winner_id, is_tie) VALUES
    (uuid_generate_v4(), (SELECT id FROM rounds WHERE bracket_id = (SELECT id FROM brackets WHERE tournament_id = (SELECT id FROM tournaments WHERE name = 'Winter Tournament'))),
     (SELECT id FROM teams WHERE name = 'Team Gamma'), NULL, 0, 0, (SELECT id FROM teams WHERE name = 'Team Gamma'), false);

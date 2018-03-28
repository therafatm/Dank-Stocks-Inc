CREATE USER seng468 WITH PASSWORD 'SENG$^*';
CREATE DATABASE transactions OWNER seng468;
GRANT ALL PRIVILEGES ON DATABASE transactions TO seng468;

\c transactions; 

CREATE TABLE Users (
    uid serial PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    money INTEGER NOT NULL
);

CREATE INDEX users_username on Users (username);

CREATE TABLE Stocks (
    sid serial PRIMARY KEY,
    username VARCHAR(50) REFERENCES Users(username) ON DELETE CASCADE ON UPDATE CASCADE,
    symbol VARCHAR(10) NOT NULL,
    shares INTEGER NOT NULL,
    UNIQUE (username, symbol)
);

CREATE INDEX stock_username_symbol on Stocks (username, symbol);

CREATE TABLE Reservations (
    rid serial PRIMARY KEY,
    username VARCHAR(50) REFERENCES Users(username) ON DELETE CASCADE ON UPDATE CASCADE,
    symbol VARCHAR(10) NOT NULL,
    type VARCHAR(10) NOT NULL,
    shares INTEGER NOT NULL,
    amount INTEGER NOT NULL,
    time BIGINT NOT NULL
);

CREATE INDEX reservations_username on Reservations(username);

CREATE TABLE Triggers (
    tid serial PRIMARY KEY,
    username VARCHAR(50) REFERENCES Users(username) ON DELETE CASCADE ON UPDATE CASCADE,
    symbol VARCHAR(10) NOT NULL,
    type VARCHAR(10) NOT NULL,
    amount INTEGER NOT NULL,
    trigger_price INTEGER NOT NULL,
    executable BOOLEAN NOT NULL,
    time BIGINT NOT NULL,
    UNIQUE (username, symbol, type)
);

CREATE INDEX triggers_username_symbol_type on Triggers (username, symbol, type);
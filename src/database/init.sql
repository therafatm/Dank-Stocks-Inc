CREATE USER seng468 WITH PASSWORD 'SENG$^*';
CREATE DATABASE transactions OWNER seng468;
GRANT ALL PRIVILEGES ON DATABASE transactions TO seng468;

\c transactions; 

CREATE TABLE Users (
    uid serial PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    money DOUBLE PRECISION NOT NULL
);

CREATE TABLE Stocks (
    sid serial PRIMARY KEY,
    username VARCHAR(50) REFERENCES Users(username),
    symbol VARCHAR(10) NOT NULL UNIQUE,
    shares INTEGER NOT NULL
);

CREATE TABLE Reservations (
    rid serial PRIMARY KEY,
    username VARCHAR(50) REFERENCES Users(username),
    symbol VARCHAR(10),
    type VARCHAR(10),
    shares INTEGER NOT NULL,
    face_value DOUBLE PRECISION NOT NULL,
    time BIGINT NOT NULL
);

CREATE TABLE Triggers (
    tid serial PRIMARY KEY,
    username VARCHAR(50) REFERENCES Users(username),
    symbol VARCHAR(10) NOT NULL,
    type VARCHAR(10) NOT NULL,
    amount DOUBLE PRECISION,
    shares INTEGER,
    trigger_price DOUBLE PRECISION
);
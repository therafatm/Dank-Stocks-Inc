CREATE USER seng468 WITH PASSWORD 'SENG$^*';
CREATE DATABASE logs OWNER seng468;
GRANT ALL PRIVILEGES ON DATABASE logs TO seng468;

\c logs; 

CREATE TABLE Errors (
    uid serial PRIMARY KEY,
    timestamp TIME NOT NULL
    server VARCHAR(20) NOT NULL
    transactionNum INTEGER NOT NULL
    command VARCHAR(20) NOT NULL
    username VARCHAR(50) NOT NULL UNIQUE,
    funds INTEGER
    errorMessage VARCHAR(200) NOT NULL
    runnumber INTEGER

);

CREATE TABLE SystemEvents (
    uid serial PRIMARY KEY,
    timestamp TIME NOT NULL
    server VARCHAR(20) NOT NULL
    command VARCHAR(20) NOT NULL
    username VARCHAR(50) NOT NULL UNIQUE,
    stocksymbol VARCHAR(4) NOT NULL
    funds INTEGER NOT NULL
    runnumber INTEGER
);

CREATE TABLE UserCommand (
    uid serial PRIMARY KEY,
    timestamp TIME NOT NULL
    server VARCHAR(20) NOT NULL
    transactionNum INTEGER NOT NULL
    command VARCHAR(20) NOT NULL
    username VARCHAR(50) NOT NULL UNIQUE,
    stocksymbol VARCHAR(4) NOT NULL
    funds INTEGER NOT NULL
    runnumber INTEGER
);

CREATE TABLE QuoteServer (
    uid serial PRIMARY KEY,
    timestamp TIME NOT NULL
    server VARCHAR(20) NOT NULL
    quoteServerTime TIME NOT NULL
    username VARCHAR(50) NOT NULL UNIQUE,
    stocksymbol VARCHAR(4) NOT NULL
    money INTEGER NOT NULL
    cryptokey VARCHAR(50) NOT NULL UNIQUE,
    runnumber INTEGER
);
CREATE TABLE AccountTransaction (
    uid serial PRIMARY KEY,
    timestamp TIME NOT NULL
    server VARCHAR(20) NOT NULL
    transactionNum INTEGER NOT NULL
    action VARCHAR(20)
    username VARCHAR(50) NOT NULL UNIQUE,
    funds INTEGER
);


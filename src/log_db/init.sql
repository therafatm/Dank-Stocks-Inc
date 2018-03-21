CREATE USER seng468 WITH PASSWORD 'SENG$^*';
CREATE DATABASE logs OWNER seng468;
GRANT ALL PRIVILEGES ON DATABASE logs TO seng468;

\c logs; 
CREATE EXTENSION citus;

CREATE TABLE Errors (
    uid serial PRIMARY KEY,
    timestamp BIGINT NOT NULL,
    server VARCHAR(20),
    transactionNum INTEGER,
    command VARCHAR(20),
    username VARCHAR(50),
    funds VARCHAR(50),
    errorMessage VARCHAR(200),
    runnumber INTEGER

);

CREATE TABLE SystemEvents (
    uid serial PRIMARY KEY,
    timestamp BIGINT NOT NULL,
    server VARCHAR(20),
    command VARCHAR(20),
    username VARCHAR(50),
    stocksymbol VARCHAR(4),
    funds VARCHAR(50),
    runnumber INTEGER
);

CREATE TABLE UserCommand (
    uid serial PRIMARY KEY,
    timestamp BIGINT NOT NULL,
    server VARCHAR(20),
    transactionNum INTEGER,
    command VARCHAR(20),
    username VARCHAR(50),
    stocksymbol VARCHAR(4),
    funds VARCHAR(50),
    runnumber INTEGER
);

CREATE TABLE QuoteServer (
    uid serial PRIMARY KEY,
    timestamp TIME NOT NULL,
    server VARCHAR(20),
    quoteServerTime INTEGER,
    username VARCHAR(50),
    stocksymbol VARCHAR(4),
    money VARCHAR(50),
    cryptokey VARCHAR(50),
    runnumber INTEGER
);
CREATE TABLE AccountTransaction (
    uid serial PRIMARY KEY,
    timestamp TIME NOT NULL,
    server VARCHAR(20),
    transactionNum INTEGER,
    action VARCHAR(20),
    username VARCHAR(50),
    funds VARCHAR(50)
);


CREATE OR REPLACE FUNCTION distribute () RETURNS void as $$
    BEGIN 
        PERFORM create_distributed_table('UserCommand', 'uid');
        PERFORM create_distributed_table('Errors', 'uid');
        PERFORM create_distributed_table('SystemEvents', 'uid');
        PERFORM create_distributed_table('QuoteServer', 'uid');
        PERFORM create_distributed_table('AccountTransaction', 'uid');
        RETURN;
    END;
$$ LANGUAGE plpgsql;

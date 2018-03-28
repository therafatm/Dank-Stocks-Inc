CREATE USER seng468 WITH PASSWORD 'SENG$^*';
CREATE DATABASE logs OWNER seng468;
GRANT ALL PRIVILEGES ON DATABASE logs TO seng468;

\c logs; 
CREATE EXTENSION citus;

CREATE TABLE Errors (
    timestamp bigint NOT NULL,
    server VARCHAR(20),
    transactionnum INTEGER,
    command VARCHAR(20),
    username VARCHAR(50),
    funds VARCHAR(50),
    errorMessage VARCHAR(200),
    runnumber INTEGER

);

CREATE TABLE SystemEvents (
    timestamp bigint NOT NULL,
    server VARCHAR(20),
    transactionnum INTEGER,
    command VARCHAR(20),
    username VARCHAR(50),
    stocksymbol VARCHAR(4),
    funds VARCHAR(50),
    runnumber INTEGER
);

CREATE TABLE UserCommand (
    timestamp bigint NOT NULL,
    server VARCHAR(20),
    transactionnum INTEGER,
    command VARCHAR(20),
    username VARCHAR(50),
    stocksymbol VARCHAR(4),
    funds VARCHAR(50),
    runnumber INTEGER
);

CREATE TABLE QuoteServer (
    timestamp bigint NOT NULL,
    server VARCHAR(20),
    transactionnum INTEGER,
    quoteServerTime bigint,
    username VARCHAR(50),
    stocksymbol VARCHAR(4),
    money VARCHAR(50),
    cryptokey VARCHAR(100),
    runnumber INTEGER
);
CREATE TABLE AccountTransaction (
    timestamp bigint NOT NULL,
    server VARCHAR(20),
    transactionnum INTEGER,
    action VARCHAR(20),
    username VARCHAR(50),
    funds VARCHAR(50)
);


CREATE OR REPLACE FUNCTION distribute () RETURNS void as $$
    BEGIN 
        PERFORM create_distributed_table('UserCommand', 'transactionnum');
        PERFORM create_distributed_table('Errors', 'transactionnum');
        PERFORM create_distributed_table('SystemEvents', 'username');
        PERFORM create_distributed_table('QuoteServer', 'transactionnum');
        PERFORM create_distributed_table('AccountTransaction', 'transactionnum');
        RETURN;
    END;
$$ LANGUAGE plpgsql;

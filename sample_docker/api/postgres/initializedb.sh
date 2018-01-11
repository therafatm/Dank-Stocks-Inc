#bin .sh/
psql postgres -c "create database disasterreliefdev;"
psql postgres -c "create extension postgis;"
psql postgres -c "create user dev with password 'dev'; alter user dev with superuser;"


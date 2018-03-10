FROM postgres:alpine
    
EXPOSE 5432
COPY postgresql.conf /tmp/postgresql.conf
COPY updateConfig.sh /docker-entrypoint-initdb.d/_updateConfig.sh
ADD init.sql /docker-entrypoint-initdb.d/


FROM postgres:9.6

RUN apt-get update && apt-get install -y \
        postgis \
        postgresql-9.6-postgis-2.3
        
#ADD initializedb.sh /usr/local/bin
#RUN chmod +x /usr/local/bin/initializedb.sh
#RUN /usr/local/bin/initializedb.sh
EXPOSE 5432
ADD init.sql /docker-entrypoint-initdb.d/


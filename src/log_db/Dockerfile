FROM citusdata/citus:7.3-alpine

RUN apk add --update --no-cache curl \
        curl-dev \
        libcurl

COPY scripts/wait-for-it.sh /usr/local/bin/wait-for-it.sh
RUN chmod +x /usr/local/bin/wait-for-it.sh

COPY log_db/connect-worker.sh /usr/local/bin/connect-worker.sh
RUN chmod +x /usr/local/bin/connect-worker.sh

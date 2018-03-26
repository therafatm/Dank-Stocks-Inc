#!/usr/bin/env bash

export MY_HOST=`hostname -i`
echo $HOSTNAME
export JSON=$(cat <<EOF
	{
		"master_host": "$LOG_DB_HOST",
		"master_port": "$LOG_DB_PORT",
		"worker_host": "$MY_HOST",
		"worker_port": "5432",
		"db": "$LOG_DB",
		"user": "$PGUSER",
		"password": "$PGPASSWORD"
	}
EOF
)


wait-for-it.sh -h manager -p 3000 -t 20
export URL="http://$MANAGER_HOST:$MANAGER_PORT"
echo "Sending join request to manager $JSON at $URL"
curl -d "$JSON" -H "Content-Type: application/json/" -X POST $URL
exec "$@"
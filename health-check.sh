#!/bin/bash

LOG_FILE="go_scanner_health_check.log"

# URL of the health check endpoint
HEALTH_URL="http://localhost:8080/health"

# Function to check health
check_health() {
    HTTP_RESPONSE=$(curl --silent --write-out "HTTPSTATUS:%{http_code}" -X GET $HEALTH_URL)
    HTTP_BODY=$(echo $HTTP_RESPONSE | sed -e 's/HTTPSTATUS\:.*//g')
    HTTP_STATUS=$(echo $HTTP_RESPONSE | tr -d '\n' | sed -e 's/.*HTTPSTATUS://')

    if [ "$HTTP_STATUS" -ne 200 ]; then
        echo "$(date): Service is down, restarting..." >> $LOG_FILE
        sudo systemctl restart go-scanner.service
    else
        echo "$(date): Service is running fine." >> $LOG_FILE
    fi
}

check_health
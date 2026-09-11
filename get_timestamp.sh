#!/bin/bash

# Usage: ./timestamp_ago.sh <seconds>
# Example: ./timestamp_ago.sh 300
# Output: timestamp from 300 seconds ago

if [ -z "$1" ]; then
    echo "Usage: $0 <seconds>"
    exit 1
fi

SECONDS_AGO=$1
CURRENT_TIME=$(date +%s)
TIMESTAMP=$((CURRENT_TIME - SECONDS_AGO))

echo "$TIMESTAMP"

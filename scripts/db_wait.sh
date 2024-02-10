#!/bin/sh

DB_HOST=$1
DB_PORT=$2

MAX_ATTEMPTS=20

SLEEP_INTERVAL=2

attempt=1
while ! nc -z $DB_HOST $DB_PORT; do
  if [ $attempt -gt $MAX_ATTEMPTS ]; then
    echo "Attempted to connect to database at $DB_HOST:$DB_PORT over $MAX_ATTEMPTS times. Giving up."
    exit 1
  fi

  echo "Waiting for database at $DB_HOST:$DB_PORT to be up (attempt: $attempt)..."
  sleep $SLEEP_INTERVAL
  attempt=$(( attempt + 1 ))
done

echo "Database at $DB_HOST:$DB_PORT is up!"

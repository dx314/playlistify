#!/bin/bash

command="up"
if [ "$1" == "down" ]; then
  command="down"
fi

# Get the directory where the script resides
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

# Get the parent directory of the script
PARENT_DIR="$(dirname "$DIR")"

if [ "$1" == "drop" ]; then
  echo "$PLAYLISTIFY_DBPASSWORD" | docker-compose exec -T postgres psql -U dev -d plailist -c "drop schema public cascade;" -p 5432 -v ON_ERROR_STOP=1
  echo "$PLAYLISTIFY_DBPASSWORD" | docker-compose exec -T postgres psql -U dev -d plailist -c "create schema public;" -p 5432 -v ON_ERROR_STOP=1
  exit 0
fi

migrate -path ${PARENT_DIR}/migrations -database postgres://${PLAYLISTIFY_DBUSERNAME}:${PLAYLISTIFY_DBPASSWORD}@${PLAYLISTIFY_DBHOST}:${PLAYLISTIFY_DBPORT}/${PLAYLISTIFY_DBNAME}?sslmode=disable $command

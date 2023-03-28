#!/bin/bash

migrate -path ./migrations -database postgres://${PLAYLISTIFY_DBUSERNAME}:${PLAYLISTIFY_DBPASSWORD}@${PLAYLISTIFY_DBHOST}:${PLAYLISTIFY_DBPORT}/${PLAYLISTIFY_DBNAME}?sslmode=disable up

#!/bin/bash

DBSTRING="host=$POSTGRES_HOST user=$POSTGRES_USER password=$POSTGRES_PASSWORD dbname=$POSTGRES_DB sslmode=disable"

goose postgres "$DBSTRING" up
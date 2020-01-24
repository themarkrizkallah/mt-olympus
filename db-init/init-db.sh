#!/bin/bash

#echo "SELECT 'CREATE DATABASE $POSTGRES_DB' WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '$POSTGRES_DB')\gexec" \
#  | psql -U "$POSTGRES_USER

#psql -v ON_ERROR_STOP=1 -U "$POSTGRES_USER" --dbname "$POSTGRES_DB" -f docker-entrypoint-initdb.d/tables.sql
#echo "lmao"
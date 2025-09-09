#!/bin/bash
set -e

# Create databases
# docker exec -it postgres psql -U omides248 -d postgres
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "postgres" <<-EOSQL
    CREATE DATABASE identity_db;
    CREATE DATABASE wallet_db;
EOSQL

# Enable pgcrypto on database
# docker exec -it postgres psql -U omides248 -d postgres
# docker exec -it postgres psql -U omides248 -d identity -c "CREATE EXTENSION IF NOT EXISTS pgcrypto;"
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "identity_db" <<-EOSQL
    CREATE EXTENSION IF NOT EXISTS pgcrypto;
EOSQL

# Enable pgcrypto on database
# docker exec -it postgres psql -U omides248 -d postgres
# docker exec -it postgres psql -U omides248 -d wallet_db -c "CREATE EXTENSION IF NOT EXISTS pgcrypto;"
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "wallet_db" <<-EOSQL
    CREATE EXTENSION IF NOT EXISTS pgcrypto;
EOSQL
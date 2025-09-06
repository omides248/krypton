
sudo -i -u postgres
psql

CREATE DATABASE identity_db;
CREATE USER omides248 WITH PASSWORD '123123';
ALTER DATABASE identity_db OWNER TO omides248;
CREATE EXTENSION IF NOT EXISTS pgcrypto;

\c identity_db
ALTER SCHEMA public OWNER TO omides248;
GRANT ALL PRIVILEGES ON SCHEMA public TO omides248;

psql -h 127.0.0.1 -U omides248 -d identity_db
CREATE USER docker WITH PASSWORD 'docker';
CREATE DATABASE docker;
GRANT ALL PRIVILEGES ON DATABASE docker TO docker;

CREATE TABLE IF NOT EXISTS users(
  id BIGSERIAL PRIMARY KEY,
  username VARCHAR(30)  NOT NULL,
  email VARCHAR(30)  NOT NULL,
  password VARCHAR(120) NOT NULL,
  pic VARCHAR(120) DEFAULT NULL,
  lvl INTEGER DEFAULT 0,
  score INTEGER DEFAULT 0
);

CREATE UNIQUE INDEX IF NOT EXISTS "users_username_uindex" ON users (username);
CREATE UNIQUE INDEX IF NOT EXISTS "users_score_uindex" ON users (score);
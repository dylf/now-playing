DROP DATABASE IF EXISTS now_playing;
CREATE DATABASE now_playing;

\c now_playing;

-- Users table
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  -- email VARCHAR(255) NOT NULL,
  info json NOT NULL
)

-- TABLE "users" --
CREATE TABLE IF NOT EXISTS "users" (
  "id"        INT PRIMARY KEY,
  "username"  TEXT,
  "firstname" TEXT,
  "lastname"  TEXT,
  "state"     INT DEFAULT 0 NOT NULL
);

CREATE TABLE IF NOT EXISTS "insta_users" (
  "username"  TEXT PRIMARY KEY,
  "followers" INT,
  "following" INT
);

CREATE TABLE IF NOT EXISTS "subscriptions" (
  "user_id"        INT PRIMARY KEY REFERENCES "users" ("id"),
  "insta_username" TEXT NOT NULL REFERENCES "insta_users" ("username")
);

CREATE TYPE "group" AS ENUM('following', 'followers');

CREATE TABLE IF NOT EXISTS "following_followers" (
  "username"       TEXT PRIMARY KEY,
  "fullname"       TEXT,
  "URL"            TEXT,
  "refer_username" TEXT NOT NULL REFERENCES "insta_users" ("username") ON DELETE CASCADE,
  "group_type"     "group"
);
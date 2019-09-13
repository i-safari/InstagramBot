-- TABLE "states" --
CREATE TABLE IF NOT EXISTS "users" (
  "id"        INT PRIMARY KEY,
  "username"  TEXT,
  "firstname" TEXT,
  "lastname"  TEXT,
  "state"     INT DEFAULT 0 NOT NULL
);

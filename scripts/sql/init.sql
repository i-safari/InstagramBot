-- TABLE "states" --
CREATE TABLE IF NOT EXISTS states (
  "user_id"     INT PRIMARY KEY,
  "state"       INT DEFAULT 0 NOT NULL
);
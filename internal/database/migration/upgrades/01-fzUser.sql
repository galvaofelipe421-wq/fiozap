-- v0 -> v1: Create fzUser table

CREATE TABLE IF NOT EXISTS "fzUser" (
    "id" VARCHAR(64) PRIMARY KEY,
    "name" VARCHAR(255) NOT NULL,
    "token" VARCHAR(255) NOT NULL UNIQUE,
    "maxSessions" INTEGER DEFAULT 5,
    "createdAt" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

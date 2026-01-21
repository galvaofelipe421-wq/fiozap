-- v1 -> v2: Create fzSession table

CREATE TABLE IF NOT EXISTS "fzSession" (
    "id" VARCHAR(64) PRIMARY KEY,
    "userId" VARCHAR(64) NOT NULL REFERENCES "fzUser"("id") ON DELETE CASCADE,
    "name" VARCHAR(255) NOT NULL,
    "jid" VARCHAR(255) DEFAULT '',
    "qrCode" TEXT DEFAULT '',
    "connected" INTEGER DEFAULT 0,
    "webhook" TEXT DEFAULT '',
    "events" TEXT DEFAULT '',
    "proxyUrl" TEXT DEFAULT '',
    "createdAt" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE("userId", "name")
);

CREATE INDEX IF NOT EXISTS "idxFzSessionUserId" ON "fzSession" ("userId");
CREATE INDEX IF NOT EXISTS "idxFzSessionConnected" ON "fzSession" ("connected");

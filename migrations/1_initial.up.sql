CREATE TABLE IF NOT EXISTS "user" (
    id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE,
    "state" VARCHAR(50),
    "status" VARCHAR(50),
    password_hash VARCHAR(255),
    "version" INTEGER,
    created_at TIMESTAMPTZ DEFAULT TIMEZONE('UTC', NOW())
);
CREATE INDEX IF NOT EXISTS "user_email_idx" ON "user" (email);

CREATE TABLE IF NOT EXISTS "user_version" (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "user_id" UUID REFERENCES "user" ON DELETE CASCADE,
    email VARCHAR(255) UNIQUE,
    "state" VARCHAR(50),
    "status" VARCHAR(50),
    password_hash VARCHAR(255),
    "version" INTEGER,
    created_at TIMESTAMPTZ DEFAULT TIMEZONE('UTC', NOW())
)

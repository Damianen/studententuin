-- Drop WebApplication and Database tables
-- These tables are being moved to a different storage system

-- Drop foreign key constraints first
ALTER TABLE "WebApplication" DROP CONSTRAINT IF EXISTS "WebApplication_userId_fkey";
ALTER TABLE "Database" DROP CONSTRAINT IF EXISTS "Database_userId_fkey";

-- Drop the tables
DROP TABLE IF EXISTS "WebApplication";
DROP TABLE IF EXISTS "Database";

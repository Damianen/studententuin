-- AlterTable
ALTER TABLE "Database" ADD COLUMN "subdomain" TEXT,
ALTER COLUMN "dockerImage" DROP NOT NULL;

-- Update existing records with default values
UPDATE "Database" SET subdomain = CONCAT('db-', id) WHERE subdomain IS NULL;

-- Make subdomain required after updating existing data
ALTER TABLE "Database" ALTER COLUMN "subdomain" SET NOT NULL;

-- CreateIndex
CREATE UNIQUE INDEX "Database_subdomain_key" ON "Database"("subdomain");

-- AlterTable
ALTER TABLE "WebApplication" DROP COLUMN "url",
ADD COLUMN "subdomain" TEXT,
ADD COLUMN "runtime" TEXT,
ALTER COLUMN "dockerImage" DROP NOT NULL;

-- Update existing records with default values
UPDATE "WebApplication" SET subdomain = CONCAT('app-', id), runtime = 'nodejs' WHERE subdomain IS NULL;

-- Make fields required after updating existing data
ALTER TABLE "WebApplication" ALTER COLUMN "subdomain" SET NOT NULL,
ALTER COLUMN "runtime" SET NOT NULL;

-- CreateIndex
CREATE UNIQUE INDEX "WebApplication_subdomain_key" ON "WebApplication"("subdomain");

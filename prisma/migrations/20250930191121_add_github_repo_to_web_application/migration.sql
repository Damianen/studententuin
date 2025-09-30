-- AlterTable
ALTER TABLE "WebApplication" ADD COLUMN     "githubBranch" TEXT DEFAULT 'main',
ADD COLUMN     "githubRepo" TEXT;

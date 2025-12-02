import postgres from "postgres"
import { readFileSync } from "fs"
import { join } from "path"

async function main() {
  const connectionString = process.env.DATABASE_URL
  if (!connectionString) {
    throw new Error("DATABASE_URL is not set")
  }

  const sql = postgres(connectionString)

  try {
    console.log("Running migration to drop WebApplication and Database tables...")

    const migrationSQL = readFileSync(
      join(process.cwd(), "drizzle", "0001_drop_webapp_and_database_tables.sql"),
      "utf-8"
    )

    await sql.unsafe(migrationSQL)

    console.log("✓ Migration completed successfully!")
    console.log("✓ WebApplication and Database tables have been dropped")
  } catch (error) {
    console.error("✗ Migration failed:", error)
    process.exit(1)
  } finally {
    await sql.end()
  }
}

main()

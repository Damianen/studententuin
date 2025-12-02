import { auth } from "@/modules/auth/infrastructure/auth"
import { NextResponse } from "next/server"
import { db } from "@/lib/db"
import { accounts } from "@/lib/db/schema"
import { eq, and } from "drizzle-orm"

export async function GET() {
    try {
        const session = await auth()

        if (!session?.user?.id) {
            return NextResponse.json({ error: "Unauthorized" }, { status: 401 })
        }

        // Get GitHub access token from the user's account
        const accountResults = await db
            .select()
            .from(accounts)
            .where(
                and(
                    eq(accounts.userId, session.user.id),
                    eq(accounts.provider, "github")
                )
            )
            .limit(1)

        const account = accountResults[0]

        if (!account?.access_token) {
            return NextResponse.json(
                { error: "GitHub account not connected" },
                { status: 400 }
            )
        }

        // Fetch repositories from GitHub API
        const response = await fetch("https://api.github.com/user/repos?per_page=100&sort=updated", {
            headers: {
                Authorization: `Bearer ${account.access_token}`,
                Accept: "application/vnd.github.v3+json",
            },
        })

        if (!response.ok) {
            const errorData = await response.json().catch(() => ({}))
            console.error("GitHub API error:", response.status, errorData)
            throw new Error(`GitHub API error: ${response.status}`)
        }

        const repos = await response.json()

        // Return simplified repo data
        return NextResponse.json(
            repos.map((repo: any) => ({
                id: repo.id,
                name: repo.name,
                full_name: repo.full_name,
                private: repo.private,
                default_branch: repo.default_branch,
            }))
        )
    } catch (error) {
        console.error("Error fetching GitHub repos:", error)
        return NextResponse.json(
            { error: "Failed to fetch repositories" },
            { status: 500 }
        )
    }
}

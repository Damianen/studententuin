"use client"

import * as React from "react"
import { IconBrandGithub, IconSearch } from "@tabler/icons-react"
import { Label } from "@/components/ui/label"
import { Input } from "@/components/ui/input"
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select"
import { Button } from "@/components/ui/button"
import { Card } from "@/components/ui/card"

type GitHubRepo = {
    id: number
    name: string
    full_name: string
    private: boolean
    default_branch: string
}

export function GitHubRepoSelector() {
    const [repos, setRepos] = React.useState<GitHubRepo[]>([])
    const [loading, setLoading] = React.useState(false)
    const [error, setError] = React.useState<string | null>(null)
    const [selectedRepo, setSelectedRepo] = React.useState<string>("")
    const [branch, setBranch] = React.useState<string>("main")
    const [searchQuery, setSearchQuery] = React.useState("")

    React.useEffect(() => {
        fetchRepos()
    }, [])

    async function fetchRepos() {
        setLoading(true)
        setError(null)
        try {
            const response = await fetch("/api/github/repos")
            const data = await response.json()

            if (!response.ok) {
                throw new Error(data.error || "Failed to fetch repositories")
            }

            setRepos(data)
        } catch (err) {
            setError(err instanceof Error ? err.message : "Failed to load repositories")
        } finally {
            setLoading(false)
        }
    }

    const filteredRepos = repos.filter((repo) =>
        repo.full_name.toLowerCase().includes(searchQuery.toLowerCase())
    )

    if (loading) {
        return (
            <Card className="p-6">
                <div className="flex items-center gap-2 text-muted-foreground">
                    <IconBrandGithub className="h-5 w-5 animate-spin" />
                    <span>Loading GitHub repositories...</span>
                </div>
            </Card>
        )
    }

    if (error) {
        const isGitHubNotConnected = error.includes("GitHub account not connected")

        return (
            <Card className="p-6 border-dashed">
                <div className="space-y-4">
                    <div className="flex items-start gap-3">
                        <IconBrandGithub className="h-5 w-5 text-muted-foreground shrink-0 mt-0.5" />
                        <div className="flex-1 space-y-2">
                            <div className="font-medium">GitHub Repository (Optional)</div>
                            {isGitHubNotConnected ? (
                                <>
                                    <div className="text-sm text-muted-foreground">
                                        Connect your GitHub account to enable repository selection and automatic deployments.
                                    </div>
                                    <div className="text-sm text-muted-foreground">
                                        You can add a repository later from your application settings.
                                    </div>
                                </>
                            ) : (
                                <>
                                    <div className="text-sm text-red-600">{error}</div>
                                    <Button onClick={fetchRepos} variant="outline" size="sm">
                                        Try Again
                                    </Button>
                                </>
                            )}
                        </div>
                    </div>
                </div>
            </Card>
        )
    }

    return (
        <div className="space-y-4">
            <div className="space-y-2">
                <Label htmlFor="githubRepo">GitHub Repository (Optional)</Label>
                <div className="relative">
                    <IconSearch className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                    <Input
                        id="githubRepoSearch"
                        placeholder="Search repositories..."
                        value={searchQuery}
                        onChange={(e) => setSearchQuery(e.target.value)}
                        className="pl-9"
                    />
                </div>
            </div>

            <div className="space-y-2">
                <Select
                    name="githubRepo"
                    value={selectedRepo}
                    onValueChange={(value) => {
                        setSelectedRepo(value)
                        const repo = repos.find((r) => r.full_name === value)
                        if (repo) {
                            setBranch(repo.default_branch)
                        }
                    }}
                >
                    <SelectTrigger>
                        <SelectValue placeholder="Select a repository..." />
                    </SelectTrigger>
                    <SelectContent>
                        {filteredRepos.length === 0 ? (
                            <div className="p-2 text-sm text-muted-foreground text-center">
                                {searchQuery ? "No repositories found" : "No repositories available"}
                            </div>
                        ) : (
                            filteredRepos.map((repo) => (
                                <SelectItem key={repo.id} value={repo.full_name}>
                                    <div className="flex items-center gap-2">
                                        <IconBrandGithub className="h-4 w-4" />
                                        <span>{repo.full_name}</span>
                                        {repo.private && (
                                            <span className="text-xs text-muted-foreground">(Private)</span>
                                        )}
                                    </div>
                                </SelectItem>
                            ))
                        )}
                    </SelectContent>
                </Select>
                <p className="text-sm text-muted-foreground">
                    Connect your GitHub repository for automatic deployments
                </p>
            </div>

            {selectedRepo && (
                <div className="space-y-2">
                    <Label htmlFor="githubBranch">Branch</Label>
                    <Input
                        id="githubBranch"
                        name="githubBranch"
                        value={branch}
                        onChange={(e) => setBranch(e.target.value)}
                        placeholder="main"
                    />
                    <p className="text-sm text-muted-foreground">
                        Branch to deploy from (default: {repos.find(r => r.full_name === selectedRepo)?.default_branch || "main"})
                    </p>
                </div>
            )}
        </div>
    )
}

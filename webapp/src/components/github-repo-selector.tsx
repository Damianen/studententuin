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
    const [repos] = React.useState<GitHubRepo[]>([
        { id: 1, name: "my-app", full_name: "username/my-app", private: false, default_branch: "main" },
        { id: 2, name: "website", full_name: "username/website", private: false, default_branch: "main" },
        { id: 3, name: "api-project", full_name: "username/api-project", private: true, default_branch: "develop" }
    ])
    const [selectedRepo, setSelectedRepo] = React.useState<string>("")
    const [branch, setBranch] = React.useState<string>("main")
    const [searchQuery, setSearchQuery] = React.useState("")

    const filteredRepos = repos.filter((repo) =>
        repo.full_name.toLowerCase().includes(searchQuery.toLowerCase())
    )

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

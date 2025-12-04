"use client"

import * as React from "react"
import {
    IconCheck,
    IconX,
    IconClock,
    IconPlayerPlay,
    IconGitBranch,
    IconUser,
    IconCalendar,
    IconDots,
} from "@tabler/icons-react"
import { Button } from "@/components/ui/button"
import {
    Card,
    CardContent,
    CardDescription,
    CardHeader,
    CardTitle,
} from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"

type DeploymentStatus = "success" | "failed" | "in_progress" | "queued"

type Deployment = {
    id: string
    status: DeploymentStatus
    commit: string
    branch: string
    message: string
    author: string
    timestamp: string
    duration: string
}

const mockDeployments: Deployment[] = [
    {
        id: "1",
        status: "success",
        commit: "a3f2c9d",
        branch: "main",
        message: "Fix authentication bug in login flow",
        author: "John Doe",
        timestamp: "2 hours ago",
        duration: "2m 34s",
    },
    {
        id: "2",
        status: "success",
        commit: "b7e4f1a",
        branch: "main",
        message: "Add rate limiting to API endpoints",
        author: "Jane Smith",
        timestamp: "5 hours ago",
        duration: "3m 12s",
    },
    {
        id: "3",
        status: "in_progress",
        commit: "c9d2e5b",
        branch: "develop",
        message: "Update dependencies and refactor components",
        author: "John Doe",
        timestamp: "Just now",
        duration: "1m 45s",
    },
    {
        id: "4",
        status: "failed",
        commit: "d1f3a7c",
        branch: "feature/new-dashboard",
        message: "Implement new dashboard layout",
        author: "Alice Johnson",
        timestamp: "1 day ago",
        duration: "45s",
    },
    {
        id: "5",
        status: "success",
        commit: "e5c8b2d",
        branch: "main",
        message: "Optimize database queries for better performance",
        author: "Bob Wilson",
        timestamp: "2 days ago",
        duration: "4m 18s",
    },
    {
        id: "6",
        status: "success",
        commit: "f7a9d4e",
        branch: "main",
        message: "Add pagination to user list",
        author: "Jane Smith",
        timestamp: "3 days ago",
        duration: "2m 56s",
    },
]

const statusConfig: Record<
    DeploymentStatus,
    { label: string; icon: React.ElementType; color: string; bgColor: string }
> = {
    success: {
        label: "Success",
        icon: IconCheck,
        color: "text-green-600",
        bgColor: "bg-green-500/10",
    },
    failed: {
        label: "Failed",
        icon: IconX,
        color: "text-red-600",
        bgColor: "bg-red-500/10",
    },
    in_progress: {
        label: "In Progress",
        icon: IconClock,
        color: "text-blue-600",
        bgColor: "bg-blue-500/10",
    },
    queued: {
        label: "Queued",
        icon: IconPlayerPlay,
        color: "text-yellow-600",
        bgColor: "bg-yellow-500/10",
    },
}

export function DeploymentsList() {
    const [deployments] = React.useState<Deployment[]>(mockDeployments)

    return (
        <div className="space-y-4">
            <Card>
                <CardHeader>
                    <div className="flex items-center justify-between">
                        <div>
                            <CardTitle>Deployment History</CardTitle>
                            <CardDescription>
                                Recent deployments and their status
                            </CardDescription>
                        </div>
                        <Button>
                            <IconPlayerPlay className="h-4 w-4 mr-2" />
                            Deploy Now
                        </Button>
                    </div>
                </CardHeader>
                <CardContent className="p-0">
                    <div className="divide-y">
                        {deployments.map((deployment) => {
                            const config = statusConfig[deployment.status]
                            const StatusIcon = config.icon

                            return (
                                <div
                                    key={deployment.id}
                                    className="p-6 hover:bg-muted/50 transition-colors"
                                >
                                    <div className="flex items-start justify-between gap-4">
                                        <div className="flex items-start gap-4 flex-1">
                                            <div className={`p-2 rounded-full ${config.bgColor}`}>
                                                <StatusIcon className={`h-5 w-5 ${config.color}`} />
                                            </div>
                                            <div className="flex-1 space-y-3">
                                                <div className="space-y-1">
                                                    <div className="flex items-center gap-2 flex-wrap">
                                                        <Badge
                                                            variant="outline"
                                                            className={`${config.bgColor} ${config.color} border-0`}
                                                        >
                                                            {config.label}
                                                        </Badge>
                                                        <span className="text-sm text-muted-foreground">
                                                            {deployment.timestamp}
                                                        </span>
                                                        <span className="text-sm text-muted-foreground">
                                                            •
                                                        </span>
                                                        <span className="text-sm text-muted-foreground">
                                                            {deployment.duration}
                                                        </span>
                                                    </div>
                                                    <p className="font-medium">
                                                        {deployment.message}
                                                    </p>
                                                </div>
                                                <div className="flex items-center gap-4 text-sm text-muted-foreground flex-wrap">
                                                    <div className="flex items-center gap-1.5">
                                                        <IconGitBranch className="h-4 w-4" />
                                                        <span>{deployment.branch}</span>
                                                    </div>
                                                    <div className="flex items-center gap-1.5">
                                                        <span className="font-mono bg-muted px-2 py-0.5 rounded">
                                                            {deployment.commit}
                                                        </span>
                                                    </div>
                                                    <div className="flex items-center gap-1.5">
                                                        <IconUser className="h-4 w-4" />
                                                        <span>{deployment.author}</span>
                                                    </div>
                                                </div>
                                            </div>
                                        </div>
                                        <DropdownMenu>
                                            <DropdownMenuTrigger asChild>
                                                <Button variant="ghost" size="icon">
                                                    <IconDots className="h-4 w-4" />
                                                </Button>
                                            </DropdownMenuTrigger>
                                            <DropdownMenuContent align="end">
                                                <DropdownMenuItem>View Logs</DropdownMenuItem>
                                                <DropdownMenuItem>Redeploy</DropdownMenuItem>
                                                <DropdownMenuItem>Rollback</DropdownMenuItem>
                                            </DropdownMenuContent>
                                        </DropdownMenu>
                                    </div>
                                </div>
                            )
                        })}
                    </div>
                </CardContent>
            </Card>

            {/* Deployment Stats */}
            <div className="grid grid-cols-1 gap-4 md:grid-cols-4">
                <Card>
                    <CardHeader className="pb-2">
                        <CardDescription>Total Deployments</CardDescription>
                        <CardTitle className="text-3xl">127</CardTitle>
                    </CardHeader>
                    <CardContent>
                        <div className="text-xs text-muted-foreground">
                            Last 30 days
                        </div>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader className="pb-2">
                        <CardDescription>Success Rate</CardDescription>
                        <CardTitle className="text-3xl">94.5%</CardTitle>
                    </CardHeader>
                    <CardContent>
                        <div className="text-xs text-muted-foreground">
                            <Badge variant="outline" className="bg-green-500/10 text-green-600">
                                Excellent
                            </Badge>
                        </div>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader className="pb-2">
                        <CardDescription>Avg Deploy Time</CardDescription>
                        <CardTitle className="text-3xl">3m 24s</CardTitle>
                    </CardHeader>
                    <CardContent>
                        <div className="text-xs text-muted-foreground">
                            Last 30 days average
                        </div>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader className="pb-2">
                        <CardDescription>Last Deployment</CardDescription>
                        <CardTitle className="text-3xl">2h ago</CardTitle>
                    </CardHeader>
                    <CardContent>
                        <div className="text-xs text-muted-foreground">
                            <Badge variant="outline" className="bg-green-500/10 text-green-600">
                                Success
                            </Badge>
                        </div>
                    </CardContent>
                </Card>
            </div>
        </div>
    )
}

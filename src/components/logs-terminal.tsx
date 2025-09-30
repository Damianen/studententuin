"use client"

import * as React from "react"
import { IconDownload, IconSearch, IconX } from "@tabler/icons-react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import {
    Card,
    CardContent,
    CardDescription,
    CardHeader,
    CardTitle,
} from "@/components/ui/card"

type LogEntry = {
    id: string
    timestamp: string
    level: "info" | "warn" | "error" | "debug"
    message: string
}

const mockLogs: LogEntry[] = [
    { id: "1", timestamp: "2025-09-30 14:32:01", level: "info", message: "Server started on port 3000" },
    { id: "2", timestamp: "2025-09-30 14:32:15", level: "info", message: "Database connection established" },
    { id: "3", timestamp: "2025-09-30 14:33:22", level: "info", message: "GET /api/users 200 45ms" },
    { id: "4", timestamp: "2025-09-30 14:34:10", level: "warn", message: "Slow query detected: SELECT * FROM users took 1234ms" },
    { id: "5", timestamp: "2025-09-30 14:35:45", level: "info", message: "POST /api/auth/login 200 123ms" },
    { id: "6", timestamp: "2025-09-30 14:36:01", level: "error", message: "Failed to connect to Redis: ECONNREFUSED 127.0.0.1:6379" },
    { id: "7", timestamp: "2025-09-30 14:36:02", level: "info", message: "Retrying Redis connection..." },
    { id: "8", timestamp: "2025-09-30 14:36:05", level: "info", message: "Redis connection successful" },
    { id: "9", timestamp: "2025-09-30 14:37:30", level: "debug", message: "Cache hit for key: user:12345" },
    { id: "10", timestamp: "2025-09-30 14:38:12", level: "info", message: "GET /api/products 200 67ms" },
    { id: "11", timestamp: "2025-09-30 14:39:45", level: "warn", message: "Rate limit exceeded for IP: 192.168.1.1" },
    { id: "12", timestamp: "2025-09-30 14:40:22", level: "error", message: "Unhandled exception in /api/orders: TypeError: Cannot read property 'id' of undefined" },
    { id: "13", timestamp: "2025-09-30 14:41:05", level: "info", message: "Scheduled task: cleanup-old-sessions started" },
    { id: "14", timestamp: "2025-09-30 14:41:15", level: "info", message: "Deleted 143 expired sessions" },
    { id: "15", timestamp: "2025-09-30 14:42:30", level: "debug", message: "Websocket connection opened: client-abc123" },
]

const levelColors = {
    info: "text-blue-400",
    warn: "text-yellow-400",
    error: "text-red-400",
    debug: "text-gray-400",
}

export function LogsTerminal() {
    const [logs] = React.useState<LogEntry[]>(mockLogs)
    const [search, setSearch] = React.useState("")
    const [selectedLevels, setSelectedLevels] = React.useState<Set<string>>(
        new Set(["info", "warn", "error", "debug"])
    )
    const terminalRef = React.useRef<HTMLDivElement>(null)

    const filteredLogs = logs.filter((log) => {
        const matchesSearch = log.message.toLowerCase().includes(search.toLowerCase())
        const matchesLevel = selectedLevels.has(log.level)
        return matchesSearch && matchesLevel
    })

    const toggleLevel = (level: string) => {
        const newLevels = new Set(selectedLevels)
        if (newLevels.has(level)) {
            newLevels.delete(level)
        } else {
            newLevels.add(level)
        }
        setSelectedLevels(newLevels)
    }

    const downloadLogs = () => {
        const logsText = logs
            .map((log) => `[${log.timestamp}] ${log.level.toUpperCase()}: ${log.message}`)
            .join("\n")
        const blob = new Blob([logsText], { type: "text/plain" })
        const url = URL.createObjectURL(blob)
        const a = document.createElement("a")
        a.href = url
        a.download = `logs-${new Date().toISOString()}.txt`
        a.click()
        URL.revokeObjectURL(url)
    }

    return (
        <Card className="overflow-hidden">
            <CardHeader>
                <div className="flex items-center justify-between">
                    <div>
                        <CardTitle>Application Logs</CardTitle>
                        <CardDescription>
                            Real-time logs from your application
                        </CardDescription>
                    </div>
                    <Button onClick={downloadLogs} variant="outline" size="sm">
                        <IconDownload className="h-4 w-4 mr-2" />
                        Download
                    </Button>
                </div>
                <div className="flex items-center gap-2 pt-4">
                    <div className="relative flex-1">
                        <IconSearch className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
                        <Input
                            placeholder="Search logs..."
                            value={search}
                            onChange={(e) => setSearch(e.target.value)}
                            className="pl-9"
                        />
                        {search && (
                            <Button
                                variant="ghost"
                                size="icon"
                                className="absolute right-0 top-0 h-full"
                                onClick={() => setSearch("")}
                            >
                                <IconX className="h-4 w-4" />
                            </Button>
                        )}
                    </div>
                    <div className="flex gap-2">
                        {["info", "warn", "error", "debug"].map((level) => (
                            <Button
                                key={level}
                                variant={selectedLevels.has(level) ? "default" : "outline"}
                                size="sm"
                                onClick={() => toggleLevel(level)}
                                className="capitalize"
                            >
                                {level}
                            </Button>
                        ))}
                    </div>
                </div>
            </CardHeader>
            <CardContent className="p-0">
                <div
                    ref={terminalRef}
                    className="bg-black text-green-400 font-mono text-sm p-4 h-[600px] overflow-auto"
                    style={{
                        fontFamily: "ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace",
                    }}
                >
                    {filteredLogs.map((log) => (
                        <div key={log.id} className="mb-1 hover:bg-gray-900/50">
                            <span className="text-gray-500">{log.timestamp}</span>
                            {" "}
                            <span className={`font-bold ${levelColors[log.level]}`}>
                                [{log.level.toUpperCase()}]
                            </span>
                            {" "}
                            <span className="text-gray-300">{log.message}</span>
                        </div>
                    ))}
                    {filteredLogs.length === 0 && (
                        <div className="text-gray-500 text-center py-8">
                            No logs found matching your criteria
                        </div>
                    )}
                </div>
            </CardContent>
        </Card>
    )
}

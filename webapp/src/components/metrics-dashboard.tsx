"use client"

import * as React from "react"
import { Area, AreaChart, CartesianGrid, XAxis } from "recharts"
import {
    Card,
    CardContent,
    CardDescription,
    CardHeader,
    CardTitle,
} from "@/components/ui/card"
import {
    ChartConfig,
    ChartContainer,
    ChartTooltip,
    ChartTooltipContent,
} from "@/components/ui/chart"
import { Badge } from "@/components/ui/badge"

// Mock data for CPU usage
const cpuData = [
    { time: "00:00", value: 20 },
    { time: "04:00", value: 25 },
    { time: "08:00", value: 45 },
    { time: "12:00", value: 65 },
    { time: "16:00", value: 55 },
    { time: "20:00", value: 35 },
    { time: "23:59", value: 30 },
]

// Mock data for Memory usage
const memoryData = [
    { time: "00:00", value: 512 },
    { time: "04:00", value: 580 },
    { time: "08:00", value: 720 },
    { time: "12:00", value: 890 },
    { time: "16:00", value: 950 },
    { time: "20:00", value: 680 },
    { time: "23:59", value: 620 },
]

// Mock data for Response time
const responseTimeData = [
    { time: "00:00", value: 45 },
    { time: "04:00", value: 52 },
    { time: "08:00", value: 89 },
    { time: "12:00", value: 125 },
    { time: "16:00", value: 98 },
    { time: "20:00", value: 67 },
    { time: "23:59", value: 55 },
]

// Mock data for Request rate
const requestRateData = [
    { time: "00:00", value: 120 },
    { time: "04:00", value: 85 },
    { time: "08:00", value: 450 },
    { time: "12:00", value: 890 },
    { time: "16:00", value: 650 },
    { time: "20:00", value: 380 },
    { time: "23:59", value: 220 },
]

const cpuChartConfig = {
    value: {
        label: "CPU %",
        color: "hsl(var(--chart-1))",
    },
} satisfies ChartConfig

const memoryChartConfig = {
    value: {
        label: "Memory (MB)",
        color: "hsl(var(--chart-2))",
    },
} satisfies ChartConfig

const responseTimeChartConfig = {
    value: {
        label: "Response Time (ms)",
        color: "hsl(var(--chart-3))",
    },
} satisfies ChartConfig

const requestRateChartConfig = {
    value: {
        label: "Requests/min",
        color: "hsl(var(--chart-4))",
    },
} satisfies ChartConfig

export function MetricsDashboard() {
    return (
        <div className="space-y-4">
            {/* Summary Cards */}
            <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-4">
                <Card>
                    <CardHeader className="pb-2">
                        <CardDescription>Average CPU</CardDescription>
                        <CardTitle className="text-3xl">38%</CardTitle>
                    </CardHeader>
                    <CardContent>
                        <div className="text-xs text-muted-foreground">
                            <Badge variant="outline" className="bg-green-500/10 text-green-600">
                                Normal
                            </Badge>
                        </div>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader className="pb-2">
                        <CardDescription>Average Memory</CardDescription>
                        <CardTitle className="text-3xl">736 MB</CardTitle>
                    </CardHeader>
                    <CardContent>
                        <div className="text-xs text-muted-foreground">
                            of 2048 MB allocated
                        </div>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader className="pb-2">
                        <CardDescription>Avg Response Time</CardDescription>
                        <CardTitle className="text-3xl">79ms</CardTitle>
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
                        <CardDescription>Request Rate</CardDescription>
                        <CardTitle className="text-3xl">397/min</CardTitle>
                    </CardHeader>
                    <CardContent>
                        <div className="text-xs text-muted-foreground">
                            Last 24 hours average
                        </div>
                    </CardContent>
                </Card>
            </div>

            {/* Charts */}
            <div className="grid grid-cols-1 gap-4 lg:grid-cols-2">
                {/* CPU Usage Chart */}
                <Card>
                    <CardHeader>
                        <CardTitle>CPU Usage</CardTitle>
                        <CardDescription>Last 24 hours</CardDescription>
                    </CardHeader>
                    <CardContent>
                        <ChartContainer config={cpuChartConfig} className="h-[250px] w-full">
                            <AreaChart data={cpuData}>
                                <defs>
                                    <linearGradient id="fillCpu" x1="0" y1="0" x2="0" y2="1">
                                        <stop
                                            offset="5%"
                                            stopColor="var(--color-value)"
                                            stopOpacity={0.8}
                                        />
                                        <stop
                                            offset="95%"
                                            stopColor="var(--color-value)"
                                            stopOpacity={0.1}
                                        />
                                    </linearGradient>
                                </defs>
                                <CartesianGrid strokeDasharray="3 3" vertical={false} />
                                <XAxis
                                    dataKey="time"
                                    tickLine={false}
                                    axisLine={false}
                                    tickMargin={8}
                                />
                                <ChartTooltip content={<ChartTooltipContent />} />
                                <Area
                                    dataKey="value"
                                    type="monotone"
                                    fill="url(#fillCpu)"
                                    stroke="var(--color-value)"
                                    strokeWidth={2}
                                />
                            </AreaChart>
                        </ChartContainer>
                    </CardContent>
                </Card>

                {/* Memory Usage Chart */}
                <Card>
                    <CardHeader>
                        <CardTitle>Memory Usage</CardTitle>
                        <CardDescription>Last 24 hours</CardDescription>
                    </CardHeader>
                    <CardContent>
                        <ChartContainer config={memoryChartConfig} className="h-[250px] w-full">
                            <AreaChart data={memoryData}>
                                <defs>
                                    <linearGradient id="fillMemory" x1="0" y1="0" x2="0" y2="1">
                                        <stop
                                            offset="5%"
                                            stopColor="var(--color-value)"
                                            stopOpacity={0.8}
                                        />
                                        <stop
                                            offset="95%"
                                            stopColor="var(--color-value)"
                                            stopOpacity={0.1}
                                        />
                                    </linearGradient>
                                </defs>
                                <CartesianGrid strokeDasharray="3 3" vertical={false} />
                                <XAxis
                                    dataKey="time"
                                    tickLine={false}
                                    axisLine={false}
                                    tickMargin={8}
                                />
                                <ChartTooltip content={<ChartTooltipContent />} />
                                <Area
                                    dataKey="value"
                                    type="monotone"
                                    fill="url(#fillMemory)"
                                    stroke="var(--color-value)"
                                    strokeWidth={2}
                                />
                            </AreaChart>
                        </ChartContainer>
                    </CardContent>
                </Card>

                {/* Response Time Chart */}
                <Card>
                    <CardHeader>
                        <CardTitle>Response Time</CardTitle>
                        <CardDescription>Last 24 hours</CardDescription>
                    </CardHeader>
                    <CardContent>
                        <ChartContainer config={responseTimeChartConfig} className="h-[250px] w-full">
                            <AreaChart data={responseTimeData}>
                                <defs>
                                    <linearGradient id="fillResponseTime" x1="0" y1="0" x2="0" y2="1">
                                        <stop
                                            offset="5%"
                                            stopColor="var(--color-value)"
                                            stopOpacity={0.8}
                                        />
                                        <stop
                                            offset="95%"
                                            stopColor="var(--color-value)"
                                            stopOpacity={0.1}
                                        />
                                    </linearGradient>
                                </defs>
                                <CartesianGrid strokeDasharray="3 3" vertical={false} />
                                <XAxis
                                    dataKey="time"
                                    tickLine={false}
                                    axisLine={false}
                                    tickMargin={8}
                                />
                                <ChartTooltip content={<ChartTooltipContent />} />
                                <Area
                                    dataKey="value"
                                    type="monotone"
                                    fill="url(#fillResponseTime)"
                                    stroke="var(--color-value)"
                                    strokeWidth={2}
                                />
                            </AreaChart>
                        </ChartContainer>
                    </CardContent>
                </Card>

                {/* Request Rate Chart */}
                <Card>
                    <CardHeader>
                        <CardTitle>Request Rate</CardTitle>
                        <CardDescription>Last 24 hours</CardDescription>
                    </CardHeader>
                    <CardContent>
                        <ChartContainer config={requestRateChartConfig} className="h-[250px] w-full">
                            <AreaChart data={requestRateData}>
                                <defs>
                                    <linearGradient id="fillRequestRate" x1="0" y1="0" x2="0" y2="1">
                                        <stop
                                            offset="5%"
                                            stopColor="var(--color-value)"
                                            stopOpacity={0.8}
                                        />
                                        <stop
                                            offset="95%"
                                            stopColor="var(--color-value)"
                                            stopOpacity={0.1}
                                        />
                                    </linearGradient>
                                </defs>
                                <CartesianGrid strokeDasharray="3 3" vertical={false} />
                                <XAxis
                                    dataKey="time"
                                    tickLine={false}
                                    axisLine={false}
                                    tickMargin={8}
                                />
                                <ChartTooltip content={<ChartTooltipContent />} />
                                <Area
                                    dataKey="value"
                                    type="monotone"
                                    fill="url(#fillRequestRate)"
                                    stroke="var(--color-value)"
                                    strokeWidth={2}
                                />
                            </AreaChart>
                        </ChartContainer>
                    </CardContent>
                </Card>
            </div>

            {/* Additional Metrics */}
            <div className="grid grid-cols-1 gap-4 md:grid-cols-3">
                <Card>
                    <CardHeader>
                        <CardTitle>Error Rate</CardTitle>
                    </CardHeader>
                    <CardContent>
                        <div className="text-3xl font-bold">0.02%</div>
                        <p className="text-xs text-muted-foreground mt-2">
                            2 errors out of 10,000 requests
                        </p>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader>
                        <CardTitle>Network I/O</CardTitle>
                    </CardHeader>
                    <CardContent>
                        <div className="space-y-2">
                            <div>
                                <div className="text-sm text-muted-foreground">Inbound</div>
                                <div className="text-2xl font-bold">45.2 MB/s</div>
                            </div>
                            <div>
                                <div className="text-sm text-muted-foreground">Outbound</div>
                                <div className="text-2xl font-bold">38.7 MB/s</div>
                            </div>
                        </div>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader>
                        <CardTitle>Active Connections</CardTitle>
                    </CardHeader>
                    <CardContent>
                        <div className="text-3xl font-bold">234</div>
                        <p className="text-xs text-muted-foreground mt-2">
                            Current active connections
                        </p>
                    </CardContent>
                </Card>
            </div>
        </div>
    )
}

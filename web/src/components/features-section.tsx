"use client"

import * as React from "react"
import {
    IconRocket,
    IconServer,
    IconDatabase,
    IconChartBar,
    IconShieldCheck,
    IconClock,
    IconCode,
    IconUsers,
    IconWorldWww,
} from "@tabler/icons-react"
import {
    Card,
    CardContent,
    CardDescription,
    CardHeader,
    CardTitle,
} from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"

const features = [
    {
        icon: IconRocket,
        title: "Instant Deployment",
        description: "Deploy your applications in seconds with our streamlined deployment process. No complex configuration needed.",
        color: "text-blue-600",
        bgColor: "bg-blue-500/10",
    },
    {
        icon: IconServer,
        title: "Multiple Runtimes",
        description: "Support for Node.js, .NET, and more. Run your applications with the stack you know and love.",
        color: "text-purple-600",
        bgColor: "bg-purple-500/10",
    },
    {
        icon: IconDatabase,
        title: "Database Hosting",
        description: "MySQL and PostgreSQL databases with automatic backups and easy management.",
        color: "text-green-600",
        bgColor: "bg-green-500/10",
    },
    {
        icon: IconChartBar,
        title: "Real-time Metrics",
        description: "Monitor CPU, memory, response times, and request rates in real-time with beautiful dashboards.",
        color: "text-orange-600",
        bgColor: "bg-orange-500/10",
    },
    {
        icon: IconShieldCheck,
        title: "Secure by Default",
        description: "Built-in security features including SSL certificates, environment variables, and access controls.",
        color: "text-red-600",
        bgColor: "bg-red-500/10",
    },
    {
        icon: IconClock,
        title: "99.9% Uptime",
        description: "Reliable infrastructure with automatic failover and health monitoring to keep your apps running.",
        color: "text-cyan-600",
        bgColor: "bg-cyan-500/10",
    },
    {
        icon: IconCode,
        title: "Developer Friendly",
        description: "Git integration, environment variables, and comprehensive logs make development a breeze.",
        color: "text-indigo-600",
        bgColor: "bg-indigo-500/10",
    },
    {
        icon: IconUsers,
        title: "Student Focused",
        description: "Special pricing and resources for students learning to code and build their first projects.",
        color: "text-pink-600",
        bgColor: "bg-pink-500/10",
    },
    {
        icon: IconWorldWww,
        title: "Custom Domains",
        description: "Connect your own domain or use our subdomains for quick prototyping and testing.",
        color: "text-teal-600",
        bgColor: "bg-teal-500/10",
    },
]

export function FeaturesSection() {
    return (
        <div className="min-h-screen py-24">
            <div className="mx-auto max-w-7xl px-6">
                {/* Header */}
                <div className="text-center mb-16">
                    <Badge className="mb-4" variant="outline">
                        Features
                    </Badge>
                    <h1 className="text-4xl font-bold mb-4 md:text-5xl lg:text-6xl">
                        Everything You Need to Deploy
                    </h1>
                    <p className="text-lg text-muted-foreground max-w-2xl mx-auto">
                        Powerful hosting platform built for students and professionals.
                        Deploy applications and databases with confidence.
                    </p>
                </div>

                {/* Features Grid */}
                <div className="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
                    {features.map((feature, index) => {
                        const Icon = feature.icon
                        return (
                            <Card key={index} className="relative overflow-hidden">
                                <CardHeader>
                                    <div className={`inline-flex p-3 rounded-lg ${feature.bgColor} mb-4 w-fit`}>
                                        <Icon className={`h-6 w-6 ${feature.color}`} />
                                    </div>
                                    <CardTitle className="text-xl">{feature.title}</CardTitle>
                                </CardHeader>
                                <CardContent>
                                    <CardDescription className="text-base">
                                        {feature.description}
                                    </CardDescription>
                                </CardContent>
                            </Card>
                        )
                    })}
                </div>

                {/* CTA Section */}
                <div className="mt-20 text-center">
                    <Card className="bg-gradient-to-br from-primary/10 to-primary/5 border-primary/20">
                        <CardHeader>
                            <CardTitle className="text-3xl mb-4">
                                Ready to Get Started?
                            </CardTitle>
                            <CardDescription className="text-lg">
                                Join thousands of students and professionals already using Studententuin
                            </CardDescription>
                        </CardHeader>
                        <CardContent className="flex flex-col sm:flex-row gap-4 justify-center">
                            <a
                                href="/login"
                                className="inline-flex items-center justify-center rounded-md bg-primary px-8 py-3 text-sm font-medium text-primary-foreground hover:bg-primary/90 transition-colors"
                            >
                                Start Building
                            </a>
                            <a
                                href="/pricing"
                                className="inline-flex items-center justify-center rounded-md border border-input bg-background px-8 py-3 text-sm font-medium hover:bg-accent hover:text-accent-foreground transition-colors"
                            >
                                View Pricing
                            </a>
                        </CardContent>
                    </Card>
                </div>
            </div>
        </div>
    )
}

"use client"

import * as React from "react"
import { IconTarget, IconHeart, IconBulb } from "@tabler/icons-react"
import {
    Card,
    CardContent,
    CardDescription,
    CardHeader,
    CardTitle,
} from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"

export function AboutSection() {
    return (
        <div className="min-h-screen py-24">
            <div className="mx-auto max-w-7xl px-6">
                {/* Header */}
                <div className="text-center mb-16">
                    <Badge className="mb-4" variant="outline">
                        About Us
                    </Badge>
                    <h1 className="text-4xl font-bold mb-4 md:text-5xl lg:text-6xl">
                        Empowering the Next Generation
                    </h1>
                    <p className="text-lg text-muted-foreground max-w-2xl mx-auto">
                        We're on a mission to make hosting accessible for students and professionals alike.
                    </p>
                </div>

                {/* Story Section */}
                <div className="mb-24">
                    <Card className="border-none shadow-lg">
                        <CardContent className="p-8 md:p-12">
                            <h2 className="text-3xl font-bold mb-6">Our Story</h2>
                            <div className="space-y-4 text-muted-foreground">
                                <p className="text-lg">
                                    Studententuin (Dutch for "student garden") was born from a simple idea:
                                    hosting shouldn't be complicated or expensive for those learning to code.
                                </p>
                                <p className="text-lg">
                                    We remember struggling with complex deployment processes and expensive hosting
                                    bills while building our first projects. That's why we created a platform that's
                                    both powerful and accessible.
                                </p>
                                <p className="text-lg">
                                    Today, Studententuin helps thousands of students and professionals deploy their
                                    applications and databases with confidence. We're proud to be part of their journey
                                    from learning to building production-ready applications.
                                </p>
                            </div>
                        </CardContent>
                    </Card>
                </div>

                {/* Values Section */}
                <div className="mb-24">
                    <h2 className="text-3xl font-bold text-center mb-12">Our Values</h2>
                    <div className="grid grid-cols-1 gap-8 md:grid-cols-3">
                        <Card>
                            <CardHeader>
                                <div className="inline-flex p-3 rounded-lg bg-blue-500/10 mb-4 w-fit">
                                    <IconTarget className="h-6 w-6 text-blue-600" />
                                </div>
                                <CardTitle>Simplicity</CardTitle>
                            </CardHeader>
                            <CardContent>
                                <CardDescription className="text-base">
                                    We believe powerful tools should be easy to use. Our platform removes
                                    complexity so you can focus on building great applications.
                                </CardDescription>
                            </CardContent>
                        </Card>

                        <Card>
                            <CardHeader>
                                <div className="inline-flex p-3 rounded-lg bg-red-500/10 mb-4 w-fit">
                                    <IconHeart className="h-6 w-6 text-red-600" />
                                </div>
                                <CardTitle>Community</CardTitle>
                            </CardHeader>
                            <CardContent>
                                <CardDescription className="text-base">
                                    We're committed to supporting the developer community. That's why we
                                    offer free hosting for students and prioritize community feedback.
                                </CardDescription>
                            </CardContent>
                        </Card>

                        <Card>
                            <CardHeader>
                                <div className="inline-flex p-3 rounded-lg bg-yellow-500/10 mb-4 w-fit">
                                    <IconBulb className="h-6 w-6 text-yellow-600" />
                                </div>
                                <CardTitle>Innovation</CardTitle>
                            </CardHeader>
                            <CardContent>
                                <CardDescription className="text-base">
                                    We constantly improve our platform with the latest technologies and
                                    features to give you the best hosting experience.
                                </CardDescription>
                            </CardContent>
                        </Card>
                    </div>
                </div>

                {/* Stats Section */}
                <div className="mb-24">
                    <Card className="bg-gradient-to-br from-primary/10 to-primary/5 border-primary/20">
                        <CardContent className="p-12">
                            <div className="grid grid-cols-1 gap-8 md:grid-cols-3 text-center">
                                <div>
                                    <div className="text-4xl font-bold mb-2">10,000+</div>
                                    <div className="text-muted-foreground">Active Users</div>
                                </div>
                                <div>
                                    <div className="text-4xl font-bold mb-2">50,000+</div>
                                    <div className="text-muted-foreground">Deployments</div>
                                </div>
                                <div>
                                    <div className="text-4xl font-bold mb-2">99.9%</div>
                                    <div className="text-muted-foreground">Uptime</div>
                                </div>
                            </div>
                        </CardContent>
                    </Card>
                </div>

                {/* CTA Section */}
                <div className="text-center">
                    <h2 className="text-3xl font-bold mb-4">Join Us Today</h2>
                    <p className="text-lg text-muted-foreground mb-8 max-w-2xl mx-auto">
                        Be part of a growing community of developers building amazing things.
                    </p>
                    <div className="flex flex-col sm:flex-row gap-4 justify-center">
                        <a
                            href="/login"
                            className="inline-flex items-center justify-center rounded-md bg-primary px-8 py-3 text-sm font-medium text-primary-foreground hover:bg-primary/90 transition-colors"
                        >
                            Get Started
                        </a>
                        <a
                            href="/features"
                            className="inline-flex items-center justify-center rounded-md border border-input bg-background px-8 py-3 text-sm font-medium hover:bg-accent hover:text-accent-foreground transition-colors"
                        >
                            Learn More
                        </a>
                    </div>
                </div>
            </div>
        </div>
    )
}

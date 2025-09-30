"use client"

import * as React from "react"
import { IconCheck } from "@tabler/icons-react"
import {
    Card,
    CardContent,
    CardDescription,
    CardHeader,
    CardTitle,
} from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"

const pricingPlans = [
    {
        name: "Student",
        price: "Free",
        description: "Perfect for students learning to code",
        features: [
            "1 Application",
            "1 Database",
            "512 MB RAM",
            "1 GB Storage",
            "Community Support",
            "SSL Certificate",
        ],
        popular: false,
        cta: "Get Started",
    },
    {
        name: "Professional",
        price: "€9.99",
        period: "/month",
        description: "For professionals building serious projects",
        features: [
            "5 Applications",
            "5 Databases",
            "2 GB RAM per app",
            "10 GB Storage",
            "Priority Support",
            "SSL Certificate",
            "Custom Domains",
            "Advanced Metrics",
        ],
        popular: true,
        cta: "Start Free Trial",
    },
    {
        name: "Enterprise",
        price: "Custom",
        description: "For teams and large-scale applications",
        features: [
            "Unlimited Applications",
            "Unlimited Databases",
            "Custom RAM",
            "Custom Storage",
            "24/7 Dedicated Support",
            "SSL Certificate",
            "Custom Domains",
            "Advanced Metrics",
            "SLA Guarantee",
            "Priority Deployment",
        ],
        popular: false,
        cta: "Contact Sales",
    },
]

export function PricingSection() {
    return (
        <div className="min-h-screen py-24">
            <div className="mx-auto max-w-7xl px-6">
                {/* Header */}
                <div className="text-center mb-16">
                    <Badge className="mb-4" variant="outline">
                        Pricing
                    </Badge>
                    <h1 className="text-4xl font-bold mb-4 md:text-5xl lg:text-6xl">
                        Simple, Transparent Pricing
                    </h1>
                    <p className="text-lg text-muted-foreground max-w-2xl mx-auto">
                        Choose the plan that fits your needs. All plans include core features.
                        No hidden fees.
                    </p>
                </div>

                {/* Pricing Cards */}
                <div className="grid grid-cols-1 gap-8 lg:grid-cols-3">
                    {pricingPlans.map((plan, index) => (
                        <Card
                            key={index}
                            className={`relative flex flex-col ${
                                plan.popular
                                    ? "border-primary shadow-lg scale-105"
                                    : ""
                            }`}
                        >
                            {plan.popular && (
                                <div className="absolute -top-4 left-0 right-0 flex justify-center">
                                    <Badge className="bg-primary text-primary-foreground">
                                        Most Popular
                                    </Badge>
                                </div>
                            )}
                            <CardHeader className="pb-8">
                                <CardTitle className="text-2xl">{plan.name}</CardTitle>
                                <CardDescription className="text-base">
                                    {plan.description}
                                </CardDescription>
                                <div className="mt-4">
                                    <span className="text-4xl font-bold">
                                        {plan.price}
                                    </span>
                                    {plan.period && (
                                        <span className="text-muted-foreground">
                                            {plan.period}
                                        </span>
                                    )}
                                </div>
                            </CardHeader>
                            <CardContent className="flex-1 space-y-4">
                                <ul className="space-y-3">
                                    {plan.features.map((feature, featureIndex) => (
                                        <li
                                            key={featureIndex}
                                            className="flex items-start gap-3"
                                        >
                                            <IconCheck className="h-5 w-5 text-green-600 shrink-0 mt-0.5" />
                                            <span className="text-sm">{feature}</span>
                                        </li>
                                    ))}
                                </ul>
                                <Button
                                    className="w-full mt-8"
                                    variant={plan.popular ? "default" : "outline"}
                                >
                                    {plan.cta}
                                </Button>
                            </CardContent>
                        </Card>
                    ))}
                </div>

                {/* FAQ Section */}
                <div className="mt-24">
                    <h2 className="text-3xl font-bold text-center mb-12">
                        Frequently Asked Questions
                    </h2>
                    <div className="grid grid-cols-1 gap-6 md:grid-cols-2 max-w-4xl mx-auto">
                        <Card>
                            <CardHeader>
                                <CardTitle className="text-lg">
                                    Can I upgrade or downgrade anytime?
                                </CardTitle>
                            </CardHeader>
                            <CardContent>
                                <CardDescription>
                                    Yes, you can change your plan at any time. Changes take effect immediately.
                                </CardDescription>
                            </CardContent>
                        </Card>
                        <Card>
                            <CardHeader>
                                <CardTitle className="text-lg">
                                    Do you offer student discounts?
                                </CardTitle>
                            </CardHeader>
                            <CardContent>
                                <CardDescription>
                                    Yes! Our Student plan is completely free for verified students.
                                </CardDescription>
                            </CardContent>
                        </Card>
                        <Card>
                            <CardHeader>
                                <CardTitle className="text-lg">
                                    What payment methods do you accept?
                                </CardTitle>
                            </CardHeader>
                            <CardContent>
                                <CardDescription>
                                    We accept all major credit cards, PayPal, and bank transfers for Enterprise plans.
                                </CardDescription>
                            </CardContent>
                        </Card>
                        <Card>
                            <CardHeader>
                                <CardTitle className="text-lg">
                                    Is there a free trial?
                                </CardTitle>
                            </CardHeader>
                            <CardContent>
                                <CardDescription>
                                    Professional plans include a 14-day free trial. No credit card required.
                                </CardDescription>
                            </CardContent>
                        </Card>
                    </div>
                </div>
            </div>
        </div>
    )
}

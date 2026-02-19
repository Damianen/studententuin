'use client'

import { Logo } from '@/components/logo'
import { Button } from '@/components/ui/button'
import Link from 'next/link'
import { useAuth } from '@/contexts/auth_context'
import { Label } from '@/components/ui/label'
import { Input } from '@/components/ui/input'
import { useRouter } from 'next/navigation'
import { type FormEvent, useEffect, useState } from 'react'

export default function LoginPage() {
    const { login, isAuthenticated } = useAuth()
    const router = useRouter()
    const [email, setEmail] = useState('')
    const [password, setPassword] = useState('')
    const [error, setError] = useState<string | null>(null)
    const [submitting, setSubmitting] = useState(false)

    useEffect(() => {
        if (isAuthenticated) {
            router.replace('/projects')
        }
    }, [isAuthenticated, router])

    const onSubmit = async (event: FormEvent<HTMLFormElement>) => {
        event.preventDefault()
        setSubmitting(true)
        setError(null)
        try {
            await login({ email, password })
            router.push('/projects')
        } catch (submitError) {
            setError(submitError instanceof Error ? submitError.message : 'Login failed')
        } finally {
            setSubmitting(false)
        }
    }

    return (
        <section className="flex bg-zinc-50 px-4 py-16 md:py-32 dark:bg-transparent">
            <form
                onSubmit={onSubmit}
                className="max-w-92 m-auto h-fit w-full">
                <div className="p-6">
                    <div>
                        <Link
                            href="/"
                            aria-label="go home">
                            <Logo />
                        </Link>
                        <h1 className="mb-1 mt-4 text-xl font-semibold">Sign In to Studententuin</h1>
                        <p>Welcome back! Sign in to continue</p>
                    </div>

                    <div className="mt-6 space-y-4">
                        <div className="space-y-2">
                            <Label htmlFor="email">Email</Label>
                            <Input
                                id="email"
                                type="email"
                                value={email}
                                onChange={(event) => setEmail(event.target.value)}
                                required
                            />
                        </div>
                        <div className="space-y-2">
                            <Label htmlFor="password">Password</Label>
                            <Input
                                id="password"
                                type="password"
                                value={password}
                                onChange={(event) => setPassword(event.target.value)}
                                required
                            />
                        </div>
                        {error ? <p className="text-sm text-red-500">{error}</p> : null}
                        <Button type="submit" className="w-full" disabled={submitting}>
                            {submitting ? 'Signing in...' : 'Sign In'}
                        </Button>
                    </div>
                </div>
            </form>
        </section>
    )
}

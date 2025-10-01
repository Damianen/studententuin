'use client'

import { Logo } from '@/components/logo'
import { Button } from '@/components/ui/button'
import Link from 'next/link'
import { signIn } from 'next-auth/react'

export default function LoginPage() {
    return (
        <section className="flex bg-zinc-50 px-4 py-16 md:py-32 dark:bg-transparent">
            <form
                action=""
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

                    <div className="mt-6">
                        <Button
                            type="button"
                            variant="outline"
                            className="w-full"
                            onClick={() => signIn('github')}>
                            <svg
                                xmlns="http://www.w3.org/2000/svg"
                                width="1em"
                                height="1em"
                                viewBox="0 0 24 24"
                                fill="currentColor"
                                aria-hidden="true">
                                <path d="M12.026 2c-5.509 0-9.974 4.465-9.974 9.974c0 4.406 2.857 8.145 6.821 9.465c.499.09.679-.217.679-.48c0-.237-.009-.868-.014-1.703c-2.775.602-3.361-1.338-3.361-1.338c-.454-1.154-1.11-1.462-1.11-1.462c-.908-.62.069-.608.069-.608c1.003.07 1.531 1.03 1.531 1.03c.892 1.529 2.341 1.087 2.91.832c.091-.647.35-1.087.637-1.337c-2.216-.252-4.549-1.108-4.549-4.932c0-1.089.39-1.982 1.029-2.68c-.103-.253-.447-1.27.098-2.647c0 0 .84-.269 2.75 1.025c.798-.222 1.654-.333 2.506-.337c.851.004 1.707.115 2.506.337c1.909-1.294 2.748-1.025 2.748-1.025c.547 1.377.203 2.394.1 2.647c.64.698 1.028 1.591 1.028 2.68c0 3.834-2.337 4.677-4.561 4.925c.359.309.678.919.678 1.852c0 1.336-.012 2.415-.012 2.741c0 .266.178.576.687.478c3.961-1.322 6.816-5.06 6.816-9.463C22 6.465 17.535 2 12.026 2" />
                            </svg>
                            <span>GitHub</span>
                        </Button>
                    </div>
                </div>
            </form>
        </section>
    )
}

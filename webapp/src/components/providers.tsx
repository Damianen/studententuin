"use client";

import { ThemeProvider } from "@/components/theme-provider";
import { HeroHeader } from "@/components/header";
import FooterSection from "@/components/footer";
import { usePathname } from "next/navigation";

export function Providers({ children }: { children: React.ReactNode }) {
    const pathname = usePathname();
    const isDashboard = pathname?.startsWith("/dashboard");

    return (
        <ThemeProvider>
            {isDashboard ? (
                children
            ) : (
                <div className="min-h-screen flex flex-col justify-between">
                    <HeroHeader />
                    {children}
                    <FooterSection />
                </div>
            )}
        </ThemeProvider>
    );
}

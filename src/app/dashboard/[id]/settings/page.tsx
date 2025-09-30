import { auth } from "@/modules/auth/infrastructure/auth";
import { prisma } from "@/lib/prisma";
import { redirect, notFound } from "next/navigation";
import { AppSidebar } from "@/components/app-sidebar";
import { SiteHeader } from "@/components/site-header";
import { EnvironmentVariables } from "@/components/environment-variables";
import {
    SidebarInset,
    SidebarProvider,
} from "@/components/ui/sidebar";

export default async function SettingsPage({
    params,
}: {
    params: Promise<{ id: string }>;
}) {
    const session = await auth();

    if (!session?.user?.id) {
        redirect("/login");
    }

    const { id } = await params;

    // Try to find as web application first
    let application = await prisma.webApplication.findUnique({
        where: { id },
    });

    // If not found, try to find as database
    let database = null;
    if (!application) {
        database = await prisma.database.findUnique({
            where: { id },
        });
    }

    // If neither found, return 404
    if (!application && !database) {
        notFound();
    }

    // Check ownership
    const resource = application || database;
    if (resource!.userId !== session.user.id) {
        redirect("/projects");
    }

    const name = resource!.name;

    return (
        <SidebarProvider
            style={
                {
                    "--sidebar-width": "calc(var(--spacing) * 72)",
                    "--header-height": "calc(var(--spacing) * 12)",
                } as React.CSSProperties
            }
        >
            <AppSidebar id={id} variant="inset" />
            <SidebarInset>
                <SiteHeader />
                <div className="flex flex-1 flex-col">
                    <div className="@container/main flex flex-1 flex-col gap-2">
                        <div className="flex flex-col gap-4 py-4 md:gap-6 md:py-6">
                            <div className="px-4 lg:px-6">
                                <h1 className="text-3xl font-bold mb-2">Settings</h1>
                                <p className="text-muted-foreground">
                                    Configure settings for {name}
                                </p>
                            </div>
                            <div className="px-4 lg:px-6">
                                <EnvironmentVariables />
                            </div>
                        </div>
                    </div>
                </div>
            </SidebarInset>
        </SidebarProvider>
    );
}

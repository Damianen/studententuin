import { auth } from "@/modules/auth/infrastructure/auth";
import { redirect, notFound } from "next/navigation";
import { AppSidebar } from "@/components/app-sidebar";
import { ChartAreaInteractive } from "@/components/chart-area-interactive";
import { SectionCards } from "@/components/section-cards";
import { SiteHeader } from "@/components/site-header";
import {
    SidebarInset,
    SidebarProvider,
} from "@/components/ui/sidebar";

import data from "../../dashboard-old/data.json";

export default async function DashboardPage({
    params,
}: {
    params: Promise<{ id: string }>;
}) {
    const session = await auth();

    if (!session?.user?.id) {
        redirect("/login");
    }

    const { id } = await params;

    // TODO: Replace with API call to fetch application/database details
    const mockData: Record<string, any> = {
        "app-1": {
            id: "app-1",
            name: "My Application",
            subdomain: "myapp",
            runtime: "nodejs",
            userId: session.user.id,
            isApplication: true,
        },
        "db-1": {
            id: "db-1",
            name: "My Database",
            subdomain: "mydb",
            type: "postgresql",
            userId: session.user.id,
            isApplication: false,
        },
    };

    const resource = mockData[id];
    if (!resource) {
        notFound();
    }

    const isApplication = resource.isApplication;
    const name = resource.name;
    const subdomain = resource.subdomain;
    const type = isApplication
        ? resource.runtime === "nodejs"
            ? "Node.js"
            : ".NET"
        : resource.type === "mysql"
            ? "MySQL"
            : "PostgreSQL";

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
                                <h1 className="text-3xl font-bold mb-2">{name}</h1>
                                <p className="text-muted-foreground">
                                    {subdomain}.studententuin • {type} • {isApplication ? "Application" : "Database"}
                                </p>
                            </div>
                            <SectionCards />
                            <div className="px-4 lg:px-6">
                                <ChartAreaInteractive />
                            </div>
                        </div>
                    </div>
                </div>
            </SidebarInset>
        </SidebarProvider>
    );
}

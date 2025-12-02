import { auth } from "@/modules/auth/infrastructure/auth";
import { redirect, notFound } from "next/navigation";
import { AppSidebar } from "@/components/app-sidebar";
import { SiteHeader } from "@/components/site-header";
import { DeploymentsList } from "@/components/deployments-list";
import {
    SidebarInset,
    SidebarProvider,
} from "@/components/ui/sidebar";

export default async function DeploymentsPage({
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
        "app-1": { id: "app-1", name: "My Application", userId: session.user.id, isApplication: true },
        "db-1": { id: "db-1", name: "My Database", userId: session.user.id, isApplication: false },
    };

    const resource = mockData[id];
    if (!resource) {
        notFound();
    }

    const isApplication = resource.isApplication;
    const name = resource.name;

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
                                <h1 className="text-3xl font-bold mb-2">Deployments</h1>
                                <p className="text-muted-foreground">
                                    View and manage deployments for {name}
                                </p>
                            </div>
                            <div className="px-4 lg:px-6">
                                <DeploymentsList />
                            </div>
                        </div>
                    </div>
                </div>
            </SidebarInset>
        </SidebarProvider>
    );
}

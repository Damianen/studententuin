"use client";

import * as React from "react";
import type { CSSProperties } from "react";
import { AppSidebar } from "@/components/app-sidebar";
import { ChartAreaInteractive } from "@/components/chart-area-interactive";
import { SectionCards } from "@/components/section-cards";
import { SiteHeader } from "@/components/site-header";
import { SidebarInset, SidebarProvider } from "@/components/ui/sidebar";
import { useDashboardResource } from "@/hooks/use-dashboard-resource";

export default function DashboardPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = React.use(params);
  const { resource, loading, error } = useDashboardResource(id);

  return (
    <SidebarProvider
      style={
        {
          "--sidebar-width": "calc(var(--spacing) * 72)",
          "--header-height": "calc(var(--spacing) * 12)",
        } as CSSProperties
      }
    >
      <AppSidebar id={id} variant="inset" />
      <SidebarInset>
        <SiteHeader />
        <div className="flex flex-1 flex-col">
          <div className="@container/main flex flex-1 flex-col gap-2">
            <div className="flex flex-col gap-4 py-4 md:gap-6 md:py-6">
              <div className="px-4 lg:px-6">
                <h1 className="text-3xl font-bold mb-2">
                  {resource?.resourceName || "Loading..."}
                </h1>
                <p className="text-muted-foreground">
                  {loading
                    ? "Loading resource details..."
                    : resource
                      ? `${resource.fullDomain} • ${resource.typeLabel} • ${resource.kind === "application" ? "Application" : "Database"}`
                      : "Resource unavailable"}
                </p>
                {error ? <p className="text-sm text-red-500 mt-2">{error}</p> : null}
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

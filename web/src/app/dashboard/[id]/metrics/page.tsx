"use client";

import * as React from "react";
import type { CSSProperties } from "react";
import { AppSidebar } from "@/components/app-sidebar";
import { SiteHeader } from "@/components/site-header";
import { MetricsDashboard } from "@/components/metrics-dashboard";
import { SidebarInset, SidebarProvider } from "@/components/ui/sidebar";
import { useDashboardResource } from "@/hooks/use-dashboard-resource";

export default function MetricsPage({
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
                <h1 className="text-3xl font-bold mb-2">Metrics</h1>
                <p className="text-muted-foreground">
                  {loading
                    ? "Loading resource details..."
                    : `Monitor performance metrics for ${resource?.resourceName || "your project"}`}
                </p>
                {error ? <p className="text-sm text-red-500 mt-2">{error}</p> : null}
              </div>
              <div className="px-4 lg:px-6">
                <MetricsDashboard />
              </div>
            </div>
          </div>
        </div>
      </SidebarInset>
    </SidebarProvider>
  );
}

"use client"

import * as React from "react"
import { useRouter } from "next/navigation"

import { NavDocuments } from "@/components/nav-documents"
import { NavMain } from "@/components/nav-main"
import { NavSecondary } from "@/components/nav-secondary"
import { NavUser } from "@/components/nav-user"
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar"
import { useAuth } from "@/contexts/auth_context"
import SubdomainController from "@/controllers/subdomain_controller"
import { mapSubdomainResources } from "@/lib/subdomain_resources"

const data = {
  user: {
    name: "User",
    email: "",
    avatar: "",
  },
  navMain: (id: string) => [
    {
      title: "Overview",
      url: `/dashboard/${id}`,
      iconType: "overview" as const,
    },
    {
      title: "Deployments",
      url: `/dashboard/${id}/deployments`,
      iconType: "deployments" as const,
    },
    {
      title: "Metrics",
      url: `/dashboard/${id}/metrics`,
      iconType: "metrics" as const,
    },
    {
      title: "Logs",
      url: `/dashboard/${id}/logs`,
      iconType: "logs" as const,
    },
    {
      title: "Settings",
      url: `/dashboard/${id}/settings`,
      iconType: "settings" as const,
    },
  ],
  navSecondary: [
    {
      title: "Documentation",
      url: "#",
      iconType: "documentation" as const,
    },
    {
      title: "Support",
      url: "#",
      iconType: "support" as const,
    },
  ],
}

export function AppSidebar({ id, ...props }: React.ComponentProps<typeof Sidebar> & { id?: string }) {
  const router = useRouter()
  const { user, isAuthenticated, logout } = useAuth()
  const [documents, setDocuments] = React.useState<Array<{
    id?: string
    name: string
    url: string
    iconType: "folder" | "world" | "database"
  }>>([
    {
      name: "All Projects",
      url: "/projects",
      iconType: "folder",
    },
  ])

  React.useEffect(() => {
    if (!isAuthenticated) {
      return
    }

    SubdomainController.getAll()
      .then((subdomains) => {
        const resourceItems = mapSubdomainResources(subdomains).map((resource) => ({
          id: resource.resourceId,
          name: resource.resourceName,
          url: `/dashboard/${resource.resourceId}`,
          iconType: resource.kind === "application" ? "world" as const : "database" as const,
        }))

        setDocuments([
          {
            name: "All Projects",
            url: "/projects",
            iconType: "folder",
          },
          ...resourceItems,
        ])
      })
      .catch(() => {
        setDocuments([
          {
            name: "All Projects",
            url: "/projects",
            iconType: "folder",
          },
        ])
      })
  }, [isAuthenticated])

  const currentUser = user
    ? {
      name: user.name || "User",
      email: user.email || "",
      avatar: "",
    }
    : data.user

  const onLogout = async () => {
    await logout()
    router.push("/login")
  }

  return (
    <Sidebar collapsible="offcanvas" {...props}>
      <SidebarHeader>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton
              asChild
              className="data-[slot=sidebar-menu-button]:!p-1.5"
            >
              <a href="/projects" className="flex items-center gap-2">
                <svg
                  viewBox="0 0 24 24"
                  fill="none"
                  xmlns="http://www.w3.org/2000/svg"
                  className="size-5">
                  <path
                    d="M12 2C10.5 4 9 6.5 9 9C9 9.5 9.1 10 9.3 10.5C8.5 10.2 7.8 10 7 10C5.3 10 4 11 4 12.5C4 13.8 4.9 14.8 6 15.2C5.4 15.4 5 15.9 5 16.5C5 17.3 5.7 18 6.5 18H11V22H13V18H17.5C18.3 18 19 17.3 19 16.5C19 15.9 18.6 15.4 18 15.2C19.1 14.8 20 13.8 20 12.5C20 11 18.7 10 17 10C16.2 10 15.5 10.2 14.7 10.5C14.9 10 15 9.5 15 9C15 6.5 13.5 4 12 2Z"
                    fill="#22c55e"
                  />
                </svg>
                <span className="text-base font-semibold">Studententuin</span>
              </a>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>
      <SidebarContent>
        <NavMain items={id ? data.navMain(id) : []} />
        <NavDocuments items={documents} currentId={id} />
        <NavSecondary items={data.navSecondary} className="mt-auto" />
      </SidebarContent>
      <SidebarFooter>
        <NavUser user={currentUser} onLogout={onLogout} />
      </SidebarFooter>
    </Sidebar>
  )
}

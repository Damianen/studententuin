import * as React from "react"
import {
  IconCamera,
  IconChartBar,
  IconDashboard,
  IconDatabase,
  IconFileAi,
  IconFileDescription,
  IconFileWord,
  IconFolder,
  IconHelp,
  IconInnerShadowTop,
  IconListDetails,
  IconReport,
  IconSearch,
  IconSettings,
  IconUsers,
  IconWorld,
} from "@tabler/icons-react"

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
import { auth } from "@/modules/auth/infrastructure/auth"

const data = {
  user: {
    name: "shadcn",
    email: "m@example.com",
    avatar: "/avatars/shadcn.jpg",
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
  navClouds: [
    {
      title: "Databases",
      icon: IconDatabase,
      isActive: true,
      url: "#",
      items: [
        {
          title: "MySQL",
          url: "#",
        },
        {
          title: "PostgreSQL",
          url: "#",
        },
      ],
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
  documents: [
    {
      name: "All Projects",
      url: "/projects",
      icon: IconFolder,
    },
  ],
}

export async function AppSidebar({ id, ...props }: React.ComponentProps<typeof Sidebar> & { id?: string }) {
  const session = await auth()

  const currentUser = session?.user
    ? {
        name: session.user.name || "User",
        email: session.user.email || "",
        avatar: session.user.image || "/avatars/default.jpg",
      }
    : data.user

  // Fetch user's applications and databases
  let documents: Array<{
    id?: string
    name: string
    url: string
    iconType: "folder" | "world" | "database"
  }> = [
    {
      name: "All Projects",
      url: "/projects",
      iconType: "folder" as const,
    },
  ]

  if (session?.user?.id) {
    // TODO: Replace with API call to fetch user's applications and databases
    const applications = [
      { id: "app-1", name: "My App", runtime: "nodejs" },
    ]
    const databases = [
      { id: "db-1", name: "My Database", type: "postgresql" },
    ]

    const userResources = [
      ...applications.map((app) => ({
        id: app.id,
        name: app.name,
        url: `/dashboard/${app.id}`,
        iconType: "world" as const,
      })),
      ...databases.map((db) => ({
        id: db.id,
        name: db.name,
        url: `/dashboard/${db.id}`,
        iconType: "database" as const,
      })),
    ]

    documents = [...documents, ...userResources]
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
        <NavUser user={currentUser} />
      </SidebarFooter>
    </Sidebar>
  )
}

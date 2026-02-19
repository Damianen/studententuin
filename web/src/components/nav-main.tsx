"use client"

import {
  IconDashboard,
  IconListDetails,
  IconChartBar,
  IconFileDescription,
  IconSettings,
} from "@tabler/icons-react"
import { usePathname } from "next/navigation"

import {
  SidebarGroup,
  SidebarGroupContent,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar"

const iconMap = {
  overview: IconDashboard,
  deployments: IconListDetails,
  metrics: IconChartBar,
  logs: IconFileDescription,
  settings: IconSettings,
}

export function NavMain({
  items,
}: {
  items: {
    title: string
    url: string
    iconType?: "overview" | "deployments" | "metrics" | "logs" | "settings"
  }[]
}) {
  const pathname = usePathname()

  return (
    <SidebarGroup>
      <SidebarGroupContent className="flex flex-col gap-2">
        <SidebarMenu>
          {items.map((item) => {
            const isActive = pathname === item.url ||
              Boolean(item.title === "Overview" && pathname && pathname.match(/^\/dashboard\/[^\/]+$/))
            const Icon = item.iconType ? iconMap[item.iconType] : null

            return (
              <SidebarMenuItem key={item.title}>
                <SidebarMenuButton tooltip={item.title} isActive={isActive} asChild>
                  <a href={item.url}>
                    {Icon && <Icon />}
                    <span>{item.title}</span>
                  </a>
                </SidebarMenuButton>
              </SidebarMenuItem>
            )
          })}
        </SidebarMenu>
      </SidebarGroupContent>
    </SidebarGroup>
  )
}

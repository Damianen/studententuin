"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/contexts/auth_context";
import SubdomainController from "@/controllers/subdomain_controller";
import {
	findResourceById,
	type ResourceSummary,
} from "@/lib/subdomain_resources";

export function useDashboardResource(resourceId: string) {
	const router = useRouter();
	const { isAuthenticated, loading } = useAuth();
	const [resource, setResource] = useState<ResourceSummary | null>(null);
	const [hasFetched, setHasFetched] = useState(false);
	const [error, setError] = useState<string | null>(null);

	useEffect(() => {
		if (loading) {
			return;
		}

		if (!isAuthenticated) {
			router.replace("/login");
			return;
		}

		SubdomainController.getAll()
			.then((subdomains) => {
				const found = findResourceById(subdomains, resourceId);
				if (!found) {
					router.replace("/projects");
					setHasFetched(true);
					return;
				}
				setResource(found);
				setError(null);
				setHasFetched(true);
			})
			.catch((fetchError) => {
				setError(
					fetchError instanceof Error
						? fetchError.message
						: "Failed to load resource",
				);
				setHasFetched(true);
			});
	}, [isAuthenticated, loading, resourceId, router]);

	return {
		resource,
		loading: loading || !hasFetched,
		error,
	};
}

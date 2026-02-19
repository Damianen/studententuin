"use client";

import { useEffect, useMemo, useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { useAuth } from "@/contexts/auth_context";
import ApplicationController from "@/controllers/application_controller";
import DatabaseController from "@/controllers/database_controller";
import SubdomainController from "@/controllers/subdomain_controller";
import { mapSubdomainResources } from "@/lib/subdomain_resources";

type ResourceCard = {
  resourceId: string;
  subdomainId: string;
  name: string;
  fullDomain: string;
  typeLabel: string;
  kind: "application" | "database";
};

export default function ProjectsPage() {
  const router = useRouter();
  const { isAuthenticated, loading } = useAuth();
  const [resources, setResources] = useState<ResourceCard[]>([]);
  const [hasFetched, setHasFetched] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [editingResourceId, setEditingResourceId] = useState<string | null>(null);
  const [editName, setEditName] = useState("");
  const [pendingKey, setPendingKey] = useState<string | null>(null);
  const [actionError, setActionError] = useState<string | null>(null);

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
        setResources(
          mapSubdomainResources(subdomains).map((resource) => ({
            resourceId: resource.resourceId,
            subdomainId: resource.subdomainId,
            name: resource.resourceName,
            fullDomain: resource.fullDomain,
            typeLabel: resource.typeLabel,
            kind: resource.kind,
          })),
        );
        setError(null);
        setHasFetched(true);
      })
      .catch((fetchError) => {
        setError(
          fetchError instanceof Error
            ? fetchError.message
            : "Failed to load projects",
        );
        setHasFetched(true);
      });
  }, [isAuthenticated, loading, router]);

  function startEditing(resource: ResourceCard) {
    setEditingResourceId(resource.resourceId);
    setEditName(resource.name);
    setActionError(null);
  }

  function cancelEditing() {
    setEditingResourceId(null);
    setEditName("");
  }

  async function saveResource(resource: ResourceCard) {
    const trimmedName = editName.trim();
    if (!trimmedName) {
      setActionError("Name cannot be empty");
      return;
    }

    setPendingKey(`save-${resource.resourceId}`);
    setActionError(null);
    try {
      if (resource.kind === "application") {
        await ApplicationController.patch(resource.subdomainId, resource.resourceId, {
          name: trimmedName,
        });
      } else {
        await DatabaseController.patch(resource.subdomainId, resource.resourceId, {
          name: trimmedName,
        });
      }

      setResources((current) =>
        current.map((item) =>
          item.resourceId === resource.resourceId
            ? {
                ...item,
                name: trimmedName,
              }
            : item,
        ),
      );
      cancelEditing();
    } catch (updateError) {
      setActionError(
        updateError instanceof Error ? updateError.message : "Failed to update resource",
      );
    } finally {
      setPendingKey(null);
    }
  }

  async function deleteResource(resource: ResourceCard) {
    const confirmed = window.confirm(
      `Delete ${resource.kind} \"${resource.name}\"? This action cannot be undone.`,
    );
    if (!confirmed) {
      return;
    }

    setPendingKey(`delete-${resource.resourceId}`);
    setActionError(null);
    try {
      if (resource.kind === "application") {
        await ApplicationController.delete(resource.subdomainId, resource.resourceId);
      } else {
        await DatabaseController.delete(resource.subdomainId, resource.resourceId);
      }

      setResources((current) =>
        current.filter((item) => item.resourceId !== resource.resourceId),
      );
      if (editingResourceId === resource.resourceId) {
        cancelEditing();
      }
    } catch (deleteError) {
      setActionError(
        deleteError instanceof Error ? deleteError.message : "Failed to delete resource",
      );
    } finally {
      setPendingKey(null);
    }
  }

  const applications = useMemo(
    () => resources.filter((resource) => resource.kind === "application"),
    [resources],
  );
  const databases = useMemo(
    () => resources.filter((resource) => resource.kind === "database"),
    [resources],
  );

  if (loading || !hasFetched) {
    return <div className="container mx-auto py-10 pt-24">Loading projects...</div>;
  }

  return (
    <div className="container mx-auto py-10 pt-24">
      <div className="mb-8">
        <h1 className="text-4xl font-bold mb-2">Projects</h1>
        <p className="text-muted-foreground">
          Manage your web applications and databases
        </p>
      </div>

      {error ? <p className="text-sm text-red-500 mb-6">{error}</p> : null}
      {actionError ? <p className="text-sm text-red-500 mb-6">{actionError}</p> : null}

      <div className="space-y-8">
        <section>
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-2xl font-semibold">Web Applications</h2>
            <Button asChild>
              <Link href="/projects/new/app">New Application</Link>
            </Button>
          </div>

          {applications.length === 0 ? (
            <Card>
              <CardContent className="pt-6">
                <p className="text-muted-foreground text-center py-8">
                  No web applications yet. Create your first one!
                </p>
              </CardContent>
            </Card>
          ) : (
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
              {applications.map((app) => (
                <Card key={app.resourceId}>
                  <CardHeader>
                    <CardTitle>{app.name}</CardTitle>
                    <CardDescription>{app.typeLabel}</CardDescription>
                  </CardHeader>
                  <CardContent>
                    <p className="text-sm text-muted-foreground mb-3">
                      {app.fullDomain}
                    </p>

                    {editingResourceId === app.resourceId ? (
                      <div className="space-y-3">
                        <Input
                          value={editName}
                          onChange={(event) => setEditName(event.target.value)}
                          placeholder="Application name"
                        />
                        <div className="flex flex-wrap gap-2">
                          <Button
                            size="sm"
                            onClick={() => saveResource(app)}
                            disabled={pendingKey !== null}
                          >
                            Save
                          </Button>
                          <Button
                            size="sm"
                            variant="outline"
                            onClick={cancelEditing}
                            disabled={pendingKey !== null}
                          >
                            Cancel
                          </Button>
                        </div>
                      </div>
                    ) : (
                      <div className="flex flex-wrap gap-2">
                        <Button asChild variant="outline" size="sm">
                          <Link href={`/dashboard/${app.resourceId}`}>View Details</Link>
                        </Button>
                        <Button
                          variant="secondary"
                          size="sm"
                          onClick={() => startEditing(app)}
                          disabled={pendingKey !== null}
                        >
                          Edit
                        </Button>
                        <Button
                          variant="destructive"
                          size="sm"
                          onClick={() => deleteResource(app)}
                          disabled={pendingKey !== null}
                        >
                          Delete
                        </Button>
                      </div>
                    )}
                  </CardContent>
                </Card>
              ))}
            </div>
          )}
        </section>

        <section>
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-2xl font-semibold">Databases</h2>
            <Button asChild>
              <Link href="/projects/new/database">New Database</Link>
            </Button>
          </div>

          {databases.length === 0 ? (
            <Card>
              <CardContent className="pt-6">
                <p className="text-muted-foreground text-center py-8">
                  No databases yet. Create your first one!
                </p>
              </CardContent>
            </Card>
          ) : (
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
              {databases.map((db) => (
                <Card key={db.resourceId}>
                  <CardHeader>
                    <CardTitle>{db.name}</CardTitle>
                    <CardDescription>{db.typeLabel}</CardDescription>
                  </CardHeader>
                  <CardContent>
                    <p className="text-sm text-muted-foreground mb-3">
                      {db.fullDomain}
                    </p>

                    {editingResourceId === db.resourceId ? (
                      <div className="space-y-3">
                        <Input
                          value={editName}
                          onChange={(event) => setEditName(event.target.value)}
                          placeholder="Database name"
                        />
                        <div className="flex flex-wrap gap-2">
                          <Button
                            size="sm"
                            onClick={() => saveResource(db)}
                            disabled={pendingKey !== null}
                          >
                            Save
                          </Button>
                          <Button
                            size="sm"
                            variant="outline"
                            onClick={cancelEditing}
                            disabled={pendingKey !== null}
                          >
                            Cancel
                          </Button>
                        </div>
                      </div>
                    ) : (
                      <div className="flex flex-wrap gap-2">
                        <Button asChild variant="outline" size="sm">
                          <Link href={`/dashboard/${db.resourceId}`}>View Details</Link>
                        </Button>
                        <Button
                          variant="secondary"
                          size="sm"
                          onClick={() => startEditing(db)}
                          disabled={pendingKey !== null}
                        >
                          Edit
                        </Button>
                        <Button
                          variant="destructive"
                          size="sm"
                          onClick={() => deleteResource(db)}
                          disabled={pendingKey !== null}
                        >
                          Delete
                        </Button>
                      </div>
                    )}
                  </CardContent>
                </Card>
              ))}
            </div>
          )}
        </section>
      </div>
    </div>
  );
}

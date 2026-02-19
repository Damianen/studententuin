"use client";

import { useEffect, useState } from "react";
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
import { Label } from "@/components/ui/label";
import { useAuth } from "@/contexts/auth_context";
import ApplicationController from "@/controllers/application_controller";
import { ensureSubdomain } from "@/lib/subdomain_lookup";

export default function NewApplicationPage() {
  const router = useRouter();
  const { isAuthenticated, loading } = useAuth();
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!loading && !isAuthenticated) {
      router.replace("/login");
    }
  }, [isAuthenticated, loading, router]);

  async function createApplication(formData: FormData) {
    const name = String(formData.get("name") || "").trim();
    const subdomain = String(formData.get("subdomain") || "").trim();
    const type = String(formData.get("type") || "").trim();
    const repoUrl = String(formData.get("repo_url") || "").trim();
    const branch = String(formData.get("branch") || "main").trim();
    const buildCommand = String(formData.get("build_command") || "").trim();
    const startCommand = String(formData.get("start_command") || "").trim();

    setSubmitting(true);
    setError(null);

    try {
      const subdomainId = await ensureSubdomain(subdomain);
      await ApplicationController.create(subdomainId, {
        name,
        type,
        repo_url: repoUrl,
        branch,
        build_command: buildCommand,
        start_command: startCommand,
      });
      router.push("/projects");
    } catch (createError) {
      setError(
        createError instanceof Error
          ? createError.message
          : "Failed to create application",
      );
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <div className="container mx-auto py-10 pt-24">
      <div className="mb-8">
        <h1 className="text-4xl font-bold mb-2">New Web Application</h1>
        <p className="text-muted-foreground">Create a new web application</p>
      </div>

      <Card className="max-w-2xl">
        <CardHeader>
          <CardTitle>Application Details</CardTitle>
          <CardDescription>Enter the details for your new web application</CardDescription>
        </CardHeader>
        <CardContent>
          <form
            action={createApplication}
            className="space-y-6"
          >
            <div className="space-y-2">
              <Label htmlFor="name">Name</Label>
              <Input id="name" name="name" placeholder="My Application" required />
            </div>

            <div className="space-y-2">
              <Label htmlFor="subdomain">Subdomain</Label>
              <div className="flex items-center gap-2">
                <Input
                  id="subdomain"
                  name="subdomain"
                  placeholder="myapp"
                  required
                  pattern="^[a-z0-9-]+$"
                  title="Only lowercase letters, numbers, and hyphens allowed"
                  className="flex-1"
                />
                <span className="text-muted-foreground">.studententuin.com</span>
              </div>
            </div>

            <div className="space-y-2">
              <Label htmlFor="type">Runtime</Label>
              <select
                id="type"
                name="type"
                required
                className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
              >
                <option value="">Select a runtime...</option>
                <option value="Nodejs">Node.js</option>
                <option value="dotnet" disabled>
                  .NET (coming soon)
                </option>
              </select>
            </div>

            <div className="space-y-2">
              <Label htmlFor="repo_url">Repository URL</Label>
              <Input id="repo_url" name="repo_url" placeholder="https://github.com/org/repo" required />
            </div>

            <div className="space-y-2">
              <Label htmlFor="branch">Branch</Label>
              <Input id="branch" name="branch" placeholder="main" defaultValue="main" required />
            </div>

            <div className="space-y-2">
              <Label htmlFor="build_command">Build Command</Label>
              <Input id="build_command" name="build_command" placeholder="npm run build" required />
            </div>

            <div className="space-y-2">
              <Label htmlFor="start_command">Start Command</Label>
              <Input id="start_command" name="start_command" placeholder="npm run start" required />
            </div>

            {error ? <p className="text-sm text-red-500">{error}</p> : null}

            <div className="flex gap-4">
              <Button type="submit" disabled={submitting}>
                {submitting ? "Creating..." : "Create Application"}
              </Button>
              <Button type="button" variant="outline" asChild>
                <a href="/projects">Cancel</a>
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}

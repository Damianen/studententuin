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
import DatabaseController from "@/controllers/database_controller";
import { ensureSubdomain } from "@/lib/subdomain_lookup";

export default function NewDatabasePage() {
  const router = useRouter();
  const { isAuthenticated, loading } = useAuth();
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!loading && !isAuthenticated) {
      router.replace("/login");
    }
  }, [isAuthenticated, loading, router]);

  async function createDatabase(formData: FormData) {
    const name = String(formData.get("name") || "").trim();
    const subdomain = String(formData.get("subdomain") || "").trim();
    const type = String(formData.get("type") || "").trim();
    const version = String(formData.get("version") || "").trim();
    const dbName = String(formData.get("db_name") || "").trim();
    const dbPassword = String(formData.get("db_password") || "").trim();

    setSubmitting(true);
    setError(null);

    try {
      const subdomainId = await ensureSubdomain(subdomain);
      await DatabaseController.create(subdomainId, {
        name,
        type,
        version,
        db_name: dbName,
        db_password: dbPassword,
      });
      router.push("/projects");
    } catch (createError) {
      setError(
        createError instanceof Error
          ? createError.message
          : "Failed to create database",
      );
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <div className="container mx-auto py-10 pt-24">
      <div className="mb-8">
        <h1 className="text-4xl font-bold mb-2">New Database</h1>
        <p className="text-muted-foreground">Create a new database instance</p>
      </div>

      <Card className="max-w-2xl">
        <CardHeader>
          <CardTitle>Database Details</CardTitle>
          <CardDescription>Enter the details for your new database</CardDescription>
        </CardHeader>
        <CardContent>
          <form action={createDatabase} className="space-y-6">
            <div className="space-y-2">
              <Label htmlFor="name">Name</Label>
              <Input id="name" name="name" placeholder="My Database" required />
            </div>

            <div className="space-y-2">
              <Label htmlFor="subdomain">Subdomain</Label>
              <div className="flex items-center gap-2">
                <Input
                  id="subdomain"
                  name="subdomain"
                  placeholder="mydb"
                  required
                  pattern="^[a-z0-9-]+$"
                  title="Only lowercase letters, numbers, and hyphens allowed"
                  className="flex-1"
                />
                <span className="text-muted-foreground">.studententuin.com</span>
              </div>
            </div>

            <div className="space-y-2">
              <Label htmlFor="type">Database Type</Label>
              <select
                id="type"
                name="type"
                required
                className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
              >
                <option value="">Select a database type...</option>
                <option value="mysql">MySQL</option>
                <option value="postgres">PostgreSQL</option>
              </select>
            </div>

            <div className="space-y-2">
              <Label htmlFor="version">Version</Label>
              <Input id="version" name="version" placeholder="8.0" required />
            </div>

            <div className="space-y-2">
              <Label htmlFor="db_name">Database Name</Label>
              <Input id="db_name" name="db_name" placeholder="app_db" required />
            </div>

            <div className="space-y-2">
              <Label htmlFor="db_password">Database Password</Label>
              <Input id="db_password" name="db_password" type="password" required />
            </div>

            {error ? <p className="text-sm text-red-500">{error}</p> : null}

            <div className="flex gap-4">
              <Button type="submit" disabled={submitting}>
                {submitting ? "Creating..." : "Create Database"}
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

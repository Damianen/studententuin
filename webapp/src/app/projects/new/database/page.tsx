import { redirect } from "next/navigation";
import { revalidatePath } from "next/cache";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

export default async function NewDatabasePage() {
  async function createDatabase(formData: FormData) {
    "use server";

    const name = formData.get("name") as string;
    const subdomain = formData.get("subdomain") as string;
    const type = formData.get("type") as string;

    // Mock creation - frontend only
    console.log("Creating database:", { name, subdomain, type, userId: "mock-user-id" });

    revalidatePath("/projects");
    redirect("/projects");
  }

  return (
    <div className="container mx-auto py-10 pt-24">
      <div className="mb-8">
        <h1 className="text-4xl font-bold mb-2">New Database</h1>
        <p className="text-muted-foreground">
          Create a new database instance
        </p>
      </div>

      <Card className="max-w-2xl">
        <CardHeader>
          <CardTitle>Database Details</CardTitle>
          <CardDescription>
            Enter the details for your new database
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form action={createDatabase} className="space-y-6">
            <div className="space-y-2">
              <Label htmlFor="name">Name</Label>
              <Input
                id="name"
                name="name"
                placeholder="My Database"
                required
              />
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
                <span className="text-muted-foreground">.studententuin</span>
              </div>
              <p className="text-sm text-muted-foreground">
                Your database will be available at subdomain.studententuin
              </p>
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

            <div className="flex gap-4">
              <Button type="submit">Create Database</Button>
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

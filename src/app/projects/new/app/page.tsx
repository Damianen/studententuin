import { auth } from "@/modules/auth/infrastructure/auth";
import { redirect } from "next/navigation";
import { prisma } from "@/lib/prisma";
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
import { GitHubRepoSelector } from "@/components/github-repo-selector";

export default async function NewApplicationPage() {
  const session = await auth();

  if (!session?.user?.id) {
    redirect("/login");
  }

  async function createApplication(formData: FormData) {
    "use server";

    const session = await auth();
    if (!session?.user?.id) {
      redirect("/login");
    }

    const name = formData.get("name") as string;
    const subdomain = formData.get("subdomain") as string;
    const runtime = formData.get("runtime") as string;
    const githubRepo = formData.get("githubRepo") as string | null;
    const githubBranch = formData.get("githubBranch") as string | null;

    await prisma.webApplication.create({
      data: {
        name,
        subdomain,
        runtime,
        githubRepo: githubRepo || null,
        githubBranch: githubBranch || "main",
        userId: session.user.id,
      },
    });

    revalidatePath("/projects");
    redirect("/projects");
  }

  return (
    <div className="container mx-auto py-10 pt-24">
      <div className="mb-8">
        <h1 className="text-4xl font-bold mb-2">New Web Application</h1>
        <p className="text-muted-foreground">
          Create a new web application
        </p>
      </div>

      <Card className="max-w-2xl">
        <CardHeader>
          <CardTitle>Application Details</CardTitle>
          <CardDescription>
            Enter the details for your new web application
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form action={createApplication} className="space-y-6">
            <div className="space-y-2">
              <Label htmlFor="name">Name</Label>
              <Input
                id="name"
                name="name"
                placeholder="My Application"
                required
              />
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
                <span className="text-muted-foreground">.studententuin</span>
              </div>
              <p className="text-sm text-muted-foreground">
                Your app will be available at subdomain.studententuin
              </p>
            </div>

            <div className="space-y-2">
              <Label htmlFor="runtime">Runtime</Label>
              <select
                id="runtime"
                name="runtime"
                required
                className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
              >
                <option value="">Select a runtime...</option>
                <option value="nodejs">Node.js</option>
                <option value="dotnet">.NET</option>
              </select>
            </div>

            <GitHubRepoSelector />

            <div className="flex gap-4">
              <Button type="submit">Create Application</Button>
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

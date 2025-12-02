import { auth } from "@/modules/auth/infrastructure/auth";
import { redirect } from "next/navigation";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import Link from "next/link";

export default async function ProjectsPage() {
  const session = await auth();

  if (!session?.user?.id) {
    redirect("/login");
  }

  // TODO: Replace with API call to fetch user's applications and databases
  const webApplications = [
    {
      id: "app-1",
      name: "My Application",
      runtime: "nodejs",
      subdomain: "myapp"
    },
  ];
  const databases = [
    {
      id: "db-1",
      name: "My Database",
      type: "postgresql",
      subdomain: "mydb"
    },
  ];

  return (
    <div className="container mx-auto py-10 pt-24">
      <div className="mb-8">
        <h1 className="text-4xl font-bold mb-2">Projects</h1>
        <p className="text-muted-foreground">
          Manage your web applications and databases
        </p>
      </div>

      <div className="space-y-8">
        {/* Web Applications Section */}
        <section>
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-2xl font-semibold">Web Applications</h2>
            <Button asChild>
              <Link href="/projects/new/app">New Application</Link>
            </Button>
          </div>

          {webApplications.length === 0 ? (
            <Card>
              <CardContent className="pt-6">
                <p className="text-muted-foreground text-center py-8">
                  No web applications yet. Create your first one!
                </p>
              </CardContent>
            </Card>
          ) : (
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
              {webApplications.map((app) => (
                <Card key={app.id}>
                  <CardHeader>
                    <CardTitle>{app.name}</CardTitle>
                    <CardDescription>
                      {app.runtime === 'nodejs' ? 'Node.js' : '.NET'}
                    </CardDescription>
                  </CardHeader>
                  <CardContent>
                    <p className="text-sm text-muted-foreground mb-4">
                      {app.subdomain}.studententuin
                    </p>
                    <Button asChild variant="outline" size="sm">
                      <Link href={`/dashboard/${app.id}`}>View Details</Link>
                    </Button>
                  </CardContent>
                </Card>
              ))}
            </div>
          )}
        </section>

        {/* Databases Section */}
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
                <Card key={db.id}>
                  <CardHeader>
                    <CardTitle>{db.name}</CardTitle>
                    <CardDescription>
                      {db.type === 'mysql' ? 'MySQL' : 'PostgreSQL'}
                    </CardDescription>
                  </CardHeader>
                  <CardContent>
                    <p className="text-sm text-muted-foreground mb-4">
                      {db.subdomain}.studententuin
                    </p>
                    <Button asChild variant="outline" size="sm">
                      <Link href={`/dashboard/${db.id}`}>View Details</Link>
                    </Button>
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

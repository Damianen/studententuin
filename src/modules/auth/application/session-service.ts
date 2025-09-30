import type { AuthSession } from "../domain/session";
import type { AuthUser } from "../domain/user";
import { auth } from "../infrastructure/auth";

function toDomainUser(session: Awaited<ReturnType<typeof auth>>): AuthUser | null {
  if (!session?.user?.id) {
    return null;
  }

  return {
    id: session.user.id,
    name: session.user.name ?? null,
    email: session.user.email ?? null,
    image: session.user.image ?? null,
  };
}

export async function getServerAuthSession(): Promise<AuthSession | null> {
  const rawSession = await auth();
  const user = toDomainUser(rawSession);

  if (!rawSession || !user) {
    return null;
  }

  return {
    user,
  };
}

export async function getAuthenticatedUser(): Promise<AuthUser | null> {
  const session = await getServerAuthSession();

  return session?.user ?? null;
}

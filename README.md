## Environment

Create a `.env` file with the required secrets:

```bash
DATABASE_URL="postgresql://USER:PASSWORD@HOST:PORT/DATABASE"
NEXTAUTH_SECRET="generate-a-strong-random-string"
GITHUB_ID="your-oauth-client-id"
GITHUB_SECRET="your-oauth-client-secret"
```

## Database & Prisma

```bash
pnpm dlx prisma generate
pnpm dlx prisma migrate dev --name init
```

The Prisma schema lives in `prisma/schema.prisma`. Generated clients use the connection string from `DATABASE_URL`.

## Auth Flow

- OAuth is powered by GitHub via NextAuth v5, backed by a Prisma adapter.
- Auth modules follow a layered structure under `src/modules/auth` separating domain, application, and infrastructure concerns.
- The `/login` page triggers `signIn('github')` and the API route at `/api/auth/[...nextauth]` handles callbacks.

## Running Locally

```bash
pnpm install
pnpm dev
```

Then visit [http://localhost:3000](http://localhost:3000).

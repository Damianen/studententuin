export type GithubOAuthConfig = {
  clientId: string;
  clientSecret: string;
};

export function getGithubOAuthConfig(): GithubOAuthConfig {
  const clientId = process.env.GITHUB_ID;
  const clientSecret = process.env.GITHUB_SECRET;

  if (!clientId || !clientSecret) {
    throw new Error("Missing GitHub OAuth environment variables. Set GITHUB_ID and GITHUB_SECRET.");
  }

  return { clientId, clientSecret };
}

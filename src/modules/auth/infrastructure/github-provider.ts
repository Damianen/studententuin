import GitHubProvider from "next-auth/providers/github";

import { getGithubOAuthConfig } from "../config/github";

export function githubProvider() {
  const { clientId, clientSecret } = getGithubOAuthConfig();

  return GitHubProvider({
    clientId,
    clientSecret,
    authorization: {
      params: {
        // Scopes explained:
        // - read:user: Read user profile data
        // - user:email: Access user email addresses
        // - repo: Full control of private repositories (includes reading public repos)
        // This gives us access to list repos, read contents, and pull code
        scope: "read:user user:email repo",
      },
    },
  });
}

# node-hello — deploy e2e fixture

A minimal Node HTTP server (no dependencies) used by the deploy-from-git
end-to-end test (`web/e2e/deploy.spec.ts`). It logs a line on startup and per
request, so the test can assert real container logs after a deploy.

The servermanager only clones public `https://github.com/...` repositories,
so this fixture must live in its own public repo. Publish it once:

```bash
cd e2e-fixtures/node-hello
git init && git add . && git commit -m "stt deploy fixture"
gh repo create Damianen/stt-sample-node --public --source=. --push
```

Then run the e2e test against it (locally or in CI):

```bash
docker compose up -d --build --wait
cd web
E2E_DEPLOY_REPO=https://github.com/Damianen/stt-sample-node npx playwright test e2e/deploy.spec.ts
```

Without `E2E_DEPLOY_REPO` the deploy spec skips itself.

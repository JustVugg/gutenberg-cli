# github — OpenClaw skill

> github should become a local, searchable operational workspace, not only an endpoint wrapper.

- **CLI:** `github`
- **MCP:** `github mcp`
- **Base URL:** https://api.github.com
- **Operations:** 1183 (616 read, 373 write, 194 destructive)

## Actions (zero-friction)
- `github meta` — GitHub API Root
- `github advisories` — List global security advisories
- `github tasks` — List tasks
- `github app` — Get the authenticated app
- `github config` — Get a webhook configuration for an app
- `github deliveries` — List deliveries for an app webhook
- `github installation-requests` — List installation requests for the authenticated app
- `github installations` — List installations for the authenticated app
- `github classrooms` — List classrooms
- `github codes-of-conduct` — Get all codes of conduct
- `github emojis` — Get emojis
- `github events` — List public events
- `github feeds` — Get feeds
- `github gists` — List gists for the authenticated user
- `github gists-2` — List public gists
- `github starred` — List starred gists
- `github templates` — Get all gitignore templates
- `github repositories` — List repositories accessible to the app installation
- `github issues` — List issues assigned to the authenticated user
- `github licenses` — Get all commonly used licenses

## All operations (first 200)
- `github call meta/root` (GET /) — read-list — GitHub API Root
- `github call security-advisories/list-global-advisories` (GET /advisories) — read-list — List global security advisories
- `github call security-advisories/get-global-advisory` (GET /advisories/{ghsa_id}) — read-detail — Get a global security advisory
- `github call agent-tasks/list-tasks-for-repo` (GET /agents/repos/{owner}/{repo}/tasks) — read — List tasks for repository
- `github call agent-tasks/create-task-in-repo` (POST /agents/repos/{owner}/{repo}/tasks) — create — Start a task
- `github call agent-tasks/get-task-by-repo-and-id` (GET /agents/repos/{owner}/{repo}/tasks/{task_id}) — read-detail — Get a task by repo
- `github call agent-tasks/list-tasks` (GET /agents/tasks) — read-list — List tasks
- `github call agent-tasks/get-task-by-id` (GET /agents/tasks/{task_id}) — read-detail — Get a task by ID
- `github call apps/get-authenticated` (GET /app) — read-list — Get the authenticated app
- `github call apps/create-from-manifest` (POST /app-manifests/{code}/conversions) — create — Create a GitHub App from a manifest
- `github call apps/get-webhook-config-for-app` (GET /app/hook/config) — read-list — Get a webhook configuration for an app
- `github call apps/update-webhook-config-for-app` (PATCH /app/hook/config) — upsert — Update a webhook configuration for an app
- `github call apps/list-webhook-deliveries` (GET /app/hook/deliveries) — read-list — List deliveries for an app webhook
- `github call apps/get-webhook-delivery` (GET /app/hook/deliveries/{delivery_id}) — read-detail — Get a delivery for an app webhook
- `github call apps/redeliver-webhook-delivery` (POST /app/hook/deliveries/{delivery_id}/attempts) — create — Redeliver a delivery for an app webhook
- `github call apps/list-installation-requests-for-authenticated-app` (GET /app/installation-requests) — read-list — List installation requests for the authenticated app
- `github call apps/list-installations` (GET /app/installations) — read-list — List installations for the authenticated app
- `github call apps/get-installation` (GET /app/installations/{installation_id}) — read-detail — Get an installation for the authenticated app
- `github call apps/delete-installation` (DELETE /app/installations/{installation_id}) — delete — Delete an installation for the authenticated app
- `github call apps/create-installation-access-token` (POST /app/installations/{installation_id}/access_tokens) — action — Create an installation access token for an app
- `github call apps/suspend-installation` (PUT /app/installations/{installation_id}/suspended) — update — Suspend an app installation
- `github call apps/unsuspend-installation` (DELETE /app/installations/{installation_id}/suspended) — delete — Unsuspend an app installation
- `github call apps/delete-authorization` (DELETE /applications/{client_id}/grant) — delete — Delete an app authorization
- `github call apps/check-token` (POST /applications/{client_id}/token) — action — Check a token
- `github call apps/reset-token` (PATCH /applications/{client_id}/token) — update — Reset a token
- `github call apps/delete-token` (DELETE /applications/{client_id}/token) — delete — Delete an app token
- `github call apps/scope-token` (POST /applications/{client_id}/token/scoped) — create — Create a scoped access token
- `github call apps/get-by-slug` (GET /apps/{app_slug}) — read-detail — Get an app
- `github call classroom/get-an-assignment` (GET /assignments/{assignment_id}) — read-detail — Get an assignment
- `github call classroom/list-accepted-assignments-for-an-assignment` (GET /assignments/{assignment_id}/accepted_assignments) — read — List accepted assignments for an assignment
- `github call classroom/get-assignment-grades` (GET /assignments/{assignment_id}/grades) — read — Get assignment grades
- `github call classroom/list-classrooms` (GET /classrooms) — read-list — List classrooms
- `github call classroom/get-a-classroom` (GET /classrooms/{classroom_id}) — read-detail — Get a classroom
- `github call classroom/list-assignments-for-a-classroom` (GET /classrooms/{classroom_id}/assignments) — read — List assignments for a classroom
- `github call codes-of-conduct/get-all-codes-of-conduct` (GET /codes_of_conduct) — read-list — Get all codes of conduct
- `github call codes-of-conduct/get-conduct-code` (GET /codes_of_conduct/{key}) — read-detail — Get a code of conduct
- `github call credentials/revoke` (POST /credentials/revoke) — create — Revoke a list of credentials
- `github call emojis/get` (GET /emojis) — read-list — Get emojis
- `github call actions/get-actions-cache-retention-limit-for-enterprise` (GET /enterprises/{enterprise}/actions/cache/retention-limit) — read — Get GitHub Actions cache retention limit for an enterprise
- `github call actions/set-actions-cache-retention-limit-for-enterprise` (PUT /enterprises/{enterprise}/actions/cache/retention-limit) — update — Set GitHub Actions cache retention limit for an enterprise
- `github call actions/get-actions-cache-storage-limit-for-enterprise` (GET /enterprises/{enterprise}/actions/cache/storage-limit) — read — Get GitHub Actions cache storage limit for an enterprise
- `github call actions/set-actions-cache-storage-limit-for-enterprise` (PUT /enterprises/{enterprise}/actions/cache/storage-limit) — update — Set GitHub Actions cache storage limit for an enterprise
- `github call oidc/list-oidc-custom-property-inclusions-for-enterprise` (GET /enterprises/{enterprise}/actions/oidc/customization/properties/repo) — read — List OIDC custom property inclusions for an enterprise
- `github call oidc/create-oidc-custom-property-inclusion-for-enterprise` (POST /enterprises/{enterprise}/actions/oidc/customization/properties/repo) — create — Create an OIDC custom property inclusion for an enterprise
- `github call oidc/delete-oidc-custom-property-inclusion-for-enterprise` (DELETE /enterprises/{enterprise}/actions/oidc/customization/properties/repo/{custom_property_name}) — delete — Delete an OIDC custom property inclusion for an enterprise
- `github call code-security/get-configurations-for-enterprise` (GET /enterprises/{enterprise}/code-security/configurations) — read — Get code security configurations for an enterprise
- `github call code-security/create-configuration-for-enterprise` (POST /enterprises/{enterprise}/code-security/configurations) — create — Create a code security configuration for an enterprise
- `github call code-security/get-default-configurations-for-enterprise` (GET /enterprises/{enterprise}/code-security/configurations/defaults) — read — Get default code security configurations for an enterprise
- `github call code-security/get-single-configuration-for-enterprise` (GET /enterprises/{enterprise}/code-security/configurations/{configuration_id}) — read-detail — Retrieve a code security configuration of an enterprise
- `github call code-security/update-enterprise-configuration` (PATCH /enterprises/{enterprise}/code-security/configurations/{configuration_id}) — update — Update a custom code security configuration for an enterprise
- … +150 more — run `github operations --json` for the full list

## Safety
Write/destructive operations dry-run by default. Append `--yes` only on explicit user confirmation.
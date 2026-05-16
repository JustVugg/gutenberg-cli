# Security Model

Gutenberg treats generated tools as operational software, not demos.

## Defaults

- API keys are read from environment variables.
- OAuth tokens are stored in `.gutenberg/<tool>-token.json` with `0600` permissions.
- Generated examples do not commit secrets.
- Headers that look like tokens are redacted in dry-run output.
- Write and destructive operations dry-run by default.
- The MCP server exposes the same guarded runtime as the CLI.

## Required Before Production

- Review provider terms of service.
- Add service-specific rate limits.
- Add audit logs for write operations.
- Add OAuth support when API keys are insufficient.
- Add integration tests against sandbox APIs.
- Store secrets in a real secret manager for hosted deployments.

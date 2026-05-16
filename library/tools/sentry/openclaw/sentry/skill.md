# sentry — OpenClaw skill

> sentry should become a local, searchable operational workspace, not only an endpoint wrapper.

- **CLI:** `sentry`
- **MCP:** `sentry mcp`
- **Base URL:** https://{region}.sentry.io
- **Operations:** 209 (112 read, 63 write, 34 destructive)

## Auth setup
```bash
export SENTRY_API_KEY=<your-key>
```

## Actions (zero-friction)
- `sentry organizations` — Return a list of organizations available to the authenticated session in a region.
This is particularly useful for requests with a user bound context. For API key-based requests this will only return the organization that belongs to the key.
- `sentry models` — Get list of actively used LLM model names from Seer.

Returns the list of AI models that are currently used in production in Seer.
This endpoint does not require authentication and can be used to discover which models Seer uses.

Requests to this endpoint should use the region-specific domain
eg. `us.sentry.io` or `de.sentry.io`

## All operations (first 200)
- `sentry call List Your Organizations` (GET /api/0/organizations/) — read-list — Return a list of organizations available to the authenticated session in a region.
This is particularly useful for requests with a user bound context. For API key-based requests this will only return the organization that belongs to the key.
- `sentry call Retrieve an Organization` (GET /api/0/organizations/{organization_id_or_slug}/) — read — Return details on an individual organization, including various details
such as membership access and teams.
- `sentry call Update an Organization` (PUT /api/0/organizations/{organization_id_or_slug}/) — update — Update various attributes and configurable settings for the given organization.
- `sentry call Get Integration Provider Information` (GET /api/0/organizations/{organization_id_or_slug}/config/integrations/) — read — Get integration provider information about all available integrations for an organization.
- `sentry call List an Organization's Custom Dashboards` (GET /api/0/organizations/{organization_id_or_slug}/dashboards/) — read — Retrieve a list of custom dashboards that are associated with the given organization.
- `sentry call Create a New Dashboard for an Organization` (POST /api/0/organizations/{organization_id_or_slug}/dashboards/) — create — Create a new dashboard for the given Organization
- `sentry call Retrieve an Organization's Custom Dashboard` (GET /api/0/organizations/{organization_id_or_slug}/dashboards/{dashboard_id}/) — read — Return details about an organization's custom dashboard.
- `sentry call Edit an Organization's Custom Dashboard` (PUT /api/0/organizations/{organization_id_or_slug}/dashboards/{dashboard_id}/) — update — Edit an organization's custom dashboard as well as any bulk
edits on widgets that may have been made. (For example, widgets
that have been rearranged, updated queries and fields, specific
display types, and so on.)
- `sentry call Delete an Organization's Custom Dashboard` (DELETE /api/0/organizations/{organization_id_or_slug}/dashboards/{dashboard_id}/) — delete — Delete an organization's custom dashboard.
- `sentry call Fetch an Organization's Monitors` (GET /api/0/organizations/{organization_id_or_slug}/detectors/) — read — ⚠️ This endpoint is currently in **beta** and may be subject to change. It is supported by [New Monitors and Alerts](/product/new-monitors-and-alerts/) and may not be viewable in the UI today.

List an Organization's Monitors
- `sentry call Mutate an Organization's Monitors` (PUT /api/0/organizations/{organization_id_or_slug}/detectors/) — update — ⚠️ This endpoint is currently in **beta** and may be subject to change. It is supported by [New Monitors and Alerts](/product/new-monitors-and-alerts/) and may not be viewable in the UI today.

Bulk enable or disable an Organization's Monitors
- `sentry call Bulk Delete Monitors` (DELETE /api/0/organizations/{organization_id_or_slug}/detectors/) — delete — ⚠️ This endpoint is currently in **beta** and may be subject to change. It is supported by [New Monitors and Alerts](/product/new-monitors-and-alerts/) and may not be viewable in the UI today.

Bulk delete Monitors for a given organization
- `sentry call Fetch a Monitor` (GET /api/0/organizations/{organization_id_or_slug}/detectors/{detector_id}/) — read — ⚠️ This endpoint is currently in **beta** and may be subject to change. It is supported by [New Monitors and Alerts](/product/new-monitors-and-alerts/) and may not be viewable in the UI today.

Return details on an individual monitor
- `sentry call Update a Monitor by ID` (PUT /api/0/organizations/{organization_id_or_slug}/detectors/{detector_id}/) — update — ⚠️ This endpoint is currently in **beta** and may be subject to change. It is supported by [New Monitors and Alerts](/product/new-monitors-and-alerts/) and may not be viewable in the UI today.

Update an existing monitor
- `sentry call Delete a Monitor` (DELETE /api/0/organizations/{organization_id_or_slug}/detectors/{detector_id}/) — delete — ⚠️ This endpoint is currently in **beta** and may be subject to change. It is supported by [New Monitors and Alerts](/product/new-monitors-and-alerts/) and may not be viewable in the UI today.

Delete a monitor
- `sentry call List an Organization's Discover Saved Queries` (GET /api/0/organizations/{organization_id_or_slug}/discover/saved/) — read — Retrieve a list of saved queries that are associated with the given organization.
- `sentry call Create a New Saved Query` (POST /api/0/organizations/{organization_id_or_slug}/discover/saved/) — search — Create a new saved query for the given organization.
- `sentry call Retrieve an Organization's Discover Saved Query` (GET /api/0/organizations/{organization_id_or_slug}/discover/saved/{query_id}/) — search — Retrieve a saved query.
- `sentry call Edit an Organization's Discover Saved Query` (PUT /api/0/organizations/{organization_id_or_slug}/discover/saved/{query_id}/) — update — Modify a saved query.
- `sentry call Delete an Organization's Discover Saved Query` (DELETE /api/0/organizations/{organization_id_or_slug}/discover/saved/{query_id}/) — delete — Delete a saved query.
- `sentry call List an Organization's Environments` (GET /api/0/organizations/{organization_id_or_slug}/environments/) — read — Lists an organization's environments.
- `sentry call Resolve an Event ID` (GET /api/0/organizations/{organization_id_or_slug}/eventids/{event_id}/) — read — This resolves an event ID to the project slug and internal issue ID and internal event ID.
- `sentry call Query Explore Events in Table Format` (GET /api/0/organizations/{organization_id_or_slug}/events/) — search — Retrieves explore data for a given organization.

**Note**: This endpoint is intended to get a table of results, and is not for doing a full export of data sent to
Sentry.

The `field` query parameter determines what fields will be selected in the `data` and `meta` keys of the endpoint response.
- The `data` key contains a list of results row by row that match the `query` made
- The `meta` key contains information about the response, including the unit or type of the fields requested
- `sentry call Query Explore Events in Timeseries Format` (GET /api/0/organizations/{organization_id_or_slug}/events-timeseries/) — read — Retrieves explore data for a given organization as a timeseries.

This endpoint can return timeseries for either 1 or many axis, and results grouped to the top events depending
on the parameters passed
- `sentry call Create an External User` (POST /api/0/organizations/{organization_id_or_slug}/external-users/) — create — Link a user from an external provider to a Sentry user.
- `sentry call Update an External User` (PUT /api/0/organizations/{organization_id_or_slug}/external-users/{external_user_id}/) — update — Update a user in an external provider that is currently linked to a Sentry user.
- `sentry call Delete an External User` (DELETE /api/0/organizations/{organization_id_or_slug}/external-users/{external_user_id}/) — delete — Delete the link between a user from an external provider and a Sentry user.
- `sentry call Retrieve Data Forwarders for an Organization` (GET /api/0/organizations/{organization_id_or_slug}/forwarding/) — read — Returns a list of data forwarders for an organization.
- `sentry call Create a Data Forwarder for an Organization` (POST /api/0/organizations/{organization_id_or_slug}/forwarding/) — create — Creates a new data forwarder for an organization.
Only one data forwarder can be created per provider for a given organization.

Project-specific overrides can only be created after creating the data forwarder.
- `sentry call Update a Data Forwarder for an Organization` (PUT /api/0/organizations/{organization_id_or_slug}/forwarding/{data_forwarder_id}/) — update — Updates a data forwarder for an organization or update a project-specific override.
Updates to the data forwarder's configuration require `org:write` permissions, and the entire
configuration to be provided, including the `project_ids` field.

To configure project-specific overrides, specify only the following fields:

  - 'project_id': The ID of the project to create/modify the override for.
  - 'overrides': Follows the same format as `config` but all provider fields are optional, since only specified fields are overridden.
  - 'is_enabled': To enable/disable the forwarder for events on the specific project.

Overrides can be performed with `project:write` permissions on the project being modified.
- `sentry call Delete a Data Forwarder for an Organization` (DELETE /api/0/organizations/{organization_id_or_slug}/forwarding/{data_forwarder_id}/) — delete — Deletes a data forwarder for an organization. All project-specific overrides will be deleted as well.
- `sentry call List an Organization's Available Integrations` (GET /api/0/organizations/{organization_id_or_slug}/integrations/) — read — Lists all the available Integrations for an Organization.
- `sentry call Retrieve an Integration for an Organization` (GET /api/0/organizations/{organization_id_or_slug}/integrations/{integration_id}/) — read — OrganizationIntegrationBaseEndpoints expect both Integration and
OrganizationIntegration DB entries to exist for a given organization and
integration_id.
- `sentry call Delete an Integration for an Organization` (DELETE /api/0/organizations/{organization_id_or_slug}/integrations/{integration_id}/) — delete — OrganizationIntegrationBaseEndpoints expect both Integration and
OrganizationIntegration DB entries to exist for a given organization and
integration_id.
- `sentry call List an Organization's Issues` (GET /api/0/organizations/{organization_id_or_slug}/issues/) — search — Return a list of issues for an organization. All parameters are supplied as query string parameters. A default query of `is:unresolved` is applied. To return all results, use an empty query value (i.e. ``?query=`). 
- `sentry call Bulk Mutate an Organization's Issues` (PUT /api/0/organizations/{organization_id_or_slug}/issues/) — update — Bulk mutate various attributes on a maxmimum of 1000 issues. 
- For non-status updates, the `id` query parameter is required. 
- For status updates, the `id` query parameter may be omitted to update issues that match the filtering. 
If any IDs are out of scope, the data won't be mutated but the endpoint will still produce a successful response. For example, if no issues were found matching the criteria, a HTTP 204 is returned.
- `sentry call Bulk Remove an Organization's Issues` (DELETE /api/0/organizations/{organization_id_or_slug}/issues/) — delete — Permanently remove the given issues. If IDs are provided, queries and filtering will be ignored. If any IDs are out of scope, the data won't be mutated but the endpoint will still produce a successful response. For example, if no issues were found matching the criteria, a HTTP 204 is returned.
- `sentry call List an Organization's Members` (GET /api/0/organizations/{organization_id_or_slug}/members/) — read — List all organization members.

Response includes pending invites that are approved by organization owners or managers but waiting to be accepted by the invitee.
- `sentry call Add a Member to an Organization` (POST /api/0/organizations/{organization_id_or_slug}/members/) — create — Add or invite a member to an organization.
- `sentry call Retrieve an Organization Member` (GET /api/0/organizations/{organization_id_or_slug}/members/{member_id}/) — read — Retrieve an organization member's details.

Response will be a pending invite if it has been approved by organization owners or managers but is waiting to be accepted by the invitee.
- `sentry call Update an Organization Member's Roles` (PUT /api/0/organizations/{organization_id_or_slug}/members/{member_id}/) — update — Update a member's [organization-level](https://docs.sentry.io/organization/membership/#organization-level-roles) and [team-level](https://docs.sentry.io/organization/membership/#team-level-roles) roles.

Note that for changing organization-roles, this endpoint is restricted to
[user auth tokens](https://docs.sentry.io/account/auth-tokens/#user-auth-tokens).
Additionally, both the original and desired organization role must have
the same or lower permissions than the role of the organization user making the request

For example, an organization Manager may change someone's role from
Member to Manager, but not to Owner.
- `sentry call Delete an Organization Member` (DELETE /api/0/organizations/{organization_id_or_slug}/members/{member_id}/) — delete — Remove an organization member.
- `sentry call Add an Organization Member to a Team` (POST /api/0/organizations/{organization_id_or_slug}/members/{member_id}/teams/{team_id_or_slug}/) — create — This request can return various success codes depending on the context of the team:
- **`201`**: The member has been successfully added.
- **`202`**: The member needs permission to join the team and an access request
has been generated.
- **`204`**: The member is already on the team.

If the team is provisioned through an identity provider, the member cannot join the
team through Sentry.

Note the permission scopes vary depending on the organization setting `"Open Membership"`
and the type of authorization token. The following table outlines the accepted scopes.
<table style="width: 100%;">
<thead>
    <tr>
    <th style="width: 33%;"></th>
    <th colspan="2" style="text-align: center; font-weight: bold; width: 33%;">Open Membership</th>
    </tr>
</thead>
<tbody>
    <tr>
    <td style="width: 34%;"></td>
    <td style="text-align: center; font-weight: bold; width: 33%;">On</td>
    <td style="text-align: center; font-weight: bold; width: 33%;">Off</td>
    </tr>
    <tr>
    <td style="text-align: center; font-weight: bold; vertical-align: middle;"><a
    href="https://docs.sentry.io/account/auth-tokens/#internal-integrations">Internal Integration Token</a></td>
    <td style="text-align: left; width: 33%;">
        <ul style="list-style-type: none; padding-left: 0;">
        <li><strong style="color: #9c5f99;">&bull; org:read</strong></li>
        </ul>
    </td>
    <td style="text-align: left; width: 33%;">
        <ul style="list-style-type: none; padding-left: 0;">
        <li><strong style="color: #9c5f99;">&bull; org:write</strong></li>
        <li><strong style="color: #9c5f99;">&bull; team:write</strong></li>
        </ul>
    </td>
    </tr>
    <tr>
    <td style="text-align: center; font-weight: bold; vertical-align: middle;"><a
    href="https://docs.sentry.io/account/auth-tokens/#user-auth-tokens">User Auth Token</a></td>
    <td style="text-align: left; width: 33%;">
        <ul style="list-style-type: none; padding-left: 0;">
        <li><strong style="color: #9c5f99;">&bull; org:read</strong></li>
        </ul>
    </td>
    <td style="text-align: left; width: 33%;">
        <ul style="list-style-type: none; padding-left: 0;">
        <li><strong style="color: #9c5f99;">&bull; org:read*</strong></li>
        <li><strong style="color: #9c5f99;">&bull; org:write</strong></li>
        <li><strong style="color: #9c5f99;">&bull; org:read +</strong></li>
        <li><strong style="color: #9c5f99;">&nbsp; &nbsp;team:write**</strong></li>
        </ul>
    </td>
    </tr>
</tbody>
</table>


*Organization members are restricted to this scope. When sending a request, it will always
return a 202 and request an invite to the team.


\*\*Team Admins must have both **`org:read`** and **`team:write`** scopes in their user
authorization token to add members to their teams.
- `sentry call Update an Organization Member's Team Role` (PUT /api/0/organizations/{organization_id_or_slug}/members/{member_id}/teams/{team_id_or_slug}/) — update — The relevant organization member must already be a part of the team.

Note that for organization admins, managers, and owners, they are
automatically granted a minimum team role of `admin` on all teams they
are part of. Read more about [team roles](https://docs.sentry.io/product/teams/roles/).
- `sentry call Delete an Organization Member from a Team` (DELETE /api/0/organizations/{organization_id_or_slug}/members/{member_id}/teams/{team_id_or_slug}/) — delete — Delete an organization member from a team.

Note the permission scopes vary depending on the type of authorization token. The following
table outlines the accepted scopes.
<table style="width: 100%;">
    <tr style="width: 50%;">
        <td style="width: 50%; text-align: center; font-weight: bold; vertical-align: middle;"><a href="https://docs.sentry.io/api/auth/#auth-tokens">Org Auth Token</a></td>
        <td style="width: 50%; text-align: left;">
            <ul style="list-style-type: none; padding-left: 0;">
                <li><strong style="color: #9c5f99;">&bull; org:write</strong></li>
                <li><strong style="color: #9c5f99;">&bull; org:admin</strong></li>
                <li><strong style="color: #9c5f99;">&bull; team:admin</strong></li>
            </ul>
        </td>
    </tr>
    <tr style="width: 50%;">
        <td style="width: 50%; text-align: center; font-weight: bold; vertical-align: middle;"><a href="https://docs.sentry.io/api/auth/#user-authentication-tokens">User Auth Token</a></td>
        <td style="width: 50%; text-align: left;">
            <ul style="list-style-type: none; padding-left: 0;">
                <li><strong style="color: #9c5f99;">&bull; org:read*</strong></li>
                <li><strong style="color: #9c5f99;">&bull; org:write</strong></li>
                <li><strong style="color: #9c5f99;">&bull; org:admin</strong></li>
                <li><strong style="color: #9c5f99;">&bull; team:admin</strong></li>
                <li><strong style="color: #9c5f99;">&bull; org:read + team:admin**</strong></li>
            </ul>
        </td>
    </tr>
</table>


\***`org:read`** can only be used to remove yourself from the teams you are a member of.


\*\*Team Admins must have both **`org:read`** and **`team:admin`** scopes in their user
authorization token to delete members from their teams.
- `sentry call Retrieve Monitors for an Organization` (GET /api/0/organizations/{organization_id_or_slug}/monitors/) — read — Lists monitors, including nested monitor environments. May be filtered to a project or environment.
- `sentry call Create a Monitor` (POST /api/0/organizations/{organization_id_or_slug}/monitors/) — create — Create a new monitor.
- `sentry call Retrieve a Monitor` (GET /api/0/organizations/{organization_id_or_slug}/monitors/{monitor_id_or_slug}/) — read — Retrieves details for a monitor.
- `sentry call Update a Monitor` (PUT /api/0/organizations/{organization_id_or_slug}/monitors/{monitor_id_or_slug}/) — update — Update a monitor.
- `sentry call Delete a Monitor or Monitor Environments` (DELETE /api/0/organizations/{organization_id_or_slug}/monitors/{monitor_id_or_slug}/) — delete — Delete a monitor or monitor environments.
- … +150 more — run `sentry operations --json` for the full list

## Safety
Write/destructive operations dry-run by default. Append `--yes` only on explicit user confirmation.
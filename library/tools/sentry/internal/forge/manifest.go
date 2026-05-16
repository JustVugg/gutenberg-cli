package forge

import "encoding/json"

const manifestJSON = `{
  "schemaVersion": "gutenberg.blueprint.v1",
  "source": "/tmp/sentry-api.json",
  "name": "sentry",
  "slug": "sentry",
  "envPrefix": "SENTRY",
  "description": "Sentry Public API",
  "version": "v0",
  "baseUrls": [
    "https://{region}.sentry.io"
  ],
  "auth": {
    "mode": "detected",
    "schemes": [
      {
        "name": "auth_token",
        "type": "http",
        "in": null,
        "header": null,
        "scheme": "bearer",
        "flows": {}
      },
      {
        "name": "dsn",
        "type": "http",
        "in": null,
        "header": null,
        "scheme": "DSN",
        "flows": {}
      }
    ],
    "env": "API_KEY",
    "oauth": false
  },
  "tags": [
    "Alerts",
    "Crons",
    "Dashboards",
    "Discover",
    "Environments",
    "Events",
    "Explore",
    "Integration",
    "Integrations",
    "Mobile Builds",
    "Monitors",
    "Organizations",
    "Prevent",
    "Projects",
    "Releases",
    "Replays",
    "SCIM",
    "Seer",
    "Teams",
    "Users"
  ],
  "operations": [
    {
      "id": "List Your Organizations",
      "method": "GET",
      "path": "/api/0/organizations/",
      "tag": "Users",
      "summary": "Return a list of organizations available to the authenticated session in a region.\nThis is particularly useful for requests with a user bound context. For API key-based requests this will only return the organization that belongs to the key.",
      "description": "Return a list of organizations available to the authenticated session in a region.\nThis is particularly useful for requests with a user bound context. For API key-based requests this will only return the organization that belongs to the key.",
      "parameters": [
        {
          "name": "owner",
          "in": "query",
          "required": false,
          "description": "Specify ` + "`" + `true` + "`" + ` to restrict results to organizations in which you are an owner.",
          "type": "boolean"
        },
        {
          "name": "cursor",
          "in": "query",
          "required": false,
          "description": "A pointer to the last object fetched and its sort order; used to retrieve the next or previous results.",
          "type": "string"
        },
        {
          "name": "query",
          "in": "query",
          "required": false,
          "description": "Filters results by using [query syntax](/product/sentry-basics/search/).\n\nValid query fields include:\n- ` + "`" + `id` + "`" + `: The organization ID\n- ` + "`" + `slug` + "`" + `: The organization slug\n- ` + "`" + `status` + "`" + `: The organization's current status (one of ` + "`" + `active` + "`" + `, ` + "`" + `pending_deletion` + "`" + `, or ` + "`" + `deletion_in_progress` + "`" + `)\n- ` + "`" + `email` + "`" + ` or ` + "`" + `member_id` + "`" + `: Filter your organizations by the emails or [organization member IDs](/api/organizations/list-an-organizations-members/) of specific members included\n- ` + "`" + `query` + "`" + `: Filter your organizations by name, slug, and members that contain this substring\n\nExample: ` + "`" + `query=(slug:foo AND status:active) OR (email:[thing-one@example.com,thing-two@example.com] AND query:bar)` + "`" + `\n",
          "type": "string"
        },
        {
          "name": "sortBy",
          "in": "query",
          "required": false,
          "description": "The field to sort results by, in descending order. If not specified the results are sorted by the date they were created.\n\nValid fields include:\n- ` + "`" + `members` + "`" + `: By number of members\n- ` + "`" + `events` + "`" + `: By number of events in the past 24 hours\n",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "401",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [],
        "query": [
          "owner",
          "cursor",
          "query",
          "sortBy"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read-list",
      "pagination": {
        "style": "cursor",
        "cursorParam": "cursor"
      },
      "related": []
    },
    {
      "id": "Retrieve an Organization",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/",
      "tag": "Organizations",
      "summary": "Return details on an individual organization, including various details\nsuch as membership access and teams.",
      "description": "Return details on an individual organization, including various details\nsuch as membership access and teams.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "detailed",
          "in": "query",
          "required": false,
          "description": "\nSpecify ` + "`" + `\"0\"` + "`" + ` to return organization details that do not include projects or teams.\n",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "401",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug"
        ],
        "query": [
          "detailed"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": null,
      "related": [
        {
          "id": "Update an Organization",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Update an Organization",
      "method": "PUT",
      "path": "/api/0/organizations/{organization_id_or_slug}/",
      "tag": "Organizations",
      "summary": "Update various attributes and configurable settings for the given organization.",
      "description": "Update various attributes and configurable settings for the given organization.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        }
      ],
      "hasRequestBody": true,
      "risk": "write",
      "cacheable": false,
      "responseCodes": [
        "200",
        "400",
        "401",
        "403",
        "404",
        "409",
        "413"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug"
        ],
        "query": [],
        "headers": [],
        "body": true
      },
      "graphql": null,
      "kind": "update",
      "pagination": null,
      "related": [
        {
          "id": "Retrieve an Organization",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Get Integration Provider Information",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/config/integrations/",
      "tag": "Integrations",
      "summary": "Get integration provider information about all available integrations for an organization.",
      "description": "Get integration provider information about all available integrations for an organization.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "providerKey",
          "in": "query",
          "required": false,
          "description": "Specific integration provider to filter by such as ` + "`" + `slack` + "`" + `. See our [Integrations Documentation](/product/integrations/) for an updated list of providers.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug"
        ],
        "query": [
          "providerKey"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": null,
      "related": []
    },
    {
      "id": "List an Organization's Custom Dashboards",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/dashboards/",
      "tag": "Dashboards",
      "summary": "Retrieve a list of custom dashboards that are associated with the given organization.",
      "description": "Retrieve a list of custom dashboards that are associated with the given organization.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "per_page",
          "in": "query",
          "required": false,
          "description": "Limit the number of rows to return in the result. Default and maximum allowed is 100.",
          "type": "integer"
        },
        {
          "name": "cursor",
          "in": "query",
          "required": false,
          "description": "A pointer to the last object fetched and its sort order; used to retrieve the next or previous results.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "400",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug"
        ],
        "query": [
          "per_page",
          "cursor"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": {
        "style": "cursor",
        "cursorParam": "cursor"
      },
      "related": [
        {
          "id": "Create a New Dashboard for an Organization",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Create a New Dashboard for an Organization",
      "method": "POST",
      "path": "/api/0/organizations/{organization_id_or_slug}/dashboards/",
      "tag": "Dashboards",
      "summary": "Create a new dashboard for the given Organization",
      "description": "Create a new dashboard for the given Organization",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        }
      ],
      "hasRequestBody": true,
      "risk": "write",
      "cacheable": false,
      "responseCodes": [
        "201",
        "400",
        "403",
        "404",
        "409"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug"
        ],
        "query": [],
        "headers": [],
        "body": true
      },
      "graphql": null,
      "kind": "create",
      "pagination": null,
      "related": [
        {
          "id": "List an Organization's Custom Dashboards",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Retrieve an Organization's Custom Dashboard",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/dashboards/{dashboard_id}/",
      "tag": "Dashboards",
      "summary": "Return details about an organization's custom dashboard.",
      "description": "Return details about an organization's custom dashboard.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "dashboard_id",
          "in": "path",
          "required": true,
          "description": "The ID of the dashboard you'd like to retrieve.",
          "type": "integer"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "dashboard_id"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": null,
      "related": [
        {
          "id": "Edit an Organization's Custom Dashboard",
          "role": "same-resource"
        },
        {
          "id": "Delete an Organization's Custom Dashboard",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Edit an Organization's Custom Dashboard",
      "method": "PUT",
      "path": "/api/0/organizations/{organization_id_or_slug}/dashboards/{dashboard_id}/",
      "tag": "Dashboards",
      "summary": "Edit an organization's custom dashboard as well as any bulk\nedits on widgets that may have been made. (For example, widgets\nthat have been rearranged, updated queries and fields, specific\ndisplay types, and so on.)",
      "description": "Edit an organization's custom dashboard as well as any bulk\nedits on widgets that may have been made. (For example, widgets\nthat have been rearranged, updated queries and fields, specific\ndisplay types, and so on.)",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "dashboard_id",
          "in": "path",
          "required": true,
          "description": "The ID of the dashboard you'd like to retrieve.",
          "type": "integer"
        }
      ],
      "hasRequestBody": true,
      "risk": "write",
      "cacheable": false,
      "responseCodes": [
        "200",
        "400",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "dashboard_id"
        ],
        "query": [],
        "headers": [],
        "body": true
      },
      "graphql": null,
      "kind": "update",
      "pagination": null,
      "related": [
        {
          "id": "Retrieve an Organization's Custom Dashboard",
          "role": "same-resource"
        },
        {
          "id": "Delete an Organization's Custom Dashboard",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Delete an Organization's Custom Dashboard",
      "method": "DELETE",
      "path": "/api/0/organizations/{organization_id_or_slug}/dashboards/{dashboard_id}/",
      "tag": "Dashboards",
      "summary": "Delete an organization's custom dashboard.",
      "description": "Delete an organization's custom dashboard.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "dashboard_id",
          "in": "path",
          "required": true,
          "description": "The ID of the dashboard you'd like to retrieve.",
          "type": "integer"
        }
      ],
      "hasRequestBody": false,
      "risk": "destructive",
      "cacheable": false,
      "responseCodes": [
        "204",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "dashboard_id"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "delete",
      "pagination": null,
      "related": [
        {
          "id": "Retrieve an Organization's Custom Dashboard",
          "role": "same-resource"
        },
        {
          "id": "Edit an Organization's Custom Dashboard",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Fetch an Organization's Monitors",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/detectors/",
      "tag": "Monitors",
      "summary": "⚠️ This endpoint is currently in **beta** and may be subject to change. It is supported by [New Monitors and Alerts](/product/new-monitors-and-alerts/) and may not be viewable in the UI today.\n\nList an Organization's Monitors",
      "description": "⚠️ This endpoint is currently in **beta** and may be subject to change. It is supported by [New Monitors and Alerts](/product/new-monitors-and-alerts/) and may not be viewable in the UI today.\n\nList an Organization's Monitors",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "project",
          "in": "query",
          "required": false,
          "description": "The IDs of projects to filter by. ` + "`" + `-1` + "`" + ` means all available projects.\nFor example, the following are valid parameters:\n- ` + "`" + `/?project=1234&project=56789` + "`" + `\n- ` + "`" + `/?project=-1` + "`" + `\n",
          "type": "array"
        },
        {
          "name": "query",
          "in": "query",
          "required": false,
          "description": "An optional search query for filtering monitors.\n\nAvailable fields are:\n- ` + "`" + `name` + "`" + `\n- ` + "`" + `type` + "`" + `: e.g. ` + "`" + `error` + "`" + `, ` + "`" + `metric_issue` + "`" + `, ` + "`" + `issue_stream` + "`" + `\n- ` + "`" + `assignee` + "`" + `: email, username, #team, me, none\n        ",
          "type": "string"
        },
        {
          "name": "sortBy",
          "in": "query",
          "required": false,
          "description": "The property to sort results by. If not specified, the results are sorted by id.\n\nAvailable fields are:\n- ` + "`" + `name` + "`" + `\n- ` + "`" + `id` + "`" + `\n- ` + "`" + `type` + "`" + `\n- ` + "`" + `connectedWorkflows` + "`" + `\n- ` + "`" + `latestGroup` + "`" + `\n- ` + "`" + `openIssues` + "`" + `\n\nPrefix with ` + "`" + `-` + "`" + ` to sort in descending order.\n        ",
          "type": "string"
        },
        {
          "name": "id",
          "in": "query",
          "required": false,
          "description": "The ID of the monitor you'd like to query.",
          "type": "array"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "400",
        "401",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug"
        ],
        "query": [
          "project",
          "query",
          "sortBy",
          "id"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": null,
      "related": [
        {
          "id": "Mutate an Organization's Monitors",
          "role": "same-resource"
        },
        {
          "id": "Bulk Delete Monitors",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Mutate an Organization's Monitors",
      "method": "PUT",
      "path": "/api/0/organizations/{organization_id_or_slug}/detectors/",
      "tag": "Monitors",
      "summary": "⚠️ This endpoint is currently in **beta** and may be subject to change. It is supported by [New Monitors and Alerts](/product/new-monitors-and-alerts/) and may not be viewable in the UI today.\n\nBulk enable or disable an Organization's Monitors",
      "description": "⚠️ This endpoint is currently in **beta** and may be subject to change. It is supported by [New Monitors and Alerts](/product/new-monitors-and-alerts/) and may not be viewable in the UI today.\n\nBulk enable or disable an Organization's Monitors",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "project",
          "in": "query",
          "required": false,
          "description": "The IDs of projects to filter by. ` + "`" + `-1` + "`" + ` means all available projects.\nFor example, the following are valid parameters:\n- ` + "`" + `/?project=1234&project=56789` + "`" + `\n- ` + "`" + `/?project=-1` + "`" + `\n",
          "type": "array"
        },
        {
          "name": "query",
          "in": "query",
          "required": false,
          "description": "An optional search query for filtering monitors.\n\nAvailable fields are:\n- ` + "`" + `name` + "`" + `\n- ` + "`" + `type` + "`" + `: e.g. ` + "`" + `error` + "`" + `, ` + "`" + `metric_issue` + "`" + `, ` + "`" + `issue_stream` + "`" + `\n- ` + "`" + `assignee` + "`" + `: email, username, #team, me, none\n        ",
          "type": "string"
        },
        {
          "name": "id",
          "in": "query",
          "required": false,
          "description": "The ID of the monitor you'd like to query.",
          "type": "array"
        }
      ],
      "hasRequestBody": true,
      "risk": "write",
      "cacheable": false,
      "responseCodes": [
        "200",
        "400",
        "401",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug"
        ],
        "query": [
          "project",
          "query",
          "id"
        ],
        "headers": [],
        "body": true
      },
      "graphql": null,
      "kind": "update",
      "pagination": null,
      "related": [
        {
          "id": "Fetch an Organization's Monitors",
          "role": "same-resource"
        },
        {
          "id": "Bulk Delete Monitors",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Bulk Delete Monitors",
      "method": "DELETE",
      "path": "/api/0/organizations/{organization_id_or_slug}/detectors/",
      "tag": "Monitors",
      "summary": "⚠️ This endpoint is currently in **beta** and may be subject to change. It is supported by [New Monitors and Alerts](/product/new-monitors-and-alerts/) and may not be viewable in the UI today.\n\nBulk delete Monitors for a given organization",
      "description": "⚠️ This endpoint is currently in **beta** and may be subject to change. It is supported by [New Monitors and Alerts](/product/new-monitors-and-alerts/) and may not be viewable in the UI today.\n\nBulk delete Monitors for a given organization",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "project",
          "in": "query",
          "required": false,
          "description": "The IDs of projects to filter by. ` + "`" + `-1` + "`" + ` means all available projects.\nFor example, the following are valid parameters:\n- ` + "`" + `/?project=1234&project=56789` + "`" + `\n- ` + "`" + `/?project=-1` + "`" + `\n",
          "type": "array"
        },
        {
          "name": "query",
          "in": "query",
          "required": false,
          "description": "An optional search query for filtering monitors.\n\nAvailable fields are:\n- ` + "`" + `name` + "`" + `\n- ` + "`" + `type` + "`" + `: e.g. ` + "`" + `error` + "`" + `, ` + "`" + `metric_issue` + "`" + `, ` + "`" + `issue_stream` + "`" + `\n- ` + "`" + `assignee` + "`" + `: email, username, #team, me, none\n        ",
          "type": "string"
        },
        {
          "name": "sortBy",
          "in": "query",
          "required": false,
          "description": "The property to sort results by. If not specified, the results are sorted by id.\n\nAvailable fields are:\n- ` + "`" + `name` + "`" + `\n- ` + "`" + `id` + "`" + `\n- ` + "`" + `type` + "`" + `\n- ` + "`" + `connectedWorkflows` + "`" + `\n- ` + "`" + `latestGroup` + "`" + `\n- ` + "`" + `openIssues` + "`" + `\n\nPrefix with ` + "`" + `-` + "`" + ` to sort in descending order.\n        ",
          "type": "string"
        },
        {
          "name": "id",
          "in": "query",
          "required": false,
          "description": "The ID of the monitor you'd like to query.",
          "type": "array"
        }
      ],
      "hasRequestBody": false,
      "risk": "destructive",
      "cacheable": false,
      "responseCodes": [
        "200",
        "204",
        "400",
        "401",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug"
        ],
        "query": [
          "project",
          "query",
          "sortBy",
          "id"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "delete",
      "pagination": null,
      "related": [
        {
          "id": "Fetch an Organization's Monitors",
          "role": "same-resource"
        },
        {
          "id": "Mutate an Organization's Monitors",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Fetch a Monitor",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/detectors/{detector_id}/",
      "tag": "Monitors",
      "summary": "⚠️ This endpoint is currently in **beta** and may be subject to change. It is supported by [New Monitors and Alerts](/product/new-monitors-and-alerts/) and may not be viewable in the UI today.\n\nReturn details on an individual monitor",
      "description": "⚠️ This endpoint is currently in **beta** and may be subject to change. It is supported by [New Monitors and Alerts](/product/new-monitors-and-alerts/) and may not be viewable in the UI today.\n\nReturn details on an individual monitor",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "detector_id",
          "in": "path",
          "required": true,
          "description": "The ID of the monitor you'd like to query.",
          "type": "integer"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "400",
        "401",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "detector_id"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": null,
      "related": [
        {
          "id": "Update a Monitor by ID",
          "role": "same-resource"
        },
        {
          "id": "Delete a Monitor",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Update a Monitor by ID",
      "method": "PUT",
      "path": "/api/0/organizations/{organization_id_or_slug}/detectors/{detector_id}/",
      "tag": "Monitors",
      "summary": "⚠️ This endpoint is currently in **beta** and may be subject to change. It is supported by [New Monitors and Alerts](/product/new-monitors-and-alerts/) and may not be viewable in the UI today.\n\nUpdate an existing monitor",
      "description": "⚠️ This endpoint is currently in **beta** and may be subject to change. It is supported by [New Monitors and Alerts](/product/new-monitors-and-alerts/) and may not be viewable in the UI today.\n\nUpdate an existing monitor",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "detector_id",
          "in": "path",
          "required": true,
          "description": "The ID of the monitor you'd like to query.",
          "type": "integer"
        }
      ],
      "hasRequestBody": true,
      "risk": "write",
      "cacheable": false,
      "responseCodes": [
        "200",
        "400",
        "401",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "detector_id"
        ],
        "query": [],
        "headers": [],
        "body": true
      },
      "graphql": null,
      "kind": "update",
      "pagination": null,
      "related": [
        {
          "id": "Fetch a Monitor",
          "role": "same-resource"
        },
        {
          "id": "Delete a Monitor",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Delete a Monitor",
      "method": "DELETE",
      "path": "/api/0/organizations/{organization_id_or_slug}/detectors/{detector_id}/",
      "tag": "Monitors",
      "summary": "⚠️ This endpoint is currently in **beta** and may be subject to change. It is supported by [New Monitors and Alerts](/product/new-monitors-and-alerts/) and may not be viewable in the UI today.\n\nDelete a monitor",
      "description": "⚠️ This endpoint is currently in **beta** and may be subject to change. It is supported by [New Monitors and Alerts](/product/new-monitors-and-alerts/) and may not be viewable in the UI today.\n\nDelete a monitor",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "detector_id",
          "in": "path",
          "required": true,
          "description": "The ID of the monitor you'd like to query.",
          "type": "integer"
        }
      ],
      "hasRequestBody": false,
      "risk": "destructive",
      "cacheable": false,
      "responseCodes": [
        "204",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "detector_id"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "delete",
      "pagination": null,
      "related": [
        {
          "id": "Fetch a Monitor",
          "role": "same-resource"
        },
        {
          "id": "Update a Monitor by ID",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "List an Organization's Discover Saved Queries",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/discover/saved/",
      "tag": "Discover",
      "summary": "Retrieve a list of saved queries that are associated with the given organization.",
      "description": "Retrieve a list of saved queries that are associated with the given organization.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "per_page",
          "in": "query",
          "required": false,
          "description": "Limit the number of rows to return in the result. Default and maximum allowed is 100.",
          "type": "integer"
        },
        {
          "name": "cursor",
          "in": "query",
          "required": false,
          "description": "A pointer to the last object fetched and its sort order; used to retrieve the next or previous results.",
          "type": "string"
        },
        {
          "name": "query",
          "in": "query",
          "required": false,
          "description": "The name of the Discover query you'd like to filter by.",
          "type": "string"
        },
        {
          "name": "sortBy",
          "in": "query",
          "required": false,
          "description": "The property to sort results by. If not specified, the results are sorted by query name.\n\nAvailable fields are:\n- ` + "`" + `name` + "`" + `\n- ` + "`" + `dateCreated` + "`" + `\n- ` + "`" + `dateUpdated` + "`" + `\n- ` + "`" + `mostPopular` + "`" + `\n- ` + "`" + `recentlyViewed` + "`" + `\n- ` + "`" + `myqueries` + "`" + `\n        ",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "400",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug"
        ],
        "query": [
          "per_page",
          "cursor",
          "query",
          "sortBy"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": {
        "style": "cursor",
        "cursorParam": "cursor"
      },
      "related": [
        {
          "id": "Create a New Saved Query",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Create a New Saved Query",
      "method": "POST",
      "path": "/api/0/organizations/{organization_id_or_slug}/discover/saved/",
      "tag": "Discover",
      "summary": "Create a new saved query for the given organization.",
      "description": "Create a new saved query for the given organization.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        }
      ],
      "hasRequestBody": true,
      "risk": "write",
      "cacheable": false,
      "responseCodes": [
        "201",
        "400",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug"
        ],
        "query": [],
        "headers": [],
        "body": true
      },
      "graphql": null,
      "kind": "search",
      "pagination": null,
      "related": [
        {
          "id": "List an Organization's Discover Saved Queries",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Retrieve an Organization's Discover Saved Query",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/discover/saved/{query_id}/",
      "tag": "Discover",
      "summary": "Retrieve a saved query.",
      "description": "Retrieve a saved query.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "query_id",
          "in": "path",
          "required": true,
          "description": "The ID of the Discover query you'd like to retrieve.",
          "type": "integer"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "query_id"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "search",
      "pagination": null,
      "related": [
        {
          "id": "Edit an Organization's Discover Saved Query",
          "role": "same-resource"
        },
        {
          "id": "Delete an Organization's Discover Saved Query",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Edit an Organization's Discover Saved Query",
      "method": "PUT",
      "path": "/api/0/organizations/{organization_id_or_slug}/discover/saved/{query_id}/",
      "tag": "Discover",
      "summary": "Modify a saved query.",
      "description": "Modify a saved query.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "query_id",
          "in": "path",
          "required": true,
          "description": "The ID of the Discover query you'd like to retrieve.",
          "type": "integer"
        }
      ],
      "hasRequestBody": true,
      "risk": "write",
      "cacheable": false,
      "responseCodes": [
        "200",
        "400",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "query_id"
        ],
        "query": [],
        "headers": [],
        "body": true
      },
      "graphql": null,
      "kind": "update",
      "pagination": null,
      "related": [
        {
          "id": "Retrieve an Organization's Discover Saved Query",
          "role": "same-resource"
        },
        {
          "id": "Delete an Organization's Discover Saved Query",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Delete an Organization's Discover Saved Query",
      "method": "DELETE",
      "path": "/api/0/organizations/{organization_id_or_slug}/discover/saved/{query_id}/",
      "tag": "Discover",
      "summary": "Delete a saved query.",
      "description": "Delete a saved query.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "query_id",
          "in": "path",
          "required": true,
          "description": "The ID of the Discover query you'd like to retrieve.",
          "type": "integer"
        }
      ],
      "hasRequestBody": false,
      "risk": "destructive",
      "cacheable": false,
      "responseCodes": [
        "204",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "query_id"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "delete",
      "pagination": null,
      "related": [
        {
          "id": "Retrieve an Organization's Discover Saved Query",
          "role": "same-resource"
        },
        {
          "id": "Edit an Organization's Discover Saved Query",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "List an Organization's Environments",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/environments/",
      "tag": "Environments",
      "summary": "Lists an organization's environments.",
      "description": "Lists an organization's environments.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "visibility",
          "in": "query",
          "required": false,
          "description": "The visibility of the environments to filter by. Defaults to ` + "`" + `visible` + "`" + `.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "400",
        "403"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug"
        ],
        "query": [
          "visibility"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": null,
      "related": []
    },
    {
      "id": "Resolve an Event ID",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/eventids/{event_id}/",
      "tag": "Organizations",
      "summary": "This resolves an event ID to the project slug and internal issue ID and internal event ID.",
      "description": "This resolves an event ID to the project slug and internal issue ID and internal event ID.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "event_id",
          "in": "path",
          "required": true,
          "description": "The event ID to look up.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "event_id"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": null,
      "related": []
    },
    {
      "id": "Query Explore Events in Table Format",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/events/",
      "tag": "Explore",
      "summary": "Retrieves explore data for a given organization.\n\n**Note**: This endpoint is intended to get a table of results, and is not for doing a full export of data sent to\nSentry.\n\nThe ` + "`" + `field` + "`" + ` query parameter determines what fields will be selected in the ` + "`" + `data` + "`" + ` and ` + "`" + `meta` + "`" + ` keys of the endpoint response.\n- The ` + "`" + `data` + "`" + ` key contains a list of results row by row that match the ` + "`" + `query` + "`" + ` made\n- The ` + "`" + `meta` + "`" + ` key contains information about the response, including the unit or type of the fields requested",
      "description": "Retrieves explore data for a given organization.\n\n**Note**: This endpoint is intended to get a table of results, and is not for doing a full export of data sent to\nSentry.\n\nThe ` + "`" + `field` + "`" + ` query parameter determines what fields will be selected in the ` + "`" + `data` + "`" + ` and ` + "`" + `meta` + "`" + ` keys of the endpoint response.\n- The ` + "`" + `data` + "`" + ` key contains a list of results row by row that match the ` + "`" + `query` + "`" + ` made\n- The ` + "`" + `meta` + "`" + ` key contains information about the response, including the unit or type of the fields requested",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "field",
          "in": "query",
          "required": true,
          "description": "The fields, functions, or equations to request for the query. At most 20 fields can be selected per request. Each field can be one of the following types:\n- A built-in key field. See possible fields in the [properties table](/concepts/search/searchable-properties/), under any field that matches the dataset passed to the dataset parameter\n    - example: ` + "`" + `field=transaction` + "`" + `\n- A tag. Tags should use the ` + "`" + `tag[{name}, {type}]` + "`" + ` formatting to avoid ambiguity with any fields,\n    - example: ` + "`" + `field=tag[isEnterprise, string]` + "`" + `\n    - example: ` + "`" + `field=tag[numberOfBytes, number]` + "`" + `\n- A function which will be in the format of ` + "`" + `function_name(parameters,...)` + "`" + `. See possible functions in the [query builder documentation](/product/discover-queries/query-builder/#stacking-functions).\n    - when a function is included, Discover will group by any tags or fields\n    - example: ` + "`" + `field=count_if(transaction.duration,greater,300)` + "`" + `\n- An equation when prefixed with ` + "`" + `equation|` + "`" + `. Read more about [equations here](/product/discover-queries/query-builder/query-equations/).\n    - example: ` + "`" + `field=equation|count_if(transaction.duration,greater,300) / count() * 100` + "`" + `\n",
          "type": "array"
        },
        {
          "name": "dataset",
          "in": "query",
          "required": true,
          "description": "Which dataset to query. The chosen dataset determines which fields are queryable.\n- ` + "`" + `errors` + "`" + ` - Error events.\n- ` + "`" + `logs` + "`" + ` - Structured log events.\n- ` + "`" + `profile_functions` + "`" + ` - Function-level Profiling data.\n- ` + "`" + `spans` + "`" + ` - Distributed tracing span events.\n- ` + "`" + `tracemetrics` + "`" + ` - Application Metrics.\n- ` + "`" + `uptime_results` + "`" + ` - Uptime monitoring check results.\n",
          "type": "string"
        },
        {
          "name": "end",
          "in": "query",
          "required": false,
          "description": "The end of the period of time for the query, expected in ISO-8601 format. For example, ` + "`" + `2001-12-14T12:34:56.7890` + "`" + `.",
          "type": "string"
        },
        {
          "name": "environment",
          "in": "query",
          "required": false,
          "description": "The name of environments to filter by.",
          "type": "array"
        },
        {
          "name": "project",
          "in": "query",
          "required": false,
          "description": "The IDs of projects to filter by. ` + "`" + `-1` + "`" + ` means all available projects.\nFor example, the following are valid parameters:\n- ` + "`" + `/?project=1234&project=56789` + "`" + `\n- ` + "`" + `/?project=-1` + "`" + `\n",
          "type": "array"
        },
        {
          "name": "start",
          "in": "query",
          "required": false,
          "description": "The start of the period of time for the query, expected in ISO-8601 format. For example, ` + "`" + `2001-12-14T12:34:56.7890` + "`" + `.",
          "type": "string"
        },
        {
          "name": "statsPeriod",
          "in": "query",
          "required": false,
          "description": "The period of time for the query, will override the start & end parameters, a number followed by one of:\n- ` + "`" + `d` + "`" + ` for days\n- ` + "`" + `h` + "`" + ` for hours\n- ` + "`" + `m` + "`" + ` for minutes\n- ` + "`" + `s` + "`" + ` for seconds\n- ` + "`" + `w` + "`" + ` for weeks\n\nFor example, ` + "`" + `24h` + "`" + `, to mean query data starting from 24 hours ago to now.",
          "type": "string"
        },
        {
          "name": "per_page",
          "in": "query",
          "required": false,
          "description": "Limit the number of rows to return in the result. Default and maximum allowed is 100.",
          "type": "integer"
        },
        {
          "name": "query",
          "in": "query",
          "required": false,
          "description": "Filters results by using [query syntax](/product/sentry-basics/search/).\n\nExample: ` + "`" + `query=(transaction:foo AND release:abc) OR (transaction:[bar,baz] AND release:def)` + "`" + `\n",
          "type": "string"
        },
        {
          "name": "sort",
          "in": "query",
          "required": false,
          "description": "What to order the results of the query by. Must be something in the ` + "`" + `field` + "`" + ` list, excluding equations.",
          "type": "string"
        },
        {
          "name": "cursor",
          "in": "query",
          "required": false,
          "description": "A pointer to the last object fetched and its sort order; used to retrieve the next or previous results.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "400",
        "403"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug"
        ],
        "query": [
          "field",
          "dataset",
          "end",
          "environment",
          "project",
          "start",
          "statsPeriod",
          "per_page",
          "query",
          "sort",
          "cursor"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "search",
      "pagination": {
        "style": "cursor",
        "cursorParam": "cursor"
      },
      "related": []
    },
    {
      "id": "Query Explore Events in Timeseries Format",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/events-timeseries/",
      "tag": "Explore",
      "summary": "Retrieves explore data for a given organization as a timeseries.\n\nThis endpoint can return timeseries for either 1 or many axis, and results grouped to the top events depending\non the parameters passed",
      "description": "Retrieves explore data for a given organization as a timeseries.\n\nThis endpoint can return timeseries for either 1 or many axis, and results grouped to the top events depending\non the parameters passed",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "dataset",
          "in": "query",
          "required": true,
          "description": "Which dataset to query. The chosen dataset determines which fields are queryable.\n- ` + "`" + `errors` + "`" + ` - Error events.\n- ` + "`" + `logs` + "`" + ` - Structured log events.\n- ` + "`" + `profile_functions` + "`" + ` - Function-level Profiling data.\n- ` + "`" + `spans` + "`" + ` - Distributed tracing span events.\n- ` + "`" + `tracemetrics` + "`" + ` - Application Metrics.\n- ` + "`" + `uptime_results` + "`" + ` - Uptime monitoring check results.\n",
          "type": "string"
        },
        {
          "name": "end",
          "in": "query",
          "required": false,
          "description": "The end of the period of time for the query, expected in ISO-8601 format. For example, ` + "`" + `2001-12-14T12:34:56.7890` + "`" + `.",
          "type": "string"
        },
        {
          "name": "environment",
          "in": "query",
          "required": false,
          "description": "The name of environments to filter by.",
          "type": "array"
        },
        {
          "name": "project",
          "in": "query",
          "required": false,
          "description": "The IDs of projects to filter by. ` + "`" + `-1` + "`" + ` means all available projects.\nFor example, the following are valid parameters:\n- ` + "`" + `/?project=1234&project=56789` + "`" + `\n- ` + "`" + `/?project=-1` + "`" + `\n",
          "type": "array"
        },
        {
          "name": "start",
          "in": "query",
          "required": false,
          "description": "The start of the period of time for the query, expected in ISO-8601 format. For example, ` + "`" + `2001-12-14T12:34:56.7890` + "`" + `.",
          "type": "string"
        },
        {
          "name": "statsPeriod",
          "in": "query",
          "required": false,
          "description": "The period of time for the query, will override the start & end parameters, a number followed by one of:\n- ` + "`" + `d` + "`" + ` for days\n- ` + "`" + `h` + "`" + ` for hours\n- ` + "`" + `m` + "`" + ` for minutes\n- ` + "`" + `s` + "`" + ` for seconds\n- ` + "`" + `w` + "`" + ` for weeks\n\nFor example, ` + "`" + `24h` + "`" + `, to mean query data starting from 24 hours ago to now.",
          "type": "string"
        },
        {
          "name": "topEvents",
          "in": "query",
          "required": false,
          "description": "The number of top event results to return, must be between 1 and 10.\nWhen TopEvents is passed, both sort and groupBy are required parameters",
          "type": "integer"
        },
        {
          "name": "comparisonDelta",
          "in": "query",
          "required": false,
          "description": "The delta in seconds to return additional offset timeseries by",
          "type": "integer"
        },
        {
          "name": "interval",
          "in": "query",
          "required": false,
          "description": "The size of the bucket for the timeseries to have, must be a value smaller than the window being\nqueried. If the interval is invalid a default interval will be selected instead",
          "type": "integer"
        },
        {
          "name": "sort",
          "in": "query",
          "required": false,
          "description": "What to order the results of the query by. Must be something in the ` + "`" + `field` + "`" + ` list, excluding equations.",
          "type": "string"
        },
        {
          "name": "groupBy",
          "in": "query",
          "required": false,
          "description": "List of fields to group by, *Required* for topEvents queries as this and sort determine what the\ntop events are",
          "type": "array"
        },
        {
          "name": "yAxis",
          "in": "query",
          "required": false,
          "description": "The aggregate field to create the timeseries for, defaults to ` + "`" + `count()` + "`" + ` when not included",
          "type": "string"
        },
        {
          "name": "query",
          "in": "query",
          "required": false,
          "description": "Filters results by using [query syntax](/product/sentry-basics/search/).\n\nExample: ` + "`" + `query=(transaction:foo AND release:abc) OR (transaction:[bar,baz] AND release:def)` + "`" + `\n",
          "type": "string"
        },
        {
          "name": "disableAggregateExtrapolation",
          "in": "query",
          "required": false,
          "description": "Whether to disable the use of extrapolation and return the sampled values, due to sampling the\nnumber returned may be less than the actual values sent to Sentry",
          "type": "string"
        },
        {
          "name": "preventMetricAggregates",
          "in": "query",
          "required": false,
          "description": "Whether to throw an error when aggregates are passed in the query or groupBy",
          "type": "string"
        },
        {
          "name": "excludeOther",
          "in": "query",
          "required": false,
          "description": "Only applicable with TopEvents, whether to include the 'other' timeseries which represents all the\nevents that aren't in the top groups.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "400",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug"
        ],
        "query": [
          "dataset",
          "end",
          "environment",
          "project",
          "start",
          "statsPeriod",
          "topEvents",
          "comparisonDelta",
          "interval",
          "sort",
          "groupBy",
          "yAxis",
          "query",
          "disableAggregateExtrapolation",
          "preventMetricAggregates",
          "excludeOther"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": null,
      "related": []
    },
    {
      "id": "Create an External User",
      "method": "POST",
      "path": "/api/0/organizations/{organization_id_or_slug}/external-users/",
      "tag": "Integrations",
      "summary": "Link a user from an external provider to a Sentry user.",
      "description": "Link a user from an external provider to a Sentry user.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        }
      ],
      "hasRequestBody": true,
      "risk": "write",
      "cacheable": false,
      "responseCodes": [
        "200",
        "201",
        "400",
        "403"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug"
        ],
        "query": [],
        "headers": [],
        "body": true
      },
      "graphql": null,
      "kind": "create",
      "pagination": null,
      "related": []
    },
    {
      "id": "Update an External User",
      "method": "PUT",
      "path": "/api/0/organizations/{organization_id_or_slug}/external-users/{external_user_id}/",
      "tag": "Integrations",
      "summary": "Update a user in an external provider that is currently linked to a Sentry user.",
      "description": "Update a user in an external provider that is currently linked to a Sentry user.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "external_user_id",
          "in": "path",
          "required": true,
          "description": "The ID of the external user object. This is returned when creating an external user.",
          "type": "integer"
        }
      ],
      "hasRequestBody": true,
      "risk": "write",
      "cacheable": false,
      "responseCodes": [
        "200",
        "400",
        "403"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "external_user_id"
        ],
        "query": [],
        "headers": [],
        "body": true
      },
      "graphql": null,
      "kind": "update",
      "pagination": null,
      "related": [
        {
          "id": "Delete an External User",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Delete an External User",
      "method": "DELETE",
      "path": "/api/0/organizations/{organization_id_or_slug}/external-users/{external_user_id}/",
      "tag": "Integrations",
      "summary": "Delete the link between a user from an external provider and a Sentry user.",
      "description": "Delete the link between a user from an external provider and a Sentry user.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "external_user_id",
          "in": "path",
          "required": true,
          "description": "The ID of the external user object. This is returned when creating an external user.",
          "type": "integer"
        }
      ],
      "hasRequestBody": false,
      "risk": "destructive",
      "cacheable": false,
      "responseCodes": [
        "204",
        "400",
        "403"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "external_user_id"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "delete",
      "pagination": null,
      "related": [
        {
          "id": "Update an External User",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Retrieve Data Forwarders for an Organization",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/forwarding/",
      "tag": "Integrations",
      "summary": "Returns a list of data forwarders for an organization.",
      "description": "Returns a list of data forwarders for an organization.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": null,
      "related": [
        {
          "id": "Create a Data Forwarder for an Organization",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Create a Data Forwarder for an Organization",
      "method": "POST",
      "path": "/api/0/organizations/{organization_id_or_slug}/forwarding/",
      "tag": "Integrations",
      "summary": "Creates a new data forwarder for an organization.\nOnly one data forwarder can be created per provider for a given organization.\n\nProject-specific overrides can only be created after creating the data forwarder.",
      "description": "Creates a new data forwarder for an organization.\nOnly one data forwarder can be created per provider for a given organization.\n\nProject-specific overrides can only be created after creating the data forwarder.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        }
      ],
      "hasRequestBody": true,
      "risk": "write",
      "cacheable": false,
      "responseCodes": [
        "201",
        "400",
        "403"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug"
        ],
        "query": [],
        "headers": [],
        "body": true
      },
      "graphql": null,
      "kind": "create",
      "pagination": null,
      "related": [
        {
          "id": "Retrieve Data Forwarders for an Organization",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Update a Data Forwarder for an Organization",
      "method": "PUT",
      "path": "/api/0/organizations/{organization_id_or_slug}/forwarding/{data_forwarder_id}/",
      "tag": "Integrations",
      "summary": "Updates a data forwarder for an organization or update a project-specific override.\nUpdates to the data forwarder's configuration require ` + "`" + `org:write` + "`" + ` permissions, and the entire\nconfiguration to be provided, including the ` + "`" + `project_ids` + "`" + ` field.\n\nTo configure project-specific overrides, specify only the following fields:\n\n  - 'project_id': The ID of the project to create/modify the override for.\n  - 'overrides': Follows the same format as ` + "`" + `config` + "`" + ` but all provider fields are optional, since only specified fields are overridden.\n  - 'is_enabled': To enable/disable the forwarder for events on the specific project.\n\nOverrides can be performed with ` + "`" + `project:write` + "`" + ` permissions on the project being modified.",
      "description": "Updates a data forwarder for an organization or update a project-specific override.\nUpdates to the data forwarder's configuration require ` + "`" + `org:write` + "`" + ` permissions, and the entire\nconfiguration to be provided, including the ` + "`" + `project_ids` + "`" + ` field.\n\nTo configure project-specific overrides, specify only the following fields:\n\n  - 'project_id': The ID of the project to create/modify the override for.\n  - 'overrides': Follows the same format as ` + "`" + `config` + "`" + ` but all provider fields are optional, since only specified fields are overridden.\n  - 'is_enabled': To enable/disable the forwarder for events on the specific project.\n\nOverrides can be performed with ` + "`" + `project:write` + "`" + ` permissions on the project being modified.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "data_forwarder_id",
          "in": "path",
          "required": true,
          "description": "The ID of the data forwarder you'd like to query.",
          "type": "integer"
        }
      ],
      "hasRequestBody": true,
      "risk": "write",
      "cacheable": false,
      "responseCodes": [
        "200",
        "400",
        "403"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "data_forwarder_id"
        ],
        "query": [],
        "headers": [],
        "body": true
      },
      "graphql": null,
      "kind": "update",
      "pagination": null,
      "related": [
        {
          "id": "Delete a Data Forwarder for an Organization",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Delete a Data Forwarder for an Organization",
      "method": "DELETE",
      "path": "/api/0/organizations/{organization_id_or_slug}/forwarding/{data_forwarder_id}/",
      "tag": "Integrations",
      "summary": "Deletes a data forwarder for an organization. All project-specific overrides will be deleted as well.",
      "description": "Deletes a data forwarder for an organization. All project-specific overrides will be deleted as well.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "data_forwarder_id",
          "in": "path",
          "required": true,
          "description": "The ID of the data forwarder you'd like to query.",
          "type": "integer"
        }
      ],
      "hasRequestBody": false,
      "risk": "destructive",
      "cacheable": false,
      "responseCodes": [
        "204",
        "403"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "data_forwarder_id"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "delete",
      "pagination": null,
      "related": [
        {
          "id": "Update a Data Forwarder for an Organization",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "List an Organization's Available Integrations",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/integrations/",
      "tag": "Integrations",
      "summary": "Lists all the available Integrations for an Organization.",
      "description": "Lists all the available Integrations for an Organization.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "providerKey",
          "in": "query",
          "required": false,
          "description": "Specific integration provider to filter by such as ` + "`" + `slack` + "`" + `. See our [Integrations Documentation](/product/integrations/) for an updated list of providers.",
          "type": "string"
        },
        {
          "name": "features",
          "in": "query",
          "required": false,
          "description": "Integration features to filter by. See our [Integrations Documentation](/product/integrations/) for an updated list of features. Current available ones are:\n- ` + "`" + `alert-rule` + "`" + `\n- ` + "`" + `chat-unfurl` + "`" + `\n- ` + "`" + `codeowners` + "`" + `\n- ` + "`" + `commits` + "`" + `\n- ` + "`" + `data-forwarding` + "`" + `\n- ` + "`" + `deployment` + "`" + `\n- ` + "`" + `enterprise-alert-rule` + "`" + `\n- ` + "`" + `enterprise-incident-management` + "`" + `\n- ` + "`" + `incident-management` + "`" + `\n- ` + "`" + `issue-basic` + "`" + `\n- ` + "`" + `issue-sync` + "`" + `\n- ` + "`" + `mobile` + "`" + `\n- ` + "`" + `serverless` + "`" + `\n- ` + "`" + `session-replay` + "`" + `\n- ` + "`" + `stacktrace-link` + "`" + `\n- ` + "`" + `ticket-rules` + "`" + `\n    ",
          "type": "array"
        },
        {
          "name": "includeConfig",
          "in": "query",
          "required": false,
          "description": "Specify ` + "`" + `True` + "`" + ` to fetch third-party integration configurations. Note that this can add several seconds to the response time.",
          "type": "boolean"
        },
        {
          "name": "cursor",
          "in": "query",
          "required": false,
          "description": "A pointer to the last object fetched and its sort order; used to retrieve the next or previous results.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug"
        ],
        "query": [
          "providerKey",
          "features",
          "includeConfig",
          "cursor"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": {
        "style": "cursor",
        "cursorParam": "cursor"
      },
      "related": []
    },
    {
      "id": "Retrieve an Integration for an Organization",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/integrations/{integration_id}/",
      "tag": "Integrations",
      "summary": "OrganizationIntegrationBaseEndpoints expect both Integration and\nOrganizationIntegration DB entries to exist for a given organization and\nintegration_id.",
      "description": "OrganizationIntegrationBaseEndpoints expect both Integration and\nOrganizationIntegration DB entries to exist for a given organization and\nintegration_id.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "integration_id",
          "in": "path",
          "required": true,
          "description": "The ID of the integration installed on the organization.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "integration_id"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": null,
      "related": [
        {
          "id": "Delete an Integration for an Organization",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Delete an Integration for an Organization",
      "method": "DELETE",
      "path": "/api/0/organizations/{organization_id_or_slug}/integrations/{integration_id}/",
      "tag": "Integrations",
      "summary": "OrganizationIntegrationBaseEndpoints expect both Integration and\nOrganizationIntegration DB entries to exist for a given organization and\nintegration_id.",
      "description": "OrganizationIntegrationBaseEndpoints expect both Integration and\nOrganizationIntegration DB entries to exist for a given organization and\nintegration_id.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "integration_id",
          "in": "path",
          "required": true,
          "description": "The ID of the integration installed on the organization.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "destructive",
      "cacheable": false,
      "responseCodes": [
        "204",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "integration_id"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "delete",
      "pagination": null,
      "related": [
        {
          "id": "Retrieve an Integration for an Organization",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "List an Organization's Issues",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/issues/",
      "tag": "Events",
      "summary": "Return a list of issues for an organization. All parameters are supplied as query string parameters. A default query of ` + "`" + `is:unresolved` + "`" + ` is applied. To return all results, use an empty query value (i.e. ` + "`" + `` + "`" + `?query=` + "`" + `). ",
      "description": "Return a list of issues for an organization. All parameters are supplied as query string parameters. A default query of ` + "`" + `is:unresolved` + "`" + ` is applied. To return all results, use an empty query value (i.e. ` + "`" + `` + "`" + `?query=` + "`" + `). ",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "environment",
          "in": "query",
          "required": false,
          "description": "The name of environments to filter by.",
          "type": "array"
        },
        {
          "name": "project",
          "in": "query",
          "required": false,
          "description": "The IDs of projects to filter by. ` + "`" + `-1` + "`" + ` means all available projects.\nFor example, the following are valid parameters:\n- ` + "`" + `/?project=1234&project=56789` + "`" + `\n- ` + "`" + `/?project=-1` + "`" + `\n",
          "type": "array"
        },
        {
          "name": "statsPeriod",
          "in": "query",
          "required": false,
          "description": "The period of time for the query, will override the start & end parameters, a number followed by one of:\n- ` + "`" + `d` + "`" + ` for days\n- ` + "`" + `h` + "`" + ` for hours\n- ` + "`" + `m` + "`" + ` for minutes\n- ` + "`" + `s` + "`" + ` for seconds\n- ` + "`" + `w` + "`" + ` for weeks\n\nFor example, ` + "`" + `24h` + "`" + `, to mean query data starting from 24 hours ago to now.",
          "type": "string"
        },
        {
          "name": "start",
          "in": "query",
          "required": false,
          "description": "The start of the period of time for the query, expected in ISO-8601 format. For example, ` + "`" + `2001-12-14T12:34:56.7890` + "`" + `.",
          "type": "string"
        },
        {
          "name": "end",
          "in": "query",
          "required": false,
          "description": "The end of the period of time for the query, expected in ISO-8601 format. For example, ` + "`" + `2001-12-14T12:34:56.7890` + "`" + `.",
          "type": "string"
        },
        {
          "name": "groupStatsPeriod",
          "in": "query",
          "required": false,
          "description": "The timeline on which stats for the groups should be presented.",
          "type": "string"
        },
        {
          "name": "shortIdLookup",
          "in": "query",
          "required": false,
          "description": "If this is set to ` + "`" + `1` + "`" + ` then the query will be parsed for issue short IDs. These may ignore other filters (e.g. projects), which is why it is an opt-in.",
          "type": "string"
        },
        {
          "name": "query",
          "in": "query",
          "required": false,
          "description": "An optional search query for filtering issues. A default query will apply if no view/query is set. For all results use this parameter with an empty string.",
          "type": "string"
        },
        {
          "name": "viewId",
          "in": "query",
          "required": false,
          "description": "The ID of the view to use. If no query is present, the view's query and filters will be applied.",
          "type": "string"
        },
        {
          "name": "sort",
          "in": "query",
          "required": false,
          "description": "The sort order of the view. Options include 'Last Seen' (` + "`" + `date` + "`" + `), 'First Seen' (` + "`" + `new` + "`" + `), 'Trends' (` + "`" + `trends` + "`" + `), 'Events' (` + "`" + `freq` + "`" + `), 'Users' (` + "`" + `user` + "`" + `), 'Date Added' (` + "`" + `inbox` + "`" + `), and 'Recommended' (` + "`" + `recommended` + "`" + `).",
          "type": "string"
        },
        {
          "name": "limit",
          "in": "query",
          "required": false,
          "description": "The maximum number of issues to affect. The maximum is 100.",
          "type": "integer"
        },
        {
          "name": "expand",
          "in": "query",
          "required": false,
          "description": "Additional data to include in the response.",
          "type": "array"
        },
        {
          "name": "collapse",
          "in": "query",
          "required": false,
          "description": "Fields to remove from the response to improve query performance.",
          "type": "array"
        },
        {
          "name": "cursor",
          "in": "query",
          "required": false,
          "description": "A pointer to the last object fetched and its sort order; used to retrieve the next or previous results.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "400",
        "401",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug"
        ],
        "query": [
          "environment",
          "project",
          "statsPeriod",
          "start",
          "end",
          "groupStatsPeriod",
          "shortIdLookup",
          "query",
          "viewId",
          "sort",
          "limit",
          "expand",
          "collapse",
          "cursor"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "search",
      "pagination": {
        "style": "cursor",
        "cursorParam": "cursor"
      },
      "related": [
        {
          "id": "Bulk Mutate an Organization's Issues",
          "role": "same-resource"
        },
        {
          "id": "Bulk Remove an Organization's Issues",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Bulk Mutate an Organization's Issues",
      "method": "PUT",
      "path": "/api/0/organizations/{organization_id_or_slug}/issues/",
      "tag": "Events",
      "summary": "Bulk mutate various attributes on a maxmimum of 1000 issues. \n- For non-status updates, the ` + "`" + `id` + "`" + ` query parameter is required. \n- For status updates, the ` + "`" + `id` + "`" + ` query parameter may be omitted to update issues that match the filtering. \nIf any IDs are out of scope, the data won't be mutated but the endpoint will still produce a successful response. For example, if no issues were found matching the criteria, a HTTP 204 is returned.",
      "description": "Bulk mutate various attributes on a maxmimum of 1000 issues. \n- For non-status updates, the ` + "`" + `id` + "`" + ` query parameter is required. \n- For status updates, the ` + "`" + `id` + "`" + ` query parameter may be omitted to update issues that match the filtering. \nIf any IDs are out of scope, the data won't be mutated but the endpoint will still produce a successful response. For example, if no issues were found matching the criteria, a HTTP 204 is returned.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "environment",
          "in": "query",
          "required": false,
          "description": "The name of environments to filter by.",
          "type": "array"
        },
        {
          "name": "project",
          "in": "query",
          "required": false,
          "description": "The IDs of projects to filter by. ` + "`" + `-1` + "`" + ` means all available projects.\nFor example, the following are valid parameters:\n- ` + "`" + `/?project=1234&project=56789` + "`" + `\n- ` + "`" + `/?project=-1` + "`" + `\n",
          "type": "array"
        },
        {
          "name": "id",
          "in": "query",
          "required": false,
          "description": "The list of issue IDs to mutate. It is optional for status updates, in which an implicit ` + "`" + `update all` + "`" + ` is assumed.",
          "type": "array"
        },
        {
          "name": "query",
          "in": "query",
          "required": false,
          "description": "An optional search query for filtering issues. A default query will apply if no view/query is set. For all results use this parameter with an empty string.",
          "type": "string"
        },
        {
          "name": "viewId",
          "in": "query",
          "required": false,
          "description": "The ID of the view to use. If no query is present, the view's query and filters will be applied.",
          "type": "string"
        },
        {
          "name": "sort",
          "in": "query",
          "required": false,
          "description": "The sort order of the view. Options include 'Last Seen' (` + "`" + `date` + "`" + `), 'First Seen' (` + "`" + `new` + "`" + `), 'Trends' (` + "`" + `trends` + "`" + `), 'Events' (` + "`" + `freq` + "`" + `), 'Users' (` + "`" + `user` + "`" + `), 'Date Added' (` + "`" + `inbox` + "`" + `), and 'Recommended' (` + "`" + `recommended` + "`" + `).",
          "type": "string"
        },
        {
          "name": "limit",
          "in": "query",
          "required": false,
          "description": "The maximum number of issues to affect. The maximum is 100.",
          "type": "integer"
        }
      ],
      "hasRequestBody": true,
      "risk": "write",
      "cacheable": false,
      "responseCodes": [
        "200",
        "204",
        "400",
        "401",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug"
        ],
        "query": [
          "environment",
          "project",
          "id",
          "query",
          "viewId",
          "sort",
          "limit"
        ],
        "headers": [],
        "body": true
      },
      "graphql": null,
      "kind": "update",
      "pagination": null,
      "related": [
        {
          "id": "List an Organization's Issues",
          "role": "same-resource"
        },
        {
          "id": "Bulk Remove an Organization's Issues",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Bulk Remove an Organization's Issues",
      "method": "DELETE",
      "path": "/api/0/organizations/{organization_id_or_slug}/issues/",
      "tag": "Events",
      "summary": "Permanently remove the given issues. If IDs are provided, queries and filtering will be ignored. If any IDs are out of scope, the data won't be mutated but the endpoint will still produce a successful response. For example, if no issues were found matching the criteria, a HTTP 204 is returned.",
      "description": "Permanently remove the given issues. If IDs are provided, queries and filtering will be ignored. If any IDs are out of scope, the data won't be mutated but the endpoint will still produce a successful response. For example, if no issues were found matching the criteria, a HTTP 204 is returned.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "environment",
          "in": "query",
          "required": false,
          "description": "The name of environments to filter by.",
          "type": "array"
        },
        {
          "name": "project",
          "in": "query",
          "required": false,
          "description": "The IDs of projects to filter by. ` + "`" + `-1` + "`" + ` means all available projects.\nFor example, the following are valid parameters:\n- ` + "`" + `/?project=1234&project=56789` + "`" + `\n- ` + "`" + `/?project=-1` + "`" + `\n",
          "type": "array"
        },
        {
          "name": "id",
          "in": "query",
          "required": false,
          "description": "The list of issue IDs to be removed. If not provided, it will attempt to remove the first 1000 issues.",
          "type": "array"
        },
        {
          "name": "query",
          "in": "query",
          "required": false,
          "description": "An optional search query for filtering issues. A default query will apply if no view/query is set. For all results use this parameter with an empty string.",
          "type": "string"
        },
        {
          "name": "viewId",
          "in": "query",
          "required": false,
          "description": "The ID of the view to use. If no query is present, the view's query and filters will be applied.",
          "type": "string"
        },
        {
          "name": "sort",
          "in": "query",
          "required": false,
          "description": "The sort order of the view. Options include 'Last Seen' (` + "`" + `date` + "`" + `), 'First Seen' (` + "`" + `new` + "`" + `), 'Trends' (` + "`" + `trends` + "`" + `), 'Events' (` + "`" + `freq` + "`" + `), 'Users' (` + "`" + `user` + "`" + `), 'Date Added' (` + "`" + `inbox` + "`" + `), and 'Recommended' (` + "`" + `recommended` + "`" + `).",
          "type": "string"
        },
        {
          "name": "limit",
          "in": "query",
          "required": false,
          "description": "The maximum number of issues to affect. The maximum is 100.",
          "type": "integer"
        }
      ],
      "hasRequestBody": false,
      "risk": "destructive",
      "cacheable": false,
      "responseCodes": [
        "204",
        "400",
        "401",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug"
        ],
        "query": [
          "environment",
          "project",
          "id",
          "query",
          "viewId",
          "sort",
          "limit"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "delete",
      "pagination": null,
      "related": [
        {
          "id": "List an Organization's Issues",
          "role": "same-resource"
        },
        {
          "id": "Bulk Mutate an Organization's Issues",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "List an Organization's Members",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/members/",
      "tag": "Organizations",
      "summary": "List all organization members.\n\nResponse includes pending invites that are approved by organization owners or managers but waiting to be accepted by the invitee.",
      "description": "List all organization members.\n\nResponse includes pending invites that are approved by organization owners or managers but waiting to be accepted by the invitee.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "cursor",
          "in": "query",
          "required": false,
          "description": "A pointer to the last object fetched and its sort order; used to retrieve the next or previous results.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "401",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug"
        ],
        "query": [
          "cursor"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": {
        "style": "cursor",
        "cursorParam": "cursor"
      },
      "related": [
        {
          "id": "Add a Member to an Organization",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Add a Member to an Organization",
      "method": "POST",
      "path": "/api/0/organizations/{organization_id_or_slug}/members/",
      "tag": "Organizations",
      "summary": "Add or invite a member to an organization.",
      "description": "Add or invite a member to an organization.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        }
      ],
      "hasRequestBody": true,
      "risk": "write",
      "cacheable": false,
      "responseCodes": [
        "201",
        "401",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug"
        ],
        "query": [],
        "headers": [],
        "body": true
      },
      "graphql": null,
      "kind": "create",
      "pagination": null,
      "related": [
        {
          "id": "List an Organization's Members",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Retrieve an Organization Member",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/members/{member_id}/",
      "tag": "Organizations",
      "summary": "Retrieve an organization member's details.\n\nResponse will be a pending invite if it has been approved by organization owners or managers but is waiting to be accepted by the invitee.",
      "description": "Retrieve an organization member's details.\n\nResponse will be a pending invite if it has been approved by organization owners or managers but is waiting to be accepted by the invitee.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "member_id",
          "in": "path",
          "required": true,
          "description": "The ID of the organization member.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "401",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "member_id"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": null,
      "related": [
        {
          "id": "Update an Organization Member's Roles",
          "role": "same-resource"
        },
        {
          "id": "Delete an Organization Member",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Update an Organization Member's Roles",
      "method": "PUT",
      "path": "/api/0/organizations/{organization_id_or_slug}/members/{member_id}/",
      "tag": "Organizations",
      "summary": "Update a member's [organization-level](https://docs.sentry.io/organization/membership/#organization-level-roles) and [team-level](https://docs.sentry.io/organization/membership/#team-level-roles) roles.\n\nNote that for changing organization-roles, this endpoint is restricted to\n[user auth tokens](https://docs.sentry.io/account/auth-tokens/#user-auth-tokens).\nAdditionally, both the original and desired organization role must have\nthe same or lower permissions than the role of the organization user making the request\n\nFor example, an organization Manager may change someone's role from\nMember to Manager, but not to Owner.",
      "description": "Update a member's [organization-level](https://docs.sentry.io/organization/membership/#organization-level-roles) and [team-level](https://docs.sentry.io/organization/membership/#team-level-roles) roles.\n\nNote that for changing organization-roles, this endpoint is restricted to\n[user auth tokens](https://docs.sentry.io/account/auth-tokens/#user-auth-tokens).\nAdditionally, both the original and desired organization role must have\nthe same or lower permissions than the role of the organization user making the request\n\nFor example, an organization Manager may change someone's role from\nMember to Manager, but not to Owner.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "member_id",
          "in": "path",
          "required": true,
          "description": "The ID of the member to update.",
          "type": "string"
        }
      ],
      "hasRequestBody": true,
      "risk": "write",
      "cacheable": false,
      "responseCodes": [
        "200",
        "400",
        "401",
        "403"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "member_id"
        ],
        "query": [],
        "headers": [],
        "body": true
      },
      "graphql": null,
      "kind": "update",
      "pagination": null,
      "related": [
        {
          "id": "Retrieve an Organization Member",
          "role": "same-resource"
        },
        {
          "id": "Delete an Organization Member",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Delete an Organization Member",
      "method": "DELETE",
      "path": "/api/0/organizations/{organization_id_or_slug}/members/{member_id}/",
      "tag": "Organizations",
      "summary": "Remove an organization member.",
      "description": "Remove an organization member.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "member_id",
          "in": "path",
          "required": true,
          "description": "The ID of the member to delete.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "destructive",
      "cacheable": false,
      "responseCodes": [
        "204",
        "401",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "member_id"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "delete",
      "pagination": null,
      "related": [
        {
          "id": "Retrieve an Organization Member",
          "role": "same-resource"
        },
        {
          "id": "Update an Organization Member's Roles",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Add an Organization Member to a Team",
      "method": "POST",
      "path": "/api/0/organizations/{organization_id_or_slug}/members/{member_id}/teams/{team_id_or_slug}/",
      "tag": "Teams",
      "summary": "This request can return various success codes depending on the context of the team:\n- **` + "`" + `201` + "`" + `**: The member has been successfully added.\n- **` + "`" + `202` + "`" + `**: The member needs permission to join the team and an access request\nhas been generated.\n- **` + "`" + `204` + "`" + `**: The member is already on the team.\n\nIf the team is provisioned through an identity provider, the member cannot join the\nteam through Sentry.\n\nNote the permission scopes vary depending on the organization setting ` + "`" + `\"Open Membership\"` + "`" + `\nand the type of authorization token. The following table outlines the accepted scopes.\n<table style=\"width: 100%;\">\n<thead>\n    <tr>\n    <th style=\"width: 33%;\"></th>\n    <th colspan=\"2\" style=\"text-align: center; font-weight: bold; width: 33%;\">Open Membership</th>\n    </tr>\n</thead>\n<tbody>\n    <tr>\n    <td style=\"width: 34%;\"></td>\n    <td style=\"text-align: center; font-weight: bold; width: 33%;\">On</td>\n    <td style=\"text-align: center; font-weight: bold; width: 33%;\">Off</td>\n    </tr>\n    <tr>\n    <td style=\"text-align: center; font-weight: bold; vertical-align: middle;\"><a\n    href=\"https://docs.sentry.io/account/auth-tokens/#internal-integrations\">Internal Integration Token</a></td>\n    <td style=\"text-align: left; width: 33%;\">\n        <ul style=\"list-style-type: none; padding-left: 0;\">\n        <li><strong style=\"color: #9c5f99;\">&bull; org:read</strong></li>\n        </ul>\n    </td>\n    <td style=\"text-align: left; width: 33%;\">\n        <ul style=\"list-style-type: none; padding-left: 0;\">\n        <li><strong style=\"color: #9c5f99;\">&bull; org:write</strong></li>\n        <li><strong style=\"color: #9c5f99;\">&bull; team:write</strong></li>\n        </ul>\n    </td>\n    </tr>\n    <tr>\n    <td style=\"text-align: center; font-weight: bold; vertical-align: middle;\"><a\n    href=\"https://docs.sentry.io/account/auth-tokens/#user-auth-tokens\">User Auth Token</a></td>\n    <td style=\"text-align: left; width: 33%;\">\n        <ul style=\"list-style-type: none; padding-left: 0;\">\n        <li><strong style=\"color: #9c5f99;\">&bull; org:read</strong></li>\n        </ul>\n    </td>\n    <td style=\"text-align: left; width: 33%;\">\n        <ul style=\"list-style-type: none; padding-left: 0;\">\n        <li><strong style=\"color: #9c5f99;\">&bull; org:read*</strong></li>\n        <li><strong style=\"color: #9c5f99;\">&bull; org:write</strong></li>\n        <li><strong style=\"color: #9c5f99;\">&bull; org:read +</strong></li>\n        <li><strong style=\"color: #9c5f99;\">&nbsp; &nbsp;team:write**</strong></li>\n        </ul>\n    </td>\n    </tr>\n</tbody>\n</table>\n\n\n*Organization members are restricted to this scope. When sending a request, it will always\nreturn a 202 and request an invite to the team.\n\n\n\\*\\*Team Admins must have both **` + "`" + `org:read` + "`" + `** and **` + "`" + `team:write` + "`" + `** scopes in their user\nauthorization token to add members to their teams.",
      "description": "This request can return various success codes depending on the context of the team:\n- **` + "`" + `201` + "`" + `**: The member has been successfully added.\n- **` + "`" + `202` + "`" + `**: The member needs permission to join the team and an access request\nhas been generated.\n- **` + "`" + `204` + "`" + `**: The member is already on the team.\n\nIf the team is provisioned through an identity provider, the member cannot join the\nteam through Sentry.\n\nNote the permission scopes vary depending on the organization setting ` + "`" + `\"Open Membership\"` + "`" + `\nand the type of authorization token. The following table outlines the accepted scopes.\n<table style=\"width: 100%;\">\n<thead>\n    <tr>\n    <th style=\"width: 33%;\"></th>\n    <th colspan=\"2\" style=\"text-align: center; font-weight: bold; width: 33%;\">Open Membership</th>\n    </tr>\n</thead>\n<tbody>\n    <tr>\n    <td style=\"width: 34%;\"></td>\n    <td style=\"text-align: center; font-weight: bold; width: 33%;\">On</td>\n    <td style=\"text-align: center; font-weight: bold; width: 33%;\">Off</td>\n    </tr>\n    <tr>\n    <td style=\"text-align: center; font-weight: bold; vertical-align: middle;\"><a\n    href=\"https://docs.sentry.io/account/auth-tokens/#internal-integrations\">Internal Integration Token</a></td>\n    <td style=\"text-align: left; width: 33%;\">\n        <ul style=\"list-style-type: none; padding-left: 0;\">\n        <li><strong style=\"color: #9c5f99;\">&bull; org:read</strong></li>\n        </ul>\n    </td>\n    <td style=\"text-align: left; width: 33%;\">\n        <ul style=\"list-style-type: none; padding-left: 0;\">\n        <li><strong style=\"color: #9c5f99;\">&bull; org:write</strong></li>\n        <li><strong style=\"color: #9c5f99;\">&bull; team:write</strong></li>\n        </ul>\n    </td>\n    </tr>\n    <tr>\n    <td style=\"text-align: center; font-weight: bold; vertical-align: middle;\"><a\n    href=\"https://docs.sentry.io/account/auth-tokens/#user-auth-tokens\">User Auth Token</a></td>\n    <td style=\"text-align: left; width: 33%;\">\n        <ul style=\"list-style-type: none; padding-left: 0;\">\n        <li><strong style=\"color: #9c5f99;\">&bull; org:read</strong></li>\n        </ul>\n    </td>\n    <td style=\"text-align: left; width: 33%;\">\n        <ul style=\"list-style-type: none; padding-left: 0;\">\n        <li><strong style=\"color: #9c5f99;\">&bull; org:read*</strong></li>\n        <li><strong style=\"color: #9c5f99;\">&bull; org:write</strong></li>\n        <li><strong style=\"color: #9c5f99;\">&bull; org:read +</strong></li>\n        <li><strong style=\"color: #9c5f99;\">&nbsp; &nbsp;team:write**</strong></li>\n        </ul>\n    </td>\n    </tr>\n</tbody>\n</table>\n\n\n*Organization members are restricted to this scope. When sending a request, it will always\nreturn a 202 and request an invite to the team.\n\n\n\\*\\*Team Admins must have both **` + "`" + `org:read` + "`" + `** and **` + "`" + `team:write` + "`" + `** scopes in their user\nauthorization token to add members to their teams.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "member_id",
          "in": "path",
          "required": true,
          "description": "The ID of the organization member to add to the team",
          "type": "string"
        },
        {
          "name": "team_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the team the resource belongs to.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "write",
      "cacheable": false,
      "responseCodes": [
        "201",
        "202",
        "204",
        "401",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "member_id",
          "team_id_or_slug"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "create",
      "pagination": null,
      "related": [
        {
          "id": "Update an Organization Member's Team Role",
          "role": "same-resource"
        },
        {
          "id": "Delete an Organization Member from a Team",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Update an Organization Member's Team Role",
      "method": "PUT",
      "path": "/api/0/organizations/{organization_id_or_slug}/members/{member_id}/teams/{team_id_or_slug}/",
      "tag": "Teams",
      "summary": "The relevant organization member must already be a part of the team.\n\nNote that for organization admins, managers, and owners, they are\nautomatically granted a minimum team role of ` + "`" + `admin` + "`" + ` on all teams they\nare part of. Read more about [team roles](https://docs.sentry.io/product/teams/roles/).",
      "description": "The relevant organization member must already be a part of the team.\n\nNote that for organization admins, managers, and owners, they are\nautomatically granted a minimum team role of ` + "`" + `admin` + "`" + ` on all teams they\nare part of. Read more about [team roles](https://docs.sentry.io/product/teams/roles/).",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "member_id",
          "in": "path",
          "required": true,
          "description": "The ID of the organization member to change",
          "type": "string"
        },
        {
          "name": "team_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the team the resource belongs to.",
          "type": "string"
        }
      ],
      "hasRequestBody": true,
      "risk": "write",
      "cacheable": false,
      "responseCodes": [
        "200",
        "400",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "member_id",
          "team_id_or_slug"
        ],
        "query": [],
        "headers": [],
        "body": true
      },
      "graphql": null,
      "kind": "update",
      "pagination": null,
      "related": [
        {
          "id": "Add an Organization Member to a Team",
          "role": "same-resource"
        },
        {
          "id": "Delete an Organization Member from a Team",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Delete an Organization Member from a Team",
      "method": "DELETE",
      "path": "/api/0/organizations/{organization_id_or_slug}/members/{member_id}/teams/{team_id_or_slug}/",
      "tag": "Teams",
      "summary": "Delete an organization member from a team.\n\nNote the permission scopes vary depending on the type of authorization token. The following\ntable outlines the accepted scopes.\n<table style=\"width: 100%;\">\n    <tr style=\"width: 50%;\">\n        <td style=\"width: 50%; text-align: center; font-weight: bold; vertical-align: middle;\"><a href=\"https://docs.sentry.io/api/auth/#auth-tokens\">Org Auth Token</a></td>\n        <td style=\"width: 50%; text-align: left;\">\n            <ul style=\"list-style-type: none; padding-left: 0;\">\n                <li><strong style=\"color: #9c5f99;\">&bull; org:write</strong></li>\n                <li><strong style=\"color: #9c5f99;\">&bull; org:admin</strong></li>\n                <li><strong style=\"color: #9c5f99;\">&bull; team:admin</strong></li>\n            </ul>\n        </td>\n    </tr>\n    <tr style=\"width: 50%;\">\n        <td style=\"width: 50%; text-align: center; font-weight: bold; vertical-align: middle;\"><a href=\"https://docs.sentry.io/api/auth/#user-authentication-tokens\">User Auth Token</a></td>\n        <td style=\"width: 50%; text-align: left;\">\n            <ul style=\"list-style-type: none; padding-left: 0;\">\n                <li><strong style=\"color: #9c5f99;\">&bull; org:read*</strong></li>\n                <li><strong style=\"color: #9c5f99;\">&bull; org:write</strong></li>\n                <li><strong style=\"color: #9c5f99;\">&bull; org:admin</strong></li>\n                <li><strong style=\"color: #9c5f99;\">&bull; team:admin</strong></li>\n                <li><strong style=\"color: #9c5f99;\">&bull; org:read + team:admin**</strong></li>\n            </ul>\n        </td>\n    </tr>\n</table>\n\n\n\\***` + "`" + `org:read` + "`" + `** can only be used to remove yourself from the teams you are a member of.\n\n\n\\*\\*Team Admins must have both **` + "`" + `org:read` + "`" + `** and **` + "`" + `team:admin` + "`" + `** scopes in their user\nauthorization token to delete members from their teams.",
      "description": "Delete an organization member from a team.\n\nNote the permission scopes vary depending on the type of authorization token. The following\ntable outlines the accepted scopes.\n<table style=\"width: 100%;\">\n    <tr style=\"width: 50%;\">\n        <td style=\"width: 50%; text-align: center; font-weight: bold; vertical-align: middle;\"><a href=\"https://docs.sentry.io/api/auth/#auth-tokens\">Org Auth Token</a></td>\n        <td style=\"width: 50%; text-align: left;\">\n            <ul style=\"list-style-type: none; padding-left: 0;\">\n                <li><strong style=\"color: #9c5f99;\">&bull; org:write</strong></li>\n                <li><strong style=\"color: #9c5f99;\">&bull; org:admin</strong></li>\n                <li><strong style=\"color: #9c5f99;\">&bull; team:admin</strong></li>\n            </ul>\n        </td>\n    </tr>\n    <tr style=\"width: 50%;\">\n        <td style=\"width: 50%; text-align: center; font-weight: bold; vertical-align: middle;\"><a href=\"https://docs.sentry.io/api/auth/#user-authentication-tokens\">User Auth Token</a></td>\n        <td style=\"width: 50%; text-align: left;\">\n            <ul style=\"list-style-type: none; padding-left: 0;\">\n                <li><strong style=\"color: #9c5f99;\">&bull; org:read*</strong></li>\n                <li><strong style=\"color: #9c5f99;\">&bull; org:write</strong></li>\n                <li><strong style=\"color: #9c5f99;\">&bull; org:admin</strong></li>\n                <li><strong style=\"color: #9c5f99;\">&bull; team:admin</strong></li>\n                <li><strong style=\"color: #9c5f99;\">&bull; org:read + team:admin**</strong></li>\n            </ul>\n        </td>\n    </tr>\n</table>\n\n\n\\***` + "`" + `org:read` + "`" + `** can only be used to remove yourself from the teams you are a member of.\n\n\n\\*\\*Team Admins must have both **` + "`" + `org:read` + "`" + `** and **` + "`" + `team:admin` + "`" + `** scopes in their user\nauthorization token to delete members from their teams.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "member_id",
          "in": "path",
          "required": true,
          "description": "The ID of the organization member to delete from the team",
          "type": "string"
        },
        {
          "name": "team_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the team the resource belongs to.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "destructive",
      "cacheable": false,
      "responseCodes": [
        "200",
        "400",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "member_id",
          "team_id_or_slug"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "delete",
      "pagination": null,
      "related": [
        {
          "id": "Add an Organization Member to a Team",
          "role": "same-resource"
        },
        {
          "id": "Update an Organization Member's Team Role",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Retrieve Monitors for an Organization",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/monitors/",
      "tag": "Crons",
      "summary": "Lists monitors, including nested monitor environments. May be filtered to a project or environment.",
      "description": "Lists monitors, including nested monitor environments. May be filtered to a project or environment.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "project",
          "in": "query",
          "required": false,
          "description": "The IDs of projects to filter by. ` + "`" + `-1` + "`" + ` means all available projects.\nFor example, the following are valid parameters:\n- ` + "`" + `/?project=1234&project=56789` + "`" + `\n- ` + "`" + `/?project=-1` + "`" + `\n",
          "type": "array"
        },
        {
          "name": "environment",
          "in": "query",
          "required": false,
          "description": "The name of environments to filter by.",
          "type": "array"
        },
        {
          "name": "owner",
          "in": "query",
          "required": false,
          "description": "The owner of the monitor, in the format ` + "`" + `user:id` + "`" + ` or ` + "`" + `team:id` + "`" + `. May be specified multiple times.",
          "type": "string"
        },
        {
          "name": "cursor",
          "in": "query",
          "required": false,
          "description": "A pointer to the last object fetched and its sort order; used to retrieve the next or previous results.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "401",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug"
        ],
        "query": [
          "project",
          "environment",
          "owner",
          "cursor"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": {
        "style": "cursor",
        "cursorParam": "cursor"
      },
      "related": [
        {
          "id": "Create a Monitor",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Create a Monitor",
      "method": "POST",
      "path": "/api/0/organizations/{organization_id_or_slug}/monitors/",
      "tag": "Crons",
      "summary": "Create a new monitor.",
      "description": "Create a new monitor.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        }
      ],
      "hasRequestBody": true,
      "risk": "write",
      "cacheable": false,
      "responseCodes": [
        "201",
        "400",
        "401",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug"
        ],
        "query": [],
        "headers": [],
        "body": true
      },
      "graphql": null,
      "kind": "create",
      "pagination": null,
      "related": [
        {
          "id": "Retrieve Monitors for an Organization",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Retrieve a Monitor",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/monitors/{monitor_id_or_slug}/",
      "tag": "Crons",
      "summary": "Retrieves details for a monitor.",
      "description": "Retrieves details for a monitor.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "monitor_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the monitor.",
          "type": "string"
        },
        {
          "name": "environment",
          "in": "query",
          "required": false,
          "description": "The name of environments to filter by.",
          "type": "array"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "401",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "monitor_id_or_slug"
        ],
        "query": [
          "environment"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": null,
      "related": [
        {
          "id": "Update a Monitor",
          "role": "same-resource"
        },
        {
          "id": "Delete a Monitor or Monitor Environments",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Update a Monitor",
      "method": "PUT",
      "path": "/api/0/organizations/{organization_id_or_slug}/monitors/{monitor_id_or_slug}/",
      "tag": "Crons",
      "summary": "Update a monitor.",
      "description": "Update a monitor.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "monitor_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the monitor.",
          "type": "string"
        }
      ],
      "hasRequestBody": true,
      "risk": "write",
      "cacheable": false,
      "responseCodes": [
        "200",
        "400",
        "401",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "monitor_id_or_slug"
        ],
        "query": [],
        "headers": [],
        "body": true
      },
      "graphql": null,
      "kind": "update",
      "pagination": null,
      "related": [
        {
          "id": "Retrieve a Monitor",
          "role": "same-resource"
        },
        {
          "id": "Delete a Monitor or Monitor Environments",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Delete a Monitor or Monitor Environments",
      "method": "DELETE",
      "path": "/api/0/organizations/{organization_id_or_slug}/monitors/{monitor_id_or_slug}/",
      "tag": "Crons",
      "summary": "Delete a monitor or monitor environments.",
      "description": "Delete a monitor or monitor environments.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "monitor_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the monitor.",
          "type": "string"
        },
        {
          "name": "environment",
          "in": "query",
          "required": false,
          "description": "The name of environments to filter by.",
          "type": "array"
        }
      ],
      "hasRequestBody": false,
      "risk": "destructive",
      "cacheable": false,
      "responseCodes": [
        "202",
        "401",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "monitor_id_or_slug"
        ],
        "query": [
          "environment"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "delete",
      "pagination": null,
      "related": [
        {
          "id": "Retrieve a Monitor",
          "role": "same-resource"
        },
        {
          "id": "Update a Monitor",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Retrieve Check-Ins for a Monitor",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/monitors/{monitor_id_or_slug}/checkins/",
      "tag": "Crons",
      "summary": "Retrieve a list of check-ins for a monitor",
      "description": "Retrieve a list of check-ins for a monitor",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "monitor_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the monitor.",
          "type": "string"
        },
        {
          "name": "cursor",
          "in": "query",
          "required": false,
          "description": "A pointer to the last object fetched and its sort order; used to retrieve the next or previous results.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "401",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "monitor_id_or_slug"
        ],
        "query": [
          "cursor"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": {
        "style": "cursor",
        "cursorParam": "cursor"
      },
      "related": []
    },
    {
      "id": "List Spike Protection Notifications",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/notifications/actions/",
      "tag": "Alerts",
      "summary": "Returns all Spike Protection Notification Actions for an organization.\n\nNotification Actions notify a set of members when an action has been triggered through a notification service such as Slack or Sentry.\nFor example, organization owners and managers can receive an email when a spike occurs.\n\nYou can use either the ` + "`" + `project` + "`" + ` or ` + "`" + `projectSlug` + "`" + ` query parameter to filter for certain projects. Note that if both are present, ` + "`" + `projectSlug` + "`" + ` takes priority.",
      "description": "Returns all Spike Protection Notification Actions for an organization.\n\nNotification Actions notify a set of members when an action has been triggered through a notification service such as Slack or Sentry.\nFor example, organization owners and managers can receive an email when a spike occurs.\n\nYou can use either the ` + "`" + `project` + "`" + ` or ` + "`" + `projectSlug` + "`" + ` query parameter to filter for certain projects. Note that if both are present, ` + "`" + `projectSlug` + "`" + ` takes priority.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "project",
          "in": "query",
          "required": false,
          "description": "The IDs of projects to filter by. ` + "`" + `-1` + "`" + ` means all available projects.\nFor example, the following are valid parameters:\n- ` + "`" + `/?project=1234&project=56789` + "`" + `\n- ` + "`" + `/?project=-1` + "`" + `\n",
          "type": "array"
        },
        {
          "name": "project_id_or_slug",
          "in": "query",
          "required": false,
          "description": "The project slugs to filter by. Use ` + "`" + `$all` + "`" + ` to include all available projects. For example, the following are valid parameters:\n- ` + "`" + `/?projectSlug=$all` + "`" + `\n- ` + "`" + `/?projectSlug=android&projectSlug=javascript-react` + "`" + `\n",
          "type": "array"
        },
        {
          "name": "triggerType",
          "in": "query",
          "required": false,
          "description": "Type of the trigger that causes the notification. The only supported value right now is: ` + "`" + `spike-protection` + "`" + `",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "201",
        "400",
        "403"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug"
        ],
        "query": [
          "project",
          "project_id_or_slug",
          "triggerType"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "search",
      "pagination": null,
      "related": [
        {
          "id": "Create a Spike Protection Notification Action",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Create a Spike Protection Notification Action",
      "method": "POST",
      "path": "/api/0/organizations/{organization_id_or_slug}/notifications/actions/",
      "tag": "Alerts",
      "summary": "Creates a new Notification Action for Spike Protection.\n\nNotification Actions notify a set of members when an action has been triggered through a notification service such as Slack or Sentry.\nFor example, organization owners and managers can receive an email when a spike occurs.",
      "description": "Creates a new Notification Action for Spike Protection.\n\nNotification Actions notify a set of members when an action has been triggered through a notification service such as Slack or Sentry.\nFor example, organization owners and managers can receive an email when a spike occurs.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        }
      ],
      "hasRequestBody": true,
      "risk": "write",
      "cacheable": false,
      "responseCodes": [
        "201",
        "400",
        "403"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug"
        ],
        "query": [],
        "headers": [],
        "body": true
      },
      "graphql": null,
      "kind": "create",
      "pagination": null,
      "related": [
        {
          "id": "List Spike Protection Notifications",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Retrieve a Spike Protection Notification Action",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/notifications/actions/{action_id}/",
      "tag": "Alerts",
      "summary": "Returns a serialized Spike Protection Notification Action object.\n\nNotification Actions notify a set of members when an action has been triggered through a notification service such as Slack or Sentry.\nFor example, organization owners and managers can receive an email when a spike occurs.",
      "description": "Returns a serialized Spike Protection Notification Action object.\n\nNotification Actions notify a set of members when an action has been triggered through a notification service such as Slack or Sentry.\nFor example, organization owners and managers can receive an email when a spike occurs.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "action_id",
          "in": "path",
          "required": true,
          "description": "ID of the notification action to retrieve",
          "type": "integer"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "action_id"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": null,
      "related": [
        {
          "id": "Update a Spike Protection Notification Action",
          "role": "same-resource"
        },
        {
          "id": "Delete a Spike Protection Notification Action",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Update a Spike Protection Notification Action",
      "method": "PUT",
      "path": "/api/0/organizations/{organization_id_or_slug}/notifications/actions/{action_id}/",
      "tag": "Alerts",
      "summary": "Updates a Spike Protection Notification Action.\n\nNotification Actions notify a set of members when an action has been triggered through a notification service such as Slack or Sentry.\nFor example, organization owners and managers can receive an email when a spike occurs.",
      "description": "Updates a Spike Protection Notification Action.\n\nNotification Actions notify a set of members when an action has been triggered through a notification service such as Slack or Sentry.\nFor example, organization owners and managers can receive an email when a spike occurs.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "action_id",
          "in": "path",
          "required": true,
          "description": "ID of the notification action to retrieve",
          "type": "integer"
        }
      ],
      "hasRequestBody": true,
      "risk": "write",
      "cacheable": false,
      "responseCodes": [
        "202",
        "400"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "action_id"
        ],
        "query": [],
        "headers": [],
        "body": true
      },
      "graphql": null,
      "kind": "update",
      "pagination": null,
      "related": [
        {
          "id": "Retrieve a Spike Protection Notification Action",
          "role": "same-resource"
        },
        {
          "id": "Delete a Spike Protection Notification Action",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Delete a Spike Protection Notification Action",
      "method": "DELETE",
      "path": "/api/0/organizations/{organization_id_or_slug}/notifications/actions/{action_id}/",
      "tag": "Alerts",
      "summary": "Deletes a Spike Protection Notification Action.\n\nNotification Actions notify a set of members when an action has been triggered through a notification service such as Slack or Sentry.\nFor example, organization owners and managers can receive an email when a spike occurs.",
      "description": "Deletes a Spike Protection Notification Action.\n\nNotification Actions notify a set of members when an action has been triggered through a notification service such as Slack or Sentry.\nFor example, organization owners and managers can receive an email when a spike occurs.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "action_id",
          "in": "path",
          "required": true,
          "description": "ID of the notification action to retrieve",
          "type": "integer"
        }
      ],
      "hasRequestBody": false,
      "risk": "destructive",
      "cacheable": false,
      "responseCodes": [
        "204"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "action_id"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "delete",
      "pagination": null,
      "related": [
        {
          "id": "Retrieve a Spike Protection Notification Action",
          "role": "same-resource"
        },
        {
          "id": "Update a Spike Protection Notification Action",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Retrieve install info for a given artifact",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/preprodartifacts/{artifact_id}/install-details/",
      "tag": "Mobile Builds",
      "summary": "Retrieve install info for a given artifact.\n\nReturns distribution and installation details for a specific preprod artifact,\nincluding whether the artifact is installable, the install URL, download count,\nand iOS-specific code signing information.",
      "description": "Retrieve install info for a given artifact.\n\nReturns distribution and installation details for a specific preprod artifact,\nincluding whether the artifact is installable, the install URL, download count,\nand iOS-specific code signing information.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "artifact_id",
          "in": "path",
          "required": true,
          "description": "The ID of the build artifact.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "artifact_id"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": null,
      "related": []
    },
    {
      "id": "Retrieve Size Analysis results for a given artifact",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/preprodartifacts/{artifact_id}/size-analysis/",
      "tag": "Mobile Builds",
      "summary": "Retrieve size analysis results for a given artifact.\n\nReturns size metrics including download size, install size, and optional insights.\nWhen a base artifact exists (either from commit comparison or via the ` + "`" + `baseArtifactId` + "`" + ` parameter),\ncomparison data showing size differences is included.\n\nThe response ` + "`" + `state` + "`" + ` field indicates the analysis status:\n- ` + "`" + `PENDING` + "`" + `: Analysis has not started yet.\n- ` + "`" + `PROCESSING` + "`" + `: Analysis is currently running.\n- ` + "`" + `FAILED` + "`" + ` / ` + "`" + `NOT_RAN` + "`" + `: Analysis did not complete; ` + "`" + `errorCode` + "`" + ` and ` + "`" + `errorMessage` + "`" + ` are included.\n- ` + "`" + `COMPLETED` + "`" + `: Analysis finished successfully with full size data.",
      "description": "Retrieve size analysis results for a given artifact.\n\nReturns size metrics including download size, install size, and optional insights.\nWhen a base artifact exists (either from commit comparison or via the ` + "`" + `baseArtifactId` + "`" + ` parameter),\ncomparison data showing size differences is included.\n\nThe response ` + "`" + `state` + "`" + ` field indicates the analysis status:\n- ` + "`" + `PENDING` + "`" + `: Analysis has not started yet.\n- ` + "`" + `PROCESSING` + "`" + `: Analysis is currently running.\n- ` + "`" + `FAILED` + "`" + ` / ` + "`" + `NOT_RAN` + "`" + `: Analysis did not complete; ` + "`" + `errorCode` + "`" + ` and ` + "`" + `errorMessage` + "`" + ` are included.\n- ` + "`" + `COMPLETED` + "`" + `: Analysis finished successfully with full size data.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "artifact_id",
          "in": "path",
          "required": true,
          "description": "The ID of the build artifact.",
          "type": "string"
        },
        {
          "name": "baseArtifactId",
          "in": "query",
          "required": false,
          "description": "Optional ID of the base artifact to compare against. If not provided, uses the default base head artifact.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "400",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "artifact_id"
        ],
        "query": [
          "baseArtifactId"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": null,
      "related": []
    },
    {
      "id": "Retrieves list of repositories for a given owner",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/prevent/owner/{owner}/repositories/",
      "tag": "Prevent",
      "summary": "Retrieves repository data for a given owner.",
      "description": "Retrieves repository data for a given owner.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "owner",
          "in": "path",
          "required": true,
          "description": "The owner of the repository.",
          "type": "string"
        },
        {
          "name": "limit",
          "in": "query",
          "required": false,
          "description": "The number of results to return. If not specified, defaults to 20.",
          "type": "integer"
        },
        {
          "name": "navigation",
          "in": "query",
          "required": false,
          "description": "Whether to get the previous or next page from paginated results. Use ` + "`" + `next` + "`" + ` for forward pagination after the cursor or ` + "`" + `prev` + "`" + ` for backward pagination before the cursor. If not specified, defaults to ` + "`" + `next` + "`" + `. If no cursor is provided, the cursor is the beginning of the result set.",
          "type": "string"
        },
        {
          "name": "cursor",
          "in": "query",
          "required": false,
          "description": "The cursor pointing to a specific position in the result set to start the query from. Results after the cursor will be returned if used with ` + "`" + `next` + "`" + ` or before the cursor if used with ` + "`" + `prev` + "`" + ` for ` + "`" + `navigation` + "`" + `.",
          "type": "string"
        },
        {
          "name": "term",
          "in": "query",
          "required": false,
          "description": "The term substring to filter name strings by using the ` + "`" + `contains` + "`" + ` operator.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "400",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "owner"
        ],
        "query": [
          "limit",
          "navigation",
          "cursor",
          "term"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": {
        "style": "cursor",
        "cursorParam": "cursor"
      },
      "related": []
    },
    {
      "id": "Gets syncing status for repositories for an integrated org",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/prevent/owner/{owner}/repositories/sync/",
      "tag": "Prevent",
      "summary": "Gets syncing status for repositories for an integrated organization.",
      "description": "Gets syncing status for repositories for an integrated organization.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "owner",
          "in": "path",
          "required": true,
          "description": "The owner of the repository.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "400",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "owner"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": null,
      "related": [
        {
          "id": "Syncs repositories from an integrated org with GitHub",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Syncs repositories from an integrated org with GitHub",
      "method": "POST",
      "path": "/api/0/organizations/{organization_id_or_slug}/prevent/owner/{owner}/repositories/sync/",
      "tag": "Prevent",
      "summary": "Syncs repositories for an integrated organization with GitHub.",
      "description": "Syncs repositories for an integrated organization with GitHub.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "owner",
          "in": "path",
          "required": true,
          "description": "The owner of the repository.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "write",
      "cacheable": false,
      "responseCodes": [
        "200",
        "400",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "owner"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "create",
      "pagination": null,
      "related": [
        {
          "id": "Gets syncing status for repositories for an integrated org",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Retrieves a paginated list of repository tokens for a given owner",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/prevent/owner/{owner}/repositories/tokens/",
      "tag": "Prevent",
      "summary": "Retrieves a paginated list of repository tokens for a given owner.",
      "description": "Retrieves a paginated list of repository tokens for a given owner.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "owner",
          "in": "path",
          "required": true,
          "description": "The owner of the repository.",
          "type": "string"
        },
        {
          "name": "limit",
          "in": "query",
          "required": false,
          "description": "The number of results to return. If not specified, defaults to 20.",
          "type": "integer"
        },
        {
          "name": "navigation",
          "in": "query",
          "required": false,
          "description": "Whether to get the previous or next page from paginated results. Use ` + "`" + `next` + "`" + ` for forward pagination after the cursor or ` + "`" + `prev` + "`" + ` for backward pagination before the cursor. If not specified, defaults to ` + "`" + `next` + "`" + `. If no cursor is provided, the cursor is the beginning of the result set.",
          "type": "string"
        },
        {
          "name": "cursor",
          "in": "query",
          "required": false,
          "description": "The cursor pointing to a specific position in the result set to start the query from. Results after the cursor will be returned if used with ` + "`" + `next` + "`" + ` or before the cursor if used with ` + "`" + `prev` + "`" + ` for ` + "`" + `navigation` + "`" + `.",
          "type": "string"
        },
        {
          "name": "sortBy",
          "in": "query",
          "required": false,
          "description": "The property to sort results by. If not specified, the default is ` + "`" + `COMMIT_DATE` + "`" + ` in descending order. Use ` + "`" + `-` + "`" + `\n        for descending order.\n\nAvailable fields are:\n- ` + "`" + `NAME` + "`" + `\n- ` + "`" + `COMMIT_DATE` + "`" + `\n        ",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "400",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "owner"
        ],
        "query": [
          "limit",
          "navigation",
          "cursor",
          "sortBy"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": {
        "style": "cursor",
        "cursorParam": "cursor"
      },
      "related": []
    },
    {
      "id": "Retrieves a single repository for a given owner",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/prevent/owner/{owner}/repository/{repository}/",
      "tag": "Prevent",
      "summary": "Retrieves repository data for a single repository.",
      "description": "Retrieves repository data for a single repository.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "owner",
          "in": "path",
          "required": true,
          "description": "The owner of the repository.",
          "type": "string"
        },
        {
          "name": "repository",
          "in": "path",
          "required": true,
          "description": "The name of the repository.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "400",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "owner",
          "repository"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": null,
      "related": []
    },
    {
      "id": "Retrieves list of branches for a given owner and repository",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/prevent/owner/{owner}/repository/{repository}/branches/",
      "tag": "Prevent",
      "summary": "Retrieves branch data for a given owner and repository.",
      "description": "Retrieves branch data for a given owner and repository.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "owner",
          "in": "path",
          "required": true,
          "description": "The owner of the repository.",
          "type": "string"
        },
        {
          "name": "repository",
          "in": "path",
          "required": true,
          "description": "The name of the repository.",
          "type": "string"
        },
        {
          "name": "limit",
          "in": "query",
          "required": false,
          "description": "The number of results to return. If not specified, defaults to 20.",
          "type": "integer"
        },
        {
          "name": "navigation",
          "in": "query",
          "required": false,
          "description": "Whether to get the previous or next page from paginated results. Use ` + "`" + `next` + "`" + ` for forward pagination after the cursor or ` + "`" + `prev` + "`" + ` for backward pagination before the cursor. If not specified, defaults to ` + "`" + `next` + "`" + `. If no cursor is provided, the cursor is the beginning of the result set.",
          "type": "string"
        },
        {
          "name": "cursor",
          "in": "query",
          "required": false,
          "description": "The cursor pointing to a specific position in the result set to start the query from. Results after the cursor will be returned if used with ` + "`" + `next` + "`" + ` or before the cursor if used with ` + "`" + `prev` + "`" + ` for ` + "`" + `navigation` + "`" + `.",
          "type": "string"
        },
        {
          "name": "term",
          "in": "query",
          "required": false,
          "description": "The term substring to filter name strings by using the ` + "`" + `contains` + "`" + ` operator.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "400",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "owner",
          "repository"
        ],
        "query": [
          "limit",
          "navigation",
          "cursor",
          "term"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": {
        "style": "cursor",
        "cursorParam": "cursor"
      },
      "related": []
    },
    {
      "id": "Retrieve paginated list of test results for repository, owner, and organization",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/prevent/owner/{owner}/repository/{repository}/test-results/",
      "tag": "Prevent",
      "summary": "Retrieves the list of test results for a given repository and owner. Also accepts a number of query parameters to filter the results.",
      "description": "Retrieves the list of test results for a given repository and owner. Also accepts a number of query parameters to filter the results.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "owner",
          "in": "path",
          "required": true,
          "description": "The owner of the repository.",
          "type": "string"
        },
        {
          "name": "repository",
          "in": "path",
          "required": true,
          "description": "The name of the repository.",
          "type": "string"
        },
        {
          "name": "sortBy",
          "in": "query",
          "required": false,
          "description": "The property to sort results by. If not specified, the default is ` + "`" + `TOTAL_FAIL_COUNT` + "`" + ` in descending order. Use ` + "`" + `-` + "`" + `\n        for descending order.\n\nAvailable fields are:\n- ` + "`" + `AVG_DURATION` + "`" + `\n- ` + "`" + `FLAKE_RATE` + "`" + `\n- ` + "`" + `FAILURE_RATE` + "`" + `\n- ` + "`" + `TOTAL_FAIL_COUNT` + "`" + `\n- ` + "`" + `UPDATED_AT` + "`" + `\n        ",
          "type": "string"
        },
        {
          "name": "filterBy",
          "in": "query",
          "required": false,
          "description": "An optional field to filter by, which will constrain the results to only include tests that match the filter.\n\nAvailable fields are:\n- ` + "`" + `FLAKY_TESTS` + "`" + `\n- ` + "`" + `FAILED_TESTS` + "`" + `\n- ` + "`" + `SLOWEST_TESTS` + "`" + `\n- ` + "`" + `SKIPPED_TESTS` + "`" + `\n        ",
          "type": "string"
        },
        {
          "name": "interval",
          "in": "query",
          "required": false,
          "description": "The time interval to search for results by.\n\nAvailable fields are:\n- ` + "`" + `INTERVAL_30_DAY` + "`" + `\n- ` + "`" + `INTERVAL_7_DAY` + "`" + `\n- ` + "`" + `INTERVAL_1_DAY` + "`" + `\n",
          "type": "string"
        },
        {
          "name": "branch",
          "in": "query",
          "required": false,
          "description": "The branch to search for results by. If not specified, the default is all branches.\n        ",
          "type": "string"
        },
        {
          "name": "limit",
          "in": "query",
          "required": false,
          "description": "The number of results to return. If not specified, defaults to 20.",
          "type": "integer"
        },
        {
          "name": "navigation",
          "in": "query",
          "required": false,
          "description": "Whether to get the previous or next page from paginated results. Use ` + "`" + `next` + "`" + ` for forward pagination after the cursor or ` + "`" + `prev` + "`" + ` for backward pagination before the cursor. If not specified, defaults to ` + "`" + `next` + "`" + `. If no cursor is provided, the cursor is the beginning of the result set.",
          "type": "string"
        },
        {
          "name": "cursor",
          "in": "query",
          "required": false,
          "description": "The cursor pointing to a specific position in the result set to start the query from. Results after the cursor will be returned if used with ` + "`" + `next` + "`" + ` or before the cursor if used with ` + "`" + `prev` + "`" + ` for ` + "`" + `navigation` + "`" + `.",
          "type": "string"
        },
        {
          "name": "term",
          "in": "query",
          "required": false,
          "description": "The term substring to filter name strings by using the ` + "`" + `contains` + "`" + ` operator.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "400",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "owner",
          "repository"
        ],
        "query": [
          "sortBy",
          "filterBy",
          "interval",
          "branch",
          "limit",
          "navigation",
          "cursor",
          "term"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "search",
      "pagination": {
        "style": "cursor",
        "cursorParam": "cursor"
      },
      "related": []
    },
    {
      "id": "Retrieve aggregated test result metrics for repository, owner, and organization",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/prevent/owner/{owner}/repository/{repository}/test-results-aggregates/",
      "tag": "Prevent",
      "summary": "Retrieves aggregated test result metrics for a given repository and owner.\nAlso accepts a query parameter to specify the time period for the metrics.",
      "description": "Retrieves aggregated test result metrics for a given repository and owner.\nAlso accepts a query parameter to specify the time period for the metrics.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "owner",
          "in": "path",
          "required": true,
          "description": "The owner of the repository.",
          "type": "string"
        },
        {
          "name": "repository",
          "in": "path",
          "required": true,
          "description": "The name of the repository.",
          "type": "string"
        },
        {
          "name": "interval",
          "in": "query",
          "required": false,
          "description": "The time interval to search for results by.\n\nAvailable fields are:\n- ` + "`" + `INTERVAL_30_DAY` + "`" + `\n- ` + "`" + `INTERVAL_7_DAY` + "`" + `\n- ` + "`" + `INTERVAL_1_DAY` + "`" + `\n",
          "type": "string"
        },
        {
          "name": "branch",
          "in": "query",
          "required": false,
          "description": "The branch to search for results by. If not specified, the default is all branches.\n        ",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "400",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "owner",
          "repository"
        ],
        "query": [
          "interval",
          "branch"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "search",
      "pagination": null,
      "related": []
    },
    {
      "id": "Retrieve test suites belonging to a repository's test results",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/prevent/owner/{owner}/repository/{repository}/test-suites/",
      "tag": "Prevent",
      "summary": "Retrieves test suites belonging to a repository's test results.\nIt accepts a list of test suites as a query parameter to specify individual test suites.",
      "description": "Retrieves test suites belonging to a repository's test results.\nIt accepts a list of test suites as a query parameter to specify individual test suites.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "owner",
          "in": "path",
          "required": true,
          "description": "The owner of the repository.",
          "type": "string"
        },
        {
          "name": "repository",
          "in": "path",
          "required": true,
          "description": "The name of the repository.",
          "type": "string"
        },
        {
          "name": "term",
          "in": "query",
          "required": false,
          "description": "The term substring to filter name strings by using the ` + "`" + `contains` + "`" + ` operator.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "400",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "owner",
          "repository"
        ],
        "query": [
          "term"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "search",
      "pagination": null,
      "related": []
    },
    {
      "id": "Regenerates a repository upload token and returns the new token",
      "method": "POST",
      "path": "/api/0/organizations/{organization_id_or_slug}/prevent/owner/{owner}/repository/{repository}/token/regenerate/",
      "tag": "Prevent",
      "summary": "Regenerates a repository upload token and returns the new token.",
      "description": "Regenerates a repository upload token and returns the new token.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "owner",
          "in": "path",
          "required": true,
          "description": "The owner of the repository.",
          "type": "string"
        },
        {
          "name": "repository",
          "in": "path",
          "required": true,
          "description": "The name of the repository.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "write",
      "cacheable": false,
      "responseCodes": [
        "200",
        "400",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "owner",
          "repository"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "create",
      "pagination": null,
      "related": []
    },
    {
      "id": "List an Organization's Client Keys",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/project-keys/",
      "tag": "Organizations",
      "summary": "Return a list of client keys (DSNs) for all projects in an organization.\n\nThis paginated endpoint lists client keys across all projects in an organization. Each key includes the project ID\nto identify which project it belongs to.\n\nQuery Parameters:\n- team: Filter by team slug or ID to get keys only for that team's projects\n- status: Filter by 'active' or 'inactive' to get keys with specific status",
      "description": "Return a list of client keys (DSNs) for all projects in an organization.\n\nThis paginated endpoint lists client keys across all projects in an organization. Each key includes the project ID\nto identify which project it belongs to.\n\nQuery Parameters:\n- team: Filter by team slug or ID to get keys only for that team's projects\n- status: Filter by 'active' or 'inactive' to get keys with specific status",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "cursor",
          "in": "query",
          "required": false,
          "description": "A pointer to the last object fetched and its sort order; used to retrieve the next or previous results.",
          "type": "string"
        },
        {
          "name": "team",
          "in": "query",
          "required": false,
          "description": "Filter keys by team slug or ID. If provided, only keys for projects belonging to this team will be returned.",
          "type": "string"
        },
        {
          "name": "status",
          "in": "query",
          "required": false,
          "description": "Filter keys by status. Options are 'active' or 'inactive'.\n\n* ` + "`" + `active` + "`" + `\n* ` + "`" + `inactive` + "`" + `",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "400",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug"
        ],
        "query": [
          "cursor",
          "team",
          "status"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "search",
      "pagination": {
        "style": "cursor",
        "cursorParam": "cursor"
      },
      "related": []
    },
    {
      "id": "List an Organization's Projects",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/projects/",
      "tag": "Organizations",
      "summary": "Return a list of projects bound to a organization.",
      "description": "Return a list of projects bound to a organization.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "cursor",
          "in": "query",
          "required": false,
          "description": "A pointer to the last object fetched and its sort order; used to retrieve the next or previous results.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "401",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug"
        ],
        "query": [
          "cursor"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": {
        "style": "cursor",
        "cursorParam": "cursor"
      },
      "related": []
    },
    {
      "id": "Create a Monitor for a Project",
      "method": "POST",
      "path": "/api/0/organizations/{organization_id_or_slug}/projects/{project_id_or_slug}/detectors/",
      "tag": "Monitors",
      "summary": "⚠️ This endpoint is currently in **beta** and may be subject to change. It is supported by [New Monitors and Alerts](/product/new-monitors-and-alerts/) and may not be viewable in the UI today.\n\nCreate a Monitor for a project",
      "description": "⚠️ This endpoint is currently in **beta** and may be subject to change. It is supported by [New Monitors and Alerts](/product/new-monitors-and-alerts/) and may not be viewable in the UI today.\n\nCreate a Monitor for a project",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the project the resource belongs to.",
          "type": "string"
        }
      ],
      "hasRequestBody": true,
      "risk": "write",
      "cacheable": false,
      "responseCodes": [
        "201",
        "400",
        "401",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "project_id_or_slug"
        ],
        "query": [],
        "headers": [],
        "body": true
      },
      "graphql": null,
      "kind": "create",
      "pagination": null,
      "related": []
    },
    {
      "id": "List an Organization's trusted Relays",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/relay_usage/",
      "tag": "Organizations",
      "summary": "Return a list of trusted relays bound to an organization.",
      "description": "Return a list of trusted relays bound to an organization.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": null,
      "related": []
    },
    {
      "id": "Retrieve Statuses of Release Thresholds (Alpha)",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/release-threshold-statuses/",
      "tag": "Releases",
      "summary": "**` + "`" + `[WARNING]` + "`" + `**: This API is an experimental Alpha feature and is subject to change!\n\nList all derived statuses of releases that fall within the provided start/end datetimes.\n\nConstructs a response key'd off \\{` + "`" + `release_version` + "`" + `\\}-\\{` + "`" + `project_slug` + "`" + `\\} that lists thresholds with their status for *specified* projects.\nEach returned enriched threshold will contain the full serialized ` + "`" + `release_threshold` + "`" + ` instance as well as it's derived health statuses.",
      "description": "**` + "`" + `[WARNING]` + "`" + `**: This API is an experimental Alpha feature and is subject to change!\n\nList all derived statuses of releases that fall within the provided start/end datetimes.\n\nConstructs a response key'd off \\{` + "`" + `release_version` + "`" + `\\}-\\{` + "`" + `project_slug` + "`" + `\\} that lists thresholds with their status for *specified* projects.\nEach returned enriched threshold will contain the full serialized ` + "`" + `release_threshold` + "`" + ` instance as well as it's derived health statuses.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "start",
          "in": "query",
          "required": true,
          "description": "The start of the time series range as an explicit datetime, either in UTC ISO8601 or epoch seconds. Use along with ` + "`" + `end` + "`" + `.",
          "type": "string"
        },
        {
          "name": "end",
          "in": "query",
          "required": true,
          "description": "The inclusive end of the time series range as an explicit datetime, either in UTC ISO8601 or epoch seconds. Use along with ` + "`" + `start` + "`" + `.",
          "type": "string"
        },
        {
          "name": "environment",
          "in": "query",
          "required": false,
          "description": "A list of environment names to filter your results by.",
          "type": "array"
        },
        {
          "name": "projectSlug",
          "in": "query",
          "required": false,
          "description": "A list of project slugs to filter your results by.",
          "type": "array"
        },
        {
          "name": "release",
          "in": "query",
          "required": false,
          "description": "A list of release versions to filter your results by.",
          "type": "array"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "400"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug"
        ],
        "query": [
          "start",
          "end",
          "environment",
          "projectSlug",
          "release"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": null,
      "related": []
    },
    {
      "id": "Retrieve an Organization's Release",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/releases/{version}/",
      "tag": "Releases",
      "summary": "Return details on an individual release.",
      "description": "Return details on an individual release.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "version",
          "in": "path",
          "required": true,
          "description": "The version identifier of the release",
          "type": "string"
        },
        {
          "name": "project_id",
          "in": "query",
          "required": false,
          "description": "The project ID to filter by.",
          "type": "string"
        },
        {
          "name": "health",
          "in": "query",
          "required": false,
          "description": "Whether or not to include health data with the release. By default, this is false.",
          "type": "boolean"
        },
        {
          "name": "adoptionStages",
          "in": "query",
          "required": false,
          "description": "Whether or not to include adoption stages with the release. By default, this is false.",
          "type": "boolean"
        },
        {
          "name": "summaryStatsPeriod",
          "in": "query",
          "required": false,
          "description": "The period of time used to query summary stats for the release. By default, this is 14d.",
          "type": "string"
        },
        {
          "name": "healthStatsPeriod",
          "in": "query",
          "required": false,
          "description": "The period of time used to query health stats for the release. By default, this is 24h if health is enabled.",
          "type": "string"
        },
        {
          "name": "sort",
          "in": "query",
          "required": false,
          "description": "The field used to sort results by. By default, this is ` + "`" + `date` + "`" + `.",
          "type": "string"
        },
        {
          "name": "status",
          "in": "query",
          "required": false,
          "description": "Release statuses that you can filter by.",
          "type": "string"
        },
        {
          "name": "query",
          "in": "query",
          "required": false,
          "description": "Filters results by using [query syntax](/product/sentry-basics/search/).\n\nExample: ` + "`" + `query=(transaction:foo AND release:abc) OR (transaction:[bar,baz] AND release:def)` + "`" + `\n",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "401",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "version"
        ],
        "query": [
          "project_id",
          "health",
          "adoptionStages",
          "summaryStatsPeriod",
          "healthStatsPeriod",
          "sort",
          "status",
          "query"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": null,
      "related": [
        {
          "id": "Update an Organization's Release",
          "role": "same-resource"
        },
        {
          "id": "Delete an Organization's Release",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Update an Organization's Release",
      "method": "PUT",
      "path": "/api/0/organizations/{organization_id_or_slug}/releases/{version}/",
      "tag": "Releases",
      "summary": "Update a release. This can change some metadata associated with\nthe release (the ref, url, and dates).",
      "description": "Update a release. This can change some metadata associated with\nthe release (the ref, url, and dates).",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "version",
          "in": "path",
          "required": true,
          "description": "The version identifier of the release",
          "type": "string"
        }
      ],
      "hasRequestBody": true,
      "risk": "write",
      "cacheable": false,
      "responseCodes": [
        "200",
        "401",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "version"
        ],
        "query": [],
        "headers": [],
        "body": true
      },
      "graphql": null,
      "kind": "update",
      "pagination": null,
      "related": [
        {
          "id": "Retrieve an Organization's Release",
          "role": "same-resource"
        },
        {
          "id": "Delete an Organization's Release",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Delete an Organization's Release",
      "method": "DELETE",
      "path": "/api/0/organizations/{organization_id_or_slug}/releases/{version}/",
      "tag": "Releases",
      "summary": "Permanently remove a release and all of its files.",
      "description": "Permanently remove a release and all of its files.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "version",
          "in": "path",
          "required": true,
          "description": "The version identifier of the release",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "destructive",
      "cacheable": false,
      "responseCodes": [
        "204",
        "401",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "version"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "delete",
      "pagination": null,
      "related": [
        {
          "id": "Retrieve an Organization's Release",
          "role": "same-resource"
        },
        {
          "id": "Update an Organization's Release",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "List a Release's Deploys",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/releases/{version}/deploys/",
      "tag": "Releases",
      "summary": "Returns a list of deploys based on the organization, version, and project.",
      "description": "Returns a list of deploys based on the organization, version, and project.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "version",
          "in": "path",
          "required": true,
          "description": "The version identifier of the release",
          "type": "string"
        },
        {
          "name": "cursor",
          "in": "query",
          "required": false,
          "description": "A pointer to the last object fetched and its sort order; used to retrieve the next or previous results.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "version"
        ],
        "query": [
          "cursor"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": {
        "style": "cursor",
        "cursorParam": "cursor"
      },
      "related": [
        {
          "id": "Create a Deploy",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Create a Deploy",
      "method": "POST",
      "path": "/api/0/organizations/{organization_id_or_slug}/releases/{version}/deploys/",
      "tag": "Releases",
      "summary": "Create a deploy for a given release.",
      "description": "Create a deploy for a given release.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "version",
          "in": "path",
          "required": true,
          "description": "The version identifier of the release",
          "type": "string"
        }
      ],
      "hasRequestBody": true,
      "risk": "write",
      "cacheable": false,
      "responseCodes": [
        "201",
        "400"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "version"
        ],
        "query": [],
        "headers": [],
        "body": true
      },
      "graphql": null,
      "kind": "create",
      "pagination": null,
      "related": [
        {
          "id": "List a Release's Deploys",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Retrieve a Count of Replays for a Given Issue or Transaction",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/replay-count/",
      "tag": "Replays",
      "summary": "Return a count of replays for a list of issue or transaction IDs.\n\nThe ` + "`" + `query` + "`" + ` parameter is required. It is a search query that includes exactly one of ` + "`" + `issue.id` + "`" + `, ` + "`" + `transaction` + "`" + `, or ` + "`" + `replay_id` + "`" + ` (string or list of strings).",
      "description": "Return a count of replays for a list of issue or transaction IDs.\n\nThe ` + "`" + `query` + "`" + ` parameter is required. It is a search query that includes exactly one of ` + "`" + `issue.id` + "`" + `, ` + "`" + `transaction` + "`" + `, or ` + "`" + `replay_id` + "`" + ` (string or list of strings).",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "environment",
          "in": "query",
          "required": false,
          "description": "The name of environments to filter by.",
          "type": "array"
        },
        {
          "name": "start",
          "in": "query",
          "required": false,
          "description": "The start of the period of time for the query, expected in ISO-8601 format. For example, ` + "`" + `2001-12-14T12:34:56.7890` + "`" + `.",
          "type": "string"
        },
        {
          "name": "end",
          "in": "query",
          "required": false,
          "description": "The end of the period of time for the query, expected in ISO-8601 format. For example, ` + "`" + `2001-12-14T12:34:56.7890` + "`" + `.",
          "type": "string"
        },
        {
          "name": "statsPeriod",
          "in": "query",
          "required": false,
          "description": "The period of time for the query, will override the start & end parameters, a number followed by one of:\n- ` + "`" + `d` + "`" + ` for days\n- ` + "`" + `h` + "`" + ` for hours\n- ` + "`" + `m` + "`" + ` for minutes\n- ` + "`" + `s` + "`" + ` for seconds\n- ` + "`" + `w` + "`" + ` for weeks\n\nFor example, ` + "`" + `24h` + "`" + `, to mean query data starting from 24 hours ago to now.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "query",
          "required": false,
          "description": "The project slugs to filter by. Use ` + "`" + `$all` + "`" + ` to include all available projects. For example, the following are valid parameters:\n- ` + "`" + `/?projectSlug=$all` + "`" + `\n- ` + "`" + `/?projectSlug=android&projectSlug=javascript-react` + "`" + `\n",
          "type": "array"
        },
        {
          "name": "query",
          "in": "query",
          "required": false,
          "description": "Filters results by using [query syntax](/product/sentry-basics/search/).\n\nExample: ` + "`" + `query=(transaction:foo AND release:abc) OR (transaction:[bar,baz] AND release:def)` + "`" + `\n",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "400",
        "403"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug"
        ],
        "query": [
          "environment",
          "start",
          "end",
          "statsPeriod",
          "project_id_or_slug",
          "query"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "search",
      "pagination": null,
      "related": []
    },
    {
      "id": "List an Organization's Selectors",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/replay-selectors/",
      "tag": "Replays",
      "summary": "Return a list of selectors for a given organization.",
      "description": "Return a list of selectors for a given organization.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "environment",
          "in": "query",
          "required": false,
          "description": "The environment to filter by.",
          "type": "array"
        },
        {
          "name": "statsPeriod",
          "in": "query",
          "required": false,
          "description": "This defines the range of the time series, relative to now. The range is given in a ` + "`" + `<number><unit>` + "`" + ` format. For example ` + "`" + `1d` + "`" + ` for a one day range. Possible units are ` + "`" + `m` + "`" + ` for minutes, ` + "`" + `h` + "`" + ` for hours, ` + "`" + `d` + "`" + ` for days and ` + "`" + `w` + "`" + ` for weeks.You must either provide a ` + "`" + `statsPeriod` + "`" + `, or a ` + "`" + `start` + "`" + ` and ` + "`" + `end` + "`" + `.",
          "type": "string"
        },
        {
          "name": "start",
          "in": "query",
          "required": false,
          "description": "This defines the start of the time series range as an explicit datetime, either in UTC ISO8601 or epoch seconds.Use along with ` + "`" + `end` + "`" + ` instead of ` + "`" + `statsPeriod` + "`" + `.",
          "type": "string"
        },
        {
          "name": "end",
          "in": "query",
          "required": false,
          "description": "This defines the inclusive end of the time series range as an explicit datetime, either in UTC ISO8601 or epoch seconds.Use along with ` + "`" + `start` + "`" + ` instead of ` + "`" + `statsPeriod` + "`" + `.",
          "type": "string"
        },
        {
          "name": "project",
          "in": "query",
          "required": false,
          "description": "The ID of the projects to filter by.",
          "type": "array"
        },
        {
          "name": "projectSlug",
          "in": "query",
          "required": false,
          "description": "A list of project slugs to filter your results by.",
          "type": "array"
        },
        {
          "name": "sort",
          "in": "query",
          "required": false,
          "description": "The field to sort the output by.",
          "type": "string"
        },
        {
          "name": "sortBy",
          "in": "query",
          "required": false,
          "description": "The field to sort the output by.",
          "type": "string"
        },
        {
          "name": "orderBy",
          "in": "query",
          "required": false,
          "description": "The field to sort the output by.",
          "type": "string"
        },
        {
          "name": "cursor",
          "in": "query",
          "required": false,
          "description": "A pointer to the last object fetched and its sort order; used to retrieve the next or previous results.",
          "type": "string"
        },
        {
          "name": "per_page",
          "in": "query",
          "required": false,
          "description": "Limit the number of rows to return in the result. Default and maximum allowed is 100.",
          "type": "integer"
        },
        {
          "name": "query",
          "in": "query",
          "required": false,
          "description": "Filters results by using [query syntax](/product/sentry-basics/search/).\n\nExample: ` + "`" + `query=(transaction:foo AND release:abc) OR (transaction:[bar,baz] AND release:def)` + "`" + `\n",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "400",
        "403"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug"
        ],
        "query": [
          "environment",
          "statsPeriod",
          "start",
          "end",
          "project",
          "projectSlug",
          "sort",
          "sortBy",
          "orderBy",
          "cursor",
          "per_page",
          "query"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": {
        "style": "cursor",
        "cursorParam": "cursor"
      },
      "related": []
    },
    {
      "id": "List an Organization's Replays",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/replays/",
      "tag": "Replays",
      "summary": "Return a list of replays belonging to an organization.",
      "description": "Return a list of replays belonging to an organization.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "statsPeriod",
          "in": "query",
          "required": false,
          "description": "\nThis defines the range of the time series, relative to now. The range is given in a\n` + "`" + `<number><unit>` + "`" + ` format. For example ` + "`" + `1d` + "`" + ` for a one day range. Possible units are ` + "`" + `m` + "`" + ` for\nminutes, ` + "`" + `h` + "`" + ` for hours, ` + "`" + `d` + "`" + ` for days and ` + "`" + `w` + "`" + ` for weeks. You must either provide a\n` + "`" + `statsPeriod` + "`" + `, or a ` + "`" + `start` + "`" + ` and ` + "`" + `end` + "`" + `.\n",
          "type": "string"
        },
        {
          "name": "start",
          "in": "query",
          "required": false,
          "description": "\nThis defines the start of the time series range as an explicit datetime, either in UTC\nISO8601 or epoch seconds. Use along with ` + "`" + `end` + "`" + ` instead of ` + "`" + `statsPeriod` + "`" + `.\n",
          "type": "string"
        },
        {
          "name": "end",
          "in": "query",
          "required": false,
          "description": "\nThis defines the inclusive end of the time series range as an explicit datetime, either in\nUTC ISO8601 or epoch seconds. Use along with ` + "`" + `start` + "`" + ` instead of ` + "`" + `statsPeriod` + "`" + `.\n",
          "type": "string"
        },
        {
          "name": "field",
          "in": "query",
          "required": false,
          "description": "Specifies a field that should be marshaled in the output. Invalid fields will be rejected.",
          "type": "array"
        },
        {
          "name": "project",
          "in": "query",
          "required": false,
          "description": "The ID of the projects to filter by.",
          "type": "array"
        },
        {
          "name": "projectSlug",
          "in": "query",
          "required": false,
          "description": "A list of project slugs to filter your results by.",
          "type": "array"
        },
        {
          "name": "environment",
          "in": "query",
          "required": false,
          "description": "The environment to filter by.",
          "type": "string"
        },
        {
          "name": "sort",
          "in": "query",
          "required": false,
          "description": "The field to sort the output by.",
          "type": "string"
        },
        {
          "name": "sortBy",
          "in": "query",
          "required": false,
          "description": "The field to sort the output by.",
          "type": "string"
        },
        {
          "name": "orderBy",
          "in": "query",
          "required": false,
          "description": "The field to sort the output by.",
          "type": "string"
        },
        {
          "name": "query",
          "in": "query",
          "required": false,
          "description": "A structured query string to filter the output by.",
          "type": "string"
        },
        {
          "name": "per_page",
          "in": "query",
          "required": false,
          "description": "Limit the number of rows to return in the result.",
          "type": "integer"
        },
        {
          "name": "cursor",
          "in": "query",
          "required": false,
          "description": "The cursor parameter is used to paginate results. See [here](https://docs.sentry.io/api/pagination/) for how to use this query parameter",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "400",
        "403"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug"
        ],
        "query": [
          "statsPeriod",
          "start",
          "end",
          "field",
          "project",
          "projectSlug",
          "environment",
          "sort",
          "sortBy",
          "orderBy",
          "query",
          "per_page",
          "cursor"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": {
        "style": "cursor",
        "cursorParam": "cursor"
      },
      "related": []
    },
    {
      "id": "Retrieve a Replay Instance",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/replays/{replay_id}/",
      "tag": "Replays",
      "summary": "Return details on an individual replay.",
      "description": "Return details on an individual replay.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "replay_id",
          "in": "path",
          "required": true,
          "description": "The ID of the replay you'd like to retrieve.",
          "type": "string"
        },
        {
          "name": "statsPeriod",
          "in": "query",
          "required": false,
          "description": "\nThis defines the range of the time series, relative to now. The range is given in a\n` + "`" + `<number><unit>` + "`" + ` format. For example ` + "`" + `1d` + "`" + ` for a one day range. Possible units are ` + "`" + `m` + "`" + ` for\nminutes, ` + "`" + `h` + "`" + ` for hours, ` + "`" + `d` + "`" + ` for days and ` + "`" + `w` + "`" + ` for weeks. You must either provide a\n` + "`" + `statsPeriod` + "`" + `, or a ` + "`" + `start` + "`" + ` and ` + "`" + `end` + "`" + `.\n",
          "type": "string"
        },
        {
          "name": "start",
          "in": "query",
          "required": false,
          "description": "\nThis defines the start of the time series range as an explicit datetime, either in UTC\nISO8601 or epoch seconds. Use along with ` + "`" + `end` + "`" + ` instead of ` + "`" + `statsPeriod` + "`" + `.\n",
          "type": "string"
        },
        {
          "name": "end",
          "in": "query",
          "required": false,
          "description": "\nThis defines the inclusive end of the time series range as an explicit datetime, either in\nUTC ISO8601 or epoch seconds. Use along with ` + "`" + `start` + "`" + ` instead of ` + "`" + `statsPeriod` + "`" + `.\n",
          "type": "string"
        },
        {
          "name": "field",
          "in": "query",
          "required": false,
          "description": "Specifies a field that should be marshaled in the output. Invalid fields will be rejected.",
          "type": "array"
        },
        {
          "name": "project",
          "in": "query",
          "required": false,
          "description": "The ID of the projects to filter by.",
          "type": "array"
        },
        {
          "name": "projectSlug",
          "in": "query",
          "required": false,
          "description": "A list of project slugs to filter your results by.",
          "type": "array"
        },
        {
          "name": "environment",
          "in": "query",
          "required": false,
          "description": "The environment to filter by.",
          "type": "string"
        },
        {
          "name": "sort",
          "in": "query",
          "required": false,
          "description": "The field to sort the output by.",
          "type": "string"
        },
        {
          "name": "sortBy",
          "in": "query",
          "required": false,
          "description": "The field to sort the output by.",
          "type": "string"
        },
        {
          "name": "orderBy",
          "in": "query",
          "required": false,
          "description": "The field to sort the output by.",
          "type": "string"
        },
        {
          "name": "query",
          "in": "query",
          "required": false,
          "description": "A structured query string to filter the output by.",
          "type": "string"
        },
        {
          "name": "per_page",
          "in": "query",
          "required": false,
          "description": "Limit the number of rows to return in the result.",
          "type": "integer"
        },
        {
          "name": "cursor",
          "in": "query",
          "required": false,
          "description": "The cursor parameter is used to paginate results. See [here](https://docs.sentry.io/api/pagination/) for how to use this query parameter",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "400",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "replay_id"
        ],
        "query": [
          "statsPeriod",
          "start",
          "end",
          "field",
          "project",
          "projectSlug",
          "environment",
          "sort",
          "sortBy",
          "orderBy",
          "query",
          "per_page",
          "cursor"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": {
        "style": "cursor",
        "cursorParam": "cursor"
      },
      "related": []
    },
    {
      "id": "List a Repository's Commits",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/repos/{repo_id}/commits/",
      "tag": "Organizations",
      "summary": "List a Repository's Commits",
      "description": "List a Repository's Commits",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "repo_id",
          "in": "path",
          "required": true,
          "description": "The repository ID.",
          "type": "string"
        },
        {
          "name": "cursor",
          "in": "query",
          "required": false,
          "description": "A pointer to the last object fetched and its sort order; used to retrieve the next or previous results.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "400",
        "401",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "repo_id"
        ],
        "query": [
          "cursor"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": {
        "style": "cursor",
        "cursorParam": "cursor"
      },
      "related": []
    },
    {
      "id": "List an Organization's Paginated Teams",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/scim/v2/Groups",
      "tag": "SCIM",
      "summary": "Returns a paginated list of teams bound to a organization with a SCIM Groups GET Request.\n\nNote that the members field will only contain up to 10,000 members.",
      "description": "Returns a paginated list of teams bound to a organization with a SCIM Groups GET Request.\n\nNote that the members field will only contain up to 10,000 members.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "startIndex",
          "in": "query",
          "required": false,
          "description": "SCIM 1-offset based index for pagination.",
          "type": "integer"
        },
        {
          "name": "count",
          "in": "query",
          "required": false,
          "description": "The maximum number of results the query should return, maximum of 100.",
          "type": "integer"
        },
        {
          "name": "filter",
          "in": "query",
          "required": false,
          "description": "A SCIM filter expression. The only operator currently supported is ` + "`" + `eq` + "`" + `.",
          "type": "string"
        },
        {
          "name": "excludedAttributes",
          "in": "query",
          "required": false,
          "description": "Fields that should be left off of return values. Right now the only supported field for this query is members.",
          "type": "array"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "401",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug"
        ],
        "query": [
          "startIndex",
          "count",
          "filter",
          "excludedAttributes"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": null,
      "related": [
        {
          "id": "Provision a New Team",
          "role": "same-resource"
        },
        {
          "id": "Query an Individual Team",
          "role": "detail-child"
        },
        {
          "id": "Update a Team's Attributes",
          "role": "detail-child"
        },
        {
          "id": "Delete an Individual Team",
          "role": "detail-child"
        }
      ]
    },
    {
      "id": "Provision a New Team",
      "method": "POST",
      "path": "/api/0/organizations/{organization_id_or_slug}/scim/v2/Groups",
      "tag": "SCIM",
      "summary": "Create a new team bound to an organization via a SCIM Groups POST\nRequest. The slug will have a normalization of uppercases/spaces to\nlowercases and dashes.\n\nNote that teams are always created with an empty member set.",
      "description": "Create a new team bound to an organization via a SCIM Groups POST\nRequest. The slug will have a normalization of uppercases/spaces to\nlowercases and dashes.\n\nNote that teams are always created with an empty member set.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        }
      ],
      "hasRequestBody": true,
      "risk": "write",
      "cacheable": false,
      "responseCodes": [
        "201",
        "401",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug"
        ],
        "query": [],
        "headers": [],
        "body": true
      },
      "graphql": null,
      "kind": "create",
      "pagination": null,
      "related": [
        {
          "id": "List an Organization's Paginated Teams",
          "role": "same-resource"
        },
        {
          "id": "Query an Individual Team",
          "role": "detail-child"
        },
        {
          "id": "Update a Team's Attributes",
          "role": "detail-child"
        },
        {
          "id": "Delete an Individual Team",
          "role": "detail-child"
        }
      ]
    },
    {
      "id": "Query an Individual Team",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/scim/v2/Groups/{team_id_or_slug}",
      "tag": "SCIM",
      "summary": "Query an individual team with a SCIM Group GET Request.\n- Note that the members field will only contain up to 10000 members.",
      "description": "Query an individual team with a SCIM Group GET Request.\n- Note that the members field will only contain up to 10000 members.",
      "parameters": [
        {
          "name": "team_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the team the resource belongs to.",
          "type": "string"
        },
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "401",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "team_id_or_slug",
          "organization_id_or_slug"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "search",
      "pagination": null,
      "related": [
        {
          "id": "List an Organization's Paginated Teams",
          "role": "list-parent"
        },
        {
          "id": "Provision a New Team",
          "role": "list-parent"
        },
        {
          "id": "Update a Team's Attributes",
          "role": "same-resource"
        },
        {
          "id": "Delete an Individual Team",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Update a Team's Attributes",
      "method": "PATCH",
      "path": "/api/0/organizations/{organization_id_or_slug}/scim/v2/Groups/{team_id_or_slug}",
      "tag": "SCIM",
      "summary": "Update a team's attributes with a SCIM Group PATCH Request.",
      "description": "Update a team's attributes with a SCIM Group PATCH Request.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "team_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the team the resource belongs to.",
          "type": "string"
        }
      ],
      "hasRequestBody": true,
      "risk": "write",
      "cacheable": false,
      "responseCodes": [
        "204",
        "401",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "team_id_or_slug"
        ],
        "query": [],
        "headers": [],
        "body": true
      },
      "graphql": null,
      "kind": "update",
      "pagination": null,
      "related": [
        {
          "id": "List an Organization's Paginated Teams",
          "role": "list-parent"
        },
        {
          "id": "Provision a New Team",
          "role": "list-parent"
        },
        {
          "id": "Query an Individual Team",
          "role": "same-resource"
        },
        {
          "id": "Delete an Individual Team",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Delete an Individual Team",
      "method": "DELETE",
      "path": "/api/0/organizations/{organization_id_or_slug}/scim/v2/Groups/{team_id_or_slug}",
      "tag": "SCIM",
      "summary": "Delete a team with a SCIM Group DELETE Request.",
      "description": "Delete a team with a SCIM Group DELETE Request.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "team_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the team the resource belongs to.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "destructive",
      "cacheable": false,
      "responseCodes": [
        "204",
        "401",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "team_id_or_slug"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "delete",
      "pagination": null,
      "related": [
        {
          "id": "List an Organization's Paginated Teams",
          "role": "list-parent"
        },
        {
          "id": "Provision a New Team",
          "role": "list-parent"
        },
        {
          "id": "Query an Individual Team",
          "role": "same-resource"
        },
        {
          "id": "Update a Team's Attributes",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "List an Organization's SCIM Members",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/scim/v2/Users",
      "tag": "SCIM",
      "summary": "Returns a paginated list of members bound to a organization with a SCIM Users GET Request.",
      "description": "Returns a paginated list of members bound to a organization with a SCIM Users GET Request.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "startIndex",
          "in": "query",
          "required": false,
          "description": "SCIM 1-offset based index for pagination.",
          "type": "integer"
        },
        {
          "name": "count",
          "in": "query",
          "required": false,
          "description": "The maximum number of results the query should return, maximum of 100.",
          "type": "integer"
        },
        {
          "name": "filter",
          "in": "query",
          "required": false,
          "description": "A SCIM filter expression. The only operator currently supported is ` + "`" + `eq` + "`" + `.",
          "type": "string"
        },
        {
          "name": "excludedAttributes",
          "in": "query",
          "required": false,
          "description": "Fields that should be left off of return values. Right now the only supported field for this query is members.",
          "type": "array"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "401",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug"
        ],
        "query": [
          "startIndex",
          "count",
          "filter",
          "excludedAttributes"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": null,
      "related": [
        {
          "id": "Provision a New Organization Member",
          "role": "same-resource"
        },
        {
          "id": "Query an Individual Organization Member",
          "role": "detail-child"
        },
        {
          "id": "Update an Organization Member's Attributes",
          "role": "detail-child"
        },
        {
          "id": "Delete an Organization Member via SCIM",
          "role": "detail-child"
        }
      ]
    },
    {
      "id": "Provision a New Organization Member",
      "method": "POST",
      "path": "/api/0/organizations/{organization_id_or_slug}/scim/v2/Users",
      "tag": "SCIM",
      "summary": "Create a new Organization Member via a SCIM Users POST Request.\n\nNote that this API does not support setting secondary emails.",
      "description": "Create a new Organization Member via a SCIM Users POST Request.\n\nNote that this API does not support setting secondary emails.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        }
      ],
      "hasRequestBody": true,
      "risk": "write",
      "cacheable": false,
      "responseCodes": [
        "201",
        "401",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug"
        ],
        "query": [],
        "headers": [],
        "body": true
      },
      "graphql": null,
      "kind": "create",
      "pagination": null,
      "related": [
        {
          "id": "List an Organization's SCIM Members",
          "role": "same-resource"
        },
        {
          "id": "Query an Individual Organization Member",
          "role": "detail-child"
        },
        {
          "id": "Update an Organization Member's Attributes",
          "role": "detail-child"
        },
        {
          "id": "Delete an Organization Member via SCIM",
          "role": "detail-child"
        }
      ]
    },
    {
      "id": "Query an Individual Organization Member",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/scim/v2/Users/{member_id}",
      "tag": "SCIM",
      "summary": "Query an individual organization member with a SCIM User GET Request.\n- The ` + "`" + `name` + "`" + ` object will contain fields ` + "`" + `firstName` + "`" + ` and ` + "`" + `lastName` + "`" + ` with the values of ` + "`" + `N/A` + "`" + `.\nSentry's SCIM API does not currently support these fields but returns them for compatibility purposes.",
      "description": "Query an individual organization member with a SCIM User GET Request.\n- The ` + "`" + `name` + "`" + ` object will contain fields ` + "`" + `firstName` + "`" + ` and ` + "`" + `lastName` + "`" + ` with the values of ` + "`" + `N/A` + "`" + `.\nSentry's SCIM API does not currently support these fields but returns them for compatibility purposes.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "member_id",
          "in": "path",
          "required": true,
          "description": "The ID of the member to query.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "401",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "member_id"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "search",
      "pagination": null,
      "related": [
        {
          "id": "List an Organization's SCIM Members",
          "role": "list-parent"
        },
        {
          "id": "Provision a New Organization Member",
          "role": "list-parent"
        },
        {
          "id": "Update an Organization Member's Attributes",
          "role": "same-resource"
        },
        {
          "id": "Delete an Organization Member via SCIM",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Update an Organization Member's Attributes",
      "method": "PATCH",
      "path": "/api/0/organizations/{organization_id_or_slug}/scim/v2/Users/{member_id}",
      "tag": "SCIM",
      "summary": "Update an organization member's attributes with a SCIM PATCH Request.",
      "description": "Update an organization member's attributes with a SCIM PATCH Request.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "member_id",
          "in": "path",
          "required": true,
          "description": "The ID of the member to update.",
          "type": "string"
        }
      ],
      "hasRequestBody": true,
      "risk": "write",
      "cacheable": false,
      "responseCodes": [
        "204",
        "401",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "member_id"
        ],
        "query": [],
        "headers": [],
        "body": true
      },
      "graphql": null,
      "kind": "update",
      "pagination": null,
      "related": [
        {
          "id": "List an Organization's SCIM Members",
          "role": "list-parent"
        },
        {
          "id": "Provision a New Organization Member",
          "role": "list-parent"
        },
        {
          "id": "Query an Individual Organization Member",
          "role": "same-resource"
        },
        {
          "id": "Delete an Organization Member via SCIM",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Delete an Organization Member via SCIM",
      "method": "DELETE",
      "path": "/api/0/organizations/{organization_id_or_slug}/scim/v2/Users/{member_id}",
      "tag": "SCIM",
      "summary": "Delete an organization member with a SCIM User DELETE Request.",
      "description": "Delete an organization member with a SCIM User DELETE Request.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "member_id",
          "in": "path",
          "required": true,
          "description": "The ID of the member to delete.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "destructive",
      "cacheable": false,
      "responseCodes": [
        "204",
        "401",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "member_id"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "delete",
      "pagination": null,
      "related": [
        {
          "id": "List an Organization's SCIM Members",
          "role": "list-parent"
        },
        {
          "id": "Provision a New Organization Member",
          "role": "list-parent"
        },
        {
          "id": "Query an Individual Organization Member",
          "role": "same-resource"
        },
        {
          "id": "Update an Organization Member's Attributes",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Retrieve the custom integrations created by an organization",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/sentry-apps/",
      "tag": "Integration",
      "summary": "Retrieve the custom integrations for an organization",
      "description": "Retrieve the custom integrations for an organization",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": null,
      "related": []
    },
    {
      "id": "Retrieve Release Health Session Statistics",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/sessions/",
      "tag": "Releases",
      "summary": "Returns a time series of release health session statistics for projects bound to an organization.\n\nThe interval and date range are subject to certain restrictions and rounding rules.\n\nThe date range is rounded to align with the interval, and is rounded to at least one\nhour. The interval can at most be one day and at least one hour currently. It has to cleanly\ndivide one day, for rounding reasons.\n\nBecause of technical limitations, this endpoint returns\nat most 10000 data points. For example, if you select a 90 day window grouped by releases,\nyou will see at most ` + "`" + `floor(10k / (90 + 1)) = 109` + "`" + ` releases. To get more results, reduce the\n` + "`" + `statsPeriod` + "`" + `.",
      "description": "Returns a time series of release health session statistics for projects bound to an organization.\n\nThe interval and date range are subject to certain restrictions and rounding rules.\n\nThe date range is rounded to align with the interval, and is rounded to at least one\nhour. The interval can at most be one day and at least one hour currently. It has to cleanly\ndivide one day, for rounding reasons.\n\nBecause of technical limitations, this endpoint returns\nat most 10000 data points. For example, if you select a 90 day window grouped by releases,\nyou will see at most ` + "`" + `floor(10k / (90 + 1)) = 109` + "`" + ` releases. To get more results, reduce the\n` + "`" + `statsPeriod` + "`" + `.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "field",
          "in": "query",
          "required": true,
          "description": "The list of fields to query.\n\nThe available fields are\n- ` + "`" + `sum(session)` + "`" + `\n- ` + "`" + `count_unique(user)` + "`" + `\n- ` + "`" + `avg` + "`" + `, ` + "`" + `p50` + "`" + `, ` + "`" + `p75` + "`" + `, ` + "`" + `p90` + "`" + `, ` + "`" + `p95` + "`" + `, ` + "`" + `p99` + "`" + `, ` + "`" + `max` + "`" + ` applied to ` + "`" + `session.duration` + "`" + `. For example, ` + "`" + `p99(session.duration)` + "`" + `. Session duration is [no longer being recorded](https://github.com/getsentry/sentry/discussions/42716) as of on Jan 12, 2023. Returned data may be incomplete.\n- ` + "`" + `crash_rate` + "`" + `, ` + "`" + `crash_free_rate` + "`" + ` applied to ` + "`" + `user` + "`" + ` or ` + "`" + `session` + "`" + `. For example, ` + "`" + `crash_free_rate(user)` + "`" + `\n",
          "type": "array"
        },
        {
          "name": "start",
          "in": "query",
          "required": false,
          "description": "The start of the period of time for the query, expected in ISO-8601 format. For example, ` + "`" + `2001-12-14T12:34:56.7890` + "`" + `.",
          "type": "string"
        },
        {
          "name": "end",
          "in": "query",
          "required": false,
          "description": "The end of the period of time for the query, expected in ISO-8601 format. For example, ` + "`" + `2001-12-14T12:34:56.7890` + "`" + `.",
          "type": "string"
        },
        {
          "name": "environment",
          "in": "query",
          "required": false,
          "description": "The name of environments to filter by.",
          "type": "array"
        },
        {
          "name": "statsPeriod",
          "in": "query",
          "required": false,
          "description": "The period of time for the query, will override the start & end parameters, a number followed by one of:\n- ` + "`" + `d` + "`" + ` for days\n- ` + "`" + `h` + "`" + ` for hours\n- ` + "`" + `m` + "`" + ` for minutes\n- ` + "`" + `s` + "`" + ` for seconds\n- ` + "`" + `w` + "`" + ` for weeks\n\nFor example, ` + "`" + `24h` + "`" + `, to mean query data starting from 24 hours ago to now.",
          "type": "string"
        },
        {
          "name": "project",
          "in": "query",
          "required": false,
          "description": "The IDs of projects to filter by. ` + "`" + `-1` + "`" + ` means all available projects.\nFor example, the following are valid parameters:\n- ` + "`" + `/?project=1234&project=56789` + "`" + `\n- ` + "`" + `/?project=-1` + "`" + `\n",
          "type": "array"
        },
        {
          "name": "per_page",
          "in": "query",
          "required": false,
          "description": "The number of groups to return per request.",
          "type": "integer"
        },
        {
          "name": "interval",
          "in": "query",
          "required": false,
          "description": "Resolution of the time series, given in the same format as ` + "`" + `statsPeriod` + "`" + `.\n\nThe default and\n        the minimum interval is ` + "`" + `1h` + "`" + `.",
          "type": "string"
        },
        {
          "name": "groupBy",
          "in": "query",
          "required": false,
          "description": "The list of properties to group by.\n\nThe available groupBy conditions are ` + "`" + `project` + "`" + `,\n        ` + "`" + `release` + "`" + `, ` + "`" + `environment` + "`" + ` and ` + "`" + `session.status` + "`" + `.",
          "type": "array"
        },
        {
          "name": "orderBy",
          "in": "query",
          "required": false,
          "description": "An optional field to order by, which must be one of the fields provided in ` + "`" + `field` + "`" + `. Use ` + "`" + `-` + "`" + `\n        for descending order, for example, ` + "`" + `-sum(session)` + "`" + `",
          "type": "string"
        },
        {
          "name": "includeTotals",
          "in": "query",
          "required": false,
          "description": "Specify ` + "`" + `0` + "`" + ` to exclude totals from the response. The default is ` + "`" + `1` + "`" + `",
          "type": "integer"
        },
        {
          "name": "includeSeries",
          "in": "query",
          "required": false,
          "description": "Specify ` + "`" + `0` + "`" + ` to exclude series from the response. The default is ` + "`" + `1` + "`" + `",
          "type": "integer"
        },
        {
          "name": "query",
          "in": "query",
          "required": false,
          "description": "Filters results by using [query syntax](/product/sentry-basics/search/).\n\nExample: ` + "`" + `query=(transaction:foo AND release:abc) OR (transaction:[bar,baz] AND release:def)` + "`" + `\n",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "400",
        "401"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug"
        ],
        "query": [
          "field",
          "start",
          "end",
          "environment",
          "statsPeriod",
          "project",
          "per_page",
          "interval",
          "groupBy",
          "orderBy",
          "includeTotals",
          "includeSeries",
          "query"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": null,
      "related": []
    },
    {
      "id": "Resolve a Short ID",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/shortids/{issue_id}/",
      "tag": "Organizations",
      "summary": "Resolve a short ID to the project slug and group details.",
      "description": "Resolve a short ID to the project slug and group details.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "issue_id",
          "in": "path",
          "required": true,
          "description": "The short ID of the issue to resolve.",
          "type": "string"
        },
        {
          "name": "collapse",
          "in": "query",
          "required": false,
          "description": "Fields to remove from the response to improve query performance.",
          "type": "array"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "400",
        "401",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "issue_id"
        ],
        "query": [
          "collapse"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": null,
      "related": []
    },
    {
      "id": "Retrieve an Organization's Events Count by Project",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/stats-summary/",
      "tag": "Organizations",
      "summary": "Query summarized event counts by project for your Organization. Also see https://docs.sentry.io/api/organizations/retrieve-event-counts-for-an-organization-v2/ for reference.",
      "description": "Query summarized event counts by project for your Organization. Also see https://docs.sentry.io/api/organizations/retrieve-event-counts-for-an-organization-v2/ for reference.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "field",
          "in": "query",
          "required": true,
          "description": "the ` + "`" + `sum(quantity)` + "`" + ` field is bytes for attachments, and all others the 'event' count for those types of events.\n\n` + "`" + `sum(times_seen)` + "`" + ` sums the number of times an event has been seen. For 'normal' event types, this will be equal to ` + "`" + `sum(quantity)` + "`" + ` for now. For sessions, quantity will sum the total number of events seen in a session, while ` + "`" + `times_seen` + "`" + ` will be the unique number of sessions. and for attachments, ` + "`" + `times_seen` + "`" + ` will be the total number of attachments, while quantity will be the total sum of attachment bytes.\n\n* ` + "`" + `sum(quantity)` + "`" + `\n* ` + "`" + `sum(times_seen)` + "`" + `",
          "type": "string"
        },
        {
          "name": "statsPeriod",
          "in": "query",
          "required": false,
          "description": "This defines the range of the time series, relative to now. The range is given in a ` + "`" + `<number><unit>` + "`" + ` format. For example ` + "`" + `1d` + "`" + ` for a one day range. Possible units are ` + "`" + `m` + "`" + ` for minutes, ` + "`" + `h` + "`" + ` for hours, ` + "`" + `d` + "`" + ` for days and ` + "`" + `w` + "`" + ` for weeks. You must either provide a ` + "`" + `statsPeriod` + "`" + `, or a ` + "`" + `start` + "`" + ` and ` + "`" + `end` + "`" + `.",
          "type": "string"
        },
        {
          "name": "interval",
          "in": "query",
          "required": false,
          "description": "This is the resolution of the time series, given in the same format as ` + "`" + `statsPeriod` + "`" + `. The default resolution is ` + "`" + `1h` + "`" + ` and the minimum resolution is currently restricted to ` + "`" + `1h` + "`" + ` as well. Intervals larger than ` + "`" + `1d` + "`" + ` are not supported, and the interval has to cleanly divide one day.",
          "type": "string"
        },
        {
          "name": "start",
          "in": "query",
          "required": false,
          "description": "This defines the start of the time series range as an explicit datetime, either in UTC ISO8601 or epoch seconds. Use along with ` + "`" + `end` + "`" + ` instead of ` + "`" + `statsPeriod` + "`" + `.",
          "type": "string"
        },
        {
          "name": "end",
          "in": "query",
          "required": false,
          "description": "This defines the inclusive end of the time series range as an explicit datetime, either in UTC ISO8601 or epoch seconds. Use along with ` + "`" + `start` + "`" + ` instead of ` + "`" + `statsPeriod` + "`" + `.",
          "type": "string"
        },
        {
          "name": "project",
          "in": "query",
          "required": false,
          "description": "The ID of the projects to filter by.",
          "type": "array"
        },
        {
          "name": "category",
          "in": "query",
          "required": false,
          "description": "If filtering by attachments, you cannot filter by any other category due to quantity values becoming nonsensical (combining bytes and event counts).\n\nIf filtering by ` + "`" + `error` + "`" + `, it will automatically add ` + "`" + `default` + "`" + ` and ` + "`" + `security` + "`" + ` as we currently roll those two categories into ` + "`" + `error` + "`" + ` for displaying.\n\n* ` + "`" + `error` + "`" + `\n* ` + "`" + `transaction` + "`" + `\n* ` + "`" + `attachment` + "`" + `\n* ` + "`" + `replays` + "`" + `\n* ` + "`" + `profiles` + "`" + `",
          "type": "string"
        },
        {
          "name": "outcome",
          "in": "query",
          "required": false,
          "description": "See https://docs.sentry.io/product/stats/ for more information on outcome statuses.\n\n* ` + "`" + `accepted` + "`" + `\n* ` + "`" + `filtered` + "`" + `\n* ` + "`" + `rate_limited` + "`" + `\n* ` + "`" + `invalid` + "`" + `\n* ` + "`" + `abuse` + "`" + `\n* ` + "`" + `client_discard` + "`" + `\n* ` + "`" + `cardinality_limited` + "`" + `",
          "type": "string"
        },
        {
          "name": "reason",
          "in": "query",
          "required": false,
          "description": "The reason field will contain why an event was filtered/dropped.",
          "type": "string"
        },
        {
          "name": "download",
          "in": "query",
          "required": false,
          "description": "Download the API response in as a csv file",
          "type": "boolean"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "401",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug"
        ],
        "query": [
          "field",
          "statsPeriod",
          "interval",
          "start",
          "end",
          "project",
          "category",
          "outcome",
          "reason",
          "download"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "search",
      "pagination": null,
      "related": []
    },
    {
      "id": "Retrieve Event Counts for an Organization (v2)",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/stats_v2/",
      "tag": "Organizations",
      "summary": "Query event counts for your Organization.\nSelect a field, define a date range, and group or filter by columns.",
      "description": "Query event counts for your Organization.\nSelect a field, define a date range, and group or filter by columns.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "groupBy",
          "in": "query",
          "required": true,
          "description": "can pass multiple groupBy parameters to group by multiple, e.g. ` + "`" + `groupBy=project&groupBy=outcome` + "`" + ` to group by multiple dimensions. Note that grouping by project can cause missing rows if the number of projects / interval is large. If you have a large number of projects, we recommend filtering and querying by them individually.Also note that grouping by projects does not currently support timeseries interval responses and will instead be a sum of the projectover the entire period specified.",
          "type": "array"
        },
        {
          "name": "field",
          "in": "query",
          "required": true,
          "description": "the ` + "`" + `sum(quantity)` + "`" + ` field is bytes for attachments, and all others the 'event' count for those types of events.\n\n` + "`" + `sum(times_seen)` + "`" + ` sums the number of times an event has been seen. For 'normal' event types, this will be equal to ` + "`" + `sum(quantity)` + "`" + ` for now. For sessions, quantity will sum the total number of events seen in a session, while ` + "`" + `times_seen` + "`" + ` will be the unique number of sessions. and for attachments, ` + "`" + `times_seen` + "`" + ` will be the total number of attachments, while quantity will be the total sum of attachment bytes.\n\n* ` + "`" + `sum(quantity)` + "`" + `\n* ` + "`" + `sum(times_seen)` + "`" + `",
          "type": "string"
        },
        {
          "name": "statsPeriod",
          "in": "query",
          "required": false,
          "description": "This defines the range of the time series, relative to now. The range is given in a ` + "`" + `<number><unit>` + "`" + ` format. For example ` + "`" + `1d` + "`" + ` for a one day range. Possible units are ` + "`" + `m` + "`" + ` for minutes, ` + "`" + `h` + "`" + ` for hours, ` + "`" + `d` + "`" + ` for days and ` + "`" + `w` + "`" + ` for weeks. You must either provide a ` + "`" + `statsPeriod` + "`" + `, or a ` + "`" + `start` + "`" + ` and ` + "`" + `end` + "`" + `.",
          "type": "string"
        },
        {
          "name": "interval",
          "in": "query",
          "required": false,
          "description": "This is the resolution of the time series, given in the same format as ` + "`" + `statsPeriod` + "`" + `. The default resolution is ` + "`" + `1h` + "`" + ` and the minimum resolution is currently restricted to ` + "`" + `1h` + "`" + ` as well. Intervals larger than ` + "`" + `1d` + "`" + ` are not supported, and the interval has to cleanly divide one day.",
          "type": "string"
        },
        {
          "name": "start",
          "in": "query",
          "required": false,
          "description": "This defines the start of the time series range as an explicit datetime, either in UTC ISO8601 or epoch seconds. Use along with ` + "`" + `end` + "`" + ` instead of ` + "`" + `statsPeriod` + "`" + `.",
          "type": "string"
        },
        {
          "name": "end",
          "in": "query",
          "required": false,
          "description": "This defines the inclusive end of the time series range as an explicit datetime, either in UTC ISO8601 or epoch seconds. Use along with ` + "`" + `start` + "`" + ` instead of ` + "`" + `statsPeriod` + "`" + `.",
          "type": "string"
        },
        {
          "name": "project",
          "in": "query",
          "required": false,
          "description": "The ID of the projects to filter by.\n\nUse ` + "`" + `-1` + "`" + ` to include all accessible projects.",
          "type": "array"
        },
        {
          "name": "category",
          "in": "query",
          "required": false,
          "description": "Filter by data category. Each category represents a different type of data:\n\n- ` + "`" + `error` + "`" + `: Error events (includes ` + "`" + `default` + "`" + ` and ` + "`" + `security` + "`" + ` categories)\n- ` + "`" + `transaction` + "`" + `: Transaction events\n- ` + "`" + `attachment` + "`" + `: File attachments (note: cannot be combined with other categories since quantity represents bytes)\n- ` + "`" + `replay` + "`" + `: Session replay events\n- ` + "`" + `profile` + "`" + `: Performance profiles\n- ` + "`" + `profile_duration` + "`" + `: Profile duration data (note: cannot be combined with other categories since quantity represents milliseconds)\n- ` + "`" + `profile_duration_ui` + "`" + `: Profile duration (UI) data (note: cannot be combined with other categories since quantity represents milliseconds)\n- ` + "`" + `profile_chunk` + "`" + `: Profile chunk data\n- ` + "`" + `profile_chunk_ui` + "`" + `: Profile chunk (UI) data\n- ` + "`" + `monitor` + "`" + `: Cron monitor events\n\n* ` + "`" + `error` + "`" + `\n* ` + "`" + `transaction` + "`" + `\n* ` + "`" + `attachment` + "`" + `\n* ` + "`" + `replay` + "`" + `\n* ` + "`" + `profile` + "`" + `\n* ` + "`" + `profile_duration` + "`" + `\n* ` + "`" + `profile_duration_ui` + "`" + `\n* ` + "`" + `profile_chunk` + "`" + `\n* ` + "`" + `profile_chunk_ui` + "`" + `\n* ` + "`" + `monitor` + "`" + `",
          "type": "string"
        },
        {
          "name": "outcome",
          "in": "query",
          "required": false,
          "description": "See https://docs.sentry.io/product/stats/ for more information on outcome statuses.\n\n* ` + "`" + `accepted` + "`" + `\n* ` + "`" + `filtered` + "`" + `\n* ` + "`" + `rate_limited` + "`" + `\n* ` + "`" + `invalid` + "`" + `\n* ` + "`" + `abuse` + "`" + `\n* ` + "`" + `client_discard` + "`" + `\n* ` + "`" + `cardinality_limited` + "`" + `",
          "type": "string"
        },
        {
          "name": "reason",
          "in": "query",
          "required": false,
          "description": "The reason field will contain why an event was filtered/dropped.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "401",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug"
        ],
        "query": [
          "groupBy",
          "field",
          "statsPeriod",
          "interval",
          "start",
          "end",
          "project",
          "category",
          "outcome",
          "reason"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "search",
      "pagination": null,
      "related": []
    },
    {
      "id": "List an Organization's Teams",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/teams/",
      "tag": "Teams",
      "summary": "Returns a list of teams bound to a organization.",
      "description": "Returns a list of teams bound to a organization.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "detailed",
          "in": "query",
          "required": false,
          "description": "\nSpecify ` + "`" + `\"0\"` + "`" + ` to return team details that do not include projects.\n",
          "type": "string"
        },
        {
          "name": "cursor",
          "in": "query",
          "required": false,
          "description": "A pointer to the last object fetched and its sort order; used to retrieve the next or previous results.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug"
        ],
        "query": [
          "detailed",
          "cursor"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": {
        "style": "cursor",
        "cursorParam": "cursor"
      },
      "related": [
        {
          "id": "Create a New Team",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Create a New Team",
      "method": "POST",
      "path": "/api/0/organizations/{organization_id_or_slug}/teams/",
      "tag": "Teams",
      "summary": "Create a new team bound to an organization. Requires at least one of the ` + "`" + `name` + "`" + `\nor ` + "`" + `slug` + "`" + ` body params to be set.",
      "description": "Create a new team bound to an organization. Requires at least one of the ` + "`" + `name` + "`" + `\nor ` + "`" + `slug` + "`" + ` body params to be set.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        }
      ],
      "hasRequestBody": true,
      "risk": "write",
      "cacheable": false,
      "responseCodes": [
        "201",
        "400",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug"
        ],
        "query": [],
        "headers": [],
        "body": true
      },
      "graphql": null,
      "kind": "create",
      "pagination": null,
      "related": [
        {
          "id": "List an Organization's Teams",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "List a User's Teams for an Organization",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/user-teams/",
      "tag": "Teams",
      "summary": "Returns a list of teams the user has access to in the specified organization.\nNote that this endpoint is restricted to [user auth tokens](https://docs.sentry.io/account/auth-tokens/#user-auth-tokens).",
      "description": "Returns a list of teams the user has access to in the specified organization.\nNote that this endpoint is restricted to [user auth tokens](https://docs.sentry.io/account/auth-tokens/#user-auth-tokens).",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "400",
        "403"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": null,
      "related": []
    },
    {
      "id": "Fetch Alerts",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/workflows/",
      "tag": "Monitors",
      "summary": "⚠️ This endpoint is currently in **beta** and may be subject to change. It is supported by [New Monitors and Alerts](/product/new-monitors-and-alerts/) and may not be viewable in the UI today.\n\nReturns a list of alerts for a given organization",
      "description": "⚠️ This endpoint is currently in **beta** and may be subject to change. It is supported by [New Monitors and Alerts](/product/new-monitors-and-alerts/) and may not be viewable in the UI today.\n\nReturns a list of alerts for a given organization",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "sortBy",
          "in": "query",
          "required": false,
          "description": "The field to sort results by. If not specified, the results are sorted by id.\n\nAvailable fields are:\n- ` + "`" + `name` + "`" + `\n- ` + "`" + `id` + "`" + `\n- ` + "`" + `dateCreated` + "`" + `\n- ` + "`" + `dateUpdated` + "`" + `\n- ` + "`" + `connectedDetectors` + "`" + `\n- ` + "`" + `actions` + "`" + `\n- ` + "`" + `priorityDetector` + "`" + `\n\nPrefix with ` + "`" + `-` + "`" + ` to sort in descending order.\n    ",
          "type": "string"
        },
        {
          "name": "query",
          "in": "query",
          "required": false,
          "description": "An optional search query for filtering alerts.",
          "type": "string"
        },
        {
          "name": "id",
          "in": "query",
          "required": false,
          "description": "The ID of the alert you'd like to query.",
          "type": "array"
        },
        {
          "name": "project",
          "in": "query",
          "required": false,
          "description": "The IDs of projects to filter by. ` + "`" + `-1` + "`" + ` means all available projects.\nFor example, the following are valid parameters:\n- ` + "`" + `/?project=1234&project=56789` + "`" + `\n- ` + "`" + `/?project=-1` + "`" + `\n",
          "type": "array"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "400",
        "401",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug"
        ],
        "query": [
          "sortBy",
          "query",
          "id",
          "project"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": null,
      "related": [
        {
          "id": "Create an Alert for an Organization",
          "role": "same-resource"
        },
        {
          "id": "Mutate an Organization's Alerts",
          "role": "same-resource"
        },
        {
          "id": "Bulk Delete Alerts",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Create an Alert for an Organization",
      "method": "POST",
      "path": "/api/0/organizations/{organization_id_or_slug}/workflows/",
      "tag": "Monitors",
      "summary": "⚠️ This endpoint is currently in **beta** and may be subject to change. It is supported by [New Monitors and Alerts](/product/new-monitors-and-alerts/) and may not be viewable in the UI today.\n\nCreates an alert for an organization",
      "description": "⚠️ This endpoint is currently in **beta** and may be subject to change. It is supported by [New Monitors and Alerts](/product/new-monitors-and-alerts/) and may not be viewable in the UI today.\n\nCreates an alert for an organization",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        }
      ],
      "hasRequestBody": true,
      "risk": "write",
      "cacheable": false,
      "responseCodes": [
        "201",
        "400",
        "401",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug"
        ],
        "query": [],
        "headers": [],
        "body": true
      },
      "graphql": null,
      "kind": "create",
      "pagination": null,
      "related": [
        {
          "id": "Fetch Alerts",
          "role": "same-resource"
        },
        {
          "id": "Mutate an Organization's Alerts",
          "role": "same-resource"
        },
        {
          "id": "Bulk Delete Alerts",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Mutate an Organization's Alerts",
      "method": "PUT",
      "path": "/api/0/organizations/{organization_id_or_slug}/workflows/",
      "tag": "Monitors",
      "summary": "⚠️ This endpoint is currently in **beta** and may be subject to change. It is supported by [New Monitors and Alerts](/product/new-monitors-and-alerts/) and may not be viewable in the UI today.\n\nBulk enable or disable alerts for a given Organization",
      "description": "⚠️ This endpoint is currently in **beta** and may be subject to change. It is supported by [New Monitors and Alerts](/product/new-monitors-and-alerts/) and may not be viewable in the UI today.\n\nBulk enable or disable alerts for a given Organization",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "query",
          "in": "query",
          "required": false,
          "description": "An optional search query for filtering alerts.",
          "type": "string"
        },
        {
          "name": "id",
          "in": "query",
          "required": false,
          "description": "The ID of the alert you'd like to query.",
          "type": "array"
        },
        {
          "name": "project",
          "in": "query",
          "required": false,
          "description": "The IDs of projects to filter by. ` + "`" + `-1` + "`" + ` means all available projects.\nFor example, the following are valid parameters:\n- ` + "`" + `/?project=1234&project=56789` + "`" + `\n- ` + "`" + `/?project=-1` + "`" + `\n",
          "type": "array"
        }
      ],
      "hasRequestBody": true,
      "risk": "write",
      "cacheable": false,
      "responseCodes": [
        "200",
        "400",
        "401",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug"
        ],
        "query": [
          "query",
          "id",
          "project"
        ],
        "headers": [],
        "body": true
      },
      "graphql": null,
      "kind": "update",
      "pagination": null,
      "related": [
        {
          "id": "Fetch Alerts",
          "role": "same-resource"
        },
        {
          "id": "Create an Alert for an Organization",
          "role": "same-resource"
        },
        {
          "id": "Bulk Delete Alerts",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Bulk Delete Alerts",
      "method": "DELETE",
      "path": "/api/0/organizations/{organization_id_or_slug}/workflows/",
      "tag": "Monitors",
      "summary": "⚠️ This endpoint is currently in **beta** and may be subject to change. It is supported by [New Monitors and Alerts](/product/new-monitors-and-alerts/) and may not be viewable in the UI today.\n\nBulk delete alerts for a given organization",
      "description": "⚠️ This endpoint is currently in **beta** and may be subject to change. It is supported by [New Monitors and Alerts](/product/new-monitors-and-alerts/) and may not be viewable in the UI today.\n\nBulk delete alerts for a given organization",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "query",
          "in": "query",
          "required": false,
          "description": "An optional search query for filtering alerts.",
          "type": "string"
        },
        {
          "name": "id",
          "in": "query",
          "required": false,
          "description": "The ID of the alert you'd like to query.",
          "type": "array"
        },
        {
          "name": "project",
          "in": "query",
          "required": false,
          "description": "The IDs of projects to filter by. ` + "`" + `-1` + "`" + ` means all available projects.\nFor example, the following are valid parameters:\n- ` + "`" + `/?project=1234&project=56789` + "`" + `\n- ` + "`" + `/?project=-1` + "`" + `\n",
          "type": "array"
        }
      ],
      "hasRequestBody": false,
      "risk": "destructive",
      "cacheable": false,
      "responseCodes": [
        "200",
        "204",
        "400",
        "401",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug"
        ],
        "query": [
          "query",
          "id",
          "project"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "delete",
      "pagination": null,
      "related": [
        {
          "id": "Fetch Alerts",
          "role": "same-resource"
        },
        {
          "id": "Create an Alert for an Organization",
          "role": "same-resource"
        },
        {
          "id": "Mutate an Organization's Alerts",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Fetch an Alert",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/workflows/{workflow_id}/",
      "tag": "Monitors",
      "summary": "⚠️ This endpoint is currently in **beta** and may be subject to change. It is supported by [New Monitors and Alerts](/product/new-monitors-and-alerts/) and may not be viewable in the UI today.\n\nReturns an alert.",
      "description": "⚠️ This endpoint is currently in **beta** and may be subject to change. It is supported by [New Monitors and Alerts](/product/new-monitors-and-alerts/) and may not be viewable in the UI today.\n\nReturns an alert.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "workflow_id",
          "in": "path",
          "required": true,
          "description": "The ID of the alert you'd like to query.",
          "type": "integer"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "400",
        "401",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "workflow_id"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": null,
      "related": [
        {
          "id": "Update an Alert by ID",
          "role": "same-resource"
        },
        {
          "id": "Delete an Alert",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Update an Alert by ID",
      "method": "PUT",
      "path": "/api/0/organizations/{organization_id_or_slug}/workflows/{workflow_id}/",
      "tag": "Monitors",
      "summary": "⚠️ This endpoint is currently in **beta** and may be subject to change. It is supported by [New Monitors and Alerts](/product/new-monitors-and-alerts/) and may not be viewable in the UI today.\n\nUpdates an alert.",
      "description": "⚠️ This endpoint is currently in **beta** and may be subject to change. It is supported by [New Monitors and Alerts](/product/new-monitors-and-alerts/) and may not be viewable in the UI today.\n\nUpdates an alert.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "workflow_id",
          "in": "path",
          "required": true,
          "description": "The ID of the alert you'd like to query.",
          "type": "integer"
        }
      ],
      "hasRequestBody": true,
      "risk": "write",
      "cacheable": false,
      "responseCodes": [
        "200",
        "400",
        "401",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "workflow_id"
        ],
        "query": [],
        "headers": [],
        "body": true
      },
      "graphql": null,
      "kind": "update",
      "pagination": null,
      "related": [
        {
          "id": "Fetch an Alert",
          "role": "same-resource"
        },
        {
          "id": "Delete an Alert",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Delete an Alert",
      "method": "DELETE",
      "path": "/api/0/organizations/{organization_id_or_slug}/workflows/{workflow_id}/",
      "tag": "Monitors",
      "summary": "⚠️ This endpoint is currently in **beta** and may be subject to change. It is supported by [New Monitors and Alerts](/product/new-monitors-and-alerts/) and may not be viewable in the UI today.\n\nDeletes an alert.",
      "description": "⚠️ This endpoint is currently in **beta** and may be subject to change. It is supported by [New Monitors and Alerts](/product/new-monitors-and-alerts/) and may not be viewable in the UI today.\n\nDeletes an alert.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "workflow_id",
          "in": "path",
          "required": true,
          "description": "The ID of the alert you'd like to query.",
          "type": "integer"
        }
      ],
      "hasRequestBody": false,
      "risk": "destructive",
      "cacheable": false,
      "responseCodes": [
        "204",
        "400",
        "401",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "workflow_id"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "delete",
      "pagination": null,
      "related": [
        {
          "id": "Fetch an Alert",
          "role": "same-resource"
        },
        {
          "id": "Update an Alert by ID",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Retrieve a Project",
      "method": "GET",
      "path": "/api/0/projects/{organization_id_or_slug}/{project_id_or_slug}/",
      "tag": "Projects",
      "summary": "Return details on an individual project.",
      "description": "Return details on an individual project.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the project the resource belongs to.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "project_id_or_slug"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": null,
      "related": [
        {
          "id": "Update a Project",
          "role": "same-resource"
        },
        {
          "id": "Delete a Project",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Update a Project",
      "method": "PUT",
      "path": "/api/0/projects/{organization_id_or_slug}/{project_id_or_slug}/",
      "tag": "Projects",
      "summary": "Update various attributes and configurable settings for the given project.\n\nNote that solely having the **` + "`" + `project:read` + "`" + `** scope restricts updatable settings to\n` + "`" + `isBookmarked` + "`" + `, ` + "`" + `autofixAutomationTuning` + "`" + `, ` + "`" + `seerScannerAutomation` + "`" + `,\n` + "`" + `preprodSizeStatusChecksEnabled` + "`" + `, ` + "`" + `preprodSizeStatusChecksRules` + "`" + `,\n` + "`" + `preprodSizeEnabledQuery` + "`" + `, ` + "`" + `preprodDistributionEnabledQuery` + "`" + `,\n` + "`" + `preprodSizeEnabledByCustomer` + "`" + `, ` + "`" + `preprodDistributionEnabledByCustomer` + "`" + `,\nand ` + "`" + `preprodDistributionPrCommentsEnabledByCustomer` + "`" + `.",
      "description": "Update various attributes and configurable settings for the given project.\n\nNote that solely having the **` + "`" + `project:read` + "`" + `** scope restricts updatable settings to\n` + "`" + `isBookmarked` + "`" + `, ` + "`" + `autofixAutomationTuning` + "`" + `, ` + "`" + `seerScannerAutomation` + "`" + `,\n` + "`" + `preprodSizeStatusChecksEnabled` + "`" + `, ` + "`" + `preprodSizeStatusChecksRules` + "`" + `,\n` + "`" + `preprodSizeEnabledQuery` + "`" + `, ` + "`" + `preprodDistributionEnabledQuery` + "`" + `,\n` + "`" + `preprodSizeEnabledByCustomer` + "`" + `, ` + "`" + `preprodDistributionEnabledByCustomer` + "`" + `,\nand ` + "`" + `preprodDistributionPrCommentsEnabledByCustomer` + "`" + `.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the project the resource belongs to.",
          "type": "string"
        }
      ],
      "hasRequestBody": true,
      "risk": "write",
      "cacheable": false,
      "responseCodes": [
        "200",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "project_id_or_slug"
        ],
        "query": [],
        "headers": [],
        "body": true
      },
      "graphql": null,
      "kind": "update",
      "pagination": null,
      "related": [
        {
          "id": "Retrieve a Project",
          "role": "same-resource"
        },
        {
          "id": "Delete a Project",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Delete a Project",
      "method": "DELETE",
      "path": "/api/0/projects/{organization_id_or_slug}/{project_id_or_slug}/",
      "tag": "Projects",
      "summary": "Schedules a project for deletion.\n\nDeletion happens asynchronously and therefore is not immediate. However once deletion has\nbegun the state of a project changes and will be hidden from most public views.",
      "description": "Schedules a project for deletion.\n\nDeletion happens asynchronously and therefore is not immediate. However once deletion has\nbegun the state of a project changes and will be hidden from most public views.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the project the resource belongs to.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "destructive",
      "cacheable": false,
      "responseCodes": [
        "204",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "project_id_or_slug"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "delete",
      "pagination": null,
      "related": [
        {
          "id": "Retrieve a Project",
          "role": "same-resource"
        },
        {
          "id": "Update a Project",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "List a Project's Environments",
      "method": "GET",
      "path": "/api/0/projects/{organization_id_or_slug}/{project_id_or_slug}/environments/",
      "tag": "Environments",
      "summary": "Lists a project's environments.",
      "description": "Lists a project's environments.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the project the resource belongs to.",
          "type": "string"
        },
        {
          "name": "visibility",
          "in": "query",
          "required": false,
          "description": "The visibility of the environments to filter by. Defaults to ` + "`" + `visible` + "`" + `.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "400",
        "401",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "project_id_or_slug"
        ],
        "query": [
          "visibility"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": null,
      "related": []
    },
    {
      "id": "Retrieve a Project Environment",
      "method": "GET",
      "path": "/api/0/projects/{organization_id_or_slug}/{project_id_or_slug}/environments/{environment}/",
      "tag": "Environments",
      "summary": "Return details on a project environment.",
      "description": "Return details on a project environment.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the project the resource belongs to.",
          "type": "string"
        },
        {
          "name": "environment",
          "in": "path",
          "required": true,
          "description": "The name of the environment.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "401",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "project_id_or_slug",
          "environment"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": null,
      "related": [
        {
          "id": "Update a Project Environment",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Update a Project Environment",
      "method": "PUT",
      "path": "/api/0/projects/{organization_id_or_slug}/{project_id_or_slug}/environments/{environment}/",
      "tag": "Environments",
      "summary": "Update the visibility for a project environment.",
      "description": "Update the visibility for a project environment.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the project the resource belongs to.",
          "type": "string"
        },
        {
          "name": "environment",
          "in": "path",
          "required": true,
          "description": "The name of the environment.",
          "type": "string"
        }
      ],
      "hasRequestBody": true,
      "risk": "write",
      "cacheable": false,
      "responseCodes": [
        "200",
        "400",
        "401",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "project_id_or_slug",
          "environment"
        ],
        "query": [],
        "headers": [],
        "body": true
      },
      "graphql": null,
      "kind": "update",
      "pagination": null,
      "related": [
        {
          "id": "Retrieve a Project Environment",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "List a Project's Error Events",
      "method": "GET",
      "path": "/api/0/projects/{organization_id_or_slug}/{project_id_or_slug}/events/",
      "tag": "Events",
      "summary": "Return a list of events bound to a project.",
      "description": "Return a list of events bound to a project.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the project the resource belongs to.",
          "type": "string"
        },
        {
          "name": "statsPeriod",
          "in": "query",
          "required": false,
          "description": "The period of time for the query, will override the start & end parameters, a number followed by one of:\n- ` + "`" + `d` + "`" + ` for days\n- ` + "`" + `h` + "`" + ` for hours\n- ` + "`" + `m` + "`" + ` for minutes\n- ` + "`" + `s` + "`" + ` for seconds\n- ` + "`" + `w` + "`" + ` for weeks\n\nFor example, ` + "`" + `24h` + "`" + `, to mean query data starting from 24 hours ago to now.",
          "type": "string"
        },
        {
          "name": "start",
          "in": "query",
          "required": false,
          "description": "The start of the period of time for the query, expected in ISO-8601 format. For example, ` + "`" + `2001-12-14T12:34:56.7890` + "`" + `.",
          "type": "string"
        },
        {
          "name": "end",
          "in": "query",
          "required": false,
          "description": "The end of the period of time for the query, expected in ISO-8601 format. For example, ` + "`" + `2001-12-14T12:34:56.7890` + "`" + `.",
          "type": "string"
        },
        {
          "name": "cursor",
          "in": "query",
          "required": false,
          "description": "A pointer to the last object fetched and its sort order; used to retrieve the next or previous results.",
          "type": "string"
        },
        {
          "name": "full",
          "in": "query",
          "required": false,
          "description": "Specify true to include the full event body, including the stacktrace, in the event payload.",
          "type": "boolean"
        },
        {
          "name": "sample",
          "in": "query",
          "required": false,
          "description": "Return events in pseudo-random order. This is deterministic so an identical query will always return the same events in the same order.",
          "type": "boolean"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "401",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "project_id_or_slug"
        ],
        "query": [
          "statsPeriod",
          "start",
          "end",
          "cursor",
          "full",
          "sample"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": {
        "style": "cursor",
        "cursorParam": "cursor"
      },
      "related": []
    },
    {
      "id": "Debug Issues Related to Source Maps for a Given Event",
      "method": "GET",
      "path": "/api/0/projects/{organization_id_or_slug}/{project_id_or_slug}/events/{event_id}/source-map-debug/",
      "tag": "Events",
      "summary": "Return a list of source map errors for a given event.",
      "description": "Return a list of source map errors for a given event.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the project the resource belongs to.",
          "type": "string"
        },
        {
          "name": "event_id",
          "in": "path",
          "required": true,
          "description": "The ID of the event.",
          "type": "string"
        },
        {
          "name": "frame_idx",
          "in": "query",
          "required": true,
          "description": "Index of the frame that should be used for source map resolution.",
          "type": "integer"
        },
        {
          "name": "exception_idx",
          "in": "query",
          "required": true,
          "description": "Index of the exception that should be used for source map resolution.",
          "type": "integer"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "401",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "project_id_or_slug",
          "event_id"
        ],
        "query": [
          "frame_idx",
          "exception_idx"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": null,
      "related": []
    },
    {
      "id": "List a Project's Data Filters",
      "method": "GET",
      "path": "/api/0/projects/{organization_id_or_slug}/{project_id_or_slug}/filters/",
      "tag": "Projects",
      "summary": "Retrieve a list of filters for a given project.\n` + "`" + `active` + "`" + ` will be either a boolean or a list for the legacy browser filters.",
      "description": "Retrieve a list of filters for a given project.\n` + "`" + `active` + "`" + ` will be either a boolean or a list for the legacy browser filters.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the project the resource belongs to.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "403"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "project_id_or_slug"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": null,
      "related": []
    },
    {
      "id": "Update an Inbound Data Filter",
      "method": "PUT",
      "path": "/api/0/projects/{organization_id_or_slug}/{project_id_or_slug}/filters/{filter_id}/",
      "tag": "Projects",
      "summary": "Update various inbound data filters for a project.",
      "description": "Update various inbound data filters for a project.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the project the resource belongs to.",
          "type": "string"
        },
        {
          "name": "filter_id",
          "in": "path",
          "required": true,
          "description": "The type of filter toggle to update. The options are:\n- ` + "`" + `browser-extensions` + "`" + ` - Filter out errors known to be caused by browser extensions.\n- ` + "`" + `localhost` + "`" + ` - Filter out events coming from localhost. This applies to both IPv4 (` + "`" + `` + "`" + `127.0.0.1` + "`" + `` + "`" + `)\nand IPv6 (` + "`" + `` + "`" + `::1` + "`" + `` + "`" + `) addresses.\n- ` + "`" + `filtered-transaction` + "`" + ` - Filter out transactions for healthcheck and ping endpoints.\n- ` + "`" + `web-crawlers` + "`" + ` - Filter out known web crawlers. Some crawlers may execute pages in incompatible\nways which cause errors that are unlikely to be seen by a normal user.\n- ` + "`" + `legacy-browser` + "`" + ` - Filter out known errors from legacy browsers. Older browsers often give less\naccurate information, and while they may report valid issues, the context to understand them is\nincorrect or missing.\n",
          "type": "string"
        }
      ],
      "hasRequestBody": true,
      "risk": "write",
      "cacheable": false,
      "responseCodes": [
        "204",
        "400",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "project_id_or_slug",
          "filter_id"
        ],
        "query": [],
        "headers": [],
        "body": true
      },
      "graphql": null,
      "kind": "update",
      "pagination": null,
      "related": []
    },
    {
      "id": "List a Project's Client Keys",
      "method": "GET",
      "path": "/api/0/projects/{organization_id_or_slug}/{project_id_or_slug}/keys/",
      "tag": "Projects",
      "summary": "Return a list of client keys bound to a project.",
      "description": "Return a list of client keys bound to a project.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the project the resource belongs to.",
          "type": "string"
        },
        {
          "name": "cursor",
          "in": "query",
          "required": false,
          "description": "A pointer to the last object fetched and its sort order; used to retrieve the next or previous results.",
          "type": "string"
        },
        {
          "name": "status",
          "in": "query",
          "required": false,
          "description": "\nFilter client keys by ` + "`" + `active` + "`" + ` or ` + "`" + `inactive` + "`" + `. Defaults to returning all\nkeys if not specified.\n",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "400",
        "403"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "project_id_or_slug"
        ],
        "query": [
          "cursor",
          "status"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": {
        "style": "cursor",
        "cursorParam": "cursor"
      },
      "related": [
        {
          "id": "Create a New Client Key",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Create a New Client Key",
      "method": "POST",
      "path": "/api/0/projects/{organization_id_or_slug}/{project_id_or_slug}/keys/",
      "tag": "Projects",
      "summary": "Create a new client key bound to a project.  The key's secret and public key\nare generated by the server.",
      "description": "Create a new client key bound to a project.  The key's secret and public key\nare generated by the server.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the project the resource belongs to.",
          "type": "string"
        }
      ],
      "hasRequestBody": true,
      "risk": "write",
      "cacheable": false,
      "responseCodes": [
        "201",
        "400",
        "403"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "project_id_or_slug"
        ],
        "query": [],
        "headers": [],
        "body": true
      },
      "graphql": null,
      "kind": "create",
      "pagination": null,
      "related": [
        {
          "id": "List a Project's Client Keys",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Retrieve a Client Key",
      "method": "GET",
      "path": "/api/0/projects/{organization_id_or_slug}/{project_id_or_slug}/keys/{key_id}/",
      "tag": "Projects",
      "summary": "Return a client key bound to a project.",
      "description": "Return a client key bound to a project.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the project the resource belongs to.",
          "type": "string"
        },
        {
          "name": "key_id",
          "in": "path",
          "required": true,
          "description": "The ID of the client key",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "project_id_or_slug",
          "key_id"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": null,
      "related": [
        {
          "id": "Update a Client Key",
          "role": "same-resource"
        },
        {
          "id": "Delete a Client Key",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Update a Client Key",
      "method": "PUT",
      "path": "/api/0/projects/{organization_id_or_slug}/{project_id_or_slug}/keys/{key_id}/",
      "tag": "Projects",
      "summary": "Update various settings for a client key.",
      "description": "Update various settings for a client key.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the project the resource belongs to.",
          "type": "string"
        },
        {
          "name": "key_id",
          "in": "path",
          "required": true,
          "description": "The ID of the key to update.",
          "type": "string"
        }
      ],
      "hasRequestBody": true,
      "risk": "write",
      "cacheable": false,
      "responseCodes": [
        "200",
        "400",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "project_id_or_slug",
          "key_id"
        ],
        "query": [],
        "headers": [],
        "body": true
      },
      "graphql": null,
      "kind": "update",
      "pagination": null,
      "related": [
        {
          "id": "Retrieve a Client Key",
          "role": "same-resource"
        },
        {
          "id": "Delete a Client Key",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Delete a Client Key",
      "method": "DELETE",
      "path": "/api/0/projects/{organization_id_or_slug}/{project_id_or_slug}/keys/{key_id}/",
      "tag": "Projects",
      "summary": "Delete a client key for a given project.",
      "description": "Delete a client key for a given project.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the project the resource belongs to.",
          "type": "string"
        },
        {
          "name": "key_id",
          "in": "path",
          "required": true,
          "description": "The ID of the key to delete.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "destructive",
      "cacheable": false,
      "responseCodes": [
        "204",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "project_id_or_slug",
          "key_id"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "delete",
      "pagination": null,
      "related": [
        {
          "id": "Retrieve a Client Key",
          "role": "same-resource"
        },
        {
          "id": "Update a Client Key",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "List a Project's Organization Members",
      "method": "GET",
      "path": "/api/0/projects/{organization_id_or_slug}/{project_id_or_slug}/members/",
      "tag": "Projects",
      "summary": "Returns a list of active organization members that belong to any team assigned to the project.",
      "description": "Returns a list of active organization members that belong to any team assigned to the project.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the project the resource belongs to.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "400",
        "403"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "project_id_or_slug"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": null,
      "related": []
    },
    {
      "id": "Retrieve a Monitor for a Project",
      "method": "GET",
      "path": "/api/0/projects/{organization_id_or_slug}/{project_id_or_slug}/monitors/{monitor_id_or_slug}/",
      "tag": "Crons",
      "summary": "Retrieves details for a monitor.",
      "description": "Retrieves details for a monitor.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the project the resource belongs to.",
          "type": "string"
        },
        {
          "name": "monitor_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the monitor.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "401",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "project_id_or_slug",
          "monitor_id_or_slug"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": null,
      "related": [
        {
          "id": "Update a Monitor for a Project",
          "role": "same-resource"
        },
        {
          "id": "Delete a Monitor or Monitor Environments for a Project",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Update a Monitor for a Project",
      "method": "PUT",
      "path": "/api/0/projects/{organization_id_or_slug}/{project_id_or_slug}/monitors/{monitor_id_or_slug}/",
      "tag": "Crons",
      "summary": "Update a monitor.",
      "description": "Update a monitor.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the project the resource belongs to.",
          "type": "string"
        },
        {
          "name": "monitor_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the monitor.",
          "type": "string"
        }
      ],
      "hasRequestBody": true,
      "risk": "write",
      "cacheable": false,
      "responseCodes": [
        "200",
        "400",
        "401",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "project_id_or_slug",
          "monitor_id_or_slug"
        ],
        "query": [],
        "headers": [],
        "body": true
      },
      "graphql": null,
      "kind": "update",
      "pagination": null,
      "related": [
        {
          "id": "Retrieve a Monitor for a Project",
          "role": "same-resource"
        },
        {
          "id": "Delete a Monitor or Monitor Environments for a Project",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Delete a Monitor or Monitor Environments for a Project",
      "method": "DELETE",
      "path": "/api/0/projects/{organization_id_or_slug}/{project_id_or_slug}/monitors/{monitor_id_or_slug}/",
      "tag": "Crons",
      "summary": "Delete a monitor or monitor environments.",
      "description": "Delete a monitor or monitor environments.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the project the resource belongs to.",
          "type": "string"
        },
        {
          "name": "monitor_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the monitor.",
          "type": "string"
        },
        {
          "name": "environment",
          "in": "query",
          "required": false,
          "description": "The name of environments to filter by.",
          "type": "array"
        }
      ],
      "hasRequestBody": false,
      "risk": "destructive",
      "cacheable": false,
      "responseCodes": [
        "202",
        "401",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "project_id_or_slug",
          "monitor_id_or_slug"
        ],
        "query": [
          "environment"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "delete",
      "pagination": null,
      "related": [
        {
          "id": "Retrieve a Monitor for a Project",
          "role": "same-resource"
        },
        {
          "id": "Update a Monitor for a Project",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Retrieve Check-Ins for a Monitor by Project",
      "method": "GET",
      "path": "/api/0/projects/{organization_id_or_slug}/{project_id_or_slug}/monitors/{monitor_id_or_slug}/checkins/",
      "tag": "Crons",
      "summary": "Retrieve a list of check-ins for a monitor",
      "description": "Retrieve a list of check-ins for a monitor",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the project the resource belongs to.",
          "type": "string"
        },
        {
          "name": "monitor_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the monitor.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "401",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "project_id_or_slug",
          "monitor_id_or_slug"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": null,
      "related": []
    },
    {
      "id": "Retrieve Ownership Configuration for a Project",
      "method": "GET",
      "path": "/api/0/projects/{organization_id_or_slug}/{project_id_or_slug}/ownership/",
      "tag": "Projects",
      "summary": "Returns details on a project's ownership configuration.",
      "description": "Returns details on a project's ownership configuration.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the project the resource belongs to.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "project_id_or_slug"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": null,
      "related": [
        {
          "id": "Update Ownership Configuration for a Project",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Update Ownership Configuration for a Project",
      "method": "PUT",
      "path": "/api/0/projects/{organization_id_or_slug}/{project_id_or_slug}/ownership/",
      "tag": "Projects",
      "summary": "Updates ownership configurations for a project. Note that only the\nattributes submitted are modified.",
      "description": "Updates ownership configurations for a project. Note that only the\nattributes submitted are modified.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the project the resource belongs to.",
          "type": "string"
        }
      ],
      "hasRequestBody": true,
      "risk": "write",
      "cacheable": false,
      "responseCodes": [
        "202",
        "400"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "project_id_or_slug"
        ],
        "query": [],
        "headers": [],
        "body": true
      },
      "graphql": null,
      "kind": "update",
      "pagination": null,
      "related": [
        {
          "id": "Retrieve Ownership Configuration for a Project",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Retrieve Size Analysis status check rules for a project",
      "method": "GET",
      "path": "/api/0/projects/{organization_id_or_slug}/{project_id_or_slug}/preprod/size-analysis/status-check-rules/",
      "tag": "Mobile Builds",
      "summary": "Retrieve the current Size Analysis status check rules configured for a project.\n\nUse this endpoint after receiving a ` + "`" + `size_analysis.completed` + "`" + ` webhook when you\nwant external CI to evaluate the same Size Analysis status check thresholds that\nSentry uses. The endpoint returns the current project configuration, not a\nhistorical snapshot from when the webhook was emitted.\n\nThe response includes whether status check enforcement is enabled and the\nnormalized rule list Sentry uses when evaluating Size Analysis thresholds.\n\nThis endpoint requires a bearer token with ` + "`" + `project:read` + "`" + ` access. Project\ndistribution tokens are not supported.\n\nResponse notes:\n\n- ` + "`" + `enabled: false` + "`" + ` means status-check enforcement is disabled for the project.\n- ` + "`" + `rules: []` + "`" + ` means there are no configured thresholds to evaluate.\n- ` + "`" + `value` + "`" + ` is returned as a string. For ` + "`" + `absolute` + "`" + ` and ` + "`" + `absolute_diff` + "`" + `\n  measurements it is a byte value; for ` + "`" + `relative_diff` + "`" + ` it is a percentage.\n- ` + "`" + `filterQuery` + "`" + ` is the original configured filter string.\n- ` + "`" + `filters` + "`" + ` is the machine-readable version of ` + "`" + `filterQuery` + "`" + `.\n- ` + "`" + `filters: []` + "`" + ` means the rule has no filters and applies to all builds.\n- ` + "`" + `filters: null` + "`" + ` means the saved filter query could not be parsed; Sentry's\n  status check trigger treats that rule as non-matching.\n\nRule evaluation semantics:\n\n- Threshold comparisons are strict: a rule triggers only when the computed value\n  is greater than the configured threshold, not greater than or equal to it.\n- ` + "`" + `absolute_diff` + "`" + ` and ` + "`" + `relative_diff` + "`" + ` require a matching base metric/build.\n- ` + "`" + `relative_diff` + "`" + ` does not trigger when the base size is zero.\n- ` + "`" + `artifactType` + "`" + ` identifies the artifact scope the rule evaluates.\n  ` + "`" + `main_artifact` + "`" + `, ` + "`" + `watch_artifact` + "`" + `, ` + "`" + `android_dynamic_feature_artifact` + "`" + `,\n  and ` + "`" + `app_clip_artifact` + "`" + ` target their matching artifact metric.\n  ` + "`" + `all_artifacts` + "`" + ` evaluates all available artifact metrics.\n- Rule filters support the keys ` + "`" + `app_id` + "`" + `, ` + "`" + `git_head_ref` + "`" + `,\n  ` + "`" + `build_configuration_name` + "`" + `, and ` + "`" + `platform_name` + "`" + `.\n- Filter objects are combined with AND. Multiple ` + "`" + `conditions` + "`" + ` inside one\n  filter object are combined with OR.\n- Each condition uses ` + "`" + `values` + "`" + `; single-value operators still return a\n  one-item array.\n- Values in ` + "`" + `filters` + "`" + ` are decoded literal values for exact/simple operators,\n  not query syntax. For example, ` + "`" + `app_id:\\*com` + "`" + ` in ` + "`" + `filterQuery` + "`" + ` becomes\n  ` + "`" + `values: [\"*com\"]` + "`" + ` with ` + "`" + `operator: \"equals\"` + "`" + `.\n- The same key can appear in more than one filter object when positive and\n  negative conditions both exist; those filter objects are still combined with\n  AND.\n- Supported filter operators are ` + "`" + `equals` + "`" + `, ` + "`" + `notEquals` + "`" + `, ` + "`" + `in` + "`" + `, ` + "`" + `notIn` + "`" + `,\n  ` + "`" + `contains` + "`" + `, ` + "`" + `notContains` + "`" + `, ` + "`" + `startsWith` + "`" + `, ` + "`" + `notStartsWith` + "`" + `, ` + "`" + `endsWith` + "`" + `,\n  ` + "`" + `notEndsWith` + "`" + `, ` + "`" + `matches` + "`" + `, and ` + "`" + `notMatches` + "`" + `.\n- ` + "`" + `matches` + "`" + ` and ` + "`" + `notMatches` + "`" + ` values use Sentry wildcard pattern syntax, not\n  regular expressions. ` + "`" + `*` + "`" + ` matches zero or more characters, escaped ` + "`" + `\\*` + "`" + `\n  matches a literal asterisk, and a pattern without ` + "`" + `*` + "`" + ` is an exact match.\n- ` + "`" + `in` + "`" + ` and ` + "`" + `notIn` + "`" + ` are evaluated as one condition against all values, matching\n  Sentry's status check trigger behavior.\n- A rule applies only when the build metadata matches all filters. If a\n  referenced metadata key is missing, the filter does not match, even for\n  negated operators.",
      "description": "Retrieve the current Size Analysis status check rules configured for a project.\n\nUse this endpoint after receiving a ` + "`" + `size_analysis.completed` + "`" + ` webhook when you\nwant external CI to evaluate the same Size Analysis status check thresholds that\nSentry uses. The endpoint returns the current project configuration, not a\nhistorical snapshot from when the webhook was emitted.\n\nThe response includes whether status check enforcement is enabled and the\nnormalized rule list Sentry uses when evaluating Size Analysis thresholds.\n\nThis endpoint requires a bearer token with ` + "`" + `project:read` + "`" + ` access. Project\ndistribution tokens are not supported.\n\nResponse notes:\n\n- ` + "`" + `enabled: false` + "`" + ` means status-check enforcement is disabled for the project.\n- ` + "`" + `rules: []` + "`" + ` means there are no configured thresholds to evaluate.\n- ` + "`" + `value` + "`" + ` is returned as a string. For ` + "`" + `absolute` + "`" + ` and ` + "`" + `absolute_diff` + "`" + `\n  measurements it is a byte value; for ` + "`" + `relative_diff` + "`" + ` it is a percentage.\n- ` + "`" + `filterQuery` + "`" + ` is the original configured filter string.\n- ` + "`" + `filters` + "`" + ` is the machine-readable version of ` + "`" + `filterQuery` + "`" + `.\n- ` + "`" + `filters: []` + "`" + ` means the rule has no filters and applies to all builds.\n- ` + "`" + `filters: null` + "`" + ` means the saved filter query could not be parsed; Sentry's\n  status check trigger treats that rule as non-matching.\n\nRule evaluation semantics:\n\n- Threshold comparisons are strict: a rule triggers only when the computed value\n  is greater than the configured threshold, not greater than or equal to it.\n- ` + "`" + `absolute_diff` + "`" + ` and ` + "`" + `relative_diff` + "`" + ` require a matching base metric/build.\n- ` + "`" + `relative_diff` + "`" + ` does not trigger when the base size is zero.\n- ` + "`" + `artifactType` + "`" + ` identifies the artifact scope the rule evaluates.\n  ` + "`" + `main_artifact` + "`" + `, ` + "`" + `watch_artifact` + "`" + `, ` + "`" + `android_dynamic_feature_artifact` + "`" + `,\n  and ` + "`" + `app_clip_artifact` + "`" + ` target their matching artifact metric.\n  ` + "`" + `all_artifacts` + "`" + ` evaluates all available artifact metrics.\n- Rule filters support the keys ` + "`" + `app_id` + "`" + `, ` + "`" + `git_head_ref` + "`" + `,\n  ` + "`" + `build_configuration_name` + "`" + `, and ` + "`" + `platform_name` + "`" + `.\n- Filter objects are combined with AND. Multiple ` + "`" + `conditions` + "`" + ` inside one\n  filter object are combined with OR.\n- Each condition uses ` + "`" + `values` + "`" + `; single-value operators still return a\n  one-item array.\n- Values in ` + "`" + `filters` + "`" + ` are decoded literal values for exact/simple operators,\n  not query syntax. For example, ` + "`" + `app_id:\\*com` + "`" + ` in ` + "`" + `filterQuery` + "`" + ` becomes\n  ` + "`" + `values: [\"*com\"]` + "`" + ` with ` + "`" + `operator: \"equals\"` + "`" + `.\n- The same key can appear in more than one filter object when positive and\n  negative conditions both exist; those filter objects are still combined with\n  AND.\n- Supported filter operators are ` + "`" + `equals` + "`" + `, ` + "`" + `notEquals` + "`" + `, ` + "`" + `in` + "`" + `, ` + "`" + `notIn` + "`" + `,\n  ` + "`" + `contains` + "`" + `, ` + "`" + `notContains` + "`" + `, ` + "`" + `startsWith` + "`" + `, ` + "`" + `notStartsWith` + "`" + `, ` + "`" + `endsWith` + "`" + `,\n  ` + "`" + `notEndsWith` + "`" + `, ` + "`" + `matches` + "`" + `, and ` + "`" + `notMatches` + "`" + `.\n- ` + "`" + `matches` + "`" + ` and ` + "`" + `notMatches` + "`" + ` values use Sentry wildcard pattern syntax, not\n  regular expressions. ` + "`" + `*` + "`" + ` matches zero or more characters, escaped ` + "`" + `\\*` + "`" + `\n  matches a literal asterisk, and a pattern without ` + "`" + `*` + "`" + ` is an exact match.\n- ` + "`" + `in` + "`" + ` and ` + "`" + `notIn` + "`" + ` are evaluated as one condition against all values, matching\n  Sentry's status check trigger behavior.\n- A rule applies only when the build metadata matches all filters. If a\n  referenced metadata key is missing, the filter does not match, even for\n  negated operators.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the project the resource belongs to.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "project_id_or_slug"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "search",
      "pagination": null,
      "related": []
    },
    {
      "id": "Get the latest installable build for a project",
      "method": "GET",
      "path": "/api/0/projects/{organization_id_or_slug}/{project_id_or_slug}/preprodartifacts/build-distribution/latest/",
      "tag": "Mobile Builds",
      "summary": "Get the latest installable build for a project.\n\nReturns the latest installable build matching filter criteria.\nWhen buildVersion is provided, also returns the current build and\nwhether an update is available.",
      "description": "Get the latest installable build for a project.\n\nReturns the latest installable build matching filter criteria.\nWhen buildVersion is provided, also returns the current build and\nwhether an update is available.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the project the resource belongs to.",
          "type": "string"
        },
        {
          "name": "appId",
          "in": "query",
          "required": true,
          "description": "App identifier (exact match).",
          "type": "string"
        },
        {
          "name": "platform",
          "in": "query",
          "required": true,
          "description": "Platform: \"apple\" or \"android\".",
          "type": "string"
        },
        {
          "name": "buildVersion",
          "in": "query",
          "required": false,
          "description": "Current build version. When provided, enables check-for-updates mode.",
          "type": "string"
        },
        {
          "name": "buildNumber",
          "in": "query",
          "required": false,
          "description": "Current build number. Either this or mainBinaryIdentifier must be provided when buildVersion is set.",
          "type": "integer"
        },
        {
          "name": "mainBinaryIdentifier",
          "in": "query",
          "required": false,
          "description": "UUID of the main binary (e.g. Mach-O UUID for Apple builds). Either this or buildNumber must be provided when buildVersion is set.",
          "type": "string"
        },
        {
          "name": "buildConfiguration",
          "in": "query",
          "required": false,
          "description": "Filter by build configuration name (exact match).",
          "type": "string"
        },
        {
          "name": "codesigningType",
          "in": "query",
          "required": false,
          "description": "Filter by code signing type.",
          "type": "string"
        },
        {
          "name": "installGroups",
          "in": "query",
          "required": false,
          "description": "Filter by install group name (repeatable for multiple groups).",
          "type": "array"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "400",
        "403"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "project_id_or_slug"
        ],
        "query": [
          "appId",
          "platform",
          "buildVersion",
          "buildNumber",
          "mainBinaryIdentifier",
          "buildConfiguration",
          "codesigningType",
          "installGroups"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": null,
      "related": []
    },
    {
      "id": "Delete a Replay Instance",
      "method": "DELETE",
      "path": "/api/0/projects/{organization_id_or_slug}/{project_id_or_slug}/replays/{replay_id}/",
      "tag": "Replays",
      "summary": "Delete a replay.",
      "description": "Delete a replay.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the project the resource belongs to.",
          "type": "string"
        },
        {
          "name": "replay_id",
          "in": "path",
          "required": true,
          "description": "The ID of the replay you'd like to retrieve.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "destructive",
      "cacheable": false,
      "responseCodes": [
        "204",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "project_id_or_slug",
          "replay_id"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "delete",
      "pagination": null,
      "related": []
    },
    {
      "id": "List Clicked Nodes",
      "method": "GET",
      "path": "/api/0/projects/{organization_id_or_slug}/{project_id_or_slug}/replays/{replay_id}/clicks/",
      "tag": "Replays",
      "summary": "Retrieve a collection of RRWeb DOM node-ids and the timestamp they were clicked.",
      "description": "Retrieve a collection of RRWeb DOM node-ids and the timestamp they were clicked.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the project the resource belongs to.",
          "type": "string"
        },
        {
          "name": "replay_id",
          "in": "path",
          "required": true,
          "description": "The ID of the replay you'd like to retrieve.",
          "type": "string"
        },
        {
          "name": "cursor",
          "in": "query",
          "required": false,
          "description": "A pointer to the last object fetched and its sort order; used to retrieve the next or previous results.",
          "type": "string"
        },
        {
          "name": "environment",
          "in": "query",
          "required": false,
          "description": "The name of environments to filter by.",
          "type": "array"
        },
        {
          "name": "per_page",
          "in": "query",
          "required": false,
          "description": "Limit the number of rows to return in the result. Default and maximum allowed is 100.",
          "type": "integer"
        },
        {
          "name": "query",
          "in": "query",
          "required": false,
          "description": "Filters results by using [query syntax](/product/sentry-basics/search/).\n\nExample: ` + "`" + `query=(transaction:foo AND release:abc) OR (transaction:[bar,baz] AND release:def)` + "`" + `\n",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "400",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "project_id_or_slug",
          "replay_id"
        ],
        "query": [
          "cursor",
          "environment",
          "per_page",
          "query"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": {
        "style": "cursor",
        "cursorParam": "cursor"
      },
      "related": []
    },
    {
      "id": "List Recording Segments",
      "method": "GET",
      "path": "/api/0/projects/{organization_id_or_slug}/{project_id_or_slug}/replays/{replay_id}/recording-segments/",
      "tag": "Replays",
      "summary": "Return a collection of replay recording segments.",
      "description": "Return a collection of replay recording segments.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the project the resource belongs to.",
          "type": "string"
        },
        {
          "name": "replay_id",
          "in": "path",
          "required": true,
          "description": "The ID of the replay you'd like to retrieve.",
          "type": "string"
        },
        {
          "name": "cursor",
          "in": "query",
          "required": false,
          "description": "A pointer to the last object fetched and its sort order; used to retrieve the next or previous results.",
          "type": "string"
        },
        {
          "name": "per_page",
          "in": "query",
          "required": false,
          "description": "Limit the number of rows to return in the result. Default and maximum allowed is 100.",
          "type": "integer"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "400",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "project_id_or_slug",
          "replay_id"
        ],
        "query": [
          "cursor",
          "per_page"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": {
        "style": "cursor",
        "cursorParam": "cursor"
      },
      "related": []
    },
    {
      "id": "Retrieve a Recording Segment",
      "method": "GET",
      "path": "/api/0/projects/{organization_id_or_slug}/{project_id_or_slug}/replays/{replay_id}/recording-segments/{segment_id}/",
      "tag": "Replays",
      "summary": "Return a replay recording segment.",
      "description": "Return a replay recording segment.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the project the resource belongs to.",
          "type": "string"
        },
        {
          "name": "replay_id",
          "in": "path",
          "required": true,
          "description": "The ID of the replay you'd like to retrieve.",
          "type": "string"
        },
        {
          "name": "segment_id",
          "in": "path",
          "required": true,
          "description": "The ID of the segment you'd like to retrieve.",
          "type": "integer"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "400",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "project_id_or_slug",
          "replay_id",
          "segment_id"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": null,
      "related": []
    },
    {
      "id": "List Users Who Have Viewed a Replay",
      "method": "GET",
      "path": "/api/0/projects/{organization_id_or_slug}/{project_id_or_slug}/replays/{replay_id}/viewed-by/",
      "tag": "Replays",
      "summary": "Return a list of users who have viewed a replay.",
      "description": "Return a list of users who have viewed a replay.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the project the resource belongs to.",
          "type": "string"
        },
        {
          "name": "replay_id",
          "in": "path",
          "required": true,
          "description": "The ID of the replay you'd like to retrieve.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "400",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "project_id_or_slug",
          "replay_id"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": null,
      "related": []
    },
    {
      "id": "List Replay Batch-Deletion Jobs",
      "method": "GET",
      "path": "/api/0/projects/{organization_id_or_slug}/{project_id_or_slug}/replays/jobs/delete/",
      "tag": "Replays",
      "summary": "Retrieve a collection of replay delete jobs.",
      "description": "Retrieve a collection of replay delete jobs.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the project the resource belongs to.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "403"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "project_id_or_slug"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": null,
      "related": [
        {
          "id": "Create Replay Batch Deletion Job",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Create Replay Batch Deletion Job",
      "method": "POST",
      "path": "/api/0/projects/{organization_id_or_slug}/{project_id_or_slug}/replays/jobs/delete/",
      "tag": "Replays",
      "summary": "Create a new replay deletion job.",
      "description": "Create a new replay deletion job.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the project the resource belongs to.",
          "type": "string"
        }
      ],
      "hasRequestBody": true,
      "risk": "write",
      "cacheable": false,
      "responseCodes": [
        "201",
        "400",
        "403"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "project_id_or_slug"
        ],
        "query": [],
        "headers": [],
        "body": true
      },
      "graphql": null,
      "kind": "create",
      "pagination": null,
      "related": [
        {
          "id": "List Replay Batch-Deletion Jobs",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Retrieve a Replay Batch-Deletion Job",
      "method": "GET",
      "path": "/api/0/projects/{organization_id_or_slug}/{project_id_or_slug}/replays/jobs/delete/{job_id}/",
      "tag": "Replays",
      "summary": "Fetch a replay delete job instance.",
      "description": "Fetch a replay delete job instance.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the project the resource belongs to.",
          "type": "string"
        },
        {
          "name": "job_id",
          "in": "path",
          "required": true,
          "description": "The ID of the replay deletion job you'd like to retrieve.",
          "type": "integer"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "project_id_or_slug",
          "job_id"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": null,
      "related": []
    },
    {
      "id": "Retrieve a Project's Symbol Sources",
      "method": "GET",
      "path": "/api/0/projects/{organization_id_or_slug}/{project_id_or_slug}/symbol-sources/",
      "tag": "Projects",
      "summary": "List custom symbol sources configured for a project.",
      "description": "List custom symbol sources configured for a project.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the project the resource belongs to.",
          "type": "string"
        },
        {
          "name": "id",
          "in": "query",
          "required": false,
          "description": "The ID of the source to look up. If this is not provided, all sources are returned.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "project_id_or_slug"
        ],
        "query": [
          "id"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": null,
      "related": [
        {
          "id": "Add a Symbol Source to a Project",
          "role": "same-resource"
        },
        {
          "id": "Update a Project's Symbol Source",
          "role": "same-resource"
        },
        {
          "id": "Delete a Symbol Source from a Project",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Add a Symbol Source to a Project",
      "method": "POST",
      "path": "/api/0/projects/{organization_id_or_slug}/{project_id_or_slug}/symbol-sources/",
      "tag": "Projects",
      "summary": "Add a custom symbol source to a project.",
      "description": "Add a custom symbol source to a project.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the project the resource belongs to.",
          "type": "string"
        }
      ],
      "hasRequestBody": true,
      "risk": "write",
      "cacheable": false,
      "responseCodes": [
        "201",
        "400",
        "403"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "project_id_or_slug"
        ],
        "query": [],
        "headers": [],
        "body": true
      },
      "graphql": null,
      "kind": "create",
      "pagination": null,
      "related": [
        {
          "id": "Retrieve a Project's Symbol Sources",
          "role": "same-resource"
        },
        {
          "id": "Update a Project's Symbol Source",
          "role": "same-resource"
        },
        {
          "id": "Delete a Symbol Source from a Project",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Update a Project's Symbol Source",
      "method": "PUT",
      "path": "/api/0/projects/{organization_id_or_slug}/{project_id_or_slug}/symbol-sources/",
      "tag": "Projects",
      "summary": "Update a custom symbol source in a project.",
      "description": "Update a custom symbol source in a project.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the project the resource belongs to.",
          "type": "string"
        },
        {
          "name": "id",
          "in": "query",
          "required": true,
          "description": "The ID of the source to update.",
          "type": "string"
        }
      ],
      "hasRequestBody": true,
      "risk": "write",
      "cacheable": false,
      "responseCodes": [
        "200",
        "400",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "project_id_or_slug"
        ],
        "query": [
          "id"
        ],
        "headers": [],
        "body": true
      },
      "graphql": null,
      "kind": "update",
      "pagination": null,
      "related": [
        {
          "id": "Retrieve a Project's Symbol Sources",
          "role": "same-resource"
        },
        {
          "id": "Add a Symbol Source to a Project",
          "role": "same-resource"
        },
        {
          "id": "Delete a Symbol Source from a Project",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Delete a Symbol Source from a Project",
      "method": "DELETE",
      "path": "/api/0/projects/{organization_id_or_slug}/{project_id_or_slug}/symbol-sources/",
      "tag": "Projects",
      "summary": "Delete a custom symbol source from a project.",
      "description": "Delete a custom symbol source from a project.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the project the resource belongs to.",
          "type": "string"
        },
        {
          "name": "id",
          "in": "query",
          "required": true,
          "description": "The ID of the source to delete.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "destructive",
      "cacheable": false,
      "responseCodes": [
        "204",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "project_id_or_slug"
        ],
        "query": [
          "id"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "delete",
      "pagination": null,
      "related": [
        {
          "id": "Retrieve a Project's Symbol Sources",
          "role": "same-resource"
        },
        {
          "id": "Add a Symbol Source to a Project",
          "role": "same-resource"
        },
        {
          "id": "Update a Project's Symbol Source",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "List a Project's Teams",
      "method": "GET",
      "path": "/api/0/projects/{organization_id_or_slug}/{project_id_or_slug}/teams/",
      "tag": "Teams",
      "summary": "Return a list of teams that have access to this project.",
      "description": "Return a list of teams that have access to this project.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the project the resource belongs to.",
          "type": "string"
        },
        {
          "name": "cursor",
          "in": "query",
          "required": false,
          "description": "A pointer to the last object fetched and its sort order; used to retrieve the next or previous results.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "project_id_or_slug"
        ],
        "query": [
          "cursor"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": {
        "style": "cursor",
        "cursorParam": "cursor"
      },
      "related": []
    },
    {
      "id": "Add a Team to a Project",
      "method": "POST",
      "path": "/api/0/projects/{organization_id_or_slug}/{project_id_or_slug}/teams/{team_id_or_slug}/",
      "tag": "Projects",
      "summary": "Give a team access to a project.",
      "description": "Give a team access to a project.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the project the resource belongs to.",
          "type": "string"
        },
        {
          "name": "team_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the team the resource belongs to.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "write",
      "cacheable": false,
      "responseCodes": [
        "201",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "project_id_or_slug",
          "team_id_or_slug"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "create",
      "pagination": null,
      "related": [
        {
          "id": "Delete a Team from a Project",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Delete a Team from a Project",
      "method": "DELETE",
      "path": "/api/0/projects/{organization_id_or_slug}/{project_id_or_slug}/teams/{team_id_or_slug}/",
      "tag": "Projects",
      "summary": "Revoke a team's access to a project.\n\nNote that Team Admins can only revoke access to teams they are admins of.",
      "description": "Revoke a team's access to a project.\n\nNote that Team Admins can only revoke access to teams they are admins of.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the project the resource belongs to.",
          "type": "string"
        },
        {
          "name": "team_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the team the resource belongs to.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "destructive",
      "cacheable": false,
      "responseCodes": [
        "200",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "project_id_or_slug",
          "team_id_or_slug"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "delete",
      "pagination": null,
      "related": [
        {
          "id": "Add a Team to a Project",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "List Seer AI Models",
      "method": "GET",
      "path": "/api/0/seer/models/",
      "tag": "Seer",
      "summary": "Get list of actively used LLM model names from Seer.\n\nReturns the list of AI models that are currently used in production in Seer.\nThis endpoint does not require authentication and can be used to discover which models Seer uses.\n\nRequests to this endpoint should use the region-specific domain\neg. ` + "`" + `us.sentry.io` + "`" + ` or ` + "`" + `de.sentry.io` + "`" + `",
      "description": "Get list of actively used LLM model names from Seer.\n\nReturns the list of AI models that are currently used in production in Seer.\nThis endpoint does not require authentication and can be used to discover which models Seer uses.\n\nRequests to this endpoint should use the region-specific domain\neg. ` + "`" + `us.sentry.io` + "`" + ` or ` + "`" + `de.sentry.io` + "`" + `",
      "parameters": [],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200"
      ],
      "inputHints": {
        "path": [],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read-list",
      "pagination": null,
      "related": []
    },
    {
      "id": "Retrieve a custom integration by ID or slug.",
      "method": "GET",
      "path": "/api/0/sentry-apps/{sentry_app_id_or_slug}/",
      "tag": "Integration",
      "summary": "Retrieve a custom integration.",
      "description": "Retrieve a custom integration.",
      "parameters": [
        {
          "name": "sentry_app_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the custom integration.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200"
      ],
      "inputHints": {
        "path": [
          "sentry_app_id_or_slug"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": null,
      "related": [
        {
          "id": "Update an existing custom integration.",
          "role": "same-resource"
        },
        {
          "id": "Delete a custom integration.",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Update an existing custom integration.",
      "method": "PUT",
      "path": "/api/0/sentry-apps/{sentry_app_id_or_slug}/",
      "tag": "Integration",
      "summary": "Update an existing custom integration.",
      "description": "Update an existing custom integration.",
      "parameters": [
        {
          "name": "sentry_app_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the custom integration.",
          "type": "string"
        }
      ],
      "hasRequestBody": true,
      "risk": "write",
      "cacheable": false,
      "responseCodes": [
        "200",
        "400",
        "403"
      ],
      "inputHints": {
        "path": [
          "sentry_app_id_or_slug"
        ],
        "query": [],
        "headers": [],
        "body": true
      },
      "graphql": null,
      "kind": "update",
      "pagination": null,
      "related": [
        {
          "id": "Retrieve a custom integration by ID or slug.",
          "role": "same-resource"
        },
        {
          "id": "Delete a custom integration.",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Delete a custom integration.",
      "method": "DELETE",
      "path": "/api/0/sentry-apps/{sentry_app_id_or_slug}/",
      "tag": "Integration",
      "summary": "Delete a custom integration.",
      "description": "Delete a custom integration.",
      "parameters": [
        {
          "name": "sentry_app_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the custom integration.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "destructive",
      "cacheable": false,
      "responseCodes": [
        "204",
        "403"
      ],
      "inputHints": {
        "path": [
          "sentry_app_id_or_slug"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "delete",
      "pagination": null,
      "related": [
        {
          "id": "Retrieve a custom integration by ID or slug.",
          "role": "same-resource"
        },
        {
          "id": "Update an existing custom integration.",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Retrieve a Team",
      "method": "GET",
      "path": "/api/0/teams/{organization_id_or_slug}/{team_id_or_slug}/",
      "tag": "Teams",
      "summary": "Return details on an individual team.",
      "description": "Return details on an individual team.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "team_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the team the resource belongs to.",
          "type": "string"
        },
        {
          "name": "expand",
          "in": "query",
          "required": false,
          "description": "\nList of strings to opt in to additional data. Supports ` + "`" + `projects` + "`" + `, ` + "`" + `externalTeams` + "`" + `.\n",
          "type": "string"
        },
        {
          "name": "collapse",
          "in": "query",
          "required": false,
          "description": "\nList of strings to opt out of certain pieces of data. Supports ` + "`" + `organization` + "`" + `.\n",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "401",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "team_id_or_slug"
        ],
        "query": [
          "expand",
          "collapse"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": null,
      "related": [
        {
          "id": "Update a Team",
          "role": "same-resource"
        },
        {
          "id": "Delete a Team",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Update a Team",
      "method": "PUT",
      "path": "/api/0/teams/{organization_id_or_slug}/{team_id_or_slug}/",
      "tag": "Teams",
      "summary": "Update various attributes and configurable settings for the given\nteam.",
      "description": "Update various attributes and configurable settings for the given\nteam.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "team_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the team the resource belongs to.",
          "type": "string"
        }
      ],
      "hasRequestBody": true,
      "risk": "write",
      "cacheable": false,
      "responseCodes": [
        "200",
        "401",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "team_id_or_slug"
        ],
        "query": [],
        "headers": [],
        "body": true
      },
      "graphql": null,
      "kind": "update",
      "pagination": null,
      "related": [
        {
          "id": "Retrieve a Team",
          "role": "same-resource"
        },
        {
          "id": "Delete a Team",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Delete a Team",
      "method": "DELETE",
      "path": "/api/0/teams/{organization_id_or_slug}/{team_id_or_slug}/",
      "tag": "Teams",
      "summary": "Schedules a team for deletion.\n\n**Note:** Deletion happens asynchronously and therefore is not\nimmediate. Teams will have their slug released while waiting for deletion.",
      "description": "Schedules a team for deletion.\n\n**Note:** Deletion happens asynchronously and therefore is not\nimmediate. Teams will have their slug released while waiting for deletion.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "team_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the team the resource belongs to.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "destructive",
      "cacheable": false,
      "responseCodes": [
        "204",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "team_id_or_slug"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "delete",
      "pagination": null,
      "related": [
        {
          "id": "Retrieve a Team",
          "role": "same-resource"
        },
        {
          "id": "Update a Team",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Create an External Team",
      "method": "POST",
      "path": "/api/0/teams/{organization_id_or_slug}/{team_id_or_slug}/external-teams/",
      "tag": "Integrations",
      "summary": "Link a team from an external provider to a Sentry team.",
      "description": "Link a team from an external provider to a Sentry team.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "team_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the team the resource belongs to.",
          "type": "string"
        }
      ],
      "hasRequestBody": true,
      "risk": "write",
      "cacheable": false,
      "responseCodes": [
        "200",
        "201",
        "400",
        "403"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "team_id_or_slug"
        ],
        "query": [],
        "headers": [],
        "body": true
      },
      "graphql": null,
      "kind": "create",
      "pagination": null,
      "related": []
    },
    {
      "id": "Update an External Team",
      "method": "PUT",
      "path": "/api/0/teams/{organization_id_or_slug}/{team_id_or_slug}/external-teams/{external_team_id}/",
      "tag": "Integrations",
      "summary": "Update a team in an external provider that is currently linked to a Sentry team.",
      "description": "Update a team in an external provider that is currently linked to a Sentry team.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "team_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the team the resource belongs to.",
          "type": "string"
        },
        {
          "name": "external_team_id",
          "in": "path",
          "required": true,
          "description": "The ID of the external team object. This is returned when creating an external team.",
          "type": "integer"
        }
      ],
      "hasRequestBody": true,
      "risk": "write",
      "cacheable": false,
      "responseCodes": [
        "200",
        "400",
        "403"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "team_id_or_slug",
          "external_team_id"
        ],
        "query": [],
        "headers": [],
        "body": true
      },
      "graphql": null,
      "kind": "update",
      "pagination": null,
      "related": [
        {
          "id": "Delete an External Team",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Delete an External Team",
      "method": "DELETE",
      "path": "/api/0/teams/{organization_id_or_slug}/{team_id_or_slug}/external-teams/{external_team_id}/",
      "tag": "Integrations",
      "summary": "Delete the link between a team from an external provider and a Sentry team.",
      "description": "Delete the link between a team from an external provider and a Sentry team.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "team_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the team the resource belongs to.",
          "type": "string"
        },
        {
          "name": "external_team_id",
          "in": "path",
          "required": true,
          "description": "The ID of the external team object. This is returned when creating an external team.",
          "type": "integer"
        }
      ],
      "hasRequestBody": false,
      "risk": "destructive",
      "cacheable": false,
      "responseCodes": [
        "204",
        "400",
        "403"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "team_id_or_slug",
          "external_team_id"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "delete",
      "pagination": null,
      "related": [
        {
          "id": "Update an External Team",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "List a Team's Members",
      "method": "GET",
      "path": "/api/0/teams/{organization_id_or_slug}/{team_id_or_slug}/members/",
      "tag": "Teams",
      "summary": "List all members on a team.\n\nThe response will not include members with pending invites.",
      "description": "List all members on a team.\n\nThe response will not include members with pending invites.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "team_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the team the resource belongs to.",
          "type": "string"
        },
        {
          "name": "cursor",
          "in": "query",
          "required": false,
          "description": "A pointer to the last object fetched and its sort order; used to retrieve the next or previous results.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "team_id_or_slug"
        ],
        "query": [
          "cursor"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": {
        "style": "cursor",
        "cursorParam": "cursor"
      },
      "related": []
    },
    {
      "id": "List a Team's Projects",
      "method": "GET",
      "path": "/api/0/teams/{organization_id_or_slug}/{team_id_or_slug}/projects/",
      "tag": "Teams",
      "summary": "Return a list of projects bound to a team.",
      "description": "Return a list of projects bound to a team.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "team_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the team the resource belongs to.",
          "type": "string"
        },
        {
          "name": "cursor",
          "in": "query",
          "required": false,
          "description": "A pointer to the last object fetched and its sort order; used to retrieve the next or previous results.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "team_id_or_slug"
        ],
        "query": [
          "cursor"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": {
        "style": "cursor",
        "cursorParam": "cursor"
      },
      "related": [
        {
          "id": "Create a New Project",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Create a New Project",
      "method": "POST",
      "path": "/api/0/teams/{organization_id_or_slug}/{team_id_or_slug}/projects/",
      "tag": "Projects",
      "summary": "Create a new project bound to a team.\n\n        Note: If your organization has disabled member project creation, the ` + "`" + `org:write` + "`" + ` or ` + "`" + `team:admin` + "`" + ` scope is required.\n        ",
      "description": "Create a new project bound to a team.\n\n        Note: If your organization has disabled member project creation, the ` + "`" + `org:write` + "`" + ` or ` + "`" + `team:admin` + "`" + ` scope is required.\n        ",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "team_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the team the resource belongs to.",
          "type": "string"
        }
      ],
      "hasRequestBody": true,
      "risk": "write",
      "cacheable": false,
      "responseCodes": [
        "201",
        "400",
        "403",
        "404",
        "409"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "team_id_or_slug"
        ],
        "query": [],
        "headers": [],
        "body": true
      },
      "graphql": null,
      "kind": "create",
      "pagination": null,
      "related": [
        {
          "id": "List a Team's Projects",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "List an Organization's Repositories",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/repos/",
      "tag": "Organizations",
      "summary": "Return a list of version control repositories for a given organization.",
      "description": "Return a list of version control repositories for a given organization.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The organization short name.",
          "type": "string"
        },
        {
          "name": "cursor",
          "in": "query",
          "required": false,
          "description": "A pointer to the last object fetched and its sort order; used to retrieve the next or previous results.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug"
        ],
        "query": [
          "cursor"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": {
        "style": "cursor",
        "cursorParam": "cursor"
      },
      "related": []
    },
    {
      "id": "List a Project's Debug Information Files",
      "method": "GET",
      "path": "/api/0/projects/{organization_id_or_slug}/{project_id_or_slug}/files/dsyms/",
      "tag": "Projects",
      "summary": "Retrieve a list of debug information files for a given project.",
      "description": "Retrieve a list of debug information files for a given project.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the file belongs to.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the project to list the DIFs of.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "project_id_or_slug"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": null,
      "related": [
        {
          "id": "Upload a New File",
          "role": "same-resource"
        },
        {
          "id": "Delete a Specific Project's Debug Information File",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Upload a New File",
      "method": "POST",
      "path": "/api/0/projects/{organization_id_or_slug}/{project_id_or_slug}/files/dsyms/",
      "tag": "Projects",
      "summary": "Upload a new debug information file for the given release.\n\nUnlike other API requests, files must be uploaded using the\ntraditional multipart/form-data content-type.\n\nRequests to this endpoint should use the region-specific domain eg. ` + "`" + `us.sentry.io` + "`" + ` or ` + "`" + `de.sentry.io` + "`" + `.\n\nThe file uploaded is a zip archive of an Apple .dSYM folder which\ncontains the individual debug images.  Uploading through this endpoint\nwill create different files for the contained images.",
      "description": "Upload a new debug information file for the given release.\n\nUnlike other API requests, files must be uploaded using the\ntraditional multipart/form-data content-type.\n\nRequests to this endpoint should use the region-specific domain eg. ` + "`" + `us.sentry.io` + "`" + ` or ` + "`" + `de.sentry.io` + "`" + `.\n\nThe file uploaded is a zip archive of an Apple .dSYM folder which\ncontains the individual debug images.  Uploading through this endpoint\nwill create different files for the contained images.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the project belongs to.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the project to upload a file to.",
          "type": "string"
        }
      ],
      "hasRequestBody": true,
      "risk": "write",
      "cacheable": false,
      "responseCodes": [
        "201",
        "400",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "project_id_or_slug"
        ],
        "query": [],
        "headers": [],
        "body": true
      },
      "graphql": null,
      "kind": "create",
      "pagination": null,
      "related": [
        {
          "id": "List a Project's Debug Information Files",
          "role": "same-resource"
        },
        {
          "id": "Delete a Specific Project's Debug Information File",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Delete a Specific Project's Debug Information File",
      "method": "DELETE",
      "path": "/api/0/projects/{organization_id_or_slug}/{project_id_or_slug}/files/dsyms/",
      "tag": "Projects",
      "summary": "Delete a debug information file for a given project.",
      "description": "Delete a debug information file for a given project.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the file belongs to.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the project to delete the DIF.",
          "type": "string"
        },
        {
          "name": "id",
          "in": "query",
          "required": true,
          "description": "The ID of the DIF to delete.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "destructive",
      "cacheable": false,
      "responseCodes": [
        "204",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "project_id_or_slug"
        ],
        "query": [
          "id"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "delete",
      "pagination": null,
      "related": [
        {
          "id": "List a Project's Debug Information Files",
          "role": "same-resource"
        },
        {
          "id": "Upload a New File",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "List a Project's Users",
      "method": "GET",
      "path": "/api/0/projects/{organization_id_or_slug}/{project_id_or_slug}/users/",
      "tag": "Projects",
      "summary": "Return a list of users seen within this project.",
      "description": "Return a list of users seen within this project.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the project.",
          "type": "string"
        },
        {
          "name": "query",
          "in": "query",
          "required": false,
          "description": "Limit results to users matching the given query. Prefixes should be used to suggest the field to match on: ` + "`" + `id` + "`" + `, ` + "`" + `email` + "`" + `, ` + "`" + `username` + "`" + `, ` + "`" + `ip` + "`" + `. For example, ` + "`" + `query=email:foo@example.com` + "`" + `",
          "type": "string"
        },
        {
          "name": "cursor",
          "in": "query",
          "required": false,
          "description": "A pointer to the last object fetched and its sort order; used to retrieve the next or previous results.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "403"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "project_id_or_slug"
        ],
        "query": [
          "query",
          "cursor"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": {
        "style": "cursor",
        "cursorParam": "cursor"
      },
      "related": []
    },
    {
      "id": "List a Tag's Values",
      "method": "GET",
      "path": "/api/0/projects/{organization_id_or_slug}/{project_id_or_slug}/tags/{key}/values/",
      "tag": "Projects",
      "summary": "Return a list of values associated with this key.  The ` + "`" + `query` + "`" + `\nparameter can be used to to perform a \"contains\" match on\nvalues. \n\nWhen [paginated](/api/pagination) can return at most 1000 values.",
      "description": "Return a list of values associated with this key.  The ` + "`" + `query` + "`" + `\nparameter can be used to to perform a \"contains\" match on\nvalues. \n\nWhen [paginated](/api/pagination) can return at most 1000 values.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the project.",
          "type": "string"
        },
        {
          "name": "key",
          "in": "path",
          "required": true,
          "description": "The tag key to look up.",
          "type": "string"
        },
        {
          "name": "cursor",
          "in": "query",
          "required": false,
          "description": "A pointer to the last object fetched and its sort order; used to retrieve the next or previous results.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "403"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "project_id_or_slug",
          "key"
        ],
        "query": [
          "cursor"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "search",
      "pagination": {
        "style": "cursor",
        "cursorParam": "cursor"
      },
      "related": []
    },
    {
      "id": "Retrieve Event Counts for a Project",
      "method": "GET",
      "path": "/api/0/projects/{organization_id_or_slug}/{project_id_or_slug}/stats/",
      "tag": "Projects",
      "summary": "Caution\nThis endpoint may change in the future without  notice.",
      "description": "Return a set of points representing a normalized timestamp and the\nnumber of events seen in the period.\n\nQuery ranges are limited to Sentry's configured time-series resolutions.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the project.",
          "type": "string"
        },
        {
          "name": "stat",
          "in": "query",
          "required": false,
          "description": "The name of the stat to query ` + "`" + `(\"received\", \"rejected\", \"blacklisted\", \"generated\")` + "`" + `.",
          "type": "string"
        },
        {
          "name": "since",
          "in": "query",
          "required": false,
          "description": "A timestamp to set the start of the query in seconds since UNIX epoch.",
          "type": "string"
        },
        {
          "name": "until",
          "in": "query",
          "required": false,
          "description": "A timestamp to set the end of the query in seconds since UNIX epoch.",
          "type": "string"
        },
        {
          "name": "resolution",
          "in": "query",
          "required": false,
          "description": "An explicit resolution to search for (one of ` + "`" + `10s` + "`" + `, ` + "`" + `1h` + "`" + `, and ` + "`" + `1d` + "`" + `).",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "403"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "project_id_or_slug"
        ],
        "query": [
          "stat",
          "since",
          "until",
          "resolution"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": null,
      "related": []
    },
    {
      "id": "List a Project's User Feedback",
      "method": "GET",
      "path": "/api/0/projects/{organization_id_or_slug}/{project_id_or_slug}/user-feedback/",
      "tag": "Projects",
      "summary": "Return a list of user feedback items within this project.\n\n*This list does not include submissions from the [User Feedback Widget](https://docs.sentry.io/product/user-feedback/#user-feedback-widget). This is because it is based on an older format called User Reports - read more [here](https://develop.sentry.dev/application/feedback-architecture/#user-reports). To return a list of user feedback items from the widget, please use the [issue API](https://docs.sentry.io/api/events/list-a-projects-issues/) with the filter ` + "`" + `issue.category:feedback` + "`" + `.*",
      "description": "Return a list of user feedback items within this project.\n\n*This list does not include submissions from the [User Feedback Widget](https://docs.sentry.io/product/user-feedback/#user-feedback-widget). This is because it is based on an older format called User Reports - read more [here](https://develop.sentry.dev/application/feedback-architecture/#user-reports). To return a list of user feedback items from the widget, please use the [issue API](https://docs.sentry.io/api/events/list-a-projects-issues/) with the filter ` + "`" + `issue.category:feedback` + "`" + `.*",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the project.",
          "type": "string"
        },
        {
          "name": "cursor",
          "in": "query",
          "required": false,
          "description": "A pointer to the last object fetched and its sort order; used to retrieve the next or previous results.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "project_id_or_slug"
        ],
        "query": [
          "cursor"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": {
        "style": "cursor",
        "cursorParam": "cursor"
      },
      "related": [
        {
          "id": "Submit User Feedback",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Submit User Feedback",
      "method": "POST",
      "path": "/api/0/projects/{organization_id_or_slug}/{project_id_or_slug}/user-feedback/",
      "tag": "Projects",
      "summary": "*This endpoint is DEPRECATED. We document it here for older SDKs and users who are still migrating to the [User Feedback Widget](https://docs.sentry.io/product/user-feedback/#user-feedback-widget) or [API](https://docs.sentry.io/platforms/javascript/user-feedback/#user-feedback-api)(multi-platform). If you are a new user, do not use this endpoint - unless you don't have a JS frontend, and your platform's SDK does not offer a feedback API.*\n\nFeedback must be received by the server no more than 30 minutes after the event was saved.\n\nAdditionally, within 5 minutes of submitting feedback it may also be overwritten. This is useful in situations where you may need to retry sending a request due to network failures.\n\nIf feedback is rejected due to a mutability threshold, a 409 status code will be returned.\n\nNote: Feedback may be submitted with DSN authentication (see auth documentation).",
      "description": "*This endpoint is DEPRECATED. We document it here for older SDKs and users who are still migrating to the [User Feedback Widget](https://docs.sentry.io/product/user-feedback/#user-feedback-widget) or [API](https://docs.sentry.io/platforms/javascript/user-feedback/#user-feedback-api)(multi-platform). If you are a new user, do not use this endpoint - unless you don't have a JS frontend, and your platform's SDK does not offer a feedback API.*\n\nFeedback must be received by the server no more than 30 minutes after the event was saved.\n\nAdditionally, within 5 minutes of submitting feedback it may also be overwritten. This is useful in situations where you may need to retry sending a request due to network failures.\n\nIf feedback is rejected due to a mutability threshold, a 409 status code will be returned.\n\nNote: Feedback may be submitted with DSN authentication (see auth documentation).",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the project.",
          "type": "string"
        }
      ],
      "hasRequestBody": true,
      "risk": "write",
      "cacheable": false,
      "responseCodes": [
        "200",
        "400",
        "403",
        "404",
        "409"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "project_id_or_slug"
        ],
        "query": [],
        "headers": [],
        "body": true
      },
      "graphql": null,
      "kind": "create",
      "pagination": null,
      "related": [
        {
          "id": "List a Project's User Feedback",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "List a Project's Service Hooks",
      "method": "GET",
      "path": "/api/0/projects/{organization_id_or_slug}/{project_id_or_slug}/hooks/",
      "tag": "Projects",
      "summary": "Return a list of service hooks bound to a project.",
      "description": "Return a list of service hooks bound to a project.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the client keys belong to.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the project the client keys belong to.",
          "type": "string"
        },
        {
          "name": "cursor",
          "in": "query",
          "required": false,
          "description": "A pointer to the last object fetched and its sort order; used to retrieve the next or previous results.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "403"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "project_id_or_slug"
        ],
        "query": [
          "cursor"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": {
        "style": "cursor",
        "cursorParam": "cursor"
      },
      "related": [
        {
          "id": "Register a New Service Hook",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Register a New Service Hook",
      "method": "POST",
      "path": "/api/0/projects/{organization_id_or_slug}/{project_id_or_slug}/hooks/",
      "tag": "Projects",
      "summary": "Register a new service hook on a project.\n\nEvents include:\n\n- event.alert: An alert is generated for an event (via rules).\n- event.created: A new event has been processed.\n\nThis endpoint requires the 'servicehooks' feature to be enabled for your project.",
      "description": "Register a new service hook on a project.\n\nEvents include:\n\n- event.alert: An alert is generated for an event (via rules).\n- event.created: A new event has been processed.\n\nThis endpoint requires the 'servicehooks' feature to be enabled for your project.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the client keys belong to.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the project the client keys belong to.",
          "type": "string"
        }
      ],
      "hasRequestBody": true,
      "risk": "write",
      "cacheable": false,
      "responseCodes": [
        "201",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "project_id_or_slug"
        ],
        "query": [],
        "headers": [],
        "body": true
      },
      "graphql": null,
      "kind": "create",
      "pagination": null,
      "related": [
        {
          "id": "List a Project's Service Hooks",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Retrieve a Service Hook",
      "method": "GET",
      "path": "/api/0/projects/{organization_id_or_slug}/{project_id_or_slug}/hooks/{hook_id}/",
      "tag": "Projects",
      "summary": "Return a service hook bound to a project.",
      "description": "Return a service hook bound to a project.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the client keys belong to.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the project the client keys belong to.",
          "type": "string"
        },
        {
          "name": "hook_id",
          "in": "path",
          "required": true,
          "description": "The GUID of the service hook.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "project_id_or_slug",
          "hook_id"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": null,
      "related": [
        {
          "id": "Update a Service Hook",
          "role": "same-resource"
        },
        {
          "id": "Remove a Service Hook",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Update a Service Hook",
      "method": "PUT",
      "path": "/api/0/projects/{organization_id_or_slug}/{project_id_or_slug}/hooks/{hook_id}/",
      "tag": "Projects",
      "summary": "Update a service hook.",
      "description": "Update a service hook.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the client keys belong to.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the project the client keys belong to.",
          "type": "string"
        },
        {
          "name": "hook_id",
          "in": "path",
          "required": true,
          "description": "The GUID of the service hook.",
          "type": "string"
        }
      ],
      "hasRequestBody": true,
      "risk": "write",
      "cacheable": false,
      "responseCodes": [
        "200",
        "400",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "project_id_or_slug",
          "hook_id"
        ],
        "query": [],
        "headers": [],
        "body": true
      },
      "graphql": null,
      "kind": "update",
      "pagination": null,
      "related": [
        {
          "id": "Retrieve a Service Hook",
          "role": "same-resource"
        },
        {
          "id": "Remove a Service Hook",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Remove a Service Hook",
      "method": "DELETE",
      "path": "/api/0/projects/{organization_id_or_slug}/{project_id_or_slug}/hooks/{hook_id}/",
      "tag": "Projects",
      "summary": "Remove a service hook.",
      "description": "Remove a service hook.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the client keys belong to.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the project the client keys belong to.",
          "type": "string"
        },
        {
          "name": "hook_id",
          "in": "path",
          "required": true,
          "description": "The GUID of the service hook.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "destructive",
      "cacheable": false,
      "responseCodes": [
        "204",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "project_id_or_slug",
          "hook_id"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "delete",
      "pagination": null,
      "related": [
        {
          "id": "Retrieve a Service Hook",
          "role": "same-resource"
        },
        {
          "id": "Update a Service Hook",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Retrieve an Event for a Project",
      "method": "GET",
      "path": "/api/0/projects/{organization_id_or_slug}/{project_id_or_slug}/events/{event_id}/",
      "tag": "Events",
      "summary": "Return details on an individual event.",
      "description": "Return details on an individual event.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the event belongs to.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the project the event belongs to.",
          "type": "string"
        },
        {
          "name": "event_id",
          "in": "path",
          "required": true,
          "description": "The ID of the event to retrieve. It is the hexadecimal ID as reported by the client.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "403"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "project_id_or_slug",
          "event_id"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": null,
      "related": []
    },
    {
      "id": "List a Project's Issues",
      "method": "GET",
      "path": "/api/0/projects/{organization_id_or_slug}/{project_id_or_slug}/issues/",
      "tag": "Events",
      "summary": "**Deprecated**: This endpoint has been replaced with the [Organization Issues endpoint](/api/events/list-an-organizations-issues/) which\nsupports filtering on project and additional functionality.\n\nReturn a list of issues (groups) bound to a project.  All parameters are supplied as query string parameters. \n\n A default query of ` + "`" + `` + "`" + `is:unresolved` + "`" + `` + "`" + ` is applied. To return results with other statuses send an new query value (i.e. ` + "`" + `` + "`" + `?query=` + "`" + `` + "`" + ` for all results).\n\nThe ` + "`" + `` + "`" + `statsPeriod` + "`" + `` + "`" + ` parameter can be used to select the timeline stats which should be present. Possible values are: ` + "`" + `` + "`" + `\"\"` + "`" + `` + "`" + ` (disable),` + "`" + `` + "`" + `\"24h\"` + "`" + `` + "`" + ` (default), ` + "`" + `` + "`" + `\"14d\"` + "`" + `` + "`" + `\n\nUser feedback items from the [User Feedback Widget](https://docs.sentry.io/product/user-feedback/#user-feedback-widget) are built off the issue platform, so to return a list of user feedback items for a specific project, filter for ` + "`" + `issue.category:feedback` + "`" + `.",
      "description": "**Deprecated**: This endpoint has been replaced with the [Organization Issues endpoint](/api/events/list-an-organizations-issues/) which\nsupports filtering on project and additional functionality.\n\nReturn a list of issues (groups) bound to a project.  All parameters are supplied as query string parameters. \n\n A default query of ` + "`" + `` + "`" + `is:unresolved` + "`" + `` + "`" + ` is applied. To return results with other statuses send an new query value (i.e. ` + "`" + `` + "`" + `?query=` + "`" + `` + "`" + ` for all results).\n\nThe ` + "`" + `` + "`" + `statsPeriod` + "`" + `` + "`" + ` parameter can be used to select the timeline stats which should be present. Possible values are: ` + "`" + `` + "`" + `\"\"` + "`" + `` + "`" + ` (disable),` + "`" + `` + "`" + `\"24h\"` + "`" + `` + "`" + ` (default), ` + "`" + `` + "`" + `\"14d\"` + "`" + `` + "`" + `\n\nUser feedback items from the [User Feedback Widget](https://docs.sentry.io/product/user-feedback/#user-feedback-widget) are built off the issue platform, so to return a list of user feedback items for a specific project, filter for ` + "`" + `issue.category:feedback` + "`" + `.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the issues belong to.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the project the issues belong to.",
          "type": "string"
        },
        {
          "name": "statsPeriod",
          "in": "query",
          "required": false,
          "description": "An optional stat period (can be one of ` + "`" + `\"24h\"` + "`" + `, ` + "`" + `\"14d\"` + "`" + `, and ` + "`" + `\"\"` + "`" + `), defaults to \"24h\" if not provided.",
          "type": "string"
        },
        {
          "name": "shortIdLookup",
          "in": "query",
          "required": false,
          "description": "If this is set to true then short IDs are looked up by this function as well. This can cause the return value of the function to return an event issue of a different project which is why this is an opt-in. Set to 1 to enable.",
          "type": "boolean"
        },
        {
          "name": "query",
          "in": "query",
          "required": false,
          "description": "An optional Sentry structured search query. If not provided an implied ` + "`" + `\"is:unresolved\"` + "`" + ` is assumed.",
          "type": "string"
        },
        {
          "name": "hashes",
          "in": "query",
          "required": false,
          "description": "A list of hashes of groups to return. Is not compatible with 'query' parameter. The maximum number of hashes that can be sent is 100. If more are sent, only the first 100 will be used.",
          "type": "string"
        },
        {
          "name": "cursor",
          "in": "query",
          "required": false,
          "description": "A pointer to the last object fetched and its sort order; used to retrieve the next or previous results.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "403"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "project_id_or_slug"
        ],
        "query": [
          "statsPeriod",
          "shortIdLookup",
          "query",
          "hashes",
          "cursor"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "search",
      "pagination": {
        "style": "cursor",
        "cursorParam": "cursor"
      },
      "related": [
        {
          "id": "Bulk Mutate a List of Issues",
          "role": "same-resource"
        },
        {
          "id": "Bulk Remove a List of Issues",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Bulk Mutate a List of Issues",
      "method": "PUT",
      "path": "/api/0/projects/{organization_id_or_slug}/{project_id_or_slug}/issues/",
      "tag": "Events",
      "summary": "Bulk mutate various attributes on issues.  The list of issues to modify is given through the ` + "`" + `id` + "`" + ` query parameter.  It is repeated for each issue that should be modified.\n\n- For non-status updates, the ` + "`" + `id` + "`" + ` query parameter is required.\n- For status updates, the ` + "`" + `id` + "`" + ` query parameter may be omitted\nfor a batch \"update all\" query.\n- An optional ` + "`" + `status` + "`" + ` query parameter may be used to restrict\nmutations to only events with the given status.\n\nThe following attributes can be modified and are supplied as JSON object in the body:\n\nIf any IDs are out of scope this operation will succeed without any data mutation.",
      "description": "Bulk mutate various attributes on issues.  The list of issues to modify is given through the ` + "`" + `id` + "`" + ` query parameter.  It is repeated for each issue that should be modified.\n\n- For non-status updates, the ` + "`" + `id` + "`" + ` query parameter is required.\n- For status updates, the ` + "`" + `id` + "`" + ` query parameter may be omitted\nfor a batch \"update all\" query.\n- An optional ` + "`" + `status` + "`" + ` query parameter may be used to restrict\nmutations to only events with the given status.\n\nThe following attributes can be modified and are supplied as JSON object in the body:\n\nIf any IDs are out of scope this operation will succeed without any data mutation.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the issues belong to.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the project the issues belong to.",
          "type": "string"
        },
        {
          "name": "id",
          "in": "query",
          "required": false,
          "description": "A list of IDs of the issues to be mutated. This parameter shall be repeated for each issue. It is optional only if a status is mutated in which case an implicit update all is assumed.",
          "type": "integer"
        },
        {
          "name": "status",
          "in": "query",
          "required": false,
          "description": "Optionally limits the query to issues of the specified status. Valid values are ` + "`" + `\"resolved\"` + "`" + `, ` + "`" + `\"reprocessing\"` + "`" + `, ` + "`" + `\"unresolved\"` + "`" + `, and ` + "`" + `\"ignored\"` + "`" + `.",
          "type": "string"
        }
      ],
      "hasRequestBody": true,
      "risk": "write",
      "cacheable": false,
      "responseCodes": [
        "200",
        "400",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "project_id_or_slug"
        ],
        "query": [
          "id",
          "status"
        ],
        "headers": [],
        "body": true
      },
      "graphql": null,
      "kind": "update",
      "pagination": null,
      "related": [
        {
          "id": "List a Project's Issues",
          "role": "same-resource"
        },
        {
          "id": "Bulk Remove a List of Issues",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Bulk Remove a List of Issues",
      "method": "DELETE",
      "path": "/api/0/projects/{organization_id_or_slug}/{project_id_or_slug}/issues/",
      "tag": "Events",
      "summary": "Permanently remove the given issues. The list of issues to modify is given through the ` + "`" + `id` + "`" + ` query parameter.  It is repeated for each issue that should be removed.\n\nOnly queries by 'id' are accepted.\n\nIf any IDs are out of scope this operation will succeed without any data mutation.",
      "description": "Permanently remove the given issues. The list of issues to modify is given through the ` + "`" + `id` + "`" + ` query parameter.  It is repeated for each issue that should be removed.\n\nOnly queries by 'id' are accepted.\n\nIf any IDs are out of scope this operation will succeed without any data mutation.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the issues belong to.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the project the issues belong to.",
          "type": "string"
        },
        {
          "name": "id",
          "in": "query",
          "required": false,
          "description": "A list of IDs of the issues to be removed. This parameter shall be repeated for each issue, e.g. ` + "`" + `?id=1&id=2&id=3` + "`" + `. If this parameter is not provided, it will attempt to remove the first 1000 issues.",
          "type": "integer"
        }
      ],
      "hasRequestBody": false,
      "risk": "destructive",
      "cacheable": false,
      "responseCodes": [
        "204",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "project_id_or_slug"
        ],
        "query": [
          "id"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "delete",
      "pagination": null,
      "related": [
        {
          "id": "List a Project's Issues",
          "role": "same-resource"
        },
        {
          "id": "Bulk Mutate a List of Issues",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "List a Tag's Values for an Issue",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/issues/{issue_id}/tags/{key}/values/",
      "tag": "Events",
      "summary": "Returns a list of values associated with this key for an issue.\nReturns at most 1000 values when paginated.",
      "description": "Returns a list of values associated with this key for an issue.\nReturns at most 1000 values when paginated.",
      "parameters": [
        {
          "name": "issue_id",
          "in": "path",
          "required": true,
          "description": "The ID of the issue you'd like to query.",
          "type": "integer"
        },
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "key",
          "in": "path",
          "required": true,
          "description": "The tag key to look the values up for.",
          "type": "string"
        },
        {
          "name": "sort",
          "in": "query",
          "required": false,
          "description": "Sort order of the resulting tag values. Prefix with '-' for descending order. Default is '-id'.",
          "type": "string"
        },
        {
          "name": "environment",
          "in": "query",
          "required": false,
          "description": "The name of environments to filter by.",
          "type": "array"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "400",
        "401",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "issue_id",
          "organization_id_or_slug",
          "key"
        ],
        "query": [
          "sort",
          "environment"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": null,
      "related": []
    },
    {
      "id": "List an Issue's Hashes",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/issues/{issue_id}/hashes/",
      "tag": "Events",
      "summary": "This endpoint lists an issue's hashes, which are the generated checksums used to aggregate individual events.",
      "description": "This endpoint lists an issue's hashes, which are the generated checksums used to aggregate individual events.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the event belongs to.",
          "type": "string"
        },
        {
          "name": "issue_id",
          "in": "path",
          "required": true,
          "description": "The ID of the issue to retrieve.",
          "type": "string"
        },
        {
          "name": "full",
          "in": "query",
          "required": false,
          "description": "If this is set to true, the event payload will include the full event body, including the stacktrace. Set to 1 to enable.",
          "type": "boolean"
        },
        {
          "name": "cursor",
          "in": "query",
          "required": false,
          "description": "A pointer to the last object fetched and its sort order; used to retrieve the next or previous results.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "403"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "issue_id"
        ],
        "query": [
          "full",
          "cursor"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": {
        "style": "cursor",
        "cursorParam": "cursor"
      },
      "related": []
    },
    {
      "id": "Retrieve an Issue",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/issues/{issue_id}/",
      "tag": "Events",
      "summary": "Return details on an individual issue. This returns the basic stats for the issue (title, last seen, first seen), some overall numbers (number of comments, user reports) as well as the summarized event data.",
      "description": "Return details on an individual issue. This returns the basic stats for the issue (title, last seen, first seen), some overall numbers (number of comments, user reports) as well as the summarized event data.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the issue belongs to.",
          "type": "string"
        },
        {
          "name": "issue_id",
          "in": "path",
          "required": true,
          "description": "The ID of the issue to retrieve.",
          "type": "string"
        },
        {
          "name": "collapse",
          "in": "query",
          "required": false,
          "description": "Fields to remove from the response to improve query performance.",
          "type": "array"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "403"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "issue_id"
        ],
        "query": [
          "collapse"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": null,
      "related": [
        {
          "id": "Update an Issue",
          "role": "same-resource"
        },
        {
          "id": "Remove an Issue",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Update an Issue",
      "method": "PUT",
      "path": "/api/0/organizations/{organization_id_or_slug}/issues/{issue_id}/",
      "tag": "Events",
      "summary": "Updates an individual issue's attributes.  Only the attributes submitted are modified.",
      "description": "Updates an individual issue's attributes.  Only the attributes submitted are modified.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the issue belongs to.",
          "type": "string"
        },
        {
          "name": "issue_id",
          "in": "path",
          "required": true,
          "description": "The ID of the group to retrieve.",
          "type": "string"
        }
      ],
      "hasRequestBody": true,
      "risk": "write",
      "cacheable": false,
      "responseCodes": [
        "200",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "issue_id"
        ],
        "query": [],
        "headers": [],
        "body": true
      },
      "graphql": null,
      "kind": "update",
      "pagination": null,
      "related": [
        {
          "id": "Retrieve an Issue",
          "role": "same-resource"
        },
        {
          "id": "Remove an Issue",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Remove an Issue",
      "method": "DELETE",
      "path": "/api/0/organizations/{organization_id_or_slug}/issues/{issue_id}/",
      "tag": "Events",
      "summary": "Removes an individual issue.",
      "description": "Removes an individual issue.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the issue belongs to.",
          "type": "string"
        },
        {
          "name": "issue_id",
          "in": "path",
          "required": true,
          "description": "The ID of the issue to delete.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "destructive",
      "cacheable": false,
      "responseCodes": [
        "202",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "issue_id"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "delete",
      "pagination": null,
      "related": [
        {
          "id": "Retrieve an Issue",
          "role": "same-resource"
        },
        {
          "id": "Update an Issue",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "List an Organization's Releases",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/releases/",
      "tag": "Releases",
      "summary": "Return a list of releases for a given organization.",
      "description": "Return a list of releases for a given organization.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization.",
          "type": "string"
        },
        {
          "name": "query",
          "in": "query",
          "required": false,
          "description": "This parameter can be used to create a \"starts with\" filter for the version.",
          "type": "string"
        },
        {
          "name": "cursor",
          "in": "query",
          "required": false,
          "description": "A pointer to the last object fetched and its sort order; used to retrieve the next or previous results.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "401",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug"
        ],
        "query": [
          "query",
          "cursor"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": {
        "style": "cursor",
        "cursorParam": "cursor"
      },
      "related": [
        {
          "id": "Create a New Release for an Organization",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Create a New Release for an Organization",
      "method": "POST",
      "path": "/api/0/organizations/{organization_id_or_slug}/releases/",
      "tag": "Releases",
      "summary": "Create a new release for the given organization.  Releases are used by\nSentry to improve its error reporting abilities by correlating\nfirst seen events with the release that might have introduced the\nproblem.\nReleases are also necessary for source maps and other debug features\nthat require manual upload for functioning well.",
      "description": "Create a new release for the given organization.  Releases are used by\nSentry to improve its error reporting abilities by correlating\nfirst seen events with the release that might have introduced the\nproblem.\nReleases are also necessary for source maps and other debug features\nthat require manual upload for functioning well.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization.",
          "type": "string"
        }
      ],
      "hasRequestBody": true,
      "risk": "write",
      "cacheable": false,
      "responseCodes": [
        "201",
        "400",
        "403"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug"
        ],
        "query": [],
        "headers": [],
        "body": true
      },
      "graphql": null,
      "kind": "create",
      "pagination": null,
      "related": [
        {
          "id": "List an Organization's Releases",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "List an Organization's Release Files",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/releases/{version}/files/",
      "tag": "Releases",
      "summary": "Return a list of files for a given release.",
      "description": "Return a list of files for a given release.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization.",
          "type": "string"
        },
        {
          "name": "version",
          "in": "path",
          "required": true,
          "description": "The version identifier of the release.",
          "type": "string"
        },
        {
          "name": "cursor",
          "in": "query",
          "required": false,
          "description": "A pointer to the last object fetched and its sort order; used to retrieve the next or previous results.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "version"
        ],
        "query": [
          "cursor"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": {
        "style": "cursor",
        "cursorParam": "cursor"
      },
      "related": [
        {
          "id": "Upload a New Organization Release File",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Upload a New Organization Release File",
      "method": "POST",
      "path": "/api/0/organizations/{organization_id_or_slug}/releases/{version}/files/",
      "tag": "Releases",
      "summary": "Upload a new file for the given release.\n\nUnlike other API requests, files must be uploaded using the traditional multipart/form-data content-type.\n\nRequests to this endpoint should use the region-specific domain eg. ` + "`" + `us.sentry.io` + "`" + ` or ` + "`" + `de.sentry.io` + "`" + `.\n\nThe optional 'name' attribute should reflect the absolute path that this file will be referenced as. For example, in the case of JavaScript you might specify the full web URI.",
      "description": "Upload a new file for the given release.\n\nUnlike other API requests, files must be uploaded using the traditional multipart/form-data content-type.\n\nRequests to this endpoint should use the region-specific domain eg. ` + "`" + `us.sentry.io` + "`" + ` or ` + "`" + `de.sentry.io` + "`" + `.\n\nThe optional 'name' attribute should reflect the absolute path that this file will be referenced as. For example, in the case of JavaScript you might specify the full web URI.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization.",
          "type": "string"
        },
        {
          "name": "version",
          "in": "path",
          "required": true,
          "description": "The version identifier of the release.",
          "type": "string"
        }
      ],
      "hasRequestBody": true,
      "risk": "write",
      "cacheable": false,
      "responseCodes": [
        "201",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "version"
        ],
        "query": [],
        "headers": [],
        "body": true
      },
      "graphql": null,
      "kind": "create",
      "pagination": null,
      "related": [
        {
          "id": "List an Organization's Release Files",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "List a Project's Release Files",
      "method": "GET",
      "path": "/api/0/projects/{organization_id_or_slug}/{project_id_or_slug}/releases/{version}/files/",
      "tag": "Releases",
      "summary": "Return a list of files for a given release.",
      "description": "Return a list of files for a given release.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the project.",
          "type": "string"
        },
        {
          "name": "version",
          "in": "path",
          "required": true,
          "description": "The version identifier of the release.",
          "type": "string"
        },
        {
          "name": "cursor",
          "in": "query",
          "required": false,
          "description": "A pointer to the last object fetched and its sort order; used to retrieve the next or previous results.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "project_id_or_slug",
          "version"
        ],
        "query": [
          "cursor"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": {
        "style": "cursor",
        "cursorParam": "cursor"
      },
      "related": [
        {
          "id": "Upload a New Project Release File",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Upload a New Project Release File",
      "method": "POST",
      "path": "/api/0/projects/{organization_id_or_slug}/{project_id_or_slug}/releases/{version}/files/",
      "tag": "Releases",
      "summary": "Upload a new file for the given release.\n\nUnlike other API requests, files must be uploaded using the traditional multipart/form-data content-type.\n\nRequests to this endpoint should use the region-specific domain eg. ` + "`" + `us.sentry.io` + "`" + ` or ` + "`" + `de.sentry.io` + "`" + `\n\nThe optional 'name' attribute should reflect the absolute path that this file will be referenced as. For example, in the case of JavaScript you might specify the full web URI.",
      "description": "Upload a new file for the given release.\n\nUnlike other API requests, files must be uploaded using the traditional multipart/form-data content-type.\n\nRequests to this endpoint should use the region-specific domain eg. ` + "`" + `us.sentry.io` + "`" + ` or ` + "`" + `de.sentry.io` + "`" + `\n\nThe optional 'name' attribute should reflect the absolute path that this file will be referenced as. For example, in the case of JavaScript you might specify the full web URI.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the project.",
          "type": "string"
        },
        {
          "name": "version",
          "in": "path",
          "required": true,
          "description": "The version identifier of the release.",
          "type": "string"
        }
      ],
      "hasRequestBody": true,
      "risk": "write",
      "cacheable": false,
      "responseCodes": [
        "201",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "project_id_or_slug",
          "version"
        ],
        "query": [],
        "headers": [],
        "body": true
      },
      "graphql": null,
      "kind": "create",
      "pagination": null,
      "related": [
        {
          "id": "List a Project's Release Files",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Retrieve an Organization Release's File",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/releases/{version}/files/{file_id}/",
      "tag": "Releases",
      "summary": "Retrieve a file for a given release.",
      "description": "Retrieve a file for a given release.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization.",
          "type": "string"
        },
        {
          "name": "version",
          "in": "path",
          "required": true,
          "description": "The version identifier of the release.",
          "type": "string"
        },
        {
          "name": "file_id",
          "in": "path",
          "required": true,
          "description": "The ID of the file to retrieve.",
          "type": "string"
        },
        {
          "name": "download",
          "in": "query",
          "required": false,
          "description": "If this is set to true, then the response payload will be the raw file contents. Otherwise, the response will be the file metadata as JSON.",
          "type": "boolean"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "version",
          "file_id"
        ],
        "query": [
          "download"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": null,
      "related": [
        {
          "id": "Update an Organization Release File",
          "role": "same-resource"
        },
        {
          "id": "Delete an Organization Release's File",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Update an Organization Release File",
      "method": "PUT",
      "path": "/api/0/organizations/{organization_id_or_slug}/releases/{version}/files/{file_id}/",
      "tag": "Releases",
      "summary": "Update an organization release file.",
      "description": "Update an organization release file.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization.",
          "type": "string"
        },
        {
          "name": "version",
          "in": "path",
          "required": true,
          "description": "The version identifier of the release.",
          "type": "string"
        },
        {
          "name": "file_id",
          "in": "path",
          "required": true,
          "description": "The ID of the file to retrieve.",
          "type": "string"
        }
      ],
      "hasRequestBody": true,
      "risk": "write",
      "cacheable": false,
      "responseCodes": [
        "200",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "version",
          "file_id"
        ],
        "query": [],
        "headers": [],
        "body": true
      },
      "graphql": null,
      "kind": "update",
      "pagination": null,
      "related": [
        {
          "id": "Retrieve an Organization Release's File",
          "role": "same-resource"
        },
        {
          "id": "Delete an Organization Release's File",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Delete an Organization Release's File",
      "method": "DELETE",
      "path": "/api/0/organizations/{organization_id_or_slug}/releases/{version}/files/{file_id}/",
      "tag": "Releases",
      "summary": "Delete a file for a given release.",
      "description": "Delete a file for a given release.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the release belongs to.",
          "type": "string"
        },
        {
          "name": "version",
          "in": "path",
          "required": true,
          "description": "The version identifier of the release.",
          "type": "string"
        },
        {
          "name": "file_id",
          "in": "path",
          "required": true,
          "description": "The ID of the file to delete.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "destructive",
      "cacheable": false,
      "responseCodes": [
        "204",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "version",
          "file_id"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "delete",
      "pagination": null,
      "related": [
        {
          "id": "Retrieve an Organization Release's File",
          "role": "same-resource"
        },
        {
          "id": "Update an Organization Release File",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Retrieve a Project Release's File",
      "method": "GET",
      "path": "/api/0/projects/{organization_id_or_slug}/{project_id_or_slug}/releases/{version}/files/{file_id}/",
      "tag": "Releases",
      "summary": "Retrieve a file for a given release.",
      "description": "Retrieve a file for a given release.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the project.",
          "type": "string"
        },
        {
          "name": "version",
          "in": "path",
          "required": true,
          "description": "The version identifier of the release.",
          "type": "string"
        },
        {
          "name": "file_id",
          "in": "path",
          "required": true,
          "description": "The ID of the file to retrieve.",
          "type": "string"
        },
        {
          "name": "download",
          "in": "query",
          "required": false,
          "description": "If this is set to true, then the response payload will be the raw file contents. Otherwise, the response will be the file metadata as JSON.",
          "type": "boolean"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "project_id_or_slug",
          "version",
          "file_id"
        ],
        "query": [
          "download"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": null,
      "related": [
        {
          "id": "Update a Project Release File",
          "role": "same-resource"
        },
        {
          "id": "Delete a Project Release's File",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Update a Project Release File",
      "method": "PUT",
      "path": "/api/0/projects/{organization_id_or_slug}/{project_id_or_slug}/releases/{version}/files/{file_id}/",
      "tag": "Releases",
      "summary": "Update a project release file.",
      "description": "Update a project release file.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the project.",
          "type": "string"
        },
        {
          "name": "version",
          "in": "path",
          "required": true,
          "description": "The version identifier of the release.",
          "type": "string"
        },
        {
          "name": "file_id",
          "in": "path",
          "required": true,
          "description": "The ID of the file to retrieve.",
          "type": "string"
        }
      ],
      "hasRequestBody": true,
      "risk": "write",
      "cacheable": false,
      "responseCodes": [
        "200",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "project_id_or_slug",
          "version",
          "file_id"
        ],
        "query": [],
        "headers": [],
        "body": true
      },
      "graphql": null,
      "kind": "update",
      "pagination": null,
      "related": [
        {
          "id": "Retrieve a Project Release's File",
          "role": "same-resource"
        },
        {
          "id": "Delete a Project Release's File",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Delete a Project Release's File",
      "method": "DELETE",
      "path": "/api/0/projects/{organization_id_or_slug}/{project_id_or_slug}/releases/{version}/files/{file_id}/",
      "tag": "Releases",
      "summary": "Delete a file for a given release.",
      "description": "Delete a file for a given release.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the release belongs to.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the project.",
          "type": "string"
        },
        {
          "name": "version",
          "in": "path",
          "required": true,
          "description": "The version identifier of the release.",
          "type": "string"
        },
        {
          "name": "file_id",
          "in": "path",
          "required": true,
          "description": "The ID of the file to delete.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "destructive",
      "cacheable": false,
      "responseCodes": [
        "204",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "project_id_or_slug",
          "version",
          "file_id"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "delete",
      "pagination": null,
      "related": [
        {
          "id": "Retrieve a Project Release's File",
          "role": "same-resource"
        },
        {
          "id": "Update a Project Release File",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "List an Organization Release's Commits",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/releases/{version}/commits/",
      "tag": "Releases",
      "summary": "List an organization release's commits.",
      "description": "List an organization release's commits.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the release belongs to.",
          "type": "string"
        },
        {
          "name": "version",
          "in": "path",
          "required": true,
          "description": "The version identifier of the release.",
          "type": "string"
        },
        {
          "name": "cursor",
          "in": "query",
          "required": false,
          "description": "A pointer to the last object fetched and its sort order; used to retrieve the next or previous results.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "version"
        ],
        "query": [
          "cursor"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": {
        "style": "cursor",
        "cursorParam": "cursor"
      },
      "related": []
    },
    {
      "id": "List a Project Release's Commits",
      "method": "GET",
      "path": "/api/0/projects/{organization_id_or_slug}/{project_id_or_slug}/releases/{version}/commits/",
      "tag": "Releases",
      "summary": "List a project release's commits.",
      "description": "List a project release's commits.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the release belongs to.",
          "type": "string"
        },
        {
          "name": "project_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the project the release belongs to.",
          "type": "string"
        },
        {
          "name": "version",
          "in": "path",
          "required": true,
          "description": "The version identifier of the release.",
          "type": "string"
        },
        {
          "name": "cursor",
          "in": "query",
          "required": false,
          "description": "A pointer to the last object fetched and its sort order; used to retrieve the next or previous results.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "project_id_or_slug",
          "version"
        ],
        "query": [
          "cursor"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": {
        "style": "cursor",
        "cursorParam": "cursor"
      },
      "related": []
    },
    {
      "id": "Retrieve Files Changed in a Release's Commits",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/releases/{version}/commitfiles/",
      "tag": "Releases",
      "summary": "Retrieve files changed in a release's commits",
      "description": "Retrieve files changed in a release's commits",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the release belongs to.",
          "type": "string"
        },
        {
          "name": "version",
          "in": "path",
          "required": true,
          "description": "The version identifier of the release.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "version"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": null,
      "related": []
    },
    {
      "id": "List an Organization's Integration Platform Installations",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/sentry-app-installations/",
      "tag": "Integration",
      "summary": "Return a list of integration platform installations for a given organization.",
      "description": "Return a list of integration platform installations for a given organization.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The organization short name.",
          "type": "string"
        },
        {
          "name": "cursor",
          "in": "query",
          "required": false,
          "description": "A pointer to the last object fetched and its sort order; used to retrieve the next or previous results.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug"
        ],
        "query": [
          "cursor"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": {
        "style": "cursor",
        "cursorParam": "cursor"
      },
      "related": []
    },
    {
      "id": "Create or update an External Issue",
      "method": "POST",
      "path": "/api/0/sentry-app-installations/{uuid}/external-issues/",
      "tag": "Integration",
      "summary": "Create or update an external issue from an integration platform integration.",
      "description": "Create or update an external issue from an integration platform integration.",
      "parameters": [
        {
          "name": "uuid",
          "in": "path",
          "required": true,
          "description": "The uuid of the integration platform integration.",
          "type": "string"
        }
      ],
      "hasRequestBody": true,
      "risk": "write",
      "cacheable": false,
      "responseCodes": [
        "200",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "uuid"
        ],
        "query": [],
        "headers": [],
        "body": true
      },
      "graphql": null,
      "kind": "create",
      "pagination": null,
      "related": []
    },
    {
      "id": "Delete an External Issue",
      "method": "DELETE",
      "path": "/api/0/sentry-app-installations/{uuid}/external-issues/{external_issue_id}/",
      "tag": "Integration",
      "summary": "Delete an external issue.",
      "description": "Delete an external issue.",
      "parameters": [
        {
          "name": "uuid",
          "in": "path",
          "required": true,
          "description": "The uuid of the integration platform integration.",
          "type": "string"
        },
        {
          "name": "external_issue_id",
          "in": "path",
          "required": true,
          "description": "The ID of the external issue.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "destructive",
      "cacheable": false,
      "responseCodes": [
        "204",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "uuid",
          "external_issue_id"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "delete",
      "pagination": null,
      "related": []
    },
    {
      "id": "Enable Spike Protection",
      "method": "POST",
      "path": "/api/0/organizations/{organization_id_or_slug}/spike-protections/",
      "tag": "Projects",
      "summary": "Enables Spike Protection feature for some of the projects within the organization.",
      "description": "Enables Spike Protection feature for some of the projects within the organization.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the projects belong to",
          "type": "string"
        }
      ],
      "hasRequestBody": true,
      "risk": "write",
      "cacheable": false,
      "responseCodes": [
        "201",
        "400",
        "403"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug"
        ],
        "query": [],
        "headers": [],
        "body": true
      },
      "graphql": null,
      "kind": "create",
      "pagination": null,
      "related": [
        {
          "id": "Disable Spike Protection",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Disable Spike Protection",
      "method": "DELETE",
      "path": "/api/0/organizations/{organization_id_or_slug}/spike-protections/",
      "tag": "Projects",
      "summary": "Disables Spike Protection feature for some of the projects within the organization.",
      "description": "Disables Spike Protection feature for some of the projects within the organization.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the projects belong to",
          "type": "string"
        }
      ],
      "hasRequestBody": true,
      "risk": "destructive",
      "cacheable": false,
      "responseCodes": [
        "200",
        "400",
        "403"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug"
        ],
        "query": [],
        "headers": [],
        "body": true
      },
      "graphql": null,
      "kind": "delete",
      "pagination": null,
      "related": [
        {
          "id": "Enable Spike Protection",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Retrieve Seer Issue Fix State",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/issues/{issue_id}/autofix/",
      "tag": "Seer",
      "summary": "Retrieve the current detailed state of an issue fix process for a specific issue including:\n\n- Current status\n- Steps performed and their outcomes\n- Repository information and permissions\n- Root Cause Analysis\n- Proposed Solution\n- Generated code changes\n\nThis endpoint although documented is still experimental and the payload may change in the future.",
      "description": "Retrieve the current detailed state of an issue fix process for a specific issue including:\n\n- Current status\n- Steps performed and their outcomes\n- Repository information and permissions\n- Root Cause Analysis\n- Proposed Solution\n- Generated code changes\n\nThis endpoint although documented is still experimental and the payload may change in the future.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "issue_id",
          "in": "path",
          "required": true,
          "description": "The ID of the issue you'd like to query.",
          "type": "integer"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "401",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "issue_id"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": null,
      "related": [
        {
          "id": "Start Seer Issue Fix",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "Start Seer Issue Fix",
      "method": "POST",
      "path": "/api/0/organizations/{organization_id_or_slug}/issues/{issue_id}/autofix/",
      "tag": "Seer",
      "summary": "Trigger a Seer Issue Fix run for a specific issue.\n\nThe issue fix process can:\n- Identify the root cause of the issue\n- Propose a solution\n- Generate code changes\n- Create a pull request with the fix\n\nThe process runs asynchronously, and you can get the state using the GET endpoint.",
      "description": "Trigger a Seer Issue Fix run for a specific issue.\n\nThe issue fix process can:\n- Identify the root cause of the issue\n- Propose a solution\n- Generate code changes\n- Create a pull request with the fix\n\nThe process runs asynchronously, and you can get the state using the GET endpoint.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "issue_id",
          "in": "path",
          "required": true,
          "description": "The ID of the issue you'd like to query.",
          "type": "integer"
        }
      ],
      "hasRequestBody": true,
      "risk": "write",
      "cacheable": false,
      "responseCodes": [
        "202",
        "400",
        "401",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "issue_id"
        ],
        "query": [],
        "headers": [],
        "body": true
      },
      "graphql": null,
      "kind": "create",
      "pagination": null,
      "related": [
        {
          "id": "Retrieve Seer Issue Fix State",
          "role": "same-resource"
        }
      ]
    },
    {
      "id": "List an Issue's Events",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/issues/{issue_id}/events/",
      "tag": "Events",
      "summary": "Return a list of error events bound to an issue",
      "description": "Return a list of error events bound to an issue",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "issue_id",
          "in": "path",
          "required": true,
          "description": "The ID of the issue you'd like to query.",
          "type": "integer"
        },
        {
          "name": "start",
          "in": "query",
          "required": false,
          "description": "The start of the period of time for the query, expected in ISO-8601 format. For example, ` + "`" + `2001-12-14T12:34:56.7890` + "`" + `.",
          "type": "string"
        },
        {
          "name": "end",
          "in": "query",
          "required": false,
          "description": "The end of the period of time for the query, expected in ISO-8601 format. For example, ` + "`" + `2001-12-14T12:34:56.7890` + "`" + `.",
          "type": "string"
        },
        {
          "name": "statsPeriod",
          "in": "query",
          "required": false,
          "description": "The period of time for the query, will override the start & end parameters, a number followed by one of:\n- ` + "`" + `d` + "`" + ` for days\n- ` + "`" + `h` + "`" + ` for hours\n- ` + "`" + `m` + "`" + ` for minutes\n- ` + "`" + `s` + "`" + ` for seconds\n- ` + "`" + `w` + "`" + ` for weeks\n\nFor example, ` + "`" + `24h` + "`" + `, to mean query data starting from 24 hours ago to now.",
          "type": "string"
        },
        {
          "name": "environment",
          "in": "query",
          "required": false,
          "description": "The name of environments to filter by.",
          "type": "array"
        },
        {
          "name": "full",
          "in": "query",
          "required": false,
          "description": "Specify true to include the full event body, including the stacktrace, in the event payload.",
          "type": "boolean"
        },
        {
          "name": "sample",
          "in": "query",
          "required": false,
          "description": "Return events in pseudo-random order. This is deterministic so an identical query will always return the same events in the same order.",
          "type": "boolean"
        },
        {
          "name": "query",
          "in": "query",
          "required": false,
          "description": "An optional search query for filtering events. See [search syntax](https://docs.sentry.io/concepts/search/) and queryable event properties at [Sentry Search Documentation](https://docs.sentry.io/concepts/search/searchable-properties/events/) for more information. An example query might be ` + "`" + `query=transaction:foo AND release:abc` + "`" + `",
          "type": "string"
        },
        {
          "name": "cursor",
          "in": "query",
          "required": false,
          "description": "A pointer to the last object fetched and its sort order; used to retrieve the next or previous results.",
          "type": "string"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "400",
        "401",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "issue_id"
        ],
        "query": [
          "start",
          "end",
          "statsPeriod",
          "environment",
          "full",
          "sample",
          "query",
          "cursor"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": {
        "style": "cursor",
        "cursorParam": "cursor"
      },
      "related": []
    },
    {
      "id": "Retrieve an Issue Event",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/issues/{issue_id}/events/{event_id}/",
      "tag": "Events",
      "summary": "Retrieves the details of an issue event.",
      "description": "Retrieves the details of an issue event.",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "issue_id",
          "in": "path",
          "required": true,
          "description": "The ID of the issue you'd like to query.",
          "type": "integer"
        },
        {
          "name": "event_id",
          "in": "path",
          "required": true,
          "description": "The ID of the event to retrieve, or 'latest', 'oldest', or 'recommended'.",
          "type": "string"
        },
        {
          "name": "environment",
          "in": "query",
          "required": false,
          "description": "The name of environments to filter by.",
          "type": "array"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "400",
        "401",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "issue_id",
          "event_id"
        ],
        "query": [
          "environment"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": null,
      "related": []
    },
    {
      "id": "Retrieve custom integration issue links for the given Sentry issue",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/issues/{issue_id}/external-issues/",
      "tag": "Integration",
      "summary": "Retrieve custom integration issue links for the given Sentry issue",
      "description": "Retrieve custom integration issue links for the given Sentry issue",
      "parameters": [
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "issue_id",
          "in": "path",
          "required": true,
          "description": "The ID of the issue you'd like to query.",
          "type": "integer"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200"
      ],
      "inputHints": {
        "path": [
          "organization_id_or_slug",
          "issue_id"
        ],
        "query": [],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": null,
      "related": []
    },
    {
      "id": "Retrieve Tag Details",
      "method": "GET",
      "path": "/api/0/organizations/{organization_id_or_slug}/issues/{issue_id}/tags/{key}/",
      "tag": "Events",
      "summary": "Return a list of values associated with this key for an issue. When paginated can return at most 1000 values.",
      "description": "Return a list of values associated with this key for an issue. When paginated can return at most 1000 values.",
      "parameters": [
        {
          "name": "issue_id",
          "in": "path",
          "required": true,
          "description": "The ID of the issue you'd like to query.",
          "type": "integer"
        },
        {
          "name": "organization_id_or_slug",
          "in": "path",
          "required": true,
          "description": "The ID or slug of the organization the resource belongs to.",
          "type": "string"
        },
        {
          "name": "key",
          "in": "path",
          "required": true,
          "description": "The tag key to look the values up for.",
          "type": "string"
        },
        {
          "name": "environment",
          "in": "query",
          "required": false,
          "description": "The name of environments to filter by.",
          "type": "array"
        }
      ],
      "hasRequestBody": false,
      "risk": "read",
      "cacheable": true,
      "responseCodes": [
        "200",
        "400",
        "401",
        "403",
        "404"
      ],
      "inputHints": {
        "path": [
          "issue_id",
          "organization_id_or_slug",
          "key"
        ],
        "query": [
          "environment"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null,
      "kind": "read",
      "pagination": null,
      "related": []
    }
  ],
  "heroes": [
    {
      "alias": "organizations",
      "operationId": "List Your Organizations",
      "summary": "Return a list of organizations available to the authenticated session in a region.\nThis is particularly useful for requests with a user bound context. For API key-based requests this will only return the organization that belongs to the key.",
      "method": "GET",
      "path": "/api/0/organizations/",
      "defaultParams": {},
      "explicit": false
    },
    {
      "alias": "models",
      "operationId": "List Seer AI Models",
      "summary": "Get list of actively used LLM model names from Seer.\n\nReturns the list of AI models that are currently used in production in Seer.\nThis endpoint does not require authentication and can be used to discover which models Seer uses.\n\nRequests to this endpoint should use the region-specific domain\neg. ` + "`" + `us.sentry.io` + "`" + ` or ` + "`" + `de.sentry.io` + "`" + `",
      "method": "GET",
      "path": "/api/0/seer/models/",
      "defaultParams": {},
      "explicit": false
    }
  ],
  "insights": {
    "thesis": "sentry should become a local, searchable operational workspace, not only an endpoint wrapper.",
    "generatedAdvantages": [
      "One command surface for humans and agents",
      "Shared core between CLI and MCP server",
      "Local cache with search and replay metadata",
      "Risk-aware dry-run defaults for write/destructive calls",
      "Machine-readable manifest for catalogs and scorecards"
    ],
    "recommendedCommands": [
      "sync: cache list endpoints locally",
      "inspect: fetch entity details by id",
      "search: query cached API responses offline",
      "plan: preview write operations before execution"
    ],
    "domainMap": [
      {
        "tag": "Alerts",
        "operations": 5
      },
      {
        "tag": "Crons",
        "operations": 10
      },
      {
        "tag": "Dashboards",
        "operations": 5
      },
      {
        "tag": "Discover",
        "operations": 5
      },
      {
        "tag": "Environments",
        "operations": 4
      },
      {
        "tag": "Events",
        "operations": 17
      },
      {
        "tag": "Explore",
        "operations": 2
      },
      {
        "tag": "Integration",
        "operations": 8
      },
      {
        "tag": "Integrations",
        "operations": 14
      },
      {
        "tag": "Mobile Builds",
        "operations": 4
      },
      {
        "tag": "Monitors",
        "operations": 14
      },
      {
        "tag": "Organizations",
        "operations": 16
      },
      {
        "tag": "Prevent",
        "operations": 10
      },
      {
        "tag": "Projects",
        "operations": 35
      },
      {
        "tag": "Releases",
        "operations": 22
      },
      {
        "tag": "Replays",
        "operations": 12
      },
      {
        "tag": "SCIM",
        "operations": 10
      },
      {
        "tag": "Seer",
        "operations": 3
      },
      {
        "tag": "Teams",
        "operations": 12
      },
      {
        "tag": "Users",
        "operations": 1
      }
    ],
    "riskNotes": [
      "34 destructive operation(s) require explicit --yes in generated clients."
    ]
  },
  "generatedAt": "2026-05-16T01:33:25.840Z",
  "provenance": {
    "schemaVersion": "gutenberg.provenance.v1",
    "generatedAt": "2026-05-16T01:33:25.936Z",
    "gutenbergVersion": "0.1.0",
    "name": "sentry",
    "spec": {
      "path": "/tmp/sentry-api.json",
      "sha256": "94d354dc45cf56de7ee81227925ab384c5e05ad0a1e9fc37778193764365043e",
      "size": 3259546
    },
    "recipe": null,
    "targets": [
      "go",
      "mcp",
      "skill",
      "openclaw"
    ]
  },
  "defaultHeaders": {},
  "policy": {
    "schemaVersion": "gutenberg.policy.v1",
    "rules": [
      {
        "risk": "read",
        "action": "allow",
        "requiresYes": false
      },
      {
        "risk": "write",
        "action": "confirm",
        "requiresYes": true
      },
      {
        "risk": "destructive",
        "action": "confirm",
        "requiresYes": true
      }
    ],
    "redaction": [
      "authorization",
      "cookie",
      "token",
      "secret",
      "api-key",
      "apikey",
      "client-secret",
      "session"
    ]
  },
  "packageName": "gutenberg.local/sentry",
  "generatedBy": "gutenberg",
  "generatedByVersion": "0.3.0",
  "language": "go"
}`

type Manifest struct {
	SchemaVersion  string            `json:"schemaVersion"`
	Source         string            `json:"source"`
	Name           string            `json:"name"`
	Slug           string            `json:"slug"`
	EnvPrefix      string            `json:"envPrefix"`
	Description    string            `json:"description"`
	Version        string            `json:"version"`
	BaseURLs       []string          `json:"baseUrls"`
	Auth           Auth              `json:"auth"`
	Tags           []string          `json:"tags"`
	Operations     []Operation       `json:"operations"`
	Heroes         []Hero            `json:"heroes"`
	DefaultHeaders map[string]string `json:"defaultHeaders"`
	Policy         Policy            `json:"policy"`
	Insights       Insights          `json:"insights"`
}

type Hero struct {
	Alias       string `json:"alias"`
	OperationID string `json:"operationId"`
	Summary     string `json:"summary"`
	Method      string `json:"method"`
	Path        string `json:"path"`
}

type Auth struct {
	Mode    string       `json:"mode"`
	Schemes []AuthScheme `json:"schemes"`
	Env     string       `json:"env"`
	OAuth   bool         `json:"oauth"`
}

type AuthScheme struct {
	Name   string               `json:"name"`
	Type   string               `json:"type"`
	In     string               `json:"in"`
	Header string               `json:"header"`
	Scheme string               `json:"scheme"`
	Flows  map[string]OAuthFlow `json:"flows"`
}

type OAuthFlow struct {
	AuthorizationURL string            `json:"authorizationUrl"`
	TokenURL         string            `json:"tokenUrl"`
	RefreshURL       string            `json:"refreshUrl"`
	Scopes           map[string]string `json:"scopes"`
}

type Operation struct {
	ID             string       `json:"id"`
	Method         string       `json:"method"`
	Path           string       `json:"path"`
	Tag            string       `json:"tag"`
	Summary        string       `json:"summary"`
	Description    string       `json:"description"`
	Parameters     []Parameter  `json:"parameters"`
	HasRequestBody bool         `json:"hasRequestBody"`
	Risk           string       `json:"risk"`
	Cacheable      bool         `json:"cacheable"`
	ResponseCodes  []string     `json:"responseCodes"`
	GraphQL        *GraphQLSpec `json:"graphql,omitempty"`
	Kind           string       `json:"kind,omitempty"`
	Pagination     *Pagination  `json:"pagination,omitempty"`
}

type Pagination struct {
	Style        string `json:"style"`
	OffsetParam  string `json:"offsetParam,omitempty"`
	LimitParam   string `json:"limitParam,omitempty"`
	CursorParam  string `json:"cursorParam,omitempty"`
	PageParam    string `json:"pageParam,omitempty"`
	PerPageParam string `json:"perPageParam,omitempty"`
}

type GraphQLSpec struct {
	Kind  string       `json:"kind"`
	Field string       `json:"field"`
	Args  []GraphQLArg `json:"args"`
}

type GraphQLArg struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Parameter struct {
	Name        string `json:"name"`
	In          string `json:"in"`
	Required    bool   `json:"required"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

type Insights struct {
	Thesis              string   `json:"thesis"`
	GeneratedAdvantages []string `json:"generatedAdvantages"`
	RecommendedCommands []string `json:"recommendedCommands"`
	RiskNotes           []string `json:"riskNotes"`
}

type Policy struct {
	SchemaVersion string       `json:"schemaVersion"`
	Rules         []PolicyRule `json:"rules"`
	Redaction     []string     `json:"redaction"`
}

type PolicyRule struct {
	Risk        string `json:"risk"`
	Action      string `json:"action"`
	RequiresYes bool   `json:"requiresYes"`
}

func LoadManifest() Manifest {
	var manifest Manifest
	if err := json.Unmarshal([]byte(manifestJSON), &manifest); err != nil {
		panic(err)
	}
	return manifest
}

func Operations() []Operation {
	return LoadManifest().Operations
}

func GetOperation(id string) (Operation, bool) {
	for _, operation := range Operations() {
		if operation.ID == id {
			return operation, true
		}
	}
	return Operation{}, false
}

func Heroes() []Hero {
	return LoadManifest().Heroes
}

func FindHero(alias string) *Hero {
	for _, hero := range Heroes() {
		if hero.Alias == alias {
			return &hero
		}
	}
	return nil
}

func PolicyFor(operation Operation) PolicyRule {
	manifest := LoadManifest()
	for _, rule := range manifest.Policy.Rules {
		if rule.Risk == operation.Risk {
			return rule
		}
	}
	if operation.Risk == "read" {
		return PolicyRule{Risk: "read", Action: "allow"}
	}
	return PolicyRule{Risk: operation.Risk, Action: "confirm", RequiresYes: true}
}

func RequiresConfirmation(operation Operation) bool {
	rule := PolicyFor(operation)
	return rule.RequiresYes || rule.Action == "confirm"
}

func PolicyDenies(operation Operation) bool {
	return PolicyFor(operation).Action == "deny"
}

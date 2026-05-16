package aggr

// Sources is the embedded list of source manifests for multi-demo.
var Sources = []SourceDescriptor{
	{
		Slug:      "espn",
		Name:      "espn",
		Operation: "getApisSiteV2SportsBasketballNbaScoreboard",
		Manifest:  `{
  "schemaVersion": "gutenberg.blueprint.v1",
  "source": "/tmp/espn.openapi.json",
  "name": "espn",
  "slug": "espn",
  "envPrefix": "ESPN",
  "description": "Generated from a browser HAR capture by Gutenberg.",
  "version": "0.1.0",
  "baseUrls": [
    "https://site.api.espn.com"
  ],
  "auth": {
    "mode": "none",
    "schemes": [],
    "env": null
  },
  "tags": [
    "apis"
  ],
  "operations": [
    {
      "id": "getApisSiteV2SportsBasketballNbaScoreboard",
      "method": "GET",
      "path": "/apis/site/v2/sports/basketball/nba/scoreboard",
      "tag": "apis",
      "summary": "GET /apis/site/v2/sports/basketball/nba/scoreboard",
      "description": "",
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
      "graphql": null
    },
    {
      "id": "getApisSiteV2SportsSoccerIta1Scoreboard",
      "method": "GET",
      "path": "/apis/site/v2/sports/soccer/ita.1/scoreboard",
      "tag": "apis",
      "summary": "GET /apis/site/v2/sports/soccer/ita.1/scoreboard",
      "description": "",
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
      "graphql": null
    },
    {
      "id": "getApisSiteV2SportsFootballNflScoreboard",
      "method": "GET",
      "path": "/apis/site/v2/sports/football/nfl/scoreboard",
      "tag": "apis",
      "summary": "GET /apis/site/v2/sports/football/nfl/scoreboard",
      "description": "",
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
      "graphql": null
    },
    {
      "id": "getApisSiteV2SportsBaseballMlbScoreboard",
      "method": "GET",
      "path": "/apis/site/v2/sports/baseball/mlb/scoreboard",
      "tag": "apis",
      "summary": "GET /apis/site/v2/sports/baseball/mlb/scoreboard",
      "description": "",
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
      "graphql": null
    }
  ],
  "insights": {
    "thesis": "espn should become a local, searchable operational workspace, not only an endpoint wrapper.",
    "generatedAdvantages": [
      "One command surface for humans and agents",
      "Shared core between CLI and MCP server",
      "Local cache with search and replay metadata",
      "Risk-aware dry-run defaults for write/destructive calls",
      "Machine-readable manifest for catalogs and scorecards"
    ],
    "recommendedCommands": [
      "sync: cache list endpoints locally",
      "search: query cached API responses offline"
    ],
    "domainMap": [
      {
        "tag": "apis",
        "operations": 4
      }
    ],
    "riskNotes": []
  },
  "generatedAt": "2026-05-10T22:29:55.445Z"
}`,
	},
	{
		Slug:      "open-meteo",
		Name:      "open-meteo",
		Operation: "getV1Forecast",
		Manifest:  `{
  "schemaVersion": "gutenberg.blueprint.v1",
  "source": "/tmp/open-meteo.openapi.json",
  "name": "open-meteo",
  "slug": "open-meteo",
  "envPrefix": "OPEN_METEO",
  "description": "Generated from a browser HAR capture by Gutenberg.",
  "version": "0.1.0",
  "baseUrls": [
    "https://api.open-meteo.com"
  ],
  "auth": {
    "mode": "none",
    "schemes": [],
    "env": null
  },
  "tags": [
    "v1"
  ],
  "operations": [
    {
      "id": "getV1Forecast",
      "method": "GET",
      "path": "/v1/forecast",
      "tag": "v1",
      "summary": "GET /v1/forecast",
      "description": "",
      "parameters": [
        {
          "name": "latitude",
          "in": "query",
          "required": false,
          "description": "",
          "type": "number"
        },
        {
          "name": "longitude",
          "in": "query",
          "required": false,
          "description": "",
          "type": "number"
        },
        {
          "name": "current",
          "in": "query",
          "required": false,
          "description": "",
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
        "path": [],
        "query": [
          "latitude",
          "longitude",
          "current"
        ],
        "headers": [],
        "body": false
      },
      "graphql": null
    }
  ],
  "insights": {
    "thesis": "open-meteo should become a local, searchable operational workspace, not only an endpoint wrapper.",
    "generatedAdvantages": [
      "One command surface for humans and agents",
      "Shared core between CLI and MCP server",
      "Local cache with search and replay metadata",
      "Risk-aware dry-run defaults for write/destructive calls",
      "Machine-readable manifest for catalogs and scorecards"
    ],
    "recommendedCommands": [
      "sync: cache list endpoints locally",
      "search: query cached API responses offline"
    ],
    "domainMap": [
      {
        "tag": "v1",
        "operations": 1
      }
    ],
    "riskNotes": []
  },
  "generatedAt": "2026-05-10T22:29:55.450Z"
}`,
	},
}

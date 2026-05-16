package forge

import "encoding/json"

const manifestJSON = `{
  "schemaVersion": "gutenberg.blueprint.v1",
  "source": "/tmp/pypi.openapi.json",
  "name": "pypi",
  "slug": "pypi",
  "envPrefix": "PYPI",
  "description": "Generated from a browser HAR capture by Gutenberg.",
  "version": "0.1.0",
  "baseUrls": [
    "https://pypi.org"
  ],
  "auth": {
    "mode": "none",
    "schemes": [],
    "env": null
  },
  "tags": [
    "pypi"
  ],
  "operations": [
    {
      "id": "getPypiRequests",
      "method": "GET",
      "path": "/pypi/requests/json",
      "tag": "pypi",
      "summary": "GET /pypi/requests/json",
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
      "graphql": null,
      "kind": "read-list",
      "pagination": null,
      "related": []
    },
    {
      "id": "getPypiNumpy",
      "method": "GET",
      "path": "/pypi/numpy/json",
      "tag": "pypi",
      "summary": "GET /pypi/numpy/json",
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
      "graphql": null,
      "kind": "read-list",
      "pagination": null,
      "related": []
    }
  ],
  "heroes": [
    {
      "alias": "requests",
      "operationId": "getPypiRequests",
      "summary": "GET /pypi/requests/json",
      "method": "GET",
      "path": "/pypi/requests/json",
      "defaultParams": {},
      "explicit": false
    },
    {
      "alias": "numpy",
      "operationId": "getPypiNumpy",
      "summary": "GET /pypi/numpy/json",
      "method": "GET",
      "path": "/pypi/numpy/json",
      "defaultParams": {},
      "explicit": false
    }
  ],
  "insights": {
    "thesis": "pypi should become a local, searchable operational workspace, not only an endpoint wrapper.",
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
        "tag": "pypi",
        "operations": 2
      }
    ],
    "riskNotes": []
  },
  "generatedAt": "2026-05-16T10:47:47.919Z",
  "provenance": {
    "schemaVersion": "gutenberg.provenance.v1",
    "generatedAt": "2026-05-16T10:47:47.943Z",
    "gutenbergVersion": "0.1.0",
    "name": "pypi",
    "spec": {
      "path": "/tmp/pypi.openapi.json",
      "sha256": "8a0e881a3b6c3daf7ad020da3013355b6031a728cb6f6aa3102b5903f4166bb3",
      "size": 920167
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
  "packageName": "gutenberg.local/pypi",
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

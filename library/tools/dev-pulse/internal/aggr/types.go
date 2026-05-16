package aggr

import (
	"encoding/json"
	"time"
)

type SourceDescriptor struct {
	Slug      string
	Name      string
	Operation string
	Manifest  string
}

func (source SourceDescriptor) parsedManifest() (*Manifest, error) {
	var mf Manifest
	if err := json.Unmarshal([]byte(source.Manifest), &mf); err != nil {
		return nil, err
	}
	return &mf, nil
}

type Manifest struct {
	Name       string      `json:"name"`
	Slug       string      `json:"slug"`
	BaseURLs   []string    `json:"baseUrls"`
	Operations []Operation `json:"operations"`
}

func (m *Manifest) findOperation(id string) *Operation {
	for i := range m.Operations {
		if m.Operations[i].ID == id {
			return &m.Operations[i]
		}
	}
	if id == "" && len(m.Operations) > 0 {
		return &m.Operations[0]
	}
	return nil
}

type Operation struct {
	ID         string      `json:"id"`
	Method     string      `json:"method"`
	Path       string      `json:"path"`
	Parameters []Parameter `json:"parameters"`
}

type Parameter struct {
	Name     string `json:"name"`
	In       string `json:"in"`
	Required bool   `json:"required"`
}

type Item struct {
	Source    string    `json:"source"`
	Operation string    `json:"operation"`
	Status    int       `json:"status,omitempty"`
	URL       string    `json:"url,omitempty"`
	FetchedAt time.Time `json:"fetchedAt,omitempty"`
	Error     string    `json:"error,omitempty"`
	Payload   any       `json:"payload,omitempty"`
}

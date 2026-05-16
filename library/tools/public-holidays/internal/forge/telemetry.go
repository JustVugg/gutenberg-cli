package forge

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// Telemetry: opt-in local JSONL log of operation calls. NEVER sent anywhere.
// Enabled when GUTENBERG_TELEMETRY=1. Override path via GUTENBERG_TELEMETRY_FILE.
//
// Audit: risky operations (write/destructive), including dry-runs, are always
// written to a local JSONL audit file. Override path via GUTENBERG_AUDIT_FILE.

type telemetryEvent struct {
	Timestamp   string `json:"ts"`
	Tool        string `json:"tool"`
	OperationID string `json:"operationId"`
	Risk        string `json:"risk,omitempty"`
	Status      int    `json:"status,omitempty"`
	ElapsedMs   int64  `json:"elapsedMs,omitempty"`
	DryRun      bool   `json:"dryRun,omitempty"`
	Error       string `json:"error,omitempty"`
}

func telemetryFile() string {
	if override := os.Getenv("GUTENBERG_TELEMETRY_FILE"); override != "" {
		return override
	}
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return filepath.Join(".gutenberg", "usage.jsonl")
	}
	return filepath.Join(home, ".gutenberg", "usage.jsonl")
}

func auditFile() string {
	if override := os.Getenv("GUTENBERG_AUDIT_FILE"); override != "" {
		return override
	}
	return filepath.Join(".gutenberg", "audit.jsonl")
}

func LogCall(operationID string, status int, elapsed time.Duration, dryRun bool, callErr error) {
	operation, ok := GetOperation(operationID)
	risk := ""
	if ok {
		risk = operation.Risk
	}
	event := telemetryEvent{
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
		Tool:        "public-holidays",
		OperationID: operationID,
		Risk:        risk,
		Status:      status,
		ElapsedMs:   elapsed.Milliseconds(),
		DryRun:      dryRun,
	}
	if callErr != nil {
		event.Error = callErr.Error()
	}
	if risk != "" && risk != "read" {
		writeEvent(auditFile(), event)
	}
	if os.Getenv("GUTENBERG_TELEMETRY") != "1" {
		return
	}
	writeEvent(telemetryFile(), event)
}

func writeEvent(path string, event telemetryEvent) {
	line, err := json.Marshal(event)
	if err != nil {
		return
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return
	}
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return
	}
	defer file.Close()
	_, _ = file.Write(append(line, '\n'))
}

package forge

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

type Cache struct {
	Version int           `json:"version"`
	Records []CacheRecord `json:"records"`
}

type CacheRecord struct {
	OperationID string           `json:"operationId"`
	Request     RequestPlan      `json:"request"`
	Response    *ResponseEnvelope `json:"response,omitempty"`
	CachedAt    string           `json:"cachedAt"`
}

type CacheStats struct {
	File        string         `json:"file"`
	Records     int            `json:"records"`
	ByOperation map[string]int `json:"byOperation"`
	FTS5        bool           `json:"fts5"`
}

func cacheFile() string {
	manifest := LoadManifest()
	if override := os.Getenv(manifest.EnvPrefix + "_SQLITE_FILE"); override != "" {
		return override
	}
	if override := os.Getenv(manifest.EnvPrefix + "_CACHE_FILE"); override != "" {
		return override
	}
	return filepath.Join(".gutenberg", manifest.Slug+".sqlite")
}

func openStore() (*sql.DB, error) {
	file := cacheFile()
	if err := os.MkdirAll(filepath.Dir(file), 0o755); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", file)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	if err := migrateStore(db); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

// migrations is an append-only list. Each entry is applied exactly once,
// in order. NEVER renumber, edit, or delete a past migration — add a new one.
var migrations = []struct {
	ID   int
	Name string
	SQL  string
}{
	{
		ID:   1,
		Name: "initial",
		SQL: `CREATE TABLE records (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			operation_id TEXT NOT NULL,
			request_json TEXT NOT NULL,
			response_json TEXT,
			cached_at TEXT NOT NULL
		);
		CREATE VIRTUAL TABLE records_fts USING fts5(operation_id, request_json, response_json, cached_at);`,
	},
}

func migrateStore(db *sql.DB) error {
	pragmas := []string{
		"PRAGMA busy_timeout=5000",
		"PRAGMA journal_mode=WAL",
		"PRAGMA synchronous=NORMAL",
	}
	for _, statement := range pragmas {
		if _, err := db.Exec(statement); err != nil {
			return err
		}
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS gutenberg_migrations (id INTEGER PRIMARY KEY, name TEXT NOT NULL, applied_at TEXT NOT NULL)`); err != nil {
		return err
	}
	applied := map[int]bool{}
	rows, err := db.Query("SELECT id FROM gutenberg_migrations")
	if err != nil {
		return err
	}
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			rows.Close()
			return err
		}
		applied[id] = true
	}
	rows.Close()
	for _, migration := range migrations {
		if applied[migration.ID] {
			continue
		}
		tx, err := db.Begin()
		if err != nil {
			return err
		}
		if _, err := tx.Exec(migration.SQL); err != nil {
			tx.Rollback()
			return fmt.Errorf("migration %d (%s) failed: %w", migration.ID, migration.Name, err)
		}
		if _, err := tx.Exec("INSERT INTO gutenberg_migrations (id, name, applied_at) VALUES (?, ?, ?)", migration.ID, migration.Name, time.Now().UTC().Format(time.RFC3339)); err != nil {
			tx.Rollback()
			return err
		}
		if err := tx.Commit(); err != nil {
			return err
		}
	}
	return nil
}

// AppliedMigrations returns the list of applied migration IDs (for diagnostics).
func AppliedMigrations() ([]int, error) {
	db, err := openStore()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	rows, err := db.Query("SELECT id FROM gutenberg_migrations ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ids := []int{}
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func ReadCache() (Cache, error) {
	db, err := openStore()
	if err != nil {
		return Cache{}, err
	}
	defer db.Close()
	rows, err := db.Query("SELECT operation_id, request_json, response_json, cached_at FROM records ORDER BY id")
	if err != nil {
		return Cache{}, err
	}
	defer rows.Close()
	cache := Cache{Version: 1, Records: []CacheRecord{}}
	for rows.Next() {
		record, err := scanCacheRecord(rows)
		if err != nil {
			return Cache{}, err
		}
		cache.Records = append(cache.Records, record)
	}
	return cache, rows.Err()
}

func WriteCache(cache Cache) error {
	db, err := openStore()
	if err != nil {
		return err
	}
	defer db.Close()
	if _, err := db.Exec("DELETE FROM records; DELETE FROM records_fts"); err != nil {
		return err
	}
	for _, record := range cache.Records {
		if _, err := saveRecordWithDB(db, record); err != nil {
			return err
		}
	}
	return nil
}

func SaveRecord(record CacheRecord) (CacheRecord, error) {
	db, err := openStore()
	if err != nil {
		return CacheRecord{}, err
	}
	defer db.Close()
	return saveRecordWithDB(db, record)
}

func saveRecordWithDB(db *sql.DB, record CacheRecord) (CacheRecord, error) {
	record.CachedAt = time.Now().UTC().Format(time.RFC3339)
	requestJSON, err := json.Marshal(record.Request)
	if err != nil {
		return CacheRecord{}, err
	}
	responseJSON, err := json.Marshal(record.Response)
	if err != nil {
		return CacheRecord{}, err
	}
	if _, err := db.Exec(
		"INSERT INTO records (operation_id, request_json, response_json, cached_at) VALUES (?, ?, ?, ?)",
		record.OperationID, string(requestJSON), string(responseJSON), record.CachedAt,
	); err != nil {
		return CacheRecord{}, err
	}
	_, _ = db.Exec(
		"INSERT INTO records_fts (operation_id, request_json, response_json, cached_at) VALUES (?, ?, ?, ?)",
		record.OperationID, string(requestJSON), string(responseJSON), record.CachedAt,
	)
	return record, nil
}

func SearchCache(query string) ([]CacheRecord, error) {
	db, err := openStore()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	if strings.TrimSpace(query) == "" {
		return latestRecords(db)
	}
	rows, err := db.Query("SELECT operation_id, request_json, response_json, cached_at FROM records_fts WHERE records_fts MATCH ? LIMIT 50", ftsQuery(query))
	if err != nil {
		rows, err = db.Query("SELECT operation_id, request_json, response_json, cached_at FROM records WHERE lower(operation_id || ' ' || request_json || ' ' || response_json) LIKE ? LIMIT 50", "%"+strings.ToLower(query)+"%")
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	results := []CacheRecord{}
	for rows.Next() {
		record, err := scanCacheRecord(rows)
		if err != nil {
			return nil, err
		}
		results = append(results, record)
	}
	return results, rows.Err()
}

func GetCacheStats() (CacheStats, error) {
	db, err := openStore()
	if err != nil {
		return CacheStats{}, err
	}
	defer db.Close()
	stats := CacheStats{File: cacheFile(), ByOperation: map[string]int{}, FTS5: true}
	if err := db.QueryRow("SELECT COUNT(*) FROM records").Scan(&stats.Records); err != nil {
		return CacheStats{}, err
	}
	rows, err := db.Query("SELECT operation_id, COUNT(*) FROM records GROUP BY operation_id ORDER BY operation_id")
	if err != nil {
		return CacheStats{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var operation string
		var count int
		if err := rows.Scan(&operation, &count); err != nil {
			return CacheStats{}, err
		}
		stats.ByOperation[operation] = count
	}
	return stats, rows.Err()
}

func latestRecords(db *sql.DB) ([]CacheRecord, error) {
	rows, err := db.Query("SELECT operation_id, request_json, response_json, cached_at FROM records ORDER BY id DESC LIMIT 50")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	results := []CacheRecord{}
	for rows.Next() {
		record, err := scanCacheRecord(rows)
		if err != nil {
			return nil, err
		}
		results = append(results, record)
	}
	return results, rows.Err()
}

func scanCacheRecord(scanner interface{ Scan(dest ...any) error }) (CacheRecord, error) {
	var operationID, requestJSON, responseJSON, cachedAt string
	if err := scanner.Scan(&operationID, &requestJSON, &responseJSON, &cachedAt); err != nil {
		return CacheRecord{}, err
	}
	record := CacheRecord{OperationID: operationID, CachedAt: cachedAt}
	if err := json.Unmarshal([]byte(requestJSON), &record.Request); err != nil {
		return CacheRecord{}, err
	}
	if responseJSON != "" && responseJSON != "null" {
		var response ResponseEnvelope
		if err := json.Unmarshal([]byte(responseJSON), &response); err != nil {
			return CacheRecord{}, err
		}
		record.Response = &response
	}
	return record, nil
}

func ftsQuery(query string) string {
	parts := strings.Fields(query)
	if len(parts) == 0 {
		return "*"
	}
	terms := []string{}
	for _, part := range parts {
		clean := strings.Map(func(r rune) rune {
			if r == '_' || r == '-' || r == '.' || r == '/' || r == ':' {
				return ' '
			}
			if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' {
				return r
			}
			return -1
		}, part)
		for _, token := range strings.Fields(clean) {
			terms = append(terms, fmt.Sprintf("%q*", token))
		}
	}
	if len(terms) == 0 {
		return "*"
	}
	return strings.Join(terms, " OR ")
}

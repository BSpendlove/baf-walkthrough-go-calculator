package history

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const maxEntries = 100

// DefaultPath returns the default history file path (~/.calc_history).
var DefaultPath = filepath.Join(os.Getenv("HOME"), ".calc_history")

// Entry represents a single history record.
type Entry struct {
	Expr   string  `json:"expr"`
	Result float64 `json:"result"`
	Time   string  `json:"time"`
}

// Store appends an entry to the history file.
func Store(path, expr string, result float64) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error: cannot open history file: %w", err)
	}
	defer f.Close()

	e := Entry{
		Expr:   expr,
		Result: result,
		Time:   time.Now().Format(time.RFC3339),
	}
	return json.NewEncoder(f).Encode(e)
}

// Load reads up to the last maxEntries entries from the history file.
func Load(path string) ([]Entry, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("error: cannot open history file: %w", err)
	}
	defer f.Close()

	var entries []Entry
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var e Entry
		if err := json.Unmarshal(scanner.Bytes(), &e); err != nil {
			continue // skip malformed lines
		}
		entries = append(entries, e)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error: reading history file: %w", err)
	}

	if len(entries) > maxEntries {
		entries = entries[len(entries)-maxEntries:]
	}
	return entries, nil
}

// Clear truncates the history file.
func Clear(path string) error {
	return os.Truncate(path, 0)
}

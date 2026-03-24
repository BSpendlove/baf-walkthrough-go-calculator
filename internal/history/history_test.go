package history

import (
	"os"
	"path/filepath"
	"testing"
)

func tempHistoryPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), ".calc_history")
}

func TestStoreAndLoad(t *testing.T) {
	path := tempHistoryPath(t)

	if err := Store(path, "2 + 3", 5); err != nil {
		t.Fatalf("Store error: %v", err)
	}
	if err := Store(path, "10 / 2", 5); err != nil {
		t.Fatalf("Store error: %v", err)
	}

	entries, err := Load(path)
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("got %d entries, want 2", len(entries))
	}
	if entries[0].Expr != "2 + 3" || entries[0].Result != 5 {
		t.Errorf("entry[0] = %+v, want expr='2 + 3' result=5", entries[0])
	}
	if entries[1].Expr != "10 / 2" || entries[1].Result != 5 {
		t.Errorf("entry[1] = %+v, want expr='10 / 2' result=5", entries[1])
	}
}

func TestLoadCapsAt100(t *testing.T) {
	path := tempHistoryPath(t)

	for i := 0; i < 120; i++ {
		if err := Store(path, "1 + 1", 2); err != nil {
			t.Fatalf("Store error: %v", err)
		}
	}

	entries, err := Load(path)
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}
	if len(entries) != 100 {
		t.Errorf("got %d entries, want 100", len(entries))
	}
}

func TestClear(t *testing.T) {
	path := tempHistoryPath(t)

	if err := Store(path, "1 + 1", 2); err != nil {
		t.Fatalf("Store error: %v", err)
	}
	if err := Clear(path); err != nil {
		t.Fatalf("Clear error: %v", err)
	}

	entries, err := Load(path)
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("got %d entries after clear, want 0", len(entries))
	}
}

func TestLoadMissingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nonexistent")

	entries, err := Load(path)
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}
	if entries != nil {
		t.Errorf("got %v, want nil for missing file", entries)
	}
}

func TestLoadEmptyFile(t *testing.T) {
	path := tempHistoryPath(t)
	if err := os.WriteFile(path, []byte{}, 0644); err != nil {
		t.Fatal(err)
	}

	entries, err := Load(path)
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}
	if entries != nil {
		t.Errorf("got %v, want nil for empty file", entries)
	}
}

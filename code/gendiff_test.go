package code

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func getInt(v any) (int64, bool) {
	switch x := v.(type) {
	case json.Number:
		n, err := x.Int64()
		return n, err == nil
	case float64:
		return int64(x), x == float64(int64(x))
	case int:
		return int64(x), true
	case int64:
		return x, true
	default:
		return 0, false
	}
}

func TestParseFile_RelativePath(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "a", "b"), 0o755); err != nil {
		t.Fatal(err)
	}
	full := filepath.Join(dir, "a", "b", "config.json")
	wantMap := map[string]any{"ok": true}

	if err := os.WriteFile(full, []byte(`{"ok": true}`), 0o644); err != nil {
		t.Fatal(err)
	}

	// рабочая директория = dir, чтобы относительный путь резолвился
	wd, _ := os.Getwd()
	defer os.Chdir(wd)
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	rel := filepath.Join("a", "b", "..", "b", "config.json")

	got, err := parseFile(rel)
	if err != nil {
		t.Fatal(err)
	}

	if got.Ext != ".json" {
		t.Fatalf("Ext = %q, want %q", got.Ext, ".json")
	}

	data, ok := got.Data.(map[string]any)
	if !ok {
		t.Fatalf("Data type = %T, want map[string]any", got.Data)
	}
	if !reflect.DeepEqual(data, wantMap) {
		t.Errorf("Data = %#v, want %#v", data, wantMap)
	}
}

func TestParseFile_AbsolutePath(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	full := filepath.Join(dir, "config.json")

	if err := os.WriteFile(full, []byte(`{"ok": 42}`), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := parseFile(full)
	if err != nil {
		t.Fatal(err)
	}

	if got.Ext != ".json" {
		t.Fatalf("Ext = %q, want %q", got.Ext, ".json")
	}

	data, ok := got.Data.(map[string]any)
	if !ok {
		t.Fatalf("Data type = %T, want map[string]any", got.Data)
	}

	val, exists := data["ok"]
	if !exists {
		t.Fatalf(`key "ok" not found`)
	}
	n, ok := getInt(val)
	if !ok || n != 42 {
		t.Errorf(`"ok" = %#v (type %T), want 42`, val, val)
	}
}

func TestParseFile_NotExists(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	_, err := parseFile(filepath.Join(dir, "missing.json"))
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("expected fs.ErrNotExist, got %v", err)
	}
}

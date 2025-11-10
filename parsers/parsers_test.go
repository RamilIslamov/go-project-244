package parsers

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

func TestParseFile_JSONDecodeError(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	p := filepath.Join(dir, "bad.json")

	if err := os.WriteFile(p, []byte(`{"a":`), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := parseFile(p); err == nil {
		t.Fatalf("expected json decode error, got nil")
	}
}

func TestParseFile_YAMLDecodeError(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	p := filepath.Join(dir, "bad.yaml")

	if err := os.WriteFile(p, []byte(": bad"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := parseFile(p); err == nil {
		t.Fatalf("expected yaml decode error, got nil")
	}
}

func TestParseFile_AbsError(t *testing.T) {
	t.Parallel()

	bad := "invalid\x00path.json"
	if _, err := parseFile(bad); err == nil {
		t.Fatalf("expected error for path with NUL byte, got nil")
	}
}

func TestParseFile_RelativePath(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "a", "b"), 0o755); err != nil {
		t.Fatal(err)
	}
	full := filepath.Join(dir, "a", "b", "config.json")
	wantMap := map[string]any{"ok": true}
	if err := os.WriteFile(full, []byte(`{"ok": true}`), 0o644); err != nil {
		t.Fatal(err)
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := os.Chdir(wd); err != nil {
			t.Fatalf("restore cwd: %v", err)
		}
	})

	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	rel := filepath.Join("a", "b", "..", "b", "config.json")
	got, err := parseFile(rel)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, wantMap) {
		t.Fatalf("data mismatch:\n got: %#v\nwant: %#v", got, wantMap)
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

	val, exists := got["ok"]
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

func TestParseFiles_OK(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	p1 := filepath.Join(dir, "a.json")
	p2 := filepath.Join(dir, "b.json")
	f1 := []byte(`{"a": 1, "b": true}`)
	f2 := []byte(`{"a": 2, "c": "x"}`)

	if err := os.WriteFile(p1, f1, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p2, f2, 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := ParseFiles(p1, p2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}

	aOne, ok := getInt(got[0]["a"])
	if !ok || aOne != 1 {
		t.Fatalf("file #1 no key 'a' or invalid val on key 'a'")
	}

	bOne, ok := got[0]["b"]
	if !ok || bOne != true {
		t.Fatalf("file #1 no key 'b' or invalid val on key 'b'")
	}

	aTwo, ok := getInt(got[1]["a"])
	if !ok || aTwo != 2 {
		t.Fatalf("file #2 no key 'a' or invalid val on key 'a'")
	}
	cTwo, ok := got[1]["c"]
	if !ok || cTwo != "x" {
		t.Fatalf("file #2 no key 'c' or invalid val on key 'c'")
	}
}

func TestParseFiles_Error(t *testing.T) {
	t.Parallel()

	_, err := ParseFiles("/no/such/file.json")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

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
	if got.Ext != ".json" {
		t.Fatalf("Ext = %q, want %q", got.Ext, ".json")
	}
	if !reflect.DeepEqual(got.Data, wantMap) {
		t.Fatalf("data mismatch:\n got: %#v\nwant: %#v", got.Data, wantMap)
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

	data := got.Data

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

func TestParseFiles_OK(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	p1 := filepath.Join(dir, "a.json")
	p2 := filepath.Join(dir, "b.json")

	// числа оставляем числами — parseFile уже приводит типы (json.Number/float64 и т.п.)
	if err := os.WriteFile(p1, []byte(`{"a": 1, "b": true}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p2, []byte(`{"a": 2, "c": "x"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := parseFiles(p1, p2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}
	if got[0].Ext != ".json" || got[1].Ext != ".json" {
		t.Fatalf("extensions = %q, %q; want .json", got[0].Ext, got[1].Ext)
	}
	// лёгкая sanity-проверка содержимого
	if _, ok := got[0].Data["a"]; !ok {
		t.Fatalf("file #1 data has no key 'a'")
	}
	if _, ok := got[1].Data["c"]; !ok {
		t.Fatalf("file #2 data has no key 'c'")
	}
}

func TestParseFiles_Error(t *testing.T) {
	t.Parallel()

	_, err := parseFiles("/no/such/file.json")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestDiff_OrderAndCases(t *testing.T) {
	t.Parallel()

	file1 := map[string]any{"a": 1, "b": 2}
	file2 := map[string]any{"a": 0, "b": 2, "c": 3}

	got := diff(file1, file2)

	want := []string{
		"  - a: 1",
		"  + a: 0",
		"    b: 2",
		"  + c: 3",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("diff mismatch:\n got: %#v\nwant: %#v", got, want)
	}
}

func TestFormatResults_NonEmpty(t *testing.T) {
	t.Parallel()

	in := []string{"  - a: 1", "    b: 2"}
	got := formatResults(in)
	want := "{\n  - a: 1\n    b: 2\n}"
	if got != want {
		t.Fatalf("formatResults:\n got: %q\nwant: %q", got, want)
	}
}

func TestFormatResults_Empty(t *testing.T) {
	t.Parallel()

	got := formatResults(nil)
	want := "{\n}"
	if got != want {
		t.Fatalf("formatResults empty:\n got: %q\nwant: %q", got, want)
	}
}

func TestGenDiff_FromTempFiles(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	p1 := filepath.Join(dir, "left.json")
	p2 := filepath.Join(dir, "right.json")

	if err := os.WriteFile(p1, []byte(`{"a":1,"b":2}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p2, []byte(`{"a":0,"b":2,"c":3}`), 0o644); err != nil {
		t.Fatal(err)
	}

	got := GenDiff(p1, p2)
	want := "{\n  - a: 1\n  + a: 0\n    b: 2\n  + c: 3\n}"

	if got != want {
		t.Fatalf("GenDiff:\n got:\n%q\nwant:\n%q", got, want)
	}
}

func TestDiff_KeyRemoved(t *testing.T) {
	t.Parallel()

	file1 := map[string]any{"only": 1}
	file2 := map[string]any{}

	got := diff(file1, file2)
	want := []string{"  - only: 1"}
	if len(got) != 1 || got[0] != want[0] {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

func TestParseFile_AbsError(t *testing.T) {
	t.Parallel()

	bad := "invalid\x00path.json"
	if _, err := parseFile(bad); err == nil {
		t.Fatalf("expected error for path with NUL byte, got nil")
	}
}

func TestParseFile_JSONDecodeError(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	p := filepath.Join(dir, "bad.json")
	// битый JSON
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

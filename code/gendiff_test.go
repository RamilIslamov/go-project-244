package code

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

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

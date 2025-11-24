package go_project_244

import (
	"os"
	"path/filepath"
	"testing"
)

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

	got, _ := GenDiff(p1, p2, "stylish")
	want := "{\n  - a: 1\n  + a: 0\n    b: 2\n  + c: 3\n}"

	if got != want {
		t.Fatalf("GenDiff:\n got:\n%q\nwant:\n%q", got, want)
	}
}

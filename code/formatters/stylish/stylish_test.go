package stylish

import (
	"github.com/RamilIslamov/go-project-244/code/ast"
	"strings"
	"testing"
)

func nl(s string) string {
	return strings.ReplaceAll(s, "\r\n", "\n")
}

func TestRender_SimpleDiff(t *testing.T) {
	nodes := []ast.Node{
		{Key: "a", Action: ast.Removed, OldVal: 1},
		{Key: "a", Action: ast.Added, NewVal: 0},
		{Key: "b", Action: ast.Unchanged, OldVal: 2},
		{Key: "c", Action: ast.Added, NewVal: 3},
	}
	got := Render(nodes)
	want := "{\n" +
		"  - a: 1\n" +
		"  + a: 0\n" +
		"    b: 2\n" +
		"  + c: 3\n" +
		"}"

	if nl(got) != nl(want) {
		t.Fatalf("mismatch\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

func TestRender_NestedAndUpdated(t *testing.T) {
	nodes := []ast.Node{
		{
			Key:    "root",
			Action: ast.Nested,
			Children: []ast.Node{
				{
					Key:    "x",
					Action: ast.Unchanged,
					OldVal: 1,
				},
				{
					Key:    "y",
					Action: ast.Updated,
					OldVal: 2,
					NewVal: 3,
				},
			},
		},
	}

	got := Render(nodes)

	want := "{\n" +
		"    root: {\n" +
		"        x: 1\n" +
		"      - y: 2\n" +
		"      + y: 3\n" +
		"    }\n" +
		"}"

	if nl(got) != nl(want) {
		t.Fatalf("nested mismatch\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

func TestStringify_Map(t *testing.T) {
	m := map[string]any{
		"b": 2,
		"a": 1,
	}
	got := stringify(m, 1)

	want := "{\n" +
		"        a: 1\n" +
		"        b: 2\n" +
		"  }"

	if nl(got) != nl(want) {
		t.Fatalf("stringify(map) mismatch\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

func TestStringify_Array(t *testing.T) {
	arr := []any{1, "x"}

	got := stringify(arr, 1)

	want := "[\n" +
		"        1\n" +
		"        x\n" +
		"  ]"

	if nl(got) != nl(want) {
		t.Fatalf("stringify([]any) mismatch\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

func TestStringify_NilAndPrimitive(t *testing.T) {
	if got := stringify(nil, 1); got != "null" {
		t.Fatalf("stringify(nil) = %q, want %q", got, "null")
	}

	if got := stringify(42, 1); got != "42" {
		t.Fatalf("stringify(42) = %q, want %q", got, "42")
	}
}

func TestIndent_EdgeCases(t *testing.T) {
	if got := indent(0); got != "" {
		t.Fatalf("indent(0) = %q, want empty string", got)
	}

	if got := indent(1); got != "  " {
		t.Fatalf("indent(1) = %q, want two spaces", got)
	}
}

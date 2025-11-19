package stylish

import (
	"code/code/ast"
	"strings"
	"testing"
)

func nl(s string) string { return strings.ReplaceAll(s, "\r\n", "\n") }

func TestRender_Stylish(t *testing.T) {
	nodes := []ast.Node{
		{Key: "a", Action: ast.Removed, OldVal: 1},
		{Key: "a", Action: ast.Added, NewVal: 0},
		{Key: "b", Action: ast.Unchanged, OldVal: 2},
		{Key: "c", Action: ast.Added, NewVal: 3},
	}
	got := Render(nodes)
	want := "{\n  - a: 1\n  + a: 0\n    b: 2\n  + c: 3\n}"
	if nl(got) != nl(want) {
		t.Fatalf("mismatch\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

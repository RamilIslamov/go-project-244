package json

import (
	"code/ast"
	stdjson "encoding/json"
	"reflect"
	"testing"
)

func TestToJSONNodes_AllActions(t *testing.T) {
	nodes := []ast.Node{
		{
			Key:    "added",
			Action: ast.Added,
			NewVal: 10,
		},
		{
			Key:    "removed",
			Action: ast.Removed,
			OldVal: "x",
		},
		{
			Key:    "updated",
			Action: ast.Updated,
			OldVal: 1,
			NewVal: 2,
		},
		{
			Key:    "unchanged",
			Action: ast.Unchanged,
			OldVal: true,
		},
		{
			Key:    "nested",
			Action: ast.Nested,
			Children: []ast.Node{
				{
					Key:    "child",
					Action: ast.Added,
					NewVal: "val",
				},
			},
		},
	}

	got := toJSONNodes(nodes)

	want := []ast.JsonNode{
		{
			Key:      "added",
			Type:     "added",
			NewValue: 10,
		},
		{
			Key:      "removed",
			Type:     "removed",
			OldValue: "x",
		},
		{
			Key:      "updated",
			Type:     "updated",
			OldValue: 1,
			NewValue: 2,
		},
		{
			Key:      "unchanged",
			Type:     "unchanged",
			OldValue: true,
		},
		{
			Key:  "nested",
			Type: "nested",
			Children: []ast.JsonNode{
				{
					Key:      "child",
					Type:     "added",
					NewValue: "val",
				},
			},
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("toJSONNodes mismatch\n got: %#v\nwant: %#v", got, want)
	}
}

func TestActionToString_AllKnown(t *testing.T) {
	cases := []struct {
		name string
		in   ast.NodeType
		want string
	}{
		{"added", ast.Added, "added"},
		{"removed", ast.Removed, "removed"},
		{"updated", ast.Updated, "updated"},
		{"nested", ast.Nested, "nested"},
		{"unchanged", ast.Unchanged, "unchanged"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := actionToString(tc.in)
			if got != tc.want {
				t.Fatalf("actionToString(%v) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestActionToString_Unknown(t *testing.T) {
	unknown := ast.NodeType("additional")

	got := actionToString(unknown)
	if got != "unknown" {
		t.Fatalf("actionToString(%v) = %q, want %q", unknown, got, "unknown")
	}
}

func TestRender_ReturnsValidJSON(t *testing.T) {
	nodes := []ast.Node{
		{
			Key:    "root",
			Action: ast.Nested,
			Children: []ast.Node{
				{Key: "a", Action: ast.Added, NewVal: 1},
				{Key: "b", Action: ast.Removed, OldVal: true},
			},
		},
	}

	out, err := Render(nodes)
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}

	if !stdjson.Valid([]byte(out)) {
		t.Fatalf("Render returned invalid JSON:\n%s", out)
	}

	var parsed struct {
		Diff []ast.JsonNode `json:"diff"`
	}
	if err := stdjson.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("unmarshal output: %v\njson: %s", err, out)
	}

	got := parsed.Diff

	if len(got) != 1 {
		t.Fatalf("len(top-level) = %d, want 1", len(got))
	}

	root := got[0]
	if root.Key != "root" || root.Type != "nested" {
		t.Fatalf("root = %#v, want key=root type=nested", root)
	}

	if len(root.Children) != 2 {
		t.Fatalf("len(root.Children) = %d, want 2", len(root.Children))
	}

	childA := root.Children[0]
	childB := root.Children[1]

	if childA.Key != "a" || childA.Type != "added" || childA.NewValue != float64(1) {
		t.Fatalf("childA = %#v, want key=a type=added newValue=1", childA)
	}
	if childB.Key != "b" || childB.Type != "removed" || childB.OldValue != true {
		t.Fatalf("childB = %#v, want key=b type=removed oldValue=true", childB)
	}
}

package ast

import (
	"testing"
)

func TestBuildDiff_FlatAndOrder(t *testing.T) {
	a := map[string]any{"a": 1, "b": 2}
	b := map[string]any{"a": 0, "b": 2, "c": 3}

	nodes := BuildDiff(a, b)

	if len(nodes) != 3 || nodes[0].Key != "a" || nodes[1].Key != "b" || nodes[2].Key != "c" {
		t.Fatalf("unexpected order/len: %#v", nodes)
	}
	if nodes[0].Action != Updated || nodes[1].Action != Unchanged || nodes[2].Action != Added {
		t.Fatalf("unexpected kinds: %#v", []Node{nodes[0], nodes[1], nodes[2]})
	}
}

func TestBuildDiff_Nested(t *testing.T) {
	a := map[string]any{"parent": map[string]any{"x": 1}}
	b := map[string]any{"parent": map[string]any{"x": 2, "y": 3}}

	nodes := BuildDiff(a, b)
	if len(nodes) != 1 || nodes[0].Action != Nested {
		t.Fatalf("want nested node, got %#v", nodes)
	}
	child := nodes[0].Children
	if len(child) != 2 || child[0].Key != "x" || child[1].Key != "y" {
		t.Fatalf("unexpected children: %#v", child)
	}
	if child[0].Action != Updated || child[1].Action != Added {
		t.Fatalf("unexpected child kinds: %#v", child)
	}
}

package plain

import (
	"github.com/RamilIslamov/go-project-244/code/ast"
	"testing"
)

func TestRenderPlain(t *testing.T) {
	nodes := []ast.Node{
		{
			Key:    "host",
			Action: ast.Updated,
			OldVal: "hexlet.io",
			NewVal: "hexlet.com",
		},
		{
			Key:    "timeout",
			Action: ast.Removed,
			OldVal: 50,
		},
		{
			Key:    "verbose",
			Action: ast.Added,
			NewVal: true,
		},
		{
			Key:    "common",
			Action: ast.Nested,
			Children: []ast.Node{
				{
					Key:    "setting6",
					Action: ast.Nested,
					Children: []ast.Node{
						{
							Key:    "doge",
							Action: ast.Nested,
							Children: []ast.Node{
								{
									Key:    "wow",
									Action: ast.Updated,
									OldVal: "",
									NewVal: "so much",
								},
							},
						},
					},
				},
			},
		},
		{
			Key:    "obj",
			Action: ast.Updated,
			OldVal: map[string]interface{}{"a": 1},
			NewVal: map[string]interface{}{"b": 2},
		},
	}

	got, _ := Render(nodes)

	want := "" +
		"Property 'host' was updated. From 'hexlet.io' to 'hexlet.com'\n" +
		"Property 'timeout' was removed\n" +
		"Property 'verbose' was added with value: true\n" +
		"Property 'common.setting6.doge.wow' was updated. From '' to 'so much'\n" +
		"Property 'obj' was updated. From [complex value] to [complex value]\n"

	if got != want {
		t.Fatalf("Render() result mismatch.\n--- got ---\n%q\n--- want ---\n%q\n", got, want)
	}
}

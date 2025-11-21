package formatters

import (
	"fmt"
	"github.com/RamilIslamov/go-project-244/code/ast"
	"github.com/RamilIslamov/go-project-244/code/formatters/json"
	"github.com/RamilIslamov/go-project-244/code/formatters/plain"
	"github.com/RamilIslamov/go-project-244/code/formatters/stylish"
)

func Render(format string, nodes []ast.Node) (string, error) {
	switch format {
	case "", "stylish":
		return stylish.Render(nodes)
	case "plain":
		return plain.Render(nodes)
	case "json":
		return json.Render(nodes)
	default:
		return "", fmt.Errorf("unknown format %q", format)
	}
}

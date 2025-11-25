package formatters

import (
	"code/ast"
	"code/formatters/json"
	"code/formatters/plain"
	"code/formatters/stylish"
	"fmt"
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

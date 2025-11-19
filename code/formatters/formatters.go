package formatters

import (
	"code/code/ast"
	"code/code/formatters/plain"
	"code/code/formatters/stylish"
	"fmt"
)

func Render(format string, nodes []ast.Node) (string, error) {
	switch format {
	case "", "stylish":
		return stylish.Render(nodes), nil
	case "plain":
		return plain.Render(nodes), nil
	default:
		return "", fmt.Errorf("unknown format %q", format)
	}
}

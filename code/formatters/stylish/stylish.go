package stylish

import (
	"fmt"
	"github.com/RamilIslamov/go-project-244/code/ast"
	"sort"
	"strings"
)

const indentSize = 4

func indent(depth int) string {
	n := depth*indentSize - 2
	if n < 0 {
		n = 0
	}
	return strings.Repeat(" ", n)
}

func Render(nodes []ast.Node) (string, error) {
	return render(nodes, 1)
}

func render(nodes []ast.Node, depth int) (string, error) {
	base := indent(depth)
	closeIndent := strings.Repeat(" ", (depth-1)*indentSize)

	var b strings.Builder
	b.WriteString("{\n")
	for _, n := range nodes {
		switch n.Action {
		case ast.Nested:
			childStr, err := render(n.Children, depth+1)
			if err != nil {
				return "", fmt.Errorf("render nested %q: %w", n.Key, err)
			}
			b.WriteString(fmt.Sprintf("%s  %s: %s\n", base, n.Key, childStr))
		case ast.Unchanged:
			b.WriteString(fmt.Sprintf("%s  %s: %s\n", base, n.Key, stringify(n.OldVal, depth+1)))
		case ast.Removed:
			b.WriteString(fmt.Sprintf("%s- %s: %s\n", base, n.Key, stringify(n.OldVal, depth+1)))
		case ast.Added:
			b.WriteString(fmt.Sprintf("%s+ %s: %s\n", base, n.Key, stringify(n.NewVal, depth+1)))
		case ast.Updated:
			b.WriteString(fmt.Sprintf("%s- %s: %s\n", base, n.Key, stringify(n.OldVal, depth+1)))
			b.WriteString(fmt.Sprintf("%s+ %s: %s\n", base, n.Key, stringify(n.NewVal, depth+1)))
		}
	}
	b.WriteString(closeIndent + "}")
	return b.String(), nil
}

func stringify(v any, depth int) string {
	if v == nil {
		return "null"
	}

	base := indent(depth)
	inner := strings.Repeat(" ", depth*indentSize+2)

	if m, ok := v.(map[string]any); ok {
		keys := make([]string, 0, len(m))
		for k := range m {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		var b strings.Builder
		b.WriteString("{\n")
		for _, k := range keys {
			b.WriteString(fmt.Sprintf("%s  %s: %s\n", inner, k, stringify(m[k], depth+1)))
		}

		b.WriteString(base + "}")
		return b.String()
	}

	if arr, ok := v.([]any); ok {
		var b strings.Builder
		b.WriteString("[\n")
		for _, el := range arr {
			b.WriteString(fmt.Sprintf("%s  %s\n", inner, stringify(el, depth+1)))
		}
		b.WriteString(base + "]")
		return b.String()
	}

	return fmt.Sprintf("%v", v)
}

package plain

import (
	"code/ast"
	"fmt"
	"strings"
)

func Render(nodes []ast.Node) (string, error) {
	return render(nodes, "")
}

func render(nodes []ast.Node, parentPath string) (string, error) {
	base := "Property"
	var b strings.Builder

	for _, n := range nodes {
		propPath := buildPath(parentPath, n.Key)

		switch n.Action {

		case ast.Nested:
			childStr, err := render(n.Children, propPath)
			if err != nil {
				return "", fmt.Errorf("render nested %q: %w", n.Key, err)
			}
			b.WriteString(childStr)

		case ast.Removed:
			fmt.Fprintf(&b, "%s '%s' was removed\n", base, propPath)

		case ast.Added:
			newValStr := formatPlainValue(n.NewVal)
			fmt.Fprintf(&b, "%s '%s' was added with value: %s\n", base, propPath, newValStr)

		case ast.Updated:
			oldValStr := formatPlainValue(n.OldVal)
			newValStr := formatPlainValue(n.NewVal)
			fmt.Fprintf(&b, "%s '%s' was updated. From %s to %s\n", base, propPath, oldValStr, newValStr)
		}
	}

	return b.String(), nil
}

func buildPath(parent, key string) string {
	if parent == "" {
		return key
	}
	return parent + "." + key
}

func formatPlainValue(v interface{}) string {
	switch vv := v.(type) {
	case map[string]interface{}, []interface{}:
		return "[complex value]"
	case string:
		return fmt.Sprintf("'%s'", vv)
	case nil:
		return "null"
	default:
		return fmt.Sprintf("%v", vv)
	}
}

package json

import (
	"code/ast"
	"encoding/json"
)

func Render(nodes []ast.Node) (string, error) {
	j := toJSONNodes(nodes)

	data, err := json.MarshalIndent(j, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func toJSONNodes(nodes []ast.Node) []ast.JsonNode {
	res := make([]ast.JsonNode, 0, len(nodes))

	for _, n := range nodes {
		j := ast.JsonNode{
			Key:  n.Key,
			Type: actionToString(n.Action),
		}

		switch n.Action {
		case ast.Nested:
			j.Children = toJSONNodes(n.Children)

		case ast.Added:
			j.NewValue = n.NewVal

		case ast.Removed:
			j.OldValue = n.OldVal

		case ast.Updated:
			j.OldValue = n.OldVal
			j.NewValue = n.NewVal

		case ast.Unchanged:
			j.OldValue = n.OldVal
		}

		res = append(res, j)
	}

	return res
}

func actionToString(a ast.NodeType) string {
	switch a {
	case ast.Added:
		return "added"
	case ast.Removed:
		return "removed"
	case ast.Updated:
		return "updated"
	case ast.Nested:
		return "nested"
	case ast.Unchanged:
		return "unchanged"
	default:
		return "unknown"
	}
}

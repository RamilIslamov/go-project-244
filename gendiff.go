package go_project_244

import (
	"code/ast"
	"code/formatters"
	"code/parsers"
	"fmt"
)

func GenDiff(path1, path2, format string) (string, error) {
	parsed, err := parsers.ParseFiles(path1, path2)
	if err != nil {
		return "", fmt.Errorf("parse files: %w", err)
	}
	nodes := ast.BuildDiff(parsed[0], parsed[1])

	return formatters.Render(format, nodes)
}

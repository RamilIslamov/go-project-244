package code

import (
	"fmt"
	"github.com/RamilIslamov/go-project-244/parsers"
	"log"
	"sort"
	"strings"
)

func GenDiff(path1, path2 string) string {
	parsed, err := parsers.ParseFiles(path1, path2)
	if err != nil {
		log.Fatalf("parse %q: %v", path1, err)
	}
	file1 := parsed[0]
	file2 := parsed[1]

	return formatResults(diff(file1, file2))
}

func diff(file1, file2 map[string]any) []string {
	union := make([]string, 0, len(file1)+len(file2))
	seen := make(map[string]struct{}, len(file1)+len(file2))
	for k := range file1 {
		seen[k] = struct{}{}
	}
	for k := range file2 {
		seen[k] = struct{}{}
	}
	for k := range seen {
		union = append(union, k)
	}
	sort.Strings(union)

	var result []string
	for _, k := range union {
		v1, ok1 := file1[k]
		v2, ok2 := file2[k]
		switch {
		case ok1 && !ok2:
			result = append(result, fmt.Sprintf("  %s: %v", "- "+k, v1))
		case !ok1 && ok2:
			result = append(result, fmt.Sprintf("  %s: %v", "+ "+k, v2))
		case ok1 && ok2 && v1 != v2:
			result = append(result, fmt.Sprintf("  %s: %v", "- "+k, v1))
			result = append(result, fmt.Sprintf("  %s: %v", "+ "+k, v2))
		case ok1 && ok2 && v1 == v2:
			result = append(result, fmt.Sprintf("  %s: %v", "  "+k, v1))
		}
	}
	return result
}

func formatResults(res []string) string {
	var b strings.Builder
	b.WriteString("{\n")
	if len(res) > 0 {
		b.WriteString(strings.Join(res, "\n"))
		b.WriteByte('\n')
	}
	b.WriteString("}")
	return b.String()
}

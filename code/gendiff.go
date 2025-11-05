package code

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Parsed struct {
	Ext  string
	Data map[string]any
}

func parseFiles(paths ...string) ([]Parsed, error) {
	res := make([]Parsed, 0, len(paths))
	for _, p := range paths {
		parsed, err := parseFile(p)
		if err != nil {
			return nil, fmt.Errorf("parse %q: %w", p, err)
		}
		res = append(res, parsed)
	}
	return res, nil
}

func GenDiff(path1, path2 string) string {
	parsed, err := parseFiles(path1, path2)
	if err != nil {
		log.Fatalf("parse %q: %v", path1, err)
	}
	file1 := parsed[0].Data
	file2 := parsed[1].Data

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
			result = append(result, fmt.Sprintf("  %s: %s", "- "+k, v1))
		case !ok1 && ok2:
			result = append(result, fmt.Sprintf("  %s: %s", "+ "+k, v2))
		case ok1 && ok2 && v1 != v2:
			result = append(result, fmt.Sprintf("  %s: %s", "- "+k, v1))
			result = append(result, fmt.Sprintf("  %s: %s", "+ "+k, v2))
		case ok1 && ok2 && v1 == v2:
			result = append(result, fmt.Sprintf("  %s: %s", "  "+k, v1))
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

func parseFile(path string) (Parsed, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return Parsed{}, fmt.Errorf("abs(%q): %w", path, err)
	}

	data, err := os.ReadFile(abs)
	if err != nil {
		return Parsed{}, fmt.Errorf("read %q: %w", abs, err)
	}

	ext := filepath.Ext(abs)

	var raw map[string]any
	switch ext {
	case ".json":
		dec := json.NewDecoder(bytes.NewReader(data))
		dec.UseNumber() // числа будут json.Number, а не float64
		if err := dec.Decode(&raw); err != nil {
			return Parsed{}, fmt.Errorf("json decode %q: %w", abs, err)
		}
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, &raw); err != nil {
			return Parsed{}, fmt.Errorf("yaml decode %q: %w", abs, err)
		}
		// yaml.v3 отдаёт числа как int/float64 — это ок для твоего getInt
	default:
		return Parsed{}, fmt.Errorf("unsupported file extension: %s", ext)
	}

	return Parsed{
		Ext:  ext,
		Data: raw,
	}, nil
}

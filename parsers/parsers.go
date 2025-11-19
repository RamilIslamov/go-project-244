package parsers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

func ParseFiles(paths ...string) ([]map[string]any, error) {
	res := make([]map[string]any, 0, len(paths))
	for _, p := range paths {
		parsed, err := parseFile(p)
		if err != nil {
			return nil, fmt.Errorf("parse %q: %w", p, err)
		}
		res = append(res, parsed)
	}
	return res, nil
}

func parseFile(path string) (map[string]any, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("abs(%q): %w", path, err)
	}

	dst := map[string]any{}
	data, err := os.ReadFile(abs)
	if err != nil {
		return nil, fmt.Errorf("read %q: %w", abs, err)
	}

	switch ext := filepath.Ext(abs); ext {
	case ".json":
		if err := parseJSON(dst, data, abs); err != nil {
			return nil, err
		}
		normalizeJSONNumbersAny(dst)
	case ".yaml", ".yml":
		if err := parseYAML(dst, data, abs); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported file extension: %s", ext)
	}

	return dst, nil
}

func parseJSON(dst map[string]any, data []byte, abs string) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()

	var tmp map[string]any
	if err := dec.Decode(&tmp); err != nil {
		return fmt.Errorf("json decode %q: %w", abs, err)
	}
	tmp = normalizeJSONNumbersAny(tmp).(map[string]any)
	deepMerge(dst, tmp)
	return nil
}

func parseYAML(dst map[string]any, data []byte, abs string) error {
	var tmp map[string]any
	if err := yaml.Unmarshal(data, &tmp); err != nil {
		return fmt.Errorf("yaml decode %q: %w", abs, err)
	}

	tmp = normalizeJSONNumbersAny(tmp).(map[string]any)
	deepMerge(dst, tmp)
	return nil
}

func deepMerge(dst, src map[string]any) {
	for k, sv := range src {
		if dv, ok := dst[k]; ok {
			dm, dok := dv.(map[string]any)
			sm, sok := sv.(map[string]any)
			if dok && sok {
				deepMerge(dm, sm)
				continue
			}
		}
		dst[k] = sv
	}
}

func normalizeJSONNumbersAny(v any) any {
	switch x := v.(type) {
	case json.Number:
		if i, err := x.Int64(); err == nil {
			return i
		}
		if f, err := x.Float64(); err == nil {
			return f
		}
		return v
	case map[string]any:
		for k, vv := range x {
			x[k] = normalizeJSONNumbersAny(vv)
		}
		return x
	case []any:
		for i, vv := range x {
			x[i] = normalizeJSONNumbersAny(vv)
		}
		return x
	default:
		return v
	}
}

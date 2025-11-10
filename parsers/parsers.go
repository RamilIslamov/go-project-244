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
		normalizeJSONNumbers(dst)
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
	mergeInto(dst, tmp)
	return nil
}

func parseYAML(dst map[string]any, data []byte, abs string) error {
	var tmp map[string]any
	if err := yaml.Unmarshal(data, &tmp); err != nil {
		return fmt.Errorf("yaml decode %q: %w", abs, err)
	}
	mergeInto(dst, tmp)
	return nil
}

func mergeInto(dst, src map[string]any) {
	for k, v := range src {
		dst[k] = v
	}
}

func normalizeJSONNumbers(m map[string]any) {
	for k, v := range m {
		switch vv := v.(type) {
		case json.Number:
			if i, err := vv.Int64(); err == nil {
				m[k] = i
				continue
			}
			if f, err := vv.Float64(); err == nil {
				m[k] = f
			}
		case map[string]any:
			normalizeJSONNumbers(vv)
		case []any:
			for i, e := range vv {
				switch ee := e.(type) {
				case json.Number:
					if ii, err := ee.Int64(); err == nil {
						vv[i] = ii
					} else if ff, err := ee.Float64(); err == nil {
						vv[i] = ff
					}
				case map[string]any:
					normalizeJSONNumbers(ee)
				}
			}
		}
	}
}

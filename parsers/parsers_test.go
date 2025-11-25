package parsers

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func getInt(v any) (int64, bool) {
	switch x := v.(type) {
	case json.Number:
		n, err := x.Int64()
		return n, err == nil
	case float64:
		return int64(x), x == float64(int64(x))
	case int:
		return int64(x), true
	case int64:
		return x, true
	default:
		return 0, false
	}
}

func TestParseFile_JSONDecodeError(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	p := filepath.Join(dir, "bad.json")

	if err := os.WriteFile(p, []byte(`{"a":`), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := parseFile(p); err == nil {
		t.Fatalf("expected json decode error, got nil")
	}
}

func TestParseFile_YAMLDecodeError(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	p := filepath.Join(dir, "bad.yaml")

	if err := os.WriteFile(p, []byte(": bad"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := parseFile(p); err == nil {
		t.Fatalf("expected yaml decode error, got nil")
	}
}

func TestParseFile_AbsError(t *testing.T) {
	t.Parallel()

	bad := "invalid\x00path.json"
	if _, err := parseFile(bad); err == nil {
		t.Fatalf("expected error for path with NUL byte, got nil")
	}
}

func TestParseFile_RelativePath(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "a", "b"), 0o755); err != nil {
		t.Fatal(err)
	}
	full := filepath.Join(dir, "a", "b", "config.json")
	wantMap := map[string]any{"ok": true}
	if err := os.WriteFile(full, []byte(`{"ok": true}`), 0o644); err != nil {
		t.Fatal(err)
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := os.Chdir(wd); err != nil {
			t.Fatalf("restore cwd: %v", err)
		}
	})

	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	rel := filepath.Join("a", "b", "..", "b", "config.json")
	got, err := parseFile(rel)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, wantMap) {
		t.Fatalf("data mismatch:\n got: %#v\nwant: %#v", got, wantMap)
	}
}

func TestParseFile_AbsolutePath(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	full := filepath.Join(dir, "config.json")

	if err := os.WriteFile(full, []byte(`{"ok": 42}`), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := parseFile(full)
	if err != nil {
		t.Fatal(err)
	}

	val, exists := got["ok"]
	if !exists {
		t.Fatalf(`key "ok" not found`)
	}
	n, ok := getInt(val)
	if !ok || n != 42 {
		t.Errorf(`"ok" = %#v (type %T), want 42`, val, val)
	}
}

func TestParseFile_NotExists(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	_, err := parseFile(filepath.Join(dir, "missing.json"))
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("expected fs.ErrNotExist, got %v", err)
	}
}

func TestParseFiles_OK(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	p1 := filepath.Join(dir, "a.json")
	p2 := filepath.Join(dir, "b.json")
	f1 := []byte(`{"a": 1, "b": true}`)
	f2 := []byte(`{"a": 2, "c": "x"}`)

	if err := os.WriteFile(p1, f1, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p2, f2, 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := ParseFiles(p1, p2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}

	aOne, ok := getInt(got[0]["a"])
	if !ok || aOne != 1 {
		t.Fatalf("file #1 no key 'a' or invalid val on key 'a'")
	}

	bOne, ok := got[0]["b"]
	if !ok || bOne != true {
		t.Fatalf("file #1 no key 'b' or invalid val on key 'b'")
	}

	aTwo, ok := getInt(got[1]["a"])
	if !ok || aTwo != 2 {
		t.Fatalf("file #2 no key 'a' or invalid val on key 'a'")
	}
	cTwo, ok := got[1]["c"]
	if !ok || cTwo != "x" {
		t.Fatalf("file #2 no key 'c' or invalid val on key 'c'")
	}
}

func TestParseFiles_Error(t *testing.T) {
	t.Parallel()

	_, err := ParseFiles("/no/such/file.json")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestDeepMerge_NestedMaps(t *testing.T) {
	t.Parallel()

	dst := map[string]any{
		"common": map[string]any{
			"setting1": "Value 1",
			"setting2": 10,
		},
	}

	src := map[string]any{
		"common": map[string]any{
			"setting2": 20,
			"setting3": true,
		},
	}

	deepMerge(dst, src)

	common, ok := dst["common"].(map[string]any)
	if !ok {
		t.Fatalf(`"common" has wrong type: %T`, dst["common"])
	}

	if common["setting1"] != "Value 1" {
		t.Fatalf(`setting1 = %#v, want "Value 1"`, common["setting1"])
	}
	if n, ok := getInt(common["setting2"]); !ok || n != 20 {
		t.Fatalf(`setting2 = %#v, want 20`, common["setting2"])
	}
	if v, ok := common["setting3"]; !ok || v != true {
		t.Fatalf(`setting3 = %#v, want true`, common["setting3"])
	}
}

func TestParseFile_YAML_OK(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	p := filepath.Join(dir, "good.yaml")

	yamlData := []byte(`
a: 1
b: true
nested:
  c: 3
`)

	if err := os.WriteFile(p, yamlData, 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := parseFile(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// проверяем, что значения действительно распарсились
	if n, ok := getInt(got["a"]); !ok || n != 1 {
		t.Fatalf(`"a" = %#v, want 1`, got["a"])
	}

	if v, ok := got["b"]; !ok || v != true {
		t.Fatalf(`"b" = %#v, want true`, v)
	}

	nested, ok := got["nested"].(map[string]any)
	if !ok {
		t.Fatalf(`"nested" has wrong type: %T`, got["nested"])
	}
	if n, ok := getInt(nested["c"]); !ok || n != 3 {
		t.Fatalf(`nested["c"] = %#v, want 3`, nested["c"])
	}
}

func TestNormalizeJSONNumbersAny_IntNumber(t *testing.T) {
	t.Parallel()

	in := json.Number("42")

	out := normalizeJSONNumbersAny(in)

	v, ok := out.(int64)
	if !ok {
		t.Fatalf("want int64, got %T (%v)", out, out)
	}
	if v != 42 {
		t.Fatalf("value = %d, want 42", v)
	}
}

func TestNormalizeJSONNumbersAny_FloatNumber(t *testing.T) {
	t.Parallel()

	in := json.Number("3.14")

	out := normalizeJSONNumbersAny(in)

	f, ok := out.(float64)
	if !ok {
		t.Fatalf("want float64, got %T (%v)", out, out)
	}
	if f != 3.14 {
		t.Fatalf("value = %v, want 3.14", f)
	}
}

func TestNormalizeJSONNumbersAny_InvalidNumber(t *testing.T) {
	t.Parallel()

	// Int64 и Float64 оба вернут ошибку -> должно вернуться исходное значение
	in := json.Number("not-a-number")

	out := normalizeJSONNumbersAny(in)

	if out != in {
		t.Fatalf("expected original json.Number to be returned, got %#v", out)
	}
}

func TestNormalizeJSONNumbersAny_MapRecursion(t *testing.T) {
	t.Parallel()

	in := map[string]any{
		"a": json.Number("1"),
		"b": "text", // попадём в default-ветку
	}

	outAny := normalizeJSONNumbersAny(in)

	out, ok := outAny.(map[string]any)
	if !ok {
		t.Fatalf("want map[string]any, got %T", outAny)
	}

	a, ok := out["a"].(int64)
	if !ok || a != 1 {
		t.Fatalf(`"a" = %#v (%T), want int64(1)`, out["a"], out["a"])
	}
	if out["b"] != "text" {
		t.Fatalf(`"b" = %#v, want "text"`, out["b"])
	}
}

func TestNormalizeJSONNumbersAny_SliceRecursion(t *testing.T) {
	t.Parallel()

	in := []any{
		json.Number("2"),
		"ok", // default
	}

	outAny := normalizeJSONNumbersAny(in)

	out, ok := outAny.([]any)
	if !ok {
		t.Fatalf("want []any, got %T", outAny)
	}

	n, ok := out[0].(int64)
	if !ok || n != 2 {
		t.Fatalf("out[0] = %#v (%T), want int64(2)", out[0], out[0])
	}
	if out[1] != "ok" {
		t.Fatalf(`out[1] = %#v, want "ok"`, out[1])
	}
}

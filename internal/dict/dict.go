package dict

import (
	"fmt"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Field holds a column's human-readable name and an optional note.
type Field struct {
	Name string `yaml:"name"`
	Note string `yaml:"note"`
}

// TableDict holds the human-readable table name and its field map.
type TableDict struct {
	Name   string           `yaml:"name"`
	Fields map[string]Field `yaml:"fields"`
}

// Dict is the parsed dictionary. Codes and Relations are map fields; use
// LookupCodes / LookupRelations to access them (Go forbids a method sharing a
// name with a struct field).
type Dict struct {
	Tables    map[string]TableDict         `yaml:"tables"`
	Codes     map[string]map[string]string `yaml:"codes"`
	Relations map[string][]string          `yaml:"relations"`

	globalFields map[string]string
}

// New parses YAML bytes into a Dict and builds the global field index.
func New(data []byte) (*Dict, error) {
	var d Dict
	if err := yaml.Unmarshal(data, &d); err != nil {
		return nil, fmt.Errorf("字典解析失败: %w", err)
	}
	d.buildGlobalFields()
	return &d, nil
}

// buildGlobalFields builds a code→Chinese-name index across all tables.
// Tables are processed in sorted order so collisions are deterministic
// (first table alphabetically wins).
func (d *Dict) buildGlobalFields() {
	d.globalFields = map[string]string{}
	names := make([]string, 0, len(d.Tables))
	for t := range d.Tables {
		names = append(names, t)
	}
	sort.Strings(names)
	for _, t := range names {
		for code, f := range d.Tables[t].Fields {
			if _, exists := d.globalFields[code]; !exists {
				d.globalFields[code] = f.Name
			}
		}
	}
}

// Table returns the TableDict for the given table name (case-insensitive).
func (d *Dict) Table(name string) (TableDict, bool) {
	t, ok := d.Tables[strings.ToUpper(name)]
	return t, ok
}

// FieldName does a global best-effort lookup of a field code (case-insensitive)
// and returns its Chinese name.
func (d *Dict) FieldName(code string) (string, bool) {
	n, ok := d.globalFields[strings.ToUpper(code)]
	return n, ok
}

// LookupCodes returns the code table for the given key (e.g. "COPTC.单别").
func (d *Dict) LookupCodes(key string) (map[string]string, bool) {
	c, ok := d.Codes[key]
	return c, ok
}

// LookupRelations returns the relation strings for the given table name
// (case-insensitive).
func (d *Dict) LookupRelations(table string) ([]string, bool) {
	r, ok := d.Relations[strings.ToUpper(table)]
	return r, ok
}

// SearchTables returns table codes whose Chinese name, field names, or field
// notes contain keyword (case-insensitive). Only dict-annotated tables match —
// it complements schema.SearchTables (which matches table codes / suffix types).
func (d *Dict) SearchTables(keyword string) []string {
	kw := strings.ToUpper(strings.TrimSpace(keyword))
	if kw == "" {
		return nil
	}
	var out []string
	for code, td := range d.Tables {
		if tableMatches(td, kw) {
			out = append(out, code)
		}
	}
	sort.Strings(out)
	return out
}

func tableMatches(td TableDict, kw string) bool {
	if strings.Contains(strings.ToUpper(td.Name), kw) {
		return true
	}
	for _, f := range td.Fields {
		if strings.Contains(strings.ToUpper(f.Name), kw) ||
			strings.Contains(strings.ToUpper(f.Note), kw) {
			return true
		}
	}
	return false
}

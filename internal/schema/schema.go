package schema

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

type Column struct {
	Name          string `json:"name"`
	Position      int    `json:"position"`
	DataType      string `json:"data_type"`
	Nullable      bool   `json:"nullable"`
	IsSystemField bool   `json:"is_system_field"`
}

type Table struct {
	RowCount   int      `json:"row_count"`
	Module     string   `json:"module"`
	SuffixType string   `json:"suffix_type"`
	Columns    []Column `json:"columns"`
}

type Module struct {
	Code       string   `json:"-"`
	Name       string   `json:"name"`
	TableCount int      `json:"table_count"`
	Tables     []string `json:"tables"`
}

type Schema struct {
	Database string            `json:"database"`
	Modules  map[string]Module `json:"modules"`
	Tables   map[string]Table  `json:"tables"`
}

func New(data []byte) (*Schema, error) {
	var s Schema
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("schema 解析失败: %w", err)
	}
	return &s, nil
}

func (s *Schema) Table(name string) (Table, bool) {
	t, ok := s.Tables[strings.ToUpper(name)]
	return t, ok
}

func (s *Schema) TablesByModule(module string) []string {
	module = strings.ToUpper(module)
	var out []string
	for name, t := range s.Tables {
		if module == "" || strings.ToUpper(t.Module) == module {
			out = append(out, name)
		}
	}
	sort.Strings(out)
	return out
}

func (s *Schema) SearchTables(keyword string) []string {
	kw := strings.ToUpper(keyword)
	var out []string
	for name, t := range s.Tables {
		if strings.Contains(strings.ToUpper(name), kw) ||
			strings.Contains(strings.ToUpper(t.SuffixType), kw) {
			out = append(out, name)
		}
	}
	sort.Strings(out)
	return out
}

func (s *Schema) ModuleList() []Module {
	var out []Module
	for code, m := range s.Modules {
		m.Code = code
		out = append(out, m)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Code < out[j].Code })
	return out
}

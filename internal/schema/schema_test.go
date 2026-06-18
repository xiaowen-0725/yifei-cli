package schema

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const sample = `{
  "database": "YDSTEST",
  "modules": {"COP": {"name": "销售管理", "table_count": 1, "tables": ["COPTC"]}},
  "tables": {
    "COPTC": {"row_count": 7032, "module": "COP", "suffix_type": "单据单头",
      "columns": [
        {"name": "TC001", "position": 8, "data_type": "char(4)", "nullable": false, "is_system_field": false},
        {"name": "TC004", "position": 11, "data_type": "char(10)", "nullable": true, "is_system_field": false}
      ]}
  }
}`

func TestTableLookupCaseInsensitive(t *testing.T) {
	s, err := New([]byte(sample))
	require.NoError(t, err)
	tbl, ok := s.Table("coptc")
	require.True(t, ok)
	require.Equal(t, 7032, tbl.RowCount)
	require.Len(t, tbl.Columns, 2)
}

func TestTablesByModule(t *testing.T) {
	s, _ := New([]byte(sample))
	require.Equal(t, []string{"COPTC"}, s.TablesByModule("COP"))
	require.Empty(t, s.TablesByModule("PUR"))
}

func TestSearchTables(t *testing.T) {
	s, _ := New([]byte(sample))
	require.Equal(t, []string{"COPTC"}, s.SearchTables("单头"))
	require.Equal(t, []string{"COPTC"}, s.SearchTables("copt"))
}

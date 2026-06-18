package dict

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const sample = `
tables:
  COPTC:
    name: 销售订单单头
    fields:
      TC003: {name: 订单日期, note: YYYYMMDD}
      TC004: {name: 客户代号, note: 关联 COPMA.MA001}
codes:
  COPTC.单别:
    "221": 销售订单
relations:
  COPTC:
    - "TC004 → COPMA.MA001 (客户)"
`

func TestTableAndFieldName(t *testing.T) {
	d, err := New([]byte(sample))
	require.NoError(t, err)
	td, ok := d.Table("coptc")
	require.True(t, ok)
	require.Equal(t, "销售订单单头", td.Name)

	name, ok := d.FieldName("TC003")
	require.True(t, ok)
	require.Equal(t, "订单日期", name)
}

func TestCodesAndRelations(t *testing.T) {
	d, _ := New([]byte(sample))
	codes, ok := d.LookupCodes("COPTC.单别")
	require.True(t, ok)
	require.Equal(t, "销售订单", codes["221"])

	rels, ok := d.LookupRelations("COPTC")
	require.True(t, ok)
	require.Len(t, rels, 1)
}

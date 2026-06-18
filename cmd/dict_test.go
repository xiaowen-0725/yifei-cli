package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xiaowen-0725/yifei-cli/internal/dict"
)

func dictDeps() Deps {
	d, _ := dict.New([]byte(`
tables:
  COPTC:
    name: 销售订单单头
    fields:
      TC004: {name: 客户代号, note: 关联 COPMA.MA001}
codes:
  COPTC.单别:
    "221": 销售订单
relations:
  COPTC:
    - "TC004 → COPMA.MA001 (客户)"
`))
	return Deps{Dict: d}
}

func TestDictTable(t *testing.T) {
	root := NewRootCmd(dictDeps())
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetArgs([]string{"dict", "table", "COPTC"})
	require.NoError(t, root.Execute())
	require.Contains(t, out.String(), "客户代号")
}

func TestDictRelations(t *testing.T) {
	root := NewRootCmd(dictDeps())
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetArgs([]string{"dict", "relations", "COPTC"})
	require.NoError(t, root.Execute())
	require.Contains(t, out.String(), "COPMA.MA001")
}

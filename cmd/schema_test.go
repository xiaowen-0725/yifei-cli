package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xiaowen-0725/yifei-cli/internal/dict"
	"github.com/xiaowen-0725/yifei-cli/internal/schema"
)

func schemaDeps() Deps {
	s, _ := schema.New([]byte(`{
	  "modules": {"COP": {"name":"销售管理","table_count":1,"tables":["COPTC"]}},
	  "tables": {"COPTC": {"row_count":7032,"module":"COP","suffix_type":"单据单头",
	    "columns":[{"name":"TC004","position":11,"data_type":"char(10)","nullable":true,"is_system_field":false}]}}
	}`))
	d, _ := dict.New([]byte("tables:\n  COPTC:\n    name: 销售订单单头\n    fields:\n      TC004: {name: 客户代号}\n"))
	return Deps{Schema: s, Dict: d}
}

func TestSchemaTable(t *testing.T) {
	root := NewRootCmd(schemaDeps())
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetArgs([]string{"schema", "table", "COPTC"})
	require.NoError(t, root.Execute())
	require.Contains(t, out.String(), "TC004")
	require.Contains(t, out.String(), "客户代号")
}

func TestSchemaSearch(t *testing.T) {
	root := NewRootCmd(schemaDeps())
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetArgs([]string{"schema", "search", "单头"})
	require.NoError(t, root.Execute())
	require.Contains(t, out.String(), "COPTC")
}

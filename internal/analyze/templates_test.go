package analyze

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOrderCountByCustomerSQL(t *testing.T) {
	tpl, ok := Get("order-count-by-customer")
	require.True(t, ok)
	sql, err := tpl.Build(map[string]string{"year": "2022"})
	require.NoError(t, err)
	require.Contains(t, sql, "COPTC")
	require.Contains(t, sql, "COPMA")
	require.Contains(t, sql, "'20220101'")
	require.True(t, strings.HasPrefix(strings.TrimSpace(strings.ToUpper(sql)), "SELECT"))
}

func TestInventoryMovesTopDefaultTop(t *testing.T) {
	tpl, _ := Get("inventory-moves-top-items")
	sql, err := tpl.Build(map[string]string{})
	require.NoError(t, err)
	require.Contains(t, sql, "TOP 20") // default
	require.Contains(t, sql, "INVLA")
}

func TestUnknownTemplate(t *testing.T) {
	_, ok := Get("does-not-exist")
	require.False(t, ok)
}

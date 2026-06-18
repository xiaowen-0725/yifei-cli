//go:build integration

package db

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Run with: go test -tags integration ./internal/db/ (requires the YDSTEST container up).
func TestQueryLive(t *testing.T) {
	dsn := "sqlserver://sa:YourStrong%40Passw0rd@localhost:1433?database=YDSTEST&encrypt=disable"
	d, err := Open(dsn)
	require.NoError(t, err)
	defer d.Close()

	cols, rows, err := d.Query("SELECT TOP 3 MB001 FROM INVMB", 2)
	require.NoError(t, err)
	require.Equal(t, []string{"MB001"}, cols)
	require.Len(t, rows, 2) // limit enforced
}

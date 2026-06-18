package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAnalyzeList(t *testing.T) {
	root := NewRootCmd(testDeps(&fakeDB{}))
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetArgs([]string{"analyze", "list"})
	require.NoError(t, root.Execute())
	require.Contains(t, out.String(), "order-count-by-customer")
}

func TestAnalyzeRun(t *testing.T) {
	cfg := writeConfig(t)
	fake := &fakeDB{}
	root := NewRootCmd(testDeps(fake))
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetArgs([]string{"--config", cfg, "--format", "json",
		"analyze", "inventory-moves-top-items"})
	require.NoError(t, root.Execute())
	require.Contains(t, out.String(), "MB001") // fakeDB returns MB001 column
}

// TestAnalyzeMaliciousYear ensures a crafted --year value is rejected before
// reaching the database. The fakeDB queried field is checked to confirm no
// query ever reached it.
func TestAnalyzeMaliciousYear(t *testing.T) {
	cfg := writeConfig(t)
	fake := &fakeDB{}
	root := NewRootCmd(testDeps(fake))
	root.SetArgs([]string{"--config", cfg,
		"analyze", "order-count-by-customer",
		"--year", "2022'; DROP TABLE COPTC--"})
	err := root.Execute()
	require.Error(t, err, "malicious --year must be rejected")
	require.Equal(t, 0, fake.lastLimit, "DB must not have been queried")
}

// TestAnalyzeMaliciousTop ensures a crafted --top value is rejected before
// reaching the database.
func TestAnalyzeMaliciousTop(t *testing.T) {
	cfg := writeConfig(t)
	fake := &fakeDB{}
	root := NewRootCmd(testDeps(fake))
	root.SetArgs([]string{"--config", cfg,
		"analyze", "inventory-moves-top-items",
		"--top", "1; DELETE FROM INVLA--"})
	err := root.Execute()
	require.Error(t, err, "malicious --top must be rejected")
	require.Equal(t, 0, fake.lastLimit, "DB must not have been queried")
}

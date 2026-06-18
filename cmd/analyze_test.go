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

package cmd

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfigInitAndShowMasksPassword(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")

	root := NewRootCmd(Deps{})
	root.SetArgs([]string{"--config", path, "config", "init",
		"--host", "localhost", "--port", "1433", "--user", "sa",
		"--password", "secret", "--database", "YDSTEST"})
	require.NoError(t, root.Execute())

	root2 := NewRootCmd(Deps{})
	var out bytes.Buffer
	root2.SetOut(&out)
	root2.SetArgs([]string{"--config", path, "config", "show"})
	require.NoError(t, root2.Execute())
	require.Contains(t, out.String(), "YDSTEST")
	require.NotContains(t, out.String(), "secret")
	require.Contains(t, out.String(), "****")
}

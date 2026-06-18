package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRootCmd_HelpRuns(t *testing.T) {
	root := NewRootCmd(Deps{})
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetArgs([]string{"--help"})
	err := root.Execute()
	require.NoError(t, err)
	require.Contains(t, out.String(), "yifei")
}

func TestRootCmd_HasGlobalFormatFlag(t *testing.T) {
	root := NewRootCmd(Deps{})
	f := root.PersistentFlags().Lookup("format")
	require.NotNil(t, f)
	require.Equal(t, "table", f.DefValue)
}

package cmd

import (
	"bytes"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xiaowen-0725/yifei-cli/internal/dict"
)

type fakeDB struct{ lastLimit int }

func (f *fakeDB) Query(sql string, limit int) ([]string, [][]any, error) {
	f.lastLimit = limit
	return []string{"MB001"}, [][]any{{"11Z122"}}, nil
}
func (f *fakeDB) Close() error { return nil }

func testDeps(fake *fakeDB) Deps {
	d, _ := dict.New([]byte("tables: {}\ncodes: {}\nrelations: {}\n"))
	return Deps{
		Dict:   d,
		OpenDB: func(string) (Querier, error) { return fake, nil },
	}
}

func writeConfig(t *testing.T) string {
	path := filepath.Join(t.TempDir(), "config.yaml")
	root := NewRootCmd(Deps{})
	root.SetArgs([]string{"--config", path, "config", "init", "--password", "x"})
	require.NoError(t, root.Execute())
	return path
}

func TestQueryRejectsWrite(t *testing.T) {
	cfg := writeConfig(t)
	root := NewRootCmd(testDeps(&fakeDB{}))
	root.SetArgs([]string{"--config", cfg, "query", "DELETE FROM COPTC"})
	err := root.Execute()
	require.Error(t, err)
	require.Contains(t, err.Error(), "只读")
}

func TestQueryRunsAndAppliesDefaultLimit(t *testing.T) {
	cfg := writeConfig(t)
	fake := &fakeDB{}
	root := NewRootCmd(testDeps(fake))
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetArgs([]string{"--config", cfg, "--format", "json", "query", "SELECT MB001 FROM INVMB"})
	require.NoError(t, root.Execute())
	require.Equal(t, 1000, fake.lastLimit)
	require.Contains(t, out.String(), "11Z122")
	_ = fmt.Sprint
}

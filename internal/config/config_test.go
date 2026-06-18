package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSaveLoadRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	in := &Config{Host: "localhost", Port: 1433, User: "sa", Password: "p@ss", Database: "YDSTEST"}
	require.NoError(t, Save(path, in))

	info, err := os.Stat(path)
	require.NoError(t, err)
	require.Equal(t, os.FileMode(0o600), info.Mode().Perm())

	out, err := Load(path)
	require.NoError(t, err)
	require.Equal(t, in.User, out.User)
	require.Equal(t, 1433, out.Port)
}

func TestLoadMissingFileGivesGuidance(t *testing.T) {
	_, err := Load(filepath.Join(t.TempDir(), "nope.yaml"))
	require.Error(t, err)
	require.Contains(t, err.Error(), "config init")
}

func TestPasswordEnvOverride(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	require.NoError(t, Save(path, &Config{Host: "h", Port: 1, User: "u", Password: "filepw", Database: "d"}))
	t.Setenv("YIFEI_PASSWORD", "envpw")
	out, err := Load(path)
	require.NoError(t, err)
	require.Equal(t, "envpw", out.Password)
}

func TestDSN(t *testing.T) {
	c := &Config{Host: "localhost", Port: 1433, User: "sa", Password: "p@ss w", Database: "YDSTEST"}
	dsn := c.DSN()
	require.Contains(t, dsn, "sqlserver://")
	require.Contains(t, dsn, "database=YDSTEST")
	require.Contains(t, dsn, "encrypt=disable")
}

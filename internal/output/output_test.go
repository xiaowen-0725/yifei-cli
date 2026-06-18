package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var cols = []string{"MB001", "MB002"}
var rows = [][]any{{"11Z122", "助拍器"}, {"16WSMZ", "云台"}}

func tr(code string) (string, bool) {
	m := map[string]string{"MB001": "品号", "MB002": "品名"}
	n, ok := m[code]
	return n, ok
}

func TestRenderJSON(t *testing.T) {
	var b bytes.Buffer
	require.NoError(t, Render(&b, "json", cols, rows, nil))
	require.Contains(t, b.String(), `"MB001": "11Z122"`)
}

func TestRenderCSVTranslated(t *testing.T) {
	var b bytes.Buffer
	require.NoError(t, Render(&b, "csv", cols, rows, tr))
	require.True(t, strings.HasPrefix(b.String(), "品号,品名"))
}

func TestRenderTable(t *testing.T) {
	var b bytes.Buffer
	require.NoError(t, Render(&b, "table", cols, rows, nil))
	require.Contains(t, b.String(), "MB001")
	require.Contains(t, b.String(), "助拍器")
}

func TestRenderUnknownFormat(t *testing.T) {
	require.Error(t, Render(&bytes.Buffer{}, "xml", cols, rows, nil))
}

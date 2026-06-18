package output

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
)

func headers(cols []string, translate func(string) (string, bool)) []string {
	if translate == nil {
		return cols
	}
	out := make([]string, len(cols))
	for i, c := range cols {
		if n, ok := translate(c); ok {
			out[i] = n
		} else {
			out[i] = c
		}
	}
	return out
}

func cell(v any) string {
	if v == nil {
		return ""
	}
	if b, ok := v.([]byte); ok {
		return string(b)
	}
	return strings.TrimSpace(fmt.Sprintf("%v", v))
}

func Render(w io.Writer, format string, cols []string, rows [][]any, translate func(string) (string, bool)) error {
	h := headers(cols, translate)
	switch format {
	case "json":
		objs := make([]map[string]string, 0, len(rows))
		for _, r := range rows {
			m := map[string]string{}
			for i, c := range h {
				if i < len(r) {
					m[c] = cell(r[i])
				}
			}
			objs = append(objs, m)
		}
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(objs)
	case "csv":
		cw := csv.NewWriter(w)
		if err := cw.Write(h); err != nil {
			return err
		}
		for _, r := range rows {
			rec := make([]string, len(h))
			for i := range h {
				if i < len(r) {
					rec[i] = cell(r[i])
				}
			}
			if err := cw.Write(rec); err != nil {
				return err
			}
		}
		cw.Flush()
		return cw.Error()
	case "table":
		tw := tabwriter.NewWriter(w, 0, 2, 2, ' ', 0)
		fmt.Fprintln(tw, strings.Join(h, "\t"))
		for _, r := range rows {
			rec := make([]string, len(h))
			for i := range h {
				if i < len(r) {
					rec[i] = cell(r[i])
				}
			}
			fmt.Fprintln(tw, strings.Join(rec, "\t"))
		}
		return tw.Flush()
	default:
		return fmt.Errorf("不支持的输出格式: %q (可选 table|json|csv)", format)
	}
}

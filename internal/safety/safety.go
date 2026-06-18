package safety

import (
	"fmt"
	"regexp"
	"strings"
)

var forbidden = []string{
	"INSERT", "UPDATE", "DELETE", "DROP", "ALTER", "TRUNCATE",
	"CREATE", "EXEC", "EXECUTE", "MERGE", "GRANT", "REVOKE", "INTO",
}

var (
	blockComment = regexp.MustCompile(`(?s)/\*.*?\*/`)
	lineComment  = regexp.MustCompile(`--[^\n]*`)
)

func stripComments(sql string) string {
	sql = blockComment.ReplaceAllString(sql, " ")
	sql = lineComment.ReplaceAllString(sql, " ")
	return sql
}

// ValidateReadOnly enforces a single read-only SELECT/WITH statement.
func ValidateReadOnly(sql string) error {
	clean := strings.TrimSpace(stripComments(sql))
	if clean == "" {
		return fmt.Errorf("SQL 为空")
	}
	trimmed := strings.TrimRight(clean, "; \t\r\n")
	if strings.Contains(trimmed, ";") {
		return fmt.Errorf("只允许单条语句,检测到多条语句(分号)")
	}
	upper := strings.ToUpper(trimmed)
	if !strings.HasPrefix(upper, "SELECT") && !strings.HasPrefix(upper, "WITH") {
		return fmt.Errorf("只允许 SELECT/WITH 只读查询")
	}
	for _, kw := range forbidden {
		if regexp.MustCompile(`(?i)\b` + kw + `\b`).MatchString(trimmed) {
			return fmt.Errorf("检测到禁止的关键字 %q,本工具仅支持只读查询", kw)
		}
	}
	return nil
}

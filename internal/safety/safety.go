package safety

import (
	"fmt"
	"regexp"
	"strings"
)

var forbiddenKeywords = []string{
	"INSERT", "UPDATE", "DELETE", "DROP", "ALTER", "TRUNCATE",
	"CREATE", "EXEC", "EXECUTE", "MERGE", "GRANT", "REVOKE", "INTO",
}

var (
	blockComment = regexp.MustCompile(`(?s)/\*.*?\*/`)
	lineComment  = regexp.MustCompile(`--[^\n]*`)

	// forbiddenRes are compiled once at init time (not per-call) for efficiency.
	forbiddenRes []*regexp.Regexp
)

func init() {
	for _, kw := range forbiddenKeywords {
		forbiddenRes = append(forbiddenRes, regexp.MustCompile(`(?i)\b`+kw+`\b`))
	}
}

func stripComments(sql string) string {
	sql = blockComment.ReplaceAllString(sql, " ")
	sql = lineComment.ReplaceAllString(sql, " ")
	return sql
}

// ValidateReadOnly enforces a single read-only SELECT/WITH statement.
//
// Note: forbidden keywords and semicolons inside string literals are
// intentionally rejected — this is conservative-by-design. The validator
// does not parse SQL; it performs a keyword scan on the stripped text.
// False positives (e.g., a column alias containing "INTO") are accepted
// as the cost of a simple, auditable read-only guarantee.
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
	for i, re := range forbiddenRes {
		if re.MatchString(trimmed) {
			return fmt.Errorf("检测到禁止的关键字 %q,本工具仅支持只读查询", forbiddenKeywords[i])
		}
	}
	return nil
}

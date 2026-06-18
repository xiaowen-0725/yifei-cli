package analyze

import "fmt"

type Template struct {
	Name  string
	Desc  string
	Build func(params map[string]string) (string, error)
}

func param(p map[string]string, key, def string) string {
	if v, ok := p[key]; ok && v != "" {
		return v
	}
	return def
}

var templates = []Template{
	{
		Name: "order-count-by-customer",
		Desc: "按客户统计销售订单数 (参数: year)",
		Build: func(p map[string]string) (string, error) {
			year := param(p, "year", "")
			where := ""
			if year != "" {
				where = fmt.Sprintf("WHERE h.TC003 BETWEEN '%s0101' AND '%s1231'", year, year)
			}
			return fmt.Sprintf(`SELECT h.TC004 AS 客户代号, c.MA002 AS 客户简称, COUNT(*) AS 订单数
FROM COPTC h
LEFT JOIN COPMA c ON h.TC004 = c.MA001
%s
GROUP BY h.TC004, c.MA002
ORDER BY 订单数 DESC`, where), nil
		},
	},
	{
		Name: "workorder-count-by-month",
		Desc: "按月统计工单数 (参数: year)",
		Build: func(p map[string]string) (string, error) {
			year := param(p, "year", "")
			where := ""
			if year != "" {
				where = fmt.Sprintf("WHERE TA003 BETWEEN '%s0101' AND '%s1231'", year, year)
			}
			return fmt.Sprintf(`SELECT SUBSTRING(TA003,1,6) AS 年月, COUNT(*) AS 工单数
FROM MOCTA
%s
GROUP BY SUBSTRING(TA003,1,6)
ORDER BY 年月`, where), nil
		},
	},
	{
		Name: "inventory-moves-top-items",
		Desc: "库存异动次数 TOP 品号 (参数: top)",
		Build: func(p map[string]string) (string, error) {
			top := param(p, "top", "20")
			return fmt.Sprintf(`SELECT TOP %s LA001 AS 品号, COUNT(*) AS 异动次数
FROM INVLA
GROUP BY LA001
ORDER BY 异动次数 DESC`, top), nil
		},
	},
}

func All() []Template { return templates }

func Get(name string) (Template, bool) {
	for _, t := range templates {
		if t.Name == name {
			return t, true
		}
	}
	return Template{}, false
}

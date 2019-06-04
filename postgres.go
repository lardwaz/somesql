package somesql

import (
	"fmt"
	"regexp"
	"strings"
)

func cleanStatement(sql string) string {
	space := regexp.MustCompile(`\s+`)
	return strings.TrimSpace(space.ReplaceAllString(sql, " "))
}

func processPlaceholders(sql string) string {
	var i int
	for _, r := range sql {
		if r == '?' {
			i++
			sql = strings.Replace(sql, "?", fmt.Sprintf("$%d", i), 1)
		}
	}

	return sql
}

func processConditions(conds []Condition) (string, []interface{}) {
	var conditions string
	values := make([]interface{}, 0)
	for i, cond := range conds {
		if i != 0 {
			switch cond.ConditionType() {
			case AndCondition:
				conditions += ` AND `
			case OrCondition:
				conditions += ` OR `
			default:
				continue
			}
		}

		c, v := cond.AsSQL()
		conditions += c
		values = append(values, v...)
	}

	return conditions, values
}

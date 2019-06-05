package somesql

import (
	"regexp"
	"strconv"
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
			sql = strings.Replace(sql, "?", "$"+strconv.Itoa(i), 1)
		}
	}

	return sql
}

func processConditions(conds []Condition) (string, []interface{}) {
	var (
		conditionsBuff strings.Builder
		values         []interface{}
	)

	for i, cond := range conds {
		if i != 0 {
			switch cond.ConditionType() {
			case AndCondition:
				conditionsBuff.WriteString(` AND `)
			case OrCondition:
				conditionsBuff.WriteString(` OR `)
			default:
				continue
			}
		}

		c, v := cond.AsSQL()
		conditionsBuff.WriteString(c)
		values = append(values, v...)
	}

	return conditionsBuff.String(), values
}

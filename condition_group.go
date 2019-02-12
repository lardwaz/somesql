package somesql

import (
	"fmt"
)

//ConditionGroup represents a group of condition (within same pair brackets)
type ConditionGroup struct {
	Type       uint8
	Conditions []Condition
}

//AndGroup creates a condition group in the format AND((condition1) OR (condition2) AND (condition3))
func AndGroup(conditions ...Condition) ConditionGroup {
	return ConditionGroup{
		Type:       AndCondition,
		Conditions: conditions,
	}
}

//OrGroup creates a condition group in the format OR((condition1) OR (condition2) AND (condition3))
func OrGroup(conditions ...Condition) ConditionGroup {
	return ConditionGroup{
		Type:       OrCondition,
		Conditions: conditions,
	}
}

//ConditionType to satisfy interface Condition
func (c ConditionGroup) ConditionType() uint8 {
	return c.Type
}

//AsSQL to satisfy interface Condition
func (c ConditionGroup) AsSQL(in ...bool) (string, []interface{}) {
	var (
		sqls   string
		values []interface{}
	)
	for i, cond := range c.Conditions {
		sql, val := cond.AsSQL()
		values = append(values, val...)

		if i == 0 {
			sqls = sql
			continue
		}
		if cond.ConditionType() == AndCondition {
			sqls = fmt.Sprintf("%s AND %s", sqls, sql)
		} else {
			sqls = fmt.Sprintf("%s OR %s", sqls, sql)
		}
	}

	if len(c.Conditions) > 0 {
		sqls = fmt.Sprintf("(%s)", sqls)
	}

	return sqls, values
}

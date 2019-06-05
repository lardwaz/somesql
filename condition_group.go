package somesql

import "strings"

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
		sqlStr  string
		sqlBuff strings.Builder
		values  []interface{}
	)
	for i, cond := range c.Conditions {
		sql, val := cond.AsSQL()
		values = append(values, val...)

		if i == 0 {
			sqlBuff.WriteString(sql)
		} else if cond.ConditionType() == AndCondition {
			sqlBuff.WriteString(" AND " + sql)
		} else {
			sqlBuff.WriteString(" OR " + sql)
		}
	}

	if sqlBuff.Len() > 0 {
		sqlStr = "(" + sqlBuff.String() + ")"
	}

	return sqlStr, values
}

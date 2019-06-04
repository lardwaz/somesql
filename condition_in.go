package somesql

import (
	"fmt"
	"strings"
)

var (
	andIn    = andOrIn(AndCondition, "IN")
	orIn     = andOrIn(OrCondition, "IN")
	andNotIn = andOrIn(AndCondition, "NOT IN")
	orNotIn  = andOrIn(OrCondition, "NOT IN")
)

// ConditionIn represents a condition in the format IN(?,?,?) / NOT IN(?,?,?)
type ConditionIn struct {
	Type          uint8
	Field         string
	FieldFunction string
	Relations     bool
	Operator      string
	Values        interface{}
	Lang          string
}

// andOrIn is a factory function
// it generates And + Or functions which are identical except for the conditionType.
// one instance of each at runtime
func andOrIn(conditionType uint8, operator string) func(lang, field string, values interface{}, isRel bool, funcs ...string) ConditionIn {
	return func(lang, field string, values interface{}, isRel bool, funcs ...string) ConditionIn {
		fieldFunction, _ := getFieldValueFunctions(funcs)

		return ConditionIn{
			Type:          conditionType,
			Field:         field,
			FieldFunction: fieldFunction,
			Relations:     isRel,
			Operator:      operator,
			Values:        values,
			Lang:          lang,
		}
	}
}

// AndIn returns a condition in the format IN(?,?,?) adjoined with AND
func AndIn(lang, field string, values interface{}, funcs ...string) ConditionIn {
	return andIn(lang, field, values, false, funcs...)
}

// OrIn returns a condition in the format IN(?,?,?) adjoined with OR
func OrIn(lang, field string, values interface{}, funcs ...string) ConditionIn {
	return orIn(lang, field, values, false, funcs...)
}

// AndNotIn returns a condition in the format NOT IN(?,?,?) adjoined with AND
func AndNotIn(lang, field string, values interface{}, funcs ...string) ConditionIn {
	return andNotIn(lang, field, values, false, funcs...)
}

// OrNotIn returns a condition in the format NOT IN(?,?,?) adjoined with OR
func OrNotIn(lang, field string, values interface{}, funcs ...string) ConditionIn {
	return orNotIn(lang, field, values, false, funcs...)
}

// AndRelIn returns a condition in the format IN(?,?,?) adjoined with AND
func AndRelIn(lang, field string, values interface{}, funcs ...string) ConditionIn {
	return andIn(lang, field, values, true, funcs...)
}

// OrRelIn returns a condition in the format IN(?,?,?) adjoined with OR
func OrRelIn(lang, field string, values interface{}, funcs ...string) ConditionIn {
	return orIn(lang, field, values, true, funcs...)
}

// AndRelNotIn returns a condition in the format NOT IN(?,?,?) adjoined with AND
func AndRelNotIn(lang, field string, values interface{}, funcs ...string) ConditionIn {
	return andNotIn(lang, field, values, true, funcs...)
}

// OrRelNotIn returns a condition in the format NOT IN(?,?,?) adjoined with OR
func OrRelNotIn(lang, field string, values interface{}, funcs ...string) ConditionIn {
	return orNotIn(lang, field, values, true, funcs...)
}

// ConditionType to satisfy interface Condition
func (c ConditionIn) ConditionType() uint8 {
	return c.Type
}

// AsSQL to satisfy interface Condition
func (c ConditionIn) AsSQL(in ...bool) (string, []interface{}) {
	var (
		lhs, rhs string
		vals     []interface{}
	)

	vals = expandValues(c.Values)

	if IsFieldMeta(c.Field) || IsFieldData(c.Field) {
		field := fmt.Sprintf(`"%s"`, c.Field)

		if c.FieldFunction == None {
			lhs = field
		} else {
			lhs = fmt.Sprintf("%s(%s)", c.FieldFunction, field)
		}

		if l := len(vals); l != 0 {
			rhs = " (" + strings.TrimSuffix(strings.Repeat("?,", l), ",") + ")"
		}
		return lhs + " " + c.Operator + rhs, vals

	} else { // fields of type JSONB
		if c.Relations {
			lhs = fmt.Sprintf(`("%s" @> `, FieldRelations)
		} else {
			lhs = fmt.Sprintf(`("%s" @> `, GetLangFieldData(c.Lang))
		}
		rhs = ")"

		var pathValue []string
		for range vals {
			pathValue = append(pathValue, fmt.Sprintf(`'{"%s":["?"]}'::JSONB`, c.Field))
		}

		if c.Operator == "NOT IN" {
			return "NOT" + lhs + strings.Join(pathValue, " OR ") + rhs, vals
		}

		return lhs + strings.Join(pathValue, " OR ") + rhs, vals
	}
}

package somesql

import "fmt"

var (
	and = andor(AndCondition)
	or  = andor(OrCondition)
)

// ConditionClause represents a single conditional clause
type ConditionClause struct {
	Type          uint8
	Field         string
	FieldFunction string
	Relations     bool
	Operator      string
	ValueFunction string
	Value         interface{}
	Lang          string
}

// andor is a factory function
// it generates And + Or functions which are identical except for the conditionType.
// one instance of each at runtime
func andor(conditionType uint8) func(lang, field, operator string, value interface{}, isRel bool, funcs ...string) ConditionClause {
	return func(lang, field, operator string, value interface{}, isRel bool, funcs ...string) ConditionClause {
		fieldFunction, valueFunction := getFieldValueFunctions(funcs)

		return ConditionClause{
			Type:          conditionType,
			Field:         field,
			FieldFunction: fieldFunction,
			Relations:     isRel,
			Operator:      operator,
			Value:         value,
			ValueFunction: valueFunction,
			Lang:          lang,
		}
	}
}

// And creates an AND conditional clause
// And("field", "=", "val", "FOO", "BAR") yields: AND FOO(field) = BAR(val)
func And(lang, field, operator string, value interface{}, funcs ...string) ConditionClause {
	return and(lang, field, operator, value, false, funcs...)
}

// Or creates an AND conditional clause
// Or("field", "=", "val", "FOO", "BAR") yields: OR FOO(field) = BAR(val)
func Or(lang, field, operator string, value interface{}, funcs ...string) ConditionClause {
	return or(lang, field, operator, value, false, funcs...)
}

// AndRel creates an AND conditional clause
// AndRel("field", "=", "val") yields: AND "relations" @> '{"field":[?]}'::JSONB
func AndRel(lang, field, operator string, value interface{}, funcs ...string) ConditionClause {
	return and(lang, field, operator, value, true, funcs...)
}

// OrRel creates an AND conditional clause
// OrRel("field", "=", "val") yields: OR "relations" @> '{"field":[?]}'::JSONB
func OrRel(lang, field, operator string, value interface{}, funcs ...string) ConditionClause {
	return or(lang, field, operator, value, true, funcs...)
}

// ConditionType to satisfy interface Condition
func (c ConditionClause) ConditionType() uint8 {
	return c.Type
}

// AsSQL to satisfy interface Condition
func (c ConditionClause) AsSQL(in ...bool) (string, []interface{}) {
	var (
		lhs, rhs, field string
	)

	if IsFieldMeta(c.Field) || IsFieldData(c.Field) {
		field = fmt.Sprintf(`"%s"`, c.Field)
	} else if c.Relations {
		field = fmt.Sprintf(`"%s"`, FieldRelations)
	} else {
		field = fmt.Sprintf(`"%s"->>'%s'`, GetFieldData(c.Lang), c.Field)
	}

	if c.FieldFunction == None || c.Relations {
		lhs = field
	} else {
		lhs = fmt.Sprintf("%s(%s)", c.FieldFunction, field)
	}

	switch c.Value.(type) {
	case bool:
		lhs = fmt.Sprintf("(%s)::BOOLEAN", lhs)
	}

	if c.Relations {
		c.Operator = " @> "
		rhs = fmt.Sprintf(`'{"%s":?}'::JSONB`, c.Field)
	} else if c.ValueFunction == None {
		rhs = "?"
	} else {
		rhs = fmt.Sprintf("%s(?)", c.ValueFunction)
	}

	if c.Relations {
		return "(" + lhs + c.Operator + rhs + ")", expandValues(c.Value)
	}

	return lhs + c.Operator + rhs, expandValues(c.Value)
}

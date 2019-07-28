package somesql

import (
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
	Operator      string
	Values        interface{}
	Lang          string
}

// andOrIn is a factory function
// it generates And + Or functions which are identical except for the conditionType.
// one instance of each at runtime
func andOrIn(conditionType uint8, operator string) func(lang, field string, values interface{}, funcs ...string) ConditionIn {
	return func(lang, field string, values interface{}, funcs ...string) ConditionIn {
		fieldFunction, _ := getFieldValueFunctions(funcs)

		return ConditionIn{
			Type:          conditionType,
			Field:         field,
			FieldFunction: fieldFunction,
			Operator:      operator,
			Values:        values,
			Lang:          lang,
		}
	}
}

// AndIn returns a condition in the format IN(?,?,?) adjoined with AND
func AndIn(lang, field string, values interface{}, funcs ...string) ConditionIn {
	return andIn(lang, field, values, funcs...)
}

// OrIn returns a condition in the format IN(?,?,?) adjoined with OR
func OrIn(lang, field string, values interface{}, funcs ...string) ConditionIn {
	return orIn(lang, field, values, funcs...)
}

// AndNotIn returns a condition in the format NOT IN(?,?,?) adjoined with AND
func AndNotIn(lang, field string, values interface{}, funcs ...string) ConditionIn {
	return andNotIn(lang, field, values, funcs...)
}

// OrNotIn returns a condition in the format NOT IN(?,?,?) adjoined with OR
func OrNotIn(lang, field string, values interface{}, funcs ...string) ConditionIn {
	return orNotIn(lang, field, values, funcs...)
}

// ConditionType to satisfy interface Condition
func (c ConditionIn) ConditionType() uint8 {
	return c.Type
}

// AsSQL to satisfy interface Condition
func (c ConditionIn) AsSQL(in ...bool) (string, []interface{}) {
	var (
		lhs, rhs, field         string
		isInnerData, isInnerRel bool
		vals                    []interface{}
		dataFieldLang           string

		rhsBuff strings.Builder
	)

	dataFieldLang = GetLangFieldData(c.Lang)

	vals, _ = expandValues(c.Values)

	if IsFieldMeta(c.Field) || IsFieldData(c.Field) || IsFieldRelations(c.Field) {
		field = `"` + c.Field + `"`

		if c.FieldFunction == None {
			lhs = field
		} else {
			lhs = c.FieldFunction + "(" + field + ")"
		}

		lhs = lhs + " " + c.Operator
	} else if innerField, ok := GetInnerField(FieldData, c.Field); ok {
		field = `"` + innerField + `"`
		lhs = `("` + dataFieldLang + `" @> `
		isInnerData = true
	} else if innerField, ok := GetInnerField(FieldRelations, c.Field); ok {
		field = `"` + innerField + `"`
		lhs = `("` + dataFieldLang + `" @> `
		isInnerRel = true
	}

	for range vals {
		if isInnerData || isInnerRel {
			rhsBuff.WriteString(`'{` + field + `:["?"]}'::JSONB OR `)
		} else {
			rhsBuff.WriteString(`?,`)
		}
	}

	if rhsBuff.Len() > 0 {
		if isInnerData || isInnerRel {
			rhs = rhsBuff.String()[:rhsBuff.Len()-4] + ")" // trim " OR "
		} else {
			rhs = " (" + rhsBuff.String()[:rhsBuff.Len()-1] + ")" // trim ","
		}
	}

	if c.Operator == "NOT IN" && (isInnerData || isInnerRel) {
		return "NOT" + lhs + rhs, vals
	}

	return lhs + rhs, vals
}

package somesql

var (
	and = andor(AndCondition)
	or  = andor(OrCondition)
)

// ConditionClause represents a single conditional clause
type ConditionClause struct {
	Type          uint8
	Field         string
	FieldFunction string
	Operator      string
	ValueFunction string
	Value         interface{}
	Lang          string
}

// andor is a factory function
// it generates And + Or functions which are identical except for the conditionType.
// one instance of each at runtime
func andor(conditionType uint8) func(lang, field, operator string, value interface{}, funcs ...string) ConditionClause {
	return func(lang, field, operator string, value interface{}, funcs ...string) ConditionClause {
		fieldFunction, valueFunction := getFieldValueFunctions(funcs)

		return ConditionClause{
			Type:          conditionType,
			Field:         field,
			FieldFunction: fieldFunction,
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
	return and(lang, field, operator, value, funcs...)
}

// Or creates an AND conditional clause
// Or("field", "=", "val", "FOO", "BAR") yields: OR FOO(field) = BAR(val)
func Or(lang, field, operator string, value interface{}, funcs ...string) ConditionClause {
	return or(lang, field, operator, value, funcs...)
}

// ConditionType to satisfy interface Condition
func (c ConditionClause) ConditionType() uint8 {
	return c.Type
}

// AsSQL to satisfy interface Condition
func (c ConditionClause) AsSQL(in ...bool) (string, []interface{}) {
	var (
		lhs, rhs, field string
		isInnerRel      bool
		dataFieldLang   = GetLangFieldData(c.Lang)
	)

	if IsFieldMeta(c.Field) || IsFieldData(c.Field) || IsFieldRelations(c.Field) {
		if IsFieldData(c.Field) {
			field = `"` + dataFieldLang + `"`
		} else {
			field = `"` + c.Field + `"`
		}
	} else if innerField, ok := GetInnerField(FieldData, c.Field); ok {
		field = `"` + dataFieldLang + `"->>'` + innerField + `'`
	} else if innerField, ok := GetInnerField(FieldRelations, c.Field); ok {
		field = `"` + FieldRelations + `"`
		c.Operator = "@>"
		rhs = `'{"` + innerField + `":?}'::JSONB`
		isInnerRel = true
	}

	if c.FieldFunction == None || isInnerRel {
		lhs = field
	} else {
		lhs = c.FieldFunction + "(" + field + ")"
	}

	switch c.Value.(type) {
	case bool:
		lhs = "(" + lhs + ")::BOOLEAN"
	}

	if isInnerRel {
		// Do nothing
	} else if c.ValueFunction == None {
		rhs = "?"
	} else {
		rhs = c.ValueFunction + "(?)"
	}

	vals, _ := expandValues(c.Value)

	if isInnerRel {
		return "(" + lhs + " " + c.Operator + " " + rhs + ")", vals
	}

	return lhs + " " + c.Operator + " " + rhs, vals
}

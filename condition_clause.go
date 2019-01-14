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
}

// andor is a factory function
// it generates And + Or functions which are identical except for the conditionType.
// one instance of each at runtime
func andor(conditionType uint8) func(field, operator string, value interface{}, funcs ...string) ConditionClause {
	return func(field, operator string, value interface{}, funcs ...string) ConditionClause {
		var fieldFunction, valueFunction string

		if l := len(funcs); l == 2 {
			fieldFunction = funcs[0]
			valueFunction = funcs[1]
		} else if l == 1 {
			fieldFunction = funcs[0]
		}

		return ConditionClause{
			Type:          conditionType,
			Field:         field,
			FieldFunction: fieldFunction,
			Operator:      operator,
			Value:         value,
			ValueFunction: valueFunction,
		}
	}
}

// And creates an AND conditional clause
// And("myfield", "=", "val", "FOO", "BAR") yields: AND FOO(myfield) = BAR(val)
func And(field, operator string, value interface{}, funcs ...string) ConditionClause {
	return and(field, operator, value, funcs...)
}

// Or creates an AND conditional clause
// Or("myfield", "=", "val", "FOO", "BAR") yields: OR FOO(myfield) = BAR(val)
func Or(field, operator string, value interface{}, funcs ...string) ConditionClause {
	return or(field, operator, value, funcs...)
}

// ConditionType to satisfy interface Condition
func (c ConditionClause) ConditionType() uint8 {
	return c.Type
}

// AsSQL to satisfy interface Condition
func (c ConditionClause) AsSQL() (string, []interface{}) {
	var (
		lhs, rhs, field string
		// values          []interface{}
	)

	switch c.Field {
	case "id", "created_at", "updated_at", "status", "owner_id", "type", "slug", "data":
		field = c.Field
	default:
		field = `"data"->>'` + c.Field + `'`
	}

	if c.FieldFunction == None {
		lhs = field
	} else {
		lhs = c.FieldFunction + "(" + field + ")"
	}

	if c.ValueFunction != None {
		rhs = c.ValueFunction + "(?)"
	} else {
		rhs = "?"
	}

	return lhs + c.Operator + rhs, []interface{}{c.Value}
}

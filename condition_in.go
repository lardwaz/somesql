package somesql

import (
	"fmt"
	"strings"
)

var (
	andIn    = andorin(AndCondition, "IN")
	orIn     = andorin(OrCondition, "IN")
	andNotIn = andorin(AndCondition, "NOT IN")
	orNotIn  = andorin(OrCondition, "NOT IN")
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

// andorin is a factory function
// it generates And + Or functions which are identical except for the conditionType.
// one instance of each at runtime
func andorin(conditionType uint8, operator string) func(lang, field string, values interface{}, funcs ...string) ConditionIn {
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
func (c ConditionIn) AsSQL() (string, []interface{}) {
	var (
		field, lhs, rhs string
		vals            []interface{}
	)

	dataField := fmt.Sprintf("data_%s", c.Lang)

	switch c.Field {
	case "id", "created_at", "updated_at", "status", "owner_id", "type", "slug", dataField:
		field = c.Field
	default:
		field = fmt.Sprintf(`"%s"->>'%s'`, dataField, c.Field)
	}

	if c.FieldFunction == None {
		lhs = field
	} else {
		lhs = fmt.Sprintf("%s(%s)", c.FieldFunction, field)
	}

	rhs, vals = inValues(c.Values)
	return lhs + " " + c.Operator + rhs, vals
}

func inValues(val interface{}) (string, []interface{}) {
	var (
		values []interface{}
		query  string
	)

	if data, ok := val.([]string); ok {
		for d := range data {
			values = append(values, data[d])
		}
	} else if data, ok := val.([]bool); ok {
		for d := range data {
			values = append(values, data[d])
		}
	} else if data, ok := val.([]byte); ok {
		for d := range data {
			values = append(values, data[d])
		}
	} else if data, ok := val.([]complex64); ok {
		for d := range data {
			values = append(values, data[d])
		}
	} else if data, ok := val.([]complex128); ok {
		for d := range data {
			values = append(values, data[d])
		}
	} else if data, ok := val.([]float32); ok {
		for d := range data {
			values = append(values, data[d])
		}
	} else if data, ok := val.([]float64); ok {
		for d := range data {
			values = append(values, data[d])
		}
	} else if data, ok := val.([]int); ok {
		for d := range data {
			values = append(values, data[d])
		}
	} else if data, ok := val.([]int8); ok {
		for d := range data {
			values = append(values, data[d])
		}
	} else if data, ok := val.([]int16); ok {
		for d := range data {
			values = append(values, data[d])
		}
	} else if data, ok := val.([]int32); ok {
		for d := range data {
			values = append(values, data[d])
		}
	} else if data, ok := val.([]int64); ok {
		for d := range data {
			values = append(values, data[d])
		}
	} else if data, ok := val.([]rune); ok {
		for d := range data {
			values = append(values, data[d])
		}
	} else if data, ok := val.([]uint); ok {
		for d := range data {
			values = append(values, data[d])
		}
	} else if data, ok := val.([]uint8); ok {
		for d := range data {
			values = append(values, data[d])
		}
	} else if data, ok := val.([]uint16); ok {
		for d := range data {
			values = append(values, data[d])
		}
	} else if data, ok := val.([]uint32); ok {
		for d := range data {
			values = append(values, data[d])
		}
	} else if data, ok := val.([]uint64); ok {
		for d := range data {
			values = append(values, data[d])
		}
	} else if data, ok := val.([]uintptr); ok {
		for d := range data {
			values = append(values, data[d])
		}
	} else {
		return "", nil
	}

	if l := len(values); l != 0 {
		query = " (" + strings.TrimSuffix(strings.Repeat("?,", l), ",") + ")"
	}

	return query, values
}

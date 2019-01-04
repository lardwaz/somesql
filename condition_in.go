package somesql

import "strings"

type inCondition struct {
	Type uint8
}

// AndIn returns a condition in the format IN(?,?,?) adjoined with AND
func AndIn(field, operator string, value interface{}, funcs ...string) Condition {
	panic("not implemented")
	return inCondition{}
}

// OrIn returns a condition in the format IN(?,?,?) adjoined with OR
func OrIn(field, operator string, value interface{}, funcs ...string) Condition {
	panic("not implemented")
	return inCondition{}
}

// AndNotIn returns a condition in the format NOT IN(?,?,?) adjoined with AND
func AndNotIn(field, operator string, value interface{}, funcs ...string) Condition {
	panic("not implemented")
	return inCondition{}
}

// OrNotIn returns a condition in the format NOT IN(?,?,?) adjoined with OR
func OrNotIn(field, operator string, value interface{}, funcs ...string) Condition {
	panic("not implemented")
	return inCondition{}
}

// ConditionType to satisfy interface Condition
func (c inCondition) ConditionType() uint8 {
	return c.Type
}

func (c inCondition) AsSQL() (string, []interface{}) {
	panic("not implemented")
	return "", nil
}

func queryIN(val interface{}) (string, []interface{}) {
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
		//panic(fmt.Sprintf("Unsupported type for IN condition. <%T> %v", val, val))
		return "", nil
	}

	if l := len(values); l != 0 {
		query = "(" + strings.TrimSuffix(strings.Repeat("?,", l), ",") + ")"
	}

	return query, values
}

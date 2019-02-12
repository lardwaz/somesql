package somesql

// getFieldValueFunctions is a helper function to get field and value funcs for sql from a list of 0 or more args
func getFieldValueFunctions(funcs []string) (string, string) {
	if l := len(funcs); l == 0 {
		return None, None
	} else if l == 1 {
		return funcs[0], None
	}
	return funcs[0], funcs[1]
}

func expandValues(val interface{}) []interface{} {
	var (
		values []interface{}
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
		return []interface{}{val}
	}

	return values
}

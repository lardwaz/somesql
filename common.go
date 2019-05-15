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
	} else if data, ok := val.([]int); ok {
		for d := range data {
			values = append(values, data[d])
		}
	} else {
		return []interface{}{val}
	}

	return values
}

// getSliceChange returns all elements that are present in sliceTwo but NOT in sliceOne
func getSliceChange(sliceOne, sliceTwo []string) []string {
	m := make(map[string]bool)

	for _, item := range sliceOne {
		m[item] = true
	}

	change := make([]string, 0)
	for _, item := range sliceTwo {
		if _, ok := m[item]; !ok {
			change = append(change, item)
		}
	}

	return change
}

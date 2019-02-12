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
	} else {
		return []interface{}{val}
	}

	return values
}

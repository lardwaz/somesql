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

func expandValues(val interface{}) ([]interface{}, bool) {
	var (
		values   []interface{}
		asserted bool
	)

	if data, ok := val.([]string); ok {
		for d := range data {
			values = append(values, data[d])
		}
		asserted = true
	} else if data, ok := val.([]bool); ok {
		for d := range data {
			values = append(values, data[d])
		}
		asserted = true
	} else if data, ok := val.([]int); ok {
		for d := range data {
			values = append(values, data[d])
		}
		asserted = true
	} else {
		return []interface{}{val}, false
	}

	return values, asserted
}

// getSliceChange returns all elements that are present in sliceTwo but NOT in sliceOne
// it can be used for several purposes. for example if we have 2 slices:
// - s1 [a, b, c]
// - s2 [c, d]
//
// 1. Say s2 is a new slice derived from s1, where elements have been added and removed
//	a) to get elements that have been added in s2
//		added = getSliceChange(s1, s2)
//		added = [d]
//	b) to get elements that have been removed from s1
//		removed = getSliceChange(s2, s1)
//		removed = [a, b]
// 2. Remove all elements from s1 that are in s2
//		s1 = getSliceChange(s2, s1)
//		s1 = [a, b]
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

func getSQLType(v interface{}) string {
	switch v.(type) {
	case uint8, uint16, uint32, uint64, int8, int16, int32, int64, int, uint:
		return "INT"
	case string:
		return "TEXT"
	case bool:
		return "BOOLEAN"
	}

	return "TEXT"
}

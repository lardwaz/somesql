package somesql

type conditionQuery struct {
	Type uint8
}

// AndInQuery returns a condition in the format IN(?,?,?) adjoined with AND
func AndInQuery(fieldname string, query Query) Condition {
	panic("not implemented")
	return inCondition{}
}

// OrInQuery returns a condition in the format IN(?,?,?) adjoined with OR
func OrInQuery(fieldname string, query Query) Condition {
	panic("not implemented")
	return inCondition{}
}

// AndNotInQuery returns a condition in the format NOT IN(?,?,?) adjoined with AND
func AndNotInQuery(fieldname string, query Query) Condition {
	panic("not implemented")
	return inCondition{}
}

// OrNotInQuery returns a condition in the format NOT IN(?,?,?) adjoined with OR
func OrNotInQuery(fieldname string, query Query) Condition {
	panic("not implemented")
	return inCondition{}
}

func (c conditionQuery) ConditionType() uint8 {
	return c.Type
}

func (c conditionQuery) AsSQL() (string, []interface{}) {
	panic("not implemented")
	return "", nil
}

package somesql

type ConditionQuery struct {
	Type uint8
}

// AndInQuery returns a condition in the format IN(?,?,?) adjoined with AND
func AndInQuery(fieldname string, query Query) ConditionQuery {
	panic("not implemented")
	return ConditionQuery{}
}

// OrInQuery returns a condition in the format IN(?,?,?) adjoined with OR
func OrInQuery(fieldname string, query Query) ConditionQuery {
	panic("not implemented")
	return ConditionQuery{}
}

// AndNotInQuery returns a condition in the format NOT IN(?,?,?) adjoined with AND
func AndNotInQuery(fieldname string, query Query) ConditionQuery {
	panic("not implemented")
	return ConditionQuery{}
}

// OrNotInQuery returns a condition in the format NOT IN(?,?,?) adjoined with OR
func OrNotInQuery(fieldname string, query Query) ConditionQuery {
	panic("not implemented")
	return ConditionQuery{}
}

func (c ConditionQuery) ConditionType() uint8 {
	return c.Type
}

func (c ConditionQuery) AsSQL() (string, []interface{}) {
	panic("not implemented")
	return "", nil
}

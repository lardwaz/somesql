package somesql

// ConditionQuery defines conditions for a query
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

// ConditionType return the condition type (or / and)
func (c ConditionQuery) ConditionType() uint8 {
	return c.Type
}

// AsSQL returns part of SQL incuding the sub-query
func (c ConditionQuery) AsSQL() (string, []interface{}) {
	panic("not implemented")
	// return "", []interface{}{}
}

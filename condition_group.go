package somesql

//ConditionGroup represents a group of condition (within same pair brackets)
type ConditionGroup struct {
	Type       uint8
	Conditions []Condition
}

//AndGroup creates a condition group in the format AND((condition1) OR (condition2) AND (condition3))
func AndGroup(conditions ...Condition) Condition {
	panic("not implemented")
	return ConditionGroup{}
}

//OrGroup creates a condition group in the format OR((condition1) OR (condition2) AND (condition3))
func OrGroup(conditions ...Condition) Condition {
	panic("not implemented")
	return ConditionGroup{}
}

//ConditionType to satisfy interface Condition
func (c ConditionGroup) ConditionType() uint8 {
	return c.Type
}

//AsSQL to satisfy interface Condition
func (c ConditionGroup) AsSQL() (string, []interface{}) {
	panic("not implemented")
	return "", []interface{}{}
}

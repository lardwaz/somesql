package somesql

const (
	//AndCondition represents a condition added to the query via AND keyword
	AndCondition uint8 = iota

	//OrCondition represents a condition added to the query via OR keyword
	OrCondition
)

const (
	//None represents a simple way of explicitly specifying no value
	None = ""
)

//Query represents a composable query
type Query interface {
	Select(...string) Query
	Where(Condition) Query
	AsSQL() (string, []interface{})
}

//Condition represents a conditional clause in a statement
type Condition interface {
	ConditionType() uint8
	AsSQL() (string, []interface{})
}

package somesql //import go.lsl.digital/gocipe/somesql

const (
	// AndCondition represents a condition added to the query via AND keyword
	AndCondition uint8 = iota
	// OrCondition represents a condition added to the query via OR keyword
	OrCondition

	// LangEN represents the english language
	LangEN string = "en"
	// LangFR represents the french language
	LangFR string = "fr"
)

const (
	// None represents a simple way of explicitly specifying no value
	None = ""
)

// Query represents a composable query
// Setters: Select(), Merge(), Delete(), Where(), SetLang(), SetLimit(), SetOffset()
// Getters: GetLang(), GetLimit(), GetOffset(), AsSQL()
type Query interface {
	Select(...string) Query
	Merge(...string) Query
	Delete() Query
	Where(Condition) Query
	SetLang(string) Query
	GetLang() string
	SetLimit(int) Query
	GetLimit() int
	SetOffset(int) Query
	GetOffset() int
	AsSQL(in ...bool) (string, []interface{})
}

// Condition represents a conditional clause in a statement
type Condition interface {
	ConditionType() uint8
	AsSQL(in ...bool) (string, []interface{})
}

package somesql

// PQQuery represents a query implementation in postgresql database backend
type PQQuery struct {
	Lang       string
	Fields     []string
	Conditions []Condition
}

// NewQuery declares a new query
func NewQuery() Query {
	var q PQQuery
	return q
}

// Select specifies which fields to retrieve data for
func (q PQQuery) Select(f ...string) Query {
	q.Fields = f
	return q
}

// Where adds a condition clause to the Query
func (q PQQuery) Where(c Condition) Query {
	q.Conditions = append(q.Conditions, c)
	return q
}

// AsSQL returns the sql query and values for the query
func (q PQQuery) AsSQL() (string, []interface{}) {
	var values []interface{}

	sql := `SELECT `

	for i := range q.Conditions {
		switch q.Conditions[i].ConditionType() {
		case AndCondition:
			sql += `AND `
		case OrCondition:
			sql += `OR `
		default:
			continue
		}

		s, v := q.Conditions[i].AsSQL()
		sql += s
		values = append(values, v...)
	}

	return sql, values
}

// SetLang is a setter for Language
func (q PQQuery) SetLang(lang string) Query {
	switch lang {
	case LangEN, LangFR:
		q.Lang = lang
	}
	return q
}

// GetLang is a getter for Language
func (q PQQuery) GetLang() string {
	if q.Lang != "" {
		return q.Lang
	}
	return LangEN
}

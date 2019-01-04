package postgres

import "github.com/fluxynet/somesql"

//Query represents a query implementation in postgresql database backend
type Query struct {
	Fields     []string
	Conditions []somesql.Condition
}

//New declares a new query
func New() somesql.Query {
	var q Query

	return q
}

//Select specifies which fields to retrieve data for
func (q Query) Select(f ...string) somesql.Query {
	q.Fields = f
	return q
}

//Where adds a condition clause to the Query
func (q Query) Where(c somesql.Condition) somesql.Query {
	q.Conditions = append(q.Conditions, c)
	return q
}

//AsSQL returns the sql query and values for the query
func (q Query) AsSQL() (string, []interface{}) {
	var values []interface{}

	sql := `SELECT `

	for i := range q.Conditions {
		switch q.Conditions[i].ConditionType() {
		case somesql.AndCondition:
			sql += `AND `
		case somesql.OrCondition:
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

package somesql

import (
	"fmt"
	"strings"
)

// PQQuery represents a query implementation in postgresql database backend
type PQQuery struct {
	Lang       string
	Fields     []string
	Conditions []Condition
	Limit 	   int
	Offset 	   int
}

// NewQuery declares a new query
func NewQuery() Query {
	var q PQQuery
	q.Fields = []string{"id", "created_at", "updated_at", "owner_id", "status", "type", "data"}
	q.Limit = 10
	return q
}

// Select specifies which fields to retrieve data for
func (q PQQuery) Select(f ...string) Query {
	if len(f) == 0 {
		return q
	}
	q.Fields = f
	return q
}

// Where adds a condition clause to the Query
func (q PQQuery) Where(c Condition) Query {
	q.Conditions = append(q.Conditions, c)
	return q
}

// AsSQL returns the sql query and values for the query
func (q PQQuery) AsSQL(in ...bool) (string, []interface{}) {
	var (
		values []interface{}
		lang   = q.GetLang()
	)

	sql := `SELECT`

	var dataFields []string
	for _, field := range q.Fields {
		switch field {
		case "id", "created_at", "updated_at", "owner_id", "status", "type":
			sql += " " + field + ","
		case "data":
			sql += " data_" + lang + ","
		default:
			dataFields = append(dataFields, field)
		}
	}

	inner := (len(in) != 0 && in[0])

	dataFieldsLen := len(dataFields)
	for i, dataField := range dataFields {
		if inner {
			sql += fmt.Sprintf(` data_%s->>'%s' "%s",`, lang, dataField, dataField)
		} else {
			if i == 0 { // Genesis
				sql += ` json_build_object(`
			}
			sql += fmt.Sprintf(`'%s', data_%s->'%s', `, dataField, lang, dataField)
			if (dataFieldsLen) == i+1 { // End
				sql = strings.TrimRight(sql, ", ")
				sql += `) "data",`
			}
		}
	}

	sql = strings.TrimRight(sql, ",")

	sql += " FROM repo"

	var conditions string
	for i, cond := range q.Conditions {
		if i != 0 {
			switch cond.ConditionType() {
			case AndCondition:
				conditions += ` AND `
			case OrCondition:
				conditions += ` OR `
			default:
				continue
			}
		}

		s, v := cond.AsSQL()
		conditions += s
		values = append(values, v...)
	}
	
	if len(conditions) != 0 {
		sql += " WHERE " + conditions
	}

	if q.GetLimit() != 0 {
		sql += fmt.Sprintf(" LIMIT %d", q.Limit) 
	}

	if q.GetOffset() != 0 {
		sql += fmt.Sprintf(" OFFSET %d", q.Offset) 
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

// SetLimit is a setter for Limit
func (q PQQuery) SetLimit(limit int) Query {
	q.Limit = limit
	return q
}

// GetLimit is a getter for Limit
func (q PQQuery) GetLimit() int {
	return q.Limit
}

// SetOffset is a setter for Offset
func (q PQQuery) SetOffset(offset int) Query {
	q.Offset = offset
	return q
}

// GetOffset is a getter for Offset
func (q PQQuery) GetOffset() int {
	return q.Offset
}
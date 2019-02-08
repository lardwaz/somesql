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
}

// NewQuery declares a new query
func NewQuery() Query {
	var q PQQuery
	q.Fields = []string{"id", "created_at", "updated_at", "owner_id", "status", "type", "data"}
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

	dataFieldsLen := len(dataFields)
	for i, dataField := range dataFields {
		if i == 0 { // Genesis
			sql += ` json_build_object(`
		}
		sql += fmt.Sprintf(`'%s', data_%s->'%s', `, dataField, lang, dataField)
		if (dataFieldsLen) == i+1 { // End
			sql = strings.TrimRight(sql, ", ")
			sql += `) "data",`
		}
	}

	sql = strings.TrimRight(sql, ",")

	sql += " FROM repo"

	for i, cond := range q.Conditions {
		if i == 0 {
			sql += ` WHERE `
		} else {
			switch cond.ConditionType() {
			case AndCondition:
				sql += ` AND `
			case OrCondition:
				sql += ` OR `
			default:
				continue
			}
		}

		s, v := cond.AsSQL()
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

package somesql

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

const (
	// UnknownQueryType represents an UNKNOWN query type
	UnknownQueryType uint8 = iota

	// SelectQueryType represents a SELECT query type
	SelectQueryType

	// SaveQueryType represents a SAVE query type
	// INSERT ON CONFLICT DO UPDATE
	SaveQueryType

	// DeleteQueryType represents a DELETE query type
	DeleteQueryType
)

// PQQuery represents a query implementation in postgresql database backend
type PQQuery struct {
	Type       uint8 // Default: UnknownQueryType
	Lang       string
	Fields     []string
	Conditions []Condition
	Limit      int
	Offset     int
}

// NewQuery declares a new query
func NewQuery() Query {
	var q PQQuery
	q.Fields = append(ReservedFields, FieldData)
	q.Limit = 10
	return q
}

// Select specifies which fields to retrieve data for
func (q PQQuery) Select(f ...string) Query {
	q.Type = SelectQueryType

	if len(f) == 0 {
		return q
	}
	q.Fields = f
	return q
}

// Save specifies a SAVE query
func (q PQQuery) Save() Query {
	q.Type = SaveQueryType
	q.Limit = 0

	return q
}

// Delete specifies a DELETE query
func (q PQQuery) Delete() Query {
	q.Type = DeleteQueryType
	q.Fields = nil
	q.Limit = 0
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
		err        error
		sql        string
		values     []interface{}
		dataFields []string
		metaFields []string

		lang      = q.GetLang()
		fieldData = GetFieldData(lang)
		inner     = (len(in) != 0 && in[0])
		t         = template.New("queries").Funcs(funcMap)
	)

	switch q.Type {
	case SelectQueryType, UnknownQueryType:
		t, err = t.Parse(selectTplStr)
	case SaveQueryType:
		t, err = t.Parse(saveTplStr)
	case DeleteQueryType:
		t, err = t.Parse(deleteTplStr)
	}
	if err != nil {
		return sql, values
	}

	for _, field := range q.Fields {
		if IsFieldMeta(field) {
			metaFields = append(metaFields, field)
		} else if field == "data" {
			metaFields = append(metaFields, fieldData)
		} else {
			dataFields = append(dataFields, field)
		}
	}

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

	var buf bytes.Buffer
	err = t.Execute(&buf, struct {
		Query         PQQuery
		MetaFields    []string
		DataFields    []string
		FieldData     string
		FieldDataLang string
		Conditions    string
		Inner         bool
	}{
		Query:         q,
		MetaFields:    metaFields,
		DataFields:    dataFields,
		FieldData:     FieldData,
		FieldDataLang: fieldData,
		Conditions:    conditions,
		Inner:         inner,
	})
	if err != nil {
		return sql, values
	}

	sql = buf.String()

	// Inner SQL we return here
	if inner {
		return sql, values
	}

	// Replace all '?' with increasing '$N' (i.e $1,$2,$3)
	var i int
	for _, r := range sql {
		if r == '?' {
			i++
			placeholder := fmt.Sprintf("$%d", i)
			sql = strings.Replace(sql, "?", placeholder, 1)
		}
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

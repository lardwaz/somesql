package somesql

import (
	"bytes"
	"database/sql"
	"fmt"
	"strings"
	"text/template"
)

const (
	// UnknownQueryType represents an UNKNOWN query type
	UnknownQueryType uint8 = iota

	// InsertQueryType represents a INSERT query type
	InsertQueryType

	// SelectQueryType represents a SELECT query type
	SelectQueryType

	// UpdateQueryType represents a UPDATE query type
	UpdateQueryType

	// DeleteQueryType represents a DELETE query type
	DeleteQueryType
)

// PQQuery represents a query implementation in postgresql database backend
type PQQuery struct {
	Type       uint8 // Default: UnknownQueryType
	Lang       string
	Fields     []string
	Conditions []Condition
	Values     []interface{}
	Limit      int
	Offset     int
	Inner      bool
	DB         *sql.DB
	Tx         *sql.Tx
}

// NewQuery declares a new query
func NewQuery(db ...*sql.DB) Query {
	var q PQQuery
	q.Fields = append(ReservedFields, FieldData)
	q.Limit = 10

	if len(db) > 0 {
		q.DB = db[0]
	}

	return q
}

// NewInnerQuery declares a new query
func NewInnerQuery() Query {
	q := NewQuery(nil)
	return q.SetInner(true)
}

// Insert specifies an INSERT query
func (q PQQuery) Insert(fieldValue FieldValuer) Query {
	q.Type = InsertQueryType

	q.Fields, q.Values = fieldValue.List()

	return q
}

// Select specifies which fields to retrieve data for
func (q PQQuery) Select(fields ...string) Query {
	q.Type = SelectQueryType

	if len(fields) == 0 {
		return q
	}
	q.Fields = fields
	return q
}

// Update specifies an UPDATE query
func (q PQQuery) Update(fieldValue FieldValuer) Query {
	q.Type = UpdateQueryType

	q.Fields, q.Values = fieldValue.List()

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

// AsSQL returns the result for the query
func (q PQQuery) AsSQL() QueryResulter {
	var (
		err        error
		sql        string
		dataFields []string
		metaFields []string

		values    = q.Values
		lang      = q.GetLang()
		fieldData = GetFieldData(lang)
		isInner   = q.IsInner()
		t         = template.New("queries").Funcs(funcMap)
	)

	switch q.Type {
	default:
		fallthrough
	case SelectQueryType, UnknownQueryType:
		t, err = t.Parse(selectTplStr)
	case InsertQueryType:
		t, err = t.Parse(insertTplStr)
	case UpdateQueryType:
		t, err = t.Parse(updateTplStr)
	case DeleteQueryType:
		t, err = t.Parse(deleteTplStr)
	}
	if err != nil {
		return NewQueryResult(q, sql, values)
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
		Inner:         isInner,
	})
	if err != nil {
		return NewQueryResult(q, sql, values)
	}

	sql = buf.String()

	// Inner SQL we return here
	if isInner {
		return NewQueryResult(q, sql, values)
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

	return NewQueryResult(q, sql, values)
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

// SetDB is a setter for sql.DB
func (q PQQuery) SetDB(db *sql.DB) Query {
	q.DB = db
	return q
}

// GetDB is a getter for sql.DB
func (q PQQuery) GetDB() *sql.DB {
	return q.DB
}

// SetInner is a setter for Limit
func (q PQQuery) SetInner(inner bool) Query {
	q.Inner = inner
	return q
}

// IsInner is a getter for inner
func (q PQQuery) IsInner() bool {
	return q.Inner
}

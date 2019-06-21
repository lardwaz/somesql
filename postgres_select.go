package somesql

import (
	"database/sql"
	"strconv"
	"strings"
)

// Select generates Postgres SELECT statement
// Implements: Accessor
type Select struct {
	fields     []string
	conditions []Condition
	inner      bool
	offset     int
	limit      int
	sql        string
	values     []interface{}
	db         *sql.DB
	lang       string
}

// NewSelect returns a new Select
func NewSelect(db ...*sql.DB) *Select {
	var s Select

	s.lang = LangEN
	s.fields = FieldsList
	s.limit = 10

	if len(db) > 0 {
		s.db = db[0]
	}

	return &s
}

// NewSelectInner returns a new inner Select
func NewSelectInner(db ...*sql.DB) *Select {
	s := NewSelect(db...)
	s.SetInner(true)

	return s
}

// NewSelectLang returns a new Delete with specific lang
func NewSelectLang(lang string, db ...*sql.DB) *Select {
	s := NewSelect(db...)
	s.SetLang(lang)

	return s
}

// NewSelectLangInner returns a new Delete with specific lang
func NewSelectLangInner(lang string, db ...*sql.DB) *Select {
	s := NewSelectInner(db...)
	s.SetLang(lang)

	return s
}

// SetDB implements Statement
func (s *Select) SetDB(db *sql.DB) {
	s.db = db
}

// GetDB implements Statement
func (s Select) GetDB() *sql.DB {
	return s.db
}

// SetLang implements Statement
func (s *Select) SetLang(lang string) {
	s.lang = lang
}

// GetLang implements Statement
func (s Select) GetLang() string {
	return s.lang
}

// GetSQL implements Statement
func (s Select) GetSQL() string {
	return s.sql
}

// GetValues implements Statement
func (s Select) GetValues() []interface{} {
	return s.values
}

// ToSQL implements Statement
func (s *Select) ToSQL() {
	var (
		fieldsStr     string
		conditionsStr string
		offsetStr     string
		limitStr      string
		isInnerQuery  = s.IsInner()
		dataFieldLang string

		fieldsBuff     strings.Builder
		metaFieldsBuff strings.Builder
		dataFieldsBuff strings.Builder
		relFieldsBuff  strings.Builder
	)

	if !IsLangValid(s.GetLang()) {
		s.SetLang(LangEN)
	}

	dataFieldLang = GetLangFieldData(s.GetLang())

	// Processing fields
	for _, f := range s.fields {
		if IsFieldMeta(f) || IsFieldData(f) || IsFieldRelations(f) {
			if f == FieldData {
				f = dataFieldLang
			}
			metaFieldsBuff.WriteString(`"` + f + `", `)
		} else if innerField, ok := GetInnerField(FieldData, f); ok {
			if isInnerQuery {
				dataFieldsBuff.WriteString(`"` + dataFieldLang + `"->>'` + innerField + `' "` + innerField + `", `)
			} else {
				dataFieldsBuff.WriteString(`'` + innerField + `', "` + dataFieldLang + `"->'` + innerField + `', `)
			}
		} else if innerField, ok := GetInnerField(FieldRelations, f); ok {
			if isInnerQuery {
				relFieldsBuff.WriteString(`"` + FieldRelations + `"->>'` + innerField + `' "` + innerField + `", `)
			} else {
				relFieldsBuff.WriteString(`'` + innerField + `', "` + FieldRelations + `"->'` + innerField + `', `)
			}
		}
	}

	// Put everything back in order
	// Meta fields
	if metaFieldsBuff.Len() > 0 {
		metaFieldsStr := metaFieldsBuff.String()[:metaFieldsBuff.Len()-2] // trim ", "
		fieldsBuff.WriteString(metaFieldsStr + `, `)
	}

	// Data fields
	if dataFieldsBuff.Len() > 0 {
		dataFieldsStr := dataFieldsBuff.String()[:dataFieldsBuff.Len()-2] // trim ", "
		if isInnerQuery {
			fieldsBuff.WriteString(dataFieldsStr + `, `)
		} else {
			fieldsBuff.WriteString(`json_build_object(` + dataFieldsStr + `) "` + FieldData + `", `)
		}
	}

	// Relationship fields
	if relFieldsBuff.Len() > 0 {
		relFieldsStr := relFieldsBuff.String()[:relFieldsBuff.Len()-2] // trim ", "
		if isInnerQuery {
			fieldsBuff.WriteString(relFieldsStr + `, `)
		} else {
			fieldsBuff.WriteString(`json_build_object(` + relFieldsStr + `) "` + FieldRelations + `", `)
		}
	}

	if fieldsBuff.Len() > 0 {
		fieldsStr = fieldsBuff.String()[:fieldsBuff.Len()-2] // trim ", "
	}

	conditions, condValues := processConditions(s.conditions)
	s.values = condValues

	if len(conditions) > 0 {
		conditionsStr = "WHERE " + conditions
	}

	if s.limit > 0 {
		limitStr = "LIMIT " + strconv.Itoa(s.limit)
	}

	if s.offset > 0 {
		offsetStr = "OFFSET " + strconv.Itoa(s.offset)
	}

	sql := "SELECT " + fieldsStr + " FROM " + Table + " " + conditionsStr + " " + limitStr + " " + offsetStr

	if !isInnerQuery {
		sql = processPlaceholders(sql)
	}

	s.sql = cleanStatement(sql)
}

// SetInner implements Accessor
func (s *Select) SetInner(inner bool) {
	s.inner = inner
}

// IsInner implements Accessor
func (s Select) IsInner() bool {
	return s.inner
}

// Rows implements Accessor
func (s Select) Rows() (*sql.Rows, error) {
	if s.GetSQL() == "" || len(s.GetValues()) == 0 {
		s.ToSQL()
	}

	return rows(s.GetSQL(), s.GetValues(), s.GetDB())
}

// Fields sets the fields for Select
func (s *Select) Fields(fields ...string) *Select {
	if len(fields) == 0 {
		return s
	}
	s.fields = fields
	return s
}

// Where adds a condition clause to the Query
func (s *Select) Where(c Condition) *Select {
	s.conditions = append(s.conditions, c)
	return s
}

// Offset sets the Offset for Select
func (s *Select) Offset(offset int) *Select {
	s.offset = offset
	return s
}

// Limit sets the Limit for Select
func (s *Select) Limit(limit int) *Select {
	s.limit = limit
	return s
}

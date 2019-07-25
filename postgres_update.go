package somesql

import (
	"database/sql"
	"encoding/json"
	"strings"
)

// Update generates Postgres UPDATE statement
// Implements: Mutator
type Update struct {
	fields     Fields
	conditions []Condition
	sql        string
	values     []interface{}
	db         *sql.DB
	lang       string
}

// NewUpdate returns a new Update
func NewUpdate(db ...*sql.DB) *Update {
	var s Update

	s.lang = LangEN

	if len(db) > 0 {
		s.db = db[0]
	}

	return &s
}

// NewUpdateLang returns a new Update with specific lang
func NewUpdateLang(lang string, db ...*sql.DB) *Update {
	s := NewUpdate(db...)
	s.SetLang(lang)

	return s
}

// SetDB implements Statement
func (s *Update) SetDB(db *sql.DB) {
	s.db = db
}

// GetDB implements Statement
func (s Update) GetDB() *sql.DB {
	return s.db
}

// SetLang implements Statement
func (s *Update) SetLang(lang string) {
	s.lang = lang
}

// GetLang implements Statement
func (s Update) GetLang() string {
	return s.lang
}

// GetSQL implements Statement
func (s Update) GetSQL() string {
	return s.sql
}

// GetValues implements Statement
func (s Update) GetValues() []interface{} {
	return s.values
}

// ToSQL implements Statement
func (s *Update) ToSQL() {
	var (
		fieldsStr     string
		conditionsStr string
		dataFieldLang string

		fieldsBuff     strings.Builder
		metaFieldsBuff strings.Builder
		dataFieldsBuff strings.Builder

		metaValues []interface{}
		dataValues []interface{}
	)

	if !IsLangValid(s.GetLang()) {
		s.SetLang(LangEN)
	}

	dataFieldLang = GetLangFieldData(s.GetLang())

	fields, values := s.fields.List()

	// Processing fields and values
	for i, f := range fields {
		if IsFieldData(f) {
			if jsonbFields, ok := values[i].(JSONBFields); ok {
				innerFields, innerValues, _ := jsonbFields.GetOrderedList()
				for idx, innerField := range innerFields {
					if _, ok := innerValues[idx].([]interface{}); ok {
						dataFieldsBuff.WriteString(`'` + innerField + `', ?::JSONB, `)
						if jsonBytes, err := json.Marshal(innerValues[idx]); err == nil {
							dataValues = append(dataValues, string(jsonBytes))
						}
					} else {

						dataFieldsBuff.WriteString(`'` + innerField + `', ?::text, `)
						dataValues = append(dataValues, innerValues[idx])
					}
				}
			}
		} else if IsFieldMeta(f) { // Check if Meta fields
			metaFieldsBuff.WriteString(`"` + f + `" = ?, `)
			metaValues = append(metaValues, values[i])
		}
	}

	s.values = make([]interface{}, 0)

	// Put everything back in order
	// Set meta fields
	if metaFieldsBuff.Len() > 0 {
		metaFieldsStr := metaFieldsBuff.String()[:metaFieldsBuff.Len()-2] // trim ", "
		fieldsBuff.WriteString(metaFieldsStr + `, `)
		s.values = append(s.values, metaValues...)
	}

	// Set data fields
	if dataFieldsBuff.Len() > 0 {
		dataFieldsStr := dataFieldsBuff.String()[:dataFieldsBuff.Len()-2] // trim ", "
		fieldsBuff.WriteString(`"` + dataFieldLang + `" = jsonb_build_object(` + dataFieldsStr + `)::JSONB, `)
		s.values = append(s.values, dataValues...)
	}

	conditions, condValues := processConditions(s.conditions)
	if len(conditions) > 0 {
		conditionsStr = " WHERE " + conditions
	}

	if fieldsBuff.Len() > 0 {
		fieldsStr = fieldsBuff.String()[:fieldsBuff.Len()-2] // trim ", "
	}

	s.values = append(s.values, condValues...)

	sql := "UPDATE " + Table + " SET " + fieldsStr + " " + conditionsStr

	s.sql = cleanStatement(processPlaceholders(sql))
}

// Exec implements Mutator
func (s Update) Exec(autocommit bool) error {
	if s.GetSQL() == "" || len(s.GetValues()) == 0 {
		s.ToSQL()
	}

	return exec(s.GetSQL(), s.GetValues(), s.GetDB(), autocommit)
}

// ExecTx implements Mutator
func (s Update) ExecTx(tx *sql.Tx, autocommit bool) error {
	if s.GetSQL() == "" || len(s.GetValues()) == 0 {
		s.ToSQL()
	}

	return execTx(s.GetSQL(), s.GetValues(), tx, autocommit)
}

// Fields sets the fields and values for Update
func (s *Update) Fields(fields Fields) *Update {
	s.fields = fields
	return s
}

// Where adds a condition clause to the Query
func (s *Update) Where(c Condition) *Update {
	s.conditions = append(s.conditions, c)
	return s
}

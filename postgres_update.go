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

		fieldsBuff           strings.Builder
		metaFieldsBuff       strings.Builder
		dataFieldsBuff       strings.Builder
		relFieldsBuff        strings.Builder
		relFieldsAddBuff     strings.Builder
		relFieldsAddBuff2    strings.Builder
		relFieldsRemoveBuff  strings.Builder
		relFieldsRemoveBuff2 strings.Builder
		relFieldsRemoveBuff3 strings.Builder
		relFieldsRemoveBuff4 strings.Builder

		metaValues      []interface{}
		dataValues      []interface{}
		relValues       []interface{}
		relValuesAdd    []interface{}
		relValuesRemove []interface{}
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
				for _, innerField := range innerFields {
					dataFieldsBuff.WriteString(`"` + innerField + `": ?, `)
				}
				dataValues = innerValues
			}
		} else if IsFieldRelations(f) {
			if jsonbFields, ok := values[i].(JSONBFields); ok {
				innerFields, innerValues, innerActions := jsonbFields.GetOrderedList()
				for idx, innerField := range innerFields {
					switch innerActions[idx] {
					case NoneJSONBArr:
						relFieldsBuff.WriteString(`"` + innerField + `": ?`)
						relValues = append(relValues, innerValues[idx])
					case JSONBArrAdd:
						relFieldsAddBuff.WriteString(`("` + FieldRelations + `" - '` + innerField + `') || `)
						relFieldsAddBuff2.WriteString(`'` + innerField + `', "` + FieldRelations + `"->'` + innerField + `' || '?'::JSONB, `)
						if jsonBytes, err := json.Marshal(innerValues[idx]); err == nil {
							relValuesAdd = append(relValuesAdd, string(jsonBytes))
						}
					case JSONBArrRemove:
						relFieldsRemoveBuff.WriteString(`("` + FieldRelations + `" - '` + innerField + `') || `)
						relFieldsRemoveBuff2.WriteString(`'` + innerField + `', JSONB_AGG(` + innerField + `Upd), `)
						relFieldsRemoveBuff3.WriteString(`JSONB_ARRAY_ELEMENTS_TEXT("` + FieldRelations + `"->'` + innerField + `') ` + innerField + `Upd, `)
						relFieldsRemoveBuff4.WriteString(innerField + `Upd NOT IN (?) AND `)
						if jsonBytes, err := json.Marshal(innerValues[idx]); err == nil {
							relValuesRemove = append(relValuesRemove, string(jsonBytes))
						}
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
		fieldsBuff.WriteString(`"` + dataFieldLang + `" = "` + dataFieldLang + `" || {` + dataFieldsStr + `}, `)
		s.values = append(s.values, dataValues...)
	}

	// Set relationship fields
	if relFieldsBuff.Len() > 0 {
		relFieldsStr := relFieldsBuff.String()[:relFieldsBuff.Len()-2] // trim ", "
		fieldsBuff.WriteString(`"` + FieldRelations + `" = "` + FieldRelations + `" || {` + relFieldsStr + `}, `)
		s.values = append(s.values, relValues...)
	}

	conditions, condValues := processConditions(s.conditions)
	if len(conditions) > 0 {
		conditionsStr = " WHERE " + conditions
	}

	// Add relationship fields
	if relFieldsAddBuff.Len() > 0 {
		relFieldsAddStr := relFieldsAddBuff.String()[:relFieldsAddBuff.Len()-4]    // trim " || "
		relFieldsAddStr2 := relFieldsAddBuff2.String()[:relFieldsAddBuff2.Len()-2] // trim ", "
		fieldsBuff.WriteString(`"` + FieldRelations + `" = relAdd.` + FieldRelations + ` FROM `)
		fieldsBuff.WriteString(`(SELECT (` + relFieldsAddStr + ` || JSONB_BUILD_OBJECT(` + relFieldsAddStr2 + `))`)
		fieldsBuff.WriteString(` "` + FieldRelations + `" FROM ` + Table + conditionsStr + `) relAdd, `)
		s.values = append(s.values, relValuesAdd...)
		s.values = append(s.values, condValues...)
	}

	// Remove relationship fields
	if relFieldsRemoveBuff.Len() > 0 {
		relFieldsRemoveStr := relFieldsRemoveBuff.String()[:relFieldsRemoveBuff.Len()-4]    // trim " || "
		relFieldsRemoveStr2 := relFieldsRemoveBuff2.String()[:relFieldsRemoveBuff2.Len()-2] // trim ", "
		relFieldsRemoveStr3 := relFieldsRemoveBuff3.String()[:relFieldsRemoveBuff3.Len()-2] // trim ", "
		relFieldsRemoveStr4 := relFieldsRemoveBuff4.String()[:relFieldsRemoveBuff4.Len()-5] // trim " AND "
		fieldsBuff.WriteString(`"` + FieldRelations + `" = updates.updRel FROM (SELECT (` + relFieldsRemoveStr + ` || JSONB_BUILD_OBJECT(` + relFieldsRemoveStr2 + `))`)
		fieldsBuff.WriteString(` "updatedRel" FROM (SELECT "` + FieldRelations + `", ` + relFieldsRemoveStr3 + ` FROM ` + Table + conditionsStr + `)`)
		fieldsBuff.WriteString(` expandedValues WHERE ` + relFieldsRemoveStr4 + ` GROUP BY "` + FieldRelations + `") updates, `)
		s.values = append(s.values, condValues...)
		s.values = append(s.values, relValuesRemove...)
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

package somesql

import (
	"database/sql"
	"encoding/json"
	"fmt"
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
func NewUpdate(lang string, db ...*sql.DB) *Update {
	var s Update

	s.lang = lang

	if len(db) > 0 {
		s.db = db[0]
	}

	return &s
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
		dataFieldLang = GetLangFieldData(s.GetLang())
	)
	fields, values := s.fields.List()

	// Processing fields and values
	// TODO: this is fugly!!! (concat on the spot?)
	metaFields := make([]string, 0)
	metaValues := make([]interface{}, 0)
	dataFields := make([]string, 0)
	dataValues := make([]interface{}, 0)
	relFields := make([]string, 0)
	relFieldsAdd := make([]string, 0)
	relFieldsAdd2 := make([]string, 0)
	relValues := make([]interface{}, 0)
	relValuesAdd := make([]interface{}, 0)
	relFieldsRemove := make([]string, 0)
	relFieldsRemove2 := make([]string, 0)
	relFieldsRemove3 := make([]string, 0)
	relFieldsRemove4 := make([]string, 0)
	relValuesRemove := make([]interface{}, 0)
	for i, f := range fields {
		if IsWholeFieldData(f) {
			if jsonbFields, ok := values[i].(JSONBFields); ok {
				innerFields, innerValues, _ := jsonbFields.GetOrderedList()
				for _, innerField := range innerFields {
					dataFields = append(dataFields, fmt.Sprintf(`"%s": ?`, innerField))
				}
				dataValues = innerValues
			}
		} else if IsWholeFieldRelations(f) {
			if jsonbFields, ok := values[i].(JSONBFields); ok {
				innerFields, innerValues, innerActions := jsonbFields.GetOrderedList()
				for idx, innerField := range innerFields {
					switch innerActions[idx] {
					case JSONBArrSet:
						relFields = append(relFields, fmt.Sprintf(`"%s": ?`, innerField))
						relValues = append(relValues, innerValues[idx])
					case JSONBArrAdd:
						relFieldsAdd = append(relFieldsAdd, fmt.Sprintf(`("%s" - '%s')`, FieldRelations, innerField))
						relFieldsAdd2 = append(relFieldsAdd2, fmt.Sprintf(`'%s', "%s"->'%s' || '?'::JSONB`, innerField, FieldRelations, innerField))
						if jsonBytes, err := json.Marshal(innerValues[idx]); err == nil {
							relValuesAdd = append(relValuesAdd, string(jsonBytes))
						}
					case JSONBArrRemove:
						relFieldsRemove = append(relFieldsRemove, fmt.Sprintf(`("%s" - '%s')`, FieldRelations, innerField))
						relFieldsRemove2 = append(relFieldsRemove2, fmt.Sprintf(`'%s', JSONB_AGG(%sUpd)`, innerField, innerField))
						relFieldsRemove3 = append(relFieldsRemove3, fmt.Sprintf(`JSONB_ARRAY_ELEMENTS_TEXT("%s"->'%s') %sUpd`, FieldRelations, innerField, innerField))
						relFieldsRemove4 = append(relFieldsRemove4, fmt.Sprintf(`%sUpd NOT IN (?)`, innerField))
						if jsonBytes, err := json.Marshal(innerValues[idx]); err == nil {
							relValuesRemove = append(relValuesRemove, string(jsonBytes))
						}
					}
				}
			}
		} else if IsFieldMeta(f) { // Check if Meta fields
			metaFields = append(metaFields, fmt.Sprintf(`"%s" = ?`, f))
			metaValues = append(metaValues, values[i])
		}
	}

	s.values = make([]interface{}, 0)
	fieldsJoined := make([]string, 0)
	if len(metaFields) > 0 {
		fieldsJoined = append(fieldsJoined, strings.Join(metaFields, ", "))
		s.values = append(s.values, metaValues...)
	}

	if len(dataFields) > 0 {
		fieldsJoined = append(fieldsJoined, fmt.Sprintf(`"%s" = "%s" || {%s}`, dataFieldLang, dataFieldLang, strings.Join(dataFields, ", ")))
		s.values = append(s.values, dataValues...)
	}

	if len(relFields) > 0 {
		fieldsJoined = append(fieldsJoined, fmt.Sprintf(`"%s" = "%s" || {%s}`, dataFieldLang, dataFieldLang, strings.Join(dataFields, ", ")))
		s.values = append(s.values, relValues...)
	}

	conditions, condValues := processConditions(s.conditions)
	if len(conditions) > 0 {
		conditionsStr = fmt.Sprintf(" WHERE %s", conditions)
	}

	if len(relFieldsAdd) > 0 {
		relFieldsAddStr := strings.Join(relFieldsAdd, " || ")
		relFieldsAddStr2 := strings.Join(relFieldsAdd2, ", ")
		relAddStr := fmt.Sprintf(`"%s" = relAdd.%s FROM (SELECT (%s || JSONB_BUILD_OBJECT(%s)) "%s" FROM %s%s) relAdd`,
			FieldRelations, FieldRelations, relFieldsAddStr, relFieldsAddStr2, FieldRelations, Table, conditionsStr,
		)
		fieldsJoined = append(fieldsJoined, relAddStr)
		s.values = append(s.values, relValuesAdd...)
		s.values = append(s.values, condValues...)
	}

	if len(relFieldsRemove) > 0 {
		relFieldsRemoveStr := strings.Join(relFieldsRemove, " || ")
		relFieldsRemoveStr2 := strings.Join(relFieldsRemove2, ", ")
		relFieldsRemoveStr3 := strings.Join(relFieldsRemove3, ", ")
		relFieldsRemoveStr4 := strings.Join(relFieldsRemove4, " AND ")
		relRemoveStr := fmt.Sprintf(`"%s" = updates.updRel FROM (SELECT (%s || JSONB_BUILD_OBJECT(%s)) "updatedRel" FROM (SELECT "%s", %s FROM %s%s) expandedValues WHERE %s GROUP BY "%s") updates`,
			FieldRelations, relFieldsRemoveStr, relFieldsRemoveStr2, FieldRelations, relFieldsRemoveStr3, Table, conditionsStr, relFieldsRemoveStr4, FieldRelations,
		)
		fieldsJoined = append(fieldsJoined, relRemoveStr)
		s.values = append(s.values, condValues...)
		s.values = append(s.values, relValuesRemove...)
	}

	fieldsStr = strings.Join(fieldsJoined, ", ")

	s.values = append(s.values, condValues...)

	sql := fmt.Sprintf(`UPDATE %s SET %s %s`, Table, fieldsStr, conditionsStr)

	s.sql = cleanStatement(processPlaceholders(sql))
}

// Exec implements Mutator
func (s Update) Exec(autocommit bool) error {
	return exec(s.GetSQL(), s.GetValues(), s.GetDB(), autocommit)
}

// ExecTx implements Mutator
func (s Update) ExecTx(tx *sql.Tx, autocommit bool) error {
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

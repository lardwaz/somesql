package somesql

import (
	"database/sql"
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

	// log.Println(fields)

	// UPDATE repo SET{{ " " -}}
	// {{- range $i, $v := .MetaFieldsList -}}
	// 	"{{ $v }}" = ?
	// 	{{- if ne (len $.MetaFieldsList) (plus $i) }}, {{ end }}
	// {{- end -}}
	// {{- if and (ne (len .MetaFieldsList) 0) (ne (len .DataFields) 0) }}, {{ end -}}
	// {{- range $i, $v := .DataFields -}}
	// 	{{ if eq $i 0 }}"{{ $.FieldDataLang }}" = "{{ $.FieldDataLang }}" || { {{- end -}}
	// 	"{{ $v }}": ?
	// 	{{- if ne (len $.DataFields) (plus $i) }}, {{ end }}
	// 	{{- if eq (len $.DataFields) (plus $i) -}} } {{- end }}
	// {{- end }}
	// {{- if ne (len .Conditions) 0 }} WHERE {{ .Conditions }}{{ end }}

	// Processing fields and values
	metaFields := make([]string, 0)
	metaValues := make([]interface{}, 0)
	dataFields := make([]string, 0)
	dataValues := make([]interface{}, 0)
	// relFields := make([]string, 0)
	for i, f := range fields {
		if IsWholeFieldData(f) {
			if jsonbFields, ok := values[i].(JSONBFields); ok {
				innerFields, innerValues, _ := jsonbFields.GetOrderedList()
				for _, innerField := range innerFields {
					dataFields = append(dataFields, fmt.Sprintf(`"%s": ?`, innerField))
				}
				dataValues = append(dataValues, innerValues...)
			}
		} else if IsWholeFieldRelations(f) {
			// Le WIP
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

	// Le WIP
	// if len(relFields) > 0 {
	// 	if isInnerQuery {
	// 		fields = append(fields, strings.Join(relFields, ", "))
	// 	} else {
	// 		fields = append(fields, fmt.Sprintf(`json_build_object(%s) "%s"`, strings.Join(relFields, ", "), FieldRelations))
	// 	}
	// }

	fieldsStr = strings.Join(fieldsJoined, ", ")

	conditions, condValues := processConditions(s.conditions)
	s.values = append(s.values, condValues...)

	if len(conditions) > 0 {
		conditionsStr = fmt.Sprintf("WHERE %s", conditions)
	}

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

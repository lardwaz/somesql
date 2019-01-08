package somesql_test

import (
	"reflect"
	"testing"

	"github.com/lsldigital/somesql"
)

func TestConditionClause(t *testing.T) {
	const (
		caseAnd = iota
		caseOr
	)

	type args struct {
		field    string
		operator string
		value    interface{}
		funcs    []string
	}

	type testcase struct {
		name     string
		args     args
		sql      string
		values   []interface{}
		caseType uint8
	}

	tests := []testcase{
		{
			"AND Normal",
			args{
				field:    "name",
				operator: "=",
				value:    "John Doe",
				funcs:    []string{},
			},
			" AND name=?",
			[]interface{}{"John Doe"},
			caseAnd,
		},
		{
			"AND with func on value",
			args{
				field:    "name",
				operator: "=",
				value:    "John Doe",
				funcs:    []string{"", "LOWER"},
			},
			" AND name=LOWER(?)",
			[]interface{}{"John Doe"},
			caseAnd,
		},
		{
			"AND with func on field",
			args{
				field:    "name",
				operator: "=",
				value:    "John Doe",
				funcs:    []string{"LOWER"},
			},
			" AND LOWER(name)=?",
			[]interface{}{"John Doe"},
			caseAnd,
		},
		{
			"AND with func on field and value",
			args{
				field:    "name",
				operator: "=",
				value:    "John Doe",
				funcs:    []string{"LOWER", "LOWER"},
			},
			" AND LOWER(name)=LOWER(?)",
			[]interface{}{"John Doe"},
			caseAnd,
		},
		{
			"AND JSONB",
			args{
				field:    "name",
				operator: "=",
				value:    "John Doe",
				funcs:    []string{},
			},
			" AND data->'name' = ?",
			[]interface{}{"John Doe"},
			caseAnd,
		},
		{
			"AND JSONB with func on value",
			args{
				field:    "name",
				operator: "=",
				value:    "John Doe",
				funcs:    []string{"", "LOWER"},
			},
			" AND data->>'name'=LOWER(?)",
			[]interface{}{"John Doe"},
			caseAnd,
		},
		{
			"AND JSONB with func on field",
			args{
				field:    "name",
				operator: "=",
				value:    "John Doe",
				funcs:    []string{"LOWER"},
			},
			" AND LOWER(data->>'name')=?",
			[]interface{}{"John Doe"},
			caseAnd,
		},
		{
			"AND JSONB with func on field and value",
			args{
				field:    "name",
				operator: "=",
				value:    "John Doe",
				funcs:    []string{"LOWER", "LOWER"},
			},
			" AND LOWER(data->>'name')=LOWER(?)",
			[]interface{}{"John Doe"},
			caseAnd,
		},

		{
			"OR Normal",
			args{
				field:    "name",
				operator: "=",
				value:    "John Doe",
				funcs:    []string{},
			},
			" OR name=?",
			[]interface{}{"John Doe"},
			caseOr,
		},
		{
			"OR with func on field",
			args{
				field:    "name",
				operator: "=",
				value:    "John Doe",
				funcs:    []string{"LOWER"},
			},
			" OR LOWER(name)=?",
			[]interface{}{"John Doe"},
			caseOr,
		},
		{
			"OR with func on value",
			args{
				field:    "name",
				operator: "=",
				value:    "John Doe",
				funcs:    []string{"", "LOWER"},
			},
			" OR name=LOWER(?)",
			[]interface{}{"John Doe"},
			caseOr,
		},
		{
			"OR with func on field and value",
			args{
				field:    "name",
				operator: "=",
				value:    "John Doe",
				funcs:    []string{"LOWER", "LOWER"},
			},
			" OR LOWER(name)=LOWER(?)",
			[]interface{}{"John Doe"},
			caseOr,
		},
		{
			"OR JSONB",
			args{
				field:    "name",
				operator: "=",
				value:    "John Doe",
				funcs:    []string{},
			},
			" OR data->'name' = ?",
			[]interface{}{"John Doe"},
			caseOr,
		},
		{
			"OR JSONB with func on value",
			args{
				field:    "name",
				operator: "=",
				value:    "John Doe",
				funcs:    []string{"", "LOWER"},
			},
			" OR data->>'name'=LOWER(?)",
			[]interface{}{"John Doe"},
			caseOr,
		},
		{
			"OR JSONB with func on field",
			args{
				field:    "name",
				operator: "=",
				value:    "John Doe",
				funcs:    []string{"LOWER"},
			},
			" OR LOWER(data->>'name')=?",
			[]interface{}{"John Doe"},
			caseOr,
		},
		{
			"OR JSONB with func on field and value",
			args{
				field:    "name",
				operator: "=",
				value:    "John Doe",
				funcs:    []string{"LOWER", "LOWER"},
			},
			" OR LOWER(data->>'name')=LOWER(?)",
			[]interface{}{"John Doe"},
			caseOr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var (
				sql    string
				values interface{}
			)

			if tt.caseType == caseAnd {
				sql, values = somesql.And(tt.args.field, tt.args.operator, tt.args.value, tt.args.funcs...).AsSQL()
			} else {
				sql, values = somesql.Or(tt.args.field, tt.args.operator, tt.args.value, tt.args.funcs...).AsSQL()
			}

			if tt.sql != sql {
				t.Errorf("Got %s, want %s", sql, tt.sql)
			}

			if !reflect.DeepEqual(values, tt.values) {
				t.Errorf("Values = %v, want %v", values, tt.values)
			}
		})
	}
}

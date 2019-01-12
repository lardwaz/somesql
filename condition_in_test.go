package somesql_test

import (
	"reflect"
	"testing"

	"github.com/lsldigital/somesql"
)

func TestConditionIn(t *testing.T) {
	const (
		caseAndIn = iota
		caseAndNotIn
		caseOrIn
		caseOrNotIn
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
			"AND IN",
			args{
				field:    "id",
				operator: "IN",
				value:    []string{"A", "B", "C"},
			},
			"id IN (?,?,?)",
			[]interface{}{"A", "B", "C"},
			caseAndIn,
		},
		{
			"AND IN (with func on field)",
			args{
				field:    "updated_at",
				operator: "IN",
				value:    []string{"2019"},
				funcs:    []string{"YEAR"},
			},
			"YEAR(updated_at) IN (?)",
			[]interface{}{"2019"},
			caseAndIn,
		},
		{
			"AND NOT IN",
			args{
				field:    "id",
				operator: "NOT IN",
				value:    []string{"A", "B", "C"},
			},
			"id NOT IN (?,?,?)",
			[]interface{}{"A", "B", "C"},
			caseAndNotIn,
		},
		{
			"AND NOT IN (with func on field)",
			args{
				field:    "updated_at",
				operator: "NOT IN",
				value:    []string{"2016"},
				funcs:    []string{"YEAR"},
			},
			"YEAR(updated_at) NOT IN (?)",
			[]interface{}{"2016"},
			caseAndNotIn,
		},
		{
			"AND IN (JSONB)",
			args{
				field:    "name",
				operator: "IN",
				value:    []string{"A", "B", "C"},
			},
			`"data"->>'name' IN (?,?,?)`,
			[]interface{}{"A", "B", "C"},
			caseAndIn,
		},
		{
			"AND IN (JSONB) (with func on field)",
			args{
				field:    "badge",
				operator: "IN",
				value:    []string{"video", "audio"},
				funcs:    []string{"LOWER"},
			},
			`LOWER("data"->>'badge') IN (?,?)`,
			[]interface{}{"video", "audio"},
			caseAndIn,
		},
		{
			"AND NOT IN (JSONB)",
			args{
				field:    "name",
				operator: "NOT IN",
				value:    []string{"A", "B", "C"},
			},
			`"data"->>'name' NOT IN (?,?,?)`,
			[]interface{}{"A", "B", "C"},
			caseAndNotIn,
		},
		{
			"AND NOT IN (JSONB) (with func on field)",
			args{
				field:    "name",
				operator: "NOT IN",
				value:    []string{"video", "audio"},
				funcs:    []string{"LOWER"},
			},
			`LOWER("data"->>'badge') NOT IN (?,?)`,
			[]interface{}{"video", "audio"},
			caseAndNotIn,
		},

		{
			"OR IN",
			args{
				field:    "id",
				operator: "IN",
				value:    []string{"A", "B", "C"},
			},
			"id IN (?,?,?)",
			[]interface{}{"A", "B", "C"},
			caseOrIn,
		},
		{
			"OR IN (with func on field)",
			args{
				field:    "updated_at",
				operator: "IN",
				value:    []string{"2019"},
				funcs:    []string{"YEAR"},
			},
			"YEAR(updated_at) IN (?)",
			[]interface{}{"2019"},
			caseOrIn,
		},
		{
			"OR NOT IN",
			args{
				field:    "id",
				operator: "IN",
				value:    []string{"A", "B", "C"},
			},
			"id NOT IN (?,?,?)",
			[]interface{}{"A", "B", "C"},
			caseOrNotIn,
		},
		{
			"OR NOT IN (with func on field)",
			args{
				field:    "updated_at",
				operator: "NOT IN",
				value:    []string{"2015"},
				funcs:    []string{"YEAR"},
			},
			"YEAR(updated_at) NOT IN (?)",
			[]interface{}{"2015"},
			caseOrNotIn,
		},
		{
			"OR IN (JSONB)",
			args{
				field:    "name",
				operator: "IN",
				value:    []string{"A", "B", "C"},
			},
			`"data"->>'name' IN (?,?,?)`,
			[]interface{}{"A", "B", "C"},
			caseOrIn,
		},
		{
			"OR IN (JSONB) (with func on field)",
			args{
				field:    "badge",
				operator: "IN",
				value:    []string{"video", "audio"},
				funcs:    []string{"LOWER"},
			},
			`LOWER("data"->>'badge') IN (?,?)`,
			[]interface{}{"video", "audio"},
			caseOrIn,
		},
		{
			"OR NOT IN (JSONB)",
			args{
				field:    "name",
				operator: "NOT IN",
				value:    []string{"A", "B", "C"},
			},
			`"data"->>'name' NOT IN (?,?,?)`,
			[]interface{}{"A", "B", "C"},
			caseOrNotIn,
		},
		{
			"OR NOT IN (JSONB) (with func on field)",
			args{
				field:    "badge",
				operator: "NOT IN",
				value:    []string{"video", "audio"},
				funcs:    []string{"LOWER"},
			},
			`LOWER("data"->>'badge') NOT IN (?,?)`,
			[]interface{}{"video", "audio"},
			caseOrNotIn,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var (
				sql    string
				values interface{}
			)

			switch tt.caseType {
			case caseAndIn:
				sql, values = somesql.AndIn(tt.args.field, tt.args.operator, tt.args.value, tt.args.funcs...).AsSQL()
			case caseOrIn:
				sql, values = somesql.OrIn(tt.args.field, tt.args.operator, tt.args.value, tt.args.funcs...).AsSQL()
			case caseAndNotIn:
				sql, values = somesql.AndNotIn(tt.args.field, tt.args.operator, tt.args.value, tt.args.funcs...).AsSQL()
			case caseOrNotIn:
				sql, values = somesql.OrNotIn(tt.args.field, tt.args.operator, tt.args.value, tt.args.funcs...).AsSQL()
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

package somesql_test

import (
	"reflect"
	"testing"

	"github.com/fluxynet/somesql"
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
				field:    "name",
				operator: "IN",
				value:    []string{"A", "B", "C"},
			},
			" AND name IN (?,?,?)", //TODO field according to db
			[]interface{}{"A", "B", "C"},
			caseAndIn,
		},
		{
			"AND NOT IN",
			args{
				field:    "name",
				operator: "NOT IN",
				value:    []string{"A", "B", "C"},
			},
			" AND name NOT IN (?,?,?)", //TODO field according to db
			[]interface{}{"A", "B", "C"},
			caseAndNotIn,
		},
		{
			"OR IN",
			args{
				field:    "name",
				operator: "IN",
				value:    []string{"A", "B", "C"},
			},
			" OR name IN (?,?,?)", //TODO field according to db
			[]interface{}{"A", "B", "C"},
			caseOrIn,
		},
		{
			"OR NOT IN",
			args{
				field:    "name",
				operator: "IN",
				value:    []string{"A", "B", "C"},
			},
			" OR name NOT IN (?,?,?)", //TODO field according to db
			[]interface{}{"A", "B", "C"},
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

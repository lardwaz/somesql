package somesql_test

import (
	"reflect"
	"testing"

	"github.com/fluxynet/somesql"
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
				funcs:    []string{"LOWER"},
			},
			" AND LOWER(name)=?", //TODO field according to db
			[]interface{}{"John Doe"},
			caseAnd,
		},
		{
			"Or Normal",
			args{
				field:    "name",
				operator: "=",
				value:    "John Doe",
				funcs:    []string{"LOWER"},
			},
			" OR LOWER(name)=?", //TODO field according to db
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

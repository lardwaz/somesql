package somesql_test

import (
	"fmt"
	"testing"

	"github.com/lsldigital/somesql"
	"github.com/stretchr/testify/assert"
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
			`"data"->>'name'=?`,
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
			`"data"->>'name'=LOWER(?)`,
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
			`LOWER("data"->>'name')=?`,
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
			`LOWER("data"->>'name')=LOWER(?)`,
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
			`"data"->>'name'=?`,
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
			`"data"->>'name'=LOWER(?)`,
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
			`LOWER("data"->>'name')=?`,
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
			`LOWER("data"->>'name')=LOWER(?)`,
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
			`"data"->>'name'=?`,
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
			`LOWER("data"->>'name')=?`,
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
			`"data"->>'name'=LOWER(?)`,
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
			`LOWER("data"->>'name')=LOWER(?)`,
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
			`"data"->>'name'=?`,
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
			`"data"->>'name'=LOWER(?)`,
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
			`LOWER("data"->>'name')=?`,
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
			`LOWER("data"->>'name')=LOWER(?)`,
			[]interface{}{"John Doe"},
			caseOr,
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var (
				sql       string
				values    interface{}
				condition somesql.Condition
			)

			if tt.caseType == caseAnd {
				condition = somesql.And(tt.args.field, tt.args.operator, tt.args.value, tt.args.funcs...)
			} else {
				condition = somesql.Or(tt.args.field, tt.args.operator, tt.args.value, tt.args.funcs...)
			}

			sql, values = condition.AsSQL()

			assert.Equal(t, tt.sql, sql, fmt.Sprintf("%d: SQL invalid", i+1))
			assert.Equal(t, tt.values, values, fmt.Sprintf("%d: Values invalid", i+1))

			if tt.caseType == caseAnd {
				assert.Equal(t, somesql.AndCondition, condition.ConditionType(), fmt.Sprintf("%d: Condition type must be AND", i+1))
			} else {
				assert.Equal(t, somesql.OrCondition, condition.ConditionType(), fmt.Sprintf("%d: Condition type must be OR", i+1))
			}
		})
	}
}

package somesql_test

import (
	"fmt"
	"testing"

	"github.com/lsldigital/somesql"
	"github.com/stretchr/testify/assert"
)

func TestConditionGroup(t *testing.T) {
	const (
		caseAnd = iota
		caseOr
	)

	type testcase struct {
		name       string
		conditions []somesql.Condition
		sql        string
		values     []interface{}
		caseType   uint8
	}

	tests := []testcase{
		{
			"AND (1)",
			[]somesql.Condition{
				somesql.And(somesql.LangEN, "id", "=", "002fd6b1-f715-4875-838b-1546f27327df"),
				somesql.And(somesql.LangEN, "status", "=", true),
			},
			"(id=? AND status=?)",
			[]interface{}{"002fd6b1-f715-4875-838b-1546f27327df", true},
			caseAnd,
		},
		{
			"AND (2)",
			[]somesql.Condition{
				somesql.Or(somesql.LangEN, "id", "=", "002fd6b1-f715-4875-838b-1546f27327df"),
				somesql.And(somesql.LangEN, "status", "=", true),
			},
			"(id=? AND status=?)",
			[]interface{}{"002fd6b1-f715-4875-838b-1546f27327df", true},
			caseAnd,
		},
		{
			"OR (1)",
			[]somesql.Condition{
				somesql.Or(somesql.LangEN, "id", "=", "002fd6b1-f715-4875-838b-1546f27327df"),
				somesql.Or(somesql.LangEN, "status", "=", true),
			},
			"(id=? OR status=?)",
			[]interface{}{"002fd6b1-f715-4875-838b-1546f27327df", true},
			caseOr,
		},
		{
			"OR (2)",
			[]somesql.Condition{
				somesql.And(somesql.LangEN, "id", "=", "002fd6b1-f715-4875-838b-1546f27327df"),
				somesql.Or(somesql.LangEN, "status", "=", true),
			},
			"(id=? OR status=?)",
			[]interface{}{"002fd6b1-f715-4875-838b-1546f27327df", true},
			caseOr,
		},
		{
			"AND JSONB (fields in 'data' only)",
			[]somesql.Condition{
				somesql.And(somesql.LangEN, "badge", "=", "video"),
				somesql.And(somesql.LangEN, "has_video", "=", true),
			},
			`("data_en"->>'badge'=? AND ("data_en"->>'has_video')::BOOLEAN=?)`,
			[]interface{}{"video", true},
			caseAnd,
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var (
				sql       string
				values    interface{}
				condition somesql.Condition
			)

			for _, cond := range tt.conditions {
				if tt.caseType == caseAnd {
					condition = somesql.AndGroup(cond)
				} else {
					condition = somesql.OrGroup(cond)
				}
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

package somesql_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.lsl.digital/lardwaz/somesql"
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
				somesql.And("en", "id", "=", "002fd6b1-f715-4875-838b-1546f27327df"),
				somesql.And("en", "type", "=", "entityA"),
			},
			`("id" = ? AND "type" = ?)`,
			[]interface{}{"002fd6b1-f715-4875-838b-1546f27327df", "entityA"},
			caseAnd,
		},
		{
			"AND (2)",
			[]somesql.Condition{
				somesql.Or("en", "id", "=", "002fd6b1-f715-4875-838b-1546f27327df"),
				somesql.And("en", "type", "=", "entityA"),
			},
			`("id" = ? AND "type" = ?)`,
			[]interface{}{"002fd6b1-f715-4875-838b-1546f27327df", "entityA"},
			caseAnd,
		},
		{
			"AND (3)",
			[]somesql.Condition{
				somesql.And("en", "id", "=", "002fd6b1-f715-4875-838b-1546f27327df"),
				somesql.Or("en", "type", "=", true),
			},
			`("id" = ? OR ("type")::BOOLEAN = ?)`,
			[]interface{}{"002fd6b1-f715-4875-838b-1546f27327df", true},
			caseAnd,
		},
		{
			"OR (1)",
			[]somesql.Condition{
				somesql.Or("en", "id", "=", "002fd6b1-f715-4875-838b-1546f27327df"),
				somesql.Or("en", "type", "=", "entityA"),
			},
			`("id" = ? OR "type" = ?)`,
			[]interface{}{"002fd6b1-f715-4875-838b-1546f27327df", "entityA"},
			caseOr,
		},
		{
			"OR (2)",
			[]somesql.Condition{
				somesql.Or("en", "id", "=", "002fd6b1-f715-4875-838b-1546f27327df"),
				somesql.And("en", "type", "=", true),
			},
			`("id" = ? AND ("type")::BOOLEAN = ?)`,
			[]interface{}{"002fd6b1-f715-4875-838b-1546f27327df", true},
			caseOr,
		},
		{
			"OR (3)",
			[]somesql.Condition{
				somesql.And("en", "id", "=", "002fd6b1-f715-4875-838b-1546f27327df"),
				somesql.Or("en", "type", "=", true),
			},
			`("id" = ? OR ("type")::BOOLEAN = ?)`,
			[]interface{}{"002fd6b1-f715-4875-838b-1546f27327df", true},
			caseOr,
		},
		{
			"AND JSONB (fields in 'data' only)",
			[]somesql.Condition{
				somesql.And("en", "data.badge", "=", "video"),
				somesql.And("en", "data.has_video", "=", true),
			},
			`("data_en"->>'badge' = ? AND ("data_en"->>'has_video')::BOOLEAN = ?)`,
			[]interface{}{"video", true},
			caseAnd,
		},

		// relations
		{
			"AND JSONB with relations",
			[]somesql.Condition{
				somesql.And("en", "relations.tags", "=", "video"),
				somesql.And("en", "data.has_video", "=", true),
			},
			`((jsonb_path_exists("data_en", '$.tags[*] £ (@  == $val)', json_object(ARRAY['val', ?])::jsonb)) AND ("data_en"->>'has_video')::BOOLEAN = ?)`,
			[]interface{}{"video", true},
			caseAnd,
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var (
				sql            string
				values         interface{}
				conditionGroup somesql.Condition
			)

			if tt.caseType == caseAnd {
				conditionGroup = somesql.AndGroup(tt.conditions...)
			} else {
				conditionGroup = somesql.OrGroup(tt.conditions...)
			}

			sql, values = conditionGroup.AsSQL()

			assert.Equal(t, tt.sql, sql, fmt.Sprintf("%d: SQL invalid", i+1))
			assert.Equal(t, tt.values, values, fmt.Sprintf("%d: Values invalid", i+1))

			if tt.caseType == caseAnd {
				assert.Equal(t, somesql.AndCondition, conditionGroup.ConditionType(), fmt.Sprintf("%d: Condition type must be AND", i+1))
			} else {
				assert.Equal(t, somesql.OrCondition, conditionGroup.ConditionType(), fmt.Sprintf("%d: Condition type must be OR", i+1))
			}
		})
	}
}

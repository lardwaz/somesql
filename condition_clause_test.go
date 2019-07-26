package somesql_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.lsl.digital/lardwaz/somesql"
)

func TestConditionClause(t *testing.T) {
	const (
		caseAnd = iota
		caseOr
	)

	type args struct {
		lang     string
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
				lang:     "en",
				field:    "id",
				operator: "=",
				value:    "002fd6b1-f715-4875-838b-1546f27327df",
				funcs:    []string{},
			},
			`"id" = ?`,
			[]interface{}{"002fd6b1-f715-4875-838b-1546f27327df"},
			caseAnd,
		},
		{
			"AND Normal (boolean)",
			args{
				lang:     "en",
				field:    "type",
				operator: "=",
				value:    true,
				funcs:    []string{},
			},
			`("type")::BOOLEAN = ?`,
			[]interface{}{true},
			caseAnd,
		},
		{
			"AND with func on value",
			args{
				lang:     "en",
				field:    "type",
				operator: "=",
				value:    "Article",
				funcs:    []string{"", "LOWER"},
			},
			`"type" = LOWER(?)`,
			[]interface{}{"Article"},
			caseAnd,
		},
		{
			"AND with func on field",
			args{
				lang:     "en",
				field:    "type",
				operator: "=",
				value:    "article",
				funcs:    []string{"LOWER"},
			},
			`LOWER("type") = ?`,
			[]interface{}{"article"},
			caseAnd,
		},
		{
			"AND with func on field and value",
			args{
				lang:     "en",
				field:    "data.slug",
				operator: "=",
				value:    "Summertime-Beat-the-heat-and-stay-active",
				funcs:    []string{"LOWER", "LOWER"},
			},
			`LOWER("data_en"->>'slug') = LOWER(?)`,
			[]interface{}{"Summertime-Beat-the-heat-and-stay-active"},
			caseAnd,
		},
		{
			"AND JSONB",
			args{
				lang:     "en",
				field:    "data.name",
				operator: "=",
				value:    "John Doe",
				funcs:    []string{},
			},
			`"data_en"->>'name' = ?`,
			[]interface{}{"John Doe"},
			caseAnd,
		},
		{
			"AND JSONB (boolean)",
			args{
				lang:     "en",
				field:    "data.has_video",
				operator: "=",
				value:    true,
				funcs:    []string{},
			},
			`("data_en"->>'has_video')::BOOLEAN = ?`,
			[]interface{}{true},
			caseAnd,
		},
		{
			"AND JSONB with func on value",
			args{
				lang:     "en",
				field:    "data.name",
				operator: "=",
				value:    "John Doe",
				funcs:    []string{"", "LOWER"},
			},
			`"data_en"->>'name' = LOWER(?)`,
			[]interface{}{"John Doe"},
			caseAnd,
		},
		{
			"AND JSONB with func on field",
			args{
				lang:     "en",
				field:    "data.name",
				operator: "=",
				value:    "John Doe",
				funcs:    []string{"LOWER"},
			},
			`LOWER("data_en"->>'name') = ?`,
			[]interface{}{"John Doe"},
			caseAnd,
		},
		{
			"AND JSONB with func on field and value",
			args{
				lang:     "en",
				field:    "data.name",
				operator: "=",
				value:    "John Doe",
				funcs:    []string{"LOWER", "LOWER"},
			},
			`LOWER("data_en"->>'name') = LOWER(?)`,
			[]interface{}{"John Doe"},
			caseAnd,
		},
		{
			"AND JSONB on relations",
			args{
				lang:  "en",
				field: "relations.name",
				value: "John Doe",
			},
			`("data_en" @> '{"name":?}'::JSONB)`,
			[]interface{}{"John Doe"},
			caseAnd,
		},

		{
			"OR Normal",
			args{
				lang:     "en",
				field:    "id",
				operator: "=",
				value:    "002fd6b1-f715-4875-838b-1546f27327df",
				funcs:    []string{},
			},
			`"id" = ?`,
			[]interface{}{"002fd6b1-f715-4875-838b-1546f27327df"},
			caseOr,
		},
		{
			"OR with func on value",
			args{
				lang:     "en",
				field:    "type",
				operator: "=",
				value:    "Article",
				funcs:    []string{"", "LOWER"},
			},
			`"type" = LOWER(?)`,
			[]interface{}{"Article"},
			caseOr,
		},
		{
			"OR with func on field",
			args{
				lang:     "en",
				field:    "type",
				operator: "=",
				value:    "article",
				funcs:    []string{"LOWER"},
			},
			`LOWER("type") = ?`,
			[]interface{}{"article"},
			caseOr,
		},
		{
			"OR with func on field and value",
			args{
				lang:     "en",
				field:    "data.slug",
				operator: "=",
				value:    "Summertime-Beat-the-heat-and-stay-active",
				funcs:    []string{"LOWER", "LOWER"},
			},
			`LOWER("data_en"->>'slug') = LOWER(?)`,
			[]interface{}{"Summertime-Beat-the-heat-and-stay-active"},
			caseOr,
		},
		{
			"OR JSONB",
			args{
				lang:     "en",
				field:    "data.name",
				operator: "=",
				value:    "John Doe",
				funcs:    []string{},
			},
			`"data_en"->>'name' = ?`,
			[]interface{}{"John Doe"},
			caseOr,
		},
		{
			"OR JSONB with func on value",
			args{
				lang:     "en",
				field:    "data.name",
				operator: "=",
				value:    "John Doe",
				funcs:    []string{"", "LOWER"},
			},
			`"data_en"->>'name' = LOWER(?)`,
			[]interface{}{"John Doe"},
			caseOr,
		},
		{
			"OR JSONB with func on field",
			args{
				lang:     "en",
				field:    "data.name",
				operator: "=",
				value:    "John Doe",
				funcs:    []string{"LOWER"},
			},
			`LOWER("data_en"->>'name') = ?`,
			[]interface{}{"John Doe"},
			caseOr,
		},
		{
			"OR JSONB with func on field and value",
			args{
				lang:     "en",
				field:    "data.name",
				operator: "=",
				value:    "John Doe",
				funcs:    []string{"LOWER", "LOWER"},
			},
			`LOWER("data_en"->>'name') = LOWER(?)`,
			[]interface{}{"John Doe"},
			caseOr,
		},
		{
			"OR JSONB on relations",
			args{
				lang:  "en",
				field: "relations.name",
				value: "John Doe",
			},
			`("data_en" @> '{"name":?}'::JSONB)`,
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
				condition = somesql.And(tt.args.lang, tt.args.field, tt.args.operator, tt.args.value, tt.args.funcs...)
			} else {
				condition = somesql.Or(tt.args.lang, tt.args.field, tt.args.operator, tt.args.value, tt.args.funcs...)
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

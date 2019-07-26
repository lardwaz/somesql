package somesql_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.lsl.digital/lardwaz/somesql"
)

func TestConditionIn(t *testing.T) {
	const (
		caseAndIn = iota
		caseAndNotIn
		caseOrIn
		caseOrNotIn
	)

	type args struct {
		lang  string
		field string
		value interface{}
		funcs []string
	}

	type testcase struct {
		name     string
		args     args
		sql      string
		values   interface{}
		caseType uint8
	}

	tests := []testcase{
		{
			"AND IN",
			args{
				lang:  "en",
				field: "id",
				value: []string{"A", "B", "C"},
			},
			`"id" IN (?,?,?)`,
			[]interface{}{"A", "B", "C"},
			caseAndIn,
		},
		{
			"AND IN (boolean values)",
			args{
				lang:  "en",
				field: "type",
				value: []bool{true, false},
			},
			`"type" IN (?,?)`,
			[]interface{}{true, false},
			caseAndIn,
		},
		{
			"AND IN (integer values)",
			args{
				lang:  "en",
				field: "type",
				value: []int{1, 2},
			},
			`"type" IN (?,?)`,
			[]interface{}{1, 2},
			caseAndIn,
		},
		{
			"AND IN (with func on field)",
			args{
				lang:  "en",
				field: "updated_at",
				value: []string{"2019"},
				funcs: []string{"YEAR"},
			},
			`YEAR("updated_at") IN (?)`,
			[]interface{}{"2019"},
			caseAndIn,
		},
		{
			"AND NOT IN",
			args{
				lang:  "en",
				field: "id",
				value: []string{"A", "B", "C"},
			},
			`"id" NOT IN (?,?,?)`,
			[]interface{}{"A", "B", "C"},
			caseAndNotIn,
		},
		{
			"AND NOT IN (with func on field)",
			args{
				lang:  "en",
				field: "updated_at",
				value: []string{"2016"},
				funcs: []string{"YEAR"},
			},
			`YEAR("updated_at") NOT IN (?)`,
			[]interface{}{"2016"},
			caseAndNotIn,
		},
		{
			"AND IN (JSONB)",
			args{
				lang:  "en",
				field: "data.name",
				value: []string{"A", "B"},
			},
			`("data_en" @> '{"name":["?"]}'::JSONB OR '{"name":["?"]}'::JSONB)`,
			[]interface{}{"A", "B"},
			caseAndIn,
		},
		{
			"AND NOT IN (JSONB)",
			args{
				lang:  "en",
				field: "data.name",
				value: []string{"A", "B"},
			},
			`NOT("data_en" @> '{"name":["?"]}'::JSONB OR '{"name":["?"]}'::JSONB)`,
			[]interface{}{"A", "B"},
			caseAndNotIn,
		},
		{
			"AND IN (JSONB) (with func on field)",
			args{
				lang:  "en",
				field: "data.badge",
				value: []string{"video", "audio"},
				funcs: []string{"LOWER"},
			},
			`("data_en" @> '{"badge":["?"]}'::JSONB OR '{"badge":["?"]}'::JSONB)`,
			[]interface{}{"video", "audio"},
			caseAndIn,
		},
		{
			"AND NOT IN (JSONB) (with func on field)",
			args{
				lang:  "en",
				field: "data.badge",
				value: []string{"video", "audio"},
				funcs: []string{"LOWER"},
			},
			`NOT("data_en" @> '{"badge":["?"]}'::JSONB OR '{"badge":["?"]}'::JSONB)`,
			[]interface{}{"video", "audio"},
			caseAndNotIn,
		},
		{
			"AND IN (JSONB) - relations",
			args{
				lang:  "en",
				field: "relations.name",
				value: []string{"A", "B"},
			},
			`("data_en" @> '{"name":["?"]}'::JSONB OR '{"name":["?"]}'::JSONB)`,
			[]interface{}{"A", "B"},
			caseAndIn,
		},
		{
			"AND NOT IN (JSONB) - relations",
			args{
				lang:  "en",
				field: "relations.name",
				value: []string{"A", "B"},
			},
			`NOT("data_en" @> '{"name":["?"]}'::JSONB OR '{"name":["?"]}'::JSONB)`,
			[]interface{}{"A", "B"},
			caseAndNotIn,
		},
		{
			"AND IN (JSONB) (with func on field) - relations",
			args{
				lang:  "en",
				field: "relations.badge",
				value: []string{"video", "audio"},
				funcs: []string{"LOWER"},
			},
			`("data_en" @> '{"badge":["?"]}'::JSONB OR '{"badge":["?"]}'::JSONB)`,
			[]interface{}{"video", "audio"},
			caseAndIn,
		},
		{
			"AND NOT IN (JSONB) (with func on field) - relations",
			args{
				lang:  "en",
				field: "relations.badge",
				value: []string{"video", "audio"},
				funcs: []string{"LOWER"},
			},
			`NOT("data_en" @> '{"badge":["?"]}'::JSONB OR '{"badge":["?"]}'::JSONB)`,
			[]interface{}{"video", "audio"},
			caseAndNotIn,
		},

		{
			"OR IN",
			args{
				lang:  "en",
				field: "id",
				value: []string{"A", "B", "C"},
			},
			`"id" IN (?,?,?)`,
			[]interface{}{"A", "B", "C"},
			caseOrIn,
		},
		{
			"OR IN (with func on field)",
			args{
				lang:  "en",
				field: "updated_at",
				value: []string{"2019"},
				funcs: []string{"YEAR"},
			},
			`YEAR("updated_at") IN (?)`,
			[]interface{}{"2019"},
			caseOrIn,
		},
		{
			"OR NOT IN",
			args{
				lang:  "en",
				field: "id",
				value: []string{"A", "B", "C"},
			},
			`"id" NOT IN (?,?,?)`,
			[]interface{}{"A", "B", "C"},
			caseOrNotIn,
		},
		{
			"OR NOT IN (with func on field)",
			args{
				lang:  "en",
				field: "updated_at",
				value: []string{"2015"},
				funcs: []string{"YEAR"},
			},
			`YEAR("updated_at") NOT IN (?)`,
			[]interface{}{"2015"},
			caseOrNotIn,
		},
		{
			"OR IN (JSONB)",
			args{
				lang:  "en",
				field: "data.name",
				value: []string{"A", "B"},
			},
			`("data_en" @> '{"name":["?"]}'::JSONB OR '{"name":["?"]}'::JSONB)`,
			[]interface{}{"A", "B"},
			caseOrIn,
		},
		{
			"OR NOT IN (JSONB)",
			args{
				lang:  "en",
				field: "data.name",
				value: []string{"A", "B"},
			},
			`NOT("data_en" @> '{"name":["?"]}'::JSONB OR '{"name":["?"]}'::JSONB)`,
			[]interface{}{"A", "B"},
			caseOrNotIn,
		},
		{
			"OR IN (JSONB) (with func on field)",
			args{
				lang:  "en",
				field: "data.badge",
				value: []string{"video", "audio"},
				funcs: []string{"LOWER"},
			},
			`("data_en" @> '{"badge":["?"]}'::JSONB OR '{"badge":["?"]}'::JSONB)`,
			[]interface{}{"video", "audio"},
			caseOrIn,
		},
		{
			"OR NOT IN (JSONB) (with func on field)",
			args{
				lang:  "en",
				field: "data.badge",
				value: []string{"video", "audio"},
				funcs: []string{"LOWER"},
			},
			`NOT("data_en" @> '{"badge":["?"]}'::JSONB OR '{"badge":["?"]}'::JSONB)`,
			[]interface{}{"video", "audio"},
			caseOrNotIn,
		},
		{
			"OR IN (JSONB) - relations",
			args{
				lang:  "en",
				field: "relations.name",
				value: []string{"A", "B"},
			},
			`("data_en" @> '{"name":["?"]}'::JSONB OR '{"name":["?"]}'::JSONB)`,
			[]interface{}{"A", "B"},
			caseOrIn,
		},
		{
			"OR NOT IN (JSONB) - relations",
			args{
				lang:  "en",
				field: "relations.name",
				value: []string{"A", "B"},
			},
			`NOT("data_en" @> '{"name":["?"]}'::JSONB OR '{"name":["?"]}'::JSONB)`,
			[]interface{}{"A", "B"},
			caseOrNotIn,
		},
		{
			"OR IN (JSONB) (with func on field) - relations",
			args{
				lang:  "en",
				field: "relations.badge",
				value: []string{"video", "audio"},
				funcs: []string{"LOWER"},
			},
			`("data_en" @> '{"badge":["?"]}'::JSONB OR '{"badge":["?"]}'::JSONB)`,
			[]interface{}{"video", "audio"},
			caseOrIn,
		},
		{
			"OR NOT IN (JSONB) (with func on field) - relations",
			args{
				lang:  "en",
				field: "relations.badge",
				value: []string{"video", "audio"},
				funcs: []string{"LOWER"},
			},
			`NOT("data_en" @> '{"badge":["?"]}'::JSONB OR '{"badge":["?"]}'::JSONB)`,
			[]interface{}{"video", "audio"},
			caseOrNotIn,
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var (
				sql    string
				values interface{}
			)

			switch tt.caseType {
			case caseAndIn:
				sql, values = somesql.AndIn(tt.args.lang, tt.args.field, tt.args.value, tt.args.funcs...).AsSQL()
			case caseOrIn:
				sql, values = somesql.OrIn(tt.args.lang, tt.args.field, tt.args.value, tt.args.funcs...).AsSQL()
			case caseAndNotIn:
				sql, values = somesql.AndNotIn(tt.args.lang, tt.args.field, tt.args.value, tt.args.funcs...).AsSQL()
			case caseOrNotIn:
				sql, values = somesql.OrNotIn(tt.args.lang, tt.args.field, tt.args.value, tt.args.funcs...).AsSQL()
			}

			assert.Equal(t, tt.sql, sql, fmt.Sprintf("%d: SQL invalid", i+1))
			assert.Equal(t, tt.values, values, fmt.Sprintf("%d: Values invalid", i+1))
		})
	}
}

package somesql_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.lsl.digital/lardwaz/somesql"
)

func TestConditionQuery(t *testing.T) {
	const (
		caseAndIn = iota
		caseAndNotIn
		caseOrIn
		caseOrNotIn
	)

	type testCase struct {
		name      string
		fieldName string
		query     *somesql.Select
		sql       string
		values    []interface{}
		caseType  uint8
		lang      string
	}

	tests := []testCase{
		// Error: cannot have more than 1 field in subquery
		{
			name:      "AndInQuery [error 1]",
			fieldName: "author_id",
			query:     somesql.NewSelectInner("").Fields("type", "data.slug").Where(somesql.And("en", "id", "=", "002fd6b1-f715-4875-838b-1546f27327df")),
			sql:       `"data_en"->>'author_id' IN (SELECT "type", "data_"->>'slug' "slug" FROM repo WHERE "id" = ? LIMIT 10)`,
			values:    []interface{}{"002fd6b1-f715-4875-838b-1546f27327df"},
			caseType:  caseAndIn,
			lang:      "en",
		},
		{
			name:      "OrInQuery [error 2]",
			fieldName: "author_id",
			query:     somesql.NewSelectInner("").Fields("type", "data.slug").Where(somesql.And("en", "id", "=", "002fd6b1-f715-4875-838b-1546f27327df")),
			sql:       `"data_en"->>'author_id' IN (SELECT "type", "data_"->>'slug' "slug" FROM repo WHERE "id" = ? LIMIT 10)`,
			values:    []interface{}{"002fd6b1-f715-4875-838b-1546f27327df"},
			caseType:  caseOrIn,
			lang:      "en",
		},
		{
			name:      "AndInQuery [1]",
			fieldName: "author_id",
			query:     somesql.NewSelectInner("en").Fields("data.author_id").Where(somesql.And("en", "id", "=", "1")),
			sql:       `"data_en"->>'author_id' IN (SELECT "data_en"->>'author_id' "author_id" FROM repo WHERE "id" = ? LIMIT 10)`,
			values:    []interface{}{"1"},
			caseType:  caseAndIn,
			lang:      "en",
		},
		{
			name:      "AndInQuery [2]",
			fieldName: "author_id",
			query:     somesql.NewSelectInner("en").Fields("data.author_id").Where(somesql.And("en", "id", "=", "1")).Where(somesql.And("en", "type", "=", "entityA")),
			sql:       `"data_en"->>'author_id' IN (SELECT "data_en"->>'author_id' "author_id" FROM repo WHERE "id" = ? AND "type" = ? LIMIT 10)`,
			values:    []interface{}{"1", "entityA"},
			caseType:  caseAndIn,
			lang:      "en",
		},
		{
			name:      "AndNotInQuery",
			fieldName: "author_id",
			query:     somesql.NewSelectInner("en").Fields("data.author_id").Where(somesql.And("en", "id", "=", "1")),
			sql:       `"data_en"->>'author_id' NOT IN (SELECT "data_en"->>'author_id' "author_id" FROM repo WHERE "id" = ? LIMIT 10)`,
			values:    []interface{}{"1"},
			caseType:  caseAndNotIn,
			lang:      "en",
		},
		{
			name:      "OrInQuery",
			fieldName: "author_id",
			query:     somesql.NewSelectInner("en").Fields("data.author_id").Where(somesql.And("en", "id", "=", "1")),
			sql:       `"data_en"->>'author_id' IN (SELECT "data_en"->>'author_id' "author_id" FROM repo WHERE "id" = ? LIMIT 10)`,
			values:    []interface{}{"1"},
			caseType:  caseOrIn,
			lang:      "en",
		},
		{
			name:      "OrNotInQuery",
			fieldName: "author_id",
			query:     somesql.NewSelectInner("en").Fields("data.author_id").Where(somesql.And("en", "id", "=", "1")),
			sql:       `"data_en"->>'author_id' NOT IN (SELECT "data_en"->>'author_id' "author_id" FROM repo WHERE "id" = ? LIMIT 10)`,
			values:    []interface{}{"1"},
			caseType:  caseOrNotIn,
			lang:      "en",
		},
		{
			name:      "OrNotInQuery",
			fieldName: "author_id",
			query:     somesql.NewSelectInner("en").Fields("data.author_id").Where(somesql.And("en", "relations.tags", "", "video")),
			sql:       `"data_en"->>'author_id' NOT IN (SELECT "data_en"->>'author_id' "author_id" FROM repo WHERE ("data_en" @> '{"tags":?}'::JSONB) LIMIT 10)`,
			values:    []interface{}{"video"},
			caseType:  caseOrNotIn,
			lang:      "en",
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var (
				cQuery somesql.ConditionQuery
				sql    string
				values interface{}
			)

			switch tt.caseType {
			case caseAndIn:
				cQuery = somesql.AndInQuery(tt.lang, tt.fieldName, tt.query)
			case caseAndNotIn:
				cQuery = somesql.AndNotInQuery(tt.lang, tt.fieldName, tt.query)
			case caseOrIn:
				cQuery = somesql.OrInQuery(tt.lang, tt.fieldName, tt.query)
			case caseOrNotIn:
				cQuery = somesql.OrNotInQuery(tt.lang, tt.fieldName, tt.query)
			}

			sql, values = cQuery.AsSQL()

			assert.Equal(t, tt.sql, sql, fmt.Sprintf("%d: SQL invalid", i+1))
			assert.Equal(t, tt.values, values, fmt.Sprintf("%d: Values invalid", i+1))

			switch tt.caseType {
			case caseAndIn, caseAndNotIn:
				assert.Equal(t, somesql.AndCondition, cQuery.ConditionType(), fmt.Sprintf("%d: Condition type must be AND", i+1))
			case caseOrIn, caseOrNotIn:
				assert.Equal(t, somesql.OrCondition, cQuery.ConditionType(), fmt.Sprintf("%d: Condition type must be OR", i+1))
			}
		})
	}
}

package somesql_test

import (
	"reflect"
	"testing"

	"github.com/fluxynet/somesql"
	"github.com/fluxynet/somesql/postgres"
)

func TestConditionQuery(t *testing.T) {
	const (
		caseAndIn = iota
		caseAndNotIn
		caseOrIn
		caseOrNotIn
	)

	type testcase struct {
		name      string
		fieldname string
		query     somesql.Query
		sql       string
		values    []interface{}
		caseType  uint8
	}

	tests := []testcase{
		{
			"AND IN",
			"foo",
			postgres.New().Select("a", "b").Where(somesql.And("bar", "=", "baz")),
			" AND name IN (...)", //TODO write subquery
			[]interface{}{"baz"},
			caseAndIn,
		},
		{
			"AND NOT IN",
			"foo",
			postgres.New().Select("a", "b").Where(somesql.And("bar", "=", "baz")),
			" AND name NOT IN (...)", //TODO write subquery
			[]interface{}{"baz"},
			caseAndNotIn,
		},
		{
			"OR IN",
			"foo",
			postgres.New().Select("a", "b").Where(somesql.And("bar", "=", "baz")),
			" OR name IN (...)", //TODO write subquery
			[]interface{}{"baz"},
			caseOrIn,
		},
		{
			"OR NOT IN",
			"foo",
			postgres.New().Select("a", "b").Where(somesql.And("bar", "=", "baz")),
			" OR name NOT IN (...)", //TODO write subquery
			[]interface{}{"baz"},
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
				sql, values = somesql.AndInQuery(tt.fieldname, tt.query).AsSQL()
			case caseAndNotIn:
				sql, values = somesql.OrInQuery(tt.fieldname, tt.query).AsSQL()
			case caseOrIn:
				sql, values = somesql.AndNotInQuery(tt.fieldname, tt.query).AsSQL()
			case caseOrNotIn:
				sql, values = somesql.OrNotInQuery(tt.fieldname, tt.query).AsSQL()

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

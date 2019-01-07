package somesql_test

import (
	"reflect"
	"testing"

	"github.com/lsldigital/somesql"
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
			"AND",
			[]somesql.Condition{},        //TODO write conditions
			" AND (x AND y)",             //TODO write conditions
			[]interface{}{"A", "B", "C"}, //TODO write values
			caseAnd,
		},
		{
			"OR",
			[]somesql.Condition{},        //TODO write conditions
			" AND (x AND y)",             //TODO write conditions
			[]interface{}{"A", "B", "C"}, //TODO write values
			caseOr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var (
				sql    string
				values interface{}
			)

			switch tt.caseType {
			case caseAnd:
				sql, values = somesql.AndGroup(tt.conditions...).AsSQL()
			case caseOr:
				sql, values = somesql.OrGroup(tt.conditions...).AsSQL()
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

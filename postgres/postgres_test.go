package postgres

import (
	"fmt"
	"testing"

	"github.com/lsldigital/somesql"
	"github.com/stretchr/testify/assert"
)

func TestQuery_AsSQL_Fields(t *testing.T) {
	type testcase struct {
		name        string
		query       Query
		expectedSQL string
	}

	tests := []testcase{
		// Select ALL
		{
			"SELECT *",
			New(),
			"SELECT id, created_at, updated_at, owner_id, status, type, data_en FROM repo",
		},
		{
			"SELECT * (langEN)",
			New().SetLang(somesql.LangEN),
			"SELECT id, created_at, updated_at, owner_id, status, type, data_en FROM repo",
		},
		{
			"SELECT * (langFR)",
			New().SetLang(somesql.LangFR),
			"SELECT id, created_at, updated_at, owner_id, status, type, data_fr FROM repo",
		},
		// Select some pre-defined fields
		{
			"SELECT id, type, data",
			New().Select("id", "type", "data"),
			`SELECT id, type, data_en FROM repo`,
		},
		{
			"SELECT id, type, data (LangEN)",
			New().Select("id", "type", "data").SetLang(somesql.LangEN),
			`SELECT id, type, data_en FROM repo`,
		},
		{
			"SELECT id, type, data (LangFR)",
			New().Select("id", "type", "data").SetLang(somesql.LangFR),
			`SELECT id, type, data_fr FROM repo`,
		},
		// Select pre-defined fields and json attributes ('data_en'/'data_fr') from data_*
		{
			"SELECT id, type, data_en",
			New().Select("id", "type", "data_en"),
			`SELECT id, type, json_build_object('data_en', data_en->'data_en') "data" FROM repo`,
		},
		{
			"SELECT id, type, data_en (LangEN)",
			New().Select("id", "type", "data_en").SetLang(somesql.LangEN),
			`SELECT id, type, json_build_object('data_en', data_en->'data_en') "data" FROM repo`,
		},
		{
			"SELECT id, type, data_en (LangFR)",
			New().Select("id", "type", "data_en").SetLang(somesql.LangFR),
			`SELECT id, type, json_build_object('data_en', data_fr->'data_en') "data" FROM repo`,
		},
		{
			"SELECT id, type, data_fr",
			New().Select("id", "type", "data_fr"),
			`SELECT id, type, json_build_object('data_fr', data_en->'data_fr') "data" FROM repo`,
		},
		{
			"SELECT id, type, data_fr (LangEN)",
			New().Select("id", "type", "data_fr").SetLang(somesql.LangEN),
			`SELECT id, type, json_build_object('data_fr', data_en->'data_fr') "data" FROM repo`,
		},
		{
			"SELECT id, type, data_fr (LangFR)",
			New().Select("id", "type", "data_fr").SetLang(somesql.LangFR),
			`SELECT id, type, json_build_object('data_fr', data_fr->'data_fr') "data" FROM repo`,
		},
		// Select pre-defined fields and json attributes (any other) from data_*
		{
			"SELECT id, type, data_en->'body'",
			New().Select("id", "type", "body"),
			`SELECT id, type, json_build_object('body', data_en->'body') "data" FROM repo`,
		},
		{
			"SELECT id, type, data_en->'body' (LangEN)",
			New().Select("id", "type", "body").SetLang(somesql.LangEN),
			`SELECT id, type, json_build_object('body', data_en->'body') "data" FROM repo`,
		},
		{
			"SELECT id, type, data_fr->'body' (LangFR)",
			New().Select("id", "type", "body").SetLang(somesql.LangFR),
			`SELECT id, type, json_build_object('body', data_fr->'body') "data" FROM repo`,
		},
		// Select pre-defined fields and json attributes (any other + compound) from data_*
		{
			"SELECT id, type, data_en->'body', data_en->'author_id'",
			New().Select("id", "type", "body", "author_id"),
			`SELECT id, type, json_build_object('body', data_en->'body', 'author_id', data_en->'author_id') "data" FROM repo`,
		},
		{
			"SELECT id, type, data_en->'body', data_en->'author_id' (LangEN)",
			New().Select("id", "type", "body", "author_id").SetLang(somesql.LangEN),
			`SELECT id, type, json_build_object('body', data_en->'body', 'author_id', data_en->'author_id') "data" FROM repo`,
		},
		{
			"SELECT id, type, data_fr->'body', data_fr->'author_id' (LangFR)",
			New().Select("id", "type", "body", "author_id").SetLang(somesql.LangFR),
			`SELECT id, type, json_build_object('body', data_fr->'body', 'author_id', data_fr->'author_id') "data" FROM repo`,
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSQL, _ := tt.query.AsSQL()
			assert.Equal(t, tt.expectedSQL, gotSQL, fmt.Sprintf("Fields %03d :: %s", i+1, tt.name))
		})
	}
}

func TestQuery_AsSQL_ConditionClause(t *testing.T) {
	type testcase struct {
		name           string
		query          Query
		expectedSQL    string
		expectedValues []interface{}
	}

	// SELECT _ WHERE c1
	// SELECT _ WHERE c1 AND c2
	// SELECT _ WHERE c1 OR c2
	// SELECT _ WHERE c1 AND c2 OR c3

	tests := []testcase{}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSQL, gotValues := tt.query.AsSQL()

			assert.Equal(t, tt.expectedSQL, gotSQL, fmt.Sprintf("Fields %03d :: invalid sql :: %s", i+1, tt.name))
			assert.Equal(t, tt.expectedValues, gotValues, fmt.Sprintf("Fields %03d :: invalid values :: %s", i+1, tt.name))
		})
	}
}

func TestQuery_AsSQL_ConditionGroup(t *testing.T) {
	type testcase struct {
		name           string
		query          Query
		expectedSQL    string
		expectedValues []interface{}
	}

	// SELECT _ WHERE c1 AND (c2 OR c3)
	// SELECT _ WHERE (c1 AND c2) OR c3

	tests := []testcase{}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSQL, gotValues := tt.query.AsSQL()

			assert.Equal(t, tt.expectedSQL, gotSQL, fmt.Sprintf("Fields %03d :: invalid sql :: %s", i+1, tt.name))
			assert.Equal(t, tt.expectedValues, gotValues, fmt.Sprintf("Fields %03d :: invalid values :: %s", i+1, tt.name))
		})
	}
}

func TestQuery_AsSQL_ConditionIN(t *testing.T) {
	type testcase struct {
		name           string
		query          Query
		expectedSQL    string
		expectedValues []interface{}
	}

	// SELECT _ WHERE c1 IN (...)
	// SELECT _ WHERE c1 NOT IN (...)
	// SELECT _ WHERE c1 IN (...) AND c2 NOT IN (...)
	// SELECT _ WHERE c1 IN (...) OR c2 NOT IN (...)
	// SELECT _ WHERE (c1 IN (...) AND c2 NOT IN (...)) AND c3
	// SELECT _ WHERE c3 OR (c1 IN (...) AND c2 NOT IN (...))

	tests := []testcase{}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSQL, gotValues := tt.query.AsSQL()

			assert.Equal(t, tt.expectedSQL, gotSQL, fmt.Sprintf("Fields %03d :: invalid sql :: %s", i+1, tt.name))
			assert.Equal(t, tt.expectedValues, gotValues, fmt.Sprintf("Fields %03d :: invalid values :: %s", i+1, tt.name))
		})
	}
}

package postgres

import (
	"fmt"
	"testing"

	"github.com/lsldigital/somesql"
	"github.com/stretchr/testify/assert"
)

func TestQuery_AsSQL_Fields(t *testing.T) {
	type testCase struct {
		name        string
		query       Query
		expectedSQL string
	}

	tests := []testCase{
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
	type testCase struct {
		name           string
		query          Query
		expectedSQL    string
		expectedValues []interface{}
	}

	tests := []testCase{
		// Regular condition clauses on pre-defined fields
		{
			"WHERE id=?",
			New().Select("data").Where(somesql.And("id", "=", "1")),
			"SELECT data_en FROM repo WHERE id=?",
			[]interface{}{"1"},
		},
		{
			"WHERE id=? AND status=?",
			New().Select("data").Where(somesql.And("id", "=", "1")).Where(somesql.And("status", "=", "published")),
			"SELECT data_en FROM repo WHERE id=? AND status=?",
			[]interface{}{"1", "published"},
		},
		{
			"WHERE id=? OR status=?",
			New().Select("data").Where(somesql.And("id", "=", "1")).Where(somesql.Or("status", "=", "published")),
			"SELECT data_en FROM repo WHERE id=? OR status=?",
			[]interface{}{"1", "published"},
		},
		{
			"WHERE id=? AND status=? OR type=?",
			New().Select("data").Where(somesql.And("id", "=", "1")).Where(somesql.And("status", "=", "published")).Where(somesql.Or("type", "=", "article")),
			"SELECT data_en FROM repo WHERE id=? AND status=? OR type=?",
			[]interface{}{"1", "published", "article"},
		},
		// Regular condition clauses on json attributes
		{
			"WHERE data_en->>'author_id'=?",
			New().Select("data").Where(somesql.And("author_id", "=", "1")),
			"SELECT data FROM repo WHERE data_en->>'author_id'=?",
			[]interface{}{"1"},
		},
		{
			"WHERE data_fr->>'author_id'=? (langFR)",
			New().Select("data").SetLang(somesql.LangFR).Where(somesql.And("author_id", "=", "1")),
			"SELECT data FROM repo WHERE data_fr->>'author_id'=?",
			[]interface{}{"1"},
		},
		{
			"WHERE data_fr->>'author_id'=? OR data_fr->>'category_id'=? (langFR)",
			New().Select("data").SetLang(somesql.LangFR).Where(somesql.And("author_id", "=", "1")).Where(somesql.Or("category_id", "=", "2")),
			"SELECT data FROM repo WHERE data_fr->>'author_id'=? OR data_fr->>'category_id'=?",
			[]interface{}{"1", "2"},
		},
		{
			"WHERE data_fr->>'author_id'=? AND data_fr->>'category_id'=? (langFR)",
			New().Select("data").SetLang(somesql.LangFR).Where(somesql.And("author_id", "=", "1")).Where(somesql.And("category_id", "=", "2")),
			"SELECT data FROM repo WHERE data_fr->>'author_id'=? AND data_fr->>'category_id'=?",
			[]interface{}{"1", "2"},
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSQL, gotValues := tt.query.AsSQL()

			assert.Equal(t, tt.expectedSQL, gotSQL, fmt.Sprintf("Fields %03d :: invalid sql :: %s", i+1, tt.name))
			assert.Equal(t, tt.expectedValues, gotValues, fmt.Sprintf("Fields %03d :: invalid values :: %s", i+1, tt.name))
		})
	}
}

func TestQuery_AsSQL_ConditionGroup(t *testing.T) {
	type testCase struct {
		name           string
		query          Query
		expectedSQL    string
		expectedValues []interface{}
	}

	// SELECT _ WHERE c1 AND (c2 OR c3)
	// SELECT _ WHERE (c1 AND c2) OR c3

	tests := []testCase{}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSQL, gotValues := tt.query.AsSQL()

			assert.Equal(t, tt.expectedSQL, gotSQL, fmt.Sprintf("Fields %03d :: invalid sql :: %s", i+1, tt.name))
			assert.Equal(t, tt.expectedValues, gotValues, fmt.Sprintf("Fields %03d :: invalid values :: %s", i+1, tt.name))
		})
	}
}

func TestQuery_AsSQL_ConditionIN(t *testing.T) {
	type testCase struct {
		name           string
		query          Query
		expectedSQL    string
		expectedValues []interface{}
	}

	// SELECT _ WHERE c1 IN (...) AND c2 NOT IN (...)
	// SELECT _ WHERE c1 IN (...) OR c2 NOT IN (...)

	// SELECT _ WHERE (c1 IN (...) AND c2 NOT IN (...)) AND c3
	// SELECT _ WHERE c3 OR (c1 IN (...) AND c2 NOT IN (...))

	tests := []testCase{
		{
			"WHERE id IN (...) - predefined field",
			New().Select("data").Where(somesql.AndIn("id", []string{"A", "B", "C"})),
			"SELECT data_en FROM repo WHERE id IN (?,?,?)",
			[]interface{}{"A", "B", "C"},
		},
		{
			"WHERE field IN (...) - JSONB",
			New().Where(somesql.AndIn("name", []string{"A", "B", "C"})),
			`SELECT id, created_at, updated_at, owner_id, status, type, data_en FROM repo WHERE "data_en"->>'name' IN (?,?,?)`,
			[]interface{}{"A", "B", "C"},
		},
		{
			"WHERE FUNC(field) IN (...) - predefined field",
			New().Select("data").Where(somesql.AndIn("updated_at", []string{"2019"}, "YEAR")),
			"SELECT data_en FROM repo WHERE YEAR(updated_at) IN (?)",
			[]interface{}{"2019"},
		},
		{
			"WHERE FUNC(field) IN (...) - JSONB",
			New().Where(somesql.AndIn("name", []string{"a"}, "LOWER")),
			`SELECT id, created_at, updated_at, owner_id, status, type, data_en FROM repo WHERE LOWER("data_en"->>'name') IN (?)`,
			[]interface{}{"a"},
		},

		{
			"WHERE id NOT IN (...) - predefined field",
			New().Select("data").Where(somesql.AndNotIn("id", []string{"A", "B", "C"})),
			"SELECT data_en FROM repo WHERE id NOT IN (?,?,?)",
			[]interface{}{"A", "B", "C"},
		},
		{
			"WHERE FUNC(field) NOT IN (...) - predefined field",
			New().Select("data").Where(somesql.AndNotIn("updated_at", []string{"2019"}, "YEAR")),
			"SELECT data_en FROM repo WHERE YEAR(updated_at) NOT IN (?)",
			[]interface{}{"2019"},
		},
		{
			"WHERE field NOT IN (...) - JSONB",
			New().Where(somesql.AndNotIn("name", []string{"A", "B", "C"})),
			`SELECT id, created_at, updated_at, owner_id, status, type, data_en FROM repo WHERE "data_en"->>'name' NOT IN (?,?,?)`,
			[]interface{}{"A", "B", "C"},
		},
		{
			"WHERE FUNC(field) NOT IN (...) - JSONB",
			New().Where(somesql.AndIn("name", []string{"a"}, "LOWER")),
			`SELECT id, created_at, updated_at, owner_id, status, type, data_en FROM repo WHERE LOWER("data_en"->>'name') NOT IN (?)`,
			[]interface{}{"a"},
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSQL, gotValues := tt.query.AsSQL()

			assert.Equal(t, tt.expectedSQL, gotSQL, fmt.Sprintf("Fields %03d :: invalid sql :: %s", i+1, tt.name))
			assert.Equal(t, tt.expectedValues, gotValues, fmt.Sprintf("Fields %03d :: invalid values :: %s", i+1, tt.name))
		})
	}
}

package postgres

import (
	"fmt"
	"github.com/lsldigital/somesql"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestQuery_AsSQL_Fields(t *testing.T) {
	type testCase struct {
		name        string
		query       somesql.Query
		expectedSQL string
	}

	tests := []testCase{
		// Select ALL
		{
			name:        "SELECT *",
			query:       New(),
			expectedSQL: "SELECT id, created_at, updated_at, owner_id, status, type, data_en FROM repo",
		},
		{
			name:        "SELECT * (langEN)",
			query:       New().SetLang(somesql.LangEN),
			expectedSQL: "SELECT id, created_at, updated_at, owner_id, status, type, data_en FROM repo",
		},
		{
			name:        "SELECT * (langFR)",
			query:       New().SetLang(somesql.LangFR),
			expectedSQL: "SELECT id, created_at, updated_at, owner_id, status, type, data_fr FROM repo",
		},
		// Select some pre-defined fields
		{
			name:        "SELECT id, type, data",
			query:       New().Select("id", "type", "data"),
			expectedSQL: `SELECT id, type, data_en FROM repo`,
		},
		{
			name:        "SELECT id, type, data (LangEN)",
			query:       New().Select("id", "type", "data").SetLang(somesql.LangEN),
			expectedSQL: `SELECT id, type, data_en FROM repo`,
		},
		{
			name:        "SELECT id, type, data (LangFR)",
			query:       New().Select("id", "type", "data").SetLang(somesql.LangFR),
			expectedSQL: `SELECT id, type, data_fr FROM repo`,
		},
		// Select pre-defined fields and json attributes ('data_en'/'data_fr') from data_*
		{
			name:        "SELECT id, type, data_en",
			query:       New().Select("id", "type", "data_en"),
			expectedSQL: `SELECT id, type, json_build_object('data_en', data_en->'data_en') "data" FROM repo`,
		},
		{
			name:        "SELECT id, type, data_en (LangEN)",
			query:       New().Select("id", "type", "data_en").SetLang(somesql.LangEN),
			expectedSQL: `SELECT id, type, json_build_object('data_en', data_en->'data_en') "data" FROM repo`,
		},
		{
			name:        "SELECT id, type, data_en (LangFR)",
			query:       New().Select("id", "type", "data_en").SetLang(somesql.LangFR),
			expectedSQL: `SELECT id, type, json_build_object('data_en', data_fr->'data_en') "data" FROM repo`,
		},
		{
			name:        "SELECT id, type, data_fr",
			query:       New().Select("id", "type", "data_fr"),
			expectedSQL: `SELECT id, type, json_build_object('data_fr', data_en->'data_fr') "data" FROM repo`,
		},
		{
			name:        "SELECT id, type, data_fr (LangEN)",
			query:       New().Select("id", "type", "data_fr").SetLang(somesql.LangEN),
			expectedSQL: `SELECT id, type, json_build_object('data_fr', data_en->'data_fr') "data" FROM repo`,
		},
		{
			name:        "SELECT id, type, data_fr (LangFR)",
			query:       New().Select("id", "type", "data_fr").SetLang(somesql.LangFR),
			expectedSQL: `SELECT id, type, json_build_object('data_fr', data_fr->'data_fr') "data" FROM repo`,
		},
		// Select pre-defined fields and json attributes (any other) from data_*
		{
			name:        "SELECT id, type, data_en->'body'",
			query:       New().Select("id", "type", "body"),
			expectedSQL: `SELECT id, type, json_build_object('body', data_en->'body') "data" FROM repo`,
		},
		{
			name:        "SELECT id, type, data_en->'body' (LangEN)",
			query:       New().Select("id", "type", "body").SetLang(somesql.LangEN),
			expectedSQL: `SELECT id, type, json_build_object('body', data_en->'body') "data" FROM repo`,
		},
		{
			name:        "SELECT id, type, data_fr->'body' (LangFR)",
			query:       New().Select("id", "type", "body").SetLang(somesql.LangFR),
			expectedSQL: `SELECT id, type, json_build_object('body', data_fr->'body') "data" FROM repo`,
		},
		// Select pre-defined fields and json attributes (any other + compound) from data_*
		{
			name:        "SELECT id, type, data_en->'body', data_en->'author_id'",
			query:       New().Select("id", "type", "body", "author_id"),
			expectedSQL: `SELECT id, type, json_build_object('body', data_en->'body', 'author_id', data_en->'author_id') "data" FROM repo`,
		},
		{
			name:        "SELECT id, type, data_en->'body', data_en->'author_id' (LangEN)",
			query:       New().Select("id", "type", "body", "author_id").SetLang(somesql.LangEN),
			expectedSQL: `SELECT id, type, json_build_object('body', data_en->'body', 'author_id', data_en->'author_id') "data" FROM repo`,
		},
		{
			name:        "SELECT id, type, data_fr->'body', data_fr->'author_id' (LangFR)",
			query:       New().Select("id", "type", "body", "author_id").SetLang(somesql.LangFR),
			expectedSQL: `SELECT id, type, json_build_object('body', data_fr->'body', 'author_id', data_fr->'author_id') "data" FROM repo`,
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
		query          somesql.Query
		expectedSQL    string
		expectedValues []interface{}
	}

	tests := []testCase{
		// Regular condition clauses on pre-defined fields
		{
			name:           "WHERE id=?",
			query:          New().Select("data").Where(somesql.And("id", "=", "1")),
			expectedSQL:    "SELECT data_en FROM repo WHERE id=?",
			expectedValues: []interface{}{"1"},
		},
		{
			name:           "WHERE id=? AND status=?",
			query:          New().Select("data").Where(somesql.And("id", "=", "1")).Where(somesql.And("status", "=", "published")),
			expectedSQL:    "SELECT data_en FROM repo WHERE id=? AND status=?",
			expectedValues: []interface{}{"1", "published"},
		},
		{
			name:           "WHERE id=? OR status=?",
			query:          New().Select("data").Where(somesql.And("id", "=", "1")).Where(somesql.Or("status", "=", "published")),
			expectedSQL:    "SELECT data_en FROM repo WHERE id=? OR status=?",
			expectedValues: []interface{}{"1", "published"},
		},
		{
			name:           "WHERE id=? AND status=? OR type=?",
			query:          New().Select("data").Where(somesql.And("id", "=", "1")).Where(somesql.And("status", "=", "published")).Where(somesql.Or("type", "=", "article")),
			expectedSQL:    "SELECT data_en FROM repo WHERE id=? AND status=? OR type=?",
			expectedValues: []interface{}{"1", "published", "article"},
		},
		// Regular condition clauses on json attributes
		{
			name:           "WHERE data_en->>'author_id'=?",
			query:          New().Select("data").Where(somesql.And("author_id", "=", "1")),
			expectedSQL:    "SELECT data FROM repo WHERE data_en->>'author_id'=?",
			expectedValues: []interface{}{"1"},
		},
		{
			name:           "WHERE data_fr->>'author_id'=? (langFR)",
			query:          New().Select("data").SetLang(somesql.LangFR).Where(somesql.And("author_id", "=", "1")),
			expectedSQL:    "SELECT data FROM repo WHERE data_fr->>'author_id'=?",
			expectedValues: []interface{}{"1"},
		},
		{
			name:           "WHERE data_fr->>'author_id'=? OR data_fr->>'category_id'=? (langFR)",
			query:          New().Select("data").SetLang(somesql.LangFR).Where(somesql.And("author_id", "=", "1")).Where(somesql.Or("category_id", "=", "2")),
			expectedSQL:    "SELECT data FROM repo WHERE data_fr->>'author_id'=? OR data_fr->>'category_id'=?",
			expectedValues: []interface{}{"1", "2"},
		},
		{
			name:           "WHERE data_fr->>'author_id'=? AND data_fr->>'category_id'=? (langFR)",
			query:          New().Select("data").SetLang(somesql.LangFR).Where(somesql.And("author_id", "=", "1")).Where(somesql.And("category_id", "=", "2")),
			expectedSQL:    "SELECT data FROM repo WHERE data_fr->>'author_id'=? AND data_fr->>'category_id'=?",
			expectedValues: []interface{}{"1", "2"},
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
		query          somesql.Query
		expectedSQL    string
		expectedValues []interface{}
	}

	tests := []testCase{
		{
			name:           "WHERE (... OR ...) [1]",
			query:          New().Select("data").Where(somesql.AndGroup(somesql.And("badge", "=", "video"), somesql.Or("badge", "=", "audio"))),
			expectedSQL:    `SELECT data_en FROM repo WHERE ("data_en"->>'badge'=? OR ("data_en"->>'badge')=?)`,
			expectedValues: []interface{}{"video", "audio"},
		},
		{
			name:           "WHERE (... OR ...) [2]",
			query:          New().Select("data").Where(somesql.OrGroup(somesql.And("badge", "=", "video"), somesql.Or("badge", "=", "audio"))),
			expectedSQL:    `SELECT data_en FROM repo WHERE ("data_en"->>'badge'=? OR ("data_en"->>'badge')=?)`,
			expectedValues: []interface{}{"video", "audio"},
		},
		{
			name:           "WHERE (... AND ...) [1]",
			query:          New().Select("data").Where(somesql.AndGroup(somesql.Or("badge", "=", "video"), somesql.And("has_video", "=", true))),
			expectedSQL:    `SELECT data_en FROM repo WHERE ("data_en"->>'badge'=? AND ("data_en"->>'has_video')::BOOLEAN=?)`,
			expectedValues: []interface{}{"video", true},
		},
		{
			name:           "WHERE (... AND ...) [2]",
			query:          New().Select("data").Where(somesql.OrGroup(somesql.Or("badge", "=", "video"), somesql.And("has_video", "=", true))),
			expectedSQL:    `SELECT data_en FROM repo WHERE ("data_en"->>'badge'=? AND ("data_en"->>'has_video')::BOOLEAN=?)`,
			expectedValues: []interface{}{"published", "video", true},
		},
		{
			name:           "WHERE (... AND ...) AND (... OR ...) [1]",
			query:          New().Select("data").Where(somesql.OrGroup(somesql.Or("badge", "=", "video"), somesql.And("has_video", "=", true))).Where(somesql.AndGroup(somesql.Or("badge", "=", "video"), somesql.Or("has_video", "=", true))),
			expectedSQL:    `SELECT data_en FROM repo WHERE ("data_en"->>'badge'=? AND ("data_en"->>'has_video')::BOOLEAN=?) AND ("data_en"->>'badge'=? OR ("data_en"->>'has_video')::BOOLEAN=?)`,
			expectedValues: []interface{}{"published", "video", true},
		},
		{
			name:           "WHERE (... AND ...) OR (... AND ...) [2]",
			query:          New().Select("data").Where(somesql.OrGroup(somesql.Or("badge", "=", "video"), somesql.And("has_video", "=", true))).Where(somesql.OrGroup(somesql.Or("badge", "=", "video"), somesql.And("has_video", "=", true))),
			expectedSQL:    `SELECT data_en FROM repo WHERE ("data_en"->>'badge'=? AND ("data_en"->>'has_video')::BOOLEAN=?) OR ("data_en"->>'badge'=? AND ("data_en"->>'has_video')::BOOLEAN=?)`,
			expectedValues: []interface{}{"published", "video", true},
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

func TestQuery_AsSQL_ConditionIN(t *testing.T) {
	type testCase struct {
		name           string
		query          somesql.Query
		expectedSQL    string
		expectedValues []interface{}
	}

	tests := []testCase{
		{
			name:           "WHERE id IN (...) - primitive field",
			query:          New().Select("data").Where(somesql.AndIn("id", []string{"A", "B", "C"})),
			expectedSQL:    "SELECT data_en FROM repo WHERE id IN (?,?,?)",
			expectedValues: []interface{}{"A", "B", "C"},
		},
		{
			name:           "WHERE id IN (...) - primitive field - LangFR",
			query:          New().Select("data").Where(somesql.AndIn("id", []string{"A", "B", "C"})).SetLang(somesql.LangFR),
			expectedSQL:    "SELECT data_fr FROM repo WHERE id IN (?,?,?)",
			expectedValues: []interface{}{"A", "B", "C"},
		},
		{
			name:           "WHERE field IN (...) - JSONB",
			query:          New().Where(somesql.AndIn("name", []string{"A", "B", "C"})),
			expectedSQL:    `SELECT id, created_at, updated_at, owner_id, status, type, data_en FROM repo WHERE "data_en"->>'name' IN (?,?,?)`,
			expectedValues: []interface{}{"A", "B", "C"},
		},
		{
			name:           "WHERE field IN (...) - JSONB - LangFR",
			query:          New().Where(somesql.AndIn("name", []string{"A", "B", "C"})).SetLang(somesql.LangFR),
			expectedSQL:    `SELECT id, created_at, updated_at, owner_id, status, type, data_fr FROM repo WHERE "data_fr"->>'name' IN (?,?,?)`,
			expectedValues: []interface{}{"A", "B", "C"},
		},
		{
			name:           "WHERE FUNC(field) IN (...) - primitive field",
			query:          New().Select("data").Where(somesql.AndIn("updated_at", []string{"2019"}, "YEAR")),
			expectedSQL:    "SELECT data_en FROM repo WHERE YEAR(updated_at) IN (?)",
			expectedValues: []interface{}{"2019"},
		},
		{
			name:           "WHERE FUNC(field) IN (...) - JSONB",
			query:          New().Where(somesql.AndIn("name", []string{"a"}, "LOWER")),
			expectedSQL:    `SELECT id, created_at, updated_at, owner_id, status, type, data_en FROM repo WHERE LOWER("data_en"->>'name') IN (?)`,
			expectedValues: []interface{}{"a"},
		},

		{
			name:           "WHERE id NOT IN (...) - primitive field",
			query:          New().Select("data").Where(somesql.AndNotIn("id", []string{"A", "B", "C"})),
			expectedSQL:    "SELECT data_en FROM repo WHERE id NOT IN (?,?,?)",
			expectedValues: []interface{}{"A", "B", "C"},
		},
		{
			name:           "WHERE id NOT IN (...) - primitive field - LangFR",
			query:          New().Select("data").Where(somesql.AndNotIn("id", []string{"A", "B", "C"})).SetLang(somesql.LangFR),
			expectedSQL:    "SELECT data_fr FROM repo WHERE id NOT IN (?,?,?)",
			expectedValues: []interface{}{"A", "B", "C"},
		},
		{
			name:           "WHERE FUNC(field) NOT IN (...) - primitive field",
			query:          New().Select("data").Where(somesql.AndNotIn("updated_at", []string{"2019"}, "YEAR")),
			expectedSQL:    "SELECT data_en FROM repo WHERE YEAR(updated_at) NOT IN (?)",
			expectedValues: []interface{}{"2019"},
		},
		{
			name:           "WHERE field NOT IN (...) - JSONB",
			query:          New().Where(somesql.AndNotIn("name", []string{"A", "B", "C"})),
			expectedSQL:    `SELECT id, created_at, updated_at, owner_id, status, type, data_en FROM repo WHERE "data_en"->>'name' NOT IN (?,?,?)`,
			expectedValues: []interface{}{"A", "B", "C"},
		},
		{
			name:           "WHERE field NOT IN (...) - JSONB",
			query:          New().Where(somesql.AndNotIn("name", []string{"A", "B", "C"})),
			expectedSQL:    `SELECT id, created_at, updated_at, owner_id, status, type, data_en FROM repo WHERE "data_en"->>'name' NOT IN (?,?,?)`,
			expectedValues: []interface{}{"A", "B", "C"},
		},
		{
			name:           "WHERE FUNC(field) NOT IN (...) - JSONB",
			query:          New().Where(somesql.AndIn("name", []string{"a"}, "LOWER")),
			expectedSQL:    `SELECT id, created_at, updated_at, owner_id, status, type, data_en FROM repo WHERE LOWER("data_en"->>'name') NOT IN (?)`,
			expectedValues: []interface{}{"a"},
		},

		{
			name:           "WHERE id IN (...) AND NOT IN (...) - primitive field",
			query:          New().Select("data").Where(somesql.AndIn("id", []string{"A", "B"})).Where(somesql.AndNotIn("id", []string{"C", "D"})),
			expectedSQL:    "SELECT data_en FROM repo WHERE id IN (?,?) AND id NOT IN (?,?)",
			expectedValues: []interface{}{"A", "B", "C", "D"},
		},
		{
			name:           "WHERE id IN (...) AND NOT IN (...) - JSONB",
			query:          New().Select("data").Where(somesql.AndIn("name", []string{"A", "B"})).Where(somesql.AndNotIn("id", []string{"C", "D"})),
			expectedSQL:    "SELECT data_en FROM repo WHERE id IN (?,?) AND id NOT IN (?,?)",
			expectedValues: []interface{}{"A", "B", "C", "D"},
		},
		{
			name:           "WHERE id IN (...) OR NOT IN (...) - primitive field",
			query:          New().Select("data").Where(somesql.AndIn("id", []string{"A", "B"})).Where(somesql.OrNotIn("id", []string{"C", "D"})),
			expectedSQL:    "SELECT data_en FROM repo WHERE id IN (?,?) OR id NOT IN (?,?)",
			expectedValues: []interface{}{"A", "B", "C", "D"},
		},
		{
			name:           "WHERE FUNC(field) IN (...) AND field NOT IN (...) - primitive field",
			query:          New().Select("data").Where(somesql.AndIn("updated_at", []string{"2019"}, "YEAR")).Where(somesql.AndNotIn("id", []string{"A", "B"})),
			expectedSQL:    "SELECT data_en FROM repo WHERE YEAR(updated_at) IN (?) AND id NOT IN (?,?)",
			expectedValues: []interface{}{"2019", "A", "B"},
		},

		{
			name:           "WHERE id IN (...) AND field = ...",
			query:          New().Select("data").Where(somesql.AndIn("type", []string{"article", "dossier"})).Where(somesql.And("status", "=", []string{"published"})),
			expectedSQL:    "SELECT data_en FROM repo WHERE type IN (?) AND status =?",
			expectedValues: []interface{}{"article", "dossier", "published"},
		},
		{
			name:           "WHERE id IN (...) AND field = ... - JSONB",
			query:          New().Select("data").Where(somesql.AndIn("type", []string{"article", "dossier"})).Where(somesql.And("status", "=", []string{"published"})),
			expectedSQL:    "SELECT data_en FROM repo WHERE type IN (?) AND status =?",
			expectedValues: []interface{}{"article", "dossier", "published"},
		},
		{
			name:           "WHERE FUNC(field) IN (...) AND field NOT IN (...) - primitive field",
			query:          New().Select("data").Where(somesql.AndIn("updated_at", []string{"2019"}, "YEAR")).Where(somesql.AndNotIn("id", []string{"A", "B"})).Where(somesql.And("status", "=", "published")),
			expectedSQL:    "SELECT data_en FROM repo WHERE YEAR(updated_at) IN (?) AND id NOT IN (?,?) AND status =?",
			expectedValues: []interface{}{"2019", "A", "B", "published"},
		},
		{
			name:           "WHERE FUNC(field) IN (...) AND field NOT IN (...) - JSONB",
			query:          New().Select("data").Where(somesql.AndIn("tag_ids", []string{"A"})).Where(somesql.AndIn("author_ids", []string{"B"})).Where(somesql.And("status", "=", "published")),
			expectedSQL:    `SELECT data_en FROM repo WHERE "data_en"->>'tag_ids' IN (?) AND "data_en"->>'author_ids' IN (?) AND status =?`,
			expectedValues: []interface{}{"A", "B", "published"},
		},

		{
			name:           "WHERE id IN (...) AND field = ...",
			query:          New().Select("data").Where(somesql.AndGroup(somesql.And("badge", "=", "video"), somesql.And("has_video", "=", true))).Where(somesql.And("status", "=", []string{"published"})),
			expectedSQL:    `SELECT data_en FROM repo WHERE ("data_en"->>'badge'=? AND ("data_en"->>'has_video')::BOOLEAN=?) AND status =?`,
			expectedValues: []interface{}{"video", true, "published"},
		},
		{
			name:           "WHERE field = ... OR (... AND ...) [1]",
			query:          New().Select("data").Where(somesql.And("status", "=", []string{"published"})).Where(somesql.OrGroup(somesql.And("badge", "=", "video"), somesql.And("has_video", "=", true))),
			expectedSQL:    `SELECT data_en FROM repo WHERE status =? OR ("data_en"->>'badge'=? AND ("data_en"->>'has_video')::BOOLEAN=?)`,
			expectedValues: []interface{}{"published", "video", true},
		},
		{
			name:           "WHERE field = ... OR (... AND ...) [2]",
			query:          New().Select("data").Where(somesql.Or("status", "=", []string{"published"})).Where(somesql.OrGroup(somesql.Or("badge", "=", "video"), somesql.And("has_video", "=", true))),
			expectedSQL:    `SELECT data_en FROM repo WHERE status =? OR ("data_en"->>'badge'=? AND ("data_en"->>'has_video')::BOOLEAN=?)`,
			expectedValues: []interface{}{"published", "video", true},
		},
		{
			name:           "WHERE field = ... AND (... OR ...) [1]",
			query:          New().Select("data").Where(somesql.Or("status", "=", []string{"published"})).Where(somesql.AndGroup(somesql.And("badge", "=", "video"), somesql.Or("has_video", "=", true))),
			expectedSQL:    `SELECT data_en FROM repo WHERE status =? AND ("data_en"->>'badge'=? OR ("data_en"->>'has_video')::BOOLEAN=?)`,
			expectedValues: []interface{}{"published", "video", true},
		},
		{
			name:           "WHERE field = ... AND (... OR ...) [2]",
			query:          New().Select("data").Where(somesql.Or("status", "=", []string{"published"})).Where(somesql.AndGroup(somesql.And("badge", "=", "video"), somesql.Or("badge", "=", "audio"))),
			expectedSQL:    `SELECT data_en FROM repo WHERE status =? AND ("data_en"->>'badge'=? OR ("data_en"->>'badge')=?)`,
			expectedValues: []interface{}{"published", "video", "audio"},
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

func TestQuesy_AsSQL_InQuery(t *testing.T) {
	type testCase struct {
		name           string
		query          somesql.Query
		expectedSQL    string
		expectedValues []interface{}
	}

	tests := []testCase{
		{
			name:           "AndInQuery",
			query:          New().Select("data").Where(somesql.AndIn("id", []string{"A", "B", "C"})).Where(somesql.AndInQuery("author_id", New().Select("author_id").Where(somesql.And("id", "=", "002fd6b1-f715-4875-838b-1546f27327df")))),
			expectedSQL:    `SELECT data_en FROM repo WHERE id IN (?,?,?) AND "data_en"->>'author_id' IN (SELECT "data_en"->>'author_id' FROM repo WHERE id =?)`,
			expectedValues: []interface{}{"A", "B", "C", "002fd6b1-f715-4875-838b-1546f27327df"},
		},
		{
			name:           "AndNotInQuery",
			query:          New().Select("data").Where(somesql.AndIn("id", []string{"A", "B", "C"})).Where(somesql.AndNotInQuery("author_id", New().Select("author_id").Where(somesql.And("id", "=", "002fd6b1-f715-4875-838b-1546f27327df")))),
			expectedSQL:    `SELECT data_en FROM repo WHERE id IN (?,?,?) AND "data_en"->>'author_id' NOT IN (SELECT "data_en"->>'author_id' FROM repo WHERE id =?)`,
			expectedValues: []interface{}{"A", "B", "C", "002fd6b1-f715-4875-838b-1546f27327df"},
		},

		{
			name:           "OrInQuery",
			query:          New().Select("data").Where(somesql.AndIn("id", []string{"A", "B", "C"})).Where(somesql.OrInQuery("author_id", New().Select("author_id").Where(somesql.And("id", "=", "002fd6b1-f715-4875-838b-1546f27327df")))),
			expectedSQL:    `SELECT data_en FROM repo WHERE id IN (?,?,?) OR "data_en"->>'author_id' IN (SELECT "data_en"->>'author_id' FROM repo WHERE id =?)`,
			expectedValues: []interface{}{"A", "B", "C", "002fd6b1-f715-4875-838b-1546f27327df"},
		},
		{
			name:           "OrNotInQuery",
			query:          New().Select("data").Where(somesql.AndIn("id", []string{"A", "B", "C"})).Where(somesql.OrNotInQuery("author_id", New().Select("author_id").Where(somesql.And("id", "=", "002fd6b1-f715-4875-838b-1546f27327df")))),
			expectedSQL:    `SELECT data_en FROM repo WHERE id IN (?,?,?) OR "data_en"->>'author_id' NOT IN (SELECT "data_en"->>'author_id' FROM repo WHERE id =?)`,
			expectedValues: []interface{}{"A", "B", "C", "002fd6b1-f715-4875-838b-1546f27327df"},
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

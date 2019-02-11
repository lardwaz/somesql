package somesql_test

import (
	"fmt"
	"testing"

	"github.com/lsldigital/somesql"
	"github.com/stretchr/testify/assert"
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
			query:       somesql.NewQuery(),
			expectedSQL: "SELECT id, created_at, updated_at, owner_id, status, type, data_en FROM repo",
		},
		{
			name:        "SELECT * (langEN)",
			query:       somesql.NewQuery().SetLang(somesql.LangEN),
			expectedSQL: "SELECT id, created_at, updated_at, owner_id, status, type, data_en FROM repo",
		},
		{
			name:        "SELECT * (langFR)",
			query:       somesql.NewQuery().SetLang(somesql.LangFR),
			expectedSQL: "SELECT id, created_at, updated_at, owner_id, status, type, data_fr FROM repo",
		},
		// Select some pre-defined fields
		{
			name:        "SELECT id, type, data",
			query:       somesql.NewQuery().Select("id", "type", "data"),
			expectedSQL: `SELECT id, type, data_en FROM repo`,
		},
		{
			name:        "SELECT id, type, data (LangEN)",
			query:       somesql.NewQuery().Select("id", "type", "data").SetLang(somesql.LangEN),
			expectedSQL: `SELECT id, type, data_en FROM repo`,
		},
		{
			name:        "SELECT id, type, data (LangFR)",
			query:       somesql.NewQuery().Select("id", "type", "data").SetLang(somesql.LangFR),
			expectedSQL: `SELECT id, type, data_fr FROM repo`,
		},
		// Select pre-defined fields and json attributes ('data_en'/'data_fr') from data_*
		{
			name:        "SELECT id, type, data_en",
			query:       somesql.NewQuery().Select("id", "type", "data_en"),
			expectedSQL: `SELECT id, type, json_build_object('data_en', data_en->'data_en') "data" FROM repo`,
		},
		{
			name:        "SELECT id, type, data_en (LangEN)",
			query:       somesql.NewQuery().Select("id", "type", "data_en").SetLang(somesql.LangEN),
			expectedSQL: `SELECT id, type, json_build_object('data_en', data_en->'data_en') "data" FROM repo`,
		},
		{
			name:        "SELECT id, type, data_en (LangFR)",
			query:       somesql.NewQuery().Select("id", "type", "data_en").SetLang(somesql.LangFR),
			expectedSQL: `SELECT id, type, json_build_object('data_en', data_fr->'data_en') "data" FROM repo`,
		},
		{
			name:        "SELECT id, type, data_fr",
			query:       somesql.NewQuery().Select("id", "type", "data_fr"),
			expectedSQL: `SELECT id, type, json_build_object('data_fr', data_en->'data_fr') "data" FROM repo`,
		},
		{
			name:        "SELECT id, type, data_fr (LangEN)",
			query:       somesql.NewQuery().Select("id", "type", "data_fr").SetLang(somesql.LangEN),
			expectedSQL: `SELECT id, type, json_build_object('data_fr', data_en->'data_fr') "data" FROM repo`,
		},
		{
			name:        "SELECT id, type, data_fr (LangFR)",
			query:       somesql.NewQuery().Select("id", "type", "data_fr").SetLang(somesql.LangFR),
			expectedSQL: `SELECT id, type, json_build_object('data_fr', data_fr->'data_fr') "data" FROM repo`,
		},
		// Select pre-defined fields and json attributes (any other) from data_*
		{
			name:        "SELECT id, type, data_en->'body'",
			query:       somesql.NewQuery().Select("id", "type", "body"),
			expectedSQL: `SELECT id, type, json_build_object('body', data_en->'body') "data" FROM repo`,
		},
		{
			name:        "SELECT id, type, data_en->'body' (LangEN)",
			query:       somesql.NewQuery().Select("id", "type", "body").SetLang(somesql.LangEN),
			expectedSQL: `SELECT id, type, json_build_object('body', data_en->'body') "data" FROM repo`,
		},
		{
			name:        "SELECT id, type, data_fr->'body' (LangFR)",
			query:       somesql.NewQuery().Select("id", "type", "body").SetLang(somesql.LangFR),
			expectedSQL: `SELECT id, type, json_build_object('body', data_fr->'body') "data" FROM repo`,
		},
		// Select pre-defined fields and json attributes (any other + compound) from data_*
		{
			name:        "SELECT id, type, data_en->'body', data_en->'author_id'",
			query:       somesql.NewQuery().Select("id", "type", "body", "author_id"),
			expectedSQL: `SELECT id, type, json_build_object('body', data_en->'body', 'author_id', data_en->'author_id') "data" FROM repo`,
		},
		{
			name:        "SELECT id, type, data_en->'body', data_en->'author_id' (LangEN)",
			query:       somesql.NewQuery().Select("id", "type", "body", "author_id").SetLang(somesql.LangEN),
			expectedSQL: `SELECT id, type, json_build_object('body', data_en->'body', 'author_id', data_en->'author_id') "data" FROM repo`,
		},
		{
			name:        "SELECT id, type, data_fr->'body', data_fr->'author_id' (LangFR)",
			query:       somesql.NewQuery().Select("id", "type", "body", "author_id").SetLang(somesql.LangFR),
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
			query:          somesql.NewQuery().Select("data").Where(somesql.And(somesql.LangEN, "id", "=", "1")),
			expectedSQL:    "SELECT data_en FROM repo WHERE id=?",
			expectedValues: []interface{}{"1"},
		},
		{
			name:           "WHERE id=? AND status=?",
			query:          somesql.NewQuery().Select("data").Where(somesql.And(somesql.LangEN, "id", "=", "1")).Where(somesql.And(somesql.LangEN, "status", "=", "published")),
			expectedSQL:    "SELECT data_en FROM repo WHERE id=? AND status=?",
			expectedValues: []interface{}{"1", "published"},
		},
		{
			name:           "WHERE id=? OR status=?",
			query:          somesql.NewQuery().Select("data").Where(somesql.And(somesql.LangEN, "id", "=", "1")).Where(somesql.Or(somesql.LangEN, "status", "=", "published")),
			expectedSQL:    "SELECT data_en FROM repo WHERE id=? OR status=?",
			expectedValues: []interface{}{"1", "published"},
		},
		{
			name:           "WHERE id=? AND status=? OR type=?",
			query:          somesql.NewQuery().Select("data").Where(somesql.And(somesql.LangEN, "id", "=", "1")).Where(somesql.And(somesql.LangEN, "status", "=", "published")).Where(somesql.Or(somesql.LangEN, "type", "=", "article")),
			expectedSQL:    "SELECT data_en FROM repo WHERE id=? AND status=? OR type=?",
			expectedValues: []interface{}{"1", "published", "article"},
		},
		// Regular condition clauses on json attributes
		{
			name:           `WHERE "data_en"->>'author_id'=?`,
			query:          somesql.NewQuery().Select("data").Where(somesql.And(somesql.LangEN, "author_id", "=", "1")),
			expectedSQL:    `SELECT data_en FROM repo WHERE "data_en"->>'author_id'=?`,
			expectedValues: []interface{}{"1"},
		},
		{
			name:           `WHERE "data_fr"->>'author_id'=? (langFR)`,
			query:          somesql.NewQuery().Select("data").SetLang(somesql.LangFR).Where(somesql.And(somesql.LangFR, "author_id", "=", "1")),
			expectedSQL:    `SELECT data_fr FROM repo WHERE "data_fr"->>'author_id'=?`,
			expectedValues: []interface{}{"1"},
		},
		{
			name:           `WHERE "data_fr"->>'author_id'=? OR "data_fr"->>'category_id'=? (langFR)`,
			query:          somesql.NewQuery().Select("data").SetLang(somesql.LangFR).Where(somesql.And(somesql.LangFR, "author_id", "=", "1")).Where(somesql.Or(somesql.LangFR, "category_id", "=", "2")),
			expectedSQL:    `SELECT data_fr FROM repo WHERE "data_fr"->>'author_id'=? OR "data_fr"->>'category_id'=?`,
			expectedValues: []interface{}{"1", "2"},
		},
		{
			name:           `WHERE "data_fr"->>'author_id'=? AND "data_fr"->>'category_id'=? (langFR)`,
			query:          somesql.NewQuery().Select("data").SetLang(somesql.LangFR).Where(somesql.And(somesql.LangFR, "author_id", "=", "1")).Where(somesql.And(somesql.LangFR, "category_id", "=", "2")),
			expectedSQL:    `SELECT data_fr FROM repo WHERE "data_fr"->>'author_id'=? AND "data_fr"->>'category_id'=?`,
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
			query:          somesql.NewQuery().Select("data").Where(somesql.AndGroup(somesql.And(somesql.LangEN, "badge", "=", "video"), somesql.Or(somesql.LangEN, "badge", "=", "audio"))),
			expectedSQL:    `SELECT data_en FROM repo WHERE ("data_en"->>'badge'=? OR "data_en"->>'badge'=?)`,
			expectedValues: []interface{}{"video", "audio"},
		},
		{
			name:           "WHERE (... OR ...) [2]",
			query:          somesql.NewQuery().Select("data").Where(somesql.OrGroup(somesql.And(somesql.LangEN, "badge", "=", "video"), somesql.Or(somesql.LangEN, "badge", "=", "audio"))),
			expectedSQL:    `SELECT data_en FROM repo WHERE ("data_en"->>'badge'=? OR "data_en"->>'badge'=?)`,
			expectedValues: []interface{}{"video", "audio"},
		},
		{
			name:           "WHERE (... AND ...) [1]",
			query:          somesql.NewQuery().Select("data").Where(somesql.AndGroup(somesql.Or(somesql.LangEN, "badge", "=", "video"), somesql.And(somesql.LangEN, "has_video", "=", true))),
			expectedSQL:    `SELECT data_en FROM repo WHERE ("data_en"->>'badge'=? AND ("data_en"->>'has_video')::BOOLEAN=?)`,
			expectedValues: []interface{}{"video", true},
		},
		{
			name:           "WHERE (... AND ...) [2]",
			query:          somesql.NewQuery().Select("data").Where(somesql.OrGroup(somesql.Or(somesql.LangEN, "badge", "=", "video"), somesql.And(somesql.LangEN, "has_video", "=", true))),
			expectedSQL:    `SELECT data_en FROM repo WHERE ("data_en"->>'badge'=? AND ("data_en"->>'has_video')::BOOLEAN=?)`,
			expectedValues: []interface{}{"video", true},
		},
		{
			name:           "WHERE (... AND ...) AND (... OR ...) [1]",
			query:          somesql.NewQuery().Select("data").Where(somesql.OrGroup(somesql.Or(somesql.LangEN, "badge", "=", "video"), somesql.And(somesql.LangEN, "has_video", "=", true))).Where(somesql.AndGroup(somesql.Or(somesql.LangEN, "badge", "=", "video"), somesql.Or(somesql.LangEN, "has_video", "=", true))),
			expectedSQL:    `SELECT data_en FROM repo WHERE ("data_en"->>'badge'=? AND ("data_en"->>'has_video')::BOOLEAN=?) AND ("data_en"->>'badge'=? OR ("data_en"->>'has_video')::BOOLEAN=?)`,
			expectedValues: []interface{}{"video", true, "video", true},
		},
		{
			name:           "WHERE (... AND ...) OR (... AND ...) [2]",
			query:          somesql.NewQuery().Select("data").Where(somesql.OrGroup(somesql.Or(somesql.LangEN, "badge", "=", "video"), somesql.And(somesql.LangEN, "has_video", "=", true))).Where(somesql.OrGroup(somesql.Or(somesql.LangEN, "badge", "=", "video"), somesql.And(somesql.LangEN, "has_video", "=", true))),
			expectedSQL:    `SELECT data_en FROM repo WHERE ("data_en"->>'badge'=? AND ("data_en"->>'has_video')::BOOLEAN=?) OR ("data_en"->>'badge'=? AND ("data_en"->>'has_video')::BOOLEAN=?)`,
			expectedValues: []interface{}{"video", true, "video", true},
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
			query:          somesql.NewQuery().Select("data").Where(somesql.AndIn(somesql.LangEN, "id", []string{"A", "B", "C"})),
			expectedSQL:    "SELECT data_en FROM repo WHERE id IN (?,?,?)",
			expectedValues: []interface{}{"A", "B", "C"},
		},
		{
			name:           "WHERE id IN (...) - primitive field - LangFR",
			query:          somesql.NewQuery().Select("data").Where(somesql.AndIn(somesql.LangEN, "id", []string{"A", "B", "C"})).SetLang(somesql.LangFR),
			expectedSQL:    "SELECT data_fr FROM repo WHERE id IN (?,?,?)",
			expectedValues: []interface{}{"A", "B", "C"},
		},
		{
			name:           "WHERE field IN (...) - JSONB",
			query:          somesql.NewQuery().Where(somesql.AndIn(somesql.LangEN, "name", []string{"A", "B", "C"})),
			expectedSQL:    `SELECT id, created_at, updated_at, owner_id, status, type, data_en FROM repo WHERE "data_en"->>'name' IN (?,?,?)`,
			expectedValues: []interface{}{"A", "B", "C"},
		},
		{
			name:           "WHERE field IN (...) - JSONB - LangFR",
			query:          somesql.NewQuery().Where(somesql.AndIn(somesql.LangFR, "name", []string{"A", "B", "C"})).SetLang(somesql.LangFR),
			expectedSQL:    `SELECT id, created_at, updated_at, owner_id, status, type, data_fr FROM repo WHERE "data_fr"->>'name' IN (?,?,?)`,
			expectedValues: []interface{}{"A", "B", "C"},
		},
		{
			name:           "WHERE FUNC(field) IN (...) - primitive field",
			query:          somesql.NewQuery().Select("data").Where(somesql.AndIn(somesql.LangEN, "updated_at", []string{"2019"}, "YEAR")),
			expectedSQL:    "SELECT data_en FROM repo WHERE YEAR(updated_at) IN (?)",
			expectedValues: []interface{}{"2019"},
		},
		{
			name:           "WHERE FUNC(field) IN (...) - JSONB",
			query:          somesql.NewQuery().Where(somesql.AndIn(somesql.LangEN, "name", []string{"a"}, "LOWER")),
			expectedSQL:    `SELECT id, created_at, updated_at, owner_id, status, type, data_en FROM repo WHERE LOWER("data_en"->>'name') IN (?)`,
			expectedValues: []interface{}{"a"},
		},

		{
			name:           "WHERE id NOT IN (...) - primitive field",
			query:          somesql.NewQuery().Select("data").Where(somesql.AndNotIn(somesql.LangEN, "id", []string{"A", "B", "C"})),
			expectedSQL:    "SELECT data_en FROM repo WHERE id NOT IN (?,?,?)",
			expectedValues: []interface{}{"A", "B", "C"},
		},
		{
			name:           "WHERE id NOT IN (...) - primitive field - LangFR",
			query:          somesql.NewQuery().Select("data").Where(somesql.AndNotIn(somesql.LangEN, "id", []string{"A", "B", "C"})).SetLang(somesql.LangFR),
			expectedSQL:    "SELECT data_fr FROM repo WHERE id NOT IN (?,?,?)",
			expectedValues: []interface{}{"A", "B", "C"},
		},
		{
			name:           "WHERE FUNC(field) NOT IN (...) - primitive field",
			query:          somesql.NewQuery().Select("data").Where(somesql.AndNotIn(somesql.LangEN, "updated_at", []string{"2019"}, "YEAR")),
			expectedSQL:    "SELECT data_en FROM repo WHERE YEAR(updated_at) NOT IN (?)",
			expectedValues: []interface{}{"2019"},
		},
		{
			name:           "WHERE field NOT IN (...) - JSONB",
			query:          somesql.NewQuery().Where(somesql.AndNotIn(somesql.LangEN, "name", []string{"A", "B", "C"})),
			expectedSQL:    `SELECT id, created_at, updated_at, owner_id, status, type, data_en FROM repo WHERE "data_en"->>'name' NOT IN (?,?,?)`,
			expectedValues: []interface{}{"A", "B", "C"},
		},
		{
			name:           "WHERE field NOT IN (...) - JSONB",
			query:          somesql.NewQuery().Where(somesql.AndNotIn(somesql.LangEN, "name", []string{"A", "B", "C"})),
			expectedSQL:    `SELECT id, created_at, updated_at, owner_id, status, type, data_en FROM repo WHERE "data_en"->>'name' NOT IN (?,?,?)`,
			expectedValues: []interface{}{"A", "B", "C"},
		},
		{
			name:           "WHERE FUNC(field) NOT IN (...) - JSONB",
			query:          somesql.NewQuery().Where(somesql.AndNotIn(somesql.LangEN, "name", []string{"a"}, "LOWER")),
			expectedSQL:    `SELECT id, created_at, updated_at, owner_id, status, type, data_en FROM repo WHERE LOWER("data_en"->>'name') NOT IN (?)`,
			expectedValues: []interface{}{"a"},
		},

		{
			name:           "WHERE id IN (...) AND NOT IN (...) - primitive field",
			query:          somesql.NewQuery().Select("data").Where(somesql.AndIn(somesql.LangEN, "id", []string{"A", "B"})).Where(somesql.AndNotIn(somesql.LangEN, "id", []string{"C", "D"})),
			expectedSQL:    "SELECT data_en FROM repo WHERE id IN (?,?) AND id NOT IN (?,?)",
			expectedValues: []interface{}{"A", "B", "C", "D"},
		},
		{
			name:           "WHERE id IN (...) AND NOT IN (...) - JSONB",
			query:          somesql.NewQuery().Select("data").Where(somesql.AndIn(somesql.LangEN, "name", []string{"A", "B"})).Where(somesql.AndNotIn(somesql.LangEN, "id", []string{"C", "D"})),
			expectedSQL:    `SELECT data_en FROM repo WHERE "data_en"->>'name' IN (?,?) AND id NOT IN (?,?)`,
			expectedValues: []interface{}{"A", "B", "C", "D"},
		},
		{
			name:           "WHERE id IN (...) OR NOT IN (...) - primitive field",
			query:          somesql.NewQuery().Select("data").Where(somesql.AndIn(somesql.LangEN, "id", []string{"A", "B"})).Where(somesql.OrNotIn(somesql.LangEN, "id", []string{"C", "D"})),
			expectedSQL:    "SELECT data_en FROM repo WHERE id IN (?,?) OR id NOT IN (?,?)",
			expectedValues: []interface{}{"A", "B", "C", "D"},
		},
		{
			name:           "WHERE FUNC(field) IN (...) AND field NOT IN (...) - primitive field",
			query:          somesql.NewQuery().Select("data").Where(somesql.AndIn(somesql.LangEN, "updated_at", []string{"2019"}, "YEAR")).Where(somesql.AndNotIn(somesql.LangEN, "id", []string{"A", "B"})),
			expectedSQL:    "SELECT data_en FROM repo WHERE YEAR(updated_at) IN (?) AND id NOT IN (?,?)",
			expectedValues: []interface{}{"2019", "A", "B"},
		},

		{
			name:           "WHERE id IN (...) AND field = ...",
			query:          somesql.NewQuery().Select("data").Where(somesql.AndIn(somesql.LangEN, "type", []string{"article", "dossier"})).Where(somesql.And(somesql.LangEN, "status", "=", []string{"published"})),
			expectedSQL:    "SELECT data_en FROM repo WHERE type IN (?,?) AND status=?",
			expectedValues: []interface{}{"article", "dossier", "published"},
		},
		{
			name:           "WHERE id IN (...) AND field = ... - JSONB",
			query:          somesql.NewQuery().Select("data").Where(somesql.AndIn(somesql.LangEN, "type", []string{"article", "dossier"})).Where(somesql.And(somesql.LangEN, "status", "=", []string{"published"})),
			expectedSQL:    "SELECT data_en FROM repo WHERE type IN (?,?) AND status=?",
			expectedValues: []interface{}{"article", "dossier", "published"},
		},
		{
			name:           "WHERE FUNC(field) IN (...) AND field NOT IN (...) - primitive field",
			query:          somesql.NewQuery().Select("data").Where(somesql.AndIn(somesql.LangEN, "updated_at", []string{"2019"}, "YEAR")).Where(somesql.AndNotIn(somesql.LangEN, "id", []string{"A", "B"})).Where(somesql.And(somesql.LangEN, "status", "=", "published")),
			expectedSQL:    "SELECT data_en FROM repo WHERE YEAR(updated_at) IN (?) AND id NOT IN (?,?) AND status=?",
			expectedValues: []interface{}{"2019", "A", "B", "published"},
		},
		{
			name:           "WHERE FUNC(field) IN (...) AND field NOT IN (...) - JSONB",
			query:          somesql.NewQuery().Select("data").Where(somesql.AndIn(somesql.LangEN, "tag_ids", []string{"A"})).Where(somesql.AndIn(somesql.LangEN, "author_ids", []string{"B"})).Where(somesql.And(somesql.LangEN, "status", "=", "published")),
			expectedSQL:    `SELECT data_en FROM repo WHERE "data_en"->>'tag_ids' IN (?) AND "data_en"->>'author_ids' IN (?) AND status=?`,
			expectedValues: []interface{}{"A", "B", "published"},
		},

		{
			name:           "WHERE id IN (...) AND field = ...",
			query:          somesql.NewQuery().Select("data").Where(somesql.AndGroup(somesql.And(somesql.LangEN, "badge", "=", "video"), somesql.And(somesql.LangEN, "has_video", "=", true))).Where(somesql.And(somesql.LangEN, "status", "=", []string{"published"})),
			expectedSQL:    `SELECT data_en FROM repo WHERE ("data_en"->>'badge'=? AND ("data_en"->>'has_video')::BOOLEAN=?) AND status=?`,
			expectedValues: []interface{}{"video", true, "published"},
		},
		{
			name:           "WHERE field = ... OR (... AND ...) [1]",
			query:          somesql.NewQuery().Select("data").Where(somesql.And(somesql.LangEN, "status", "=", []string{"published"})).Where(somesql.OrGroup(somesql.And(somesql.LangEN, "badge", "=", "video"), somesql.And(somesql.LangEN, "has_video", "=", true))),
			expectedSQL:    `SELECT data_en FROM repo WHERE status=? OR ("data_en"->>'badge'=? AND ("data_en"->>'has_video')::BOOLEAN=?)`,
			expectedValues: []interface{}{"published", "video", true},
		},
		{
			name:           "WHERE field = ... OR (... AND ...) [2]",
			query:          somesql.NewQuery().Select("data").Where(somesql.Or(somesql.LangEN, "status", "=", []string{"published"})).Where(somesql.OrGroup(somesql.Or(somesql.LangEN, "badge", "=", "video"), somesql.And(somesql.LangEN, "has_video", "=", true))),
			expectedSQL:    `SELECT data_en FROM repo WHERE status=? OR ("data_en"->>'badge'=? AND ("data_en"->>'has_video')::BOOLEAN=?)`,
			expectedValues: []interface{}{"published", "video", true},
		},
		{
			name:           "WHERE field = ... AND (... OR ...) [1]",
			query:          somesql.NewQuery().Select("data").Where(somesql.Or(somesql.LangEN, "status", "=", []string{"published"})).Where(somesql.AndGroup(somesql.And(somesql.LangEN, "badge", "=", "video"), somesql.Or(somesql.LangEN, "has_video", "=", true))),
			expectedSQL:    `SELECT data_en FROM repo WHERE status=? AND ("data_en"->>'badge'=? OR ("data_en"->>'has_video')::BOOLEAN=?)`,
			expectedValues: []interface{}{"published", "video", true},
		},
		{
			name:           "WHERE field = ... AND (... OR ...) [2]",
			query:          somesql.NewQuery().Select("data").Where(somesql.Or(somesql.LangEN, "status", "=", []string{"published"})).Where(somesql.AndGroup(somesql.And(somesql.LangEN, "badge", "=", "video"), somesql.Or(somesql.LangEN, "badge", "=", "audio"))),
			expectedSQL:    `SELECT data_en FROM repo WHERE status=? AND ("data_en"->>'badge'=? OR "data_en"->>'badge'=?)`,
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

	// Examples
	// SELECT * FROM repo WHERE type IN (SELECT type FROM repo WHERE data->>'author_id' = '002fd6b1-f715-4875-838b-1546f27327df');
	// SELECT * FROM repo WHERE data->>'author_id' IN (SELECT data->>'author_id' FROM repo WHERE data->>'author_id' = '002fd6b1-f715-4875-838b-1546f27327df');

	tests := []testCase{
		// AndInQuery -> somesql.NewQuery().Select("type", "slug") : cannot have more than 1 field in subquery (throw an error?)
		{
			name:           "AndInQuery",
			query:          somesql.NewQuery().Select("data").Where(somesql.AndInQuery("type", somesql.NewQuery().Select("type", "slug").Where(somesql.And(somesql.LangEN, "id", "=", "002fd6b1-f715-4875-838b-1546f27327df")))),
			expectedSQL:    "",
			expectedValues: []interface{}{},
		},
		{
			name:           "AndInQuery",
			query:          somesql.NewQuery().Select("data").Where(somesql.AndInQuery("type", somesql.NewQuery().Select("type").Where(somesql.And(somesql.LangEN, "id", "=", "002fd6b1-f715-4875-838b-1546f27327df")))),
			expectedSQL:    `SELECT data_en FROM repo WHERE type IN (SELECT type FROM repo WHERE id =?)`,
			expectedValues: []interface{}{"002fd6b1-f715-4875-838b-1546f27327df"},
		},
		{
			name:           "AndInQuery",
			query:          somesql.NewQuery().Select("data").Where(somesql.AndInQuery("author_id", somesql.NewQuery().Select("author_id").Where(somesql.And(somesql.LangEN, "id", "=", "002fd6b1-f715-4875-838b-1546f27327df")))),
			expectedSQL:    `SELECT data_en FROM repo WHERE "data_en"->>'author_id' IN (SELECT data_en->'author_id' FROM repo WHERE id =?)`,
			expectedValues: []interface{}{"002fd6b1-f715-4875-838b-1546f27327df"},
		},
		{
			name:           "AndInQuery",
			query:          somesql.NewQuery().Select("data").Where(somesql.AndIn(somesql.LangEN, "id", []string{"A", "B", "C"})).Where(somesql.AndInQuery("author_id", somesql.NewQuery().Select("author_id").Where(somesql.And(somesql.LangEN, "id", "=", "002fd6b1-f715-4875-838b-1546f27327df")))),
			expectedSQL:    `SELECT data_en FROM repo WHERE id IN (?,?,?) AND "data_en"->>'author_id' IN (SELECT data_en->'author_id' FROM repo WHERE id =?)`,
			expectedValues: []interface{}{"A", "B", "C", "002fd6b1-f715-4875-838b-1546f27327df"},
		},
		{
			name:           "AndNotInQuery",
			query:          somesql.NewQuery().Select("data").Where(somesql.AndIn(somesql.LangEN, "id", []string{"A", "B", "C"})).Where(somesql.AndNotInQuery("author_id", somesql.NewQuery().Select("author_id").Where(somesql.And(somesql.LangEN, "id", "=", "002fd6b1-f715-4875-838b-1546f27327df")))),
			expectedSQL:    `SELECT data_en FROM repo WHERE id IN (?,?,?) AND "data_en"->>'author_id' NOT IN (SELECT data_en->'author_id' FROM repo WHERE id =?)`,
			expectedValues: []interface{}{"A", "B", "C", "002fd6b1-f715-4875-838b-1546f27327df"},
		},

		// OrInQuery -> somesql.NewQuery().Select("type", "slug") : cannot have more than 1 field in subquery (throw an error?)
		{
			name:           "OrInQuery",
			query:          somesql.NewQuery().Select("data").Where(somesql.AndInQuery("type", somesql.NewQuery().Select("type", "slug").Where(somesql.And(somesql.LangEN, "id", "=", "002fd6b1-f715-4875-838b-1546f27327df")))),
			expectedSQL:    "",
			expectedValues: []interface{}{},
		},
		{
			name:           "OrInQuery",
			query:          somesql.NewQuery().Select("data").Where(somesql.OrInQuery("type", somesql.NewQuery().Select("type").Where(somesql.And(somesql.LangEN, "id", "=", "002fd6b1-f715-4875-838b-1546f27327df")))),
			expectedSQL:    `SELECT data_en FROM repo WHERE type IN (SELECT type FROM repo WHERE id =?)`,
			expectedValues: []interface{}{"002fd6b1-f715-4875-838b-1546f27327df"},
		},
		{
			name:           "OrInQuery",
			query:          somesql.NewQuery().Select("data").Where(somesql.OrInQuery("author_id", somesql.NewQuery().Select("author_id").Where(somesql.And(somesql.LangEN, "id", "=", "002fd6b1-f715-4875-838b-1546f27327df")))),
			expectedSQL:    `SELECT data_en FROM repo WHERE "data_en"->>'author_id' IN (SELECT data_en->'author_id' FROM repo WHERE id =?)`,
			expectedValues: []interface{}{"002fd6b1-f715-4875-838b-1546f27327df"},
		},
		{
			name:           "OrInQuery",
			query:          somesql.NewQuery().Select("data").Where(somesql.AndIn(somesql.LangEN, "id", []string{"A", "B", "C"})).Where(somesql.OrInQuery("author_id", somesql.NewQuery().Select("author_id").Where(somesql.And(somesql.LangEN, "id", "=", "002fd6b1-f715-4875-838b-1546f27327df")))),
			expectedSQL:    `SELECT data_en FROM repo WHERE id IN (?,?,?) OR "data_en"->>'author_id' IN (SELECT data_en->'author_id' FROM repo WHERE id =?)`,
			expectedValues: []interface{}{"A", "B", "C", "002fd6b1-f715-4875-838b-1546f27327df"},
		},
		{
			name:           "OrNotInQuery",
			query:          somesql.NewQuery().Select("data").Where(somesql.AndIn(somesql.LangEN, "id", []string{"A", "B", "C"})).Where(somesql.OrNotInQuery("author_id", somesql.NewQuery().Select("author_id").Where(somesql.And(somesql.LangEN, "id", "=", "002fd6b1-f715-4875-838b-1546f27327df")))),
			expectedSQL:    `SELECT data_en FROM repo WHERE id IN (?,?,?) OR "data_en"->>'author_id' NOT IN (SELECT data_en->'author_id' FROM repo WHERE id =?)`,
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

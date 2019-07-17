package somesql_test

import (
	"fmt"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"go.lsl.digital/lardwaz/somesql"
)

func TestQuery_AsSQL_Select(t *testing.T) {
	type testCase struct {
		name        string
		query       *somesql.Select
		expectedSQL string

		checkValues    bool
		expectedValues []interface{}
	}

	tests := []testCase{
		// Select ALL
		{
			name:        "SELECT *",
			query:       somesql.NewSelect(),
			expectedSQL: `SELECT "id", "created_at", "updated_at", "owner_id", "type", "data_en", "relations" FROM repo LIMIT 10`,
		},
		{
			name:        "SELECT * NO LIMIT",
			query:       somesql.NewSelect().Limit(0),
			expectedSQL: `SELECT "id", "created_at", "updated_at", "owner_id", "type", "data_en", "relations" FROM repo`,
		},
		{
			name:        "SELECT * LIMIT 30",
			query:       somesql.NewSelect().Limit(30),
			expectedSQL: `SELECT "id", "created_at", "updated_at", "owner_id", "type", "data_en", "relations" FROM repo LIMIT 30`,
		},
		{
			name:        "SELECT * OFFSET 10",
			query:       somesql.NewSelect().Offset(10),
			expectedSQL: `SELECT "id", "created_at", "updated_at", "owner_id", "type", "data_en", "relations" FROM repo LIMIT 10 OFFSET 10`,
		},
		{
			name:        "SELECT * LIMIT 30 OFFSET 20",
			query:       somesql.NewSelect().Limit(30).Offset(20),
			expectedSQL: `SELECT "id", "created_at", "updated_at", "owner_id", "type", "data_en", "relations" FROM repo LIMIT 30 OFFSET 20`,
		},
		{
			name:        "SELECT * (langEN)",
			query:       somesql.NewSelect(),
			expectedSQL: `SELECT "id", "created_at", "updated_at", "owner_id", "type", "data_en", "relations" FROM repo LIMIT 10`,
		},
		{
			name:        "SELECT * (langFR)",
			query:       somesql.NewSelectLang(somesql.LangFR),
			expectedSQL: `SELECT "id", "created_at", "updated_at", "owner_id", "type", "data_fr", "relations" FROM repo LIMIT 10`,
		},
		// Select some pre-defined fields
		{
			name:        "SELECT id, type, data",
			query:       somesql.NewSelect().Fields("id", "type", "data"),
			expectedSQL: `SELECT "id", "type", "data_en" FROM repo LIMIT 10`,
		},
		{
			name:        "SELECT (EMPTY)",
			query:       somesql.NewSelect().Fields(),
			expectedSQL: `SELECT "id", "created_at", "updated_at", "owner_id", "type", "data_en", "relations" FROM repo LIMIT 10`,
		},
		{
			name:        "SELECT id, type, data (LangEN)",
			query:       somesql.NewSelect().Fields("id", "type", "data"),
			expectedSQL: `SELECT "id", "type", "data_en" FROM repo LIMIT 10`,
		},
		{
			name:        "SELECT id, type, data (LangFR)",
			query:       somesql.NewSelectLang(somesql.LangFR).Fields("id", "type", "data"),
			expectedSQL: `SELECT "id", "type", "data_fr" FROM repo LIMIT 10`,
		},
		// Select pre-defined fields and json attributes ('data_en'/'data_fr') from data_*
		{
			name:        "SELECT id, type, data_en",
			query:       somesql.NewSelect().Fields("id", "type", "data.data_en"),
			expectedSQL: `SELECT "id", "type", json_build_object('data_en', "data_en"->'data_en') "data" FROM repo LIMIT 10`,
		},
		{
			name:        "SELECT id, type, data_en (LangEN)",
			query:       somesql.NewSelect().Fields("id", "type", "data.data_en"),
			expectedSQL: `SELECT "id", "type", json_build_object('data_en', "data_en"->'data_en') "data" FROM repo LIMIT 10`,
		},
		{
			name:        "SELECT id, type, data_en (LangFR)",
			query:       somesql.NewSelectLang(somesql.LangFR).Fields("id", "type", "data.data_en"),
			expectedSQL: `SELECT "id", "type", json_build_object('data_en', "data_fr"->'data_en') "data" FROM repo LIMIT 10`,
		},
		{
			name:        "SELECT id, type, data_fr",
			query:       somesql.NewSelect().Fields("id", "type", "data.data_fr"),
			expectedSQL: `SELECT "id", "type", json_build_object('data_fr', "data_en"->'data_fr') "data" FROM repo LIMIT 10`,
		},
		{
			name:        "SELECT id, type, data_fr (LangEN)",
			query:       somesql.NewSelect().Fields("id", "type", "data.data_fr"),
			expectedSQL: `SELECT "id", "type", json_build_object('data_fr', "data_en"->'data_fr') "data" FROM repo LIMIT 10`,
		},
		{
			name:        "SELECT id, type, data_fr (LangFR)",
			query:       somesql.NewSelectLang(somesql.LangFR).Fields("id", "type", "data.data_fr"),
			expectedSQL: `SELECT "id", "type", json_build_object('data_fr', "data_fr"->'data_fr') "data" FROM repo LIMIT 10`,
		},
		// Select pre-defined fields and json attributes (any other) from data_*
		{
			name:        "SELECT id, type, data_en->'body'",
			query:       somesql.NewSelect().Fields("id", "type", "data.body"),
			expectedSQL: `SELECT "id", "type", json_build_object('body', "data_en"->'body') "data" FROM repo LIMIT 10`,
		},
		{
			name:        "SELECT id, type, data_en->'body' (LangEN)",
			query:       somesql.NewSelect().Fields("id", "type", "data.body"),
			expectedSQL: `SELECT "id", "type", json_build_object('body', "data_en"->'body') "data" FROM repo LIMIT 10`,
		},
		{
			name:        "SELECT id, type, data_fr->'body' (LangFR)",
			query:       somesql.NewSelectLang(somesql.LangFR).Fields("id", "type", "data.body"),
			expectedSQL: `SELECT "id", "type", json_build_object('body', "data_fr"->'body') "data" FROM repo LIMIT 10`,
		},
		// Select pre-defined fields and json attributes (any other + compound) from data_*
		{
			name:        "SELECT id, type, data_en->'body', data_en->'author_id'",
			query:       somesql.NewSelect().Fields("id", "type", "data.body", "data.author_id"),
			expectedSQL: `SELECT "id", "type", json_build_object('body', "data_en"->'body', 'author_id', "data_en"->'author_id') "data" FROM repo LIMIT 10`,
		},
		{
			name:        "SELECT id, type, data_en->'body', data_en->'author_id' (LangEN)",
			query:       somesql.NewSelect().Fields("id", "type", "data.body", "data.author_id"),
			expectedSQL: `SELECT "id", "type", json_build_object('body', "data_en"->'body', 'author_id', "data_en"->'author_id') "data" FROM repo LIMIT 10`,
		},
		{
			name:        "SELECT id, type, data_fr->'body', data_fr->'author_id' (LangFR)",
			query:       somesql.NewSelectLang(somesql.LangFR).Fields("id", "type", "data.body", "data.author_id"),
			expectedSQL: `SELECT "id", "type", json_build_object('body', "data_fr"->'body', 'author_id', "data_fr"->'author_id') "data" FROM repo LIMIT 10`,
		},
		// SELECT with conditions
		{
			name:           "SELECT * with condition",
			query:          somesql.NewSelect().Where(somesql.And(somesql.LangEN, "id", "=", "uuid")),
			expectedSQL:    `SELECT "id", "created_at", "updated_at", "owner_id", "type", "data_en", "relations" FROM repo WHERE "id" = $1 LIMIT 10`,
			checkValues:    true,
			expectedValues: []interface{}{"uuid"},
		},
		// SELECT relations
		{
			name:        "SELECT id, relations->rel",
			query:       somesql.NewSelect().Fields("id", "relations.author", "relations.tags"),
			expectedSQL: `SELECT "id", json_build_object('author', "relations"->'author', 'tags', "relations"->'tags') "relations" FROM repo LIMIT 10`,
		},
		{
			name:           "SELECT * + relations->rel, with conditions",
			query:          somesql.NewSelect().Fields("relations.author", "relations.tags").Where(somesql.And(somesql.LangEN, "id", "=", "uuid")),
			expectedSQL:    `SELECT json_build_object('author', "relations"->'author', 'tags', "relations"->'tags') "relations" FROM repo WHERE "id" = $1 LIMIT 10`,
			checkValues:    true,
			expectedValues: []interface{}{"uuid"},
		},
		{
			name:           "SELECT relations->rel only, with conditions",
			query:          somesql.NewSelect().Fields("relations.author", "relations.tags").Where(somesql.And(somesql.LangEN, "id", "=", "uuid")),
			expectedSQL:    `SELECT json_build_object('author', "relations"->'author', 'tags', "relations"->'tags') "relations" FROM repo WHERE "id" = $1 LIMIT 10`,
			checkValues:    true,
			expectedValues: []interface{}{"uuid"},
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.query.ToSQL()
			gotSQL, gotValues := tt.query.GetSQL(), tt.query.GetValues()

			assert.Equal(t, tt.expectedSQL, gotSQL, fmt.Sprintf("Fields %03d :: %s", i+1, tt.name))

			if tt.checkValues {
				assert.Equal(t, tt.expectedValues, gotValues, fmt.Sprintf("Fields %03d :: invalid values :: %s", i+1, tt.name))
			}
		})
	}
}

func TestQuery_AsSQL_Insert(t *testing.T) {
	type testCase struct {
		name           string
		query          *somesql.Insert
		expectedSQL    string
		expectedValues []interface{}
	}

	tests := []testCase{
		// Insert
		{
			name:           "INSERT defaults",
			query:          somesql.NewInsert().Fields(somesql.NewFields().UseDefaults().ID("1").CreatedAt(time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)).UpdatedAt(time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC))),
			expectedSQL:    `INSERT INTO repo ("id", "created_at", "updated_at", "owner_id", "type", "data_en", "relations") VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			expectedValues: []interface{}{"1", time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC), time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC), uuid.Nil.String(), "", "{}", "{}"},
		},
		{
			name:           "INSERT no defaults",
			query:          somesql.NewInsert().Fields(somesql.NewFields().ID("1").Type("entityA")),
			expectedSQL:    `INSERT INTO repo ("id", "type") VALUES ($1, $2)`,
			expectedValues: []interface{}{"1", "entityA"},
		},
		{
			name:           "INSERT no default + 1 relation",
			query:          somesql.NewInsert().Fields(somesql.NewFields().Type("entityA").Set("relations.tags", []string{"a", "b", "c"})),
			expectedSQL:    `INSERT INTO repo ("type", "relations") VALUES ($1, $2)`,
			expectedValues: []interface{}{"entityA", `{"tags":["a","b","c"]}`},
		},
		{
			name:           "INSERT no default + 2 relations",
			query:          somesql.NewInsert().Fields(somesql.NewFields().Type("entityA").Set("relations.author", []string{"x"}).Set("relations.tags", []string{"a", "b", "c"})),
			expectedSQL:    `INSERT INTO repo ("type", "relations") VALUES ($1, $2)`,
			expectedValues: []interface{}{"entityA", `{"author":["x"],"tags":["a","b","c"]}`},
		},
		{
			name:           "INSERT no default + 1 relation one value",
			query:          somesql.NewInsert().Fields(somesql.NewFields().Type("entityA").Set("relations.tags", "a")),
			expectedSQL:    `INSERT INTO repo ("type", "relations") VALUES ($1, $2)`,
			expectedValues: []interface{}{"entityA", `{"tags":["a"]}`},
		},
		{
			name:           "INSERT no default + 1 relation multiple seperate values",
			query:          somesql.NewInsert().Fields(somesql.NewFields().Type("entityA").Set("relations.tags", "a").Add("relations.tags", "b").Add("relations.tags", []string{"c"})),
			expectedSQL:    `INSERT INTO repo ("type", "relations") VALUES ($1, $2)`,
			expectedValues: []interface{}{"entityA", `{"tags":["a","b","c"]}`},
		},
		{
			name:           "INSERT no default + 1 relation multiple seperate values 2",
			query:          somesql.NewInsert().Fields(somesql.NewFields().Type("entityA").Set("relations.tags", []string{"a"}).Add("relations.tags", "b").Add("relations.tags", []string{"c"})),
			expectedSQL:    `INSERT INTO repo ("type", "relations") VALUES ($1, $2)`,
			expectedValues: []interface{}{"entityA", `{"tags":["a","b","c"]}`},
		},
		{
			name:           "INSERT no default + 1 relation multiple seperate values 3",
			query:          somesql.NewInsert().Fields(somesql.NewFields().Type("entityA").Add("relations.tags", "a").Add("relations.tags", "b").Add("relations.tags", "c")),
			expectedSQL:    `INSERT INTO repo ("type", "relations") VALUES ($1, $2)`,
			expectedValues: []interface{}{"entityA", `{"tags":["a","b","c"]}`},
		},
		{
			name:           "INSERT no default + 1 relation multiple seperate values overwrite 4",
			query:          somesql.NewInsert().Fields(somesql.NewFields().Type("entityA").Add("relations.tags", "a").Add("relations.tags", "b").Set("relations.tags", "c")),
			expectedSQL:    `INSERT INTO repo ("type", "relations") VALUES ($1, $2)`,
			expectedValues: []interface{}{"entityA", `{"tags":["c"]}`},
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.query.ToSQL()
			gotSQL, gotValues := tt.query.GetSQL(), tt.query.GetValues()

			assert.Equal(t, tt.expectedSQL, gotSQL, fmt.Sprintf("Fields %03d :: invalid sql :: %s", i+1, tt.name))
			assert.Equal(t, tt.expectedValues, gotValues, fmt.Sprintf("Fields %03d :: invalid values :: %s", i+1, tt.name))
		})
	}
}

func TestQuery_AsSQL_Update(t *testing.T) {
	type testCase struct {
		name           string
		query          *somesql.Update
		expectedSQL    string
		expectedValues []interface{}
	}

	tests := []testCase{
		// Update
		{
			name:           "UPDATE meta fields",
			query:          somesql.NewUpdate().Fields(somesql.NewFields().ID("1").Type("entityA")),
			expectedSQL:    `UPDATE repo SET "id" = $1, "type" = $2`,
			expectedValues: []interface{}{"1", "entityA"},
		},
		{
			name:           "UPDATE data fields",
			query:          somesql.NewUpdate().Fields(somesql.NewFields().Set("data.body", "body value").Set("data.author_id", "123")),
			expectedSQL:    `UPDATE repo SET "data_en" = jsonb_build_object('body', $1::text, 'author_id', $2::text)::JSONB`,
			expectedValues: []interface{}{"body value", "123"},
		},
		{
			name:           "UPDATE data fields (LangFR)",
			query:          somesql.NewUpdateLang(somesql.LangFR).Fields(somesql.NewFields().Set("data.body", "body value").Set("data.author_id", "123")),
			expectedSQL:    `UPDATE repo SET "data_fr" = jsonb_build_object('body', $1::text, 'author_id', $2::text)::JSONB`,
			expectedValues: []interface{}{"body value", "123"},
		},
		{
			name:           "UPDATE meta + data fields",
			query:          somesql.NewUpdate().Fields(somesql.NewFields().ID("1").Type("entityA").Set("data.body", "body value").Set("data.author_id", "123")),
			expectedSQL:    `UPDATE repo SET "id" = $1, "type" = $2, "data_en" = jsonb_build_object('body', $3::text, 'author_id', $4::text)::JSONB`,
			expectedValues: []interface{}{"1", "entityA", "body value", "123"},
		},
		{
			name:           "UPDATE meta + data fields (LangFR)",
			query:          somesql.NewUpdateLang(somesql.LangFR).Fields(somesql.NewFields().ID("1").Type("entityA").Set("data.body", "body value").Set("data.author_id", "123")),
			expectedSQL:    `UPDATE repo SET "id" = $1, "type" = $2, "data_fr" = jsonb_build_object('body', $3::text, 'author_id', $4::text)::JSONB`,
			expectedValues: []interface{}{"1", "entityA", "body value", "123"},
		},
		{
			name:           "UPDATE meta + data fields conditions",
			query:          somesql.NewUpdate().Fields(somesql.NewFields().ID("1").Type("entityA").Set("data.body", "body value").Set("data.author_id", "123")).Where(somesql.And(somesql.LangEN, "id", "=", "234")),
			expectedSQL:    `UPDATE repo SET "id" = $1, "type" = $2, "data_en" = jsonb_build_object('body', $3::text, 'author_id', $4::text)::JSONB WHERE "id" = $5`,
			expectedValues: []interface{}{"1", "entityA", "body value", "123", "234"},
		},
		{
			name:           "UPDATE meta + data fields conditions 2",
			query:          somesql.NewUpdate().Fields(somesql.NewFields().ID("1").Type("entityA").Set("data.body", "body value").Set("data.author_id", "123")).Where(somesql.And(somesql.LangEN, "data.author_id", "=", "234")),
			expectedSQL:    `UPDATE repo SET "id" = $1, "type" = $2, "data_en" = jsonb_build_object('body', $3::text, 'author_id', $4::text)::JSONB WHERE "data_en"->>'author_id' = $5`,
			expectedValues: []interface{}{"1", "entityA", "body value", "123", "234"},
		},
		// Update relations : add
		{
			name:           "UPDATE add 1 relation only",
			query:          somesql.NewUpdate().Fields(somesql.NewFields().Add("relations.tags", []string{"a", "b"})),
			expectedSQL:    `UPDATE repo SET "relations" = relAdd.relations FROM (SELECT (("relations" - 'tags') || JSONB_BUILD_OBJECT('tags', COALESCE("relations"->'tags' || $1::JSONB, $2::JSONB))) "relations" FROM repo) relAdd`,
			expectedValues: []interface{}{`["a","b"]`, `["a","b"]`},
		},
		{
			name:           "UPDATE add 1 or more relations only",
			query:          somesql.NewUpdate().Fields(somesql.NewFields().Add("relations.tags", []string{"a", "b"}).Add("relations.author", []string{"x"})),
			expectedSQL:    `UPDATE repo SET "relations" = relAdd.relations FROM (SELECT (("relations" - 'tags') || ("relations" - 'author') || JSONB_BUILD_OBJECT('tags', COALESCE("relations"->'tags' || $1::JSONB, $2::JSONB), 'author', COALESCE("relations"->'author' || $3::JSONB, $4::JSONB))) "relations" FROM repo) relAdd`,
			expectedValues: []interface{}{`["a","b"]`, `["a","b"]`, `["x"]`, `["x"]`},
		},
		{
			name:           "UPDATE add 1 relation only + conditions",
			query:          somesql.NewUpdate().Fields(somesql.NewFields().Add("relations.author", []string{"x"})).Where(somesql.And(somesql.LangEN, "id", "=", "uuid")),
			expectedSQL:    `UPDATE repo SET "relations" = relAdd.relations FROM (SELECT (("relations" - 'author') || JSONB_BUILD_OBJECT('author', COALESCE("relations"->'author' || $1::JSONB, $2::JSONB))) "relations" FROM repo WHERE "id" = $3) relAdd WHERE "id" = $4`,
			expectedValues: []interface{}{`["x"]`, `["x"]`, "uuid", "uuid"},
		},
		// // Update relations : remove
		{
			name:           "UPDATE remove 1 relation only",
			query:          somesql.NewUpdate().Fields(somesql.NewFields().Remove("relations.tags", []string{"a", "b"})),
			expectedSQL:    `UPDATE repo SET "relations" = updates.updRel FROM (SELECT (("relations" - 'tags') || JSONB_BUILD_OBJECT('tags', JSONB_AGG(tagsUpd))) "updatedRel" FROM (SELECT "relations", JSONB_ARRAY_ELEMENTS_TEXT("relations"->'tags') tagsUpd FROM repo) expandedValues WHERE tagsUpd NOT IN ($1) GROUP BY "relations") updates`,
			expectedValues: []interface{}{`["a","b"]`},
		},
		{
			name:           "UPDATE remove 1 or more relations only",
			query:          somesql.NewUpdate().Fields(somesql.NewFields().Remove("relations.tags", []string{"a", "b"}).Remove("relations.author", []string{"x"})),
			expectedSQL:    `UPDATE repo SET "relations" = updates.updRel FROM (SELECT (("relations" - 'tags') || ("relations" - 'author') || JSONB_BUILD_OBJECT('tags', JSONB_AGG(tagsUpd), 'author', JSONB_AGG(authorUpd))) "updatedRel" FROM (SELECT "relations", JSONB_ARRAY_ELEMENTS_TEXT("relations"->'tags') tagsUpd, JSONB_ARRAY_ELEMENTS_TEXT("relations"->'author') authorUpd FROM repo) expandedValues WHERE tagsUpd NOT IN ($1) AND authorUpd NOT IN ($2) GROUP BY "relations") updates`,
			expectedValues: []interface{}{`["a","b"]`, `["x"]`},
		},
		{
			name:           "UPDATE remove 1 relation only + conditions",
			query:          somesql.NewUpdate().Fields(somesql.NewFields().Remove("relations.author", []string{"x"})).Where(somesql.And(somesql.LangEN, "id", "=", "uuid")),
			expectedSQL:    `UPDATE repo SET "relations" = updates.updRel FROM (SELECT (("relations" - 'author') || JSONB_BUILD_OBJECT('author', JSONB_AGG(authorUpd))) "updatedRel" FROM (SELECT "relations", JSONB_ARRAY_ELEMENTS_TEXT("relations"->'author') authorUpd FROM repo WHERE "id" = $1) expandedValues WHERE authorUpd NOT IN ($2) GROUP BY "relations") updates WHERE "id" = $3`,
			expectedValues: []interface{}{"uuid", `["x"]`, "uuid"},
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.query.ToSQL()
			gotSQL, gotValues := tt.query.GetSQL(), tt.query.GetValues()

			assert.Equal(t, tt.expectedSQL, gotSQL, fmt.Sprintf("Fields %03d :: invalid sql :: %s", i+1, tt.name))
			assert.Equal(t, tt.expectedValues, gotValues, fmt.Sprintf("Fields %03d :: invalid values :: %s", i+1, tt.name))
		})
	}
}

func TestQuery_AsSQL_Delete(t *testing.T) {
	type testCase struct {
		name           string
		query          *somesql.Delete
		expectedSQL    string
		checkValues    bool
		expectedValues []interface{}
	}

	tests := []testCase{
		{
			name:        "DELETE * NO LIMIT",
			query:       somesql.NewDelete(),
			expectedSQL: `DELETE FROM repo`,
		},
		{
			name:        "DELETE * LIMIT",
			query:       somesql.NewDelete().Limit(20),
			expectedSQL: `DELETE FROM repo LIMIT 20`,
		},
		{
			name:        "DELETE * OFFSET 10",
			query:       somesql.NewDelete().Offset(10),
			expectedSQL: `DELETE FROM repo OFFSET 10`,
		},
		{
			name:        "DELETE * LIMIT 20 OFFSET 10",
			query:       somesql.NewDelete().Limit(20).Offset(10),
			expectedSQL: `DELETE FROM repo LIMIT 20 OFFSET 10`,
		},
		{
			name:           "DELETE with condition",
			query:          somesql.NewDelete().Where(somesql.And(somesql.LangEN, "id", "=", "uuid")),
			expectedSQL:    `DELETE FROM repo WHERE "id" = $1`,
			checkValues:    true,
			expectedValues: []interface{}{"uuid"},
		},
		{
			name:           "DELETE with conditions + relations",
			query:          somesql.NewDelete().Where(somesql.And("", "relations.article", "=", "uuid")),
			expectedSQL:    `DELETE FROM repo WHERE ("relations" @> '{"article":$1}'::JSONB)`,
			checkValues:    true,
			expectedValues: []interface{}{"uuid"},
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.query.ToSQL()
			gotSQL, gotValues := tt.query.GetSQL(), tt.query.GetValues()

			assert.Equal(t, tt.expectedSQL, gotSQL, fmt.Sprintf("Fields %03d :: invalid sql :: %s", i+1, tt.name))
			if tt.checkValues {
				assert.Equal(t, tt.expectedValues, gotValues, fmt.Sprintf("Fields %03d :: invalid values :: %s", i+1, tt.name))
			}
		})
	}
}

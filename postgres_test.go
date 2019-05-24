package somesql_test

import (
	"fmt"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"go.lsl.digital/gocipe/somesql"
)

func TestQuery_AsSQL_Fields(t *testing.T) {
	type testCase struct {
		name        string
		query       somesql.Query
		expectedSQL string

		checkValues    bool
		expectedValues []interface{}
	}

	tests := []testCase{
		// Select ALL
		{
			name:        "SELECT *",
			query:       somesql.NewQuery(),
			expectedSQL: `SELECT "id", "created_at", "updated_at", "owner_id", "status", "type", "data_en" FROM repo LIMIT 10`,
		},
		{
			name:        "SELECT * NO LIMIT",
			query:       somesql.NewQuery().SetLimit(0),
			expectedSQL: `SELECT "id", "created_at", "updated_at", "owner_id", "status", "type", "data_en" FROM repo`,
		},
		{
			name:        "SELECT * LIMIT 30",
			query:       somesql.NewQuery().SetLimit(30),
			expectedSQL: `SELECT "id", "created_at", "updated_at", "owner_id", "status", "type", "data_en" FROM repo LIMIT 30`,
		},
		{
			name:        "SELECT * OFFSET 10",
			query:       somesql.NewQuery().SetOffset(10),
			expectedSQL: `SELECT "id", "created_at", "updated_at", "owner_id", "status", "type", "data_en" FROM repo LIMIT 10 OFFSET 10`,
		},
		{
			name:        "SELECT * LIMIT 30 OFFSET 20",
			query:       somesql.NewQuery().SetLimit(30).SetOffset(20),
			expectedSQL: `SELECT "id", "created_at", "updated_at", "owner_id", "status", "type", "data_en" FROM repo LIMIT 30 OFFSET 20`,
		},
		{
			name:        "SELECT * (langEN)",
			query:       somesql.NewQuery().SetLang(somesql.LangEN),
			expectedSQL: `SELECT "id", "created_at", "updated_at", "owner_id", "status", "type", "data_en" FROM repo LIMIT 10`,
		},
		{
			name:        "SELECT * (langFR)",
			query:       somesql.NewQuery().SetLang(somesql.LangFR),
			expectedSQL: `SELECT "id", "created_at", "updated_at", "owner_id", "status", "type", "data_fr" FROM repo LIMIT 10`,
		},
		// Select some pre-defined fields
		{
			name:        "SELECT id, type, data",
			query:       somesql.NewQuery().Select("id", "type", "data"),
			expectedSQL: `SELECT "id", "type", "data_en" FROM repo LIMIT 10`,
		},
		{
			name:        "SELECT (EMPTY)",
			query:       somesql.NewQuery().Select(),
			expectedSQL: `SELECT "id", "created_at", "updated_at", "owner_id", "status", "type", "data_en" FROM repo LIMIT 10`,
		},
		{
			name:        "SELECT id, type, data (LangEN)",
			query:       somesql.NewQuery().Select("id", "type", "data").SetLang(somesql.LangEN),
			expectedSQL: `SELECT "id", "type", "data_en" FROM repo LIMIT 10`,
		},
		{
			name:        "SELECT id, type, data (LangFR)",
			query:       somesql.NewQuery().Select("id", "type", "data").SetLang(somesql.LangFR),
			expectedSQL: `SELECT "id", "type", "data_fr" FROM repo LIMIT 10`,
		},
		// Select pre-defined fields and json attributes ('data_en'/'data_fr') from data_*
		{
			name:        "SELECT id, type, data_en",
			query:       somesql.NewQuery().Select("id", "type", "data_en"),
			expectedSQL: `SELECT "id", "type", json_build_object('data_en', "data_en"->'data_en') "data" FROM repo LIMIT 10`,
		},
		{
			name:        "SELECT id, type, data_en (LangEN)",
			query:       somesql.NewQuery().Select("id", "type", "data_en").SetLang(somesql.LangEN),
			expectedSQL: `SELECT "id", "type", json_build_object('data_en', "data_en"->'data_en') "data" FROM repo LIMIT 10`,
		},
		{
			name:        "SELECT id, type, data_en (LangFR)",
			query:       somesql.NewQuery().Select("id", "type", "data_en").SetLang(somesql.LangFR),
			expectedSQL: `SELECT "id", "type", json_build_object('data_en', "data_fr"->'data_en') "data" FROM repo LIMIT 10`,
		},
		{
			name:        "SELECT id, type, data_fr",
			query:       somesql.NewQuery().Select("id", "type", "data_fr"),
			expectedSQL: `SELECT "id", "type", json_build_object('data_fr', "data_en"->'data_fr') "data" FROM repo LIMIT 10`,
		},
		{
			name:        "SELECT id, type, data_fr (LangEN)",
			query:       somesql.NewQuery().Select("id", "type", "data_fr").SetLang(somesql.LangEN),
			expectedSQL: `SELECT "id", "type", json_build_object('data_fr', "data_en"->'data_fr') "data" FROM repo LIMIT 10`,
		},
		{
			name:        "SELECT id, type, data_fr (LangFR)",
			query:       somesql.NewQuery().Select("id", "type", "data_fr").SetLang(somesql.LangFR),
			expectedSQL: `SELECT "id", "type", json_build_object('data_fr', "data_fr"->'data_fr') "data" FROM repo LIMIT 10`,
		},
		// Select pre-defined fields and json attributes (any other) from data_*
		{
			name:        "SELECT id, type, data_en->'body'",
			query:       somesql.NewQuery().Select("id", "type", "body"),
			expectedSQL: `SELECT "id", "type", json_build_object('body', "data_en"->'body') "data" FROM repo LIMIT 10`,
		},
		{
			name:        "SELECT id, type, data_en->'body' (LangEN)",
			query:       somesql.NewQuery().Select("id", "type", "body").SetLang(somesql.LangEN),
			expectedSQL: `SELECT "id", "type", json_build_object('body', "data_en"->'body') "data" FROM repo LIMIT 10`,
		},
		{
			name:        "SELECT id, type, data_fr->'body' (LangFR)",
			query:       somesql.NewQuery().Select("id", "type", "body").SetLang(somesql.LangFR),
			expectedSQL: `SELECT "id", "type", json_build_object('body', "data_fr"->'body') "data" FROM repo LIMIT 10`,
		},
		// Select pre-defined fields and json attributes (any other + compound) from data_*
		{
			name:        "SELECT id, type, data_en->'body', data_en->'author_id'",
			query:       somesql.NewQuery().Select("id", "type", "body", "author_id"),
			expectedSQL: `SELECT "id", "type", json_build_object('body', "data_en"->'body', 'author_id', "data_en"->'author_id') "data" FROM repo LIMIT 10`,
		},
		{
			name:        "SELECT id, type, data_en->'body', data_en->'author_id' (LangEN)",
			query:       somesql.NewQuery().Select("id", "type", "body", "author_id").SetLang(somesql.LangEN),
			expectedSQL: `SELECT "id", "type", json_build_object('body', "data_en"->'body', 'author_id', "data_en"->'author_id') "data" FROM repo LIMIT 10`,
		},
		{
			name:        "SELECT id, type, data_fr->'body', data_fr->'author_id' (LangFR)",
			query:       somesql.NewQuery().Select("id", "type", "body", "author_id").SetLang(somesql.LangFR),
			expectedSQL: `SELECT "id", "type", json_build_object('body', "data_fr"->'body', 'author_id', "data_fr"->'author_id') "data" FROM repo LIMIT 10`,
		},
		// SELECT with conditions
		{
			name:           "SELECT * with condition",
			query:          somesql.NewQuery().Where(somesql.And(somesql.LangEN, "id", "=", "uuid")),
			expectedSQL:    `SELECT "id", "created_at", "updated_at", "owner_id", "status", "type", "data_en" FROM repo WHERE "id"=$1 LIMIT 10`,
			checkValues:    true,
			expectedValues: []interface{}{"uuid"},
		},
		// SELECT relations
		{
			name:        "SELECT id, relations->rel",
			query:       somesql.NewQuery().Select("id").SelectRel("author", "tags"),
			expectedSQL: `SELECT "id", json_build_object('author', "relations"->'author', 'tags', "relations"->'tags') "relations" FROM repo LIMIT 10`,
		},
		{
			name:           "SELECT * + relations->rel, with conditions",
			query:          somesql.NewQuery().SelectRel("author", "tags").Where(somesql.And(somesql.LangEN, "id", "=", "uuid")),
			expectedSQL:    `SELECT "id", "created_at", "updated_at", "owner_id", "status", "type", "data_en", json_build_object('author', "relations"->'author', 'tags', "relations"->'tags') "relations" FROM repo WHERE "id"=$1 LIMIT 10`,
			checkValues:    true,
			expectedValues: []interface{}{"uuid"},
		},
		{
			name:           "SELECT relations->rel only, with conditions",
			query:          somesql.NewQuery().Select("").SelectRel("author", "tags").Where(somesql.And(somesql.LangEN, "id", "=", "uuid")),
			expectedSQL:    `SELECT json_build_object('author', "relations"->'author', 'tags', "relations"->'tags') "relations" FROM repo WHERE "id"=$1 LIMIT 10`,
			checkValues:    true,
			expectedValues: []interface{}{"uuid"},
		},
		// DELETE
		{
			name:        "DELETE * NO LIMIT",
			query:       somesql.NewQuery().Delete(),
			expectedSQL: `DELETE FROM repo`,
		},
		{
			name:        "DELETE * LIMIT",
			query:       somesql.NewQuery().Delete().SetLimit(20),
			expectedSQL: `DELETE FROM repo LIMIT 20`,
		},
		{
			name:        "DELETE * OFFSET 10",
			query:       somesql.NewQuery().Delete().SetOffset(10),
			expectedSQL: `DELETE FROM repo OFFSET 10`,
		},
		{
			name:        "DELETE * LIMIT 20 OFFSET 10",
			query:       somesql.NewQuery().Delete().SetLimit(20).SetOffset(10),
			expectedSQL: `DELETE FROM repo LIMIT 20 OFFSET 10`,
		},
		{
			name:           "DELETE with condition",
			query:          somesql.NewQuery().Delete().Where(somesql.And(somesql.LangEN, "id", "=", "uuid")),
			expectedSQL:    `DELETE FROM repo WHERE "id"=$1`,
			checkValues:    true,
			expectedValues: []interface{}{"uuid"},
		},
		{
			name:           "DELETE with conditions + relations",
			query:          somesql.NewQuery().Delete().Where(somesql.AndRel("", "article", "=", "uuid")),
			expectedSQL:    `DELETE FROM repo WHERE ("relations" @> '{"article":$1}'::JSONB)`,
			checkValues:    true,
			expectedValues: []interface{}{"uuid"},
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queryResult := tt.query.AsSQL()

			gotSQL := queryResult.GetSQL()
			assert.Equal(t, tt.expectedSQL, gotSQL, fmt.Sprintf("Fields %03d :: %s", i+1, tt.name))

			if tt.checkValues {
				gotValues := queryResult.GetValues()
				assert.Equal(t, tt.expectedValues, gotValues, fmt.Sprintf("Fields %03d :: invalid values :: %s", i+1, tt.name))
			}
		})
	}
}

func TestQuery_AsSQL_Insert(t *testing.T) {
	type testCase struct {
		name           string
		query          somesql.Query
		expectedSQL    string
		expectedValues []interface{}
	}

	tests := []testCase{
		// Insert
		{
			name:           "INSERT defaults",
			query:          somesql.NewQuery().Insert(somesql.NewFieldValue().UseDefaults().ID("1").CreatedAt(time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)).UpdatedAt(time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)).Status("published")),
			expectedSQL:    `INSERT INTO repo ("id", "created_at", "updated_at", "owner_id", "status", "type", "data_en") VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			expectedValues: []interface{}{"1", time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC), time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC), uuid.Nil.String(), "published", "", "{}"},
		},
		{
			name:           "INSERT no defaults",
			query:          somesql.NewQuery().Insert(somesql.NewFieldValue().ID("1").Status("published")),
			expectedSQL:    `INSERT INTO repo ("id", "status") VALUES ($1, $2)`,
			expectedValues: []interface{}{"1", "published"},
		},
		{
			name:           "INSERT no default + 1 relation",
			query:          somesql.NewQuery().Insert(somesql.NewFieldValue().Status("published")).InsertRel("tags", []string{"a", "b", "c"}),
			expectedSQL:    `INSERT INTO repo ("status", "relations") VALUES ($1, $2)`,
			expectedValues: []interface{}{"published", `{"tags":["a","b","c"]}`},
		},
		{
			name:           "INSERT no default + 2 relations",
			query:          somesql.NewQuery().Insert(somesql.NewFieldValue().Status("published")).InsertRel("author", []string{"x"}).InsertRel("tags", []string{"a", "b", "c"}),
			expectedSQL:    `INSERT INTO repo ("status", "relations") VALUES ($1, $2)`,
			expectedValues: []interface{}{"published", `{"author":["x"],"tags":["a","b","c"]}`},
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queryResult := tt.query.AsSQL()
			gotSQL, gotValues := queryResult.GetSQL(), queryResult.GetValues()

			assert.Equal(t, tt.expectedSQL, gotSQL, fmt.Sprintf("Fields %03d :: invalid sql :: %s", i+1, tt.name))
			assert.Equal(t, tt.expectedValues, gotValues, fmt.Sprintf("Fields %03d :: invalid values :: %s", i+1, tt.name))
		})
	}
}

func TestQuery_AsSQL_Update(t *testing.T) {
	type testCase struct {
		name           string
		query          somesql.Query
		expectedSQL    string
		expectedValues []interface{}
	}

	tests := []testCase{
		// Update
		{
			name:           "UPDATE meta fields",
			query:          somesql.NewQuery().Update(somesql.NewFieldValue().ID("1").Status("published")),
			expectedSQL:    `UPDATE repo SET "id" = $1, "status" = $2`,
			expectedValues: []interface{}{"1", "published"},
		},
		{
			name:           "UPDATE data fields",
			query:          somesql.NewQuery().Update(somesql.NewFieldValue().Set("body", "body value").Set("author_id", "123")),
			expectedSQL:    `UPDATE repo SET "data_en" = "data_en" || {"body": $1, "author_id": $2}`,
			expectedValues: []interface{}{"body value", "123"},
		},
		{
			name:           "UPDATE data fields (LangFR)",
			query:          somesql.NewQuery().Update(somesql.NewFieldValue().Set("body", "body value").Set("author_id", "123")).SetLang(somesql.LangFR),
			expectedSQL:    `UPDATE repo SET "data_fr" = "data_fr" || {"body": $1, "author_id": $2}`,
			expectedValues: []interface{}{"body value", "123"},
		},
		{
			name:           "UPDATE meta + data fields",
			query:          somesql.NewQuery().Update(somesql.NewFieldValue().ID("1").Status("published").Set("body", "body value").Set("author_id", "123")),
			expectedSQL:    `UPDATE repo SET "id" = $1, "status" = $2, "data_en" = "data_en" || {"body": $3, "author_id": $4}`,
			expectedValues: []interface{}{"1", "published", "body value", "123"},
		},
		{
			name:           "UPDATE meta + data fields (LangFR)",
			query:          somesql.NewQuery().Update(somesql.NewFieldValue().ID("1").Status("published").Set("body", "body value").Set("author_id", "123")).SetLang(somesql.LangFR),
			expectedSQL:    `UPDATE repo SET "id" = $1, "status" = $2, "data_fr" = "data_fr" || {"body": $3, "author_id": $4}`,
			expectedValues: []interface{}{"1", "published", "body value", "123"},
		},
		{
			name:           "UPDATE meta + data fields conditions",
			query:          somesql.NewQuery().Update(somesql.NewFieldValue().ID("1").Status("published").Set("body", "body value").Set("author_id", "123")).Where(somesql.And(somesql.LangEN, "id", "=", "234")),
			expectedSQL:    `UPDATE repo SET "id" = $1, "status" = $2, "data_en" = "data_en" || {"body": $3, "author_id": $4} WHERE "id"=$5`,
			expectedValues: []interface{}{"1", "published", "body value", "123", "234"},
		},
		{
			name:           "UPDATE meta + data fields conditions 2",
			query:          somesql.NewQuery().Update(somesql.NewFieldValue().ID("1").Status("published").Set("body", "body value").Set("author_id", "123")).Where(somesql.And(somesql.LangEN, "author_id", "=", "234")),
			expectedSQL:    `UPDATE repo SET "id" = $1, "status" = $2, "data_en" = "data_en" || {"body": $3, "author_id": $4} WHERE "data_en"->>'author_id'=$5`,
			expectedValues: []interface{}{"1", "published", "body value", "123", "234"},
		},
		{
			name:           "UPDATE meta + data json",
			query:          somesql.NewQuery().Update(somesql.NewFieldValue().ID("1").Data(`{"body": 'body value', "author_id": 1}`)),
			expectedSQL:    `UPDATE repo SET "id" = $1, "data_en" = $2`,
			expectedValues: []interface{}{"1", `{"body": 'body value', "author_id": 1}`},
		},
		// Update relations : add
		{
			name:           "UPDATE add 1 relation only",
			query:          somesql.NewQuery(nil).Update(somesql.NewFieldValue()).AddRel("tags", []string{"a", "b"}),
			expectedSQL:    `UPDATE repo SET "relations" = relAdd.relations FROM (SELECT (("relations" - 'tags') || JSONB_BUILD_OBJECT('tags', "relations"->'tags' || '$1'::JSONB)) "relations" FROM repo) relAdd`,
			expectedValues: []interface{}{`["a","b"]`},
		},
		{
			name:           "UPDATE add 1 or more relations only",
			query:          somesql.NewQuery(nil).Update(somesql.NewFieldValue()).AddRel("tags", []string{"a", "b"}).AddRel("author", []string{"x"}),
			expectedSQL:    `UPDATE repo SET "relations" = relAdd.relations FROM (SELECT (("relations" - 'tags') || ("relations" - 'author') || JSONB_BUILD_OBJECT('tags', "relations"->'tags' || '$1'::JSONB, 'author', "relations"->'author' || '$2'::JSONB)) "relations" FROM repo) relAdd`,
			expectedValues: []interface{}{`["a","b"]`, `["x"]`},
		},
		{
			name:           "UPDATE add 1 relation only + conditions",
			query:          somesql.NewQuery(nil).Update(somesql.NewFieldValue()).AddRel("author", []string{"x"}).Where(somesql.And(somesql.LangEN, "id", "=", "uuid")),
			expectedSQL:    `UPDATE repo SET "relations" = relAdd.relations FROM (SELECT (("relations" - 'author') || JSONB_BUILD_OBJECT('author', "relations"->'author' || '$1'::JSONB)) "relations" FROM repo WHERE "id"=$2) relAdd WHERE "id"=$3`,
			expectedValues: []interface{}{`["x"]`, "uuid", "uuid"},
		},
		// Update relations : remove
		{
			name:           "UPDATE remove 1 relation only",
			query:          somesql.NewQuery(nil).Update(somesql.NewFieldValue()).RemoveRel("tags", []string{"a", "b"}),
			expectedSQL:    `UPDATE repo SET "relations" = updates.updRel FROM (SELECT (("relations" - 'tags') || JSONB_BUILD_OBJECT('tags', JSONB_AGG(tagsUpd))) "updatedRel" FROM (SELECT "relations", JSONB_ARRAY_ELEMENTS_TEXT("relations"->'tags') tagsUpd FROM repo) expandedValues WHERE tagsUpd NOT IN ($1) GROUP BY "relations") updates`,
			expectedValues: []interface{}{`["a","b"]`},
		},
		{
			name:           "UPDATE remove 1 or more relations only",
			query:          somesql.NewQuery(nil).Update(somesql.NewFieldValue()).RemoveRel("tags", []string{"a", "b"}).RemoveRel("author", []string{"x"}),
			expectedSQL:    `UPDATE repo SET "relations" = updates.updRel FROM (SELECT (("relations" - 'tags') || ("relations" - 'author') || JSONB_BUILD_OBJECT('tags', JSONB_AGG(tagsUpd), 'author', JSONB_AGG(authorUpd))) "updatedRel" FROM (SELECT "relations", JSONB_ARRAY_ELEMENTS_TEXT("relations"->'tags') tagsUpd, JSONB_ARRAY_ELEMENTS_TEXT("relations"->'author') authorUpd FROM repo) expandedValues WHERE tagsUpd NOT IN ($1) AND authorUpd NOT IN ($2) GROUP BY "relations") updates`,
			expectedValues: []interface{}{`["a","b"]`, `["x"]`},
		},
		{
			name:           "UPDATE remove 1 relation only + conditions",
			query:          somesql.NewQuery(nil).Update(somesql.NewFieldValue()).RemoveRel("author", []string{"x"}).Where(somesql.And(somesql.LangEN, "id", "=", "uuid")),
			expectedSQL:    `UPDATE repo SET "relations" = updates.updRel FROM (SELECT (("relations" - 'author') || JSONB_BUILD_OBJECT('author', JSONB_AGG(authorUpd))) "updatedRel" FROM (SELECT "relations", JSONB_ARRAY_ELEMENTS_TEXT("relations"->'author') authorUpd FROM repo WHERE "id"=$1) expandedValues WHERE authorUpd NOT IN ($2) GROUP BY "relations") updates WHERE "id"=$3`,
			expectedValues: []interface{}{"uuid", `["x"]`, "uuid"},
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queryResult := tt.query.AsSQL()
			gotSQL, gotValues := queryResult.GetSQL(), queryResult.GetValues()

			assert.Equal(t, tt.expectedSQL, gotSQL, fmt.Sprintf("Fields %03d :: invalid sql :: %s", i+1, tt.name))
			assert.Equal(t, tt.expectedValues, gotValues, fmt.Sprintf("Fields %03d :: invalid values :: %s", i+1, tt.name))
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
			expectedSQL:    `SELECT "data_en" FROM repo WHERE "id"=$1 LIMIT 10`,
			expectedValues: []interface{}{"1"},
		},
		{
			name:           "WHERE id=? AND status=?",
			query:          somesql.NewQuery().Select("data").Where(somesql.And(somesql.LangEN, "id", "=", "1")).Where(somesql.And(somesql.LangEN, "status", "=", "published")),
			expectedSQL:    `SELECT "data_en" FROM repo WHERE "id"=$1 AND "status"=$2 LIMIT 10`,
			expectedValues: []interface{}{"1", "published"},
		},
		{
			name:           "WHERE id=? OR status=?",
			query:          somesql.NewQuery().Select("data").Where(somesql.And(somesql.LangEN, "id", "=", "1")).Where(somesql.Or(somesql.LangEN, "status", "=", "published")),
			expectedSQL:    `SELECT "data_en" FROM repo WHERE "id"=$1 OR "status"=$2 LIMIT 10`,
			expectedValues: []interface{}{"1", "published"},
		},
		{
			name:           "WHERE id=? AND status=? OR type=?",
			query:          somesql.NewQuery().Select("data").Where(somesql.And(somesql.LangEN, "id", "=", "1")).Where(somesql.And(somesql.LangEN, "status", "=", "published")).Where(somesql.Or(somesql.LangEN, "type", "=", "article")),
			expectedSQL:    `SELECT "data_en" FROM repo WHERE "id"=$1 AND "status"=$2 OR "type"=$3 LIMIT 10`,
			expectedValues: []interface{}{"1", "published", "article"},
		},
		// Regular condition clauses on json attributes
		{
			name:           `WHERE "data_en"->>'author_id'=?`,
			query:          somesql.NewQuery().Select("data").Where(somesql.And(somesql.LangEN, "author_id", "=", "1")),
			expectedSQL:    `SELECT "data_en" FROM repo WHERE "data_en"->>'author_id'=$1 LIMIT 10`,
			expectedValues: []interface{}{"1"},
		},
		{
			name:           `WHERE "data_fr"->>'author_id'=? (langFR)`,
			query:          somesql.NewQuery().Select("data").SetLang(somesql.LangFR).Where(somesql.And(somesql.LangFR, "author_id", "=", "1")),
			expectedSQL:    `SELECT "data_fr" FROM repo WHERE "data_fr"->>'author_id'=$1 LIMIT 10`,
			expectedValues: []interface{}{"1"},
		},
		{
			name:           `WHERE "data_fr"->>'author_id'=? OR "data_fr"->>'category_id'=? (langFR)`,
			query:          somesql.NewQuery().Select("data").SetLang(somesql.LangFR).Where(somesql.And(somesql.LangFR, "author_id", "=", "1")).Where(somesql.Or(somesql.LangFR, "category_id", "=", "2")),
			expectedSQL:    `SELECT "data_fr" FROM repo WHERE "data_fr"->>'author_id'=$1 OR "data_fr"->>'category_id'=$2 LIMIT 10`,
			expectedValues: []interface{}{"1", "2"},
		},
		{
			name:           `WHERE "data_fr"->>'author_id'=? AND "data_fr"->>'category_id'=? (langFR)`,
			query:          somesql.NewQuery().Select("data").SetLang(somesql.LangFR).Where(somesql.And(somesql.LangFR, "author_id", "=", "1")).Where(somesql.And(somesql.LangFR, "category_id", "=", "2")),
			expectedSQL:    `SELECT "data_fr" FROM repo WHERE "data_fr"->>'author_id'=$1 AND "data_fr"->>'category_id'=$2 LIMIT 10`,
			expectedValues: []interface{}{"1", "2"},
		},
		// DELETE
		{
			name:           "DELETE WHERE id=?",
			query:          somesql.NewQuery().Delete().Where(somesql.And(somesql.LangEN, "id", "=", "1")),
			expectedSQL:    `DELETE FROM repo WHERE "id"=$1`,
			expectedValues: []interface{}{"1"},
		},
		{
			name:           "DELETE WHERE id=? AND status=? OR type=?",
			query:          somesql.NewQuery().Delete().Where(somesql.And(somesql.LangEN, "id", "=", "1")).Where(somesql.And(somesql.LangEN, "status", "=", "published")).Where(somesql.Or(somesql.LangEN, "type", "=", "article")),
			expectedSQL:    `DELETE FROM repo WHERE "id"=$1 AND "status"=$2 OR "type"=$3`,
			expectedValues: []interface{}{"1", "published", "article"},
		},
		{
			name:           `DELETE WHERE "data_en"->>'author_id'=?`,
			query:          somesql.NewQuery().Delete().Where(somesql.And(somesql.LangEN, "author_id", "=", "1")),
			expectedSQL:    `DELETE FROM repo WHERE "data_en"->>'author_id'=$1`,
			expectedValues: []interface{}{"1"},
		},
		{
			name:           `DELETE WHERE "data_fr"->>'author_id'=? AND "data_fr"->>'category_id'=? (langFR)`,
			query:          somesql.NewQuery().Delete().SetLang(somesql.LangFR).Where(somesql.And(somesql.LangFR, "author_id", "=", "1")).Where(somesql.And(somesql.LangFR, "category_id", "=", "2")),
			expectedSQL:    `DELETE FROM repo WHERE "data_fr"->>'author_id'=$1 AND "data_fr"->>'category_id'=$2`,
			expectedValues: []interface{}{"1", "2"},
		},
		{
			name:           `DELETE WHERE "relations"->>'article'`,
			query:          somesql.NewQuery().Delete().Where(somesql.And("", "type", "=", "category")).Where(somesql.AndRel("", "article", "=", "uuid")),
			expectedSQL:    `DELETE FROM repo WHERE "type"=$1 AND ("relations" @> '{"article":$2}'::JSONB)`,
			expectedValues: []interface{}{"category", "uuid"},
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queryResult := tt.query.AsSQL()
			gotSQL, gotValues := queryResult.GetSQL(), queryResult.GetValues()

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
			expectedSQL:    `SELECT "data_en" FROM repo WHERE ("data_en"->>'badge'=$1 OR "data_en"->>'badge'=$2) LIMIT 10`,
			expectedValues: []interface{}{"video", "audio"},
		},
		{
			name:           "WHERE (... OR ...) [2]",
			query:          somesql.NewQuery().Select("data").Where(somesql.OrGroup(somesql.And(somesql.LangEN, "badge", "=", "video"), somesql.Or(somesql.LangEN, "badge", "=", "audio"))),
			expectedSQL:    `SELECT "data_en" FROM repo WHERE ("data_en"->>'badge'=$1 OR "data_en"->>'badge'=$2) LIMIT 10`,
			expectedValues: []interface{}{"video", "audio"},
		},
		{
			name:           "WHERE (... AND ...) [1]",
			query:          somesql.NewQuery().Select("data").Where(somesql.AndGroup(somesql.Or(somesql.LangEN, "badge", "=", "video"), somesql.And(somesql.LangEN, "has_video", "=", true))),
			expectedSQL:    `SELECT "data_en" FROM repo WHERE ("data_en"->>'badge'=$1 AND ("data_en"->>'has_video')::BOOLEAN=$2) LIMIT 10`,
			expectedValues: []interface{}{"video", true},
		},
		{
			name:           "WHERE (... AND ...) [2]",
			query:          somesql.NewQuery().Select("data").Where(somesql.OrGroup(somesql.Or(somesql.LangEN, "badge", "=", "video"), somesql.And(somesql.LangEN, "has_video", "=", true))),
			expectedSQL:    `SELECT "data_en" FROM repo WHERE ("data_en"->>'badge'=$1 AND ("data_en"->>'has_video')::BOOLEAN=$2) LIMIT 10`,
			expectedValues: []interface{}{"video", true},
		},
		{
			name:           "WHERE (... AND ...) AND (... OR ...) [1]",
			query:          somesql.NewQuery().Select("data").Where(somesql.OrGroup(somesql.Or(somesql.LangEN, "badge", "=", "video"), somesql.And(somesql.LangEN, "has_video", "=", true))).Where(somesql.AndGroup(somesql.Or(somesql.LangEN, "badge", "=", "video"), somesql.Or(somesql.LangEN, "has_video", "=", true))),
			expectedSQL:    `SELECT "data_en" FROM repo WHERE ("data_en"->>'badge'=$1 AND ("data_en"->>'has_video')::BOOLEAN=$2) AND ("data_en"->>'badge'=$3 OR ("data_en"->>'has_video')::BOOLEAN=$4) LIMIT 10`,
			expectedValues: []interface{}{"video", true, "video", true},
		},
		{
			name:           "WHERE (... AND ...) OR (... AND ...) [2]",
			query:          somesql.NewQuery().Select("data").Where(somesql.OrGroup(somesql.Or(somesql.LangEN, "badge", "=", "video"), somesql.And(somesql.LangEN, "has_video", "=", true))).Where(somesql.OrGroup(somesql.Or(somesql.LangEN, "badge", "=", "video"), somesql.And(somesql.LangEN, "has_video", "=", true))),
			expectedSQL:    `SELECT "data_en" FROM repo WHERE ("data_en"->>'badge'=$1 AND ("data_en"->>'has_video')::BOOLEAN=$2) OR ("data_en"->>'badge'=$3 AND ("data_en"->>'has_video')::BOOLEAN=$4) LIMIT 10`,
			expectedValues: []interface{}{"video", true, "video", true},
		},
		// DELETE
		{
			name:           "DELETE WHERE (... OR ...)",
			query:          somesql.NewQuery().Delete().Where(somesql.AndGroup(somesql.And(somesql.LangEN, "badge", "=", "video"), somesql.Or(somesql.LangEN, "badge", "=", "audio"))),
			expectedSQL:    `DELETE FROM repo WHERE ("data_en"->>'badge'=$1 OR "data_en"->>'badge'=$2)`,
			expectedValues: []interface{}{"video", "audio"},
		},
		{
			name:           "DELETE WHERE (... AND ...)",
			query:          somesql.NewQuery().Delete().Where(somesql.OrGroup(somesql.Or(somesql.LangEN, "badge", "=", "video"), somesql.And(somesql.LangEN, "has_video", "=", true))),
			expectedSQL:    `DELETE FROM repo WHERE ("data_en"->>'badge'=$1 AND ("data_en"->>'has_video')::BOOLEAN=$2)`,
			expectedValues: []interface{}{"video", true},
		},
		{
			name:           "DELETE WHERE (... AND ...) AND (... OR ...)",
			query:          somesql.NewQuery().Delete().Where(somesql.OrGroup(somesql.Or(somesql.LangEN, "badge", "=", "video"), somesql.And(somesql.LangEN, "has_video", "=", true))).Where(somesql.AndGroup(somesql.Or(somesql.LangEN, "badge", "=", "video"), somesql.Or(somesql.LangEN, "has_video", "=", true))),
			expectedSQL:    `DELETE FROM repo WHERE ("data_en"->>'badge'=$1 AND ("data_en"->>'has_video')::BOOLEAN=$2) AND ("data_en"->>'badge'=$3 OR ("data_en"->>'has_video')::BOOLEAN=$4)`,
			expectedValues: []interface{}{"video", true, "video", true},
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queryResult := tt.query.AsSQL()
			gotSQL, gotValues := queryResult.GetSQL(), queryResult.GetValues()

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
			expectedSQL:    `SELECT "data_en" FROM repo WHERE "id" IN ($1,$2,$3) LIMIT 10`,
			expectedValues: []interface{}{"A", "B", "C"},
		},
		{
			name:           "WHERE id IN (...) - primitive field - LangFR",
			query:          somesql.NewQuery().Select("data").Where(somesql.AndIn(somesql.LangEN, "id", []string{"A", "B", "C"})).SetLang(somesql.LangFR),
			expectedSQL:    `SELECT "data_fr" FROM repo WHERE "id" IN ($1,$2,$3) LIMIT 10`,
			expectedValues: []interface{}{"A", "B", "C"},
		},
		{
			name:           "WHERE field IN (...) - JSONB",
			query:          somesql.NewQuery().Where(somesql.AndIn(somesql.LangEN, "name", []string{"A", "B", "C"})),
			expectedSQL:    `SELECT "id", "created_at", "updated_at", "owner_id", "status", "type", "data_en" FROM repo WHERE ("data_en" @> '{"name":["$1"]}'::JSONB OR '{"name":["$2"]}'::JSONB OR '{"name":["$3"]}'::JSONB) LIMIT 10`,
			expectedValues: []interface{}{"A", "B", "C"},
		},
		{
			name:           "WHERE field IN (...) - JSONB - LangFR",
			query:          somesql.NewQuery().Where(somesql.AndIn(somesql.LangFR, "name", []string{"A", "B", "C"})).SetLang(somesql.LangFR),
			expectedSQL:    `SELECT "id", "created_at", "updated_at", "owner_id", "status", "type", "data_fr" FROM repo WHERE ("data_fr" @> '{"name":["$1"]}'::JSONB OR '{"name":["$2"]}'::JSONB OR '{"name":["$3"]}'::JSONB) LIMIT 10`,
			expectedValues: []interface{}{"A", "B", "C"},
		},
		{
			name:           "WHERE FUNC(field) IN (...) - primitive field",
			query:          somesql.NewQuery().Select("data").Where(somesql.AndIn(somesql.LangEN, "updated_at", []string{"2019"}, "YEAR")),
			expectedSQL:    `SELECT "data_en" FROM repo WHERE YEAR("updated_at") IN ($1) LIMIT 10`,
			expectedValues: []interface{}{"2019"},
		},
		{
			name:           "WHERE FUNC(field) IN (...) - JSONB",
			query:          somesql.NewQuery().Where(somesql.AndIn(somesql.LangEN, "name", []string{"a"}, "LOWER")),
			expectedSQL:    `SELECT "id", "created_at", "updated_at", "owner_id", "status", "type", "data_en" FROM repo WHERE ("data_en" @> '{"name":["$1"]}'::JSONB) LIMIT 10`,
			expectedValues: []interface{}{"a"},
		},

		{
			name:           "WHERE id NOT IN (...) - primitive field",
			query:          somesql.NewQuery().Select("data").Where(somesql.AndNotIn(somesql.LangEN, "id", []string{"A", "B", "C"})),
			expectedSQL:    `SELECT "data_en" FROM repo WHERE "id" NOT IN ($1,$2,$3) LIMIT 10`,
			expectedValues: []interface{}{"A", "B", "C"},
		},
		{
			name:           "WHERE id NOT IN (...) - primitive field - LangFR",
			query:          somesql.NewQuery().Select("data").Where(somesql.AndNotIn(somesql.LangEN, "id", []string{"A", "B", "C"})).SetLang(somesql.LangFR),
			expectedSQL:    `SELECT "data_fr" FROM repo WHERE "id" NOT IN ($1,$2,$3) LIMIT 10`,
			expectedValues: []interface{}{"A", "B", "C"},
		},
		{
			name:           "WHERE FUNC(field) NOT IN (...) - primitive field",
			query:          somesql.NewQuery().Select("data").Where(somesql.AndNotIn(somesql.LangEN, "updated_at", []string{"2019"}, "YEAR")),
			expectedSQL:    `SELECT "data_en" FROM repo WHERE YEAR("updated_at") NOT IN ($1) LIMIT 10`,
			expectedValues: []interface{}{"2019"},
		},
		{
			name:           "WHERE field NOT IN (...) - JSONB",
			query:          somesql.NewQuery().Where(somesql.AndNotIn(somesql.LangEN, "name", []string{"A", "B", "C"})),
			expectedSQL:    `SELECT "id", "created_at", "updated_at", "owner_id", "status", "type", "data_en" FROM repo WHERE NOT("data_en" @> '{"name":["$1"]}'::JSONB OR '{"name":["$2"]}'::JSONB OR '{"name":["$3"]}'::JSONB) LIMIT 10`,
			expectedValues: []interface{}{"A", "B", "C"},
		},
		{
			name:           "WHERE FUNC(field) NOT IN (...) - JSONB",
			query:          somesql.NewQuery().Where(somesql.AndNotIn(somesql.LangEN, "name", []string{"a"}, "LOWER")),
			expectedSQL:    `SELECT "id", "created_at", "updated_at", "owner_id", "status", "type", "data_en" FROM repo WHERE NOT("data_en" @> '{"name":["$1"]}'::JSONB) LIMIT 10`,
			expectedValues: []interface{}{"a"},
		},

		{
			name:           "WHERE id IN (...) AND NOT IN (...) - primitive field",
			query:          somesql.NewQuery().Select("data").Where(somesql.AndIn(somesql.LangEN, "id", []string{"A", "B"})).Where(somesql.AndNotIn(somesql.LangEN, "id", []string{"C", "D"})),
			expectedSQL:    `SELECT "data_en" FROM repo WHERE "id" IN ($1,$2) AND "id" NOT IN ($3,$4) LIMIT 10`,
			expectedValues: []interface{}{"A", "B", "C", "D"},
		},
		{
			name:           "WHERE id IN (...) AND NOT IN (...) - JSONB",
			query:          somesql.NewQuery().Select("data").Where(somesql.AndIn(somesql.LangEN, "name", []string{"A", "B"})).Where(somesql.AndNotIn(somesql.LangEN, "id", []string{"C", "D"})),
			expectedSQL:    `SELECT "data_en" FROM repo WHERE ("data_en" @> '{"name":["$1"]}'::JSONB OR '{"name":["$2"]}'::JSONB) AND "id" NOT IN ($3,$4) LIMIT 10`,
			expectedValues: []interface{}{"A", "B", "C", "D"},
		},
		{
			name:           "WHERE id IN (...) OR NOT IN (...) - primitive field",
			query:          somesql.NewQuery().Select("data").Where(somesql.AndIn(somesql.LangEN, "id", []string{"A", "B"})).Where(somesql.OrNotIn(somesql.LangEN, "id", []string{"C", "D"})),
			expectedSQL:    `SELECT "data_en" FROM repo WHERE "id" IN ($1,$2) OR "id" NOT IN ($3,$4) LIMIT 10`,
			expectedValues: []interface{}{"A", "B", "C", "D"},
		},
		{
			name:           "WHERE FUNC(field) IN (...) AND field NOT IN (...) - primitive field",
			query:          somesql.NewQuery().Select("data").Where(somesql.AndIn(somesql.LangEN, "updated_at", []string{"2019"}, "YEAR")).Where(somesql.AndNotIn(somesql.LangEN, "id", []string{"A", "B"})),
			expectedSQL:    `SELECT "data_en" FROM repo WHERE YEAR("updated_at") IN ($1) AND "id" NOT IN ($2,$3) LIMIT 10`,
			expectedValues: []interface{}{"2019", "A", "B"},
		},

		{
			name:           "WHERE id IN (...) AND field = ...",
			query:          somesql.NewQuery().Select("data").Where(somesql.AndIn(somesql.LangEN, "type", []string{"article", "dossier"})).Where(somesql.And(somesql.LangEN, "status", "=", []string{"published"})),
			expectedSQL:    `SELECT "data_en" FROM repo WHERE "type" IN ($1,$2) AND "status"=$3 LIMIT 10`,
			expectedValues: []interface{}{"article", "dossier", "published"},
		},
		{
			name:           "WHERE id IN (...) AND field = ... - JSONB",
			query:          somesql.NewQuery().Select("data").Where(somesql.AndIn(somesql.LangEN, "type", []string{"article", "dossier"})).Where(somesql.And(somesql.LangEN, "status", "=", []string{"published"})),
			expectedSQL:    `SELECT "data_en" FROM repo WHERE "type" IN ($1,$2) AND "status"=$3 LIMIT 10`,
			expectedValues: []interface{}{"article", "dossier", "published"},
		},
		{
			name:           "WHERE FUNC(field) IN (...) AND field NOT IN (...) - primitive field",
			query:          somesql.NewQuery().Select("data").Where(somesql.AndIn(somesql.LangEN, "updated_at", []string{"2019"}, "YEAR")).Where(somesql.AndNotIn(somesql.LangEN, "id", []string{"A", "B"})).Where(somesql.And(somesql.LangEN, "status", "=", "published")),
			expectedSQL:    `SELECT "data_en" FROM repo WHERE YEAR("updated_at") IN ($1) AND "id" NOT IN ($2,$3) AND "status"=$4 LIMIT 10`,
			expectedValues: []interface{}{"2019", "A", "B", "published"},
		},
		{
			name:           "WHERE FUNC(field) IN (...) AND field NOT IN (...) - JSONB",
			query:          somesql.NewQuery().Select("data").Where(somesql.AndIn(somesql.LangEN, "tag_ids", []string{"A"})).Where(somesql.AndIn(somesql.LangEN, "author_ids", []string{"B"})).Where(somesql.And(somesql.LangEN, "status", "=", "published")),
			expectedSQL:    `SELECT "data_en" FROM repo WHERE ("data_en" @> '{"tag_ids":["$1"]}'::JSONB) AND ("data_en" @> '{"author_ids":["$2"]}'::JSONB) AND "status"=$3 LIMIT 10`,
			expectedValues: []interface{}{"A", "B", "published"},
		},

		{
			name:           "WHERE id IN (...) AND field = ...",
			query:          somesql.NewQuery().Select("data").Where(somesql.AndGroup(somesql.And(somesql.LangEN, "badge", "=", "video"), somesql.And(somesql.LangEN, "has_video", "=", true))).Where(somesql.And(somesql.LangEN, "status", "=", []string{"published"})),
			expectedSQL:    `SELECT "data_en" FROM repo WHERE ("data_en"->>'badge'=$1 AND ("data_en"->>'has_video')::BOOLEAN=$2) AND "status"=$3 LIMIT 10`,
			expectedValues: []interface{}{"video", true, "published"},
		},
		{
			name:           "WHERE field = ... OR (... AND ...) [1]",
			query:          somesql.NewQuery().Select("data").Where(somesql.And(somesql.LangEN, "status", "=", []string{"published"})).Where(somesql.OrGroup(somesql.And(somesql.LangEN, "badge", "=", "video"), somesql.And(somesql.LangEN, "has_video", "=", true))),
			expectedSQL:    `SELECT "data_en" FROM repo WHERE "status"=$1 OR ("data_en"->>'badge'=$2 AND ("data_en"->>'has_video')::BOOLEAN=$3) LIMIT 10`,
			expectedValues: []interface{}{"published", "video", true},
		},
		{
			name:           "WHERE field = ... OR (... AND ...) [2]",
			query:          somesql.NewQuery().Select("data").Where(somesql.Or(somesql.LangEN, "status", "=", []string{"published"})).Where(somesql.OrGroup(somesql.Or(somesql.LangEN, "badge", "=", "video"), somesql.And(somesql.LangEN, "has_video", "=", true))),
			expectedSQL:    `SELECT "data_en" FROM repo WHERE "status"=$1 OR ("data_en"->>'badge'=$2 AND ("data_en"->>'has_video')::BOOLEAN=$3) LIMIT 10`,
			expectedValues: []interface{}{"published", "video", true},
		},
		{
			name:           "WHERE field = ... AND (... OR ...) [1]",
			query:          somesql.NewQuery().Select("data").Where(somesql.Or(somesql.LangEN, "status", "=", []string{"published"})).Where(somesql.AndGroup(somesql.And(somesql.LangEN, "badge", "=", "video"), somesql.Or(somesql.LangEN, "has_video", "=", true))),
			expectedSQL:    `SELECT "data_en" FROM repo WHERE "status"=$1 AND ("data_en"->>'badge'=$2 OR ("data_en"->>'has_video')::BOOLEAN=$3) LIMIT 10`,
			expectedValues: []interface{}{"published", "video", true},
		},
		{
			name:           "WHERE field = ... AND (... OR ...) [2]",
			query:          somesql.NewQuery().Select("data").Where(somesql.Or(somesql.LangEN, "status", "=", []string{"published"})).Where(somesql.AndGroup(somesql.And(somesql.LangEN, "badge", "=", "video"), somesql.Or(somesql.LangEN, "badge", "=", "audio"))),
			expectedSQL:    `SELECT "data_en" FROM repo WHERE "status"=$1 AND ("data_en"->>'badge'=$2 OR "data_en"->>'badge'=$3) LIMIT 10`,
			expectedValues: []interface{}{"published", "video", "audio"},
		},
		// DELETE
		{
			name:           "DELETE WHERE id IN (...) - primitive field",
			query:          somesql.NewQuery().Delete().Where(somesql.AndIn(somesql.LangEN, "id", []string{"A", "B", "C"})),
			expectedSQL:    `DELETE FROM repo WHERE "id" IN ($1,$2,$3)`,
			expectedValues: []interface{}{"A", "B", "C"},
		},
		{
			name:           "DELETE WHERE field IN (...) - JSONB - LangFR",
			query:          somesql.NewQuery().Delete().Where(somesql.AndIn(somesql.LangFR, "name", []string{"A", "B", "C"})).SetLang(somesql.LangFR),
			expectedSQL:    `DELETE FROM repo WHERE ("data_fr" @> '{"name":["$1"]}'::JSONB OR '{"name":["$2"]}'::JSONB OR '{"name":["$3"]}'::JSONB)`,
			expectedValues: []interface{}{"A", "B", "C"},
		},
		{
			name:           "DELETE WHERE FUNC(field) IN (...) AND field NOT IN (...) - primitive field",
			query:          somesql.NewQuery().Delete().Where(somesql.AndIn(somesql.LangEN, "updated_at", []string{"2019"}, "YEAR")).Where(somesql.AndNotIn(somesql.LangEN, "id", []string{"A", "B"})),
			expectedSQL:    `DELETE FROM repo WHERE YEAR("updated_at") IN ($1) AND "id" NOT IN ($2,$3)`,
			expectedValues: []interface{}{"2019", "A", "B"},
		},
		{
			name:           "DELETE WHERE FUNC(field) IN (...) AND field NOT IN (...) - JSONB",
			query:          somesql.NewQuery().Delete().Where(somesql.AndIn(somesql.LangEN, "tag_ids", []string{"A"})).Where(somesql.AndIn(somesql.LangEN, "author_ids", []string{"B"})).Where(somesql.And(somesql.LangEN, "status", "=", "published")),
			expectedSQL:    `DELETE FROM repo WHERE ("data_en" @> '{"tag_ids":["$1"]}'::JSONB) AND ("data_en" @> '{"author_ids":["$2"]}'::JSONB) AND "status"=$3`,
			expectedValues: []interface{}{"A", "B", "published"},
		},
		{
			name:           "DELETE WHERE field = ... OR (... AND ...)",
			query:          somesql.NewQuery().Delete().Where(somesql.Or(somesql.LangEN, "status", "=", []string{"published"})).Where(somesql.OrGroup(somesql.Or(somesql.LangEN, "badge", "=", "video"), somesql.And(somesql.LangEN, "has_video", "=", true))),
			expectedSQL:    `DELETE FROM repo WHERE "status"=$1 OR ("data_en"->>'badge'=$2 AND ("data_en"->>'has_video')::BOOLEAN=$3)`,
			expectedValues: []interface{}{"published", "video", true},
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queryResult := tt.query.AsSQL()
			gotSQL, gotValues := queryResult.GetSQL(), queryResult.GetValues()

			assert.Equal(t, tt.expectedSQL, gotSQL, fmt.Sprintf("Fields %03d :: invalid sql :: %s", i+1, tt.name))
			assert.Equal(t, tt.expectedValues, gotValues, fmt.Sprintf("Fields %03d :: invalid values :: %s", i+1, tt.name))
		})
	}
}

func TestQuery_AsSQL_InQuery(t *testing.T) {
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
		{
			name:           "AndInQuery 2 Fields",
			query:          somesql.NewQuery().Select("data").Where(somesql.AndInQuery(somesql.LangEN, "type", somesql.NewInnerQuery().Select("type", "slug").Where(somesql.And(somesql.LangEN, "id", "=", "002fd6b1-f715-4875-838b-1546f27327df")))),
			expectedSQL:    `SELECT "data_en" FROM repo WHERE "type" IN (SELECT "type", "data_en"->>'slug' "slug" FROM repo WHERE "id"=$1 LIMIT 10) LIMIT 10`,
			expectedValues: []interface{}{"002fd6b1-f715-4875-838b-1546f27327df"},
		},
		{
			name:           "AndInQuery",
			query:          somesql.NewQuery().Select("data").Where(somesql.AndInQuery(somesql.LangEN, "type", somesql.NewInnerQuery().Select("type").Where(somesql.And(somesql.LangEN, "id", "=", "002fd6b1-f715-4875-838b-1546f27327df")))),
			expectedSQL:    `SELECT "data_en" FROM repo WHERE "type" IN (SELECT "type" FROM repo WHERE "id"=$1 LIMIT 10) LIMIT 10`,
			expectedValues: []interface{}{"002fd6b1-f715-4875-838b-1546f27327df"},
		},
		{
			name:           "AndInQuery INNER NO LIMIT",
			query:          somesql.NewQuery().Select("data").Where(somesql.AndInQuery(somesql.LangEN, "type", somesql.NewInnerQuery().Select("type", "slug", "brief").Where(somesql.And(somesql.LangEN, "id", "=", "002fd6b1-f715-4875-838b-1546f27327df")).SetLimit(0))),
			expectedSQL:    `SELECT "data_en" FROM repo WHERE "type" IN (SELECT "type", "data_en"->>'slug' "slug", "data_en"->>'brief' "brief" FROM repo WHERE "id"=$1) LIMIT 10`,
			expectedValues: []interface{}{"002fd6b1-f715-4875-838b-1546f27327df"},
		},
		{
			name:           "AndInQuery INNER LIMIT 20",
			query:          somesql.NewQuery().Select("data").Where(somesql.AndInQuery(somesql.LangEN, "type", somesql.NewInnerQuery().Select("type", "slug").Where(somesql.And(somesql.LangEN, "id", "=", "002fd6b1-f715-4875-838b-1546f27327df")).SetLimit(20))),
			expectedSQL:    `SELECT "data_en" FROM repo WHERE "type" IN (SELECT "type", "data_en"->>'slug' "slug" FROM repo WHERE "id"=$1 LIMIT 20) LIMIT 10`,
			expectedValues: []interface{}{"002fd6b1-f715-4875-838b-1546f27327df"},
		},
		{
			name:           "AndInQuery INNER OFFSET 20",
			query:          somesql.NewQuery().Select("data").Where(somesql.AndInQuery(somesql.LangEN, "type", somesql.NewInnerQuery().Select("type", "slug").Where(somesql.And(somesql.LangEN, "id", "=", "002fd6b1-f715-4875-838b-1546f27327df")).SetOffset(20))),
			expectedSQL:    `SELECT "data_en" FROM repo WHERE "type" IN (SELECT "type", "data_en"->>'slug' "slug" FROM repo WHERE "id"=$1 LIMIT 10 OFFSET 20) LIMIT 10`,
			expectedValues: []interface{}{"002fd6b1-f715-4875-838b-1546f27327df"},
		},
		{
			name:           "AndInQuery",
			query:          somesql.NewQuery().Select("data").Where(somesql.AndInQuery(somesql.LangEN, "author_id", somesql.NewInnerQuery().Select("author_id").Where(somesql.And(somesql.LangEN, "id", "=", "002fd6b1-f715-4875-838b-1546f27327df")))),
			expectedSQL:    `SELECT "data_en" FROM repo WHERE "data_en"->>'author_id' IN (SELECT "data_en"->>'author_id' "author_id" FROM repo WHERE "id"=$1 LIMIT 10) LIMIT 10`,
			expectedValues: []interface{}{"002fd6b1-f715-4875-838b-1546f27327df"},
		},
		{
			name:           "AndInQuery",
			query:          somesql.NewQuery().Select("data").Where(somesql.AndIn(somesql.LangEN, "id", []string{"A", "B", "C"})).Where(somesql.AndInQuery(somesql.LangEN, "author_id", somesql.NewInnerQuery().Select("author_id").Where(somesql.And(somesql.LangEN, "id", "=", "002fd6b1-f715-4875-838b-1546f27327df")))),
			expectedSQL:    `SELECT "data_en" FROM repo WHERE "id" IN ($1,$2,$3) AND "data_en"->>'author_id' IN (SELECT "data_en"->>'author_id' "author_id" FROM repo WHERE "id"=$4 LIMIT 10) LIMIT 10`,
			expectedValues: []interface{}{"A", "B", "C", "002fd6b1-f715-4875-838b-1546f27327df"},
		},
		{
			name:           "AndNotInQuery",
			query:          somesql.NewQuery().Select("data").Where(somesql.AndIn(somesql.LangEN, "id", []string{"A", "B", "C"})).Where(somesql.AndNotInQuery(somesql.LangEN, "author_id", somesql.NewInnerQuery().Select("author_id").Where(somesql.And(somesql.LangEN, "id", "=", "002fd6b1-f715-4875-838b-1546f27327df")))),
			expectedSQL:    `SELECT "data_en" FROM repo WHERE "id" IN ($1,$2,$3) AND "data_en"->>'author_id' NOT IN (SELECT "data_en"->>'author_id' "author_id" FROM repo WHERE "id"=$4 LIMIT 10) LIMIT 10`,
			expectedValues: []interface{}{"A", "B", "C", "002fd6b1-f715-4875-838b-1546f27327df"},
		},

		// OrInQuery -> somesql.NewQuery().Select("type", "slug") : cannot have more than 1 field in subquery (throw an error?)
		{
			name:           "OrInQuery",
			query:          somesql.NewQuery().Select("data").Where(somesql.AndInQuery(somesql.LangEN, "type", somesql.NewInnerQuery().Select("type", "slug").Where(somesql.And(somesql.LangEN, "id", "=", "002fd6b1-f715-4875-838b-1546f27327df")))),
			expectedSQL:    `SELECT "data_en" FROM repo WHERE "type" IN (SELECT "type", "data_en"->>'slug' "slug" FROM repo WHERE "id"=$1 LIMIT 10) LIMIT 10`,
			expectedValues: []interface{}{"002fd6b1-f715-4875-838b-1546f27327df"},
		},
		{
			name:           "OrInQuery",
			query:          somesql.NewQuery().Select("data").Where(somesql.OrInQuery(somesql.LangEN, "type", somesql.NewInnerQuery().Select("type").Where(somesql.And(somesql.LangEN, "id", "=", "002fd6b1-f715-4875-838b-1546f27327df")))),
			expectedSQL:    `SELECT "data_en" FROM repo WHERE "type" IN (SELECT "type" FROM repo WHERE "id"=$1 LIMIT 10) LIMIT 10`,
			expectedValues: []interface{}{"002fd6b1-f715-4875-838b-1546f27327df"},
		},
		{
			name:           "OrInQuery",
			query:          somesql.NewQuery().Select("data").Where(somesql.OrInQuery(somesql.LangEN, "author_id", somesql.NewInnerQuery().Select("author_id").Where(somesql.And(somesql.LangEN, "id", "=", "002fd6b1-f715-4875-838b-1546f27327df")))),
			expectedSQL:    `SELECT "data_en" FROM repo WHERE "data_en"->>'author_id' IN (SELECT "data_en"->>'author_id' "author_id" FROM repo WHERE "id"=$1 LIMIT 10) LIMIT 10`,
			expectedValues: []interface{}{"002fd6b1-f715-4875-838b-1546f27327df"},
		},
		{
			name:           "OrInQuery",
			query:          somesql.NewQuery().Select("data").Where(somesql.AndIn(somesql.LangEN, "id", []string{"A", "B", "C"})).Where(somesql.OrInQuery(somesql.LangEN, "author_id", somesql.NewInnerQuery().Select("author_id").Where(somesql.And(somesql.LangEN, "id", "=", "002fd6b1-f715-4875-838b-1546f27327df")))),
			expectedSQL:    `SELECT "data_en" FROM repo WHERE "id" IN ($1,$2,$3) OR "data_en"->>'author_id' IN (SELECT "data_en"->>'author_id' "author_id" FROM repo WHERE "id"=$4 LIMIT 10) LIMIT 10`,
			expectedValues: []interface{}{"A", "B", "C", "002fd6b1-f715-4875-838b-1546f27327df"},
		},
		{
			name:           "OrNotInQuery",
			query:          somesql.NewQuery().Select("data").Where(somesql.AndIn(somesql.LangEN, "id", []string{"A", "B", "C"})).Where(somesql.OrNotInQuery(somesql.LangEN, "author_id", somesql.NewInnerQuery().Select("author_id").Where(somesql.And(somesql.LangEN, "id", "=", "002fd6b1-f715-4875-838b-1546f27327df")))),
			expectedSQL:    `SELECT "data_en" FROM repo WHERE "id" IN ($1,$2,$3) OR "data_en"->>'author_id' NOT IN (SELECT "data_en"->>'author_id' "author_id" FROM repo WHERE "id"=$4 LIMIT 10) LIMIT 10`,
			expectedValues: []interface{}{"A", "B", "C", "002fd6b1-f715-4875-838b-1546f27327df"},
		},
		// DELETE
		{
			name:           "DELETE AndInQuery 2 Fields",
			query:          somesql.NewQuery().Delete().Where(somesql.AndInQuery(somesql.LangEN, "type", somesql.NewInnerQuery().Select("type", "slug").Where(somesql.And(somesql.LangEN, "id", "=", "002fd6b1-f715-4875-838b-1546f27327df")))),
			expectedSQL:    `DELETE FROM repo WHERE "type" IN (SELECT "type", "data_en"->>'slug' "slug" FROM repo WHERE "id"=$1 LIMIT 10)`,
			expectedValues: []interface{}{"002fd6b1-f715-4875-838b-1546f27327df"},
		},
		{
			name:           "DELETE OrNotInQuery",
			query:          somesql.NewQuery().Delete().Where(somesql.AndIn(somesql.LangEN, "id", []string{"A", "B", "C"})).Where(somesql.OrNotInQuery(somesql.LangEN, "author_id", somesql.NewInnerQuery().Select("author_id").Where(somesql.And(somesql.LangEN, "id", "=", "002fd6b1-f715-4875-838b-1546f27327df")))),
			expectedSQL:    `DELETE FROM repo WHERE "id" IN ($1,$2,$3) OR "data_en"->>'author_id' NOT IN (SELECT "data_en"->>'author_id' "author_id" FROM repo WHERE "id"=$4 LIMIT 10)`,
			expectedValues: []interface{}{"A", "B", "C", "002fd6b1-f715-4875-838b-1546f27327df"},
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queryResult := tt.query.AsSQL()
			gotSQL, gotValues := queryResult.GetSQL(), queryResult.GetValues()

			assert.Equal(t, tt.expectedSQL, gotSQL, fmt.Sprintf("Fields %03d :: invalid sql :: %s", i+1, tt.name))
			assert.Equal(t, tt.expectedValues, gotValues, fmt.Sprintf("Fields %03d :: invalid values :: %s", i+1, tt.name))
		})
	}
}

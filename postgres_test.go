package somesql_test

import (
	"fmt"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"go.lsl.digital/gocipe/somesql"
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
			query:       somesql.NewSelect(somesql.LangEN, false),
			expectedSQL: `SELECT "id", "created_at", "updated_at", "owner_id", "status", "type", "data_en", "relations" FROM repo LIMIT 10`,
		},
		{
			name:        "SELECT * NO LIMIT",
			query:       somesql.NewSelect(somesql.LangEN, false).Limit(0),
			expectedSQL: `SELECT "id", "created_at", "updated_at", "owner_id", "status", "type", "data_en", "relations" FROM repo`,
		},
		{
			name:        "SELECT * LIMIT 30",
			query:       somesql.NewSelect(somesql.LangEN, false).Limit(30),
			expectedSQL: `SELECT "id", "created_at", "updated_at", "owner_id", "status", "type", "data_en", "relations" FROM repo LIMIT 30`,
		},
		{
			name:        "SELECT * OFFSET 10",
			query:       somesql.NewSelect(somesql.LangEN, false).Offset(10),
			expectedSQL: `SELECT "id", "created_at", "updated_at", "owner_id", "status", "type", "data_en", "relations" FROM repo LIMIT 10 OFFSET 10`,
		},
		{
			name:        "SELECT * LIMIT 30 OFFSET 20",
			query:       somesql.NewSelect(somesql.LangEN, false).Limit(30).Offset(20),
			expectedSQL: `SELECT "id", "created_at", "updated_at", "owner_id", "status", "type", "data_en", "relations" FROM repo LIMIT 30 OFFSET 20`,
		},
		{
			name:        "SELECT * (langEN)",
			query:       somesql.NewSelect(somesql.LangEN, false),
			expectedSQL: `SELECT "id", "created_at", "updated_at", "owner_id", "status", "type", "data_en", "relations" FROM repo LIMIT 10`,
		},
		{
			name:        "SELECT * (langFR)",
			query:       somesql.NewSelect(somesql.LangFR, false),
			expectedSQL: `SELECT "id", "created_at", "updated_at", "owner_id", "status", "type", "data_fr", "relations" FROM repo LIMIT 10`,
		},
		// Select some pre-defined fields
		{
			name:        "SELECT id, type, data",
			query:       somesql.NewSelect(somesql.LangEN, false).Fields("id", "type", "data"),
			expectedSQL: `SELECT "id", "type", "data_en" FROM repo LIMIT 10`,
		},
		{
			name:        "SELECT (EMPTY)",
			query:       somesql.NewSelect(somesql.LangEN, false).Fields(),
			expectedSQL: `SELECT "id", "created_at", "updated_at", "owner_id", "status", "type", "data_en", "relations" FROM repo LIMIT 10`,
		},
		{
			name:        "SELECT id, type, data (LangEN)",
			query:       somesql.NewSelect(somesql.LangEN, false).Fields("id", "type", "data"),
			expectedSQL: `SELECT "id", "type", "data_en" FROM repo LIMIT 10`,
		},
		{
			name:        "SELECT id, type, data (LangFR)",
			query:       somesql.NewSelect(somesql.LangFR, false).Fields("id", "type", "data"),
			expectedSQL: `SELECT "id", "type", "data_fr" FROM repo LIMIT 10`,
		},
		// Select pre-defined fields and json attributes ('data_en'/'data_fr') from data_*
		{
			name:        "SELECT id, type, data_en",
			query:       somesql.NewSelect(somesql.LangEN, false).Fields("id", "type", "data.data_en"),
			expectedSQL: `SELECT "id", "type", json_build_object('data_en', "data_en"->'data_en') "data" FROM repo LIMIT 10`,
		},
		{
			name:        "SELECT id, type, data_en (LangEN)",
			query:       somesql.NewSelect(somesql.LangEN, false).Fields("id", "type", "data.data_en"),
			expectedSQL: `SELECT "id", "type", json_build_object('data_en', "data_en"->'data_en') "data" FROM repo LIMIT 10`,
		},
		{
			name:        "SELECT id, type, data_en (LangFR)",
			query:       somesql.NewSelect(somesql.LangFR, false).Fields("id", "type", "data.data_en"),
			expectedSQL: `SELECT "id", "type", json_build_object('data_en', "data_fr"->'data_en') "data" FROM repo LIMIT 10`,
		},
		{
			name:        "SELECT id, type, data_fr",
			query:       somesql.NewSelect(somesql.LangEN, false).Fields("id", "type", "data.data_fr"),
			expectedSQL: `SELECT "id", "type", json_build_object('data_fr', "data_en"->'data_fr') "data" FROM repo LIMIT 10`,
		},
		{
			name:        "SELECT id, type, data_fr (LangEN)",
			query:       somesql.NewSelect(somesql.LangEN, false).Fields("id", "type", "data.data_fr"),
			expectedSQL: `SELECT "id", "type", json_build_object('data_fr', "data_en"->'data_fr') "data" FROM repo LIMIT 10`,
		},
		{
			name:        "SELECT id, type, data_fr (LangFR)",
			query:       somesql.NewSelect(somesql.LangFR, false).Fields("id", "type", "data.data_fr"),
			expectedSQL: `SELECT "id", "type", json_build_object('data_fr', "data_fr"->'data_fr') "data" FROM repo LIMIT 10`,
		},
		// Select pre-defined fields and json attributes (any other) from data_*
		{
			name:        "SELECT id, type, data_en->'body'",
			query:       somesql.NewSelect(somesql.LangEN, false).Fields("id", "type", "data.body"),
			expectedSQL: `SELECT "id", "type", json_build_object('body', "data_en"->'body') "data" FROM repo LIMIT 10`,
		},
		{
			name:        "SELECT id, type, data_en->'body' (LangEN)",
			query:       somesql.NewSelect(somesql.LangEN, false).Fields("id", "type", "data.body"),
			expectedSQL: `SELECT "id", "type", json_build_object('body', "data_en"->'body') "data" FROM repo LIMIT 10`,
		},
		{
			name:        "SELECT id, type, data_fr->'body' (LangFR)",
			query:       somesql.NewSelect(somesql.LangFR, false).Fields("id", "type", "data.body"),
			expectedSQL: `SELECT "id", "type", json_build_object('body', "data_fr"->'body') "data" FROM repo LIMIT 10`,
		},
		// Select pre-defined fields and json attributes (any other + compound) from data_*
		{
			name:        "SELECT id, type, data_en->'body', data_en->'author_id'",
			query:       somesql.NewSelect(somesql.LangEN, false).Fields("id", "type", "data.body", "data.author_id"),
			expectedSQL: `SELECT "id", "type", json_build_object('body', "data_en"->'body', 'author_id', "data_en"->'author_id') "data" FROM repo LIMIT 10`,
		},
		{
			name:        "SELECT id, type, data_en->'body', data_en->'author_id' (LangEN)",
			query:       somesql.NewSelect(somesql.LangEN, false).Fields("id", "type", "data.body", "data.author_id"),
			expectedSQL: `SELECT "id", "type", json_build_object('body', "data_en"->'body', 'author_id', "data_en"->'author_id') "data" FROM repo LIMIT 10`,
		},
		{
			name:        "SELECT id, type, data_fr->'body', data_fr->'author_id' (LangFR)",
			query:       somesql.NewSelect(somesql.LangFR, false).Fields("id", "type", "data.body", "data.author_id"),
			expectedSQL: `SELECT "id", "type", json_build_object('body', "data_fr"->'body', 'author_id', "data_fr"->'author_id') "data" FROM repo LIMIT 10`,
		},
		// SELECT with conditions
		{
			name:           "SELECT * with condition",
			query:          somesql.NewSelect(somesql.LangEN, false).Where(somesql.And(somesql.LangEN, "id", "=", "uuid")),
			expectedSQL:    `SELECT "id", "created_at", "updated_at", "owner_id", "status", "type", "data_en", "relations" FROM repo WHERE "id"=$1 LIMIT 10`,
			checkValues:    true,
			expectedValues: []interface{}{"uuid"},
		},
		// SELECT relations
		{
			name:        "SELECT id, relations->rel",
			query:       somesql.NewSelect(somesql.LangEN, false).Fields("id", "relations.author", "relations.tags"),
			expectedSQL: `SELECT "id", json_build_object('author', "relations"->'author', 'tags', "relations"->'tags') "relations" FROM repo LIMIT 10`,
		},
		{
			name:           "SELECT * + relations->rel, with conditions",
			query:          somesql.NewSelect(somesql.LangEN, false).Fields("relations.author", "relations.tags").Where(somesql.And(somesql.LangEN, "id", "=", "uuid")),
			expectedSQL:    `SELECT json_build_object('author', "relations"->'author', 'tags', "relations"->'tags') "relations" FROM repo WHERE "id"=$1 LIMIT 10`,
			checkValues:    true,
			expectedValues: []interface{}{"uuid"},
		},
		{
			name:           "SELECT relations->rel only, with conditions",
			query:          somesql.NewSelect(somesql.LangEN, false).Fields("relations.author", "relations.tags").Where(somesql.And(somesql.LangEN, "id", "=", "uuid")),
			expectedSQL:    `SELECT json_build_object('author', "relations"->'author', 'tags', "relations"->'tags') "relations" FROM repo WHERE "id"=$1 LIMIT 10`,
			checkValues:    true,
			expectedValues: []interface{}{"uuid"},
		},
		// DELETE
		// {
		// 	name:        "DELETE * NO LIMIT",
		// 	query:       somesql.NewQuery().Delete(),
		// 	expectedSQL: `DELETE FROM repo`,
		// },
		// {
		// 	name:        "DELETE * LIMIT",
		// 	query:       somesql.NewQuery().Delete().SetLimit(20),
		// 	expectedSQL: `DELETE FROM repo LIMIT 20`,
		// },
		// {
		// 	name:        "DELETE * OFFSET 10",
		// 	query:       somesql.NewQuery().Delete().SetOffset(10),
		// 	expectedSQL: `DELETE FROM repo OFFSET 10`,
		// },
		// {
		// 	name:        "DELETE * LIMIT 20 OFFSET 10",
		// 	query:       somesql.NewQuery().Delete().SetLimit(20).SetOffset(10),
		// 	expectedSQL: `DELETE FROM repo LIMIT 20 OFFSET 10`,
		// },
		// {
		// 	name:           "DELETE with condition",
		// 	query:          somesql.NewQuery().Delete().Where(somesql.And(somesql.LangEN, "id", "=", "uuid")),
		// 	expectedSQL:    `DELETE FROM repo WHERE "id"=$1`,
		// 	checkValues:    true,
		// 	expectedValues: []interface{}{"uuid"},
		// },
		// {
		// 	name:           "DELETE with conditions + relations",
		// 	query:          somesql.NewQuery().Delete().Where(somesql.AndRel("", "article", "=", "uuid")),
		// 	expectedSQL:    `DELETE FROM repo WHERE ("relations" @> '{"article":$1}'::JSONB)`,
		// 	checkValues:    true,
		// 	expectedValues: []interface{}{"uuid"},
		// },
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
			query:          somesql.NewInsert(somesql.LangEN).Fields(somesql.NewFields().UseDefaults().ID("1").CreatedAt(time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)).UpdatedAt(time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)).Status("published")),
			expectedSQL:    `INSERT INTO repo ("id", "created_at", "updated_at", "owner_id", "status", "type", "data_en", "relations") VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
			expectedValues: []interface{}{"1", time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC), time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC), uuid.Nil.String(), "published", "", "{}", "{}"},
		},
		{
			name:           "INSERT no defaults",
			query:          somesql.NewInsert(somesql.LangEN).Fields(somesql.NewFields().ID("1").Status("published")),
			expectedSQL:    `INSERT INTO repo ("id", "status") VALUES ($1, $2)`,
			expectedValues: []interface{}{"1", "published"},
		},
		{
			name:           "INSERT no default + 1 relation",
			query:          somesql.NewInsert(somesql.LangEN).Fields(somesql.NewFields().Status("published").Set("relations.tags", []string{"a", "b", "c"})),
			expectedSQL:    `INSERT INTO repo ("status", "relations") VALUES ($1, $2)`,
			expectedValues: []interface{}{"published", `{"tags":["a","b","c"]}`},
		},
		{
			name:           "INSERT no default + 2 relations",
			query:          somesql.NewInsert(somesql.LangEN).Fields(somesql.NewFields().Status("published").Set("relations.author", []string{"x"}).Set("relations.tags", []string{"a", "b", "c"})),
			expectedSQL:    `INSERT INTO repo ("status", "relations") VALUES ($1, $2)`,
			expectedValues: []interface{}{"published", `{"author":["x"],"tags":["a","b","c"]}`},
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

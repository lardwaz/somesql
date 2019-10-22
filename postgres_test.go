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
			query:       somesql.NewSelect("en"),
			expectedSQL: `SELECT "id", "created_at", "updated_at", "owner_id", "type", "data_en" FROM repo LIMIT 10`,
		},
		{
			name:        "SELECT * NO LIMIT",
			query:       somesql.NewSelect("en").Limit(0),
			expectedSQL: `SELECT "id", "created_at", "updated_at", "owner_id", "type", "data_en" FROM repo`,
		},
		{
			name:        "SELECT * LIMIT 30",
			query:       somesql.NewSelect("en").Limit(30),
			expectedSQL: `SELECT "id", "created_at", "updated_at", "owner_id", "type", "data_en" FROM repo LIMIT 30`,
		},
		{
			name:        "SELECT * OFFSET 10",
			query:       somesql.NewSelect("en").Offset(10),
			expectedSQL: `SELECT "id", "created_at", "updated_at", "owner_id", "type", "data_en" FROM repo LIMIT 10 OFFSET 10`,
		},
		{
			name:        "SELECT * LIMIT 30 OFFSET 20",
			query:       somesql.NewSelect("en").Limit(30).Offset(20),
			expectedSQL: `SELECT "id", "created_at", "updated_at", "owner_id", "type", "data_en" FROM repo LIMIT 30 OFFSET 20`,
		},
		{
			name:        "SELECT * (langEN)",
			query:       somesql.NewSelect("en"),
			expectedSQL: `SELECT "id", "created_at", "updated_at", "owner_id", "type", "data_en" FROM repo LIMIT 10`,
		},
		{
			name:        "SELECT * (langFR)",
			query:       somesql.NewSelect("fr"),
			expectedSQL: `SELECT "id", "created_at", "updated_at", "owner_id", "type", "data_fr" FROM repo LIMIT 10`,
		},
		// Select some pre-defined fields
		{
			name:        "SELECT id, type, data",
			query:       somesql.NewSelect("en").Fields("id", "type", "data"),
			expectedSQL: `SELECT "id", "type", "data_en" FROM repo LIMIT 10`,
		},
		{
			name:        "SELECT (EMPTY)",
			query:       somesql.NewSelect("en").Fields(),
			expectedSQL: `SELECT "id", "created_at", "updated_at", "owner_id", "type", "data_en" FROM repo LIMIT 10`,
		},
		{
			name:        "SELECT id, type, data (LangEN)",
			query:       somesql.NewSelect("en").Fields("id", "type", "data"),
			expectedSQL: `SELECT "id", "type", "data_en" FROM repo LIMIT 10`,
		},
		{
			name:        "SELECT id, type, data (LangFR)",
			query:       somesql.NewSelect("fr").Fields("id", "type", "data"),
			expectedSQL: `SELECT "id", "type", "data_fr" FROM repo LIMIT 10`,
		},
		// Order by
		{
			name:        "SELECT id, type, data ORDER BY id ASC",
			query:       somesql.NewSelect("en").Fields("id", "type", "data").Order("id", true),
			expectedSQL: `SELECT "id", "type", "data_en" FROM repo ORDER BY id ASC LIMIT 10`,
		},
		{
			name:        "SELECT id, type, data ORDER BY id ASC, type DESC",
			query:       somesql.NewSelect("en").Fields("id", "type", "data").Order("id", true).Order("type", false),
			expectedSQL: `SELECT "id", "type", "data_en" FROM repo ORDER BY id ASC, type DESC LIMIT 10`,
		},
		{
			name:        "SELECT id, type, data ORDER BY id ASC, type DESC, data_en ASC",
			query:       somesql.NewSelect("en").Fields("id", "type", "data").Order("id", true).Order("type", false).Order("data", true),
			expectedSQL: `SELECT "id", "type", "data_en" FROM repo ORDER BY id ASC, type DESC, data_en ASC LIMIT 10`,
		},
		{
			name:        "SELECT id, type, data ORDER BY id ASC, type DESC, data_en ASC, data_en.age DESC",
			query:       somesql.NewSelect("en").Fields("id", "type", "data").Order("id", true).Order("type", false).Order("data", true).Order("age", false),
			expectedSQL: `SELECT "id", "type", "data_en" FROM repo ORDER BY id ASC, type DESC, data_en ASC, "data_en"->>'age' DESC LIMIT 10`,
		},
		{
			name:        "SELECT id, type, data ORDER BY id ASC, type DESC, data_en ASC, data_en.age DESC, data_en.time ASC",
			query:       somesql.NewSelect("en").Fields("id", "type", "data").Order("id", true).Order("type", false).Order("data", true).Order("age", false).Order("time", true),
			expectedSQL: `SELECT "id", "type", "data_en" FROM repo ORDER BY id ASC, type DESC, data_en ASC, "data_en"->>'age' DESC, "data_en"->>'time' ASC LIMIT 10`,
		},
		// Select pre-defined fields and json attributes ('data_en'/'data_fr') from data_*
		{
			name:        "SELECT id, type, data_en",
			query:       somesql.NewSelect("en").Fields("id", "type", "data.data_en"),
			expectedSQL: `SELECT "id", "type", json_build_object('data_en', "data_en"->'data_en') "data" FROM repo LIMIT 10`,
		},
		{
			name:        "SELECT id, type, data_en (LangEN)",
			query:       somesql.NewSelect("en").Fields("id", "type", "data.data_en"),
			expectedSQL: `SELECT "id", "type", json_build_object('data_en', "data_en"->'data_en') "data" FROM repo LIMIT 10`,
		},
		{
			name:        "SELECT id, type, data_en (LangFR)",
			query:       somesql.NewSelect("fr").Fields("id", "type", "data.data_en"),
			expectedSQL: `SELECT "id", "type", json_build_object('data_en', "data_fr"->'data_en') "data" FROM repo LIMIT 10`,
		},
		{
			name:        "SELECT id, type, data_fr",
			query:       somesql.NewSelect("en").Fields("id", "type", "data.data_fr"),
			expectedSQL: `SELECT "id", "type", json_build_object('data_fr', "data_en"->'data_fr') "data" FROM repo LIMIT 10`,
		},
		{
			name:        "SELECT id, type, data_fr (LangEN)",
			query:       somesql.NewSelect("en").Fields("id", "type", "data.data_fr"),
			expectedSQL: `SELECT "id", "type", json_build_object('data_fr', "data_en"->'data_fr') "data" FROM repo LIMIT 10`,
		},
		{
			name:        "SELECT id, type, data_fr (LangFR)",
			query:       somesql.NewSelect("fr").Fields("id", "type", "data.data_fr"),
			expectedSQL: `SELECT "id", "type", json_build_object('data_fr', "data_fr"->'data_fr') "data" FROM repo LIMIT 10`,
		},
		// Select pre-defined fields and json attributes (any other) from data_*
		{
			name:        "SELECT id, type, data_en->'body'",
			query:       somesql.NewSelect("en").Fields("id", "type", "data.body"),
			expectedSQL: `SELECT "id", "type", json_build_object('body', "data_en"->'body') "data" FROM repo LIMIT 10`,
		},
		{
			name:        "SELECT id, type, data_en->'body' (LangEN)",
			query:       somesql.NewSelect("en").Fields("id", "type", "data.body"),
			expectedSQL: `SELECT "id", "type", json_build_object('body', "data_en"->'body') "data" FROM repo LIMIT 10`,
		},
		{
			name:        "SELECT id, type, data_fr->'body' (LangFR)",
			query:       somesql.NewSelect("fr").Fields("id", "type", "data.body"),
			expectedSQL: `SELECT "id", "type", json_build_object('body', "data_fr"->'body') "data" FROM repo LIMIT 10`,
		},
		// Select pre-defined fields and json attributes (any other + compound) from data_*
		{
			name:        "SELECT id, type, data_en->'body', data_en->'author_id'",
			query:       somesql.NewSelect("en").Fields("id", "type", "data.body", "data.author_id"),
			expectedSQL: `SELECT "id", "type", json_build_object('body', "data_en"->'body', 'author_id', "data_en"->'author_id') "data" FROM repo LIMIT 10`,
		},
		{
			name:        "SELECT id, type, data_en->'body', data_en->'author_id' (LangEN)",
			query:       somesql.NewSelect("en").Fields("id", "type", "data.body", "data.author_id"),
			expectedSQL: `SELECT "id", "type", json_build_object('body', "data_en"->'body', 'author_id', "data_en"->'author_id') "data" FROM repo LIMIT 10`,
		},
		{
			name:        "SELECT id, type, data_fr->'body', data_fr->'author_id' (LangFR)",
			query:       somesql.NewSelect("fr").Fields("id", "type", "data.body", "data.author_id"),
			expectedSQL: `SELECT "id", "type", json_build_object('body', "data_fr"->'body', 'author_id', "data_fr"->'author_id') "data" FROM repo LIMIT 10`,
		},
		// SELECT with conditions
		{
			name:           "SELECT * with condition",
			query:          somesql.NewSelect("en").Where(somesql.And("en", "id", "=", "uuid")),
			expectedSQL:    `SELECT "id", "created_at", "updated_at", "owner_id", "type", "data_en" FROM repo WHERE "id" = $1 LIMIT 10`,
			checkValues:    true,
			expectedValues: []interface{}{"uuid"},
		},
		{
			name:           "SELECT * with condition ORDER",
			query:          somesql.NewSelect("en").Where(somesql.And("en", "id", "=", "uuid")).Order("name", true),
			expectedSQL:    `SELECT "id", "created_at", "updated_at", "owner_id", "type", "data_en" FROM repo WHERE "id" = $1 ORDER BY "data_en"->>'name' ASC LIMIT 10`,
			checkValues:    true,
			expectedValues: []interface{}{"uuid"},
		},
		// SELECT relations
		{
			name:        "SELECT id, data->rel",
			query:       somesql.NewSelect("en").Fields("id", "relations.author", "relations.tags"),
			expectedSQL: `SELECT "id", json_build_object('author', "data_en"->'author', 'tags', "data_en"->'tags') "data" FROM repo LIMIT 10`,
		},
		{
			name:           "SELECT * + relations->rel, with conditions",
			query:          somesql.NewSelect("en").Fields("relations.author", "relations.tags").Where(somesql.And("en", "id", "=", "uuid")),
			expectedSQL:    `SELECT json_build_object('author', "data_en"->'author', 'tags', "data_en"->'tags') "data" FROM repo WHERE "id" = $1 LIMIT 10`,
			checkValues:    true,
			expectedValues: []interface{}{"uuid"},
		},
		{
			name:           "SELECT relations->rel only, with conditions",
			query:          somesql.NewSelect("en").Fields("relations.author", "relations.tags").Where(somesql.And("en", "id", "=", "uuid")),
			expectedSQL:    `SELECT json_build_object('author', "data_en"->'author', 'tags', "data_en"->'tags') "data" FROM repo WHERE "id" = $1 LIMIT 10`,
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
			query:          somesql.NewInsert("en").Fields(somesql.NewFields().UseDefaults().ID("1").CreatedAt(time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)).UpdatedAt(time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC))),
			expectedSQL:    `INSERT INTO repo ("id", "created_at", "updated_at", "owner_id", "type", "data_en") VALUES ($1, $2, $3, $4, $5, $6)`,
			expectedValues: []interface{}{"1", time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC), time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC), uuid.Nil.String(), "", "{}"},
		},
		{
			name:           "INSERT no defaults",
			query:          somesql.NewInsert("en").Fields(somesql.NewFields().ID("1").Type("entityA")),
			expectedSQL:    `INSERT INTO repo ("id", "type") VALUES ($1, $2)`,
			expectedValues: []interface{}{"1", "entityA"},
		},
		{
			name:           "INSERT no default + 1 data + 1 relation",
			query:          somesql.NewInsert("en").Fields(somesql.NewFields().Type("entityA").Set("data.body", "abc").Set("relations.tags", []string{"a", "b", "c"})),
			expectedSQL:    `INSERT INTO repo ("type", "data_en") VALUES ($1, $2)`,
			expectedValues: []interface{}{"entityA", `{"body":"abc","tags":["a","b","c"]}`},
		},
		{
			name:           "INSERT no default + 1 relation",
			query:          somesql.NewInsert("en").Fields(somesql.NewFields().Type("entityA").Set("relations.tags", []string{"a", "b", "c"})),
			expectedSQL:    `INSERT INTO repo ("type", "data_en") VALUES ($1, $2)`,
			expectedValues: []interface{}{"entityA", `{"tags":["a","b","c"]}`},
		},
		{
			name:           "INSERT no default + 2 relations",
			query:          somesql.NewInsert("en").Fields(somesql.NewFields().Type("entityA").Set("relations.author", []string{"x"}).Set("relations.tags", []string{"a", "b", "c"})),
			expectedSQL:    `INSERT INTO repo ("type", "data_en") VALUES ($1, $2)`,
			expectedValues: []interface{}{"entityA", `{"author":["x"],"tags":["a","b","c"]}`},
		},
		{
			name:           "INSERT no default + 1 relation one value",
			query:          somesql.NewInsert("en").Fields(somesql.NewFields().Type("entityA").Set("relations.tags", "a")),
			expectedSQL:    `INSERT INTO repo ("type", "data_en") VALUES ($1, $2)`,
			expectedValues: []interface{}{"entityA", `{"tags":["a"]}`},
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
			query:          somesql.NewUpdate("en").Fields(somesql.NewFields().ID("1").Type("entityA")),
			expectedSQL:    `UPDATE repo SET "id" = $1, "type" = $2`,
			expectedValues: []interface{}{"1", "entityA"},
		},
		{
			name:           "UPDATE data fields",
			query:          somesql.NewUpdate("en").Fields(somesql.NewFields().Set("data.body", "body value").Set("data.author_id", "123")),
			expectedSQL:    `UPDATE repo SET "data_en" = jsonb_build_object('body', $1::text, 'author_id', $2::text)::JSONB`,
			expectedValues: []interface{}{"body value", "123"},
		},
		{
			name:           "UPDATE data fields (LangFR)",
			query:          somesql.NewUpdate("fr").Fields(somesql.NewFields().Set("data.body", "body value").Set("data.author_id", "123")),
			expectedSQL:    `UPDATE repo SET "data_fr" = jsonb_build_object('body', $1::text, 'author_id', $2::text)::JSONB`,
			expectedValues: []interface{}{"body value", "123"},
		},
		{
			name:           "UPDATE meta + data fields",
			query:          somesql.NewUpdate("en").Fields(somesql.NewFields().ID("1").Type("entityA").Set("data.body", "body value").Set("data.author_id", "123")),
			expectedSQL:    `UPDATE repo SET "id" = $1, "type" = $2, "data_en" = jsonb_build_object('body', $3::text, 'author_id', $4::text)::JSONB`,
			expectedValues: []interface{}{"1", "entityA", "body value", "123"},
		},
		{
			name:           "UPDATE meta + data fields (LangFR)",
			query:          somesql.NewUpdate("fr").Fields(somesql.NewFields().ID("1").Type("entityA").Set("data.body", "body value").Set("data.author_id", "123")),
			expectedSQL:    `UPDATE repo SET "id" = $1, "type" = $2, "data_fr" = jsonb_build_object('body', $3::text, 'author_id', $4::text)::JSONB`,
			expectedValues: []interface{}{"1", "entityA", "body value", "123"},
		},
		{
			name:           "UPDATE meta + data fields conditions",
			query:          somesql.NewUpdate("en").Fields(somesql.NewFields().ID("1").Type("entityA").Set("data.body", "body value").Set("data.author_id", "123")).Where(somesql.And("en", "id", "=", "234")),
			expectedSQL:    `UPDATE repo SET "id" = $1, "type" = $2, "data_en" = jsonb_build_object('body', $3::text, 'author_id', $4::text)::JSONB WHERE "id" = $5`,
			expectedValues: []interface{}{"1", "entityA", "body value", "123", "234"},
		},
		{
			name:           "UPDATE meta + data fields conditions 2",
			query:          somesql.NewUpdate("en").Fields(somesql.NewFields().ID("1").Type("entityA").Set("data.body", "body value").Set("data.author_id", "123")).Where(somesql.And("en", "data.author_id", "=", "234")),
			expectedSQL:    `UPDATE repo SET "id" = $1, "type" = $2, "data_en" = jsonb_build_object('body', $3::text, 'author_id', $4::text)::JSONB WHERE "data_en"->>'author_id' = $5`,
			expectedValues: []interface{}{"1", "entityA", "body value", "123", "234"},
		},
		// Update relations
		{
			name:           "UPDATE set relation only",
			query:          somesql.NewUpdate("en").Fields(somesql.NewFields().Set("relations.tags", []string{"a", "b"})),
			expectedSQL:    `UPDATE repo SET "data_en" = jsonb_build_object('tags', $1::JSONB)::JSONB`,
			expectedValues: []interface{}{`["a","b"]`},
		},
		{
			name:           "UPDATE set relation with data",
			query:          somesql.NewUpdate("en").Fields(somesql.NewFields().Set("data.body", "body value").Set("relations.tags", []string{"a", "b"})),
			expectedSQL:    `UPDATE repo SET "data_en" = jsonb_build_object('body', $1::text, 'tags', $2::JSONB)::JSONB`,
			expectedValues: []interface{}{"body value", `["a","b"]`},
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
			query:       somesql.NewDelete("en"),
			expectedSQL: `DELETE FROM repo`,
		},
		{
			name:        "DELETE * LIMIT",
			query:       somesql.NewDelete("en").Limit(20),
			expectedSQL: `DELETE FROM repo LIMIT 20`,
		},
		{
			name:        "DELETE * OFFSET 10",
			query:       somesql.NewDelete("en").Offset(10),
			expectedSQL: `DELETE FROM repo OFFSET 10`,
		},
		{
			name:        "DELETE * LIMIT 20 OFFSET 10",
			query:       somesql.NewDelete("en").Limit(20).Offset(10),
			expectedSQL: `DELETE FROM repo LIMIT 20 OFFSET 10`,
		},
		{
			name:           "DELETE with condition",
			query:          somesql.NewDelete("en").Where(somesql.And("en", "id", "=", "uuid")),
			expectedSQL:    `DELETE FROM repo WHERE "id" = $1`,
			checkValues:    true,
			expectedValues: []interface{}{"uuid"},
		},
		{
			name:           "DELETE with conditions + relations",
			query:          somesql.NewDelete("en").Where(somesql.And("en", "relations.article", "=", "uuid")),
			expectedSQL:    `DELETE FROM repo WHERE (jsonb_path_exists("data_en", '$.article[*] ? (@ == $val)', json_object(ARRAY['val', $1])::jsonb))`,
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

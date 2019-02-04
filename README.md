# SomeSQL

SomeSQL is a simple Query builder which stores data objects as JSON in a supported DBMS backend.

## DB Tables

- Repo
- Slugs
- Archives
- Cards
- CardSchedules

```postgresql
DROP TABLE IF EXISTS repo;
CREATE TABLE repo (
    "id"         UUID      PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT (NOW() AT TIME ZONE 'Indian/Mauritius'::TEXT),
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT (NOW() AT TIME ZONE 'Indian/Mauritius'::TEXT),
    "owner_id"   UUID      NOT NULL DEFAULT uuid_nil(),
    "status"     TEXT      NOT NULL DEFAULT 'draft',
    "type"       TEXT      NOT NULL DEFAULT 'article',
    "data_en"    JSONB     NOT NULL DEFAULT '[{}]',
    "data_fr"    JSONB     NOT NULL DEFAULT '[{}]'
);
CREATE INDEX repo__status ON repo("status");
CREATE INDEX repo__type   ON repo("type"); 

DROP TABLE IF EXISTS slugs;
CREATE TABLE slugs (
    "repo_id" UUID NOT NULL DEFAULT uuid_nil(),
    "path"    TEXT NOT NULL DEFAULT '',
    "lang"    TEXT NOT NULL DEFAULT 'en'
);
CREATE INDEX slugs__path ON slugs("path");

DROP TABLE IF EXISTS archives;
CREATE TABLE archives (
    "id"          UUID      PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    "repo_id"     UUID      NOT NULL DEFAULT uuid_nil(),
    "archived_at" TIMESTAMP WITH TIME ZONE DEFAULT (NOW() AT TIME ZONE 'Indian/Mauritius'::TEXT),
    "archive"     JSONB     NOT NULL DEFAULT '[{}]'
);
CREATE INDEX archives__repo_id ON archives("repo_id");

DROP TABLE IF EXISTS cards;
CREATE TABLE cards (
	"id"                UUID    PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
	"status"            TEXT    NOT NULL DEFAULT 'draft',
	"deck_machine_name" TEXT    NOT NULL,
	"position"          INTEGER NOT NULL,
	"entity"            UUID    NOT NULL DEFAULT uuid_nil(),
	"entity_type"       TEXT    NOT NULL DEFAULT 'article'
);
CREATE INDEX cards__deck_machine_name ON cards("deck_machine_name");
CREATE INDEX cards__entity            ON cards("entity");

DROP TABLE IF EXISTS cardschedules;
CREATE TABLE cardschedules (
	"id"        UUID      PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
	"status"    TEXT      NOT NULL DEFAULT 'draft',
	"date_time" TIMESTAMP WITH TIME ZONE DEFAULT (NOW() AT TIME ZONE 'Indian/Mauritius'::TEXT),
	"action"    TEXT      NOT NULL DEFAULT '',
	"card_id"   UUID      NOT NULL DEFAULT uuid_nil()
);
CREATE INDEX cardschedules__card_id ON cardschedules("card_id");

ALTER TABLE slugs         ADD CONSTRAINT slugs__repo_id_fk         FOREIGN KEY ("repo_id") REFERENCES repo("id");
ALTER TABLE archives      ADD CONSTRAINT archives__repo_id_fk      FOREIGN KEY ("repo_id") REFERENCES repo("id");
ALTER TABLE cards         ADD CONSTRAINT cards__entity_fk          FOREIGN KEY ("entity")  REFERENCES repo("id");
ALTER TABLE cardschedules ADD CONSTRAINT cardschedules__card_id_fk FOREIGN KEY ("card_id") REFERENCES cards("id");
```

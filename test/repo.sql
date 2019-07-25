CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

DROP TABLE IF EXISTS repo;
CREATE TABLE repo (
    "id"         UUID       PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT (NOW() AT TIME ZONE 'Indian/Mauritius'::TEXT),
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT (NOW() AT TIME ZONE 'Indian/Mauritius'::TEXT),
    "owner_id"   UUID       NOT NULL DEFAULT uuid_nil(),
    "type"       TEXT       NOT NULL DEFAULT 'article',
    "data_en"    JSONB      NOT NULL DEFAULT '{}'
);

-- Authors
INSERT INTO repo ("id", "created_at", "updated_at", "owner_id", "type", "data_en") VALUES ('002fd6b1-f715-4875-838b-1546f27327df', NOW(), NOW(), uuid_nil(), 'author', '{"name":"John Doe", "position":"Content writer"}, "article":["99885fcb-9c64-4fb8-87a8-c88a96054325", "0f62965b-8b00-4f7b-8112-de3544ff242a", "8d38b883-0227-4433-83b9-a8c4dfa02798", "bcf57169-1d64-4800-a813-3f845831f12f", "0296f909-ac98-4ab8-a10f-b12de02bf22b"]');
INSERT INTO repo ("id", "created_at", "updated_at", "owner_id", "type", "data_en") VALUES ('d0dbc1d8-3f76-47d8-95f3-d20b1335f8cc', NOW(), NOW(), uuid_nil(), 'author', '{"name":"Jane Doe", "position":"Content writer"}, "article":["51639d71-0e56-4c75-b0ab-c2d9b6a80fac", "b498568c-a306-41bd-b2b5-1a06927b65e2", "62cda2cb-c095-424c-9bbb-96cc02c4eedc", "8722ef1a-2965-40aa-a143-e86d4328f36b", "4eaecffd-95b6-4c21-ad66-179785229faa"]');

-- Categories
INSERT INTO repo ("id", "created_at", "updated_at", "owner_id", "type", "data_en") VALUES ('90d84d12-5fac-4048-b661-54d0784b09b4', NOW(), NOW(), uuid_nil(), 'category', '{"parent":"00000000-0000-0000-0000-000000000000", "name":"Politique", "article":["99885fcb-9c64-4fb8-87a8-c88a96054325", "8722ef1a-2965-40aa-a143-e86d4328f36b"]}');
INSERT INTO repo ("id", "created_at", "updated_at", "owner_id", "type", "data_en") VALUES ('c58eacdb-0133-48d5-87bb-d8cbf73ee8e4', NOW(), NOW(), uuid_nil(), 'category', '{"parent":"00000000-0000-0000-0000-000000000000", "name":"Société", "article":["99885fcb-9c64-4fb8-87a8-c88a96054325", "51639d71-0e56-4c75-b0ab-c2d9b6a80fac", "0296f909-ac98-4ab8-a10f-b12de02bf22b"]}');
INSERT INTO repo ("id", "created_at", "updated_at", "owner_id", "type", "data_en") VALUES ('73b8780e-de03-4bfd-8d56-da2420ca7060', NOW(), NOW(), uuid_nil(), 'category', '{"parent":"00000000-0000-0000-0000-000000000000", "name":"Économie", "article":["51639d71-0e56-4c75-b0ab-c2d9b6a80fac", "0f62965b-8b00-4f7b-8112-de3544ff242a", "4eaecffd-95b6-4c21-ad66-179785229faa"]}');
INSERT INTO repo ("id", "created_at", "updated_at", "owner_id", "type", "data_en") VALUES ('e9e8c53c-2609-4903-bca5-6051ed8da97d', NOW(), NOW(), uuid_nil(), 'category', '{"parent":"00000000-0000-0000-0000-000000000000", "name":"Régions", "article":["0f62965b-8b00-4f7b-8112-de3544ff242a", "b498568c-a306-41bd-b2b5-1a06927b65e2"]}');
INSERT INTO repo ("id", "created_at", "updated_at", "owner_id", "type", "data_en") VALUES ('6ba5e86f-7c80-4d47-a5bd-9a21228de2d6', NOW(), NOW(), uuid_nil(), 'category', '{"parent":"00000000-0000-0000-0000-000000000000", "name":"Idées", "article":["b498568c-a306-41bd-b2b5-1a06927b65e2", "8d38b883-0227-4433-83b9-a8c4dfa02798"]}');
INSERT INTO repo ("id", "created_at", "updated_at", "owner_id", "type", "data_en") VALUES ('ae2fd1ec-801f-43c8-b2fc-06e6ff7fe882', NOW(), NOW(), uuid_nil(), 'category', '{"parent":"00000000-0000-0000-0000-000000000000", "name":"Sport", "article":["8d38b883-0227-4433-83b9-a8c4dfa02798", "62cda2cb-c095-424c-9bbb-96cc02c4eedc"]}');
INSERT INTO repo ("id", "created_at", "updated_at", "owner_id", "type", "data_en") VALUES ('770e77e6-179b-4bbf-a289-c5e3a46e9769', NOW(), NOW(), uuid_nil(), 'category', '{"parent":"00000000-0000-0000-0000-000000000000", "name":"International", "article":["62cda2cb-c095-424c-9bbb-96cc02c4eedc", "bcf57169-1d64-4800-a813-3f845831f12f"]}');
INSERT INTO repo ("id", "created_at", "updated_at", "owner_id", "type", "data_en") VALUES ('b2147e43-3f38-44b4-8d0e-7d7b08e26527', NOW(), NOW(), uuid_nil(), 'category', '{"parent":"00000000-0000-0000-0000-000000000000", "name":"Vous", "article":["bcf57169-1d64-4800-a813-3f845831f12f", "8722ef1a-2965-40aa-a143-e86d4328f36b"]}');
INSERT INTO repo ("id", "created_at", "updated_at", "owner_id", "type", "data_en") VALUES ('ae56d75e-1511-4c1b-9808-f60b43fe8f23', NOW(), NOW(), uuid_nil(), 'category', '{"parent":"00000000-0000-0000-0000-000000000000", "name":"Multimédia", "article":["0296f909-ac98-4ab8-a10f-b12de02bf22b"]}');
INSERT INTO repo ("id", "created_at", "updated_at", "owner_id", "type", "data_en") VALUES ('34f85fd4-e213-40a5-9403-c4043347580b', NOW(), NOW(), uuid_nil(), 'category', '{"parent":"00000000-0000-0000-0000-000000000000", "name":"Petites Annonces", "article":["4eaecffd-95b6-4c21-ad66-179785229faa"]}');

-- Articles
INSERT INTO repo ("id", "created_at", "updated_at", "owner_id", "type", "data_en") VALUES ('99885fcb-9c64-4fb8-87a8-c88a96054325', NOW(), NOW(), uuid_nil(), 'article', '{"title":"Title of article 01", "index":0,"type":"Legacy","name":"Legacy content, nid=1001","body":"Lorem ipsum dolor sit amet, consectetur adipi", "author":["002fd6b1-f715-4875-838b-1546f27327df"], "category":["90d84d12-5fac-4048-b661-54d0784b09b4", "c58eacdb-0133-48d5-87bb-d8cbf73ee8e4"]}');
INSERT INTO repo ("id", "created_at", "updated_at", "owner_id", "type", "data_en") VALUES ('51639d71-0e56-4c75-b0ab-c2d9b6a80fac', NOW(), NOW(), uuid_nil(), 'article', '{"title":"Title of article 02", "index":0,"type":"Legacy","name":"Legacy content, nid=1002","body":"Sed hendrerit lorem vel nunc rutrum, eu digni", "author":["d0dbc1d8-3f76-47d8-95f3-d20b1335f8cc"], "category":["c58eacdb-0133-48d5-87bb-d8cbf73ee8e4", "73b8780e-de03-4bfd-8d56-da2420ca7060"]}');
INSERT INTO repo ("id", "created_at", "updated_at", "owner_id", "type", "data_en") VALUES ('0f62965b-8b00-4f7b-8112-de3544ff242a', NOW(), NOW(), uuid_nil(), 'article', '{"title":"Title of article 03", "index":0,"type":"Legacy","name":"Legacy content, nid=1003","body":"Donec eu tortor semper, mattis turpis commodo", "author":["002fd6b1-f715-4875-838b-1546f27327df"], "category":["73b8780e-de03-4bfd-8d56-da2420ca7060", "e9e8c53c-2609-4903-bca5-6051ed8da97d"]}');
INSERT INTO repo ("id", "created_at", "updated_at", "owner_id", "type", "data_en") VALUES ('b498568c-a306-41bd-b2b5-1a06927b65e2', NOW(), NOW(), uuid_nil(), 'article', '{"title":"Title of article 04", "index":0,"type":"Legacy","name":"Legacy content, nid=1004","body":"Donec in lacus fermentum, iaculis orci rutrum", "author":["d0dbc1d8-3f76-47d8-95f3-d20b1335f8cc"], "category":["e9e8c53c-2609-4903-bca5-6051ed8da97d", "6ba5e86f-7c80-4d47-a5bd-9a21228de2d6"]}');
INSERT INTO repo ("id", "created_at", "updated_at", "owner_id", "type", "data_en") VALUES ('8d38b883-0227-4433-83b9-a8c4dfa02798', NOW(), NOW(), uuid_nil(), 'article', '{"title":"Title of article 05", "index":0,"type":"Legacy","name":"Legacy content, nid=1005","body":"Pellentesque eget nisl eget nulla lacinia sus", "author":["002fd6b1-f715-4875-838b-1546f27327df"], "category":["6ba5e86f-7c80-4d47-a5bd-9a21228de2d6", "ae2fd1ec-801f-43c8-b2fc-06e6ff7fe882"]}');
INSERT INTO repo ("id", "created_at", "updated_at", "owner_id", "type", "data_en") VALUES ('62cda2cb-c095-424c-9bbb-96cc02c4eedc', NOW(), NOW(), uuid_nil(), 'article', '{"title":"Title of article 06", "index":0,"type":"Legacy","name":"Legacy content, nid=1006","body":"Aenean et turpis sed mauris sagittis tincidun", "author":["d0dbc1d8-3f76-47d8-95f3-d20b1335f8cc"], "category":["ae2fd1ec-801f-43c8-b2fc-06e6ff7fe882", "770e77e6-179b-4bbf-a289-c5e3a46e9769"]}');
INSERT INTO repo ("id", "created_at", "updated_at", "owner_id", "type", "data_en") VALUES ('bcf57169-1d64-4800-a813-3f845831f12f', NOW(), NOW(), uuid_nil(), 'article', '{"title":"Title of article 07", "index":0,"type":"Legacy","name":"Legacy content, nid=1007","body":"Ut tempus ante vitae nisi laoreet, a variusos", "author":["002fd6b1-f715-4875-838b-1546f27327df"], "category":["770e77e6-179b-4bbf-a289-c5e3a46e9769", "b2147e43-3f38-44b4-8d0e-7d7b08e26527"]}');
INSERT INTO repo ("id", "created_at", "updated_at", "owner_id", "type", "data_en") VALUES ('8722ef1a-2965-40aa-a143-e86d4328f36b', NOW(), NOW(), uuid_nil(), 'article', '{"title":"Title of article 08", "index":0,"type":"Legacy","name":"Legacy content, nid=1008","body":"Duis efficitur eros sit amet nisi hendreritoi", "author":["d0dbc1d8-3f76-47d8-95f3-d20b1335f8cc"], "category":["b2147e43-3f38-44b4-8d0e-7d7b08e26527", "90d84d12-5fac-4048-b661-54d0784b09b4"]}');
INSERT INTO repo ("id", "created_at", "updated_at", "owner_id", "type", "data_en") VALUES ('0296f909-ac98-4ab8-a10f-b12de02bf22b', NOW(), NOW(), uuid_nil(), 'article', '{"title":"Title of article 09", "index":0,"type":"Legacy","name":"Legacy content, nid=1009","body":"Phasellus consectetur sapien varius diam laor", "author":["002fd6b1-f715-4875-838b-1546f27327df"], "category":["ae56d75e-1511-4c1b-9808-f60b43fe8f23", "c58eacdb-0133-48d5-87bb-d8cbf73ee8e4"]}');
INSERT INTO repo ("id", "created_at", "updated_at", "owner_id", "type", "data_en") VALUES ('4eaecffd-95b6-4c21-ad66-179785229faa', NOW(), NOW(), uuid_nil(), 'article', '{"title":"Title of article 10", "index":0,"type":"Legacy","name":"Legacy content, nid=1010","body":"Pellentesque a nibh eget nisi aliquam ultrici", "author":["d0dbc1d8-3f76-47d8-95f3-d20b1335f8cc"], "category":["34f85fd4-e213-40a5-9403-c4043347580b", "73b8780e-de03-4bfd-8d56-da2420ca7060"]}');

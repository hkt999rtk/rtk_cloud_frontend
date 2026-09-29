-- Schema documentation metadata. Keep changes beside schema changes.
CREATE TABLE IF NOT EXISTS schema_metadata (
  object_kind TEXT NOT NULL CHECK (object_kind IN ('group','table','column','logical_ref')),
  schema_name TEXT NOT NULL DEFAULT 'main',
  object_name TEXT NOT NULL,
  meta_version INTEGER NOT NULL DEFAULT 1 CHECK (meta_version = 1),
  details TEXT NOT NULL CHECK (json_valid(details)),
  PRIMARY KEY (object_kind, schema_name, object_name)
);
INSERT INTO schema_metadata (object_kind,schema_name,object_name,details) VALUES ('group','main','Website search index','{"purpose":"Searchable documents, text chunks, and embedding data.","scenario":"Used to rebuild the index, search site documents, or return results.","service":"Cloud Frontend"}') ON CONFLICT (object_kind,schema_name,object_name) DO UPDATE SET details=excluded.details;
INSERT INTO schema_metadata (object_kind,schema_name,object_name,details) VALUES ('table','main','search_chunks','{"group":"Website search index","purpose":"Search-document text chunks and embedding JSON.","scenario":"Used when building a semantic search index or retrieving similar passages.","service":"Cloud Frontend","source":"repos/rtk_cloud_frontend/internal/search/repository.go"}') ON CONFLICT (object_kind,schema_name,object_name) DO UPDATE SET details=excluded.details;
INSERT INTO schema_metadata (object_kind,schema_name,object_name,details) VALUES ('column','main','search_chunks.id','{"critical":"key","description":"Primary key of search_chunks."}') ON CONFLICT (object_kind,schema_name,object_name) DO UPDATE SET details=excluded.details;
INSERT INTO schema_metadata (object_kind,schema_name,object_name,details) VALUES ('column','main','search_chunks.document_id','{"critical":"key","description":"References search_documents.id within this database."}') ON CONFLICT (object_kind,schema_name,object_name) DO UPDATE SET details=excluded.details;
INSERT INTO schema_metadata (object_kind,schema_name,object_name,details) VALUES ('table','main','search_documents','{"group":"Website search index","purpose":"Searchable document titles, URLs, locales, sources, and bodies.","scenario":"Imported or updated for the website index and displayed in search results.","service":"Cloud Frontend","source":"repos/rtk_cloud_frontend/internal/search/repository.go"}') ON CONFLICT (object_kind,schema_name,object_name) DO UPDATE SET details=excluded.details;
INSERT INTO schema_metadata (object_kind,schema_name,object_name,details) VALUES ('column','main','search_documents.id','{"critical":"key","description":"Primary key of search_documents."}') ON CONFLICT (object_kind,schema_name,object_name) DO UPDATE SET details=excluded.details;

-- Schema documentation metadata. Keep changes beside schema changes.
CREATE TABLE IF NOT EXISTS schema_metadata (
  object_kind TEXT NOT NULL CHECK (object_kind IN ('group','table','column','logical_ref')),
  schema_name TEXT NOT NULL DEFAULT 'main',
  object_name TEXT NOT NULL,
  meta_version INTEGER NOT NULL DEFAULT 1 CHECK (meta_version = 1),
  details TEXT NOT NULL CHECK (json_valid(details)),
  PRIMARY KEY (object_kind, schema_name, object_name)
);
INSERT INTO schema_metadata (object_kind,schema_name,object_name,details) VALUES ('group','main','Website interactions and leads','{"purpose":"Website analytics events, SDK-term acceptances, and visitor contact leads.","scenario":"Used when visitors browse the site, download an SDK, or submit a contact form.","service":"Cloud Frontend"}') ON CONFLICT (object_kind,schema_name,object_name) DO UPDATE SET details=excluded.details;
INSERT INTO schema_metadata (object_kind,schema_name,object_name,details) VALUES ('table','main','leads','{"group":"Website interactions and leads","purpose":"Visitor names, companies, contact details, requirements, and messages.","scenario":"Created when someone submits a contact or sales inquiry.","service":"Cloud Frontend","source":"repos/rtk_cloud_frontend/internal/leads/repository.go"}') ON CONFLICT (object_kind,schema_name,object_name) DO UPDATE SET details=excluded.details;
INSERT INTO schema_metadata (object_kind,schema_name,object_name,details) VALUES ('column','main','leads.id','{"critical":"key","description":"Primary key of leads."}') ON CONFLICT (object_kind,schema_name,object_name) DO UPDATE SET details=excluded.details;

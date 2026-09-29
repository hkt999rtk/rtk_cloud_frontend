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
INSERT INTO schema_metadata (object_kind,schema_name,object_name,details) VALUES ('table','main','analytics_events','{"group":"Website interactions and leads","purpose":"Website page, call-to-action, dwell-time, and referral analytics events.","scenario":"Records visitor behavior while they browse or interact with website elements.","service":"Cloud Frontend","source":"repos/rtk_cloud_frontend/internal/analytics/repository.go"}') ON CONFLICT (object_kind,schema_name,object_name) DO UPDATE SET details=excluded.details;
INSERT INTO schema_metadata (object_kind,schema_name,object_name,details) VALUES ('column','main','analytics_events.id','{"critical":"key","description":"Primary key of analytics_events."}') ON CONFLICT (object_kind,schema_name,object_name) DO UPDATE SET details=excluded.details;
INSERT INTO schema_metadata (object_kind,schema_name,object_name,details) VALUES ('table','main','sdk_download_acceptances','{"group":"Website interactions and leads","purpose":"SDK download terms versions, packages, sessions, and request IDs.","scenario":"Records acceptance before providing an SDK package.","service":"Cloud Frontend","source":"repos/rtk_cloud_frontend/internal/analytics/repository.go"}') ON CONFLICT (object_kind,schema_name,object_name) DO UPDATE SET details=excluded.details;
INSERT INTO schema_metadata (object_kind,schema_name,object_name,details) VALUES ('column','main','sdk_download_acceptances.id','{"critical":"key","description":"Primary key of sdk_download_acceptances."}') ON CONFLICT (object_kind,schema_name,object_name) DO UPDATE SET details=excluded.details;

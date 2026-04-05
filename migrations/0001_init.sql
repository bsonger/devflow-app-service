CREATE TABLE IF NOT EXISTS projects (
  id UUID PRIMARY KEY,
  key TEXT NOT NULL,
  name TEXT NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  namespace TEXT NOT NULL,
  owner TEXT NOT NULL DEFAULT '',
  labels JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  deleted_at TIMESTAMPTZ NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_projects_key_active
  ON projects (key)
  WHERE deleted_at IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS uq_projects_name_active
  ON projects (name)
  WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS applications (
  id UUID PRIMARY KEY,
  project_id UUID NOT NULL,
  name TEXT NOT NULL,
  repo_address TEXT NOT NULL,
  active_manifest_id UUID NULL,
  labels JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  deleted_at TIMESTAMPTZ NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_applications_project_name_active
  ON applications (project_id, name)
  WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_applications_project_id
  ON applications (project_id);

CREATE INDEX IF NOT EXISTS idx_applications_active_manifest_id
  ON applications (active_manifest_id);

CREATE TABLE IF NOT EXISTS services (
  id UUID PRIMARY KEY,
  application_id UUID NOT NULL,
  name TEXT NOT NULL,
  exposure TEXT NOT NULL,
  ports JSONB NOT NULL DEFAULT '[]'::jsonb,
  status TEXT NOT NULL DEFAULT 'active',
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  deleted_at TIMESTAMPTZ NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_services_name_active
  ON services (application_id, name)
  WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_services_application_id
  ON services (application_id);

CREATE TABLE IF NOT EXISTS environments (
  id UUID PRIMARY KEY,
  key TEXT NOT NULL,
  name TEXT NOT NULL,
  cluster TEXT NOT NULL,
  namespace TEXT NOT NULL,
  labels JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  deleted_at TIMESTAMPTZ NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_environments_key_active
  ON environments (key)
  WHERE deleted_at IS NULL;

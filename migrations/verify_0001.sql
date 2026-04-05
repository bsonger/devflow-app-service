DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'projects') THEN
    RAISE EXCEPTION 'missing table: projects';
  END IF;
  IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'applications') THEN
    RAISE EXCEPTION 'missing table: applications';
  END IF;
  IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'services') THEN
    RAISE EXCEPTION 'missing table: services';
  END IF;
  IF NOT EXISTS (
    SELECT 1
    FROM information_schema.columns
    WHERE table_name = 'applications' AND column_name = 'repo_address'
  ) THEN
    RAISE EXCEPTION 'missing column: applications.repo_address';
  END IF;
  IF NOT EXISTS (
    SELECT 1
    FROM information_schema.columns
    WHERE table_name = 'services' AND column_name = 'application_id'
  ) THEN
    RAISE EXCEPTION 'missing column: services.application_id';
  END IF;
  IF NOT EXISTS (
    SELECT 1
    FROM pg_indexes
    WHERE schemaname = 'public' AND indexname = 'uq_applications_project_name_active'
  ) THEN
    RAISE EXCEPTION 'missing index: uq_applications_project_name_active';
  END IF;
  IF NOT EXISTS (
    SELECT 1
    FROM pg_indexes
    WHERE schemaname = 'public' AND indexname = 'uq_services_name_active'
  ) THEN
    RAISE EXCEPTION 'missing index: uq_services_name_active';
  END IF;
END $$;

DO $$
BEGIN
  IF EXISTS (
    SELECT 1
    FROM information_schema.columns
    WHERE table_name = 'applications' AND column_name = 'active_manifest_id'
  ) AND NOT EXISTS (
    SELECT 1
    FROM information_schema.columns
    WHERE table_name = 'applications' AND column_name = 'active_image_id'
  ) THEN
    ALTER TABLE applications
      RENAME COLUMN active_manifest_id TO active_image_id;
  END IF;
END $$;

ALTER INDEX IF EXISTS idx_applications_active_manifest_id
  RENAME TO idx_applications_active_image_id;

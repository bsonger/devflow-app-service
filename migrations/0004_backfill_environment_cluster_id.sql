DO $$
BEGIN
  IF EXISTS (
    SELECT 1
    FROM information_schema.columns
    WHERE table_name = 'environments' AND column_name = 'cluster'
  ) AND NOT EXISTS (
    SELECT 1
    FROM information_schema.columns
    WHERE table_name = 'environments' AND column_name = 'cluster_id'
  ) THEN
    ALTER TABLE environments ADD COLUMN cluster_id UUID;
  END IF;
END $$;

DO $$
BEGIN
  IF EXISTS (
    SELECT 1
    FROM information_schema.columns
    WHERE table_name = 'environments' AND column_name = 'cluster'
  ) THEN
    EXECUTE $sql$
      UPDATE environments e
      SET cluster_id = c.id
      FROM clusters c
      WHERE e.cluster_id IS NULL
        AND c.name = e.cluster
    $sql$;
  END IF;
END $$;

DO $$
BEGIN
  IF EXISTS (
    SELECT 1
    FROM information_schema.columns
    WHERE table_name = 'environments' AND column_name = 'cluster_id'
  ) THEN
    IF EXISTS (
      SELECT 1 FROM environments WHERE cluster_id IS NULL AND deleted_at IS NULL
    ) THEN
      RAISE EXCEPTION 'cannot finalize environments.cluster_id migration: active environments still missing cluster mapping';
    END IF;

    ALTER TABLE environments
      ALTER COLUMN cluster_id SET NOT NULL;

    CREATE INDEX IF NOT EXISTS idx_environments_cluster_id
      ON environments (cluster_id);
  END IF;
END $$;

DO $$
BEGIN
  IF EXISTS (
    SELECT 1
    FROM information_schema.columns
    WHERE table_name = 'environments' AND column_name = 'key'
  ) THEN
    DROP INDEX IF EXISTS uq_environments_key_active;
    ALTER TABLE environments DROP COLUMN key;
  END IF;

  IF EXISTS (
    SELECT 1
    FROM information_schema.columns
    WHERE table_name = 'environments' AND column_name = 'cluster'
  ) THEN
    ALTER TABLE environments DROP COLUMN cluster;
  END IF;

  IF EXISTS (
    SELECT 1
    FROM information_schema.columns
    WHERE table_name = 'environments' AND column_name = 'namespace'
  ) THEN
    ALTER TABLE environments DROP COLUMN namespace;
  END IF;
END $$;

ALTER TABLE environments
  ALTER COLUMN labels SET DEFAULT '[]'::jsonb;

UPDATE environments
SET labels = '[]'::jsonb
WHERE labels = '{}'::jsonb;

DROP INDEX IF EXISTS uq_environments_name_active;
CREATE UNIQUE INDEX IF NOT EXISTS uq_environments_name_active
  ON environments (name)
  WHERE deleted_at IS NULL;

-- Drop the foreign key constraint first if explicitly named (optional but good practice)
-- ALTER TABLE watch_providers DROP CONSTRAINT IF EXISTS fk_watch_providers_anime;

-- Drop the watch_providers table
DROP TABLE IF EXISTS watch_providers;

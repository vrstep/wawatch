-- Drop the anime_caches table
DROP TABLE IF EXISTS watch_providers;

-- Indexes are typically dropped automatically when the table is dropped,
-- but you could explicitly drop them first if needed:
-- DROP INDEX IF EXISTS idx_watch_providers_anime_id;
-- DROP INDEX IF EXISTS idx_anime_caches_title_romaji;
-- DROP INDEX IF EXISTS idx_anime_caches_title_english;
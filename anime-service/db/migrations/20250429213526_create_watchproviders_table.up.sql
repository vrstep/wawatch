-- Create the watch_providers table
CREATE TABLE IF NOT EXISTS watch_providers (
    id SERIAL PRIMARY KEY, -- Use SERIAL for auto-incrementing ID specific to this table
    anime_id BIGINT NOT NULL,
    provider_name VARCHAR(255) NOT NULL,
    provider_url TEXT NOT NULL,
    provider_logo_url TEXT, -- Optional logo URL
    region VARCHAR(100), -- e.g., "US", "JP", "Global"
    type VARCHAR(50), -- e.g., "Stream", "Rent", "Buy"
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    -- Foreign key constraint linking to the anime_caches table
    CONSTRAINT fk_watch_providers_anime
        FOREIGN KEY (anime_id)
        REFERENCES anime_caches(id)
        ON DELETE CASCADE -- If an anime is removed from cache, remove its providers too
);

-- Add index for faster lookups by anime_id
CREATE INDEX IF NOT EXISTS idx_watch_providers_anime_id ON watch_providers (anime_id);

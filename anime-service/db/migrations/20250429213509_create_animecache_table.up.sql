-- Create the anime_caches table to store basic anime info fetched from external APIs
CREATE TABLE IF NOT EXISTS anime_caches (
    id BIGINT PRIMARY KEY, -- Assuming this corresponds to the external API's ID (e.g., AniList ID)
    title_romaji VARCHAR(255),
    title_english VARCHAR(255),
    title_native VARCHAR(255),
    type VARCHAR(50),
    format VARCHAR(50),
    status VARCHAR(50),
    description TEXT,
    start_date DATE,
    end_date DATE,
    season VARCHAR(50),
    season_year INT,
    episodes INT,
    duration INT,
    country_of_origin VARCHAR(10),
    source VARCHAR(50),
    cover_image_large TEXT,
    cover_image_medium TEXT,
    banner_image TEXT,
    average_score INT,
    mean_score INT,
    popularity INT,
    favourites INT,
    is_adult BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Add index for faster lookups if needed, e.g., by title
CREATE INDEX IF NOT EXISTS idx_anime_caches_title_romaji ON anime_caches (title_romaji);
CREATE INDEX IF NOT EXISTS idx_anime_caches_title_english ON anime_caches (title_english);

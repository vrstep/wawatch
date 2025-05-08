-- From: anime-service/db/migrations/000001_create_anime_tables.up.sql

CREATE TABLE IF NOT EXISTS anime_details (
    id BIGINT PRIMARY KEY, -- Corresponds to AnimeDetails.ID
    title_romaji VARCHAR(255), -- Corresponds to AnimeDetails.Title.Romaji
    title_english VARCHAR(255), -- Corresponds to AnimeDetails.Title.English
    title_native VARCHAR(255), -- Corresponds to AnimeDetails.Title.Native
    type VARCHAR(50), -- Not directly in AnimeDetails, maybe add if needed
    format VARCHAR(50), -- Corresponds to AnimeDetails.Format
    status VARCHAR(50), -- Corresponds to AnimeDetails.Status
    description TEXT, -- Corresponds to AnimeDetails.Description
    start_date DATE, -- Derived from AnimeDetails.StartDate
    end_date DATE, -- Derived from AnimeDetails.EndDate
    season VARCHAR(50), -- Corresponds to AnimeDetails.Season
    season_year INT, -- Corresponds to AnimeDetails.SeasonYear
    episodes INT, -- Corresponds to AnimeDetails.Episodes
    duration INT, -- Corresponds to AnimeDetails.Duration
    country_of_origin VARCHAR(10), -- Not directly in AnimeDetails, maybe add if needed
    source VARCHAR(50), -- Not directly in AnimeDetails, maybe add if needed
    cover_image_large TEXT, -- Corresponds to AnimeDetails.CoverImage.Large
    cover_image_medium TEXT, -- Corresponds to AnimeDetails.CoverImage.Medium
    banner_image TEXT, -- Corresponds to AnimeDetails.BannerImage
    average_score INT, -- Corresponds to AnimeDetails.AverageScore
    mean_score INT, -- Not directly in AnimeDetails, maybe add if needed
    popularity INT, -- Corresponds to AnimeDetails.Popularity
    favourites INT, -- Not directly in AnimeDetails, maybe add if needed
    is_adult BOOLEAN DEFAULT FALSE, -- Not directly in AnimeDetails, maybe add if needed
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
    -- Note: Genres and Studios are often handled in separate related tables if needed for querying,
    -- or stored as JSON/array types if the DB supports it and complex querying isn't needed.
);

-- ... indexes ...
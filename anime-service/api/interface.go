package api

import "github.com/vrstep/wawatch-backend/models"

// AniListAPI defines the interface for AniList client operations
type AniListAPI interface {
	GetAnimeByID(id int) (*models.AnimeDetails, error)
	SearchAnime(query string, page int, perPage int) ([]models.AnimeCache, int, error)
	GetPopularAnime(page int, perPage int) ([]models.AnimeCache, int, error)
	GetTrendingAnime(page int, perPage int) ([]models.AnimeCache, int, error)
	GetAnimeBySeason(year int, season string, page int, perPage int) ([]models.AnimeCache, int, error)
	GetUpcomingAnime(page int, perPage int) ([]models.AnimeCache, int, error)
	GetRecentlyReleasedAnime(page int, perPage int) ([]models.AnimeCache, int, error)
	GetAnimeByTags(tags []string, page int, perPage int) ([]models.AnimeCache, int, error) // Modified
}

// Ensure the real client implements the interface
var _ AniListAPI = (*AniListClient)(nil)

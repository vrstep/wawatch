// backend/client/interface.go
package client

import "github.com/vrstep/wawatch-backend/models"

type AnimeServiceAPIClient interface {
	WithRequestID(requestID string) AnimeServiceAPIClient
	GetAnimeDetailsAndProviders(animeID int) (*models.AnimeDetails, []models.WatchProvider, error)
	GetAnimeByID(animeID int) (*models.AnimeDetails, error)
	SearchAnime(query string, page, perPage int) ([]models.AnimeCache, int, error)
	GetPopularAnime(page, perPage int) ([]models.AnimeCache, int, error)
	GetTrendingAnime(page, perPage int) ([]models.AnimeCache, int, error)
	GetAnimeBySeason(year int, season string, page, perPage int) ([]models.AnimeCache, int, error)
	GetAnimeRecommendations(page, perPage int) ([]models.AnimeCache, int, error)
	ExploreAnime(tags []string, page, perPage int) ([]models.AnimeCache, int, error) // The missing one
	GetUpcomingAnime(page, perPage int) ([]models.AnimeCache, int, error)
	GetRecentlyReleasedAnime(page, perPage int) ([]models.AnimeCache, int, error)
}

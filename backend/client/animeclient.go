package client

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	// Assuming models are defined within this backend/user-service module
	"github.com/vrstep/wawatch-backend/models"
)

var _ AnimeServiceAPIClient = (*AnimeClient)(nil)

type AnimeClient struct {
	client    *resty.Client
	baseURL   string
	requestID string // Store requestID per instance for forwarding
}

// NewAnimeClient creates a new client for interacting with the anime-service.
// The baseURL for the anime-service should be configurable.
func NewAnimeClient() *AnimeClient {
	baseURL := os.Getenv("ANIME_SERVICE_URL") // This will now be "http://anime-service:8082"
	if baseURL == "" {
		// This default is only for running user-service outside Docker directly on host,
		// and anime-service also directly on host on port 8081.
		// When in Docker, ANIME_SERVICE_URL MUST be set.
		baseURL = "http://localhost:8081" // Fallback for local non-Docker dev, matches anime-service host port
		log.Printf("Warning: ANIME_SERVICE_URL environment variable not set. Using default for local dev: %s. THIS WILL FAIL IN DOCKER IF ANIME_SERVICE_URL IS NOT SET TO THE DOCKER SERVICE NAME.", baseURL)
	}
	baseURL = strings.TrimSuffix(baseURL, "/")

	client := resty.New().
		SetTimeout(15 * time.Second). // Slightly longer timeout for inter-service calls
		SetRetryCount(2).
		SetRetryWaitTime(300 * time.Millisecond).
		SetRetryMaxWaitTime(2 * time.Second)

	return &AnimeClient{
		client:  client,
		baseURL: baseURL,
	}
}

func (c *AnimeClient) WithRequestID(requestID string) AnimeServiceAPIClient {
	newC := *c // shallow copy
	newC.requestID = requestID
	return &newC // Return the concrete type, which satisfies the interface
}

// Helper to prepare a request with common settings like RequestID
func (c *AnimeClient) R() *resty.Request {
	req := c.client.R()
	if c.requestID != "" {
		req.SetHeader("X-Request-ID", c.requestID)
	}
	return req
}

// GetAnimeDetailsAndProviders fetches anime details and providers from the anime-service.
// The anime-service returns {"anime": ..., "providers": ...}
func (c *AnimeClient) GetAnimeDetailsAndProviders(animeID int) (*models.AnimeDetails, []models.WatchProvider, error) {
	var result struct {
		Anime     *models.AnimeDetails   `json:"anime"`
		Providers []models.WatchProvider `json:"providers"`
	}

	resp, err := c.R().
		SetResult(&result).
		Get(fmt.Sprintf("%s/anime/%d", c.baseURL, animeID))

	if err != nil {
		log.Printf("Error calling anime-service for details (ID: %d): %v", animeID, err)
		return nil, nil, fmt.Errorf("anime-service call failed: %w", err)
	}

	if !resp.IsSuccess() {
		log.Printf("anime-service returned error for details (ID: %d) - Status: %s, Body: %s", animeID, resp.Status(), resp.String())
		return nil, nil, fmt.Errorf("anime-service error (status %s): %s", resp.Status(), resp.String())
	}
	if result.Anime == nil {
		return nil, nil, fmt.Errorf("anime-service returned no anime data in expected structure for ID %d", animeID)
	}

	return result.Anime, result.Providers, nil
}

// GetAnimeByID is used by user_animelist_controller. It should get details from anime-service.
// It calls the same endpoint as GetAnimeDetailsAndProviders but extracts only AnimeDetails.
func (c *AnimeClient) GetAnimeByID(animeID int) (*models.AnimeDetails, error) {
	animeDetails, _, err := c.GetAnimeDetailsAndProviders(animeID)
	return animeDetails, err
}

// Helper for paged results from anime-service that use {"data": ..., "meta": ...} structure
type pagedAnimeCacheResult struct {
	Data []models.AnimeCache `json:"data"`
	Meta struct {
		Total       int  `json:"total"`
		Page        int  `json:"page"`
		PerPage     int  `json:"perPage"`
		TotalPages  int  `json:"totalPages"`
		HasNextPage bool `json:"hasNextPage"`
	} `json:"meta"`
}

// SearchAnime searches for anime through the anime-service
func (c *AnimeClient) SearchAnime(query string, page, perPage int) ([]models.AnimeCache, int, error) {
	var result pagedAnimeCacheResult
	resp, err := c.R().
		SetQueryParams(map[string]string{
			"q":       query,
			"page":    fmt.Sprintf("%d", page),
			"perPage": fmt.Sprintf("%d", perPage),
		}).
		SetResult(&result).
		Get(fmt.Sprintf("%s/anime/search", c.baseURL))

	if err != nil {
		return nil, 0, fmt.Errorf("anime-service call failed for search: %w", err)
	}
	if !resp.IsSuccess() {
		return nil, 0, fmt.Errorf("anime-service error on search (status %s): %s", resp.Status(), resp.String())
	}
	return result.Data, result.Meta.Total, nil
}

// GetPopularAnime from anime-service
func (c *AnimeClient) GetPopularAnime(page, perPage int) ([]models.AnimeCache, int, error) {
	var result pagedAnimeCacheResult
	resp, err := c.R().
		SetQueryParams(map[string]string{"page": fmt.Sprintf("%d", page), "perPage": fmt.Sprintf("%d", perPage)}).
		SetResult(&result).
		Get(fmt.Sprintf("%s/anime/popular", c.baseURL))

	if err != nil {
		return nil, 0, fmt.Errorf("anime-service call failed for popular: %w", err)
	}
	if !resp.IsSuccess() {
		return nil, 0, fmt.Errorf("anime-service error on popular (status %s): %s", resp.Status(), resp.String())
	}
	return result.Data, result.Meta.Total, nil
}

// GetTrendingAnime from anime-service
func (c *AnimeClient) GetTrendingAnime(page, perPage int) ([]models.AnimeCache, int, error) {
	var result pagedAnimeCacheResult
	resp, err := c.R().
		SetQueryParams(map[string]string{"page": fmt.Sprintf("%d", page), "perPage": fmt.Sprintf("%d", perPage)}).
		SetResult(&result).
		Get(fmt.Sprintf("%s/anime/trending", c.baseURL))

	if err != nil {
		return nil, 0, fmt.Errorf("anime-service call failed for trending: %w", err)
	}
	if !resp.IsSuccess() {
		return nil, 0, fmt.Errorf("anime-service error on trending (status %s): %s", resp.Status(), resp.String())
	}
	return result.Data, result.Meta.Total, nil
}

// GetAnimeBySeason from anime-service
// Note: anime-service takes /:year/:season in path, client was sending as query params. Correcting.
func (c *AnimeClient) GetAnimeBySeason(year int, season string, page, perPage int) ([]models.AnimeCache, int, error) {
	var result pagedAnimeCacheResult
	resp, err := c.R().
		SetQueryParams(map[string]string{"page": fmt.Sprintf("%d", page), "perPage": fmt.Sprintf("%d", perPage)}).
		SetResult(&result).
		Get(fmt.Sprintf("%s/anime/season/%d/%s", c.baseURL, year, season))

	if err != nil {
		return nil, 0, fmt.Errorf("anime-service call failed for season: %w", err)
	}
	if !resp.IsSuccess() {
		return nil, 0, fmt.Errorf("anime-service error on season (status %s): %s", resp.Status(), resp.String())
	}
	return result.Data, result.Meta.Total, nil
}

// GetAnimeRecommendations (placeholder, adapt if anime-service implements it properly)
// The anime-service's placeholder doesn't use userID.
func (c *AnimeClient) GetAnimeRecommendations(page, perPage int) ([]models.AnimeCache, int, error) {
	var result pagedAnimeCacheResult
	resp, err := c.R().
		SetQueryParams(map[string]string{"page": fmt.Sprintf("%d", page), "perPage": fmt.Sprintf("%d", perPage)}).
		SetResult(&result).
		Get(fmt.Sprintf("%s/anime/recommendations", c.baseURL))
	if err != nil {
		return nil, 0, fmt.Errorf("anime-service call failed for recommendations: %w", err)
	}
	if !resp.IsSuccess() {
		return nil, 0, fmt.Errorf("anime-service error on recommendations (status %s): %s", resp.Status(), resp.String())
	}
	return result.Data, result.Meta.Total, nil
}

// AddWatchProvider - This function might be for admin purposes.
// The user-service typically wouldn't directly tell anime-service to add a generic provider
// unless it's a "suggestion" feature. The endpoint in anime-service is POST /api/v1/anime/:animeId/providers
func (c *AnimeClient) AddWatchProviderToAnime(animeID int, providerData models.WatchProvider) (*models.WatchProvider, error) {
	var result models.WatchProvider // Assuming anime-service returns the created provider
	resp, err := c.R().
		SetBody(providerData).
		SetResult(&result).
		Post(fmt.Sprintf("%s/anime/%d/providers", c.baseURL, animeID))

	if err != nil {
		return nil, fmt.Errorf("anime-service call failed for adding provider: %w", err)
	}
	if resp.StatusCode() != 201 { // Assuming 201 Created
		return nil, fmt.Errorf("anime-service error adding provider (status %s): %s", resp.Status(), resp.String())
	}
	return &result, nil
}

// ExploreAnime calls the anime-service's explore endpoint
func (c *AnimeClient) ExploreAnime(tags []string, page, perPage int) ([]models.AnimeCache, int, error) {
	var result pagedAnimeCacheResult
	resp, err := c.R().
		SetQueryParams(map[string]string{
			"tags":    strings.Join(tags, ","), // anime-service expects comma-separated
			"page":    fmt.Sprintf("%d", page),
			"perPage": fmt.Sprintf("%d", perPage),
		}).
		SetResult(&result).
		Get(fmt.Sprintf("%s/anime/explore", c.baseURL))

	if err != nil {
		return nil, 0, fmt.Errorf("anime-service call failed for explore: %w", err)
	}
	if !resp.IsSuccess() {
		return nil, 0, fmt.Errorf("anime-service error on explore (status %s): %s", resp.Status(), resp.String())
	}
	return result.Data, result.Meta.Total, nil
}

// GetUpcomingAnime from anime-service
func (c *AnimeClient) GetUpcomingAnime(page, perPage int) ([]models.AnimeCache, int, error) {
	var result pagedAnimeCacheResult // Assuming pagedAnimeCacheResult is defined as before
	resp, err := c.R().
		SetQueryParams(map[string]string{"page": fmt.Sprintf("%d", page), "perPage": fmt.Sprintf("%d", perPage)}).
		SetResult(&result).
		Get(fmt.Sprintf("%s/anime/upcoming", c.baseURL)) // Ensure this path matches anime-service

	if err != nil {
		return nil, 0, fmt.Errorf("anime-service call failed for upcoming: %w", err)
	}
	if !resp.IsSuccess() {
		return nil, 0, fmt.Errorf("anime-service error on upcoming (status %s): %s", resp.Status(), resp.String())
	}
	return result.Data, result.Meta.Total, nil
}

// GetRecentlyReleasedAnime from anime-service
func (c *AnimeClient) GetRecentlyReleasedAnime(page, perPage int) ([]models.AnimeCache, int, error) {
	var result pagedAnimeCacheResult // Assuming pagedAnimeCacheResult is defined
	resp, err := c.R().
		SetQueryParams(map[string]string{"page": fmt.Sprintf("%d", page), "perPage": fmt.Sprintf("%d", perPage)}).
		SetResult(&result).
		Get(fmt.Sprintf("%s/anime/recently-released", c.baseURL)) // Ensure this path matches anime-service

	if err != nil {
		return nil, 0, fmt.Errorf("anime-service call failed for recently-released: %w", err)
	}
	if !resp.IsSuccess() {
		return nil, 0, fmt.Errorf("anime-service error on recently-released (status %s): %s", resp.Status(), resp.String())
	}
	return result.Data, result.Meta.Total, nil
}

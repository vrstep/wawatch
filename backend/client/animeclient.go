package client

import (
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/vrstep/wawatch-backend/models"
)

type AnimeClient struct {
	client  *resty.Client
	baseURL string
}

func NewAnimeClient(baseURL string) *AnimeClient {
	client := resty.New()
	client.SetTimeout(10 * time.Second)
	client.SetRetryCount(3)

	return &AnimeClient{
		client:  client,
		baseURL: baseURL,
	}
}

// ForwardRequestID forwards the request ID to downstream services
func (c *AnimeClient) ForwardRequestID(requestID string) *AnimeClient {
	c.client.SetHeader("X-Request-ID", requestID)
	return c
}

// GetAnimeDetails fetches anime details from the anime service
func (c *AnimeClient) GetAnimeDetails(animeID int) (*models.AnimeDetails, error) {
	var result struct {
		Anime     *models.AnimeDetails   `json:"anime"`
		Providers []models.WatchProvider `json:"providers"`
	}

	resp, err := c.client.R().
		SetResult(&result).
		Get(fmt.Sprintf("%s/anime/%d", c.baseURL, animeID))

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to get anime details: %s", resp.String())
	}

	return result.Anime, nil
}

func (c *AnimeClient) GetAnimeByID(id int) (*models.AnimeDetails, error) {
	var result struct {
		Data *models.AnimeDetails `json:"data"`
	}

	resp, err := c.client.R().
		SetResult(&result).
		Get(fmt.Sprintf("%s/anime/%d", c.baseURL, id))

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to get anime by ID: %s", resp.String())
	}

	return result.Data, nil
}

// SearchAnime searches for anime through the anime service
func (c *AnimeClient) SearchAnime(query string, page, perPage int) ([]models.AnimeCache, int, error) {
	var result struct {
		Data []models.AnimeCache `json:"data"`
		Meta struct {
			Total int `json:"total"`
		} `json:"meta"`
	}

	resp, err := c.client.R().
		SetQueryParams(map[string]string{
			"q":       query,
			"page":    fmt.Sprintf("%d", page),
			"perPage": fmt.Sprintf("%d", perPage),
		}).
		SetResult(&result).
		Get(fmt.Sprintf("%s/anime/search", c.baseURL))

	if err != nil {
		return nil, 0, err
	}

	if resp.StatusCode() != 200 {
		return nil, 0, fmt.Errorf("failed to search anime: %s", resp.String())
	}

	return result.Data, result.Meta.Total, nil
}

// Additional methods for other anime operations...

func (c *AnimeClient) GetPopularAnime(page, perPage int) ([]models.AnimeCache, int, error) {
	var result struct {
		Data []models.AnimeCache `json:"data"`
		Meta struct {
			Total int `json:"total"`
		} `json:"meta"`
	}

	resp, err := c.client.R().
		SetQueryParams(map[string]string{
			"page":    fmt.Sprintf("%d", page),
			"perPage": fmt.Sprintf("%d", perPage),
		}).
		SetResult(&result).
		Get(fmt.Sprintf("%s/anime/popular", c.baseURL))

	if err != nil {
		return nil, 0, err
	}

	if resp.StatusCode() != 200 {
		return nil, 0, fmt.Errorf("failed to get popular anime: %s", resp.String())
	}

	return result.Data, result.Meta.Total, nil
}
func (c *AnimeClient) GetTrendingAnime(page, perPage int) ([]models.AnimeCache, int, error) {
	var result struct {
		Data []models.AnimeCache `json:"data"`
		Meta struct {
			Total int `json:"total"`
		} `json:"meta"`
	}

	resp, err := c.client.R().
		SetQueryParams(map[string]string{
			"page":    fmt.Sprintf("%d", page),
			"perPage": fmt.Sprintf("%d", perPage),
		}).
		SetResult(&result).
		Get(fmt.Sprintf("%s/anime/trending", c.baseURL))

	if err != nil {
		return nil, 0, err
	}

	if resp.StatusCode() != 200 {
		return nil, 0, fmt.Errorf("failed to get trending anime: %s", resp.String())
	}

	return result.Data, result.Meta.Total, nil
}
func (c *AnimeClient) GetAnimeBySeason(year int, season string, page, perPage int) ([]models.AnimeCache, int, error) {
	var result struct {
		Data []models.AnimeCache `json:"data"`
		Meta struct {
			Total int `json:"total"`
		} `json:"meta"`
	}

	resp, err := c.client.R().
		SetQueryParams(map[string]string{
			"year":    fmt.Sprintf("%d", year),
			"season":  season,
			"page":    fmt.Sprintf("%d", page),
			"perPage": fmt.Sprintf("%d", perPage),
		}).
		SetResult(&result).
		Get(fmt.Sprintf("%s/anime/season", c.baseURL))

	if err != nil {
		return nil, 0, err
	}

	if resp.StatusCode() != 200 {
		return nil, 0, fmt.Errorf("failed to get anime by season: %s", resp.String())
	}

	return result.Data, result.Meta.Total, nil
}
func (c *AnimeClient) GetAnimeRecommendations(userID int) ([]models.AnimeCache, error) {
	var result struct {
		Data []models.AnimeCache `json:"data"`
	}

	resp, err := c.client.R().
		SetQueryParams(map[string]string{
			"userID": fmt.Sprintf("%d", userID),
		}).
		SetResult(&result).
		Get(fmt.Sprintf("%s/anime/recommendations", c.baseURL))

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to get anime recommendations: %s", resp.String())
	}

	return result.Data, nil
}

func (c *AnimeClient) AddWatchProvider(provider models.WatchProvider) error {
	resp, err := c.client.R().
		SetBody(provider).
		Post(fmt.Sprintf("%s/anime/providers", c.baseURL))

	if err != nil {
		return err
	}

	if resp.StatusCode() != 201 {
		return fmt.Errorf("failed to add watch provider: %s", resp.String())
	}

	return nil
}

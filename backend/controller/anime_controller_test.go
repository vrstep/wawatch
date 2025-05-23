// backend/controller/anime_controller_test.go (or mocks_test.go)
package controller_test // Or "package controller" if in the same package

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vrstep/wawatch-backend/client" // The interface
	"github.com/vrstep/wawatch-backend/config"
	"github.com/vrstep/wawatch-backend/controller"
	"github.com/vrstep/wawatch-backend/models" // Your models
)

type MockAnimeServiceClient struct {
	mock.Mock
}

func (m *MockAnimeServiceClient) WithRequestID(requestID string) client.AnimeServiceAPIClient {
	args := m.Called(requestID)
	// If WithRequestID in your actual client is used for chaining and returns the client itself,
	// the mock should also return itself to allow m.WithRequestID(...).SearchAnime(...).
	// If it's just for setting a header internally, it might not need to return anything specific in the mock's signature,
	// but to match the interface, it must return client.AnimeServiceAPIClient.
	if ret := args.Get(0); ret != nil {
		return ret.(client.AnimeServiceAPIClient)
	}
	return m // Return self for chaining
}

func (m *MockAnimeServiceClient) SearchAnime(query string, page int, perPage int) ([]models.AnimeCache, int, error) {
	args := m.Called(query, page, perPage)
	var resData []models.AnimeCache
	if args.Get(0) != nil {
		resData = args.Get(0).([]models.AnimeCache)
	}
	return resData, args.Int(1), args.Error(2)
}

func (m *MockAnimeServiceClient) GetAnimeDetailsAndProviders(animeID int) (*models.AnimeDetails, []models.WatchProvider, error) {
	args := m.Called(animeID)
	var ad *models.AnimeDetails
	var wp []models.WatchProvider
	if args.Get(0) != nil {
		ad = args.Get(0).(*models.AnimeDetails)
	}
	if args.Get(1) != nil {
		wp = args.Get(1).([]models.WatchProvider)
	}
	return ad, wp, args.Error(2)
}

func (m *MockAnimeServiceClient) GetAnimeByID(animeID int) (*models.AnimeDetails, error) {
	args := m.Called(animeID)
	var ad *models.AnimeDetails
	if args.Get(0) != nil {
		ad = args.Get(0).(*models.AnimeDetails)
	}
	return ad, args.Error(1)
}

func (m *MockAnimeServiceClient) GetPopularAnime(page int, perPage int) ([]models.AnimeCache, int, error) {
	args := m.Called(page, perPage)
	var resData []models.AnimeCache
	if args.Get(0) != nil {
		resData = args.Get(0).([]models.AnimeCache)
	}
	return resData, args.Int(1), args.Error(2)
}

func (m *MockAnimeServiceClient) GetTrendingAnime(page int, perPage int) ([]models.AnimeCache, int, error) {
	args := m.Called(page, perPage)
	var resData []models.AnimeCache
	if args.Get(0) != nil {
		resData = args.Get(0).([]models.AnimeCache)
	}
	return resData, args.Int(1), args.Error(2)
}

func (m *MockAnimeServiceClient) GetAnimeBySeason(year int, season string, page int, perPage int) ([]models.AnimeCache, int, error) {
	args := m.Called(year, season, page, perPage)
	var resData []models.AnimeCache
	if args.Get(0) != nil {
		resData = args.Get(0).([]models.AnimeCache)
	}
	return resData, args.Int(1), args.Error(2)
}

func (m *MockAnimeServiceClient) GetAnimeRecommendations(page int, perPage int) ([]models.AnimeCache, int, error) {
	args := m.Called(page, perPage)
	var resData []models.AnimeCache
	if args.Get(0) != nil {
		resData = args.Get(0).([]models.AnimeCache)
	}
	return resData, args.Int(1), args.Error(2)
}

// **** ADDED/COMPLETED MISSING METHODS ****
func (m *MockAnimeServiceClient) ExploreAnime(tags []string, page int, perPage int) ([]models.AnimeCache, int, error) {
	args := m.Called(tags, page, perPage)
	var resData []models.AnimeCache
	if args.Get(0) != nil {
		resData = args.Get(0).([]models.AnimeCache)
	}
	return resData, args.Int(1), args.Error(2)
}

func (m *MockAnimeServiceClient) GetUpcomingAnime(page int, perPage int) ([]models.AnimeCache, int, error) {
	args := m.Called(page, perPage)
	var resData []models.AnimeCache
	if args.Get(0) != nil {
		resData = args.Get(0).([]models.AnimeCache)
	}
	return resData, args.Int(1), args.Error(2)
}

func (m *MockAnimeServiceClient) GetRecentlyReleasedAnime(page int, perPage int) ([]models.AnimeCache, int, error) {
	args := m.Called(page, perPage)
	var resData []models.AnimeCache
	if args.Get(0) != nil {
		resData = args.Get(0).([]models.AnimeCache)
	}
	return resData, args.Int(1), args.Error(2)
}

func TestSearchAnimePassThrough_Success(t *testing.T) {
	_, token := createAndLoginTestUser(config.DB, "searchuser", "password") // Create a user for auth

	mockClient := new(MockAnimeServiceClient)
	controller.SetAnimeServiceClientForTest(mockClient)

	query := "Love Ru"
	expectedCaches := []models.AnimeCache{{ID: 123, Title: "To Love Ru"}}
	expectedTotal := 1
	page := 1
	perPage := 10

	mockClient.On("SearchAnime", query, 1, 20).Return(expectedCaches, expectedTotal, nil).Once()

	rr := performAuthRequest("GET", fmt.Sprintf("/ext/anime/search?q=%s&page=%d&perPage=%d", query, page, perPage), nil, token, testRouter) // USE performAuthRequest

	assert.Equal(t, http.StatusOK, rr.Code)
	// Assertions on response body based on what mockClient returns
	var responseBody struct {
		Data []models.AnimeCache    `json:"data"`
		Meta map[string]interface{} `json:"meta"`
	}
	json.Unmarshal(rr.Body.Bytes(), &responseBody)
	assert.Len(t, responseBody.Data, 1)
	assert.Equal(t, "To Love Ru", responseBody.Data[0].Title)
	assert.Equal(t, float64(1), responseBody.Meta["total"]) // JSON numbers are float64

	mockClient.AssertExpectations(t)
}

package controller_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin" // For the mock client and interface
	"github.com/stretchr/testify/assert"
	"github.com/vrstep/wawatch-backend/config"
	"github.com/vrstep/wawatch-backend/controller" // For SetAnimeServiceClientForTest
	"github.com/vrstep/wawatch-backend/models"
	"gorm.io/gorm"
)

// testRouter and testDB are assumed to be initialized by TestMain
// MockAnimeServiceClient struct is defined (as fixed in Part 1)

func TestAddAnimeToList_NewItem_Success(t *testing.T) {
	user, token := createAndLoginTestUser(config.DB, "listadd", "password")

	mockClient := new(MockAnimeServiceClient)
	controller.SetAnimeServiceClientForTest(mockClient) // Inject mock

	animeIDToAdd := 101
	expectedAnimeDetailsFromService := &models.AnimeDetails{
		ID: animeIDToAdd,
		Title: struct {
			Romaji  string `json:"romaji"`
			English string `json:"english"`
			Native  string `json:"native"`
		}{
			Romaji:  "Remote Anime",
			English: "Remote Anime",
			Native:  "",
		},
		Episodes: 12,
		Format:   "TV",
		CoverImage: struct {
			Large  string `json:"large"`
			Medium string `json:"medium"`
		}{
			Large:  "remote.jpg",
			Medium: "",
		},
	}

	// Expect GetAnimeByID to be called because item is not in local User-Service cache
	mockClient.On("GetAnimeByID", animeIDToAdd).Return(expectedAnimeDetailsFromService, nil).Once()

	payload := gin.H{"anime_id": animeIDToAdd, "status": models.Planned}
	rr := performAuthRequest("POST", "/api/v1/me/animelist/", payload, token, testRouter)

	assert.Equal(t, http.StatusCreated, rr.Code)
	var responseData struct{ Data models.UserAnimeList }
	json.Unmarshal(rr.Body.Bytes(), &responseData) // Assuming structure { "message": "...", "data": ... }
	assert.Equal(t, animeIDToAdd, responseData.Data.AnimeExternalID)
	assert.Equal(t, models.Planned, responseData.Data.Status)

	// Verify in DB
	var listEntry models.UserAnimeList
	config.DB.Where("user_id = ? AND anime_external_id = ?", user.ID, animeIDToAdd).First(&listEntry)
	assert.Equal(t, models.Planned, listEntry.Status)

	// Verify local AnimeCache in user-service DB
	var cacheEntry models.AnimeCache
	config.DB.First(&cacheEntry, animeIDToAdd)
	assert.Equal(t, expectedAnimeDetailsFromService.Title.English, cacheEntry.Title)

	mockClient.AssertExpectations(t)
}

func TestAddAnimeToList_UpdateExisting_Success(t *testing.T) {
	user, token := createAndLoginTestUser(config.DB, "listupdate", "password")
	animeIDToUpdate := 102

	// Pre-populate local AnimeCache for this test to skip anime-service call
	config.DB.Create(&models.AnimeCache{ID: animeIDToUpdate, Title: "Cached Anime", Format: "TV"})
	// Pre-populate UserAnimeList entry
	initialEntry := models.UserAnimeList{UserID: user.ID, AnimeExternalID: animeIDToUpdate, Status: models.Planned, Progress: 1}
	config.DB.Create(&initialEntry)

	// No mock setup needed for animeServiceClient.GetAnimeByID() as it should hit local cache

	payload := gin.H{"anime_id": animeIDToUpdate, "status": models.Watching, "progress": 5}
	rr := performAuthRequest("POST", "/api/v1/me/animelist/", payload, token, testRouter)

	assert.Equal(t, http.StatusOK, rr.Code) // Should be update
	var responseData struct{ Data models.UserAnimeList }
	json.Unmarshal(rr.Body.Bytes(), &responseData)
	assert.Equal(t, models.Watching, responseData.Data.Status)
	assert.Equal(t, 5, responseData.Data.Progress)

	var updatedEntry models.UserAnimeList
	config.DB.First(&updatedEntry, initialEntry.ID)
	assert.Equal(t, models.Watching, updatedEntry.Status)
	assert.Equal(t, 5, updatedEntry.Progress)
}

func TestUpdateAnimeInList_Success(t *testing.T) {
	user, token := createAndLoginTestUser(config.DB, "entryupdate", "password")
	entry := models.UserAnimeList{
		UserID:          user.ID,
		AnimeExternalID: 777,
		Status:          models.Watching,
		Progress:        5,
	}
	config.DB.Create(&entry) // entry.ID will be populated

	payload := gin.H{"status": models.Completed, "progress": 12, "score": 9}
	rr := performAuthRequest("PATCH", fmt.Sprintf("/api/v1/me/animelist/entry/%d", entry.ID), payload, token, testRouter)

	assert.Equal(t, http.StatusOK, rr.Code)
	var updatedEntry models.UserAnimeList
	config.DB.First(&updatedEntry, entry.ID)
	assert.Equal(t, models.Completed, updatedEntry.Status)
	assert.Equal(t, 12, updatedEntry.Progress)
	assert.NotNil(t, updatedEntry.Score)
	assert.Equal(t, 9, *updatedEntry.Score)
}

func TestRemoveAnimeFromList_Success(t *testing.T) {
	user, token := createAndLoginTestUser(config.DB, "entrydelete", "password")
	entry := models.UserAnimeList{UserID: user.ID, AnimeExternalID: 888, Status: models.Dropped}
	config.DB.Create(&entry)

	rr := performAuthRequest("DELETE", fmt.Sprintf("/api/v1/me/animelist/entry/%d", entry.ID), nil, token, testRouter)
	assert.Equal(t, http.StatusOK, rr.Code)

	var deletedEntry models.UserAnimeList
	err := config.DB.First(&deletedEntry, entry.ID).Error
	assert.Error(t, err) // Should be gorm.ErrRecordNotFound
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}

func TestShowUserAnimeList_Success(t *testing.T) {
	user, token := createAndLoginTestUser(config.DB, "showlist", "password")
	// Pre-populate some list items and their cache entries
	anime1ID, anime2ID := 901, 902
	config.DB.Create(&models.AnimeCache{ID: anime1ID, Title: "Anime One"})
	config.DB.Create(&models.AnimeCache{ID: anime2ID, Title: "Anime Two"})
	config.DB.Create(&models.UserAnimeList{UserID: user.ID, AnimeExternalID: anime1ID, Status: models.Watching})
	config.DB.Create(&models.UserAnimeList{UserID: user.ID, AnimeExternalID: anime2ID, Status: models.Planned})

	rr := performAuthRequest("GET", "/api/v1/me/animelist/", nil, token, testRouter)
	assert.Equal(t, http.StatusOK, rr.Code)

	var response struct {
		Data []controller.UserAnimeListResponse `json:"data"` // Use the response struct from controller
		Meta map[string]interface{}             `json:"meta"`
	}
	json.Unmarshal(rr.Body.Bytes(), &response)
	assert.Len(t, response.Data, 2)
	// Further assertions on content and meta
}

// Add tests for:
// - AddToAnimeList when anime-service fails to find anime
// - UpdateList with invalid entry ID, or entry not belonging to user
// - RemoveList with invalid entry ID, or entry not belonging to user
// - GetUserAnimeList with status filter, pagination

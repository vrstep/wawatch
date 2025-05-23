package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vrstep/wawatch-backend/config"
	"github.com/vrstep/wawatch-backend/models"
)

// UserViewHistoryResponse combines UserViewHistory with its local AnimeCache details
type UserViewHistoryResponse struct {
	models.UserViewHistory
	AnimeDetails *models.AnimeCache `json:"anime_details,omitempty"`
}

// GetUserViewHistory retrieves the authenticated user's view history, paginated.
func GetUserViewHistory(c *gin.Context) {
	currentUser := c.MustGet("user").(models.User)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("perPage", "30")) // Default to 30 history items
	offset := (page - 1) * perPage

	var historyEntries []models.UserViewHistory
	var totalEntries int64

	query := config.DB.Model(&models.UserViewHistory{}).Where("user_id = ?", currentUser.ID)

	// Count total entries for pagination meta
	if err := query.Count(&totalEntries).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count user view history"})
		return
	}

	// Fetch paginated UserViewHistory entries
	err := query.Order("last_viewed_at DESC").Limit(perPage).Offset(offset).Find(&historyEntries).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user view history"})
		return
	}

	responseItems := make([]UserViewHistoryResponse, 0, len(historyEntries))
	if len(historyEntries) > 0 {
		animeIDs := make([]int, len(historyEntries))
		for i, entry := range historyEntries {
			animeIDs[i] = entry.AnimeExternalID
		}

		// Fetch corresponding AnimeCache details from user-service's local cache
		var animeCaches []models.AnimeCache
		if dbErr := config.DB.Where("id IN ?", animeIDs).Find(&animeCaches).Error; dbErr != nil {
			// Log error but continue; some anime details might be missing from cache
			log.Printf("Error fetching some anime cache details for view history: %v", dbErr)
		}

		cacheMap := make(map[int]models.AnimeCache)
		for _, ac := range animeCaches {
			cacheMap[ac.ID] = ac
		}

		for _, entry := range historyEntries {
			respItem := UserViewHistoryResponse{UserViewHistory: entry}
			if cachedAnime, found := cacheMap[entry.AnimeExternalID]; found {
				respItem.AnimeDetails = &cachedAnime
			}
			responseItems = append(responseItems, respItem)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": responseItems,
		"meta": gin.H{
			"total":       totalEntries,
			"page":        page,
			"perPage":     perPage,
			"totalPages":  (totalEntries + int64(perPage) - 1) / int64(perPage),
			"hasNextPage": int64(page*perPage) < totalEntries,
		},
	})
}

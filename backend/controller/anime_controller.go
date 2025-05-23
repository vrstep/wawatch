package controller

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vrstep/wawatch-backend/client" // Client to call anime-service
	"github.com/vrstep/wawatch-backend/config"
	"github.com/vrstep/wawatch-backend/models"
	"gorm.io/gorm"
	// Models might not be directly needed here if client handles types,
	// but often used for request/response structuring if not fully handled by client.
)

// animeServiceClient is initialized once and reused.
var animeServiceClient *client.AnimeClient

func InitAnimeServiceClient() { // Call this from main.go after env load
	animeServiceClient = client.NewAnimeClient()
}

// Helper function to get animeServiceClient with RequestID
func getClientWithRequestID(c *gin.Context) *client.AnimeClient {
	if animeServiceClient == nil {
		// This should not happen if InitAnimeServiceClient is called at startup
		log.Println("Error: animeServiceClient not initialized!")
		InitAnimeServiceClient() // Fallback initialization
	}
	reqID, exists := c.Get("RequestID")
	if exists {
		if idStr, ok := reqID.(string); ok {
			return animeServiceClient.WithRequestID(idStr)
		}
	}
	return animeServiceClient // Return base client if no RequestID
}

// SearchAnime forwards search to anime-service
func SearchAnime(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("perPage", "20"))

	client := getClientWithRequestID(c)
	results, total, err := client.SearchAnime(query, page, perPage)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Failed to search anime via anime-service: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": results, "meta": gin.H{"total": total, "page": page, "perPage": perPage, "totalPages": (total + perPage - 1) / perPage, "hasNextPage": page*perPage < total}})
}

// GetAnimeDetails forwards request to anime-service
func GetAnimeDetails(c *gin.Context) {
	idParam := c.Param("id")
	animeID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid anime ID"})
		return
	}

	client := getClientWithRequestID(c)
	animeDetailsFromService, providers, err := client.GetAnimeDetailsAndProviders(animeID)
	if err != nil {
		if strings.Contains(err.Error(), "status 404") || strings.Contains(err.Error(), "no anime data") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Anime not found via anime-service: " + err.Error()})
		} else {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Failed to fetch anime details via anime-service: " + err.Error()})
		}
		return
	}

	// If successfully fetched details and user is authenticated, record/update view history
	userInterface, userExists := c.Get("user")
	if userExists && animeDetailsFromService != nil { // Ensure animeDetailsFromService is not nil
		currentUser := userInterface.(models.User)

		var historyEntry models.UserViewHistory
		dbErr := config.DB.Where("user_id = ? AND anime_external_id = ?", currentUser.ID, animeDetailsFromService.ID).First(&historyEntry).Error
		if dbErr != nil {
			if dbErr == gorm.ErrRecordNotFound {
				historyEntry = models.UserViewHistory{
					UserID:          currentUser.ID,
					AnimeExternalID: animeDetailsFromService.ID,
					LastViewedAt:    time.Now(),
					ViewCount:       1,
				}
				if createErr := config.DB.Create(&historyEntry).Error; createErr != nil {
					log.Printf("Error creating view history for user %d, anime %d: %v", currentUser.ID, animeDetailsFromService.ID, createErr)
				}
			} else {
				log.Printf("Error fetching existing view history for user %d, anime %d: %v", currentUser.ID, animeDetailsFromService.ID, dbErr)
			}
		} else {
			historyEntry.LastViewedAt = time.Now()
			historyEntry.ViewCount += 1
			if updateErr := config.DB.Save(&historyEntry).Error; updateErr != nil {
				log.Printf("Error updating view history for user %d, anime %d: %v", currentUser.ID, animeDetailsFromService.ID, updateErr)
			}
		}

		// Ensure this anime is in the local AnimeCache of user-service
		var localAnimeCache models.AnimeCache
		if cacheErr := config.DB.First(&localAnimeCache, animeDetailsFromService.ID).Error; cacheErr != nil && cacheErr == gorm.ErrRecordNotFound {
			// ---- CORRECTED SECTION ----
			// Directly call ToAnimeCache() if models.AnimeDetails (in user-service) has this method.
			// The variable animeDetailsFromService is already of type *models.AnimeDetails.
			convertedCacheEntry := animeDetailsFromService.ToAnimeCache() // Assumes ToAnimeCache() exists on *models.AnimeDetails or models.AnimeDetails

			if createCacheErr := config.DB.Create(&convertedCacheEntry).Error; createCacheErr != nil {
				log.Printf("Error saving anime %d to user-service local cache from GetAnimeDetails: %v", animeDetailsFromService.ID, createCacheErr)
			} else {
				log.Printf("Successfully cached anime %d in user-service local cache from GetAnimeDetails.", animeDetailsFromService.ID)
			}
			// ---- END CORRECTED SECTION ----
		}
	}

	c.JSON(http.StatusOK, gin.H{"anime": animeDetailsFromService, "providers": providers})
}

// GetPopularAnime forwards to anime-service
func GetPopularAnime(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("perPage", "20"))
	client := getClientWithRequestID(c)
	results, total, err := client.GetPopularAnime(page, perPage)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Failed to fetch popular anime: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": results, "meta": gin.H{"total": total, "page": page, "perPage": perPage, "totalPages": (total + perPage - 1) / perPage, "hasNextPage": page*perPage < total}})
}

// GetTrendingAnime forwards to anime-service
func GetTrendingAnime(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("perPage", "20"))
	client := getClientWithRequestID(c)
	results, total, err := client.GetTrendingAnime(page, perPage)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Failed to fetch trending anime: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": results, "meta": gin.H{"total": total, "page": page, "perPage": perPage, "totalPages": (total + perPage - 1) / perPage, "hasNextPage": page*perPage < total}})
}

// GetAnimeBySeason forwards to anime-service
func GetAnimeBySeason(c *gin.Context) {
	yearParam := c.Param("year")
	seasonParam := strings.ToUpper(c.Param("season")) // Match anime-service validation
	year, err := strconv.Atoi(yearParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year format"})
		return
	}
	validSeasons := map[string]bool{"WINTER": true, "SPRING": true, "SUMMER": true, "FALL": true}
	if !validSeasons[seasonParam] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid season. Use WINTER, SPRING, SUMMER, or FALL"})
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("perPage", "20"))

	client := getClientWithRequestID(c)
	results, total, err := client.GetAnimeBySeason(year, seasonParam, page, perPage)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Failed to fetch anime by season: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": results, "meta": gin.H{"total": total, "page": page, "perPage": perPage, "totalPages": (total + perPage - 1) / perPage, "hasNextPage": page*perPage < total}})
}

// GetAnimeRecommendations forwards to anime-service
func GetAnimeRecommendations(c *gin.Context) {
	// Assuming recommendations in anime-service don't currently depend on user ID
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("perPage", "10"))

	client := getClientWithRequestID(c)
	results, total, err := client.GetAnimeRecommendations(page, perPage) // Client method updated
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Failed to fetch recommendations: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": results, "meta": gin.H{"total": total, "page": page, "perPage": perPage, "totalPages": (total + perPage - 1) / perPage, "hasNextPage": page*perPage < total}})
}

// ExploreAnime forwards to anime-service
func ExploreAnime(c *gin.Context) {
	tagsQuery := c.Query("tags")
	if tagsQuery == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tags query parameter is required"})
		return
	}
	tags := strings.Split(tagsQuery, ",")
	for i, tag := range tags {
		tags[i] = strings.TrimSpace(tag)
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("perPage", "20"))

	client := getClientWithRequestID(c)
	results, total, err := client.ExploreAnime(tags, page, perPage)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Failed to explore anime: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": results, "meta": gin.H{"total": total, "page": page, "perPage": perPage, "totalPages": (total + perPage - 1) / perPage, "hasNextPage": page*perPage < total}})
}

// GetUpcomingAnime forwards to anime-service
func GetUpcomingAnime(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("perPage", "20"))

	client := getClientWithRequestID(c)
	results, total, err := client.GetUpcomingAnime(page, perPage)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Failed to fetch upcoming anime: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": results, "meta": gin.H{"total": total, "page": page, "perPage": perPage, "totalPages": (total + perPage - 1) / perPage, "hasNextPage": page*perPage < total}})
}

// GetRecentlyReleasedAnime forwards to anime-service
func GetRecentlyReleasedAnime(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("perPage", "20"))

	client := getClientWithRequestID(c)
	results, total, err := client.GetRecentlyReleasedAnime(page, perPage)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Failed to fetch recently released anime: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": results, "meta": gin.H{"total": total, "page": page, "perPage": perPage, "totalPages": (total + perPage - 1) / perPage, "hasNextPage": page*perPage < total}})
}

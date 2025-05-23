package controller

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vrstep/wawatch-backend/api"
	"github.com/vrstep/wawatch-backend/config"
	"github.com/vrstep/wawatch-backend/models"
)

var anilistClient api.AniListAPI

func SetAniListClient(client api.AniListAPI) {
	anilistClient = client
}

func init() {
	// Initialize with the real client by default when the package loads.
	// Ensure NewAniListClient() is accessible, or initialize it in main and pass it.
	// For simplicity, keeping init(), but dependency injection in main() is often preferred.
	SetAniListClient(api.NewAniListClient())
}

// SearchAnime handles anime search requests
func SearchAnime(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("perPage", "20"))

	results, total, err := anilistClient.SearchAnime(query, page, perPage)
	if err != nil {
		log.Printf("Error searching anime (query: %s): %v", query, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search anime"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": results,
		"meta": gin.H{"total": total, "page": page, "perPage": perPage, "totalPages": (total + perPage - 1) / perPage, "hasNextPage": page*perPage < total},
	})
}

// GetAnimeDetails fetches detailed information about an anime and its watch providers
func GetAnimeDetails(c *gin.Context) {
	idParam := c.Param("id")
	animeID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid anime ID format"})
		return
	}

	// Get detailed info from AniList
	animeDetails, err := anilistClient.GetAnimeByID(animeID)
	if err != nil {
		if strings.Contains(err.Error(), "no anime data returned") || strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Anime not found on AniList"})
		} else {
			log.Printf("Error fetching anime details from AniList (ID: %d): %v", animeID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch anime details from AniList"})
		}
		return
	}

	// Update or create cache entry in local DB
	cacheEntry := animeDetails.ToAnimeCache()
	if err := config.DB.Save(&cacheEntry).Error; err != nil {
		log.Printf("Warning: Failed to save anime (ID: %d) to cache: %v", animeID, err)
		// Continue even if cache save fails, priority is serving AniList data
	}

	// Get watch providers from local DB
	var providers []models.WatchProvider
	if err := config.DB.Where("anime_id = ?", animeID).Find(&providers).Error; err != nil {
		log.Printf("Error fetching watch providers for anime ID %d: %v", animeID, err)
		// Don't fail the request, just return empty providers
		providers = []models.WatchProvider{}
	}

	c.JSON(http.StatusOK, gin.H{
		"anime":     animeDetails,
		"providers": providers,
	})
}

// GetPopularAnime fetches popular anime from AniList
func GetPopularAnime(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("perPage", "20"))
	results, total, err := anilistClient.GetPopularAnime(page, perPage)
	if err != nil {
		log.Printf("Error fetching popular anime: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch popular anime"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": results, "meta": gin.H{"total": total, "page": page, "perPage": perPage, "totalPages": (total + perPage - 1) / perPage, "hasNextPage": page*perPage < total}})
}

// GetTrendingAnime fetches trending anime from AniList
func GetTrendingAnime(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("perPage", "20"))
	results, total, err := anilistClient.GetTrendingAnime(page, perPage)
	if err != nil {
		log.Printf("Error fetching trending anime: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch trending anime"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": results, "meta": gin.H{"total": total, "page": page, "perPage": perPage, "totalPages": (total + perPage - 1) / perPage, "hasNextPage": page*perPage < total}})
}

// GetAnimeRecommendations fetches recommendations (placeholder)
func GetAnimeRecommendations(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("perPage", "10"))
	results, total, err := anilistClient.GetPopularAnime(page, perPage) // Placeholder
	if err != nil {
		log.Printf("Error fetching recommendations (placeholder): %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recommendations"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": results, "meta": gin.H{"total": total, "page": page, "perPage": perPage, "totalPages": (total + perPage - 1) / perPage, "hasNextPage": page*perPage < total}})
}

// GetUpcomingAnime fetches upcoming anime from AniList
func GetUpcomingAnime(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("perPage", "20"))
	results, total, err := anilistClient.GetUpcomingAnime(page, perPage)
	if err != nil {
		log.Printf("Error fetching upcoming anime: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch upcoming anime"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": results, "meta": gin.H{"total": total, "page": page, "perPage": perPage, "totalPages": (total + perPage - 1) / perPage, "hasNextPage": page*perPage < total}})
}

// GetRecentlyReleasedAnime fetches recently released anime
func GetRecentlyReleasedAnime(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("perPage", "20"))
	results, total, err := anilistClient.GetRecentlyReleasedAnime(page, perPage)
	if err != nil {
		log.Printf("Error fetching recently released anime: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recently released anime"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": results, "meta": gin.H{"total": total, "page": page, "perPage": perPage, "totalPages": (total + perPage - 1) / perPage, "hasNextPage": page*perPage < total}})
}

// GetAnimeBySeason fetches anime by season
func GetAnimeBySeason(c *gin.Context) {
	yearParam := c.Param("year")
	seasonParam := strings.ToUpper(c.Param("season"))
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
	results, total, err := anilistClient.GetAnimeBySeason(year, seasonParam, page, perPage)
	if err != nil {
		log.Printf("Error fetching anime by season (Year: %d, Season: %s): %v", year, seasonParam, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch anime by season"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": results, "meta": gin.H{"total": total, "page": page, "perPage": perPage, "totalPages": (total + perPage - 1) / perPage, "hasNextPage": page*perPage < total}})
}

// ExploreAnime fetches anime by a comma-separated list of tags
func ExploreAnime(c *gin.Context) {
	tagsQuery := c.Query("tags")
	if tagsQuery == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tags query parameter is required"})
		return
	}
	tags := strings.Split(tagsQuery, ",")
	for i, tag := range tags { // Trim whitespace from tags
		tags[i] = strings.TrimSpace(tag)
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("perPage", "20"))

	results, total, err := anilistClient.GetAnimeByTags(tags, page, perPage)
	if err != nil {
		log.Printf("Error fetching anime by tags (%v): %v", tags, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch anime by tags"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": results, "meta": gin.H{"total": total, "page": page, "perPage": perPage, "totalPages": (total + perPage - 1) / perPage, "hasNextPage": page*perPage < total}})
}

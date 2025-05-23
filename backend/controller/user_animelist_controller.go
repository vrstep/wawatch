package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vrstep/wawatch-backend/config" // User-service's DB
	"github.com/vrstep/wawatch-backend/models"
	"gorm.io/gorm"
	// AnimeServiceClient is now initialized in anime_controller.go or main
)

// AddToAnimeList adds or updates an anime in the user's list
func AddToAnimeList(c *gin.Context) {
	currentUser := c.MustGet("user").(models.User) // From RequireAuth middleware

	var input struct {
		AnimeID      int        `json:"anime_id" binding:"required"` // This is AnimeExternalID
		Status       string     `json:"status" binding:"required"`
		Score        *int       `json:"score,omitempty"`
		Progress     *int       `json:"progress,omitempty"` // Use pointer for optional 0
		StartDate    *time.Time `json:"start_date,omitempty"`
		EndDate      *time.Time `json:"end_date,omitempty"`
		Notes        *string    `json:"notes,omitempty"`
		RewatchCount *int       `json:"rewatch_count,omitempty"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	validStatuses := map[string]bool{
		models.Watching: true, models.Completed: true, models.Planned: true,
		models.Dropped: true, models.Paused: true, models.Rewatching: true,
	}
	if !validStatuses[input.Status] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status: " + input.Status})
		return
	}

	// Check if anime exists in local User-Service cache. If not, fetch from Anime-Service and cache locally.
	var localAnimeCache models.AnimeCache
	err := config.DB.First(&localAnimeCache, input.AnimeID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("AnimeID %d not in user-service cache. Fetching from anime-service.", input.AnimeID)
			client := getClientWithRequestID(c)                                // Get client with RequestID
			remoteAnimeDetails, fetchErr := client.GetAnimeByID(input.AnimeID) // Uses the client
			if fetchErr != nil || remoteAnimeDetails == nil {
				log.Printf("Failed to fetch anime %d from anime-service: %v", input.AnimeID, fetchErr)
				c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Anime with ID %d not found via anime-service", input.AnimeID)})
				return
			}
			// Convert AnimeDetails from anime-service to local AnimeCache model and save
			// This assumes AnimeDetails.ToAnimeCache() exists and returns user-service's models.AnimeCache
			// If AnimeDetails is also a model in user-service, the conversion is direct.
			convertedCacheEntry := remoteAnimeDetails.ToAnimeCache() // Ensure this method exists on user-service's AnimeDetails model

			if createErr := config.DB.Create(&convertedCacheEntry).Error; createErr != nil {
				log.Printf("Failed to save anime %d to user-service cache: %v", input.AnimeID, createErr)
				// Continue, but this means list items might lack details later if this fails.
			}
			localAnimeCache = convertedCacheEntry
		} else {
			log.Printf("DB error checking user-service anime cache for %d: %v", input.AnimeID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
	}

	// Prepare UserAnimeList entry
	listEntry := models.UserAnimeList{
		UserID:          currentUser.ID,
		AnimeExternalID: input.AnimeID,
		Status:          input.Status,
	}
	// Apply optional fields from input
	if input.Score != nil {
		listEntry.Score = input.Score
	}
	if input.Progress != nil {
		listEntry.Progress = *input.Progress
	} else {
		listEntry.Progress = 0
	}
	if input.StartDate != nil {
		listEntry.StartDate = input.StartDate
	}
	if input.EndDate != nil {
		listEntry.EndDate = input.EndDate
	}
	if input.Notes != nil {
		listEntry.Notes = *input.Notes
	}
	if input.RewatchCount != nil {
		listEntry.RewatchCount = *input.RewatchCount
	} else {
		listEntry.RewatchCount = 0
	}

	// Upsert logic: Find existing or create new
	var existingEntry models.UserAnimeList
	tx := config.DB.Where("user_id = ? AND anime_external_id = ?", currentUser.ID, input.AnimeID).First(&existingEntry)

	if tx.Error == nil { // Entry exists, update it
		existingEntry.Status = listEntry.Status
		existingEntry.Score = listEntry.Score
		existingEntry.Progress = listEntry.Progress
		existingEntry.StartDate = listEntry.StartDate
		existingEntry.EndDate = listEntry.EndDate
		existingEntry.Notes = listEntry.Notes
		existingEntry.RewatchCount = listEntry.RewatchCount
		if errSave := config.DB.Save(&existingEntry).Error; errSave != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update list entry", "details": errSave.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "List entry updated", "data": existingEntry})
	} else if tx.Error == gorm.ErrRecordNotFound { // Entry does not exist, create it
		if errCreate := config.DB.Create(&listEntry).Error; errCreate != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add to list", "details": errCreate.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "Anime added to list", "data": listEntry})
	} else { // Other DB error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error finding list entry", "details": tx.Error.Error()})
	}
}

// UpdateListEntry updates specific fields of an entry in the user's anime list
func UpdateListEntry(c *gin.Context) {
	currentUser := c.MustGet("user").(models.User)
	entryDBIDStr := c.Param("id") // This is the DB ID of the UserAnimeList record
	entryDBID, err := strconv.Atoi(entryDBIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid list entry ID format"})
		return
	}

	var entry models.UserAnimeList // This will hold the state of the entry *before* updates
	// 1. Fetch the entry to ensure it exists and belongs to the user
	if err := config.DB.Where("id = ? AND user_id = ?", entryDBID, currentUser.ID).First(&entry).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "List entry not found or access denied"})
		} else {
			log.Printf("Error fetching list entry ID %d for user %d: %v", entryDBID, currentUser.ID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error fetching list entry"})
		}
		return
	}

	// This 'entry' variable now holds the current state from the DB.
	// We will compare the input to this and build an updateMap.

	var input struct {
		Status       *string    `json:"status,omitempty"`
		Score        *int       `json:"score,omitempty"`
		Progress     *int       `json:"progress,omitempty"`
		StartDate    *time.Time `json:"start_date,omitempty"`
		EndDate      *time.Time `json:"end_date,omitempty"`
		Notes        *string    `json:"notes,omitempty"`
		RewatchCount *int       `json:"rewatch_count,omitempty"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	log.Printf("UpdateListEntry: Input received for entry ID %d: %+v", entryDBID, input)

	updateMap := make(map[string]interface{})
	changed := false // Flag to see if any actual changes are requested

	// Status
	if input.Status != nil && *input.Status != entry.Status { // Only update if different
		validStatuses := map[string]bool{
			models.Watching: true, models.Completed: true, models.Planned: true,
			models.Dropped: true, models.Paused: true, models.Rewatching: true,
		}
		if !validStatuses[*input.Status] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status provided: " + *input.Status})
			return
		}
		updateMap["status"] = *input.Status
		changed = true
	}
	currentStatus := entry.Status // Default to existing status
	if val, ok := updateMap["status"].(string); ok {
		currentStatus = val // Use new status if it's being updated
	}

	// Score
	if input.Score != nil { // If score is present in input
		if entry.Score == nil || *input.Score != *entry.Score { // Update if different or original was nil
			updateMap["score"] = input.Score // Assign pointer directly to handle NULL
			changed = true
		}
	} else if c.Request.ContentLength > 0 && strings.Contains(string(c.ContentType()), "application/json") {
		// Check if "score": null was explicitly sent to clear it
		// This requires more complex JSON parsing or a specific flag.
		// For now, if input.Score is nil, it means "don't update score" or "score was null".
		// If you want to clear a score, send "score": null. The current logic handles that by
		// `updateMap["score"] = input.Score` which would be `updateMap["score"] = nil`.
	}

	// Progress
	if input.Progress != nil && *input.Progress != entry.Progress {
		updateMap["progress"] = *input.Progress
		changed = true
	}
	currentProgress := entry.Progress
	if val, ok := updateMap["progress"].(int); ok {
		currentProgress = val
	}

	// StartDate
	if input.StartDate != nil {
		if entry.StartDate == nil || !(*input.StartDate).Equal(*entry.StartDate) {
			updateMap["start_date"] = input.StartDate
			changed = true
		}
	}

	// EndDate Logic
	if input.EndDate != nil { // User explicitly sets EndDate (could be value or null)
		if entry.EndDate == nil || (input.EndDate != nil && !(*input.EndDate).Equal(*entry.EndDate)) || (input.EndDate == nil && entry.EndDate != nil) {
			updateMap["end_date"] = input.EndDate // Assign pointer directly to handle NULL
			changed = true
		}
	} else if currentStatus == models.Completed && entry.EndDate == nil { // Auto-set EndDate if status is now Completed and EndDate was not set AND not provided
		var animeCache models.AnimeCache
		if config.DB.First(&animeCache, entry.AnimeExternalID).Error == nil && animeCache.TotalEpisodes != nil {
			if currentProgress >= *animeCache.TotalEpisodes {
				now := time.Now()
				updateMap["end_date"] = &now
				changed = true
			}
		}
	} else if currentStatus != models.Completed && entry.EndDate != nil { // If status is no longer completed, and EndDate was set
		updateMap["end_date"] = gorm.Expr("NULL") // Clear it
		changed = true
	}

	// Notes
	if input.Notes != nil && *input.Notes != entry.Notes {
		updateMap["notes"] = *input.Notes
		changed = true
	}

	// RewatchCount
	if input.RewatchCount != nil && *input.RewatchCount != entry.RewatchCount {
		updateMap["rewatch_count"] = *input.RewatchCount
		changed = true
	}

	log.Printf("UpdateListEntry: Assembled updateMap for entry ID %d: %+v. Changed: %t", entryDBID, updateMap, changed)

	if !changed && len(updateMap) == 0 { // If no actual differing values were provided to update
		log.Printf("UpdateListEntry: No actual changes to apply for entry ID %d. Returning current data.", entryDBID)
		// Re-fetch to be safe, though 'entry' should be current state if no changes.
		var currentDBEntry models.UserAnimeList
		config.DB.First(&currentDBEntry, entry.ID)
		c.JSON(http.StatusOK, gin.H{"message": "No changes applied to the entry.", "data": currentDBEntry})
		return
	}

	// Only proceed with DB update if there are fields in updateMap
	if len(updateMap) > 0 {
		dbResult := config.DB.Model(&entry).Updates(updateMap)
		if dbResult.Error != nil {
			log.Printf("UpdateListEntry: DB.Updates failed for entry ID %d. Error: %v", entryDBID, dbResult.Error)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update entry", "details": dbResult.Error.Error()})
			return
		}
		log.Printf("UpdateListEntry: DB.Updates successful for entry ID %d. Rows affected: %d", entryDBID, dbResult.RowsAffected)

		if dbResult.RowsAffected == 0 && changed {
			// This is a strange case: we determined a change was needed, but DB said 0 rows affected.
			// Could happen if the record was deleted between the .First() and .Updates() call (race condition).
			log.Printf("UpdateListEntry: WARNING - Determined changes were needed, but DB reported 0 rows affected for entry ID %d.", entryDBID)
		}
	} else if changed {
		log.Printf("UpdateListEntry: WARNING - 'changed' is true but updateMap is empty for entry ID %d. This is unexpected.", entryDBID)
	}

	// Re-fetch the entry to return the fully updated record regardless of RowsAffected,
	// as GORM might not update the in-memory 'entry' struct directly with .Updates(map)
	var updatedEntry models.UserAnimeList
	if err := config.DB.First(&updatedEntry, entry.ID).Error; err != nil {
		log.Printf("UpdateListEntry: Failed to retrieve updated entry ID %d post-update attempt: %v", entry.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated entry state"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Entry update processed", "data": updatedEntry})
}

// UserAnimeListResponse combines UserAnimeList with its local AnimeCache details
type UserAnimeListResponse struct {
	models.UserAnimeList                    // Embed UserAnimeList fields
	AnimeDetails         *models.AnimeCache `json:"anime_details,omitempty"`
}

// GetUserAnimeList retrieves a user's anime list, paginated
func GetUserAnimeList(c *gin.Context) {
	currentUser := c.MustGet("user").(models.User)
	statusFilter := c.Query("status")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("perPage", "20"))
	offset := (page - 1) * perPage

	var userListItems []models.UserAnimeList
	var totalItems int64

	query := config.DB.Model(&models.UserAnimeList{}).Where("user_id = ?", currentUser.ID)
	if statusFilter != "" {
		query = query.Where("status = ?", statusFilter)
	}

	query.Count(&totalItems) // Get total count for pagination
	err := query.Order("updated_at DESC").Limit(perPage).Offset(offset).Find(&userListItems).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user anime list"})
		return
	}

	responseItems := make([]UserAnimeListResponse, 0, len(userListItems))
	if len(userListItems) > 0 {
		animeIDs := make([]int, len(userListItems))
		for i, item := range userListItems {
			animeIDs[i] = item.AnimeExternalID
		}

		var animeCaches []models.AnimeCache
		config.DB.Where("id IN ?", animeIDs).Find(&animeCaches) // Fetch from user-service's local cache

		cacheMap := make(map[int]models.AnimeCache)
		for _, ac := range animeCaches {
			cacheMap[ac.ID] = ac
		}

		for _, item := range userListItems {
			respItem := UserAnimeListResponse{UserAnimeList: item}
			if cachedAnime, found := cacheMap[item.AnimeExternalID]; found {
				respItem.AnimeDetails = &cachedAnime
			}
			responseItems = append(responseItems, respItem)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": responseItems,
		"meta": gin.H{
			"total":       totalItems,
			"page":        page,
			"perPage":     perPage,
			"totalPages":  (totalItems + int64(perPage) - 1) / int64(perPage),
			"hasNextPage": int64(page*perPage) < totalItems,
		},
	})
}

// DeleteListEntry removes an anime from the user's list by UserAnimeList DB ID
func DeleteListEntry(c *gin.Context) {
	currentUser := c.MustGet("user").(models.User)
	entryDBIDStr := c.Param("id") // DB ID of UserAnimeList record
	entryDBID, err := strconv.Atoi(entryDBIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid list entry ID"})
		return
	}

	result := config.DB.Where("id = ? AND user_id = ?", entryDBID, currentUser.ID).Delete(&models.UserAnimeList{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete entry", "details": result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Entry not found or not authorized to delete"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Entry deleted successfully"})
}

// GetUserAnimeListStats calculates and returns statistics for the user's anime list
func GetUserAnimeListStats(c *gin.Context) {
	currentUser := c.MustGet("user").(models.User)
	var list []models.UserAnimeList
	if err := config.DB.Where("user_id = ?", currentUser.ID).Find(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve anime list for stats"})
		return
	}

	stats := gin.H{
		"total_anime":      len(list),
		"episodes_watched": 0,
		"mean_score":       0.0,
		"status_counts": map[string]int{
			models.Watching: 0, models.Completed: 0, models.Planned: 0,
			models.Dropped: 0, models.Paused: 0, models.Rewatching: 0,
		},
	}
	totalScore, scoredCount, totalEpisodesWatched := 0, 0, 0

	for _, item := range list {
		stats["status_counts"].(map[string]int)[item.Status]++
		totalEpisodesWatched += item.Progress
		if item.Score != nil && *item.Score > 0 {
			totalScore += *item.Score
			scoredCount++
		}
	}
	stats["episodes_watched"] = totalEpisodesWatched
	if scoredCount > 0 {
		stats["mean_score"] = float64(totalScore) / float64(scoredCount)
	}
	c.JSON(http.StatusOK, stats)
}

// GetAnimeInUserList checks if a specific anime (by external ID) is in the user's list
func GetAnimeInUserList(c *gin.Context) {
	currentUser := c.MustGet("user").(models.User)
	animeExternalIDStr := c.Param("animeExternalID") // Changed from "id" to be specific
	animeExternalID, err := strconv.Atoi(animeExternalIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid anime external ID"})
		return
	}

	var entry models.UserAnimeList
	err = config.DB.Where("user_id = ? AND anime_external_id = ?", currentUser.ID, animeExternalID).First(&entry).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"in_list": false})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"in_list": true, "status": entry.Status, "progress": entry.Progress, "score": entry.Score, "list_entry_id": entry.ID})
}

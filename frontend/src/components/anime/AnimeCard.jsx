// src/components/anime/AnimeCard.js

import { useAuth } from "../../contexts/AuthContext";
import { useState } from "react";
import Button from "../ui/Button";
import LoadingSpinner from "../ui/LoadingSpinner";
import apiService from "../../api/api";
import { ChevronDown } from "lucide-react";

const AnimeCard = ({ anime, onStatusChange }) => {
  const { user } = useAuth();
  const [showStatusDropdown, setShowStatusDropdown] = useState(false);
  // Initialize currentStatus from anime.user_status if your API sends that for items already on a list
  // For now, assuming it might come as anime.userStatus (you may need to adjust if the field name from API is different for list status)
  const [currentStatus, setCurrentStatus] = useState(anime.userStatus || null);
  const [isLoadingStatus, setIsLoadingStatus] = useState(false);

  const handleStatusUpdate = async (newStatus) => {
    if (!user) return;
    setIsLoadingStatus(true);
    try {
      await apiService.addOrUpdateAnimeInList({
        anime_id: anime.id, // This uses the correct lowercase 'id' from your API
        status: newStatus,
      });
      setCurrentStatus(newStatus === "REMOVING" ? null : newStatus); // If removing, set status to null
      if (onStatusChange) onStatusChange(anime.id, newStatus);
    } catch (error) {
      console.error("Failed to update status for anime_id:", anime.id, error);
      // TODO: Show error to user via a toast or message
    } finally {
      setIsLoadingStatus(false);
      setShowStatusDropdown(false);
    }
  };

  // --- Corrected Title ---
  // API provides title as a direct string
  const displayTitle = anime.title || "Untitled";

  // --- Corrected Cover Image ---
  // API provides cover_image as a direct string URL
  const imageUrl =
    anime.cover_image || // Use the direct URL from API
    `https://placehold.co/300x450/1f2937/ffffff?text=${encodeURIComponent(
      displayTitle.substring(0, 15)
    )}`;

  return (
    <div className="bg-slate-200 dark:bg-slate-800 rounded-lg shadow-lg overflow-hidden transform transition-all hover:scale-105 group">
      <a href={`#/anime/${anime.id}`} className="block">
        <img
          src={imageUrl}
          alt={displayTitle}
          className="w-full h-72 object-cover"
          onError={(e) => {
            e.target.onerror = null; // prevent infinite loop if fallback also fails
            e.target.src = `https://placehold.co/300x450/1f2937/ffffff?text=${encodeURIComponent(
              displayTitle.substring(0, 15)
            )}`;
          }}
        />
      </a>
      <div className="p-4">
        <h3
          className="text-lg font-semibold text-slate-800 dark:text-slate-100 mb-1 truncate group-hover:whitespace-normal group-hover:overflow-visible"
          title={displayTitle}
        >
          <a href={`#/anime/${anime.id}`}>{displayTitle}</a>
        </h3>
        <p className="text-xs text-slate-600 dark:text-slate-400 mb-2">
          {/* --- Format is correct --- */}
          {anime.format || "N/A"}
          {/* --- Anime Airing Status: Assuming API might not provide it, or under a different name --- */}
          {/* If your API provides anime's airing status (e.g., as anime.airing_status), use that here instead of anime.status */}
          &bull; {anime.status || "N/A"}{" "}
          {/* If API provides this, use it. Otherwise, it will be N/A */}
          {/* --- Corrected Episodes --- */}
          {anime.total_episodes && ` &bull; ${anime.total_episodes} eps`}
        </p>

        {user && (
          <div className="relative mt-2">
            <Button
              variant={currentStatus ? "primary" : "secondary"}
              size="sm"
              className="w-full"
              onClick={() => setShowStatusDropdown(!showStatusDropdown)}
              disabled={isLoadingStatus}
            >
              {isLoadingStatus ? (
                <LoadingSpinner size={16} />
              ) : currentStatus ? (
                currentStatus.charAt(0) + currentStatus.slice(1).toLowerCase() // Display nicely
              ) : (
                "Add to List"
              )}
              <ChevronDown
                className={`ml-2 h-4 w-4 transition-transform ${
                  showStatusDropdown ? "rotate-180" : ""
                }`}
              />
            </Button>
            {showStatusDropdown && (
              <div className="absolute bottom-full mb-1 w-full bg-white dark:bg-slate-700 rounded-md shadow-lg z-10 border border-slate-300 dark:border-slate-600">
                {["WATCHING", "PLANNED", "COMPLETED", "PAUSED", "DROPPED"].map(
                  (statusOption) => (
                    <button
                      key={statusOption}
                      onClick={() => handleStatusUpdate(statusOption)}
                      className="block w-full text-left px-4 py-2 text-sm text-slate-700 dark:text-slate-200 hover:bg-slate-100 dark:hover:bg-slate-600"
                    >
                      {statusOption.charAt(0) +
                        statusOption.slice(1).toLowerCase()}
                    </button>
                  )
                )}
                {currentStatus && (
                  <button
                    onClick={() => handleStatusUpdate("REMOVING")} // Assuming your backend handles "REMOVING" to delete the entry
                    className="block w-full text-left px-4 py-2 text-sm text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900"
                  >
                    Remove from list
                  </button>
                )}
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  );
};

export default AnimeCard;

import React, { useState } from "react";
// import { useAuth } from "../contexts/AuthContext";
import ErrorMessage from "../components/ui/ErrorMessage";
import LoadingSpinner from "../components/ui/LoadingSpinner";
import NotFoundPage from "./NotFoundPage";
import apiService from "../api/api"; // Adjust the import path as necessary
import cn from "../utils/cn"; // Adjust path if necessary: from ui/ to utils/ is ../../

import { useEffect } from "react";

const AnimeDetailPage = ({ id }) => {
  const [anime, setAnime] = useState(null);
  const [providers, setProviders] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [activeTab, setActiveTab] = useState("info");
  // const { user } = useAuth();

  useEffect(() => {
    if (!id) return;
    const fetchDetails = async () => {
      try {
        setLoading(true);
        setError(null);
        // The getAnimeDetails from apiService now calls the /ext/anime/:id endpoint
        // which in user-service calls anime-service and records history.
        const data = await apiService.getAnimeDetails(id);
        setAnime(data.anime);
        setProviders(data.providers || []);
      } catch (err) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };
    fetchDetails();
  }, [id]);

  if (loading)
    return (
      <div className="flex justify-center items-center h-screen">
        <LoadingSpinner size={48} />
      </div>
    );
  if (error) return <ErrorMessage message={error} />;
  if (!anime) return <NotFoundPage message="Anime details not found." />;

  const title = anime.title.english || anime.title.romaji || anime.title.native;

  return (
    <div className="container mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div className="md:flex md:space-x-8">
        <div className="md:w-1/3 lg:w-1/4 mb-6 md:mb-0">
          <img
            src={
              anime.coverImage?.large ||
              `https://placehold.co/400x600/1f2937/ffffff?text=${encodeURIComponent(
                title.substring(0, 10)
              )}`
            }
            alt={title}
            className="rounded-lg shadow-xl w-full"
            onError={(e) => {
              e.target.onerror = null;
              e.target.src = `https://placehold.co/400x600/1f2937/ffffff?text=${encodeURIComponent(
                title.substring(0, 10)
              )}`;
            }}
          />
          {/* Add to List / Status Button Here if desired */}
        </div>
        <div className="md:w-2/3 lg:w-3/4">
          <h1 className="text-3xl md:text-4xl font-bold text-slate-800 dark:text-slate-100 mb-2">
            {title}
          </h1>
          {anime.title.english && anime.title.english !== title && (
            <p className="text-lg text-slate-600 dark:text-slate-400 mb-1">
              {anime.title.romaji}
            </p>
          )}
          {anime.title.native && (
            <p className="text-md text-slate-500 dark:text-slate-500 mb-4">
              {anime.title.native}
            </p>
          )}

          <div className="flex space-x-2 mb-4">
            {anime.genres?.map((genre) => (
              <span
                key={genre}
                className="px-2 py-0.5 bg-indigo-100 text-indigo-700 dark:bg-indigo-900 dark:text-indigo-300 text-xs font-semibold rounded-full"
              >
                {genre}
              </span>
            ))}
          </div>

          <div className="mb-6 border-b border-slate-300 dark:border-slate-700">
            <nav className="-mb-px flex space-x-6" aria-label="Tabs">
              {["info", "streams", "related", "characters"].map((tabName) => (
                <button
                  key={tabName}
                  onClick={() => setActiveTab(tabName)}
                  className={cn(
                    "whitespace-nowrap py-3 px-1 border-b-2 font-medium text-sm capitalize",
                    activeTab === tabName
                      ? "border-indigo-500 text-indigo-600 dark:text-indigo-400"
                      : "border-transparent text-slate-500 hover:text-slate-700 hover:border-slate-300 dark:text-slate-400 dark:hover:text-slate-200 dark:hover:border-slate-600"
                  )}
                >
                  {tabName}
                </button>
              ))}
            </nav>
          </div>

          {activeTab === "info" && (
            <div className="space-y-4 text-slate-700 dark:text-slate-300">
              <p
                className="leading-relaxed"
                dangerouslySetInnerHTML={{
                  __html: anime.description || "No description available.",
                }}
              ></p>
              <div className="grid grid-cols-2 gap-4 text-sm">
                <div>
                  <strong>Format:</strong> {anime.format}
                </div>
                <div>
                  <strong>Episodes:</strong> {anime.episodes || "N/A"}
                </div>
                <div>
                  <strong>Duration:</strong>{" "}
                  {anime.duration ? `${anime.duration} min` : "N/A"}
                </div>
                <div>
                  <strong>Status:</strong> {anime.status}
                </div>
                <div>
                  <strong>Season:</strong>{" "}
                  {anime.season ? `${anime.season} ${anime.seasonYear}` : "N/A"}
                </div>
                <div>
                  <strong>Avg Score:</strong>{" "}
                  {anime.averageScore ? `${anime.averageScore}/100` : "N/A"}
                </div>
                <div>
                  <strong>Popularity:</strong> {anime.popularity || "N/A"}
                </div>
                <div>
                  <strong>Studios:</strong>{" "}
                  {anime.studios?.nodes?.map((s) => s.name).join(", ") || "N/A"}
                </div>
              </div>
            </div>
          )}
          {activeTab === "streams" && (
            <div>
              <h3 className="text-xl font-semibold mb-3 text-slate-800 dark:text-slate-100">
                Streaming Platforms
              </h3>
              {providers.length > 0 ? (
                <ul className="space-y-2">
                  {providers.map((p) => (
                    <li
                      key={p.ID}
                      className="p-3 bg-slate-100 dark:bg-slate-800 rounded-md"
                    >
                      <a
                        href={p.url}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="font-medium text-indigo-600 dark:text-indigo-400 hover:underline"
                      >
                        {p.name}
                      </a>{" "}
                      ({p.type}, {p.region || "Global"})
                    </li>
                  ))}
                </ul>
              ) : (
                <p className="text-slate-600 dark:text-slate-400">
                  No streaming providers listed for this anime yet.
                </p>
              )}
            </div>
          )}
          {/* Placeholder for other tabs */}
          {activeTab === "related" && (
            <p>Related anime section (Not Implemented)</p>
          )}
          {activeTab === "characters" && (
            <p>Characters section (Not Implemented)</p>
          )}
        </div>
      </div>
    </div>
  );
};

export default AnimeDetailPage;

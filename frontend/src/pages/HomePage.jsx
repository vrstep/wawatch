import { useState, useEffect } from "react";
import ErrorMessage from "../components/ui/ErrorMessage";
import LoadingSpinner from "../components/ui/LoadingSpinner";
import apiService from "../api/api";
import SectionTitle from "../components/layout/SectionTitle";
import { TrendingUp } from "lucide-react";
import { Star } from "lucide-react";
import { Clock } from "lucide-react";
import { Calendar } from "lucide-react";
import AnimeGrid from "../components/anime/AnimeGrid";

const HomePage = () => {
  const [popular, setPopular] = useState([]);
  const [trending, setTrending] = useState([]);
  const [upcoming, setUpcoming] = useState([]);
  const [recentlyReleased, setRecentlyReleased] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        setLoading(true);
        setError(null);
        const [popRes, trendRes, upRes, recentRes] = await Promise.all([
          apiService.getPopularAnime(1, 12),
          apiService.getTrendingAnime(1, 12),
          apiService.getUpcomingAnime(1, 12),
          apiService.getRecentlyReleasedAnime(1, 12),
        ]);

        // --- 👇 Add these console logs 👇 ---
        console.log("Popular Anime API Response:", popRes);
        console.log("Trending Anime API Response:", trendRes);
        console.log("Upcoming Anime API Response:", upRes);
        console.log("Recently Released API Response:", recentRes);
        // --- End of console logs ---

        // Original state setting (we might need to change this)
        setPopular(popRes.data);
        setTrending(trendRes.data);
        setUpcoming(upRes.data);
        setRecentlyReleased(recentRes.data);
      } catch (err) {
        setError(err.message);
        console.error("Full error object in HomePage:", err); // Log the full error object
      } finally {
        setLoading(false);
      }
    };
    fetchData();
  }, []);

  if (loading)
    return (
      <div className="flex justify-center items-center h-screen">
        <LoadingSpinner size={48} />
      </div>
    );
  if (error) return <ErrorMessage message={error} />;

  return (
    <div className="container mx-auto px-4 sm:px-6 lg:px-8 py-8 space-y-12">
      <section>
        <SectionTitle title="Trending Now" icon={TrendingUp} />
        <AnimeGrid animes={trending} />
      </section>
      <section>
        <SectionTitle title="Popular Anime" icon={Star} />
        <AnimeGrid animes={popular} />
      </section>
      <section>
        <SectionTitle title="Recently Released" icon={Clock} />
        <AnimeGrid animes={recentlyReleased} />
      </section>
      <section>
        <SectionTitle title="Upcoming Anime" icon={Calendar} />
        <AnimeGrid animes={upcoming} />
      </section>
    </div>
  );
};

export default HomePage;

import { useState, useEffect } from "react";
import apiService from "../api/api";
import AnimeGrid from "../components/anime/AnimeGrid";
import ErrorMessage from "../components/ui/ErrorMessage";
import LoadingSpinner from "../components/ui/LoadingSpinner";
import Button from "../components/ui/Button";

const SearchResultsPage = ({ query }) => {
  const [results, setResults] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const perPage = 24;

  useEffect(() => {
    if (!query) {
      setResults([]);
      setTotal(0);
      return;
    }
    const fetchResults = async () => {
      setLoading(true);
      setError(null);
      try {
        const data = await apiService.searchAnime(query, page, perPage);
        setResults((prevResults) =>
          page === 1 ? data.data : [...prevResults, ...data.data]
        );
        setTotal(data.meta.total);
      } catch (err) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };
    fetchResults();
  }, [query, page]);

  const loadMore = () => {
    if (results.length < total) {
      setPage((prevPage) => prevPage + 1);
    }
  };

  return (
    <div className="container mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <h1 className="text-3xl font-bold mb-6 text-slate-800 dark:text-slate-100">
        Search Results for: "{query}"
      </h1>
      {error && <ErrorMessage message={error} />}
      <AnimeGrid animes={results} />
      {loading && (
        <div className="flex justify-center py-6">
          <LoadingSpinner size={36} />
        </div>
      )}
      {!loading && results.length > 0 && results.length < total && (
        <div className="text-center mt-8">
          <Button onClick={loadMore} variant="primary">
            Load More
          </Button>
        </div>
      )}
      {!loading && results.length === 0 && !error && (
        <p>No results found for "{query}".</p>
      )}
    </div>
  );
};

export default SearchResultsPage;

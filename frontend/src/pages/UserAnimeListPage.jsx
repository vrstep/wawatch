const UserAnimeListPage = () => {
  const { user, loading: authLoading } = useAuth();
  const [animeList, setAnimeList] = useState([]);
  const [statusFilter, setStatusFilter] = useState(""); // ALL, WATCHING, COMPLETED, etc.
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const perPage = 20;

  const statuses = [
    "ALL",
    "WATCHING",
    "COMPLETED",
    "PLANNED",
    "PAUSED",
    "DROPPED",
    "REWATCHING",
  ];

  const fetchUserList = useCallback(
    async (currentPage, currentStatus) => {
      if (!user) return;
      setLoading(true);
      setError(null);
      try {
        const effectiveStatus = currentStatus === "ALL" ? "" : currentStatus;
        const data = await apiService.getUserAnimeList(
          effectiveStatus,
          currentPage,
          perPage
        );
        setAnimeList((prev) =>
          currentPage === 1 ? data.data : [...prev, ...data.data]
        );
        setTotal(data.meta.total);
      } catch (err) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    },
    [user, perPage]
  );

  useEffect(() => {
    setPage(1); // Reset page when filter changes
    fetchUserList(1, statusFilter);
  }, [statusFilter, fetchUserList]);

  useEffect(() => {
    if (page > 1) {
      // Only fetch if it's not the initial load triggered by statusFilter change
      fetchUserList(page, statusFilter);
    }
  }, [page, fetchUserList, statusFilter]);

  if (authLoading)
    return (
      <div className="flex justify-center items-center h-screen">
        <LoadingSpinner size={48} />
      </div>
    );
  if (!user) {
    useEffect(() => {
      window.location.hash = "#/login";
    }, []);
    return null;
  }

  const handleStatusChange = (newStatus) => {
    setStatusFilter(newStatus);
  };

  const handleAnimeStatusUpdateInList = (animeId, newStatus) => {
    // Refetch list to reflect changes, or update locally for better UX
    setAnimeList((prevList) =>
      prevList
        .map((item) =>
          item.AnimeDetails && item.AnimeDetails.id === animeId
            ? { ...item, Status: newStatus, userStatus: newStatus }
            : item
        )
        .filter(
          (item) =>
            newStatus !== "REMOVING" ||
            (item.AnimeDetails && item.AnimeDetails.id !== animeId)
        )
    );
    // A full refetch might be simpler to ensure data consistency:
    // fetchUserList(1, statusFilter);
  };

  const loadMore = () => {
    if (animeList.length < total) {
      setPage((prevPage) => prevPage + 1);
    }
  };

  return (
    <div className="container mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <SectionTitle title="My Anime List" icon={List} />
      <div className="mb-6 flex items-center space-x-2">
        <Label htmlFor="status-filter">Filter by status:</Label>
        <select
          id="status-filter"
          value={statusFilter}
          onChange={(e) => handleStatusChange(e.target.value)}
          className="bg-slate-100 dark:bg-slate-700 border border-slate-300 dark:border-slate-600 text-slate-900 dark:text-slate-50 text-sm rounded-lg focus:ring-indigo-500 focus:border-indigo-500 block p-2.5"
        >
          {statuses.map((s) => (
            <option key={s} value={s}>
              {s.charAt(0) + s.slice(1).toLowerCase()}
            </option>
          ))}
        </select>
      </div>
      {error && <ErrorMessage message={error} />}
      <AnimeGrid
        animes={animeList.map((item) => ({
          ...item.AnimeDetails,
          userStatus: item.Status,
        }))}
        onStatusChange={handleAnimeStatusUpdateInList}
      />
      {loading && page === 1 && (
        <div className="flex justify-center py-6">
          <LoadingSpinner size={36} />
        </div>
      )}
      {!loading && animeList.length > 0 && animeList.length < total && (
        <div className="text-center mt-8">
          <Button onClick={loadMore} variant="primary" disabled={loading}>
            {loading && page > 1 ? (
              <LoadingSpinner size={20} className="mr-2" />
            ) : null}{" "}
            Load More
          </Button>
        </div>
      )}
      {!loading && animeList.length === 0 && !error && (
        <p>Your list for "{statusFilter}" is empty.</p>
      )}
    </div>
  );
};
export default UserAnimeListPage;

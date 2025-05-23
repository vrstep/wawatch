const UserHistoryPage = () => {
  const { user, loading: authLoading } = useAuth();
  const [history, setHistory] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const perPage = 30;

  useEffect(() => {
    if (!user) return;
    const fetchHistory = async () => {
      setLoading(true);
      setError(null);
      try {
        const data = await apiService.getUserHistory(page, perPage);
        setHistory((prev) =>
          page === 1 ? data.data : [...prev, ...data.data]
        );
        setTotal(data.meta.total);
      } catch (err) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };
    fetchHistory();
  }, [user, page, perPage]);

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

  const loadMore = () => {
    if (history.length < total) {
      setPage((prevPage) => prevPage + 1);
    }
  };

  return (
    <div className="container mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <SectionTitle title="My View History" icon={History} />
      {error && <ErrorMessage message={error} />}
      {history.length > 0 ? (
        // AnimeGrid expects slightly different structure, so we adapt or make a new component
        <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-4 md:gap-6">
          {history.map((entry) =>
            entry.AnimeDetails ? (
              <AnimeCard key={entry.ID} anime={entry.AnimeDetails} />
            ) : null
          )}
        </div>
      ) : (
        !loading && (
          <p className="text-center text-slate-500 dark:text-slate-400 py-8">
            Your viewing history is empty.
          </p>
        )
      )}
      {loading && page === 1 && (
        <div className="flex justify-center py-6">
          <LoadingSpinner size={36} />
        </div>
      )}
      {!loading && history.length > 0 && history.length < total && (
        <div className="text-center mt-8">
          <Button onClick={loadMore} variant="primary" disabled={loading}>
            {loading && page > 1 ? (
              <LoadingSpinner size={20} className="mr-2" />
            ) : null}{" "}
            Load More
          </Button>
        </div>
      )}
    </div>
  );
};
export default UserHistoryPage;

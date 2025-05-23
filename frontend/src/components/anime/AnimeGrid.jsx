import AnimeCard from "./AnimeCard";

const AnimeGrid = ({ animes, onStatusChange }) => {
  if (!animes || animes.length === 0) {
    return (
      <p className="text-center text-slate-500 dark:text-slate-400 py-8">
        No anime found.
      </p>
    );
  }
  return (
    <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-4 md:gap-6">
      {animes.map((anime) => (
        <AnimeCard
          key={anime.id}
          anime={anime}
          onStatusChange={onStatusChange}
        />
      ))}
    </div>
  );
};
export default AnimeGrid;

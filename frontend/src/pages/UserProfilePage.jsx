import { useEffect } from "react";
import { useAuth } from "../contexts/AuthContext";
import Button from "../components/ui/Button";
import LoadingSpinner from "../components/ui/LoadingSpinner";

const UserProfilePage = () => {
  const { user, loading: authLoading } = useAuth();
  // Add states for user's lists, stats etc.
  // For now, just displays basic user info.
  // This page could be expanded to show user's anime lists, stats, history tabs.

  if (authLoading)
    return (
      <div className="flex justify-center items-center h-screen">
        <LoadingSpinner size={48} />
      </div>
    );
  if (!user) {
    // Redirect to login if not authenticated
    useEffect(() => {
      window.location.hash = "#/login";
    }, []);
    return null;
  }

  return (
    <div className="container mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div className="max-w-2xl mx-auto bg-slate-100 dark:bg-slate-800 rounded-lg shadow-xl p-8">
        <div className="flex flex-col items-center">
          <img
            src={
              user.profile_picture ||
              `https://placehold.co/128x128/7c3aed/ffffff?text=${user.username
                .charAt(0)
                .toUpperCase()}`
            }
            alt={user.username}
            className="w-32 h-32 rounded-full object-cover mb-6 border-4 border-indigo-500"
            onError={(e) => {
              e.target.onerror = null;
              e.target.src = `https://placehold.co/128x128/7c3aed/ffffff?text=${user.username
                .charAt(0)
                .toUpperCase()}`;
            }}
          />
          <h1 className="text-3xl font-bold text-slate-800 dark:text-slate-100">
            {user.username}
          </h1>
          <p className="text-slate-600 dark:text-slate-400">{user.email}</p>
          <p className="text-xs text-slate-500 dark:text-slate-500 mt-1">
            Joined: {new Date(user.created_at).toLocaleDateString()}
          </p>
          <div className="mt-6">
            <Button
              onClick={() => (window.location.hash = "#/me/animelist")}
              className="mr-2"
            >
              My Anime List
            </Button>
            <Button
              onClick={() => (window.location.hash = "#/me/history")}
              variant="secondary"
              className="mr-2"
            >
              View History
            </Button>
            <Button
              onClick={() => (window.location.hash = "#/me/settings")}
              variant="ghost"
            >
              <Settings className="h-5 w-5" />
            </Button>
          </div>
        </div>
        {/* TODO: Add tabs for Watchlist, Stats, History here if desired */}
      </div>
    </div>
  );
};

export default UserProfilePage;

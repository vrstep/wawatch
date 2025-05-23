import { useAuth } from "../../contexts/AuthContext";
import { useState } from "react";
import { useTheme } from "../../contexts/ThemeContext";
import Input from "../ui/Input";
import Button from "../ui/Button";
import { Search, LogIn, LogOut, Sun, Moon } from "lucide-react";
import { Tv } from "lucide-react";

const Navbar = () => {
  const { user, logout } = useAuth();
  const [searchTerm, setSearchTerm] = useState("");
  const { theme, toggleTheme } = useTheme();

  const handleSearch = (e) => {
    e.preventDefault();
    if (searchTerm.trim()) {
      window.location.hash = `#/search?q=${encodeURIComponent(
        searchTerm.trim()
      )}`;
    }
  };

  return (
    <nav className="bg-slate-100 dark:bg-slate-900 shadow-md sticky top-0 z-50">
      <div className="container mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-16">
          <div className="flex items-center">
            <a
              href="#/"
              className="text-2xl font-bold text-indigo-600 dark:text-indigo-400 flex items-center"
            >
              <Tv className="h-8 w-8 mr-2" /> WaWatch
            </a>
          </div>
          <div className="flex-1 max-w-xl mx-4">
            <form onSubmit={handleSearch} className="relative">
              <Input
                type="search"
                placeholder="Search anime..."
                className="w-full !pr-10"
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
              />
              <Button
                type="submit"
                variant="ghost"
                size="icon"
                className="absolute right-0 top-0 h-full"
              >
                <Search className="h-5 w-5" />
              </Button>
            </form>
          </div>
          <div className="flex items-center space-x-3">
            <Button
              variant="ghost"
              size="icon"
              onClick={toggleTheme}
              aria-label="Toggle theme"
            >
              {theme === "dark" ? (
                <Sun className="h-5 w-5" />
              ) : (
                <Moon className="h-5 w-5" />
              )}
            </Button>
            {user ? (
              <>
                <a href="#/me/profile">
                  <img
                    src={
                      user.profile_picture ||
                      `https://placehold.co/40x40/7c3aed/ffffff?text=${user.username
                        .charAt(0)
                        .toUpperCase()}`
                    }
                    alt={user.username}
                    className="h-8 w-8 rounded-full object-cover"
                    onError={(e) => {
                      e.target.onerror = null;
                      e.target.src = `https://placehold.co/40x40/7c3aed/ffffff?text=${user.username
                        .charAt(0)
                        .toUpperCase()}`;
                    }}
                  />
                </a>
                <Button onClick={logout} variant="secondary" size="sm">
                  <LogOut className="h-4 w-4 mr-2" /> Logout
                </Button>
              </>
            ) : (
              <>
                <Button
                  onClick={() => (window.location.hash = "#/login")}
                  variant="ghost"
                  size="sm"
                >
                  <LogIn className="h-4 w-4 mr-2" /> Login
                </Button>
                <Button
                  onClick={() => (window.location.hash = "#/signup")}
                  variant="primary"
                  size="sm"
                >
                  Sign Up
                </Button>
              </>
            )}
          </div>
        </div>
      </div>
    </nav>
  );
};
export default Navbar;

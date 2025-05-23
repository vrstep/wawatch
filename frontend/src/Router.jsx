import { useState, useEffect } from "react";
import HomePage from "./pages/HomePage";
import AnimeDetailPage from "./pages/AnimeDetailPage";
import LoginPage from "./pages/LoginPage";
import SignupPage from "./pages/SignupPage";
import UserProfilePage from "./pages/UserProfilePage";
import UserAnimeListPage from "./pages/UserAnimeListPage";
import UserHistoryPage from "./pages/UserHistoryPage";
import SearchResultsPage from "./pages/SearchResultPage";
import SettingsPage from "./pages/SettingsPage";
import NotFoundPage from "./pages/NotFoundPage";

const Router = () => {
  const [route, setRoute] = useState(window.location.hash || "#/");

  useEffect(() => {
    const handleHashChange = () => setRoute(window.location.hash || "#/");
    window.addEventListener("hashchange", handleHashChange);
    return () => window.removeEventListener("hashchange", handleHashChange);
  }, []);

  let Component;
  const params = {};

  if (route.startsWith("#/anime/")) {
    Component = AnimeDetailPage;
    params.id = route.split("/")[2];
  } else if (route.startsWith("#/search")) {
    Component = SearchResultsPage;
    const searchParams = new URLSearchParams(route.split("?")[1] || "");
    params.query = searchParams.get("q");
  } else if (route === "#/login") {
    Component = LoginPage;
  } else if (route === "#/signup") {
    Component = SignupPage;
  } else if (route === "#/me/profile" || route === "#/profile") {
    // Alias /profile to /me/profile
    Component = UserProfilePage;
  } else if (route === "#/me/animelist") {
    Component = UserAnimeListPage;
  } else if (route === "#/me/history") {
    Component = UserHistoryPage;
  } else if (route === "#/me/settings") {
    Component = SettingsPage;
  } else if (route === "#/" || route === "") {
    Component = HomePage;
  } else {
    Component = NotFoundPage;
  }

  return Component ? <Component {...params} /> : <NotFoundPage />;
};

export default Router;

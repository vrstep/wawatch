import { createContext, useState, useEffect } from "react"; // Added createContext here
import { useContext, useMemo } from "react";

import apiService from "../api/api"; // Corrected import path
// --- AUTH CONTEXT ---
const AuthContext = createContext(null);

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true); // For initial auth check

  useEffect(() => {
    const validateUser = async () => {
      try {
        const data = await apiService.validate();
        setUser(data.user);
      } catch (error) {
        setUser(null);
        localStorage.removeItem("authToken"); // Example if using token in LS
      } finally {
        setLoading(false);
      }
    };
    validateUser();
  }, []);

  const login = async (credentials) => {
    const data = await apiService.login(credentials); // Backend sets HttpOnly cookie
    setUser(data.user);
    // localStorage.setItem('authToken', data.token); // If also storing token for non-HttpOnly use
    return data;
  };

  const signup = async (userData) => {
    return await apiService.signup(userData);
  };

  const logout = () => {
    // No backend logout endpoint provided, so client-side only for now
    // If backend had /logout that invalidates cookie, call it here:
    // await apiService.logout();
    setUser(null);
    localStorage.removeItem("authToken");
    // Navigate to home or login page
    window.location.hash = "#/login";
  };

  const value = useMemo(
    () => ({ user, setUser, login, logout, signup, loading }),
    [user, loading]
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

export const useAuth = () => useContext(AuthContext);

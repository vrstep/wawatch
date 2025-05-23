// src/api/apiService.js

const API_BASE_URL =
  import.meta.env.VITE_API_BASE_URL || "http://localhost:8080";
const ANIME_API_BASE_URL = `${API_BASE_URL}/ext/anime`;

const apiService = {
  async request(
    endpoint,
    method = "GET",
    data = null,
    requiresAuth = false, // This will be 'true' for anime calls now
    customHeaders = {}
  ) {
    let url = endpoint;
    if (!endpoint.startsWith("http://") && !endpoint.startsWith("https://")) {
      url = `${API_BASE_URL}${endpoint}`;
    }
    // Log the final URL and auth requirement for clarity
    console.log(
      "[API Service] Requesting:",
      method,
      url,
      "Requires Auth:",
      requiresAuth
    );

    const headers = {
      "Content-Type": "application/json",
      ...customHeaders,
    };

    const config = {
      method,
      headers,
      credentials:
        requiresAuth || endpoint.startsWith("/api/v1/auth") // Use "include" if requiresAuth is true or for auth paths
          ? "include"
          : "same-origin", // Otherwise, 'same-origin' is a reasonable default
    };

    if (data && method !== "GET") {
      config.body = JSON.stringify(data);
    }

    try {
      const response = await fetch(url, config);

      if (response.status === 204) return null;

      const responseText = await response.text(); // Get raw text of the response

      if (!response.ok) {
        let errorData = null;
        try {
          if (responseText) errorData = JSON.parse(responseText);
        } catch (e) {
          // Parsing JSON error body failed, errorData will use raw text
        }
        console.error(
          `API Error (${url}): ${response.status}`,
          errorData || responseText
        );
        const message =
          (errorData && (errorData.error || errorData.message)) ||
          (responseText && responseText.length > 0 && responseText.length < 200
            ? responseText
            : `Request failed with status ${response.status}`);
        const error = new Error(message);
        error.status = response.status;
        error.data = errorData || { rawError: responseText };
        throw error;
      }

      // If response is OK
      if (!responseText) {
        console.warn(
          `[API Service] Successful response for ${url} (status ${response.status}) but body is empty.`
        );
        return null; // Or based on API contract for empty successful responses
      }

      try {
        return JSON.parse(responseText);
      } catch (e) {
        console.error(
          `Malformed JSON in successful response for ${url}:`,
          responseText,
          e
        );
        const error = new Error("Malformed JSON response from server.");
        error.status = response.status;
        error.data = { rawResponse: responseText };
        throw error;
      }
    } catch (error) {
      // Log the error with more details if available
      const errorMsg = `Network or other error for ${url}: ${error.message}${
        error.status ? ` (Status: ${error.status})` : ""
      }`;
      console.error(errorMsg, error.data ? error.data : error);
      if (!error.status) {
        error.message = error.message || "Network error or server unreachable.";
      }
      throw error; // Re-throw the possibly augmented error
    }
  },

  // Auth Endpoints
  signup: (userData) =>
    apiService.request("/api/v1/auth/signup", "POST", userData),
  login: (credentials) =>
    apiService.request("/api/v1/auth/login", "POST", credentials, true), // Login also likely needs to handle cookies
  validate: () =>
    apiService.request("/api/v1/auth/validate", "GET", null, true),
  logout: () => {
    // Implement client-side logout (clear context, local storage)
    // Optionally call a backend logout endpoint if it exists:
    // return apiService.request("/api/v1/auth/logout", "POST", null, true);
  },

  // ... other user-specific endpoints should likely have requiresAuth = true ...

  // Anime Data (Pass-through) Endpoints - MODIFIED
  // Assuming these data points might be user-specific or you want to ensure the user session is validated by the gateway
  searchAnime: (query, page = 1, perPage = 20) =>
    apiService.request(
      `${ANIME_API_BASE_URL}/search?q=${encodeURIComponent(
        query
      )}&page=${page}&perPage=${perPage}`,
      "GET",
      null,
      true // CHANGED
    ),
  getPopularAnime: (
    page = 1,
    perPage = 12 // Matched perPage to your logs
  ) =>
    apiService.request(
      `${ANIME_API_BASE_URL}/popular?page=${page}&perPage=${perPage}`,
      "GET",
      null,
      true // CHANGED
    ),
  getTrendingAnime: (page = 1, perPage = 12) =>
    apiService.request(
      `${ANIME_API_BASE_URL}/trending?page=${page}&perPage=${perPage}`,
      "GET",
      null,
      true // CHANGED
    ),
  getUpcomingAnime: (page = 1, perPage = 12) =>
    apiService.request(
      `${ANIME_API_BASE_URL}/upcoming?page=${page}&perPage=${perPage}`,
      "GET",
      null,
      true // CHANGED
    ),
  getRecentlyReleasedAnime: (page = 1, perPage = 12) =>
    apiService.request(
      `${ANIME_API_BASE_URL}/recently-released?page=${page}&perPage=${perPage}`,
      "GET",
      null,
      true // CHANGED
    ),
  // getAnimeRecommendations and getAnimeDetails already have requiresAuth = true, which is good.
  getAnimeRecommendations: (page = 1, perPage = 10) =>
    apiService.request(
      `${ANIME_API_BASE_URL}/recommendations?page=${page}&perPage=${perPage}`,
      "GET",
      null,
      true
    ),
  getAnimeDetails: (animeId) =>
    apiService.request(`${ANIME_API_BASE_URL}/${animeId}`, "GET", null, true),
};

export default apiService;

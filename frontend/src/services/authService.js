import axios from 'axios';

const API_URL = 'http://localhost:8080/api';

/**
 * Authentication service to handle JWT verification
 * Since the JWT is stored in a HttpOnly cookie, we need to
 * verify authentication through API calls instead of reading the cookie directly
 */
const authService = {
  /**
   * Check if the user is authenticated by verifying the token with the backend
   * @returns {Promise<boolean>} True if authenticated, false otherwise
   */
  async isAuthenticated() {
    try {
      const response = await axios.get(`${API_URL}/verify`, {
        withCredentials: true // Important to send cookies with the request
      });
      return response.status === 200;
    } catch (error) {
      console.error('Authentication check failed:', error);
      return false;
    }
  },

  /**
   * Login user
   * @param {Object} credentials User credentials
   * @returns {Promise<Object>} Login response
   */
  async login(credentials) {
    const response = await axios.post(`${API_URL}/login`, credentials, {
      withCredentials: true
    });
    return response.data;
  },

  /**
   * Logout user
   */
  async logout() {
    await axios.post(`${API_URL}/logout`, {}, {
      withCredentials: true
    });
  }
};

export default authService; 
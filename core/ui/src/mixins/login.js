export default {
  name: "LoginService",
  methods: {
    executeLogin(username, password, otp) {
      return this.axios.post(`${this.$root.apiUrl}/login`, {
        username,
        password,
        otp,
      });
    },
    executeLogout() {
      const loginInfo = this.getFromStorage("loginInfo");
      const token = loginInfo ? loginInfo.token : "";

      // send the refresh token too: the server ends the session even if
      // the access token is expired
      const body = loginInfo?.refreshToken
        ? { refresh_token: loginInfo.refreshToken }
        : {};

      return this.axios.post(`${this.$root.apiUrl}/logout`, body, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
    },
    executeRefresh(refreshToken) {
      return this.axios.post(`${this.$root.apiUrl}/refresh`, {
        refresh_token: refreshToken,
      });
    },
    executeLogoutOthers() {
      const token = this.getFromStorage("loginInfo")
        ? this.getFromStorage("loginInfo").token
        : "";

      return this.axios.post(
        `${this.$root.apiUrl}/logout-others`,
        {},
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );
    },
  },
};

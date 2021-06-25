export default {
  name: "LoginService",
  methods: {
    executeLogin(username, password) {
      return this.axios.post(`${this.$root.apiUrl}/login`, {
        username,
        password,
      });
    },
    executeLogout() {
      const token = this.getFromStorage("loginInfo")
        ? this.getFromStorage("loginInfo").token
        : "";

      return this.axios.post(
        `${this.$root.apiUrl}/logout`,
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

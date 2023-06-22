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

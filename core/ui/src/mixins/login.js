export default {
  name: "LoginService",
  computed: {
    apiUrl() {
      return this.$root.config.API_SCHEME + this.$root.config.API_ENDPOINT;
    },
  },
  methods: {
    executeLogin(username, password) {
      return this.axios.post(this.apiUrl + "/login", { username, password });
    },
    executeLogout() {
      const token = this.getFromStorage("loginInfo")
        ? this.getFromStorage("loginInfo").token
        : "";

      return this.axios.post(
        this.apiUrl + "/logout",
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

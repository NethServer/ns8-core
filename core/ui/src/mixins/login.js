export default {
  name: "LoginService",
  computed: {
    apiUrl() {
      return this.$root.config.API_SCHEME + this.$root.config.API_ENDPOINT;
    },
  },
  methods: {
    login(username, password) {
      return this.axios.post(this.apiUrl + "/login", { username, password });
    },
  },
};

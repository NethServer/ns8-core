import StorageService from "@/mixins/storage";

export default {
  name: "NodeService",
  mixins: [StorageService],
  methods: {
    getNodes() {
      const token = this.getFromStorage("loginInfo")
        ? this.getFromStorage("loginInfo").token
        : "";
      return this.axios.get(`${this.$root.apiUrl}/nodes`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
    },
  },
};

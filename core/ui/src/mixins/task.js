import StorageService from "@/mixins/storage";

export default {
  name: "TaskService",
  mixins: [StorageService],
  computed: {
    apiUrl() {
      return this.$root.config.API_SCHEME + this.$root.config.API_ENDPOINT;
    },
  },
  methods: {
    getTaskContext(taskPath) {
      const token = this.getFromStorage("loginInfo")
        ? this.getFromStorage("loginInfo").token
        : "";
      return this.axios.get(`${this.apiUrl}/tasks/${taskPath}/context`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
    },
    getClusterTasks() {
      const token = this.getFromStorage("loginInfo")
        ? this.getFromStorage("loginInfo").token
        : "";
      return this.axios.get(this.apiUrl + "/tasks/cluster", {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
    },
    createTask(taskData) {
      const token = this.getFromStorage("loginInfo")
        ? this.getFromStorage("loginInfo").token
        : "";
      return this.axios.post(this.apiUrl + "/tasks/cluster", taskData, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
    },
  },
};

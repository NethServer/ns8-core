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
      return this.axios.get(`${this.apiUrl}/${taskPath}/context`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
    },
    getTaskStatus(taskPath) {
      const token = this.getFromStorage("loginInfo")
        ? this.getFromStorage("loginInfo").token
        : "";
      return this.axios.get(`${this.apiUrl}/${taskPath}/status`, {
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
      return this.axios.post(this.apiUrl + "/cluster/tasks", taskData, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
    },
    getTaskTitle(task) {
      if (
        task &&
        task.context &&
        task.context.data &&
        task.context.data.title
      ) {
        return task.context.data.title;
      } else {
        return "";
      }
    },
    getTaskKind(task) {
      switch (task.status) {
        case "aborted":
        case "validation-failed":
          return "error";
        case "completed":
          return "success";
        case "pending":
          return "warning";
        default:
          return "info";
      }
    },
    getTaskStatusDescription(task, rootTask = true) {
      const taskAction = task.context.action;
      const taskOrSubtask = rootTask ? "task" : "subtask";

      console.log("!!! task, taskOrSubtask", task, taskOrSubtask); ////

      switch (task.status) {
        case "aborted":
          return this.$t("task." + taskOrSubtask + "_failed", {
            action: taskAction,
          });
        case "validation-failed":
          return this.$t("task." + taskOrSubtask + "_failed_validation", {
            action: taskAction,
          });
        case "completed":
          return this.$t("task." + taskOrSubtask + "_completed", {
            action: taskAction,
          });
        case "pending":
          return this.$t("task." + taskOrSubtask + "_pending", {
            action: taskAction,
          });
        default:
          return "";
      }
    },
  },
};

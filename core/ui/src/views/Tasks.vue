<template>
  <div>
    <h1>Tasks page</h1>
    <InlineNotification
      v-if="error.clusterTasks"
      kind="error"
      title="Cannot retrieve tasks:"
      :description="error.clusterTasks"
      low-contrast
    />
  </div>
</template>

<script>
import TaskService from "@/mixins/task";
import to from "await-to-js";
import InlineNotification from "@/components/InlineNotification";
import ErrorService from "@/mixins/error";

//// delete file?

export default {
  name: "Tasks",
  components: { InlineNotification },
  mixins: [TaskService, ErrorService],
  data() {
    return {
      error: {
        clusterTasks: "",
      },
    };
  },
  mounted() {
    this.retrieveTasks();
  },
  methods: {
    async retrieveTasks() {
      const [taskError, response] = await to(this.getClusterTasks());

      if (taskError) {
        this.error.clusterTasks = this.getErrorMessage(taskError);
        return;
      }

      console.log("tasks", response); ////
    },
  },
};
</script>

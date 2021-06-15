<template>
  <div class="indent">
    <div v-for="(subTask, index) in subTasks" :key="index">
      <InlineNotification
        :kind="getTaskKind(subTask)"
        :title="getTaskStatusDescription(subTask, false)"
        low-contrast
      />
      <TaskHierarchy
        v-if="subTask.subTasks.length"
        :subTasks="subTask.subTasks"
      />
    </div>
  </div>
</template>

<script>
import TaskService from "@/mixins/task";
import InlineNotification from "@/components/InlineNotification";

export default {
  name: "TaskHierarchy",
  mixins: [TaskService],
  components: { InlineNotification },
  props: {
    subTasks: {
      type: Array,
      required: true,
    },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";

.indent {
  padding-left: $spacing-09;
}
</style>

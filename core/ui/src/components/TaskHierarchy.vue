<template>
  <div class="indent">
    <div v-for="(subTask, index) in subTasks" :key="index">
      <NsInlineNotification
        :kind="getTaskKind(subTask)"
        :title="getTaskStatusDescription(subTask, false)"
        low-contrast
        :showCloseButton="false"
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
import NsInlineNotification from "@/components/NsInlineNotification";

export default {
  name: "TaskHierarchy",
  mixins: [TaskService],
  components: { NsInlineNotification },
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

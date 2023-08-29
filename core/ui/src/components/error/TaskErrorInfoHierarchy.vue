<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <div class="mg-left-07">
    <div v-for="task in subTasks" :key="task.context.id">
      <TaskErrorInfo :task="task" />
    </div>
  </div>
</template>

<script>
import { TaskService } from "@nethserver/ns8-ui-lib";

export default {
  name: "TaskErrorInfoHierarchy",
  mixins: [TaskService],
  components: {},
  props: {
    subTasks: {
      type: Array,
      required: true,
    },
  },
  // needed to solve circular dependency between TaskErrorInfo & TaskErrorInfoHierarchy:
  beforeCreate: function () {
    this.$options.components.TaskErrorInfo =
      require("./TaskErrorInfo.vue").default;
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";

.mg-left-07 {
  margin-left: $spacing-07;
}
</style>

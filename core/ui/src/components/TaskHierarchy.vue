<template>
  <ul>
    <li v-for="(subTask, index) in subTasks" :key="index">
      <div class="task-item">
        <component
          :is="getTaskIcon(subTask)"
          :class="`bx--inline-notification__icon ${subTask.status}`"
        />
        <span v-html="getTaskStatusDescription(subTask, false)"></span>
      </div>
      <TaskHierarchy
        v-if="subTask.subTasks.length"
        :subTasks="subTask.subTasks"
      />
    </li>
  </ul>
</template>

<script>
import TaskService from "@/mixins/task";

export default {
  name: "TaskHierarchy",
  mixins: [TaskService],
  components: {},
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

ul {
  padding: 0;
  margin: 0;
  list-style-type: none;
  position: relative;
}

li {
  list-style-type: none;
  border-left: 1px solid #000;
  margin-left: 1.3rem;
}

li div {
  padding-left: 1em;
  position: relative;
}

li div::before {
  content: "";
  position: absolute;
  top: 0;
  left: -1px;
  bottom: 50%;
  width: 0.75em;
  border: 1px solid #000;
  border-top: 0 none transparent;
  border-right: 0 none transparent;
}

ul > li:last-child {
  border-left: 1px solid transparent;
}
</style>

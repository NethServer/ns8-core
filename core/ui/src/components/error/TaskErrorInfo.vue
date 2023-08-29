<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <div>
    <h6 class="mg-top-lg mg-bottom-sm">
      {{ taskTitle }}
    </h6>
    <div>
      {{ $t("task.task_id") }}:
      <cv-tooltip
        alignment="center"
        direction="top"
        :tip="
          justCopied
            ? $t('common.copied_to_clipboard')
            : $t('common.copy_to_clipboard')
        "
      >
        <cv-link
          v-clipboard:copy="task.context.id"
          v-clipboard:success="onCopy"
          v-clipboard:error="onCopyError"
        >
          {{ task.context.id }}
        </cv-link></cv-tooltip
      >
    </div>
    <div
      v-if="task.result && task.result.error && task.result.exit_code != 0"
      class="mg-top-sm"
    >
      <NsCodeSnippet
        :copyTooltip="$t('common.copy_to_clipboard')"
        :copy-feedback="$t('common.copied_to_clipboard')"
        :feedback-aria-label="$t('common.copied_to_clipboard')"
        :wrap-text="true"
        :moreText="$t('common.show_more')"
        :lessText="$t('common.show_less')"
        light
        >{{ task.result.error }}</NsCodeSnippet
      >
    </div>
    <TaskErrorInfoHierarchy
      v-if="task.subTasks.length"
      :subTasks="task.subTasks"
    >
    </TaskErrorInfoHierarchy>
  </div>
</template>

<script>
import { TaskService } from "@nethserver/ns8-ui-lib";
import TaskErrorInfoHierarchy from "./TaskErrorInfoHierarchy";

export default {
  name: "TaskErrorInfo",
  mixins: [TaskService],
  components: { TaskErrorInfoHierarchy },
  props: {
    task: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      justCopied: false,
    };
  },
  computed: {
    taskTitle() {
      const agent = this.task.context.queue
        ? this.task.context.queue.split("/tasks")[0]
        : "cluster";

      return `${agent}/${this.task.context.action}`;
    },
  },
  methods: {
    onCopy() {
      this.justCopied = true;

      setTimeout(() => {
        this.justCopied = false;
      }, 3000);
    },
    onCopyError() {
      console.error("cannot copy to clipboard");
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";
</style>

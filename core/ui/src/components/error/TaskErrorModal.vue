<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <NsModal
    size="default"
    :visible="isTaskErrorShown"
    @modal-hidden="taskErrorModalHidden"
    @primary-click="onPrimaryClick"
    class="task-error-modal"
  >
    <template slot="title">{{ getTaskTitle(taskErrorToShow) }}</template>
    <template slot="content" v-if="taskErrorToShow">
      <div class="task-item">
        <component
          :is="getTaskIcon(taskErrorToShow)"
          :class="`bx--inline-notification__icon ${taskErrorToShow.status} task-icon`"
        />
        <div>
          <span>{{ getTaskStatusDescription(taskErrorToShow, true) }}</span>
        </div>
      </div>
      <div v-if="subTasks.length">
        <TaskTreeHierarchy
          :subTasks="subTasks"
          :isMoreInfoShown="isMoreInfoShown"
          class="task-hierarchy"
        />
      </div>
      <!-- copy task json task data to clipboard -->
      <cv-tooltip
        alignment="center"
        direction="right"
        :tip="
          justCopied
            ? $t('common.copied_to_clipboard')
            : $t('common.copy_to_clipboard')
        "
        class="mg-top-md"
      >
        <cv-link
          v-clipboard:copy="JSON.stringify(taskErrorToShow)"
          v-clipboard:success="onCopy"
          v-clipboard:error="onCopyError"
        >
          {{ $t("task.copy_task_data_to_clipboard") }}
        </cv-link></cv-tooltip
      >
      <cv-toggle
        value="moreInfoValue"
        :form-item="true"
        v-model="isMoreInfoShown"
        hide-label
        class="mg-top-sm mg-bottom-md"
      >
        <template slot="text-left">{{ $t("common.more_info") }}</template>
        <template slot="text-right">{{ $t("common.more_info") }}</template>
      </cv-toggle>

      <template v-if="isMoreInfoShown">
        <!-- copy to clipboard hint -->
        <span class="hint hint-copy-to-clipboard">
          <cv-interactive-tooltip
            alignment="end"
            direction="bottom"
            :visible="isCopyClipboardHintShown"
          >
            <template slot="trigger">
              <span></span>
            </template>
            <template slot="content">
              <p>
                {{ $t("hint.copy_to_clipboard") }}
              </p>
              <NsButton
                kind="primary"
                size="small"
                @click="isCopyClipboardHintShown = false"
                class="hint-button"
                >{{ $t("common.got_it") }}</NsButton
              >
            </template>
          </cv-interactive-tooltip>
        </span>
        <TaskErrorInfo :task="taskErrorToShow" />
      </template>
    </template>
    <template v-if="isLogAvailable">
      <template slot="secondary-button">{{ $t("common.close") }}</template>
      <template slot="primary-button">{{
        $t("common.view_instance_logs", {
          instance: taskErrorToShow.context.extra.logs.instance,
        })
      }}</template>
    </template>
    <template v-else>
      <template slot="primary-button">{{ $t("common.close") }}</template>
    </template>
  </NsModal>
</template>

<script>
import { mapState } from "vuex";
import TaskTreeHierarchy from "./TaskTreeHierarchy";
import TaskErrorInfoHierarchy from "./TaskErrorInfoHierarchy";
import TaskErrorInfo from "./TaskErrorInfo";
import {
  TaskService,
  StorageService,
  IconService,
} from "@nethserver/ns8-ui-lib";

export default {
  name: "TaskErrorModal",
  components: { TaskTreeHierarchy, TaskErrorInfoHierarchy, TaskErrorInfo },
  mixins: [TaskService, StorageService, IconService],
  props: {},
  data() {
    return {
      isCopyClipboardHintShown: false,
      isMoreInfoShown: false,
      justCopied: false,
    };
  },
  computed: {
    ...mapState(["taskErrorToShow", "isTaskErrorShown"]),
    subTasks: function () {
      if (this.taskErrorToShow && this.taskErrorToShow.subTasks) {
        return this.taskErrorToShow.subTasks;
      } else {
        return [];
      }
    },
    isLogAvailable: function () {
      return (
        this.taskErrorToShow &&
        this.taskErrorToShow.context &&
        this.taskErrorToShow.context.extra &&
        this.taskErrorToShow.context.extra.logs
      );
    },
  },
  watch: {
    isMoreInfoShown: function () {
      if (this.isMoreInfoShown) {
        this.showCopyClipboardHint();
      }
    },
    taskErrorToShow: function () {
      if (this.taskErrorToShow) {
        // more info is initially hidden
        this.isMoreInfoShown = false;
      }
    },
  },
  methods: {
    taskErrorModalHidden() {
      this.$emit("hide");
    },
    getTaskError() {
      if (this.taskErrorToShow && this.taskErrorToShow.subTasks) {
        return this.taskErrorToShow.result.error;
      } else {
        return "";
      }
    },
    showCopyClipboardHint() {
      setTimeout(() => {
        //// TODO FIX
        // const isCopyClipboardHintShown = this.getFromStorage(
        //   "isCopyClipboardHintShown"
        // );
        //
        // if (!isCopyClipboardHintShown) {
        //   this.isCopyClipboardHintShown = true;
        //   this.saveToStorage("isCopyClipboardHintShown", true);
        // }
      }, 400);
    },
    onCopy() {
      this.justCopied = true;

      setTimeout(() => {
        this.justCopied = false;
      }, 3000);
    },
    onCopyError() {
      console.error("cannot copy to clipboard");
    },
    goToLogs() {
      this.$router.push(
        "/system-logs" + this.taskErrorToShow.context.extra.logs.path
      );
      this.$emit("hide");
    },
    onPrimaryClick() {
      if (this.isLogAvailable) {
        this.goToLogs();
      } else {
        this.$emit("hide");
      }
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";

.task-hierarchy {
  margin-left: -0.8rem;
}

.hint-copy-to-clipboard {
  top: 1.7rem;
  right: 2.3rem;
  float: right;
}
</style>

<style lang="scss">
@import "../../styles/carbon-utils";

// global styles

.task-error-modal .bx--modal-content {
  margin-top: $spacing-05;
}

.task-error-modal .bx--inline-notification__icon {
  margin-top: 0;
  margin-right: $spacing-02;
  fill: $support-04; // info color (low contrast)
}

.task-error-modal .bx--inline-notification__icon.aborted,
.task-error-modal .bx--inline-notification__icon.validation-failed {
  fill: $danger-01; // error color (low contrast)
}

.task-error-modal .bx--inline-notification__icon.completed {
  fill: $support-02; // success color (low contrast)
}

.task-error-modal .bx--inline-notification__icon.pending {
  fill: $support-03; // warning color (low contrast)
}

.task-error-modal .task-item {
  display: flex;
  align-items: center;
  padding-top: 0.5rem;
  padding-bottom: 0.5rem;
}
</style>

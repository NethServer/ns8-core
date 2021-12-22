<template>
  <cv-modal
    size="default"
    :visible="!!taskErrorToShow"
    @modal-hidden="taskErrorModalHidden"
    class="task-error-modal"
  >
    <template slot="title">{{ getTaskTitle(taskErrorToShow) }}</template>
    <template slot="content" v-if="taskErrorToShow">
      <div class="task-item">
        <component
          :is="getTaskIcon(taskErrorToShow)"
          :class="`bx--inline-notification__icon ${taskErrorToShow.status} task-icon`"
        />
        <span v-html="getTaskStatusDescription(taskErrorToShow, true)"></span>
      </div>
      <div v-if="subTasks.length">
        <TaskHierarchy :subTasks="subTasks" class="task-hierarchy" />
      </div>
      <div @click="showCopyClipboardHint">
        <cv-accordion ref="accordion" class="task-more-info">
          <cv-accordion-item :open="toggleAccordion[0]">
            <template slot="title">{{ $t("common.more_info") }}</template>
            <template slot="content">
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
              <div class="code-snippet-wrapper">
                <NsCodeSnippet
                  :copyTooltip="$t('common.copy_to_clipboard')"
                  :copy-feedback="$t('common.copied_to_clipboard')"
                  :feedback-aria-label="$t('common.copied_to_clipboard')"
                  :wrap-text="true"
                  :moreText="$t('common.show_more')"
                  :lessText="$t('common.show_less')"
                  light
                  expanded
                  >{{ taskErrorToShow }}</NsCodeSnippet
                >
              </div>
            </template>
          </cv-accordion-item>
        </cv-accordion>
      </div>
    </template>
    <template slot="secondary-button">{{ $t("common.close") }}</template>
  </cv-modal>
</template>

<script>
import { mapState, mapActions } from "vuex";
import TaskHierarchy from "@/components/TaskHierarchy";
import { TaskService, StorageService } from "@nethserver/ns8-ui-lib";

export default {
  name: "TaskErrorModal",
  components: { TaskHierarchy },
  mixins: [TaskService, StorageService],
  props: {},
  data() {
    return {
      isCopyClipboardHintShown: false,
    };
  },
  computed: {
    ...mapState(["taskErrorToShow"]),
    subTasks: function () {
      if (this.taskErrorToShow && this.taskErrorToShow.subTasks) {
        return this.taskErrorToShow.subTasks;
      } else {
        return [];
      }
    },
  },
  methods: {
    ...mapActions(["setTaskErrorToShowInStore"]),
    taskErrorModalHidden() {
      // use a delay to show a smooth animation for modal closing
      setTimeout(() => this.setTaskErrorToShowInStore(null), 300);
    },
    toggleAccordion(ev) {
      this.$refs.accordion.state.map(
        (item, index) => index === ev.changedIndex
      );
    },
    showCopyClipboardHint() {
      setTimeout(() => {
        const isCopyClipboardHintShown = this.getFromStorage(
          "isCopyClipboardHintShown"
        );

        if (!isCopyClipboardHintShown) {
          this.isCopyClipboardHintShown = true;
          this.saveToStorage("isCopyClipboardHintShown", true);
        }
      }, 400);
    },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";

.code-snippet-wrapper {
  max-height: 11rem;
}

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
@import "../styles/carbon-utils";

// global styles

.task-error-modal .bx--modal-content {
  margin-top: $spacing-05;
}

.task-more-info {
  margin-top: $spacing-04;
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

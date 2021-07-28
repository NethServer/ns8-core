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
      <cv-accordion ref="accordion" class="task-more-info">
        <cv-accordion-item :open="toggleAccordion[0]">
          <template slot="title">{{ $t("common.more_info") }}</template>
          <template slot="content">
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
    </template>
    <template slot="secondary-button">{{ $t("common.close") }}</template>
  </cv-modal>
</template>

<script>
import { mapState } from "vuex";
import { mapActions } from "vuex";
import TaskHierarchy from "@/components/TaskHierarchy";
import { TaskService } from "andrelib"; ////

export default {
  name: "TaskErrorModal",
  components: { TaskHierarchy },
  mixins: [TaskService],
  props: {},
  data() {
    return {};
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
</style>

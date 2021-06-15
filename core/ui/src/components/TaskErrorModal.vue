<template>
  <cv-modal
    size="default"
    :visible="!!taskErrorToShow"
    @modal-hidden="taskErrorModalHidden"
  >
    <template slot="title">{{ getTaskTitle(taskErrorToShow) }}</template>
    <template slot="content" v-if="taskErrorToShow">
      <InlineNotification
        :kind="getTaskKind(taskErrorToShow)"
        :title="getTaskStatusDescription(taskErrorToShow, true)"
        low-contrast
      />
      <div v-if="subTasks.length">
        <TaskHierarchy :subTasks="subTasks" />
      </div>

      <cv-accordion ref="acc" class="technical-details">
        <cv-accordion-item :open="toggleAccordion[0]">
          <template slot="title">{{ $t("task.technical_details") }}</template>
          <template slot="content">
            <cv-code-snippet
              kind="multiline"
              :copy-feedback="$t('common.copied_to_clipboard')"
              :feedback-aria-label="$t('copied_to_clipboard')"
              :less-text="$t('common.show_less')"
              :more-text="$t('common.show_more')"
              :wrap-text="true"
              >{{ taskErrorToShow }}</cv-code-snippet
            >
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
import InlineNotification from "@/components/InlineNotification";
import TaskHierarchy from "@/components/TaskHierarchy";
import TaskService from "@/mixins/task";

export default {
  name: "TaskErrorModal",
  components: { InlineNotification, TaskHierarchy },
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
      this.$refs.acc.state.map((item, index) => index === ev.changedIndex);
    },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";
</style>

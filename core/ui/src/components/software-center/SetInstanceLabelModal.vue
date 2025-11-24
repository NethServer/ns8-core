<template>
  <NsModal
    size="default"
    :visible="visible"
    :isLoading="loading"
    @modal-hidden="onModalHidden"
    @primary-click="setInstanceLabel"
    :primary-button-disabled="inputLabelError ? true : false"
  >
    <template slot="title">
      {{ $t("software_center.edit_instance_label") }}
    </template>
    <template slot="content">
      <template v-if="currentInstance">
        <NsTextInput
          ref="inputLabel"
          :label="
            $t('software_center.instance_label') +
            ' (' +
            $t('common.optional') +
            ')'
          "
          v-model="inputLabel"
          :placeholder="$t('common.no_label')"
          :helper-text="$t('software_center.instance_label_tooltip')"
          :disabled="loading"
          :invalid-message="$t(inputLabelError)"
        />
        <div v-if="error.setInstanceLabel">
          <NsInlineNotification
            kind="error"
            :title="$t('action.set-instance-label')"
            :description="error.setInstanceLabel"
            :showCloseButton="false"
          />
        </div>
      </template>
    </template>
    <template slot="secondary-button">{{ $t("common.cancel") }}</template>
    <template slot="primary-button">{{ $t("common.save") }}</template>
  </NsModal>
</template>
<script>
import { UtilService, TaskService, IconService } from "@nethserver/ns8-ui-lib";
import to from "await-to-js";
export default {
  name: "SetInstanceLabelModal",
  mixins: [UtilService, TaskService, IconService],
  props: {
    visible: Boolean,
    currentInstance: Object,
    newInstanceLabel: String,
  },
  data() {
    return {
      inputLabel: this.newInstanceLabel || "",
      loading: false,
      error: {
        setInstanceLabel: "",
      },
    };
  },
  watch: {
    newInstanceLabel(newVal) {
      this.inputLabel = newVal || "";
    },
    visible(newVal) {
      // When modal is opened, reset inputLabel to prop value
      if (newVal) {
        this.inputLabel = this.newInstanceLabel || "";
      }
    },
  },
  computed: {
    inputLabelError() {
      if (this.inputLabel.length > 24) {
        return "software_center.instance_label_too_long";
      }
      const alphanumRegex = /^[a-zA-Z0-9 ]*$/;
      if (!alphanumRegex.test(this.inputLabel)) {
        return "software_center.instance_label_alphanum_only";
      }
      return "";
    },
  },
  methods: {
    async setInstanceLabel() {
      this.error.setInstanceLabel = "";
      this.loading = true;
      const taskAction = "set-name";
      const eventId = this.getUuid();

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.setInstanceLabelCompleted
      );
      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.setInstanceLabelAborted
      );

      const res = await to(
        this.createModuleTaskForApp(this.currentInstance.id, {
          action: taskAction,
          data: {
            name: this.inputLabel,
          },
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
            eventId,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.setInstanceLabel = this.getErrorMessage(err);
        return;
      }
      // emit event to close modal
      this.$emit("hide");
    },
    setInstanceLabelAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.setInstanceLabel = this.$t("error.generic_error");
      this.loading = false;
    },
    setInstanceLabelCompleted() {
      this.loading = false;
      this.$emit("hide");
      // update instance label in app drawer
      this.$root.$emit("reloadAppDrawer");
      // reload instances and highlight new instance
      this.$emit("setInstanceLabelCompleted");
    },
    onModalHidden() {
      this.clearErrors();
      // Reset inputLabel to prop value when modal is hidden (cancelled)
      this.inputLabel = this.newInstanceLabel || "";
      this.$emit("hide");
    },
  },
};
</script>

<template>
  <NsModal
    size="default"
    :visible="visible"
    :isLoading="loading.restartModule"
    @modal-hidden="onModalHidden"
    @primary-click="restartModule"
  >
    <template slot="title">
      {{
        $t("software_center.app_restart", {
          app: instanceToRestart ? instanceToRestart.id : "",
        })
      }}
    </template>
    <template slot="content">
      <cv-form @submit.prevent="setInstanceLabel">
        <NsInlineNotification
          kind="warning"
          :title="$t('common.please_read_carefully')"
          :description="$t('software_center.restart_module_warning')"
          :showCloseButton="false"
        />
        <div>
          {{
            $t("software_center.restart_app", {
              name: instanceToRestart ? instanceToRestart.id : "",
            })
          }}
        </div>
        <div v-if="error.restartModule">
          <NsInlineNotification
            kind="error"
            :title="$t('action.restart-module')"
            :description="error.restartModule"
            :showCloseButton="false"
          />
        </div>
      </cv-form>
    </template>
    <template slot="secondary-button">{{ $t("common.cancel") }}</template>
    <template slot="primary-button">{{
      $t("software_center.restart_instance")
    }}</template>
  </NsModal>
</template>

<script>
import { UtilService, TaskService, IconService } from "@nethserver/ns8-ui-lib";
import to from "await-to-js";
export default {
  name: "RestartModuleModal",
  mixins: [UtilService, TaskService, IconService],
  props: {
    visible: {
      type: Boolean,
      required: true,
    },
    instanceToRestart: {
      type: Object,
      default: null,
    },
  },
  data() {
    return { error: { restartModule: "" }, loading: { restartModule: false } };
  },
  methods: {
    async restartModule() {
      this.error.restartModule = "";
      this.loading.restartModule = true;
      const taskAction = "restart-module";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.restartModuleAborted
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.restartModuleCompleted
      );

      const res = await to(
        this.createNodeTask(this.instanceToRestart.node, {
          action: taskAction,
          data: {
            module_id: this.instanceToRestart.id,
          },
          extra: {
            title: this.$t("applications.restart_instance_name", {
              instance: this.instanceToRestart.ui_name
                ? this.instanceToRestart.ui_name
                : this.instanceToRestart.id,
            }),
            description: this.$t("applications.restarting"),
            eventId,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.restartModule = this.getErrorMessage(err);
        this.loading.restartModule = false;
        return;
      }
      // emit event to close modal
      this.$emit("hide");
    },
    restartModuleAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.restartModule = this.$t("error.generic_error");
      this.loading.restartModule = false;
    },
    restartModuleCompleted() {
      this.loading.restartModule = false;
      this.$emit("hide");
    },
    onModalHidden() {
      this.clearErrors();
      this.$emit("hide");
    },
  },
};
</script>

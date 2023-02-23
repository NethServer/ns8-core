<template>
  <NsModal
    size="default"
    :visible="isShown"
    @modal-hidden="$emit('hide')"
    @primary-click="downloadClusterConfigurationBackup"
    :primary-button-disabled="!password || !confirmRead"
  >
    <template slot="title">{{
      $t("backup.download_cluster_configuration_backup")
    }}</template>
    <template slot="content">
      <cv-form @submit.prevent="!password || !confirmRead">
        <NsInlineNotification
          kind="info"
          :title="$t('backup.cluster_backup_password_info')"
          :description="$t('backup.cluster_backup_password_info_description')"
          :showCloseButton="false"
        />
        <NsTextInput
            v-model="password"
            type="password"
            :label="$t('backup.password')"
            :passwordHideLabel="$t('password.hide_password')"
            :passwordShowLabel="$t('password.show_password')"
            class="confirm-password"
            ref="confirmPassword"
        />
        <NsCheckbox
          :label="$t('backup.i_have_stored_encryption_password')"
          v-model="confirmRead"
          value="confirmRead"
          class="mg-top-lg"
        />
      </cv-form>
    </template>
    <template slot="secondary-button">{{ $t("common.cancel") }}</template>
    <template slot="primary-button">{{ $t("backup.download") }}</template>
  </NsModal>
</template>

<script>
import { UtilService, TaskService, IconService, DateTimeService} from "@nethserver/ns8-ui-lib";
import to from "await-to-js";

export default {
  name: "DownloadBackupModal",
  mixins: [UtilService, TaskService, IconService, DateTimeService],
  props: {
    isShown: {
      type: Boolean,
      default: true,
    },
  },
  data() {
    return {
      name: "",
      password: "",
      confirmRead: false,
      loading: {
        downloadClusterBackup: false
      }
    };
  },
  watch: {
    isShown: function () {
      if (this.isShown) {
        this.password = "";
        this.confirmRead = false;
      }
    },
  },
  methods: {
    async downloadClusterConfigurationBackup() {
      this.loading.downloadClusterBackup = true;
      const taskAction = "download-cluster-backup";

      // register to task error
      this.$root.$off(taskAction + "-aborted");
      this.$root.$once(
        taskAction + "-aborted",
        this.downloadClusterBackupAborted
      );

      // register to task completion
      this.$root.$off(taskAction + "-completed");
      this.$root.$once(
        taskAction + "-completed",
        this.downloadClusterBackupCompleted
      );

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: { password: this.password },
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.downloadClusterBackup = this.getErrorMessage(err);
        return;
      }
    },
    downloadClusterBackupAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.loading.downloadClusterBackup = false;
    },
    downloadClusterBackupCompleted(taskContext, taskResult) {
      this.loading.downloadClusterBackup = false;
      const downloadUrl = `${window.location.protocol}//${window.location.hostname}/cluster-admin/backup/${taskResult.output.path}`;

      const fileName =
        "cluster-configuration-backup " +
        this.formatDate(new Date(), "yyyy-MM-dd HH.mm") +
        ".json.gz.gpg";

      this.axios({
        url: downloadUrl,
        method: "GET",
        responseType: "blob",
      }).then((response) => {
        const url = window.URL.createObjectURL(new Blob([response.data]));
        const link = document.createElement("a");
        link.href = url;
        link.setAttribute("download", fileName);
        document.body.appendChild(link);
        link.click();
      });
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";
</style>

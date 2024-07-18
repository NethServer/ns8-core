<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <NsModal
    size="default"
    :visible="isShown"
    @modal-hidden="$emit('hide')"
    @primary-click="setBackupPassword"
    :primary-button-disabled="!clusterPassword || !confirmRead"
  >
    <template slot="title">{{
      $t("backup.set_cluster_backup_password")
    }}</template>
    <template slot="content">
      <cv-form @submit.prevent="!clusterPassword || !confirmRead">
        <NsInlineNotification
          kind="info"
          :description="$t('backup.cluster_backup_password_info_description')"
          :showCloseButton="false"
        />
        <NsTextInput
          v-model="clusterPassword"
          type="password"
          :label="$t('backup.password')"
          :passwordHideLabel="$t('password.hide_password')"
          :passwordShowLabel="$t('password.show_password')"
          data-modal-primary-focus
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
    <template slot="primary-button">{{
      $t("backup.set_cluster_backup_password")
    }}</template>
  </NsModal>
</template>

<script>
import {
  UtilService,
  TaskService,
  IconService,
  DateTimeService,
} from "@nethserver/ns8-ui-lib";
import to from "await-to-js";

export default {
  name: "PasswordBackupModal",
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
      confirmRead: false,
      loading: {
        setBackupPassword: false,
      },
      clusterPassword: "",
      error: {
        setBackupPassword: "",
      },
    };
  },
  watch: {
    isShown: function () {
      if (this.isShown) {
        this.confirmRead = false;
        this.clusterPassword = "";
      }
    },
  },
  methods: {
    async setBackupPassword() {
      this.loading.downloadClusterBackup = true;
      const taskAction = "download-cluster-backup";

      // register to task error
      this.$root.$off(taskAction + "-aborted");
      this.$root.$once(taskAction + "-aborted", this.setBackupPasswordAborted);

      // register to task completion
      this.$root.$off(taskAction + "-completed");
      this.$root.$once(
        taskAction + "-completed",
        this.setBackupPasswordCompleted
      );

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: { password: this.clusterPassword },
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
          },
        })
      );
      const err = res[0];
      this.$emit("hide");

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.setBackupPassword = this.getErrorMessage(err);
        return;
      }
    },
    setBackupPasswordAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.loading.setBackupPassword = false;
    },
    setBackupPasswordCompleted() {
      this.loading.setBackupPassword = false;
      this.$emit("password-set");
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";
</style>

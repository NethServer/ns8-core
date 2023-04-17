<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <NsModal
    size="default"
    :visible="isShown"
    :primary-button-disabled="loading.alterUser"
    @modal-hidden="onModalHidden"
    @primary-click="alterUser"
  >
    <template v-if="admin" slot="title">{{
      $t("settings_cluster_admins.edit_display_name_for_admin", {
        admin: admin.user,
      })
    }}</template>
    <template slot="content">
      <cv-form @submit.prevent="alterUser">
        <NsTextInput
          :label="
            $t('settings_cluster_admins.display_name') +
            ' (' +
            $t('common.optional') +
            ')'
          "
          :placeholder="$t('settings_cluster_admins.no_display_name')"
          v-model.trim="newDisplayName"
          :invalid-message="error.display_name"
          :disabled="loading.alterUser"
          data-modal-primary-focus
          ref="display_name"
        />
        <div v-if="error.alterUser">
          <NsInlineNotification
            kind="error"
            :title="$t('action.alter-user')"
            :description="error.alterUser"
            :showCloseButton="false"
          />
        </div>
      </cv-form>
    </template>
    <template slot="secondary-button">{{ $t("common.cancel") }}</template>
    <template slot="primary-button">{{
      $t("settings_cluster_admins.edit_display_name")
    }}</template>
  </NsModal>
</template>

<script>
import to from "await-to-js";
import { UtilService, TaskService } from "@nethserver/ns8-ui-lib";

export default {
  name: "EditAdminDisplayNameModal",
  mixins: [UtilService, TaskService],
  props: {
    isShown: Boolean,
    admin: {
      type: [Object, null],
    },
  },
  data() {
    return {
      newDisplayName: "",
      loading: {
        alterUser: false,
      },
      error: {
        alterUser: "",
        display_name: "",
      },
    };
  },
  watch: {
    isShown: function () {
      if (this.isShown) {
        this.newDisplayName = this.admin.display_name;
      }
    },
  },
  methods: {
    onModalHidden() {
      this.clearErrors();
      this.$emit("hide");
    },
    async alterUser() {
      // no validation needed
      this.loading.alterUser = true;
      this.error.alterUser = "";
      const taskAction = "alter-user";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.alterUserAborted
      );

      // register to task validation
      this.$root.$once(
        `${taskAction}-validation-ok-${eventId}`,
        this.alterUserValidationOk
      );
      this.$root.$once(
        `${taskAction}-validation-failed-${eventId}`,
        this.alterUserValidationFailed
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.alterUserCompleted
      );

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            user: this.admin.user,
            set: { display_name: this.newDisplayName },
            revoke: [],
            grant: [],
          },
          extra: {
            title: this.$t(
              "settings_cluster_admins.edit_display_name_for_admin",
              {
                admin: this.admin.user,
              }
            ),
            description: this.$t("common.processing"),
            eventId,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.alterUser = this.getErrorMessage(err);
        this.loading.alterUser = false;
        return;
      }
    },
    alterUserAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.loading.alterUser = false;

      // hide modal so that user can see error notification
      this.$emit("hide");
    },
    alterUserValidationOk() {
      this.loading.alterUser = false;

      // hide modal after validation
      this.$emit("hide");
    },
    alterUserValidationFailed(validationErrors) {
      this.loading.alterUser = false;
      let focusAlreadySet = false;

      for (const validationError of validationErrors) {
        const param = validationError.parameter;

        // set i18n error message
        this.error[param] = this.getI18nStringWithFallback(
          "settings_cluster_admins." + validationError.error,
          "error." + validationError.error
        );

        if (!focusAlreadySet) {
          this.focusElement(param);
          focusAlreadySet = true;
        }
      }
    },
    alterUserCompleted() {
      this.loading.alterUser = false;

      // reload admins
      this.$emit("displayNameUpdated");
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";
</style>

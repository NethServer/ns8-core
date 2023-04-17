<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <NsModal
    size="default"
    :visible="isShown"
    @modal-hidden="onModalHidden"
    @primary-click="alterUser"
    :primary-button-disabled="loading.alterUser"
  >
    <template v-if="admin" slot="title">{{
      $t("settings_cluster_admins.change_password_for_admin", {
        admin: admin.user,
      })
    }}</template>
    <template slot="content">
      <cv-form @submit.prevent="alterUser">
        <NsPasswordInput
          :newPasswordLabel="$t('password.new_password')"
          :confirmPasswordLabel="$t('password.re_enter_new_password')"
          v-model="newPassword"
          @passwordValidation="onPasswordValidation"
          :newPasswordInvalidMessage="$t(error.newPassword)"
          :confirmPasswordInvalidMessage="$t(error.confirmPassword)"
          :passwordHideLabel="$t('password.hide_password')"
          :passwordShowLabel="$t('password.show_password')"
          :lengthLabel="$t('password.long_enough')"
          :lowercaseLabel="$t('password.lowercase_letter')"
          :uppercaseLabel="$t('password.uppercase_letter')"
          :numberLabel="$t('password.number')"
          :symbolLabel="$t('password.symbol')"
          :equalLabel="$t('password.equal')"
          :focus="focusPasswordField"
          :clearConfirmPasswordCommand="clearConfirmPasswordCommand"
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
      $t("settings_cluster_admins.change_password")
    }}</template>
  </NsModal>
</template>

<script>
import { UtilService, TaskService, IconService } from "@nethserver/ns8-ui-lib";
import to from "await-to-js";

export default {
  name: "ChangeAdminPasswordModal",
  mixins: [UtilService, TaskService, IconService],
  props: {
    isShown: Boolean,
    admin: { type: [Object, null] },
  },
  data() {
    return {
      newPassword: "",
      passwordValidation: null,
      focusPasswordField: { element: "" },
      clearConfirmPasswordCommand: 0,
      loading: {
        alterUser: false,
      },
      error: {
        alterUser: "",
        newPassword: "",
        confirmPassword: "",
      },
    };
  },
  watch: {
    isShown: function () {
      if (this.isShown) {
        // clear password fields
        this.newPassword = "";
        this.clearConfirmPasswordCommand++;

        setTimeout(() => {
          this.focusPasswordField = { element: "newPassword" };
        }, 400);
      }
    },
  },
  methods: {
    validateAlterUser() {
      this.clearErrors();
      let isValidationOk = true;

      // password validation

      if (!this.newPassword) {
        this.error.newPassword = this.$t("common.required");

        if (isValidationOk) {
          this.focusPasswordField = { element: "newPassword" };
          isValidationOk = false;
        }
      } else {
        if (this.currentPassword === this.newPassword) {
          if (!this.error.newPassword) {
            this.error.newPassword = this.$t(
              "password.old_new_passwords_must_be_different"
            );
          }

          if (isValidationOk) {
            this.focusPasswordField = { element: "newPassword" };
            isValidationOk = false;
          }
        }

        if (
          !this.passwordValidation.isLengthOk ||
          !this.passwordValidation.isLowercaseOk ||
          !this.passwordValidation.isUppercaseOk ||
          !this.passwordValidation.isNumberOk ||
          !this.passwordValidation.isSymbolOk
        ) {
          if (!this.error.newPassword) {
            this.error.newPassword = this.$t("password.password_not_secure");
          }

          if (isValidationOk) {
            this.focusPasswordField = { element: "newPassword" };
            isValidationOk = false;
          }
        }

        if (!this.passwordValidation.isEqualOk) {
          if (!this.error.newPassword) {
            this.error.newPassword = this.$t("password.passwords_do_not_match");
          }

          if (!this.error.confirmPassword) {
            this.error.confirmPassword = this.$t(
              "password.passwords_do_not_match"
            );
          }

          if (isValidationOk) {
            this.focusPasswordField = { element: "confirmPassword" };
            isValidationOk = false;
          }
        }
      }
      return isValidationOk;
    },
    async alterUser() {
      if (!this.validateAlterUser()) {
        return;
      }
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
            set: { password: this.newPassword },
            revoke: [],
            grant: [],
          },
          extra: {
            title: this.$t(
              "settings_cluster_admins.change_password_for_admin",
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
        this.error[param] = this.$t(
          "settings_cluster_admins." + validationError.error
        );

        if (!focusAlreadySet) {
          this.focusElement(param);
          focusAlreadySet = true;
        }
      }
    },
    alterUserCompleted() {
      this.loading.alterUser = false;
    },
    onModalHidden() {
      this.clearErrors();
      this.$emit("hide");
    },
    onPasswordValidation(passwordValidation) {
      this.passwordValidation = passwordValidation;
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";
</style>

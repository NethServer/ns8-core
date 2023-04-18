<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <NsModal
    size="default"
    :visible="isShown"
    :primary-button-disabled="loading.addUser"
    :isLoading="loading.addUser"
    @modal-hidden="onModalHidden"
    @primary-click="addUser"
    class="no-pad-modal"
  >
    <template slot="title">{{
      $t("settings_cluster_admins.create_admin")
    }}</template>
    <template slot="content">
      <cv-form @submit.prevent="addUser">
        <NsTextInput
          v-model.trim="username"
          :label="$t('settings_cluster_admins.username')"
          :invalid-message="error.user"
          :disabled="loading.addUser"
          data-modal-primary-focus
          ref="user"
        />
        <NsTextInput
          v-model.trim="displayName"
          :label="
            $t('settings_cluster_admins.display_name') +
            ' (' +
            $t('common.optional') +
            ')'
          "
          :invalid-message="error.display_name"
          :disabled="loading.addUser"
          ref="display_name"
        />
        <NsPasswordInput
          :newPasswordLabel="$t('password.password')"
          :confirmPasswordLabel="$t('password.re_enter_password')"
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
        <!-- need to wrap error notification inside a div: custom elements like NsInlineNotification don't have scrollIntoView() function -->
        <div ref="addUserError">
          <NsInlineNotification
            v-if="error.addUser"
            kind="error"
            :title="$t('settings_cluster_admins.create_admin')"
            :description="error.addUser"
            :showCloseButton="false"
          />
        </div>
      </cv-form>
    </template>
    <template slot="secondary-button">{{ $t("common.cancel") }}</template>
    <template slot="primary-button">{{
      $t("settings_cluster_admins.create_admin")
    }}</template>
  </NsModal>
</template>

<script>
import { UtilService, TaskService, IconService } from "@nethserver/ns8-ui-lib";
import to from "await-to-js";

export default {
  name: "CreateAdminModal",
  mixins: [UtilService, TaskService, IconService],
  props: {
    isShown: Boolean,
    admin: { type: [Object, null] },
  },
  data() {
    return {
      username: "",
      displayName: "",
      selectedGroups: [],
      newPassword: "",
      passwordValidation: null,
      focusPasswordField: { element: "" },
      clearConfirmPasswordCommand: 0,
      loading: {
        addUser: false,
      },
      error: {
        addUser: "",
        user: "",
        display_name: "",
        newPassword: "",
        confirmPassword: "",
        groups: "",
      },
    };
  },
  watch: {
    isShown: function () {
      if (this.isShown) {
        this.clearErrors();

        this.username = "";
        this.displayName = "";
        this.newPassword = "";
        this.clearConfirmPasswordCommand++;
      }
    },
    "error.addUser": function () {
      if (this.error.addUser) {
        // scroll to notification error

        this.$nextTick(() => {
          const el = this.$refs.addUserError;
          this.scrollToElement(el);
        });
      }
    },
  },
  methods: {
    validateAddUser() {
      this.clearErrors();
      let isValidationOk = true;

      // username

      if (!this.username) {
        this.error.user = this.$t("common.required");

        if (isValidationOk) {
          this.focusElement("user");
          isValidationOk = false;
        }
      }

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
    async addUser() {
      if (!this.validateAddUser()) {
        return;
      }
      this.loading.addUser = true;
      this.error.addUser = "";
      const taskAction = "add-user";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(`${taskAction}-aborted-${eventId}`, this.addUserAborted);

      // register to task validation
      this.$root.$once(
        `${taskAction}-validation-ok-${eventId}`,
        this.addUserValidationOk
      );
      this.$root.$once(
        `${taskAction}-validation-failed-${eventId}`,
        this.addUserValidationFailed
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.addUserCompleted
      );

      const [errPasswordHash, passwordHash] = await to(
        this.getSha256(this.newPassword)
      );

      if (errPasswordHash) {
        console.error(`error getting password sha256`, errPasswordHash);
        this.error.addUser = this.getErrorMessage(errPasswordHash);
        this.loading.addUser = false;
        return;
      }

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            user: this.username,
            set: {
              display_name: this.displayName,
            },
            password_hash: passwordHash,
            grant: [
              {
                role: "owner",
                on: "*",
              },
            ],
          },
          extra: {
            title: this.$t("settings_cluster_admins.create_admin_admin", {
              admin: this.username,
            }),
            description: this.$t("common.processing"),
            eventId,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.addUser = this.getErrorMessage(err);
        this.loading.addUser = false;
        return;
      }
    },
    addUserAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.loading.addUser = false;

      // hide modal so that user can see error notification
      this.$emit("hide");
    },
    addUserValidationOk() {
      this.loading.addUser = false;

      // hide modal after validation
      this.$emit("hide");
    },
    addUserValidationFailed(validationErrors) {
      this.loading.addUser = false;
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
    addUserCompleted() {
      this.loading.addUser = false;

      // reload admins
      this.$emit("adminCreated");
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

<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <NsModal
    size="default"
    :visible="isShown"
    @modal-hidden="onModalHidden"
    @primary-click="changeUserPassword"
    :primary-button-disabled="loading.alterUser || loading.listPasswordPolicy"
    class="no-pad-modal"
  >
    <template v-if="user" slot="title">{{
      $t("domain_users.change_password_for_user", { user: user.user })
    }}</template>
    <template slot="content">
      <cv-form @submit.prevent="changeUserPassword">
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
          :minLength="
            policy.strength.enforced ? policy.strength.password_min_length : 0
          "
          :checkComplexity="
            policy.strength.enforced ? policy.strength.complexity_check : false
          "
        />
        <div v-if="error.alterUser">
          <NsInlineNotification
            kind="error"
            :title="$t('action.alter-user')"
            :description="error.alterUser"
            :showCloseButton="false"
          />
        </div>
        <NsInlineNotification
          v-if="error.listPasswordPolicy"
          kind="error"
          :title="$t('action.get-password-policy')"
          :description="error.listPasswordPolicy"
          :showCloseButton="false"
        />
      </cv-form>
    </template>
    <template slot="secondary-button">{{ $t("common.cancel") }}</template>
    <template slot="primary-button">{{
      $t("domain_users.change_password")
    }}</template>
  </NsModal>
</template>

<script>
import { UtilService, TaskService, IconService } from "@nethserver/ns8-ui-lib";
import to from "await-to-js";

export default {
  name: "ChangeUserPasswordModal",
  mixins: [UtilService, TaskService, IconService],
  props: {
    isShown: Boolean,
    domain: { type: Object },
    user: { type: [Object, null] },
  },
  data() {
    return {
      newPassword: "",
      passwordValidation: null,
      focusPasswordField: { element: "" },
      clearConfirmPasswordCommand: 0,
      policy: {
        strength: {
          complexity_check: false,
          enforced: false,
          password_min_length: 8,
        },
      },
      loading: {
        alterUser: false,
        listPasswordPolicy: false,
      },
      error: {
        alterUser: "",
        newPassword: "",
        confirmPassword: "",
        listPasswordPolicy: "",
      },
    };
  },
  computed: {
    mainProvider() {
      return this.domain.providers[0].id;
    },
  },
  watch: {
    isShown: function () {
      if (this.isShown) {
        // clear password fields
        this.newPassword = "";
        this.clearConfirmPasswordCommand++;
        this.listPasswordPolicy();
        setTimeout(() => {
          this.focusPasswordField = { element: "newPassword" };
        }, 400);
      }
    },
  },
  methods: {
    async listPasswordPolicy() {
      this.loading.listPasswordPolicy = true;
      this.error.listPasswordPolicy = "";
      const taskAction = "get-password-policy";
      const eventId = this.getUuid();
      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.listPasswordPolicyCompleted
      );
      const res = await to(
        this.createModuleTaskForApp(this.mainProvider, {
          action: taskAction,
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
        this.error.listPasswordPolicy = this.getErrorMessage(err);
        return;
      }
    },
    listPasswordPolicyCompleted(taskContext, taskResult) {
      const config = taskResult.output;
      this.policy.strength.enforced = config.strength.enforced;
      this.policy.strength.complexity_check = config.strength.complexity_check;
      this.policy.strength.password_min_length =
        config.strength.password_min_length.toString();

      this.loading.listPasswordPolicy = false;
    },
    validateChangeUserPassword() {
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
    async changeUserPassword() {
      if (!this.validateChangeUserPassword()) {
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
        this.createModuleTaskForApp(this.mainProvider, {
          action: taskAction,
          data: {
            user: this.user.user,
            password: this.newPassword,
          },
          extra: {
            title: this.$t("domain_users.change_password_for_user", {
              user: this.user.user,
            }),
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
        this.error[param] = this.$t("domain_users." + validationError.error);

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

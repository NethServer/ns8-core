<template>
  <NsModal
    size="default"
    :visible="isShown"
    :primary-button-disabled="isPrimaryButtonDisabled"
    @modal-hidden="onModalHidden"
    @primary-click="createOrEditUser"
    class="no-pad-modal"
  >
    <template slot="title">{{
      isEditing ? $t("domain_users.edit_user") : $t("domain_users.create_user")
    }}</template>
    <template slot="content">
      <cv-form @submit.prevent="createOrEditUser">
        <NsTextInput
          v-model.trim="username"
          :label="$t('domain_users.username')"
          :invalid-message="error.username"
          :disabled="isEditing || loading.addUser || loading.alterUser"
          data-modal-primary-focus
          ref="username"
        />
        <NsTextInput
          v-model.trim="fullName"
          :label="$t('domain_users.full_name')"
          :invalid-message="error.fullName"
          :disabled="loading.addUser || loading.alterUser"
          ref="fullName"
        />
        <NsPasswordInput
          v-if="!isEditing"
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
        <cv-multi-select
          v-model="selectedGroups"
          :options="allGroupsForSelect"
          :title="
            $t('domain_users.groups') + ' (' + $t('common.optional') + ')'
          "
          :label="
            allGroups.length
              ? $t('domain_users.select_groups')
              : $t('domain_users.no_group')
          "
          :helper-text="groupsHelperText"
          :filterable="!!allGroups.length"
          :auto-filter="true"
          :auto-highlight="true"
          :disabled="!allGroups.length || loading.getDomainUser"
          :class="{ 'mg-bottom-14': isEditing }"
          ref="groups"
        >
        </cv-multi-select>
        <NsInlineNotification
          v-if="error.addUser"
          kind="error"
          :title="$t('action.add-user')"
          :description="error.addUser"
          :showCloseButton="false"
        />
        <NsInlineNotification
          v-if="error.alterUser"
          kind="error"
          :title="$t('action.alter-user')"
          :description="error.alterUser"
          :showCloseButton="false"
        />
      </cv-form>
    </template>
    <template slot="secondary-button">{{ $t("common.cancel") }}</template>
    <template slot="primary-button">{{
      isEditing ? $t("domain_users.edit_user") : $t("domain_users.create_user")
    }}</template>
  </NsModal>
</template>

<script>
import { UtilService, TaskService, IconService } from "@nethserver/ns8-ui-lib";
import to from "await-to-js";

export default {
  name: "CreateOrEditUserModal",
  mixins: [UtilService, TaskService, IconService],
  props: {
    isShown: Boolean,
    isEditing: {
      type: Boolean,
      default: false,
    },
    domain: { type: Object },
    user: { type: [Object, null] },
    allGroups: { type: Array, required: true },
    provider: { type: String },
  },
  data() {
    return {
      username: "",
      fullName: "",
      selectedGroups: [],
      newPassword: "",
      passwordValidation: null,
      focusPasswordField: { element: "" },
      clearConfirmPasswordCommand: 0,
      loading: {
        addUser: false,
        alterUser: false,
        getDomainUser: false,
      },
      error: {
        addUser: "",
        alterUser: "",
        getDomainUser: "",
        username: "",
        fullName: "",
        newPassword: "",
        confirmPassword: "",
      },
    };
  },
  computed: {
    allGroupsForSelect() {
      return this.allGroups.map((group) => {
        return {
          value: group.group,
          label: group.group,
          name: group.group,
        };
      });
    },
    groupsHelperText() {
      if (!this.selectedGroups.length) {
        return "";
      } else {
        return (
          this.$t("common.selected") + ": " + this.selectedGroups.join(", ")
        );
      }
    },
    isPrimaryButtonDisabled() {
      return (
        this.loading.addUser ||
        this.loading.alterUser ||
        this.loading.getDomainUser
      );
    },
  },
  watch: {
    isShown: function () {
      if (this.isShown) {
        if (!this.isEditing) {
          // create user
          this.username = "";
          this.fullName = "";

          // clear password fields
          this.newPassword = "";
          this.clearConfirmPasswordCommand++;
        } else {
          // edit user
          this.username = this.user.user;
          this.fullName = this.user.full_name;
          this.getDomainUser();
        }
        this.selectedGroups = [];
      }
    },
  },
  methods: {
    createOrEditUser() {
      if (this.isEditing) {
        this.alterUser();
      } else {
        this.addUser();
      }
    },
    validateAddUser() {
      this.clearErrors();
      let isValidationOk = true;

      if (!this.username) {
        this.error.username = this.$t("common.required");

        if (isValidationOk) {
          this.focusElement("username");
          isValidationOk = false;
        }
      }

      //// check if username is already taken

      ////check illegal characters an syntax

      if (!this.fullName) {
        this.error.fullName = this.$t("common.required");

        if (isValidationOk) {
          this.focusElement("fullName");
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
    addUser() {
      if (!this.validateAddUser()) {
        return;
      }

      console.log("validation ok"); ////

      // const taskAction = "..."; ////
    },
    validateAlterUser() {
      this.clearErrors();
      let isValidationOk = true;

      if (!this.fullName) {
        this.error.fullName = this.$t("common.required");

        if (isValidationOk) {
          this.focusElement("fullName");
          isValidationOk = false;
        }
      }
      return isValidationOk;
    },
    alterUser() {
      if (!this.validateAlterUser()) {
        return;
      }

      console.log("validation ok"); ////

      // const taskAction = "..."; ////
    },
    onModalHidden() {
      this.clearErrors();
      this.$emit("hide");
    },
    onPasswordValidation(passwordValidation) {
      this.passwordValidation = passwordValidation;
    },
    async getDomainUser() {
      this.loading.getDomainUser = true;
      this.error.getDomainUser = "";
      const taskAction = "get-domain-user";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.getDomainUserAborted
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.getDomainUserCompleted
      );

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            domain: this.domain.name,
            user: this.user.user,
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
        this.error.getDomainUser = this.getErrorMessage(err);
        return;
      }
    },
    getDomainUserAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.loading.getDomainUser = false;
      this.error.getDomainUser = this.$t("error.generic_error");
    },
    getDomainUserCompleted(taskContext, taskResult) {
      console.log("getDomainUserCompleted", taskResult.output); ////

      const groups = taskResult.output.user.groups;

      this.selectedGroups = groups.map((g) => g.group);

      console.log("selectedGroups", this.selectedGroups); ////
      this.loading.getDomainUser = false;
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";

.mg-bottom-14 {
  margin-bottom: 14rem;
}
</style>

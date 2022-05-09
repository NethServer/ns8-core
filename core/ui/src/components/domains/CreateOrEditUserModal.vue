<template>
  <NsModal
    size="default"
    :visible="isShown"
    :primary-button-disabled="isPrimaryButtonDisabled"
    :isLoading="loading.addUser || loading.alterUser"
    @modal-hidden="onModalHidden"
    @primary-click="createOrEditUser"
    class="no-pad-modal"
  >
    <template slot="title">{{
      isEditing ? $t("domain_users.edit_user") : $t("domain_users.create_user")
    }}</template>
    <template slot="content">
      <cv-form @submit.prevent="createOrEditUser">
        <NsInlineNotification
          v-if="error.getDomainUser"
          kind="error"
          :title="$t('action.get-domain-user')"
          :description="error.getDomainUser"
          :showCloseButton="false"
        />
        <NsTextInput
          v-model.trim="username"
          :label="$t('domain_users.username')"
          :invalid-message="error.user"
          :disabled="isEditing || loading.addUser || loading.alterUser"
          data-modal-primary-focus
          ref="user"
        />
        <NsTextInput
          v-model.trim="displayName"
          :label="$t('domain_users.display_name')"
          :invalid-message="error.displayName"
          :disabled="loading.addUser || loading.alterUser"
          ref="displayName"
        />
        <cv-multi-select
          v-model="selectedGroups"
          :options="allGroupsForSelect"
          :title="
            $t('domain_users.groups') + ' (' + $t('common.optional') + ')'
          "
          :label="selectGroupsLabel"
          :helper-text="groupsHelperText"
          :filterable="!!allGroups.length"
          :auto-filter="true"
          :auto-highlight="true"
          :invalid-message="error.groups"
          :disabled="
            !allGroups.length ||
            loading.getDomainUser ||
            loading.addUser ||
            loading.alterUser
          "
          :class="{ 'mg-bottom-14': isEditing }"
          ref="groups"
        >
        </cv-multi-select>
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
        alterUser: false,
        getDomainUser: false,
      },
      error: {
        addUser: "",
        alterUser: "",
        getDomainUser: "",
        user: "",
        displayName: "",
        newPassword: "",
        confirmPassword: "",
        groups: "",
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
        this.loading.getDomainUser ||
        !!this.error.getDomainUser
      );
    },
    mainProvider() {
      return this.domain.providers[0].id;
    },
    selectGroupsLabel() {
      if (this.loading.getDomainUser) {
        return this.$t("common.loading");
      } else if (this.allGroups.length) {
        return this.$t("domain_users.select_groups");
      } else {
        return this.$t("domain_users.no_group");
      }
    },
  },
  watch: {
    isShown: function () {
      if (this.isShown) {
        this.clearErrors();

        if (!this.isEditing) {
          // create user
          this.username = "";
          this.displayName = "";

          // clear password fields
          this.newPassword = "";
          this.clearConfirmPasswordCommand++;
        } else {
          // edit user
          this.username = this.user.user;
          this.displayName = this.user.display_name;
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

      const res = await to(
        this.createModuleTaskForApp(this.mainProvider, {
          action: taskAction,
          data: {
            user: this.username,
            display_name: this.displayName,
            password: this.newPassword,
            locked: false,
            groups: this.selectedGroups,
          },
          extra: {
            title: this.$t("domain_users.create_user_user", {
              user: this.username,
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
        this.error[param] = this.$t("domain_users." + validationError.error);

        if (!focusAlreadySet) {
          this.focusElement(param);
          focusAlreadySet = true;
        }
      }
    },
    addUserCompleted() {
      this.loading.addUser = false;

      // reload users
      this.$emit("reloadUsers");
    },
    async alterUser() {
      this.loading.alterUser = true;
      this.error.alterUser = "";
      const taskAction = "alter-user";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.alterUserAborted
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
            display_name: this.displayName,
            groups: this.selectedGroups,
          },
          extra: {
            title: this.$t("domain_users.edit_user_user", {
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

      // hide modal
      this.$emit("hide");
    },
    alterUserAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.loading.alterUser = false;
    },
    alterUserCompleted() {
      this.loading.alterUser = false;

      // reload users
      this.$emit("reloadUsers");
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
      this.error.getDomainUser = this.$t("error.generic_error");
      this.loading.getDomainUser = false;
    },
    getDomainUserCompleted(taskContext, taskResult) {
      const groups = taskResult.output.user.groups;
      this.selectedGroups = groups.map((g) => g.group);
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

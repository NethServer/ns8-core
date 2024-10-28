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
      <cv-form @submit.prevent="createOrEditUser" autocomplete="off">
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
          autocomplete="off"
          ref="user"
        />
        <NsTextInput
          v-model.trim="displayName"
          :label="$t('domain_users.display_name')"
          :invalid-message="error.display_name"
          :disabled="loading.addUser || loading.alterUser"
          autocomplete="off"
          ref="display_name"
        />
        <NsMultiSelect
          v-model="selectedGroups"
          :options="allGroupsForSelect"
          :title="
            $t('domain_users.groups') + ' (' + $t('common.optional') + ')'
          "
          :label="selectGroupsLabel"
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
          showSelectedItems
          :unselectAriaLabel="$t('common.unselect')"
          :clearSelectionAriaLabel="$t('common.clear_selection')"
          :clearFilterLabel="$t('common.clear_filter')"
          :selectedLabel="$t('common.selected_l')"
          :userInputLabel="$t('common.user_input_l')"
          :class="{ 'mg-bottom-14': isEditing }"
          ref="groups"
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
          :minLength="
            policy.strength.enforced ? policy.strength.password_min_length : 0
          "
          :checkComplexity="
            policy.strength.enforced ? policy.strength.complexity_check : false
          "
          autocomplete="new-password"
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
      isEditing ? $t("domain_users.edit_user") : $t("domain_users.create_user")
    }}</template>
  </NsModal>
</template>

<script>
import { ref, reactive, computed, watch, onMounted, onBeforeUnmount } from "vue";
import { useStore } from "vuex";
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
  setup(props, { emit }) {
    const store = useStore();
    const username = ref("");
    const displayName = ref("");
    const selectedGroups = ref([]);
    const newPassword = ref("");
    const passwordValidation = ref(null);
    const focusPasswordField = ref({ element: "" });
    const clearConfirmPasswordCommand = ref(0);
    const policy = reactive({
      strength: {
        complexity_check: false,
        enforced: false,
        password_min_length: 8,
      },
    });
    const loading = reactive({
      addUser: false,
      alterUser: false,
      getDomainUser: false,
      listPasswordPolicy: false,
    });
    const error = reactive({
      addUser: "",
      alterUser: "",
      getDomainUser: "",
      user: "",
      display_name: "",
      newPassword: "",
      confirmPassword: "",
      groups: "",
      listPasswordPolicy: "",
    });

    const allGroupsForSelect = computed(() => {
      return props.allGroups.map((group) => {
        return {
          value: group.group,
          label: group.group,
          name: group.group,
        };
      });
    });

    const isPrimaryButtonDisabled = computed(() => {
      return (
        loading.addUser ||
        loading.alterUser ||
        loading.getDomainUser ||
        loading.listPasswordPolicy ||
        !!error.getDomainUser
      );
    });

    const mainProvider = computed(() => {
      return props.domain.providers[0].id;
    });

    const selectGroupsLabel = computed(() => {
      if (loading.getDomainUser) {
        return this.$t("common.loading");
      } else if (props.allGroups.length) {
        return this.$t("domain_users.select_groups");
      } else {
        return this.$t("domain_users.no_group");
      }
    });

    watch(
      () => props.isShown,
      (newVal) => {
        if (newVal) {
          clearErrors();
          listPasswordPolicy();
          if (!props.isEditing) {
            // create user

            // clear password fields
            newPassword.value = "";
            clearConfirmPasswordCommand.value++;
          } else {
            // edit user
            username.value = props.user.user;
            displayName.value = props.user.display_name;
            selectedGroups.value = [];
            getDomainUser();
          }
        } else {
          // hiding modal

          if (props.isEditing) {
            clearFields();
          }
        }
      }
    );

    const listPasswordPolicy = async () => {
      loading.listPasswordPolicy = true;
      error.listPasswordPolicy = "";
      const taskAction = "get-password-policy";
      const eventId = getUuid();
      // register to task completion
      store.$root.$once(
        `${taskAction}-completed-${eventId}`,
        listPasswordPolicyCompleted
      );
      const res = await to(
        createModuleTaskForApp(mainProvider.value, {
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
        error.listPasswordPolicy = getErrorMessage(err);
        return;
      }
    };

    const listPasswordPolicyCompleted = (taskContext, taskResult) => {
      const config = taskResult.output;
      policy.strength.enforced = config.strength.enforced;
      policy.strength.complexity_check = config.strength.complexity_check;
      policy.strength.password_min_length =
        config.strength.password_min_length.toString();

      loading.listPasswordPolicy = false;
    };

    const createOrEditUser = () => {
      if (props.isEditing) {
        alterUser();
      } else {
        addUser();
      }
    };

    const validateAddUser = () => {
      clearErrors();
      let isValidationOk = true;

      if (!username.value) {
        error.user = this.$t("common.required");

        if (isValidationOk) {
          focusElement("user");
          isValidationOk = false;
        }
      }
      // displayName is required
      if (!displayName.value) {
        error.display_name = this.$t("common.required");

        if (isValidationOk) {
          focusElement("display_name");
          isValidationOk = false;
        }
      }

      // password validation

      if (!newPassword.value) {
        error.newPassword = this.$t("common.required");

        if (isValidationOk) {
          focusPasswordField.value = { element: "newPassword" };
          isValidationOk = false;
        }
      } else {
        if (currentPassword === newPassword.value) {
          if (!error.newPassword) {
            error.newPassword = this.$t(
              "password.old_new_passwords_must_be_different"
            );
          }

          if (isValidationOk) {
            focusPasswordField.value = { element: "newPassword" };
            isValidationOk = false;
          }
        }

        if (
          !passwordValidation.value.isLengthOk ||
          !passwordValidation.value.isLowercaseOk ||
          !passwordValidation.value.isUppercaseOk ||
          !passwordValidation.value.isNumberOk ||
          !passwordValidation.value.isSymbolOk
        ) {
          if (!error.newPassword) {
            error.newPassword = this.$t("password.password_not_secure");
          }

          if (isValidationOk) {
            focusPasswordField.value = { element: "newPassword" };
            isValidationOk = false;
          }
        }

        if (!passwordValidation.value.isEqualOk) {
          if (!error.newPassword) {
            error.newPassword = this.$t("password.passwords_do_not_match");
          }

          if (!error.confirmPassword) {
            error.confirmPassword = this.$t(
              "password.passwords_do_not_match"
            );
          }

          if (isValidationOk) {
            focusPasswordField.value = { element: "confirmPassword" };
            isValidationOk = false;
          }
        }
      }
      return isValidationOk;
    };

    const addUser = async () => {
      if (!validateAddUser()) {
        return;
      }
      loading.addUser = true;
      error.addUser = "";
      const taskAction = "add-user";
      const eventId = getUuid();

      // register to task error
      store.$root.$once(`${taskAction}-aborted-${eventId}`, addUserAborted);

      // register to task validation
      store.$root.$once(
        `${taskAction}-validation-ok-${eventId}`,
        addUserValidationOk
      );
      store.$root.$once(
        `${taskAction}-validation-failed-${eventId}`,
        addUserValidationFailed
      );

      // register to task completion
      store.$root.$once(
        `${taskAction}-completed-${eventId}`,
        addUserCompleted
      );

      const res = await to(
        createModuleTaskForApp(mainProvider.value, {
          action: taskAction,
          data: {
            user: username.value,
            display_name: displayName.value,
            password: newPassword.value,
            locked: false,
            groups: selectedGroups.value,
          },
          extra: {
            title: this.$t("domain_users.create_user_user", {
              user: username.value,
            }),
            description: this.$t("common.processing"),
            eventId,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        error.addUser = getErrorMessage(err);
        loading.addUser = false;
        return;
      }
    };

    const addUserAborted = (taskResult, taskContext) => {
      console.error(`${taskContext.action} aborted`, taskResult);
      loading.addUser = false;

      // hide modal so that user can see error notification
      emit("hide");
    };

    const addUserValidationOk = () => {
      loading.addUser = false;

      // hide modal after validation
      emit("hide");
    };

    const addUserValidationFailed = (validationErrors) => {
      loading.addUser = false;
      let focusAlreadySet = false;
      for (const validationError of validationErrors) {
        const param = validationError.parameter;
        // set i18n error message
        error[param] = this.$t("domain_users." + validationError.error, {
          tok: validationError.value,
        });

        if (!focusAlreadySet) {
          focusElement(param);
          focusAlreadySet = true;
        }
      }
    };

    const addUserCompleted = () => {
      clearFields();
      loading.addUser = false;

      // reload users
      emit("reloadUsers");
    };

    const alterUser = async () => {
      loading.alterUser = true;
      error.alterUser = "";
      const taskAction = "alter-user";
      const eventId = getUuid();

      // register to task error
      store.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        alterUserAborted
      );

      // register to task completion
      store.$root.$once(
        `${taskAction}-completed-${eventId}`,
        alterUserCompleted
      );

      const res = await to(
        createModuleTaskForApp(mainProvider.value, {
          action: taskAction,
          data: {
            user: props.user.user,
            display_name: displayName.value,
            groups: selectedGroups.value,
          },
          extra: {
            title: this.$t("domain_users.edit_user_user", {
              user: props.user.user,
            }),
            description: this.$t("common.processing"),
            eventId,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        error.alterUser = getErrorMessage(err);
        loading.alterUser = false;
        return;
      }

      // hide modal
      emit("hide");
    };

    const alterUserAborted = (taskResult, taskContext) => {
      console.error(`${taskContext.action} aborted`, taskResult);
      loading.alterUser = false;
    };

    const alterUserCompleted = () => {
      loading.alterUser = false;

      // reload users
      emit("reloadUsers");
    };

    const onModalHidden = () => {
      clearErrors();
      emit("hide");
    };

    const onPasswordValidation = (passwordValidation) => {
      passwordValidation.value = passwordValidation;
    };

    const getDomainUser = async () => {
      loading.getDomainUser = true;
      error.getDomainUser = "";
      const taskAction = "get-domain-user";
      const eventId = getUuid();

      // register to task error
      store.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        getDomainUserAborted
      );

      // register to task completion
      store.$root.$once(
        `${taskAction}-completed-${eventId}`,
        getDomainUserCompleted
      );

      const res = await to(
        createClusterTask({
          action: taskAction,
          data: {
            domain: props.domain.name,
            user: props.user.user,
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
        error.getDomainUser = getErrorMessage(err);
        return;
      }
    };

    const getDomainUserAborted = (taskResult, taskContext) => {
      console.error(`${taskContext.action} aborted`, taskResult);
      error.getDomainUser = this.$t("error.generic_error");
      loading.getDomainUser = false;
    };

    const getDomainUserCompleted = (taskContext, taskResult) => {
      const groups = taskResult.output.user.groups;
      selectedGroups.value = groups.map((g) => g.group);
      loading.getDomainUser = false;
    };

    const clearFields = () => {
      username.value = "";
      displayName.value = "";
      selectedGroups.value = [];

      // clear password fields
      newPassword.value = "";
      clearConfirmPasswordCommand.value++;
    };

    return {
      username,
      displayName,
      selectedGroups,
      newPassword,
      passwordValidation,
      focusPasswordField,
      clearConfirmPasswordCommand,
      policy,
      loading,
      error,
      allGroupsForSelect,
      isPrimaryButtonDisabled,
      mainProvider,
      selectGroupsLabel,
      listPasswordPolicy,
      listPasswordPolicyCompleted,
      createOrEditUser,
      validateAddUser,
      addUser,
      addUserAborted,
      addUserValidationOk,
      addUserValidationFailed,
      addUserCompleted,
      alterUser,
      alterUserAborted,
      alterUserCompleted,
      onModalHidden,
      onPasswordValidation,
      getDomainUser,
      getDomainUserAborted,
      getDomainUserCompleted,
      clearFields,
    };
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";

.mg-bottom-14 {
  margin-bottom: 14rem;
}
</style>

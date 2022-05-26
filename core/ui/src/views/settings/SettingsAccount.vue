<template>
  <div>
    <cv-grid fullWidth>
      <cv-row>
        <cv-column>
          <cv-breadcrumb
            aria-label="breadcrumb"
            :no-trailing-slash="true"
            class="breadcrumb"
          >
            <cv-breadcrumb-item>
              <cv-link to="/settings">{{ $t("settings.title") }}</cv-link>
            </cv-breadcrumb-item>
            <cv-breadcrumb-item>
              <span>{{ $t("settings_account.title") }}</span>
            </cv-breadcrumb-item>
          </cv-breadcrumb>
        </cv-column>
      </cv-row>
      <cv-row>
        <cv-column class="subpage-title">
          <h3>{{ $t("settings_account.title") }}</h3>
        </cv-column>
      </cv-row>
      <cv-row>
        <cv-column>
          <cv-tile light>
            <h6 class="mg-bottom-md">
              {{ $t("settings_account.change_password") }}
            </h6>
            <cv-form @submit.prevent="changeUserPassword">
              <NsTextInput
                :label="
                  $t('settings_account.current_user_password', {
                    user: currentUsername,
                  })
                "
                v-model="currentPassword"
                :invalid-message="$t(error.currentPassword)"
                type="password"
                ref="currentPassword"
              >
              </NsTextInput>
              <NsPasswordInput
                :newPasswordLabel="
                  $t('settings_account.new_user_password', {
                    user: currentUsername,
                  })
                "
                :confirmPasswordLabel="
                  $t('settings_account.new_user_password_confirm', {
                    user: currentUsername,
                  })
                "
                v-model="newPassword"
                @passwordValidation="onPasswordValidation"
                :newPasswordInvalidMessage="error.newPassword"
                :confirmPasswordInvalidMessage="error.confirmPassword"
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
                v-if="error.changeUserPassword"
                kind="error"
                :title="$t('action.change-user-password')"
                :description="error.changeUserPassword"
                :showCloseButton="false"
              />
              <NsButton
                kind="primary"
                :loading="loading.changeUserPassword"
                :disabled="loading.changeUserPassword"
                :icon="Password20"
                >{{ $t("settings_account.change_password") }}
              </NsButton>
            </cv-form>
          </cv-tile>
        </cv-column>
      </cv-row>
    </cv-grid>
  </div>
</template>

<script>
import to from "await-to-js";
import {
  QueryParamService,
  UtilService,
  TaskService,
  IconService,
} from "@nethserver/ns8-ui-lib";

export default {
  name: "SettingsAccount",
  mixins: [TaskService, UtilService, IconService, QueryParamService],
  pageTitle() {
    return this.$t("settings_account.title");
  },
  data() {
    return {
      currentUsername: "",
      currentPassword: "",
      newPassword: "",
      passwordValidation: null,
      focusPasswordField: { element: "" },
      clearConfirmPasswordCommand: 0,
      loading: {
        changeUserPassword: false,
      },
      error: {
        changeUserPassword: "",
        currentPassword: "",
        newPassword: "",
        confirmPassword: "",
      },
    };
  },
  created() {
    const loginInfo = this.getFromStorage("loginInfo");

    if (loginInfo) {
      this.currentUsername = loginInfo.username;
    }
  },
  beforeRouteEnter(to, from, next) {
    next((vm) => {
      vm.watchQueryData(vm);
      vm.queryParamsToDataForCore(vm, to.query);
    });
  },
  beforeRouteUpdate(to, from, next) {
    this.queryParamsToDataForCore(this, to.query);
    next();
  },
  methods: {
    onPasswordValidation(passwordValidation) {
      this.passwordValidation = passwordValidation;
    },
    validateChangeUserPassword() {
      this.clearErrors(this);
      let isValidationOk = true;

      if (!this.currentPassword) {
        this.error.currentPassword = this.$t("common.required");

        if (isValidationOk) {
          this.focusElement("currentPassword");
          isValidationOk = false;
        }
      }

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
      this.loading.changeUserPassword = true;
      const taskAction = "change-user-password";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.changeUserPasswordAborted
      );

      // register to task validation
      this.$root.$once(
        `${taskAction}-validation-failed-${eventId}`,
        this.changeUserPasswordValidationFailed
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.changeUserPasswordCompleted
      );

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            user: this.currentUsername,
            current_password: this.currentPassword,
            new_password: this.newPassword,
          },
          extra: {
            title: this.$t("action." + taskAction),
            description: this.$t("common.processing"),
            eventId,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.changeUserPassword = this.getErrorMessage(err);
        this.loading.changeUserPassword = false;
        return;
      }
    },
    changeUserPasswordAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.changeUserPassword = this.$t("error.generic_error");
      this.loading.changeUserPassword = false;
    },
    changeUserPasswordValidationFailed(validationErrors) {
      this.loading.changeUserPassword = false;

      for (const validationError of validationErrors) {
        // set i18n error message
        this.error.currentPassword = this.$t(
          "password." + validationError.error
        );
        this.focusElement("currentPassword");
      }
    },
    changeUserPasswordCompleted() {
      this.currentPassword = "";
      this.newPassword = "";
      this.clearConfirmPasswordCommand++;
      this.loading.changeUserPassword = false;
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";

.justify-flex-end {
  display: flex;
  justify-content: flex-end;
}

.icon-and-text {
  justify-content: flex-start;
}
</style>

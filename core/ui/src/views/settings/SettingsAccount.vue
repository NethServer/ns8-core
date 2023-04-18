<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
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
          <NsInlineNotification
            v-if="error.userInfo"
            kind="error"
            :title="$t('action.get-user-info')"
            :description="error.userInfo"
            :showCloseButton="false"
          />
        </cv-column>
      </cv-row>
      <cv-row>
        <cv-column>
          <cv-tile light>
            <cv-skeleton-text
              v-if="loading.getUserInfo"
              paragraph
              :line-count="3"
              width="70%"
            ></cv-skeleton-text>
            <div v-else class="user-info">
              <NsPictogram title="user avatar" class="avatar">
                <UserProfilePictogram />
              </NsPictogram>
              <div class="user-details">
                <div v-if="userInfo.display_name" class="display-name">
                  {{ userInfo.display_name }}
                </div>
                <div>{{ currentUsername }}</div>
              </div>
            </div>
            <NsButton
              kind="secondary"
              :icon="Edit20"
              @click="showEditDisplayNameModal"
              :disabled="loading.getUserInfo"
              class="mg-top-lg"
              >{{ $t("settings_account.edit_display_name") }}
            </NsButton>
          </cv-tile>
        </cv-column>
      </cv-row>
      <cv-row>
        <cv-column>
          <cv-tile light>
            <h4 class="mg-bottom-md">
              {{ $t("settings_account.change_password") }}
            </h4>
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
          <cv-tile light>
            <!-- two-factor authentication -->
            <h4>
              {{ $t("settings_account.2fa_title_with_acronym") }}
            </h4>
            <div class="title-description mg-bottom-xlg">
              {{ $t("settings_account.2fa_description") }}
            </div>
            <NsInlineNotification
              v-if="error.get2FaStatus"
              kind="error"
              :title="$t('settings_account.cannot_retrieve_2fa_status')"
              :description="error.get2FaStatus"
              :showCloseButton="false"
            />
            <cv-button-skeleton v-else-if="loading.get2FaStatus" />
            <template v-else>
              <template v-if="twoFaEnabled">
                <!-- 2fa enabled -->
                <div class="icon-and-text mg-bottom-lg">
                  <NsSvg :svg="CheckmarkFilled16" class="icon ns-success" />
                  <span>{{ $t("settings_account.2fa_is_enabled") }}</span>
                </div>
                <div class="mg-bottom-lg">
                  <span>
                    {{ $t("settings_account.need_to_reconfigure_2fa") }}
                  </span>
                  <cv-link @click="show2FaQrCodeModal">
                    {{ $t("settings_account.show_qr_code") }}
                  </cv-link>
                </div>
                <div>
                  <NsInlineNotification
                    v-if="error.revoke2Fa"
                    kind="error"
                    :title="$t('settings_account.cannot_revoke_2fa')"
                    :description="error.revoke2Fa"
                    :showCloseButton="false"
                  />
                  <NsButton
                    kind="danger"
                    :icon="Power20"
                    @click="showConfirmRevoke2FaModal"
                  >
                    {{ $t("settings_account.revoke_2fa") }}
                  </NsButton>
                </div>
              </template>
              <template v-else>
                <!-- 2fa disabled -->
                <div class="icon-and-text mg-bottom-lg">
                  <NsSvg :svg="InformationFilled16" class="icon ns-info" />
                  <span>{{ $t("settings_account.2fa_is_disabled") }}</span>
                </div>
                <NsButton
                  kind="secondary"
                  :icon="Power20"
                  @click="showEnable2FaModal"
                >
                  {{ $t("settings_account.enable_2fa") }}
                </NsButton>
              </template>
            </template>
          </cv-tile>
        </cv-column>
      </cv-row>
    </cv-grid>
    <!-- enable 2fa modal -->
    <Enable2FaModal
      :isShown="isShownEnable2FaModal"
      @hide="hideEnable2FaModal"
      @twoFaEnabled="onTwoFaEnabled"
    />
    <!-- 2fa qr code modal -->
    <TwoFaQrCodeModal
      :isShown="isShown2FaQrCodeModal"
      @hide="hide2FaQrCodeModal"
    />
    <!-- confirm revoke 2fa modal -->
    <ConfirmRevoke2FaModal
      :isShown="isShownConfirmRevoke2FaModal"
      @hide="hideConfirmRevoke2FaModal"
      @confirm="disable2Fa"
    />
    <!-- edit display name modal -->
    <EditAdminDisplayNameModal
      :isShown="isShownEditDisplayNameModal"
      :admin="account"
      @hide="hideEditDisplayNameModal"
      @displayNameUpdated="getUserInfo"
    />
  </div>
</template>

<script>
import to from "await-to-js";
import {
  QueryParamService,
  UtilService,
  TaskService,
  IconService,
  PageTitleService,
} from "@nethserver/ns8-ui-lib";
import TwoFaQrCodeModal from "@/components/settings/2FaQrCodeModal";
import Enable2FaModal from "@/components/settings/Enable2FaModal";
import ConfirmRevoke2FaModal from "@/components/settings/ConfirmRevoke2FaModal";
import TwoFaService from "@/mixins/2fa";
import NotificationService from "@/mixins/notification";
import EditAdminDisplayNameModal from "@/components/settings/EditAdminDisplayNameModal.vue";

export default {
  name: "SettingsAccount",
  components: {
    TwoFaQrCodeModal,
    Enable2FaModal,
    ConfirmRevoke2FaModal,
    EditAdminDisplayNameModal,
  },
  mixins: [
    TaskService,
    UtilService,
    IconService,
    QueryParamService,
    PageTitleService,
    TwoFaService,
    NotificationService,
  ],
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
      isShown2FaQrCodeModal: false,
      isShownEnable2FaModal: false,
      twoFaEnabled: false,
      isShownConfirmRevoke2FaModal: false,
      userInfo: {},
      isShownEditDisplayNameModal: false,
      loading: {
        changeUserPassword: false,
        get2FaStatus: false,
        revoke2Fa: false,
        getUserInfo: false,
      },
      error: {
        changeUserPassword: "",
        currentPassword: "",
        newPassword: "",
        confirmPassword: "",
        get2FaStatus: "",
        revoke2Fa: "",
        getUserInfo: "",
      },
    };
  },
  computed: {
    account() {
      return {
        user: this.currentUsername,
        display_name: this.userInfo ? this.userInfo.display_name : "",
      };
    },
  },
  created() {
    const loginInfo = this.getFromStorage("loginInfo");

    if (loginInfo) {
      this.currentUsername = loginInfo.username;
    }
    this.retrieve2FaStatus();
    this.getUserInfo();
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
    show2FaQrCodeModal() {
      this.isShown2FaQrCodeModal = true;
    },
    hide2FaQrCodeModal() {
      this.isShown2FaQrCodeModal = false;
    },
    showEnable2FaModal() {
      this.isShownEnable2FaModal = true;
    },
    hideEnable2FaModal() {
      this.isShownEnable2FaModal = false;
    },
    showConfirmRevoke2FaModal() {
      this.isShownConfirmRevoke2FaModal = true;
    },
    hideConfirmRevoke2FaModal() {
      this.isShownConfirmRevoke2FaModal = false;
    },
    async retrieve2FaStatus() {
      this.loading.get2FaStatus = true;
      this.error.get2FaStatus = "";
      const [err2FaStatus, response2FaStatus] = await to(this.get2FaStatus());
      this.loading.get2FaStatus = false;

      if (err2FaStatus) {
        console.error("error retrieving 2FA status", err2FaStatus);
        this.error.get2FaStatus = this.getErrorMessage(err2FaStatus);
        return;
      }
      this.twoFaEnabled = response2FaStatus.data.data;
    },
    async disable2Fa() {
      this.hideConfirmRevoke2FaModal();
      this.loading.revoke2Fa = true;
      this.error.revoke2Fa = "";

      const errRevoke = await to(this.revoke2Fa())[0];
      this.loading.revoke2Fa = false;

      if (errRevoke) {
        console.error("error revoking 2FA", errRevoke);
        this.error.revoke2Fa = this.getErrorMessage(errRevoke);
        return;
      }

      const notification = {
        title: this.$t("settings_account.2fa_disabled"),
        description: this.$t("task.completed"),
        type: "success",
      };
      this.createNotification(notification);

      this.retrieve2FaStatus();
    },
    onTwoFaEnabled() {
      const notification = {
        title: this.$t("settings_account.2fa_enabled"),
        description: this.$t("task.completed"),
        type: "success",
      };
      this.createNotification(notification);
      this.isShownEnable2FaModal = false;
      this.retrieve2FaStatus();
    },
    async getUserInfo() {
      this.loading.getUserInfo = true;
      this.error.getUserInfo = "";
      const taskAction = "get-user-info";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.getUserInfoAborted
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.getUserInfoCompleted
      );

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            user: this.currentUsername,
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
        this.error.getUserInfo = this.getErrorMessage(err);
        return;
      }
    },
    getUserInfoAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.getUserInfo = this.$t("error.generic_error");
      this.loading.getUserInfo = false;
    },
    getUserInfoCompleted(taskContext, taskResult) {
      this.userInfo = taskResult.output;
      this.loading.getUserInfo = false;
    },
    showEditDisplayNameModal() {
      this.isShownEditDisplayNameModal = true;
    },
    hideEditDisplayNameModal() {
      this.isShownEditDisplayNameModal = false;
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";

.icon-and-text {
  justify-content: flex-start;
}

.user-info {
  display: flex;
  align-items: center;

  .avatar {
    margin-right: $spacing-05;
  }

  .display-name {
    font-weight: bold;
    margin-bottom: $spacing-03;
  }
}
</style>

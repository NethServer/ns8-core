<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <NsModal
    kind="danger"
    size="default"
    :visible="isShown"
    :isLoading="loading.removeInternalProvider"
    @modal-hidden="onModalHidden"
    @primary-click="removeInternalProvider"
  >
    <template slot="title">{{
      $t("domain_configuration.delete_provider")
    }}</template>
    <template slot="content">
      <NsInlineNotification
        kind="warning"
        :title="$t('common.please_read_carefully')"
        :showCloseButton="false"
      />
      <div>
        {{
          $t("domain_configuration.delete_provider_confirm", {
            name: provider.id,
          })
        }}
      </div>
      <cv-form @submit.prevent="removeInternalProvider" class="mg-top-2">
        <!-- samba username -->
        <NsTextInput
          :label="$t('samba.adminuser')"
          v-model.trim="adminuser"
          :helper-text="$t('samba.enter_samba_admin_username')"
          :invalid-message="$t(error.adminuser)"
          :disabled="loading.removeInternalProvider"
          ref="adminuser"
        >
        </NsTextInput>
        <!-- samba password -->
        <NsTextInput
          :label="$t('samba.adminpass')"
          type="password"
          v-model="adminpass"
          :helper-text="$t('samba.enter_samba_admin_password')"
          :invalid-message="$t(error.adminpass)"
          :disabled="loading.removeInternalProvider"
          :password-hide-label="$t('password.hide_password')"
          :password-show-label="$t('password.show_password')"
          ref="adminpass"
          name="adminpass"
        ></NsTextInput>
      </cv-form>
      <NsInlineNotification
        v-if="error.removeInternalProvider"
        kind="error"
        :title="$t('action.remove-internal-provider')"
        :description="error.removeInternalProvider"
        :showCloseButton="false"
      />
    </template>
    <template slot="secondary-button">{{ $t("common.cancel") }}</template>
    <template slot="primary-button">{{
      $t("common.understood_delete")
    }}</template>
  </NsModal>
</template>

<script>
import { TaskService, UtilService } from "@nethserver/ns8-ui-lib";
import to from "await-to-js";

export default {
  name: "DeleteSambaProviderModal",
  mixins: [TaskService, UtilService],
  props: {
    isShown: {
      type: Boolean,
      default: true,
    },
    provider: { type: Object },
  },
  data() {
    return {
      adminuser: "",
      adminpass: "",
      loading: {
        removeInternalProvider: false,
      },
      error: {
        adminuser: "",
        adminpass: "",
        removeInternalProvider: "",
      },
    };
  },
  watch: {
    isShown: function () {
      if (this.isShown) {
        this.adminuser = "";
        this.adminpass = "";

        setTimeout(() => {
          this.focusElement("adminuser");
        }, 300);
      }
    },
  },
  methods: {
    onModalHidden() {
      this.$emit("hide");
    },
    validateRemoveInternalProvider() {
      this.clearErrors();
      let isValidationOk = true;

      // samba username

      if (!this.adminuser) {
        this.error.adminuser = "common.required";

        if (isValidationOk) {
          this.focusElement("adminuser");
          isValidationOk = false;
        }
      }

      // samba password

      if (!this.adminpass) {
        this.error.adminpass = "common.required";

        if (isValidationOk) {
          this.focusElement("adminpass");
          isValidationOk = false;
        }
      }
      return isValidationOk;
    },
    async removeInternalProvider() {
      const isValidationOk = this.validateRemoveInternalProvider();
      if (!isValidationOk) {
        return;
      }
      this.loading.removeInternalProvider = true;
      this.error.removeInternalProvider = "";
      const taskAction = "remove-internal-provider";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.removeInternalProviderAborted
      );

      // register to task validation
      this.$root.$once(
        `${taskAction}-validation-failed-${eventId}`,
        this.removeInternalProviderValidationFailed
      );
      this.$root.$once(
        `${taskAction}-validation-ok-${eventId}`,
        this.removeInternalProviderValidationOk
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.removeInternalProviderCompleted
      );

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            module_id: this.provider.id,
            adminuser: this.adminuser,
            adminpass: this.adminpass,
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
        this.error.removeInternalProvider = this.getErrorMessage(err);
        this.loading.removeInternalProvider = false;
        return;
      }
    },
    removeInternalProviderAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.loading.removeInternalProvider = false;

      // hide modal so that user can see error notification
      this.$emit("hide");
    },
    removeInternalProviderValidationOk() {
      // hide modal
      this.$emit("hide");
    },
    removeInternalProviderValidationFailed(validationErrors) {
      this.loading.removeInternalProvider = false;
      let focusAlreadySet = false;

      for (const validationError of validationErrors) {
        const param = validationError.parameter;
        // set i18n error message
        this.error.samba[param] = "domains." + validationError.error;

        if (!focusAlreadySet) {
          this.focusElement(param);
          focusAlreadySet = true;
        }
      }
    },
    removeInternalProviderCompleted() {
      this.loading.removeInternalProvider = false;

      // hide modal
      this.$emit("hide");

      // reload domains
      this.$emit("reloadDomains");

      // reload app drawer (samba file server might have disappeared)
      this.$root.$emit("reloadAppDrawer");
    },
  },
};
</script>

<style scoped lang="scss">
.mg-top-2 {
  margin-top: 2rem;
}
</style>

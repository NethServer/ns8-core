<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <NsModal
    size="default"
    :visible="isShown"
    :primary-button-disabled="loading.setAcmeServer"
    :isLoading="loading.setAcmeServer"
    @modal-hidden="onModalHidden"
    @primary-click="setAcmeServer"
  >
    <template slot="title">{{
      $t("settings_acme_servers.edit_acme_server_for_node", {
        node: server ? server.node : "",
      })
    }}</template>
    <template slot="content">
      <cv-form @submit.prevent="setAcmeServer">
        <NsInlineNotification
          v-if="server && server.nodeId == leaderNode.id"
          kind="warning"
          :title="$t('settings_acme_servers.reload_page')"
          :description="
            $t('settings_acme_servers.edit_acme_server_leader_node_warning')
          "
          :showCloseButton="false"
        />
        <NsTextInput
          v-model.trim="url"
          :label="$t('settings_acme_servers.url')"
          :invalid-message="error.url"
          :disabled="loading.setAcmeServer"
          data-modal-primary-focus
          ref="url"
        />
        <NsInlineNotification
          v-if="error.setAcmeServer"
          kind="error"
          :title="$t('action.set-acme-server')"
          :description="error.setAcmeServer"
          :showCloseButton="false"
        />
      </cv-form>
    </template>
    <template slot="secondary-button">{{ $t("common.cancel") }}</template>
    <template slot="primary-button">{{
      $t("settings_acme_servers.edit_acme_server")
    }}</template>
  </NsModal>
</template>

<script>
import to from "await-to-js";
import { UtilService, TaskService } from "@nethserver/ns8-ui-lib";
import { mapGetters } from "vuex";

export default {
  name: "EditAcmeServerModal",
  mixins: [UtilService, TaskService],
  props: {
    isShown: Boolean,
    server: {
      type: [Object, null],
    },
  },
  data() {
    return {
      url: "",
      loading: {
        setAcmeServer: false,
      },
      error: {
        setAcmeServer: "",
        url: "",
      },
    };
  },
  computed: {
    ...mapGetters(["leaderNode"]),
  },
  watch: {
    isShown: function () {
      if (this.isShown) {
        this.url = this.server.url;
      }
    },
  },
  methods: {
    onModalHidden() {
      this.clearErrors();
      this.$emit("hide");
    },
    validateSetAcmeServer() {
      this.clearErrors();
      let isValidationOk = true;

      // url

      if (!this.url) {
        this.error.url = this.$t("common.required");

        if (isValidationOk) {
          this.focusElement("url");
          isValidationOk = false;
        }
      }
      return isValidationOk;
    },
    async setAcmeServer() {
      if (!this.validateSetAcmeServer()) {
        return;
      }
      this.loading.setAcmeServer = true;
      this.error.setAcmeServer = "";
      const taskAction = "set-acme-server";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.setAcmeServerAborted
      );

      // register to task validation
      this.$root.$once(
        `${taskAction}-validation-ok-${eventId}`,
        this.setAcmeServerValidationOk
      );
      this.$root.$once(
        `${taskAction}-validation-failed-${eventId}`,
        this.setAcmeServerValidationFailed
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.setAcmeServerCompleted
      );

      const res = await to(
        this.createModuleTaskForApp(this.server.traefikInstance, {
          action: taskAction,
          data: {
            url: this.url,
          },
          extra: {
            title: this.$t("settings_acme_servers.edit_acme_server"),
            description: this.$t("common.processing"),
            eventId,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.setAcmeServer = this.getErrorMessage(err);
        this.loading.setAcmeServer = false;
        return;
      }
    },
    setAcmeServerAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.loading.setAcmeServer = false;

      // hide modal so that user can see error notification
      this.$emit("hide");
    },
    setAcmeServerValidationOk() {
      this.loading.setAcmeServer = false;

      // hide modal after validation
      this.$emit("hide");
    },
    setAcmeServerValidationFailed(validationErrors) {
      this.loading.setAcmeServer = false;
      let focusAlreadySet = false;

      for (const validationError of validationErrors) {
        const param = validationError.parameter;

        // set i18n error message
        this.error[param] = this.getI18nStringWithFallback(
          "settings_acme_servers." + validationError.error,
          "error." + validationError.error
        );

        if (!focusAlreadySet) {
          this.focusElement(param);
          focusAlreadySet = true;
        }
      }
    },
    setAcmeServerCompleted() {
      this.loading.setAcmeServer = false;

      // reload servers
      this.$emit("reloadServers");
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";
</style>

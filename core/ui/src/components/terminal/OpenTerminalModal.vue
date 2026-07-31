<!--
  Copyright (C) 2026 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <NsModal
    size="default"
    :visible="isShown"
    :isLoading="loading.open"
    :primaryButtonDisabled="!canOpen"
    @modal-hidden="onModalHidden"
    @primary-click="openTerminal"
  >
    <template slot="title">{{
      $t("terminal.open_on_node", { node: nodeLabel })
    }}</template>
    <template slot="content">
      <cv-skeleton-text
        v-if="loading.probe"
        :paragraph="true"
        :line-count="4"
        heading
      ></cv-skeleton-text>
      <template v-else>
        <NsInlineNotification
          v-if="error.probe"
          kind="error"
          :title="$t('action.probe-terminal-access')"
          :description="error.probe"
          :showCloseButton="false"
        />
        <!-- sshd is not reachable on the VPN address: nothing will work -->
        <NsInlineNotification
          v-else-if="!probe.listen_wg0"
          kind="error"
          :title="$t('terminal.sshd_not_listening')"
          :description="$t('terminal.sshd_not_listening_description')"
          :showCloseButton="false"
        />
        <!-- root by password is refused: say so before a password is typed -->
        <NsInlineNotification
          v-else-if="!probe.permit_root_login"
          kind="warning"
          :title="$t('terminal.root_password_refused')"
          :description="$t('terminal.root_password_refused_description')"
          :showCloseButton="false"
        />
        <NsInlineNotification
          v-else-if="!probe.password_auth"
          kind="warning"
          :title="$t('terminal.password_auth_disabled')"
          :description="$t('terminal.password_auth_disabled_description')"
          :showCloseButton="false"
        />
        <cv-form @submit.prevent="openTerminal">
          <cv-text-input
            :label="$t('terminal.system_user')"
            v-model.trim="username"
            :invalid-message="error.username"
            :disabled="loading.open"
            ref="username"
          >
          </cv-text-input>
          <cv-text-input
            :label="$t('terminal.system_password')"
            v-model="password"
            type="password"
            :invalid-message="error.password"
            :disabled="loading.open"
            :password-hide-label="$t('password.hide_password')"
            :password-show-label="$t('password.show_password')"
            ref="password"
          >
          </cv-text-input>
        </cv-form>
        <NsInlineNotification
          kind="info"
          :title="$t('terminal.credentials_notice')"
          :description="$t('terminal.credentials_notice_description')"
          :showCloseButton="false"
        />
        <NsInlineNotification
          v-if="error.open"
          kind="error"
          :title="$t('terminal.cannot_open')"
          :description="error.open"
          :showCloseButton="false"
        />
      </template>
    </template>
    <template slot="primary-button">{{ $t("terminal.open") }}</template>
  </NsModal>
</template>

<script>
import to from "await-to-js";
import { TaskService, UtilService, IconService } from "@nethserver/ns8-ui-lib";

export default {
  name: "OpenTerminalModal",
  mixins: [TaskService, UtilService, IconService],
  props: {
    isShown: {
      type: Boolean,
      required: true,
    },
    nodeId: {
      type: [String, Number],
      required: true,
    },
    nodeLabel: {
      type: String,
      default: "",
    },
  },
  data() {
    return {
      username: "root",
      password: "",
      probe: {
        permit_root_login: false,
        password_auth: false,
        listen_wg0: false,
        port: 22,
      },
      loading: {
        probe: false,
        open: false,
      },
      error: {
        probe: "",
        open: "",
        username: "",
        password: "",
      },
    };
  },
  computed: {
    canOpen() {
      return !this.loading.probe && !this.loading.open;
    },
  },
  watch: {
    isShown: function () {
      if (this.isShown) {
        this.clearErrors();
        this.username = "root";
        this.password = "";
        this.probeAccess();
      }
    },
  },
  methods: {
    clearErrors() {
      this.error.probe = "";
      this.error.open = "";
      this.error.username = "";
      this.error.password = "";
    },
    onModalHidden() {
      // Do not leave the password in a reactive store once the modal closes.
      this.password = "";
      this.$emit("hide");
    },
    /**
     * Ask the node what sshd will accept, before the user types anything. The
     * probe passes a connection spec to sshd -T, so it reports the policy that
     * applies to the terminal rather than the global one.
     */
    async probeAccess() {
      this.loading.probe = true;
      const taskAction = "probe-terminal-access";
      const eventId = this.getUuid();

      this.$root.$once(`${taskAction}-aborted-${eventId}`, this.probeAborted);
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.probeCompleted
      );

      const res = await to(
        this.createNodeTask(this.nodeId, {
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
        this.error.probe = this.getErrorMessage(err);
        this.loading.probe = false;
      }
    },
    probeCompleted(taskContext, taskResult) {
      this.probe = taskResult.output;
      this.loading.probe = false;

      setTimeout(() => {
        this.focusElement("password");
      }, 300);
    },
    probeAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.probe = this.$t("terminal.probe_failed");
      this.loading.probe = false;
    },
    validate() {
      this.error.username = "";
      this.error.password = "";
      let isValidationOk = true;

      if (!this.username) {
        this.error.username = this.$t("common.required");
        this.focusElement("username");
        isValidationOk = false;
      }
      if (!this.password) {
        this.error.password = this.$t("common.required");
        if (isValidationOk) {
          this.focusElement("password");
        }
        isValidationOk = false;
      }
      return isValidationOk;
    },
    openTerminal() {
      this.error.open = "";
      if (!this.validate()) {
        return;
      }
      // The parent owns the websocket: it sends these once and drops them.
      this.$emit("open", {
        username: this.username,
        password: this.password,
      });
      this.password = "";
    },
    setOpenError(message) {
      this.error.open = message;
      this.loading.open = false;
    },
  },
};
</script>

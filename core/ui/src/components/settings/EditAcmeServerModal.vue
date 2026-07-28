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
      $t("settings_acme_servers.edit_acme_server")
    }}</template>
    <template slot="content">
      <cv-form @submit.prevent="setAcmeServer">
        <div class="mg-bottom-md">
          {{
            $t("settings_acme_servers.acme_settings_for_node", {
              node: server ? server.node : "",
            })
          }}
        </div>
        <!-- set-acme-server restarts Traefik on the edited node, whichever it is -->
        <NsInlineNotification
          v-if="server"
          kind="warning"
          :title="$t('settings_tls_certificates.traefik_will_be_restarted')"
          :description="
            $t('settings_acme_servers.acme_restart_message', {
              node: server.node,
            })
          "
          :showCloseButton="false"
        />
        <NsTextInput
          v-model.trim="url"
          :label="$t('settings_acme_servers.url')"
          :helper-text="$t('settings_acme_servers.url_helper')"
          :invalid-message="error.url"
          :disabled="loading.setAcmeServer"
          data-modal-primary-focus
          ref="url"
        />
        <!-- an older module reports no challenge type: do not invent one -->
        <!-- ref on the wrapper: a ref inside v-for is an array, focusElement would throw -->
        <div
          v-if="isChallengeSupported"
          class="mg-top-md"
          tabindex="-1"
          ref="challenge"
        >
          <label id="acme-challenge-label" class="bx--label">{{
            $t("settings_acme_servers.challenge")
          }}</label>
          <div role="group" aria-labelledby="acme-challenge-label">
            <cv-radio-group :vertical="true">
              <!-- disabled must be per button: cv-radio-group has no such prop -->
              <cv-radio-button
                v-for="challengeType in challengeTypes"
                :key="challengeType.value"
                :label="$t(challengeType.labelKey)"
                :value="challengeType.value"
                name="acme-challenge"
                v-model="challenge"
                :disabled="loading.setAcmeServer"
              />
            </cv-radio-group>
          </div>
          <div
            v-if="error.challenge"
            class="bx--form-requirement challenge-error"
          >
            {{ error.challenge }}
          </div>
        </div>
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
    <template slot="primary-button">{{ $t("common.save") }}</template>
  </NsModal>
</template>

<script>
import to from "await-to-js";
import { UtilService, TaskService } from "@nethserver/ns8-ui-lib";
import { mapState } from "vuex";

// DNS-01 is not implemented yet: adding an entry here is enough
export const ACME_CHALLENGE_TYPES = [
  {
    value: "HTTP-01",
    labelKey: "settings_acme_servers.challenge_http_01",
    tagKind: "blue",
  },
  {
    value: "TLS-ALPN-01",
    labelKey: "settings_acme_servers.challenge_tls_alpn_01",
    tagKind: "purple",
  },
];

export const DEFAULT_ACME_CHALLENGE_TYPE = "HTTP-01";

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
      challenge: DEFAULT_ACME_CHALLENGE_TYPE,
      loading: {
        setAcmeServer: false,
      },
      error: {
        setAcmeServer: "",
        url: "",
        challenge: "",
      },
      // [eventName, handler] pairs to remove in beforeDestroy
      taskListeners: [],
    };
  },
  computed: {
    ...mapState(["isWebsocketConnected"]),
    challengeTypes() {
      return ACME_CHALLENGE_TYPES;
    },
    isChallengeSupported() {
      return !!this.server && this.server.challenge !== undefined;
    },
  },
  watch: {
    isWebsocketConnected: function (isConnected) {
      // safety net: nothing else clears the pending state if the restart cut
      // the request off before validation-ok could be delivered
      if (isConnected && this.loading.setAcmeServer) {
        this.loading.setAcmeServer = false;
        this.$emit("hide");
      }
    },
    isShown: function () {
      if (this.isShown) {
        this.url = this.server.url;
        this.challenge = this.server.challenge || DEFAULT_ACME_CHALLENGE_TYPE;
      }
    },
  },
  beforeDestroy() {
    // a save completing after destroy runs focusElement on a deleted ref
    this.taskListeners.forEach(([eventName, handler]) => {
      this.$root.$off(eventName, handler);
    });
    this.taskListeners = [];
  },
  methods: {
    registerTaskListener(eventName, handler) {
      this.$root.$once(eventName, handler);
      this.taskListeners.push([eventName, handler]);
    },
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

      // challenge

      if (this.isChallengeSupported && !this.challenge) {
        this.error.challenge = this.$t("common.required");

        if (isValidationOk) {
          this.focusElement("challenge");
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
      this.registerTaskListener(
        `${taskAction}-aborted-${eventId}`,
        this.setAcmeServerAborted
      );

      // register to task validation
      this.registerTaskListener(
        `${taskAction}-validation-ok-${eventId}`,
        this.setAcmeServerValidationOk
      );
      this.registerTaskListener(
        `${taskAction}-validation-failed-${eventId}`,
        this.setAcmeServerValidationFailed
      );

      // register to task completion
      this.registerTaskListener(
        `${taskAction}-completed-${eventId}`,
        this.setAcmeServerCompleted
      );

      // the action sets additionalProperties false: send only reported fields
      const data = { url: this.url };

      if (this.isChallengeSupported) {
        data.challenge = this.challenge;
      }

      if (this.server.email !== undefined) {
        // the action resets the email when the field is missing
        data.email = this.server.email;
      }

      const res = await to(
        this.createModuleTaskForApp(this.server.traefikInstance, {
          action: taskAction,
          data: data,
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
        const message = this.getI18nStringWithFallback(
          "settings_acme_servers." + validationError.error,
          "error." + validationError.error
        );

        // focusElement would throw on a parameter with no field of its own
        if (param !== "setAcmeServer" && param in this.error) {
          this.error[param] = message;

          if (!focusAlreadySet) {
            this.focusElement(param);
            focusAlreadySet = true;
          }
        } else {
          this.error.setAcmeServer = message;
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

.bx--inline-notification {
  max-width: 38rem;
}

// Carbon reveals .bx--form-requirement only next to a [data-invalid] wrapper,
// which a radio group never is; $text-error is not reachable from carbon-utils
.challenge-error {
  display: block;
  max-height: none;
  overflow: visible;
  color: #da1e28;
}
</style>

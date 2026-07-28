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
          :helper-text="$t('settings_acme_servers.url_helper')"
          :invalid-message="error.url"
          :disabled="loading.setAcmeServer"
          data-modal-primary-focus
          ref="url"
        />
        <!-- hidden when the module did not report a challenge type: it is too
             old to know about it, so the UI must not invent a value -->
        <!-- ref is on the wrapper: a ref inside v-for is an array, and
             focusElement() calls .focus() on it without checking -->
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
              <!-- disabled goes on each button: cv-radio-group has no such prop,
                   while cv-radio-button forwards attributes to its input -->
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
          <!-- .bx--form-requirement is display:none unless it follows a
               [data-invalid] field wrapper, which a radio group never is -->
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
import { mapGetters } from "vuex";

// challenge types supported by the traefik set-acme-server action.
// DNS-01 is not implemented yet: adding it here is enough to expose it
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
      // [eventName, handler] pairs registered on $root, removed in beforeDestroy
      taskListeners: [],
    };
  },
  computed: {
    ...mapGetters(["leaderNode"]),
    challengeTypes() {
      return ACME_CHALLENGE_TYPES;
    },
    isChallengeSupported() {
      return !!this.server && this.server.challenge !== undefined;
    },
  },
  watch: {
    isShown: function () {
      if (this.isShown) {
        this.url = this.server.url;
        this.challenge = this.server.challenge || DEFAULT_ACME_CHALLENGE_TYPE;
      }
    },
  },
  beforeDestroy() {
    // remove task listeners still registered on $root, otherwise a save
    // completing after this modal is gone runs focusElement() on a deleted ref
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

      // only echo back the fields the module actually reported: the action
      // schema sets additionalProperties false, so sending a key an older
      // traefik does not know about would fail the whole request
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

        // a parameter with no field of its own, such as email, would land on a
        // non-reactive key and make focusElement() throw on an undefined ref
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

// Carbon only reveals .bx--form-requirement next to a [data-invalid] field
// wrapper, and a radio group never gets one, so the message needs both the
// visibility and the error color. $text-error is not reachable from
// carbon-utils in this Carbon version, hence the literal value
.challenge-error {
  display: block;
  max-height: none;
  overflow: visible;
  color: #da1e28;
}
</style>

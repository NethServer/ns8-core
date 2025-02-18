<!--
  Copyright (C) 2025 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <NsModal
    size="default"
    :visible="isShown"
    @modal-hidden="$emit('hide')"
    @primary-click="alterExternalDomain"
  >
    <template slot="title">{{ $t("domains.edit_external_domain") }}</template>
    <template slot="content">
      <cv-form @submit.prevent="alterExternalDomain">
        <NsTextInput
          :label="$t('domains.bind_dn')"
          v-model.trim="bind_dn"
          :invalid-message="$t(error.bind_dn)"
          :disabled="loading.alterExternalDomain"
          ref="bind_dn"
        >
        </NsTextInput>
        <NsTextInput
          type="password"
          :label="$t('domains.bind_password')"
          v-model="bind_password"
          :invalid-message="$t(error.bind_password)"
          :password-hide-label="$t('password.hide_password')"
          :password-show-label="$t('password.show_password')"
          :disabled="loading.alterExternalDomain"
          ref="bind_password"
        >
        </NsTextInput>
        <NsToggle
          :label="$t('domains.tls')"
          value="tlsValue"
          :form-item="true"
          v-model="tls"
          :disabled="loading.alterExternalDomain"
          ref="tls"
        >
          <template slot="text-left">{{ $t("common.disabled") }}</template>
          <template slot="text-right">{{ $t("common.enabled") }}</template>
        </NsToggle>
        <NsInlineNotification
          v-if="error.tls"
          kind="error"
          :title="$t(error.tls)"
          :showCloseButton="false"
        />
        <template v-if="tls">
          <NsToggle
            :label="$t('domains.tls_verify')"
            value="tls_verifyValue"
            :form-item="true"
            v-model="tls_verify"
            :disabled="loading.alterExternalDomain"
            ref="tls_verify"
          >
            <template slot="text-left">{{ $t("common.disabled") }}</template>
            <template slot="text-right">{{ $t("common.enabled") }}</template>
          </NsToggle>
          <NsInlineNotification
            v-if="error.tls_verify"
            kind="error"
            :title="$t(error.tls_verify)"
            :showCloseButton="false"
          />
        </template>
      </cv-form>
    </template>
    <template slot="secondary-button">{{ $t("common.cancel") }}</template>
    <template slot="primary-button">{{
      $t("domains.edit_password_policy")
    }}</template>
  </NsModal>
</template>

<script>
import { UtilService, TaskService, IconService } from "@nethserver/ns8-ui-lib";
import to from "await-to-js";

export default {
  name: "EditExternalDomainModal",
  mixins: [UtilService, TaskService, IconService],
  props: {
    isShown: Boolean,
    domain: { type: [Object, null] },
  },
  data() {
    return {
      loading: {
        alterExternalDomain: false,
      },
      bind_dn: "",
      bind_password: "",
      tls: "",
      tls_verify: "",
      error: {
        bind_dn: "",
        bind_password: "",
        tls: "",
        tls_verify: "",
      },
    };
  },
  watch: {
    isShown: function () {
      if (this.isShown) {
        this.clearErrors(this);

        setTimeout(() => {
          this.focusElement("bind_dn");
        }, 300);
      }
      this.bind_dn = this.domain.bind_dn;
      this.bind_password = this.domain.bind_password;
      this.tls = this.domain.tls;
      this.tls_verify = this.domain.tls_verify;
    },
  },
  methods: {
    validateAlterExternalDomain() {
      this.clearErrors(this);
      let isValidationOk = true;

      if (!this.bind_dn) {
        this.error.bind_dn = "common.required";

        if (isValidationOk) {
          this.focusElement("bind_dn");
          isValidationOk = false;
        }
      }

      if (!this.bind_password) {
        this.error.bind_password = "common.required";

        if (isValidationOk) {
          this.focusElement("bind_password");
          isValidationOk = false;
        }
      }
      return isValidationOk;
    },
    async alterExternalDomain() {
      if (!this.validateAlterExternalDomain()) {
        return;
      }
      this.loading.alterExternalDomain = true;
      this.error.alterExternalDomain = "";
      const taskAction = "alter-external-domain";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.alterExternalDomainAborted
      );
      //register to task validation error
      this.$root.$once(
        `${taskAction}-validation-failed-${eventId}`,
        this.alterExternalDomainValidationFailed
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.alterExternalDomainCompleted
      );

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            domain: this.domain.name,
            protocol: this.domain.protocol,
            bind_dn: this.bind_dn,
            bind_password: this.bind_password,
            tls: this.tls,
            tls_verify: this.tls_verify,
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
        this.error.alterExternalDomain = this.getErrorMessage(err);
        this.loading.alterExternalDomain = false;
        return;
      }
    },
    alterExternalDomainAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.loading.alterExternalDomain = false;
    },
    alterExternalDomainValidationFailed(validationErrors) {
      this.loading.alterExternalDomain = false;

      for (const validationError of validationErrors) {
        const param = validationError.parameter;

        // set i18n error message
        this.error[param] = this.$t("domains." + validationError.error);
      }
    },
    alterExternalDomainCompleted() {
      this.loading.alterExternalDomain = false;
      this.$emit("reload");
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";
.left-margin {
  margin-left: 1rem;
}
</style>

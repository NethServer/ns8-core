<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <NsModal
    v-if="domain"
    size="default"
    :visible="isShown"
    @modal-hidden="$emit('hide')"
    @primary-click="addExternalProvider"
    :primary-button-disabled="loading.addExternalProvider"
  >
    <template slot="title">{{
      $t("domain_configuration.add_provider")
    }}</template>
    <template slot="content">
      <cv-form @submit.prevent="addExternalProvider">
        <cv-text-input
          :label="$t('domains.host')"
          v-model.trim="host"
          :invalid-message="$t(error.host)"
          :disabled="loading.addExternalProvider"
          ref="host"
        >
        </cv-text-input>
        <cv-text-input
          :label="$t('domains.port')"
          v-model.trim="port"
          :invalid-message="$t(error.port)"
          type="number"
          :disabled="loading.addExternalProvider"
          ref="port"
        >
        </cv-text-input>
      </cv-form>
      <NsInlineNotification
        v-if="error.addExternalProvider"
        kind="error"
        :title="$t('action.add-external-provider')"
        :description="error.addExternalProvider"
        :showCloseButton="false"
      />
    </template>
    <template slot="secondary-button">{{ $t("common.cancel") }}</template>
    <template slot="primary-button">{{
      $t("domain_configuration.add_provider")
    }}</template>
  </NsModal>
</template>

<script>
import { UtilService, TaskService, IconService } from "@nethserver/ns8-ui-lib";
import to from "await-to-js";

export default {
  name: "AddExternalProviderModal",
  mixins: [UtilService, TaskService, IconService],
  props: {
    isShown: {
      type: Boolean,
      default: true,
    },
    domain: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      host: "",
      port: "",
      loading: {
        addExternalProvider: false,
      },
      error: {
        host: "",
        port: "",
        addExternalProvider: "",
      },
    };
  },
  watch: {
    isShown: function () {
      if (this.isShown) {
        this.clearErrors(this);

        setTimeout(() => {
          this.focusElement("host");
        }, 300);
      }
    },
  },
  methods: {
    validateAddExternalProvider() {
      this.clearErrors(this);
      let isValidationOk = true;

      if (!this.host) {
        this.error.host = "common.required";

        if (isValidationOk) {
          this.focusElement("host");
          isValidationOk = false;
        }
      }

      if (!this.port) {
        this.error.port = "common.required";

        if (isValidationOk) {
          this.focusElement("port");
          isValidationOk = false;
        }
      }
      return isValidationOk;
    },
    async addExternalProvider() {
      if (!this.validateAddExternalProvider()) {
        return;
      }
      this.loading.addExternalProvider = true;
      const taskAction = "add-external-provider";

      // register to task validation
      this.$root.$off(taskAction + "-validation-failed");
      this.$root.$once(
        taskAction + "-validation-failed",
        this.addExternalProviderValidationFailed
      );

      // register to task completion
      this.$root.$once(
        taskAction + "-completed",
        this.addExternalProviderCompleted
      );

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            domain: this.domain.name,
            protocol: "ldap",
            host: this.host,
            port: parseInt(this.port),
          },
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.addExternalProvider = this.getErrorMessage(err);
        return;
      }
    },
    addExternalProviderValidationFailed(validationErrors) {
      this.loading.addExternalProvider = false;
      let focusAlreadySet = false;

      for (const validationError of validationErrors) {
        const param = validationError.parameter;
        // set i18n error message
        this.error[param] = "domains." + validationError.error;

        if (!focusAlreadySet) {
          this.focusElement(param);
          focusAlreadySet = true;
        }
      }
    },
    addExternalProviderCompleted() {
      // hide modal after validation
      this.$emit("hide");
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";
</style>

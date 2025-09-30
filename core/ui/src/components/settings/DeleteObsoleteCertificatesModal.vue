<!--
  Copyright (C) 2025 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <NsModal
    kind="danger"
    size="default"
    :visible="isShown"
    :primary-button-disabled="loading.deleteCertificate"
    :isLoading="loading.deleteCertificate"
    @modal-hidden="onModalHidden"
    @primary-click="deleteCertificate"
  >
    <template slot="title">{{
      $t("settings_tls_certificates.delete_obsolete_certificates")
    }}</template>
    <template slot="content">
      <cv-form @submit.prevent="deleteCertificate">
        <p class="mb-2">
          {{ $t("settings_tls_certificates.obsolete_certificates_message_1") }}
        </p>
        <p class="mb-6">
          {{ $t("settings_tls_certificates.obsolete_certificates_message_2") }}
        </p>
        <NsComboBox
          v-model="selectedNodeId"
          :label="$t('common.choose')"
          :title="$t('common.node')"
          :auto-filter="true"
          :auto-highlight="true"
          :options="nodes"
          :disabled="loading.deleteCertificate"
          :invalid-message="error.node"
          light
          ref="node"
        />
        <NsInlineNotification
          kind="warning"
          :title="$t('settings_tls_certificates.traefik_will_be_restarted')"
          :description="
            $t('settings_tls_certificates.delete_obsolete_certificates_message')
          "
          :showCloseButton="false"
        />
        <NsInlineNotification
          v-if="error.deleteCertificate"
          kind="error"
          :title="
            $t('settings_tls_certificates.cannot_delete_obsolete_certificates')
          "
          :description="error.deleteCertificate"
          :showCloseButton="false"
        />
      </cv-form>
    </template>
    <template slot="secondary-button">{{ $t("common.cancel") }}</template>
    <template slot="primary-button">{{
      $t("settings_tls_certificates.delete_obsolete_certificates")
    }}</template>
  </NsModal>
</template>

<script>
import to from "await-to-js";
import { UtilService, TaskService } from "@nethserver/ns8-ui-lib";

export default {
  name: "DeleteObsoleteCertificatesModal",
  mixins: [UtilService, TaskService],
  props: {
    isShown: Boolean,
    nodes: {
      type: Array,
      required: true,
    },
  },
  data() {
    return {
      names: "",
      selectedNodeId: null,
      loading: {
        deleteCertificate: false,
      },
      error: {
        deleteCertificate: "",
        node: "",
      },
    };
  },
  watch: {
    isShown: function () {
      if (this.isShown) {
        this.clearErrors();
        this.clearFields();
      }
    },
  },
  methods: {
    onModalHidden() {
      this.clearErrors();
      this.$emit("hide");
    },
    validateDeleteCertificate() {
      this.clearErrors();

      let isValidationOk = true;

      // node

      if (!this.selectedNodeId) {
        this.error.node = this.$t("common.required");

        if (isValidationOk) {
          this.focusElement("node");
          isValidationOk = false;
        }
      }
      return isValidationOk;
    },
    async deleteCertificate() {
      this.error.deleteCertificate = "";

      if (!this.validateDeleteCertificate()) {
        return;
      }
      this.loading.deleteCertificate = true;
      const taskAction = "delete-certificate";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.deleteCertificateAborted
      );

      // register to task validation
      this.$root.$once(
        `${taskAction}-validation-ok-${eventId}`,
        this.deleteCertificateValidationOk
      );
      this.$root.$once(
        `${taskAction}-validation-failed-${eventId}`,
        this.deleteCertificateValidationFailed
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.deleteCertificateCompleted
      );

      const selectedNode = this.nodes.find(
        (node) => node.value === this.selectedNodeId
      );
      const traefikInstance = selectedNode.traefikInstance;
      const node = selectedNode.label;

      const res = await to(
        this.createModuleTaskForApp(traefikInstance, {
          action: taskAction,
          data: {
            obsolete: true,
            type: "internal",
          },
          extra: {
            title: this.$t(
              "settings_tls_certificates.delete_obsolete_certificates_for_node",
              {
                node,
              }
            ),
            description: this.$t("common.processing"),
            eventId,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.deleteCertificate = this.getErrorMessage(err);
        this.loading.deleteCertificate = false;
        return;
      }
    },
    deleteCertificateAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.deleteCertificate = this.$t("error.generic_error");
      this.loading.deleteCertificate = false;

      // reload certificates
      this.$root.$emit("reloadCertificates");
    },
    deleteCertificateValidationOk() {
      this.loading.deleteCertificate = false;
      this.clearFields();

      // hide modal after validation
      this.$emit("hide");
    },
    deleteCertificateValidationFailed(validationErrors) {
      this.loading.deleteCertificate = false;
      let focusAlreadySet = false;

      for (const validationError of validationErrors) {
        const param = validationError.parameter;

        // set i18n error message
        this.error[param] = this.getI18nStringWithFallback(
          "settings_tls_certificates." + validationError.error,
          "error." + validationError.error,
          null,
          { value: validationError.value }
        );

        if (!focusAlreadySet) {
          this.focusElement(param);
          focusAlreadySet = true;
        }
      }
    },
    deleteCertificateCompleted() {
      this.loading.deleteCertificate = false;
      this.clearFields();

      // reload certificates
      this.$root.$emit("reloadCertificates");
    },
    clearFields() {
      this.names = "";
      this.selectedNodeId = "";

      if (this.$refs.node) {
        this.$refs.node.clearValue();
      }
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";
</style>

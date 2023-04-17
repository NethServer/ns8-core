<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <NsModal
    size="default"
    :visible="isShown"
    :primary-button-disabled="loading.setCertificate"
    :isLoading="loading.setCertificate"
    @modal-hidden="onModalHidden"
    @primary-click="setCertificate"
  >
    <template slot="title">{{
      $t("settings_tls_certificates.request_certificate")
    }}</template>
    <template slot="content">
      <cv-form @submit.prevent="setCertificate">
        <NsTextInput
          v-model.trim="fqdn"
          :label="$t('settings_tls_certificates.fqdn')"
          :helper-text="$t('settings_tls_certificates.fqdn_helper')"
          :invalid-message="error.fqdn"
          :disabled="loading.setCertificate"
          data-modal-primary-focus
          ref="fqdn"
        />
        <cv-combo-box
          v-model="selectedNodeId"
          :label="$t('common.choose')"
          :title="$t('common.node')"
          :auto-filter="true"
          :auto-highlight="true"
          :options="nodes"
          :disabled="loading.setCertificate"
          :invalid-message="error.node"
          light
          class="mg-bottom-12"
          ref="node"
        >
        </cv-combo-box>
        <NsInlineNotification
          v-if="error.setCertificate"
          kind="error"
          :title="$t('action.set-certificate')"
          :description="error.setCertificate"
          :showCloseButton="false"
        />
      </cv-form>
    </template>
    <template slot="secondary-button">{{ $t("common.cancel") }}</template>
    <template slot="primary-button">{{
      $t("settings_tls_certificates.request_certificate")
    }}</template>
  </NsModal>
</template>

<script>
import to from "await-to-js";
import {
  UtilService,
  TaskService,
  DateTimeService,
} from "@nethserver/ns8-ui-lib";
import { mapActions } from "vuex";

export default {
  name: "RequestTlsCertificateModal",
  mixins: [UtilService, TaskService, DateTimeService],
  props: {
    isShown: Boolean,
    nodes: {
      type: Array,
      required: true,
    },
    defaultNodeId: {
      type: String,
      default: "",
    },
    allCertificates: {
      type: Array,
      required: true,
    },
  },
  data() {
    return {
      fqdn: "",
      selectedNodeId: "",
      loading: {
        setCertificate: false,
      },
      error: {
        setCertificate: "",
        fqdn: "",
        node: "",
      },
    };
  },
  watch: {
    defaultNodeId: function () {
      this.updateSelectedNodeId();
    },
  },
  created() {
    this.updateSelectedNodeId();
  },
  methods: {
    ...mapActions([
      "addPendingTlsCertificateInStore",
      "removePendingTlsCertificateInStore",
    ]),
    updateSelectedNodeId() {
      if (this.defaultNodeId != "all") {
        this.selectedNodeId = this.defaultNodeId;
      } else {
        this.selectedNodeId = "";
      }
    },
    onModalHidden() {
      this.clearErrors();
      this.$emit("hide");
    },
    validateSetCertificate() {
      this.clearErrors();

      let isValidationOk = true;

      // fqdn

      if (!this.fqdn) {
        this.error.fqdn = this.$t("common.required");

        if (isValidationOk) {
          this.focusElement("fqdn");
          isValidationOk = false;
        }
      }

      // check if fqdn already exists
      const duplicatedFqdn = this.allCertificates.find(
        (cert) => cert.fqdn === this.fqdn
      );

      if (duplicatedFqdn) {
        this.error.fqdn = this.$t(
          "settings_tls_certificates.fqdn_already_exists"
        );

        if (isValidationOk) {
          this.focusElement("fqdn");
          isValidationOk = false;
        }
      }

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
    async setCertificate() {
      if (!this.validateSetCertificate()) {
        return;
      }

      this.loading.setCertificate = true;
      this.error.setCertificate = "";
      const taskAction = "set-certificate";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.setCertificateAborted
      );

      // register to task validation
      this.$root.$once(
        `${taskAction}-validation-ok-${eventId}`,
        this.setCertificateValidationOk
      );
      this.$root.$once(
        `${taskAction}-validation-failed-${eventId}`,
        this.setCertificateValidationFailed
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.setCertificateCompleted
      );

      const selectedNode = this.nodes.find(
        (node) => node.value === this.selectedNodeId
      );
      const traefikInstance = selectedNode.traefikInstance;

      const logsStartDate = this.formatDate(new Date(), "yyyy-MM-dd");
      const logsHours = this.formatDate(new Date(), "HH");
      const logsMins = this.formatDate(new Date(), "mm");
      const logsStartTime = `${logsHours}%3A${logsMins}`;

      const res = await to(
        this.createModuleTaskForApp(traefikInstance, {
          action: taskAction,
          data: {
            fqdn: this.fqdn,
            sync: true,
          },
          extra: {
            title: this.$t(
              "settings_tls_certificates.request_certificate_for_fqdn",
              {
                fqdn: this.fqdn,
              }
            ),
            description: this.$t("common.processing"),
            eventId,
            logs: {
              path: `?searchQuery=&context=module&selectedAppId=${traefikInstance}&followLogs=false&startDate=${logsStartDate}&startTime=${logsStartTime}&autoStartSearch=true`,
              instance: traefikInstance,
            },
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.setCertificate = this.getErrorMessage(err);
        this.loading.setCertificate = false;
        return;
      }

      // add pending certificate
      this.addPendingTlsCertificateInStore(this.fqdn);

      // reload certificates
      this.$root.$emit("reloadCertificates");
    },
    setCertificateAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.loading.setCertificate = false;

      // remove pending certificate
      this.removePendingTlsCertificateInStore(taskContext.data.fqdn);

      // reload certificates
      this.$root.$emit("reloadCertificates");
    },
    setCertificateValidationOk() {
      this.loading.setCertificate = false;
      this.clearFields();

      // hide modal after validation
      this.$emit("hide");
    },
    setCertificateValidationFailed(validationErrors) {
      this.loading.setCertificate = false;
      let focusAlreadySet = false;

      for (const validationError of validationErrors) {
        const param = validationError.parameter;

        // set i18n error message
        this.error[param] = this.getI18nStringWithFallback(
          "settings_tls_certificates." + validationError.error,
          "error." + validationError.error
        );

        if (!focusAlreadySet) {
          this.focusElement(param);
          focusAlreadySet = true;
        }
      }
    },
    setCertificateCompleted(taskContext) {
      this.loading.setCertificate = false;
      this.clearFields();

      // remove pending certificate
      this.removePendingTlsCertificateInStore(taskContext.data.fqdn);

      // reload certificates
      this.$root.$emit("reloadCertificates");
    },
    clearFields() {
      this.fqdn = "";
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";

.mg-bottom-12 {
  margin-bottom: 12rem;
}
</style>

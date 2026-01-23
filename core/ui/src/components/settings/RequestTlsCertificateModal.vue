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
      certificate
        ? $t("settings_tls_certificates.manage_names")
        : $t("settings_tls_certificates.request_certificate")
    }}</template>
    <template slot="content">
      <cv-form @submit.prevent="setCertificate">
        <div v-show="loading.setCertificate">
          <NsEmptyState
            :title="
              $t('settings_tls_certificates.checking_certificate_requirements')
            "
            :animationData="GearsLottie"
            animationTitle="gears"
            :loop="true"
            class="checking-certificate"
          />
        </div>
        <div v-show="!loading.setCertificate">
          <NsInlineNotification
            v-if="certificate"
            kind="info"
            :showCloseButton="false"
          >
            <template #description>
              <p>
                {{ $t("settings_tls_certificates.manage_names_message_1") }}
              </p>
              <ul class="unordered-list">
                <li>
                  {{ $t("settings_tls_certificates.manage_names_message_2") }}
                </li>
                <li>
                  {{ $t("settings_tls_certificates.manage_names_message_3") }}
                </li>
              </ul>
            </template>
          </NsInlineNotification>
          <p v-else class="mb-6">
            {{
              $t("settings_tls_certificates.request_certificate_description")
            }}
          </p>
          <NsTextInput
            v-if="certificate"
            v-model="certificate.node"
            :label="$t('common.node')"
            disabled
            light
          />
          <NsComboBox
            v-else
            v-model="selectedNodeId"
            :label="$t('common.choose_a_node')"
            :title="$t('common.node')"
            :auto-filter="true"
            :auto-highlight="true"
            :options="nodes"
            :disabled="loading.setCertificate || !!certificate"
            :invalid-message="error.node"
            tooltipAlignment="start"
            tooltipDirection="right"
            light
            ref="node"
          >
            <template #tooltip>
              {{ $t("settings_tls_certificates.node_tooltip") }}
            </template>
          </NsComboBox>
          <div class="names">
            <label for="names" class="bx--label mb-1">
              <div class="label-with-tooltip">
                {{ $t("settings_tls_certificates.names") }}
                <cv-interactive-tooltip
                  alignment="start"
                  direction="right"
                  class="info"
                >
                  <template slot="content">
                    {{ $t("settings_tls_certificates.names_tooltip") }}
                  </template>
                </cv-interactive-tooltip>
              </div>
            </label>
            <cv-text-area
              id="names"
              v-model.trim="names"
              :invalid-message="$t(error.names)"
              :helper-text="$t('settings_tls_certificates.names_helper')"
              class="names"
              ref="names"
            >
            </cv-text-area>
          </div>
          <NsInlineNotification
            v-if="error.setCertificate"
            kind="error"
            :title="$t('settings_tls_certificates.cannot_obtain_certificate')"
            :description="error.setCertificate"
            :showCloseButton="false"
          />
          <NsInlineNotification
            v-if="validationErrorDetails.length"
            kind="error"
            :title="$t('settings_tls_certificates.cannot_obtain_certificate')"
            :showCloseButton="false"
          >
            <template #description>
              <div class="flex flex-col gap-2">
                <p>
                  {{ $t("settings_tls_certificates.validation_error_message") }}
                </p>
                <p
                  v-for="(detail, index) in validationErrorDetails"
                  :key="index"
                >
                  {{ detail }}
                </p>
              </div>
            </template>
          </NsInlineNotification>
        </div>
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
  LottieService,
} from "@nethserver/ns8-ui-lib";

export default {
  name: "RequestTlsCertificateModal",
  mixins: [UtilService, TaskService, DateTimeService, LottieService],
  props: {
    isShown: Boolean,
    nodes: {
      type: Array,
      required: true,
    },
    allCertificates: {
      type: Array,
      required: true,
    },
    certificate: {
      type: Object,
      default: null,
    },
  },
  data() {
    return {
      names: "",
      selectedNodeId: null,
      validationErrorDetails: [],
      loading: {
        setCertificate: false,
      },
      error: {
        setCertificate: "",
        names: "",
        node: "",
      },
    };
  },
  watch: {
    isShown: function () {
      if (this.isShown) {
        this.clearErrors();
        this.validationErrorDetails = [];

        if (this.certificate) {
          this.names = this.certificate.names.join("\n");

          this.$nextTick(() => {
            this.focusElement("names");
          });
        } else {
          this.clearFields();
        }
      }
    },
  },
  methods: {
    onModalHidden() {
      this.clearErrors();
      this.validationErrorDetails = [];
      this.$emit("hide");
    },
    validateSetCertificate() {
      this.clearErrors();
      this.validationErrorDetails = [];

      let isValidationOk = true;

      // names

      if (!this.names) {
        this.error.names = this.$t("common.required");

        if (isValidationOk) {
          this.focusElement("names");
          isValidationOk = false;
        }
      }

      // node

      if (!this.certificate && !this.selectedNodeId) {
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
      const taskAction = "set-default-certificate";
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

      let traefikInstance = null;
      let node = null;

      if (this.certificate) {
        traefikInstance = this.certificate.traefikInstance;
        node = this.certificate.longNodeLabel;
      } else {
        const selectedNode = this.nodes.find(
          (node) => node.value === this.selectedNodeId
        );
        traefikInstance = selectedNode.traefikInstance;
        node = selectedNode.label;
      }

      const logsStartDate = this.formatDate(new Date(), "yyyy-MM-dd");
      const logsHours = this.formatDate(new Date(), "HH");
      const logsMins = this.formatDate(new Date(), "mm");
      const logsStartTime = `${logsHours}%3A${logsMins}`;

      const res = await to(
        this.createModuleTaskForApp(traefikInstance, {
          action: taskAction,
          data: {
            names: this.names
              .split("\n")
              .filter((detail) => detail.trim() !== ""),
          },
          extra: {
            title: this.$t(
              "settings_tls_certificates.request_certificate_for_node",
              {
                node,
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
    },
    setCertificateAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.setCertificate = this.getSetCertificateErrorMessage(
        taskResult.error
      );
      this.loading.setCertificate = false;

      // reload certificates
      this.$root.$emit("reloadCertificates");
    },
    getSetCertificateErrorMessage(errorMessage) {
      if (errorMessage.includes("Invalid identifiers requested")) {
        return this.$t(
          "settings_tls_certificates.invalid_identifiers_requested_error"
        );
      } else if (
        /no valid A records found|check that a DNS record exists/.test(
          errorMessage
        )
      ) {
        return this.$t("settings_tls_certificates.dns_record_error");
      } else if (
        errorMessage.includes("refuses to issue a certificate for this domain")
      ) {
        return this.$t(
          "settings_tls_certificates.refuses_to_issue_certificate_error"
        );
      } else if (errorMessage.includes("tls: unrecognized name")) {
        return this.$t(
          "settings_tls_certificates.ensure_node_matches_fqdn_error"
        );
      } else {
        return this.$t(
          "settings_tls_certificates.set_certificate_unknown_error"
        );
      }
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

        if (validationError.details) {
          // show inline error notification with details
          this.validationErrorDetails = validationError.details
            .split("\n")
            .filter((detail) => detail.trim() !== "");
        } else {
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
      }
    },
    setCertificateCompleted() {
      this.loading.setCertificate = false;
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

.names {
  margin-bottom: 4rem;
}

.checking-certificate {
  margin-top: 4rem;
  margin-bottom: 4rem;
}

.label-with-tooltip {
  display: flex;
  align-items: baseline;
}

@media (min-width: $breakpoint-medium) {
  // limit text area width on large devices
  .cv-text-area {
    max-width: 38rem;
  }
}
</style>

<style lang="scss">
@import "../../styles/carbon-utils";

.names textarea {
  min-height: 7rem;
}

.names .bx--tooltip__label .bx--tooltip__trigger {
  margin-left: 0.25rem;
}
</style>

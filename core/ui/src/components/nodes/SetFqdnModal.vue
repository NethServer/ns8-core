<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <NsModal
    size="default"
    :visible="isShown"
    :isLoading="loading.setFqdn"
    :primaryButtonDisabled="loading.getFqdn || loading.setFqdn"
    @modal-hidden="onModalHidden"
    @primary-click="setFqdn"
  >
    <template slot="title">{{ $t("nodes.set_fqdn") }}</template>
    <template slot="content">
      <NsInlineNotification
        v-if="isLeaderNode"
        kind="warning"
        :title="$t('nodes.vpn_endpoint_warning_title')"
        :description="$t('nodes.vpn_endpoint_warning_description')"
        :showCloseButton="false"
      />
      <NsInlineNotification
        v-if="error.getFqdn"
        kind="error"
        :title="$t('action.get-fqdn')"
        :description="error.getFqdn"
        :showCloseButton="false"
      />
      <cv-skeleton-text
        v-if="loading.getFqdn"
        :paragraph="true"
        :line-count="5"
        heading
      ></cv-skeleton-text>
      <cv-form v-else @submit.prevent="setFqdn">
        <template>
          <cv-text-input
            :label="
              $t('nodes.node_name_hostname', {
                node: this.nodeLabel,
              })
            "
            v-model.trim="hostname"
            :invalid-message="error.hostname"
            :disabled="loading.getFqdn || loading.setFqdn"
            ref="hostname"
          >
          </cv-text-input>
          <cv-text-input
            :label="
              $t('nodes.node_name_domain', {
                node: this.nodeLabel,
              })
            "
            v-model.trim="domain"
            placeholder="example.org"
            :invalid-message="error.domain"
            :disabled="loading.getFqdn || loading.setFqdn"
            ref="domain"
          >
          </cv-text-input>
          <!-- check node connectivity -->
          <NsCheckbox
            v-if="isLeaderNode"
            :label="$t('nodes.check_node_connectivity')"
            v-model="checkNodeConnectivity"
            :disabled="loading.getFqdn || loading.setFqdn"
            value="checkNodeConnectivity"
          />
        </template>
      </cv-form>
      <NsInlineNotification
        v-if="error.checkNodeConnectivity"
        kind="error"
        :title="$t('nodes.connectivity_check_failed')"
        :description="error.checkNodeConnectivity"
        :showCloseButton="false"
      />
      <NsInlineNotification
        v-if="error.setFqdn"
        kind="error"
        :title="$t('action.set-fqdn')"
        :description="error.setFqdn"
        :showCloseButton="false"
      />
    </template>
    <template slot="secondary-button">{{ $t("common.cancel") }}</template>
    <template slot="primary-button">{{ $t("nodes.set_fqdn") }}</template>
  </NsModal>
</template>

<script>
import { UtilService, IconService, TaskService } from "@nethserver/ns8-ui-lib";
import to from "await-to-js";
import { mapState } from "vuex";

export default {
  name: "SetFqdnModal",
  components: {},
  mixins: [UtilService, IconService, TaskService],
  props: {
    isShown: {
      type: Boolean,
      default: true,
    },
    node: { type: Object },
  },
  data() {
    return {
      hostname: "",
      domain: "",
      checkNodeConnectivity: true,
      loading: {
        getFqdn: false,
        setFqdn: false,
      },
      error: {
        getFqdn: "",
        setFqdn: "",
        hostname: "",
        domain: "",
        checkNodeConnectivity: "",
      },
    };
  },
  computed: {
    ...mapState(["clusterNodes"]),
    isLeaderNode() {
      return this.node && this.node.local;
    },
    nodeLabel() {
      if (this.node) {
        return this.getShortNodeLabel(this.node);
      } else {
        return "";
      }
    },
  },
  watch: {
    isShown: function () {
      if (this.isShown) {
        this.clearErrors();
        this.getFqdn();
        this.checkNodeConnectivity = true;
      }
    },
  },
  methods: {
    onModalHidden() {
      this.$emit("hide");
    },
    async getFqdn() {
      this.error.getFqdn = "";
      this.loading.getFqdn = true;
      const taskAction = "get-fqdn";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(`${taskAction}-aborted-${eventId}`, this.getFqdnAborted);

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.getFqdnCompleted
      );

      const res = await to(
        this.createNodeTask(this.node.id, {
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
        this.error.getFqdn = this.getErrorMessage(err);
        return;
      }
    },
    getFqdnCompleted(taskContext, taskResult) {
      this.hostname = taskResult.output.hostname;
      this.domain = taskResult.output.domain;
      this.loading.getFqdn = false;

      setTimeout(() => {
        this.focusElement("hostname");
      }, 300);
    },
    getFqdnAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.loading.getFqdn = false;

      // hide modal
      this.$emit("hide");
    },
    validateSetFqdn() {
      this.error.hostname = "";
      this.error.domain = "";
      let isValidationOk = true;

      // hostname

      if (!this.hostname) {
        this.error.hostname = this.$t("common.required");

        if (isValidationOk) {
          this.focusElement("hostname");
          isValidationOk = false;
        }
      }

      // domain

      if (!this.domain) {
        this.error.domain = this.$t("common.required");

        if (isValidationOk) {
          this.focusElement("domain");
          isValidationOk = false;
        }
      }
      return isValidationOk;
    },
    async setFqdn() {
      if (!this.validateSetFqdn()) {
        return;
      }
      this.loading.setFqdn = true;
      this.error.setFqdn = "";
      const taskAction = "set-fqdn";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(`${taskAction}-aborted-${eventId}`, this.setFqdnAborted);

      // register to task validation
      this.$root.$once(
        `${taskAction}-validation-failed-${eventId}`,
        this.setFqdnValidationFailed
      );
      this.$root.$once(
        `${taskAction}-validation-ok-${eventId}`,
        this.setFqdnValidationOk
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.setFqdnCompleted
      );

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            node: this.node.id,
            hostname: this.hostname,
            domain: this.domain,
            reachability_check: this.checkNodeConnectivity,
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
        this.error.setFqdn = this.getErrorMessage(err);
        this.loading.setFqdn = false;
        return;
      }
    },
    setFqdnAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.loading.setFqdn = false;

      // hide modal
      this.$emit("hide");
    },
    setFqdnValidationOk() {
      // hide modal
      this.$emit("hide");
    },
    setFqdnValidationFailed(validationErrors) {
      this.loading.setFqdn = false;
      let focusAlreadySet = false;

      for (const validationError of validationErrors) {
        const param = validationError.parameter;

        if (validationError.error == "reachability_check_failed") {
          // show error notification for connectivity check
          const failedNodeId = validationError.value;
          const failedNode = this.clusterNodes.find(
            (node) => node.id == failedNodeId
          );

          if (failedNode) {
            const failedNodeLabel = this.getNodeLabel(failedNode);
            this.error.checkNodeConnectivity = this.$t(
              "nodes.connectivity_check_failed_description",
              { fqdn: `${this.hostname}.${this.domain}`, node: failedNodeLabel }
            );
          }
        } else {
          // set i18n error message
          this.error[param] = this.getI18nStringWithFallback(
            "init." + validationError.error,
            "error." + validationError.error
          );

          if (!focusAlreadySet) {
            this.focusElement(param);
            focusAlreadySet = true;
          }
        }
      }
    },
    setFqdnCompleted() {
      this.loading.setFqdn = false;

      // hide modal
      this.$emit("hide");
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";
</style>

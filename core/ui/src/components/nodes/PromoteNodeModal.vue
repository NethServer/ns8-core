<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <NsModal
    size="default"
    :visible="isShown"
    :isLoading="loading.promoteNode"
    @modal-hidden="onModalHidden"
    @primary-click="promoteNode"
    :primary-button-disabled="userInput !== 'node' + (node ? node.id : '')"
  >
    <template slot="title">{{ $t("nodes.promote_node_to_leader") }}</template>
    <template slot="content">
      <NsInlineNotification
        kind="info"
        :title="$t('common.remember')"
        :description="$t('nodes.reset_backup_encryption_password_info')"
        :showCloseButton="false"
      />
      <!-- need to wrap error notification inside a div: custom elements like NsInlineNotification don't have scrollIntoView() function -->
      <div ref="getNodeInfoError">
        <NsInlineNotification
          v-if="error.getNodeInfo"
          kind="error"
          :title="$t('nodes.cannot_retrieve_node_info')"
          :description="error.getNodeInfo"
          :showCloseButton="false"
          class="mg-top-lg"
        />
      </div>
      <div>
        {{
          $t("nodes.promote_node_confirm", {
            name: node ? this.getNodeLabel(node) : "",
          })
        }}
      </div>
      <cv-skeleton-text
        v-if="loading.getNodeInfo"
        :paragraph="true"
        :line-count="10"
        heading
        class="mg-top-lg"
      ></cv-skeleton-text>
      <cv-form v-else @submit.prevent="promoteNode" class="mg-top-lg">
        <NsCheckbox
          :label="
            $t('nodes.check_node_connectivity_promote', {
              fqdn: this.vpnEndpointAddress,
            })
          "
          v-model="checkNodeConnectivity"
          :disabled="loading.promoteNode"
          value="checkNodeConnectivity"
        />
        <div class="mg-top-xlg">
          {{
            $t("nodes.type_to_confirm_promotion", {
              name: "node" + (node ? node.id : ""),
            })
          }}
        </div>
        <NsTextInput
          v-model="userInput"
          :disabled="loading.promoteNode"
          class="mg-bottom-md"
          ref="userInput"
        >
        </NsTextInput>
      </cv-form>
      <!-- need to wrap error notification inside a div: custom elements like NsInlineNotification don't have scrollIntoView() function -->
      <div ref="promoteNodeError">
        <NsInlineNotification
          v-if="error.promoteNode"
          kind="error"
          :title="$t('nodes.cannot_promote_node')"
          :description="error.promoteNode"
          :showCloseButton="false"
          class="mg-top-lg"
        />
      </div>
    </template>
    <template slot="secondary-button">{{ $t("common.cancel") }}</template>
    <template slot="primary-button">{{
      $t("nodes.understood_promote_node")
    }}</template>
  </NsModal>
</template>

<script>
import { TaskService, UtilService } from "@nethserver/ns8-ui-lib";
import to from "await-to-js";

export default {
  name: "PromoteNodeModal",
  mixins: [TaskService, UtilService],
  props: {
    isShown: {
      type: Boolean,
      default: true,
    },
    node: { type: Object },
  },
  data() {
    return {
      vpnEndpointAddress: "",
      vpnEndpointPort: "",
      userInput: "",
      checkNodeConnectivity: true,
      loading: {
        promoteNode: false,
        getNodeInfo: false,
      },
      error: {
        promoteNode: "",
        getNodeInfo: "",
        endpoint_address: "",
        endpoint_port: "",
      },
    };
  },
  watch: {
    isShown: function () {
      if (this.isShown) {
        this.userInput = "";
        this.clearErrors();
        this.getNodeInfo();
      }
    },
    "error.promoteNode": function () {
      if (this.error.promoteNode) {
        // scroll to notification error

        this.$nextTick(() => {
          const el = this.$refs.promoteNodeError;
          this.scrollToElement(el);
        });
      }
    },
    "error.getNodeInfo": function () {
      if (this.error.getNodeInfo) {
        // scroll to notification error

        this.$nextTick(() => {
          const el = this.$refs.getNodeInfoError;
          this.scrollToElement(el);
        });
      }
    },
  },
  methods: {
    onModalHidden() {
      this.$emit("hide");
    },
    async getNodeInfo() {
      this.error.getNodeInfo = "";
      this.loading.getNodeInfo = true;
      const taskAction = "get-info";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.getNodeInfoAborted
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.getNodeInfoCompleted
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
        this.error.getNodeInfo = this.getErrorMessage(err);
        return;
      }
    },
    getNodeInfoAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.getNodeInfo = this.$t("error.generic_error");
      this.loading.getNodeInfo = false;
    },
    getNodeInfoCompleted(taskContext, taskResult) {
      this.vpnEndpointAddress = taskResult.output.hostname;
      this.vpnEndpointPort = taskResult.output.vpn_port;
      this.loading.getNodeInfo = false;
    },
    async promoteNode() {
      if (this.userInput !== "node" + this.node.id) {
        return;
      }
      this.clearErrors();
      this.loading.promoteNode = true;
      this.error.promoteNode = "";
      const taskAction = "promote-node";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.promoteNodeAborted
      );

      // register to task validation
      this.$root.$once(
        `${taskAction}-validation-failed-${eventId}`,
        this.promoteNodeValidationFailed
      );
      this.$root.$once(
        `${taskAction}-validation-ok-${eventId}`,
        this.promoteNodeValidationOk
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.promoteNodeCompleted
      );

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            node_id: this.node.id,
            endpoint_address: this.vpnEndpointAddress,
            endpoint_port: this.vpnEndpointPort,
            endpoint_validation: this.checkNodeConnectivity,
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
        this.error.promoteNode = this.getErrorMessage(err);
        this.loading.promoteNode = false;
        return;
      }
    },
    promoteNodeAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.loading.promoteNode = false;
    },
    promoteNodeValidationOk() {
      // hide modal
      this.$emit("hide");
    },
    promoteNodeValidationFailed(validationErrors) {
      this.loading.promoteNode = false;
      let focusAlreadySet = false;

      for (const validationError of validationErrors) {
        if (
          validationError.parameter === "node_id" ||
          validationError.error === "endpoint_address_not_reachable"
        ) {
          this.error.promoteNode = this.$t(
            "nodes.validation_" + validationError.error
          );
        } else {
          const param = validationError.parameter;

          this.error[param] = this.$t(
            "nodes.validation_" + validationError.error
          );

          if (!focusAlreadySet) {
            this.focusElement(param);
            focusAlreadySet = true;
          }
        }
      }
    },
    promoteNodeCompleted() {
      this.loading.promoteNode = false;

      // hide modal
      this.$emit("hide");

      // show new leader modal
      this.$emit("nodePromoted", this.vpnEndpointAddress);
    },
  },
};
</script>

<style scoped lang="scss"></style>

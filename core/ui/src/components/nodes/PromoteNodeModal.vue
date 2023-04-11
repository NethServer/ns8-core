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
      <div
        v-html="
          $t('nodes.promote_node_confirm', {
            name: node ? this.getNodeLabel(node) : '',
          })
        "
      ></div>
      <cv-form @submit.prevent="promoteNode" class="mg-top-lg">
        <NsTextInput
          :label="
            $t('nodes.node_vpn_endpoint_address', {
              node: node ? this.getNodeLabel(node) : '',
            })
          "
          v-model.trim="internalVpnEndpointAddress"
          :invalid-message="$t(error.vpnEndpointAddress)"
          :disabled="loading.promoteNode"
          ref="vpnEndpointAddress"
        >
          <template #tooltip>{{
            $t("nodes.node_vpn_endpoint_address_tooltip")
          }}</template>
        </NsTextInput>
        <NsTextInput
          :label="
            $t('nodes.node_vpn_endpoint_port', {
              node: node ? this.getNodeLabel(node) : '',
            })
          "
          v-model.trim="vpnEndpointPort"
          :invalid-message="$t(error.vpnEndpointPort)"
          :disabled="loading.promoteNode"
          type="number"
          class="narrow"
          ref="vpnEndpointPort"
        >
        </NsTextInput>
        <NsCheckbox
          :label="$t('nodes.check_node_connectivity')"
          v-model="checkNodeConnectivity"
          :disabled="loading.promoteNode"
          value="checkNodeConnectivity"
        />
        <div
          class="mg-top-xlg"
          v-html="
            $t('nodes.type_to_confirm_promotion', {
              name: 'node' + (node ? node.id : ''),
            })
          "
        ></div>
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
    vpnEndpointAddress: { type: String, required: true },
  },
  data() {
    return {
      internalVpnEndpointAddress: "",
      vpnEndpointPort: "55820",
      userInput: "",
      checkNodeConnectivity: true,
      loading: {
        promoteNode: false,
      },
      error: {
        promoteNode: "",
        vpnEndpointAddress: "",
        vpnEndpointPort: "",
      },
    };
  },
  watch: {
    isShown: function () {
      if (this.isShown) {
        this.userInput = "";
        this.internalVpnEndpointAddress = this.vpnEndpointAddress;
        this.clearErrors();

        setTimeout(() => {
          this.focusElement("vpnEndpointAddress");
        }, 300);
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
  },
  methods: {
    onModalHidden() {
      this.$emit("hide");
    },
    validatePromoteNode() {
      this.clearErrors();
      let isValidationOk = true;

      // VPN endpoint address

      if (!this.internalVpnEndpointAddress) {
        this.error.vpnEndpointAddress = this.$t("common.required");

        if (isValidationOk) {
          this.focusElement("vpnEndpointAddress");
          isValidationOk = false;
        }
      }

      // VPN endpoint port

      if (!this.vpnEndpointPort) {
        this.error.vpnEndpointPort = this.$t("common.required");

        if (isValidationOk) {
          this.focusElement("vpnEndpointPort");
          isValidationOk = false;
        }
      }
      return isValidationOk;
    },
    async promoteNode() {
      if (this.userInput !== "node" + this.node.id) {
        return;
      }

      const isValidationOk = this.validatePromoteNode();
      if (!isValidationOk) {
        return;
      }
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
            endpoint_address: this.internalVpnEndpointAddress,
            endpoint_port: Number(this.vpnEndpointPort),
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

      for (const validationError of validationErrors) {
        // set i18n error message
        this.error.promoteNode = this.$t(
          "nodes.validation_" + validationError.error
        );
      }
    },
    promoteNodeCompleted() {
      this.loading.promoteNode = false;

      // hide modal
      this.$emit("hide");

      // show new leader modal
      this.$emit("nodePromoted", this.internalVpnEndpointAddress);
    },
  },
};
</script>

<style scoped lang="scss"></style>

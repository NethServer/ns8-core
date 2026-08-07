<!--
  Copyright (C) 2026 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <NsModal
    size="default"
    :visible="isShown"
    :primary-button-disabled="loading.setTrustedProxies"
    :isLoading="loading.setTrustedProxies"
    @modal-hidden="onModalHidden"
    @primary-click="setTrustedProxies"
  >
    <template slot="title">{{
      isEditing
        ? $t("settings_http_routes.edit_frontend_proxy_node", {
            node: proxyConfig ? proxyConfig.nodeLabel : "",
          })
        : $t("settings_http_routes.add_frontend_proxy")
    }}</template>
    <template slot="content">
      <cv-form @submit.prevent="setTrustedProxies">
        <!-- when editing, the node is fixed and already named in the modal title -->
        <!-- the key forces a remount: NsComboBox keeps its previous text when the
             value is cleared, since it only repaints on a matching option -->
        <NsComboBox
          v-if="!isEditing"
          :key="nodeComboBoxKey"
          v-model="selectedNodeId"
          :label="$t('settings_http_routes.choose_node')"
          :title="$t('common.node')"
          :auto-filter="true"
          :auto-highlight="true"
          :options="nodes"
          :disabled="loading.setTrustedProxies"
          :invalid-message="error.node"
          light
          ref="node"
        />
        <cv-text-area
          v-model.trim="proxies_str"
          :placeholder="$t('settings_http_routes.frontend_proxies_placeholder')"
          :label="$t('settings_http_routes.frontend_proxies')"
          :helper-text="$t('settings_http_routes.frontend_proxies_helper')"
          :invalid-message="error.proxies"
          :disabled="loading.setTrustedProxies"
          rows="4"
          ref="proxies"
        />
        <!-- cv-number-input has no tooltip slot: label markup copied from NsTextInput -->
        <div class="bx--label label-with-tooltip">
          {{ $t("settings_http_routes.trust_depth") }}
          <cv-interactive-tooltip
            alignment="start"
            direction="top"
            class="info"
          >
            <template slot="content">
              {{ $t("settings_http_routes.trust_depth_tooltip") }}
            </template>
          </cv-interactive-tooltip>
        </div>
        <cv-number-input
          v-model="depth"
          :label="''"
          :aria-label="$t('settings_http_routes.trust_depth')"
          :helper-text="$t('settings_http_routes.trust_depth_helper')"
          :invalid-message="error.depth"
          :disabled="loading.setTrustedProxies"
          :min="0"
          :step="1"
          class="trust-depth"
          ref="depth"
        />
        <NsInlineNotification
          v-if="selectedNodeLabel"
          kind="warning"
          :title="$t('settings_http_routes.traefik_will_be_restarted')"
          :description="
            $t('settings_http_routes.frontend_proxy_restart_message', {
              node: selectedNodeLabel,
            })
          "
          :showCloseButton="false"
        />
        <!-- need to wrap error notification inside a div: custom elements like NsInlineNotification don't have scrollIntoView() function -->
        <div ref="setTrustedProxiesError">
          <NsInlineNotification
            v-if="error.setTrustedProxies"
            kind="error"
            :title="$t('action.set-trusted-proxies')"
            :description="error.setTrustedProxies"
            :showCloseButton="false"
          />
        </div>
      </cv-form>
    </template>
    <template slot="secondary-button">{{ $t("common.cancel") }}</template>
    <template slot="primary-button">{{
      isEditing ? $t("common.save") : $t("common.add")
    }}</template>
  </NsModal>
</template>

<script>
import to from "await-to-js";
import { UtilService, TaskService } from "@nethserver/ns8-ui-lib";
import { mapState } from "vuex";

// set-trusted-proxies validates with ipaddress.ip_address(): no CIDR ranges, and
// no leading-zero octets, which Python rejects since 3.9.5
const IPV4_OCTET = "(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])";
const IPV4 = `${IPV4_OCTET}(\\.${IPV4_OCTET}){3}`;
const IPV4_PATTERN = new RegExp(`^${IPV4}$`);
const IPV6_PATTERN = new RegExp(
  `^(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}${IPV4}|([0-9a-fA-F]{1,4}:){1,4}:${IPV4})$`
);

export default {
  name: "ConfigureFrontendProxyModal",
  mixins: [UtilService, TaskService],
  props: {
    isShown: Boolean,
    nodes: {
      type: Array,
      required: true,
    },
    proxyConfig: { type: [Object, null], default: null },
    isEditing: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      selectedNodeId: "",
      nodeComboBoxKey: 0,
      proxies_str: "",
      depth: 1,
      loading: {
        setTrustedProxies: false,
      },
      error: {
        node: "",
        proxies: "",
        depth: "",
        setTrustedProxies: "",
      },
    };
  },
  computed: {
    ...mapState(["isWebsocketConnected"]),
    selectedNode() {
      return this.nodes.find((node) => node.value === this.selectedNodeId);
    },
    selectedNodeLabel() {
      return this.selectedNode ? this.selectedNode.label : "";
    },
  },
  watch: {
    isWebsocketConnected: function (isConnected) {
      // safety net: nothing else clears the pending state if the request was cut off
      if (isConnected && this.loading.setTrustedProxies) {
        this.loading.setTrustedProxies = false;
        this.$emit("hide");
      }
    },
    isShown: function () {
      if (this.isShown) {
        this.clearErrors();

        if (this.isEditing && this.proxyConfig) {
          this.selectedNodeId = this.proxyConfig.nodeId;
          this.proxies_str = this.proxyConfig.proxies_str;
          // show the depth as returned by the backend, even if it is 0
          this.depth = this.proxyConfig.depth;
        } else {
          this.selectedNodeId = "";
          this.proxies_str = "";
          this.depth = 1;
          this.nodeComboBoxKey++;
        }
      }
    },
    "error.setTrustedProxies": function () {
      if (this.error.setTrustedProxies) {
        // scroll to notification error

        this.$nextTick(() => {
          const el = this.$refs.setTrustedProxiesError;
          this.scrollToElement(el);
        });
      }
    },
  },
  methods: {
    onModalHidden() {
      this.clearErrors();
      this.$emit("hide");
    },
    getProxyList() {
      return this.proxies_str
        .split("\n")
        .map((proxy) => proxy.trim())
        .filter((proxy) => proxy !== "");
    },
    validateSetTrustedProxies() {
      this.clearErrors();

      let isValidationOk = true;

      // the node field is only rendered when adding
      // check selectedNode, not selectedNodeId: the parent drops nodes from the
      // options as soon as they get a configuration, even while the modal is open

      if (!this.isEditing && !this.selectedNode) {
        this.error.node = this.$t("common.required");

        if (isValidationOk) {
          this.focusElement("node");
          isValidationOk = false;
        }
      }

      // frontend proxies

      const proxyList = this.getProxyList();

      if (!proxyList.length) {
        this.error.proxies = this.$t("common.required");

        if (isValidationOk) {
          this.focusElement("proxies");
          isValidationOk = false;
        }
      } else {
        for (const proxy of proxyList) {
          if (!IPV4_PATTERN.test(proxy) && !IPV6_PATTERN.test(proxy)) {
            const isIPv6Like = proxy.includes(":");
            this.error.proxies = this.$t(
              isIPv6Like
                ? "settings_http_routes.invalid_ipv6"
                : "settings_http_routes.invalid_ipv4",
              { ip: proxy }
            );

            if (isValidationOk) {
              this.focusElement("proxies");
              isValidationOk = false;
            }
            break;
          }
        }
      }

      // trust depth

      const depth = Number(this.depth);

      // cv-number-input emits NaN, not "", when a numeric value is cleared
      if (this.depth === "" || this.depth === null || Number.isNaN(depth)) {
        this.error.depth = this.$t("common.required");

        if (isValidationOk) {
          this.focusElement("depth");
          isValidationOk = false;
        }
      } else if (!Number.isInteger(depth) || depth < 0) {
        // the input schema requires an integer >= 0
        this.error.depth = this.$t("settings_http_routes.trust_depth_invalid");

        if (isValidationOk) {
          this.focusElement("depth");
          isValidationOk = false;
        }
      }

      return isValidationOk;
    },
    async setTrustedProxies() {
      if (!this.validateSetTrustedProxies()) {
        return;
      }
      // when editing, the node combo box is not rendered: read the instance from the row
      const traefikInstanceId = this.isEditing
        ? this.proxyConfig.traefikInstance
        : this.selectedNode.traefikInstance;
      this.loading.setTrustedProxies = true;
      this.error.setTrustedProxies = "";
      const taskAction = "set-trusted-proxies";
      const eventId = this.getUuid();

      // register to task events

      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.setTrustedProxiesAborted
      );

      this.$root.$once(
        `${taskAction}-validation-failed-${eventId}`,
        this.setTrustedProxiesValidationFailed
      );

      this.$root.$once(
        `${taskAction}-validation-ok-${eventId}`,
        this.setTrustedProxiesValidationOk
      );

      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.setTrustedProxiesCompleted
      );

      const res = await to(
        this.createModuleTaskForApp(traefikInstanceId, {
          action: taskAction,
          // always send depth: the backend resets it to 1 if omitted
          data: {
            proxies: this.getProxyList(),
            depth: Number(this.depth),
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
        this.error.setTrustedProxies = this.getErrorMessage(err);
        this.loading.setTrustedProxies = false;
        return;
      }
      // stay open until validation completes: closing now would hide the field errors
    },
    setTrustedProxiesAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.loading.setTrustedProxies = false;

      // an abort after validation is reported by the task notification
      this.$emit("hide");
    },
    setTrustedProxiesValidationOk() {
      // the restart kills the websocket, so no completion event will arrive
      this.loading.setTrustedProxies = false;
      this.$emit("hide");
    },
    setTrustedProxiesValidationFailed(validationErrors) {
      this.loading.setTrustedProxies = false;
      let focusAlreadySet = false;

      for (const validationError of validationErrors) {
        const param = validationError.parameter;

        // set i18n error message
        this.error[param] = this.getI18nStringWithFallback(
          "settings_http_routes." + validationError.error,
          "error." + validationError.error
        );

        if (!focusAlreadySet) {
          this.focusElement(param);
          focusAlreadySet = true;
        }
      }
    },
    setTrustedProxiesCompleted() {
      this.loading.setTrustedProxies = false;
      this.$emit("reloadTrustedProxies");
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";

.label-with-tooltip {
  display: flex;
  align-items: center;
}

// the inline svg sits on the text baseline: make the trigger a flex box
.label-with-tooltip ::v-deep .bx--tooltip__trigger {
  display: flex;
  margin-left: 0.25rem;
}

// align with the form fields, capped at the same width globally
.bx--inline-notification {
  max-width: 38rem;
}

// wide enough to render the helper text on a single line
.trust-depth {
  max-width: 15rem;
}

// Carbon renders number inputs in monospace 0.75rem: match the other fields
.trust-depth ::v-deep input[type="number"] {
  font-family: inherit;
  font-size: 0.875rem;
  line-height: 1.42857;
  letter-spacing: 0.16px;
}
</style>

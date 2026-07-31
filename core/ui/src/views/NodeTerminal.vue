<!--
  Copyright (C) 2026 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <div>
    <cv-grid fullWidth>
      <cv-row>
        <cv-column class="page-title">
          <h2>{{ $t("terminal.title") }}</h2>
        </cv-column>
      </cv-row>
      <cv-row>
        <cv-column>
          <NsInlineNotification
            v-if="error.listNodes"
            kind="error"
            :title="$t('action.list-nodes')"
            :description="error.listNodes"
            :showCloseButton="false"
          />
          <cv-skeleton-text
            v-if="loading.listNodes"
            :paragraph="true"
            :line-count="3"
            heading
          ></cv-skeleton-text>
          <template v-else>
            <cv-form @submit.prevent="requestSession">
              <cv-select
                v-model="selectedNodeId"
                :label="$t('terminal.node')"
                :disabled="isConnected"
              >
                <cv-select-option
                  v-for="node in nodes"
                  :key="node.node_id"
                  :value="String(node.node_id)"
                >
                  {{ nodeLabelOf(node) }}
                </cv-select-option>
              </cv-select>
            </cv-form>

            <!-- The switch lives next to the terminal on purpose: the person
                 able to fix sshd is the one looking at this page. -->
            <template v-if="selectedNode">
              <NsInlineNotification
                v-if="!selectedNode.terminal_enabled"
                kind="info"
                :title="$t('terminal.disabled_on_node')"
                :description="$t('terminal.enable_warning')"
                :showCloseButton="false"
              />
              <cv-toggle
                :value="'terminal-enabled'"
                v-model="terminalEnabled"
                :disabled="loading.toggle || isConnected"
                @change="onToggleTerminal"
              >
                <template slot="text-left">{{ $t("terminal.off") }}</template>
                <template slot="text-right">{{ $t("terminal.on") }}</template>
              </cv-toggle>
              <NsInlineNotification
                v-if="error.toggle"
                kind="error"
                :title="$t('terminal.cannot_change_state')"
                :description="error.toggle"
                :showCloseButton="false"
              />
            </template>

            <div class="terminal-actions">
              <NsButton
                kind="primary"
                :icon="Terminal20"
                :disabled="!canRequest"
                @click="requestSession"
                v-if="!isConnected"
                >{{ $t("terminal.open") }}</NsButton
              >
              <NsButton
                kind="danger--tertiary"
                :icon="Close20"
                @click="closeSession"
                v-else
                >{{ $t("terminal.close") }}</NsButton
              >
            </div>

            <NsInlineNotification
              v-if="error.session"
              kind="error"
              :title="$t('terminal.cannot_open')"
              :description="error.session"
              :showCloseButton="false"
            />
            <NsInlineNotification
              v-if="closedReason"
              kind="warning"
              :title="$t('terminal.session_closed')"
              :description="closedReason"
              :showCloseButton="false"
            />

            <XtermPane
              v-if="socket && isConnected"
              ref="pane"
              :socket="socket"
            />
          </template>
        </cv-column>
      </cv-row>
    </cv-grid>

    <OpenTerminalModal
      v-if="selectedNode"
      :isShown="isModalShown"
      :nodeId="selectedNodeId"
      :nodeLabel="nodeLabelOf(selectedNode)"
      @open="onCredentials"
      @hide="onModalHidden"
    />
  </div>
</template>

<script>
import to from "await-to-js";
import Terminal20 from "@carbon/icons-vue/es/terminal/20";
import Close20 from "@carbon/icons-vue/es/close/20";
import {
  TaskService,
  UtilService,
  StorageService,
  IconService,
  PageTitleService,
} from "@nethserver/ns8-ui-lib";
import NodeService from "@/mixins/node";
import OpenTerminalModal from "@/components/terminal/OpenTerminalModal";
import XtermPane from "@/components/terminal/XtermPane";

export default {
  name: "NodeTerminal",
  components: { OpenTerminalModal, XtermPane },
  mixins: [
    TaskService,
    UtilService,
    StorageService,
    IconService,
    PageTitleService,
    NodeService,
  ],
  data() {
    return {
      Terminal20,
      Close20,
      nodes: [],
      selectedNodeId: "",
      terminalEnabled: [],
      socket: null,
      isModalShown: false,
      isConnected: false,
      closedReason: "",
      loading: {
        listNodes: false,
        toggle: false,
      },
      error: {
        listNodes: "",
        toggle: "",
        session: "",
      },
    };
  },
  computed: {
    selectedNode() {
      return this.nodes.find(
        (node) => String(node.node_id) === String(this.selectedNodeId)
      );
    },
    canRequest() {
      return (
        !!this.selectedNode &&
        this.selectedNode.terminal_enabled &&
        !this.loading.toggle
      );
    },
  },
  beforeRouteLeave(to, from, next) {
    this.closeSession();
    next();
  },
  created() {
    this.listNodes();
  },
  beforeDestroy() {
    this.closeSession();
  },
  methods: {
    nodeLabelOf(node) {
      const name = node.ui_name || node.fqdn || "";
      return name ? `${node.node_id} - ${name}` : String(node.node_id);
    },
    async listNodes() {
      this.loading.listNodes = true;
      this.error.listNodes = "";

      const taskAction = "list-nodes";
      const eventId = this.getUuid();

      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.listNodesCompleted
      );
      this.$root.$once(`${taskAction}-aborted-${eventId}`, (taskResult) => {
        console.error(`${taskAction} aborted`, taskResult);
        this.error.listNodes = this.$t("terminal.cannot_list_nodes");
        this.loading.listNodes = false;
      });

      const res = await to(
        this.createClusterTask({
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
        this.error.listNodes = this.getErrorMessage(err);
        this.loading.listNodes = false;
      }
    },
    listNodesCompleted(taskContext, taskResult) {
      this.nodes = (taskResult.output.nodes || []).filter(
        (node) => node.role !== "ns7migration"
      );
      if (!this.selectedNodeId && this.nodes.length) {
        this.selectedNodeId = String(this.nodes[0].node_id);
      }
      this.syncToggle();
      this.loading.listNodes = false;
    },
    syncToggle() {
      this.terminalEnabled =
        this.selectedNode && this.selectedNode.terminal_enabled
          ? ["terminal-enabled"]
          : [];
    },
    async onToggleTerminal() {
      const wanted = this.terminalEnabled.includes("terminal-enabled");
      const taskAction = wanted
        ? "enable-node-terminal"
        : "disable-node-terminal";

      this.loading.toggle = true;
      this.error.toggle = "";

      const eventId = this.getUuid();
      this.$root.$once(`${taskAction}-completed-${eventId}`, () => {
        this.loading.toggle = false;
        this.listNodes();
      });
      this.$root.$once(`${taskAction}-aborted-${eventId}`, (taskResult) => {
        console.error(`${taskAction} aborted`, taskResult);
        // The flag may be set while sshd was not reconfigured: reload rather
        // than trusting the switch.
        this.error.toggle = this.$t("terminal.toggle_failed");
        this.loading.toggle = false;
        this.listNodes();
      });

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: { node_id: Number(this.selectedNodeId) },
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
        this.error.toggle = this.getErrorMessage(err);
        this.loading.toggle = false;
      }
    },
    requestSession() {
      this.error.session = "";
      this.closedReason = "";
      this.isModalShown = true;
    },
    onModalHidden() {
      this.isModalShown = false;
    },
    /**
     * Two steps on purpose. The REST call only authorizes and returns a
     * one-shot ticket; credentials travel on the websocket, where the handshake
     * runs. That leaves room for a browser-held signer later without changing
     * this flow.
     */
    async onCredentials(credentials) {
      this.isModalShown = false;
      this.error.session = "";

      const token = this.getFromStorage("loginInfo")
        ? this.getFromStorage("loginInfo").token
        : "";

      const res = await to(
        this.axios.post(
          `${this.$root.apiUrl}/node/${this.selectedNodeId}/terminal-sessions`,
          { action: "open-terminal" },
          { headers: { Authorization: `Bearer ${token}` } }
        )
      );
      const err = res[0];
      const response = res[1];

      if (err) {
        this.error.session = this.getErrorMessage(err);
        return;
      }

      this.openSocket(response.data.data.ticket, credentials);
    },
    openSocket(ticket, credentials) {
      const socket = new WebSocket(`${this.$root.config.WS_ENDPOINT}/terminal`);
      socket.binaryType = "arraybuffer";
      this.socket = socket;

      socket.addEventListener("open", () => {
        socket.send(JSON.stringify({ type: "ticket", ticket }));
      });

      socket.addEventListener("message", (event) => {
        if (typeof event.data !== "string") {
          // Terminal bytes: XtermPane consumes them.
          return;
        }
        this.onControlFrame(JSON.parse(event.data), credentials);
      });

      socket.addEventListener("close", () => {
        this.isConnected = false;
        this.socket = null;
      });

      socket.addEventListener("error", () => {
        this.error.session = this.$t("terminal.transport_error");
      });
    },
    onControlFrame(frame, credentials) {
      switch (frame.type) {
        case "ticket-accepted": {
          const size = this.$refs.pane
            ? this.$refs.pane.currentSize()
            : { rows: 24, cols: 80 };
          this.socket.send(
            JSON.stringify({
              type: "credentials",
              username: credentials.username,
              password: credentials.password,
              rows: size.rows,
              cols: size.cols,
            })
          );
          // Drop the password as soon as it is on the wire.
          credentials.password = "";
          break;
        }
        case "ready":
          this.isConnected = true;
          break;
        case "auth-error":
          this.error.session = frame.retry_after
            ? this.$t("terminal.retry_after", { seconds: frame.retry_after })
            : frame.message;
          break;
        case "closed":
          this.closedReason = frame.reason;
          if (this.$refs.pane) {
            this.$refs.pane.writeNotice(frame.reason);
          }
          break;
      }
    },
    closeSession() {
      if (this.socket) {
        this.socket.close();
        this.socket = null;
      }
      this.isConnected = false;
    },
  },
};
</script>

<style scoped lang="scss">
.terminal-actions {
  margin: 1rem 0;
}
</style>

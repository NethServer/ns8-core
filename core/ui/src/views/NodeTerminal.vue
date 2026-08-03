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

              <!-- What sshd will accept, asked before anything is typed. The
                   probe also republishes the host keys the handshake pins. -->
              <template v-if="selectedNode.terminal_enabled">
                <cv-skeleton-text
                  v-if="loading.probe"
                  :paragraph="true"
                  :line-count="2"
                ></cv-skeleton-text>
                <template v-else>
                  <NsInlineNotification
                    v-if="error.probe"
                    kind="error"
                    :title="$t('action.probe-terminal-access')"
                    :description="error.probe"
                    :showCloseButton="false"
                  />
                  <NsInlineNotification
                    v-else-if="!probe.listen_wg0"
                    kind="error"
                    :title="$t('terminal.sshd_not_listening')"
                    :description="$t('terminal.sshd_not_listening_description')"
                    :showCloseButton="false"
                  />
                  <!-- Most operators log in as root: say it before they try -->
                  <NsInlineNotification
                    v-else-if="!probe.permit_root_login"
                    kind="warning"
                    :title="$t('terminal.root_password_refused')"
                    :description="
                      $t('terminal.root_password_refused_description')
                    "
                    :showCloseButton="false"
                  />
                  <NsInlineNotification
                    v-else-if="!probe.password_auth"
                    kind="warning"
                    :title="$t('terminal.password_auth_disabled')"
                    :description="
                      $t('terminal.password_auth_disabled_description')
                    "
                    :showCloseButton="false"
                  />
                </template>
              </template>
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
                @click="confirmClose"
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

            <!-- Mounted as soon as the socket is up, not once connected: the
                 login prompts are drawn inside this pane. -->
            <XtermPane v-if="socket" ref="pane" :socket="socket" />
          </template>
        </cv-column>
      </cv-row>
    </cv-grid>

    <NsModal
      size="default"
      kind="danger"
      :visible="isCloseModalShown"
      @modal-hidden="isCloseModalShown = false"
      @primary-click="onCloseConfirmed"
    >
      <template slot="title">{{ $t("terminal.close_confirm") }}</template>
      <template slot="content">
        <p>{{ $t("terminal.close_confirm_description") }}</p>
      </template>
      <template slot="secondary-button">{{ $t("common.cancel") }}</template>
      <template slot="primary-button">{{ $t("terminal.close") }}</template>
    </NsModal>
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
import XtermPane from "@/components/terminal/XtermPane";

export default {
  name: "NodeTerminal",
  components: { XtermPane },
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
      pendingOutput: [],
      isConnected: false,
      isCloseModalShown: false,
      closedReason: "",
      probe: {
        permit_root_login: false,
        password_auth: false,
        listen_wg0: false,
        port: 22,
      },
      loading: {
        listNodes: false,
        toggle: false,
        probe: false,
      },
      error: {
        listNodes: "",
        toggle: "",
        session: "",
        probe: "",
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
        !this.loading.toggle &&
        !this.loading.probe
      );
    },
  },
  watch: {
    // Everything shown below the selector belongs to one node. Without this the
    // switch keeps the previous node's state, since syncToggle only ran after
    // the node list was reloaded.
    selectedNodeId: function () {
      this.syncToggle();
      this.error.toggle = "";
      this.error.probe = "";

      // The probe is not only informative: it republishes the host keys the
      // handshake pins, so it has to run before a session is requested.
      if (this.selectedNode && this.selectedNode.terminal_enabled) {
        this.probeAccess();
      }
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
    /**
     * Ask the node what sshd will accept. The probe passes a connection spec to
     * sshd -T, so it reports the policy that applies to the terminal rather
     * than the global one.
     */
    async probeAccess() {
      this.loading.probe = true;
      this.error.probe = "";
      const taskAction = "probe-terminal-access";
      const eventId = this.getUuid();

      this.$root.$once(`${taskAction}-aborted-${eventId}`, this.probeAborted);
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.probeCompleted
      );

      const res = await to(
        this.createNodeTask(this.selectedNodeId, {
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
        this.error.probe = this.getErrorMessage(err);
        this.loading.probe = false;
      }
    },
    probeCompleted(taskContext, taskResult) {
      this.probe = taskResult.output;
      this.loading.probe = false;
    },
    probeAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.probe = this.$t("terminal.probe_failed");
      this.loading.probe = false;
    },
    /**
     * Two steps on purpose. The REST call only authorizes and returns a
     * one-shot ticket; the login prompts then run inside the terminal, driven
     * by api-server for the user name and by sshd for the password.
     */
    async requestSession() {
      this.error.session = "";
      this.closedReason = "";

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

      this.openSocket(response.data.data.ticket);
    },
    openSocket(ticket) {
      const socket = new WebSocket(`${this.$root.config.WS_ENDPOINT}/terminal`);
      socket.binaryType = "arraybuffer";
      this.socket = socket;

      socket.addEventListener("open", () => {
        socket.send(JSON.stringify({ type: "ticket", ticket }));
      });

      socket.addEventListener("message", (event) => {
        if (typeof event.data !== "string") {
          // Terminal bytes: XtermPane consumes them once mounted. Until then
          // they are kept here, otherwise the login prompt would be lost.
          if (!this.$refs.pane) {
            this.pendingOutput.push(event.data);
          }
          return;
        }
        this.onControlFrame(JSON.parse(event.data));
      });

      this.$nextTick(this.flushPendingOutput);

      socket.addEventListener("close", () => {
        this.isConnected = false;
        this.socket = null;
      });

      socket.addEventListener("error", () => {
        this.error.session = this.$t("terminal.transport_error");
      });
    },
    onControlFrame(frame) {
      switch (frame.type) {
        case "ticket-accepted": {
          // Nothing to answer: the server owns the login dialogue from here.
          // Only the size is worth sending, so the shell it opens next matches
          // the pane the prompts are already drawn in.
          const size = this.$refs.pane
            ? this.$refs.pane.currentSize()
            : { rows: 24, cols: 80 };
          this.socket.send(
            JSON.stringify({ type: "resize", rows: size.rows, cols: size.cols })
          );
          break;
        }
        case "ready":
          this.isConnected = true;
          break;
        case "auth-error":
          this.error.session = frame.message;
          break;
        case "closed":
          this.closedReason = frame.reason;
          if (this.$refs.pane) {
            this.$refs.pane.writeNotice(frame.reason);
          }
          break;
      }
    },
    flushPendingOutput() {
      if (!this.$refs.pane) {
        return;
      }
      for (const data of this.pendingOutput) {
        this.$refs.pane.writeBytes(data);
      }
      this.pendingOutput = [];
    },
    // Only worth asking once a shell is up: at the login prompt there is
    // nothing running to lose.
    confirmClose() {
      if (this.isConnected) {
        this.isCloseModalShown = true;
      } else {
        this.closeSession();
      }
    },
    onCloseConfirmed() {
      this.isCloseModalShown = false;
      this.closeSession();
    },
    closeSession() {
      this.isCloseModalShown = false;
      if (this.socket) {
        this.socket.close();
        this.socket = null;
      }
      this.pendingOutput = [];
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

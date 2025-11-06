<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <cv-grid fullWidth>
    <cv-row>
      <cv-column class="page-title">
        <h2>{{ $t("nodes.title") }}</h2>
      </cv-column>
    </cv-row>
    <cv-row v-if="error.getClusterStatus">
      <cv-column>
        <NsInlineNotification
          kind="error"
          :title="$t('action.get-cluster-status')"
          :description="error.getClusterStatus"
          :showCloseButton="false"
        />
      </cv-column>
    </cv-row>
    <cv-row v-if="nodesOffline.length">
      <cv-column>
        <NsInlineNotification
          kind="error"
          :title="
            $tc('nodes.nodes_offline_c', nodesOffline.length, {
              num: nodesOffline.length,
            })
          "
          :description="
            $tc('nodes.nodes_offline_description_c', nodesOffline.length)
          "
          :showCloseButton="false"
        />
      </cv-column>
    </cv-row>
    <cv-row v-if="error.listNodes">
      <cv-column>
        <NsInlineNotification
          kind="error"
          :title="$t('action.list-nodes')"
          :description="error.listNodes"
          :showCloseButton="false"
        />
      </cv-column>
    </cv-row>
    <cv-row v-if="error.listAlerts">
      <cv-column>
        <NsInlineNotification
          kind="error"
          :title="$t('action.list-alerts')"
          :description="error.listAlerts"
          :showCloseButton="false"
        />
      </cv-column>
    </cv-row>
    <cv-row class="toolbar">
      <cv-column>
        <NsButton
          kind="secondary"
          :icon="Add20"
          @click="q.isShownAddNodeModal = true"
          >{{ $t("nodes.add_node_to_cluster") }}</NsButton
        >
      </cv-column>
    </cv-row>
    <cv-row v-if="stillLoading">
      <cv-column>
        <!-- skeleton card grid -->
        <div
          class="card-grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-3 3xl:grid-cols-4"
        >
          <cv-tile v-for="index in 2" :key="index" light>
            <cv-skeleton-text
              :paragraph="true"
              :line-count="9"
              heading
            ></cv-skeleton-text>
          </cv-tile>
        </div>
      </cv-column>
    </cv-row>
    <cv-row v-else>
      <cv-column>
        <!-- card grid -->
        <div
          class="card-grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-3 3xl:grid-cols-4"
        >
          <template v-for="node in clusterNodes">
            <NodeCard
              v-if="!node.online"
              :key="node.id"
              :nodeId="node.id"
              :nodeLabel="node.ui_name"
              :isLeader="leaderNode && node.id === leaderNode.id"
              light
              loading
              :online="node.online"
            >
              <template #menu>
                <cv-overflow-menu
                  v-if="node.id !== leaderNode.id"
                  :flip-menu="true"
                  tip-position="top"
                  tip-alignment="end"
                  class="top-right-overflow-menu"
                >
                  <cv-overflow-menu-item
                    danger
                    @click="showRemoveNodeModal(node)"
                  >
                    <NsMenuItem
                      :icon="TrashCan20"
                      :label="$t('nodes.remove_from_cluster')"
                    />
                  </cv-overflow-menu-item>
                </cv-overflow-menu>
              </template>
            </NodeCard>
            <NodeCard
              v-else
              :key="node.id + '_'"
              :nodeId="node.id"
              :nodeLabel="node.ui_name"
              :isLeader="leaderNode && node.id === leaderNode.id"
              :role="node.role"
              :leaderLabel="$t('nodes.leader')"
              :workerLabel="$t('nodes.worker')"
              :ns7MigrationLabel="$t('nodes.ns7_migration')"
              :diskUsageLabel="$t('nodes.usage')"
              :online="node.online"
              :fqdn="node.hostname"
              :fqdnLabel="$t('nodes.fqdn')"
              :node_id="node.id"
              :ip_address="node.ip_address"
              :ip_addressLabel="$t('nodes.ip_address')"
              :applications="node.applications"
              :applicationsLabel="$t('nodes.applications')"
              :alertsCount="alertsCountByNode[node.id] || 0"
              :alertsLabel="
                $tc('nodes.alerts', alertsCountByNode[node.id] || 0)
              "
              :class="{ 'node-card--ns7': node.role === 'ns7migration' }"
              light
            >
              <template #menu>
                <cv-overflow-menu
                  v-if="node.role !== 'ns7migration'"
                  :flip-menu="true"
                  tip-position="top"
                  tip-alignment="end"
                  class="top-right-overflow-menu"
                >
                  <cv-overflow-menu-item @click="showSetNodeLabelModal(node)">
                    <NsMenuItem
                      :icon="Edit20"
                      :label="$t('nodes.edit_node_label')"
                    />
                  </cv-overflow-menu-item>
                  <cv-overflow-menu-item @click="showSetNodeFqdnModal(node)">
                    <NsMenuItem :icon="Edit20" :label="$t('nodes.set_fqdn')" />
                  </cv-overflow-menu-item>
                  <cv-overflow-menu-item @click="goToApplications(node)">
                    <NsMenuItem
                      :icon="ArrowRight20"
                      :label="$t('nodes.go_to_applications')"
                    />
                  </cv-overflow-menu-item>
                  <cv-overflow-menu-item @click="goToHttpRoutes(node)">
                    <NsMenuItem
                      :icon="ArrowRight20"
                      :label="$t('nodes.go_to_http_routes')"
                    />
                  </cv-overflow-menu-item>
                  <cv-overflow-menu-item @click="goToTlsCertificates(node)">
                    <NsMenuItem
                      :icon="ArrowRight20"
                      :label="$t('nodes.go_to_tls_certificates')"
                    />
                  </cv-overflow-menu-item>
                  <cv-overflow-menu-item @click="goToFirewall(node)">
                    <NsMenuItem
                      :icon="ArrowRight20"
                      :label="$t('nodes.go_to_firewall')"
                    />
                  </cv-overflow-menu-item>
                  <cv-overflow-menu-item
                    v-if="leaderNode && node.id !== leaderNode.id"
                    @click="showPromoteNodeModal(node)"
                  >
                    <NsMenuItem
                      :icon="Recommend20"
                      :label="$t('nodes.promote_to_leader')"
                      ><Recommend20 />
                    </NsMenuItem>
                  </cv-overflow-menu-item>
                  <cv-overflow-menu-item
                    v-if="leaderNode && node.id !== leaderNode.id"
                    danger
                    @click="showRemoveNodeModal(node)"
                  >
                    <NsMenuItem
                      :icon="TrashCan20"
                      :label="$t('nodes.remove_from_cluster')"
                    />
                  </cv-overflow-menu-item>
                </cv-overflow-menu>
              </template>
              <template v-if="node.role !== 'ns7migration'" #content>
                <NsButton
                  kind="ghost"
                  :icon="ArrowRight20"
                  @click="goToNodeDetail(node.id)"
                  >{{ $t("common.see_details") }}</NsButton
                >
              </template>
            </NodeCard>
          </template>
        </div>
      </cv-column>
    </cv-row>
    <!-- add node modal -->
    <NsModal
      size="default"
      :visible="q.isShownAddNodeModal"
      @modal-hidden="q.isShownAddNodeModal = false"
    >
      <template slot="title">{{ $t("nodes.add_node_to_cluster") }}</template>
      <template slot="content">
        <ol class="mg-bottom-md">
          <li>
            {{
              $t("nodes.add_node_to_cluster_step_1", {
                product: this.$root.config.PRODUCT_NAME,
              })
            }}
          </li>
          <li>{{ $t("nodes.add_node_to_cluster_step_2") }}</li>
          <li>{{ $t("nodes.add_node_to_cluster_step_3") }}</li>
          <li>{{ $t("nodes.add_node_to_cluster_step_4") }}</li>
          <li>{{ $t("nodes.add_node_to_cluster_step_5") }}</li>
          <li>{{ $t("nodes.add_node_to_cluster_step_6") }}</li>
        </ol>
        <cv-skeleton-text
          v-if="loading.getFqdn || !leaderNode"
          :paragraph="true"
          :line-count="5"
          heading
        ></cv-skeleton-text>
        <template v-else>
          <span class="join-code">{{ $t("common.join_code") }}</span>
          <NsCodeSnippet
            :copyTooltip="$t('common.copy_to_clipboard')"
            :copy-feedback="$t('common.copied_to_clipboard')"
            :feedback-aria-label="$t('common.copied_to_clipboard')"
            :wrap-text="true"
            :moreText="$t('common.show_more')"
            :lessText="$t('common.show_less')"
            light
            expanded
            hideExpandButton
            >{{ joinCode }}</NsCodeSnippet
          >
        </template>
      </template>
      <template slot="secondary-button">{{ $t("common.close") }}</template>
    </NsModal>
    <!-- set node label modal -->
    <NsModal
      size="default"
      :visible="isShownSetNodeLabelModal"
      @modal-hidden="hideSetNodeLabelModal"
      @primary-click="setNodeLabel"
    >
      <template slot="title">{{ $t("nodes.edit_node_label") }}</template>
      <template slot="content">
        <template v-if="currentNode">
          <cv-form @submit.prevent="setNodeLabel">
            <cv-text-input
              :label="
                $t('nodes.node_label') + ' (' + $t('common.optional') + ')'
              "
              v-model.trim="newNodeLabel"
              :placeholder="$t('common.no_label')"
              :helper-text="$t('nodes.node_label_tooltip')"
              maxlength="24"
              ref="newNodeLabel"
            >
            </cv-text-input>
            <div v-if="error.setNodeLabel">
              <NsInlineNotification
                kind="error"
                :title="$t('action.set-name')"
                :description="error.setNodeLabel"
                :showCloseButton="false"
              />
            </div>
          </cv-form>
        </template>
      </template>
      <template slot="secondary-button">{{ $t("common.cancel") }}</template>
      <template slot="primary-button">{{ $t("common.save") }}</template>
    </NsModal>
    <!-- remove node modal -->
    <RemoveNodeModal
      :isShown="isShownRemoveNodeModal"
      :node="nodeToRemove"
      @hide="hideRemoveNodeModal"
    />
    <!-- promote node modal -->
    <PromoteNodeModal
      :isShown="isShownPromoteNodeModal"
      :node="nodeToPromote"
      @hide="hidePromoteNodeModal"
      @nodePromoted="onNodePromoted"
    />
    <!-- new leader modal -->
    <NewLeaderModal
      :isShown="isShownNewLeaderModal"
      :endpointAddress="newLeaderEndpointAddress"
      @hide="hideNewLeaderModal"
    />
    <!-- set node FQDN modal -->
    <SetFqdnModal
      :isShown="isShownSetFqdnModal"
      :node="currentNode"
      @hide="hideSetFqdnModal"
    />
  </cv-grid>
</template>

<script>
import {
  QueryParamService,
  UtilService,
  TaskService,
  IconService,
  PageTitleService,
} from "@nethserver/ns8-ui-lib";
import to from "await-to-js";
import { mapState, mapActions } from "vuex";
import NodeCard from "@/components/nodes/NodeCard";
import RemoveNodeModal from "@/components/nodes/RemoveNodeModal";
import PromoteNodeModal from "@/components/nodes/PromoteNodeModal";
import NewLeaderModal from "@/components/nodes/NewLeaderModal";
import SetFqdnModal from "@/components/nodes/SetFqdnModal";
import Recommend20 from "@carbon/icons-vue/es/recommend/20";

export default {
  name: "Nodes",
  components: {
    NodeCard,
    RemoveNodeModal,
    PromoteNodeModal,
    NewLeaderModal,
    SetFqdnModal,
    Recommend20,
  },
  mixins: [
    TaskService,
    UtilService,
    IconService,
    QueryParamService,
    PageTitleService,
  ],
  pageTitle() {
    return this.$t("nodes.title");
  },
  data() {
    return {
      REFRESH_DATA_TIME_INTERVAL: 30000,
      q: {
        isShownAddNodeModal: false,
      },
      Recommend20,
      joinCode: "",
      nodes: [],
      clusterNodes: [],
      refreshDataInterval: null,
      currentNode: null,
      newNodeLabel: "",
      isShownSetNodeLabelModal: false,
      isShownRemoveNodeModal: false,
      nodeToRemove: null,
      isShownPromoteNodeModal: false,
      nodeToPromote: null,
      isShownNewLeaderModal: false,
      newLeaderEndpointAddress: "",
      isShownSetFqdnModal: false,
      alertsCountByNode: {},
      loading: {
        getClusterStatus: true,
        setNodeLabel: false,
        getFqdn: false,
        listAlerts: true,
        listNodes: true,
      },
      error: {
        getClusterStatus: "",
        setNodeLabel: "",
        getFqdn: "",
        listNodes: "",
        listAlerts: "",
      },
    };
  },
  computed: {
    ...mapState(["leaderListenPort"]),
    leaderNode: function () {
      if (!this.nodes) {
        return null;
      }

      return this.nodes.find((node) => node.local);
    },
    nodesOffline() {
      return this.nodes.filter((n) => n.online == false);
    },
    stillLoading() {
      return (
        this.loading.getClusterStatus ||
        this.loading.listAlerts ||
        this.loading.listNodes
      );
    },
  },
  watch: {
    "q.isShownAddNodeModal": function () {
      if (this.q.isShownAddNodeModal && this.leaderNode) {
        this.getLeaderFqdn();
      }
    },
    leaderNode: function (newLeaderNode, oldLeaderNode) {
      if (
        newLeaderNode &&
        this.q.isShownAddNodeModal &&
        newLeaderNode.id != oldLeaderNode?.id
      ) {
        this.getLeaderFqdn();
      }
    },
  },
  beforeRouteEnter(to, from, next) {
    next((vm) => {
      vm.watchQueryData(vm);
      vm.queryParamsToDataForCore(vm, to.query);
    });
  },
  beforeRouteUpdate(to, from, next) {
    this.queryParamsToDataForCore(this, to.query);
    next();
  },
  created() {
    this.retrieveData();

    // periodically retrieve cluster and nodes status
    this.refreshDataInterval = setInterval(
      this.retrieveData,
      this.REFRESH_DATA_TIME_INTERVAL
    );
  },
  beforeDestroy() {
    clearInterval(this.refreshDataInterval);
  },
  methods: {
    ...mapActions(["setClusterNodesInStore"]),
    retrieveData() {
      this.listAlerts();
      // cluster status must be retrieved before list-nodes to have authoritative online/local flags
      this.getClusterStatus();
    },
    async getLeaderFqdn() {
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
        this.createNodeTask(this.leaderNode.id, {
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
        this.loading.getFqdn = false;
        return;
      }
    },
    getFqdnCompleted(taskContext, taskResult) {
      const fqdn = `${taskResult.output.hostname}.${taskResult.output.domain}`;
      const loginInfo = this.getFromStorage("loginInfo");

      if (loginInfo && loginInfo.token) {
        const endpoint = window.location.protocol + "//" + fqdn;

        // join code is obtained by concatenating FQDN and auth token with pipe character
        this.joinCode = btoa(endpoint + "|" + loginInfo.token);
      }
      this.loading.getFqdn = false;
    },
    getFqdnAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.loading.getFqdn = false;
    },
    async getClusterStatus() {
      this.error.getClusterStatus = "";
      const taskAction = "get-cluster-status";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.getClusterStatusAborted
      );
      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.getClusterStatusCompleted
      );

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
        this.error.getClusterStatus = this.getErrorMessage(err);
        return;
      }
    },
    getClusterStatusAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.getClusterStatus = this.$t("error.generic_error");
      this.loading.getClusterStatus = false;
    },
    getClusterStatusCompleted(taskContext, taskResult) {
      const clusterStatus = taskResult.output;

      this.nodes = clusterStatus.nodes.sort(this.sortByProperty("id"));
      this.setClusterNodesInStore(this.nodes);

      this.loading.getClusterStatus = false;
      // now that authoritative online/local flags are present, fetch list-nodes
      this.listNodes();
    },
    async listNodes() {
      this.error.listNodes = "";
      const taskAction = "list-nodes";
      const eventId = this.getUuid();
      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.listNodesAborted
      );
      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.listNodesCompleted
      );

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
        return;
      }
    },
    listNodesAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.listNodes = this.$t("error.generic_error");
      this.loading.listNodes = false;
    },
    listNodesCompleted(taskContext, taskResult) {
      const nodesData = taskResult.output.nodes;

      // Lookups from cluster status (may be empty if somehow race still happens)
      // Guard against missing nodes
      const clusterOnlineById = Object.fromEntries(
        (this.nodes || []).map((node) => [node.id, !!node.online])
      );

      const transformedNodes = nodesData.map((node) => {
        const baseOnline =
          node.node_id in clusterOnlineById
            ? clusterOnlineById[node.node_id]
            : !!(node.fqdn && node.fqdn !== ""); // fallback if cluster status not available, fqdn presence means online

        // ns7migration nodes must be shown as online even without fqdn
        const online = node.role === "ns7migration" ? true : baseOnline;

        return {
          id: node.node_id,
          hostname: node.fqdn || "",
          ui_name: node.ui_name || "",
          role: node.role || "worker",
          online: online,
          ip_address: node.main_ip || "",
          applications: node.app_count || 0,
        };
      });

      this.clusterNodes = transformedNodes.sort(this.sortByProperty("id"));
      this.loading.listNodes = false;
    },
    async listAlerts() {
      this.error.listAlerts = "";
      const taskAction = "list-alerts";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.listAlertsAborted
      );
      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.listAlertsCompleted
      );

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
        this.error.listAlerts = this.getErrorMessage(err);
        return;
      }
    },
    listAlertsAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.listAlerts = this.$t("error.generic_error");
      this.loading.listAlerts = false;
    },
    listAlertsCompleted(taskContext, taskResult) {
      const alerts = taskResult.output.alerts;
      const alertsCountByNode = {};
      alerts.forEach((alert) => {
        if (alert.node_id in alertsCountByNode) {
          alertsCountByNode[alert.node_id] += 1;
        } else {
          alertsCountByNode[alert.node_id] = 1;
        }
      });
      this.alertsCountByNode = alertsCountByNode;
      this.loading.listAlerts = false;
    },
    goToNodeDetail(nodeId) {
      this.$router.push({
        name: "NodeDetail",
        params: { nodeId },
      });
    },
    goToFirewall(node) {
      this.$router.push({
        name: "NodeFirewall",
        params: { nodeId: node.id },
      });
    },
    showSetNodeLabelModal(node) {
      this.currentNode = node;
      this.newNodeLabel = node.ui_name;
      this.isShownSetNodeLabelModal = true;
      setTimeout(() => {
        this.focusElement("newNodeLabel");
      }, 300);
    },
    showSetNodeFqdnModal(node) {
      this.currentNode = node;
      this.isShownSetFqdnModal = true;
    },
    hideSetNodeLabelModal() {
      this.isShownSetNodeLabelModal = false;
    },
    async setNodeLabel() {
      this.error.setNodeLabel = "";
      this.loading.setNodeLabel = true;
      const taskAction = "set-name";

      // register to task completion
      this.$root.$once(taskAction + "-completed", this.setNodeLabelCompleted);

      const res = await to(
        this.createNodeTask(this.currentNode.id, {
          action: taskAction,
          data: {
            name: this.newNodeLabel,
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
        this.error.setNodeLabel = this.getErrorMessage(err);
        this.loading.setNodeLabel = false;
        return;
      }
    },
    setNodeLabelCompleted() {
      this.loading.setNodeLabel = false;
      this.hideSetNodeLabelModal();
      this.getClusterStatus();
    },
    goToHttpRoutes(node) {
      this.$router.push({
        path: "/settings/http-routes",
        query: { selectedNodeId: node.id },
      });
    },
    goToTlsCertificates(node) {
      this.$router.push({
        path: "/settings/tls-certificates",
        query: { selectedNodeId: node.id },
      });
    },
    goToApplications(node) {
      this.$router.push({
        path: "/applications",
        query: { selectedNodeId: node.id },
      });
    },
    showPromoteNodeModal(node) {
      this.nodeToPromote = node;
      this.isShownPromoteNodeModal = true;
    },
    hidePromoteNodeModal() {
      this.isShownPromoteNodeModal = false;
    },
    showNewLeaderModal() {
      this.isShownNewLeaderModal = true;
    },
    hideNewLeaderModal() {
      this.isShownNewLeaderModal = false;
    },
    hideSetFqdnModal() {
      this.isShownSetFqdnModal = false;
    },
    onNodePromoted(newLeaderEndpointAddress) {
      this.newLeaderEndpointAddress = newLeaderEndpointAddress;
      this.showNewLeaderModal();
    },
    showRemoveNodeModal(node) {
      this.nodeToRemove = node;
      this.isShownRemoveNodeModal = true;
    },
    hideRemoveNodeModal() {
      this.isShownRemoveNodeModal = false;
    },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";

ol {
  list-style-type: decimal;
  padding-left: $spacing-05;
}

.join-code {
  font-weight: bold;
}
</style>

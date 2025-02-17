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
    <cv-row v-if="error.getNodeStatus">
      <cv-column>
        <NsInlineNotification
          kind="error"
          :title="$t('action.get-node-status')"
          :description="error.getNodeStatus"
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
    <cv-row v-if="loading.nodes">
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
          <template v-for="node in nodes">
            <NodeCard
              v-if="!nodesStatus[node.id]"
              :key="node.id"
              :nodeId="node.id"
              :nodeLabel="node.ui_name"
              :isLeader="node.id == leaderNode.id"
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
              :isLeader="node.id == leaderNode.id"
              :leaderLabel="$t('nodes.leader')"
              :workerLabel="$t('nodes.worker')"
              :cpuUsageLabel="$t('nodes.cpu_usage')"
              :cpuLoadLabel="$t('nodes.cpu_load')"
              :cpuLoadTooltip="$t('nodes.cpu_load_tooltip')"
              :memoryUsageLabel="$t('nodes.memory_usage')"
              :swapUsageLabel="$t('nodes.swap_usage')"
              :diskUsageLabel="$t('nodes.usage')"
              :cpuUsage="nodesStatus[node.id].cpu.usage"
              :cpuUsageWarningTh="90"
              :load1Min="nodesStatus[node.id].load['1min']"
              :load5Min="nodesStatus[node.id].load['5min']"
              :load15Min="nodesStatus[node.id].load['15min']"
              :cpuLoadWarningTh="90"
              :memoryUsage="nodesStatus[node.id].memoryUsage"
              :memoryWarningTh="90"
              :swapUsage="nodesStatus[node.id].swapUsage"
              :swapWarningTh="90"
              :disksUsage="nodesStatus[node.id].disksUsage"
              :diskWarningTh="90"
              :online="node.online"
              light
            >
              <template #menu>
                <cv-overflow-menu
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
                    <NsMenuItem :icon="Wikis20" :label="$t('nodes.set_fqdn')" />
                  </cv-overflow-menu-item>
                  <cv-overflow-menu-item @click="goToHttpRoutes(node)">
                    <NsMenuItem
                      :icon="Router20"
                      :label="$t('settings_http_routes.title')"
                    />
                  </cv-overflow-menu-item>
                  <cv-overflow-menu-item @click="goToTlsCertificates(node)">
                    <NsMenuItem
                      :icon="Certificate20"
                      :label="$t('settings_tls_certificates.title')"
                    />
                  </cv-overflow-menu-item>
                  <cv-overflow-menu-item @click="goToFirewall(node)">
                    <NsMenuItem
                      :icon="Firewall20"
                      :label="$t('firewall.title')"
                    />
                  </cv-overflow-menu-item>
                  <cv-overflow-menu-item
                    v-if="node.id !== leaderNode.id"
                    @click="showPromoteNodeModal(node)"
                  >
                    <NsMenuItem
                      :icon="Recommend20"
                      :label="$t('nodes.promote_to_leader')"
                      ><Recommend20 />
                    </NsMenuItem>
                  </cv-overflow-menu-item>
                  <cv-overflow-menu-item
                    v-if="node.id !== leaderNode.id"
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
              <template #content>
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
      <template slot="primary-button">{{
        $t("nodes.edit_node_label")
      }}</template>
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
      nodesStatus: {},
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
      loading: {
        nodes: true,
        setNodeLabel: false,
        getFqdn: false,
      },
      error: {
        getClusterStatus: "",
        getNodeStatus: "",
        setNodeLabel: "",
        getFqdn: "",
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
    this.getClusterStatus();

    // periodically retrieve cluster and nodes status
    this.refreshDataInterval = setInterval(
      this.getClusterStatus,
      this.REFRESH_DATA_TIME_INTERVAL
    );
  },
  beforeDestroy() {
    clearInterval(this.refreshDataInterval);
  },
  methods: {
    ...mapActions(["setClusterNodesInStore"]),
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

      // register to task completion
      this.$root.$once(
        taskAction + "-completed",
        this.getClusterStatusCompleted
      );

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
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
    getClusterStatusCompleted(taskContext, taskResult) {
      const clusterStatus = taskResult.output;

      // update nodes in vuex store
      this.nodes = clusterStatus.nodes.sort(this.sortByProperty("id"));
      this.setClusterNodesInStore(this.nodes);

      this.loading.nodes = false;
      this.retrieveNodesStatus();
    },
    async retrieveNodesStatus() {
      this.error.getNodeStatus = "";

      for (const node of this.nodes) {
        if (!node.online) {
          // node is offline
          continue;
        }

        const nodeId = node.id;
        const taskAction = "get-node-status";
        const eventId = this.getUuid();

        // register to task events
        this.$root.$once(
          `${taskAction}-completed-${eventId}`,
          this.getNodeStatusCompleted
        );

        const res = await to(
          this.createNodeTask(nodeId, {
            action: taskAction,
            extra: {
              title: this.$t("action." + taskAction),
              isNotificationHidden: true,
              node: nodeId,
              eventId,
            },
          })
        );
        const err = res[0];

        if (err) {
          console.error(`error creating task ${taskAction}`, err);
          this.error.getNodeStatus = this.getErrorMessage(err);
        }
      }
    },
    getNodeStatusCompleted(taskContext, taskResult) {
      const nodeId = taskContext.extra.node;
      const nodeStatus = taskResult.output;

      nodeStatus.cpu.usage = Math.round(nodeStatus.cpu.usage);
      nodeStatus.load["1min"] = Math.round(nodeStatus.load["1min"]);
      nodeStatus.load["5min"] = Math.round(nodeStatus.load["5min"]);
      nodeStatus.load["15min"] = Math.round(nodeStatus.load["15min"]);

      // memory and swap usage
      nodeStatus.memoryUsage = Math.round(
        (nodeStatus.memory.used / nodeStatus.memory.total) * 100
      );
      nodeStatus.swapUsage = Math.round(
        (nodeStatus.swap.used / nodeStatus.swap.total) * 100
      );

      // disks usage
      const disksUsage = [];
      let diskIndex = 1;

      for (const disk of nodeStatus.disks) {
        const diskUsage = Math.round((disk.used / disk.total) * 100);
        const diskId = this.$t("nodes.disk") + ` ${diskIndex}`;
        disksUsage.push({ diskId: diskId, usage: diskUsage });
        diskIndex++;
      }
      nodeStatus.disksUsage = disksUsage;

      // needed for reactivity (see https://vuejs.org/v2/guide/reactivity.html#For-Objects)
      this.$set(this.nodesStatus, nodeId, nodeStatus);
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

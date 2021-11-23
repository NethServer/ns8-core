<template>
  <div class="bx--grid bx--grid--full-width">
    <div class="bx--row">
      <div class="bx--col-lg-16 page-title">
        <h2>{{ $t("nodes.title") }}</h2>
      </div>
    </div>
    <div v-if="error.getClusterStatus" class="bx--row">
      <div class="bx--col">
        <NsInlineNotification
          kind="error"
          :title="$t('action.get-cluster-status')"
          :description="error.getClusterStatus"
          :showCloseButton="false"
        />
      </div>
    </div>
    <div v-if="error.getNodeStatus" class="bx--row">
      <div class="bx--col">
        <NsInlineNotification
          kind="error"
          :title="$t('action.get-node-status')"
          :description="error.getNodeStatus"
          :showCloseButton="false"
        />
      </div>
    </div>
    <div class="bx--row">
      <div class="bx--col">
        <NsButton
          kind="secondary"
          :icon="Add20"
          @click="q.isShownAddNodeModal = true"
          >{{ $t("nodes.add_node_to_cluster") }}</NsButton
        >
      </div>
    </div>
    <div
      class="bx--row loader-large loader-theme nodes-loader"
      v-if="loading.nodes"
    ></div>
    <div class="bx--row" v-else>
      <div
        v-for="node in nodes"
        :key="node.id"
        class="bx--col-md-4 bx--col-max-4"
      >
        <NodeCard
          v-if="!nodesStatus[node.id]"
          :nodeId="node.id"
          :nodeLabel="node.ui_name"
          :isLeader="node.id == leaderNode.id"
          light
          loading
        />
        <NodeCard
          v-else
          :nodeId="node.id"
          :nodeLabel="node.ui_name"
          :isLeader="node.id == leaderNode.id"
          :leaderLabel="$t('nodes.leader')"
          :workerLabel="$t('nodes.worker')"
          :cpuUsageLabel="$t('nodes.cpu_usage')"
          :cpuLoadLabel="$t('nodes.cpu_load')"
          :cpuPressureLabel="$t('nodes.cpu_pressure')"
          :cpuLoadTooltip="$t('nodes.cpu_load_tooltip')"
          :cpuPressureTooltip="$t('nodes.cpu_pressure_tooltip')"
          :memoryUsageLabel="$t('nodes.memory_usage')"
          :swapUsageLabel="$t('nodes.swap_usage')"
          :diskUsageLabel="$t('nodes.usage')"
          :cpuUsage="nodesStatus[node.id].cpu.usage"
          :cpuUsageWarningTh="90"
          :load1Min="nodesStatus[node.id].load['1min']"
          :load5Min="nodesStatus[node.id].load['5min']"
          :load15Min="nodesStatus[node.id].load['15min']"
          :pressure10Sec="nodesStatus[node.id].pressure['10sec']"
          :pressure1min="nodesStatus[node.id].pressure['1min']"
          :pressure5Min="nodesStatus[node.id].pressure['5min']"
          :cpuLoadWarningTh="90"
          :cpuPressureWarningTh="90"
          :memoryUsage="nodesStatus[node.id].memoryUsage"
          :memoryWarningTh="90"
          :swapUsage="nodesStatus[node.id].swapUsage"
          :swapWarningTh="90"
          :disksUsage="nodesStatus[node.id].disksUsage"
          :diskWarningTh="90"
          light
        >
          <NsButton
            kind="ghost"
            :icon="ZoomIn20"
            @click="goToNodeDetail(node.id)"
            >{{ $t("common.details") }}</NsButton
          >
          <cv-overflow-menu
            :flip-menu="true"
            tip-position="top"
            tip-alignment="end"
            class="top-right-overflow-menu"
          >
            <cv-overflow-menu-item @click="showSetNodeLabelModal(node)">{{
              $t("nodes.edit_node_label")
            }}</cv-overflow-menu-item>
          </cv-overflow-menu>
        </NodeCard>
      </div>
    </div>
    <!-- add node modal -->
    <cv-modal
      size="default"
      :visible="q.isShownAddNodeModal"
      @modal-hidden="q.isShownAddNodeModal = false"
    >
      <template slot="title">{{ $t("nodes.add_node_to_cluster") }}</template>
      <template slot="content">
        <ol class="mg-bottom-md">
          <li>{{ $t("nodes.add_node_to_cluster_step_1") }}</li>
          <li v-html="$t('nodes.add_node_to_cluster_step_2')"></li>
          <li v-html="$t('nodes.add_node_to_cluster_step_3')"></li>
          <li>{{ $t("nodes.add_node_to_cluster_step_4") }}</li>
        </ol>
        <span class="join-code">{{ $t("common.join_code") }}</span>
        <!-- copy to clipboard hint -->
        <span class="hint hint-copy-to-clipboard">
          <cv-interactive-tooltip
            alignment="end"
            direction="bottom"
            :visible="isCopyClipboardHintShown"
          >
            <template slot="trigger">
              <span></span>
            </template>
            <template slot="content">
              <p>
                {{ $t("hint.copy_to_clipboard") }}
              </p>
              <NsButton
                kind="primary"
                size="small"
                @click="isCopyClipboardHintShown = false"
                class="hint-button"
                >{{ $t("common.got_it") }}</NsButton
              >
            </template>
          </cv-interactive-tooltip>
        </span>
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
          >{{ joinCode }}
        </NsCodeSnippet>
      </template>
      <template slot="secondary-button">{{ $t("common.close") }}</template>
    </cv-modal>
    <!-- set node label modal -->
    <cv-modal
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
            <div v-if="error.setNodeLabel" class="bx--row">
              <div class="bx--col">
                <NsInlineNotification
                  kind="error"
                  :title="$t('action.set-name')"
                  :description="error.setNodeLabel"
                  :showCloseButton="false"
                />
              </div>
            </div>
          </cv-form>
        </template>
      </template>
      <template slot="secondary-button">{{ $t("common.cancel") }}</template>
      <template slot="primary-button">{{
        $t("nodes.edit_node_label")
      }}</template>
    </cv-modal>
  </div>
</template>

<script>
import {
  QueryParamService,
  UtilService,
  TaskService,
  IconService,
} from "@nethserver/ns8-ui-lib";
import to from "await-to-js";
import NodeService from "@/mixins/node";
import { mapState } from "vuex";
import NodeCard from "@/components/NodeCard";

export default {
  name: "Nodes",
  components: { NodeCard },
  mixins: [
    TaskService,
    UtilService,
    IconService,
    QueryParamService,
    NodeService,
  ],
  pageTitle() {
    return this.$t("nodes.title");
  },
  data() {
    return {
      NODE_STATUS_TIME_INTERVAL: 5000,
      q: {
        isShownAddNodeModal: false,
      },
      joinCode: "",
      isCopyClipboardHintShown: false,
      nodes: [],
      nodesStatus: {},
      nodesStatusInterval: null,
      currentNode: null,
      newNodeLabel: "",
      isShownSetNodeLabelModal: false,
      loading: {
        nodes: true,
        setNodeLabel: false,
      },
      error: {
        getClusterStatus: "",
        getNodeStatus: "",
        setNodeLabel: "",
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
  },
  watch: {
    "q.isShownAddNodeModal": function () {
      if (this.q.isShownAddNodeModal) {
        this.retrieveJoinCode();
        this.showCopyClipboardHint();
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
    this.retrieveClusterStatus();
  },
  beforeDestroy() {
    clearInterval(this.nodesStatusInterval);
  },
  methods: {
    retrieveJoinCode() {
      const loginInfo = this.getFromStorage("loginInfo");

      if (loginInfo && loginInfo.token) {
        const endpoint =
          window.location.protocol + "//" + window.location.hostname;

        console.log("endpoint", endpoint); ////
        console.log("leaderListenPort", this.leaderListenPort); ////

        // join code is obtained by concatenating endpoint, leader VPN port and auth token with pipe character
        this.joinCode = btoa(
          endpoint + "|" + this.leaderListenPort + "|" + loginInfo.token
        );
      }
    },
    async retrieveClusterStatus() {
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
      this.nodes = clusterStatus.nodes.sort(this.sortByProperty("id"));
      this.loading.nodes = false;

      // immediately retrieve nodes status
      this.retrieveNodesStatus();

      // periodically retrieve nodes status
      this.nodesStatusInterval = setInterval(
        this.retrieveNodesStatus,
        this.NODE_STATUS_TIME_INTERVAL
      );
    },
    showCopyClipboardHint() {
      setTimeout(() => {
        const isCopyClipboardHintShown = this.getFromStorage(
          "isCopyClipboardHintShown"
        );

        if (!isCopyClipboardHintShown) {
          this.isCopyClipboardHintShown = true;
          this.saveToStorage("isCopyClipboardHintShown", true);
        }
      }, 1000);
    },
    async retrieveNodesStatus() {
      this.error.getNodeStatus = "";

      for (const node of this.nodes) {
        const nodeId = node.id;
        const taskAction = "get-node-status";

        // register to task events
        this.$root.$once(
          taskAction + "-completed-node-" + nodeId,
          this.getNodeStatusCompleted
        );

        const res = await to(
          this.createNodeTask(nodeId, {
            action: taskAction,
            extra: {
              title: this.$t("action." + taskAction),
              isNotificationHidden: true,
              node: nodeId,
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
      nodeStatus.pressure["10sec"] = Math.round(nodeStatus.pressure["10sec"]);
      nodeStatus.pressure["1min"] = Math.round(nodeStatus.pressure["1min"]);
      nodeStatus.pressure["5min"] = Math.round(nodeStatus.pressure["5min"]);

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
    showSetNodeLabelModal(node) {
      this.currentNode = node;
      this.newNodeLabel = node.ui_name;
      this.isShownSetNodeLabelModal = true;
      setTimeout(() => {
        this.focusElement("newNodeLabel");
      }, 300);
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
      this.retrieveClusterStatus();
    },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";

.nodes-loader {
  margin: $spacing-05 auto;
}

ol {
  list-style-type: decimal;
  padding-left: $spacing-05;
}

.join-code {
  font-weight: bold;
}

.hint-copy-to-clipboard {
  top: 2.5rem;
  right: 2.5rem;
  float: right;
}
</style>

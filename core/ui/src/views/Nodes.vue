<template>
  <div class="bx--grid bx--grid--full-width">
    <div class="bx--row">
      <div class="bx--col-lg-16 page-title">
        <h2>{{ $t("nodes.title") }}</h2>
      </div>
    </div>
    <div v-if="error.nodes" class="bx--row">
      <div class="bx--col">
        <NsInlineNotification
          kind="error"
          :title="$t('error.cannot_retrieve_cluster_nodes')"
          :description="error.nodes"
          :showCloseButton="false"
        />
      </div>
    </div>
    <div class="bx--row">
      <div class="bx--col">
        <NsButton
          kind="secondary"
          :icon="Add20"
          @click="showAddNodeModal"
          :disabled="!!error.nodes"
          >{{ $t("nodes.add_node_to_cluster") }}</NsButton
        >
      </div>
    </div>
    <div class="bx--row loader-large nodes-loader" v-if="!nodes.length"></div>
    <div class="bx--row" v-else>
      <div
        v-for="node in nodes"
        :key="node.id"
        class="bx--col-md-4 bx--col-max-4"
      >
        <div v-if="!nodesStatus[node.id]">
          <cv-tile light>
            <cv-skeleton-text
              :paragraph="true"
              :line-count="12"
            ></cv-skeleton-text>
          </cv-tile>
        </div>
        <NsNodeCard
          v-else
          :nodeId="node.id.toString()"
          :nodeLabel="$t('common.node')"
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
          :cpuUsageWarningTh="80"
          :load1Min="nodesStatus[node.id].load['1min']"
          :load5Min="nodesStatus[node.id].load['5min']"
          :load15Min="nodesStatus[node.id].load['15min']"
          :cpuLoadWarningTh="90"
          :memoryUsage="nodesStatus[node.id].memoryUsage"
          :memoryWarningTh="80"
          :swapUsage="nodesStatus[node.id].swapUsage"
          :swapWarningTh="80"
          :disksUsage="nodesStatus[node.id].disksUsage"
          :diskWarningTh="80"
          light
        >
          <NsButton
            kind="ghost"
            :icon="ZoomIn20"
            @click="goToNodeDetail(node.id)"
            >{{ $t("common.details") }}</NsButton
          >
        </NsNodeCard>
      </div>
    </div>
    <!-- add node modal -->
    <cv-modal
      size="default"
      :visible="isShownAddNodeModal"
      @modal-hidden="isShownAddNodeModal = false"
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

export default {
  name: "Nodes",
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
      isShownAddNodeModal: false,
      joinCode: "",
      isCopyClipboardHintShown: false,
      nodes: [],
      nodesStatus: {},
      nodesStatusInterval: null,
      useNodesStatusMock: true, //// remove
      loading: {
        nodes: true,
      },
      error: {
        nodes: "",
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
        let endpoint = //// use const
          window.location.protocol + "//" + window.location.hostname;

        //// remove
        //endpoint = "https://192.168.122.220";

        console.log("endpoint", endpoint); ////
        console.log("leaderListenPort", this.leaderListenPort); ////

        // join code is obtained by concatenating endpoint, leader VPN port and auth token with pipe character
        this.joinCode = btoa(
          endpoint + "|" + this.leaderListenPort + "|" + loginInfo.token
        );

        console.log("joinCode", this.joinCode); ////
      }
    },
    async retrieveClusterStatus() {
      const taskAction = "get-cluster-status";
      const completedCallback = this.getClusterStatusCompleted;

      // register to task completion
      this.$root.$once(taskAction + "-completed", completedCallback);

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
        this.createErrorNotification(
          err,
          this.$t("task.cannot_create_task", { action: taskAction })
        );
        return;
      }
    },
    getClusterStatusCompleted(taskContext, taskResult) {
      console.log("getClusterStatusCompleted"); ////
      console.log("taskResult", taskResult); ////

      const clusterStatus = taskResult.output;
      console.log("clusterStatus", clusterStatus); ////

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
    showAddNodeModal() {
      this.retrieveJoinCode();
      this.isShownAddNodeModal = true;
      this.showCopyClipboardHint();
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
      console.log("retrieveNodesStatus"); ////

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
          this.createErrorNotification(
            err,
            this.$t("task.cannot_create_task", { action: taskAction })
          );
        }

        // //// remove mock
        // if (this.useNodesStatusMock) {
        //   const nodeStatus = {
        //     leader: nodeId == "1",
        //     cpu: { usage: Math.ceil(Math.random() * 100) },
        //     load: {
        //       "1min": Math.ceil(Math.random() * 200) / 100,
        //       "5min": Math.ceil(Math.random() * 200) / 100,
        //       "15min": Math.ceil(Math.random() * 200) / 100,
        //     },
        //     memoryUsage: Math.ceil(Math.random() * 100),
        //     swapUsage: Math.ceil(Math.random() * 100),
        //     disksUsage: [
        //       {
        //         name: "/dev/vda1",
        //         usage: Math.ceil(Math.random() * 100),
        //       },
        //       {
        //         name: "/dev/vda2",
        //         usage: Math.ceil(Math.random() * 100),
        //       },
        //     ],
        //   };

        //   // needed for reactivity (see https://vuejs.org/v2/guide/reactivity.html#For-Objects)
        //   this.$set(this.nodesStatus, nodeId, nodeStatus);
        // }
      }
    },
    getNodeStatusCompleted(taskContext, taskResult) {
      console.log("getNodeStatusCompleted, nodeId", taskContext.extra.node); ////

      const nodeId = taskContext.extra.node;
      const nodeStatus = taskResult.output;

      // round cpu load (sometimes it has roundoff error)
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
    // getNodeStatusAborted(taskContext, taskResult) { ////
    //   console.log("getNodeStatusAborted", taskResult); ////
    // },
    goToNodeDetail(nodeId) {
      this.$router.push({
        name: "NodeDetail",
        params: { nodeId },
      });
    },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";

.nodes-loader {
  margin: $spacing-05 auto;
  color: $interactive-01;
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

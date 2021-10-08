<template>
  <div>
    <div class="bx--grid bx--grid--full-width">
      <div class="bx--row">
        <div class="bx--col-lg-16">
          <cv-breadcrumb
            aria-label="breadcrumb"
            :no-trailing-slash="true"
            class="breadcrumb"
          >
            <cv-breadcrumb-item>
              <cv-link to="/nodes">{{ $t("nodes.title") }}</cv-link>
            </cv-breadcrumb-item>
            <cv-breadcrumb-item>
              <span>{{ $t("common.node") }} {{ nodeId }}</span>
            </cv-breadcrumb-item>
          </cv-breadcrumb>
        </div>
      </div>
      <div class="bx--row">
        <div class="bx--col-lg-16 page-subtitle title-and-role">
          <h3 class="title">{{ $t("common.node") }} {{ nodeId }}</h3>
          <cv-tag
            v-if="isLeader"
            kind="green"
            :label="$t('nodes.leader')"
          ></cv-tag>
          <cv-tag v-else kind="blue" :label="$t('nodes.worker')"></cv-tag>
        </div>
      </div>
      <div class="bx--row">
        <div class="bx--col-md-4 bx--col-max-4">
          <cv-tile :light="true" class="content-tile same-height-tile">
            <h4>{{ $t("node_detail.cpu") }}</h4>
            <MeterChart
              :label="$t('node_detail.usage')"
              :value="loading.nodeStatus ? 0 : nodeStatus.cpu.usage"
              class="mg-bottom-md"
            />
            <div class="mg-bottom-sm">
              <span class="label"
                >{{ $t("node_detail.load") }}
                <cv-tooltip
                  alignment="start"
                  direction="bottom"
                  :tip="$t('nodes.cpu_load_tooltip')"
                  class="info"
                >
                  <Information16 />
                </cv-tooltip>
              </span>
              <span v-if="loading.nodeStatus">- / - / -</span>
              <template v-else>
                <span
                  :class="{
                    warning:
                      nodeStatus.load['1min'] >= this.CPU_LOAD_WARNING_TH,
                  }"
                  >{{ nodeStatus.load["1min"] }}%</span
                >
                /
                <span
                  :class="{
                    warning:
                      nodeStatus.load['5min'] >= this.CPU_LOAD_WARNING_TH,
                  }"
                  >{{ nodeStatus.load["5min"] }}%</span
                >
                /
                <span
                  :class="{
                    warning:
                      nodeStatus.load['15min'] >= this.CPU_LOAD_WARNING_TH,
                  }"
                  >{{ nodeStatus.load["15min"] }}%</span
                >
              </template>
            </div>
            <template v-if="!loading.nodeStatus">
              <div
                v-for="(core, index) in nodeStatus.cpu.info"
                :key="index"
                class="mg-bottom-sm"
              >
                <span class="label"
                  >{{ $t("node_detail.core") }} {{ index + 1 }}</span
                >
                <span
                  >{{ nodeStatus.cpu.info[index].vendor }} -
                  {{ nodeStatus.cpu.info[index].model }}</span
                >
              </div>
            </template>
          </cv-tile>
        </div>
        <div class="bx--col-md-4 bx--col-max-4">
          <cv-tile :light="true" class="content-tile same-height-tile">
            <h4>{{ $t("node_detail.memory") }}</h4>
            <MeterChart
              :label="$t('node_detail.usage')"
              :value="loading.nodeStatus ? 0 : nodeStatus.memory.usage"
              class="mg-bottom-md"
            />
            <div class="mg-bottom-sm">
              <span class="label">{{ $t("node_detail.total") }}</span>
              <span>{{
                loading.nodeStatus ? "-" : nodeStatus.memory.total | byteFormat
              }}</span>
            </div>
            <div class="mg-bottom-sm">
              <span class="label">{{ $t("node_detail.used") }}</span>
              <span>{{
                loading.nodeStatus ? "-" : nodeStatus.memory.used | byteFormat
              }}</span>
            </div>
            <div class="mg-bottom-sm">
              <span class="label">{{ $t("node_detail.free") }}</span>
              <span>{{
                loading.nodeStatus ? "-" : nodeStatus.memory.free | byteFormat
              }}</span>
            </div>
          </cv-tile>
        </div>
        <div class="bx--col-md-4 bx--col-max-4">
          <cv-tile :light="true" class="content-tile same-height-tile">
            <h4>{{ $t("node_detail.swap") }}</h4>
            <MeterChart
              :label="$t('node_detail.usage')"
              :value="loading.nodeStatus ? 0 : nodeStatus.swap.usage"
              class="mg-bottom-md"
            />
            <div class="mg-bottom-sm">
              <span class="label">{{ $t("node_detail.total") }}</span>
              <span>{{
                loading.nodeStatus ? "-" : nodeStatus.swap.total | byteFormat
              }}</span>
            </div>
            <div class="mg-bottom-sm">
              <span class="label">{{ $t("node_detail.used") }}</span>
              <span>{{
                loading.nodeStatus ? "-" : nodeStatus.swap.used | byteFormat
              }}</span>
            </div>
            <div class="mg-bottom-sm">
              <span class="label">{{ $t("node_detail.free") }}</span>
              <span>{{
                loading.nodeStatus ? "-" : nodeStatus.swap.free | byteFormat
              }}</span>
            </div>
          </cv-tile>
        </div>
        <div class="bx--col-md-4 bx--col-max-4">
          <cv-tile :light="true" class="content-tile same-height-tile">
            <h4>{{ $t("node_detail.vpn") }}</h4>
            <template v-if="!loading.clusterStatus && vpnInfo.endpoint">
              <div class="mg-bottom-sm">
                <span class="label">{{ $t("node_detail.endpoint") }}</span>
                <span>{{ vpnInfo.endpoint }}</span>
              </div>
            </template>
            <div class="mg-bottom-sm">
              <span class="label">{{ $t("node_detail.ip_address") }}</span>
              <span>{{
                loading.clusterStatus ? "-" : vpnInfo.ip_address
              }}</span>
            </div>
            <div class="mg-bottom-sm">
              <span class="label">{{ $t("node_detail.listen_port") }}</span>
              <span>{{
                loading.clusterStatus ? "-" : vpnInfo.listen_port
              }}</span>
            </div>
            <template v-if="!loading.clusterStatus && vpnInfo.rcvd">
              <div class="mg-bottom-sm">
                <span class="label">{{ $t("node_detail.data_received") }}</span>
                <span>{{ vpnInfo.rcvd | byteFormat }}</span>
              </div>
            </template>
            <template v-if="!loading.clusterStatus && vpnInfo.sent">
              <div class="mg-bottom-sm">
                <span class="label">{{ $t("node_detail.data_sent") }}</span>
                <span>{{ vpnInfo.sent | byteFormat }}</span>
              </div>
            </template>
            <template v-if="!loading.clusterStatus && vpnInfo.last_seen">
              <div class="mg-bottom-sm">
                <span class="label">{{ $t("node_detail.last_seen") }}</span>
                <span
                  :class="{
                    warning:
                      new Date().getTime() - vpnInfo.last_seen * 1000 >=
                      this.LAST_SEEN_WARNING_TH,
                  }"
                >
                  {{
                    formatDateDistance(vpnInfo.last_seen * 1000, new Date(), {
                      addSuffix: true,
                    })
                  }}
                </span>
              </div>
            </template>
          </cv-tile>
        </div>
      </div>
      <div class="bx--row">
        <div class="bx--col-lg-16">
          <h4>{{ $t("node_detail.disks") }}</h4>
        </div>
      </div>
      <div class="bx--row">
        <div v-if="loading.nodeStatus" class="bx--col-md-4 bx--col-max-4">
          <cv-tile light>
            <cv-skeleton-text
              :paragraph="true"
              :line-count="10"
            ></cv-skeleton-text>
          </cv-tile>
        </div>
        <template v-else>
          <div
            v-for="(disk, index) in nodeStatus.disks"
            :key="index"
            class="bx--col-md-4 bx--col-max-4"
          >
            <cv-tile :light="true" class="content-tile">
              <h4>{{ $t("node_detail.disk") }} {{ index + 1 }}</h4>
              <MeterChart
                :label="$t('node_detail.usage')"
                :value="disk.usage"
                class="mg-bottom-md"
              />
              <div class="mg-bottom-sm">
                <span class="label">{{ $t("node_detail.device") }}</span>
                <span>{{ disk.device }}</span>
              </div>
              <div class="mg-bottom-sm">
                <span class="label">{{ $t("node_detail.mount_point") }}</span>
                <span>{{ disk.mountpoint }}</span>
              </div>
              <div class="mg-bottom-sm">
                <span class="label">{{ $t("node_detail.file_system") }}</span>
                <span>{{ disk.fstype }}</span>
              </div>
              <div class="mg-bottom-sm">
                <span class="label">{{ $t("node_detail.total") }}</span>
                <span>{{ disk.total | byteFormat }}</span>
              </div>
              <div class="mg-bottom-sm">
                <span class="label">{{ $t("node_detail.used") }}</span>
                <span>{{ disk.used | byteFormat }}</span>
              </div>
              <div class="mg-bottom-sm">
                <span class="label">{{ $t("node_detail.free") }}</span>
                <span>{{ disk.free | byteFormat }}</span>
              </div>
            </cv-tile>
          </div>
        </template>
      </div>
    </div>
  </div>
</template>

<script>
import {
  QueryParamService,
  UtilService,
  TaskService,
  IconService,
  DateTimeService,
} from "@nethserver/ns8-ui-lib";
import MeterChart from "@/components/MeterChart";
import to from "await-to-js";
import Information16 from "@carbon/icons-vue/es/information/16";

export default {
  name: "NodeDetail",
  components: { MeterChart, Information16 },
  mixins: [
    TaskService,
    UtilService,
    IconService,
    QueryParamService,
    DateTimeService,
  ],
  pageTitle() {
    return this.$t("node_detail.title");
  },
  data() {
    return {
      q: {},
      NODE_STATUS_TIME_INTERVAL: 5000,
      CLUSTER_STATUS_TIME_INTERVAL: 5000,
      CPU_LOAD_WARNING_TH: 90,
      LAST_SEEN_WARNING_TH: 5 * 60 * 1000, // milliseconds: 5 minutes
      nodeId: "",
      nodeStatus: {},
      nodeStatusInterval: null,
      clusterStatusInterval: null,
      isLeader: false,
      vpnInfo: {},
      loading: {
        nodeStatus: true,
        clusterStatus: true,
      },
    };
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
    this.nodeId = this.$route.params.nodeId;
    this.retrieveNodeStatus();
    this.retrieveClusterStatus();

    // periodically retrieve nodes status
    this.nodeStatusInterval = setInterval(
      this.retrieveNodeStatus,
      this.NODE_STATUS_TIME_INTERVAL
    );

    // periodically retrieve cluster status to get "last seen"
    this.clusterStatusInterval = setInterval(
      this.retrieveClusterStatus,
      this.CLUSTER_STATUS_TIME_INTERVAL
    );
  },
  beforeDestroy() {
    clearInterval(this.nodeStatusInterval);
    clearInterval(this.clusterStatusInterval);
  },
  methods: {
    async retrieveNodeStatus() {
      const taskAction = "get-node-status";

      // register to task events
      this.$root.$once(
        taskAction + "-completed-node-" + this.nodeId,
        this.getNodeStatusCompleted
      );

      const res = await to(
        this.createNodeTask(this.nodeId, {
          action: taskAction,
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
            node: this.nodeId,
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
    },
    getNodeStatusCompleted(taskContext, taskResult) {
      console.log("getNodeStatusCompleted", taskResult.output); ////

      const nodeStatus = taskResult.output;

      // round cpu load (sometimes it has roundoff error)
      nodeStatus.load["1min"] = Math.round(nodeStatus.load["1min"]);
      nodeStatus.load["5min"] = Math.round(nodeStatus.load["5min"]);
      nodeStatus.load["15min"] = Math.round(nodeStatus.load["15min"]);

      nodeStatus.memory.usage = Math.round(
        (nodeStatus.memory.used / nodeStatus.memory.total) * 100
      );
      nodeStatus.swap.usage = Math.round(
        (nodeStatus.swap.used / nodeStatus.swap.total) * 100
      );

      for (const disk of nodeStatus.disks) {
        disk.usage = Math.round((disk.used / disk.total) * 100);
      }
      this.nodeStatus = nodeStatus;
      this.loading.nodeStatus = false;
    },
    async retrieveClusterStatus() {
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
        this.createErrorNotification(
          err,
          this.$t("task.cannot_create_task", { action: taskAction })
        );
        return;
      }
    },
    getClusterStatusCompleted(taskContext, taskResult) {
      console.log("getClusterStatusCompleted"); ////

      const clusterStatus = taskResult.output;

      const currentNode = clusterStatus.nodes.find(
        (node) => node.id == this.nodeId
      );
      this.isLeader = currentNode.local;
      this.vpnInfo = currentNode.vpn;
      this.loading.clusterStatus = false;
    },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";

.label {
  font-weight: bold;
  margin-right: 0.5rem;
}

.warning {
  color: $danger-01;
}

.title-and-role {
  display: flex;

  .title {
    margin-right: $spacing-05;
  }
}

.same-height-tile {
  min-height: 14rem;
}
</style>

<style lang="scss">
@import "../styles/carbon-utils";

// global styles

// hide success status icon for meter charts
.bx--cc--meter-title .status-indicator.status--success {
  display: none;
}
</style>

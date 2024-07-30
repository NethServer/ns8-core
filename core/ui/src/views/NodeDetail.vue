<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
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
              <span>{{
                nodeLabel
                  ? nodeLabel + " (" + $t("common.node") + " " + nodeId + ")"
                  : $t("common.node") + " " + nodeId
              }}</span>
            </cv-breadcrumb-item>
          </cv-breadcrumb>
        </div>
      </div>
      <div class="bx--row">
        <div class="bx--col-lg-16 subpage-title title-and-role">
          <h3 class="title">
            {{
              nodeLabel
                ? nodeLabel + " (" + $t("common.node") + " " + nodeId + ")"
                : $t("common.node") + " " + nodeId
            }}
          </h3>
          <cv-tag
            v-if="isLeader"
            kind="green"
            :label="$t('nodes.leader')"
          ></cv-tag>
          <cv-tag v-else kind="blue" :label="$t('nodes.worker')"></cv-tag>
        </div>
      </div>
      <div v-if="!isOnline" class="bx--row">
        <div class="bx--col">
          <NsInlineNotification
            kind="error"
            :title="$t('nodes.node_is_offline')"
            :description="$t('nodes.node_is_offline_description')"
            :showCloseButton="false"
          />
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
        <div class="bx--col-lg-16">
          <h4 class="mg-bottom-md">{{ $t("node_detail.system") }}</h4>
        </div>
      </div>
      <div class="bx--row">
        <cv-column>
          <!-- card grid -->
          <div
            class="card-grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-3 3xl:grid-cols-4"
          >
            <cv-tile :light="true">
              <h4 class="mg-bottom-md">{{ $t("node_detail.cpu") }}</h4>
              <NsMeterChart
                :label="$t('node_detail.usage')"
                :value="loading.nodeStatus ? 0 : nodeStatus.cpu.usage"
                :loading="loading.nodeStatus"
                :warningThreshold="70"
                :dangerThreshold="90"
                progressBarHeight="10px"
                class="mg-bottom-lg"
              />
              <div class="mg-bottom-sm">
                <span class="label"
                  >{{ $t("node_detail.load") }}
                  <cv-interactive-tooltip
                    alignment="start"
                    direction="bottom"
                    class="info"
                  >
                    <template slot="trigger">
                      <Information16 />
                    </template>
                    <template slot="content">
                      {{ $t("nodes.cpu_load_tooltip") }}
                    </template>
                  </cv-interactive-tooltip>
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
                <template v-if="nodeStatus.cpu.info.length">
                  <!-- single core -->
                  <span
                    v-if="nodeStatus.cpu.info.length == 1"
                    class="label mg-bottom-xs"
                    >{{ $t("node_detail.core") }}</span
                  >
                  <!-- multi core -->
                  <span v-else class="label mg-bottom-xs">{{
                    $t("node_detail.num_core", {
                      num: nodeStatus.cpu.info.length,
                    })
                  }}</span>
                  <span
                    >{{ nodeStatus.cpu.info[0].vendor }} -
                    {{ nodeStatus.cpu.info[0].model }}</span
                  >
                </template>
              </template>
            </cv-tile>
            <cv-tile :light="true">
              <h4 class="mg-bottom-md">{{ $t("node_detail.memory") }}</h4>
              <NsMeterChart
                :label="$t('node_detail.usage')"
                :value="loading.nodeStatus ? 0 : nodeStatus.memory.usage"
                :loading="loading.nodeStatus"
                :warningThreshold="70"
                :dangerThreshold="90"
                progressBarHeight="10px"
                class="mg-bottom-lg"
              />
              <div class="mg-bottom-sm">
                <span class="label">{{ $t("node_detail.total") }}</span>
                <span>{{
                  loading.nodeStatus
                    ? "-"
                    : nodeStatus.memory.total | byteFormat
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
            <cv-tile :light="true">
              <h4 class="mg-bottom-md">{{ $t("node_detail.swap") }}</h4>
              <NsMeterChart
                :label="$t('node_detail.usage')"
                :value="
                  loading.nodeStatus || Number.isNaN(nodeStatus.swap.usage)
                    ? '-'
                    : nodeStatus.swap.usage
                "
                :loading="loading.nodeStatus"
                :warningThreshold="70"
                :dangerThreshold="90"
                progressBarHeight="10px"
                class="mg-bottom-lg"
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
            <cv-tile :light="true">
              <h4 class="mg-bottom-md">{{ $t("node_detail.vpn") }}</h4>
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
              <template v-if="!loading.clusterStatus && vpnInfo.listen_port">
                <div class="mg-bottom-sm">
                  <span class="label">{{ $t("node_detail.listen_port") }}</span>
                  <span>{{ vpnInfo.listen_port }}</span>
                </div>
              </template>
              <template v-if="!loading.clusterStatus && vpnInfo.rcvd">
                <div class="mg-bottom-sm">
                  <span class="label">{{
                    $t("node_detail.data_received")
                  }}</span>
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
        </cv-column>
      </div>
      <div class="bx--row">
        <div class="bx--col-lg-16">
          <h4 class="mg-bottom-md">{{ $t("node_detail.disks") }}</h4>
        </div>
      </div>
      <div class="bx--row">
        <div v-if="loading.nodeStatus" class="bx--col-md-4 bx--col-max-4">
          <cv-tile light>
            <cv-skeleton-text
              :paragraph="true"
              :line-count="7"
              heading
            ></cv-skeleton-text>
          </cv-tile>
        </div>
        <template v-else>
          <cv-column>
            <!-- card grid -->
            <div
              class="card-grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-3 3xl:grid-cols-4"
            >
              <cv-tile
                v-for="(disk, index) in nodeStatus.disks"
                :key="index"
                :light="true"
              >
                <h4 class="mg-bottom-md">
                  {{ $t("node_detail.disk") }} {{ index + 1 }}
                </h4>
                <NsMeterChart
                  :label="$t('node_detail.usage')"
                  :value="disk.usage"
                  :loading="loading.nodeStatus"
                  :warningThreshold="70"
                  :dangerThreshold="90"
                  progressBarHeight="10px"
                  class="mg-bottom-lg"
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
          </cv-column>
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
  PageTitleService,
} from "@nethserver/ns8-ui-lib";
import to from "await-to-js";
import Information16 from "@carbon/icons-vue/es/information/16";
import { mapActions } from "vuex";

export default {
  name: "NodeDetail",
  components: { Information16 },
  mixins: [
    TaskService,
    UtilService,
    IconService,
    QueryParamService,
    DateTimeService,
    PageTitleService,
  ],
  pageTitle() {
    return this.$t("node_detail.title");
  },
  data() {
    return {
      q: {},
      NODE_STATUS_TIME_INTERVAL: 30000,
      CLUSTER_STATUS_TIME_INTERVAL: 30000,
      CPU_LOAD_WARNING_TH: 90,
      LAST_SEEN_WARNING_TH: 5 * 60 * 1000, // milliseconds: 5 minutes
      nodeId: "",
      nodeStatus: {},
      nodeStatusInterval: null,
      clusterStatusInterval: null,
      isLeader: false,
      vpnInfo: {},
      nodeLabel: "-",
      isOnline: true,
      loading: {
        nodeStatus: true,
        clusterStatus: true,
      },
      error: {
        getNodeStatus: "",
        getClusterStatus: "",
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
    this.getClusterStatus();

    // periodically retrieve nodes status
    this.nodeStatusInterval = setInterval(
      this.retrieveNodeStatus,
      this.NODE_STATUS_TIME_INTERVAL
    );

    // periodically retrieve cluster status to get "last seen"
    this.clusterStatusInterval = setInterval(
      this.getClusterStatus,
      this.CLUSTER_STATUS_TIME_INTERVAL
    );
  },
  beforeDestroy() {
    clearInterval(this.nodeStatusInterval);
    clearInterval(this.clusterStatusInterval);
  },
  methods: {
    ...mapActions(["setClusterNodesInStore"]),
    async retrieveNodeStatus() {
      if (!this.isOnline) {
        return;
      }

      this.error.getNodeStatus = "";
      const taskAction = "get-node-status";
      const eventId = this.getUuid();

      // register to task events
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.getNodeStatusCompleted
      );

      const res = await to(
        this.createNodeTask(this.nodeId, {
          action: taskAction,
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
            node: this.nodeId,
            eventId,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.getNodeStatus = this.getErrorMessage(err);
        return;
      }
    },
    getNodeStatusCompleted(taskContext, taskResult) {
      const nodeStatus = taskResult.output;
      nodeStatus.cpu.usage = Math.round(nodeStatus.cpu.usage);
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
      const currentNode = clusterStatus.nodes.find(
        (node) => node.id == this.nodeId
      );
      this.isLeader = currentNode.local;
      this.vpnInfo = currentNode.vpn;
      this.nodeLabel = currentNode.ui_name;
      this.isOnline = currentNode.online;
      this.loading.clusterStatus = false;

      // update nodes in vuex store
      const nodes = clusterStatus.nodes.sort(this.sortByProperty("id"));
      this.setClusterNodesInStore(nodes);
    },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";

.label {
  font-weight: bold;
  display: inline-block;
  margin-right: $spacing-03;
}

.warning {
  color: $danger-01;
}

.title-and-role {
  display: flex;

  .title {
    margin-right: $spacing-05;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
}

@media (max-width: $breakpoint-medium) {
  .title-and-role .title {
    max-width: 15rem;
  }
}
</style>

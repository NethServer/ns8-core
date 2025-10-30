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
      <div v-if="error.listAlerts" class="bx--row">
        <div class="bx--col">
          <NsInlineNotification
            kind="error"
            :title="$t('action.list-alerts')"
            :description="error.listAlerts"
            :showCloseButton="false"
          />
        </div>
      </div>
      <div class="bx--row">
        <div class="bx--col-lg-16">
          <h4 class="mg-bottom-md">{{ $t("node_detail.node_information") }}</h4>
        </div>
      </div>
      <div class="bx--row">
        <cv-column>
          <!-- card grid -->
          <div
            class="card-grid nodeinfo-grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-3 3xl:grid-cols-4"
          >
            <cv-tile :light="true">
              <h4 class="mg-bottom-lg">{{ $t("node_detail.node_details") }}</h4>
              <template v-if="!loading.nodeStatus && nodeStatus.fqdn">
                <div class="mg-bottom-lg fqdn-row">
                  <span class="label">{{ $t("node_detail.fqdn") }}</span>
                  <span class="fqdn-text">{{ nodeStatus.fqdn }}</span>
                </div>
              </template>
              <div class="mg-bottom-lg">
                <span class="label">{{ $t("node_detail.node_id") }}</span>
                <span>{{ loading.nodeStatus ? "-" : nodeStatus.node_id }}</span>
              </div>
              <template v-if="!loading.nodeStatus && nodeStatus.app_count">
                <div class="mg-bottom-lg">
                  <span class="label">{{
                    $t("node_detail.applications_count")
                  }}</span>
                  <cv-link @click.prevent="goToApplications">
                    {{ nodeStatus.app_count }}
                  </cv-link>
                </div>
              </template>
              <template
                v-if="!loading.nodeStatus && nodeStatus.network_interface_count"
              >
                <div class="mg-bottom-lg">
                  <span class="label">{{
                    $t("node_detail.network_interface_count")
                  }}</span>
                  <cv-link @click.prevent="goToFirewall">
                    {{ nodeStatus.network_interface_count }}
                  </cv-link>
                </div>
              </template>
            </cv-tile>
            <cv-tile :light="true">
              <h4 class="mg-bottom-lg">{{ $t("node_detail.vpn") }}</h4>
              <template v-if="!loading.clusterStatus && vpnInfo.endpoint">
                <div class="mg-bottom-lg">
                  <span class="label">{{ $t("node_detail.endpoint") }}</span>
                  <span>{{ vpnInfo.endpoint }}</span>
                </div>
              </template>
              <div class="mg-bottom-lg">
                <span class="label">{{ $t("node_detail.ip_address") }}</span>
                <span>{{
                  loading.clusterStatus ? "-" : vpnInfo.ip_address
                }}</span>
              </div>
              <template
                v-if="!loading.nodeStatus && nodeStatus.vpn_listen_port"
              >
                <div class="mg-bottom-lg">
                  <span class="label">{{ $t("node_detail.listen_port") }}</span>
                  <span>{{ nodeStatus.vpn_listen_port }}</span>
                </div>
              </template>
            </cv-tile>
            <cv-tile class="alerts-tile" :light="true">
              <h4 class="mg-bottom-lg">{{ $t("node_detail.alarms") }}</h4>
              <NsDataTable
                :allRows="alerts"
                :columns="i18nTableColumnsAlarms"
                :rawColumns="tableColumnsAlarms"
                :sortable="true"
                :pageSizes="[5]"
                :overflow-menu="true"
                :isSearchable="false"
                :searchClearLabel="$t('common.clear_search')"
                :noSearchResultsLabel="$t('common.no_search_results')"
                :noSearchResultsDescription="
                  $t('common.no_search_results_description')
                "
                :isLoading="loading.alerts"
                :skeletonRows="5"
                :isErrorShown="!!error.listAlerts"
                :errorTitle="$t('action.list-alerts')"
                :errorDescription="error.listAlerts"
                :itemsPerPageLabel="$t('pagination.items_per_page')"
                :rangeOfTotalItemsLabel="$t('pagination.range_of_total_items')"
                :ofTotalPagesLabel="$t('pagination.of_total_pages')"
                :backwardText="$t('pagination.previous_page')"
                :forwardText="$t('pagination.next_page')"
                :pageNumberLabel="$t('pagination.page_number')"
                @updatePage="tablePage = $event"
              >
                <template slot="empty-state">
                  <template v-if="!alerts.length">
                    <NsEmptyState :title="$t('node_detail.no_alarms')">
                      <template #pictogram>
                        <CircleCheckPictogram />
                      </template>
                    </NsEmptyState>
                  </template>
                </template>
                <template slot="data">
                  <cv-data-table-row
                    v-for="(row, rowIndex) in tablePage"
                    :key="`${rowIndex}`"
                    :value="`${rowIndex}`"
                  >
                    <cv-data-table-cell class="severity-cell">
                      <span class="severity-content">
                        <NsSvg
                          v-if="row.labels.severity === 'warning'"
                          :svg="Warning16"
                          class="icon ns-warning"
                        />
                        <NsSvg
                          v-else-if="row.labels.severity === 'critical'"
                          :svg="ErrorFilled16"
                          class="icon ns-error"
                        >
                        </NsSvg>
                        <NsSvg
                          v-else
                          :svg="InformationFilled16"
                          class="icon ns-info"
                        />
                        <span class="summary-text">{{ row.summary }}</span>
                      </span>
                    </cv-data-table-cell>
                    <cv-data-table-cell class="log-cell">
                      <span class="log-cell-content">
                        {{ formatAlertDate(row.startsAt) }}
                        <cv-button
                          kind="ghost"
                          size="sm"
                          @click.prevent="goToLog"
                        >
                          {{ $t("node_detail.go_to_logs") }}
                        </cv-button>
                      </span>
                    </cv-data-table-cell>
                  </cv-data-table-row>
                </template>
              </NsDataTable>
            </cv-tile>
          </div>
        </cv-column>
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
              <h4 class="mg-bottom-lg">{{ $t("node_detail.cpu") }}</h4>
              <div class="mg-bottom-lg">
                {{ $t("node_detail.usage") }}
                {{ loading.nodeStatus ? 0 : nodeStatus.cpu.usage }}%
              </div>
              <NsProgressBar
                :value="loading.nodeStatus ? 0 : nodeStatus.cpu.usage"
                :loading="loading.nodeStatus"
                :warningThreshold="70"
                :dangerThreshold="90"
                :height="'10px'"
                class="mg-bottom-lg"
                :useHealthyColor="false"
              />
              <div class="mg-bottom-lg">
                <span class="label">
                  {{ $t("node_detail.load") }}
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
                        nodeStatus.load['1min'] / (nodeStatus.cpu.count || 1) >
                        5,
                    }"
                  >
                    {{ nodeStatus.load["1min"] }}
                  </span>
                  /
                  <span
                    :class="{
                      warning:
                        nodeStatus.load['5min'] / (nodeStatus.cpu.count || 1) >
                        5,
                    }"
                  >
                    {{ nodeStatus.load["5min"] }}
                  </span>
                  /
                  <span
                    :class="{
                      warning:
                        nodeStatus.load['15min'] / (nodeStatus.cpu.count || 1) >
                        5,
                    }"
                  >
                    {{ nodeStatus.load["15min"] }}
                  </span>
                </template>
              </div>
              <template v-if="!loading.nodeStatus">
                <template v-if="nodeStatus.cpu.info.length">
                  <!-- single core -->
                  <span
                    v-if="nodeStatus.cpu.info.count == 1"
                    class="label mg-bottom-xs"
                    >{{ $t("node_detail.core") }}</span
                  >
                  <!-- multi core -->
                  <span v-else class="label mg-bottom-xs">{{
                    $t("node_detail.num_core", {
                      num: nodeStatus.cpu.count,
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
              <h4 class="mg-bottom-lg">{{ $t("node_detail.memory") }}</h4>
              <div class="mg-bottom-lg">
                {{ $t("node_detail.usage") }}
                {{ loading.nodeStatus ? 0 : nodeStatus.memory.usage }}%
              </div>
              <NsProgressBar
                :value="loading.nodeStatus ? 0 : nodeStatus.memory.usage"
                :loading="loading.nodeStatus"
                :warningThreshold="70"
                :dangerThreshold="90"
                :height="'10px'"
                class="mg-bottom-lg"
                :useHealthyColor="false"
              />
              <div class="mg-bottom-lg">
                <span class="label">{{ $t("node_detail.total") }}</span>
                <span>{{
                  loading.nodeStatus
                    ? "-"
                    : nodeStatus.memory.total | byteFormat
                }}</span>
              </div>
              <div class="mg-bottom-lg">
                <span class="label">{{ $t("node_detail.used") }}</span>
                <span>{{
                  loading.nodeStatus ? "-" : nodeStatus.memory.used | byteFormat
                }}</span>
              </div>
              <div class="mg-bottom-lg">
                <span class="label">{{ $t("node_detail.free") }}</span>
                <span>{{
                  loading.nodeStatus ? "-" : nodeStatus.memory.free | byteFormat
                }}</span>
              </div>
            </cv-tile>
            <cv-tile :light="true">
              <h4 class="mg-bottom-lg">{{ $t("node_detail.swap") }}</h4>
              <div class="mg-bottom-lg">
                {{ $t("node_detail.usage") }}
                {{ loading.nodeStatus ? 0 : nodeStatus.swap.usage }}%
              </div>
              <NsProgressBar
                :value="
                  loading.nodeStatus || Number.isNaN(nodeStatus.swap.usage)
                    ? '-'
                    : nodeStatus.swap.usage
                "
                :loading="loading.nodeStatus"
                :warningThreshold="70"
                :dangerThreshold="90"
                :height="'10px'"
                class="mg-bottom-lg"
                :useHealthyColor="false"
              />
              <div class="mg-bottom-lg">
                <span class="label">{{ $t("node_detail.total") }}</span>
                <span>{{
                  loading.nodeStatus ? "-" : nodeStatus.swap.total | byteFormat
                }}</span>
              </div>
              <div class="mg-bottom-lg">
                <span class="label">{{ $t("node_detail.used") }}</span>
                <span>{{
                  loading.nodeStatus ? "-" : nodeStatus.swap.used | byteFormat
                }}</span>
              </div>
              <div class="mg-bottom-lg">
                <span class="label">{{ $t("node_detail.free") }}</span>
                <span>{{
                  loading.nodeStatus ? "-" : nodeStatus.swap.free | byteFormat
                }}</span>
              </div>
            </cv-tile>
          </div>
        </cv-column>
      </div>
      <div class="bx--row">
        <div class="bx--col-lg-16">
          <h4 class="mg-bottom-lg">{{ $t("node_detail.disks") }}</h4>
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
            <div>
              <cv-tile :light="true">
                <NsDataTable
                  :allRows="nodeStatus.disks"
                  :columns="i18nTableColumnsDisks"
                  :rawColumns="tableColumnsDisks"
                  :sortable="true"
                  :pageSizes="[5, 10, 15, 20]"
                  :overflow-menu="true"
                  :isSearchable="false"
                  :searchClearLabel="$t('common.clear_search')"
                  :noSearchResultsLabel="$t('common.no_search_results')"
                  :noSearchResultsDescription="
                    $t('common.no_search_results_description')
                  "
                  :isLoading="loading.nodeStatus"
                  :skeletonRows="5"
                  :isErrorShown="!!error.listNodes"
                  :errorTitle="$t('action.list-nodes')"
                  :errorDescription="error.listNodes"
                  :itemsPerPageLabel="$t('pagination.items_per_page')"
                  :rangeOfTotalItemsLabel="
                    $t('pagination.range_of_total_items')
                  "
                  :ofTotalPagesLabel="$t('pagination.of_total_pages')"
                  :backwardText="$t('pagination.previous_page')"
                  :forwardText="$t('pagination.next_page')"
                  :pageNumberLabel="$t('pagination.page_number')"
                  @updatePage="tablePageDisks = $event"
                >
                  <template slot="empty-state">
                    <template v-if="!nodeStatus.disks.length">
                      <NsEmptyState :title="$t('node_detail.no_disks_info')">
                        <template #pictogram>
                          <ExclamationMarkPictogram />
                        </template>
                      </NsEmptyState>
                    </template>
                  </template>
                  <template slot="data">
                    <cv-data-table-row
                      v-for="(row, rowIndex) in tablePageDisks"
                      :key="`${rowIndex}`"
                      :value="`${rowIndex}`"
                    >
                      <cv-data-table-cell>
                        {{ $t("node_detail.disk") }} {{ rowIndex + 1 }}
                      </cv-data-table-cell>
                      <cv-data-table-cell>
                        {{ row.device }}
                      </cv-data-table-cell>
                      <cv-data-table-cell>
                        {{ row.mountpoint }}
                      </cv-data-table-cell>
                      <cv-data-table-cell>
                        {{ row.fstype }}
                      </cv-data-table-cell>
                      <cv-data-table-cell>
                        {{ row.total | byteFormat }}
                      </cv-data-table-cell>
                      <cv-data-table-cell>
                        {{ row.used | byteFormat }} ({{ row.usage }}%)
                        <NsProgressBar
                          :value="row.usage"
                          :warningThreshold="70"
                          :dangerThreshold="90"
                          :loading="loading.nodeStatus"
                          :useHealthyColor="false"
                        />
                      </cv-data-table-cell>
                      <cv-data-table-cell>
                        {{ row.free | byteFormat }}
                      </cv-data-table-cell>
                    </cv-data-table-row>
                  </template>
                </NsDataTable>
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
      LAST_SEEN_WARNING_TH: 5 * 60 * 1000, // milliseconds: 5 minutes
      nodeId: "",
      nodeStatus: {},
      alerts: [],
      listNodesStatusInterval: null,
      listAlertsInterval: null,
      clusterStatusInterval: null,
      isLeader: false,
      vpnInfo: {},
      nodeLabel: "-",
      isOnline: true,
      loading: {
        nodeStatus: true,
        clusterStatus: true,
        alerts: true,
      },
      error: {
        getNodeStatus: "",
        getClusterStatus: "",
        listNodes: "",
        listAlerts: "",
      },
      tablePage: [],
      tablePageDisks: [],
      tableColumnsAlarms: ["summary", "startsAt"],
      tableColumnsDisks: [
        "disks",
        "devices",
        "mountpoints",
        "file_system",
        "total_space",
        "used_space",
        "free_space",
      ],
    };
  },
  computed: {
    i18nTableColumnsAlarms() {
      return this.tableColumnsAlarms.map((column) => {
        return this.$t("node_detail.col_" + column);
      });
    },
    i18nTableColumnsDisks() {
      return this.tableColumnsDisks.map((column) => {
        return this.$t("node_detail.col_" + column);
      });
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
    this.nodeId = this.$route.params.nodeId;
    this.getClusterStatus();
    this.listNodesStatus();
    this.listAlerts();

    // periodically retrieve nodes status
    this.listNodesStatusInterval = setInterval(
      this.listNodesStatus,
      this.NODE_STATUS_TIME_INTERVAL
    );
    // periodically retrieve alerts
    this.listAlertsInterval = setInterval(
      this.listAlerts,
      this.NODE_STATUS_TIME_INTERVAL
    );
    // periodically retrieve cluster status to get "last seen"
    this.clusterStatusInterval = setInterval(
      this.getClusterStatus,
      this.CLUSTER_STATUS_TIME_INTERVAL
    );
  },
  beforeDestroy() {
    clearInterval(this.listNodesStatusInterval);
    clearInterval(this.listAlertsInterval);
    clearInterval(this.clusterStatusInterval);
  },
  methods: {
    ...mapActions(["setClusterNodesInStore"]),
    goToFirewall() {
      this.$router.replace("/settings/firewall/" + this.nodeId);
    },
    goToApplications() {
      this.$router.replace({
        name: "applications",
        params: { nodeId: this.nodeId },
      });
    },
    goToLog() {
      this.$router.replace("/system-logs");
    },
    async listNodesStatus() {
      if (!this.isOnline) {
        return;
      }
      this.error.getNodeStatus = "";
      const taskAction = "list-nodes";

      // register to task completion
      this.$root.$once(taskAction + "-completed", this.listNodesCompleted);

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
        this.error.listNodes = this.getErrorMessage(err);
        return;
      }
    },
    listNodesCompleted(taskContext, taskResult) {
      // Filter by this.nodeId
      const nodeStatus = taskResult.output.nodes.find(
        (node) => Number(node.node_id) == Number(this.nodeId)
      );
      // Defensive defaults
      nodeStatus.cpu = nodeStatus.cpu || {};
      // Ensure usage numeric
      const usageVal = Number(nodeStatus.cpu.usage);
      nodeStatus.cpu.usage = Number.isFinite(usageVal)
        ? Math.round(usageVal * (usageVal < 1 ? 100 : 1)) // if backend sends fraction (e.g. 0.12), convert to %
        : 0;
      // Normalize load values as float (no percentage)
      nodeStatus.load = nodeStatus.load || { "1min": 0, "5min": 0, "15min": 0 };
      for (const k of ["1min", "5min", "15min"]) {
        const v = Number(nodeStatus.load[k]);
        nodeStatus.load[k] = Number.isFinite(v) ? v : 0;
      }

      // Build cpu.info array expected by template
      if (!Array.isArray(nodeStatus.cpu.info)) {
        nodeStatus.cpu.info = [
          {
            vendor: nodeStatus.cpu.vendor || "",
            model: nodeStatus.cpu.model_name || nodeStatus.cpu.model || "",
            count: nodeStatus.cpu.count || null,
          },
        ];
      }

      nodeStatus.memory = nodeStatus.memory || {
        used: 0,
        total: 0,
        free: 0,
        usage: 0,
      };
      nodeStatus.swap = nodeStatus.swap || {
        used: 0,
        total: 0,
        free: 0,
        usage: 0,
      };
      nodeStatus.disks = Array.isArray(nodeStatus.disks)
        ? nodeStatus.disks
        : [];

      // Memory %
      if (Number(nodeStatus.memory.total) > 0) {
        nodeStatus.memory.usage = Math.round(
          (nodeStatus.memory.used / nodeStatus.memory.total) * 100
        );
        nodeStatus.memory.free =
          nodeStatus.memory.total - nodeStatus.memory.used;
      } else {
        nodeStatus.memory.usage = 0;
      }
      // Swap %
      if (Number(nodeStatus.swap.total) > 0) {
        // Guard null used
        const swapUsed = Number(nodeStatus.swap.used) || 0;
        nodeStatus.swap.usage = Math.round(
          (swapUsed / nodeStatus.swap.total) * 100
        );
        nodeStatus.swap.free = nodeStatus.swap.total - swapUsed;
      } else {
        nodeStatus.swap.usage = 0;
      }

      // Disk usage %
      for (const disk of nodeStatus.disks) {
        const total = Number(disk.total);
        const used = Number(disk.used);
        if (total > 0 && Number.isFinite(used)) {
          disk.usage = Math.round((used / total) * 100);
        } else {
          disk.usage = 0;
        }
      }

      this.nodeStatus = nodeStatus;
      this.loading.nodeStatus = false;
    },
    async listAlerts() {
      this.error.listAlerts = "";
      const taskAction = "list-alerts";

      // register to task completion
      this.$root.$once(taskAction + "-completed", this.listAlertsCompleted);

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
        this.error.listAlerts = this.getErrorMessage(err);
        return;
      }
    },
    listAlertsCompleted(taskContext, taskResult) {
      // Filter by this.nodeId
      const filtered = taskResult.output.alerts.filter(
        (alert) => Number(alert.node_id) === Number(this.nodeId)
      );

      // Map to simplified objects (flatten annotations)
      const simplified = filtered.map((alert) => ({
        description: alert.annotations?.description || alert.description || "",
        summary: alert.annotations?.summary || alert.summary || "",
        endsAt: alert.endsAt || null,
        fingerprint: alert.fingerprint || null,
        receivers: alert.receivers || [],
        startsAt: alert.startsAt || null,
        status: alert.status || {},
        updatedAt: alert.updatedAt || null,
        generatorURL: alert.generatorURL || "",
        labels: alert.labels || {},
        name: alert.name || alert.labels?.alertname || "",
        category: alert.category || "",
        node_id: alert.node_id ?? Number(alert.labels?.node) ?? null,
        raw: alert, // keep original
      }));

      // Sort newest first
      simplified.sort((a, b) => {
        const ta = Date.parse(b.startsAt || 0);
        const tb = Date.parse(a.startsAt || 0);
        return ta - tb;
      });

      this.alerts = simplified;
      this.loading.alerts = false;
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
    formatAlertDate(dt) {
      if (!dt) return "-";
      const d = new Date(dt);
      if (isNaN(d)) return dt; // fallback if parse fails
      return new Intl.DateTimeFormat(navigator.language, {
        dateStyle: "medium",
        timeStyle: "short",
      }).format(d);
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

.severity-cell .icon {
  display: inline-block;
  width: 3rem; // increased (≈22px)
  height: 1.4rem;
  vertical-align: middle;
}

.severity-cell .severity-content {
  display: inline-flex;
  align-items: flex-start; // or center
  gap: 0.4rem;
  line-height: 1.2;
}

.severity-cell .summary-text {
  line-height: 1.3;
  padding-top: 2px; // slight vertical adjustment
}

@media (max-width: $breakpoint-medium) {
  .title-and-role .title {
    max-width: 15rem;
  }
}

.alerts-tile {
  min-width: 0; // avoid overflow
  overflow-x: auto;
}

.nodeinfo-grid {
  display: grid;
  grid-template-columns: 0.7fr 0.7fr 1.6fr; // 3 columns
  gap: 2rem; // space between tiles
}

.log-cell-content {
  display: inline-flex;
  align-items: center;
  gap: 0.75rem; // espace entre la date et le bouton
}

.log-cell .cv-button {
  margin: 0;
  padding: 0 0.5rem;
  font-size: 1em;
  height: 2rem;
  line-height: 1.2;
}

.fqdn-row {
  display: flex;
  align-items: baseline;
  flex-wrap: wrap; // allows wrapping only if really needed
  gap: 0.5rem;
}

.fqdn-text {
  white-space: normal;
  word-break: break-all;
  overflow-wrap: anywhere;
  max-width: 100%;
}
</style>

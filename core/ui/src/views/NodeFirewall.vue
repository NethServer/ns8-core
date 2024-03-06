<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <div>
    <cv-grid fullWidth>
      <cv-row>
        <cv-column>
          <cv-breadcrumb
            aria-label="breadcrumb"
            :no-trailing-slash="true"
            class="breadcrumb"
            v-if="route.includes('/nodes/')"
          >
            <cv-breadcrumb-item>
              <cv-link to="/nodes">{{ $t("nodes.title") }}</cv-link>
            </cv-breadcrumb-item>
            <cv-breadcrumb-item>
              <span>{{ title }}</span>
            </cv-breadcrumb-item>
          </cv-breadcrumb>
          <cv-breadcrumb
            aria-label="breadcrumb"
            :no-trailing-slash="true"
            class="breadcrumb"
            v-else
          >
            <cv-breadcrumb-item>
              <cv-link to="/settings">{{ $t("common.settings") }}</cv-link>
            </cv-breadcrumb-item>
            <cv-breadcrumb-item>
              <cv-link to="/settings/firewall">{{
                $t("firewall.title")
              }}</cv-link>
            </cv-breadcrumb-item>
            <cv-breadcrumb-item>
              <span>{{ title }}</span>
            </cv-breadcrumb-item>
          </cv-breadcrumb>
        </cv-column>
      </cv-row>
      <cv-row>
        <cv-column class="page-title">
          <h2>
            {{ title }}
          </h2>
        </cv-column>
      </cv-row>
      <cv-row class="landscape-warning">
        <cv-column>
          <NsInlineNotification
            kind="warning"
            :title="$t('common.use_landscape_mode')"
            :description="$t('common.use_landscape_mode_description')"
          />
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
      <cv-row v-if="error.getFirewallStatus">
        <cv-column>
          <NsInlineNotification
            kind="error"
            :title="$t('action.get-firewall-status')"
            :description="error.getFirewallStatus"
            :showCloseButton="false"
          />
        </cv-column>
      </cv-row>
      <cv-row>
        <cv-column>
          <h4 class="mg-bottom-md">
            {{ $t("firewall.public_zone_services") }}
          </h4>
        </cv-column>
      </cv-row>
      <cv-row>
        <cv-column>
          <cv-tile light>
            <cv-grid class="no-padding">
              <cv-row>
                <cv-column>
                  <NsDataTable
                    :allRows="servicesAndPorts"
                    :columns="i18nTableColumns"
                    :rawColumns="tableColumns"
                    :sortable="true"
                    :pageSizes="[10, 25, 50, 100]"
                    isSearchable
                    :searchPlaceholder="$t('firewall.search_service_or_port')"
                    :searchClearLabel="$t('common.clear_search')"
                    :noSearchResultsLabel="$t('common.no_search_results')"
                    :noSearchResultsDescription="
                      $t('common.no_search_results_description')
                    "
                    :isLoading="loading.getFirewallStatus"
                    :skeletonRows="5"
                    :isErrorShown="!!error.getFirewallStatus"
                    :errorTitle="$t('action.get-firewall-status')"
                    :errorDescription="error.getFirewallStatus"
                    :itemsPerPageLabel="$t('pagination.items_per_page')"
                    :rangeOfTotalItemsLabel="
                      $t('pagination.range_of_total_items')
                    "
                    :ofTotalPagesLabel="$t('pagination.of_total_pages')"
                    :backwardText="$t('pagination.previous_page')"
                    :forwardText="$t('pagination.next_page')"
                    :pageNumberLabel="$t('pagination.page_number')"
                    @updatePage="tablePage = $event"
                  >
                    <template slot="empty-state">
                      <NsEmptyState :title="$t('firewall.no_service')">
                        <template #description>
                          <div>
                            {{ $t("firewall.no_service_description") }}
                          </div>
                        </template>
                      </NsEmptyState>
                    </template>
                    <template slot="data">
                      <cv-data-table-row
                        v-for="(row, rowIndex) in tablePage"
                        :key="`${rowIndex}`"
                        :value="`${rowIndex}`"
                      >
                        <cv-data-table-cell>
                          <span>
                            {{ row.name }}
                          </span>
                        </cv-data-table-cell>
                        <cv-data-table-cell>
                          <span>{{ row.tcp_ports }}</span>
                        </cv-data-table-cell>
                        <cv-data-table-cell>
                          <span>{{ row.udp_ports }}</span>
                        </cv-data-table-cell>
                      </cv-data-table-row>
                    </template>
                  </NsDataTable>
                </cv-column>
              </cv-row>
            </cv-grid>
          </cv-tile>
        </cv-column>
      </cv-row>
      <cv-row>
        <cv-column>
          <h4 class="mg-bottom-md">
            {{ $t("firewall.network_interfaces") }}
          </h4>
        </cv-column>
      </cv-row>
      <cv-row>
        <!-- loading interfaces -->
        <template v-if="loading.getFirewallStatus">
          <cv-column v-for="index in 2" :key="index" :md="4" :max="4">
            <cv-tile kind="standard" light>
              <cv-skeleton-text
                :paragraph="true"
                :line-count="5"
              ></cv-skeleton-text>
            </cv-tile>
          </cv-column>
        </template>
        <!-- no interface -->
        <cv-column v-else-if="!interfaces.length">
          <cv-tile light>
            <NsEmptyState :title="$t('firewall.no_interface')"></NsEmptyState>
          </cv-tile>
        </cv-column>
        <!-- interfaces -->
        <template v-else>
          <cv-column v-for="iface in interfaces" :key="iface" :md="4" :max="4">
            <NsInfoCard
              light
              :title="iface"
              :icon="Network_232"
              :loading="loading.getFirewallStatus"
              :isErrorShown="error.getFirewallStatus"
              :errorTitle="$t('action.get-firewall-status')"
              :errorDescription="error.getFirewallStatus"
            >
            </NsInfoCard>
          </cv-column>
        </template>
      </cv-row>
    </cv-grid>
  </div>
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
import { mapActions } from "vuex";

export default {
  name: "NodeFirewall",
  mixins: [
    TaskService,
    UtilService,
    IconService,
    QueryParamService,
    PageTitleService,
  ],
  pageTitle() {
    return this.$t("firewall.title");
  },
  data() {
    return {
      q: {},
      nodeId: "",
      nodeLabel: "-",
      servicesAndPorts: [],
      interfaces: [],
      tablePage: [],
      tableColumns: ["name", "tcp_ports", "udp_ports"],
      loading: {
        getFirewallStatus: true,
        getClusterStatus: true,
      },
      error: {
        getFirewallStatus: "",
        getClusterStatus: "",
      },
    };
  },
  computed: {
    title() {
      return this.$t("firewall.node_firewall", {
        node: this.nodeLabel
          ? this.nodeLabel
          : this.$t("common.node") + " " + this.nodeId,
      });
    },
    route() {
      return this.$route.path;
    },
    i18nTableColumns() {
      return this.tableColumns.map((column) => {
        return this.$t("firewall." + column);
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
    this.getFirewallStatus();
  },
  methods: {
    ...mapActions(["setClusterNodesInStore"]),
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
      const currentNode = clusterStatus.nodes.find(
        (node) => node.id == this.nodeId
      );
      this.nodeLabel = currentNode.ui_name;
      this.loading.clusterStatus = false;

      // update nodes in vuex store
      const nodes = clusterStatus.nodes.sort(this.sortByProperty("id"));
      this.setClusterNodesInStore(nodes);
    },
    async getFirewallStatus() {
      this.error.getFirewallStatus = "";
      const taskAction = "get-firewall-status";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.getFirewallStatusAborted
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.getFirewallStatusCompleted
      );

      const res = await to(
        this.createNodeTask(this.nodeId, {
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
        this.error.getFirewallStatus = this.getErrorMessage(err);
        return;
      }
    },
    getFirewallStatusAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.getFirewallStatus = this.$t("error.generic_error");
      this.loading.getFirewallStatus = false;
    },
    getFirewallStatusCompleted(taskContext, taskResult) {
      const firewallStatus = taskResult.output;
      this.interfaces = firewallStatus.interfaces;
      const tcpPortRegex = /^(.+)\/tcp$/;
      const udpPortRegex = /^(.+)\/udp$/;

      // services

      for (const service of firewallStatus.services) {
        const tcpPorts = [];
        const udpPorts = [];

        for (const port of service.ports) {
          if (port.includes("tcp")) {
            const match = tcpPortRegex.exec(port);
            const tcpPort = match[1];
            tcpPorts.push(tcpPort);
          } else if (port.includes("udp")) {
            const match = udpPortRegex.exec(port);
            const updPort = match[1];
            udpPorts.push(updPort);
          }
        }
        service.tcp_ports = tcpPorts.join(", ");
        service.udp_ports = udpPorts.join(", ");
      }
      const servicesAndPorts = firewallStatus.services;

      // ports with no service associated

      for (const port of firewallStatus.ports) {
        if (port.includes("tcp")) {
          const match = tcpPortRegex.exec(port);
          const tcpPort = match[1];
          servicesAndPorts.push({
            name: "-",
            tcp_ports: tcpPort,
            udp_ports: "",
          });
        } else if (port.includes("udp")) {
          const match = udpPortRegex.exec(port);
          const updPort = match[1];
          servicesAndPorts.push({
            name: "-",
            tcp_ports: "",
            udp_ports: updPort,
          });
        }
      }
      this.servicesAndPorts = servicesAndPorts;
      this.loading.getFirewallStatus = false;
    },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";
</style>

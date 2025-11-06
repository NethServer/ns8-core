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
          </cv-tile>
        </cv-column>
      </cv-row>
      <cv-row>
        <cv-column>
          <h4 class="mg-bottom-md">
            {{ $t("firewall.active_network_interfaces") }}
          </h4>
        </cv-column>
      </cv-row>
      <cv-row>
        <cv-column>
          <cv-tile light>
            <cv-row>
              <cv-column>
                <NsDataTable
                  :allRows="ipaddr"
                  :columns="i18nTableColumnsNetworkInterfaces"
                  :rawColumns="tableColumnsNetworkInterfaces"
                  :sortable="true"
                  :pageSizes="[5, 10, 15, 20]"
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
                  @updatePage="tablePageNetworkInterfaces = $event"
                >
                  <template slot="data">
                    <cv-data-table-row
                      v-for="(row, rowIndex) in tablePageNetworkInterfaces"
                      :key="`${rowIndex}`"
                      :value="`${rowIndex}`"
                    >
                      <cv-data-table-cell>
                        <span>
                          {{ row.ifname }}
                        </span>
                      </cv-data-table-cell>
                      <cv-data-table-cell>
                        <span v-if="row.addresses.length <= 3">
                          <span v-for="(addr, idx) in row.addresses" :key="idx">
                            {{ addr
                            }}<br v-if="idx < row.addresses.length - 1" />
                          </span>
                        </span>
                        <span v-else>
                          <span
                            v-for="(addr, idx) in row.addresses.slice(0, 3)"
                            :key="idx"
                          >
                            {{ addr }}<br />
                          </span>
                          <cv-interactive-tooltip
                            alignment="center"
                            direction="top"
                          >
                            <template slot="trigger">
                              <a>
                                {{
                                  $tc(
                                    "common.plus_others",
                                    row.addresses.length - 3,
                                    {
                                      num: row.addresses.length - 3,
                                    }
                                  )
                                }}
                              </a>
                            </template>
                            <template slot="content">
                              <div>
                                <span
                                  v-for="(addr, idx) in row.addresses.slice(3)"
                                  :key="idx"
                                >
                                  {{ addr
                                  }}<br
                                    v-if="
                                      idx < row.addresses.slice(3).length - 1
                                    "
                                  />
                                </span>
                              </div>
                            </template>
                          </cv-interactive-tooltip>
                        </span>
                      </cv-data-table-cell>
                    </cv-data-table-row>
                  </template>
                </NsDataTable>
              </cv-column>
            </cv-row>
          </cv-tile>
        </cv-column>
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
      ipaddr: [],
      tablePage: [],
      tableColumns: ["name", "tcp_ports", "udp_ports"],
      tablePageNetworkInterfaces: [],
      tableColumnsNetworkInterfaces: ["interface", "ip_addresses"],
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
    i18nTableColumnsNetworkInterfaces() {
      return this.tableColumnsNetworkInterfaces.map((column) => {
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

      // network interfaces - ipaddr parsing
      function parseIpaddr(ipaddr) {
        return ipaddr.map((iface) => ({
          ifname: iface.ifname,
          addresses: Array.isArray(iface.addresses)
            ? iface.addresses
                .map((addr) => `${addr.address}/${addr.prefixlen}`)
                .sort()
            : [],
        }));
      }
      this.ipaddr = parseIpaddr(firewallStatus.ipaddr);

      // interfaces
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

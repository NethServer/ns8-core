<!--
  Copyright (C) 2026 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <div>
    <cv-tile light>
      <!-- fullWidth: without it Carbon caps the grid at 99rem and centers it -->
      <cv-grid fullWidth class="no-padding">
        <cv-row>
          <cv-column>
            <p class="title-description mg-bottom-md">
              {{ $t("settings_http_routes.frontend_proxies_description") }}
            </p>
          </cv-column>
        </cv-row>
        <cv-row>
          <cv-column>
            <div class="data-table-filters">
              <cv-search
                :label="$t('common.search')"
                :placeholder="
                  $t('settings_http_routes.filter_frontend_proxies')
                "
                :clear-aria-label="$t('common.clear_search')"
                v-model="filter.text"
                :disabled="loadingProxies"
                size="large"
                ref="tableSearch"
                class="self-end filter-field"
              >
              </cv-search>
              <NsComboBox
                v-model="filter.nodeId"
                :label="$t('common.choose')"
                :title="$t('common.node')"
                :auto-filter="true"
                :auto-highlight="true"
                :options="nodesForFilter"
                :disabled="loadingProxies"
                class="filter-field"
              >
              </NsComboBox>
              <cv-link @click="clearFilters()" class="self-end mb-3 shrink-0"
                >{{ $t("common.clear_filters") }}
              </cv-link>
            </div>
          </cv-column>
        </cv-row>
        <!-- with no proxy configured, the button lives in the empty state -->
        <cv-row v-if="proxyConfigs.length" class="toolbar">
          <cv-column>
            <NsButton
              v-if="unconfiguredNodes.length || loadingProxies"
              kind="primary"
              :icon="Add20"
              :disabled="loadingProxies"
              @click="showAddProxyModal"
              data-test-id="add-frontend-proxy"
              >{{ $t("common.add_frontend_proxy") }}
            </NsButton>
            <cv-tooltip
              v-else
              alignment="start"
              direction="top"
              :tip="
                $t('settings_http_routes.add_frontend_proxy_disabled_message')
              "
            >
              <NsButton
                kind="primary"
                :icon="Add20"
                disabled
                data-test-id="add-frontend-proxy"
                >{{ $t("common.add_frontend_proxy") }}
              </NsButton>
            </cv-tooltip>
          </cv-column>
        </cv-row>
        <cv-row>
          <cv-column>
            <!-- NsDataTable replaces the rows with its error: report a partial
                 failure here, so the nodes that did answer stay readable -->
            <NsInlineNotification
              v-if="tableError && proxyConfigs.length"
              kind="error"
              :title="tableErrorTitle"
              :description="tableErrorDescription"
              :showCloseButton="false"
            />
            <NsDataTable
              :allRows="filteredProxyConfigs"
              :columns="i18nTableColumns"
              :rawColumns="tableColumns"
              :sortable="true"
              :pageSizes="[10, 25, 50, 100]"
              :overflow-menu="true"
              :isSearchable="false"
              :noSearchResultsLabel="$t('common.no_search_results')"
              :noSearchResultsDescription="
                $t('common.no_search_results_description')
              "
              :isLoading="loadingProxies"
              :skeletonRows="3"
              :isErrorShown="!!tableError && !proxyConfigs.length"
              :errorTitle="tableErrorTitle"
              :errorDescription="tableErrorDescription"
              :itemsPerPageLabel="$t('pagination.items_per_page')"
              :rangeOfTotalItemsLabel="$t('pagination.range_of_total_items')"
              :ofTotalPagesLabel="$t('pagination.of_total_pages')"
              :backwardText="$t('pagination.previous_page')"
              :forwardText="$t('pagination.next_page')"
              :pageNumberLabel="$t('pagination.page_number')"
              @updatePage="tablePage = $event"
            >
              <template slot="empty-state">
                <template v-if="hasActiveFilters && proxyConfigs.length">
                  <!-- no search results -->
                  <NsEmptyState
                    :title="$t('common.no_search_results')"
                    key="no-results-empty-state"
                  >
                    <template #description>
                      <div class="flex flex-col items-center text-center gap-2">
                        <p>
                          {{ $t("common.no_search_results_description") }}
                        </p>
                        <NsButton
                          kind="ghost"
                          size="field"
                          @click="clearFilters()"
                          >{{ $t("common.clear_filters") }}
                        </NsButton>
                      </div>
                    </template>
                  </NsEmptyState>
                </template>
                <template v-else>
                  <!-- no frontend proxy configured -->
                  <NsEmptyState
                    :title="$t('settings_http_routes.no_frontend_proxy')"
                    key="no-frontend-proxy-empty-state"
                  >
                    <template #pictogram>
                      <NetworkPictogram />
                    </template>
                    <template #description>
                      <div class="flex flex-col items-center gap-4">
                        <p>
                          {{
                            $t(
                              "settings_http_routes.no_frontend_proxy_description"
                            )
                          }}
                        </p>
                        <NsButton
                          v-if="unconfiguredNodes.length || loadingProxies"
                          kind="primary"
                          :icon="Add20"
                          :disabled="loadingProxies"
                          @click="showAddProxyModal"
                          data-test-id="add-frontend-proxy-empty-state"
                          >{{ $t("common.add_frontend_proxy") }}
                        </NsButton>
                        <!-- no node to configure: the modal would open on an empty combo box -->
                        <cv-tooltip
                          v-else
                          alignment="start"
                          direction="top"
                          :tip="
                            $t(
                              'settings_http_routes.add_frontend_proxy_no_node_message'
                            )
                          "
                        >
                          <NsButton
                            kind="primary"
                            :icon="Add20"
                            disabled
                            data-test-id="add-frontend-proxy-empty-state"
                            >{{ $t("common.add_frontend_proxy") }}
                          </NsButton>
                        </cv-tooltip>
                      </div>
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
                  <cv-data-table-cell>
                    <span>{{ row.node }}</span>
                  </cv-data-table-cell>
                  <cv-data-table-cell>
                    <div class="flex flex-col">
                      <span v-for="proxy in row.proxies" :key="proxy">
                        {{ proxy }}
                      </span>
                    </div>
                  </cv-data-table-cell>
                  <cv-data-table-cell>
                    <span>{{ row.depth }}</span>
                  </cv-data-table-cell>
                  <cv-data-table-cell class="table-overflow-menu-cell">
                    <cv-overflow-menu
                      flip-menu
                      class="table-overflow-menu"
                      :data-test-id="row.node + '-menu'"
                    >
                      <cv-overflow-menu-item
                        @click="showEditProxyModal(row)"
                        :data-test-id="row.node + '-edit'"
                      >
                        <NsMenuItem :icon="Edit20" :label="$t('common.edit')" />
                      </cv-overflow-menu-item>
                      <cv-overflow-menu-item
                        danger
                        @click="showDeleteProxyModal(row)"
                        :data-test-id="row.node + '-delete'"
                      >
                        <NsMenuItem
                          :icon="TrashCan20"
                          :label="$t('common.delete')"
                        />
                      </cv-overflow-menu-item>
                    </cv-overflow-menu>
                  </cv-data-table-cell>
                </cv-data-table-row>
              </template>
            </NsDataTable>
          </cv-column>
        </cv-row>
      </cv-grid>
    </cv-tile>
    <ConfigureFrontendProxyModal
      :isShown="isShownConfigureProxyModal"
      :nodes="isEditingProxy ? internalNodes : unconfiguredNodes"
      :proxyConfig="currentProxyConfig"
      :isEditing="isEditingProxy"
      @hide="hideConfigureProxyModal"
      @reloadTrustedProxies="getTrustedProxies"
    />
    <!-- delete frontend proxy modal -->
    <NsDangerDeleteModal
      :isShown="isShownDeleteProxyModal"
      :name="proxyConfigToDelete ? 'node' + proxyConfigToDelete.nodeId : ''"
      :title="
        $t('settings_http_routes.delete_frontend_proxy_node', {
          node: proxyConfigToDelete ? proxyConfigToDelete.node : '',
        })
      "
      :warning="$t('common.please_read_carefully')"
      :description="
        $t('settings_http_routes.delete_frontend_proxy_description', {
          node: proxyConfigToDelete ? proxyConfigToDelete.node : '',
        })
      "
      :typeToConfirm="
        $t('common.type_to_confirm', {
          name: proxyConfigToDelete ? 'node' + proxyConfigToDelete.nodeId : '',
        })
      "
      :loading="loading.deleteTrustedProxies"
      :isErrorShown="!!error.deleteTrustedProxies"
      :errorTitle="$t('action.set-trusted-proxies')"
      :errorDescription="error.deleteTrustedProxies"
      :isWarningShown="false"
      @hide="hideDeleteProxyModal"
      @confirmDelete="deleteTrustedProxies"
      data-test-id="delete-frontend-proxy-modal"
    >
      <template #explanation>
        <div class="mt-4">
          <NsInlineNotification
            kind="warning"
            :title="$t('settings_http_routes.traefik_will_be_restarted')"
            :description="
              $t('settings_http_routes.frontend_proxy_restart_message', {
                node: proxyConfigToDelete ? proxyConfigToDelete.node : '',
              })
            "
            :showCloseButton="false"
          />
        </div>
      </template>
    </NsDangerDeleteModal>
  </div>
</template>

<script>
import to from "await-to-js";
import { UtilService, TaskService, IconService } from "@nethserver/ns8-ui-lib";
import _cloneDeep from "lodash/cloneDeep";
import { mapState } from "vuex";
import ConfigureFrontendProxyModal from "@/components/settings/ConfigureFrontendProxyModal.vue";

export default {
  name: "HttpFrontendProxies",
  components: {
    ConfigureFrontendProxyModal,
  },
  mixins: [TaskService, UtilService, IconService],
  props: {
    // built by the parent view: nodes with a traefik instance carry traefikInstance
    nodes: {
      type: Array,
      required: true,
    },
    traefikInstances: {
      type: Array,
      required: true,
    },
    isLoadingInstances: {
      type: Boolean,
      default: false,
    },
    // the instance lookup belongs to the parent: without it there is nothing to read
    instancesError: {
      type: String,
      default: "",
    },
  },
  data() {
    return {
      tablePage: [],
      // must match row properties: NsDataTable searches and sorts on them
      tableColumns: ["node", "proxies", "depth"],
      proxyConfigs: [],
      // local on purpose: the parent binds q.selectedNodeId to the Routes tab filter
      filter: {
        text: "",
        // set to the "all" sentinel by repaintNodeFilter(), once the options exist
        nodeId: "",
      },
      currentProxyConfig: null,
      proxyConfigToDelete: null,
      isShownConfigureProxyModal: false,
      isShownDeleteProxyModal: false,
      isEditingProxy: false,
      currentErrorAction: "",
      currentErrorDescription: "",
      // [eventName, handler] pairs registered on $root
      readListeners: [],
      deleteListeners: [],
      loading: {
        getTrustedProxiesNum: 0,
        deleteTrustedProxies: false,
      },
      error: {
        getTrustedProxies: "",
        deleteTrustedProxies: "",
      },
    };
  },
  computed: {
    ...mapState(["isWebsocketConnected"]),
    i18nTableColumns() {
      return this.tableColumns.map((column) => {
        return this.$t("settings_http_routes." + column);
      });
    },
    loadingProxies() {
      return this.isLoadingInstances || this.loading.getTrustedProxiesNum > 0;
    },
    tableError() {
      return this.instancesError || this.error.getTrustedProxies;
    },
    tableErrorTitle() {
      return this.instancesError
        ? this.$t("action.list-installed-modules")
        : this.currentErrorAction;
    },
    tableErrorDescription() {
      return this.instancesError || this.currentErrorDescription;
    },
    // only nodes running a traefik instance can be configured
    internalNodes() {
      return this.nodes.filter((node) => node.traefikInstance);
    },
    unconfiguredNodes() {
      const configuredNodeIds = this.proxyConfigs.map(
        (proxyConfig) => proxyConfig.nodeId
      );

      return this.internalNodes.filter(
        (node) => !configuredNodeIds.includes(node.value)
      );
    },
    nodesForFilter() {
      if (!this.internalNodes.length) {
        return [];
      }

      // add "Any node" at the beginning of internalNodes array
      const nodes = _cloneDeep(this.internalNodes);

      // same sentinel as the node filter of the routes table
      nodes.unshift({
        name: "all",
        label: this.$t("common.any_node"),
        value: "all",
      });
      return nodes;
    },
    filteredProxyConfigs() {
      let proxyConfigs = this.proxyConfigs;

      // filter by node
      if (this.filter.nodeId && this.filter.nodeId !== "all") {
        proxyConfigs = proxyConfigs.filter((proxyConfig) => {
          return proxyConfig.nodeId === this.filter.nodeId;
        });
      }

      // filter by text
      if (this.filter.text.trim()) {
        const searchText = this.filter.text.toLowerCase().trim();
        proxyConfigs = proxyConfigs.filter((proxyConfig) => {
          return (
            proxyConfig.node.toLowerCase().includes(searchText) ||
            proxyConfig.proxies.some((proxy) =>
              proxy.toLowerCase().includes(searchText)
            ) ||
            proxyConfig.depth.toString().includes(searchText)
          );
        });
      }
      return proxyConfigs;
    },
    hasActiveFilters() {
      return (
        (this.filter.text && this.filter.text.trim() !== "") ||
        (this.filter.nodeId && this.filter.nodeId !== "all")
      );
    },
  },
  watch: {
    // NsComboBox paints its text from the options only when the value changes,
    // and the parent publishes the nodes well after this component is created
    nodesForFilter: function (nodesForFilter) {
      if (nodesForFilter.length) {
        this.repaintNodeFilter();
      }
    },
    // the parent republishes the instances after a Traefik restart killed the
    // websocket, which is also how this table recovers
    traefikInstances: function () {
      this.getTrustedProxies();
    },
    isWebsocketConnected: function (isConnected) {
      // safety net: nothing else clears the pending state if the request was cut off
      if (isConnected && this.loading.deleteTrustedProxies) {
        this.loading.deleteTrustedProxies = false;
        this.hideDeleteProxyModal();
      }
    },
  },
  created() {
    if (this.traefikInstances.length) {
      this.getTrustedProxies();
    }
  },
  beforeDestroy() {
    // a task completing after destroy would run a handler on a dead instance
    this.clearListeners(this.readListeners);
    this.clearListeners(this.deleteListeners);
  },
  methods: {
    registerListener(listeners, eventName, handler) {
      this.$root.$once(eventName, handler);
      listeners.push([eventName, handler]);
    },
    clearListeners(listeners) {
      listeners.forEach(([eventName, handler]) => {
        this.$root.$off(eventName, handler);
      });
      listeners.splice(0);
    },
    clearTableError() {
      this.error.getTrustedProxies = "";
      this.currentErrorAction = "";
      this.currentErrorDescription = "";
    },
    clearFilters() {
      this.filter.text = "";
      // "all", not "": NsComboBox only repaints its text on a matching option
      this.filter.nodeId = "all";
    },
    repaintNodeFilter() {
      const nodeId = this.filter.nodeId || "all";

      // the value has to change for NsComboBox to read the options again
      this.filter.nodeId = "";
      this.$nextTick(() => {
        this.filter.nodeId = nodeId;
      });
    },
    showAddProxyModal() {
      this.isEditingProxy = false;
      this.currentProxyConfig = null;
      this.isShownConfigureProxyModal = true;
    },
    showEditProxyModal(proxyConfig) {
      this.isEditingProxy = true;
      this.currentProxyConfig = proxyConfig;
      this.isShownConfigureProxyModal = true;
    },
    hideConfigureProxyModal() {
      this.isShownConfigureProxyModal = false;
    },
    showDeleteProxyModal(proxyConfig) {
      this.error.deleteTrustedProxies = "";
      this.proxyConfigToDelete = proxyConfig;
      this.isShownDeleteProxyModal = true;
    },
    hideDeleteProxyModal() {
      this.isShownDeleteProxyModal = false;
    },
    async getTrustedProxies() {
      // a handler of a previous round would decrement the counter of this one
      this.clearListeners(this.readListeners);
      this.proxyConfigs = [];
      // otherwise a transient error sticks after the nodes answer again
      this.clearTableError();
      // count the whole batch upfront: the counter must not reach zero between
      // two iterations
      this.loading.getTrustedProxiesNum = this.traefikInstances.length;

      for (const traefikInstance of this.traefikInstances) {
        const taskAction = "get-trusted-proxies";
        const eventId = this.getUuid();

        // register to task events

        this.registerListener(
          this.readListeners,
          `${taskAction}-aborted-${eventId}`,
          this.getTrustedProxiesAborted
        );

        this.registerListener(
          this.readListeners,
          `${taskAction}-completed-${eventId}`,
          this.getTrustedProxiesCompleted
        );

        const res = await to(
          this.createModuleTaskForApp(traefikInstance.id, {
            action: taskAction,
            extra: {
              title: this.$t("action." + taskAction),
              isNotificationHidden: true,
              traefikInstance: traefikInstance,
              eventId,
            },
          })
        );
        const err = res[0];

        if (err) {
          console.error(`error creating task ${taskAction}`, err);
          const errMessage = this.getErrorMessage(err);
          this.error.getTrustedProxies = errMessage;
          this.currentErrorAction = this.$t("action." + taskAction);
          this.currentErrorDescription = errMessage;
          this.loading.getTrustedProxiesNum--;
        }
      }
    },
    getTrustedProxiesAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.getTrustedProxies = this.$t("error.generic_error");
      this.currentErrorAction = this.$t("action." + taskContext.action);
      this.currentErrorDescription = this.$t("error.generic_error");
      this.loading.getTrustedProxiesNum--;
    },
    getTrustedProxiesCompleted(taskContext, taskResult) {
      const traefikId = taskContext.extra.traefikInstance.id;
      const nodeId = taskContext.extra.traefikInstance.node;
      const nodeUiName = taskContext.extra.traefikInstance.node_ui_name;
      // get-trusted-proxies returns a Python set: sort for a stable display
      const proxies = (taskResult.output.proxies || []).sort();

      // drop the previous row of this node: a late event must not duplicate it
      this.proxyConfigs = this.proxyConfigs.filter(
        (proxyConfig) => proxyConfig.nodeId !== nodeId
      );

      // nodes without any frontend proxy are not listed
      if (proxies.length) {
        const node = { id: nodeId, ui_name: nodeUiName };

        this.proxyConfigs.push({
          nodeId: nodeId,
          node: this.getShortNodeLabel(node) + ` (${traefikId})`,
          nodeLabel: this.getShortNodeLabel(node),
          traefikInstance: traefikId,
          proxies: proxies,
          proxies_str: proxies.join("\n"),
          depth: taskResult.output.depth || 0,
        });
        this.proxyConfigs.sort(this.sortByProperty("node"));
      }
      this.loading.getTrustedProxiesNum--;
    },
    async deleteTrustedProxies() {
      this.loading.deleteTrustedProxies = true;
      this.error.deleteTrustedProxies = "";
      const taskAction = "set-trusted-proxies";
      const eventId = this.getUuid();

      // register to task events

      this.registerListener(
        this.deleteListeners,
        `${taskAction}-aborted-${eventId}`,
        this.deleteTrustedProxiesAborted
      );

      this.registerListener(
        this.deleteListeners,
        `${taskAction}-completed-${eventId}`,
        this.deleteTrustedProxiesCompleted
      );

      const res = await to(
        this.createModuleTaskForApp(this.proxyConfigToDelete.traefikInstance, {
          action: taskAction,
          // an empty list clears the configuration and resets the depth to 0
          data: {
            proxies: [],
            depth: 0,
          },
          extra: {
            title: this.$t("settings_http_routes.delete_frontend_proxy_node", {
              node: this.proxyConfigToDelete.node,
            }),
            description: this.$t("common.processing"),
            eventId,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.deleteTrustedProxies = this.getErrorMessage(err);
        this.loading.deleteTrustedProxies = false;
        return;
      }
      // stay open until the task ends: closing now would hide an abort error.
      // If the Traefik restart kills the websocket first, the reconnection
      // watcher closes the modal and reads the configuration again.
    },
    deleteTrustedProxiesAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.deleteTrustedProxies = this.$t("error.generic_error");
      this.loading.deleteTrustedProxies = false;

      // the node may have been left untouched: show what it holds now
      this.getTrustedProxies();
    },
    deleteTrustedProxiesCompleted() {
      this.loading.deleteTrustedProxies = false;
      this.hideDeleteProxyModal();
      this.getTrustedProxies();
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";

// Carbon makes these full width: cap them to keep the design proportions
.filter-field {
  max-width: 18rem;
  width: 100%;
}
</style>

<!--
  Copyright (C) 2026 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <div>
    <cv-tile light>
      <!-- fullWidth: without it Carbon caps the grid at 99rem and centers it -->
      <cv-grid fullWidth class="no-padding">
        <cv-row class="toolbar">
          <cv-column>
            <NsButton
              kind="primary"
              :icon="Add20"
              @click="showCreateRouteModal"
              data-test-id="create-route"
              >{{ $t("settings_http_routes.create_route") }}
            </NsButton>
          </cv-column>
        </cv-row>
        <cv-row>
          <cv-column>
            <div class="data-table-filters">
              <cv-search
                :label="$t('common.search')"
                :placeholder="$t('settings_http_routes.search_route')"
                :clear-aria-label="$t('common.clear_search')"
                v-model="filter.text"
                :disabled="loadingRoutes"
                size="large"
                ref="tableSearch"
                class="self-end filter-field"
              >
              </cv-search>
              <NsComboBox
                v-model="internalSelectedNodeId"
                :label="$t('common.choose')"
                :title="$t('common.node')"
                :auto-filter="false"
                :auto-highlight="true"
                :options="nodesForFilter"
                :disabled="isLoadingInstances"
                class="filter-field"
              >
              </NsComboBox>
              <cv-link @click="clearFilters()" class="self-end mb-3 shrink-0"
                >{{ $t("common.clear_filters") }}
              </cv-link>
            </div>
          </cv-column>
        </cv-row>
        <cv-row>
          <cv-column>
            <!-- NsDataTable replaces the rows with its error: report a
                 partial failure here, so the nodes that did answer
                 stay readable -->
            <NsInlineNotification
              v-if="tableError && routes.length"
              kind="error"
              :title="tableErrorTitle"
              :description="tableErrorDescription"
              :showCloseButton="false"
            />
            <NsDataTable
              :allRows="filteredRoutes"
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
              :isLoading="loadingRoutes"
              :skeletonRows="5"
              :isErrorShown="!!tableError && !routes.length"
              :errorTitle="tableErrorTitle"
              :errorDescription="tableErrorDescription"
              :itemsPerPageLabel="$t('pagination.items_per_page')"
              :rangeOfTotalItemsLabel="$t('pagination.range_of_total_items')"
              :ofTotalPagesLabel="$t('pagination.of_total_pages')"
              :backwardText="$t('pagination.previous_page')"
              :forwardText="$t('pagination.next_page')"
              :pageNumberLabel="$t('pagination.page_number')"
              :customSortTable="sortHttpRoutes"
              @updatePage="tablePage = $event"
            >
              <template slot="empty-state">
                <template v-if="hasActiveFilters && routes.length">
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
                  <!-- no route configured -->
                  <NsEmptyState
                    :title="$t('settings_http_routes.no_http_route')"
                    key="no-http-route-empty-state"
                  >
                    <template #pictogram>
                      <NetworkPictogram />
                    </template>
                    <template #description>
                      <div>
                        {{
                          $t("settings_http_routes.no_http_route_description")
                        }}
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
                    <div class="flex items-center gap-2">
                      <span v-if="row.host && row.path">
                        {{ row.host }}{{ row.path }}
                      </span>
                      <span v-else>
                        {{ row.host || row.path }}
                      </span>
                      <cv-interactive-tooltip
                        v-if="row.lets_encrypt_status === 'pending'"
                        alignment="center"
                        direction="top"
                        class="shrink-0"
                      >
                        <template #trigger>
                          <WarningAltFilled16 class="ns-warning" />
                        </template>
                        <template #content>
                          {{
                            $t(
                              "settings_http_routes.cannot_obtain_tls_certificate"
                            )
                          }}
                          <cv-link @click="goToPendingCertificateLogs(row)">
                            {{ $t("settings_http_routes.show_logs") }}
                          </cv-link>
                        </template>
                      </cv-interactive-tooltip>
                    </div>
                  </cv-data-table-cell>
                  <cv-data-table-cell>
                    <cv-link @click="showRouteDetailModal(row)">
                      {{ row.name }}
                    </cv-link>
                  </cv-data-table-cell>
                  <cv-data-table-cell>
                    <span>
                      {{ $t(`settings_http_routes.${row.type}`) }}
                    </span>
                  </cv-data-table-cell>
                  <cv-data-table-cell>
                    <span>{{ row.node }}</span>
                  </cv-data-table-cell>
                  <cv-data-table-cell>
                    <NsTag
                      v-if="!row.user_created"
                      kind="gray"
                      size="sm"
                      :label="$t('settings_http_routes.automatic')"
                    />
                    <NsTag
                      v-if="
                        row.ip_allowlist !== undefined &&
                        row.ip_allowlist.length > 0
                      "
                      kind="high-contrast"
                      size="sm"
                      :label="$t('settings_http_routes.restricted')"
                      class="no-margin"
                    />
                  </cv-data-table-cell>
                  <cv-data-table-cell class="table-overflow-menu-cell">
                    <cv-overflow-menu
                      flip-menu
                      class="table-overflow-menu"
                      :data-test-id="row.name + '-menu'"
                    >
                      <cv-overflow-menu-item
                        @click="showRouteDetailModal(row)"
                        :data-test-id="row.name + '-details'"
                      >
                        <NsMenuItem
                          :icon="ArrowRight20"
                          :label="$t('common.see_details')"
                        />
                      </cv-overflow-menu-item>
                      <cv-overflow-menu-item
                        @click="showEditRouteModal(row)"
                        :data-test-id="row.name + '-edit'"
                      >
                        <NsMenuItem :icon="Edit20" :label="$t('common.edit')" />
                      </cv-overflow-menu-item>
                      <cv-overflow-menu-item
                        danger
                        @click="showDeleteRouteModal(row)"
                        :disabled="!row.user_created"
                        :data-test-id="row.name + '-delete'"
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
    <HttpRouteDetailModal
      :isShown="isShownRouteDetailModal"
      :route="currentRoute"
      @hide="hideRouteDetailModal"
    />
    <CreateOrEditHttpRouteModal
      :isShown="isShownCreateOrEditRouteModal"
      :nodes="nodes"
      :defaultNodeId="internalSelectedNodeId"
      :allRoutes="routes"
      :route="currentRoute"
      :isEditing="isEditingRoute"
      @hide="hideCreateRouteModal"
      @reloadRoutes="onReloadRoutes"
    />
    <!-- delete route modal -->
    <NsDangerDeleteModal
      :isShown="isShownDeleteRouteModal"
      :name="routeToDelete ? routeToDelete.name : ''"
      :title="
        $t('settings_http_routes.delete_route_route', {
          route: routeToDelete ? routeToDelete.name : '',
        })
      "
      :warning="$t('common.please_read_carefully')"
      :description="
        $t('settings_http_routes.delete_route_description', {
          name: routeToDelete ? routeToDelete.name : '',
        })
      "
      :typeToConfirm="
        $t('common.type_to_confirm', {
          name: routeToDelete ? routeToDelete.name : '',
        })
      "
      :isErrorShown="!!error.deleteRoute"
      :errorTitle="$t('action.delete-route')"
      :errorDescription="error.deleteRoute"
      :isWarningShown="false"
      @hide="hideDeleteRouteModal"
      @confirmDelete="deleteRoute"
      data-test-id="delete-route-modal"
    >
      <!-- traefik will be restarted if the route requested a certificate -->
      <template v-if="routeToDelete && routeToDelete.lets_encrypt" #explanation>
        <div class="mt-4">
          <NsInlineNotification
            kind="warning"
            :title="$t('settings_http_routes.traefik_will_be_restarted')"
            :description="
              $t('settings_http_routes.delete_route_with_certificate_message', {
                node: routeToDelete.node,
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
import {
  UtilService,
  TaskService,
  IconService,
  DateTimeService,
} from "@nethserver/ns8-ui-lib";
import _cloneDeep from "lodash/cloneDeep";
import HttpRouteDetailModal from "@/components/settings/HttpRouteDetailModal.vue";
import CreateOrEditHttpRouteModal from "@/components/settings/CreateOrEditHttpRouteModal.vue";
import WarningAltFilled16 from "@carbon/icons-vue/es/warning--alt--filled/16";

export default {
  name: "HttpRoutesPanel",
  components: {
    HttpRouteDetailModal,
    CreateOrEditHttpRouteModal,
    WarningAltFilled16,
  },
  mixins: [TaskService, UtilService, IconService, DateTimeService],
  props: {
    // cluster nodes, with traefikInstance set on those running an instance
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
    instancesError: {
      type: String,
      default: "",
    },
    // the node filter is a query param: the view owns it, this panel mirrors it
    selectedNodeId: {
      type: String,
      default: "",
    },
  },
  data() {
    return {
      // repainting the combo box needs a synchronous write, which a round trip
      // through the query param cannot guarantee
      internalSelectedNodeId: "",
      filter: {
        text: "",
      },
      tablePage: [],
      tableColumns: ["route", "name", "type", "node", "attributes"],
      routes: [],
      isShownCreateOrEditRouteModal: false,
      isShownRouteDetailModal: false,
      isShownDeleteRouteModal: false,
      currentErrorAction: "",
      currentErrorDescription: "",
      currentRoute: null,
      routeToDelete: null,
      isEditingRoute: false,
      // [eventName, handler] pairs registered on $root, cleared on destroy
      readListeners: [],
      // bumped by every read batch: an older one bails out when it resumes
      readGeneration: 0,
      deleteListeners: [],
      pendingCertificatesLogsPath: {},
      loading: {
        listRoutesNum: 0,
        deleteRoute: false,
      },
      error: {
        listRoutes: "",
        deleteRoute: "",
      },
    };
  },
  computed: {
    i18nTableColumns() {
      return this.tableColumns.map((column) => {
        return this.$t("settings_http_routes." + column);
      });
    },
    loadingRoutes() {
      return this.isLoadingInstances || this.loading.listRoutesNum > 0;
    },
    tableError() {
      return this.instancesError || this.error.listRoutes;
    },
    tableErrorTitle() {
      return this.instancesError
        ? this.$t("action.list-installed-modules")
        : this.currentErrorAction;
    },
    tableErrorDescription() {
      return this.instancesError || this.currentErrorDescription;
    },
    filteredRoutes() {
      let routes = this.routes;

      // filter by node
      if (
        this.internalSelectedNodeId &&
        this.internalSelectedNodeId !== "all"
      ) {
        routes = routes.filter((route) => {
          return route.nodeId === this.internalSelectedNodeId;
        });
      }

      // filter by text
      if (this.filter.text.trim()) {
        // the route reads as one string on screen: match it that way, ignoring
        // the separators the user is unlikely to type
        const cleanRegex = /[^a-zA-Z0-9]/g;
        const queryText = this.filter.text.replace(cleanRegex, "");
        const searchFields = ["host", "path", "name", "type", "node"];

        routes = routes.filter((route) => {
          if (route.host && route.path) {
            const hostAndPath = `${route.host}${route.path}`;

            if (
              new RegExp(queryText, "i").test(
                hostAndPath.replace(cleanRegex, "")
              )
            ) {
              return true;
            }
          }

          return searchFields.some((searchField) => {
            const searchValue = route[searchField];

            return searchValue
              ? new RegExp(queryText, "i").test(
                  searchValue.replace(cleanRegex, "")
                )
              : false;
          });
        });
      }
      return routes;
    },
    hasActiveFilters() {
      return (
        (this.filter.text && this.filter.text.trim() !== "") ||
        (this.internalSelectedNodeId && this.internalSelectedNodeId !== "all")
      );
    },
    nodesForFilter() {
      if (!this.nodes.length) {
        return [];
      }

      // add "All nodes" at the beginning of internalNodes array
      const nodes = _cloneDeep(this.nodes);

      nodes.unshift({
        name: "all",
        label: this.$t("common.any_node"),
        value: "all",
      });
      return nodes;
    },
  },
  watch: {
    selectedNodeId: function (selectedNodeId) {
      if (selectedNodeId !== this.internalSelectedNodeId) {
        this.internalSelectedNodeId = selectedNodeId;
      }
    },
    internalSelectedNodeId: function (internalSelectedNodeId) {
      // "" is a repaint artifact, never a user choice: keep it out of the URL
      if (
        internalSelectedNodeId &&
        internalSelectedNodeId !== this.selectedNodeId
      ) {
        this.$emit("update:selectedNodeId", internalSelectedNodeId);
      }
    },
    // the filter has no option to select until both the nodes and the instances
    // landed, in either order
    nodesForFilter: function (nodesForFilter) {
      if (nodesForFilter.length) {
        this.repaintNodeFilter();
      }
    },
    // the view republishes the instances after a Traefik restart killed the
    // websocket: reload on a new array, never on an in-place mutation
    traefikInstances: function () {
      this.listRoutes();
    },
    instancesError: function (instancesError) {
      // the routes chain will not run: release its skeleton rows
      if (instancesError) {
        this.clearListeners(this.readListeners);
        this.loading.listRoutesNum = 0;
      }
    },
  },
  created() {
    this.internalSelectedNodeId = this.selectedNodeId;

    // the view may have loaded the instances before this panel was created
    if (this.traefikInstances.length) {
      this.listRoutes();
    }
  },
  beforeDestroy() {
    // a task completing after destroy would re-fire the read chain
    this.clearListeners(this.readListeners);
    this.clearListeners(this.deleteListeners);
  },
  methods: {
    showCreateRouteModal() {
      this.isEditingRoute = false;
      this.isShownCreateOrEditRouteModal = true;
    },
    showEditRouteModal(route) {
      this.isEditingRoute = true;
      this.currentRoute = route;
      this.isShownCreateOrEditRouteModal = true;
    },
    hideCreateRouteModal() {
      this.isShownCreateOrEditRouteModal = false;
    },
    showRouteDetailModal(route) {
      this.currentRoute = route;
      this.isShownRouteDetailModal = true;
    },
    hideRouteDetailModal() {
      this.isShownRouteDetailModal = false;
    },
    showDeleteRouteModal(route) {
      this.error.deleteRoute = "";
      this.routeToDelete = route;
      this.isShownDeleteRouteModal = true;
    },
    hideDeleteRouteModal() {
      this.isShownDeleteRouteModal = false;
    },
    repaintNodeFilter() {
      this.$nextTick(() => {
        if (!this.internalSelectedNodeId) {
          // initially show all nodes
          this.internalSelectedNodeId = "all";
          return;
        }
        const nodeId = this.internalSelectedNodeId;

        // cv-combo-box only repaints its text when the value changes
        this.internalSelectedNodeId = "";
        this.$nextTick(() => {
          this.internalSelectedNodeId = nodeId;
        });
      });
    },
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
    async listRoutes() {
      // a handler of a previous chain would decrement the counter of this one
      this.clearListeners(this.readListeners);
      const generation = ++this.readGeneration;

      this.routes = [];
      // otherwise a transient error sticks after the nodes answer again
      this.error.listRoutes = "";
      this.currentErrorAction = "";
      this.currentErrorDescription = "";
      // count the whole batch upfront: the counter must not reach zero between
      // two iterations
      this.loading.listRoutesNum = this.traefikInstances.length;

      for (const traefikInstance of this.traefikInstances) {
        const taskAction = "list-routes";
        const eventId = this.getUuid();

        // register to task events

        this.registerListener(
          this.readListeners,
          `${taskAction}-aborted-${eventId}`,
          this.listRoutesAborted
        );

        this.registerListener(
          this.readListeners,
          `${taskAction}-completed-${eventId}`,
          this.listRoutesCompleted
        );

        const res = await to(
          this.createModuleTaskForApp(traefikInstance.id, {
            action: taskAction,
            data: {
              expand_list: true,
            },
            extra: {
              title: this.$t("action." + taskAction),
              isNotificationHidden: true,
              traefikInstance: traefikInstance,
              eventId,
            },
          })
        );
        // clearListeners cannot stop a batch suspended here: without this the
        // remaining iterations would decrement the counter of the newer one
        if (generation !== this.readGeneration) {
          return;
        }
        const err = res[0];

        if (err) {
          console.error(`error creating task ${taskAction}`, err);
          const errMessage = this.getErrorMessage(err);
          this.error.listRoutes = errMessage;
          this.currentErrorAction = this.$t("action." + taskAction);
          this.currentErrorDescription = errMessage;
          this.loading.listRoutesNum--;
        }
      }
    },
    listRoutesAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.listRoutes = this.$t("error.generic_error");
      this.currentErrorAction = this.$t("action." + taskContext.action);
      this.currentErrorDescription = this.$t("error.generic_error");
      this.loading.listRoutesNum--;
    },
    listRoutesCompleted(taskContext, taskResult) {
      const routes = [];
      const traefikId = taskContext.extra.traefikInstance.id;
      const nodeId = taskContext.extra.traefikInstance.node;
      const nodeUiName = taskContext.extra.traefikInstance.node_ui_name;

      for (let route of taskResult.output) {
        route.name = route.instance;

        let type = "";
        if (route.host && route.path) {
          type = "host_and_path";
        } else if (route.host) {
          type = "host";
        } else if (route.path) {
          type = "path";
        }
        route.type = type;
        route.ip_allowlist = route.ip_allowlist ? route.ip_allowlist : [];
        route.ip_allowlist_str = route.ip_allowlist
          ? route.ip_allowlist.join("\n")
          : "";

        const node = { id: nodeId, ui_name: nodeUiName };
        const nodeLabel = this.getShortNodeLabel(node) + ` (${traefikId})`;
        route.node = nodeLabel;
        route.nodeId = nodeId;
        route.longNodeLabel = this.getNodeLabel(node);
        route.traefikInstance = traefikId;
        routes.push(route);
      }
      // drop the previous routes of this instance: a late event must not duplicate them
      this.routes = this.routes
        .filter((route) => route.traefikInstance !== traefikId)
        .concat(routes)
        .sort(this.sortByProperty("name"));
      this.loading.listRoutesNum--;

      // compute logs path for pending certificates

      const now = new Date();
      const threeDaysAgo = new Date(now.getTime() - 72 * 60 * 60 * 1000);
      const startDate = this.formatDate(threeDaysAgo, "yyyy-MM-dd");
      const startHours = this.formatDate(threeDaysAgo, "HH");
      const startMins = this.formatDate(threeDaysAgo, "mm");
      const startTime = `${startHours}%3A${startMins}`;
      const endDate = this.formatDate(now, "yyyy-MM-dd");
      const endHours = this.formatDate(now, "HH");
      const endMins = this.formatDate(now, "mm");
      const endTime = `${endHours}%3A${endMins}`;
      const maxLines = 10;

      this.pendingCertificatesLogsPath[
        traefikId
      ] = `?searchQuery=acmeCA%3D&context=module&selectedAppId=${traefikId}&followLogs=false&startDate=${startDate}&startTime=${startTime}&endDate=${endDate}&endTime=${endTime}&maxLines=${maxLines}&autoStartSearch=true`;
    },
    clearFilters() {
      this.filter.text = "";
      // "all", not "": NsComboBox only repaints its text on a matching option
      this.internalSelectedNodeId = "all";
    },
    onReloadRoutes() {
      this.listRoutes();
    },
    async deleteRoute() {
      this.loading.deleteRoute = true;
      this.error.deleteRoute = "";
      const taskAction = "delete-route";
      const eventId = this.getUuid();

      // register to task error
      this.registerListener(
        this.deleteListeners,
        `${taskAction}-aborted-${eventId}`,
        this.deleteRouteAborted
      );

      // register to task completion
      this.registerListener(
        this.deleteListeners,
        `${taskAction}-completed-${eventId}`,
        this.deleteRouteCompleted
      );

      const res = await to(
        this.createModuleTaskForApp(this.routeToDelete.traefikInstance, {
          action: taskAction,
          data: {
            instance: this.routeToDelete.name,
            lets_encrypt_cleanup: true,
          },
          extra: {
            title: this.$t("settings_http_routes.delete_route_route", {
              route: this.routeToDelete.name,
            }),
            description: this.$t("common.processing"),
            eventId,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.deleteRoute = this.getErrorMessage(err);
        this.loading.deleteRoute = false;
        return;
      }
      this.hideDeleteRouteModal();
    },
    deleteRouteAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.loading.deleteRoute = false;
    },
    deleteRouteCompleted() {
      this.loading.deleteRoute = false;

      // reload routes
      this.listRoutes();
    },
    sortHttpRoutes(sortProperty) {
      if (sortProperty != "route") {
        // default sort
        return this.sortByProperty(sortProperty);
      }

      // sort by route (host or path)

      return function (a, b) {
        const valueA = a.host || a.path;
        const valueB = b.host || b.path;

        if (valueA < valueB) {
          return -1;
        }
        if (valueA > valueB) {
          return 1;
        }
        return 0;
      };
    },
    goToPendingCertificateLogs(route) {
      this.$router.push(
        `/system-logs${this.pendingCertificatesLogsPath[route.traefikInstance]}`
      );
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

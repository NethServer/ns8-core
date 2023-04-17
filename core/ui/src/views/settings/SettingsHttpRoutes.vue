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
          >
            <cv-breadcrumb-item>
              <cv-link to="/settings">{{ $t("settings.title") }}</cv-link>
            </cv-breadcrumb-item>
            <cv-breadcrumb-item>
              <span>{{ $t("settings_http_routes.title") }}</span>
            </cv-breadcrumb-item>
          </cv-breadcrumb>
        </cv-column>
      </cv-row>
      <cv-row>
        <cv-column class="subpage-title">
          <h3>{{ $t("settings_http_routes.title") }}</h3>
        </cv-column>
      </cv-row>
      <cv-row>
        <cv-column>
          <NsInlineNotification
            kind="warning"
            :title="$t('common.use_landscape_mode')"
            :description="$t('common.use_landscape_mode_description')"
            class="landscape-warning"
          />
        </cv-column>
      </cv-row>
      <cv-row>
        <cv-column>
          <NsInlineNotification
            v-if="
              q.selectedNodeId &&
              q.selectedNodeId !== 'all' &&
              selectedNodeLabel
            "
            kind="info"
            :title="$t('settings_http_routes.routes_filtered')"
            :description="
              $t('settings_http_routes.routes_filtered_description', {
                node: selectedNodeLabel,
              })
            "
            :actionLabel="$t('settings_http_routes.show_all')"
            @action="showAllNodes"
            :showCloseButton="false"
          />
        </cv-column>
      </cv-row>
      <cv-row>
        <cv-column>
          <cv-tile light>
            <cv-grid class="no-padding">
              <cv-row>
                <cv-column :md="4">
                  <cv-combo-box
                    v-model="q.selectedNodeId"
                    :label="$t('common.choose')"
                    :title="$t('common.show')"
                    :auto-filter="true"
                    :auto-highlight="true"
                    :options="nodesForFilter"
                    :disabled="loading.listInstalledModules"
                    class="mg-bottom-xlg"
                  >
                  </cv-combo-box>
                </cv-column>
              </cv-row>
              <cv-row class="toolbar">
                <cv-column>
                  <NsButton
                    kind="secondary"
                    :icon="Add20"
                    @click="showCreateRouteModal"
                    >{{ $t("settings_http_routes.create_route") }}
                  </NsButton>
                </cv-column>
              </cv-row>
              <cv-row>
                <cv-column>
                  <NsDataTable
                    :allRows="filteredRoutes"
                    :columns="i18nTableColumns"
                    :rawColumns="tableColumns"
                    :sortable="true"
                    :pageSizes="[10, 25, 50, 100]"
                    :overflow-menu="true"
                    isSearchable
                    :searchPlaceholder="$t('settings_http_routes.search_route')"
                    :searchClearLabel="$t('common.clear_search')"
                    :noSearchResultsLabel="$t('common.no_search_results')"
                    :noSearchResultsDescription="
                      $t('common.no_search_results_description')
                    "
                    :isLoading="loadingRoutes"
                    :skeletonRows="5"
                    :isErrorShown="
                      !!error.listInstalledModules || !!error.listRoutes
                    "
                    :errorTitle="currentErrorAction"
                    :errorDescription="currentErrorDescription"
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
                      <NsEmptyState
                        :title="$t('settings_http_routes.no_http_route')"
                      >
                        <template #description>
                          <div>
                            {{
                              $t(
                                "settings_http_routes.no_http_route_description"
                              )
                            }}
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
                              :disabled="!row.user_created"
                              :data-test-id="row.name + '-edit'"
                            >
                              <NsMenuItem
                                :icon="Edit20"
                                :label="$t('common.edit')"
                              />
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
        </cv-column>
      </cv-row>
    </cv-grid>
    <HttpRouteDetailModal
      :isShown="isShownRouteDetailModal"
      :route="currentRoute"
      @hide="hideRouteDetailModal"
    />
    <CreateOrEditHttpRouteModal
      :isShown="isShownCreateOrEditRouteModal"
      :nodes="internalNodes"
      :defaultNodeId="q.selectedNodeId"
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
      @hide="hideDeleteRouteModal"
      @confirmDelete="deleteRoute"
      data-test-id="delete-route-modal"
    />
  </div>
</template>

<script>
import to from "await-to-js";
import {
  QueryParamService,
  UtilService,
  TaskService,
  IconService,
  PageTitleService,
} from "@nethserver/ns8-ui-lib";
import { mapState } from "vuex";
import HttpRouteDetailModal from "@/components/settings/HttpRouteDetailModal.vue";
import CreateOrEditHttpRouteModal from "@/components/settings/CreateOrEditHttpRouteModal.vue";
import _cloneDeep from "lodash/cloneDeep";

export default {
  name: "SettingsHttpRoutes",
  components: { HttpRouteDetailModal, CreateOrEditHttpRouteModal },
  mixins: [
    TaskService,
    UtilService,
    IconService,
    QueryParamService,
    PageTitleService,
  ],
  pageTitle() {
    return this.$t("settings_http_routes.title");
  },
  data() {
    return {
      q: {
        selectedNodeId: "",
      },
      tablePage: [],
      tableColumns: ["name", "type", "node"],
      routes: [],
      internalNodes: [],
      isShownCreateOrEditRouteModal: false,
      isShownRouteDetailModal: false,
      isShownDeleteRouteModal: false,
      currentErrorAction: "",
      currentErrorDescription: "",
      traefikInstances: [],
      currentRoute: null,
      routeToDelete: null,
      isEditingRoute: false,
      loading: {
        listInstalledModules: false,
        listRoutesNum: 0,
        deleteRoute: false,
      },
      error: {
        listInstalledModules: "",
        listRoutes: "",
        deleteRoute: "",
      },
    };
  },
  computed: {
    ...mapState(["clusterNodes"]),
    i18nTableColumns() {
      return this.tableColumns.map((column) => {
        return this.$t("settings_http_routes." + column);
      });
    },
    loadingRoutes() {
      return (
        this.loading.listInstalledModules || this.loading.listRoutesNum > 0
      );
    },
    filteredRoutes() {
      if (this.q.selectedNodeId === "all" || !this.q.selectedNodeId) {
        return this.routes;
      }

      return this.routes.filter((route) => {
        return route.nodeId === this.q.selectedNodeId;
      });
    },
    selectedNodeLabel() {
      if (this.q.selectedNodeId === "all" || !this.q.selectedNodeId) {
        return "";
      }

      const node = this.internalNodes.find(
        (node) => node.value === this.q.selectedNodeId
      );

      if (node) {
        return node.label;
      }
      return "";
    },
    nodesForFilter() {
      if (!this.internalNodes.length) {
        return [];
      }

      // add "All nodes" at the beginning of internalNodes array
      const nodes = _cloneDeep(this.internalNodes);

      nodes.unshift({
        name: "all",
        label: this.$t("common.all_nodes"),
        value: "all",
      });
      return nodes;
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
    this.listInstalledModules();
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
      this.routeToDelete = route;
      this.isShownDeleteRouteModal = true;
    },
    hideDeleteRouteModal() {
      this.isShownDeleteRouteModal = false;
    },
    async listInstalledModules() {
      this.loading.listInstalledModules = true;
      const taskAction = "list-installed-modules";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.listInstalledModulesAborted
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.listInstalledModulesCompleted
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
        const errMessage = this.getErrorMessage(err);
        this.error.listInstalledModules = errMessage;
        this.currentErrorAction = this.$t("action." + taskAction);
        this.currentErrorDescription = errMessage;
        return;
      }
    },
    listInstalledModulesAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.listInstalledModules = this.$t("error.generic_error");
      this.currentErrorAction = this.$t("action." + taskContext.action);
      this.currentErrorDescription = this.$t("error.generic_error");
      this.loading.listInstalledModules = false;
    },
    listInstalledModulesCompleted(taskContext, taskResult) {
      // init nodes
      let nodes = [];

      for (let node of this.clusterNodes) {
        nodes.push({
          name: node.id.toString(),
          label: this.getShortNodeLabel(node),
          value: node.id.toString(),
        });
      }

      let traefikInstances = [];

      for (let instanceList of Object.values(taskResult.output)) {
        for (let instance of instanceList) {
          if (instance.id.startsWith("traefik")) {
            traefikInstances.push(instance);

            // update nodes labels
            const node = nodes.find((node) => node.value === instance.node);

            if (node) {
              node.label += ` (${instance.id})`;
              node.traefikInstance = instance.id;
            }
          }
        }
      }
      this.internalNodes = nodes;
      this.traefikInstances = traefikInstances;
      this.loading.listInstalledModules = false;

      this.$nextTick(() => {
        if (!this.q.selectedNodeId) {
          // initially show all nodes
          this.q.selectedNodeId = "all";
        } else {
          const nodeId = this.q.selectedNodeId;

          // workaround to update combo box
          this.q.selectedNodeId = "";
          this.$nextTick(() => {
            this.q.selectedNodeId = nodeId;
          });
        }
      });

      this.listRoutes();
    },
    async listRoutes() {
      this.routes = [];

      for (const traefikInstance of this.traefikInstances) {
        const taskAction = "list-routes";
        const eventId = this.getUuid();
        this.loading.listRoutesNum++;

        // register to task events

        this.$root.$once(
          `${taskAction}-aborted-${eventId}`,
          this.listRoutesAborted
        );

        this.$root.$once(
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

        const traefikId = taskContext.extra.traefikInstance.id;
        const nodeId = taskContext.extra.traefikInstance.node;
        const nodeUiName = taskContext.extra.traefikInstance.node_ui_name;
        const node = { id: nodeId, ui_name: nodeUiName };
        const nodeLabel = this.getShortNodeLabel(node) + ` (${traefikId})`;
        route.node = nodeLabel;
        route.nodeId = nodeId;
        route.longNodeLabel = this.getNodeLabel(node);
        route.traefikInstance = traefikId;
        routes.push(route);
      }
      this.routes = this.routes
        .concat(routes)
        .sort(this.sortByProperty("name"));
      this.loading.listRoutesNum--;
    },
    showAllNodes() {
      this.q.selectedNodeId = "all";
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
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.deleteRouteAborted
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.deleteRouteCompleted
      );

      const res = await to(
        this.createModuleTaskForApp(this.routeToDelete.traefikInstance, {
          action: taskAction,
          data: {
            instance: this.routeToDelete.name,
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
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";
</style>

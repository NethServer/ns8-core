<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <div>
    <cv-grid fullWidth>
      <cv-row>
        <cv-column :md="4" :xlg="10">
          <cv-breadcrumb
            aria-label="breadcrumb"
            :no-trailing-slash="true"
            class="breadcrumb"
          >
            <cv-breadcrumb-item>
              <cv-link to="/software-center">{{
                $t("software_center.title")
              }}</cv-link>
            </cv-breadcrumb-item>
            <cv-breadcrumb-item>
              <span>{{ $t("software_center.core_apps") }}</span>
            </cv-breadcrumb-item>
          </cv-breadcrumb>
        </cv-column>
      </cv-row>
      <cv-row>
        <cv-column class="subpage-title">
          <h2>{{ $t("software_center.core_apps") }}</h2>
        </cv-column>
      </cv-row>
      <cv-row>
        <cv-column>
          <NsInlineNotification
            v-if="isCoreUpdatable"
            kind="warning"
            :title="$t('software_center.core_app_update_available')"
            :description="
              $t('software_center.core_app_update_available_description')
            "
            :showCloseButton="false"
          />
        </cv-column>
      </cv-row>
      <cv-tile light>
        <cv-row>
          <cv-column>
            <NsButton
              class="mg-bottom-md"
              v-if="isCoreUpdatable"
              kind="primary"
              :icon="Upgrade20"
              :disabled="updateCoreTimeout || isUpdateInProgress"
              :loading="isUpdateInProgress"
              @click="onWillUpdateCore"
              >{{ $t("software_center.update_core") }}
            </NsButton>
            <div class="mg-bottom-md">
              {{
                $t("software_center.core_app_table_description", {
                  productName: $root.config.PRODUCT_NAME,
                })
              }}
            </div>

            <NsDataTable
              :allRows="coreInstances"
              :columns="i18nTableColumns"
              :rawColumns="tableColumns"
              :sortable="true"
              :overflow-menu="false"
              :isLoading="loading.listCoreModules"
              :skeletonRows="5"
              isSearchable
              :searchPlaceholder="$t('software_center.search_core_apps')"
              :searchClearLabel="$t('common.clear_search')"
              :noSearchResultsLabel="$t('common.no_search_results')"
              :noSearchResultsDescription="
                $t('common.no_search_results_description')
              "
              :itemsPerPageLabel="$t('pagination.items_per_page')"
              :rangeOfTotalItemsLabel="$t('pagination.range_of_total_items')"
              :ofTotalPagesLabel="$t('pagination.of_total_pages')"
              :backwardText="$t('pagination.previous_page')"
              :forwardText="$t('pagination.next_page')"
              :pageNumberLabel="$t('pagination.page_number')"
              @updatePage="tablePage = $event"
            >
              <template slot="data">
                <cv-data-table-row
                  v-for="(row, rowIndex) in tablePage"
                  :key="`${rowIndex}`"
                  :value="`${rowIndex}`"
                >
                  <cv-data-table-cell>
                    <a @click="showAppInfo(row)">{{
                      row.id
                    }}</a></cv-data-table-cell
                  >
                  <cv-data-table-cell>{{ row.node_id }}</cv-data-table-cell>
                  <cv-data-table-cell>{{ row.version }}</cv-data-table-cell>
                  <cv-data-table-cell v-if="isCoreUpdatable">
                    {{ row.update || "-" }}
                    <cv-link
                      v-if="row.docs.relnotes_url && row.update"
                      :href="row.docs.relnotes_url"
                      target="_blank"
                    >
                      <NsButton kind="ghost" :icon="Launch20">
                        {{ $t("common.release_notes") }}
                      </NsButton></cv-link
                    >
                  </cv-data-table-cell>
                </cv-data-table-row>
              </template>
            </NsDataTable>
          </cv-column>
        </cv-row>
      </cv-tile>
    </cv-grid>
    <AppInfoModal
      :app="appInfo.app"
      :isShown="appInfo.isShown"
      @close="onClose"
    />
  </div>
</template>

<script>
import {
  QueryParamService,
  UtilService,
  TaskService,
  IconService,
  LottieService,
  PageTitleService,
} from "@nethserver/ns8-ui-lib";
import { mapState, mapActions } from "vuex";
import to from "await-to-js";
import AppInfoModal from "@/components/software-center/AppInfoModal";

export default {
  name: "SoftwareCenterCoreApps",
  components: { AppInfoModal },
  mixins: [
    IconService,
    QueryParamService,
    UtilService,
    TaskService,
    LottieService,
    PageTitleService,
  ],
  pageTitle() {
    return this.$t("settings_sw_repositories.title");
  },
  data() {
    return {
      tablePage: [],
      tableColumns: [],
      updateCoreTimeout: 0,
      coreApp: null,
      appInfo: {
        isShown: false,
        app: null,
      },
      loading: {
        cleanRepositoriesCache: false,
        listCoreModules: true,
      },
      error: {
        cleanRepositoriesCache: "",
        listCoreModules: "",
        updateCore: "",
        updateModules: "",
      },
    };
  },
  computed: {
    ...mapState(["isUpdateInProgress", "clusterNodes"]),
    nodeIds() {
      return this.clusterNodes.map((node) => node.id);
    },
    isCoreUpdatable() {
      // check if at least a core module has an update available
      if (this.coreApp) {
        for (const coreApp of this.coreApp.installed) {
          for (const instance of coreApp.instances) {
            if (instance.update) {
              return true;
            }
          }
        }
      }
      return false;
    },
    coreInstances() {
      if (this.coreApp) {
        const coreInstances = [];

        for (const coreApp of this.coreApp.installed) {
          coreInstances.push(...coreApp.instances);
        }
        coreInstances.sort(this.sortByProperty("id"));
        return coreInstances;
      }
      return [];
    },
    i18nTableColumns() {
      return this.tableColumns.map((column) => {
        return this.$t("software_center." + column);
      });
    },
  },
  watch: {
    coreApp: function () {
      this.updateTableData();
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
    this.listCoreModules();
    this.updateTableData();
  },
  methods: {
    ...mapActions(["setUpdateInProgressInStore"]),
    onClose() {
      const context = this;

      // needed to reset modal scroll to top
      setTimeout(() => {
        context.appInfo.isShown = false;
      }, 250);
    },
    showAppInfo(app) {
      this.appInfo.isShown = true;
      this.appInfo.app = app;
    },
    updateTableData() {
      this.tableColumns = this.isCoreUpdatable
        ? ["id", "node_id", "version", "update"]
        : ["id", "node_id", "version"];
    },
    async listCoreModules() {
      this.error.listCoreModules = "";
      this.loading.listCoreModules = true;
      const taskAction = "list-core-modules";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.listCoreModulesAborted
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.listCoreModulesCompleted
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
        this.error.listCoreModules = this.getErrorMessage(err);
        this.loading.listCoreModules = false;
        return;
      }
    },
    listCoreModulesAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.listCoreModules = this.$t("error.generic_error");
      this.loading.listCoreModules = false;
    },
    listCoreModulesCompleted(taskContext, taskResult) {
      let coreModules = taskResult.output;
      // update node_id with this.getNodeLabel() to  implement the search by node_id and node_ui_name
      coreModules = coreModules.map((coreApp) => {
        coreApp.instances = coreApp.instances.map((instance) => {
          instance.node_id = this.getNodeLabel({
            id: instance.node_id,
            ui_name: instance.node_ui_name,
          });
          return instance;
        });
        return coreApp;
      });

      let isCoreUpdateAvailable = false;

      this.coreApp = {
        id: "core",
        name: this.$t("software_center.core_app_name", {
          productName: this.$root.config.PRODUCT_NAME,
        }),
        description: {
          en: this.$t("software_center.core_app_description", {
            productName: this.$root.config.PRODUCT_NAME,
          }),
        },
        installed: coreModules,
      };

      // check for core updates

      for (const coreApp of coreModules) {
        for (const coreInstance of coreApp.instances) {
          if (coreInstance.update) {
            isCoreUpdateAvailable = true;
            break;
          }
        }
      }

      if (isCoreUpdateAvailable) {
        this.isCoreUpdateAvailable = true;
      } else {
        this.isCoreUpdateAvailable = false;
      }
      this.coreModules = coreModules;
      this.loading.listCoreModules = false;
    },

    onWillUpdateCore() {
      const notification = {
        id: "updateCore",
        title: this.$t("software_center.core_update_will_start_in_a_moment"),
        type: "info",
        toastTimeout: this.DELETE_DELAY,
        actionLabel: this.$t("common.cancel"),
        action: {
          type: "callback",
          callback: this.cancelDeleteAddress,
        },
      };
      this.createNotificationForApp(notification);

      this.updateCoreTimeout = setTimeout(() => {
        // remove notification from drawer
        this.deleteNotificationForApp("updateCore");

        // delete timeout
        this.updateCoreTimeout = 0;

        // call api to update core
        this.updateCore();
      }, this.DELETE_DELAY);
    },
    cancelDeleteAddress() {
      clearTimeout(this.updateCoreTimeout);

      // remove notification from drawer
      this.deleteNotificationForApp("updateCore");

      // delete timeout
      this.updateCoreTimeout = 0;
    },

    async updateCore() {
      this.error.updateCore = "";
      this.setUpdateInProgressInStore(true);
      const taskAction = "update-core";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.updateCoreAborted
      );

      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.updateCoreCompleted
      );

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            nodes: this.nodeIds,
          },
          extra: {
            title: this.$t("software_center.update_core"),
            description: this.$t("common.processing"),
            eventId,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.updateCore = this.getErrorMessage(err);
        this.setUpdateInProgressInStore(false);
        return;
      }
    },
    updateCoreAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.setUpdateInProgressInStore(false);
    },
    updateCoreCompleted() {
      // add a brief delay: API server may not be ready to accept requests immediately
      setTimeout(() => {
        this.setUpdateInProgressInStore(false);
        this.listCoreModules();
      }, 5000);
    },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";
</style>

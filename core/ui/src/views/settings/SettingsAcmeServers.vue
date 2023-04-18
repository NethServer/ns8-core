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
              <cv-link to="/settings/tls-certificates">{{
                $t("settings_tls_certificates.title")
              }}</cv-link>
            </cv-breadcrumb-item>
            <cv-breadcrumb-item>
              <span>{{ $t("settings_acme_servers.title") }}</span>
            </cv-breadcrumb-item>
          </cv-breadcrumb>
        </cv-column>
      </cv-row>
      <cv-row>
        <cv-column :md="4" :xlg="10" class="subpage-title">
          <h3>{{ $t("settings_acme_servers.title") }}</h3>
        </cv-column>
        <cv-column :md="4" :xlg="6">
          <div class="page-toolbar">
            <NsButton
              kind="tertiary"
              size="field"
              :icon="Certificate20"
              @click="goToTlsCertificates()"
              class="subpage-toolbar-item"
              >{{ $t("settings_tls_certificates.title") }}
            </NsButton>
          </div>
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
          <cv-tile light>
            <cv-grid class="no-padding">
              <cv-row>
                <cv-column>
                  <NsDataTable
                    :allRows="servers"
                    :columns="i18nTableColumns"
                    :rawColumns="tableColumns"
                    :sortable="true"
                    :pageSizes="[10, 25, 50, 100]"
                    :overflow-menu="true"
                    isSearchable
                    :searchPlaceholder="
                      $t('settings_acme_servers.search_acme_server')
                    "
                    :searchClearLabel="$t('common.clear_search')"
                    :noSearchResultsLabel="$t('common.no_search_results')"
                    :noSearchResultsDescription="
                      $t('common.no_search_results_description')
                    "
                    :isLoading="loadingServers"
                    :skeletonRows="5"
                    :isErrorShown="
                      !!error.listInstalledModules || !!error.getAcmeServer
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
                        :title="$t('settings_acme_servers.no_acme_server')"
                      >
                        <template #description>
                          <div>
                            {{
                              $t(
                                "settings_acme_servers.no_acme_server_description"
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
                          <span>
                            {{ row.url }}
                          </span>
                        </cv-data-table-cell>
                        <cv-data-table-cell>
                          <span>{{ row.node }}</span>
                        </cv-data-table-cell>
                        <cv-data-table-cell>
                          <div class="justify-flex-end">
                            <NsButton
                              kind="secondary"
                              :icon="Edit20"
                              size="small"
                              @click="showEditServerModal(row)"
                              :data-test-id="row.url + '-delete'"
                              >{{ $t("common.edit") }}
                            </NsButton>
                          </div>
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
    <EditAcmeServerModal
      :isShown="isShownEditServerModal"
      :server="currentServer"
      @hide="hideEditServerModal"
      @reloadServers="onReloadServers"
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
import EditAcmeServerModal from "@/components/settings/EditAcmeServerModal.vue";

export default {
  name: "SettingsAcmeServers",
  components: { EditAcmeServerModal },
  mixins: [
    TaskService,
    UtilService,
    IconService,
    QueryParamService,
    PageTitleService,
  ],
  pageTitle() {
    return this.$t("settings_acme_servers.title");
  },
  data() {
    return {
      tablePage: [],
      tableColumns: ["url", "node"],
      servers: [],
      internalNodes: [],
      isShownEditServerModal: false,
      currentErrorAction: "",
      currentErrorDescription: "",
      traefikInstances: [],
      currentServer: null,
      loading: {
        listInstalledModules: false,
        getAcmeServerNum: 0,
      },
      error: {
        listInstalledModules: "",
        getAcmeServer: "",
      },
    };
  },
  computed: {
    ...mapState(["clusterNodes"]),
    i18nTableColumns() {
      return this.tableColumns.map((column) => {
        return this.$t("settings_acme_servers." + column);
      });
    },
    loadingServers() {
      return (
        this.loading.listInstalledModules || this.loading.getAcmeServerNum > 0
      );
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
    showEditServerModal(server) {
      this.currentServer = server;
      this.isShownEditServerModal = true;
    },
    hideEditServerModal() {
      this.isShownEditServerModal = false;
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
      this.getAcmeServer();
    },
    async getAcmeServer() {
      this.servers = [];

      for (const traefikInstance of this.traefikInstances) {
        const taskAction = "get-acme-server";
        const eventId = this.getUuid();
        this.loading.getAcmeServerNum++;

        // register to task events

        this.$root.$once(
          `${taskAction}-aborted-${eventId}`,
          this.getAcmeServerAborted
        );

        this.$root.$once(
          `${taskAction}-completed-${eventId}`,
          this.getAcmeServerCompleted
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
          this.error.getAcmeServer = errMessage;
          this.currentErrorAction = this.$t("action." + taskAction);
          this.currentErrorDescription = errMessage;
          this.loading.getAcmeServerNum--;
        }
      }
    },
    getAcmeServerAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.getAcmeServer = this.$t("error.generic_error");
      this.currentErrorAction = this.$t("action." + taskContext.action);
      this.currentErrorDescription = this.$t("error.generic_error");
      this.loading.getAcmeServerNum--;
    },
    getAcmeServerCompleted(taskContext, taskResult) {
      const server = taskResult.output;
      const traefikId = taskContext.extra.traefikInstance.id;
      const nodeId = taskContext.extra.traefikInstance.node;
      const nodeUiName = taskContext.extra.traefikInstance.node_ui_name;
      const node = { id: nodeId, ui_name: nodeUiName };
      const nodeLabel = this.getShortNodeLabel(node) + ` (${traefikId})`;
      server.node = nodeLabel;
      server.nodeId = nodeId;
      server.longNodeLabel = this.getNodeLabel(node);
      server.traefikInstance = traefikId;
      this.servers.push(server);
      this.servers.sort(this.sortByProperty("url"));
      this.loading.getAcmeServerNum--;
    },
    onReloadServers() {
      this.getAcmeServer();
    },
    goToTlsCertificates() {
      this.$router.push("/settings/tls-certificates");
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";

.justify-flex-end {
  display: flex;
  justify-content: flex-end;
}
</style>

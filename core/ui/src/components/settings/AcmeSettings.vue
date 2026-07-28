<!--
  Copyright (C) 2026 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <div>
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
              :rangeOfTotalItemsLabel="$t('pagination.range_of_total_items')"
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
                        $t("settings_acme_servers.no_acme_server_description")
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
                    <span>{{ row.node }}</span>
                  </cv-data-table-cell>
                  <cv-data-table-cell>
                    <span>
                      {{ row.url }}
                    </span>
                  </cv-data-table-cell>
                  <cv-data-table-cell>
                    <NsTag
                      v-if="row.challenge"
                      :kind="getChallengeTagKind(row.challenge)"
                      size="sm"
                      :label="getChallengeLabel(row.challenge)"
                    />
                  </cv-data-table-cell>
                  <cv-data-table-cell>
                    <div class="justify-flex-end">
                      <NsButton
                        kind="ghost"
                        :icon="Edit20"
                        size="small"
                        @click="showEditServerModal(row)"
                        :data-test-id="row.nodeId + '-edit'"
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
import { UtilService, TaskService, IconService } from "@nethserver/ns8-ui-lib";
import EditAcmeServerModal, {
  ACME_CHALLENGE_TYPES,
} from "@/components/settings/EditAcmeServerModal.vue";

export default {
  name: "AcmeSettings",
  components: { EditAcmeServerModal },
  mixins: [TaskService, UtilService, IconService],
  data() {
    return {
      tablePage: [],
      tableColumns: ["node", "url", "challenge"],
      servers: [],
      isShownEditServerModal: false,
      currentErrorAction: "",
      currentErrorDescription: "",
      traefikInstances: [],
      currentServer: null,
      // [eventName, handler] pairs registered on $root, removed in beforeDestroy
      taskListeners: [],
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
  created() {
    this.listInstalledModules();
  },
  beforeDestroy() {
    // remove task listeners still registered on $root: a task completing after
    // the component is gone would run its handler on a dead instance, and
    // listInstalledModulesCompleted() would create one task per traefik node.
    // $off only touches $root._events, the notification drawer is store-driven
    this.taskListeners.forEach(([eventName, handler]) => {
      this.$root.$off(eventName, handler);
    });
    this.taskListeners = [];
  },
  methods: {
    registerTaskListener(eventName, handler) {
      this.$root.$once(eventName, handler);
      this.taskListeners.push([eventName, handler]);
    },
    getChallengeType(challenge) {
      return ACME_CHALLENGE_TYPES.find((type) => type.value === challenge);
    },
    getChallengeTagKind(challenge) {
      const challengeType = this.getChallengeType(challenge);
      // fallback needed: cv-tag validates the kind against a fixed list
      return challengeType ? challengeType.tagKind : "gray";
    },
    getChallengeLabel(challenge) {
      const challengeType = this.getChallengeType(challenge);
      return challengeType ? this.$t(challengeType.labelKey) : challenge;
    },
    showEditServerModal(server) {
      this.currentServer = server;
      this.isShownEditServerModal = true;
    },
    hideEditServerModal() {
      this.isShownEditServerModal = false;
    },
    clearTableError() {
      this.error.listInstalledModules = "";
      this.error.getAcmeServer = "";
      this.currentErrorAction = "";
      this.currentErrorDescription = "";
    },
    async listInstalledModules() {
      this.loading.listInstalledModules = true;
      this.clearTableError();
      const taskAction = "list-installed-modules";
      const eventId = this.getUuid();

      // register to task error
      this.registerTaskListener(
        `${taskAction}-aborted-${eventId}`,
        this.listInstalledModulesAborted
      );

      // register to task completion
      this.registerTaskListener(
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
        this.loading.listInstalledModules = false;
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
      let traefikInstances = [];

      for (let instanceList of Object.values(taskResult.output)) {
        for (let instance of instanceList) {
          if (instance.id.startsWith("traefik")) {
            traefikInstances.push(instance);
          }
        }
      }
      this.traefikInstances = traefikInstances;
      this.loading.listInstalledModules = false;
      this.getAcmeServer();
    },
    async getAcmeServer() {
      this.servers = [];
      // without this a transient failure keeps the table on the error state
      // even once every node answers again
      this.clearTableError();

      for (const traefikInstance of this.traefikInstances) {
        const taskAction = "get-acme-server";
        const eventId = this.getUuid();
        this.loading.getAcmeServerNum++;

        // register to task events

        this.registerTaskListener(
          `${taskAction}-aborted-${eventId}`,
          this.getAcmeServerAborted
        );

        this.registerTaskListener(
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
      this.servers.sort(this.sortByProperty("node"));
      this.loading.getAcmeServerNum--;
    },
    onReloadServers() {
      this.getAcmeServer();
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

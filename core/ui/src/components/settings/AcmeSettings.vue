<!--
  Copyright (C) 2026 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <div>
    <cv-tile light>
      <!-- no nested cv-grid: bx--grid caps at 99rem and centers, which leaves
           empty gutters on screens wider than 1584px -->
      <!-- NsDataTable replaces the rows with its error: report a partial
           failure here, so the nodes that did answer stay readable -->
      <NsInlineNotification
        v-if="tableError && servers.length"
        kind="error"
        :title="tableErrorTitle"
        :description="tableErrorDescription"
        :showCloseButton="false"
      />
      <NsDataTable
        :allRows="servers"
        :columns="i18nTableColumns"
        :rawColumns="tableColumns"
        :sortable="true"
        :pageSizes="[10, 25, 50, 100]"
        :overflow-menu="true"
        isSearchable
        :searchPlaceholder="$t('settings_acme_servers.search_acme_settings')"
        :searchClearLabel="$t('common.clear_search')"
        :noSearchResultsLabel="$t('common.no_search_results')"
        :noSearchResultsDescription="$t('common.no_search_results_description')"
        :isLoading="loadingServers"
        :skeletonRows="5"
        :isErrorShown="!!tableError && !servers.length"
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
          <NsEmptyState :title="$t('settings_acme_servers.no_acme_settings')">
            <template #description>
              <div>
                {{ $t("settings_acme_servers.no_acme_settings_description") }}
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
  props: {
    // the parent page already lists them: do not run the task twice
    traefikInstances: {
      type: Array,
      default: () => [],
    },
    isLoadingInstances: Boolean,
    instancesError: {
      type: String,
      default: "",
    },
  },
  data() {
    return {
      tablePage: [],
      tableColumns: ["node", "url", "challenge"],
      // the raw columns are row property names: labels need their own keys
      tableColumnLabelKeys: [
        "settings_acme_servers.node",
        "settings_acme_servers.acme_directory_url",
        "settings_acme_servers.challenge",
      ],
      servers: [],
      isShownEditServerModal: false,
      currentErrorAction: "",
      currentErrorDescription: "",
      currentServer: null,
      // [eventName, handler] pairs registered on $root
      taskListeners: [],
      loading: {
        getAcmeServerNum: 0,
      },
      error: {
        getAcmeServer: "",
      },
    };
  },
  computed: {
    i18nTableColumns() {
      return this.tableColumnLabelKeys.map((key) => this.$t(key));
    },
    loadingServers() {
      return this.isLoadingInstances || this.loading.getAcmeServerNum > 0;
    },
    tableError() {
      return this.instancesError || this.error.getAcmeServer;
    },
    tableErrorTitle() {
      return this.instancesError
        ? this.$t("action.list-installed-modules")
        : this.currentErrorAction;
    },
    tableErrorDescription() {
      return this.instancesError || this.currentErrorDescription;
    },
  },
  watch: {
    // the parent republishes the instances after a Traefik restart killed the
    // websocket, which is also how this table recovers
    traefikInstances: function () {
      this.getAcmeServer();
    },
  },
  created() {
    if (this.traefikInstances.length) {
      this.getAcmeServer();
    }
  },
  beforeDestroy() {
    // a task completing after destroy would re-fire the whole task chain
    this.clearTaskListeners();
  },
  methods: {
    registerTaskListener(eventName, handler) {
      this.$root.$once(eventName, handler);
      this.taskListeners.push([eventName, handler]);
    },
    clearTaskListeners() {
      this.taskListeners.forEach(([eventName, handler]) => {
        this.$root.$off(eventName, handler);
      });
      this.taskListeners = [];
    },
    getChallengeType(challenge) {
      return ACME_CHALLENGE_TYPES.find((type) => type.value === challenge);
    },
    getChallengeTagKind(challenge) {
      const challengeType = this.getChallengeType(challenge);
      // cv-tag validates the kind against a fixed list
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
      this.error.getAcmeServer = "";
      this.currentErrorAction = "";
      this.currentErrorDescription = "";
    },
    async getAcmeServer() {
      // a handler of a previous round would decrement the counter of this one
      this.clearTaskListeners();
      this.servers = [];
      // otherwise a transient error sticks after the nodes answer again
      this.clearTableError();
      // count the whole batch upfront: the counter must not reach zero between
      // two iterations
      this.loading.getAcmeServerNum = this.traefikInstances.length;

      for (const traefikInstance of this.traefikInstances) {
        const taskAction = "get-acme-server";
        const eventId = this.getUuid();

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

      // replace instead of append: an event delivered after a reconnection
      // reset would otherwise duplicate the row
      const index = this.servers.findIndex(
        (existing) => existing.traefikInstance === traefikId
      );

      if (index > -1) {
        this.servers.splice(index, 1, server);
      } else {
        this.servers.push(server);
      }
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

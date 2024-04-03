<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <cv-grid fullWidth class="system-logs">
    <cv-row>
      <cv-column :md="4" :xlg="10" class="page-title">
        <h2>{{ $t("system_logs.title") }}</h2>
      </cv-column>
      <cv-column :md="4" :xlg="6">
        <div class="page-toolbar">
          <NsButton
            kind="secondary"
            size="field"
            :icon="Add20"
            @click="addSearch()"
            :disabled="searches.length > 1"
            class="page-toolbar-item"
            >{{ $t("system_logs.add_search") }}</NsButton
          >
          <!-- <template v-else>
            <cv-content-switcher
              @selected="onLayoutSelected"
              class="page-toolbar-item"
            >
              <cv-content-switcher-button
                owner-id="horizontal"
                :selected="csbHorizontalLayoutSelected"
                >{{
                  $t("system_logs.horizontal_layout")
                }}</cv-content-switcher-button
              >
              <cv-content-switcher-button
                owner-id="vertical"
                :selected="csbVerticalLayoutSelected"
                >{{
                  $t("system_logs.vertical_layout")
                }}</cv-content-switcher-button
              >
            </cv-content-switcher> -->
          <NsIconMenu
            :flipMenu="true"
            tipPosition="top"
            tipAlignment="end"
            class="page-toolbar-item"
          >
            <cv-overflow-menu-item @click="collapseAllFilters()">
              <NsMenuItem
                :icon="RowCollapse20"
                :label="$t('system_logs.collapse_filters')"
              />
            </cv-overflow-menu-item>
            <cv-overflow-menu-item
              @click="setHorizontalLayout()"
              :disabled="searches.length < 2 || !verticalLayout"
              class="toggle-layout"
            >
              <NsMenuItem
                :icon="Row20"
                :label="$t('system_logs.horizontal_layout')"
              />
            </cv-overflow-menu-item>
            <cv-overflow-menu-item
              @click="setVerticalLayout()"
              :disabled="searches.length < 2 || verticalLayout"
              class="toggle-layout"
            >
              <NsMenuItem
                :icon="Column20"
                :label="$t('system_logs.vertical_layout')"
              />
            </cv-overflow-menu-item>
          </NsIconMenu>
          <!-- </template> -->
        </div>
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
    <cv-row v-if="error.listInstalledModules">
      <cv-column>
        <NsInlineNotification
          kind="error"
          :title="$t('action.list-installed-modules')"
          :description="error.listInstalledModules"
          :showCloseButton="false"
        />
      </cv-column>
    </cv-row>
    <cv-row>
      <cv-column
        v-for="(searchId, index) in searches"
        :key="searchId"
        :md="verticalLayout ? 4 : 8"
      >
        <LogSearch
          :searchId="searchId"
          :nodes="internalNodes"
          :apps="apps"
          :lokiInstances="lokiInstances"
          :loadingApps="loading.listInstalledModules"
          :loadingLoki="loading.lokiInstances"
          :verticalLayout="verticalLayout"
          :numSearches="searches.length"
          :mainSearch="index == 0"
          :searchQuery="q.searchQuery"
          :context="q.context"
          :selectedNodeId="q.selectedNodeId"
          :selectedAppId="q.selectedAppId"
          :selectedLokiId="q.selectedLokiId"
          :followLogs="q.followLogs"
          :maxLines="q.maxLines"
          :startDate="q.startDate"
          :startTime="q.startTime"
          :endDate="q.endDate"
          :endTime="q.endTime"
          :startSearchCommand="startSearchCommand"
          :closeAriaLabel="$t('common.close')"
          @close="closeSearch"
          @updateSearchQuery="onUpdateSearchQuery"
          @updateContext="onUpdateContext"
          @updateSelectedNodeId="onUpdateSelectedNodeId"
          @updateSelectedAppId="onUpdateSelectedAppId"
          @updateFollowLogs="onUpdateFollowLogs"
          @updateMaxLines="onUpdateMaxLines"
          @updateStartDate="onUpdateStartDate"
          @updateStartTime="onUpdateStartTime"
          @updateEndDate="onUpdateEndDate"
          @updateEndTime="onUpdateEndTime"
          @updateTimezone="onUpdateTimezone"
          :light="true"
          :ref="'search-' + searchId"
          :timezone="q.timezone"
        />
      </cv-column>
    </cv-row>
  </cv-grid>
</template>

<script>
import to from "await-to-js";
import {
  UtilService,
  IconService,
  TaskService,
  QueryParamService,
  DateTimeService,
  PageTitleService,
} from "@nethserver/ns8-ui-lib";
import { v4 as uuidv4 } from "uuid";
import LogSearch from "@/components/system-logs/LogSearch.vue";
import { mapState } from "vuex";

export default {
  name: "SystemLogs",
  components: { LogSearch },
  mixins: [
    IconService,
    UtilService,
    TaskService,
    QueryParamService,
    DateTimeService,
    PageTitleService,
  ],
  pageTitle() {
    return this.$t("system_logs.title");
  },
  data() {
    return {
      q: {
        searchQuery: "",
        context: "cluster",
        selectedNodeId: "",
        selectedAppId: "",
        selectedLokiId: "",
        followLogs: false,
        maxLines: "500",
        startDate: "",
        startTime: "",
        endDate: "",
        endTime: "",
        autoStartSearch: false,
        timezone: "local",
      },
      internalNodes: [],
      apps: [],
      lokiInstances: [],
      searches: [],
      verticalLayout: false,
      startSearchCommand: 0,
      loading: {
        listInstalledModules: false,
        lokiInstances: false,
      },
      error: {
        listInstalledModules: "",
        lokiInstances: "",
      },
    };
  },
  computed: {
    ...mapState(["clusterNodes"]),
  },
  watch: {
    clusterNodes: function () {
      if (this.clusterNodes.length) {
        this.initNodes();
      }
    },
    apps: function () {
      if (this.apps.length && this.q.autoStartSearch) {
        this.$nextTick(() => {
          this.startSearchCommand++;
        });
      }
    },
    lokiInstances: function () {
      if (this.lokiInstances.length && this.q.autoStartSearch) {
        this.$nextTick(() => {
          this.startSearchCommand++;
        });
      }
    },
  },
  created() {
    this.initTimeFilters();
    this.initNodes();
    this.listInstalledModules();
    this.getClusterLokiInstances();
    this.addSearch();
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
  methods: {
    initTimeFilters() {
      this.q.startDate = this.formatDate(
        this.subDays(new Date(), 1),
        "yyyy-MM-dd"
      );

      this.q.endDate = this.formatDate(new Date(), "yyyy-MM-dd");
    },
    addSearch() {
      const searchId = uuidv4();
      this.searches.push(searchId);

      // scroll to new search
      this.$nextTick(() => {
        const el = this.$refs["search-" + searchId][0].$el;
        this.scrollToElement(el);
      });
    },
    initNodes() {
      let nodes = [];

      for (let node of this.clusterNodes) {
        nodes.push({
          name: node.id.toString(),
          label: node.ui_name
            ? node.ui_name + " (" + this.$t("common.node") + " " + node.id + ")"
            : this.$t("common.node") + " " + node.id,
          value: node.id.toString(),
        });
      }
      nodes.sort(this.sortByProperty("label"));
      this.internalNodes = nodes;
    },
    async listInstalledModules() {
      this.loading.listInstalledModules = true;
      const taskAction = "list-installed-modules";

      // register to task error
      this.$root.$off(taskAction + "-aborted");
      this.$root.$once(
        taskAction + "-aborted",
        this.listInstalledModulesAborted
      );

      // register to task completion
      this.$root.$once(
        taskAction + "-completed",
        this.listInstalledModulesCompleted
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
        this.error.listInstalledModules = this.getErrorMessage(err);
        return;
      }
    },
    listInstalledModulesAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.listInstalledModules = this.$t("error.generic_error");
      this.loading.listInstalledModules = false;
    },
    listInstalledModulesCompleted(taskContext, taskResult) {
      let apps = [];
      for (let instanceList of Object.values(taskResult.output)) {
        for (let instance of instanceList) {
          apps.push({
            name: instance.id,
            label:
              (instance.ui_name
                ? instance.ui_name + " (" + instance.id + ")"
                : instance.id) +
              " - " +
              (instance.node_ui_name
                ? instance.node_ui_name
                : this.$t("common.node") + " " + instance.node),
            value: instance.id,
          });
        }
      }
      apps.sort(this.sortByProperty("label"));
      this.apps = apps;
      this.loading.listInstalledModules = false;
    },
    async getClusterLokiInstances() {
      this.loading.lokiInstances = true;
      this.error.getClusterLokiInstances = "";
      const taskAction = "list-loki-instances";

      this.$root.$once(
        taskAction + "-aborted",
        this.getClusterLokiInstancesAborted
      );

      // register to task completion
      this.$root.$once(
        taskAction + "-completed",
        this.getClusterLokiInstancesCompleted
      );

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          extra: {
            isNotificationHidden: true,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.getClusterLokiInstances = this.getErrorMessage(err);
        return;
      }
    },
    getClusterLokiInstancesAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.lokiInstances = this.$t("error.generic_error");
      this.loading.lokiInstances = false;
    },
    getClusterLokiInstancesCompleted(taskContext, taskResult) {
      const lokiOutput = taskResult.output.instances.filter(
        (instance) => !instance.offline
      );
      let lokiInstances = [];
      for (let instance of lokiOutput) {
        let currentInstance = "";
        if (instance.active) {
          currentInstance =
            " (" + this.$t("system_logs.current_instance") + ")";
        }
        lokiInstances.push({
          name: instance.instance_id,
          label: instance.instance_label
            ? instance.instance_label +
              " - " +
              instance.instance_id +
              currentInstance
            : instance.instance_id + currentInstance,
          value: instance.instance_id,
        });
      }
      this.lokiInstances = lokiInstances;
      this.q.selectedLokiId = lokiOutput.find(
        (instance) => instance.active === true
      ).instance_id;
      this.loading.lokiInstances = false;
    },
    setHorizontalLayout() {
      this.verticalLayout = false;
    },
    setVerticalLayout() {
      this.verticalLayout = true;
    },
    collapseAllFilters() {
      for (const searchId of this.searches) {
        this.$root.$emit(`collapseSystemLogsFilters-${searchId}`);
      }
    },
    closeSearch(searchId) {
      this.searches = this.searches.filter((s) => s !== searchId);
      this.verticalLayout = false;
      this.$root.$off(`logsStart-${this.searchId}`);
    },
    onUpdateSearchQuery(searchQuery) {
      this.q.searchQuery = searchQuery;
    },
    onUpdateContext(context) {
      this.q.context = context;
    },
    onUpdateSelectedNodeId(selectedNodeId) {
      this.q.selectedNodeId = selectedNodeId;
    },
    onUpdateSelectedAppId(selectedAppId) {
      this.q.selectedAppId = selectedAppId;
    },
    onUpdateFollowLogs(followLogs) {
      this.q.followLogs = followLogs;
    },
    onUpdateMaxLines(maxLines) {
      this.q.maxLines = maxLines;
    },
    onUpdateStartDate(startDate) {
      this.q.startDate = startDate;
    },
    onUpdateEndDate(endDate) {
      this.q.endDate = endDate;
    },
    onUpdateStartTime(startTime) {
      this.q.startTime = startTime;
    },
    onUpdateEndTime(endTime) {
      this.q.endTime = endTime;
    },
    onUpdateTimezone(timezone) {
      this.q.timezone = timezone;
    },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";

// cannot switch to vertical layout on small screens
@media (max-width: $breakpoint-x-large) {
  .toggle-layout {
    display: none;
  }
}
</style>

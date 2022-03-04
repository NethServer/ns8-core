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
        v-for="searchId in searches"
        :key="searchId"
        :md="verticalLayout ? 4 : 8"
      >
        <LogSearch
          :searchId="searchId"
          :nodes="nodes"
          :apps="apps"
          :loadingNodes="loading.getClusterStatus"
          :loadingApps="loading.apps"
          :verticalLayout="verticalLayout"
          :numSearches="searches.length"
          @close="closeSearch"
          :closeAriaLabel="$t('common.close')"
          :light="true"
          :ref="'search-' + searchId"
        />
      </cv-column>
    </cv-row>
  </cv-grid>
</template>

<script>
import to from "await-to-js";
import { UtilService, IconService, TaskService } from "@nethserver/ns8-ui-lib";
import { v4 as uuidv4 } from "uuid";
import LogSearch from "../components/LogSearch.vue";

export default {
  name: "SystemLogs",
  components: { LogSearch },
  mixins: [IconService, UtilService, TaskService],
  pageTitle() {
    return this.$t("system_logs.title");
  },
  data() {
    return {
      nodes: [],
      apps: [],
      searches: [],
      verticalLayout: false,
      loading: {
        getClusterStatus: false,
        listInstalledModules: false,
      },
      error: {
        getClusterStatus: "",
        listInstalledModules: "",
      },
    };
  },
  computed: {
    // csbHorizontalLayoutSelected() { ////
    //   return !this.verticalLayout;
    // },
    // csbVerticalLayoutSelected() {
    //   return this.verticalLayout;
    // },
  },
  created() {
    this.getClusterStatus();
    this.listInstalledModules();
    this.addSearch();
  },
  methods: {
    addSearch() {
      const searchId = uuidv4();
      this.searches.push(searchId);

      // scroll to new search
      this.$nextTick(() => {
        const el = this.$refs["search-" + searchId][0].$el;
        this.scrollToElement(el);
      });
    },
    async getClusterStatus() {
      this.error.getClusterStatus = "";
      this.loading.getClusterStatus = true;
      const taskAction = "get-cluster-status";

      // register to task error
      this.$root.$off(taskAction + "-aborted");
      this.$root.$once(taskAction + "-aborted", this.getClusterStatusAborted);

      // register to task completion
      this.$root.$once(
        taskAction + "-completed",
        this.getClusterStatusCompleted
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

      let nodes = [];
      for (let node of clusterStatus.nodes) {
        nodes.push({
          name: node.hostname,
          label: node.ui_name
            ? node.ui_name + " (" + this.$t("common.node") + " " + node.id + ")"
            : this.$t("common.node") + " " + node.id,
          value: node.hostname,
        });
      }
      nodes.sort(this.sortByProperty("label"));
      this.nodes = nodes;
      this.loading.getClusterStatus = false;
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
            label: instance.ui_name
              ? instance.ui_name + " (" + instance.id + ")"
              : instance.id,
            value: instance.id,
          });
        }
      }
      apps.sort(this.sortByProperty("label"));
      this.apps = apps;
      this.loading.listInstalledModules = false;
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

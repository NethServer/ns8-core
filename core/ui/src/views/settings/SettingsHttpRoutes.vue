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
          <NsTabs
            :container="false"
            :aria-label="$t('common.tab_navigation')"
            :noDefaultToFirst="true"
            @tab-selected="tabSelected"
          >
            <cv-tab
              id="tab-1"
              :label="$t('settings_http_routes.routes')"
              :selected="q.view === 'routes'"
            >
              <HttpRoutesPanel
                :nodes="internalNodes"
                :traefikInstances="traefikInstances"
                :isLoadingInstances="loading.listInstalledModules"
                :instancesError="error.listInstalledModules"
                :selectedNodeId.sync="q.selectedNodeId"
              />
            </cv-tab>
            <cv-tab
              id="tab-2"
              :label="$t('settings_http_routes.frontend_proxies')"
              :selected="q.view === 'frontend-proxies'"
            >
              <HttpFrontendProxiesPanel
                :nodes="internalNodes"
                :traefikInstances="traefikInstances"
                :isLoadingInstances="loading.listInstalledModules"
                :instancesError="error.listInstalledModules"
              />
            </cv-tab>
          </NsTabs>
        </cv-column>
      </cv-row>
    </cv-grid>
  </div>
</template>

<script>
import to from "await-to-js";
import {
  QueryParamService,
  UtilService,
  TaskService,
  PageTitleService,
} from "@nethserver/ns8-ui-lib";
import { mapState } from "vuex";
import HttpRoutesPanel from "@/components/settings/HttpRoutesPanel.vue";
import HttpFrontendProxiesPanel from "@/components/settings/HttpFrontendProxiesPanel.vue";

// tab order: the view query param is the tab selector
const VIEWS = ["routes", "frontend-proxies"];

export default {
  name: "SettingsHttpRoutes",
  components: {
    HttpRoutesPanel,
    HttpFrontendProxiesPanel,
  },
  mixins: [TaskService, UtilService, QueryParamService, PageTitleService],
  pageTitle() {
    return this.$t("settings_http_routes.title");
  },
  data() {
    return {
      q: {
        view: "",
        selectedNodeId: "",
      },
      traefikInstances: [],
      // [eventName, handler] pairs registered on $root for the read chain
      readListeners: [],
      loading: {
        listInstalledModules: false,
      },
      error: {
        listInstalledModules: "",
      },
    };
  },
  computed: {
    ...mapState(["clusterNodes", "isWebsocketConnected"]),
    // clusterNodes can land after list-installed-modules, on a page reload it
    // usually does: recompute the list instead of snapshotting it once
    internalNodes() {
      return this.clusterNodes.map((node) => {
        const nodeId = node.id.toString();
        const traefikInstance = this.traefikInstances.find(
          (instance) => instance.node === nodeId
        );
        const internalNode = {
          name: nodeId,
          label: this.getShortNodeLabel(node),
          value: nodeId,
        };

        if (traefikInstance) {
          internalNode.label += ` (${traefikInstance.id})`;
          internalNode.traefikInstance = traefikInstance.id;
        }
        return internalNode;
      });
    },
  },
  watch: {
    isWebsocketConnected: function (isConnected) {
      // a Traefik restart kills this websocket: pending task events are lost.
      // Always restart from the instance list, it chains both tabs
      if (isConnected) {
        this.listInstalledModules();
      }
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
  beforeDestroy() {
    // a task completing after destroy would re-fire the whole read chain
    this.clearListeners(this.readListeners);
  },
  mounted() {
    // the tabs have no default selection: an unknown view would leave them all closed
    if (!VIEWS.includes(this.q.view)) {
      this.q.view = "routes";
    }
  },
  methods: {
    tabSelected(tabNum) {
      this.q.view = VIEWS[tabNum];
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
    async listInstalledModules() {
      // a handler of a previous chain would decrement the counter of this one
      this.clearListeners(this.readListeners);
      this.loading.listInstalledModules = true;
      this.error.listInstalledModules = "";
      const taskAction = "list-installed-modules";
      const eventId = this.getUuid();

      // register to task error
      this.registerListener(
        this.readListeners,
        `${taskAction}-aborted-${eventId}`,
        this.listInstalledModulesAborted
      );

      // register to task completion
      this.registerListener(
        this.readListeners,
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
        this.error.listInstalledModules = this.getErrorMessage(err);
        // otherwise the tables keep their skeleton rows for good
        this.loading.listInstalledModules = false;
      }
    },
    listInstalledModulesAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.listInstalledModules = this.$t("error.generic_error");
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
      // a new array, not an in-place mutation: both panels reload from a
      // traefikInstances watcher
      this.traefikInstances = traefikInstances;
      this.loading.listInstalledModules = false;
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";
</style>

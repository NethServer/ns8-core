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
              <span>{{ $t("settings_tls_certificates.title") }}</span>
            </cv-breadcrumb-item>
          </cv-breadcrumb>
        </cv-column>
      </cv-row>
      <cv-row>
        <cv-column class="subpage-title">
          <h3>{{ $t("settings_tls_certificates.title") }}</h3>
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
    </cv-grid>
    <cv-grid fullWidth>
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
              :label="$t('settings_tls_certificates.certificates')"
              :selected="q.view === 'certificates'"
            >
              <cv-row>
                <cv-column>
                  <TlsCertificatesPanel
                    :traefikInstances="traefikInstances"
                    :internalNodes="internalNodes"
                    :isLoadingInstances="loading.listInstalledModules"
                    :instancesError="error.listInstalledModules"
                    :selectedNodeId.sync="q.selectedNodeId"
                  />
                </cv-column>
              </cv-row>
            </cv-tab>
            <cv-tab
              id="tab-2"
              :label="$t('settings_acme_servers.acme_settings')"
              :selected="q.view === 'acme'"
            >
              <cv-row>
                <cv-column>
                  <AcmeSettingsPanel
                    v-if="acmeTabVisited"
                    :traefikInstances="traefikInstances"
                    :isLoadingInstances="loading.listInstalledModules"
                    :instancesError="error.listInstalledModules"
                  />
                </cv-column>
              </cv-row>
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
import TlsCertificatesPanel from "@/components/settings/TlsCertificatesPanel.vue";
import AcmeSettingsPanel from "@/components/settings/AcmeSettingsPanel.vue";

// order must match the cv-tab order: tabSelected maps the tab index
const TAB_VIEWS = ["certificates", "acme"];

export default {
  name: "SettingsTlsCertificates",
  components: {
    TlsCertificatesPanel,
    AcmeSettingsPanel,
  },
  mixins: [TaskService, UtilService, QueryParamService, PageTitleService],
  pageTitle() {
    return this.$t("settings_tls_certificates.title");
  },
  data() {
    return {
      q: {
        view: "",
        selectedNodeId: "",
      },
      // mount once and keep alive: a v-if would re-run the task chain
      acmeTabVisited: false,
      internalNodes: [],
      traefikInstances: [],
      taskListeners: [],
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
  },
  beforeRouteEnter(to, from, next) {
    next((vm) => {
      vm.watchQueryData(vm);
      vm.queryParamsToDataForCore(vm, to.query);
      vm.normalizeView();
    });
  },
  beforeRouteUpdate(to, from, next) {
    this.queryParamsToDataForCore(this, to.query);

    if (typeof to.query.view === "undefined") {
      this.q.view = TAB_VIEWS[0];
    }
    this.normalizeView();
    next();
  },
  watch: {
    isWebsocketConnected: function (isConnected) {
      if (isConnected) {
        this.listInstalledModules();
      }
    },
    "q.view": function (view) {
      if (view === "acme") {
        this.acmeTabVisited = true;
      }
    },
  },
  created() {
    this.listInstalledModules();
  },
  beforeDestroy() {
    this.clearTaskListeners();
  },
  methods: {
    registerTaskListener(eventName, handler) {
      this.$root.$once(eventName, handler);
      this.taskListeners.push([eventName, handler]);
    },
    clearTaskListeners() {
      for (const [eventName, handler] of this.taskListeners) {
        this.$root.$off(eventName, handler);
      }
      this.taskListeners = [];
    },
    normalizeView() {
      if (!TAB_VIEWS.includes(this.q.view)) {
        this.q.view = TAB_VIEWS[0];
      }
    },
    async listInstalledModules() {
      // a handler of a previous round would decrement the counters of this one
      this.clearTaskListeners();
      this.loading.listInstalledModules = true;
      this.error.listInstalledModules = "";
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
        this.loading.listInstalledModules = false;
        return;
      }
    },
    listInstalledModulesAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.listInstalledModules = this.$t("error.generic_error");
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
              instance.node_id = node.name;
            }
          }
        }
      }
      this.internalNodes = nodes;
      this.traefikInstances = traefikInstances;
      this.loading.listInstalledModules = false;

      this.$nextTick(() => {
        if (!this.q.selectedNodeId) {
          // initially show any node
          this.q.selectedNodeId = "any";
        } else {
          const nodeId = this.q.selectedNodeId;

          // workaround to update combo box
          this.q.selectedNodeId = "";
          this.$nextTick(() => {
            this.q.selectedNodeId = nodeId;
          });
        }
      });
    },
    tabSelected(tabNum) {
      if (tabNum == 0) {
        this.q.view = "certificates";
      } else if (tabNum == 1) {
        this.q.view = "acme";
      }
    },
  },
};
</script>

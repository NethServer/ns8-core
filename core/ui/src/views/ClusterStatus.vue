<template>
  <cv-grid fullWidth v-if="isClusterInitialized">
    <cv-row>
      <cv-column class="page-title">
        <h2>
          {{ $t("cluster_status.title") }}
          <cv-interactive-tooltip
            alignment="center"
            direction="bottom"
            class="info"
          >
            <template slot="trigger">
              <Information16 />
            </template>
            <template slot="content">
              {{ $t("common.global_shortcut") + ": CTRL+SHIFT+S" }}
            </template>
          </cv-interactive-tooltip>
        </h2>
      </cv-column>
    </cv-row>
    <cv-row>
      <cv-column :md="4" :max="4">
        <NsInfoCard
          light
          :title="moduleUpdates.length.toString()"
          :description="
            $tc('cluster_status.updated_available_c', moduleUpdates.length)
          "
          :icon="Upgrade32"
          :loading="loading.listModules"
          :isErrorShown="error.listModules"
          :errorTitle="$t('error.cannot_retrieve_available_apps')"
          :errorDescription="error.listModules"
          class="min-height-card"
        >
          <template slot="content">
            <NsButton
              kind="ghost"
              :icon="ArrowRight20"
              @click="$router.push('/software-center?search=&view=updates')"
            >
              {{ $t("cluster_status.go_to_software_center") }}
            </NsButton>
          </template>
        </NsInfoCard>
      </cv-column>
      <cv-column :md="4" :max="4">
        <NsInfoCard
          light
          :title="installedModules.length.toString()"
          :description="
            $tc('cluster_status.apps_installed_c', installedModules.length)
          "
          :icon="Application32"
          :loading="loading.listModules"
          :isErrorShown="error.listModules"
          :errorTitle="$t('error.cannot_retrieve_available_apps')"
          :errorDescription="error.listModules"
          class="min-height-card"
        >
          <template slot="content">
            <NsButton
              kind="ghost"
              :icon="ArrowRight20"
              @click="$router.push('/software-center?search=&view=installed')"
            >
              {{ $t("cluster_status.go_to_software_center") }}
            </NsButton>
          </template>
        </NsInfoCard>
      </cv-column>
      <cv-column :md="4" :max="4">
        <NsInfoCard
          light
          :title="nodes.length.toString()"
          :description="$tc('common.nodes_c', nodes.length)"
          :icon="Chip32"
          :loading="loading.getClusterStatus"
          :isErrorShown="error.getClusterStatus"
          :errorTitle="$t('error.cannot_retrieve_cluster_nodes')"
          :errorDescription="error.getClusterStatus"
          class="min-height-card"
        >
          <template slot="content">
            <NsButton
              kind="ghost"
              :icon="ArrowRight20"
              @click="$router.push('/nodes')"
            >
              {{ $t("cluster_status.go_to_nodes") }}
            </NsButton>
          </template>
        </NsInfoCard>
      </cv-column>
      <cv-column :md="4" :max="4">
        <NsInfoCard
          light
          :title="backups.length.toString()"
          :description="$tc('backup.backups_scheduled_c', backups.length)"
          :icon="Save32"
          :loading="loading.listBackups"
          :isErrorShown="error.listBackups"
          :errorTitle="$t('error.cannot_retrieve_backups')"
          :errorDescription="error.listBackups"
          class="min-height-card"
        >
          <template slot="content">
            <div>
              <div
                v-if="!loading.listBackups && !error.listBackups"
                class="card-row mg-top-sm"
              >
                <template
                  v-if="erroredBackups.length || instancesNotBackedUp.length"
                >
                  <div
                    v-if="erroredBackups.length"
                    class="card-row icon-and-text"
                  >
                    <NsSvg :svg="ErrorFilled16" class="icon ns-error" />
                    <span>
                      {{
                        $tc(
                          "cluster_status.backups_failed_c",
                          erroredBackups.length,
                          { num: erroredBackups.length }
                        )
                      }}
                    </span>
                  </div>
                  <div
                    v-if="instancesNotBackedUp.length"
                    class="card-row icon-and-text"
                  >
                    <NsSvg :svg="Warning16" class="icon ns-warning" />
                    <span>
                      {{
                        $tc(
                          "cluster_status.instances_not_backed_up_c",
                          instancesNotBackedUp.length,
                          { num: instancesNotBackedUp.length }
                        )
                      }}
                    </span>
                  </div>
                </template>
                <template v-else>
                  <div class="card-row icon-and-text">
                    <NsSvg :svg="CheckmarkFilled16" class="icon ns-success" />
                    <span>{{ $t("common.all_good") }}</span>
                  </div>
                </template>
              </div>
              <div class="card-row">
                <NsButton
                  kind="ghost"
                  :icon="ArrowRight20"
                  @click="$router.push('/backup')"
                >
                  {{ $t("cluster_status.go_to_backup") }}
                </NsButton>
              </div>
            </div>
          </template>
        </NsInfoCard>
      </cv-column>
    </cv-row>
  </cv-grid>
</template>

<script>
import NotificationService from "@/mixins/notification";
import Information16 from "@carbon/icons-vue/es/information/16";
import { mapState, mapActions } from "vuex";
import to from "await-to-js";
import WebSocketService from "@/mixins/websocket";
import {
  QueryParamService,
  UtilService,
  TaskService,
  IconService,
  StorageService,
} from "@nethserver/ns8-ui-lib";

export default {
  name: "ClusterStatus",
  components: { Information16 },
  mixins: [
    NotificationService,
    QueryParamService,
    UtilService,
    TaskService,
    WebSocketService,
    IconService,
    StorageService,
  ],
  pageTitle() {
    return this.$t("cluster_status.title");
  },
  data() {
    return {
      q: {
        testInput: "",
        testToggle: false,
        moduleToAdd: "",
        name: "",
        email: "",
      },
      nodes: [],
      apps: [],
      updates: [],
      backups: [],
      instancesNotBackedUp: [],
      loading: {
        getClusterStatus: true,
        listModules: true,
        listBackups: true,
      },
      error: {
        name: "",
        email: "",
        getClusterStatus: "",
        listModules: "",
        listBackups: "",
      },
      toastVisible: true,
      toastTitle: "Toast title",
      toastSubTitle: "Toast subtitle",
      toastCaption: "Toast caption",
      lowContrast: false,
      snippet: 'printf("A short bit of code.");',
    };
  },
  computed: {
    ...mapState([
      "notifications",
      "loggedUser",
      "isWebsocketConnected",
      "isClusterInitialized",
    ]),
    erroredBackups() {
      return this.backups.filter((b) => b.errorInstances.length);
    },
    installedModules() {
      let installedModules = [];

      for (const app of this.apps) {
        for (const installedModule of app.installed) {
          if (
            !installedModule.flags.includes("account_provider") &&
            !installedModule.flags.includes("core_module")
          ) {
            installedModules.push(installedModule);
          }
        }
      }
      return installedModules;
    },
    moduleUpdates() {
      let moduleUpdates = [];

      for (const appUpdate of this.updates) {
        for (const moduleUpdate of appUpdate.updates) {
          if (
            !moduleUpdate.flags.includes("account_provider") &&
            !moduleUpdate.flags.includes("core_module")
          ) {
            moduleUpdates.push(moduleUpdate);
          }
        }
      }
      return moduleUpdates;
    },
  },
  watch: {
    isWebsocketConnected: function () {
      if (this.isWebsocketConnected && this.loggedUser) {
        // retrieve initial data
        this.retrieveData();
      }
    },
    loggedUser: function () {
      if (this.isWebsocketConnected && this.loggedUser) {
        // retrieve initial data
        this.retrieveData();
      }
    },
  },
  created() {
    if (this.isWebsocketConnected && this.loggedUser) {
      // retrieve initial data
      this.retrieveData();
    }
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
    ...mapActions(["setClusterNodesInStore"]),
    retrieveData() {
      this.getClusterStatus();
      this.listModules();
      //// TODO this.listCoreModules();
      this.listBackups();
    },
    async getClusterStatus() {
      this.error.getClusterStatus = "";
      this.loading.getClusterStatus = true;
      const taskAction = "get-cluster-status";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.getClusterStatusAborted
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.getClusterStatusCompleted
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
        this.error.getClusterStatus = this.getErrorMessage(err);
        this.loading.getClusterStatus = false;
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
      const nodes = clusterStatus.nodes.sort(this.sortByProperty("id"));
      this.setClusterNodesInStore(nodes);
      this.nodes = nodes;

      // update nodes in vuex store
      this.setClusterNodesInStore(this.nodes);

      this.loading.getClusterStatus = false;
    },
    async listModules() {
      this.loading.modules = true;
      this.error.listModules = "";
      const taskAction = "list-modules";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.listModulesAborted
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.listModulesCompleted
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
        this.error.listModules = this.getErrorMessage(err);
        this.loading.listModules = false;
        return;
      }
    },
    listModulesAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.listModules = this.$t("error.generic_error");
      this.loading.listModules = false;
    },
    listModulesCompleted(taskContext, taskResult) {
      let apps = taskResult.output;
      apps.sort(this.sortByProperty("name"));
      let updates = [];

      for (const app of apps) {
        if (app.updates.length) {
          updates.push(app);
        }

        // sort installed instances
        app.installed.sort(this.sortModuleInstances());
      }
      this.updates = updates;
      this.apps = apps;
      this.loading.listModules = false;
    },
    async listBackups() {
      this.loading.listBackups = true;
      this.error.listBackups = "";
      const taskAction = "list-backups";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.listBackupsAborted
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.listBackupsCompleted
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
        this.error.listBackups = this.getErrorMessage(err);
        return;
      }
    },
    listBackupsAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.listBackups = this.$t("error.generic_error");
      this.loading.listBackups = false;
    },
    listBackupsCompleted(taskContext, taskResult) {
      this.instancesNotBackedUp = taskResult.output.unconfigured_instances;
      let backups = taskResult.output.backups;

      for (const backup of backups) {
        backup.errorInstances = backup.instances.filter(
          (i) => i.status == false
        );
      }
      this.backups = backups;
      this.loading.listBackups = false;
    },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";

.card-row {
  margin-bottom: $spacing-05;
  display: flex;
  justify-content: center;
}

.card-row:last-child {
  margin-bottom: 0;
}
</style>

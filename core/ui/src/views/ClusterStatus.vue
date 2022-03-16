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
          :title="apps.length.toString()"
          :description="$tc('common.installed_apps', apps.length)"
          :icon="Application32"
          :loading="loading.listInstalledModules"
          :isErrorShown="error.listInstalledModules"
          :errorTitle="$t('error.cannot_retrieve_installed_apps')"
          :errorDescription="error.listInstalledModules"
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
          :loading="loading.nodes"
          :isErrorShown="error.nodes"
          :errorTitle="$t('error.cannot_retrieve_cluster_nodes')"
          :errorDescription="error.nodes"
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
import NodeService from "@/mixins/node";
import Information16 from "@carbon/icons-vue/es/information/16";
import { mapState } from "vuex";
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
    NodeService,
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
      backups: [],
      instancesNotBackedUp: [],
      loading: {
        nodes: true,
        listInstalledModules: true,
        listBackups: true,
      },
      error: {
        name: "",
        email: "",
        nodes: "",
        listInstalledModules: "",
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
    retrieveData() {
      this.retrieveClusterNodes();
      this.listInstalledModules();
      this.listBackups();
    },
    async retrieveClusterNodes() {
      //// retrieve cluster status instead of retrieveClusterNodes? It provides more info
      this.loading.nodes = true;
      const [errNodes, responseNodes] = await to(this.getNodes());

      if (errNodes) {
        console.error("error retrieving cluster nodes", errNodes);
        this.error.nodes = this.getErrorMessage(errNodes);
        this.loading.nodes = false;
        return;
      }

      this.nodes = responseNodes.data.data.list;
      this.loading.nodes = false;
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
        this.error.listInstalledModules = this.getErrorMessage(err);
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
      let apps = [];

      for (let instanceList of Object.values(taskResult.output)) {
        for (let instance of instanceList) {
          if (
            !instance.flags.includes("account_provider") &&
            !instance.flags.includes("core_module")
          ) {
            apps.push(instance);
          }
        }
      }
      this.apps = apps;
      this.loading.listInstalledModules = false;
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
}

.card-row:last-child {
  margin-bottom: 0;
}
</style>

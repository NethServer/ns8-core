<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
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
      <cv-column>
        <!-- card grid -->
        <div
          class="card-grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-3 3xl:grid-cols-4"
        >
          <NsInfoCard
            light
            :title="numUpdates.toString()"
            :description="$tc('cluster_status.updated_available_c', numUpdates)"
            :icon="Upgrade32"
            :loading="loading.listModules || loading.listCoreModules"
            :isErrorShown="error.listCoreModules"
            :errorTitle="$t('error.cannot_retrieve_available_apps')"
            :errorDescription="error.listCoreModules"
            class="min-height-card"
          >
            <template slot="content">
              <div class="card-rows">
                <div
                  v-if="isCoreUpdateAvailable"
                  class="card-row mg-top-sm icon-and-text"
                >
                  <NsSvg :svg="Warning16" class="icon ns-warning" />
                  <span>{{ $t("cluster_status.core_update_available") }}</span>
                </div>
                <div class="card-row">
                  <NsButton
                    kind="ghost"
                    :icon="ArrowRight20"
                    @click="
                      $router.push('/software-center?search=&view=updates')
                    "
                  >
                    {{ $t("cluster_status.go_to_software_center") }}
                  </NsButton>
                </div>
              </div>
            </template>
          </NsInfoCard>
          <NsInfoCard
            light
            :title="$t('cluster_status.subscription_status')"
            :icon="Badge32"
            :loading="loading.getSubscription"
            :isErrorShown="error.getSubscription"
            :errorTitle="$t('error.cannot_retrieve_subscription_status')"
            :errorDescription="error.getSubscription"
            class="min-height-card"
          >
            <template slot="content">
              <div class="card-rows">
                <div class="card-row">
                  <div v-if="!loading.getSubscription" class="card-row">
                    <cv-tag
                      v-if="subscription_status === 'active'"
                      kind="green"
                      :label="$t('common.active')"
                    ></cv-tag>
                    <cv-tag
                      v-else
                      kind="high-contrast"
                      :label="$t('common.not_active')"
                    ></cv-tag>
                  </div>
                </div>
                <div
                  v-if="support_active && !loading.getSubscription"
                  class="card-row"
                >
                  <div class="mg-top-sm icon-and-text">
                    <NsSvg :svg="InformationFilled16" class="icon ns-info" />
                    <span>{{
                      $t("settings_subscription.remote_support_in_progress")
                    }}</span>
                  </div>
                </div>
                <div v-if="!loading.getSubscription" class="card-row">
                  <NsButton
                    kind="ghost"
                    :icon="ArrowRight20"
                    @click="$router.push('/settings/subscription')"
                  >
                    {{ $t("cluster_status.go_to_subscription") }}
                  </NsButton>
                </div>
              </div>
            </template>
          </NsInfoCard>
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
              <div class="card-rows">
                <div
                  class="card-row"
                  v-if="!loading.getClusterStatus && nodes.length > 1"
                >
                  <div class="mg-top-sm">
                    <div
                      v-if="nodesOffline.length"
                      class="card-row icon-and-text"
                    >
                      <NsSvg :svg="ErrorFilled16" class="icon ns-error" />
                      <span>
                        {{
                          $tc("nodes.nodes_offline_c", nodesOffline.length, {
                            num: nodesOffline.length,
                          })
                        }}
                      </span>
                    </div>
                    <div v-else class="card-row icon-and-text">
                      <NsSvg :svg="CheckmarkFilled16" class="icon ns-success" />
                      <span>
                        {{ $t("nodes.nodes_online_cluster_status") }}
                      </span>
                    </div>
                  </div>
                </div>
                <div class="card-row">
                  <NsButton
                    kind="ghost"
                    :icon="ArrowRight20"
                    @click="$router.push('/nodes')"
                  >
                    {{ $t("cluster_status.go_to_nodes") }}
                  </NsButton>
                </div>
              </div>
            </template>
          </NsInfoCard>
          <NsInfoCard
            light
            :title="enabledBackupSchedules.length.toString()"
            :description="$tc('backup.backups_scheduled_c', backups.length)"
            :icon="Save32"
            :loading="loading.listBackups"
            :isErrorShown="error.listBackups"
            :errorTitle="$t('error.cannot_retrieve_backups')"
            :errorDescription="error.listBackups"
            class="min-height-card"
          >
            <template slot="content">
              <div class="card-rows mg-top-sm">
                <div
                  v-if="!loading.listBackups && !error.listBackups"
                  class="mg-bottom-sm"
                >
                  <template
                    v-if="
                      erroredBackups.length ||
                      disabledBackupSchedules.length ||
                      instancesNotBackedUp.length
                    "
                  >
                    <!-- some backups failed -->
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
                    <!-- some backups are disabled -->
                    <div
                      v-if="disabledBackupSchedules.length"
                      class="card-row icon-and-text"
                    >
                      <NsSvg :svg="Warning16" class="icon ns-warning" />
                      <span>
                        {{
                          $tc(
                            "cluster_status.backup_disabled_c",
                            disabledBackupSchedules.length,
                            { num: disabledBackupSchedules.length }
                          )
                        }}
                      </span>
                    </div>
                    <!-- some instances are not backed up -->
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
          <NsInfoCard
            light
            :title="$t('cluster_status.email_notification')"
            :icon="Email32"
            :loading="loading.getEmailNotification"
            :isErrorShown="error.getEmailNotification"
            :errorTitle="$t('error.cannot_retrieve_email_notification_status')"
            :errorDescription="error.getEmailNotification"
            class="min-height-card"
          >
            <template slot="content">
              <div class="card-rows">
                <div class="card-row">
                  <div v-if="!loading.getEmailNotification" class="card-row">
                    <cv-tag
                      v-if="emailNotificationEnabled"
                      kind="green"
                      :label="$t('common.enabled')"
                    ></cv-tag>
                    <cv-tag
                      v-else
                      kind="high-contrast"
                      :label="$t('common.disabled')"
                    ></cv-tag>
                  </div>
                </div>
                <div class="card-row">
                  <NsButton
                    kind="ghost"
                    :icon="ArrowRight20"
                    @click="$router.push('/settings/smarthost')"
                  >
                    {{ $t("cluster_status.go_to_settings") }}
                  </NsButton>
                </div>
              </div>
            </template>
          </NsInfoCard>
          <NsInfoCard
            light
            :title="$t('cluster_status.system_logs')"
            :description="
              $tc('cluster_status.active_instance', activeLokiInstance, {
                instance: activeLokiInstance,
              })
            "
            :icon="Catalog32"
            :loading="loading.lokiInstance"
            :isErrorShown="error.lokiInstance"
            :errorTitle="
              $t('system_logs.loki.cannot_retrieve_system_logs_status')
            "
            :errorDescription="error.lokiInstance"
            class="min-height-card"
          >
            <template slot="content">
              <div class="card-rows">
                <div
                  v-if="cloudLogManagerForwarderStatus == 'active'"
                  class="card-row mg-top-sm icon-and-text"
                >
                  <NsSvg :svg="CheckmarkFilled16" class="icon ns-success" />
                  <span>{{
                    $t(
                      "cloud_log_manager_forwarder.cloud_log_manager_export_enabled"
                    )
                  }}</span>
                </div>
                <div
                  v-if="cloudLogManagerForwarderStatus == 'failed'"
                  class="card-row mg-top-sm icon-and-text"
                >
                  <NsSvg :svg="ErrorFilled16" class="icon ns-error" />
                  <span>{{
                    $t(
                      "cloud_log_manager_forwarder.cloud_log_manager_export_failed"
                    )
                  }}</span>
                </div>
                <div
                  v-if="syslogForwarderStatus == 'active'"
                  class="card-row mg-top-sm icon-and-text"
                >
                  <NsSvg :svg="CheckmarkFilled16" class="icon ns-success" />
                  <span>{{
                    $t("syslog_forwarder.syslog_export_enabled")
                  }}</span>
                </div>
                <div
                  v-if="syslogForwarderStatus == 'failed'"
                  class="card-row mg-top-sm icon-and-text"
                >
                  <NsSvg :svg="ErrorFilled16" class="icon ns-error" />
                  <span>{{ $t("syslog_forwarder.syslog_export_failed") }}</span>
                </div>
                <div class="card-row">
                  <NsButton
                    kind="ghost"
                    :icon="ArrowRight20"
                    @click="$router.push('/settings/system-logs')"
                  >
                    {{ $t("cluster_status.go_to_system_logs") }}
                  </NsButton>
                </div>
              </div>
            </template>
          </NsInfoCard>
        </div>
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
import { mapGetters } from "vuex";
import {
  QueryParamService,
  UtilService,
  TaskService,
  IconService,
  StorageService,
  PageTitleService,
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
    PageTitleService,
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
      isCoreUpdateAvailable: false,
      subscription_status: "inactive",
      emailNotificationEnabled: false,
      support_active: false,
      activeLokiInstance: "",
      cloudLogManagerForwarderStatus: "inactive",
      syslogForwarderStatus: "inactive",
      loading: {
        getClusterStatus: true,
        listModules: true,
        listBackups: true,
        listCoreModules: true,
        getSubscription: true,
        getEmailNotification: true,
        getSupportSession: true,
        lokiInstance: true,
      },
      error: {
        name: "",
        email: "",
        getClusterStatus: "",
        listModules: "",
        listBackups: "",
        listCoreModules: "",
        getSubscription: "",
        getEmailNotification: "",
        lokiInstance: "",
      },
    };
  },
  computed: {
    ...mapState([
      "notifications",
      "loggedUser",
      "isWebsocketConnected",
      "isClusterInitialized",
    ]),
    ...mapGetters(["leaderNode"]),
    erroredBackups() {
      const erroredInstances = [];

      for (const backup of this.backups) {
        if (!backup.enabled) {
          continue;
        }

        for (const instance of backup.instances) {
          if (instance.status && instance.status.success == false) {
            erroredInstances.push(instance.module_id);
          }
        }
      }
      return erroredInstances;
    },
    disabledBackupSchedules() {
      return this.backups.filter((b) => !b.enabled);
    },
    enabledBackupSchedules() {
      return this.backups.filter((b) => b.enabled);
    },
    nodesOffline() {
      return this.nodes.filter((n) => n.online == false);
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
    numUpdates() {
      if (this.isCoreUpdateAvailable) {
        return this.moduleUpdates.length + 1;
      } else {
        return this.moduleUpdates.length;
      }
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
      this.getActiveLokiInstance();
      this.getClusterStatus();
      this.listModules();
      this.listCoreModules();
      this.listBackups();
      this.getSubscription();
      this.getEmailNotification();
    },
    async getSupportSession() {
      if (!this.leaderNode) {
        return;
      }
      this.error.getSupportSession = "";
      this.loading.getSupportSession = true;
      const taskAction = "get-support-session";

      // register to task completion
      this.$root.$once(
        taskAction + "-completed",
        this.getSupportSessionCompleted
      );

      const res = await to(
        this.createNodeTask(this.leaderNode.id, {
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
        this.error.getSupportSession = this.getErrorMessage(err);
        this.loading.getSupportSession = false;
        return;
      }
    },
    getSupportSessionCompleted(taskContext, taskResult) {
      const output = taskResult.output;
      this.support_active = output.active;
      this.loading.getSupportSession = false;
    },
    async getSubscription() {
      this.clearErrors();
      this.loading.getSubscription = true;
      const taskAction = "get-subscription";

      // register to task completion
      this.$root.$once(
        taskAction + "-completed",
        this.getSubscriptionCompleted
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
        this.error.getSubscription = this.getErrorMessage(err);
        this.loading.getSubscription = false;
        return;
      }
    },
    getSubscriptionCompleted(taskContext, taskResult) {
      const output = taskResult.output;
      if (output.subscription == null) {
        this.status = "inactive";
        this.loading.getSubscription = false;
        return;
      }
      this.getSupportSession();
      this.subscription_status = output.subscription.status;
      this.loading.getSubscription = false;
    },
    async getEmailNotification() {
      this.clearErrors();
      this.loading.getEmailNotification = true;
      const taskAction = "get-smarthost";

      // register to task completion
      this.$root.$once(
        taskAction + "-completed",
        this.getEmailNotificationCompleted
      );

      // register to task error
      this.$root.$once(
        taskAction + "-aborted",
        this.getEmailNotificationAborted
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
        this.error.getEmailNotification = this.getErrorMessage(err);
        this.loading.getEmailNotification = false;
        return;
      }
    },
    getEmailNotificationCompleted(taskContext, taskResult) {
      this.emailNotificationEnabled = taskResult.output.enabled;
      this.getSupportSession();
      this.loading.getEmailNotification = false;
    },
    getEmailNotificationAborted(taskContext, taskResult) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.getEmailNotification = this.$t("error.generic_error");
      this.loading.getEmailNotification = false;
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
        const hasStableUpdate = app.updates.some((update) => update.update);

        if (hasStableUpdate) {
          updates.push(app);
        }

        // sort installed instances
        app.installed.sort(this.sortModuleInstances());
      }
      this.updates = updates;
      this.apps = apps;
      this.loading.listModules = false;
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
      const coreApps = taskResult.output;
      let isCoreUpdateAvailable = false;

      for (const coreApp of coreApps) {
        for (const coreInstance of coreApp.instances) {
          if (coreInstance.update) {
            isCoreUpdateAvailable = true;
          }
        }
      }
      this.isCoreUpdateAvailable = isCoreUpdateAvailable;
      this.loading.listCoreModules = false;
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
      this.backups = taskResult.output.backups;
      this.loading.listBackups = false;
    },
    async getActiveLokiInstance() {
      this.loading.lokiInstance = true;
      this.error.lokiInstance = "";
      const taskAction = "list-loki-instances";

      // register to task abortion
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
        this.error.lokiInstance = this.getErrorMessage(err);
        return;
      }
    },
    getClusterLokiInstancesAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.lokiInstance = this.$t("error.generic_error");
      this.loading.lokiInstance = false;
    },
    getClusterLokiInstancesCompleted(taskContext, taskResult) {
      for (const instance of taskResult.output.instances) {
        if (instance.active) {
          if (instance.instance_label != "") {
            this.activeLokiInstance = instance.instance_label;
          } else {
            this.activeLokiInstance = instance.instance_id;
          }
          this.cloudLogManagerForwarderStatus =
            instance.cloud_log_manager.status;
          this.syslogForwarderStatus = instance.syslog.status;
        }
      }
      this.loading.lokiInstance = false;
    },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";

.card-rows {
  display: flex;
  flex-direction: column;
}

.card-row {
  margin-bottom: $spacing-05;
  display: flex;
  justify-content: center;
}

.card-row:last-child {
  margin-bottom: 0;
}
</style>

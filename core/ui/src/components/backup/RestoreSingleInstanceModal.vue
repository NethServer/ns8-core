<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <NsWizard
    size="default"
    :visible="isShown"
    :cancelLabel="$t('common.cancel')"
    :previousLabel="$t('common.previous')"
    :nextLabel="isLastStep ? $t('backup.restore') : $t('common.next')"
    :isPreviousDisabled="isFirstStep || loading.restoreModule"
    :isNextDisabled="isNextButtonDisabled"
    :isNextLoading="loading.restoreModule"
    @modal-hidden="$emit('hide')"
    @cancel="$emit('hide')"
    @previousStep="previousStep"
    @nextStep="nextStep"
  >
    <template slot="title">{{ $t("backup.restore_app") }}</template>
    <template slot="content">
      <cv-form @submit.prevent="nextStep">
        <template v-if="step == 'instance'">
          <div class="mg-bottom-md">
            {{ $t("backup.select_instance_to_restore") }}
          </div>
          <NsInlineNotification
            v-if="error.readBackupRepositories"
            kind="error"
            :title="$t('action.read-backup-repositories')"
            :description="error.readBackupRepositories"
            :showCloseButton="false"
          />
          <cv-grid class="mg-top-xlg mg-bottom-md no-padding">
            <cv-row>
              <cv-column>
                <RestoreSingleInstanceSelector
                  :instances="instances"
                  :loading="loading.readBackupRepositories"
                  :light="true"
                  @select="onSelectInstance"
                />
              </cv-column>
            </cv-row>
            <cv-row>
              <cv-column>
                <cv-checkbox
                  :label="
                    $t('backup.replace_existing_app', {
                      app: instanceToReplace,
                    })
                  "
                  v-model="replaceExistingApp"
                  :disabled="!instanceToReplace"
                  value="checkReplaceExistingApp"
                  class="mg-top-xlg"
                />
              </cv-column>
            </cv-row>
          </cv-grid>
        </template>
        <template v-if="step == 'snapshot'">
          <div>
            {{ $t("backup.select_backup_snapshot") }}
          </div>
          <NsInlineNotification
            v-if="error.readBackupSnapshots"
            kind="error"
            :title="$t('action.read-backup-snapshots')"
            :description="error.readBackupSnapshots"
            :showCloseButton="false"
          />
          <cv-grid class="mg-top-xlg mg-bottom-md no-padding">
            <cv-row>
              <cv-column>
                <RestoreSingleInstanceSnapshotSelector
                  :snapshots="backupSnapshots"
                  :loading="loading.readBackupSnapshots"
                  :light="true"
                  @select="onSelectSnapshot"
                />
              </cv-column>
            </cv-row>
          </cv-grid>
        </template>
        <template v-if="step == 'node'">
          <div>
            {{ $t("backup.select_restore_node") }}
          </div>
          <NsInlineNotification
            v-if="error.determineRestoreEligibility"
            kind="error"
            :title="$t('action.determine-restore-eligibility')"
            :description="error.determineRestoreEligibility"
            :showCloseButton="false"
          />
          <NsInlineNotification
            v-if="clusterNodes.length == disabledNodes.length"
            kind="info"
            :title="$t('backup.no_node_eligible_for_instance_restoration')"
            :showCloseButton="false"
          />
          <NodeSelector
            @selectNode="onSelectNode"
            :disabledNodes="disabledNodes"
            :loading="loading.determineRestoreEligibility"
            class="mg-top-xlg"
          >
            <template v-for="(nodeInfoMessage, nodeId) in nodesInfo">
              <template :slot="`node-${nodeId}`">
                {{ nodeInfoMessage }}
              </template>
            </template>
          </NodeSelector>
          <NsInlineNotification
            v-if="error.restoreModule"
            kind="error"
            :title="$t('action.restore-module')"
            :description="error.restoreModule"
            :showCloseButton="false"
          />
        </template>
      </cv-form>
    </template>
  </NsWizard>
</template>

<script>
import { UtilService, TaskService, IconService } from "@nethserver/ns8-ui-lib";
import to from "await-to-js";
import NodeSelector from "@/components/nodes/NodeSelector";
import RestoreSingleInstanceSelector from "@/components/backup/RestoreSingleInstanceSelector";
import { mapState } from "vuex";
import RestoreSingleInstanceSnapshotSelector from "./RestoreSingleInstanceSnapshotSelector.vue";

export default {
  name: "RestoreSingleInstanceModal",
  components: {
    NodeSelector,
    RestoreSingleInstanceSelector,
    RestoreSingleInstanceSnapshotSelector,
  },
  mixins: [UtilService, TaskService, IconService],
  props: {
    isShown: {
      type: Boolean,
      default: true,
    },
  },
  data() {
    return {
      step: "",
      steps: ["instance", "snapshot", "node"],
      instances: [],
      backupSnapshots: [],
      selectedNode: null,
      selectedInstance: null,
      selectedSnapshot: null,
      replaceExistingApp: false,
      installDestinations: [],
      loading: {
        readBackupRepositories: true,
        restoreModule: false,
        readBackupSnapshots: false,
        determineRestoreEligibility: false,
      },
      error: {
        readBackupRepositories: "",
        restoreModule: "",
        readBackupSnapshots: "",
        determineRestoreEligibility: "",
      },
    };
  },
  computed: {
    ...mapState(["clusterNodes"]),
    stepIndex() {
      return this.steps.indexOf(this.step);
    },
    isFirstStep() {
      return this.stepIndex == 0;
    },
    isLastStep() {
      return this.stepIndex == this.steps.length - 1;
    },
    isNextButtonDisabled() {
      return (
        this.loading.restoreModule ||
        (this.step == "instance" && !this.selectedInstance) ||
        (this.step == "snapshot" && !this.selectedSnapshot) ||
        (this.step == "node" &&
          (this.loading.determineRestoreEligibility ||
            !this.selectedNode ||
            this.clusterNodes.length == this.disabledNodes.length))
      );
    },
    instanceToReplace() {
      if (
        !(this.selectedInstance && this.selectedInstance.installed_instance)
      ) {
        return "";
      }

      if (this.selectedInstance.installed_instance_ui_name) {
        return (
          this.selectedInstance.installed_instance_ui_name +
          " (" +
          this.selectedInstance.installed_instance +
          ")"
        );
      } else {
        return this.selectedInstance.installed_instance;
      }
    },
    nodesInfo() {
      const nodesInfo = {};

      for (const nodeInfo of this.installDestinations) {
        if (!nodeInfo.eligible) {
          // show reason why node is not eligible
          const rejectReason = nodeInfo.reject_reason;

          if (rejectReason.message === "max_per_node_limit") {
            const numMaxInstances = rejectReason.parameter;
            nodesInfo[nodeInfo.node_id] = this.$tc(
              `software_center.reason_${rejectReason.message}`,
              numMaxInstances,
              { param: numMaxInstances }
            );
          } else {
            nodesInfo[nodeInfo.node_id] = this.$t(
              `software_center.reason_${rejectReason.message}`,
              { param: rejectReason.parameter }
            );
          }
        } else if (nodeInfo.instances) {
          // show number of instances installed
          nodesInfo[nodeInfo.node_id] = this.$tc(
            "software_center.num_instances_installed",
            nodeInfo.instances,
            { num: nodeInfo.instances }
          );
        }
      }
      return nodesInfo;
    },
    disabledNodes() {
      return this.installDestinations
        .filter((nodeInfo) => !nodeInfo.eligible)
        .map((nodeInfo) => nodeInfo.node_id);
    },
  },
  watch: {
    isShown: function () {
      if (this.isShown) {
        // show first step
        this.step = this.steps[0];
        this.selectedInstance = null;
        this.replaceExistingApp = false;
        this.readBackupRepositories();
      }
    },
    step: function () {
      if (this.step == "instance") {
        this.selectedInstance = null;
        this.replaceExistingApp = false;
      } else if (this.step == "snapshot") {
        this.selectedSnapshot = null;
        this.readBackupSnapshots();
      } else if (this.step == "node") {
        this.determineRestoreEligibility();
      }
    },
  },
  methods: {
    nextStep() {
      if (this.isNextButtonDisabled) {
        return;
      }

      if (this.isLastStep) {
        this.restoreModule();
      } else {
        this.step = this.steps[this.stepIndex + 1];
      }
    },
    previousStep() {
      if (!this.isFirstStep) {
        this.step = this.steps[this.stepIndex - 1];
      }
    },
    async readBackupRepositories() {
      this.error.readBackupRepositories = "";
      this.loading.readBackupRepositories = true;
      const taskAction = "read-backup-repositories";

      // register to task error
      this.$root.$off(taskAction + "-aborted");
      this.$root.$once(
        taskAction + "-aborted",
        this.readBackupRepositoriesAborted
      );

      // register to task completion
      this.$root.$off(taskAction + "-completed");
      this.$root.$once(
        taskAction + "-completed",
        this.readBackupRepositoriesCompleted
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
        this.error.readBackupRepositories = this.getErrorMessage(err);
        return;
      }
    },
    readBackupRepositoriesAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.readBackupRepositories = this.$t("error.generic_error");
      this.loading.readBackupRepositories = false;
    },
    readBackupRepositoriesCompleted(taskContext, taskResult) {
      let instances = taskResult.output;

      // sort instances by timestamp
      instances.sort(this.sortByProperty("timestamp")).reverse();
      this.instances = instances;
      this.loading.readBackupRepositories = false;
    },
    async readBackupSnapshots() {
      this.error.readBackupSnapshots = "";
      this.loading.readBackupSnapshots = true;
      const taskAction = "read-backup-snapshots";

      // register to task error
      this.$root.$off(taskAction + "-aborted");
      this.$root.$once(
        taskAction + "-aborted",
        this.readBackupSnapshotsAborted
      );

      // register to task completion
      this.$root.$off(taskAction + "-completed");
      this.$root.$once(
        taskAction + "-completed",
        this.readBackupSnapshotsCompleted
      );

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            repository: this.selectedInstance.repository_id,
            path: this.selectedInstance.path,
          },
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.readBackupSnapshots = this.getErrorMessage(err);
        return;
      }
    },
    readBackupSnapshotsAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);

      if (taskResult.error.includes("wrong password or no key found")) {
        this.error.readBackupSnapshots = this.$t(
          "backup.wrong_password_or_no_key_found"
        );
      } else {
        this.error.readBackupSnapshots = this.$t("error.generic_error");
      }
      this.loading.readBackupSnapshots = false;
    },
    readBackupSnapshotsCompleted(taskContext, taskResult) {
      let snapshots = taskResult.output;
      // sort snapshots by timestamp descending (most recent first)
      snapshots.sort(this.sortByProperty("timestamp")).reverse();
      this.backupSnapshots = snapshots;
      this.loading.readBackupSnapshots = false;
    },
    async restoreModule() {
      this.loading.restoreModule = true;
      this.error.restoreModule = "";
      const taskAction = "restore-module";

      // register to task error
      this.$root.$off(taskAction + "-aborted");
      this.$root.$once(taskAction + "-aborted", this.restoreModuleAborted);

      // register to task completion
      this.$root.$off(taskAction + "-completed");
      this.$root.$once(taskAction + "-completed", this.restoreModuleCompleted);

      const app = this.selectedInstance.path.split("/")[0];
      const nodeName =
        this.selectedNode.ui_name ||
        this.$t("common.node") + ` ${this.selectedNode.id}`;

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            repository: this.selectedInstance.repository_id,
            path: this.selectedInstance.path,
            snapshot: this.selectedSnapshot.id,
            node: this.selectedNode.id,
            replace: this.replaceExistingApp,
          },
          extra: {
            title: this.$t("action." + taskAction),
            description: this.$t("backup.restoring_app_to_node", {
              app: app,
              node: nodeName,
            }),
            node: nodeName,
            completion: {
              i18nString: "backup.instance_restored_to_node",
              extraTextParams: ["node"],
              outputTextParams: ["module_id"],
            },
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.restoreModule = this.getErrorMessage(err);
        return;
      }

      // close modal immediately, no validation needed
      this.$emit("hide");
    },
    restoreModuleAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.restoreModule = this.$t("error.generic_error");
      this.loading.restoreModule = false;
      this.$emit("hide");
    },
    restoreModuleCompleted() {
      this.loading.restoreModule = false;

      // show restored instance in app drawer
      this.$root.$emit("reloadAppDrawer");
    },
    onSelectNode(selectedNode) {
      this.selectedNode = selectedNode;
    },
    onSelectInstance(selectedInstance) {
      this.selectedInstance = selectedInstance;
      this.replaceExistingApp = false;
    },
    onSelectSnapshot(selectedSnapshot) {
      this.selectedSnapshot = selectedSnapshot;
    },
    async determineRestoreEligibility() {
      this.loading.determineRestoreEligibility = true;
      this.error.determineRestoreEligibility = "";
      this.selectedNode = null;
      const taskAction = "determine-restore-eligibility";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.determineRestoreEligibilityAborted
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.determineRestoreEligibilityCompleted
      );

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            repository: this.selectedInstance.repository_id,
            path: this.selectedInstance.path,
            snapshot: this.selectedSnapshot.id,
          },
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
        this.error.determineRestoreEligibility = this.getErrorMessage(err);
        return;
      }
    },
    determineRestoreEligibilityAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.determineRestoreEligibility = this.$t("error.generic_error");
      this.loading.determineRestoreEligibility = false;
    },
    determineRestoreEligibilityCompleted(taskContext, taskResult) {
      this.installDestinations = taskResult.output.install_destinations;
      this.loading.determineRestoreEligibility = false;
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";

.instances {
  margin-top: $spacing-07;
}
</style>

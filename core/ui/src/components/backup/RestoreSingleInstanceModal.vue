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
    :nextLabel="
      isLastStep || shouldShowRestoreLabel
        ? $t('backup.restore')
        : $t('common.next')
    "
    :isPreviousDisabled="isFirstStep || loading.restoreModule"
    :isNextDisabled="isNextButtonDisabled"
    @modal-hidden="$emit('hide')"
    @modal-shown="onModalShown"
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
                <NsCheckbox
                  :label="
                    $t('backup.replace_existing_app', {
                      app: instanceToReplace,
                    })
                  "
                  v-model="replaceExistingApp"
                  :disabled="!instanceToReplace || replaceExistingDisabled"
                  value="checkReplaceExistingApp"
                  class="mg-top-xlg"
                  tooltipAlignment="start"
                  tooltipDirection="right"
                >
                  <template slot="tooltip">
                    <div
                      v-if="
                        selectedInstance &&
                        instanceToReplace &&
                        replaceExistingDisabled
                      "
                    >
                      {{
                        $t(
                          "backup.replace_existing_app_is_not_available_tooltip"
                        )
                      }}
                    </div>
                    <div v-else-if="!selectedInstance">
                      {{
                        $t("backup.replace_existing_app_select_an_app_tooltip")
                      }}
                    </div>
                    <div v-else>
                      {{ $t("backup.replace_existing_app_tooltip") }}
                    </div>
                  </template>
                </NsCheckbox>
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
            v-if="error.getClusterStatus"
            kind="error"
            :title="$t('action.get-cluster-status')"
            :description="error.getClusterStatus"
            :showCloseButton="false"
          />
          <NsInlineNotification
            v-if="error.listModules"
            kind="error"
            :title="$t('action.list-modules')"
            :description="error.listModules"
            :showCloseButton="false"
          />
          <NsInlineNotification
            v-if="error.nodesList"
            kind="error"
            :title="$t('action.list-nodes')"
            :description="error.nodesList"
            :showCloseButton="false"
          />
          <NsInlineNotification
            v-if="
              clusterNodesCount == disabledNodes.length && !isLoadingNodeData
            "
            kind="info"
            :title="$t('backup.no_node_eligible_for_instance_restoration')"
            :showCloseButton="false"
          />
          <NodeSelector
            :nodesWithAdditionalStorage="nodesWithAdditionalStorage"
            @selectNode="onSelectNode"
            :disabledNodes="disabledNodes"
            :loading="isLoadingNodeData"
            class="mg-top-xlg"
          >
            <template v-for="(nodeInfoMessage, nodeId) in nodesInfo">
              <template :slot="`node-${nodeId}`">
                <template v-if="Array.isArray(nodeInfoMessage)">
                  {{ nodeInfoMessage.join(", ") }}
                </template>
                <template v-else>
                  {{ nodeInfoMessage }}
                </template>
              </template>
            </template>
          </NodeSelector>
        </template>
        <template v-else-if="step == 'volumes' && selectedNode">
          <NsInlineNotification
            v-if="error.listMountPoints"
            kind="error"
            :title="$t('action.list-mountpoints')"
            :description="error.listMountPoints"
            :showCloseButton="false"
          />
          <div>
            {{
              $t("software_center.select_node_volume_for_installation", {
                node: this.getNodeLabel(selectedNode) || "",
              })
            }}
          </div>
          <!-- additional volumes -->
          <AdditionalVolumesSelector
            @selectVolume="onSelectVolume"
            :volumes="additionalVolumes"
            :loading="loading.listMountPoints"
            :light="true"
            class="mg-top-lg"
          />
        </template>
        <NsInlineNotification
          v-if="error.restoreModule"
          kind="error"
          :title="$t('action.restore-module')"
          :description="error.restoreModule"
          :showCloseButton="false"
        />
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
import AdditionalVolumesSelector from "@/components/software-center/AdditionalVolumesSelector.vue";

export default {
  name: "RestoreSingleInstanceModal",
  components: {
    NodeSelector,
    RestoreSingleInstanceSelector,
    RestoreSingleInstanceSnapshotSelector,
    AdditionalVolumesSelector,
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
      steps: ["instance", "snapshot", "node", "volumes"],
      instances: [],
      backupSnapshots: [],
      selectedNode: null,
      selectedInstance: null,
      selectedSnapshot: null,
      replaceExistingApp: false,
      installDestinations: [],
      clusterStatus: [],
      additionalVolumes: [],
      nodesList: [],
      selectedVolume: {},
      modules: [],
      loading: {
        readBackupRepositories: true,
        restoreModule: false,
        readBackupSnapshots: false,
        determineRestoreEligibility: false,
        getClusterStatus: true,
        listMountPoints: false,
        listModules: true,
        nodesList: true,
      },
      error: {
        readBackupRepositories: "",
        restoreModule: "",
        readBackupSnapshots: "",
        determineRestoreEligibility: "",
        getClusterStatus: "",
        listMountPoints: "",
        listModules: "",
      },
    };
  },
  computed: {
    ...mapState(["clusterNodes"]),
    isLoadingNodeData() {
      return (
        this.loading.determineRestoreEligibility ||
        this.loading.nodesList ||
        this.loading.listModules ||
        this.loading.getClusterStatus
      );
    },
    clusterNodesCount() {
      return this.clusterNodes.length;
    },
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
          (this.isLoadingNodeData ||
            !this.selectedNode ||
            (this.clusterNodesCount == this.disabledNodes.length &&
              !this.shouldShowRestoreLabel))) ||
        (this.step == "volumes" &&
          (this.loading.listMountPoints ||
            !this.selectedNode ||
            !this.selectedVolume ||
            Object.keys(this.selectedVolume).length === 0))
      );
    },
    selectedModule() {
      if (
        !this.selectedInstance ||
        !this.modules ||
        this.modules.length === 0
      ) {
        return null;
      }
      return this.modules.find(
        (module) => module.id === this.selectedInstance.name
      );
    },
    appVolumes() {
      if (
        this.selectedModule &&
        this.selectedModule.versions &&
        this.selectedModule.versions.length &&
        this.selectedModule.versions[0].labels &&
        this.selectedModule.versions[0].labels["org.nethserver.volumes"]
      ) {
        const volumesString =
          this.selectedModule.versions[0].labels["org.nethserver.volumes"];
        return volumesString.split(" ").map((v) => v.trim());
      }
      return [];
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
    replaceExistingDisabled() {
      // FIXME: remove hardcoded values due to the static port, the API should give us this info
      //check if the selected instance is not in an array ['loki']
      const notAllowed = [
        "samba",
        "loki",
        "mail",
        "ejabberd",
        "nethvoice-proxy",
      ];
      return notAllowed.includes(this.selectedInstance.name) ? true : false;
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
      // Add offline nodes message
      for (const node of this.clusterStatus) {
        if (!node.online) {
          nodesInfo[node.id] = this.$t("software_center.node_offline");
        }
      }
      return nodesInfo;
    },
    disabledNodes() {
      // Get non-eligible nodes from install_destinations
      const nonEligibleNodeIds = this.installDestinations
        .filter((nodeInfo) => !nodeInfo.eligible)
        .map((nodeInfo) => nodeInfo.node_id);
      // Get offline nodes from clusterNodes
      const offlineNodeIds = this.clusterStatus
        .filter((node) => !node.online)
        .map((node) => node.id);
      // Combine and remove duplicates
      return [...new Set([...nonEligibleNodeIds, ...offlineNodeIds])];
    },
    nodesWithAdditionalStorage() {
      const nodeIds = [];
      // If app doesn't require volumes, no nodes have additional storage for it
      if (this.appVolumes.length === 0) {
        return nodeIds;
      }
      if (this.nodesList && Array.isArray(this.nodesList)) {
        for (const node of this.nodesList) {
          if (node.additional_disk_count > 0) {
            nodeIds.push(node.node_id);
          }
        }
      }
      return nodeIds;
    },
    shouldShowRestoreLabel() {
      return (
        this.step == "node" &&
        this.selectedNode &&
        (!this.nodesWithAdditionalStorage.includes(this.selectedNode.id) ||
          this.appVolumes.length === 0)
      );
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
        this.listNodes();
        this.listModules();
        this.getClusterStatus();
      } else if (this.step == "volumes") {
        this.listMountPoints();
      }
    },
  },
  methods: {
    nextStep() {
      if (this.isNextButtonDisabled) {
        return;
      }

      if (this.isLastStep) {
        this.restoreModule(this.selectedVolume);
      } else {
        // Check if we should skip the volumes step
        if (
          this.step == "node" &&
          (!this.nodesWithAdditionalStorage.includes(this.selectedNode.id) ||
            this.appVolumes.length === 0)
        ) {
          // Skip volumes step, go directly to installation
          this.restoreModule();
        } else {
          this.step = this.steps[this.stepIndex + 1];
        }
      }
    },
    previousStep() {
      if (!this.isFirstStep) {
        this.step = this.steps[this.stepIndex - 1];
      }
    },
    onModalShown() {
      // reset state before showing modal
      this.clearErrors();
      this.clusterStatus = [];
      this.nodesList = [];
      this.selectedVolume = {};
      this.instances = [];
      this.backupSnapshots = [];
      this.selectedInstance = null;
      this.selectedSnapshot = null;
      this.replaceExistingApp = false;
      this.installDestinations = [];
      this.additionalVolumes = [];
      if (this.clusterNodesCount == 1) {
        const firstNode = this.clusterNodes[0];
        this.selectedNode = { ...firstNode, selected: true };
      } else {
        this.selectedNode = null;
      }
    },
    async listModules() {
      this.loading.listModules = true;
      this.error.listModules = "";
      this.modules = [];
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
      this.loading.listModules = false;
      let modules = taskResult.output;
      this.modules = modules;
      this.loading.listModules = false;
    },
    async listNodes() {
      this.error.nodesList = "";
      this.loading.nodesList = true;
      const taskAction = "list-nodes";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.listNodesAborted
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.listNodesCompleted
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
        this.error.nodesList = this.getErrorMessage(err);
        return;
      }
    },
    listNodesAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.nodesList = this.$t("error.generic_error");
      this.loading.nodesList = false;
    },
    listNodesCompleted(taskContext, taskResult) {
      this.nodesList = taskResult.output.nodes;
      this.loading.nodesList = false;
    },
    async listMountPoints() {
      this.error.listMountPoints = "";
      this.loading.listMountPoints = true;
      const taskAction = "list-mountpoints";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.listMountPointsAborted
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.listMountPointsCompleted
      );

      const res = await to(
        this.createNodeTask(this.selectedNode.id, {
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
        this.error.listMountPoints = this.getErrorMessage(err);
        this.loading.listMountPoints = false;
        return;
      }
    },
    listMountPointsAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.listMountPoints = this.$t("error.generic_error");
      this.loading.listMountPoints = false;
    },
    listMountPointsCompleted(taskContext, taskResult) {
      const additionalVolumes = taskResult.output.mountpoints;
      // Add default disk at the end, push default property
      if (taskResult.output.default_disk) {
        additionalVolumes.push({
          ...taskResult.output.default_disk,
          default: true, // mark as default disk
        });
      }
      this.additionalVolumes = additionalVolumes;
      this.loading.listMountPoints = false;
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
    async restoreModule(volumes) {
      this.loading.restoreModule = true;
      this.error.restoreModule = "";
      const taskAction = "restore-module";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.restoreModuleAborted
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.restoreModuleCompleted
      );

      const app = this.selectedInstance.path.split("/")[0];
      const nodeName =
        this.selectedNode.ui_name ||
        this.$t("common.node") + ` ${this.selectedNode.id}`;

      const data = {
        repository: this.selectedInstance.repository_id,
        path: this.selectedInstance.path,
        snapshot: this.selectedSnapshot.id,
        node: this.selectedNode.id,
        replace: this.replaceExistingApp,
      };

      // Add volumes path if any and if app requires volumes
      if (volumes?.path && this.appVolumes.length) {
        data.volumes = {};
        this.appVolumes.forEach((volume) => {
          data.volumes[volume] = volumes.path;
        });
      }

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data,
          extra: {
            eventId,
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
    onSelectVolume(selectedVolume) {
      this.selectedVolume = selectedVolume;
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
    async getClusterStatus() {
      this.loading.getClusterStatus = true;
      this.error.getClusterStatus = "";
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
        return;
      }
    },
    getClusterStatusAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.getClusterStatus = this.$t("error.generic_error");
      this.loading.getClusterStatus = false;
    },
    getClusterStatusCompleted(taskContext, taskResult) {
      this.clusterStatus = taskResult.output.nodes;
      this.loading.getClusterStatus = false;
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

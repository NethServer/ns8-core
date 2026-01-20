<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <NsWizard
    size="default"
    :visible="isShown"
    @primary-click="cloneOrMove"
    @modal-shown="onModalShown"
    class="no-pad-modal"
    :cancelLabel="$t('common.cancel')"
    :previousLabel="$t('common.previous')"
    :nextLabel="
      isLastStep || shouldShowCloneOrMoveLabel
        ? isClone
          ? $t('software_center.clone_app')
          : $t('software_center.move_app')
        : $t('common.next')
    "
    :isPreviousDisabled="isFirstStep || loading.getClusterStatus"
    @previousStep="previousStep"
    @nextStep="nextStep"
    @cancel="onModalHidden"
    :isNextDisabled="isNextButtonDisabled"
    :isNextLoading="loading.getClusterStatus || loading.listMountPoints"
    :isPreviousShown="hasAdditionalStorageAvailable"
  >
    <template slot="title">{{
      isClone ? $t("software_center.clone_app") : $t("software_center.move_app")
    }}</template>
    <template slot="content">
      <cv-form @submit.prevent="nextStep">
        <template v-if="step == 'node'">
          <div class="mg-bottom-md">
            {{
              isClone
                ? $t("software_center.clone_app_description", {
                    instanceLabel,
                  })
                : $t("software_center.move_app_description", { instanceLabel })
            }}
          </div>
          <NsInlineNotification
            v-if="error.getClusterStatus"
            kind="error"
            :title="$t('action.get-cluster-status')"
            :description="error.getClusterStatus"
            :showCloseButton="false"
          />
          <NsInlineNotification
            v-if="error.listNodes"
            kind="error"
            :title="$t('action.list-nodes')"
            :description="error.listNodes"
            :showCloseButton="false"
          />
          <div>{{ $t("software_center.select_destination_node") }}:</div>
          <NsInlineNotification
            v-if="
              clusterStatus.length == disabledNodes.length &&
              !loading.getClusterStatus
            "
            kind="info"
            :title="
              isClone
                ? $t('software_center.no_node_eligible_for_instance_cloning')
                : $t('software_center.no_node_eligible_for_instance_migration')
            "
            :showCloseButton="false"
          />
          <NodeSelector
            class="mg-top-xlg"
            @selectNode="onSelectNode"
            :disabledNodes="disabledNodes"
            :loading="loading.getClusterStatus"
            :nodesWithAdditionalStorage="nodesWithAdditionalStorage"
          >
            <template v-for="(nodeMessages, nodeId) in nodesInfo">
              <template :slot="`node-${nodeId}`">
                <div
                  v-for="(nodeMessage, index) in nodeMessages"
                  :key="index"
                  class="node-message"
                >
                  {{ nodeMessage }}
                </div>
              </template>
            </template>
          </NodeSelector>
        </template>
        <template v-else-if="step == 'volumes' && selectedNode">
          <div>
            {{
              $t("software_center.select_node_volume_for_installation", {
                node: selectedNode.ui_name
                  ? selectedNode.ui_name +
                    " (" +
                    $t("common.node") +
                    ` ${selectedNode.id})`
                  : ` ${selectedNode.id}`,
              })
            }}
          </div>
          <NsInlineNotification
            v-if="error.listMountPoints"
            kind="error"
            :title="$t('action.list-mountpoints')"
            :description="error.listMountPoints"
            :showCloseButton="false"
          />
          <!-- additonal volumes -->
          <AdditionnalVolumesSelector
            @selectVolume="onSelectVolume"
            :volumes="additionnalVolumes"
            :loading="loading.listMountPoints"
            :light="true"
            class="mg-top-lg"
          />
        </template>
        <!-- clone-module error -->
        <NsInlineNotification
          v-if="error.cloneModule"
          kind="error"
          :title="$t('action.clone-module')"
          :description="error.cloneModule"
          :showCloseButton="false"
        />
      </cv-form>
    </template>
  </NsWizard>
</template>

<script>
import { UtilService, IconService, TaskService } from "@nethserver/ns8-ui-lib";
import { mapState } from "vuex";
import NodeSelector from "@/components/nodes/NodeSelector";
import AdditionnalVolumesSelector from "@/components/software-center/AdditionnalVolumesSelector.vue";
import to from "await-to-js";

export default {
  name: "CloneOrMoveAppModal",
  components: { NodeSelector, AdditionnalVolumesSelector },
  mixins: [UtilService, IconService, TaskService],
  props: {
    isShown: {
      type: Boolean,
      default: false,
    },
    isClone: {
      type: Boolean,
      required: true,
    },
    instanceId: {
      type: String,
      required: true,
    },
    instanceUiName: {
      type: String,
      default: "",
    },
    installationNode: {
      type: Number,
      required: true,
    },
    app: {
      type: Object,
      default: null,
    },
  },
  data() {
    return {
      selectedNode: null,
      step: "node",
      steps: ["node", "volumes"],
      clusterStatus: [],
      listNodes: [],
      additionnalVolumes: [],
      selectedVolume: {},
      loading: {
        cloneModule: false,
        getClusterStatus: true,
        listMountPoints: false,
      },
      error: {
        cloneModule: "",
        getClusterStatus: "",
        listNodes: "",
        listMountPoints: "",
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
    shouldShowCloneOrMoveLabel() {
      return (
        this.step == "node" &&
        this.selectedNode &&
        !this.nodesWithAdditionalStorage.includes(this.selectedNode.id)
      );
    },
    nodesWithAdditionalStorage() {
      const nodeIds = [];
      if (this.listNodes && Array.isArray(this.listNodes)) {
        for (const node of this.listNodes) {
          if (node.additional_disk_count > 0) {
            nodeIds.push(node.node_id);
          }
        }
      }
      return nodeIds;
    },
    hasAdditionalStorageAvailable() {
      return this.nodesWithAdditionalStorage.length > 0;
    },
    isNextButtonDisabled() {
      if (this.step == "node") {
        return (
          !this.selectedNode ||
          (!this.isClone && this.installationNode == this.selectedNode.id) ||
          this.loading.cloneModule ||
          this.loading.getClusterStatus ||
          this.clusterStatus.length == this.disabledNodes.length
        );
      } else if (this.step == "volumes") {
        return (
          this.loading.listMountPoints ||
          !this.selectedNode ||
          !this.selectedVolume ||
          Object.keys(this.selectedVolume).length === 0
        );
      }
      return false;
    },
    instanceLabel() {
      if (this.instanceUiName) {
        return `${this.instanceUiName} (${this.instanceId})`;
      } else {
        return this.instanceId;
      }
    },
    nodesInfo() {
      if (!this.app) {
        return {};
      }
      const nodesInfo = {};

      if (this.app) {
        for (const nodeInfo of this.app.install_destinations) {
          const nodeMessages = [];

          if (nodeInfo.node_id === this.installationNode) {
            // show "Current node"
            nodeMessages.push(this.$t("software_center.current_node"));
          }

          if (!nodeInfo.eligible) {
            // show reason why node is not eligible
            const rejectReason = nodeInfo.reject_reason;

            if (rejectReason.message === "max_per_node_limit") {
              const numMaxInstances = rejectReason.parameter;
              nodeMessages.push(
                this.$tc(
                  `software_center.reason_${rejectReason.message}`,
                  numMaxInstances,
                  { param: numMaxInstances },
                ),
              );
            } else {
              nodeMessages.push(
                this.$t(`software_center.reason_${rejectReason.message}`, {
                  param: rejectReason.parameter,
                }),
              );
            }
          } else if (nodeInfo.instances) {
            // show number of instances installed
            nodeMessages.push(
              this.$tc(
                "software_center.num_instances_installed",
                nodeInfo.instances,
                { num: nodeInfo.instances },
              ),
            );
          }
          nodesInfo[nodeInfo.node_id] = nodeMessages;
        }
        // Add offline nodes message
        for (const node of this.clusterStatus) {
          if (!node.online) {
            if (!nodesInfo[node.id]) {
              nodesInfo[node.id] = [];
            }
            nodesInfo[node.id].push(this.$t("software_center.node_offline"));
          }
        }
      }
      return nodesInfo;
    },
    disabledNodes() {
      if (!this.app) {
        return [];
      }

      // Get non-eligible nodes from install_destinations
      const ineligibleNodes = this.app.install_destinations
        .filter((nodeInfo) => !nodeInfo.eligible)
        .map((nodeInfo) => nodeInfo.node_id);

      // Get offline nodes from clusterStatus
      const offlineNodeIds = this.clusterStatus
        .filter((node) => !node.online)
        .map((node) => node.id);

      if (this.isClone) {
        // cloning app: combine ineligible and offline nodes
        return [...new Set([...ineligibleNodes, ...offlineNodeIds])];
      } else {
        // moving app: add current node to ineligible/offline nodes and remove possible duplicates
        return [
          ...new Set([
            ...ineligibleNodes.concat(this.installationNode),
            ...offlineNodeIds,
          ]),
        ];
      }
    },
    appVolumes() {
      if (
        this.app &&
        this.app.versions &&
        this.app.versions.length &&
        this.app.versions[0].labels &&
        this.app.versions[0].labels["org.nethserver.volumes"]
      ) {
        const volumesString =
          this.app.versions[0].labels["org.nethserver.volumes"];
        return volumesString.split(" ").map((v) => v.trim());
      }
      return [];
    },
  },
  watch: {
    step: function () {
      if (this.step == "volumes") {
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
        this.cloneOrMove(this.selectedVolume);
      } else {
        // Check if we should skip the volumes step
        if (
          this.step == "node" &&
          !this.nodesWithAdditionalStorage.includes(this.selectedNode.id)
        ) {
          // Skip volumes step, go directly to installation
          this.cloneOrMove();
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
      // Force selection to node 1 if only available
      if (this.clusterNodes.length == 1) {
        this.selectedNode = this.clusterNodes[0];
        this.clusterNodes[0].selected = true;
      } else {
        this.selectedNode = null;
      }
      this.clusterStatus = [];
      this.fetchListNodes();
    },
    onModalHidden() {
      // reset state
      this.$emit("hide");
      this.clearErrors();
      // reset state
      this.selectedNode = null;
      this.selectedVolume = {};
      this.clusterStatus = [];
      this.listNodes = [];
      this.additionnalVolumes = [];
      this.step = this.steps[0];
      this.loading.getClusterStatus = true;
    },
    async fetchListNodes() {
      this.error.listNodes = "";
      const taskAction = "list-nodes";

      // register to task error
      this.$root.$off(taskAction + "-aborted");
      this.$root.$once(taskAction + "-aborted", this.listNodesAborted);

      // register to task completion
      this.$root.$off(taskAction + "-completed");
      this.$root.$once(taskAction + "-completed", this.ListNodesCompleted);

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
          },
        }),
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.listNodes = this.getErrorMessage(err);
        return;
      }
    },
    listNodesAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.listNodes = this.$t("error.generic_error");
      this.loading.getClusterStatus = false;
    },
    ListNodesCompleted(taskContext, taskResult) {
      this.listNodes = taskResult.output.nodes;
      this.getClusterStatus();
    },
    async getClusterStatus() {
      this.error.getClusterStatus = "";
      const taskAction = "get-cluster-status";

      // register to task error
      this.$root.$off(taskAction + "-aborted");
      this.$root.$once(taskAction + "-aborted", this.getClusterStatusAborted);

      // register to task completion
      this.$root.$off(taskAction + "-completed");
      this.$root.$once(
        taskAction + "-completed",
        this.getClusterStatusCompleted,
      );

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
          },
        }),
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
    async listMountPoints() {
      this.error.listMountPoints = "";
      this.loading.listMountPoints = true;
      const taskAction = "list-mountpoints";

      // register to task error
      this.$root.$off(taskAction + "-aborted");
      this.$root.$once(taskAction + "-aborted", this.listMountPointsAborted);

      // register to task completion
      this.$root.$off(taskAction + "-completed");
      this.$root.$once(
        taskAction + "-completed",
        this.listMountPointsCompleted,
      );
      const res = await to(
        this.createNodeTask(this.selectedNode.id, {
          action: taskAction,
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
          },
        }),
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.listNodes = this.getErrorMessage(err);
        return;
      }
    },
    listMountPointsAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.listMountPoints = this.$t("error.generic_error");
      this.loading.listMountPoints = false;
    },
    listMountPointsCompleted(taskContext, taskResult) {
      this.additionnalVolumes = taskResult.output.mountpoints;
      // Add default disk at the end, push default property
      if (taskResult.output.default_disk) {
        this.additionnalVolumes.push({
          ...taskResult.output.default_disk,
          default: true, // mark as default disk
        });
      }
      this.loading.listMountPoints = false;
    },
    async cloneOrMove(volumes) {
      this.loading.cloneModule = true;
      const taskAction = "clone-module";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.cloneModuleAborted,
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.cloneModuleCompleted,
      );

      const shortInstanceLabel = this.instanceUiName || this.instanceId;
      const shortNodeLabel =
        this.selectedNode.ui_name ||
        this.$t("common.node") + ` ${this.selectedNode.id}`;

      const notificationTitle = this.isClone
        ? this.$t("software_center.clone_app")
        : this.$t("software_center.move_app");

      const notificationDescription = this.isClone
        ? this.$t("software_center.cloning_app_to_node", {
            instance: shortInstanceLabel,
            node: shortNodeLabel,
          })
        : this.$t("software_center.moving_app_to_node", {
            instance: shortInstanceLabel,
            node: shortNodeLabel,
          });

      let completionString = "";

      if (this.instanceUiName) {
        if (this.isClone) {
          completionString = "software_center.app_with_ui_name_cloned_to_node";
        } else {
          completionString = "software_center.app_with_ui_name_moved_to_node";
        }
      } else {
        if (this.isClone) {
          completionString =
            "software_center.app_without_ui_name_cloned_to_node";
        } else {
          completionString =
            "software_center.app_without_ui_name_moved_to_node";
        }
      }
      const data = {
        module: this.instanceId,
        node: parseInt(this.selectedNode.id),
        replace: !this.isClone,
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
            title: notificationTitle,
            description: notificationDescription,
            instance: shortInstanceLabel,
            node: shortNodeLabel,
            completion: {
              i18nString: completionString,
              extraTextParams: ["node", "instance"],
              outputTextParams: ["module_id"],
            },
            eventId,
          },
        }),
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.cloneModule = this.getErrorMessage(err);
        this.loading.cloneModule = false;
        return;
      }
      this.loading.cloneModule = false;

      // hide modal
      this.$emit("hide");
    },
    cloneModuleAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.loading.cloneModule = false;
    },
    cloneModuleCompleted(taskContext, taskResult) {
      const newModuleId = taskResult.output.module_id;
      this.clearErrors();
      // reset state
      this.selectedNode = null;
      this.selectedVolume = {};
      this.clusterStatus = [];
      this.listNodes = [];
      this.additionnalVolumes = [];
      this.step = this.steps[0];
      this.loading.getClusterStatus = true;

      // reload instances and highlight cloned/moved instance
      this.$emit("cloneOrMoveCompleted", newModuleId);
    },
    onSelectNode(selectedNode) {
      this.selectedNode = selectedNode;
    },
    onSelectVolume(selectedVolume) {
      this.selectedVolume = selectedVolume;
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";

.node-message {
  margin-bottom: $spacing-05;
}

.node-message:last-child {
  margin-bottom: 0;
}
</style>

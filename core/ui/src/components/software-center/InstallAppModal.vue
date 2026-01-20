<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <NsWizard
    size="default"
    :visible="isShown"
    @modal-hidden="onModalHidden"
    @modal-shown="onModalShown"
    class="no-pad-modal"
    :cancelLabel="$t('common.cancel')"
    :previousLabel="$t('common.previous')"
    :nextLabel="
      isLastStep || shouldShowInstallLabel
        ? $t('software_center.install')
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
    <template v-if="app" slot="title">{{
      $t("software_center.app_installation", { app: app.name })
    }}</template>

    <template v-if="app" slot="content">
      <cv-form @submit.prevent="nextStep">
        <template v-if="step == 'node'">
          <!-- warning for rootfull app -->
          <NsInlineNotification
            v-if="app.rootfull && app.certification_level < 3"
            kind="warning"
            :title="$t('software_center.rootfull_app_warning_title')"
            :description="
              $t('software_center.rootfull_app_warning_description', {
                appName: app.name,
              })
            "
            :showCloseButton="false"
          />
          <!-- node selection -->
          <div>
            {{
              $t("software_center.choose_node_for_installation", {
                app: app.name,
                version: appVersion,
              })
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
          <NsInlineNotification
            v-if="
              clusterStatus.length == disabledNodes.length &&
              !loading.getClusterStatus
            "
            kind="info"
            :title="$t('software_center.no_node_eligible_for_app_installation')"
            :showCloseButton="false"
          />
          <NodeSelector
            @selectNode="onSelectNode"
            :disabledNodes="disabledNodes"
            :nodesWithAdditionalStorage="nodesWithAdditionalStorage"
            :loading="loading.getClusterStatus"
            class="mg-top-lg"
          >
            <template v-for="(nodeInfoMessage, nodeId) in nodesInfo">
              <template :slot="`node-${nodeId}`">
                {{ nodeInfoMessage }}
              </template>
            </template>
          </NodeSelector>
          <!-- terms and conditions -->
          <NsCheckbox
            v-if="app.docs.terms_url"
            v-model="agreeTerms"
            value="checkTerms"
          >
            <template slot="label">
              <i18n
                path="software_center.agree_terms_before_install"
                tag="span"
              >
                <template v-slot:appName>
                  {{ app.name }}
                </template>
                <template v-slot:terms>
                  <cv-link
                    :href="app.docs.terms_url"
                    target="_blank"
                    rel="noreferrer"
                    class="inline"
                  >
                    {{ $t("common.terms_and_conditions") }}
                  </cv-link>
                </template>
              </i18n>
            </template>
          </NsCheckbox>
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
        <!-- add-module error -->
        <NsInlineNotification
          v-if="error.addModule"
          kind="error"
          :title="$t('action.add-module')"
          :description="error.addModule"
          :showCloseButton="false"
        />
      </cv-form>
    </template>
  </NsWizard>
</template>

<script>
import { UtilService, TaskService, IconService } from "@nethserver/ns8-ui-lib";
import to from "await-to-js";
import NotificationService from "@/mixins/notification";
import NodeSelector from "@/components/nodes/NodeSelector";
import AdditionnalVolumesSelector from "@/components/software-center/AdditionnalVolumesSelector.vue";
import { mapState } from "vuex";

export default {
  name: "InstallAppModal",
  components: { NodeSelector, AdditionnalVolumesSelector },
  mixins: [UtilService, TaskService, IconService, NotificationService],
  props: {
    isShown: Boolean,
    app: { type: [Object, null] },
  },
  data() {
    return {
      selectedNode: null,
      agreeTerms: false,
      step: "node",
      steps: ["node", "volumes"],
      loading: {
        getClusterStatus: true,
        listMountPoints: false,
      },
      clusterStatus: [],
      listNodes: [],
      selectedVolume: {},
      additionnalVolumes: [],
      error: {
        addModule: "",
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
    shouldShowInstallLabel() {
      return (
        this.step == "node" &&
        this.selectedNode &&
        !this.nodesWithAdditionalStorage.includes(this.selectedNode.id)
      );
    },
    isNextButtonDisabled() {
      if (this.step == "node") {
        return (
          this.loading.getClusterStatus ||
          !this.selectedNode ||
          this.clusterStatus.length == this.disabledNodes.length ||
          (this.app && this.app.docs.terms_url && !this.agreeTerms)
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
    firstNodeLabel() {
      if (this.clusterNodes.length && this.clusterNodes[0]) {
        if (this.clusterNodes[0].ui_name) {
          return `${this.clusterNodes[0].ui_name} (${this.$t("common.node")} ${
            this.clusterNodes[0].id
          })`;
        } else {
          return `${this.$t("common.node")} ${this.clusterNodes[0].id}`;
        }
      } else {
        return "";
      }
    },
    appVersion() {
      if (this.app.versions.length) {
        return this.app.versions[0].tag;
      } else {
        return "-"; // this value is never sent. It would fail the backend validation
      }
    },
    canInstallOnSingleNode() {
      const installDestination = this.app.install_destinations.find(
        (nodeInfo) => nodeInfo.node_id == this.clusterNodes[0].id,
      );
      return installDestination?.eligible;
    },
    nodesInfo() {
      const nodesInfo = {};

      if (this.app) {
        for (const nodeInfo of this.app.install_destinations) {
          if (!nodeInfo.eligible) {
            // show reason why node is not eligible
            const rejectReason = nodeInfo.reject_reason;

            if (rejectReason.message === "max_per_node_limit") {
              const numMaxInstances = rejectReason.parameter;
              nodesInfo[nodeInfo.node_id] = this.$tc(
                `software_center.reason_${rejectReason.message}`,
                numMaxInstances,
                { param: numMaxInstances },
              );
            } else {
              nodesInfo[nodeInfo.node_id] = this.$t(
                `software_center.reason_${rejectReason.message}`,
                { param: rejectReason.parameter },
              );
            }
          } else if (nodeInfo.instances) {
            // show number of instances installed
            nodesInfo[nodeInfo.node_id] = this.$tc(
              "software_center.num_instances_installed",
              nodeInfo.instances,
              { num: nodeInfo.instances },
            );
          }
        }
        // Add offline nodes message
        for (const node of this.clusterStatus) {
          if (!node.online) {
            nodesInfo[node.id] = this.$t("software_center.node_offline");
          }
        }
      }
      return nodesInfo;
    },
    disabledNodes() {
      // Get non-eligible nodes from install_destinations
      const nonEligibleNodeIds = this.app.install_destinations
        .filter((nodeInfo) => !nodeInfo.eligible)
        .map((nodeInfo) => nodeInfo.node_id);
      // Get offline nodes from clusterStatus
      const offlineNodeIds = this.clusterStatus
        .filter((node) => !node.online)
        .map((node) => node.id);
      // Combine and remove duplicates
      return [...new Set([...nonEligibleNodeIds, ...offlineNodeIds])];
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
        this.installInstance(this.selectedVolume);
      } else {
        // Check if we should skip the volumes step
        if (
          this.step == "node" &&
          !this.nodesWithAdditionalStorage.includes(this.selectedNode.id)
        ) {
          // Skip volumes step, go directly to installation
          this.installInstance();
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
      this.agreeTerms = false;
      // reset state before showing modal
      this.clearErrors();
      this.clusterStatus = [];
      this.listNodes = [];
      this.selectedVolume = {};
      // Force selection to node 1 if only available
      if (this.clusterNodes.length == 1 && this.canInstallOnSingleNode) {
        this.selectedNode = this.clusterNodes[0];
        this.clusterNodes[0].selected = true;
      } else {
        this.selectedNode = null;
      }
      this.fetchListNodes();
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
    async installInstance(volumes) {
      this.error.addModule = "";
      const taskAction = "add-module";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.addModuleAborted,
      );

      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.addModuleCompleted,
      );

      const nodeName =
        this.selectedNode.ui_name ||
        this.$t("common.node") + ` ${this.selectedNode.id}`;

      const data = {
        image: this.app.source + ":" + this.appVersion,
        node: parseInt(this.selectedNode.id),
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
            title: this.$t("software_center.app_installation", {
              app: this.app.name,
            }),
            description: this.$t("software_center.installing_on_node", {
              node: nodeName,
            }),
            node: nodeName,
            eventId,
            completion: {
              i18nString: "software_center.instance_installed_on_node",
              extraTextParams: ["node"],
              outputTextParams: ["module_id"],
            },
          },
        }),
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.addModule = this.getErrorMessage(err);
        return;
      }

      // emit event to close modal
      this.$emit("close");
    },
    addModuleAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
    },
    addModuleCompleted(taskContext, taskResult) {
      const moduleId = taskResult.output.module_id;

      // reload instances and highlight new instance
      this.$emit("installationCompleted", moduleId);

      // show new app in app drawer
      this.$root.$emit("reloadAppDrawer");
      // backup notification

      setTimeout(() => {
        const notification = {
          title: this.$t("backup.schedule_backup"),
          description: this.$t("backup.schedule_backup_for_instance", {
            instance: moduleId,
          }),
          type: "info",
          actionLabel: this.$t("backup.schedule_action"),
          action: {
            type: "changeRoute",
            url: `/backup`,
          },
        };
        this.createNotification(notification);
      }, 15000);
    },
    onModalHidden() {
      this.$emit("close");
      this.clearErrors(this);
      // reset state
      this.selectedNode = null;
      this.selectedVolume = {};
      this.clusterStatus = [];
      this.listNodes = [];
      this.additionnalVolumes = [];
      this.step = this.steps[0];
      this.loading.getClusterStatus = true;
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
</style>

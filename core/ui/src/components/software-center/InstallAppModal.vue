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
            :selectedNode="2"
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
          <!-- add-module error -->
          <NsInlineNotification
            v-if="error.addModule"
            kind="error"
            :title="$t('action.add-module')"
            :description="error.addModule"
            :showCloseButton="false"
          />
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
          <!-- add-module error -->
          <NsInlineNotification
            v-if="error.addModule"
            kind="error"
            :title="$t('action.add-module')"
            :description="error.addModule"
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
        restoreModule: false,
        determineRestoreEligibility: false,
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
        (nodeInfo) => nodeInfo.node_id == this.clusterNodes[0].id
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
      if (this.listNodes && Array.isArray(this.listNodes.nodes)) {
        for (const node of this.listNodes.nodes) {
          if (node.additionnal_storage === true) {
            nodeIds.push(node.node_id);
          }
        }
      }
      return nodeIds;
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
        })
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
    // ListNodesCompleted(taskResult, taskContext) {
    ListNodesCompleted() {
      // this.listNodes = taskResult.output.nodes;
      this.listNodes = {
        nodes: [
          {
            app_count: 0,
            cpu: {
              count: 0,
              family: "",
              microcode: "",
              model: "",
              model_name: "",
              package: "",
              stepping: "",
              usage: 0,
              vendor: "",
            },
            disks: [],
            fqdn: "",
            load: {
              "15min": 0,
              "1min": 0,
              "5min": 0,
            },
            main_ip: "127.0.0.1",
            memory: {
              free: 0,
              total: 0,
              used: 0,
            },
            network_interface_count: 0,
            node_id: 4,
            additionnal_storage: false,
            os_release: {
              name: "",
              version: "",
            },
            role: "worker",
            swap: {
              free: 0,
              total: 0,
              used: 0,
            },
            ui_name: "",
            vpn_endpoint: "",
            vpn_ip_address: "10.5.4.4",
            vpn_listen_port: "",
          },
          {
            app_count: 2,
            cpu: {
              count: 8,
              family: "15",
              microcode: "0x1000065",
              model: "107",
              model_name: "QEMU Virtual CPU version 2.5+",
              package: "0",
              stepping: "1",
              usage: 0.0021164197510120664,
              vendor: "AuthenticAMD",
            },
            disks: [
              {
                device: "/dev/mapper/rl-root",
                free: 52930420736,
                fstype: "xfs",
                mountpoint: "/",
                total: 60747141120,
                used: 7816720384,
              },
              {
                device: "/dev/mapper/rl-root",
                free: 52930420736,
                fstype: "xfs",
                mountpoint: "/var/lib/containers/storage/overlay",
                total: 60747141120,
                used: 7816720384,
              },
              {
                device: "/dev/sda1",
                free: 63864328192,
                fstype: "ext4",
                mountpoint: "/mnt/data",
                total: 67317051392,
                used: 3452723200,
              },
              {
                device: "/dev/sdb1",
                free: 619704320,
                fstype: "xfs",
                mountpoint: "/boot",
                total: 1063256064,
                used: 443551744,
              },
            ],
            fqdn: "R1-pve.rocky9-pve.org",
            load: {
              "15min": 0,
              "1min": 0.01,
              "5min": 0.02,
            },
            main_ip: "192.168.12.110",
            memory: {
              free: 6578237440,
              total: 8052805632,
              used: 1474568192,
            },
            network_interface_count: 2,
            node_id: 1,
            additionnal_storage: true,
            os_release: {
              name: "Rocky Linux",
              version: "9.7",
            },
            role: "leader",
            swap: {
              free: 6874460160,
              total: 6874460160,
              used: 0,
            },
            ui_name: "",
            vpn_endpoint: "R1-pve.rocky9-pve.org:55820",
            vpn_ip_address: "10.5.4.1",
            vpn_listen_port: "55820",
          },
        ],
      };
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
        this.listMountPointsCompleted
      );
      const res = await to(
        this.createNodeTask(this.selectedNode.id, {
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
      console.log("additionnalVolumes", this.additionnalVolumes);
      console.log("selectedNode", this.selectedNode);
      this.loading.listMountPoints = false;
    },
    async installInstance(volumes) {
      this.error.addModule = "";
      const taskAction = "add-module";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.addModuleAborted
      );

      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.addModuleCompleted
      );

      const nodeName =
        this.selectedNode.ui_name ||
        this.$t("common.node") + ` ${this.selectedNode.id}`;

      const data = {
        image: this.app.source + ":" + this.appVersion,
        node: parseInt(this.selectedNode.id),
      };

      // Add volumes if provided
      if (volumes && Object.keys(volumes).length > 0) {
        data.volumes = {
          shares: volumes.path,
        };
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
        })
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

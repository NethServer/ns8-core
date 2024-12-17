<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <NsModal
    size="default"
    :visible="isShown"
    @modal-hidden="$emit('hide')"
    @primary-click="cloneOrMove"
    :primary-button-disabled="
      !selectedNode ||
      (!isClone && installationNode == selectedNode.id) ||
      loading.cloneModule
    "
  >
    <template slot="title">{{
      isClone ? $t("software_center.clone_app") : $t("software_center.move_app")
    }}</template>
    <template slot="content">
      <div class="mg-bottom-md">
        {{
          isClone
            ? $t("software_center.clone_app_description", {
                instanceLabel,
              })
            : $t("software_center.move_app_description", { instanceLabel })
        }}
      </div>
      <div>{{ $t("software_center.select_destination_node") }}:</div>
      <NsInlineNotification
        v-if="clusterNodes.length == disabledNodes.length"
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
      <NsInlineNotification
        v-if="error.cloneModule"
        kind="error"
        :title="$t('action.clone-module')"
        :description="error.cloneModule"
        :showCloseButton="false"
      />
    </template>
    <template slot="secondary-button">{{ $t("common.cancel") }}</template>
    <template slot="primary-button">{{
      isClone ? $t("software_center.clone_app") : $t("software_center.move_app")
    }}</template>
  </NsModal>
</template>

<script>
import { UtilService, IconService, TaskService } from "@nethserver/ns8-ui-lib";
import { mapState } from "vuex";
import NodeSelector from "@/components/nodes/NodeSelector";
import to from "await-to-js";

export default {
  name: "CloneOrMoveAppModal",
  components: { NodeSelector },
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
      loading: {
        cloneModule: false,
      },
      error: {
        cloneModule: "",
      },
    };
  },
  computed: {
    ...mapState(["clusterNodes"]),
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
                  { param: numMaxInstances }
                )
              );
            } else {
              nodeMessages.push(
                this.$t(`software_center.reason_${rejectReason.message}`, {
                  param: rejectReason.parameter,
                })
              );
            }
          } else if (nodeInfo.instances) {
            // show number of instances installed
            nodeMessages.push(
              this.$tc(
                "software_center.num_instances_installed",
                nodeInfo.instances,
                { num: nodeInfo.instances }
              )
            );
          }
          nodesInfo[nodeInfo.node_id] = nodeMessages;
        }
      }
      return nodesInfo;
    },
    disabledNodes() {
      if (!this.app) {
        return [];
      }

      const ineligibleNodes = this.app.install_destinations
        .filter((nodeInfo) => !nodeInfo.eligible)
        .map((nodeInfo) => nodeInfo.node_id);

      if (this.isClone) {
        // cloning app
        return ineligibleNodes;
      } else {
        // moving app: add current node to ineligible nodes and remove possible duplicates
        return [...new Set(ineligibleNodes.concat(this.installationNode))];
      }
    },
  },
  methods: {
    async cloneOrMove() {
      this.loading.cloneModule = true;
      const taskAction = "clone-module";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.cloneModuleAborted
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.cloneModuleCompleted
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

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            module: this.instanceId,
            node: this.selectedNode.id,
            replace: !this.isClone,
          },
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
        })
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

      // reload instances and highlight cloned/moved instance
      this.$emit("cloneOrMoveCompleted", newModuleId);
    },
    onSelectNode(selectedNode) {
      this.selectedNode = selectedNode;
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

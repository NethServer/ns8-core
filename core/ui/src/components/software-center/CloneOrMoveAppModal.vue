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
      <div>{{ $t("software_center.select_destination_node") }}</div>
      <NodeSelector
        class="mg-top-xlg"
        @selectNode="onSelectNode"
        :extraInfoNode="installationNode"
        :extraInfoLabel="$t('software_center.current_node')"
        :disabledNodes="isClone ? [] : [installationNode]"
      />
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
</style>

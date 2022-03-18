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
      isClone
        ? $t("software_center.clone_instance")
        : $t("software_center.move_instance")
    }}</template>
    <template slot="content">
      <div
        class="mg-bottom-md"
        v-html="
          isClone
            ? $t('software_center.clone_instance_description', {
                instanceLabel,
              })
            : $t('software_center.move_instance_description', { instanceLabel })
        "
      ></div>
      <div>{{ $t("software_center.select_destination_node") }}</div>
      <NodeSelector
        class="mg-top-lg"
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
      isClone
        ? $t("software_center.clone_instance")
        : $t("software_center.move_instance")
    }}</template>
  </NsModal>
</template>

<script>
import { UtilService, IconService, TaskService } from "@nethserver/ns8-ui-lib";
import { mapState } from "vuex";
import NodeSelector from "@/components/misc/NodeSelector";
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
    ...mapState(["nodes"]),
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

      // register to task validation
      this.$root.$once(
        `${taskAction}-validation-ok-${eventId}`,
        this.cloneModuleValidationOk
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.cloneModuleCompleted
      );

      const shortInstanceLabel = this.instanceUiName || this.instanceId;
      const shortNodeLabel =
        this.selectedNode.ui_name ||
        this.$t("common.node") + this.selectedNode.id;

      const notificationTitle = this.isClone
        ? this.$t("software_center.clone_instance_to_node", {
            instance: shortInstanceLabel,
            node: shortNodeLabel,
          })
        : this.$t("software_center.move_instance_to_node", {
            instance: shortInstanceLabel,
            node: shortNodeLabel,
          });

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
            description: this.$t("common.processing"),
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
    },
    cloneModuleAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.loading.cloneModule = false;
    },
    cloneModuleValidationOk() {
      this.loading.cloneModule = false;

      // hide modal
      this.$emit("hide");
    },
    cloneModuleCompleted() {
      // reload instances
      this.$emit("cloneOrMoveCompleted");
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

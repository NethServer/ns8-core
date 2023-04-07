<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <NsModal
    kind="danger"
    size="default"
    :visible="isShown"
    :isLoading="loading.removeNode"
    @modal-hidden="onModalHidden"
    @primary-click="removeNode"
    :primary-button-disabled="userInput !== 'node' + (node ? node.id : '')"
  >
    <template slot="title">{{ $t("nodes.remove_node_from_cluster") }}</template>
    <template slot="content">
      <NsInlineNotification
        kind="warning"
        :title="$t('common.please_read_carefully')"
        :showCloseButton="false"
      />
      <div
        v-html="
          $t('nodes.remove_node_confirm', {
            name: node ? this.getNodeLabel(node) : '',
          })
        "
      ></div>
      <div
        class="mg-top-xlg"
        v-html="
          $t('common.type_to_confirm', { name: 'node' + (node ? node.id : '') })
        "
      ></div>
      <cv-form @submit.prevent="removeNode">
        <NsTextInput
          v-model="userInput"
          :disabled="loading.removeNode"
          class="mg-bottom-md"
          ref="userInput"
        >
        </NsTextInput>
      </cv-form>
      <NsInlineNotification
        v-if="error.removeNode"
        kind="error"
        :title="$t('action.remove-node')"
        :description="error.removeNode"
        :showCloseButton="false"
        class="mg-top-lg"
      />
    </template>
    <template slot="secondary-button">{{ $t("common.cancel") }}</template>
    <template slot="primary-button">{{
      $t("nodes.understood_remove_node")
    }}</template>
  </NsModal>
</template>

<script>
import { TaskService, UtilService } from "@nethserver/ns8-ui-lib";
import to from "await-to-js";

export default {
  name: "RemoveNodeModal",
  mixins: [TaskService, UtilService],
  props: {
    isShown: {
      type: Boolean,
      default: true,
    },
    node: { type: Object },
  },
  data() {
    return {
      userInput: "",
      loading: {
        removeNode: false,
      },
      error: {
        removeNode: "",
      },
    };
  },
  watch: {
    isShown: function () {
      if (this.isShown) {
        this.userInput = "";
        this.clearErrors();

        setTimeout(() => {
          this.focusElement("userInput");
        }, 300);
      }
    },
  },
  methods: {
    onModalHidden() {
      this.$emit("hide");
    },
    async removeNode() {
      if (this.userInput !== "node" + this.node.id) {
        return;
      }
      this.loading.removeNode = true;
      this.error.removeNode = "";
      const taskAction = "remove-node";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.removeNodeAborted
      );

      // register to task validation
      this.$root.$once(
        `${taskAction}-validation-failed-${eventId}`,
        this.removeNodeValidationFailed
      );
      this.$root.$once(
        `${taskAction}-validation-ok-${eventId}`,
        this.removeNodeValidationOk
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.removeNodeCompleted
      );

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            node_id: this.node.id,
          },
          extra: {
            title: this.$t("action." + taskAction),
            description: this.$t("common.processing"),
            eventId,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.removeNode = this.getErrorMessage(err);
        this.loading.removeNode = false;
        return;
      }
    },
    removeNodeAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.loading.removeNode = false;
    },
    removeNodeValidationOk() {
      // hide modal
      this.$emit("hide");
    },
    removeNodeValidationFailed(validationErrors) {
      this.loading.removeNode = false;

      for (const validationError of validationErrors) {
        // set i18n error message

        if (validationError.error === "node_in_use") {
          // some apps are installed on the node
          const apps = validationError.value.join(", ");
          this.error.removeNode = this.$tc(
            "nodes.validation_node_in_use",
            apps.length,
            { apps }
          );
        } else {
          // other validation errors
          this.error.removeNode = this.$t(
            "nodes.validation_" + validationError.error
          );
        }
      }
    },
    removeNodeCompleted() {
      this.loading.removeNode = false;

      // hide modal
      this.$emit("hide");
      this.$emit("nodeRemoved");
    },
  },
};
</script>

<style scoped lang="scss"></style>

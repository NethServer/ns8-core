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
      <cv-skeleton-text
        v-if="loading.listInstalledModules"
        :paragraph="true"
        :line-count="4"
        heading
      ></cv-skeleton-text>
      <template v-else-if="!error.listInstalledModules">
        <NsInlineNotification
          kind="warning"
          :title="$t('common.please_read_carefully')"
          :description="
            nodeApps.length
              ? $tc('nodes.remove_node_apps_confirm', nodeApps.length, {
                  node: node ? this.getNodeLabel(node) : '',
                  nodeApps: nodeApps.join(', '),
                })
              : ''
          "
          :showCloseButton="false"
        />
        <NsInlineNotification
          v-if="nodeApps.length"
          kind="info"
          :title="$t('common.remember')"
          :description="
            $t('nodes.shutdown_node_after_removal', {
              node: node ? this.getNodeLabel(node) : '',
            })
          "
          :showCloseButton="false"
        />
        <div>
          {{
            $t("nodes.remove_node_confirm", {
              name: node ? this.getNodeLabel(node) : "",
            })
          }}
        </div>
        <div class="mg-top-xlg">
          {{
            $t("common.type_to_confirm", {
              name: "node" + (node ? node.id : ""),
            })
          }}
        </div>
        <cv-form @submit.prevent="removeNode">
          <NsTextInput
            v-model="userInput"
            :disabled="loading.removeNode"
            class="mg-bottom-md hide-label"
            ref="userInput"
          >
          </NsTextInput>
        </cv-form>
      </template>
      <!-- need to wrap error notification inside a div: custom elements like NsInlineNotification don't have scrollIntoView() function -->
      <div ref="listInstalledModulesError">
        <NsInlineNotification
          v-if="error.listInstalledModules"
          kind="error"
          :title="$t('action.list-installed-modules')"
          :description="error.listInstalledModules"
          :showCloseButton="false"
          class="mg-top-lg"
        />
      </div>
      <div ref="removeNodeError">
        <NsInlineNotification
          v-if="error.removeNode"
          kind="error"
          :title="$t('action.remove-node')"
          :description="error.removeNode"
          :showCloseButton="false"
          class="mg-top-lg"
        />
      </div>
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
import NodeService from "@/mixins/node";

export default {
  name: "RemoveNodeModal",
  mixins: [TaskService, UtilService, NodeService],
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
      nodeApps: [],
      loading: {
        removeNode: false,
        listInstalledModules: false,
      },
      error: {
        removeNode: "",
        listInstalledModules: "",
      },
    };
  },
  watch: {
    isShown: function () {
      if (this.isShown) {
        this.userInput = "";
        this.clearErrors();
        this.listInstalledModules();
      }
    },
    "error.listInstalledModules": function () {
      if (this.error.listInstalledModules) {
        // scroll to notification error

        this.$nextTick(() => {
          const el = this.$refs.listInstalledModulesError;
          this.scrollToElement(el);
        });
      }
    },
    "error.removeNode": function () {
      if (this.error.removeNode) {
        // scroll to notification error

        this.$nextTick(() => {
          const el = this.$refs.removeNodeError;
          this.scrollToElement(el);
        });
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

      // reload app drawer (any apps on the node have been removed)
      this.$root.$emit("reloadAppDrawer");
    },
    async listInstalledModules() {
      this.loading.listInstalledModules = true;
      const taskAction = "list-installed-modules";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.listInstalledModulesAborted
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.listInstalledModulesCompleted
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
        const errMessage = this.getErrorMessage(err);
        this.error.listInstalledModules = errMessage;
        this.loading.listInstalledModules = false;
        return;
      }
    },
    listInstalledModulesAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.listInstalledModules = this.$t("error.generic_error");
      this.loading.listInstalledModules = false;
    },
    listInstalledModulesCompleted(taskContext, taskResult) {
      const appsByNode = this.getInstalledAppsByNode(taskResult.output);
      let nodeApps = [];
      if (this.node.id in appsByNode) {
        nodeApps = appsByNode[this.node.id];
      }
      nodeApps = nodeApps.filter((app) => {
        // exclude core apps (but keep account providers)
        return (
          !app.flags.includes("core_module") ||
          app.flags.includes("account_provider")
        );
      });

      this.nodeApps = nodeApps.map((app) => {
        return app.ui_name ? `${app.ui_name} (${app.id})` : app.id;
      });

      this.loading.listInstalledModules = false;
      this.focusElement("userInput");
    },
  },
};
</script>

<style scoped lang="scss"></style>

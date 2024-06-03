<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <NsModal
    size="default"
    :visible="isShown"
    @modal-hidden="onModalHidden"
    @primary-click="installInstance"
    class="no-pad-modal"
    :primary-button-disabled="!selectedNode"
  >
    <template v-if="app" slot="title">{{
      $t("software_center.app_installation", { app: app.name })
    }}</template>
    <template v-if="app" slot="content">
      <cv-form @submit.prevent="installInstance">
        <template v-if="clusterNodes.length > 1">
          <div>
            {{
              $t("software_center.choose_node_for_installation", {
                app: app.name,
                version: appVersion,
              })
            }}
          </div>
          <NodeSelector @selectNode="onSelectNode" class="mg-top-xlg" />
        </template>
        <div v-else>
          {{
            $t("software_center.about_to_install_app", {
              app: app.name,
              version: appVersion,
              node: firstNodeLabel,
            })
          }}
        </div>
        <div v-if="error.addModule">
          <NsInlineNotification
            kind="error"
            :title="$t('action.add-module')"
            :description="error.addModule"
            :showCloseButton="false"
          />
        </div>
      </cv-form>
    </template>
    <template slot="secondary-button">{{ $t("common.cancel") }}</template>
    <template slot="primary-button">{{
      $t("software_center.install")
    }}</template>
  </NsModal>
</template>

<script>
import { UtilService, TaskService, IconService } from "@nethserver/ns8-ui-lib";
import to from "await-to-js";
import NotificationService from "@/mixins/notification";
import NodeSelector from "@/components/nodes/NodeSelector";
import { mapState } from "vuex";

export default {
  name: "InstallAppModal",
  components: { NodeSelector },
  mixins: [UtilService, TaskService, IconService, NotificationService],
  props: {
    isShown: Boolean,
    app: { type: [Object, null] },
  },
  data() {
    return {
      selectedNode: null,
      error: {
        addModule: "",
      },
    };
  },
  computed: {
    ...mapState(["clusterNodes"]),
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
        return "latest";
      }
    },
  },
  created() {
    if (this.clusterNodes.length == 1) {
      this.selectedNode = this.clusterNodes[0];
    }
  },
  methods: {
    async installInstance() {
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

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            image: this.app.source + ":" + this.appVersion,
            node: parseInt(this.selectedNode.id),
          },
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
      this.clearErrors(this);
      this.$emit("close");
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

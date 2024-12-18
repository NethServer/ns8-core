<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <NsModal
    size="default"
    :visible="isShown"
    @modal-hidden="onModalHidden"
    @modal-shown="onModalShown"
    @primary-click="installInstance"
    class="no-pad-modal"
    :primary-button-disabled="
      !selectedNode || (app && app.docs.terms_url && !agreeTerms)
    "
  >
    <template v-if="app" slot="title">{{
      $t("software_center.app_installation", { app: app.name })
    }}</template>
    <template v-if="app" slot="content">
      <cv-form @submit.prevent="installInstance">
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
        <template v-if="clusterNodes.length > 1">
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
            v-if="clusterNodes.length == disabledNodes.length"
            kind="info"
            :title="$t('software_center.no_node_eligible_for_app_installation')"
            :showCloseButton="false"
          />
          <NodeSelector
            @selectNode="onSelectNode"
            :disabledNodes="disabledNodes"
            class="mg-top-lg"
          >
            <template v-for="(nodeInfoMessage, nodeId) in nodesInfo">
              <template :slot="`node-${nodeId}`">
                {{ nodeInfoMessage }}
              </template>
            </template>
          </NodeSelector>
        </template>
        <!-- single node -->
        <div v-else class="mg-bottom-lg">
          <template v-if="canInstallOnSingleNode">
            {{
              $t("software_center.about_to_install_app", {
                app: app.name,
                version: appVersion,
                node: firstNodeLabel,
              })
            }}
          </template>
          <template v-else>
            {{
              $t("software_center.cannot_install_app_on_node", {
                app: app.name,
                version: appVersion,
                node: firstNodeLabel,
              })
            }}:
            {{ nodesInfo[clusterNodes[0].id] }}
          </template>
        </div>
        <!-- terms and conditions -->
        <NsCheckbox
          v-if="app.docs.terms_url"
          v-model="agreeTerms"
          value="checkTerms"
        >
          <template slot="label">
            <i18n path="software_center.agree_terms_before_install" tag="span">
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
      agreeTerms: false,
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
      }
      return nodesInfo;
    },
    disabledNodes() {
      return this.app.install_destinations
        .filter((nodeInfo) => !nodeInfo.eligible)
        .map((nodeInfo) => nodeInfo.node_id);
    },
  },
  methods: {
    onModalShown() {
      this.agreeTerms = false;

      if (this.clusterNodes.length == 1 && this.canInstallOnSingleNode) {
        this.selectedNode = this.clusterNodes[0];
      } else {
        this.selectedNode = null;
      }
    },
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

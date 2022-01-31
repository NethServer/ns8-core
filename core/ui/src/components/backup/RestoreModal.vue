<template>
  <NsWizard
    size="default"
    :visible="isShown"
    :cancelLabel="$t('common.cancel')"
    :previousLabel="$t('common.previous')"
    :nextLabel="isLastStep ? $t('backup.restore') : $t('common.next')"
    :isPreviousDisabled="isFirstStep || loading.restoreModule"
    :isNextDisabled="isNextButtonDisabled"
    :isNextLoading="loading.restoreModule"
    @modal-hidden="$emit('hide')"
    @cancel="$emit('hide')"
    @previousStep="previousStep"
    @nextStep="nextStep"
  >
    <template slot="title">{{ $t("backup.restore") }}</template>
    <template slot="content">
      <cv-form>
        <NsInlineNotification
          v-if="error.getClusterStatus"
          kind="error"
          :title="$t('action.get-cluster-status')"
          :description="error.getClusterStatus"
          :showCloseButton="false"
        />
        <NsInlineNotification
          v-if="error.readBackupRepository"
          kind="error"
          :title="$t('action.read-backup-repository')"
          :description="error.readBackupRepository"
          :showCloseButton="false"
        />
        <template v-if="step == 'instance'">
          <div class="mg-bottom-md">
            {{ $t("backup.select_instance_to_restore") }}
          </div>
          <cv-grid class="instances mg-bottom-md no-padding">
            <cv-row>
              <cv-column v-if="loading.readBackupRepository" :lg="16">
                <cv-tile light>
                  <cv-skeleton-text
                    :paragraph="true"
                    :line-count="4"
                  ></cv-skeleton-text>
                </cv-tile>
              </cv-column>
              <cv-column v-else-if="!instances.length" :lg="16">
                <NsEmptyState :title="$t('backup.no_instance_to_restore')">
                  <template #description>{{
                    $t("backup.no_instance_to_restore_description")
                  }}</template>
                </NsEmptyState>
              </cv-column>
              <cv-column
                v-else
                v-for="(instance, index) in instances"
                :key="index"
                :lg="16"
              >
                <NsTile
                  :light="true"
                  kind="selectable"
                  v-model="instance.selected"
                  value="nodeValue"
                  @click="deselectOtherInstances(instance)"
                  class="mg-bottom-sm"
                >
                  <h6>
                    {{ instance.path }}
                  </h6>
                </NsTile>
              </cv-column>
            </cv-row>
          </cv-grid>
        </template>
        <template v-if="step == 'node'">
          <div class="mg-bottom-md">
            {{ $t("backup.select_restore_node") }}
          </div>
          <div class="bx--grid nodes mg-bottom-md no-padding">
            <div class="bx--row">
              <div
                v-for="(node, index) in nodes"
                :key="index"
                class="bx--col-md-4 bx--col-max-4"
              >
                <NsTile
                  :light="true"
                  kind="selectable"
                  v-model="node.selected"
                  value="nodeValue"
                  :footerIcon="Chip20"
                  @click="deselectOtherNodes(node)"
                >
                  <h6>
                    {{
                      node.ui_name
                        ? node.ui_name
                        : $t("common.node") + " " + node.id
                    }}
                  </h6>
                  <div v-if="node.ui_name" class="node-id">
                    {{ $t("common.node") }} {{ node.id }}
                  </div>
                </NsTile>
              </div>
            </div>
          </div>
          <NsInlineNotification
            v-if="error.restoreModule"
            kind="error"
            :title="$t('action.restore-module')"
            :description="error.restoreModule"
            :showCloseButton="false"
          />
        </template>
      </cv-form>
    </template>
  </NsWizard>
</template>

<script>
import { UtilService, TaskService, IconService } from "@nethserver/ns8-ui-lib";
import to from "await-to-js";

export default {
  name: "RestoreModal",
  mixins: [UtilService, TaskService, IconService],
  props: {
    isShown: {
      type: Boolean,
      default: true,
    },
    repositoryId: {
      type: String,
      required: true,
    },
  },
  data() {
    return {
      step: "",
      steps: ["instance", "node"],
      instance: {}, //// check var type
      nodes: [],
      instances: [],
      loading: {
        getClusterStatus: true,
        readBackupRepository: true,
        restoreModule: false,
      },
      error: {
        getClusterStatus: "",
        readBackupRepository: "",
        restoreModule: "",
      },
    };
  },
  computed: {
    stepIndex() {
      return this.steps.indexOf(this.step);
    },
    isFirstStep() {
      return this.stepIndex == 0;
    },
    isLastStep() {
      return this.stepIndex == this.steps.length - 1;
    },
    isNextButtonDisabled() {
      return (
        this.loading.restoreModule ||
        (this.step == "instance" && !this.selectedInstance) ||
        (this.step == "node" && !this.selectedNode)
      );
    },
    selectedInstance() {
      return this.instances.find((i) => i.selected);
    },
    selectedNode() {
      return this.nodes.find((n) => n.selected);
    },
  },
  watch: {
    isShown: function () {
      if (this.isShown) {
        // show first step
        this.step = this.steps[0];

        // ensure repositoryId prop is updated
        this.$nextTick(() => {
          console.log("inside nexttick", this.repositoryId); ////

          if (this.repositoryId) {
            this.readBackupRepository();
            this.getClusterStatus();
          }
        });
      }
    },
  },
  methods: {
    nextStep() {
      if (this.isNextButtonDisabled) {
        return;
      }

      if (this.isLastStep) {
        this.restoreModule();
      } else {
        this.step = this.steps[this.stepIndex + 1];
      }
    },
    previousStep() {
      if (!this.isFirstStep) {
        this.step = this.steps[this.stepIndex - 1];
      }
    },
    async getClusterStatus() {
      this.error.getClusterStatus = "";
      this.loading.getClusterStatus = true;
      const taskAction = "get-cluster-status";

      // register to task completion
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
    getClusterStatusCompleted(taskContext, taskResult) {
      const clusterStatus = taskResult.output;
      let nodes = clusterStatus.nodes.sort(this.sortByProperty("id"));

      //// remove mock
      // nodes.push({ id: 2, local: false, ui_name: "" }); ////
      // nodes.push({ id: 3, local: false, ui_name: "" }); ////

      for (const node of nodes) {
        node.selected = false;
      }

      if (nodes.length == 1) {
        // select node if there's only one
        nodes[0].selected = true;
      }

      this.nodes = nodes;
      this.loading.getClusterStatus = false;
    },
    async readBackupRepository() {
      this.error.readBackupRepository = "";
      this.loading.readBackupRepository = true;
      const taskAction = "read-backup-repository";

      // register to task completion
      this.$root.$once(
        taskAction + "-completed",
        this.readBackupRepositoryCompleted
      );

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            repository: this.repositoryId,
          },
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.readBackupRepository = this.getErrorMessage(err);
        return;
      }
    },
    readBackupRepositoryCompleted(taskContext, taskResult) {
      let instances = [];

      for (const [app, instancePaths] of Object.entries(taskResult.output)) {
        for (const instancePath of instancePaths) {
          const instance = {
            path: `${app}/${instancePath}`,
            selected: false,
          };
          instances.push(instance);
        }
      }
      this.instances = instances;
      this.loading.readBackupRepository = false;
    },
    async restoreModule() {
      this.loading.restoreModule = true;
      this.error.restoreModule = "";
      const taskAction = "restore-module";

      // register to task error
      this.$root.$off(taskAction + "-aborted");
      this.$root.$once(taskAction + "-aborted", this.restoreModuleAborted);

      // register to task validation
      // this.$root.$off(taskAction + "-validation-ok"); ////
      // this.$root.$once(
      //   taskAction + "-validation-ok",
      //   this.restoreModuleValidationOk
      // );

      // register to task completion
      this.$root.$off(taskAction + "-completed");
      this.$root.$once(taskAction + "-completed", this.restoreModuleCompleted);

      const app = this.selectedInstance.path.split("/")[0];
      const nodeName =
        this.selectedNode.ui_name ||
        this.$t("common.node") + ` ${this.selectedNode.id}`;

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            repository: this.repositoryId,
            path: this.selectedInstance.path,
            snapshot: "", // latest
            node: this.selectedNode.id,
          },
          extra: {
            title: this.$t("action." + taskAction),
            description: this.$t("backup.restoring_app_to_node", {
              app: app,
              node: nodeName,
            }),
            node: nodeName,
            completion: {
              i18nString: "backup.instance_restored_to_node",
              extraTextParams: ["node"],
              outputTextParams: ["module_id"],
            },
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.restoreModule = this.getErrorMessage(err);
        return;
      }

      // disable restore command on UI until completion
      this.$emit("enableRestore", false);

      // close modal immediately, no validation needed
      this.$emit("hide");
    },
    // restoreModuleValidationOk() { ////
    //   // hide modal after validation
    //   this.$emit("hide");
    // },
    restoreModuleAborted(taskResult) {
      console.error("restore-module aborted", taskResult);
      this.loading.restoreModule = false;
      this.$emit("hide");

      // re-enable restore command on UI
      this.$emit("enableRestore", true);
    },
    restoreModuleCompleted() {
      this.loading.restoreModule = false;

      // re-enable restore command on UI
      this.$emit("enableRestore", true);

      // show restored instance in app drawer
      this.$root.$emit("reloadAppDrawer");
    },
    deselectOtherNodes(node) {
      for (let n of this.nodes) {
        if (n.id !== node.id) {
          n.selected = false;
        }
      }
    },
    deselectOtherInstances(instance) {
      for (let i of this.instances) {
        if (i.path !== instance.path) {
          i.selected = false;
        }
      }
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";

.instances {
  margin-top: $spacing-07;
}

.nodes {
  margin-top: $spacing-07;
}

.node-id {
  margin-top: $spacing-05;
}
</style>

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
    <template slot="title">{{ $t("backup.restore_app") }}</template>
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
          v-if="error.readBackupRepositories"
          kind="error"
          :title="$t('action.read-backup-repositories')"
          :description="error.readBackupRepositories"
          :showCloseButton="false"
        />
        <template v-if="step == 'instance'">
          <div class="mg-bottom-md">
            {{ $t("backup.select_instance_to_restore") }}
          </div>
          <cv-grid class="instances mg-bottom-md no-padding">
            <cv-row>
              <cv-column
                v-if="!loading.readBackupRepositories && !instances.length"
              >
                <NsEmptyState :title="$t('backup.no_instance_to_restore')">
                  <template #description>{{
                    $t("backup.no_instance_to_restore_description")
                  }}</template>
                </NsEmptyState>
              </cv-column>
            </cv-row>
            <cv-row>
              <cv-column>
                <RestoreSingleInstanceSelector
                  :instances="instances"
                  :loading="loading.readBackupRepositories"
                  :light="true"
                  @select="onSelectInstance"
                />
              </cv-column>
            </cv-row>
            <cv-row>
              <cv-column>
                <cv-checkbox
                  :label="
                    $t('backup.replace_existing_app', {
                      app: instanceToReplace,
                    })
                  "
                  v-model="replaceExistingApp"
                  :disabled="!instanceToReplace"
                  value="checkReplaceExistingApp"
                  class="mg-top-lg"
                />
              </cv-column>
            </cv-row>
          </cv-grid>
        </template>
        <template v-if="step == 'node'">
          <div>
            {{ $t("backup.select_restore_node") }}
          </div>
          <NodeSelector
            :nodes="nodes"
            @selectNode="onSelectNode"
            class="mg-top-lg"
          />
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
import NodeSelector from "@/components/NodeSelector";
import RestoreSingleInstanceSelector from "@/components/backup/RestoreSingleInstanceSelector";

export default {
  name: "RestoreSingleInstanceModal",
  components: { NodeSelector, RestoreSingleInstanceSelector },
  mixins: [UtilService, TaskService, IconService],
  props: {
    isShown: {
      type: Boolean,
      default: true,
    },
  },
  data() {
    return {
      step: "",
      steps: ["instance", "node"],
      nodes: [],
      instances: [],
      selectedNode: null,
      selectedInstance: null,
      replaceExistingApp: false,
      loading: {
        getClusterStatus: true,
        readBackupRepositories: true,
        restoreModule: false,
      },
      error: {
        getClusterStatus: "",
        readBackupRepositories: "",
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
    instanceToReplace() {
      if (
        !(this.selectedInstance && this.selectedInstance.installed_instance)
      ) {
        return "";
      }
      return this.selectedInstance.installed_instance;
    },
  },
  watch: {
    isShown: function () {
      if (this.isShown) {
        // show first step
        this.step = this.steps[0];
        this.selectedInstance = null;
        this.replaceExistingApp = false;
        this.readBackupRepositories();
        this.getClusterStatus();
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
    async readBackupRepositories() {
      this.error.readBackupRepositories = "";
      this.loading.readBackupRepositories = true;
      const taskAction = "read-backup-repositories";

      // register to task error
      this.$root.$off(taskAction + "-aborted");
      this.$root.$once(
        taskAction + "-aborted",
        this.readBackupRepositoriesAborted
      );

      // register to task completion
      this.$root.$off(taskAction + "-completed");
      this.$root.$once(
        taskAction + "-completed",
        this.readBackupRepositoriesCompleted
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
        this.error.readBackupRepositories = this.getErrorMessage(err);
        return;
      }
    },
    readBackupRepositoriesAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.loading.readBackupRepositories = false;
    },
    readBackupRepositoriesCompleted(taskContext, taskResult) {
      let instances = taskResult.output;

      for (const instance of instances) {
        // prepare data for instance selector
        instance.selected = false;
      }
      this.instances = instances;
      this.loading.readBackupRepositories = false;
    },
    async restoreModule() {
      this.loading.restoreModule = true;
      this.error.restoreModule = "";
      const taskAction = "restore-module";

      // register to task error
      this.$root.$off(taskAction + "-aborted");
      this.$root.$once(taskAction + "-aborted", this.restoreModuleAborted);

      // register to task completion
      this.$root.$off(taskAction + "-completed");
      this.$root.$once(taskAction + "-completed", this.restoreModuleCompleted);

      const app = this.selectedInstance.path.split("/")[0];
      const nodeName =
        this.selectedNode.ui_name ||
        this.$t("common.node_lc") + ` ${this.selectedNode.id}`;

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            repository: this.selectedInstance.repository_id,
            path: this.selectedInstance.path,
            snapshot: "", // latest
            node: this.selectedNode.id,
            replace: this.replaceExistingApp,
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

      // close modal immediately, no validation needed
      this.$emit("hide");
    },
    restoreModuleAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.loading.restoreModule = false;
      this.$emit("hide");
    },
    restoreModuleCompleted() {
      this.loading.restoreModule = false;

      // show restored instance in app drawer
      this.$root.$emit("reloadAppDrawer");
    },
    onSelectNode(selectedNode) {
      this.selectedNode = selectedNode;
    },
    onSelectInstance(selectedInstance) {
      this.selectedInstance = selectedInstance;
      this.replaceExistingApp = false;
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";

.instances {
  margin-top: $spacing-07;
}
</style>

<template>
  <cv-modal
    size="default"
    :visible="isShown"
    @modal-hidden="$emit('hide')"
    class="wizard-modal"
  >
    <template slot="title">{{ $t("backup.schedule_backup") }}</template>
    <template slot="content">
      <cv-form @submit.prevent="nextStep">
        <div v-show="step == 'instances'">
          <div class="mg-bottom-md">
            {{ $t("backup.choose_app_instances_to_backup") }}
          </div>
          <InstanceSelector
            :instances="installedModules"
            :selection="instanceSelection"
            :instancesNotBackedUp="instancesNotBackedUp"
            :loading="loading.listInstalledModules"
            @select="onSelectInstances"
          />
        </div>
        <div v-show="step == 'repository'">
          <div class="mg-bottom-md">
            {{ $t("backup.choose_destination_repository") }}
          </div>
          <div class="bx--grid">
            <div class="bx--row">
              <div
                v-for="repo in internalRepositories"
                :key="'col-' + repo.id"
                class="bx--col-md-4"
              >
                <NsTile
                  :light="true"
                  kind="selectable"
                  :value="repo.id"
                  :footerIcon="DataBase20"
                  v-model="repo.selected"
                  @click="deselectOtherRepos(repo)"
                  class="same-height-tile"
                >
                  <h6>
                    {{ repo.name }}
                  </h6>
                  <div class="mg-top-sm">
                    {{ $t("backup." + repo.provider) }}
                  </div>
                  <div class="mg-top-sm ellipsis">
                    {{ repo.url }}
                  </div>
                </NsTile>
              </div>
            </div>
          </div>
        </div>
        <div v-show="step == 'settings'">
          <cv-text-input
            :label="$t('backup.frequency')"
            v-model.trim="schedule"
            :invalid-message="$t(error.schedule)"
            :disabled="loading.addBackup"
            ref="schedule"
          >
          </cv-text-input>
          <cv-text-input
            :label="$t('backup.retention')"
            v-model.trim="retention"
            :invalid-message="$t(error.retention)"
            :disabled="loading.addBackup"
            ref="retention"
          >
          </cv-text-input>
          <cv-text-input
            :label="$t('backup.backup_name')"
            v-model.trim="name"
            :invalid-message="$t(error.name)"
            :disabled="loading.addBackup"
            ref="name"
          >
          </cv-text-input>
          <cv-toggle
            :label="$t('common.status')"
            value="statusValue"
            :form-item="true"
            v-model="enabled"
          >
            <template slot="text-left">{{ $t("common.disabled") }}</template>
            <template slot="text-right">{{ $t("common.enabled") }}</template>
          </cv-toggle>
          <NsInlineNotification
            v-if="error.addBackup"
            kind="error"
            :title="$t('action.add-backup')"
            :description="error.addBackup"
            :showCloseButton="false"
          />
        </div>
        <div class="wizard-buttons">
          <NsButton
            kind="secondary"
            :icon="Close20"
            @click="$emit('hide')"
            type="button"
            class="wizard-button"
            >{{ $t("common.cancel") }}
          </NsButton>
          <NsButton
            kind="secondary"
            :icon="ChevronLeft20"
            @click="previousStep"
            :disabled="isFirstStep || loading.addBackup"
            type="button"
            class="wizard-button"
            >{{ $t("common.previous") }}
          </NsButton>
          <NsButton
            kind="primary"
            :icon="ChevronRight20"
            :disabled="isNextStepDisabled"
            :loading="loading.addBackup"
            type="submit"
            class="wizard-button"
            ref="wizardNext"
            >{{ isLastStep ? $t("common.finish") : $t("common.next") }}
          </NsButton>
        </div>
      </cv-form>
    </template>
  </cv-modal>
</template>

<script>
import { UtilService, TaskService, IconService } from "@nethserver/ns8-ui-lib";
import to from "await-to-js";
import InstanceSelector from "@/components/backup/InstanceSelector";
import _cloneDeep from "lodash/cloneDeep";

export default {
  name: "CreateBackupModal",
  components: { InstanceSelector },
  mixins: [UtilService, TaskService, IconService],
  props: {
    isShown: {
      type: Boolean,
      default: true,
    },
    repositories: {
      type: Array,
      required: true,
    },
    instanceSelection: {
      type: String,
      default: "",
    },
    instancesNotBackedUp: {
      type: Array,
      default: () => [],
    },
  },
  data() {
    return {
      step: "",
      steps: ["instances", "repository", "settings"],
      instances: [],
      name: "",
      schedule: "daily", ////
      retention: "7d", ////
      enabled: true,
      installedModules: [],
      internalRepositories: [],
      loading: {
        addBackup: false,
        listInstalledModules: true,
      },
      error: {
        name: "",
        repository: "",
        schedule: "",
        retention: "",
        addBackup: "",
        listInstalledModules: "",
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
    isNextStepDisabled() {
      return (
        this.loading.addBackup ||
        !this.installedModules.length ||
        (this.step == "instances" && !this.instances.length) ||
        (this.step == "repository" && !this.selectedRepo)
      );
    },
    selectedRepo() {
      return this.internalRepositories.find((r) => r.selected);
    },
  },
  watch: {
    isShown: function () {
      if (this.isShown) {
        // show first step
        this.step = this.steps[0];
        this.clearFields();

        // load installed moudules
        this.listInstalledModules();
      }
    },
    repositories: function () {
      this.updateInternalRepositories();
    },
  },
  created() {
    this.updateInternalRepositories();
  },
  methods: {
    clearFields() {
      this.instances = [];
      this.name = "";
      this.schedule = "daily"; ////
      this.retention = "7d"; ////
      this.enabled = true;

      for (let repo of this.internalRepositories) {
        repo.selected = false;
      }
    },
    nextStep() {
      if (this.isNextStepDisabled) {
        return;
      }

      if (this.isLastStep) {
        this.addBackup();
      } else {
        this.step = this.steps[this.stepIndex + 1];
      }
    },
    previousStep() {
      if (!this.isFirstStep) {
        this.step = this.steps[this.stepIndex - 1];
      }
    },
    updateInternalRepositories() {
      // deep copy (needed to avoid reactivity issues)
      let internalRepositories = _cloneDeep(this.repositories);

      for (const repo of internalRepositories) {
        repo.selected = false;
      }
      this.internalRepositories = internalRepositories;
    },
    async addBackup() {
      //// validation
      this.error.addBackup = "";
      this.loading.addBackup = true;
      const taskAction = "add-backup";

      // register to task validation
      this.$root.$off(taskAction + "-validation-failed");
      this.$root.$once(
        taskAction + "-validation-failed",
        this.addBackupValidationFailed
      );
      this.$root.$off(taskAction + "-validation-ok");
      this.$root.$once(
        taskAction + "-validation-ok",
        this.addBackupValidationOk
      );

      // register to task completion
      this.$root.$once(taskAction + "-completed", this.addBackupCompleted);

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            name: this.name,
            repository: this.selectedRepo.id,
            schedule: this.schedule,
            retention: this.retention,
            instances: this.instances,
            enabled: this.enabled,
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
        this.error.addBackup = this.getErrorMessage(err);
        return;
      }
    },
    addBackupValidationFailed(validationErrors) {
      this.loading.addBackup = false;
      let focusAlreadySet = false;

      for (const validationError of validationErrors) {
        const param = validationError.parameter;

        // set i18n error message
        this.error[param] = "backup." + validationError.error;

        if (!focusAlreadySet) {
          this.focusElement(param);
          focusAlreadySet = true;
        }
      }
    },
    addBackupValidationOk() {
      console.log("addBackupValidationOk"); ////

      // hide modal after validation
      this.$emit("hide");
    },
    addBackupCompleted(taskContext, taskResult) {
      console.log("addBackupRepositoryCompleted", taskResult.output); ////

      this.loading.addBackup = false;

      // hide modal
      // this.$emit("hide"); ////

      // reload backup configuration
      this.$emit("backupCreated");
    },
    async listInstalledModules() {
      this.loading.listInstalledModules = true;
      const taskAction = "list-installed-modules";

      // register to task completion
      this.$root.$once(
        taskAction + "-completed",
        this.listInstalledModulesCompleted
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
        this.error.listInstalledModules = this.getErrorMessage(err);
        this.createErrorNotification(
          err,
          this.$t("task.cannot_create_task", { action: taskAction })
        );
        return;
      }
    },
    listInstalledModulesCompleted(taskContext, taskResult) {
      this.loading.listInstalledModules = false;
      let apps = [];

      for (let instanceList of Object.values(taskResult.output)) {
        for (let instance of instanceList) {
          apps.push(instance);
        }
      }

      apps.sort(this.sortModuleInstances());

      console.log("appsss", apps); ////

      this.installedModules = apps;
    },
    onSelectInstances(instances) {
      console.log("onSelectInstances", instances); ////

      this.instances = instances;
    },
    deselectOtherRepos(repo) {
      for (let r of this.internalRepositories) {
        if (r.id !== repo.id) {
          r.selected = false;
        }
      }
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";

.same-height-tile {
  min-height: 9rem;
}
</style>

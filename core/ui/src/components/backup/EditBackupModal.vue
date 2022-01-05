<template>
  <cv-modal
    size="default"
    :visible="isShown"
    @modal-hidden="$emit('hide')"
    class="wizard-modal"
  >
    <template slot="title">{{ $t("backup.edit_backup") }}</template>
    <template slot="content">
      <cv-form @submit.prevent="nextStep">
        <div v-show="step == 'instances'">
          <div class="mg-bottom-md">
            {{ $t("backup.choose_app_instances_to_backup") }}
          </div>
          <InstanceSelector
            :instances="installedModules"
            :selection="backup.instances"
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
                  class="min-height-card"
                >
                  <h6>
                    {{ repo.name }}
                  </h6>
                  <div class="mg-top-md">
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
            :label="$t('backup.schedule')"
            v-model.trim="schedule"
            :invalid-message="$t(error.schedule)"
            :disabled="loading.alterBackup"
            ref="schedule"
          >
          </cv-text-input>
          <cv-text-input
            :label="$t('backup.retention')"
            v-model.trim="retention"
            :invalid-message="$t(error.retention)"
            :disabled="loading.alterBackup"
            ref="retention"
          >
          </cv-text-input>
          <cv-text-input
            :label="$t('backup.backup_name')"
            v-model.trim="name"
            :invalid-message="$t(error.name)"
            :disabled="loading.alterBackup"
            ref="name"
          >
          </cv-text-input>
          <!-- <cv-toggle ////
            :label="$t('common.status')"
            value="statusValue"
            :form-item="true"
            v-model="enabled"
          >
            <template slot="text-left">{{ $t("common.disabled") }}</template>
            <template slot="text-right">{{ $t("common.enabled") }}</template>
          </cv-toggle> -->
          <NsInlineNotification
            v-if="error.alterBackup"
            kind="error"
            :title="$t('action.alter-backup')"
            :description="error.alterBackup"
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
            :disabled="isFirstStep || loading.alterBackup"
            type="button"
            class="wizard-button"
            >{{ $t("common.previous") }}
          </NsButton>
          <NsButton
            kind="primary"
            :icon="ChevronRight20"
            :disabled="isNextStepDisabled"
            :loading="loading.alterBackup"
            type="submit"
            class="wizard-button"
            ref="wizardNext"
            >{{ isLastStep ? $t("common.save") : $t("common.next") }}
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
  name: "EditBackupModal",
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
    instancesNotBackedUp: {
      type: Array,
      default: () => [],
    },
    backup: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      step: "",
      steps: ["instances", "repository", "settings"],
      instances: [],
      name: "",
      schedule: "", ////
      retention: "", ////
      // enabled: false, ////
      installedModules: [],
      internalRepositories: [],
      loading: {
        alterBackup: false,
        listInstalledModules: true,
      },
      error: {
        name: "",
        repository: "",
        schedule: "",
        retention: "",
        alterBackup: "",
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
        this.loading.alterBackup ||
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

        // load installed moudules
        this.listInstalledModules();
      }
    },
    repositories: function () {
      this.updateInternalRepositories();
    },
    backup: function () {
      this.name = this.backup.name;
      this.schedule = this.backup.schedule;
      this.retention = this.backup.retention;
      // this.enabled = this.backup.enabled; ////
      this.updateInternalRepositories();
    },
  },
  created() {
    this.updateInternalRepositories();
  },
  methods: {
    nextStep() {
      if (this.isNextStepDisabled) {
        return;
      }

      if (this.isLastStep) {
        this.alterBackup();
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
        if (this.backup.repository == repo.id) {
          repo.selected = true;
        } else {
          repo.selected = false;
        }
      }

      this.internalRepositories = internalRepositories;
    },
    validateAlterBackup() {
      this.clearErrors(this);
      let isValidationOk = true;

      //// todo validate schedule and retention

      if (!this.name) {
        this.error.name = "common.required";

        if (isValidationOk) {
          this.focusElement("name");
          isValidationOk = false;
        }
      }

      return isValidationOk;
    },
    async alterBackup() {
      if (!this.validateAlterBackup()) {
        return;
      }
      this.error.alterBackup = "";
      this.loading.alterBackup = true;
      const taskAction = "alter-backup";

      // register to task validation
      this.$root.$off(taskAction + "-validation-failed");
      this.$root.$once(
        taskAction + "-validation-failed",
        this.alterBackupValidationFailed
      );

      // register to task completion
      this.$root.$off(taskAction + "-completed");
      this.$root.$once(taskAction + "-completed", this.alterBackupCompleted);

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            id: this.backup.id,
            name: this.name,
            // repository: this.selectedRepo.id, //// ?
            schedule: this.schedule,
            retention: this.retention,
            instances: this.instances,
            enabled: this.backup.enabled, ////
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
        this.error.alterBackup = this.getErrorMessage(err);
        return;
      }
    },
    alterBackupValidationFailed(validationErrors) {
      this.loading.alterBackup = false;
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
    alterBackupCompleted() {
      this.loading.alterBackup = false;

      // hide modal
      this.$emit("hide");

      // reload backup configuration
      this.$emit("backupAltered");
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
      this.installedModules = apps;
    },
    onSelectInstances(instances) {
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
</style>

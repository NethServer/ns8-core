<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <div class="create-backup-modal">
    <NsWizard
      size="default"
      :visible="isShown"
      :cancelLabel="$t('common.cancel')"
      :previousLabel="$t('common.previous')"
      :nextLabel="nextButtonLabel"
      :isPreviousDisabled="isFirstStep || loading.addBackup"
      :isNextDisabled="isNextButtonDisabled"
      :isNextLoading="loading.addBackup"
      @modal-hidden="$emit('hide')"
      @cancel="$emit('hide')"
      @previousStep="previousStep"
      @nextStep="nextStep"
    >
      <template slot="title">{{
        isEditing ? $t("backup.edit_backup") : $t("backup.schedule_backup")
      }}</template>
      <template slot="content">
        <cv-form>
          <div v-show="step == 'instances'">
            <div class="mg-bottom-md">
              {{ $t("backup.choose_app_instances_to_backup") }}
              <cv-interactive-tooltip
                alignment="start"
                direction="right"
                class="info"
              >
                <template slot="trigger">
                  <Information16 />
                </template>
                <template slot="content">
                  {{ $t("backup.cluster_configuration_info") }}
                </template>
              </cv-interactive-tooltip>
            </div>
            <BackupInstanceSelector
              :instances="installedModules"
              :selection="isEditing ? backup.instances : instanceSelection"
              :instancesNotBackedUp="instancesNotBackedUp"
              :loading="loading.listInstalledModules"
              @select="onSelectInstances"
            />
          </div>
          <div v-show="step == 'repository'">
            <div class="mg-bottom-md">
              {{ $t("backup.choose_destination_repository") }}
            </div>
            <cv-grid class="no-padding">
              <cv-row>
                <cv-column
                  v-for="repo in internalRepositories"
                  :key="'col-' + repo.id"
                  :md="4"
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
                </cv-column>
              </cv-row>
            </cv-grid>
          </div>
          <div v-show="step == 'settings'">
            <cv-grid class="schedule-container no-padding">
              <cv-row>
                <cv-column :md="2" :max="4">
                  <!-- schedule interval -->
                  <cv-select
                    v-model="schedule.interval"
                    :label="$t('backup.schedule')"
                    :disabled="loading.addBackup || loading.alterBackup"
                    :invalid-message="$t(error.schedule)"
                    class="schedule-interval"
                  >
                    <cv-select-option value="hourly">
                      {{ $t("backup.hourly") }}
                    </cv-select-option>
                    <cv-select-option value="daily">
                      {{ $t("backup.daily") }}
                    </cv-select-option>
                    <cv-select-option value="weekly">
                      {{ $t("backup.weekly") }}
                    </cv-select-option>
                    <cv-select-option value="monthly">
                      {{ $t("backup.monthly") }}
                    </cv-select-option>
                  </cv-select>
                </cv-column>
                <cv-column
                  v-if="schedule.interval == 'hourly'"
                  :md="2"
                  :max="4"
                >
                  <!-- hourly schedule -->
                  <cv-text-input
                    :label="$t('backup.minute')"
                    v-model.trim="schedule.minute"
                    type="number"
                    min="0"
                    max="59"
                    :disabled="loading.addBackup || loading.alterBackup"
                    ref="schedule-minute"
                  >
                  </cv-text-input>
                </cv-column>
                <!-- daily schedule -->
                <cv-column v-if="schedule.interval == 'daily'" :md="2" :max="4">
                  <NsTimePicker
                    hideClearButton
                    v-model="schedule.time"
                    :label="$t('backup.at')"
                  ></NsTimePicker>
                </cv-column>
                <!-- weekly schedule -->
                <template v-if="schedule.interval == 'weekly'">
                  <cv-column :md="2" :max="4">
                    <cv-select
                      v-model="schedule.weekDay"
                      :label="$t('backup.on')"
                      :disabled="loading.addBackup || loading.alterBackup"
                    >
                      <cv-select-option value="monday">
                        {{ $t("calendar.monday") }}
                      </cv-select-option>
                      <cv-select-option value="tuesday">
                        {{ $t("calendar.tuesday") }}
                      </cv-select-option>
                      <cv-select-option value="wednesday">
                        {{ $t("calendar.wednesday") }}
                      </cv-select-option>
                      <cv-select-option value="thursday">
                        {{ $t("calendar.thursday") }}
                      </cv-select-option>
                      <cv-select-option value="friday">
                        {{ $t("calendar.friday") }}
                      </cv-select-option>
                      <cv-select-option value="saturday">
                        {{ $t("calendar.saturday") }}
                      </cv-select-option>
                      <cv-select-option value="sunday">
                        {{ $t("calendar.sunday") }}
                      </cv-select-option>
                    </cv-select>
                  </cv-column>
                  <cv-column :md="2" :max="4">
                    <NsTimePicker
                      hideClearButton
                      v-model="schedule.time"
                      :label="$t('backup.at')"
                    ></NsTimePicker>
                  </cv-column>
                </template>
                <!-- monthly schedule -->
                <template v-if="schedule.interval == 'monthly'">
                  <cv-column :md="2" :max="4">
                    <cv-text-input
                      :label="$t('backup.on_day')"
                      v-model.trim="schedule.monthDay"
                      type="number"
                      min="1"
                      max="31"
                      :disabled="loading.addBackup || loading.alterBackup"
                    >
                    </cv-text-input>
                  </cv-column>
                  <cv-column :md="2" :max="4">
                    <NsTimePicker
                      hideClearButton
                      v-model="schedule.time"
                      :label="$t('backup.at')"
                    ></NsTimePicker>
                  </cv-column>
                </template>
                <cv-column
                  v-if="schedule.interval == 'custom'"
                  :md="6"
                  :max="12"
                >
                  <!-- custom schedule -->
                  <NsTextInput
                    v-model.trim="schedule.custom"
                    :label="$t('backup.calendar_event_expression')"
                    :invalid-message="$t(error.schedule)"
                    :disabled="loading.addBackup || loading.alterBackup"
                    tooltipAlignment="center"
                    tooltipDirection="bottom"
                    ref="schedule-custom"
                  >
                    <template #tooltip>
                      <div class="mg-bottom-sm">
                        {{ $t("backup.custom_schedule_tooltip_description") }}
                      </div>
                      <div class="mg-bottom-sm">
                        {{ $t("backup.custom_schedule_tooltip_more_info") }}
                        <cv-link
                          href="https://www.freedesktop.org/software/systemd/man/systemd.time.html#Calendar%20Events"
                          target="_blank"
                        >
                          {{
                            $t(
                              "backup.custom_schedule_tooltip_systemd_time_documentation"
                            )
                          }}
                        </cv-link>
                      </div>
                    </template>
                  </NsTextInput>
                </cv-column>
              </cv-row>
              <cv-row
                v-if="schedule.interval !== 'custom'"
                class="mg-bottom-xlg"
              >
                <!-- schedule description -->
                <cv-column class="schedule-description">
                  {{ getBackupScheduleDescription(schedule) }}
                </cv-column>
              </cv-row>
              <cv-row>
                <cv-column :md="2" :max="4">
                  <NsTextInput
                    :label="$t('backup.retention')"
                    v-model.trim="retention"
                    type="number"
                    min="1"
                    :helper-text="$t('backup.backup_snapshots')"
                    :invalid-message="$t(error.retention)"
                    :disabled="loading.addBackup || loading.alterBackup"
                    tooltipAlignment="end"
                    tooltipDirection="right"
                    ref="retention"
                    class="retention"
                  >
                    <template #tooltip>{{
                      $t("backup.retention_tooltip")
                    }}</template>
                  </NsTextInput>
                </cv-column>
              </cv-row>
            </cv-grid>
          </div>
          <div v-show="step == 'name'">
            <cv-text-input
              :label="$t('backup.backup_name')"
              v-model.trim="name"
              :helper-text="$t('backup.backup_name_helper')"
              :invalid-message="$t(error.name)"
              :disabled="loading.addBackup || loading.alterBackup"
              class="mg-bottom-xlg"
              ref="name"
            >
            </cv-text-input>
            <!-- run backup on finish -->
            <cv-checkbox
              :label="$t('backup.run_backup_now')"
              v-model="runBackupOnFinish"
              :disabled="loading.addBackup || loading.alterBackup"
              value="checkRunBackupOnFinish"
            />
            <NsInlineNotification
              v-if="error.addBackup"
              kind="error"
              :title="$t('action.add-backup')"
              :description="error.addBackup"
              :showCloseButton="false"
            />
            <NsInlineNotification
              v-if="error.alterBackup"
              kind="error"
              :title="$t('action.alter-backup')"
              :description="error.alterBackup"
              :showCloseButton="false"
            />
          </div>
        </cv-form>
      </template>
    </NsWizard>
  </div>
</template>

<script>
import { UtilService, TaskService, IconService } from "@nethserver/ns8-ui-lib";
import to from "await-to-js";
import BackupInstanceSelector from "@/components/backup/BackupInstanceSelector";
import _cloneDeep from "lodash/cloneDeep";
import _capitalize from "lodash/capitalize";
import Information16 from "@carbon/icons-vue/es/information/16";

const DEFAULT_SCHEDULE_INTERVAL = "daily";
const DEFAULT_SCHEDULE_MINUTE = "0";
const DEFAULT_SCHEDULE_TIME = "00:00";
const DEFAULT_SCHEDULE_WEEK_DAY = "sunday";
const DEFAULT_SCHEDULE_MONTH_DAY = "1";
const DEFAULT_SCHEDULE_CUSTOM = "";
const DEFAULT_RETENTION = "5";

export default {
  name: "CreateOrEditBackupModal",
  components: { BackupInstanceSelector, Information16 },
  mixins: [UtilService, TaskService, IconService],
  props: {
    isShown: {
      type: Boolean,
      default: true,
    },
    isEditing: {
      type: Boolean,
      default: false,
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
    backup: {
      type: Object,
      required: true,
    },
    backups: {
      type: Array,
      required: true,
    },
  },
  data() {
    return {
      step: "",
      instances: [],
      name: "",
      schedule: {
        interval: DEFAULT_SCHEDULE_INTERVAL,
        minute: DEFAULT_SCHEDULE_MINUTE,
        time: DEFAULT_SCHEDULE_TIME,
        weekDay: DEFAULT_SCHEDULE_WEEK_DAY,
        monthDay: DEFAULT_SCHEDULE_MONTH_DAY,
        custom: DEFAULT_SCHEDULE_CUSTOM,
      },
      retention: DEFAULT_RETENTION,
      installedModules: [],
      internalRepositories: [],
      runBackupOnFinish: false,
      loading: {
        addBackup: false,
        alterBackup: false,
        listInstalledModules: true,
      },
      error: {
        name: "",
        repository: "",
        schedule: "",
        retention: "",
        addBackup: "",
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
    isNextButtonDisabled() {
      return (
        this.loading.addBackup ||
        this.loading.alterBackup ||
        !this.installedModules.length ||
        (this.step == "instances" && !this.instances.length) ||
        (this.step == "repository" && !this.selectedRepo)
      );
    },
    nextButtonLabel() {
      if (this.isLastStep) {
        if (this.isEditing) {
          return this.$t("common.save");
        } else {
          return this.$t("backup.schedule_backup");
        }
      } else {
        return this.$t("common.next");
      }
    },
    steps() {
      if (this.isEditing) {
        // cannot change repository
        return ["instances", "settings", "name"];
      } else {
        return ["instances", "repository", "settings", "name"];
      }
    },
    selectedRepo() {
      return this.internalRepositories.find((r) => r.selected);
    },
    isScheduleValid() {
      switch (this.schedule.interval) {
        case "hourly":
          return this.schedule.minute >= 0 && this.schedule.minute <= 59;
        case "daily":
          return this.time24HourPattern.test(this.schedule.time);
        case "weekly":
          return (
            this.schedule.weekDay &&
            this.time24HourPattern.test(this.schedule.time)
          );
        case "monthly":
          return (
            this.schedule.monthDay >= 1 &&
            this.schedule.monthDay <= 31 &&
            this.time24HourPattern.test(this.schedule.time)
          );
      }
      return false;
    },
    scheduleExpression() {
      switch (this.schedule.interval) {
        case "hourly": {
          const minutes = this.schedule.minute.padStart(2, "0");
          return `*-*-* *:${minutes}:00`;
        }
        case "daily":
          return `*-*-* ${this.schedule.time}:00`;
        case "weekly": {
          const shortDay = _capitalize(this.schedule.weekDay.substring(0, 3));
          return `${shortDay} *-*-* ${this.schedule.time}:00`;
        }
        case "monthly": {
          const monthDay = this.schedule.monthDay.padStart(2, "0");
          return `*-*-${monthDay} ${this.schedule.time}:00`;
        }
        case "custom":
          return this.schedule.custom;
      }
      return null;
    },
  },
  watch: {
    isShown: function () {
      if (this.isShown) {
        // show first step
        this.step = this.steps[0];

        // ensure backup prop is updated as well
        this.$nextTick(() => {
          if (this.isEditing) {
            this.name = this.backup.name;
            this.schedule = _cloneDeep(this.backup.schedule);
            this.retention = this.backup.retention.toString();
            this.runBackupOnFinish = false;
            this.updateInternalRepositories();
          } else {
            this.clearFields();
          }
        });

        // load installed moudules
        this.listInstalledModules();
      }
    },
    repositories: function () {
      this.updateInternalRepositories();
    },
    step: function () {
      if (this.step == "settings" && !this.isEditing) {
        // prefill backup name

        let backupName = this.$t("backup.default_backup_name", {
          repository: this.selectedRepo.name,
        });

        // ensure we don't generate an already existing backup name
        let isBackupNameDuplicated = this.backups.find(
          (b) => b.name == backupName
        );
        let backupNameSuffix = 2;

        while (isBackupNameDuplicated) {
          backupName =
            this.$t("backup.default_backup_name", {
              repository: `${this.selectedRepo.name}`,
            }) + ` (${backupNameSuffix})`;

          isBackupNameDuplicated = this.backups.find(
            (b) => b.name == backupName
          );
          backupNameSuffix++;
        }

        this.name = backupName;
      }
    },
  },
  created() {
    this.updateInternalRepositories();
  },
  methods: {
    clearFields() {
      this.instances = [];
      this.name = "";
      this.schedule.interval = DEFAULT_SCHEDULE_INTERVAL;
      this.schedule.minute = DEFAULT_SCHEDULE_MINUTE;
      this.schedule.time = DEFAULT_SCHEDULE_TIME;
      this.schedule.weekDay = DEFAULT_SCHEDULE_WEEK_DAY;
      this.schedule.monthDay = DEFAULT_SCHEDULE_MONTH_DAY;
      this.schedule.custom = DEFAULT_SCHEDULE_CUSTOM;
      this.retention = DEFAULT_RETENTION;
      this.runBackupOnFinish = false;

      for (let repo of this.internalRepositories) {
        repo.selected = false;
      }

      // preselect if there is only one repo
      if (this.internalRepositories.length == 1) {
        this.internalRepositories[0].selected = true;
      }
    },
    nextStep() {
      if (this.isNextButtonDisabled) {
        return;
      }

      if (!this.validateStep()) {
        return;
      }

      if (this.isLastStep) {
        if (this.isEditing) {
          this.alterBackup();
        } else {
          this.addBackup();
        }
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

      if (this.isEditing) {
        // edit backup
        for (const repo of internalRepositories) {
          if (this.backup.repository == repo.id) {
            repo.selected = true;
          } else {
            repo.selected = false;
          }
        }
      } else {
        // create backup
        for (const repo of internalRepositories) {
          repo.selected = false;
        }

        // preselect if there is only one repo
        if (internalRepositories.length == 1) {
          internalRepositories[0].selected = true;
        }
      }

      this.internalRepositories = internalRepositories;
    },
    validateSettingsStep() {
      let isValidationOk = true;
      this.error.retention = "";

      if (!this.retention || this.retention < 1) {
        this.error.retention = "error.invalid_value";

        if (isValidationOk) {
          this.focusElement("retention");
          isValidationOk = false;
        }
      }

      return isValidationOk;
    },
    validateNameStep() {
      let isValidationOk = true;
      this.error.name = "";

      if (!this.name) {
        this.error.name = "common.required";

        if (isValidationOk) {
          this.focusElement("name");
          isValidationOk = false;
        }
      }

      return isValidationOk;
    },
    validateStep() {
      switch (this.step) {
        case "settings":
          return this.validateSettingsStep();
        case "name":
          return this.validateNameStep();
        default:
          return true;
      }
    },
    async addBackup() {
      // validation has already been performed by validateStep

      this.error.addBackup = "";
      this.loading.addBackup = true;
      const taskAction = "add-backup";

      // register to task error
      this.$root.$off(taskAction + "-aborted");
      this.$root.$once(taskAction + "-aborted", this.addBackupAborted);

      // register to task validation
      this.$root.$off(taskAction + "-validation-failed");
      this.$root.$once(
        taskAction + "-validation-failed",
        this.addBackupValidationFailed
      );

      // register to task completion
      this.$root.$off(taskAction + "-completed");
      this.$root.$once(taskAction + "-completed", this.addBackupCompleted);

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            name: this.name,
            repository: this.selectedRepo.id,
            schedule: this.scheduleExpression,
            schedule_hint: this.schedule,
            retention: parseInt(this.retention),
            instances: this.instances,
            enabled: true,
          },
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
            runBackupOnFinish: this.runBackupOnFinish,
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
    addBackupAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.loading.addBackup = false;

      // hide modal
      this.$emit("hide");
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
    addBackupCompleted(taskContext, taskResult) {
      this.loading.addBackup = false;

      // hide modal
      this.$emit("hide");

      const runBackupOnFinish = taskContext.extra.runBackupOnFinish;
      const backupId = taskResult.output;

      const createdBackup = {
        id: backupId,
        name: taskContext.data.name,
      };

      // reload backup configuration
      this.$emit("backupCreated", runBackupOnFinish, createdBackup);
    },
    async alterBackup() {
      // validation has already been performed by validateStep

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
            repository: this.backup.repository,
            schedule: this.scheduleExpression,
            schedule_hint: this.schedule,
            retention: parseInt(this.retention),
            instances: this.instances,
            enabled: this.backup.enabled,
          },
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
            runBackupOnFinish: this.runBackupOnFinish,
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
    alterBackupCompleted(taskContext) {
      this.loading.alterBackup = false;

      // hide modal
      this.$emit("hide");

      const runBackupOnFinish = taskContext.extra.runBackupOnFinish;

      const alteredBackup = {
        id: taskContext.data.id,
        name: taskContext.data.name,
      };

      // reload backup configuration
      this.$emit("backupAltered", runBackupOnFinish, alteredBackup);
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

.schedule-description {
  font-weight: bold;
}

.calendar-event-expression-tooltip {
  display: inline-block;
}

.calendar-event-expression-tooltip-title {
  margin-bottom: $spacing-03;
}
</style>

<style lang="scss">
@import "../../styles/carbon-utils";

// global styles

.calendar-event-expression-tooltip .bx--tooltip__trigger {
  margin-left: 0 !important;
}
</style>

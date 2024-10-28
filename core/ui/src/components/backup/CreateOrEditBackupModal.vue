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
import { ref, reactive, onMounted, onBeforeUnmount } from "vue";
import { useStore } from "vuex";
import { mapState, mapActions } from "vuex";
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
  setup(props, { emit }) {
    const store = useStore();
    const step = ref("");
    const instances = ref([]);
    const name = ref("");
    const schedule = reactive({
      interval: DEFAULT_SCHEDULE_INTERVAL,
      minute: DEFAULT_SCHEDULE_MINUTE,
      time: DEFAULT_SCHEDULE_TIME,
      weekDay: DEFAULT_SCHEDULE_WEEK_DAY,
      monthDay: DEFAULT_SCHEDULE_MONTH_DAY,
      custom: DEFAULT_SCHEDULE_CUSTOM,
    });
    const retention = ref(DEFAULT_RETENTION);
    const installedModules = ref([]);
    const internalRepositories = ref([]);
    const runBackupOnFinish = ref(false);
    const loading = reactive({
      addBackup: false,
      alterBackup: false,
      listInstalledModules: true,
    });
    const error = reactive({
      name: "",
      repository: "",
      schedule: "",
      retention: "",
      addBackup: "",
      alterBackup: "",
      listInstalledModules: "",
    });

    const stepIndex = computed(() => steps.indexOf(step.value));
    const isFirstStep = computed(() => stepIndex.value == 0);
    const isLastStep = computed(() => stepIndex.value == steps.length - 1);
    const isNextButtonDisabled = computed(
      () =>
        loading.addBackup ||
        loading.alterBackup ||
        !installedModules.value.length ||
        (step.value == "instances" && !instances.value.length) ||
        (step.value == "repository" && !selectedRepo.value)
    );
    const nextButtonLabel = computed(() => {
      if (isLastStep.value) {
        if (props.isEditing) {
          return this.$t("common.save");
        } else {
          return this.$t("backup.schedule_backup");
        }
      } else {
        return this.$t("common.next");
      }
    });
    const steps = computed(() => {
      if (props.isEditing) {
        // cannot change repository
        return ["instances", "settings", "name"];
      } else {
        return ["instances", "repository", "settings", "name"];
      }
    });
    const selectedRepo = computed(() =>
      internalRepositories.value.find((r) => r.selected)
    );
    const isScheduleValid = computed(() => {
      switch (schedule.interval) {
        case "hourly":
          return schedule.minute >= 0 && schedule.minute <= 59;
        case "daily":
          return this.time24HourPattern.test(schedule.time);
        case "weekly":
          return (
            schedule.weekDay &&
            this.time24HourPattern.test(schedule.time)
          );
        case "monthly":
          return (
            schedule.monthDay >= 1 &&
            schedule.monthDay <= 31 &&
            this.time24HourPattern.test(schedule.time)
          );
      }
      return false;
    });
    const scheduleExpression = computed(() => {
      switch (schedule.interval) {
        case "hourly": {
          const minutes = schedule.minute.padStart(2, "0");
          return `*-*-* *:${minutes}:00`;
        }
        case "daily":
          return `*-*-* ${schedule.time}:00`;
        case "weekly": {
          const shortDay = _capitalize(schedule.weekDay.substring(0, 3));
          return `${shortDay} *-*-* ${schedule.time}:00`;
        }
        case "monthly": {
          const monthDay = schedule.monthDay.padStart(2, "0");
          return `*-*-${monthDay} ${schedule.time}:00`;
        }
        case "custom":
          return schedule.custom;
      }
      return null;
    });

    watch(
      () => props.isShown,
      (newVal) => {
        if (newVal) {
          // show first step
          step.value = steps.value[0];

          // ensure backup prop is updated as well
          nextTick(() => {
            if (props.isEditing) {
              name.value = props.backup.name;
              Object.assign(schedule, _cloneDeep(props.backup.schedule));
              retention.value = props.backup.retention.toString();
              runBackupOnFinish.value = false;
              updateInternalRepositories();
            } else {
              clearFields();
            }
          });

          // load installed modules
          listInstalledModules();
        }
      }
    );

    watch(
      () => props.repositories,
      () => {
        updateInternalRepositories();
      }
    );

    watch(
      () => step.value,
      (newVal) => {
        if (newVal == "settings" && !props.isEditing) {
          // prefill backup name

          let backupName = this.$t("backup.default_backup_name", {
            repository: selectedRepo.value.name,
          });

          // ensure we don't generate an already existing backup name
          let isBackupNameDuplicated = props.backups.find(
            (b) => b.name == backupName
          );
          let backupNameSuffix = 2;

          while (isBackupNameDuplicated) {
            backupName =
              this.$t("backup.default_backup_name", {
                repository: `${selectedRepo.value.name}`,
              }) + ` (${backupNameSuffix})`;

            isBackupNameDuplicated = props.backups.find(
              (b) => b.name == backupName
            );
            backupNameSuffix++;
          }

          name.value = backupName;
        }
      }
    );

    const clearFields = () => {
      instances.value = [];
      name.value = "";
      schedule.interval = DEFAULT_SCHEDULE_INTERVAL;
      schedule.minute = DEFAULT_SCHEDULE_MINUTE;
      schedule.time = DEFAULT_SCHEDULE_TIME;
      schedule.weekDay = DEFAULT_SCHEDULE_WEEK_DAY;
      schedule.monthDay = DEFAULT_SCHEDULE_MONTH_DAY;
      schedule.custom = DEFAULT_SCHEDULE_CUSTOM;
      retention.value = DEFAULT_RETENTION;
      runBackupOnFinish.value = false;

      for (let repo of internalRepositories.value) {
        repo.selected = false;
      }

      // preselect if there is only one repo
      if (internalRepositories.value.length == 1) {
        internalRepositories.value[0].selected = true;
      }
    };

    const nextStep = () => {
      if (isNextButtonDisabled.value) {
        return;
      }

      if (!validateStep()) {
        return;
      }

      if (isLastStep.value) {
        if (props.isEditing) {
          alterBackup();
        } else {
          addBackup();
        }
      } else {
        step.value = steps.value[stepIndex.value + 1];
      }
    };

    const previousStep = () => {
      if (!isFirstStep.value) {
        step.value = steps.value[stepIndex.value - 1];
      }
    };

    const updateInternalRepositories = () => {
      // deep copy (needed to avoid reactivity issues)
      let internalRepos = _cloneDeep(props.repositories);

      if (props.isEditing) {
        // edit backup
        for (const repo of internalRepos) {
          if (props.backup.repository == repo.id) {
            repo.selected = true;
          } else {
            repo.selected = false;
          }
        }
      } else {
        // create backup
        for (const repo of internalRepos) {
          repo.selected = false;
        }

        // preselect if there is only one repo
        if (internalRepos.length == 1) {
          internalRepos[0].selected = true;
        }
      }

      internalRepositories.value = internalRepos;
    };

    const validateSettingsStep = () => {
      let isValidationOk = true;
      error.retention = "";

      if (!retention.value || retention.value < 1) {
        error.retention = "error.invalid_value";

        if (isValidationOk) {
          focusElement("retention");
          isValidationOk = false;
        }
      }

      return isValidationOk;
    };

    const validateNameStep = () => {
      let isValidationOk = true;
      error.name = "";

      if (!name.value) {
        error.name = "common.required";

        if (isValidationOk) {
          focusElement("name");
          isValidationOk = false;
        }
      }

      return isValidationOk;
    };

    const validateStep = () => {
      switch (step.value) {
        case "settings":
          return validateSettingsStep();
        case "name":
          return validateNameStep();
        default:
          return true;
      }
    };

    const addBackup = async () => {
      // validation has already been performed by validateStep

      error.addBackup = "";
      loading.addBackup = true;
      const taskAction = "add-backup";

      // register to task error
      store.$off(taskAction + "-aborted");
      store.$once(taskAction + "-aborted", addBackupAborted);

      // register to task validation
      store.$off(taskAction + "-validation-failed");
      store.$once(
        taskAction + "-validation-failed",
        addBackupValidationFailed
      );

      // register to task completion
      store.$off(taskAction + "-completed");
      store.$once(taskAction + "-completed", addBackupCompleted);

      const res = await to(
        createClusterTask({
          action: taskAction,
          data: {
            name: name.value,
            repository: selectedRepo.value.id,
            schedule: scheduleExpression.value,
            schedule_hint: schedule,
            retention: parseInt(retention.value),
            instances: instances.value,
            enabled: true,
          },
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
            runBackupOnFinish: runBackupOnFinish.value,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        error.addBackup = this.getErrorMessage(err);
        return;
      }
    };

    const addBackupAborted = (taskResult, taskContext) => {
      console.error(`${taskContext.action} aborted`, taskResult);
      loading.addBackup = false;

      // hide modal
      emit("hide");
    };

    const addBackupValidationFailed = (validationErrors) => {
      loading.addBackup = false;
      let focusAlreadySet = false;

      for (const validationError of validationErrors) {
        const param = validationError.parameter;

        // set i18n error message
        error[param] = "backup." + validationError.error;

        if (!focusAlreadySet) {
          focusElement(param);
          focusAlreadySet = true;
        }
      }
    };

    const addBackupCompleted = (taskContext, taskResult) => {
      loading.addBackup = false;

      // hide modal
      emit("hide");

      const runBackupOnFinish = taskContext.extra.runBackupOnFinish;
      const backupId = taskResult.output;

      const createdBackup = {
        id: backupId,
        name: taskContext.data.name,
      };

      // reload backup configuration
      emit("backupCreated", runBackupOnFinish, createdBackup);
    };

    const alterBackup = async () => {
      // validation has already been performed by validateStep

      error.alterBackup = "";
      loading.alterBackup = true;
      const taskAction = "alter-backup";

      // register to task validation
      store.$off(taskAction + "-validation-failed");
      store.$once(
        taskAction + "-validation-failed",
        alterBackupValidationFailed
      );

      // register to task completion
      store.$off(taskAction + "-completed");
      store.$once(taskAction + "-completed", alterBackupCompleted);

      const res = await to(
        createClusterTask({
          action: taskAction,
          data: {
            id: props.backup.id,
            name: name.value,
            repository: props.backup.repository,
            schedule: scheduleExpression.value,
            schedule_hint: schedule,
            retention: parseInt(retention.value),
            instances: instances.value,
            enabled: props.backup.enabled,
          },
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
            runBackupOnFinish: runBackupOnFinish.value,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        error.alterBackup = this.getErrorMessage(err);
        return;
      }
    };

    const alterBackupValidationFailed = (validationErrors) => {
      loading.alterBackup = false;
      let focusAlreadySet = false;

      for (const validationError of validationErrors) {
        const param = validationError.parameter;

        // set i18n error message
        error[param] = "backup." + validationError.error;

        if (!focusAlreadySet) {
          focusElement(param);
          focusAlreadySet = true;
        }
      }
    };

    const alterBackupCompleted = (taskContext) => {
      loading.alterBackup = false;

      // hide modal
      emit("hide");

      const runBackupOnFinish = taskContext.extra.runBackupOnFinish;

      const alteredBackup = {
        id: taskContext.data.id,
        name: taskContext.data.name,
      };

      // reload backup configuration
      emit("backupAltered", runBackupOnFinish, alteredBackup);
    };

    const listInstalledModules = async () => {
      loading.listInstalledModules = true;
      const taskAction = "list-installed-modules";

      // register to task completion
      store.$once(
        taskAction + "-completed",
        listInstalledModulesCompleted
      );

      const res = await to(
        createClusterTask({
          action: taskAction,
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
          },
        })
      );
      const err = res[0];

      if (err) {
        error.listInstalledModules = this.getErrorMessage(err);
        createErrorNotification(
          err,
          this.$t("task.cannot_create_task", { action: taskAction })
        );
        return;
      }
    };

    const listInstalledModulesCompleted = (taskContext, taskResult) => {
      loading.listInstalledModules = false;
      let apps = [];

      for (let instanceList of Object.values(taskResult.output)) {
        for (let instance of instanceList) {
          apps.push(instance);
        }
      }
      apps.sort(this.sortModuleInstances());
      installedModules.value = apps;
    };

    const onSelectInstances = (instances) => {
      instances.value = instances;
    };

    const deselectOtherRepos = (repo) => {
      for (let r of internalRepositories.value) {
        if (r.id !== repo.id) {
          r.selected = false;
        }
      }
    };

    return {
      step,
      instances,
      name,
      schedule,
      retention,
      installedModules,
      internalRepositories,
      runBackupOnFinish,
      loading,
      error,
      stepIndex,
      isFirstStep,
      isLastStep,
      isNextButtonDisabled,
      nextButtonLabel,
      steps,
      selectedRepo,
      isScheduleValid,
      scheduleExpression,
      clearFields,
      nextStep,
      previousStep,
      updateInternalRepositories,
      validateSettingsStep,
      validateNameStep,
      validateStep,
      addBackup,
      addBackupAborted,
      addBackupValidationFailed,
      addBackupCompleted,
      alterBackup,
      alterBackupValidationFailed,
      alterBackupCompleted,
      listInstalledModules,
      listInstalledModulesCompleted,
      onSelectInstances,
      deselectOtherRepos,
    };
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

<template>
  <!-- //// use standard modal instead of wizard? (only one step) -->
  <cv-modal
    size="default"
    :visible="isShown"
    @modal-hidden="$emit('hide')"
    class="wizard-modal"
  >
    <template slot="title">{{ $t("backup.edit_backup_repository") }}</template>
    <template slot="content">
      <cv-form @submit.prevent="nextStep">
        <template v-if="step == 'settings'">
          <cv-text-input
            :label="$t('backup.url')"
            v-model.trim="repository.url"
            disabled
            ref="url"
          >
          </cv-text-input>
          <!-- backblaze -->
          <template v-if="repository.provider == 'backblaze'">
            <cv-text-input
              :label="$t('backup.b2_account_id')"
              v-model.trim="repository.parameters.b2_account_id"
              :invalid-message="$t(error.backblaze.b2_account_id)"
              :disabled="loading.alterBackupRepository"
              ref="b2_account_id"
            >
            </cv-text-input>
            <cv-text-input
              :label="$t('backup.b2_account_key')"
              v-model.trim="repository.parameters.b2_account_key"
              :invalid-message="$t(error.backblaze.b2_account_key)"
              :disabled="loading.alterBackupRepository"
              ref="b2_account_key"
            >
            </cv-text-input>
          </template>
          <!-- amazon s3 -->
          <template v-if="repository.provider == 'aws'">
            <cv-text-input
              :label="$t('backup.aws_access_key_id')"
              v-model.trim="repository.parameters.aws_access_key_id"
              :invalid-message="$t(error.aws.aws_access_key_id)"
              :disabled="loading.alterBackupRepository"
              ref="aws_access_key_id"
            >
            </cv-text-input>
            <cv-text-input
              :label="$t('backup.aws_default_region')"
              v-model.trim="repository.parameters.aws_default_region"
              :invalid-message="$t(error.aws.aws_default_region)"
              :disabled="loading.alterBackupRepository"
              ref="aws_default_region"
            >
            </cv-text-input>
            <cv-text-input
              :label="$t('backup.aws_secret_access_key')"
              v-model.trim="repository.parameters.aws_secret_access_key"
              :invalid-message="$t(error.aws.aws_secret_access_key)"
              :disabled="loading.alterBackupRepository"
              ref="aws_secret_access_key"
            >
            </cv-text-input>
          </template>
          <!-- generic s3 -->
          <template v-if="repository.provider == 'genericS3'">
            generic s3 ////
          </template>
          <!-- azure -->
          <template v-if="repository.provider == 'azure'">
            azure ////
          </template>
          <!-- //// handle ALL providers -->
          <cv-text-input
            :label="$t('backup.repository_name')"
            v-model.trim="repository.name"
            :helper-text="$t('backup.repository_name_helper')"
            :invalid-message="$t(error.name)"
            :disabled="loading.alterBackupRepository"
            ref="name"
          >
          </cv-text-input>
          <NsInlineNotification
            v-if="error.addBackupRepository"
            kind="error"
            :title="$t('action.add-backup-repository')"
            :description="error.alterBackupRepository"
            :showCloseButton="false"
          />
        </template>
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
            :disabled="isFirstStep || loading.alterBackupRepository"
            type="button"
            class="wizard-button"
            >{{ $t("common.previous") }}
          </NsButton>
          <NsButton
            kind="primary"
            :icon="ChevronRight20"
            :disabled="isNextStepDisabled"
            :loading="loading.alterBackupRepository"
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

//// review all (copy/paste)

export default {
  name: "EditRepositoryModal",
  mixins: [UtilService, TaskService, IconService],
  props: {
    isShown: {
      type: Boolean,
      default: true,
    },
    repository: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      step: "",
      steps: ["settings"],
      name: "",
      backblaze: {
        b2_account_id: "",
        b2_account_key: "",
      },
      aws: {
        aws_access_key_id: "",
        aws_default_region: "",
        aws_secret_access_key: "",
      },
      genericS3: {},
      azure: {},
      //// handle all providers
      loading: {
        alterBackupRepository: false,
      },
      error: {
        name: "",
        backblaze: {
          b2_account_id: "",
          b2_account_key: "",
        },
        aws: {
          aws_access_key_id: "",
          aws_default_region: "",
          aws_secret_access_key: "",
        },
        genericS3: {},
        azure: {},
        //// handle all providers
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
      return this.loading.alterBackupRepository;
    },
  },
  watch: {
    isShown: function () {
      if (this.isShown) {
        // show first step
        this.step = this.steps[0];
        this.clearFields();
      }
    },
  },
  methods: {
    clearFields() {
      this.name = "";

      this.backblaze.b2_account_id = "";
      this.backblaze.b2_account_key = "";

      this.aws.aws_access_key_id = "";
      this.aws.aws_default_region = "";
      this.aws.aws_secret_access_key = "";

      //// handle ALL providers
    },
    nextStep() {
      if (this.isNextStepDisabled) {
        return;
      }

      if (this.isLastStep) {
        this.alterBackupRepository();
      } else {
        this.step = this.steps[this.stepIndex + 1];
      }
    },
    previousStep() {
      if (!this.isFirstStep) {
        this.step = this.steps[this.stepIndex - 1];
      }
    },
    buildRepositoryParameters() {
      switch (this.repository.provider) {
        case "backblaze":
          return {
            b2_account_id: this.backblaze.b2_account_id,
            b2_account_key: this.backblaze.b2_account_key,
          };
        case "aws":
          return {
            aws_default_region: this.aws.aws_default_region,
            aws_access_key_id: this.aws.aws_access_key_id,
            aws_secret_access_key: this.aws.aws_secret_access_key,
          };
        case "azure":
          return {
            azure_account_name: "", ////
            azure_account_key: "",
          };
      }
      //// handle all providers
    },
    async alterBackupRepository() {
      //// validation
      this.loading.alterBackupRepository = true;
      this.error.alterBackupRepository = "";
      const taskAction = "alter-backup-repository";

      // register to task validation
      this.$root.$off(taskAction + "-validation-failed");
      this.$root.$once(
        taskAction + "-validation-failed",
        this.alterBackupRepositoryValidationFailed
      );

      // register to task completion
      this.$root.$once(
        taskAction + "-completed",
        this.alterBackupRepositoryCompleted
      );

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            id: this.repository.id,
            name: this.name,
            parameters: this.buildRepositoryParameters(),
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
        this.error.alterBackupRepository = this.getErrorMessage(err);
        return;
      }
    },
    alterBackupRepositoryValidationFailed(validationErrors) {
      this.loading.alterBackupRepository = false;
      let focusAlreadySet = false;

      for (const validationError of validationErrors) {
        const param = validationError.parameter;
        console.log("param", param); ////

        // set i18n error message
        this.error[param] = "backup." + validationError.error;

        if (!focusAlreadySet) {
          this.focusElement(param);
          focusAlreadySet = true;
        }
      }
    },
    alterBackupRepositoryCompleted(taskContext, taskResult) {
      console.log("alterBackupRepositoryCompleted", taskResult.output); ////

      this.loading.alterBackupRepository = false;
      this.$emit("repoEdited");
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";

.min-height-card {
  min-height: 8rem;
}
</style>

<template>
  <NsWizard
    size="default"
    :visible="isShown"
    :cancelLabel="$t('common.cancel')"
    :previousLabel="$t('common.previous')"
    :nextLabel="isLastStep ? $t('backup.add_repository') : $t('common.next')"
    :isPreviousDisabled="isFirstStep || loading.addBackupRepository"
    :isNextDisabled="isNextButtonDisabled"
    :isNextLoading="loading.addBackupRepository"
    @modal-hidden="$emit('hide')"
    @cancel="$emit('hide')"
    @previousStep="previousStep"
    @nextStep="nextStep"
  >
    <template slot="title">{{ $t("backup.add_repository") }}</template>
    <template slot="content">
      <cv-form>
        <template v-if="step == 'provider'">
          <div class="mg-bottom-md">
            {{ $t("backup.select_backup_provider") }}
          </div>
          <div class="bx--grid">
            <div class="bx--row">
              <!-- backblaze -->
              <div class="bx--col-md-4">
                <NsTile
                  :light="true"
                  kind="selectable"
                  v-model="isBackblazeSelected"
                  value="providerValue"
                  @click="selectBackblaze()"
                  class="provider-card"
                >
                  <div class="provider-card-content">
                    <div class="provider-icon">
                      <img
                        :src="require('@/assets/backblaze.png')"
                        alt="backblaze logo"
                      />
                    </div>
                    <h6>
                      {{ $t("backup.backblaze") }}
                    </h6>
                  </div>
                </NsTile>
              </div>
              <!-- azure -->
              <div class="bx--col-md-4">
                <NsTile
                  :light="true"
                  kind="selectable"
                  v-model="isAzureSelected"
                  value="providerValue"
                  @click="selectAzure()"
                  class="provider-card"
                  disabled
                >
                  <div class="provider-card-content">
                    <div class="provider-icon">
                      <img
                        :src="require('@/assets/azure.png')"
                        alt="azure logo"
                      />
                    </div>
                    <h6>
                      {{ $t("backup.azure") }}
                    </h6>
                  </div>
                </NsTile>
              </div>
              <!-- amazon s3 -->
              <div class="bx--col-md-4">
                <NsTile
                  :light="true"
                  kind="selectable"
                  v-model="isAmazonS3Selected"
                  value="providerValue"
                  @click="selectAmazonS3()"
                  class="provider-card"
                >
                  <div class="provider-card-content">
                    <div class="provider-icon">
                      <img
                        :src="require('@/assets/s3.png')"
                        alt="amazon s3 logo"
                      />
                    </div>
                    <h6>
                      {{ $t("backup.aws") }}
                    </h6>
                  </div>
                </NsTile>
              </div>
              <!-- generic s3 -->
              <div class="bx--col-md-4">
                <NsTile
                  :light="true"
                  kind="selectable"
                  v-model="isGenericS3Selected"
                  value="providerValue"
                  @click="selectGenericS3()"
                  class="provider-card"
                  disabled
                >
                  <div class="provider-card-content">
                    <div class="provider-icon">
                      <img
                        :src="require('@/assets/s3-generic.png')"
                        alt="generic s3 logo"
                      />
                    </div>
                    <h6>
                      {{ $t("backup.generic_s3") }}
                    </h6>
                  </div>
                </NsTile>
              </div>
            </div>
          </div>
        </template>
        <template v-if="step == 'settings'">
          <NsInlineNotification
            v-if="error.repoConnection"
            kind="error"
            :title="$t('backup.backup_repository_auth_error')"
            :description="$t('backup.backup_repository_auth_error_description')"
            :showCloseButton="false"
          />
          <cv-text-input
            :label="$t('backup.url')"
            v-model.trim="url"
            :invalid-message="error.url"
            :placeholder="urlPlaceholder"
            :disabled="loading.addBackupRepository"
            ref="url"
          >
          </cv-text-input>
          <!-- backblaze -->
          <template v-if="isBackblazeSelected">
            <cv-text-input
              :label="$t('backup.b2_account_id')"
              v-model.trim="backblaze.b2_account_id"
              :invalid-message="error.backblaze.b2_account_id"
              :disabled="loading.addBackupRepository"
              ref="b2_account_id"
            >
            </cv-text-input>
            <cv-text-input
              :label="$t('backup.b2_account_key')"
              type="password"
              v-model.trim="backblaze.b2_account_key"
              :disabled="loading.addBackupRepository"
              :invalid-message="error.backblaze.b2_account_key"
              :password-hide-label="$t('password.hide_password')"
              :password-show-label="$t('password.show_password')"
              ref="b2_account_key"
            ></cv-text-input>
          </template>
          <!-- amazon s3 -->
          <template v-if="isAmazonS3Selected">
            <cv-text-input
              :label="$t('backup.aws_access_key_id')"
              v-model.trim="aws.aws_access_key_id"
              :invalid-message="error.aws.aws_access_key_id"
              :disabled="loading.addBackupRepository"
              ref="aws_access_key_id"
            >
            </cv-text-input>
            <cv-text-input
              :label="$t('backup.aws_default_region')"
              v-model.trim="aws.aws_default_region"
              :invalid-message="error.aws.aws_default_region"
              :disabled="loading.addBackupRepository"
              ref="aws_default_region"
            >
            </cv-text-input>
            <cv-text-input
              :label="$t('backup.aws_secret_access_key')"
              type="password"
              v-model.trim="aws.aws_secret_access_key"
              :disabled="loading.addBackupRepository"
              :invalid-message="error.aws.aws_secret_access_key"
              :password-hide-label="$t('password.hide_password')"
              :password-show-label="$t('password.show_password')"
              ref="aws_secret_access_key"
            ></cv-text-input>
          </template>
          <!-- azure -->
          <template v-if="isAzureSelected"> azure //// </template>
          <!-- generic s3 -->
          <template v-if="isGenericS3Selected"> generic s3 //// </template>
          <!-- //// handle ALL providers -->
          <cv-text-input
            :label="$t('backup.repository_name')"
            v-model.trim="name"
            :helper-text="$t('backup.repository_name_helper')"
            :invalid-message="error.name"
            :disabled="loading.addBackupRepository"
            ref="name"
          >
          </cv-text-input>
          <!-- advanced options -->
          <cv-accordion ref="accordion">
            <cv-accordion-item :open="toggleAccordion[0]">
              <template slot="title">{{ $t("common.advanced") }}</template>
              <template slot="content">
                <cv-text-input
                  :label="$t('backup.repository_password')"
                  v-model.trim="password"
                  :helper-text="$t('backup.repository_password_helper')"
                  :invalid-message="error.password"
                  :disabled="loading.addBackupRepository"
                  ref="password"
                >
                </cv-text-input>
              </template>
            </cv-accordion-item>
          </cv-accordion>
          <NsInlineNotification
            v-if="error.addBackupRepository"
            kind="error"
            :title="$t('action.add-backup-repository')"
            :description="error.addBackupRepository"
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
  name: "AddRepositoryModal",
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
  },
  data() {
    return {
      BACKBLAZE_PROTOCOL: "b2:",
      AWS_PROTOCOL: "s3:",
      AZURE_PROTOCOL: "azure:",
      step: "",
      steps: ["provider", "settings"],
      isBackblazeSelected: false,
      isAzureSelected: false,
      isAmazonS3Selected: false,
      isGenericS3Selected: false,
      name: "",
      url: "",
      password: "",
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
        addBackupRepository: false,
      },
      error: {
        name: "",
        url: "",
        addBackupRepository: "",
        repoConnection: "",
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
    isNextButtonDisabled() {
      return this.loading.addBackupRepository || !this.selectedProvider;
    },
    selectedProvider() {
      if (this.isBackblazeSelected) {
        return "backblaze";
      } else if (this.isAmazonS3Selected) {
        return "aws";
      } else if (this.isGenericS3Selected) {
        return "genericS3";
      } else if (this.isAzureSelected) {
        return "azure";
      } else {
        return null;
      }
      //// handle all providers
    },
    urlPlaceholder() {
      const suffix = this.$t("backup.url_placeholder");

      if (this.isBackblazeSelected) {
        return this.BACKBLAZE_PROTOCOL + suffix;
      } else if (this.isAmazonS3Selected || this.isGenericS3Selected) {
        return this.AWS_PROTOCOL + suffix;
      } else if (this.isAzureSelected) {
        return this.AZURE_PROTOCOL + suffix;
      } else {
        return "";
      }
    },
  },
  watch: {
    isShown: function () {
      if (this.isShown) {
        // show first step
        this.step = this.steps[0];
        this.clearFields();
        this.clearErrors();
      }
    },
    step: function () {
      if (this.step == "settings") {
        // prefill repository name

        let repoName = this.$t("backup.default_repository_name", {
          provider: this.getProviderShortName(),
        });

        // ensure we don't generate an already existing repo name
        let isRepoNameDuplicated = this.repositories.find(
          (b) => b.name == repoName
        );
        let repoNameSuffix = 2;

        while (isRepoNameDuplicated) {
          repoName =
            this.$t("backup.default_repository_name", {
              provider: `${this.getProviderShortName()}`,
            }) + ` (${repoNameSuffix})`;

          isRepoNameDuplicated = this.repositories.find(
            (b) => b.name == repoName
          );
          repoNameSuffix++;
        }

        this.name = repoName;
      }
    },
  },
  methods: {
    clearFields() {
      this.isBackblazeSelected = false;
      this.isAmazonS3Selected = false;
      this.isAzureSelected = false;
      this.name = "";
      this.url = "";

      this.backblaze.b2_account_id = "";
      this.backblaze.b2_account_key = "";

      this.aws.aws_access_key_id = "";
      this.aws.aws_default_region = "";
      this.aws.aws_secret_access_key = "";

      //// handle ALL providers
    },
    nextStep() {
      if (this.isNextButtonDisabled) {
        return;
      }

      if (this.isLastStep) {
        this.addBackupRepository();
      } else {
        this.step = this.steps[this.stepIndex + 1];
      }
    },
    previousStep() {
      if (!this.isFirstStep) {
        this.step = this.steps[this.stepIndex - 1];
      }
    },
    selectBackblaze() {
      //// handle ALL providers
      this.isAzureSelected = false;
      this.isAmazonS3Selected = false;
      this.isGenericS3Selected = false;
    },
    selectAmazonS3() {
      //// handle ALL providers
      this.isBackblazeSelected = false;
      this.isAzureSelected = false;
      this.isGenericS3Selected = false;
    },
    selectGenericS3() {
      //// handle ALL providers
      this.isBackblazeSelected = false;
      this.isAzureSelected = false;
      this.isAmazonS3Selected = false;
    },
    selectAzure() {
      //// handle ALL providers
      this.isBackblazeSelected = false;
      this.isAmazonS3Selected = false;
      this.isGenericS3Selected = false;
    },
    buildRepositoryParameters() {
      switch (this.selectedProvider) {
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
    validateAddBackblazeRepository() {
      // clear errors
      this.error.name = "";
      this.error.url = ""; //// ?
      this.error.repoConnection = "";

      this.error.backblaze.b2_account_id = "";
      this.error.backblaze.b2_account_key = "";

      let isValidationOk = true;

      if (!this.url) {
        this.error.url = this.$t("common.required");

        if (isValidationOk) {
          this.focusElement("url");
          isValidationOk = false;
        }
      } else if (!this.url.startsWith(this.BACKBLAZE_PROTOCOL)) {
        // wrong url protocol
        this.error.url = this.$t("backup.invalid_url_protocol", {
          protocol: this.BACKBLAZE_PROTOCOL,
        });

        if (isValidationOk) {
          this.focusElement("url");
          isValidationOk = false;
        }
      }

      if (!this.backblaze.b2_account_id) {
        this.error.backblaze.b2_account_id = this.$t("common.required");

        if (isValidationOk) {
          this.focusElement("b2_account_id");
          isValidationOk = false;
        }
      }

      if (!this.backblaze.b2_account_key) {
        this.error.backblaze.b2_account_key = this.$t("common.required");

        if (isValidationOk) {
          this.focusElement("b2_account_key");
          isValidationOk = false;
        }
      }

      if (!this.name) {
        this.error.name = this.$t("common.required");

        if (isValidationOk) {
          this.focusElement("name");
          isValidationOk = false;
        }
      }
      return isValidationOk;
    },
    validateAmazonS3Repository() {
      // clear errors
      this.error.name = "";
      this.error.url = ""; //// ?
      this.error.repoConnection = "";

      this.error.aws.aws_access_key_id = "";
      this.error.aws.aws_default_region = "";
      this.error.aws.aws_secret_access_key = "";

      let isValidationOk = true;

      if (!this.url) {
        this.error.url = this.$t("common.required");

        if (isValidationOk) {
          this.focusElement("url");
          isValidationOk = false;
        }
      } else if (!this.url.startsWith(this.AWS_PROTOCOL)) {
        // wrong url protocol
        this.error.url = this.$t("backup.invalid_url_protocol", {
          protocol: this.AWS_PROTOCOL,
        });

        if (isValidationOk) {
          this.focusElement("url");
          isValidationOk = false;
        }
      }

      if (!this.aws.aws_access_key_id) {
        this.error.aws.aws_access_key_id = this.$t("common.required");

        if (isValidationOk) {
          this.focusElement("aws_access_key_id");
          isValidationOk = false;
        }
      }

      if (!this.aws.aws_default_region) {
        this.error.aws.aws_default_region = this.$t("common.required");

        if (isValidationOk) {
          this.focusElement("aws_default_region");
          isValidationOk = false;
        }
      }

      if (!this.aws.aws_secret_access_key) {
        this.error.aws.aws_secret_access_key = this.$t("common.required");

        if (isValidationOk) {
          this.focusElement("aws_secret_access_key");
          isValidationOk = false;
        }
      }

      if (!this.name) {
        this.error.name = this.$t("common.required");

        if (isValidationOk) {
          this.focusElement("name");
          isValidationOk = false;
        }
      }
      return isValidationOk;
    },
    validateAddGenericS3Repository() {
      // clear errors
      this.error.name = "";
      this.error.url = ""; //// ?
      this.error.repoConnection = "";

      ////
      console.error("not implemented"); ////
      return false;
    },
    validateAddAzureRepository() {
      // clear errors
      this.error.name = "";
      this.error.url = ""; //// ?
      this.error.repoConnection = "";

      ////
      console.error("not implemented"); ////
      return false;
    },
    validateAddBackupRepository() {
      switch (this.selectedProvider) {
        case "backblaze":
          return this.validateAddBackblazeRepository();
        case "azure":
          return this.validateAddAzureRepository();
        case "aws":
          return this.validateAmazonS3Repository();
        case "genericS3":
          return this.validateAddGenericS3Repository();
      }
    },
    async addBackupRepository() {
      if (!this.validateAddBackupRepository()) {
        return;
      }
      this.error.addBackupRepository = "";
      this.loading.addBackupRepository = true;
      const taskAction = "add-backup-repository";

      // register to task validation
      this.$root.$off(taskAction + "-validation-failed");
      this.$root.$once(
        taskAction + "-validation-failed",
        this.addBackupRepositoryValidationFailed
      );

      // register to task completion
      this.$root.$once(
        taskAction + "-completed",
        this.addBackupRepositoryCompleted
      );

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            name: this.name,
            provider: this.selectedProvider,
            url: this.url,
            parameters: this.buildRepositoryParameters(),
            password: this.password,
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
        this.error.addBackupRepository = this.getErrorMessage(err);
        return;
      }
    },
    addBackupRepositoryValidationFailed(validationErrors) {
      this.loading.addBackupRepository = false;
      let focusAlreadySet = false;

      for (const validationError of validationErrors) {
        const param = validationError.parameter;

        if (validationError.error == "backup_repository_not_accessible") {
          // show error nontification
          this.error.repoConnection = "error";
        } else {
          // set i18n error message
          this.error[param] = this.$t("backup." + validationError.error);

          if (!focusAlreadySet) {
            this.focusElement(param);
            focusAlreadySet = true;
          }
        }
      }
    },
    addBackupRepositoryCompleted() {
      this.loading.addBackupRepository = false;
      this.$emit("repoCreated");
      this.$emit("hide");
    },
    getProviderShortName() {
      return this.$t(`backup.${this.selectedProvider}_short`);
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";

.provider-card {
  min-height: 9rem;
  padding-right: $spacing-05;
}

.provider-card-content {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.provider-icon {
  width: 3rem;
  height: 3rem;
  margin-bottom: $spacing-06;
}

.provider-icon img {
  width: 100%;
  height: 100%;
}
</style>

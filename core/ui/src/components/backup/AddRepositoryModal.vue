<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
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
          <cv-grid>
            <cv-row>
              <!-- backblaze -->
              <cv-column :md="4">
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
              </cv-column>
              <!-- azure -->
              <cv-column :md="4">
                <NsTile
                  :light="true"
                  kind="selectable"
                  v-model="isAzureSelected"
                  value="providerValue"
                  @click="selectAzure()"
                  class="provider-card"
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
              </cv-column>
              <!-- amazon s3 -->
              <cv-column :md="4">
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
              </cv-column>
              <!-- generic s3 -->
              <cv-column :md="4">
                <NsTile
                  :light="true"
                  kind="selectable"
                  v-model="isGenericS3Selected"
                  value="providerValue"
                  @click="selectGenericS3()"
                  class="provider-card"
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
              </cv-column>
              <!-- samba -->
              <cv-column :md="4">
                <NsTile
                  :light="true"
                  kind="selectable"
                  v-model="isSambaSelected"
                  value="providerValue"
                  @click="selectSamba()"
                  class="provider-card"
                >
                  <div class="provider-card-content">
                    <div class="provider-icon">
                      <img
                        :src="require('@/assets/samba.png')"
                        alt="samba logo"
                      />
                    </div>
                    <h6>
                      {{ $t("backup.samba") }}
                    </h6>
                  </div>
                </NsTile>
              </cv-column>
            </cv-row>
          </cv-grid>
        </template>
        <template v-if="step == 'settings'">
          <NsInlineNotification
            v-if="error.repoConnection"
            kind="error"
            :title="$t('backup.backup_repository_auth_error')"
            :description="$t('backup.backup_repository_auth_error_description')"
            :showCloseButton="false"
          />
          <NsTextInput
            :label="$t('backup.url')"
            v-model.trim="url"
            :invalid-message="error.url"
            :placeholder="$t('backup.' + selectedProviderHelper)"
            :disabled="loading.addBackupRepository"
            :prefix="selectedProviderPrefix"
            ref="url"
          >
          </NsTextInput>
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
          <template v-if="isAzureSelected">
            <cv-text-input
              :label="$t('backup.azure_account_name')"
              v-model.trim="azure.azure_account_name"
              :invalid-message="error.azure.azure_account_name"
              :disabled="loading.addBackupRepository"
              ref="azure_account_name"
            >
            </cv-text-input>
            <cv-text-input
              :label="$t('backup.azure_account_key')"
              type="password"
              v-model.trim="azure.azure_account_key"
              :disabled="loading.addBackupRepository"
              :invalid-message="error.aws.azure_account_key"
              :password-hide-label="$t('password.hide_password')"
              :password-show-label="$t('password.show_password')"
              ref="azure_account_key"
            ></cv-text-input>
          </template>
          <!-- generic s3 -->
          <template v-if="isGenericS3Selected">
            <cv-text-input
              :label="$t('backup.genericS3_access_key_id')"
              v-model.trim="genericS3.aws_access_key_id"
              :invalid-message="error.genericS3.aws_access_key_id"
              :disabled="loading.addBackupRepository"
              ref="genericS3_access_key_id"
            >
            </cv-text-input>
            <cv-text-input
              :label="$t('backup.genericS3_secret_access_key')"
              type="password"
              v-model.trim="genericS3.aws_secret_access_key"
              :disabled="loading.addBackupRepository"
              :invalid-message="error.genericS3.aws_secret_access_key"
              :password-hide-label="$t('password.hide_password')"
              :password-show-label="$t('password.show_password')"
              ref="genericS3_secret_access_key"
            ></cv-text-input>
          </template>
          <!-- samba -->
          <template v-if="isSambaSelected">
            <cv-text-input
              :label="$t('backup.smb_host')"
              v-model.trim="smb.smb_host"
              :invalid-message="error.smb.smb_host"
              :disabled="loading.addBackupRepository"
              ref="smb_host"
            >
            </cv-text-input>
            <cv-text-input
              :label="$t('backup.smb_domain')"
              v-model.trim="smb.smb_domain"
              :invalid-message="error.smb.smb_domain"
              :disabled="loading.addBackupRepository"
              ref="smb_domain"
            >
            </cv-text-input>
            <cv-text-input
              :label="$t('backup.smb_user')"
              v-model.trim="smb.smb_user"
              :invalid-message="error.smb.smb_user"
              :disabled="loading.addBackupRepository"
              ref="smb_user"
            >
            </cv-text-input>
            <cv-text-input
              :label="$t('backup.smb_pass')"
              type="password"
              v-model.trim="smb.smb_pass"
              :disabled="loading.addBackupRepository"
              :invalid-message="error.smb.smb_pass"
              :password-hide-label="$t('password.hide_password')"
              :password-show-label="$t('password.show_password')"
              ref="smb_pass"
            ></cv-text-input>
          </template>
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
                  v-model="password"
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
      isSambaSelected: false,
      name: "",
      url: "",
      password: "",
      backblaze: {
        b2_account_id: "",
        b2_account_key: "",
        repoPrefix: "b2:",
      },
      aws: {
        aws_access_key_id: "",
        aws_default_region: "",
        aws_secret_access_key: "",
        repoPrefix: "s3:s3.amazonaws.com/",
      },
      smb: {
        smb_host: "",
        smb_user: "",
        smb_pass: "",
        smb_domain: "",
        repoPrefix: "smb:",
      },
      genericS3: {
        aws_access_key_id: "",
        aws_secret_access_key: "",
        repoPrefix: "s3:",
      },
      azure: {
        azure_account_key: "",
        azure_account_name: "",
        repoPrefix: "azure:",
      },
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
        smb: {
          smb_host: "",
          smb_user: "",
          smb_pass: "",
          smb_domain: "",
        },
        genericS3: {
          aws_access_key_id: "",
          aws_default_region: "",
        },
        azure: {
          azure_account_key: "",
          azure_account_name: "",
        },
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
      } else if (this.isSambaSelected) {
        return "smb";
      } else {
        return null;
      }
      //// handle all providers
    },
    selectedProviderPrefix() {
      return this[this.selectedProvider].repoPrefix;
    },
    selectedProviderHelper() {
      if (this.isBackblazeSelected) {
        return "backblaze_placeholder";
      } else if (this.isAmazonS3Selected) {
        return "aws_placeholder";
      } else if (this.isGenericS3Selected) {
        return "s3_placeholder";
      } else if (this.isAzureSelected) {
        return "azure_placeholder";
      } else if (this.isSambaSelected) {
        return "samba_placeholder";
      } else {
        return "url_placeholder";
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
        this.focusElement("url");
      }
    },
  },
  methods: {
    clearFields() {
      this.isBackblazeSelected = false;
      this.isAmazonS3Selected = false;
      this.isAzureSelected = false;
      this.isSambaSelected = false;
      this.name = "";
      this.url = "";

      this.backblaze.b2_account_id = "";
      this.backblaze.b2_account_key = "";

      this.aws.aws_access_key_id = "";
      this.aws.aws_default_region = "";
      this.aws.aws_secret_access_key = "";

      this.genericS3.aws_access_key_id = "";
      this.genericS3.aws_secret_access_key = "";

      this.azure.azure_account_name = "";
      this.azure.azure_account_key = "";

      this.smb.smb_host = "";
      this.smb.smb_user = "";
      this.smb.smb_pass = "";
      this.smb.smb_domain = "";

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
      this.isSambaSelected = false;
    },
    selectAmazonS3() {
      //// handle ALL providers
      this.isBackblazeSelected = false;
      this.isAzureSelected = false;
      this.isGenericS3Selected = false;
      this.isSambaSelected = false;
    },
    selectSamba() {
      //// handle ALL providers
      this.isBackblazeSelected = false;
      this.isAzureSelected = false;
      this.isGenericS3Selected = false;
      this.isAmazonS3Selected = false;
    },
    selectGenericS3() {
      //// handle ALL providers
      this.isBackblazeSelected = false;
      this.isAzureSelected = false;
      this.isAmazonS3Selected = false;
      this.isSambaSelected = false;
    },
    selectAzure() {
      //// handle ALL providers
      this.isBackblazeSelected = false;
      this.isAmazonS3Selected = false;
      this.isGenericS3Selected = false;
      this.isSambaSelected = false;
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
        case "genericS3":
          return {
            aws_access_key_id: this.genericS3.aws_access_key_id,
            aws_secret_access_key: this.genericS3.aws_secret_access_key,
          };
        case "azure":
          return {
            azure_account_name: this.azure.azure_account_name,
            azure_account_key: this.azure.azure_account_key,
          };
        case "smb":
          return {
            smb_host: this.smb.smb_host,
            smb_user: this.smb.smb_user,
            smb_pass: this.smb.smb_pass,
            smb_domain: this.smb.smb_domain,
          };
      }
      //// handle all providers
    },
    validateAddBackblazeRepository() {
      // clear errors
      this.error.name = "";
      this.error.url = "";
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
      } else if (this.url.includes(" ")) {
        // wrong url protocol
        this.error.url = this.$t("backup.invalid_url");

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
      this.error.url = "";
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
      } else if (this.url.includes(" ")) {
        // wrong url protocol
        this.error.url = this.$t("backup.invalid_url");

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
      this.error.url = "";
      this.error.repoConnection = "";

      this.error.genericS3.aws_access_key_id = "";
      this.error.genericS3.aws_default_region = "";
      this.error.genericS3.aws_secret_access_key = "";

      let isValidationOk = true;

      if (!this.url) {
        this.error.url = this.$t("common.required");

        if (isValidationOk) {
          this.focusElement("url");
          isValidationOk = false;
        }
      } else if (this.url.includes(" ")) {
        // wrong url protocol
        this.error.url = this.$t("backup.invalid_url");

        if (isValidationOk) {
          this.focusElement("url");
          isValidationOk = false;
        }
      }

      if (!this.genericS3.aws_access_key_id) {
        this.error.genericS3.aws_access_key_id = this.$t("common.required");

        if (isValidationOk) {
          this.focusElement("genericS3_access_key_id");
          isValidationOk = false;
        }
      }

      if (!this.genericS3.aws_secret_access_key) {
        this.error.genericS3.aws_secret_access_key = this.$t("common.required");

        if (isValidationOk) {
          this.focusElement("genericS3_secret_access_key");
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
    validateAddAzureRepository() {
      // clear errors
      this.error.name = "";
      this.error.url = "";
      this.error.repoConnection = "";

      this.error.azure.azure_account_name = "";
      this.error.azure.azure_account_key = "";

      let isValidationOk = true;

      if (!this.url) {
        this.error.url = this.$t("common.required");

        if (isValidationOk) {
          this.focusElement("url");
          isValidationOk = false;
        }
      } else if (this.url.includes(" ")) {
        // wrong url protocol
        this.error.url = this.$t("backup.invalid_url");

        if (isValidationOk) {
          this.focusElement("url");
          isValidationOk = false;
        }
      }

      if (!this.azure.azure_account_name) {
        this.error.azure.azure_account_name = this.$t("common.required");

        if (isValidationOk) {
          this.focusElement("azure_account_name");
          isValidationOk = false;
        }
      }

      if (!this.azure.azure_account_key) {
        this.error.azure.azure_account_key = this.$t("common.required");

        if (isValidationOk) {
          this.focusElement("azure_account_key");
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
    validateAddSambaRepository() {
      // clear errors
      this.error.name = "";
      this.error.url = "";
      this.error.repoConnection = "";

      this.error.smb.smb_host = "";
      this.error.smb.smb_user = "";
      this.error.smb.smb_pass = "";
      this.error.smb.smb_domain = "";

      let isValidationOk = true;

      if (!this.url) {
        this.error.url = this.$t("common.required");

        if (isValidationOk) {
          this.focusElement("url");
          isValidationOk = false;
        }
      } else if (this.url.includes(" ")) {
        // wrong url protocol
        this.error.url = this.$t("backup.invalid_url");

        if (isValidationOk) {
          this.focusElement("url");
          isValidationOk = false;
        }
      }

      if (!this.smb.smb_host) {
        this.error.smb.smb_host = this.$t("common.required");

        if (isValidationOk) {
          this.focusElement("smb_host");
          isValidationOk = false;
        }
      }

      if (!this.smb.smb_user) {
        this.error.smb.smb_user = this.$t("common.required");

        if (isValidationOk) {
          this.focusElement("smb_user");
          isValidationOk = false;
        }
      }

      if (!this.smb.smb_domain) {
        this.error.smb.smb_domain = this.$t("common.required");

        if (isValidationOk) {
          this.focusElement("smb_domain");
          isValidationOk = false;
        }
      }

      if (!this.smb.smb_pass) {
        this.error.smb.smb_pass = this.$t("common.required");

        if (isValidationOk) {
          this.focusElement("smb_pass");
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
        case "smb":
          return this.validateAddSambaRepository();
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
            provider:
              this.selectedProvider == "genericS3"
                ? "generic-s3"
                : this.selectedProvider,
            url: this.selectedProviderPrefix + this.url,
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

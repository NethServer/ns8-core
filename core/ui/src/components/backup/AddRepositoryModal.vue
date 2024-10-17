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
                      {{ $t("backup.generic-s3") }}
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
                      <cv-interactive-tooltip
                        alignment="center"
                        direction="bottom"
                        class="info"
                      >
                        <template slot="content">
                          {{ $t("backup.samba_tooltips") }}
                        </template>
                      </cv-interactive-tooltip>
                    </h6>
                  </div>
                </NsTile>
              </cv-column>
              <!-- Cluster -->
              <cv-column :md="4">
                <NsTile
                  :light="true"
                  kind="selectable"
                  v-model="isClusterSelected"
                  value="providerValue"
                  @click="selectCluster()"
                  class="provider-card"
                >
                  <div class="provider-card-content">
                    <div class="provider-icon">
                      <img
                        :src="require('@/assets/logo.png')"
                        alt="cluster logo"
                      />
                    </div>
                    <h6>
                      {{ $t("backup.cluster") }}
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
          <template v-if="isClusterSelected">
            <NsComboBox
              :options="endpoints"
              v-model.trim="url"
              :autoFilter="true"
              :autoHighlight="true"
              :title="$t('backup.node')"
              :label="$t('backup.cluster_placeholder')"
              :acceptUserInput="false"
              :showItemType="true"
              :invalid-message="$t(error.url)"
              tooltipAlignment="start"
              tooltipDirection="top"
              ref="url"
            >
            </NsComboBox>
          </template>
          <!-- backblaze -->
          <template v-if="isBackblazeSelected">
            <cv-text-input
              :label="$t('backup.backblaze_url_label')"
              v-model.trim="url"
              :invalid-message="error.url"
              :disabled="loading.addBackupRepository"
              ref="url"
            >
            </cv-text-input>
            <cv-text-input
              :label="$t('backup.b2_account_id')"
              v-model.trim="backblaze.b2_account_id"
              :invalid-message="error.backblaze.b2_account_id"
              :disabled="loading.addBackupRepository"
              autocomplete="off"
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
              autocomplete="new-password"
              ref="b2_account_key"
            ></cv-text-input>
          </template>
          <!-- amazon s3 -->
          <template v-if="isAmazonS3Selected">
            <cv-text-input
              :label="$t('backup.aws_url_label')"
              v-model.trim="url"
              :invalid-message="error.url"
              :disabled="loading.addBackupRepository"
              ref="url"
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
              :label="$t('backup.aws_access_key_id')"
              v-model.trim="aws.aws_access_key_id"
              :invalid-message="error.aws.aws_access_key_id"
              :disabled="loading.addBackupRepository"
              autocomplete="off"
              ref="aws_access_key_id"
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
              autocomplete="new-password"
              ref="aws_secret_access_key"
            ></cv-text-input>
          </template>
          <!-- azure -->
          <template v-if="isAzureSelected">
            <cv-text-input
              :label="$t('backup.azure_url_label')"
              v-model.trim="url"
              :invalid-message="error.url"
              :disabled="loading.addBackupRepository"
              ref="url"
            >
            </cv-text-input>
            <cv-text-input
              :label="$t('backup.azure_account_name')"
              v-model.trim="azure.azure_account_name"
              :invalid-message="error.azure.azure_account_name"
              :disabled="loading.addBackupRepository"
              autocomplete="off"
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
              autocomplete="new-password"
              ref="azure_account_key"
            ></cv-text-input>
          </template>
          <!-- generic s3 -->
          <template v-if="isGenericS3Selected">
            <NsTextInput
              :label="$t('backup.genericS3_url_label')"
              v-model.trim="url"
              :invalid-message="error.url"
              :placeholder="$t('backup.genericS3_placeholder')"
              :disabled="loading.addBackupRepository"
              ref="url"
            >
              <template slot="tooltip">
                <span>{{ $t("backup.genericS3_url_tooltip") }}</span>
              </template>
            </NsTextInput>
            <cv-text-input
              :label="$t('backup.genericS3_access_key_id')"
              v-model.trim="genericS3.aws_access_key_id"
              :invalid-message="error.genericS3.aws_access_key_id"
              :disabled="loading.addBackupRepository"
              autocomplete="off"
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
              autocomplete="new-password"
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
            <NsTextInput
              :label="$t('backup.smb_url_label')"
              v-model.trim="url"
              :invalid-message="error.url"
              :disabled="loading.addBackupRepository"
              :placeholder="$t('backup.smb_url_placeholder')"
              ref="url"
            >
              <template slot="tooltip">
                <span>{{ $t("backup.smb_url_tooltip") }}</span>
              </template>
            </NsTextInput>
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
              autocomplete="off"
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
              autocomplete="new-password"
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
                <NsTextInput
                  :label="
                    $t('backup.data_encryption_key') +
                    ' (' +
                    $t('common.optional') +
                    ')'
                  "
                  v-model="password"
                  :helper-text="$t('backup.repository_password_helper')"
                  :invalid-message="error.password"
                  :disabled="loading.addBackupRepository"
                  autocomplete="new-password"
                  ref="password"
                >
                  <template slot="tooltip">
                    {{ $t("backup.data_encryption_key_tooltip") }}
                  </template>
                </NsTextInput>
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
import last from "lodash/last";

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
      isClusterSelected: false,
      endpoints: [],
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
      cluster: {
        repoPrefix: "",
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
      } else if (this.isClusterSelected) {
        return "cluster";
      } else {
        return null;
      }
      //// handle all providers
    },
    selectedProviderPrefix() {
      return this[this.selectedProvider].repoPrefix;
    },
  },
  watch: {
    isShown: function () {
      if (this.isShown) {
        // show first step
        this.step = this.steps[0];
        this.clearFields();
        this.clearErrors();
        this.listBackupEndpoints();
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
        if (this.selectedProvider == "smb") {
          this.focusElement("smb_host");
        } else {
          this.focusElement("url");
        }
      }
    },
  },
  methods: {
    clearFields() {
      this.isBackblazeSelected = false;
      this.isAmazonS3Selected = false;
      this.isAzureSelected = false;
      this.isSambaSelected = false;
      this.isClusterSelected = false;
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
      this.isClusterSelected = false;
    },
    selectAmazonS3() {
      //// handle ALL providers
      this.isBackblazeSelected = false;
      this.isAzureSelected = false;
      this.isGenericS3Selected = false;
      this.isSambaSelected = false;
      this.isClusterSelected = false;
    },
    selectSamba() {
      //// handle ALL providers
      this.isBackblazeSelected = false;
      this.isAzureSelected = false;
      this.isGenericS3Selected = false;
      this.isAmazonS3Selected = false;
      this.isClusterSelected = false;
    },
    selectGenericS3() {
      //// handle ALL providers
      this.isBackblazeSelected = false;
      this.isAzureSelected = false;
      this.isAmazonS3Selected = false;
      this.isSambaSelected = false;
      this.isClusterSelected = false;
    },
    selectAzure() {
      //// handle ALL providers
      this.isBackblazeSelected = false;
      this.isAmazonS3Selected = false;
      this.isGenericS3Selected = false;
      this.isSambaSelected = false;
      this.isClusterSelected = false;
    },
    selectCluster() {
      //// handle ALL providers
      this.isBackblazeSelected = false;
      this.isAmazonS3Selected = false;
      this.isGenericS3Selected = false;
      this.isSambaSelected = false;
      this.isAzureSelected = false;
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
        case "cluster":
          return {};
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
    validateAddClusterRepository() {
      // clear errors
      this.error.name = "";
      this.error.url = "";
      this.error.repoConnection = "";

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
        case "cluster":
          return this.validateAddClusterRepository();
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
    getValidationErrorField(validationError) {
      // error field could be "parameters.fieldName", let's take "fieldName" only
      const fieldTokens = validationError.field.split(".");
      return last(fieldTokens);
    },
    addBackupRepositoryValidationFailed(validationErrors) {
      this.loading.addBackupRepository = false;
      let focusAlreadySet = false;

      for (const validationError of validationErrors) {
        let field = this.getValidationErrorField(validationError);

        if (validationError.error == "backup_repository_not_accessible") {
          // show error notification
          this.error.repoConnection = "error";
        } else if (field !== "(root)") {
          // set i18n error message
          if (
            // we could have error.property and error.provider.property
            field === "name" ||
            field === "url" ||
            field === "repoConnection"
          ) {
            this.error[field] = this.$t("backup." + validationError.error);
          } else {
            this.error[this.selectedProvider][field] = this.$t(
              "backup." + validationError.error
            );
          }
          if (!focusAlreadySet) {
            this.focusElement(field);
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
    async listBackupEndpoints() {
      this.loading.listBackupEndpoints = true;
      this.error.listBackupEndpoints = "";
      const taskAction = "list-cluster-backup-endpoints";

      // register to task completion
      this.$root.$once(
        taskAction + "-completed",
        this.listBackupEndpointsCompleted
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
        this.error.listBackupEndpoints = this.getErrorMessage(err);
        return;
      }
    },
    listBackupEndpointsCompleted(taskContext, taskResult) {
      let endpoints = taskResult.output.endpoints.sort(
        this.sortByProperty("ui_label")
      );

      // Function to convert each JSON object
      function convertToObject(endpoints) {
        return {
          name: endpoints.ui_label,
          label: endpoints.ui_label,
          value: endpoints.url,
        };
      }

      // Iterate through the array and create new objects
      this.endpoints = endpoints.map(convertToObject);
      this.loading.listBackupEndpoints = false;
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

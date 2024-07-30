<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <NsModal
    size="default"
    :visible="isShown"
    @modal-hidden="$emit('hide')"
    @primary-click="alterBackupRepository"
    :primary-button-disabled="loading.alterBackupRepository"
  >
    <template slot="title">{{ $t("backup.edit_backup_repository") }}</template>
    <template slot="content">
      <cv-form @submit.prevent="alterBackupRepository">
        <NsInlineNotification
          v-if="error.repoConnection"
          kind="error"
          :title="$t('backup.backup_repository_auth_error')"
          :description="$t('backup.backup_repository_auth_error_description')"
          :showCloseButton="false"
        />
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
            v-model.trim="backblaze.b2_account_id"
            :invalid-message="$t(error.backblaze.b2_account_id)"
            :disabled="loading.alterBackupRepository"
            ref="b2_account_id"
          >
          </cv-text-input>
          <cv-text-input
            :label="$t('backup.b2_account_key')"
            v-model.trim="backblaze.b2_account_key"
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
            v-model.trim="aws.aws_access_key_id"
            :invalid-message="$t(error.aws.aws_access_key_id)"
            :disabled="loading.alterBackupRepository"
            ref="aws_access_key_id"
          >
          </cv-text-input>
          <cv-text-input
            :label="$t('backup.aws_default_region')"
            v-model.trim="aws.aws_default_region"
            :invalid-message="$t(error.aws.aws_default_region)"
            :disabled="loading.alterBackupRepository"
            ref="aws_default_region"
          >
          </cv-text-input>
          <cv-text-input
            :label="$t('backup.aws_secret_access_key')"
            v-model.trim="aws.aws_secret_access_key"
            :invalid-message="$t(error.aws.aws_secret_access_key)"
            :disabled="loading.alterBackupRepository"
            ref="aws_secret_access_key"
          >
          </cv-text-input>
        </template>
        <!-- samba -->
        <template v-if="repository.provider == 'smb'">
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
        <!-- generic s3 -->
        <!-- <template v-if="repository.provider == 'genericS3'">
          generic s3 ////
        </template> -->
        <!-- azure -->
        <!-- <template v-if="repository.provider == 'azure'"> azure //// </template> -->
        <!-- //// handle ALL providers -->
        <cv-text-input
          :label="$t('backup.repository_name')"
          v-model.trim="name"
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
      </cv-form>
    </template>
    <template slot="secondary-button">{{ $t("common.cancel") }}</template>
    <template slot="primary-button">{{ $t("common.save") }}</template>
  </NsModal>
</template>

<script>
import { UtilService, TaskService, IconService } from "@nethserver/ns8-ui-lib";
import to from "await-to-js";

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
      cluster: {},
      smb: {
        smb_host: "",
        smb_user: "",
        smb_pass: "",
        smb_domain: "",
      },
      //// handle all providers
      loading: {
        alterBackupRepository: false,
      },
      error: {
        name: "",
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
        cluster: {},
        smb: {
          smb_host: "",
          smb_user: "",
          smb_pass: "",
          smb_domain: "",
        },
        //// handle all providers
      },
    };
  },
  watch: {
    isShown: function () {
      if (this.isShown) {
        // ensure repository prop is updated as well
        this.$nextTick(() => {
          this.clearErrors();
          this.name = this.repository.name;

          switch (this.repository.provider) {
            case "backblaze":
              this.backblaze.b2_account_id =
                this.repository.parameters.b2_account_id;
              this.backblaze.b2_account_key =
                this.repository.parameters.b2_account_key;
              break;
            case "aws":
              this.aws.aws_access_key_id =
                this.repository.parameters.aws_access_key_id;
              this.aws.aws_default_region =
                this.repository.parameters.aws_default_region;
              this.aws.aws_secret_access_key =
                this.repository.parameters.aws_secret_access_key;
              break;
            case "smb":
              this.smb.smb_host = this.repository.parameters.smb_host;
              this.smb.smb_domain = this.repository.parameters.smb_domain;
              this.smb.smb_user = this.repository.parameters.smb_user;
              this.smb.smb_pass = this.repository.parameters.smb_pass;
              break;
            //// handle ALL providers
          }
        });
      }
    },
  },
  methods: {
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
            azure_account_name: "",
            azure_account_key: "",
          };
        case "cluster":
          return {};
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
    validateAlterBackblazeRepository() {
      // clear errors
      this.error.name = "";
      this.error.repoConnection = "";

      this.error.backblaze.b2_account_id = "";
      this.error.backblaze.b2_account_key = "";

      let isValidationOk = true;

      if (!this.backblaze.b2_account_id) {
        this.error.backblaze.b2_account_id = "common.required";

        if (isValidationOk) {
          this.focusElement("b2_account_id");
          isValidationOk = false;
        }
      }

      if (!this.backblaze.b2_account_key) {
        this.error.backblaze.b2_account_key = "common.required";

        if (isValidationOk) {
          this.focusElement("b2_account_key");
          isValidationOk = false;
        }
      }

      if (!this.name) {
        this.error.name = "common.required";

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
      this.error.repoConnection = "";

      this.error.aws.aws_access_key_id = "";
      this.error.aws.aws_default_region = "";
      this.error.aws.aws_secret_access_key = "";

      let isValidationOk = true;

      if (!this.aws.aws_access_key_id) {
        this.error.aws.aws_access_key_id = "common.required";

        if (isValidationOk) {
          this.focusElement("aws_access_key_id");
          isValidationOk = false;
        }
      }

      if (!this.aws.aws_default_region) {
        this.error.aws.aws_default_region = "common.required";

        if (isValidationOk) {
          this.focusElement("aws_default_region");
          isValidationOk = false;
        }
      }

      if (!this.aws.aws_secret_access_key) {
        this.error.aws.aws_secret_access_key = "common.required";

        if (isValidationOk) {
          this.focusElement("aws_secret_access_key");
          isValidationOk = false;
        }
      }

      if (!this.name) {
        this.error.name = "common.required";

        if (isValidationOk) {
          this.focusElement("name");
          isValidationOk = false;
        }
      }
      return isValidationOk;
    },
    validateAlterSambaRepository() {
      // clear errors
      this.error.name = "";
      this.error.url = "";
      this.error.repoConnection = "";

      this.error.smb.smb_host = "";
      this.error.smb.smb_user = "";
      this.error.smb.smb_pass = "";
      this.error.smb.smb_domain = "";

      let isValidationOk = true;

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
    validateAlterGenericS3Repository() {
      // clear errors
      this.error.name = "";
      this.error.repoConnection = "";
      console.error("not implemented");
      return false;
    },
    validateAlterAzureRepository() {
      // clear errors
      this.error.name = "";
      this.error.repoConnection = "";
      console.error("not implemented");
      return false;
    },
    validateAlterClusterRepository() {
      // clear errors
      this.error.name = "";
      this.error.url = "";
      this.error.repoConnection = "";

      let isValidationOk = true;

      if (!this.name) {
        this.error.name = this.$t("common.required");

        if (isValidationOk) {
          this.focusElement("name");
          isValidationOk = false;
        }
      }
      return isValidationOk;
    },
    validateAlterBackupRepository() {
      switch (this.repository.provider) {
        case "backblaze":
          return this.validateAlterBackblazeRepository();
        case "azure":
          return this.validateAlterAzureRepository();
        case "aws":
          return this.validateAmazonS3Repository();
        case "genericS3":
          return this.validateAlterGenericS3Repository();
        case "smb":
          return this.validateAlterSambaRepository();
        case "cluster":
          return this.validateAlterClusterRepository();
      }
    },
    async alterBackupRepository() {
      if (!this.validateAlterBackupRepository()) {
        return;
      }
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
            provider: this.repository.provider,
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

        if (validationError.error == "backup_repository_not_accessible") {
          // show error notification
          this.error.repoConnection = "error";
        } else {
          // set i18n error message
          this.error[param] = "backup." + validationError.error;

          if (!focusAlreadySet) {
            this.focusElement(param);
            focusAlreadySet = true;
          }
        }
      }
    },
    alterBackupRepositoryCompleted() {
      this.loading.alterBackupRepository = false;
      this.$emit("repoAltered");
      this.$emit("hide");
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";
</style>

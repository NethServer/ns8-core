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
    :primary-button-disabled="
      loading.alterBackupRepository || loading.fileReading
    "
  >
    <template slot="title">{{ $t("backup.edit_destination") }}</template>
    <template slot="content">
      <cv-form @submit.prevent="alterBackupRepository">
        <NsInlineNotification
          v-if="error.repoConnection"
          kind="error"
          :title="$t('backup.backup_repository_connection_error')"
          :description="$t('backup.' + error.repoConnection)"
          :showCloseButton="false"
        />
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
            :placeholder="$t('common.unchanged_password')"
            :label="$t('backup.b2_account_key')"
            v-model.trim="backblaze.b2_account_key"
            :invalid-message="$t(error.backblaze.b2_account_key)"
            :disabled="loading.alterBackupRepository"
            ref="b2_account_key"
            type="password"
            :password-hide-label="$t('password.hide_password')"
            :password-show-label="$t('password.show_password')"
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
            :placeholder="$t('common.unchanged_password')"
            :label="$t('backup.aws_secret_access_key')"
            v-model.trim="aws.aws_secret_access_key"
            :invalid-message="$t(error.aws.aws_secret_access_key)"
            :disabled="loading.alterBackupRepository"
            type="password"
            :password-hide-label="$t('password.hide_password')"
            :password-show-label="$t('password.show_password')"
            ref="aws_secret_access_key"
          >
          </cv-text-input>
        </template>
        <!-- samba -->
        <template v-if="repository.provider == 'smb'">
          <NsInlineNotification
            kind="info"
            :title="$t('backup.smb_protocol_warning_title')"
            :description="$t('backup.smb_protocol_warning_description')"
            :showCloseButton="false"
          />
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
            :placeholder="$t('common.unchanged_password')"
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
        <template v-if="repository.provider == 'generic-s3'">
          <cv-text-input
            :label="$t('backup.genericS3_access_key_id')"
            v-model.trim="genericS3.aws_access_key_id"
            :invalid-message="$t(error.genericS3.aws_access_key_id)"
            :disabled="loading.alterBackupRepository"
            ref="genericS3_access_key_id"
          >
          </cv-text-input>
          <cv-text-input
            :placeholder="$t('common.unchanged_password')"
            :label="$t('backup.genericS3_secret_access_key')"
            v-model.trim="genericS3.aws_secret_access_key"
            :invalid-message="$t(error.genericS3.aws_secret_access_key)"
            :disabled="loading.alterBackupRepository"
            type="password"
            :password-hide-label="$t('password.hide_password')"
            :password-show-label="$t('password.show_password')"
            ref="genericS3_secret_access_key"
          >
          </cv-text-input>
        </template>
        <!-- rclone -->
        <template v-if="repository.provider == 'rclone'">
          <div>
            <label class="bx--label">{{
              $t("backup.rclone_configuration_source")
            }}</label>
            <cv-radio-group
              :vertical="true"
              :disabled="loading.alterBackupRepository"
            >
              <cv-radio-button
                :label="
                  $t('backup.rclone_configuration_source_keep_configuration')
                "
                value="keep"
                v-model="rclone.configuration_source"
              />
              <cv-radio-button
                :label="
                  $t('backup.rclone_configuration_source_paste_configuration')
                "
                value="paste"
                v-model="rclone.configuration_source"
              />
              <cv-radio-button
                :label="
                  $t('backup.rclone_configuration_source_upload_configuration')
                "
                value="upload"
                v-model="rclone.configuration_source"
              />
            </cv-radio-group>
          </div>
          <template v-if="rclone.configuration_source == 'paste'">
            <cv-text-area
              class="mg-top-md"
              :label="$t('backup.rclone_configuration_content')"
              v-model.trim="rclone.rclone_conf_secret"
              :invalid-message="error.rclone.rclone_conf_secret"
              :helper-text="
                $t(
                  'backup.rclone_configuration_content_single_destination_helper_text'
                )
              "
              :value="rclone.rclone_conf_secret"
              rows="5"
              ref="rclone_conf_secret"
              :disabled="loading.alterBackupRepository"
            ></cv-text-area>
          </template>
          <template v-else-if="rclone.configuration_source == 'upload'">
            <div
              :class="{
                'file-uploader-error': error.rclone.configuration_file,
              }"
            >
              <cv-file-uploader
                :key="componentKey"
                :label="$t('backup.rclone_configuration_file')"
                :helper-text="$t('backup.rclone_configuration_file_helper')"
                :multiple="false"
                :removable="true"
                :clear-on-reselect="true"
                :drop-target-label="
                  $t('common.drag_and_drop_or_click_to_upload')
                "
                v-model="rclone.configuration_uploaded_file"
                accept=".conf,.txt"
                ref="backup_file"
              ></cv-file-uploader>
              <div
                v-if="error.rclone.configuration_file"
                class="validation-failed-invalid-message"
              >
                {{ error.rclone.configuration_file }}
              </div>
            </div>
          </template>
          <cv-text-input
            :label="
              $t('backup.rclone_base_path') + ' (' + $t('common.optional') + ')'
            "
            v-model.trim="rclone.base_path"
            :invalid-message="error.rclone.base_path"
            :disabled="loading.alterBackupRepository"
            ref="rclone_base_path"
          ></cv-text-input>
        </template>
        <cv-text-input
          :label="$t('backup.repository_name')"
          v-model.trim="name"
          :helper-text="$t('backup.repository_name_helper')"
          :invalid-message="error.name"
          :disabled="loading.alterBackupRepository"
          ref="name"
        >
        </cv-text-input>
        <!-- advanced options (read-only) -->
        <cv-accordion>
          <cv-accordion-item :open="false">
            <template slot="title">{{ $t("common.advanced") }}</template>
            <template slot="content">
              <div>
                <label class="bx--label">
                  <span>
                    {{ $t("backup.data_encryption_key") }}
                  </span>
                  <cv-interactive-tooltip
                    alignment="center"
                    direction="right"
                    class="info"
                  >
                    <template slot="content">
                      <div>
                        {{ $t("backup.data_encryption_key_tooltip") }}
                      </div>
                    </template>
                  </cv-interactive-tooltip>
                </label>
                <NsCodeSnippet
                  light
                  copyTipPosition="left"
                  copyTipAlignment="center"
                  :copyTooltip="$t('common.copy_to_clipboard')"
                  :copy-feedback="$t('common.copied_to_clipboard')"
                  :feedback-aria-label="$t('common.copied_to_clipboard')"
                  :wrap-text="true"
                  :moreText="$t('common.show_more')"
                  :lessText="$t('common.show_less')"
                  hideExpandButton
                  class="mg-top-sm"
                  >{{ repository.password }}</NsCodeSnippet
                >
              </div>
            </template>
          </cv-accordion-item>
        </cv-accordion>
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
      genericS3: {
        aws_access_key_id: "",
        aws_secret_access_key: "",
      },
      cluster: {},
      smb: {
        smb_host: "",
        smb_user: "",
        smb_pass: "",
        smb_domain: "",
      },
      rclone: {
        configuration_source: "keep",
        rclone_conf_secret: "",
        base_path: "",
        configuration_uploaded_file: null,
      },
      componentKey: Date.now(),
      //// handle all providers
      loading: {
        alterBackupRepository: false,
        fileReading: false,
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
        genericS3: {
          aws_access_key_id: "",
          aws_secret_access_key: "",
        },
        cluster: {},
        smb: {
          smb_host: "",
          smb_user: "",
          smb_pass: "",
          smb_domain: "",
        },
        rclone: {
          rclone_conf_secret: "",
          base_path: "",
          configuration_file: "",
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
              this.backblaze.b2_account_key = "";
              break;
            case "aws":
              this.aws.aws_access_key_id =
                this.repository.parameters.aws_access_key_id;
              this.aws.aws_default_region =
                this.repository.parameters.aws_default_region;
              this.aws.aws_secret_access_key = "";
              break;
            case "generic-s3":
              this.genericS3.aws_access_key_id =
                this.repository.parameters.aws_access_key_id;
              this.genericS3.aws_secret_access_key = "";
              break;
            case "smb":
              this.smb.smb_host = this.repository.parameters.smb_host;
              this.smb.smb_domain = this.repository.parameters.smb_domain;
              this.smb.smb_user = this.repository.parameters.smb_user;
              this.smb.smb_pass = "";
              break;
            case "rclone":
              this.rclone.configuration_source = "keep";
              this.rclone.rclone_conf_secret = "";
              this.rclone.base_path = this.repository.basepath || "";
              this.rclone.configuration_uploaded_file = null;
              this.componentKey = Date.now();
              break;
          }
        });
      }
    },
    "rclone.configuration_uploaded_file": function (newFiles) {
      if (newFiles && newFiles.length > 0 && newFiles[0].file) {
        const reader = new FileReader();
        this.loading.fileReading = true;
        reader.onload = (e) => {
          this.rclone.rclone_conf_secret = e.target.result.trim();
          this.loading.fileReading = false;
        };
        reader.onerror = () => {
          this.loading.fileReading = false;
          this.error.rclone.configuration_file = this.$t(
            "backup.rclone_conf_error"
          );
        };
        reader.readAsText(newFiles[0].file);
      } else {
        this.rclone.rclone_conf_secret = "";
        this.loading.fileReading = false;
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
        case "generic-s3":
          return {
            aws_access_key_id: this.genericS3.aws_access_key_id,
            aws_secret_access_key: this.genericS3.aws_secret_access_key,
          };
        case "rclone": {
          const secret = this.rclone.rclone_conf_secret;
          return {
            rclone_conf_secret:
              secret && !secret.endsWith("\n") ? secret + "\n" : secret,
            basepath: this.rclone.base_path,
          };
        }
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

      if (!this.name) {
        this.error.name = this.$t("common.required");

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

      this.error.genericS3.aws_access_key_id = "";
      this.error.genericS3.aws_secret_access_key = "";

      let isValidationOk = true;

      if (!this.genericS3.aws_access_key_id) {
        this.error.genericS3.aws_access_key_id = "common.required";

        if (isValidationOk) {
          this.focusElement("genericS3_access_key_id");
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
    validateAlterClusterRepository() {
      // clear errors
      this.error.name = "";
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
    validateAlterRcloneRepository() {
      // clear errors
      this.error.name = "";
      this.error.repoConnection = "";
      this.error.rclone.rclone_conf_secret = "";
      this.error.rclone.base_path = "";
      this.error.rclone.configuration_file = "";

      let isValidationOk = true;

      if (this.rclone.configuration_source === "paste") {
        if (!this.rclone.rclone_conf_secret) {
          this.error.rclone.rclone_conf_secret = this.$t("common.required");

          if (isValidationOk) {
            this.focusElement("rclone_conf_secret");
            isValidationOk = false;
          }
        }
      } else if (this.rclone.configuration_source === "upload") {
        if (
          !this.rclone.configuration_uploaded_file ||
          this.rclone.configuration_uploaded_file.length === 0
        ) {
          this.error.rclone.configuration_file = this.$t("common.required");
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
    validateAlterBackupRepository() {
      switch (this.repository.provider) {
        case "backblaze":
          return this.validateAlterBackblazeRepository();
        case "aws":
          return this.validateAmazonS3Repository();
        case "generic-s3":
          return this.validateAlterGenericS3Repository();
        case "smb":
          return this.validateAlterSambaRepository();
        case "cluster":
          return this.validateAlterClusterRepository();
        case "rclone":
          return this.validateAlterRcloneRepository();
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

        const connectionErrors = [
          "backup_repository_not_accessible",
          "no_read_permission",
          "no_write_permission",
          "no_delete_permission",
        ];
        if (connectionErrors.includes(validationError.error)) {
          // show error notification
          this.error.repoConnection = validationError.error;
        } else {
          // set i18n error message for rclone nested fields
          if (
            this.repository.provider === "rclone" &&
            (param === "rclone_conf_secret" ||
              param === "basepath" ||
              param === "parameters")
          ) {
            const key =
              param === "basepath"
                ? "base_path"
                : param === "parameters"
                ? "rclone_conf_secret"
                : param;
            // use generic rclone_conf_error for any error on rclone_conf_secret
            // (anyOf schema branches produce untranslatable codes like parameters_pattern)
            const translationKey =
              key === "rclone_conf_secret"
                ? "rclone_conf_error"
                : validationError.error;
            // show error on upload widget or paste textarea depending on active mode
            const errorKey =
              key === "rclone_conf_secret" &&
              this.rclone.configuration_source === "upload"
                ? "configuration_file"
                : key;
            this.error.rclone[errorKey] = this.$t("backup." + translationKey);
          } else if (param !== "provider" && param !== "(root)") {
            this.error[param] = this.$t("backup." + validationError.error);
          }

          if (!focusAlreadySet) {
            const refKey =
              this.repository.provider === "rclone" && param === "basepath"
                ? "rclone_base_path"
                : param;
            this.$nextTick(() => {
              const el = this.$refs[refKey];
              if (el) {
                if (typeof el.focus === "function") {
                  el.focus();
                } else if (el.$el) {
                  const focusable = el.$el.querySelector("input, textarea");
                  if (focusable) focusable.focus();
                }
              }
            });
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

// Full width and white background for the rclone file uploader
::v-deep .bx--file {
  width: 100%;
  max-width: 38rem;
}
::v-deep .bx--file__drop-container {
  width: 100%;
  max-width: 38rem;
}
::v-deep .bx--file-container {
  width: 100%;
  max-width: 38rem;
}
::v-deep .bx--file__selected-file {
  background-color: #ffffff;
  max-width: 38rem;
}

.file-uploader-error ::v-deep .bx--file__drop-container {
  outline: 2px solid $danger-01;
  outline-offset: -2px;
}

.file-uploader-error .validation-failed-invalid-message {
  margin-top: -0.75rem;
  margin-bottom: 0.75rem;
  color: $danger-01;
  font-size: 0.75rem;
}
</style>

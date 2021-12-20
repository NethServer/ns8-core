<template>
  <cv-modal
    size="default"
    :visible="isShown"
    @modal-hidden="$emit('hide')"
    class="wizard-modal"
  >
    <template slot="title">{{ $t("backup.add_backup_repository") }}</template>
    <template slot="content">
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
                class="min-height-card"
              >
                <!--  //// provider icon -->
                <h6>
                  {{ $t("backup.backblaze") }}
                </h6>
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
                class="min-height-card"
              >
                <!--  //// provider icon -->
                <h6>
                  {{ $t("backup.amazon_s3") }}
                </h6>
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
                class="min-height-card"
              >
                <!--  //// provider icon -->
                <h6>
                  {{ $t("backup.azure") }}
                </h6>
              </NsTile>
            </div>
          </div>
        </div>
      </template>
      <template v-if="step == 'settings'">
        <!-- backblaze -->
        <cv-form>
          <template v-if="isBackblazeSelected">
            <!-- <cv-text-input ////
            :label="$t('backup.domain')"
            v-model.trim="external.domain"
            :helper-text="$t('domains.enter_external_domain_name')"
            :invalid-message="$t(error.external.domain)"
            :disabled="loading.external.addExternalDomain"
            ref="domain"
          >
          </cv-text-input> -->
            backblaze ////
          </template>
          <!-- amazon s3 -->
          <template v-if="isAmazonS3Selected"> amazon s3 //// </template>
          <!-- azure -->
          <template v-if="isAzureSelected"> azure //// </template>
          <!-- //// handle ALL providers -->
          <cv-text-input
            :label="$t('backup.url')"
            v-model.trim="url"
            :invalid-message="$t(error.url)"
            :disabled="loading.addBackupRepository"
            ref="url"
          >
          </cv-text-input>
          <cv-text-input
            :label="$t('backup.repository_name')"
            v-model.trim="name"
            :helper-text="$t('backup.repository_name_helper')"
            :invalid-message="$t(error.name)"
            :disabled="loading.addBackupRepository"
            ref="name"
          >
          </cv-text-input>
        </cv-form>
        <NsInlineNotification
          v-if="error.addBackupRepository"
          kind="error"
          :title="$t('action.add-backup-repository')"
          :description="error.addBackupRepository"
          :showCloseButton="false"
        />
      </template>
      <div class="wizard-buttons">
        <NsButton
          kind="secondary"
          :icon="Close20"
          @click="$emit('hide')"
          class="wizard-button"
          >{{ $t("common.cancel") }}
        </NsButton>
        <NsButton
          kind="secondary"
          :icon="ChevronLeft20"
          @click="previousStep"
          :disabled="loading.addBackupRepository || step == 'provider'"
          class="wizard-button"
          >{{ $t("common.previous") }}
        </NsButton>
        <NsButton
          kind="primary"
          :icon="ChevronRight20"
          @click="nextStep"
          :disabled="loading.addBackupRepository || !selectedProvider"
          :loading="loading.addBackupRepository"
          class="wizard-button"
          ref="wizardNext"
          >{{ nextStepLabel }}
        </NsButton>
      </div>
    </template>
  </cv-modal>
</template>

<script>
import { UtilService, TaskService, IconService } from "@nethserver/ns8-ui-lib";
import to from "await-to-js";

export default {
  name: "AddBackupRepositoryModal",
  mixins: [UtilService, TaskService, IconService],
  props: {
    isShown: {
      type: Boolean,
      default: true,
    },
  },
  data() {
    return {
      step: "provider",
      isBackblazeSelected: false,
      isAmazonS3Selected: false,
      isAzureSelected: false,
      name: "",
      url: "",
      loading: {
        addBackupRepository: false,
      },
      error: {
        name: "",
        url: "",
        addBackupRepository: "",
      },
    };
  },
  computed: {
    nextStepLabel() {
      if (this.step == "settings") {
        return this.$t("backup.add_backup_repository");
      } else {
        return this.$t("common.next");
      }
    },
    selectedProvider() {
      if (this.isBackblazeSelected) {
        return "backblaze";
      } else if (this.isAmazonS3Selected) {
        return "amazon";
      } else if (this.isAzureSelected) {
        return "azure";
      } else {
        return null;
      }
      //// handle all providers
    },
  },
  watch: {
    isShown: function () {
      if (this.isShown) {
        this.step = "provider";
      }
    },
    step: function () {
      // if (this.step == "internalConfig") { ////
      //   if (this.isOpenLdapSelected) {
      //     //// focus first input field
      //   } else if (this.isSambaSelected) {
      //     setTimeout(() => {
      //       this.focusElement("realm");
      //     }, 300);
      //   }
      // }
    },
  },
  created() {
    // this.newProviderId = this.providerId; ////
    // this.isOpenLdapSelected = this.isOpenLdap;
    // this.isSambaSelected = this.isSamba;
  },
  methods: {
    nextStep() {
      switch (this.step) {
        case "provider":
          this.step = "settings";
          break;
        case "settings":
          this.addBackupRepository();
          break;
      }
    },
    previousStep() {
      switch (this.step) {
        case "settings":
          this.step = "provider";
          break;
      }
    },
    selectBackblaze() {
      //// handle ALL providers
      this.isAmazonS3Selected = false;
      this.isAzureSelected = false;
    },
    selectAmazonS3() {
      //// handle ALL providers
      this.isBackblazeSelected = false;
      this.isAzureSelected = false;
    },
    selectAzure() {
      //// handle ALL providers
      this.isBackblazeSelected = false;
      this.isAmazonS3Selected = false;
    },
    buildRepositoryParameters() {
      switch (this.selectedProvider) {
        ////
        case "backblaze":
          return {
            b2_account_id: "",
            b2_account_key: "",
          };
        case "amazon":
          return {
            aws_default_region: "",
            aws_access_key_id: "",
            aws_secret_access_key: "",
          };
        case "azure":
          return {
            azure_account_name: "",
            azure_account_key: "",
          };
      }
      //// handle all providers
    },
    async addBackupRepository() {
      //// validation
      this.error.addBackupRepository = "";
      const taskAction = "add-backup-repository";

      // register to task completion
      this.$root.$once(
        taskAction + "-completed",
        this.addBackupRepositoryCompleted
      );

      const parameters = this.buildRepositoryParameters();

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            name: this.name,
            provider: this.selectedProvider,
            url: this.url,
            parameters,
            //// password
          },
          extra: {
            title: this.$t("action." + taskAction),
            description: this.$t("common.processing"),
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
    addBackupRepositoryCompleted(taskContext, taskResult) {
      console.log("addBackupRepositoryCompleted", taskResult.output); ////

      //// show repo password
    },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";

.min-height-card {
  min-height: 8rem;
}
</style>

<template>
  <div class="bx--grid bx--grid--full-width">
    <div class="bx--row">
      <div class="bx--col-lg-16 page-title">
        <h2>{{ $t("settings.title") }}</h2>
      </div>
    </div>
    <!-- sample settings page -->
    <!-- <div v-if="error.getConfiguration" class="bx--row">
      <div class="bx--col">
        <NsInlineNotification
          kind="error"
          :title="$t('action.get-configuration')"
          :description="error.getConfiguration"
          :showCloseButton="false"
        />
      </div>
    </div>
    <div class="bx--row">
      <div class="bx--col-lg-16">
        <cv-tile :light="true">
          <cv-form @submit.prevent="saveSettings"> -->
    <!-- TODO configuration fields -->
    <!-- <div v-if="error.configureModule" class="bx--row">
              <div class="bx--col">
                <NsInlineNotification
                  kind="error"
                  :title="$t('action.configure-module')"
                  :description="error.configureModule"
                  :showCloseButton="false"
                />
              </div>
            </div>
            <NsButton
              kind="primary"
              :icon="Save20"
              :loading="loading.settings"
              :disabled="loading.settings"
              >{{ $t("settings.save") }}</NsButton
            >
          </cv-form>
        </cv-tile>
      </div>
    </div> -->
  </div>
</template>

<script>
// import to from "await-to-js";
import { mapState } from "vuex";
import {
  QueryParamService,
  UtilService,
  TaskService,
  IconService,
} from "@nethserver/ns8-ui-lib";

export default {
  name: "Settings",
  mixins: [TaskService, IconService, UtilService, QueryParamService],
  pageTitle() {
    return this.$t("settings.title") + " - " + this.appName;
  },
  data() {
    return {
      q: {
        page: "settings",
      },
      urlCheckInterval: null,
      loading: {
        settings: true,
      },
      error: {
        getConfiguration: "",
        configureModule: "",
      },
    };
  },
  computed: {
    ...mapState(["instanceName", "ns8Core", "appName"]),
  },
  created() {
    // this.getConfiguration();
  },
  beforeRouteEnter(to, from, next) {
    next((vm) => {
      vm.watchQueryData(vm);
      vm.urlCheckInterval = vm.initUrlBindingForApp(vm, vm.q.page);
    });
  },
  beforeRouteLeave(to, from, next) {
    clearInterval(this.urlCheckInterval);
    next();
  },
  methods: {
    // async getConfiguration() {
    //   this.loading.settings = true;
    //   this.error.getConfiguration = "";
    //   const taskAction = "get-configuration";
    //   // register to task completion
    //   this.ns8Core.$root.$once(
    //     taskAction + "-completed",
    //     this.getConfigurationCompleted
    //   );
    //   const res = await to(
    //     this.createModuleTaskForApp(this.instanceName, {
    //       action: taskAction,
    //       extra: {
    //         title: this.$t("action." + taskAction),
    //         isNotificationHidden: true,
    //       },
    //     })
    //   );
    //   const err = res[0];
    //   if (err) {
    //     console.error(`error creating task ${taskAction}`, err);
    //     this.error.getConfiguration = this.getErrorMessage(err);
    //     return;
    //   }
    // },
    // validateSaveSettings() {
    //   this.clearErrors(this);
    //   let isValidationOk = true;
    //   // TODO validate settings
    //   return isValidationOk;
    // },
    // saveSettingsValidationFailed(validationErrors) {
    //   this.loading.settings = false;
    //   for (const validationError of validationErrors) {
    //     const param = validationError.parameter;
    //     // set i18n error message
    //     this.error[param] = "settings." + validationError.error;
    //   }
    // },
    // async saveSettings() {
    //   const isValidationOk = this.validateSaveSettings();
    //   if (!isValidationOk) {
    //     return;
    //   }
    //   this.loading.settings = true;
    //   const taskAction = "configure-module";
    //   // register to task validation
    //   this.ns8Core.$root.$off(taskAction + "-validation-failed");
    //   this.ns8Core.$root.$once(
    //     taskAction + "-validation-failed",
    //     this.saveSettingsValidationFailed
    //   );
    //   // register to task completion
    //   this.ns8Core.$root.$off(taskAction + "-completed");
    //   this.ns8Core.$root.$once(
    //     taskAction + "-completed",
    //     this.saveSettingsCompleted
    //   );
    //   const res = await to(
    //     this.createModuleTaskForApp(this.instanceName, {
    //       action: taskAction,
    //       data: {
    //         // TODO
    //       },
    //       extra: {
    //         title: this.$t("settings.instance_configuration", {
    //           instance: this.instanceName,
    //         }),
    //         description: this.$t("settings.configuring"),
    //       },
    //     })
    //   );
    //   const err = res[0];
    //   if (err) {
    //     console.error(`error creating task ${taskAction}`, err);
    //     this.error.configureModule = this.getErrorMessage(err);
    //     this.loading.settings = false;
    //     return;
    //   }
    // },
    // getConfigurationCompleted(taskContext, taskResult) {
    //   const config = taskResult.output;
    //   // TODO
    //   this.focusElement("wikiName");
    // },
    // saveSettingsCompleted() {
    //   this.loading.settings = false;
    //   // reload configuration
    //   this.getConfiguration();
    // },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";
</style>

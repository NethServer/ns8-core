<template>
  <div class="bx--grid bx--grid--full-width">
    <div class="bx--row">
      <div class="bx--col-lg-16 page-title">
        <h2>{{ $t("settings.title") }}</h2>
      </div>
    </div>
    <div class="bx--row">
      <div class="bx--col-lg-16">
        <cv-tile :light="true" class="content-tile">
          <cv-form @submit.prevent="saveSettings">
            <cv-text-input
              :label="$t('settings.wiki_name')"
              v-model.trim="wikiName"
              class="mg-bottom"
              :invalid-message="$t(error.wiki_name)"
              :disabled="loading.settings"
              ref="wikiName"
            >
            </cv-text-input>
            <cv-text-input
              :label="$t('settings.admin_username')"
              v-model.trim="username"
              class="mg-bottom"
              :invalid-message="$t(error.username)"
              :disabled="loading.settings"
              ref="username"
            >
            </cv-text-input>
            <cv-text-input
              :label="$t('settings.admin_password')"
              v-model.trim="password"
              type="password"
              :password-show-label="$t('settings.show_password')"
              :password-hide-label="$t('settings.hide_password')"
              class="mg-bottom"
              :invalid-message="$t(error.password)"
              :disabled="loading.settings"
              ref="password"
            >
            </cv-text-input>
            <cv-text-input
              :label="$t('settings.admin_email')"
              placeholder="admin@example.com"
              v-model.trim="email"
              class="mg-bottom"
              :invalid-message="$t(error.email)"
              :disabled="loading.settings"
              ref="email"
            >
            </cv-text-input>
            <cv-text-input
              :label="$t('settings.admin_full_name')"
              v-model.trim="userFullName"
              class="mg-bottom"
              :invalid-message="$t(error.user_full_name)"
              :disabled="loading.settings"
              ref="userFullName"
            >
            </cv-text-input>
            <cv-text-input
              :label="$t('settings.wiki_fqdn')"
              placeholder="mywiki.example.org"
              v-model.trim="host"
              class="mg-bottom"
              :invalid-message="$t(error.host)"
              :disabled="loading.settings"
              ref="host"
            >
            </cv-text-input>
            <cv-toggle
              value="letsEncrypt"
              :label="$t('settings.lets_encrypt')"
              v-model="isLetsEncryptEnabled"
              :disabled="loading.settings"
              class="mg-bottom"
            >
              <template slot="text-left">{{
                $t("settings.disabled")
              }}</template>
              <template slot="text-right">{{
                $t("settings.enabled")
              }}</template>
            </cv-toggle>
            <cv-toggle
              value="httpToHttps"
              :label="$t('settings.http_to_https')"
              v-model="isHttpToHttpsEnabled"
              :disabled="loading.settings"
              class="mg-bottom"
            >
              <template slot="text-left">{{
                $t("settings.disabled")
              }}</template>
              <template slot="text-right">{{
                $t("settings.enabled")
              }}</template>
            </cv-toggle>
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
    </div>
  </div>
</template>

<script>
import to from "await-to-js";
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
      wikiName: "",
      username: "",
      password: "",
      email: "",
      userFullName: "",
      host: "",
      isLetsEncryptEnabled: false,
      isHttpToHttpsEnabled: false,
      loading: {
        settings: true,
      },
      error: {
        wiki_name: "",
        username: "",
        password: "",
        user_full_name: "",
        email: "",
        host: "",
        lets_encrypt: "",
        http2https: "",
      },
    };
  },
  computed: {
    ...mapState(["instanceName", "ns8Core", "appName"]),
  },
  created() {
    this.getConfiguration();
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
    async getConfiguration() {
      this.loading.settings = true;
      const taskAction = "get-configuration";

      // register to task completion
      this.ns8Core.$root.$on(
        taskAction + "-completed",
        this.getConfigurationCompleted
      );

      const res = await to(
        this.createModuleTaskForApp(this.instanceName, {
          action: taskAction,
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
          },
        })
      );
      const err = res[0];

      if (err) {
        this.createErrorNotificationForApp(
          err,
          this.$t("task.cannot_create_task", { action: taskAction })
        );
        return;
      }
    },
    validateSaveSettings() {
      this.clearErrors(this);

      let isValidationOk = true;

      if (!this.wikiName) {
        this.error.wiki_name = "common.required";

        if (isValidationOk) {
          this.focusElement("wikiName");
        }
        isValidationOk = false;
      }

      if (!this.username) {
        this.error.username = "common.required";

        if (isValidationOk) {
          this.focusElement("username");
        }
        isValidationOk = false;
      }

      if (!this.password) {
        this.error.password = "common.required";

        if (isValidationOk) {
          this.focusElement("password");
        }
        isValidationOk = false;
      }

      if (!this.email) {
        this.error.email = "common.required";

        if (isValidationOk) {
          this.focusElement("email");
        }
        isValidationOk = false;
      }

      if (this.email && !/^\S+@\S+$/.test(this.email)) {
        this.error.email = "settings.email_format";

        if (isValidationOk) {
          this.focusElement("email");
        }
        isValidationOk = false;
      }

      if (!this.userFullName) {
        this.error.user_full_name = "common.required";

        if (isValidationOk) {
          this.focusElement("userFullName");
        }
        isValidationOk = false;
      }

      if (!this.host) {
        this.error.host = "common.required";

        if (isValidationOk) {
          this.focusElement("host");
        }
        isValidationOk = false;
      }

      return isValidationOk;
    },
    saveSettingsValidationFailed(validationErrors) {
      this.loading.settings = false;

      for (const validationError of validationErrors) {
        const param = validationError.parameter;
        // set i18n error message
        this.error[param] = "settings." + validationError.error;
      }
    },
    async saveSettings() {
      const isValidationOk = this.validateSaveSettings();
      if (!isValidationOk) {
        return;
      }

      this.loading.settings = true;
      const taskAction = "configure-module";

      // register to task validation
      this.ns8Core.$root.$on(
        taskAction + "-validation-failed",
        this.saveSettingsValidationFailed
      );

      // register to task completion
      this.ns8Core.$root.$on(
        taskAction + "-completed",
        this.saveSettingsCompleted
      );

      const res = await to(
        this.createModuleTaskForApp(this.instanceName, {
          action: taskAction,
          data: {
            wiki_name: this.wikiName,
            username: this.username,
            password: this.password,
            user_full_name: this.userFullName,
            email: this.email,
            host: this.host,
            lets_encrypt: this.isLetsEncryptEnabled,
            http2https: this.isHttpToHttpsEnabled,
          },
          extra: {
            title: this.$t("settings.instance_configuration", {
              instance: this.instanceName,
            }),
            description: this.$t("settings.configuring"),
          },
        })
      );
      const err = res[0];

      if (err) {
        this.createErrorNotificationForApp(
          err,
          this.$t("task.cannot_create_task", { action: taskAction })
        );
        return;
      }
    },
    getConfigurationCompleted(taskContext, taskResult) {
      // unregister from event
      this.ns8Core.$root.$off("get-configuration-completed");

      const config = taskResult.output;
      this.wikiName = config.wiki_name;
      this.username = config.username;
      this.password = config.password;
      this.userFullName = config.user_full_name;
      this.email = config.email;
      this.host = config.host;
      this.isLetsEncryptEnabled = config.lets_encrypt;
      this.isHttpToHttpsEnabled = config.http2https;
      this.loading.settings = false;
      this.focusElement("wikiName");
    },
    saveSettingsCompleted() {
      // unregister from event
      this.ns8Core.$root.$off("configure-module-completed");

      this.loading.settings = false;

      // reload configuration
      this.getConfiguration();
    },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";
.mg-bottom {
  margin-bottom: $spacing-06;
}
</style>

<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <cv-grid fullWidth>
    <cv-row>
      <cv-column>
        <cv-breadcrumb
          aria-label="breadcrumb"
          :no-trailing-slash="true"
          class="breadcrumb"
        >
          <cv-breadcrumb-item>
            <cv-link to="/settings">{{ $t("settings.title") }}</cv-link>
          </cv-breadcrumb-item>
          <cv-breadcrumb-item>
            <span>{{ $t("smarthost.title") }}</span>
          </cv-breadcrumb-item>
        </cv-breadcrumb>
      </cv-column>
    </cv-row>
    <cv-row>
      <cv-column class="subpage-title">
        <h3>{{ $t("smarthost.title") }}</h3>
      </cv-column>
    </cv-row>
    <cv-row v-if="error.getSmarthost">
      <cv-column>
        <NsInlineNotification
          kind="error"
          :title="$t('action.get-smarthost')"
          :description="error.getSmarthost"
          :showCloseButton="false"
        />
      </cv-column>
    </cv-row>
    <cv-row>
      <cv-column>
        <cv-tile :light="true">
          <cv-form @submit.prevent="saveSettings">
            <cv-toggle
              :label="$t('smarthost.enabled')"
              value="statusValue"
              :form-item="true"
              v-model="enabled"
              :disabled="loading.getSmarthost || loading.setSmarthost"
            >
              <template slot="text-left">
                {{ $t("common.disabled") }}
              </template>
              <template slot="text-right">
                {{ $t("common.enabled") }}
              </template>
            </cv-toggle>
            <template v-if="enabled === true">
              <NsTextInput
                v-model.trim="host"
                :placeholder="$t('smarthost.hostname_placeholder')"
                :label="$t('smarthost.hostname_label')"
                :invalid-message="error.host"
                :disabled="loading.getSmarthost || loading.setSmarthost"
                ref="host"
              >
              </NsTextInput>
              <NsTextInput
                v-model.trim="port"
                :label="$t('smarthost.port_label')"
                :placeholder="$t('smarthost.port_placeholder')"
                :invalid-message="error.port"
                :disabled="loading.getSmarthost || loading.setSmarthost"
                ref="port"
              >
              </NsTextInput>
              <NsComboBox
                v-model="encrypt_smtp"
                :autoFilter="true"
                :autoHighlight="true"
                :title="$t('smarthost.encrypt_smtp')"
                :label="$t('smarthost.choose')"
                :options="options"
                :acceptUserInput="false"
                :showItemType="false"
                :invalid-message="error.encrypt_smtp"
                :disabled="loading.getSmarthost || loading.setSmarthost"
                ref="encrypt_smtp"
              />
              <template v-if="encrypt_smtp != 'none'">
                <cv-toggle
                  :label="$t('smarthost.tls_verify')"
                  value="statusValue"
                  :form-item="true"
                  v-model="tls_verify"
                  :disabled="loading.getSmarthost || loading.setSmarthost"
                >
                  <template slot="text-left">
                    {{ $t("common.disabled") }}
                  </template>
                  <template slot="text-right">
                    {{ $t("common.enabled") }}
                  </template>
                </cv-toggle>
              </template>
              <NsTextInput
                v-model.trim="username"
                :placeholder="$t('smarthost.username_placeholder')"
                :label="$t('smarthost.username_label')"
                :helper-text="$t('smarthost.username_helper')"
                :invalid-message="error.username"
                :disabled="loading.getSmarthost || loading.setSmarthost"
                ref="username"
              >
              </NsTextInput>
              <NsTextInput
                :label="$t('smarthost.password_label')"
                v-model="password"
                :invalid-message="$t(error.password)"
                type="password"
                ref="password"
              >
              </NsTextInput>
              <div ref="setSmarthostError">
                <cv-row v-if="error.test_smarthost">
                  <cv-column>
                    <NsInlineNotification
                      kind="error"
                      :title="$t('action.set-smarthost')"
                      :description="error.test_smarthost"
                      :showCloseButton="false"
                    />
                  </cv-column>
                </cv-row>
              </div>
            </template>
            <NsButton
              kind="primary"
              :icon="Save20"
              :loading="loading.setSmarthost"
              :disabled="isLoadingSettings"
              >{{ $t("common.save_settings") }}</NsButton
            >
          </cv-form>
        </cv-tile>
      </cv-column>
    </cv-row>
  </cv-grid>
</template>

<script>
import to from "await-to-js";
import {
  QueryParamService,
  UtilService,
  TaskService,
  IconService,
  PageTitleService,
} from "@nethserver/ns8-ui-lib";
import { mapActions } from "vuex";

export default {
  name: "SettingsSmartHost",
  mixins: [
    TaskService,
    UtilService,
    IconService,
    QueryParamService,
    PageTitleService,
  ],
  pageTitle() {
    return this.$t("smarthost.title");
  },
  data() {
    return {
      q: {},
      port: "587",
      host: "",
      username: "",
      password: "",
      enabled: false,
      encrypt_smtp: "STARTTLS",
      tls_verify: true,
      options: [
        {
          name: "none",
          label: this.$t("smarthost.none"),
          value: "none",
        },
        {
          name: "starttls",
          label: this.$t("smarthost.starttls"),
          value: "starttls",
        },
        {
          name: "tls",
          label: this.$t("smarthost.tls"),
          value: "tls",
        },
      ],
      loading: {
        getSmarthost: true,
        setSmarthost: false,
      },
      error: {
        getSmarthost: "",
        port: "",
        host: "",
        username: "",
        password: "",
        enabled: "",
        tls_verify: "",
        setSmarthost: "",
        test_smarthost: "",
        encrypt_smtp: "",
      },
    };
  },
  computed: {
    isLoadingSettings() {
      return this.loading.getSmarthost || this.loading.setSmarthost;
    },
  },
  beforeRouteEnter(to, from, next) {
    next((vm) => {
      vm.watchQueryData(vm);
      vm.queryParamsToDataForCore(vm, to.query);
    });
  },
  beforeRouteUpdate(to, from, next) {
    this.queryParamsToDataForCore(this, to.query);
    next();
  },
  watch: {
    "error.setSmarthost": function () {
      if (this.error.setSmarthost) {
        // scroll to notification error

        this.$nextTick(() => {
          const el = this.$refs.setSmarthostError;
          this.scrollToElement(el);
        });
      }
    },
  },
  created() {
    this.getSmarthost();
  },
  methods: {
    ...mapActions(["setSmarthostInStore"]),
    async getSmarthost() {
      for (const key of Object.keys(this.error)) {
        this.error[key] = "";
      }
      // this.error.getSmarthost = "";
      this.loading.getSmarthost = true;
      const taskAction = "get-smarthost";

      // register to task completion
      this.$root.$once(taskAction + "-completed", this.getSmarthostCompleted);

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
        this.error.getSmarthost = this.getErrorMessage(err);
        this.loading.getSmarthost = false;
        return;
      }
    },
    getSmarthostCompleted(taskContext, taskResult) {
      const smarthost = taskResult.output;
      this.host = smarthost.host;
      this.username = smarthost.username;
      this.password = smarthost.password;
      this.port = smarthost.port.toString();
      this.encrypt_smtp = smarthost.encrypt_smtp;
      this.tls_verify = smarthost.tls_verify;
      this.enabled = smarthost.enabled;
      this.loading.getSmarthost = false;
    },
    saveSettings() {
      this.setSmarthost();
    },
    async setSmarthost() {
      this.error.setSmarthost = false;
      this.error.test_smarthost = false;
      const isValidationOk = this.validateConfigureModule();
      if (!isValidationOk) {
        return;
      }
      this.loading.setSmarthost = true;

      const taskAction = "set-smarthost";

      // register to task completion
      this.$root.$once(taskAction + "-completed", this.setSmarthostCompleted);
      // register to task error
      this.$root.$once(taskAction + "-aborted", this.setSmarthostAborted);
      this.$root.$once(
        taskAction + "-validation-failed",
        this.setSmarthostValidationFailed
      );
      // register to task completion
      this.$root.$once(taskAction + "-completed", this.setSmarthostCompleted);

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            host: this.enabled ? this.host : "",
            username: this.enabled ? this.username : "",
            password: this.enabled ? this.password : "",
            port: this.enabled ? parseInt(this.port) : 587,
            encrypt_smtp: this.enabled ? this.encrypt_smtp : "starttls",
            tls_verify: this.enabled ? this.tls_verify : true,
            enabled: this.enabled,
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
        this.error.setSmarthost = this.getErrorMessage(err);
        this.loading.setSmarthost = false;
        return;
      }
    },
    setSmarthostValidationFailed(validationErrors) {
      this.loading.setSmarthost = false;
      let focusAlreadySet = false;

      for (const validationError of validationErrors) {
        const param = validationError.parameter;

        // set i18n error message
        this.error[param] = this.getI18nStringWithFallback(
          "smarthost." + validationError.error,
          "error." + validationError.error
        );

        if (!focusAlreadySet) {
          this.focusElement(param);
          focusAlreadySet = true;
        }
      }
    },
    setSmarthostCompleted() {
      this.loading.setSmarthost = false;
      this.getSmarthost();
    },
    setSmarthostAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.setSmarthost = this.$t("error.generic_error");
      this.loading.setSmarthost = false;
    },
    validateConfigureModule() {
      for (const key of Object.keys(this.error)) {
        this.error[key] = "";
      }

      let isValidationOk = true;
      if (!this.host) {
        this.error.host = this.getI18nStringWithFallback("smarthost.required");
        if (isValidationOk) {
          this.focusElement("host");
        }
        isValidationOk = false;
      }
      return isValidationOk;
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";
</style>

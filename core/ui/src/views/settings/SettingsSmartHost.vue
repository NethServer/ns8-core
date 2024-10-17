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
            <template v-if="enabled">
              <label class="bx--label">
                {{ $t("smarthost.configuration") }}
              </label>
              <cv-radio-group vertical>
                <cv-radio-button
                  v-model="manual_configuration"
                  value="manual"
                  :label="$t('smarthost.manual_configuration')"
                  ref="manual_configuration"
                  checked
                ></cv-radio-button>
                <cv-radio-button
                  v-model="manual_configuration"
                  value="app"
                  :label="$t('smarthost.mail_app_instance')"
                  ref="manual_configuration"
                ></cv-radio-button>
              </cv-radio-group>
            </template>
            <template v-if="enabled && manual_configuration === 'app'">
              <NsInlineNotification
                v-if="!mail_server.length > 0"
                kind="info"
                :title="$t('smarthost.no_mail_app_instance')"
                :description="$t('smarthost.no_mail_app_instance_description')"
                :showCloseButton="false"
                :actionLabel="$t('smarthost.go_to_software_center')"
                @action="goToSoftwareCenter"
              />
              <NsComboBox
                key="app"
                :options="combobox_mail_server"
                v-model="selected_mail_id"
                :autoFilter="true"
                :autoHighlight="true"
                :title="$t('smarthost.mail_app_instance')"
                :label="$t('smarthost.choose_instance')"
                :acceptUserInput="false"
                :showItemType="false"
                :invalid-message="$t(error.selected_mail_id)"
                :disabled="
                  loading.getSmarthost ||
                  loading.setSmarthost ||
                  !mail_server.length > 0
                "
                ref="selected_mail_id"
              />
            </template>
            <template v-else-if="enabled && manual_configuration === 'manual'">
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
                key="manual"
                v-model="encrypt_smtp"
                :autoFilter="true"
                :autoHighlight="true"
                :title="$t('smarthost.encrypt_smtp')"
                :label="$t('smarthost.choose')"
                :options="combobox_encrypt_smtp"
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
                autocomplete="off"
                ref="username"
              >
              </NsTextInput>
              <NsTextInput
                :label="$t('smarthost.password_label')"
                v-model="password"
                :invalid-message="$t(error.password)"
                type="password"
                autocomplete="new-password"
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
              :disabled="
                isLoadingSettings ||
                (!mail_server.length > 0 && manual_configuration === 'app')
              "
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
      manual_configuration: "manual",
      mail_server: [],
      combobox_mail_server: [],
      selected_mail_id: "",
      combobox_encrypt_smtp: [
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
        selected_mail_id: "",
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
    goToSoftwareCenter() {
      this.$router.push("/software-center");
    },
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
    convertToComboboxObject(server) {
      const label = `${server.mail_name ? server.mail_name : server.mail_id} (${
        server.node_name
          ? server.node_name
          : this.$t("smarthost.node", { nodeId: server.node })
      })`;
      return {
        name: label,
        label: label,
        value: server.mail_id,
      };
    },
    getSmarthostCompleted(taskContext, taskResult) {
      const smarthost = taskResult.output;

      this.manual_configuration = smarthost.manual_configuration
        ? "manual"
        : "app";
      this.host = smarthost.host;
      this.username = smarthost.username;
      this.password = smarthost.password;
      this.port = smarthost.port.toString();
      this.tls_verify = smarthost.tls_verify;
      this.enabled = smarthost.enabled;

      this.$nextTick(() => {
        this.encrypt_smtp = smarthost.encrypt_smtp;
      });

      this.mail_server = smarthost.mail_server;

      this.combobox_mail_server = smarthost.mail_server.map(
        this.convertToComboboxObject
      );
      this.$nextTick(() => {
        if (smarthost.mail_server) {
          const selected_mail_server = smarthost.mail_server.filter(
            (mail) => mail.host === smarthost.host
          )[0];
          this.selected_mail_id = selected_mail_server
            ? selected_mail_server.mail_id
            : "";
        }
      });

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

      let data = {};
      if (!this.enabled) {
        data = {
          host: "",
          username: "",
          password: "",
          port: 587,
          encrypt_smtp: "starttls",
          enabled: false,
          tls_verify: true,
        };
      } else if (this.manual_configuration === "manual") {
        data = {
          host: this.host,
          username: this.username,
          password: this.password,
          port: parseInt(this.port),
          encrypt_smtp: this.encrypt_smtp,
          enabled: true,
          tls_verify: this.encrypt_smtp === "none" ? false : this.tls_verify,
        };
      } else {
        const server_host = this.mail_server.filter(
          (mail) => mail.mail_id === this.selected_mail_id
        )[0].host;
        data = {
          host: server_host,
          username: "",
          password: "",
          port: 25,
          encrypt_smtp: "none",
          enabled: true,
          tls_verify: false,
        };
      }

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: data,
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
      if (this.manual_configuration === "manual" && !this.host) {
        this.error.host = this.getI18nStringWithFallback("smarthost.required");
        if (isValidationOk) {
          this.focusElement("host");
        }
        isValidationOk = false;
      }

      if (this.manual_configuration === "app" && !this.selected_mail_id) {
        this.error.selected_mail_id =
          this.getI18nStringWithFallback("smarthost.required");
        if (isValidationOk) {
          this.focusElement("selected_mail_id");
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

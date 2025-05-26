<!--
  Copyright (C) 2025 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <div>
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
              <span>{{ $t("settings_metrics.title") }}</span>
            </cv-breadcrumb-item>
          </cv-breadcrumb>
        </cv-column>
      </cv-row>
      <cv-row>
        <cv-column class="subpage-title">
          <h3>{{ $t("settings_metrics.title") }}</h3>
        </cv-column>
      </cv-row>
      <cv-row>
        <cv-column>
          <NsInlineNotification
            v-if="error.getMetricsId"
            kind="error"
            :title="$t('action.get-metrics-ids')"
            :description="error.getMetricsId"
            :showCloseButton="false"
          />
        </cv-column>
        <cv-column>
          <NsInlineNotification
            v-if="error.getConfiguration"
            kind="error"
            :title="$t('action.get-metrics-configuration')"
            :description="error.getConfiguration"
            :showCloseButton="false"
          />
        </cv-column>
        <cv-column>
          <NsInlineNotification
            v-if="error.getSubscription"
            kind="error"
            :title="$t('action.get-subscription')"
            :description="error.getSubscription"
            :showCloseButton="false"
          />
        </cv-column>
        <cv-column>
          <NsInlineNotification
            v-if="error.getSmarthost"
            kind="error"
            :title="$t('action.get-smarthost')"
            :description="error.getSmarthost"
            :showCloseButton="false"
          />
        </cv-column>
      </cv-row>
      <cv-row>
        <cv-column>
          <cv-tile light>
            <h4 class="mg-bottom-md">
              {{ $t("settings_metrics.grafana") }}
            </h4>
            <div class="title-description mg-bottom-md">
              {{ $t("settings_metrics.grafana_description") }}
            </div>
            <div
              v-if="grafana_path && enable_grafana"
              class="title-description mg-bottom-md"
            >
              <a
                :href="'https://' + hostname"
                target="_blank"
                class="external-link-button"
              >
                <NsButton
                  :disabled="isGrafanaSectionLoading"
                  kind="ghost"
                  :icon="Launch20"
                >
                  {{ $t("settings_metrics.grafana_link") }}
                </NsButton>
              </a>
            </div>
            <cv-form @submit.prevent="setGrafana">
              <NsInlineNotification
                v-if="error.setGrafana"
                kind="error"
                :title="$t('action.set-metrics-configuration')"
                :description="error.setGrafana"
                :showCloseButton="false"
              />
              <NsToggle
                :label="$t('settings_metrics.access_grafana')"
                value="enable_grafana"
                :form-item="true"
                v-model="enable_grafana"
                :disabled="isGrafanaSectionLoading"
                ref="enable_grafana"
              >
                <template slot="tooltip">
                  <span>{{ $t("settings_metrics.grafana_tooltip") }}</span>
                </template>
                <template slot="text-left">{{
                  $t("common.disabled")
                }}</template>
                <template slot="text-right">{{
                  $t("common.enabled")
                }}</template>
              </NsToggle>
              <template v-if="enable_grafana">
                <div class="mg-top-sm icon-and-text mg-bottom-lg">
                  <NsSvg :svg="InformationFilled16" class="icon ns-info" />
                  <span>{{
                    $t(
                      "settings_metrics.note_on_grafana_cluster_admin_credentials"
                    )
                  }}</span>
                </div>
              </template>
              <NsInlineNotification
                v-if="error.setGrafana"
                kind="error"
                :title="$t('action.set-metrics-configuration')"
                :description="error.setGrafana"
                :showCloseButton="false"
              />
              <NsButton
                kind="secondary"
                :disabled="isGrafanaSectionLoading"
                :loading="isGrafanaSectionLoading"
                :icon="Save20"
                >{{ $t("settings_metrics.save") }}
              </NsButton>
            </cv-form>
          </cv-tile>
        </cv-column>
      </cv-row>
      <cv-row>
        <cv-column>
          <cv-tile light>
            <h4 class="mg-bottom-md">
              {{ $t("settings_metrics.alert_notifications") }}
            </h4>
            <div class="title-description mg-bottom-md">
              {{ $t("settings_metrics.alert_notifications_description") }}
            </div>
            <div class="maxwidth" v-if="!smarthost_enabled">
              <NsInlineNotification
                kind="info"
                :title="$t('settings_metrics.disabled_email_notification')"
                :description="
                  $t(
                    'settings_metrics.disabled_email_notifications_description'
                  )
                "
                :showCloseButton="false"
                @click="goToSettingsSmarthost"
                :actionLabel="$t('settings_metrics.go_to_email_notifications')"
              />
            </div>
            <div class="maxwidth" v-if="subscription.status === 'active'">
              <NsInlineNotification
                kind="info"
                :title="
                  $t('settings_metrics.notification_enabled_on_my_nethesis')
                "
                :description="
                  $t('settings_metrics.linked_to_an_active_subscription')
                "
                :showCloseButton="false"
              />
            </div>
            <cv-form @submit.prevent="setNotifications">
              <NsToggle
                v-if="subscription.status !== 'active'"
                :label="$t('settings_metrics.status_notifications')"
                value="status_notifications"
                :form-item="true"
                v-model="status_notifications"
                :disabled="isNotificationsSectionLoading || !smarthost_enabled"
                ref="status_notifications"
              >
                <template slot="tooltip">
                  <span>{{
                    $t("settings_metrics.notifications_tooltip")
                  }}</span>
                </template>
                <template slot="text-left">{{
                  $t("common.disabled")
                }}</template>
                <template slot="text-right">{{
                  $t("common.enabled")
                }}</template>
              </NsToggle>
              <template
                v-if="status_notifications || subscription.status === 'active'"
              >
                <NsTextInput
                  v-model="mail_from"
                  :label="$t('settings_metrics.mail_from')"
                  :form-item="true"
                  :placeholder="$t('settings_metrics.mail_from_placeholder')"
                  :invalid-message="error.mail_from"
                  :disabled="
                    isNotificationsSectionLoading || !smarthost_enabled
                  "
                  ref="mail_from"
                >
                </NsTextInput>
                <cv-text-area
                  :label="$t('settings_metrics.mail_to')"
                  v-model.trim="mail_to"
                  :invalid-message="error.mail_to"
                  :helper-text="$t('settings_metrics.mail_to_helper_text')"
                  :value="mail_to"
                  class="maxwidth textarea"
                  ref="mail_to"
                  :disabled="
                    isNotificationsSectionLoading || !smarthost_enabled
                  "
                >
                </cv-text-area>
              </template>
              <NsInlineNotification
                v-if="error.setNotifications"
                kind="error"
                :title="$t('action.set-metrics-notifications')"
                :description="error.setNotifications"
                :showCloseButton="false"
              />
              <NsButton
                kind="secondary"
                :loading="isNotificationsSectionLoading"
                :disabled="isNotificationsSectionLoading || !smarthost_enabled"
                :icon="Save20"
                >{{ $t("settings_metrics.save") }}
              </NsButton>
            </cv-form>
          </cv-tile>
        </cv-column>
      </cv-row>
    </cv-grid>
  </div>
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
import { mapGetters } from "vuex";

import NotificationService from "@/mixins/notification";
export default {
  name: "SettingsMetrics",
  components: {},
  mixins: [
    TaskService,
    UtilService,
    IconService,
    QueryParamService,
    PageTitleService,
    NotificationService,
  ],
  pageTitle() {
    return this.$t("settings_metrics.title");
  },
  data() {
    return {
      metricsId: "metrics1",
      prometheus_path: "",
      grafana_path: "",
      mail_to: [],
      mail_from: "",
      enable_grafana: false,
      status_notifications: false,
      hostname: "",
      smarthost_enabled: true,
      subscription: {
        auth_token: "",
        system_id: "",
        plan_name: { name: "", description: "" },
        expires: true,
        expire_date: "",
        status: "inactive",
        with_remote_support: false,
        error: "",
      },
      loading: {
        getGrafana: false,
        setGrafana: false,
        setNotifications: false,
        getConfiguration: false,
        getSubscription: false,
        getSmarthost: false,
      },
      error: {
        getGrafana: "",
        setGrafana: "",
        prometheus_path: "",
        grafana_path: "",
        mail_to: "",
        mail_from: "",
        enable_grafana: false,
        status_notifications: false,
        hostname: "",
        getConfiguration: "",
        getMetricsId: "",
        setNotifications: "",
        getSmarthost: "",
      },
    };
  },
  computed: {
    ...mapGetters(["leaderNode"]),
    isGrafanaSectionLoading() {
      return (
        this.loading.getGrafana ||
        this.loading.setGrafana ||
        this.loading.getSubscription ||
        this.loading.getSmarthost
      );
    },
    isNotificationsSectionLoading() {
      return (
        this.loading.getGrafana ||
        this.loading.setNotifications ||
        this.loading.getSubscription ||
        this.loading.getSmarthost
      );
    },
  },
  created() {
    this.getMetricsId();
    this.getSubscription();
    this.getSmarthost();
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
  methods: {
    goToSettingsSmarthost() {
      this.$router.push("/settings/smarthost");
    },
    async getSmarthost() {
      this.error.getSmarthost = "";
      this.loading.getSmarthost = true;
      const taskAction = "get-smarthost";
      const eventId = this.getUuid();

      // register to task completion
      this.$root.$once(
        taskAction + "-completed-" + eventId,
        this.getSmarthostCompleted
      );
      // register to task error
      this.$root.$once(
        taskAction + "-aborted-" + eventId,
        this.getSmarthostAborted
      );
      const res = await to(
        this.createClusterTask({
          action: taskAction,
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
            eventId,
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
    getSmarthostAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.getSmarthost = this.core.$t("error.generic_error");
      this.loading.getSmarthost = false;
    },
    getSmarthostCompleted(taskContext, taskResult) {
      const smarthost = taskResult.output;
      this.smarthost_enabled = smarthost.enabled;
      this.loading.getSmarthost = false;
    },
    async getSubscription() {
      this.error.getSubscription = "";
      this.loading.getSubscription = true;
      const taskAction = "get-subscription";
      const eventId = this.getUuid();
      // register to task completion
      this.$root.$once(
        taskAction + "-completed-" + eventId,
        this.getSubscriptionCompleted
      );
      // register to task error
      this.$root.$once(
        taskAction + "-aborted-" + eventId,
        this.getSubscriptionAborted
      );
      const res = await to(
        this.createClusterTask({
          action: taskAction,
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
            eventId,
          },
        })
      );

      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.getSubscription = this.getErrorMessage(err);
        this.loading.getSubscription = false;
        return;
      }
    },
    getSubscriptionAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.getSubscription = this.core.$t("error.generic_error");
      this.loading.getSubscription = false;
    },
    getSubscriptionCompleted(taskContext, taskResult) {
      const output = taskResult.output;
      this.termsUrl = output.terms_url;

      if (output.subscription == null) {
        this.subscription.status = "inactive";
        this.subscription.error = output.error;
        this.loading.getSubscription = false;
        return;
      }
      this.subscription = output.subscription;
      // set the notification toggle to true if the subscription is active
      this.status_notifications = true;
      this.loading.getSubscription = false;
    },
    async getMetricsId() {
      this.error.getGrafana = "";
      this.loading.getGrafana = true;
      const taskAction = "get-metrics-id";
      const eventId = this.getUuid();

      // register to task completion
      this.$root.$once(
        taskAction + "-completed-" + eventId,
        this.getMetricsIdCompleted
      );
      // register to task error
      this.$root.$once(
        taskAction + "-aborted-" + eventId,
        this.getMetricsIdAborted
      );
      const res = await to(
        this.createClusterTask({
          action: taskAction,
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
            eventId,
          },
        })
      );

      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.getMetricsId = this.getErrorMessage(err);
        this.loading.getGrafana = false;
        return;
      }
    },
    getMetricsIdAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.getMetricsId = this.core.$t("error.generic_error");
      this.loading.getMetricsId = false;
    },
    getMetricsIdCompleted(taskContext, taskResult) {
      const output = taskResult.output;
      this.metricsId = output.metrics_id;
      this.getConfiguration();
    },
    async getConfiguration() {
      this.error.getConfiguration = "";
      const taskAction = "get-configuration";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        taskAction + "-aborted-" + eventId,
        this.getConfigurationAborted
      );

      // register to task completion
      this.$root.$once(
        taskAction + "-completed-" + eventId,
        this.getConfigurationCompleted
      );

      const res = await to(
        this.createModuleTaskForApp(this.metricsId, {
          action: taskAction,
          extra: {
            title: this.$t("action.get-metrics-configuration"),
            isNotificationHidden: true,
            eventId,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.getConfiguration = this.getErrorMessage(err);
        this.loading.getGrafana = false;
        return;
      }
    },
    getConfigurationAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.getConfiguration = this.core.$t("error.generic_error");
      this.loading.getGrafana = false;
    },
    getConfigurationCompleted(taskContext, taskResult) {
      const config = taskResult.output;
      this.grafana_path = config.grafana_path;
      // build grafana hostname, grafana_path is a random string generated by the server
      // grafana is accessible at https://<hostname>/<grafana_path>
      this.hostname = window.location.hostname + "/" + config.grafana_path;
      // set enable_grafana toggle to true if grafana_path is not empty
      this.enable_grafana = this.grafana_path ? true : false;
      this.prometheus_path = config.prometheus_path;
      // build textarea content from mail_to array
      this.mail_to = config.mail_to ? config.mail_to.join("\n") : [];
      // status_notifications toggle is true if mail_to is not empty
      this.status_notifications = this.mail_to.length > 0;
      // if mail_from is empty, set it with default value from leader hostname
      const leaderHostname = this.leaderNode.vpn.endpoint.split(":")[0];
      this.mail_from = config.mail_from
        ? config.mail_from
        : "alertmanager@" + leaderHostname;
      this.loading.getGrafana = false;
    },
    async setGrafana() {
      this.error.setGrafana = "";
      this.loading.setGrafana = true;
      const taskAction = "configure-module";
      const eventId = this.getUuid();

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.setGrafanaCompleted
      );

      // register to task aborted
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.setGrafanaAborted
      );

      // register to validation error
      this.$root.$once(
        `${taskAction}-validation-failed-${eventId}`,
        this.setGrafanaValidationFailed
      );

      const res = await to(
        this.createModuleTaskForApp(this.metricsId, {
          action: taskAction,
          data: {
            grafana_path: this.enable_grafana ? "grafana" : "",
            prometheus_path: "",
          },
          extra: {
            title: this.$t("action.set-metrics-configuration"),
            isNotificationHidden: false,
            eventId,
            description: this.$t("common.processing"),
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.setGrafana = this.getErrorMessage(err);
        this.loading.setGrafana = false;
        return;
      }
    },
    setGrafanaAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.setGrafana = this.$t("error.generic_error");
      this.loading.setGrafana = false;
    },
    setGrafanaValidationFailed(validationErrors) {
      this.loading.setGrafana = false;
      let focusAlreadySet = false;

      for (const validationError of validationErrors) {
        const param = validationError.parameter;

        // set i18n error message
        this.error[param] = this.getI18nStringWithFallback(
          "settings_metrics." + validationError.error
        );

        if (!focusAlreadySet) {
          this.focusElement(param);
          focusAlreadySet = true;
        }
      }
    },
    setGrafanaCompleted() {
      this.loading.setGrafana = false;
      this.getMetricsId();
    },
    validateSetNotifications() {
      this.clearErrors(this);
      let isValidationOk = true;

      if (this.mail_to.length === 0) {
        this.error.mail_to = this.$t("common.required");

        if (isValidationOk) {
          this.focusElement("mail_to");
          isValidationOk = false;
        }
      }

      if (!this.mail_from) {
        this.error.mail_from = this.$t("common.required");

        if (isValidationOk) {
          this.focusElement("mail_from");
          isValidationOk = false;
        }
      }
      return isValidationOk;
    },
    async setNotifications() {
      const isValidationOk = this.validateSetNotifications();
      if (!isValidationOk) {
        return;
      }
      this.error.setNotifications = "";
      this.loading.setNotifications = true;
      const taskAction = "configure-module";
      const eventId = this.getUuid();

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.setNotificationsCompleted
      );

      // register to task aborted
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.setNotificationsAborted
      );

      // register to validation error
      this.$root.$once(
        `${taskAction}-validation-failed-${eventId}`,
        this.setNotificationsValidationFailed
      );

      const res = await to(
        this.createModuleTaskForApp(this.metricsId, {
          action: taskAction,
          data: {
            grafana_path: this.grafana_path ? this.grafana_path : "",
            prometheus_path: "",
            mail_to:
              this.status_notifications || this.subscription.status === "active"
                ? this.mail_to.split("\n")
                : [],
            mail_from: this.mail_from,
          },
          extra: {
            title: this.$t("action.set-metrics-notifications"),
            isNotificationHidden: false,
            eventId,
            description: this.$t("common.processing"),
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.setNotifications = this.getErrorMessage(err);
        this.loading.setNotifications = false;
        return;
      }
    },
    setNotificationsAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.setNotifications = this.$t("error.generic_error");
      this.loading.setNotifications = false;
    },
    setNotificationsValidationFailed(validationErrors) {
      this.loading.setNotifications = false;

      let focusAlreadySet = false;

      for (const validationError of validationErrors) {
        const param = validationError.parameter;

        // set i18n error message
        this.error[param] = this.getI18nStringWithFallback(
          "settings_metrics." + validationError.error
        );

        if (!focusAlreadySet) {
          this.focusElement(param);
          focusAlreadySet = true;
        }
      }
    },
    setNotificationsCompleted() {
      this.loading.setNotifications = false;
      this.getMetricsId();
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";

.icon-and-text {
  justify-content: flex-start;
}
.maxwidth {
  max-width: 38rem;
}
</style>

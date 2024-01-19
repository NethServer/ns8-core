<!--
  Copyright (C) 2023 Nethesis S.r.l.
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
              <span>{{ $t("settings_subscription.title") }}</span>
            </cv-breadcrumb-item>
          </cv-breadcrumb>
        </cv-column>
      </cv-row>
      <cv-row>
        <cv-column class="subpage-title">
          <h3>{{ $t("settings_subscription.title") }}</h3>
        </cv-column>
      </cv-row>
      <cv-row>
        <cv-column>
          <NsInlineNotification
            v-if="error.status"
            kind="error"
            :title="$t('action.get-subscription')"
            :description="error.status"
            :showCloseButton="false"
          />
        </cv-column>
      </cv-row>
      <cv-row>
        <cv-column>
          <cv-tile light>
            <h4 class="mg-bottom-md">
              {{ $t("settings_subscription.cluster_subscription") }}
            </h4>
            <div class="title-description mg-bottom-xlg">
              {{ $t("settings_subscription.cluster_subscription_description") }}
            </div>
            <cv-form @submit.prevent="setSubscription">
              <template v-if="status == 'inactive'">
                <NsTextInput
                  :label="$t('settings_subscription.authentication_token')"
                  v-model="auth_token"
                  :invalid-message="$t(error.auth_token)"
                  :placeholder="
                    $t('settings_subscription.authentication_token_placeholder')
                  "
                  ref="auth_token"
                  :helper-text="
                    $t('settings_subscription.authentication_token_helper')
                  "
                  required
                >
                </NsTextInput>
              </template>
              <template v-else>
                <div class="key-value-setting">
                  <span class="label">{{
                    $t("settings_subscription.system_id")
                  }}</span>
                  <span class="value">{{ system_id }}</span>
                </div>
                <div class="key-value-setting">
                  <span class="label">{{
                    $t("settings_subscription.plan_name")
                  }}</span>
                  <span class="value">{{ plan_name }}</span>
                </div>
                <div class="key-value-setting">
                  <span class="label">{{
                    $t("settings_subscription.expire_date")
                  }}</span>
                  <span class="value">{{ expire_date }}</span>
                </div>
                <div class="key-value-setting">
                  <span class="label">{{
                    $t("settings_subscription.status")
                  }}</span>
                  <span class="value">
                    <cv-tag
                      v-if="status"
                      kind="green"
                      :label="$t('common.active')"
                      size="sm"
                      class="no-margin"
                    ></cv-tag>
                    <cv-tag
                      v-else
                      kind="high-contrast"
                      :label="$t('common.not_active')"
                      size="sm"
                      class="no-margin"
                    ></cv-tag>
                  </span>
                </div>
              </template>
              <NsInlineNotification
                v-if="error.getSubscription"
                kind="error"
                :title="$t('action.register-cluster-subscription')"
                :description="error.getSubscription"
                :showCloseButton="false"
              />
              <NsButton
                kind="primary"
                v-if="status == 'inactive'"
                :loading="loading.getSubscription || loading.setSubscription"
                :disabled="
                  loading.register_cluster_subsciption ||
                  loading.setSubscription
                "
                :icon="Badge20"
                >{{ $t("settings_subscription.request_subscription") }}
              </NsButton>
              <NsButton
                v-if="status == 'active'"
                kind="tertiary"
                :loading="loading.getSubscription || loading.setSubscription"
                :disabled="
                  loading.register_cluster_subsciption ||
                  loading.setSubscription
                "
                :icon="TrashCan20"
                >{{ $t("settings_subscription.remove_subscription") }}
              </NsButton>
            </cv-form>
          </cv-tile>
        </cv-column>
      </cv-row>
      <cv-row v-if="with_remote_support">
        <cv-column>
          <cv-tile light>
            <h4>
              {{ $t("settings_subscription.remote_support") }}
            </h4>
            <div class="title-description mg-bottom-xlg">
              {{ $t("settings_subscription.remote_support_description") }}
            </div>
            <NsInlineNotification
              v-if="error.get_support"
              kind="error"
              :title="
                $t('settings_subscription.cannot_retrieve_get_support_status')
              "
              :description="error.get_support"
              :showCloseButton="false"
            />
            <cv-skeleton-text
              v-if="loading.getSubscription"
              :paragraph="true"
              :line-count="5"
              heading
            ></cv-skeleton-text>
            <template v-else>
              <template v-if="!active">
                <div>
                  <NsInlineNotification
                    v-if="error.request_support"
                    kind="error"
                    :title="$t('settings_subscription.cannot_request_support')"
                    :description="error.request_support"
                    :showCloseButton="false"
                  />

                  <NsButton
                    kind="secondary"
                    :icon="Play20"
                    @click="startSessionSupport"
                  >
                    {{ $t("settings_subscription.start_session_support") }}
                  </NsButton>
                </div>
              </template>
              <template v-else>
                <span class="label">
                  {{ $t("settings_subscription.session_id") }}
                </span>
                <NsCodeSnippet
                  :copyTooltip="$t('common.copy_to_clipboard')"
                  :copy-feedback="$t('common.copied_to_clipboard')"
                  :feedback-aria-label="$t('common.copied_to_clipboard')"
                  :wrap-text="true"
                  :moreText="$t('common.show_more')"
                  :lessText="$t('common.show_less')"
                  :helper-text="$t('settings_tls_certificates.fqdn_helper')"
                  hideExpandButton
                  class="mg-top-sm"
                  >{{ session_id }}</NsCodeSnippet
                >
                <div class="mg-top-sm mg-bottom-lg helper-text">
                  {{
                    $t("settings_subscription.paste_session_id_to_support_team")
                  }}
                </div>
                <div class="mg-top-sm icon-and-text mg-bottom-lg">
                  <NsSvg :svg="InformationFilled16" class="icon ns-info" />
                  <span>{{
                    $t("settings_subscription.remote_support_in_progress")
                  }}</span>
                </div>
                <NsButton
                  class="mg-top-sm mg-bottom-lg"
                  kind="primary"
                  :icon="Stop20"
                  @click="stopSessionSupport"
                >
                  {{ $t("settings_subscription.stop_session_support") }}
                </NsButton>
              </template>
            </template>
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
import { mapActions } from "vuex";

import NotificationService from "@/mixins/notification";
import Play20 from "@carbon/icons-vue/es/play/20";
import Stop20 from "@carbon/icons-vue/es/stop/20";
export default {
  name: "SettingsSubscription",
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
    return this.$t("settings_subscription.title");
  },
  data() {
    return {
      Play20,
      Stop20,
      auth_token: "",
      system_id: "",
      vpn_cert_cn:
        "C=IT, ST=PU, L=Pesaro, O=Nethesis, OU=Support, CN=Nethesis CA, name=sos, emailAddress=support@nethesis.it",
      dartagnan_url: "https://my.nethserver.com/api",
      support_user: "nethsupport",
      provider: "nscom",
      system_url: "https://my.nethserver.com/servers/5329",
      plan_name: { name: "", description: "" },
      expires: true,
      expire_date: "",
      status: "inactive",
      active: false,
      session_id: "",
      with_remote_support: false,
      loading: {
        getSubscription: false,
        get_support: false,
        setSubscription: false,
      },
      error: {
        status: "",
        auth_token: "",
        getSubscription: "",
        get_support: "",
        setSubscription: "",
        request_support: "",
        startSessionSupport: "",
      },
    };
  },
  computed: {},
  created() {
    this.nodeId = this.$route.params.nodeId;

    this.getSubscription();
    this.getSupportSession();
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
    ...mapActions(["setSubscriptionInStore"]),
    async getSubscription() {
      this.clearErrors();
      this.loading.getSubscription = true;
      const taskAction = "get-subscription";

      // register to task completion
      this.$root.$once(
        taskAction + "-completed",
        this.getSubscriptionCompleted
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
        this.error.getSubscription = this.getErrorMessage(err);
        this.loading.getSubscription = false;
        return;
      }
    },
    getSubscriptionCompleted(taskContext, taskResult) {
      const output = taskResult.output;
      if (output.subscription == null) {
        this.status = "inactive";
        this.loading.getSubscription = false;
        return;
      }
      this.auth_token = output.subscription.auth_token;
      this.system_id = output.subscription.system_id;
      this.vpn_cert_cn = output.subscription.vpn_cert_cn;
      this.dartagnan_url = output.subscription.dartagnan_url;
      this.support_user = output.subscription.support_user;
      this.provider = output.subscription.provider;
      this.system_url = output.subscription.system_url;
      this.plan_name = output.subscription.plan_name;
      this.expires = output.subscription.expires;
      this.expire_date = output.subscription.expire_date;
      this.status = output.subscription.status;
      this.with_remote_support = output.subscription.with_remote_support;
      this.loading.getSubscription = false;
    },
    async setSubscription() {
      this.error.getSubscription = "";
      this.error.setSubscription = "";
      this.loading.setSubscription = true;

      const taskAction = "set-subscription";

      // register to task completion
      this.$root.$once(
        taskAction + "-completed",
        this.setSubscriptionCompleted
      );

      // register to task error
      this.$root.$once(taskAction + "-aborted", this.setSubscriptionFailed);

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            subscription:
              this.status === "inactive"
                ? { auth_token: this.auth_token }
                : null,
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
        this.error.getSubscription = this.getErrorMessage(err);
        this.loading.getSubscription = false;
        return;
      }
    },
    setSubscriptionFailed(validationErrors) {
      this.loading.setSubscription = false;
      this.status = "inactive";

      let focusAlreadySet = false;

      for (const validationError of validationErrors) {
        const param = validationError.parameter;

        // set i18n error message
        this.error[param] = this.getI18nStringWithFallback(
          "subscription." + validationError.error,
          "error." + validationError.error
        );

        if (!focusAlreadySet) {
          this.focusElement(param);
          focusAlreadySet = true;
        }
      }
    },

    setSubscriptionCompleted() {
      this.status = "active";
      this.auth_token = "";
      this.loading.setSubscription = false;
      this.getSubscription();
    },

    async getSupportSession() {
      // this.error.getSupportSession = "";
      this.active = false;
      const taskAction = "get-support-session";

      // register to task completion
      this.$root.$once(
        taskAction + "-completed",
        this.getSupportSessionCompleted
      );

      const res = await to(
        this.createNodeTask(1, {
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
        this.error.getSupportSession = this.getErrorMessage(err);
        this.loading.getSupportSession = false;
        return;
      }
    },
    getSupportSessionCompleted(taskContext, taskResult) {
      const output = taskResult.output;
      this.session_id = output.session_id;
      this.active = output.active;
    },

    async startSessionSupport() {
      this.error.request_support = "";
      this.loading.startSessionSupport = true;

      const taskAction = "start-support-session";

      // register to task completion
      this.$root.$once(
        taskAction + "-completed",
        this.startSessionSupportCompleted
      );

      // register to task error
      this.$root.$once(taskAction + "-aborted", this.startSessionSupportFailed);

      const res = await to(
        this.createNodeTask(1, {
          action: taskAction,
          data: {
            session_id: this.session_id,
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
        this.error.getSubscription = this.getErrorMessage(err);
        this.loading.getSubscription = false;
        return;
      }
    },
    startSessionSupportFailed(validationErrors) {
      this.loading.startSessionSupport = false;
      this.status = false;

      let focusAlreadySet = false;

      for (const validationError of validationErrors) {
        const param = validationError.parameter;

        // set i18n error message
        this.error[param] = this.getI18nStringWithFallback(
          "subscription." + validationError.error,
          "error." + validationError.error
        );

        if (!focusAlreadySet) {
          this.focusElement(param);
          focusAlreadySet = true;
        }
      }
    },

    startSessionSupportCompleted() {
      this.loading.startSessionSupport = false;
      this.getSupportSession();
    },

    async stopSessionSupport() {
      this.error.request_support = "";
      this.loading.stopSessionSupport = true;

      const taskAction = "stop-support-session";

      // register to task completion
      this.$root.$once(
        taskAction + "-completed",
        this.stopSessionSupportCompleted
      );

      // register to task error
      this.$root.$once(taskAction + "-aborted", this.stopSessionSupportFailed);

      const res = await to(
        this.createNodeTask(1, {
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
        this.error.getSubscription = this.getErrorMessage(err);
        this.loading.getSubscription = false;
        return;
      }
    },
    stopSessionSupportFailed(validationErrors) {
      this.loading.stopSessionSupport = false;
      this.status = false;

      let focusAlreadySet = false;

      for (const validationError of validationErrors) {
        const param = validationError.parameter;

        // set i18n error message
        this.error[param] = this.getI18nStringWithFallback(
          "subscription." + validationError.error,
          "error." + validationError.error
        );

        if (!focusAlreadySet) {
          this.focusElement(param);
          focusAlreadySet = true;
        }
      }
    },

    stopSessionSupportCompleted() {
      this.loading.stopSessionSupport = false;
      this.getSupportSession();
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";

.icon-and-text {
  justify-content: flex-start;
}
.helper-text {
  font-size: 0.75rem;
  line-height: 1.33333;
  letter-spacing: 0.32px;
  z-index: 0;
  width: 100%;
  margin-top: 0.25rem;
  color: #525252;
  opacity: 1;
}
</style>

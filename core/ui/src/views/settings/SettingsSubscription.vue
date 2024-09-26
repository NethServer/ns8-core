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
            v-if="error.getSubscription"
            kind="error"
            :title="$t('action.get-subscription')"
            :description="error.getSubscription"
            :showCloseButton="false"
          />
        </cv-column>
      </cv-row>
      <cv-row>
        <cv-column>
          <NsInlineNotification
            v-if="subscription.error === 'os_not_supported'"
            kind="info"
            :title="$t('settings_subscription.subscription_cannot_be_enabled')"
            :description="$t('settings_subscription.os_not_supported')"
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
            <div class="title-description mg-bottom-md">
              <i18n
                path="settings_subscription.cluster_subscription_description"
                tag="p"
              >
                <template v-slot:enterprise>
                  <cv-link
                    href="https://my.nethesis.it"
                    target="_blank"
                    rel="noreferrer"
                  >
                    my.nethesis.it
                  </cv-link>
                </template>
                <template v-slot:subscription>
                  <cv-link
                    href="https://my.nethserver.com"
                    target="_blank"
                    rel="noreferrer"
                  >
                    my.nethserver.com
                  </cv-link>
                </template>
              </i18n>
            </div>
            <cv-form @submit.prevent="setSubscription">
              <cv-skeleton-text
                v-if="loading.getSubscription"
                :paragraph="true"
                :line-count="6"
                width="80%"
              ></cv-skeleton-text>
              <template
                v-if="
                  !loading.getSubscription && subscription.status == 'inactive'
                "
              >
                <NsTextInput
                  :disabled="subscription.error === 'os_not_supported'"
                  :label="$t('settings_subscription.authentication_token')"
                  v-model="subscription.auth_token"
                  :invalid-message="error.auth_token"
                  :placeholder="
                    $t('settings_subscription.authentication_token_placeholder')
                  "
                  ref="auth_token"
                  :helper-text="
                    $t('settings_subscription.authentication_token_helper')
                  "
                >
                </NsTextInput>
                <!-- terms and conditions -->
                <div class="mg-bottom-lg">
                  <NsCheckbox
                    v-model="agreeTerms"
                    value="checkTerms"
                    class="no-mg-bottom"
                  >
                    <template slot="label">
                      <i18n
                        path="settings_subscription.agree_terms_before_register"
                        tag="span"
                      >
                        <template v-slot:terms>
                          <cv-link
                            :href="termsUrl"
                            target="_blank"
                            rel="noreferrer"
                            class="inline"
                          >
                            {{ $t("common.terms_and_conditions") }}
                          </cv-link>
                        </template>
                      </i18n>
                    </template>
                  </NsCheckbox>
                  <!-- invalid message for terms and conditions -->
                  <div
                    v-if="error.agreeTerms"
                    class="validation-failed-invalid-message"
                  >
                    {{ error.agreeTerms }}
                  </div>
                </div>
              </template>
              <template
                v-if="
                  !loading.getSubscription && subscription.status == 'active'
                "
              >
                <div class="key-value-setting">
                  <span class="label">{{
                    $t("settings_subscription.system_id")
                  }}</span>
                  <span class="value">{{ subscription.system_id }}</span>
                </div>
                <div class="key-value-setting">
                  <span class="label">{{
                    $t("settings_subscription.plan_name")
                  }}</span>
                  <span class="value">{{ subscription.plan_name }}</span>
                </div>
                <div class="key-value-setting">
                  <span class="label">{{
                    $t("settings_subscription.expire_date")
                  }}</span>
                  <span class="value">{{
                    subscription.expire_date === "-1"
                      ? $t("settings_subscription.no_expiration")
                      : formatDate(
                          new Date(subscription.expire_date),
                          "yyyy-MM-dd"
                        )
                  }}</span>
                </div>
                <div class="key-value-setting">
                  <span class="label">{{
                    $t("settings_subscription.status")
                  }}</span>
                  <span class="value">
                    <cv-tag
                      v-if="subscription.status"
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
                v-if="error.setSubscription"
                kind="error"
                :title="$t('action.register-cluster-subscription')"
                :description="error.setSubscription"
                :showCloseButton="false"
              />
              <NsButton
                kind="primary"
                v-if="
                  !loading.getSubscription && subscription.status == 'inactive'
                "
                :loading="loading.getSubscription || loading.setSubscription"
                :disabled="
                  loading.getSubscription ||
                  loading.setSubscription ||
                  subscription.error === 'os_not_supported'
                "
                :icon="Badge20"
                >{{ $t("settings_subscription.request_subscription") }}
              </NsButton>
            </cv-form>
            <NsButton
              v-if="subscription.status == 'active'"
              kind="tertiary"
              :loading="loading.getSubscription || loading.setSubscription"
              :disabled="loading.getSubscription || loading.setSubscription"
              @click="showRemoveSubcriptionModal"
              :icon="TrashCan20"
              >{{ $t("settings_subscription.remove_subscription") }}
            </NsButton>
          </cv-tile>
        </cv-column>
      </cv-row>
      <cv-row
        v-if="
          subscription.with_remote_support && subscription.status == 'active'
        "
      >
        <cv-column>
          <cv-tile light>
            <h4>
              {{ $t("settings_subscription.remote_support") }}
            </h4>
            <div class="title-description mg-bottom-xlg">
              {{
                $t("settings_subscription.remote_support_description", {
                  productName: $root.config.PRODUCT_NAME,
                })
              }}
            </div>
            <NsInlineNotification
              v-if="error.getSupportSession"
              kind="error"
              :title="
                $t('settings_subscription.cannot_retrieve_get_support_status')
              "
              :description="error.getSupportSession"
              :showCloseButton="false"
            />
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
                  :loading="loading.startSessionSupport"
                  :disabled="loading.startSessionSupport"
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
                :loading="loading.stopSessionSupport"
                :disabled="loading.stopSessionSupport"
                class="mg-top-sm mg-bottom-lg"
                kind="primary"
                :icon="Stop20"
                @click="stopSessionSupport"
              >
                {{ $t("settings_subscription.stop_session_support") }}
              </NsButton>
            </template>
          </cv-tile>
        </cv-column>
      </cv-row>
    </cv-grid>
    <NsModal
      size="default"
      kind="danger"
      :visible="isShownRemoveSubcription"
      @modal-hidden="hideRemoveSubcriptionModal"
      @primary-click="removeSubscription"
    >
      <template slot="title">{{
        $t("settings_subscription.remove_cluster_subscription_title")
      }}</template>
      <template slot="content">
        <div>
          {{
            $t("settings_subscription.remove_cluster_subscription_description")
          }}
        </div>
      </template>
      <template slot="secondary-button">{{ $t("common.cancel") }}</template>
      <template slot="primary-button">{{ $t("common.remove") }}</template>
    </NsModal>
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
  DateTimeService,
} from "@nethserver/ns8-ui-lib";
import { mapGetters } from "vuex";

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
    DateTimeService,
  ],
  pageTitle() {
    return this.$t("settings_subscription.title");
  },
  data() {
    return {
      isShownRemoveSubcription: false,
      Play20,
      Stop20,
      enterprise_link:
        "<a href='https://my.nethesis.it' target=_blank'>my.nethesis.it</a>",
      subscription_link:
        "<a href='https://my.nethserver.com' target=_blank'>my.nethserver.com</a>",
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
      active: false,
      session_id: "",
      termsUrl: "",
      agreeTerms: false,
      loading: {
        getSubscription: false,
        setSubscription: false,
        startSessionSupport: false,
        stopSessionSupport: false,
      },
      error: {
        status: "",
        auth_token: "",
        getSubscription: "",
        setSubscription: "",
        request_support: "",
        startSessionSupport: "",
        removeSubscription: "",
        unknown_token: "",
        agreeTerms: "",
      },
    };
  },
  computed: { ...mapGetters(["leaderNode"]) },
  created() {
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
    showRemoveSubcriptionModal() {
      this.isShownRemoveSubcription = true;
    },
    hideRemoveSubcriptionModal() {
      this.isShownRemoveSubcription = false;
    },
    async getSubscription() {
      this.error.getSubscription = "";
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
      this.termsUrl = output.terms_url;

      if (output.subscription == null) {
        this.subscription.status = "inactive";
        this.subscription.error = output.error;
        this.loading.getSubscription = false;
        return;
      }
      this.subscription = output.subscription;
      this.loading.getSubscription = false;
    },
    validateSetSubscription() {
      this.clearErrors(this);
      let isValidationOk = true;

      // auth token required

      if (!this.subscription.auth_token) {
        this.error.auth_token = this.$t("common.required");

        if (isValidationOk) {
          this.focusElement("auth_token");
          isValidationOk = false;
        }
      } else {
        // validate auth_token if length is <= 32 and > 128

        if (
          this.subscription.auth_token.length <= 32 ||
          this.subscription.auth_token.length > 128
        ) {
          this.error.auth_token = this.$t(
            "settings_subscription.must_be_32_chars_but_less_than_128"
          );

          if (isValidationOk) {
            this.focusElement("auth_token");
            isValidationOk = false;
          }
        }
      }

      // need to agree to terms and conditions

      if (!this.agreeTerms) {
        this.error.agreeTerms = this.$t("common.terms_required");
        isValidationOk = false;
      }
      return isValidationOk;
    },
    async setSubscription() {
      const isValidationOk = this.validateSetSubscription();
      if (!isValidationOk) {
        return;
      }
      this.error.setSubscription = "";
      this.error.auth_token = "";
      this.loading.setSubscription = true;
      const taskAction = "set-subscription";

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed`,
        this.setSubscriptionCompleted
      );

      // register to task aborted
      this.$root.$once(`${taskAction}-aborted`, this.setSubscriptionAborted);

      // register to task error
      this.$root.$once(
        `${taskAction}-validation-failed`,
        this.setSubscriptionFailed
      );

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            subscription: { auth_token: this.subscription.auth_token },
          },
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: false,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.setSubscription = this.getErrorMessage(err);
        this.loading.setSubscription = false;
        return;
      }
    },
    setSubscriptionAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.setSubscription = this.$t("error.generic_error");
      this.loading.setSubscription = false;
    },
    setSubscriptionFailed(validationErrors) {
      this.loading.setSubscription = false;
      this.subscription.status = "inactive";

      let focusAlreadySet = false;

      for (const validationError of validationErrors) {
        const param = validationError.parameter;

        // set i18n error message
        this.error[param] = this.getI18nStringWithFallback(
          "settings_subscription." + validationError.error
        );

        if (!focusAlreadySet) {
          this.focusElement(param);
          focusAlreadySet = true;
        }
      }
    },
    setSubscriptionCompleted() {
      this.subscription.auth_token = "";
      this.loading.setSubscription = false;
      this.agreeTerms = false;
      this.getSubscription();
    },
    async removeSubscription() {
      this.error.removeSubscription = "";
      this.loading.setSubscription = true;

      const taskAction = "set-subscription";

      // register to task completion
      this.$root.$once(
        taskAction + "-completed",
        this.removeSubscriptionCompleted
      );

      // register to task error
      this.$root.$once(
        taskAction + "validation-failed",
        this.removeSubscriptionFailed
      );

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            subscription: null,
          },
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: false,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.removeSubscription = this.getErrorMessage(err);
        this.loading.setSubscription = false;
        return;
      }
    },
    removeSubscriptionFailed(validationErrors) {
      this.loading.setSubscription = false;

      let focusAlreadySet = false;

      for (const validationError of validationErrors) {
        const param = validationError.parameter;

        // set i18n error message
        this.error[param] = this.getI18nStringWithFallback(
          "settings_subscription." + validationError.error,
          "error." + validationError.error
        );

        if (!focusAlreadySet) {
          this.focusElement(param);
          focusAlreadySet = true;
        }
      }
    },
    removeSubscriptionCompleted() {
      // if inactive, we remove subscription and force stop session support
      this.stopSessionSupport();
      this.subscription.auth_token = "";
      this.loading.setSubscription = false;
      this.hideRemoveSubcriptionModal();
      this.getSubscription();
    },

    async getSupportSession() {
      this.error.getSupportSession = "";
      this.active = false;
      const taskAction = "get-support-session";

      // register to task completion
      this.$root.$once(
        taskAction + "-completed",
        this.getSupportSessionCompleted
      );

      const res = await to(
        this.createNodeTask(this.leaderNode.id, {
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
      this.$root.$once(
        taskAction + "validation-failed",
        this.startSessionSupportFailed
      );

      const res = await to(
        this.createNodeTask(this.leaderNode.id, {
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
        this.error.request_support = this.getErrorMessage(err);
        this.loading.startSessionSupport = false;
        return;
      }
    },
    startSessionSupportFailed(validationErrors) {
      this.loading.startSessionSupport = false;

      let focusAlreadySet = false;

      for (const validationError of validationErrors) {
        const param = validationError.parameter;

        // set i18n error message
        this.error[param] = this.getI18nStringWithFallback(
          "settings_subscription." + validationError.error,
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
      this.$root.$once(
        taskAction + "validation-failed",
        this.stopSessionSupportFailed
      );

      const res = await to(
        this.createNodeTask(this.leaderNode.id, {
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
        this.error.request_support = this.getErrorMessage(err);
        this.loading.stopSessionSupport = false;
        return;
      }
    },
    stopSessionSupportFailed(validationErrors) {
      this.loading.stopSessionSupport = false;

      let focusAlreadySet = false;

      for (const validationError of validationErrors) {
        const param = validationError.parameter;

        // set i18n error message
        this.error[param] = this.getI18nStringWithFallback(
          "settings_subscription." + validationError.error,
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

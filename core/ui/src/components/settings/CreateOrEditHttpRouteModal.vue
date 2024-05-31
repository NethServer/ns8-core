<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <NsModal
    size="default"
    :visible="isShown"
    :primary-button-disabled="loading.setRoute"
    :isLoading="loading.setRoute"
    @modal-hidden="onModalHidden"
    @primary-click="createOrEditRoute"
  >
    <template slot="title">{{
      isEditing
        ? $t("settings_http_routes.edit_route_route", { route: route.instance })
        : $t("settings_http_routes.create_route")
    }}</template>
    <template slot="content">
      <cv-form @submit.prevent="createOrEditRoute">
        <NsTextInput
          v-model.trim="instance"
          :placeholder="$t('settings_http_routes.name_placeholder')"
          :label="$t('settings_http_routes.name')"
          :invalid-message="error.instance"
          :disabled="isEditing || loading.setRoute"
          data-modal-primary-focus
          ref="instance"
        />
        <NsComboBox
          v-model="selectedNodeId"
          :label="$t('common.choose')"
          :title="$t('common.node')"
          :auto-filter="true"
          :auto-highlight="true"
          :options="nodes"
          :disabled="isEditing || loading.setRoute"
          :invalid-message="error.node"
          light
          tooltipAlignment="start"
          tooltipDirection="top"
          ref="node"
        >
          <template slot="tooltip">
            {{ $t("settings_http_routes.node_helper") }}
          </template>
        </NsComboBox>
        <NsTextInput
          v-model.trim="url"
          :placeholder="$t('settings_http_routes.url_placeholder')"
          :label="$t('settings_http_routes.url')"
          :helper-text="$t('settings_http_routes.url_helper')"
          :invalid-message="error.url"
          :disabled="loading.setRoute"
          ref="url"
        >
          <template slot="tooltip">
            <span>{{ $t("settings_http_routes.url_tooltip") }}</span>
          </template>
        </NsTextInput>
        <NsToggle
          :label="$t('settings_http_routes.skip_cert_verify')"
          value="stripPrefixValue"
          :form-item="true"
          v-model="skip_cert_verify"
          :disabled="loading.setRoute || !url.startsWith('https://')"
          ref="skip_cert_verify"
        >
          <template slot="tooltip">
            <span>{{ $t("settings_http_routes.skipCertVerify_tooltip") }}</span>
          </template>
          <template slot="text-left">{{ $t("common.disabled") }}</template>
          <template slot="text-right">{{ $t("common.enabled") }}</template>
        </NsToggle>
        <NsTextInput
          v-model.trim="host"
          :placeholder="$t('settings_http_routes.host_placeholder')"
          :label="
            $t('settings_http_routes.host') + ' (' + $t('common.optional') + ')'
          "
          :helper-text="$t('settings_http_routes.host_helper')"
          :invalid-message="error.host"
          :disabled="loading.setRoute"
          ref="host"
        >
          <template slot="tooltip">
            <span>{{ $t("settings_http_routes.host_tooltip") }}</span>
          </template>
        </NsTextInput>
        <NsTextInput
          v-model.trim="path"
          :placeholder="$t('settings_http_routes.path_placeholder')"
          :label="
            $t('settings_http_routes.path') + ' (' + $t('common.optional') + ')'
          "
          :invalid-message="error.path"
          :disabled="loading.setRoute"
          ref="path"
        >
          <template slot="tooltip">
            <span>{{ $t("settings_http_routes.path_tooltip") }}</span>
          </template>
        </NsTextInput>
        <NsToggle
          :label="$t('settings_http_routes.strip_prefix')"
          value="stripPrefixValue"
          :form-item="true"
          v-model="strip_prefix"
          :disabled="loading.setRoute"
          ref="tls"
        >
          <template slot="tooltip">
            <i18n path="settings_http_routes.strip_prefix_tooltip" tag="span">
              <template v-slot:traefikPrefixesOption>
                <cv-link
                  href="https://doc.traefik.io/traefik/middlewares/http/stripprefix/"
                  target="_blank"
                  rel="noreferrer"
                >
                  {{ $t("settings_http_routes.traefik_prefixes_option") }}
                </cv-link>
              </template>
            </i18n>
          </template>
          <template slot="text-left">{{ $t("common.disabled") }}</template>
          <template slot="text-right">{{ $t("common.enabled") }}</template>
        </NsToggle>
        <NsToggle
          :label="$t('settings_http_routes.request_lets_encrypt_certificate')"
          value="letsEncryptValue"
          :form-item="true"
          v-model="lets_encrypt"
          :disabled="loading.setRoute"
          ref="lets_encrypt"
        >
          <template slot="text-left">{{ $t("common.disabled") }}</template>
          <template slot="text-right">{{ $t("common.enabled") }}</template>
        </NsToggle>
        <NsToggle
          :label="$t('settings_http_routes.http2https')"
          value="http2httpsValue"
          :form-item="true"
          v-model="http2https"
          :disabled="loading.setRoute"
          ref="http2https"
        >
          <template slot="text-left">{{ $t("common.disabled") }}</template>
          <template slot="text-right">{{ $t("common.enabled") }}</template>
        </NsToggle>
        <!-- need to wrap error notification inside a div: custom elements like NsInlineNotification don't have scrollIntoView() function -->
        <div ref="setRouteError">
          <NsInlineNotification
            v-if="error.setRoute"
            kind="error"
            :title="$t('action.set-route')"
            :description="error.setRoute"
            :showCloseButton="false"
          />
        </div>
      </cv-form>
    </template>
    <template slot="secondary-button">{{ $t("common.cancel") }}</template>
    <template slot="primary-button">{{
      isEditing
        ? $t("settings_http_routes.edit_route")
        : $t("settings_http_routes.create_route")
    }}</template>
  </NsModal>
</template>

<script>
import to from "await-to-js";
import { UtilService, TaskService } from "@nethserver/ns8-ui-lib";

export default {
  name: "CreateOrEditHttpRouteModal",
  mixins: [UtilService, TaskService],
  props: {
    isShown: Boolean,
    nodes: {
      type: Array,
      required: true,
    },
    defaultNodeId: {
      type: String,
      default: "",
    },
    allRoutes: {
      type: Array,
      required: true,
    },
    route: { type: [Object, null] },
    isEditing: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      instance: "",
      selectedNodeId: "",
      url: "",
      host: "",
      path: "",
      skip_cert_verify: false,
      lets_encrypt: false,
      http2https: false,
      strip_prefix: false,
      loading: {
        setRoute: false,
      },
      error: {
        skip_cert_verify: "",
        setRoute: "",
        instance: "",
        node: "",
        url: "",
        host: "",
        path: "",
      },
    };
  },
  watch: {
    defaultNodeId: function () {
      this.updateSelectedNodeId();
    },
    url(newUrl) {
      if (!newUrl.startsWith("https://")) {
        this.skip_cert_verify = false;
      }
    },
    isShown: function () {
      if (this.isShown) {
        this.clearErrors();

        if (this.isEditing) {
          // edit route
          this.instance = this.route.instance;
          this.selectedNodeId = this.route.nodeId;
          this.url = this.route.url;
          this.host = this.route.host;
          this.path = this.route.path;
          this.strip_prefix = this.route.strip_prefix;
          this.lets_encrypt = this.route.lets_encrypt;
          this.http2https = this.route.http2https;
          this.skip_cert_verify = this.route.skip_cert_verify;
        }
      } else {
        // closing modal
        if (this.isEditing) {
          this.clearFields();
        }
      }
    },
    "error.setRoute": function () {
      if (this.error.setRoute) {
        // scroll to notification error

        this.$nextTick(() => {
          const el = this.$refs.setRouteError;
          this.scrollToElement(el);
        });
      }
    },
  },
  created() {
    this.updateSelectedNodeId();
  },
  methods: {
    updateSelectedNodeId() {
      if (this.defaultNodeId != "all") {
        this.selectedNodeId = this.defaultNodeId;
      } else {
        this.selectedNodeId = "";
      }
    },
    onModalHidden() {
      this.clearErrors();
      this.$emit("hide");
    },
    validateSetRoute(isEditingRoute) {
      this.clearErrors();

      let isValidationOk = true;

      // route name

      if (!this.instance) {
        this.error.instance = this.$t("common.required");

        if (isValidationOk) {
          this.focusElement("instance");
          isValidationOk = false;
        }
      }

      // check if route name already exists if user is creating a route
      if (!isEditingRoute) {
        const duplicatedRoute = this.allRoutes.find(
          (route) => route.instance === this.instance
        );

        if (duplicatedRoute) {
          this.error.instance = this.$t(
            "settings_http_routes.route_already_exists"
          );

          if (isValidationOk) {
            this.focusElement("instance");
            isValidationOk = false;
          }
        }
      }

      // node

      if (!this.selectedNodeId) {
        this.error.node = this.$t("common.required");

        if (isValidationOk) {
          this.focusElement("node");
          isValidationOk = false;
        }
      }

      // url

      if (!this.url) {
        this.error.url = this.$t("common.required");

        if (isValidationOk) {
          this.focusElement("url");
          isValidationOk = false;
        }
      }

      // host or path

      if (!this.host && !this.path) {
        this.error.host = this.$t("settings_http_routes.host_or_path_required");
        this.error.path = this.$t("settings_http_routes.host_or_path_required");

        if (isValidationOk) {
          this.focusElement("host");
          isValidationOk = false;
        }
      } else {
        if (this.path && !this.path.startsWith("/")) {
          // path must start with "/"
          this.error.path = this.$t(
            "settings_http_routes.path_must_start_with_slash"
          );

          if (isValidationOk) {
            this.focusElement("path");
            isValidationOk = false;
          }
        }
      }
      return isValidationOk;
    },
    createOrEditRoute() {
      if (!this.isEditing) {
        // create route
        this.setRoute(false);
      } else {
        // edit route
        this.setRoute(true);
      }
    },
    async setRoute(isEditingRoute) {
      if (!this.validateSetRoute(isEditingRoute)) {
        return;
      }

      this.loading.setRoute = true;
      this.error.setRoute = "";
      const taskAction = "set-route";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.setRouteAborted
      );

      // register to task validation
      this.$root.$once(
        `${taskAction}-validation-ok-${eventId}`,
        this.setRouteValidationOk
      );
      this.$root.$once(
        `${taskAction}-validation-failed-${eventId}`,
        this.setRouteValidationFailed
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.setRouteCompleted
      );

      const selectedNode = this.nodes.find(
        (node) => node.value === this.selectedNodeId
      );
      const traefikInstance = selectedNode.traefikInstance;

      let setRouteData = {
        instance: this.instance,
        url: this.url,
        lets_encrypt: this.lets_encrypt,
        http2https: this.http2https,
        user_created: true,
        skip_cert_verify: this.skip_cert_verify,
      };

      if (this.host) {
        setRouteData.host = this.host;
      }

      if (this.path) {
        setRouteData.path = this.path;
        setRouteData.strip_prefix = this.strip_prefix;
      }

      const notificationTitle = isEditingRoute
        ? this.$t("settings_http_routes.edit_route_route", {
            route: this.instance,
          })
        : this.$t("settings_http_routes.create_route_route", {
            route: this.instance,
          });

      const res = await to(
        this.createModuleTaskForApp(traefikInstance, {
          action: taskAction,
          data: setRouteData,
          extra: {
            title: notificationTitle,
            description: this.$t("common.processing"),
            eventId,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.setRoute = this.getErrorMessage(err);
        this.loading.setRoute = false;
        return;
      }
    },
    setRouteAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.loading.setRoute = false;

      // hide modal so that user can see error notification
      this.$emit("hide");
    },
    setRouteValidationOk() {
      this.loading.setRoute = false;

      // hide modal after validation
      this.$emit("hide");
    },
    setRouteValidationFailed(validationErrors) {
      this.loading.setRoute = false;
      let focusAlreadySet = false;

      for (const validationError of validationErrors) {
        const param = validationError.parameter;

        // set i18n error message
        this.error[param] = this.getI18nStringWithFallback(
          "settings_http_routes." + validationError.error,
          "error." + validationError.error
        );

        if (!focusAlreadySet) {
          this.focusElement(param);
          focusAlreadySet = true;
        }
      }
    },
    setRouteCompleted() {
      this.loading.setRoute = false;
      this.clearFields();

      // reload routes
      this.$emit("reloadRoutes");
    },
    clearFields() {
      this.instance = "";
      this.url = "";
      this.host = "";
      this.path = "";
      this.lets_encrypt = false;
      this.http2https = false;
      this.strip_prefix = false;
      this.skip_cert_verify = false;
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";
</style>

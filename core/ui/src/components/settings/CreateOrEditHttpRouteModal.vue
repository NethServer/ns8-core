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
          v-if="user_created"
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
          v-if="user_created"
          v-model.trim="url"
          :placeholder="$t('settings_http_routes.url_placeholder')"
          :label="$t('settings_http_routes.url')"
          :helper-text="$t('settings_http_routes.url_helper')"
          :invalid-message="error.url"
          :disabled="loading.setRoute || !user_created"
          ref="url"
        >
          <template slot="tooltip">
            <span>{{ $t("settings_http_routes.url_tooltip") }}</span>
          </template>
        </NsTextInput>
        <NsToggle
          v-if="user_created"
          class="toggle-dependent"
          :label="$t('settings_http_routes.skip_cert_verify')"
          value="stripPrefixValue"
          :form-item="true"
          v-model="skip_cert_verify"
          :disabled="
            loading.setRoute || !url.startsWith('https://') || !user_created
          "
          ref="skip_cert_verify"
        >
          <template slot="tooltip">
            <span>{{ $t("settings_http_routes.skipCertVerify_tooltip") }}</span>
          </template>
          <template slot="text-left">{{ $t("common.disabled") }}</template>
          <template slot="text-right">{{ $t("common.enabled") }}</template>
        </NsToggle>
        <NsTextInput
          v-if="user_created"
          v-model.trim="host"
          :placeholder="$t('settings_http_routes.host_placeholder')"
          :label="
            $t('settings_http_routes.host') + ' (' + $t('common.optional') + ')'
          "
          :helper-text="$t('settings_http_routes.host_helper')"
          :invalid-message="error.host"
          :disabled="loading.setRoute || !user_created"
          ref="host"
        >
          <template slot="tooltip">
            <span>{{ $t("settings_http_routes.host_tooltip") }}</span>
          </template>
        </NsTextInput>
        <NsToggle
          v-if="user_created"
          :label="$t('settings_http_routes.request_lets_encrypt_certificate')"
          value="letsEncryptValue"
          class="toggle-dependent"
          :form-item="true"
          v-model="lets_encrypt"
          :disabled="loading.setRoute || !user_created || !host"
          ref="lets_encrypt"
        >
          <template slot="text-left">{{ $t("common.disabled") }}</template>
          <template slot="text-right">{{ $t("common.enabled") }}</template>
        </NsToggle>
        <NsTextInput
          v-if="user_created"
          v-model.trim="path"
          :placeholder="$t('settings_http_routes.path_placeholder')"
          :label="
            $t('settings_http_routes.path') + ' (' + $t('common.optional') + ')'
          "
          :invalid-message="error.path"
          :disabled="loading.setRoute || !user_created"
          ref="path"
        >
          <template slot="tooltip">
            <span>{{ $t("settings_http_routes.path_tooltip") }}</span>
          </template>
        </NsTextInput>
        <NsToggle
          v-if="user_created"
          :label="$t('settings_http_routes.strip_prefix')"
          value="stripPrefixValue"
          class="toggle-dependent"
          :form-item="true"
          v-model="strip_prefix"
          :disabled="loading.setRoute || !user_created || !path"
          ref="strip_prefix"
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
          v-if="user_created"
          :label="$t('settings_http_routes.slash_redirect')"
          value="slashRedirectValue"
          class="toggle-dependent"
          :form-item="true"
          v-model="slash_redirect"
          :disabled="loading.setRoute || !user_created || !path"
          ref="slash_redirect"
        >
          <template slot="tooltip">
            {{ $t("settings_http_routes.slash_redirect_tooltip") }}
          </template>
          <template slot="text-left">{{ $t("common.disabled") }}</template>
          <template slot="text-right">{{ $t("common.enabled") }}</template>
        </NsToggle>
        <!-- advanced options -->
        <template v-if="user_created">
          <cv-accordion ref="accordion" class="maxwidth mg-bottom">
            <cv-accordion-item
              :open="isAccordionOpen"
              @click="toggleAccordionStatus"
            >
              <template slot="title">{{ $t("common.advanced") }}</template>
              <template slot="content">
                <NsToggle
                  v-if="user_created"
                  :label="$t('settings_http_routes.http2https')"
                  value="http2httpsValue"
                  :form-item="true"
                  v-model="http2https"
                  :disabled="loading.setRoute || !user_created"
                  ref="http2https"
                >
                  <template slot="text-left">{{
                    $t("common.disabled")
                  }}</template>
                  <template slot="text-right">{{
                    $t("common.enabled")
                  }}</template>
                </NsToggle>
                <cv-text-area
                  v-model.trim="ip_allowlist_str"
                  :placeholder="
                    $t('settings_http_routes.ip_allowlist_placeholder')
                  "
                  :label="$t('settings_http_routes.ip_allowlist')"
                  :invalid-message="error.ip_allowlist"
                  :disabled="loading.setRoute"
                  ref="ip_allowlist"
                  :helper-text="$t('settings_http_routes.ip_allowlist_helper')"
              /></template>
            </cv-accordion-item>
          </cv-accordion>
        </template>
        <!-- if user_created is false, show ip_allowlist as text area -->
        <template v-else>
          <cv-text-area
            v-model.trim="ip_allowlist_str"
            :placeholder="$t('settings_http_routes.ip_allowlist_placeholder')"
            :label="$t('settings_http_routes.ip_allowlist')"
            :invalid-message="error.ip_allowlist"
            :disabled="loading.setRoute"
            ref="ip_allowlist"
            :helper-text="$t('settings_http_routes.ip_allowlist_helper')"
        /></template>
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
        ? $t("settings_http_routes.save_route")
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
      http2https: true,
      strip_prefix: false,
      slash_redirect: false,
      ip_allowlist: "",
      ip_allowlist_str: "",
      user_created: true,
      isAccordionOpen: false,
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
        ip_allowlist: "",
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
    host(newHost) {
      if (!newHost) {
        this.lets_encrypt = false; // disable lets_encrypt by default when host is not set
      }
    },
    path(newPath) {
      if (!newPath) {
        this.strip_prefix = false; // disable strip_prefix by default when path is not set
        this.slash_redirect = false; // disable slash_redirect by default when path is not set
      }
    },
    isShown: function () {
      if (this.isShown) {
        this.clearErrors();
        this.isAccordionOpen = false;

        if (this.isEditing) {
          // edit route
          this.instance = this.route.instance;
          this.selectedNodeId = this.route.nodeId;
          this.url = this.route.url;
          this.host = this.route.host;
          this.path = this.route.path;
          this.strip_prefix = this.route.strip_prefix;
          this.slash_redirect = this.route.slash_redirect;
          this.lets_encrypt = this.route.lets_encrypt;
          this.http2https = this.route.http2https;
          this.skip_cert_verify = this.route.skip_cert_verify;
          this.ip_allowlist_str = this.route.ip_allowlist_str;
          this.user_created = this.route.user_created;
          // open accordion if http2https or ip_allowlist_str is set
          if (this.http2https || this.ip_allowlist_str) {
            this.$nextTick(() => {
              this.isAccordionOpen = true;
            });
          }
        } else {
          // create route
          this.clearFields();
          this.isAccordionOpen = false;
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
    toggleAccordionStatus() {
      this.isAccordionOpen = !this.isAccordionOpen;
    },
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

      // IP allow list
      if (this.ip_allowlist_str) {
        const ipList = this.ip_allowlist_str.split("\n").map((ip) => ip.trim());
        // IPv4 and IPv4 CIDR pattern
        const ipv4Pattern =
          /^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\/([0-9]|[1-2][0-9]|3[0-2]))?$/;

        for (const ip of ipList) {
          if (!ipv4Pattern.test(ip)) {
            this.error.ip_allowlist = this.$t(
              "settings_http_routes.invalid_ipv4",
              { ip: ip }
            );
            if (isValidationOk) {
              this.focusElement("ip_allowlist");
              isValidationOk = false;
              break;
            }
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
        lets_encrypt_check: true,
        lets_encrypt_cleanup: true,
        http2https: this.http2https,
        user_created: this.user_created,
        skip_cert_verify: this.skip_cert_verify,
        ip_allowlist: this.ip_allowlist_str
          ? this.ip_allowlist_str.split("\n").map((ip) => ip.trim())
          : [],
      };

      if (this.host) {
        setRouteData.host = this.host;
      } else {
        setRouteData.host = null;
      }

      if (this.path) {
        setRouteData.path = this.path;
        setRouteData.strip_prefix = this.strip_prefix;
        setRouteData.slash_redirect = this.slash_redirect;
      } else {
        setRouteData.path = null;
        setRouteData.strip_prefix = false;
        setRouteData.slash_redirect = false;
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
      this.selectedNodeId = "";
      this.url = "";
      this.host = "";
      this.path = "";
      this.lets_encrypt = false;
      this.http2https = true;
      this.strip_prefix = false;
      this.slash_redirect = false;
      this.skip_cert_verify = false;
      this.ip_allowlist_str = "";
      this.user_created = true;
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";
.maxwidth {
  max-width: 38rem;
}
.toggle-dependent {
  margin-bottom: $spacing-07;
}
</style>

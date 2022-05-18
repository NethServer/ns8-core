<template>
  <NsModal
    size="default"
    :visible="isShown"
    :primary-button-disabled="loading.setRoute"
    :isLoading="loading.setRoute"
    @modal-hidden="onModalHidden"
    @primary-click="setRoute"
  >
    <template slot="title">{{
      $t("settings_http_routes.create_route")
    }}</template>
    <template slot="content">
      <cv-form @submit.prevent="setRoute">
        <NsInlineNotification
          v-if="error.setRoute"
          kind="error"
          :title="$t('action.set-route')"
          :description="error.setRoute"
          :showCloseButton="false"
        />
        <NsTextInput
          v-model.trim="instance"
          :label="$t('settings_http_routes.name')"
          :invalid-message="error.instance"
          :disabled="loading.setRoute"
          data-modal-primary-focus
          ref="instance"
        />
        <cv-combo-box
          v-model="selectedNodeId"
          :label="$t('common.choose')"
          :title="$t('common.node')"
          :auto-filter="true"
          :auto-highlight="true"
          :options="nodes"
          :disabled="loading.setRoute"
          :invalid-message="error.node"
          ref="node"
        >
        </cv-combo-box>
        <NsTextInput
          v-model.trim="url"
          :label="$t('settings_http_routes.url')"
          :helper-text="$t('settings_http_routes.url_helper')"
          :invalid-message="error.url"
          :disabled="loading.setRoute"
          ref="url"
        >
          <template slot="tooltip">
            <span v-html="$t('settings_http_routes.url_tooltip')"></span>
          </template>
        </NsTextInput>
        <NsTextInput
          v-model.trim="host"
          :label="
            $t('settings_http_routes.host') + ' (' + $t('common.optional') + ')'
          "
          :helper-text="$t('settings_http_routes.host_helper')"
          :invalid-message="error.host"
          :disabled="loading.setRoute"
          ref="host"
        >
          <template slot="tooltip">
            <span v-html="$t('settings_http_routes.host_tooltip')"></span>
          </template>
        </NsTextInput>
        <NsTextInput
          v-model.trim="path"
          :label="
            $t('settings_http_routes.path') + ' (' + $t('common.optional') + ')'
          "
          :invalid-message="error.path"
          :disabled="loading.setRoute"
          ref="path"
        >
          <template slot="tooltip">
            <span v-html="$t('settings_http_routes.path_tooltip')"></span>
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
            <span
              v-html="$t('settings_http_routes.strip_prefix_tooltip')"
            ></span>
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
      </cv-form>
    </template>
    <template slot="secondary-button">{{ $t("common.cancel") }}</template>
    <template slot="primary-button">{{
      $t("settings_http_routes.create_route")
    }}</template>
  </NsModal>
</template>

<script>
import to from "await-to-js";
import { UtilService, TaskService } from "@nethserver/ns8-ui-lib";

export default {
  name: "CreateHttpRouteModal",
  mixins: [UtilService, TaskService],
  props: {
    isShown: Boolean,
    nodes: {
      type: Array,
      required: true,
    },
  },
  data() {
    return {
      instance: "",
      selectedNodeId: "",
      url: "",
      host: "",
      path: "",
      lets_encrypt: false,
      http2https: false,
      strip_prefix: false,
      loading: {
        setRoute: false,
      },
      error: {
        setRoute: "",
        instance: "",
        node: "",
        url: "",
        host: "",
        path: "",
      },
    };
  },
  methods: {
    onModalHidden() {
      this.clearErrors();
      this.$emit("hide");
    },
    validateSetRoute() {
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
      }
      return isValidationOk;
    },
    async setRoute() {
      if (!this.validateSetRoute()) {
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

      const res = await to(
        this.createModuleTaskForApp(traefikInstance, {
          action: taskAction,
          data: {
            instance: this.instance,
            url: this.url,
            host: this.host,
            path: this.path,
            strip_prefix: this.strip_prefix,
            lets_encrypt: this.lets_encrypt,
            http2https: this.http2https,
          },
          extra: {
            title: this.$t("settings_http_routes.create_route_route", {
              route: this.instance,
            }),
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
        this.error[param] = this.$t(
          "settings_http_routes." + validationError.error
        );

        if (!focusAlreadySet) {
          this.focusElement(param);
          focusAlreadySet = true;
        }
      }
    },
    setRouteCompleted() {
      this.loading.setRoute = false;

      // reload users
      this.$emit("reloadRoutes");
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";

.cv-form .bx--form-item {
  margin-bottom: $spacing-06;
}
</style>

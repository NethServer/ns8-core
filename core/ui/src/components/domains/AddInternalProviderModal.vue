<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <NsWizard
    size="default"
    :visible="isShown"
    :cancelLabel="$t('common.cancel')"
    :previousLabel="$t('common.previous')"
    :nextLabel="nextButtonLabel"
    :isPreviousDisabled="true"
    :isNextDisabled="isNextButtonDisabled"
    :isCancelDisabled="isCancelButtonDisabled"
    :isNextLoading="false"
    @modal-hidden="$emit('hide')"
    @cancel="$emit('hide')"
    @nextStep="nextStep"
  >
    <template slot="title">{{
      $t("domain_configuration.add_provider")
    }}</template>
    <template slot="content">
      <template v-if="step == 'node'">
        <div>
          {{ $t("domains.choose_node_for_account_provider_installation") }}
        </div>
        <NsInlineNotification
          v-if="clusterNodes.length == disabledNodes.length"
          kind="info"
          :title="$t('domains.no_node_eligible_for_provider_installation')"
          :showCloseButton="false"
        />
        <NodeSelector
          :disabledNodes="disabledNodes"
          @selectNode="onSelectNode"
          class="mg-top-xlg"
        >
          <template v-for="node in clusterNodes">
            <template :slot="`node-${node.id}`">
              <span v-if="disabledNodes.includes(node.id)" :key="node.id">
                {{ $t("domains.provider_already_installed") }}
              </span>
            </template>
          </template>
        </NodeSelector>
      </template>
      <template v-if="step == 'installingProvider'">
        <NsInlineNotification
          v-if="error.addInternalProvider"
          kind="error"
          :title="$t('action.add-internal-provider')"
          :description="error.addInternalProvider"
          :showCloseButton="false"
        />
        <NsEmptyState
          :title="$t('domains.installing_account_provider')"
          :animationData="GearsLottie"
          animationTitle="gears"
          :loop="true"
        />
        <NsProgressBar
          :value="installProviderProgress"
          :indeterminate="!installProviderProgress"
          class="mg-bottom-md"
        />
      </template>
      <template v-if="step == 'internalConfig'">
        <!-- openldap -->
        <template v-if="isOpenLdap">
          <NsInlineNotification
            v-if="error.openldap.getDefaults"
            kind="error"
            :title="$t('action.get-defaults')"
            :description="error.openldap.getDefaults"
            :showCloseButton="false"
          />
          <cv-form>
            <NsTextInput
              :label="$t('openldap.admuser')"
              v-model.trim="openldap.admuser"
              :helper-text="$t('openldap.enter_openldap_admin_username')"
              :invalid-message="$t(error.openldap.admuser)"
              :disabled="
                loading.openldap.configureModule || loading.openldap.getDefaults
              "
              ref="admuser"
            >
            </NsTextInput>
            <NsTextInput
              :label="$t('openldap.admpass')"
              type="password"
              v-model="openldap.admpass"
              :helper-text="$t('openldap.enter_openldap_admin_password')"
              :invalid-message="$t(error.openldap.admpass)"
              :disabled="
                loading.openldap.configureModule || loading.openldap.getDefaults
              "
              :password-hide-label="$t('password.hide_password')"
              :password-show-label="$t('password.show_password')"
              ref="admpass"
              name="admpass"
            ></NsTextInput>
          </cv-form>
          <NsInlineNotification
            v-if="error.openldap.configureModule"
            kind="error"
            :title="$t('action.configure-module')"
            :description="error.openldap.configureModule"
            :showCloseButton="false"
          />
        </template>
        <!-- samba -->
        <template v-if="isSamba">
          <NsInlineNotification
            v-if="error.samba.getDefaults"
            kind="error"
            :title="$t('action.get-defaults')"
            :description="error.samba.getDefaults"
            :showCloseButton="false"
          />
          <cv-form>
            <NsTextInput
              :label="$t('samba.adminuser')"
              v-model.trim="samba.adminuser"
              :helper-text="$t('samba.enter_samba_admin_username')"
              :invalid-message="$t(error.samba.adminuser)"
              :disabled="
                loading.samba.configureModule || loading.samba.getDefaults
              "
              ref="adminuser"
            >
            </NsTextInput>
            <NsTextInput
              :label="$t('samba.adminpass')"
              type="password"
              v-model="samba.adminpass"
              :helper-text="$t('samba.enter_samba_admin_password')"
              :invalid-message="$t(error.samba.adminpass)"
              :disabled="
                loading.samba.configureModule || loading.samba.getDefaults
              "
              :password-hide-label="$t('password.hide_password')"
              :password-show-label="$t('password.show_password')"
              ref="adminpass"
              name="adminpass"
            ></NsTextInput>
            <NsTextInput
              :label="$t('samba.hostname')"
              v-model.trim="samba.hostname"
              :invalid-message="$t(error.samba.hostname)"
              :disabled="
                loading.samba.configureModule || loading.samba.getDefaults
              "
              ref="hostname"
            >
            </NsTextInput>
            <NsInlineNotification
              v-if="!samba.canInstallFileServer"
              kind="info"
              :description="$t('samba.enable_file_server_disabled_info')"
              :showCloseButton="false"
            />
            <NsToggle
              :label="$t('samba.enable_file_server_label')"
              value="enableFileServer"
              :form-item="true"
              v-model="samba.enableFileServer"
              :disabled="
                !samba.canInstallFileServer ||
                loading.samba.configureModule ||
                loading.samba.getDefaults
              "
              tooltipAlignment="start"
              tooltipDirection="bottom"
              ref="enableFileServer"
            >
              <template slot="tooltip">{{
                $t("samba.enable_file_server_tooltip")
              }}</template>
              <template slot="text-left">{{ $t("common.disabled") }}</template>
              <template slot="text-right">{{ $t("common.enabled") }}</template>
            </NsToggle>
            <NsInlineNotification
              v-if="samba.enableFileServer"
              kind="warning"
              :title="$t('common.warning')"
              :description="$t('samba.choose_ip_address_warning')"
              :showCloseButton="false"
            />
            <NsComboBox
              v-if="samba.enableFileServer"
              v-model="samba.ipaddress"
              :label="$t('common.choose')"
              :options="filteredIpAddressOptions"
              auto-highlight
              :title="$t('samba.provider_ipaddress')"
              :invalid-message="$t(error.samba.ipaddress)"
              :disabled="
                loading.samba.configureModule || loading.samba.getDefaults
              "
              light
              class="mg-bottom-5"
              ref="ipaddress"
            >
            </NsComboBox>
            <NsInlineNotification
              v-if="!samba.enableFileServer && vpnInterface"
              kind="info"
              :description="
                $t('samba.provider_on_vpn_address_message', {
                  vpnIpAddress: vpnInterface.value,
                })
              "
              :showCloseButton="false"
            />
          </cv-form>
          <NsInlineNotification
            v-if="error.samba.configureModule"
            kind="error"
            :title="$t('action.configure-module')"
            :description="error.samba.configureModule"
            :showCloseButton="false"
          />
          <NsInlineNotification
            v-if="error.samba.cannotReachDc"
            kind="error"
            :title="$t('domains.cannot_connect_to_domain_controller')"
            :description="
              $t('domains.cannot_connect_to_domain_controller_description')
            "
            :showCloseButton="false"
          />
        </template>
      </template>
      <template v-if="step == 'configuringProvider'">
        <NsEmptyState
          :title="$t('domains.configuring_account_provider')"
          :animationData="GearsLottie"
          animationTitle="gears"
          :loop="true"
        />
        <NsProgressBar
          :value="configureProviderProgress"
          :indeterminate="!configureProviderProgress"
          class="mg-bottom-md"
        />
      </template>
    </template>
  </NsWizard>
</template>

<script>
import {
  UtilService,
  TaskService,
  IconService,
  LottieService,
} from "@nethserver/ns8-ui-lib";
import to from "await-to-js";
import NodeSelector from "@/components/nodes/NodeSelector";
import { mapState } from "vuex";

export default {
  name: "AddInternalProviderModal",
  components: { NodeSelector },
  mixins: [UtilService, TaskService, IconService, LottieService],
  props: {
    isShown: {
      type: Boolean,
      default: true,
    },
    domain: {
      type: Object,
      required: true,
    },
    isResumeConfiguration: {
      type: Boolean,
      default: false,
    },
    providerId: {
      type: String,
      default: "",
    },
    domains: {
      type: Array,
      default: () => [],
    },
  },
  data() {
    return {
      step: "node",
      installProviderProgress: 0,
      configureProviderProgress: 0,
      selectedNode: null,
      openldap: {
        admuser: "",
        admpass: "",
      },
      samba: {
        adminuser: "",
        adminpass: "",
        ipaddress: "",
        ipAddressOptions: [],
        hostname: "",
        canInstallFileServer: false,
        enableFileServer: false,
      },
      loading: {
        openldap: {
          getDefaults: false,
          configureModule: false,
        },
        samba: {
          getDefaults: false,
          configureModule: false,
        },
      },
      error: {
        addInternalProvider: "",
        openldap: {
          admuser: "",
          admpass: "",
          configureModule: "",
          getDefaults: "",
        },
        samba: {
          adminuser: "",
          adminpass: "",
          ipaddress: "",
          hostname: "",
          configureModule: "",
          getDefaults: "",
          cannotReachDc: false,
        },
      },
    };
  },
  computed: {
    ...mapState(["clusterNodes"]),
    isOpenLdap() {
      return this.domain.schema == "rfc2307";
    },
    isSamba() {
      return this.domain.schema == "ad";
    },
    nextButtonLabel() {
      if (this.step == "node") {
        return this.$t("domains.install_provider");
      } else if (this.step == "internalConfig") {
        return this.$t("domains.configure_provider");
      } else {
        return this.$t("common.next");
      }
    },
    isNextButtonDisabled() {
      return (
        this.step == "installingProvider" ||
        this.step == "configuringProvider" ||
        (this.step == "node" && !this.selectedNode) ||
        this.loading.samba.configureModule ||
        this.loading.openldap.configureModule
      );
    },
    isCancelButtonDisabled() {
      return (
        this.step == "installingProvider" || this.step == "configuringProvider"
      );
    },
    vpnInterface() {
      return this.samba.ipAddressOptions.find((option) => {
        return option.iface === "wg0";
      });
    },
    filteredIpAddressOptions() {
      if (this.samba.enableFileServer) {
        // remove VPN IP address from options
        return this.samba.ipAddressOptions.filter((option) => {
          return option != this.vpnInterface;
        });
      } else {
        return this.samba.ipAddressOptions;
      }
    },
    disabledNodes() {
      if (this.isOpenLdap) {
        // openldap supports multiple instances on the same node
        return [];
      } else {
        // samba allows only one instance per node
        const disabledNodes = [];

        for (const domain of this.domains) {
          for (const provider of domain.providers) {
            if (provider.id.startsWith("samba")) {
              disabledNodes.push(provider.node);
            }
          }
        }
        // remove possible duplicates
        return [...new Set(disabledNodes)];
      }
    },
  },
  watch: {
    isShown: function () {
      if (this.isShown) {
        this.clearErrors();

        if (!this.isResumeConfiguration) {
          // start wizard from first step
          this.step = "node";
        } else {
          // resume configuration
          this.step = "internalConfig";

          if (this.isOpenLdap) {
            this.getOpenLdapDefaults();
          } else if (this.isSamba) {
            this.getSambaDefaults();
          }
        }
      }
    },
    providerId: function () {
      this.newProviderId = this.providerId;
    },
  },
  created() {
    this.newProviderId = this.providerId;
  },
  methods: {
    nextStep() {
      switch (this.step) {
        case "node":
          this.step = "installingProvider";
          this.installProvider();
          break;
        case "internalConfig":
          if (this.isSamba) {
            this.configureSambaModule();
          } else if (this.isOpenLdap) {
            this.configureOpenLdapModule();
          }
          break;
      }
    },
    onSelectNode(selectedNode) {
      this.selectedNode = selectedNode;
    },
    async installProvider() {
      this.error.addInternalProvider = "";
      const taskAction = "add-internal-provider";
      const eventId = this.getUuid();
      this.installProviderProgress = 0;

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.addInternalProviderCompleted
      );

      // register to task progress to update progress bar
      this.$root.$on(
        `${taskAction}-progress-${eventId}`,
        this.addInternalProviderProgress
      );

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.addInternalProviderAborted
      );

      const providerImage = this.isOpenLdap ? "openldap" : "samba";

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            image: providerImage,
            node: parseInt(this.selectedNode.id),
            domain: this.domain.name,
          },
          extra: {
            title: this.$t("action." + taskAction),
            node: this.selectedNode.id,
            isNotificationHidden: true,
            isProgressNotified: true,
            eventId,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.addInternalProvider = this.getErrorMessage(err);
        return;
      }
    },
    addInternalProviderAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);

      // unregister to task progress
      this.$root.$off(
        `${taskContext.action}-progress-${taskContext.extra.eventId}`
      );

      // hide modal so that user can see error notification
      this.$emit("hide");
    },
    addInternalProviderCompleted(taskContext, taskResult) {
      // unregister to task progress
      this.$root.$off(
        `${taskContext.action}-progress-${taskContext.extra.eventId}`
      );

      this.step = "internalConfig";
      this.newProviderId = taskResult.output.module_id;

      if (this.isSamba) {
        this.getSambaDefaults();
      } else {
        this.getOpenLdapDefaults();
      }

      // reload domains
      this.$emit("reloadDomains");

      // show new app in app drawer
      this.$root.$emit("reloadAppDrawer");
    },
    addInternalProviderProgress(progress) {
      this.installProviderProgress = progress;
    },
    async getOpenLdapDefaults() {
      this.loading.openldap.getDefaults = true;
      this.error.openldap.getDefaults = "";
      const taskAction = "get-defaults";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.getOpenLdapDefaultsAborted
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.getOpenLdapDefaultsCompleted
      );

      const res = await to(
        this.createModuleTaskForApp(this.newProviderId, {
          action: taskAction,
          data: {
            provision: "join-domain",
            domain: this.domain.name,
          },
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
        this.error.openldap.getDefaults = this.getErrorMessage(err);
        return;
      }
    },
    getOpenLdapDefaultsAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.loading.openldap.getDefaults = false;

      // hide modal so that user can see error notification
      this.$emit("hide");
    },
    getOpenLdapDefaultsCompleted(taskContext, taskResult) {
      this.loading.openldap.getDefaults = false;
      const defaults = taskResult.output;
      this.openldap.admuser = defaults.admuser;

      // clear password
      this.openldap.admpass = "";

      // focus on first field
      this.$nextTick(() => {
        this.focusElement("admuser");
      });
    },
    async getSambaDefaults() {
      this.loading.samba.getDefaults = true;
      this.error.samba.getDefaults = "";
      const taskAction = "get-defaults";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.getSambaDefaultsAborted
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.getSambaDefaultsCompleted
      );

      const res = await to(
        this.createModuleTaskForApp(this.newProviderId, {
          action: taskAction,
          data: {
            provision: "join-domain",
            realm: this.domain.name,
          },
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
        this.error.samba.getDefaults = this.getErrorMessage(err);
        return;
      }
    },
    getSambaDefaultsAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.loading.samba.getDefaults = false;

      // hide modal so that user can see error notification
      this.$emit("hide");
    },
    getSambaDefaultsCompleted(taskContext, taskResult) {
      this.loading.samba.getDefaults = false;
      const defaults = taskResult.output;

      this.samba.adminuser = defaults.adminuser;
      this.samba.adminpass = "";
      this.samba.hostname = defaults.hostname;
      this.samba.canInstallFileServer = defaults.can_install_file_server;

      // ip address combo box
      let index = 0;
      for (const ipObject of defaults.ipaddress_list) {
        const option = {
          name: ipObject.ipaddress.replace(/\W/g, "_"),
          label: ipObject.ipaddress + " - " + ipObject.label,
          value: ipObject.ipaddress,
          iface: ipObject.label,
        };
        this.$set(this.samba.ipAddressOptions, index, option);
        index++;
      }

      // focus first field
      this.$nextTick(() => {
        this.focusElement("adminuser");
      });
    },
    clearOpenLdapErrors() {
      for (const key of Object.keys(this.error.openldap)) {
        this.error.openldap[key] = "";
      }
    },
    clearSambaErrors() {
      for (const key of Object.keys(this.error.samba)) {
        this.error.samba[key] = "";
      }
    },
    validateConfigureOpenLdapModule() {
      this.clearOpenLdapErrors();
      let isValidationOk = true;

      // openldap admin user

      if (!this.openldap.admuser) {
        this.error.openldap.admuser = "common.required";

        if (isValidationOk) {
          this.focusElement("admuser");
          isValidationOk = false;
        }
      }

      // openldap admin password

      if (!this.openldap.admpass) {
        this.error.openldap.admpass = "common.required";

        if (isValidationOk) {
          this.focusElement("admpass");
          isValidationOk = false;
        }
      }
      return isValidationOk;
    },
    validateConfigureSambaModule() {
      this.clearSambaErrors();
      let isValidationOk = true;

      // samba admin user

      if (!this.samba.adminuser) {
        this.error.samba.adminuser = "common.required";

        if (isValidationOk) {
          this.focusElement("adminuser");
          isValidationOk = false;
        }
      }

      // samba admin password

      if (!this.samba.adminpass) {
        this.error.samba.adminpass = "common.required";

        if (isValidationOk) {
          this.focusElement("adminpass");
          isValidationOk = false;
        }
      }

      // samba hostname

      if (!this.samba.hostname) {
        this.error.samba.hostname = "common.required";

        if (isValidationOk) {
          this.focusElement("hostname");
          isValidationOk = false;
        }
      }

      // samba ip address

      if (this.samba.enableFileServer && !this.samba.ipaddress) {
        this.error.samba.ipaddress = "common.required";

        if (isValidationOk) {
          isValidationOk = false;
        }
      }
      return isValidationOk;
    },
    async configureSambaModule() {
      const isValidationOk = this.validateConfigureSambaModule();
      if (!isValidationOk) {
        return;
      }

      this.loading.samba.configureModule = true;
      const taskAction = "configure-module";
      const eventId = this.getUuid();
      this.configureProviderProgress = 0;

      if (!this.samba.enableFileServer) {
        // set VPN address as IP address
        this.samba.ipaddress = this.vpnInterface.value;
      }

      // register to task validation
      this.$root.$once(
        `${taskAction}-validation-failed-${eventId}`,
        this.configureSambaModuleValidationFailed
      );
      this.$root.$once(
        `${taskAction}-validation-ok-${eventId}`,
        this.configureSambaModuleValidationOk
      );

      // register to task progress to update progress bar
      this.$root.$on(
        `${taskAction}-progress-${eventId}`,
        this.configureSambaModuleProgress
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.configureSambaModuleCompleted
      );

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.configureSambaModuleAborted
      );

      const res = await to(
        this.createModuleTaskForApp(this.newProviderId, {
          action: taskAction,
          data: {
            adminuser: this.samba.adminuser,
            adminpass: this.samba.adminpass,
            realm: this.domain.name,
            ipaddress: this.samba.ipaddress,
            hostname: this.samba.hostname,
            provision: "join-domain",
          },
          extra: {
            title: this.$t("samba.samba_configuration"),
            isNotificationHidden: true,
            isProgressNotified: true,
            eventId,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.samba.configureModule = this.getErrorMessage(err);
        this.loading.samba.configureModule = false;
        return;
      }
    },
    configureSambaModuleValidationOk() {
      this.step = "configuringProvider";
    },
    configureSambaModuleValidationFailed(validationErrors, taskContext) {
      this.loading.samba.configureModule = false;

      // unregister to task progress
      this.$root.$off(
        `${taskContext.action}-progress-${taskContext.extra.eventId}`
      );

      let focusAlreadySet = false;

      for (const validationError of validationErrors) {
        if (
          [
            "realm_dc_reachable_check_failed",
            "realm_dc_avail_check_failed",
          ].includes(validationError.error)
        ) {
          this.error.samba.cannotReachDc = true;
        } else {
          const param = validationError.parameter;
          // set i18n error message
          this.error.samba[param] = "domains." + validationError.error;

          if (!focusAlreadySet && param != "ipaddress") {
            this.focusElement(param);
            focusAlreadySet = true;
          }
        }
      }
    },
    configureSambaModuleProgress(progress) {
      this.configureProviderProgress = progress;
    },
    configureSambaModuleCompleted(taskContext) {
      this.loading.samba.configureModule = false;

      // unregister to task progress
      this.$root.$off(
        `${taskContext.action}-progress-${taskContext.extra.eventId}`
      );

      // hide modal
      this.$emit("hide");

      // reload domains
      this.$emit("reloadDomains");

      // reload app drawer (samba file server might have appeared)
      this.$root.$emit("reloadAppDrawer");
    },
    configureSambaModuleAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.loading.samba.configureModule = false;

      // unregister to task progress
      this.$root.$off(
        `${taskContext.action}-progress-${taskContext.extra.eventId}`
      );

      // hide modal so that user can see error notification
      this.$emit("hide");
    },
    async configureOpenLdapModule() {
      const isValidationOk = this.validateConfigureOpenLdapModule();
      if (!isValidationOk) {
        return;
      }

      this.loading.openldap.configureModule = true;
      const taskAction = "configure-module";
      const eventId = this.getUuid();
      this.configureProviderProgress = 0;

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.configureOpenLdapModuleAborted
      );

      // register to task validation
      this.$root.$once(
        `${taskAction}-validation-failed-${eventId}`,
        this.configureOpenLdapModuleValidationFailed
      );
      this.$root.$once(
        `${taskAction}-validation-ok-${eventId}`,
        this.configureOpenLdapModuleValidationOk
      );

      // register to task progress to update progress bar
      this.$root.$on(
        `${taskAction}-progress-${eventId}`,
        this.configureOpenLdapModuleProgress
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.configureOpenLdapModuleCompleted
      );

      const res = await to(
        this.createModuleTaskForApp(this.newProviderId, {
          action: taskAction,
          data: {
            admuser: this.openldap.admuser,
            admpass: this.openldap.admpass,
            domain: this.domain.name,
            provision: "join-domain",
          },
          extra: {
            title: this.$t("openldap.openldap_configuration"),
            isNotificationHidden: true,
            isProgressNotified: true,
            eventId,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.openldap.configureModule = this.getErrorMessage(err);
        this.loading.openldap.configureModule = false;
        return;
      }
    },
    configureOpenLdapModuleAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.loading.openldap.configureModule = false;

      // unregister to task progress
      this.$root.$off(
        `${taskContext.action}-progress-${taskContext.extra.eventId}`
      );

      // hide modal so that user can see error notification
      this.$emit("hide");
    },
    configureOpenLdapModuleValidationOk() {
      this.step = "configuringProvider";
    },
    configureOpenLdapModuleValidationFailed(validationErrors, taskContext) {
      this.loading.openldap.configureModule = false;

      // unregister to task progress
      this.$root.$off(
        `${taskContext.action}-progress-${taskContext.extra.eventId}`
      );

      let focusAlreadySet = false;

      for (const validationError of validationErrors) {
        const param = validationError.parameter;
        // set i18n error message
        this.error.openldap[param] = "domains." + validationError.error;

        if (!focusAlreadySet) {
          this.focusElement(param);
          focusAlreadySet = true;
        }
      }
    },
    configureOpenLdapModuleProgress(progress) {
      this.configureProviderProgress = progress;
    },
    configureOpenLdapModuleCompleted(taskContext) {
      this.loading.openldap.configureModule = false;

      // unregister to task progress
      this.$root.$off(
        `${taskContext.action}-progress-${taskContext.extra.eventId}`
      );

      // hide modal
      this.$emit("hide");

      // reload domains
      this.$emit("reloadDomains");
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";
.mg-top {
  margin-top: $spacing-05;
}

.row {
  margin-bottom: $spacing-03;
}

.label {
  padding-right: 0.5rem;
  font-weight: bold;
}

.bx--form {
  .bx--form-item,
  .cv-combo-box {
    margin-bottom: $spacing-06;
  }
}

.mg-bottom-5 {
  margin-bottom: 5rem !important;
}
</style>

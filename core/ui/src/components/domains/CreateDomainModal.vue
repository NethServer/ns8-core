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
    :isPreviousDisabled="isPreviousButtonDisabled"
    :isNextDisabled="isNextButtonDisabled"
    :isCancelDisabled="isCancelButtonDisabled"
    :isNextLoading="loading.external.addExternalDomain"
    @modal-hidden="$emit('hide')"
    @cancel="$emit('hide')"
    @previousStep="previousStep"
    @nextStep="nextStep"
  >
    <template slot="title">{{ $t("domains.create_domain") }}</template>
    <template slot="content">
      <template v-if="step == 'location'">
        <div class="mg-bottom-md">
          {{ $t("domains.select_domain_location") }}
        </div>
        <div class="bx--grid no-padding">
          <div class="bx--row">
            <div class="bx--col-md-4">
              <NsTile
                :light="true"
                kind="selectable"
                v-model="isInternalSelected"
                value="locationValue"
                @click="isExternalSelected = false"
                class="min-height-card"
              >
                <h6 class="mg-bottom-md">
                  {{ $t("domains.internal") }}
                </h6>
                <div>
                  {{ $t("domains.internal_description") }}
                </div>
              </NsTile>
            </div>
            <div class="bx--col-md-4">
              <NsTile
                :light="true"
                kind="selectable"
                v-model="isExternalSelected"
                value="locationValue"
                @click="isInternalSelected = false"
                class="min-height-card"
              >
                <h6 class="mg-bottom-md">
                  {{ $t("domains.external") }}
                </h6>
                <div>
                  {{ $t("domains.external_description") }}
                </div>
              </NsTile>
            </div>
          </div>
        </div>
      </template>
      <template v-if="step == 'instance'">
        <div class="mg-bottom-md">
          {{ $t("domains.select_prodiver_module") }}
        </div>
        <div class="bx--grid no-padding">
          <div class="bx--row">
            <div class="bx--col-md-4">
              <NsTile
                :light="true"
                kind="selectable"
                v-model="isOpenLdapSelected"
                value="instanceValue"
                @click="isSambaSelected = false"
                class="min-height-card"
              >
                <h6 class="mg-bottom-md">
                  {{ $t("domains.openldap") }}
                </h6>
              </NsTile>
            </div>
            <div class="bx--col-md-4">
              <NsTile
                :light="true"
                kind="selectable"
                v-model="isSambaSelected"
                value="instanceValue"
                @click="isOpenLdapSelected = false"
                class="min-height-card"
              >
                <h6 class="mg-bottom-md">
                  {{ $t("domains.samba") }}
                </h6>
              </NsTile>
            </div>
          </div>
        </div>
      </template>
      <template v-if="step == 'externalConfig'">
        <cv-form>
          <cv-text-input
            :label="$t('domains.domain')"
            v-model.trim="external.domain"
            :helper-text="$t('domains.enter_external_domain_name')"
            :invalid-message="$t(error.external.domain)"
            :disabled="loading.external.addExternalDomain"
            ref="domain"
          >
          </cv-text-input>
          <cv-text-input
            :label="$t('domains.host')"
            v-model.trim="external.host"
            :invalid-message="$t(error.external.host)"
            :disabled="loading.external.addExternalDomain"
            ref="host"
          >
          </cv-text-input>
          <cv-text-input
            :label="$t('domains.port')"
            v-model.trim="external.port"
            :invalid-message="$t(error.external.port)"
            :disabled="loading.external.addExternalDomain"
            type="number"
            ref="port"
          >
          </cv-text-input>
          <cv-text-input
            :label="$t('domains.bind_dn')"
            v-model.trim="external.bind_dn"
            :invalid-message="$t(error.external.bind_dn)"
            :disabled="loading.external.addExternalDomain"
            ref="bind_dn"
          >
          </cv-text-input>
          <cv-text-input
            :label="$t('domains.bind_password')"
            v-model="external.bind_password"
            :invalid-message="$t(error.external.bind_password)"
            :disabled="loading.external.addExternalDomain"
            ref="bind_password"
          >
          </cv-text-input>
          <cv-text-input
            :label="$t('domains.base_dn') + ' (' + $t('common.optional') + ')'"
            v-model.trim="external.base_dn"
            :invalid-message="$t(error.external.base_dn)"
            :disabled="loading.external.addExternalDomain"
            ref="base_dn"
          >
          </cv-text-input>
          <cv-toggle
            :label="$t('domains.tls')"
            value="tlsValue"
            :form-item="true"
            v-model="external.tls"
            :disabled="loading.external.addExternalDomain"
            ref="tls"
          >
            <template slot="text-left">{{ $t("common.disabled") }}</template>
            <template slot="text-right">{{ $t("common.enabled") }}</template>
          </cv-toggle>
          <NsInlineNotification
            v-if="error.external.tls"
            kind="error"
            :title="$t(error.external.tls)"
            :showCloseButton="false"
          />
          <cv-toggle
            :label="$t('domains.tls_verify')"
            value="tls_verifyValue"
            :form-item="true"
            v-model="external.tls_verify"
            :disabled="loading.external.addExternalDomain"
            ref="tls_verify"
          >
            <template slot="text-left">{{ $t("common.disabled") }}</template>
            <template slot="text-right">{{ $t("common.enabled") }}</template>
          </cv-toggle>
          <NsInlineNotification
            v-if="error.external.tls_verify"
            kind="error"
            :title="$t(error.external.tls_verify)"
            :showCloseButton="false"
          />
        </cv-form>
      </template>
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
        <template v-if="isOpenLdapSelected">
          <NsInlineNotification
            v-if="error.openldap.getDefaults"
            kind="error"
            :title="$t('action.get-defaults')"
            :description="error.openldap.getDefaults"
            :showCloseButton="false"
          />
          <cv-form>
            <cv-text-input
              :label="$t('openldap.domain')"
              v-model.trim="openldap.domain"
              :invalid-message="$t(error.openldap.domain)"
              :disabled="
                loading.openldap.configureModule || loading.openldap.getDefaults
              "
              ref="domain"
            >
            </cv-text-input>
            <cv-text-input
              :label="$t('openldap.admuser')"
              v-model.trim="openldap.admuser"
              :invalid-message="$t(error.openldap.admuser)"
              :disabled="
                loading.openldap.configureModule || loading.openldap.getDefaults
              "
              ref="admuser"
            >
            </cv-text-input>
            <NsPasswordInput
              :newPasswordLabel="$t('openldap.admpass')"
              :confirmPasswordLabel="$t('openldap.admpass_confirm')"
              v-model="openldap.admpass"
              @passwordValidation="onNewOpenLdapPasswordValidation"
              :newPaswordHelperText="
                $t('openldap.choose_openldap_admin_password')
              "
              :newPasswordInvalidMessage="$t(error.openldap.admpass)"
              :confirmPasswordInvalidMessage="
                $t(error.openldap.confirmPassword)
              "
              :passwordHideLabel="$t('password.hide_password')"
              :passwordShowLabel="$t('password.show_password')"
              :lengthLabel="$t('password.long_enough')"
              :lowercaseLabel="$t('password.lowercase_letter')"
              :uppercaseLabel="$t('password.uppercase_letter')"
              :numberLabel="$t('password.number')"
              :symbolLabel="$t('password.symbol')"
              :equalLabel="$t('password.equal')"
              :focus="openldap.focusPasswordField"
              :disabled="
                loading.openldap.configureModule || loading.openldap.getDefaults
              "
              light
              class="new-provider-password"
            />
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
        <template v-if="isSambaSelected">
          <NsInlineNotification
            v-if="error.samba.getDefaults"
            kind="error"
            :title="$t('action.get-defaults')"
            :description="error.samba.getDefaults"
            :showCloseButton="false"
          />
          <cv-form>
            <cv-text-input
              :label="$t('samba.realm')"
              v-model.trim="samba.realm"
              :invalid-message="$t(error.samba.realm)"
              :disabled="
                loading.samba.configureModule || loading.samba.getDefaults
              "
              ref="realm"
            >
            </cv-text-input>
            <cv-text-input
              :label="$t('samba.nbdomain')"
              v-model.trim="samba.nbdomain"
              :invalid-message="$t(error.samba.nbdomain)"
              :disabled="
                loading.samba.configureModule || loading.samba.getDefaults
              "
              ref="nbdomain"
            >
            </cv-text-input>
            <NsTextInput
              :label="$t('samba.choose_samba_admin_username')"
              v-model.trim="samba.adminuser"
              :invalid-message="$t(error.samba.adminuser)"
              :disabled="
                loading.samba.configureModule || loading.samba.getDefaults
              "
              tooltipAlignment="start"
              tooltipDirection="bottom"
              ref="adminuser"
            >
              <template slot="tooltip">
                <div class="mg-bottom-sm">
                  {{ $t("samba.samba_admin_username_tooltip_1") }}
                </div>
                <ul class="unordered-list">
                  <li class="mg-bottom-sm">
                    {{ $t("samba.samba_admin_username_tooltip_2") }}
                  </li>
                  <li>{{ $t("samba.samba_admin_username_tooltip_3") }}</li>
                </ul>
              </template>
            </NsTextInput>
            <NsPasswordInput
              :newPasswordLabel="$t('samba.choose_samba_admin_password')"
              :confirmPasswordLabel="$t('samba.adminpass_confirm')"
              v-model="samba.adminpass"
              @passwordValidation="onNewSambaPasswordValidation"
              :newPasswordInvalidMessage="$t(error.samba.adminpass)"
              :confirmPasswordInvalidMessage="$t(error.samba.confirmPassword)"
              :passwordHideLabel="$t('password.hide_password')"
              :passwordShowLabel="$t('password.show_password')"
              :lengthLabel="$t('password.long_enough')"
              :lowercaseLabel="$t('password.lowercase_letter')"
              :uppercaseLabel="$t('password.uppercase_letter')"
              :numberLabel="$t('password.number')"
              :symbolLabel="$t('password.symbol')"
              :equalLabel="$t('password.equal')"
              :focus="samba.focusPasswordField"
              :disabled="
                loading.samba.configureModule || loading.samba.getDefaults
              "
              light
              class="new-provider-password"
            />
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
            <NsToggle
              :label="$t('samba.enable_file_server_label')"
              value="enableFileServer"
              :form-item="true"
              v-model="samba.enableFileServer"
              :disabled="
                loading.samba.configureModule || loading.samba.getDefaults
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
  name: "CreateDomainModal",
  components: { NodeSelector },
  mixins: [UtilService, TaskService, IconService, LottieService],
  props: {
    isShown: {
      type: Boolean,
      default: true,
    },
    isResumeConfiguration: {
      type: Boolean,
      default: false,
    },
    providerId: {
      type: String,
      default: "",
    },
    isOpenLdap: {
      type: Boolean,
      default: false,
    },
    isSamba: {
      type: Boolean,
      default: false,
    },
    domains: {
      type: Array,
      default: () => [],
    },
  },
  data() {
    return {
      step: "",
      isInternalSelected: false,
      isExternalSelected: false,
      newProviderId: "",
      isOpenLdapSelected: false,
      isSambaSelected: false,
      installProviderProgress: 0,
      configureProviderProgress: 0,
      selectedNode: null,
      samba: {
        adminuser: "",
        adminpass: "",
        realm: "",
        ipaddress: "",
        ipAddressOptions: [],
        hostname: "",
        nbdomain: "",
        passwordValidation: null,
        focusPasswordField: { element: "" },
        enableFileServer: false,
      },
      openldap: {
        domain: "",
        admuser: "",
        admpass: "",
        passwordValidation: null,
        focusPasswordField: { element: "" },
      },
      external: {
        domain: "",
        host: "",
        port: "",
        bind_dn: "",
        bind_password: "",
        base_dn: "",
        tls: true,
        tls_verify: true,
      },
      loading: {
        samba: {
          getDefaults: false,
          configureModule: false,
        },
        openldap: {
          getDefaults: false,
          configureModule: false,
        },
        external: {
          addExternalDomain: false,
        },
      },
      error: {
        addInternalProvider: "",
        samba: {
          adminuser: "",
          adminpass: "",
          confirmPassword: "",
          realm: "",
          ipaddress: "",
          hostname: "",
          nbdomain: "",
          configureModule: "",
          getDefaults: "",
          addExternalDomain: "",
          cannotReachDc: false,
        },
        openldap: {
          getDefaults: "",
          domain: "",
          admuser: "",
          admpass: "",
          confirmPassword: "",
        },
        external: {
          domain: "",
          host: "",
          port: "",
          bind_dn: "",
          bind_password: "",
          base_dn: "",
          tls: "",
          tls_verify: "",
        },
      },
    };
  },
  computed: {
    ...mapState(["clusterNodes"]),
    nextButtonLabel() {
      if (this.step == "node") {
        return this.$t("domains.install_provider");
      } else if (
        this.step == "internalConfig" ||
        this.step == "externalConfig"
      ) {
        return this.$t("domains.configure_domain");
      } else {
        return this.$t("common.next");
      }
    },
    isNextButtonDisabled() {
      return (
        this.loading.openldap.configureModule ||
        this.loading.samba.configureModule ||
        this.loading.external.addExternalDomain ||
        this.step == "installingProvider" ||
        this.step == "configuringProvider" ||
        (this.step == "location" &&
          !this.isInternalSelected &&
          !this.isExternalSelected) ||
        (this.step == "instance" &&
          !this.isOpenLdapSelected &&
          !this.isSambaSelected) ||
        (this.step == "node" && !this.selectedNode)
      );
    },
    isCancelButtonDisabled() {
      return (
        this.step == "installingProvider" || this.step == "configuringProvider"
      );
    },
    isPreviousButtonDisabled() {
      return (
        this.isResumeConfiguration ||
        this.loading.samba.configureModule ||
        this.loading.openldap.configureModule ||
        this.loading.external.addExternalDomain ||
        [
          "location",
          "installingProvider",
          "internalConfig",
          "configureProviderProgress",
        ].includes(this.step)
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
      if (this.isOpenLdapSelected) {
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
          this.step = "location";
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
    isOpenLdap: function () {
      this.isOpenLdapSelected = this.isOpenLdap;
    },
    isSamba: function () {
      this.isSambaSelected = this.isSamba;
    },
  },
  created() {
    this.newProviderId = this.providerId;
    this.isOpenLdapSelected = this.isOpenLdap;
    this.isSambaSelected = this.isSamba;
  },
  methods: {
    nextStep() {
      switch (this.step) {
        case "location":
          if (this.isInternalSelected) {
            // internal
            this.step = "instance";
          } else {
            // external
            this.step = "externalConfig";
          }
          break;
        case "instance":
          this.step = "node";
          break;
        case "externalConfig":
          this.addExternalDomain();
          break;
        case "node":
          this.step = "installingProvider";
          this.installProvider();
          break;
        case "internalConfig":
          if (this.isSambaSelected) {
            this.configureSambaModule();
          } else if (this.isOpenLdapSelected) {
            this.configureOpenLdapModule();
          }
          break;
      }
    },
    previousStep() {
      switch (this.step) {
        case "instance":
        case "externalConfig":
          this.step = "location";
          break;
        case "node":
          this.step = "instance";
          break;
      }
    },
    async installProvider() {
      this.error.addInternalProvider = "";
      const taskAction = "add-internal-provider";
      const eventId = this.getUuid();
      this.installProviderProgress = 0;

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.addInternalProviderAborted
      );

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

      const selectedNodeId = this.selectedNode
        ? parseInt(this.selectedNode.id)
        : 1;

      const providerImage = this.isOpenLdapSelected ? "openldap" : "samba";

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            image: providerImage,
            node: selectedNodeId,
          },
          extra: {
            title: this.$t("action." + taskAction),
            node: selectedNodeId,
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

      if (this.isSambaSelected) {
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
            provision: "new-domain",
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
      this.openldap.domain = defaults.domain;
      this.openldap.admuser = defaults.admuser;

      // clear password
      this.openldap.admpass = "";

      // focus on first field
      this.$nextTick(() => {
        this.focusElement("domain");
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
            provision: "new-domain",
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
      this.samba.hostname = defaults.hostname;
      this.samba.nbdomain = defaults.nbdomain;
      this.samba.realm = defaults.realm;

      // clear other fields
      this.samba.adminpass = "";
      this.samba.ipaddress = "";
      this.samba.enableFileServer = false;

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

      // focus on first field
      this.$nextTick(() => {
        this.focusElement("realm");
      });
    },
    clearSambaErrors() {
      for (const key of Object.keys(this.error.samba)) {
        this.error.samba[key] = "";
      }
    },
    validateConfigureSambaModule() {
      this.clearSambaErrors();
      let isValidationOk = true;

      // samba realm

      if (!this.samba.realm) {
        this.error.samba.realm = "common.required";

        if (isValidationOk) {
          this.focusElement("realm");
          isValidationOk = false;
        }
      }

      // samba nbdomain

      if (!this.samba.nbdomain) {
        this.error.samba.nbdomain = "common.required";

        if (isValidationOk) {
          this.focusElement("nbdomain");
          isValidationOk = false;
        }
      }

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
          this.samba.focusPasswordField = { element: "newPassword" };
          isValidationOk = false;
        }
      } else {
        if (
          !this.samba.passwordValidation.isLengthOk ||
          !this.samba.passwordValidation.isLowercaseOk ||
          !this.samba.passwordValidation.isUppercaseOk ||
          !this.samba.passwordValidation.isNumberOk ||
          !this.samba.passwordValidation.isSymbolOk
        ) {
          if (!this.error.samba.adminpass) {
            this.error.samba.adminpass = "password.password_not_secure";
          }

          if (isValidationOk) {
            this.samba.focusPasswordField = { element: "newPassword" };
            isValidationOk = false;
          }
        }

        if (!this.samba.passwordValidation.isEqualOk) {
          if (!this.error.samba.adminpass) {
            this.error.samba.adminpass = "password.passwords_do_not_match";
          }

          if (!this.error.samba.confirmPassword) {
            this.error.samba.confirmPassword =
              "password.passwords_do_not_match";
          }

          if (isValidationOk) {
            this.samba.focusPasswordField = { element: "confirmPassword" };
            isValidationOk = false;
          }
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

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.configureSambaModuleAborted
      );

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

      const res = await to(
        this.createModuleTaskForApp(this.newProviderId, {
          action: taskAction,
          data: {
            adminuser: this.samba.adminuser,
            adminpass: this.samba.adminpass,
            realm: this.samba.realm,
            ipaddress: this.samba.ipaddress,
            hostname: this.samba.hostname,
            nbdomain: this.samba.nbdomain,
            provision: "new-domain",
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

          if (!focusAlreadySet) {
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
    clearOpenLdapErrors() {
      for (const key of Object.keys(this.error.openldap)) {
        this.error.openldap[key] = "";
      }
    },
    validateConfigureOpenLdapModule() {
      this.clearOpenLdapErrors();
      let isValidationOk = true;

      // openldap domain

      if (!this.openldap.domain) {
        this.error.openldap.domain = "common.required";

        if (isValidationOk) {
          this.focusElement("domain");
          isValidationOk = false;
        }
      }

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
          this.openldap.focusPasswordField = { element: "newPassword" };
          isValidationOk = false;
        }
      } else {
        if (
          !this.openldap.passwordValidation.isLengthOk ||
          !this.openldap.passwordValidation.isLowercaseOk ||
          !this.openldap.passwordValidation.isUppercaseOk ||
          !this.openldap.passwordValidation.isNumberOk ||
          !this.openldap.passwordValidation.isSymbolOk
        ) {
          if (!this.error.openldap.admpass) {
            this.error.openldap.admpass = "password.password_not_secure";
          }

          if (isValidationOk) {
            this.openldap.focusPasswordField = { element: "newPassword" };
            isValidationOk = false;
          }
        }

        if (!this.openldap.passwordValidation.isEqualOk) {
          if (!this.error.openldap.admpass) {
            this.error.openldap.admpass = "password.passwords_do_not_match";
          }

          if (!this.error.openldap.confirmPassword) {
            this.error.openldap.confirmPassword =
              "password.passwords_do_not_match";
          }

          if (isValidationOk) {
            this.openldap.focusPasswordField = { element: "confirmPassword" };
            isValidationOk = false;
          }
        }
      }
      return isValidationOk;
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
            domain: this.openldap.domain,
            admuser: this.openldap.admuser,
            admpass: this.openldap.admpass,
            provision: "new-domain",
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
    onNewSambaPasswordValidation(passwordValidation) {
      this.samba.passwordValidation = passwordValidation;
    },
    onNewOpenLdapPasswordValidation(passwordValidation) {
      this.openldap.passwordValidation = passwordValidation;
    },
    validateAddExternalDomain() {
      this.clearExternalDomainErrors();
      let isValidationOk = true;

      // domain
      if (!this.external.domain) {
        this.error.external.domain = "common.required";

        if (isValidationOk) {
          this.focusElement("domain");
          isValidationOk = false;
        }
      }

      // host
      if (!this.external.host) {
        this.error.external.host = "common.required";

        if (isValidationOk) {
          this.focusElement("host");
          isValidationOk = false;
        }
      }

      // port
      if (!this.external.port) {
        this.error.external.port = "common.required";

        if (isValidationOk) {
          this.focusElement("port");
          isValidationOk = false;
        }
      }

      // bind_dn
      if (!this.external.bind_dn) {
        this.error.external.bind_dn = "common.required";

        if (isValidationOk) {
          this.focusElement("bind_dn");
          isValidationOk = false;
        }
      }

      // bind_password
      if (!this.external.bind_password) {
        this.error.external.bind_password = "common.required";

        if (isValidationOk) {
          this.focusElement("bind_password");
          isValidationOk = false;
        }
      }

      return isValidationOk;
    },
    async addExternalDomain() {
      const isValidationOk = this.validateAddExternalDomain();
      if (!isValidationOk) {
        return;
      }

      this.loading.external.addExternalDomain = true;
      const taskAction = "add-external-domain";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.addExternalDomainAborted
      );

      // register to task validation
      this.$root.$once(
        `${taskAction}-validation-failed-${eventId}`,
        this.addExternalDomainValidationFailed
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.addExternalDomainCompleted
      );

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            domain: this.external.domain,
            protocol: "ldap",
            host: this.external.host,
            port: parseInt(this.external.port),
            bind_dn: this.external.bind_dn,
            bind_password: this.external.bind_password,
            base_dn: this.external.base_dn,
            tls: this.external.tls,
            tls_verify: this.external.tls_verify,
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
        this.error.addExternalDomain = this.getErrorMessage(err);
        return;
      }
    },
    addExternalDomainValidationFailed(validationErrors) {
      this.loading.external.addExternalDomain = false;
      let focusAlreadySet = false;

      for (const validationError of validationErrors) {
        const param = validationError.parameter;
        // set i18n error message
        this.error.external[param] = "domains." + validationError.error;

        if (!focusAlreadySet) {
          this.focusElement(param);
          focusAlreadySet = true;
        }
      }
    },
    addExternalDomainAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.loading.external.addExternalDomain = false;
      this.$emit("hide");
    },
    addExternalDomainCompleted() {
      this.loading.external.addExternalDomain = false;

      // hide modal
      this.$emit("hide");

      // reload domains
      this.$emit("reloadDomains");
    },
    clearExternalDomainErrors() {
      for (const key of Object.keys(this.error.external)) {
        this.error.external[key] = "";
      }
    },
    onSelectNode(selectedNode) {
      this.selectedNode = selectedNode;
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";
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

<style lang="scss">
@import "../../styles/carbon-utils";

// global styles

.new-provider-password .new-password-container {
  margin-bottom: $spacing-06 !important;
}
</style>

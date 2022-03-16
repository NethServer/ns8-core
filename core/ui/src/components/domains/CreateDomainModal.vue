<template>
  <NsWizard
    size="default"
    :visible="isShown"
    :cancelLabel="$t('common.cancel')"
    :previousLabel="$t('common.previous')"
    :nextLabel="nextButtonLabel"
    :isPreviousDisabled="isPreviousButtonDisabled"
    :isNextDisabled="isNextButtonDisabled"
    :isNextLoading="
      loading.samba.configureModule || loading.external.addExternalDomain
    "
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
                disabled
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
          <!-- //// remove, schema is automatically detected -->
          <!-- <label class="bx--label">{{ $t("domains.type") }}</label>
          <cv-radio-group :vertical="false" ref="schema">
            <cv-radio-button
              name="schema-group"
              :label="$t('domains.openldap')"
              value="rfc2307"
              :checked="true"
              v-model="external.schema"
              :disabled="loading.external.addExternalDomain"
            />
            <cv-radio-button
              name="schema-group"
              :label="$t('domains.samba')"
              value="ad"
              v-model="external.schema"
              :disabled="loading.external.addExternalDomain"
            />
          </cv-radio-group> -->
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
        <NodeSelector
          :nodes="nodes"
          class="mg-top-lg"
          @selectNode="onSelectNode"
        />
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
        <cv-form v-if="isOpenLdapSelected"> //// openldap config </cv-form>
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
            <cv-text-input
              :label="$t('samba.adminuser')"
              v-model.trim="samba.adminuser"
              :helper-text="$t('samba.choose_samba_admin_username')"
              :invalid-message="$t(error.samba.adminuser)"
              :disabled="
                loading.samba.configureModule || loading.samba.getDefaults
              "
              ref="adminuser"
            >
            </cv-text-input>
            <NsPasswordInput
              :newPasswordLabel="$t('samba.adminpass')"
              :confirmPasswordLabel="$t('samba.adminpass_confirm')"
              v-model="samba.adminpass"
              @passwordValidation="onNewSambaPasswordValidation"
              :newPaswordHelperText="$t('samba.choose_samba_admin_password')"
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
              class="new-samba-password"
            />
            <cv-combo-box
              v-model="samba.ipaddress"
              :options="samba.ipAddressOptions"
              auto-highlight
              :title="$t('samba.ipaddress')"
              :invalid-message="$t(error.samba.ipaddress)"
              :disabled="
                loading.samba.configureModule || loading.samba.getDefaults
              "
              light
              ref="ipaddress"
            >
            </cv-combo-box>
            <cv-text-input
              :label="$t('samba.hostname')"
              v-model.trim="samba.hostname"
              :invalid-message="$t(error.samba.hostname)"
              :disabled="
                loading.samba.configureModule || loading.samba.getDefaults
              "
              ref="hostname"
            >
            </cv-text-input>
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
import NodeSelector from "@/components/misc/NodeSelector";

export default {
  name: "CreateDomainModal",
  components: { NodeSelector },
  mixins: [UtilService, TaskService, IconService, LottieService],
  props: {
    isShown: {
      type: Boolean,
      default: true,
    },
    nodes: {
      type: Array,
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
    isOpenLdap: {
      type: Boolean,
      default: false,
    },
    isSamba: {
      type: Boolean,
      default: false,
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
    nextButtonLabel() {
      if (
        (this.nodes.length == 1 && this.step == "instance") ||
        this.step == "node"
      ) {
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
        this.loading.samba.configureModule ||
        this.loading.external.addExternalDomain ||
        this.step == "installingProvider" ||
        (this.step == "location" &&
          !this.isInternalSelected &&
          !this.isExternalSelected) ||
        (this.step == "instance" &&
          !this.isOpenLdapSelected &&
          !this.isSambaSelected) ||
        (this.step == "node" && !this.selectedNode)
      );
    },
    isPreviousButtonDisabled() {
      return (
        this.isResumeConfiguration ||
        this.loading.samba.configureModule ||
        this.loading.external.addExternalDomain ||
        ["location", "installingProvider", "internalConfig"].includes(this.step)
      );
    },
  },
  watch: {
    isShown: function () {
      if (this.isShown) {
        // this.clearOpenLdapErrors(); ////
        this.clearSambaErrors();

        if (!this.isResumeConfiguration) {
          // start wizard from first step
          this.step = "location";
        } else {
          // resume configuration
          this.step = "internalConfig";

          if (this.isOpenLdap) {
            // this.getOpenLdapDefaults(); ////
          } else if (this.isSamba) {
            this.getSambaDefaults();
          }
        }

        if (this.step !== "internalConfig") {
          // // set focus to next button //// not working with NsWizard
          // setTimeout(() => {
          //   console.log("this.$refs", this.$refs); ////
          //   const element = this.$refs["wizardNext"].$el;
          //   element.focus();
          // }, 300);
        } else {
          if (this.isOpenLdapSelected) {
            //// focus first input field
          } else if (this.isSambaSelected) {
            setTimeout(() => {
              this.focusElement("realm");
            }, 300);
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
    step: function () {
      if (this.step == "internalConfig") {
        if (this.isOpenLdapSelected) {
          //// focus first input field
        } else if (this.isSambaSelected) {
          setTimeout(() => {
            this.focusElement("realm");
          }, 300);
        }
      }
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
          if (this.nodes.length > 1) {
            this.step = "node";
          } else {
            this.step = "installingProvider";
            this.installProvider();
          }
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

      // latest stable version
      let version = "ns8-stable";

      const taskAction = "add-internal-provider";

      // register to task completion
      this.$root.$once(
        taskAction + "-completed",
        this.addInternalProviderCompleted
      );

      // register to task progress to update progress bar
      this.$root.$on(
        taskAction + "-progress",
        this.addInternalProviderProgress
      );

      // register to task error
      this.$root.$off(taskAction + "-aborted");
      this.$root.$once(
        taskAction + "-aborted",
        this.addInternalProviderAborted
      );

      const selectedNodeId = this.selectedNode
        ? parseInt(this.selectedNode.id)
        : 1;

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            image: "ghcr.io/nethserver/samba:" + version, ////
            node: selectedNodeId,
          },
          extra: {
            title: this.$t("action." + taskAction),
            node: selectedNodeId,
            isNotificationHidden: true,
            isProgressNotified: true,
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

      // hide modal so that user can see error notification
      this.$emit("hide");
    },
    addInternalProviderCompleted(taskContext, taskResult) {
      // unregister to task progress
      this.$root.$off("add-internal-provider-progress");

      this.step = "internalConfig";

      if (this.isSambaSelected) {
        this.newProviderId = taskResult.output.module_id;
        this.getSambaDefaults();
      } //// else openldap

      // reload domains
      this.$emit("reloadDomains");

      // show new app in app drawer
      this.$root.$emit("reloadAppDrawer");
    },
    addInternalProviderProgress(progress) {
      this.installProviderProgress = progress;
    },
    async getSambaDefaults() {
      this.loading.samba.getDefaults = true;
      this.error.samba.getDefaults = "";
      const taskAction = "get-defaults";

      // register to task completion
      this.$root.$once(
        taskAction + "-completed",
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
    getSambaDefaultsCompleted(taskContext, taskResult) {
      this.loading.samba.getDefaults = false;
      const defaults = taskResult.output;

      this.samba.adminuser = defaults.adminuser;
      this.samba.hostname = defaults.hostname;
      this.samba.nbdomain = defaults.nbdomain;
      this.samba.realm = defaults.realm;

      // clear password
      this.samba.adminpass = "";

      // ip address combo box
      let index = 0;
      for (const ipObject of defaults.ipaddress_list) {
        const option = {
          name: ipObject.ipaddress.replace(/\W/g, "_"),
          label: ipObject.ipaddress + " - " + ipObject.label,
          value: ipObject.ipaddress,
        };
        this.$set(this.samba.ipAddressOptions, index, option);
        index++;
      }
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

      // samba ip address

      if (!this.samba.ipaddress) {
        this.error.samba.ipaddress = "common.required";

        if (isValidationOk) {
          this.focusElement("ipaddress");
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
      return isValidationOk;
    },
    async configureSambaModule() {
      const isValidationOk = this.validateConfigureSambaModule();
      if (!isValidationOk) {
        return;
      }

      this.loading.samba.configureModule = true;
      const taskAction = "configure-module";

      // register to task validation
      this.$root.$off(taskAction + "-validation-failed");
      this.$root.$once(
        taskAction + "-validation-failed",
        this.configureSambaModuleValidationFailed
      );

      // register to task completion
      this.$root.$off(taskAction + "-completed");
      this.$root.$once(
        taskAction + "-completed",
        this.configureSambaModuleCompleted
      );

      // register to task error
      this.$root.$off(taskAction + "-aborted");
      this.$root.$once(
        taskAction + "-aborted",
        this.configureSambaModuleAborted
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
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
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
    configureSambaModuleValidationFailed(validationErrors) {
      this.loading.samba.configureModule = false;
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
    configureSambaModuleCompleted() {
      this.loading.samba.configureModule = false;

      // hide modal
      this.$emit("hide");

      // reload domains
      this.$emit("reloadDomains");
    },
    configureSambaModuleAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.loading.samba.configureModule = false;

      // hide modal so that user can see error notification
      this.$emit("hide");
    },
    async configureOpenLdapModule() {
      console.log("configureOpenLdapModule"); ////
    },
    onNewSambaPasswordValidation(passwordValidation) {
      this.samba.passwordValidation = passwordValidation;
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

      // register to task validation
      this.$root.$off(taskAction + "-validation-failed");
      this.$root.$once(
        taskAction + "-validation-failed",
        this.addExternalDomainValidationFailed
      );

      // register to task error
      this.$root.$off(taskAction + "-aborted");
      this.$root.$once(taskAction + "-aborted", this.addExternalDomainAborted);

      // register to task completion
      this.$root.$off(taskAction + "-completed");
      this.$root.$once(
        taskAction + "-completed",
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
</style>

<style lang="scss">
@import "../../styles/carbon-utils";

// global styles

.new-samba-password .new-password-container {
  margin-bottom: $spacing-06 !important;
}
</style>
